package collab

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 64 * 1024 // 64KB
	sendBufferSize = 256
)

// Client represents a WebSocket client
//
// KONKURENSI: satu Client disentuh tiga goroutine sekaligus — ReadPump-nya
// sendiri, loop Hub, dan goroutine HTTP (mis. update role). Semua field yang
// bisa berubah setelah registrasi karena itu dijaga `mu` dan hanya boleh
// diakses lewat getter/setter di bawah. `mu` SELALU lock terdalam: jangan
// pernah mengambil room.mu atau hub.mu sambil memegangnya.
type Client struct {
	conn     *websocket.Conn
	hub      *Hub
	RoomID   string
	UserID   string
	Username string
	Role     string // JWT app role (admin/respondent) — immutable setelah dibuat
	// RoleFromInvite menandai role yang berasal dari token undangan. Role
	// semacam itu tidak boleh dinaikkan oleh aturan "orang pertama jadi owner":
	// undangan viewer harus tetap viewer meski dia kebetulan masuk saat room
	// sedang kosong. Di-set sebelum client masuk hub, jadi immutable.
	RoleFromInvite bool

	mu              sync.RWMutex
	roomRole        string                 // per-room: owner|editor|viewer
	followingUserID string                 // Figma-like follow target
	viewport        map[string]interface{} // viewport terakhir, dipush ke follower
	chatBubble      map[string]interface{} // cursor chat aktif: {x, y, text}; nil kalau tidak ada

	// sendMu menjaga `send` dari send-on-closed-channel. Channel ini ditutup
	// hub saat unregister, sementara peserta lain bisa sedang mengirim pesan
	// ke sini (follow/unfollow, update role) — tanpa gerbang ini prosesnya
	// panic dan seluruh server ikut mati.
	sendMu   sync.RWMutex
	sendOpen bool
	send     chan *Message

	lastMessageTime time.Time // hanya disentuh ReadPump
	messageCount    int       // hanya disentuh ReadPump
	stopCh          chan struct{}
}

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn, hub *Hub, roomID, userID, username, role, roomRole string) *Client {
	if roomRole == "" {
		// Fail-closed: role yang tidak diketahui hanya boleh menonton.
		roomRole = RoomRoleViewer
	}
	return &Client{
		conn:     conn,
		hub:      hub,
		RoomID:   roomID,
		UserID:   userID,
		Username: username,
		Role:     role,
		roomRole: roomRole,
		viewport: map[string]interface{}{},
		send:     make(chan *Message, sendBufferSize),
		sendOpen: true,
		stopCh:   make(chan struct{}),
	}
}

// ─── Accessor aman-konkuren ──────────────────────────────────────────────────

// RoomRole - role client di room saat ini (owner|editor|viewer)
func (c *Client) RoomRole() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.roomRole
}

// SetRoomRole - ubah role client di room
func (c *Client) SetRoomRole(role string) {
	c.mu.Lock()
	c.roomRole = role
	c.mu.Unlock()
}

// FollowingUserID - user yang sedang diikuti client ini ("" kalau tidak ada)
func (c *Client) FollowingUserID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.followingUserID
}

// SetFollowingUserID - set/lepas target follow
func (c *Client) SetFollowingUserID(userID string) {
	c.mu.Lock()
	c.followingUserID = userID
	c.mu.Unlock()
}

// SetViewport - simpan viewport terakhir client
func (c *Client) SetViewport(vp map[string]interface{}) {
	c.mu.Lock()
	c.viewport = vp
	c.mu.Unlock()
}

// ViewportCopy - salinan viewport terakhir; aman dibaca goroutine lain
func (c *Client) ViewportCopy() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.viewport) == 0 {
		return nil
	}
	out := make(map[string]interface{}, len(c.viewport))
	for k, v := range c.viewport {
		out[k] = v
	}
	return out
}

// SetChatBubble - set/hapus bubble cursor chat client ini
func (c *Client) SetChatBubble(bubble map[string]interface{}) {
	c.mu.Lock()
	c.chatBubble = bubble
	c.mu.Unlock()
}

