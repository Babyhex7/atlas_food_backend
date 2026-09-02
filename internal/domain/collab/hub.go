package collab

import (
	"log"
	"sync"
	"time"
)

// Hub manages all WebSocket rooms and clients (in-memory; Redis deferred).
type Hub struct {
	rooms      map[string]*Room
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	locks      *LockManager
	invites    *InviteStore
	mu         sync.RWMutex
	stopCh     chan struct{}
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]*Room),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		broadcast:  make(chan *Message, 1024),
		locks:      NewLockManager(),
		invites:    NewInviteStore(),
		stopCh:     make(chan struct{}),
	}
}

// Locks returns the in-memory lock manager.
func (h *Hub) Locks() *LockManager {
	return h.locks
}

// Invites returns invite token store.
func (h *Hub) Invites() *InviteStore {
	return h.invites
}

// Run starts the hub's main event loop
func (h *Hub) Run() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				h.cleanupInactiveRooms()
			case <-h.stopCh:
				return
			}
		}
	}()

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		case <-h.stopCh:
			log.Println("Hub stopping...")
			return
		}
	}
}

// GetOrCreateRoom gets existing room or creates new one
func (h *Hub) GetOrCreateRoom(roomID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.rooms[roomID]
	if !exists {
		room = NewRoom(roomID, h)
		h.rooms[roomID] = room
		go room.Run()
		log.Printf("📡 Created new room: %s", roomID)
	}
	return room
}

// registerClient - daftarkan client ke room, kirim presence/state/follow ke yang baru join.
//
// Satu user boleh punya banyak socket sekaligus (multi-tab / multi-device), persis
// seperti Figma. Socket lama TIDAK ditendang: kalau ditendang, tab lama akan
// otomatis reconnect lalu menendang tab baru, dan keduanya saling tendang tanpa
// henti. Duplikasi ditangani di tempat yang benar: presence list dedupe by user_id
// dan event "joined" hanya disiarkan untuk socket pertama milik user tersebut.
func (h *Hub) registerClient(client *Client) {
	room := h.GetOrCreateRoom(client.RoomID)

	room.mu.Lock()
	isFirst := len(room.clients) == 0
	// alreadyPresent = user ini sudah punya socket lain di room (tab kedua)
	alreadyPresent := false
	for existing := range room.clients {
		if existing.UserID != "" && existing.UserID == client.UserID && existing != client {
			// Tab baru mewarisi room role tab lama supaya owner tidak turun
			// hanya karena membuka tab kedua. Satu orang di satu room selalu
			// punya satu role — kecuali koneksi ini memang datang membawa
			// undangan yang menetapkan role lain (invite token menang).
			if !client.RoleFromInvite {
				client.RoomRole = existing.RoomRole
			}
			alreadyPresent = true
		}
	}

	// Urutan prioritas role (tertinggi ke terendah):
	// 1. Invite token (RoleFromInvite=true) — niat eksplisit pemilik room, selalu menang.
	// 2. Remembered role — cegah viewer naik jadi editor hanya karena pindah halaman/reconnect.
	// 3. Orang pertama di room tanpa invite → owner.
	// 4. Fallback → viewer (fail-closed).
	if !client.RoleFromInvite {
		// Invite TIDAK ada: gunakan role yang pernah tercatat (jika ada).
		if remembered := room.roles[client.UserID]; remembered != "" {
			client.RoomRole = remembered
		} else if isFirst {
			// Orang pertama di room tanpa undangan jadi owner.
			client.RoomRole = RoomRoleOwner
		}
	}
	// Jika setelah semua pengecekan di atas role masih kosong, paksa viewer.
	if client.RoomRole == "" {
		client.RoomRole = RoomRoleViewer
	}
	// Simpan role yang berlaku — termasuk yang datang dari invite — ke memori room.
	// Ini penting agar user yang baru di-upgrade ke editor tetap editor saat pindah halaman.
	if client.UserID != "" {
		room.roles[client.UserID] = client.RoomRole
	}
	room.clients[client] = true
	room.mu.Unlock()

	log.Printf("✅ Client %s joined room %s as %s (total: %d, tab_kedua=%v, from_invite=%v)", client.UserID, client.RoomID, client.RoomRole, room.GetClientCount(), alreadyPresent, client.RoleFromInvite)

	// Sync state to joining client
	client.sendQuiet(h.buildPresenceList(room))
	client.sendQuiet(h.buildStateSync(room, client))
	client.sendQuiet(h.buildFollowState(room))

	// Presence list ke room selalu di-refresh (dedupe by user)
	h.broadcastToRoom(room, h.buildPresenceList(room), client)

	// Tab kedua user yang sama: cukup refresh presence, jangan umumkan "bergabung" lagi
	if alreadyPresent {
		return
	}

	joinPayload := map[string]interface{}{
		"user_id":      client.UserID,
		"username":     client.Username,
		"role":         client.Role,
		"room_role":    client.RoomRole,
		"display_name": client.Username,
		"color":        colorForUser(client.UserID),
		"timestamp":    time.Now().Unix(),
	}

	h.broadcastToRoom(room, newMessage(MsgUserJoined, client.RoomID, client.UserID, client.Username, joinPayload), client)
	h.broadcastToRoom(room, newMessage(MsgPresenceJoined, client.RoomID, client.UserID, client.Username, joinPayload), client)
	h.broadcastToRoom(room, newMessage(MsgActivityLog, client.RoomID, client.UserID, client.Username, map[string]interface{}{
		"action":  "joined",
		"details": client.Username + " bergabung (" + client.RoomRole + ")",
	}), client)
}

