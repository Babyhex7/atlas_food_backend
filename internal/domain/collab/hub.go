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
		stopCh:     make(chan struct{}),
	}
}

// Locks returns the in-memory lock manager.
func (h *Hub) Locks() *LockManager {
	return h.locks
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

func (h *Hub) registerClient(client *Client) {
	room := h.GetOrCreateRoom(client.RoomID)

	// Satu user = satu socket per room. Reconnect / Strict Mode sering buka
	// koneksi baru tanpa menutup lama — tanpa ini presence & activity "join"
	// menumpuk meski orangnya cuma satu.
	room.mu.Lock()
	replaced := false
	for existing := range room.clients {
		if existing.UserID != "" && existing.UserID == client.UserID && existing != client {
			delete(room.clients, existing)
			close(existing.send)
			replaced = true
		}
	}
	room.clients[client] = true
	room.mu.Unlock()

	log.Printf("✅ Client %s joined room %s (total: %d, reconnect=%v)", client.UserID, client.RoomID, room.GetClientCount(), replaced)

	// Sync state to joining client
	client.sendQuiet(h.buildPresenceList(room))
	client.sendQuiet(h.buildStateSync(room))

	// Presence list ke room selalu di-refresh (dedupe by user)
	h.broadcastToRoom(room, h.buildPresenceList(room), client)

	if replaced {
		return
	}

	joinPayload := map[string]interface{}{
		"user_id":      client.UserID,
		"username":     client.Username,
		"role":         client.Role,
		"display_name": client.Username,
		"color":        colorForUser(client.UserID),
		"timestamp":    time.Now().Unix(),
	}

	h.broadcastToRoom(room, newMessage(MsgUserJoined, client.RoomID, client.UserID, client.Username, joinPayload), client)
	h.broadcastToRoom(room, newMessage(MsgPresenceJoined, client.RoomID, client.UserID, client.Username, joinPayload), client)
	h.broadcastToRoom(room, newMessage(MsgActivityLog, client.RoomID, client.UserID, client.Username, map[string]interface{}{
		"action":  "joined",
		"details": client.Username + " bergabung",
	}), client)
}

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
	room.mu.Unlock()

	if !ok {
		return
	}

	log.Printf("❌ Client %s left room %s (remaining: %d)", client.UserID, client.RoomID, remaining)

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

	if remaining == 0 {
		log.Printf("🗑️  Room %s is now empty", client.RoomID)
	}
}

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
			case MsgPresenceList, MsgHistory, MsgPong, MsgError, MsgStateSync:
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
			"color":        colorForUser(c.UserID),
		})
	}
	return newMessage(MsgPresenceList, room.ID, "", "", map[string]interface{}{
		"users": users,
	})
}

func (h *Hub) buildStateSync(room *Room) *Message {
	return newMessage(MsgStateSync, room.ID, "", "", map[string]interface{}{
		"locks":    h.locks.Snapshot(),
		"history":  room.GetHistory(30),
		"room_id":  room.ID,
	})
}

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
			"user_id":  client.UserID,
			"username": client.Username,
			"role":     client.Role,
			"color":    colorForUser(client.UserID),
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