// ChatBubbleCopy - salinan bubble aktif; nil kalau tidak ada bubble terbuka
func (c *Client) ChatBubbleCopy() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.chatBubble == nil {
		return nil
	}
	out := make(map[string]interface{}, len(c.chatBubble))
	for k, v := range c.chatBubble {
		out[k] = v
	}
	return out
}

// SetChatBubbleText - perbarui teks bubble yang sedang terbuka.
// Mengembalikan false kalau tidak ada bubble aktif.
func (c *Client) SetChatBubbleText(text string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.chatBubble == nil {
		return false
	}
	c.chatBubble["text"] = text
	return true
}

// closeSend - tutup channel kirim satu kali saja, aman dipanggil berulang.
// Setelah ini semua sendQuiet/trySend jadi no-op, bukan panic.
func (c *Client) closeSend() {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()
	if !c.sendOpen {
		return
	}
	c.sendOpen = false
	close(c.send)
}

// trySend - satu-satunya jalan menulis ke c.send.
// Mengembalikan false kalau channel sudah ditutup atau buffer penuh.
func (c *Client) trySend(msg *Message) bool {
	c.sendMu.RLock()
	defer c.sendMu.RUnlock()
	if !c.sendOpen {
		return false
	}
	select {
	case c.send <- msg:
		return true
	default:
		return false
	}
}

// canEdit - cek apakah client boleh mengubah data (owner/editor). Viewer selalu false
func (c *Client) canEdit() bool {
	role := c.RoomRole()
	return role == RoomRoleOwner || role == RoomRoleEditor
}

// sendQuiet - kirim pesan ke channel client tanpa blocking; drop + log kalau buffer penuh atau sudah tertutup
func (c *Client) sendQuiet(msg *Message) {
	if !c.trySend(msg) {
		log.Printf("Failed to send to client %s (buffer penuh atau koneksi tertutup)", c.UserID)
	}
}

// sendError - kirim pesan error (code + message) balik ke client yang bersangkutan
func (c *Client) sendError(code, message string) {
	c.sendQuiet(newMessage(MsgError, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"code":    code,
		"message": message,
	}))
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	// Panic di goroutine ini TIDAK tertangkap gin.Recovery (goroutine terpisah),
	// jadi tanpa recover satu bug di jalur pesan mematikan seluruh proses API.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC di ReadPump client %s room %s: %v", c.UserID, c.RoomID, r)
		}
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-c.stopCh:
			return
		default:
			_, messageBytes, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}

			if !c.checkRateLimit() {
				c.sendError("RATE_LIMITED", "Terlalu banyak pesan — tunggu sebentar")
				continue
			}

			var rawMsg map[string]interface{}
			if err := json.Unmarshal(messageBytes, &rawMsg); err != nil {
				c.sendError("INVALID_JSON", "Format pesan tidak valid")
				continue
			}

			msgType, _ := rawMsg["type"].(string)
			payload, _ := rawMsg["payload"].(map[string]interface{})
			msg := newMessage(msgType, c.RoomID, c.UserID, c.Username, payload)
			c.handleMessage(msg)
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC di WritePump client %s room %s: %v", c.UserID, c.RoomID, r)
		}
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			messageJSON, err := json.Marshal(message)
			if err != nil {
				continue
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.stopCh:
			return
		}
	}
}