// unregisterClient - keluarkan satu socket dari room.
//
// Event "keluar" hanya disiarkan bila itu socket TERAKHIR milik user tersebut —
// menutup satu dari dua tab bukan berarti orangnya pergi. Relasi follow pun baru
// dilepas kalau leader-nya benar-benar sudah tidak punya socket lagi.
func (h *Hub) unregisterClient(client *Client) {
	h.mu.RLock()
	room, exists := h.rooms[client.RoomID]
	h.mu.RUnlock()
	if !exists {
		return
	}

	room.mu.Lock()
	_, ok := room.clients[client]
	if ok {
		delete(room.clients, client)
		close(client.send)
	}
	remaining := len(room.clients)

	// Masih ada socket lain milik user yang sama? Berarti dia belum benar-benar keluar.
	stillOnline := false
	for other := range room.clients {
		if other.UserID != "" && other.UserID == client.UserID {
			stillOnline = true
			break
		}
	}
	var detachedFollowers []*Client
	if !stillOnline {
		// Leader benar-benar pergi — lepas follower + siapkan follow_stopped ke FE
		for other := range room.clients {
			if other.FollowingUserID == client.UserID {
				other.FollowingUserID = ""
				detachedFollowers = append(detachedFollowers, other)
			}
		}
	}
	room.mu.Unlock()

	if !ok {
		return
	}

	// Tab lain user ini masih terbuka: cukup refresh presence, jangan umumkan "keluar"
	if stillOnline {
		log.Printf("➖ Satu tab %s ditutup di room %s (socket tersisa: %d)", client.UserID, client.RoomID, remaining)
		h.broadcastToRoom(room, h.buildPresenceList(room), nil)
		return
	}

	log.Printf("❌ Client %s left room %s (remaining: %d)", client.UserID, client.RoomID, remaining)

	// Bubble cursor chat yang masih terbuka harus ikut hilang di layar peer —
	// tanpa ini bubble jadi "hantu" tertinggal selamanya kalau user disconnect
	// mendadak (tutup tab/refresh) di tengah mengetik.
	if client.ChatBubble != nil {
		h.broadcastToRoom(room, newMessage(MsgCursorChatClosed, client.RoomID, client.UserID, client.Username, map[string]interface{}{}), nil)
	}

	// Beritahu follower secara eksplisit agar banner "Following…" ikut hilang di FE
	for _, follower := range detachedFollowers {
		stopped := newMessage(MsgFollowStopped, client.RoomID, follower.UserID, follower.Username, map[string]interface{}{
			"follower_id": follower.UserID,
			"leader_id":   client.UserID,
			"reason":      "leader_left",
		})
		follower.sendQuiet(stopped)
	}

	leavePayload := map[string]interface{}{
		"user_id":   client.UserID,
		"username":  client.Username,
		"timestamp": time.Now().Unix(),
	}
	h.broadcastToRoom(room, newMessage(MsgUserLeft, client.RoomID, client.UserID, client.Username, leavePayload), nil)
	h.broadcastToRoom(room, newMessage(MsgPresenceLeft, client.RoomID, client.UserID, client.Username, leavePayload), nil)
	h.broadcastToRoom(room, newMessage(MsgActivityLog, client.RoomID, client.UserID, client.Username, map[string]interface{}{
		"action":  "left",
		"details": client.Username + " keluar",
	}), nil)
	h.broadcastToRoom(room, h.buildFollowState(room), nil)

	if remaining == 0 {
		log.Printf("🗑️  Room %s is now empty", client.RoomID)
	}
}

// broadcastMessage - simpan pesan ke history room lalu sebarkan ke semua anggota room
func (h *Hub) broadcastMessage(message *Message) {
	h.mu.RLock()
	room, exists := h.rooms[message.RoomID]
	h.mu.RUnlock()
	if !exists {
		return
	}

	room.addToHistory(message)
	h.broadcastToRoom(room, message, nil)
}

