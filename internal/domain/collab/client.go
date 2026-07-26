package collab

import (
	"encoding/json"
	"log"
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
type Client struct {
	conn            *websocket.Conn
	hub             *Hub
	RoomID          string
	UserID          string
	Username        string
	Role            string // JWT app role (admin/respondent)
	RoomRole        string // per-room: owner|editor|viewer
	FollowingUserID string // Figma-like follow target
	Viewport        map[string]interface{}
	send            chan *Message
	lastMessageTime time.Time
	messageCount    int
	stopCh          chan struct{}
}

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn, hub *Hub, roomID, userID, username, role, roomRole string) *Client {
	if roomRole == "" {
		roomRole = RoomRoleEditor
	}
	return &Client{
		conn:     conn,
		hub:      hub,
		RoomID:   roomID,
		UserID:   userID,
		Username: username,
		Role:     role,
		RoomRole: roomRole,
		Viewport: map[string]interface{}{},
		send:     make(chan *Message, sendBufferSize),
		stopCh:   make(chan struct{}),
	}
}

// canEdit - cek apakah client boleh mengubah data (owner/editor). Viewer selalu false
func (c *Client) canEdit() bool {
	return c.RoomRole == RoomRoleOwner || c.RoomRole == RoomRoleEditor
}

// sendQuiet - kirim pesan ke channel client tanpa blocking; drop + log kalau buffer penuh
func (c *Client) sendQuiet(msg *Message) {
	select {
	case c.send <- msg:
	default:
		log.Printf("Failed to send to client %s", c.UserID)
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
	defer func() {
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

	case MsgChatMessage:
		c.handleChatMessage(msg)

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
	payload["room_role"] = c.RoomRole
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
	c.Viewport = map[string]interface{}{
		"page":       payloadString(payload, "page"),
		"scroll_x":   payload["scroll_x"],
		"scroll_y":   payload["scroll_y"],
		"step":       payloadString(payload, "step"),
		"path":       payloadString(payload, "path"),
		"zoom":       payload["zoom"],
		"timestamp":  time.Now().UnixMilli(),
	}
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

	c.FollowingUserID = targetID

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

	// Push viewport leader saat ini agar follower langsung mirror
	if len(target.Viewport) > 0 {
		vp := map[string]interface{}{}
		for k, v := range target.Viewport {
			vp[k] = v
		}
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
	if c.FollowingUserID == "" {
		return
	}
	prev := c.FollowingUserID
	c.FollowingUserID = ""

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

	lock, ok, existing := c.hub.Locks().TryLock(entityType, entityID, c.UserID, c.Username, version)
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
	lock := c.hub.Locks().Get(entityType, entityID)
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

	lock := c.hub.Locks().Get(entityType, entityID)
	if lock == nil || lock.LockedBy != c.UserID {
		c.sendError("NOT_LOCK_OWNER", "Anda tidak memegang lock entity ini")
		return
	}
	if clientVersion > 0 && clientVersion < lock.Version {
		c.sendError("VERSION_CONFLICT", "Versi data sudah berubah — refresh lalu coba lagi")
		return
	}

	bumped, ok := c.hub.Locks().BumpVersion(entityType, entityID, c.UserID)
	if !ok {
		c.sendError("VERSION_CONFLICT", "Gagal menyimpan versi")
		return
	}

	c.hub.Locks().Release(entityType, entityID, c.UserID)
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
	if !c.hub.Locks().Release(entityType, entityID, c.UserID) {
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

// Stop stops the client
func (c *Client) Stop() {
	close(c.stopCh)
}