// handleMessage - router utama pesan masuk: cek izin role lalu teruskan ke handler sesuai tipe pesan
func (c *Client) handleMessage(msg *Message) {
	switch msg.Type {
	case MsgPresenceJoin:
		// Already registered on connect; re-broadcast presence list
		room := c.hub.GetOrCreateRoom(c.RoomID)
		c.sendQuiet(c.hub.buildPresenceList(room))
		c.sendQuiet(c.hub.buildFollowState(room))

	case MsgPresenceLeave:
		c.hub.unregister <- c

	case MsgCursorMove:
		c.handleCursorMove(msg)

	case MsgViewportUpdate:
		c.handleViewportUpdate(msg)

	case MsgFollowUser:
		c.handleFollowUser(msg)

	case MsgUnfollowUser:
		c.handleUnfollowUser()

	case MsgCursorChatOpen:
		c.handleCursorChatOpen(msg)

	case MsgCursorChatUpdate:
		c.handleCursorChatUpdate(msg)

	case MsgCursorChatClose:
		c.handleCursorChatClose()

	case MsgFoodSearch:
		if !c.canEdit() {
			c.sendError("FORBIDDEN", "Viewer hanya bisa mengikuti — tidak bisa mencari bersama")
			return
		}
		c.handleFoodSearch(msg)

	case MsgFoodSelect:
		if !c.canEdit() {
			c.sendError("FORBIDDEN", "Viewer tidak dapat memilih makanan")
			return
		}
		c.handleFoodSelect(msg)

	case MsgMealAdd:
		if !c.canEdit() {
			c.sendError("FORBIDDEN", "Viewer tidak dapat menambah makanan")
			return
		}
		c.handleMealAdd(msg)

	case MsgPortionSet, MsgPortionSelect:
		if !c.canEdit() {
			c.sendError("FORBIDDEN", "Viewer tidak dapat mengatur porsi")
			return
		}
		c.handlePortionSet(msg)

	case MsgReviewSubmit:
		if !c.canEdit() {
			c.sendError("FORBIDDEN", "Viewer tidak dapat mengirim laporan")
			return
		}
		c.handleReviewSubmit(msg)

	case MsgDBEditStart:
		if !c.canEdit() {
			c.sendError("FORBIDDEN", "Viewer tidak dapat mengedit")
			return
		}
		c.handleDBEditStart(msg)

	case MsgDBEditField:
		if !c.canEdit() {
			c.sendError("FORBIDDEN", "Viewer tidak dapat mengedit")
			return
		}
		c.handleDBEditField(msg)

	case MsgDBEditSave:
		if !c.canEdit() {
			c.sendError("FORBIDDEN", "Viewer tidak dapat mengedit")
			return
		}
		c.handleDBEditSave(msg)

	case MsgDBEditCancel:
		if !c.canEdit() {
			c.sendError("FORBIDDEN", "Viewer tidak dapat mengedit")
			return
		}
		c.handleDBEditCancel(msg)

	case MsgCanvasDrawStart:
		if !c.canEdit() {
			c.sendError("FORBIDDEN", "Viewer role tidak dapat menggambar di canvas")
			return
		}
		c.handleCanvasDrawStart(msg)

	case MsgCanvasDrawMove:
		if !c.canEdit() {
			c.sendError("FORBIDDEN", "Viewer role tidak dapat menggambar di canvas")
			return
		}
		c.handleCanvasDrawMove(msg)

	case MsgCanvasDrawEnd:
		if !c.canEdit() {
			c.sendError("FORBIDDEN", "Viewer role tidak dapat menggambar di canvas")
			return
		}
		c.handleCanvasDrawEnd(msg)

	case MsgCanvasLaserMove:
		if !c.canEdit() {
			c.sendError("FORBIDDEN", "Viewer role tidak dapat menggunakan laser pointer")
			return
		}
		c.handleCanvasLaserMove(msg)

	case MsgCanvasClear:
		if !c.canEdit() {
			c.sendError("FORBIDDEN", "Viewer role tidak dapat menghapus canvas")
			return
		}
		c.handleCanvasClear(msg)

	case MsgChatMessage:
		c.handleChatMessage(msg)

	case MsgUpdateUserRole:
		c.handleUpdateUserRole(msg)

	case MsgGetHistory:
		c.sendRoomHistory()

	case MsgPing:
		c.sendQuiet(newMessage(MsgPong, c.RoomID, c.UserID, c.Username, map[string]interface{}{}))

	default:
		log.Printf("Unknown message type: %s", msg.Type)
		c.sendError("UNKNOWN_TYPE", "Tipe pesan tidak dikenal: "+msg.Type)
	}
}