// BroadcastExcept sends to all clients in room except skip (nil = all including sender filtered by UserID skip logic below).
func (h *Hub) broadcastToRoom(room *Room, message *Message, skip *Client) {
	room.mu.RLock()
	defer room.mu.RUnlock()

	for client := range room.clients {
		if skip != nil && client == skip {
			continue
		}
		// Skip sender for activity broadcasts that carry UserID (unless message is directed to self via empty UserID filter)
		if skip == nil && message.UserID != "" && client.UserID == message.UserID {
			// Still allow presence_list / history / pong / error / state_sync to reach sender when UserID matches
			switch message.Type {
			case MsgPresenceList, MsgHistory, MsgPong, MsgError, MsgStateSync, MsgFollowState, MsgFollowStarted, MsgFollowStopped:
				// deliver
			default:
				continue
			}
		}

		select {
		case client.send <- message:
		default:
			log.Printf("⚠️  Skipped slow client %s in room %s", client.UserID, message.RoomID)
		}
	}
}

// buildPresenceList - susun daftar user aktif di room (unik per user_id) beserta warna & role-nya
func (h *Hub) buildPresenceList(room *Room) *Message {
	room.mu.RLock()
	defer room.mu.RUnlock()

	seen := make(map[string]bool, len(room.clients))
	users := make([]map[string]interface{}, 0, len(room.clients))
	for c := range room.clients {
		if c.UserID == "" || seen[c.UserID] {
			continue
		}
		seen[c.UserID] = true
		users = append(users, map[string]interface{}{
			"user_id":      c.UserID,
			"username":     c.Username,
			"display_name": c.Username,
			"role":         c.Role,
			"room_role":    c.RoomRole,
			"following":    c.FollowingUserID,
			"color":        colorForUser(c.UserID),
		})
	}
	return newMessage(MsgPresenceList, room.ID, "", "", map[string]interface{}{
		"users": users,
	})
}

// buildStateSync - susun snapshot state room (lock aktif + 30 pesan terakhir) untuk client yang baru join.
//
// Menyertakan blok "self" berisi identitas client itu sendiri. Frontend tidak boleh
// menebak user_id-nya dari auth store: store itu tidak dipersist, jadi setelah refresh
// atau di tab baru nilainya null dan avatar sendiri jadi ikut bisa di-Follow.
// Server adalah satu-satunya sumber kebenaran identitas di dalam room.
func (h *Hub) buildStateSync(room *Room, client *Client) *Message {
	return newMessage(MsgStateSync, room.ID, "", "", map[string]interface{}{
		"locks":          h.locks.Snapshot(),
		"history":        room.GetHistory(30),
		"canvas_strokes": room.GetCanvasStrokes(),
		"room_id":        room.ID,
		"self": map[string]interface{}{
			"user_id":      client.UserID,
			"username":     client.Username,
			"display_name": client.Username,
			"role":         client.Role,
			"room_role":    client.RoomRole,
			"color":        colorForUser(client.UserID),
		},
	})
}

// buildFollowState - susun graf follow room (pasangan follower -> leader) agar UI bisa gambar garis follow
func (h *Hub) buildFollowState(room *Room) *Message {
	room.mu.RLock()
	defer room.mu.RUnlock()

	pairs := make([]map[string]interface{}, 0)
	seen := make(map[string]bool)
	for c := range room.clients {
		if c.UserID == "" || c.FollowingUserID == "" || seen[c.UserID] {
			continue
		}
		seen[c.UserID] = true
		pairs = append(pairs, map[string]interface{}{
			"follower_id": c.UserID,
			"leader_id":   c.FollowingUserID,
		})
	}
	return newMessage(MsgFollowState, room.ID, "", "", map[string]interface{}{
		"follows": pairs,
	})
}

// findClientInRoom - cari client di room berdasarkan userID; nil kalau tidak ada
func (h *Hub) findClientInRoom(room *Room, userID string) *Client {
	room.mu.RLock()
	defer room.mu.RUnlock()
	for c := range room.clients {
		if c.UserID == userID {
			return c
		}
	}
	return nil
}

// broadcastToFollowers mengirim viewport_sync hanya ke client yang mengikuti leaderID.
func (h *Hub) broadcastToFollowers(roomID, leaderID string, message *Message) {
	h.mu.RLock()
	room, exists := h.rooms[roomID]
	h.mu.RUnlock()
	if !exists {
		return
	}
	room.mu.RLock()
	defer room.mu.RUnlock()
	for client := range room.clients {
		if client.FollowingUserID != leaderID {
			continue
		}
		select {
		case client.send <- message:
		default:
			log.Printf("⚠️  Skipped slow follower %s", client.UserID)
		}
	}
}