// handleCursorMove - siarkan posisi kursor client ke room (di-batch oleh Room agar hemat bandwidth)
func (c *Client) handleCursorMove(msg *Message) {
	payload := msg.Payload
	if payload == nil {
		payload = map[string]interface{}{}
	}
	payload["color"] = colorForUser(c.UserID)
	payload["page"] = payloadString(payload, "page")
	payload["room_role"] = c.RoomRole()
	update := newMessage(MsgCursorUpdate, c.RoomID, c.UserID, c.Username, payload)
	room := c.hub.GetOrCreateRoom(c.RoomID)
	room.AddMessage(update)
}

// handleViewportUpdate - simpan viewport terakhir client lalu kirim ke follower yang mengikutinya
func (c *Client) handleViewportUpdate(msg *Message) {
	payload := msg.Payload
	if payload == nil {
		payload = map[string]interface{}{}
	}
	// Simpan viewport leader untuk follower yang baru join follow
	c.SetViewport(map[string]interface{}{
		"page":      payloadString(payload, "page"),
		"scroll_x":  payload["scroll_x"],
		"scroll_y":  payload["scroll_y"],
		"step":      payloadString(payload, "step"),
		"path":      payloadString(payload, "path"),
		"zoom":      payload["zoom"],
		"timestamp": time.Now().UnixMilli(),
	})
	payload["color"] = colorForUser(c.UserID)
	payload["display_name"] = c.Username

	sync := newMessage(MsgViewportSync, c.RoomID, c.UserID, c.Username, payload)
	// Kirim hanya ke follower yang sedang mengikuti user ini
	c.hub.broadcastToFollowers(c.RoomID, c.UserID, sync)
}

// handleFollowUser - mulai follow user lain (mode Figma): set leader, beritahu kedua pihak, lalu push viewport leader saat ini
func (c *Client) handleFollowUser(msg *Message) {
	targetID := payloadString(msg.Payload, "user_id")
	if targetID == "" || targetID == c.UserID {
		c.sendError("INVALID_PAYLOAD", "user_id target follow tidak valid")
		return
	}

	room := c.hub.GetOrCreateRoom(c.RoomID)
	target := c.hub.findClientInRoom(room, targetID)
	if target == nil {
		c.sendError("NOT_FOUND", "User target tidak ada di room")
		return
	}

	c.SetFollowingUserID(targetID)

	started := newMessage(MsgFollowStarted, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"follower_id":   c.UserID,
		"follower_name": c.Username,
		"leader_id":     targetID,
		"leader_name":   target.Username,
		"leader_color":  colorForUser(targetID),
	})
	c.sendQuiet(started)
	// Beritahu leader seseorang mengikuti
	target.sendQuiet(started)
	c.hub.broadcastToRoom(room, c.hub.buildFollowState(room), nil)

	// Push viewport leader saat ini agar follower langsung mirror.
	// ViewportCopy() menyalin di bawah lock — membaca map leader langsung
	// dari goroutine ini adalah data race.
	if vp := target.ViewportCopy(); vp != nil {
		vp["color"] = colorForUser(targetID)
		vp["display_name"] = target.Username
		c.sendQuiet(newMessage(MsgViewportSync, c.RoomID, targetID, target.Username, vp))
	}

	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "follow",
		"details": c.Username + " mengikuti " + target.Username,
	}))
}

// handleUnfollowUser - berhenti mengikuti leader dan beritahu room bahwa follow state berubah
func (c *Client) handleUnfollowUser() {
	prev := c.FollowingUserID()
	if prev == "" {
		return
	}
	c.SetFollowingUserID("")

	room := c.hub.GetOrCreateRoom(c.RoomID)
	stopped := newMessage(MsgFollowStopped, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"follower_id": c.UserID,
		"leader_id":   prev,
	})
	c.sendQuiet(stopped)
	if leader := c.hub.findClientInRoom(room, prev); leader != nil {
		leader.sendQuiet(stopped)
	}
	c.hub.broadcastToRoom(room, c.hub.buildFollowState(room), nil)
}

// clampCursorChatText - trim & potong teks bubble ke batas aman sebelum disimpan/disiarkan
func clampCursorChatText(raw string) string {
	if len(raw) > cursorChatMaxTextLen {
		raw = raw[:cursorChatMaxTextLen]
	}
	return raw
}

// broadcastCursorChatState - siarkan state bubble client ini (posisi + teks tersimpan di c.ChatBubble)
func (c *Client) broadcastCursorChatState() {
	bubble := c.ChatBubbleCopy()
	if bubble == nil {
		return
	}
	payload := map[string]interface{}{
		"x":         bubble["x"],
		"y":         bubble["y"],
		"text":      bubble["text"],
		"color":     colorForUser(c.UserID),
		"room_role": c.RoomRole(),
	}
	c.hub.Publish(newMessage(MsgCursorChatUpdated, c.RoomID, c.UserID, c.Username, payload))
}

// handleCursorChatOpen - buka bubble chat baru di posisi kursor saat ini (dipicu tombol "/" di FE).
// Posisi di-anchor sekali di sini dan tidak berubah lagi selama bubble ini terbuka.
func (c *Client) handleCursorChatOpen(msg *Message) {
	c.SetChatBubble(map[string]interface{}{
		"x":    msg.Payload["x"],
		"y":    msg.Payload["y"],
		"text": clampCursorChatText(payloadString(msg.Payload, "text")),
	})
	c.broadcastCursorChatState()
}

// handleCursorChatUpdate - perbarui teks bubble yang sedang terbuka; diabaikan kalau belum ada open
func (c *Client) handleCursorChatUpdate(msg *Message) {
	if !c.SetChatBubbleText(clampCursorChatText(payloadString(msg.Payload, "text"))) {
		return
	}
	c.broadcastCursorChatState()
}

// handleCursorChatClose - tutup bubble (Enter/Esc/timeout di FE) dan beritahu room agar bubble dihapus
func (c *Client) handleCursorChatClose() {
	if c.ChatBubbleCopy() == nil {
		return
	}
	c.SetChatBubble(nil)
	c.hub.Publish(newMessage(MsgCursorChatClosed, c.RoomID, c.UserID, c.Username, map[string]interface{}{}))
}

// handleFoodSearch - siarkan kata kunci pencarian makanan ke seluruh anggota room + catat activity log
func (c *Client) handleFoodSearch(msg *Message) {
	query := payloadString(msg.Payload, "query")
	c.hub.Publish(newMessage(MsgUserSearching, c.RoomID, c.UserID, c.Username, msg.Payload))
	c.hub.Publish(newMessage(MsgFoodSearchShared, c.RoomID, c.UserID, c.Username, msg.Payload))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "food_search",
		"details": c.Username + ` mencari "` + query + `"`,
		"query":   query,
	}))
}

// handleFoodSelect - siarkan makanan yang dipilih client ke seluruh anggota room
func (c *Client) handleFoodSelect(msg *Message) {
	foodName := payloadString(msg.Payload, "food_name")
	c.hub.Publish(newMessage(MsgFoodSelected, c.RoomID, c.UserID, c.Username, msg.Payload))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "food_select",
		"details": c.Username + " memilih " + foodName,
	}))
}

// handleMealAdd - siarkan penambahan makanan ke slot meal tertentu (sarapan/makan siang/dll)
func (c *Client) handleMealAdd(msg *Message) {
	foodName := payloadString(msg.Payload, "food_name")
	mealType := payloadString(msg.Payload, "meal_type")
	c.hub.Publish(newMessage(MsgMealUpdated, c.RoomID, c.UserID, c.Username, msg.Payload))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "meal_add",
		"details": c.Username + " menambah " + foodName + " ke " + mealType,
	}))
}