// ResolveRoomRole menentukan role join dengan urutan prioritas:
//
//  1. invite token yang valid — niat eksplisit pemilik room
//  2. role yang sudah pernah tercatat untuk user ini di room tsb
//  3. room masih kosong → owner
//  4. selain itu → editor
//
// Langkah 2 yang membuat viewer tetap viewer saat pindah halaman atau reconnect,
// meski query ?invite= sudah tidak ada lagi di URL.
func (h *Hub) ResolveRoomRole(roomID, userID, inviteToken string) string {
	if inviteToken != "" {
		if inv, ok := h.invites.Get(inviteToken); ok && inv.RoomID == roomID {
			return inv.Role
		}
	}
	h.mu.RLock()
	room, exists := h.rooms[roomID]
	h.mu.RUnlock()
	if !exists {
		return RoomRoleOwner // akan di-set owner di register bila first
	}
	if remembered := room.RememberedRole(userID); remembered != "" {
		return remembered
	}
	if room.GetClientCount() == 0 {
		return RoomRoleOwner
	}
	// Masuk hanya berbekal ?room= tanpa undangan = hanya boleh menonton.
	// Sebelumnya default-nya editor, sehingga link room yang diteruskan ke
	// siapa pun langsung memberi hak mengubah data — hak edit sekarang harus
	// datang dari undangan bertoken yang dibuat owner/editor.
	return RoomRoleViewer
}

// RoomRoleOf - role seseorang di sebuah room; "" kalau dia bukan anggota.
//
// Dipakai endpoint HTTP (invite/revoke) untuk memeriksa izin, karena di sana
// tidak ada *Client yang bisa ditanya.
func (h *Hub) RoomRoleOf(roomID, userID string) string {
	if roomID == "" || userID == "" {
		return ""
	}
	h.mu.RLock()
	room, exists := h.rooms[roomID]
	h.mu.RUnlock()
	if !exists {
		return ""
	}
	return room.RememberedRole(userID)
}

// cleanupInactiveRooms - hapus room yang sudah kosong (dipanggil ticker tiap 30 detik) agar memori tidak bocor
func (h *Hub) cleanupInactiveRooms() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for roomID, room := range h.rooms {
		room.mu.RLock()
		empty := len(room.clients) == 0
		room.mu.RUnlock()
		if empty {
			close(room.stopCh)
			delete(h.rooms, roomID)
			log.Printf("🗑️  Cleaned up empty room: %s", roomID)
		}
	}
}

// Stop stops the hub
func (h *Hub) Stop() {
	close(h.stopCh)
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, room := range h.rooms {
		close(room.stopCh)
	}
}

// GetRoomInfo returns information about a room
func (h *Hub) GetRoomInfo(roomID string) map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, exists := h.rooms[roomID]
	if !exists {
		return nil
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	users := make([]map[string]string, 0, len(room.clients))
	for client := range room.clients {
		users = append(users, map[string]string{
			"user_id":   client.UserID,
			"username":  client.Username,
			"role":      client.Role,
			"room_role": client.RoomRole,
			"color":     colorForUser(client.UserID),
		})
	}

	return map[string]interface{}{
		"room_id":      roomID,
		"client_count": len(room.clients),
		"users":        users,
		"locks":        h.locks.Snapshot(),
	}
}

// GetStats returns hub statistics
func (h *Hub) GetStats() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	totalClients := 0
	for _, room := range h.rooms {
		room.mu.RLock()
		totalClients += len(room.clients)
		room.mu.RUnlock()
	}

	return map[string]interface{}{
		"total_rooms":   len(h.rooms),
		"total_clients": totalClients,
		"active_locks":  len(h.locks.Snapshot()),
	}
}

// Publish broadcasts a message to a room (used by client handlers).
func (h *Hub) Publish(msg *Message) {
	select {
	case h.broadcast <- msg:
	default:
		log.Printf("⚠️  Broadcast channel full, dropping message type=%s", msg.Type)
	}
}

// colorForUser - tentukan warna kursor user secara deterministik dari hash userID (8 warna palette)
func colorForUser(userID string) string {
	palette := []string{
		"#E11D48", "#EA580C", "#CA8A04", "#16A34A",
		"#0891B2", "#2563EB", "#7C3AED", "#DB2777",
	}
	if userID == "" {
		return palette[0]
	}
	hash := 0
	for i := 0; i < len(userID); i++ {
		hash = (hash*31 + int(userID[i])) % len(palette)
	}
	if hash < 0 {
		hash = -hash
	}
	return palette[hash%len(palette)]
}