// handlePortionSet - siarkan perubahan porsi makanan ke seluruh anggota room
func (c *Client) handlePortionSet(msg *Message) {
	c.hub.Publish(newMessage(MsgPortionUpdated, c.RoomID, c.UserID, c.Username, msg.Payload))
	c.hub.Publish(newMessage(MsgPortionSelected, c.RoomID, c.UserID, c.Username, msg.Payload))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "portion_set",
		"details": c.Username + " mengatur porsi",
	}))
}

// handleReviewSubmit - siarkan bahwa client mengirim review/submit survey
func (c *Client) handleReviewSubmit(msg *Message) {
	c.hub.Publish(newMessage(MsgReviewSubmitted, c.RoomID, c.UserID, c.Username, msg.Payload))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "review_submit",
		"details": c.Username + " mengirim review/submit survey",
	}))
}

// handleDBEditStart - ambil lock entity sebelum edit; kalau sudah dikunci user lain kirim error LOCKED
func (c *Client) handleDBEditStart(msg *Message) {
	entityType := payloadString(msg.Payload, "entity_type")
	entityID := payloadString(msg.Payload, "entity_id")
	version := payloadInt(msg.Payload, "version")
	if entityType == "" || entityID == "" {
		c.sendError("INVALID_PAYLOAD", "entity_type dan entity_id wajib")
		return
	}

	lock, ok, existing := c.hub.Locks().TryLock(c.RoomID, entityType, entityID, c.UserID, c.Username, version)
	if !ok {
		c.sendError("LOCKED", "Sedang diedit oleh "+existing.Username)
		c.sendQuiet(newMessage(MsgDBLocked, c.RoomID, existing.LockedBy, existing.Username, map[string]interface{}{
			"entity_type": existing.EntityType,
			"entity_id":   existing.EntityID,
			"locked_by":   existing.LockedBy,
			"username":    existing.Username,
			"version":     existing.Version,
		}))
		return
	}

	payload := map[string]interface{}{
		"entity_type": lock.EntityType,
		"entity_id":   lock.EntityID,
		"locked_by":   lock.LockedBy,
		"username":    lock.Username,
		"version":     lock.Version,
	}
	c.hub.Publish(newMessage(MsgDBLocked, c.RoomID, c.UserID, c.Username, payload))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "db_edit_start",
		"details": c.Username + " mulai edit " + entityType + " " + entityID,
	}))
}

// handleDBEditField - siarkan perubahan per-field secara live; hanya boleh oleh pemegang lock
func (c *Client) handleDBEditField(msg *Message) {
	entityType := payloadString(msg.Payload, "entity_type")
	entityID := payloadString(msg.Payload, "entity_id")
	lock := c.hub.Locks().Get(c.RoomID, entityType, entityID)
	if lock == nil || lock.LockedBy != c.UserID {
		c.sendError("NOT_LOCK_OWNER", "Anda tidak memegang lock entity ini")
		return
	}
	c.hub.Publish(newMessage(MsgDBFieldUpdated, c.RoomID, c.UserID, c.Username, msg.Payload))
}

// handleDBEditSave - simpan hasil edit: cek lock + versi (optimistic locking), naikkan versi, lalu lepas lock
func (c *Client) handleDBEditSave(msg *Message) {
	entityType := payloadString(msg.Payload, "entity_type")
	entityID := payloadString(msg.Payload, "entity_id")
	clientVersion := payloadInt(msg.Payload, "version")

	lock := c.hub.Locks().Get(c.RoomID, entityType, entityID)
	if lock == nil || lock.LockedBy != c.UserID {
		c.sendError("NOT_LOCK_OWNER", "Anda tidak memegang lock entity ini")
		return
	}
	if clientVersion > 0 && clientVersion < lock.Version {
		c.sendError("VERSION_CONFLICT", "Versi data sudah berubah — refresh lalu coba lagi")
		return
	}

	bumped, ok := c.hub.Locks().BumpVersion(c.RoomID, entityType, entityID, c.UserID)
	if !ok {
		c.sendError("VERSION_CONFLICT", "Gagal menyimpan versi")
		return
	}

	c.hub.Locks().Release(c.RoomID, entityType, entityID, c.UserID)
	c.hub.Publish(newMessage(MsgDBEditSaved, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"entity_type": entityType,
		"entity_id":   entityID,
		"version":     bumped.Version,
		"changes":     msg.Payload["changes"],
	}))
	c.hub.Publish(newMessage(MsgDBUnlocked, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"entity_type": entityType,
		"entity_id":   entityID,
	}))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "db_edit_save",
		"details": c.Username + " menyimpan " + entityType + " " + entityID,
	}))
}

// handleDBEditCancel - batalkan edit dan lepaskan lock entity
func (c *Client) handleDBEditCancel(msg *Message) {
	entityType := payloadString(msg.Payload, "entity_type")
	entityID := payloadString(msg.Payload, "entity_id")
	if !c.hub.Locks().Release(c.RoomID, entityType, entityID, c.UserID) {
		c.sendError("NOT_LOCK_OWNER", "Tidak bisa melepas lock")
		return
	}
	c.hub.Publish(newMessage(MsgDBUnlocked, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"entity_type": entityType,
		"entity_id":   entityID,
	}))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "db_edit_cancel",
		"details": c.Username + " membatalkan edit",
	}))
}

// handleChatMessage - teruskan pesan chat ke seluruh anggota room
func (c *Client) handleChatMessage(msg *Message) {
	c.hub.Publish(newMessage(MsgChatMessage, c.RoomID, c.UserID, c.Username, msg.Payload))
}

// sendRoomHistory - kirim 50 pesan terakhir room ke client yang baru minta history
func (c *Client) sendRoomHistory() {
	room := c.hub.GetOrCreateRoom(c.RoomID)
	c.sendQuiet(newMessage(MsgHistory, c.RoomID, "", "", map[string]interface{}{
		"messages": room.GetHistory(50),
	}))
}

// checkRateLimit - batasi maksimal 50 pesan per detik per client agar hub tidak dibanjiri
func (c *Client) checkRateLimit() bool {
	now := time.Now()
	if now.Sub(c.lastMessageTime) > time.Second {
		c.messageCount = 0
		c.lastMessageTime = now
	}
	c.messageCount++
	return c.messageCount <= 50
}

// handleCanvasDrawStart - inisiasi stroke baru di canvas: simpan ke memory room & broadcast
func (c *Client) handleCanvasDrawStart(msg *Message) {
	payload := msg.Payload
	if payload == nil {
		c.sendError("INVALID_PAYLOAD", "Payload canvas_draw_start tidak boleh kosong")
		return
	}
	strokeID := payloadString(payload, "stroke_id")
	if strokeID == "" {
		c.sendError("INVALID_PAYLOAD", "stroke_id wajib diisi")
		return
	}
	tool := payloadString(payload, "tool")
	if tool == "" {
		tool = "pencil"
	}
	color := payloadString(payload, "color")
	if color == "" {
		color = colorForUser(c.UserID)
	}
	width := payloadFloat64(payload, "width", 3.0)
	if width <= 0 {
		width = 3
	}
	targetImageID := payloadString(payload, "target_image_id")

	x := payloadFloat64(payload, "x", 0.0)
	y := payloadFloat64(payload, "y", 0.0)

	initialPoints := [][]float64{{x, y}}

	room := c.hub.GetOrCreateRoom(c.RoomID)
	room.AddCanvasStroke(&CanvasStrokeItem{
		StrokeID:      strokeID,
		UserID:        c.UserID,
		Username:      c.Username,
		Tool:          tool,
		Color:         color,
		Width:         width,
		TargetImageID: targetImageID,
		Points:        initialPoints,
		Timestamp:     time.Now().UnixMilli(),
	})

	outPayload := map[string]interface{}{
		"stroke_id":       strokeID,
		"user_id":         c.UserID,
		"username":        c.Username,
		"tool":            tool,
		"color":           color,
		"width":           width,
		"target_image_id": targetImageID,
		"x":               x,
		"y":               y,
		"points":          initialPoints,
	}

	c.hub.Publish(newMessage(MsgCanvasStrokeStarted, c.RoomID, c.UserID, c.Username, outPayload))
}

// handleCanvasDrawMove - perbarui stroke dengan sekumpulan koordinat titik baru
func (c *Client) handleCanvasDrawMove(msg *Message) {
	payload := msg.Payload
	if payload == nil {
		return
	}
	strokeID := payloadString(payload, "stroke_id")
	if strokeID == "" {
		return
	}

	var points [][]float64
	if rawPoints, ok := payload["points"].([]interface{}); ok {
		for _, pt := range rawPoints {
			if ptArr, ok := pt.([]interface{}); ok && len(ptArr) >= 2 {
				x, _ := ptArr[0].(float64)
				y, _ := ptArr[1].(float64)
				points = append(points, []float64{x, y})
			}
		}
	}

	if len(points) > 0 {
		room := c.hub.GetOrCreateRoom(c.RoomID)
		room.AppendCanvasPoints(strokeID, points)
	}

	outPayload := map[string]interface{}{
		"stroke_id": strokeID,
		"user_id":   c.UserID,
		"points":    points,
	}
	room := c.hub.GetOrCreateRoom(c.RoomID)
	room.AddMessage(newMessage(MsgCanvasStrokeUpdated, c.RoomID, c.UserID, c.Username, outPayload))
}

// handleCanvasDrawEnd - selesaikan stroke canvas
func (c *Client) handleCanvasDrawEnd(msg *Message) {
	payload := msg.Payload
	strokeID := payloadString(payload, "stroke_id")
	outPayload := map[string]interface{}{
		"stroke_id": strokeID,
		"user_id":   c.UserID,
	}
	c.hub.Publish(newMessage(MsgCanvasStrokeEnded, c.RoomID, c.UserID, c.Username, outPayload))
}

// handleCanvasLaserMove - siarkan posisi laser pointer
func (c *Client) handleCanvasLaserMove(msg *Message) {
	payload := msg.Payload
	if payload == nil {
		return
	}
	payload["user_id"] = c.UserID
	payload["username"] = c.Username
	if payloadString(payload, "color") == "" {
		payload["color"] = colorForUser(c.UserID)
	}
	room := c.hub.GetOrCreateRoom(c.RoomID)
	room.AddMessage(newMessage(MsgCanvasLaserUpdated, c.RoomID, c.UserID, c.Username, payload))
}

// handleCanvasClear - hapus seluruh coretan canvas di room
func (c *Client) handleCanvasClear(msg *Message) {
	targetImageID := payloadString(msg.Payload, "target_image_id")
	room := c.hub.GetOrCreateRoom(c.RoomID)
	room.ClearCanvasStrokes(targetImageID)

	c.hub.Publish(newMessage(MsgCanvasCleared, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"target_image_id": targetImageID,
		"cleared_by":      c.UserID,
	}))
}

// handleUpdateUserRole - memproses permintaan update role dari owner room
func (c *Client) handleUpdateUserRole(msg *Message) {
	if c.RoomRole() != RoomRoleOwner {
		c.sendError("FORBIDDEN", "Hanya owner room yang bisa mengubah role peserta")
		return
	}
	targetID := payloadString(msg.Payload, "target_user_id")
	newRole := payloadString(msg.Payload, "new_role")
	if targetID == "" || (newRole != RoomRoleEditor && newRole != RoomRoleViewer) {
		c.sendError("INVALID_PAYLOAD", "target_user_id dan new_role (editor|viewer) wajib diisi")
		return
	}
	if err := c.hub.UpdateUserRole(c.RoomID, c.UserID, targetID, newRole); err != nil {
		c.sendError("UPDATE_ROLE_FAILED", err.Error())
		return
	}
}

// Stop stops the client
func (c *Client) Stop() {
	close(c.stopCh)
}
