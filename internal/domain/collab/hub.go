package collab

import (
	"log"
	"sync"
	"time"
)

// Hub manages all WebSocket rooms and clients
// Optimized with goroutines and channels for concurrent operations
type Hub struct {
	// Registered rooms
	rooms map[string]*Room

	// Register requests from rooms
	register chan *Client

	// Unregister requests from rooms
	unregister chan *Client

	// Broadcast messages
	broadcast chan *Message

	// Mutex for thread-safe room operations
	mu sync.RWMutex

	// Stop channel
	stopCh chan struct{}
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]*Room),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		broadcast:  make(chan *Message, 1024),
		stopCh:     make(chan struct{}),
	}
}

// Run starts the hub's main event loop
// Handles all WebSocket operations concurrently
func (h *Hub) Run() {
	// Cleanup goroutine - removes inactive rooms every 30 seconds
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

	// Main event loop
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
		
		// Start room in goroutine
		go room.Run()
		
		log.Printf("📡 Created new room: %s", roomID)
	}

	return room
}

// registerClient registers a client to a room
func (h *Hub) registerClient(client *Client) {
	room := h.GetOrCreateRoom(client.RoomID)
	room.clients[client] = true
	
	log.Printf("✅ Client %s joined room %s (total: %d)", client.UserID, client.RoomID, len(room.clients))
	
	// Broadcast join notification to other users in room
	joinMsg := &Message{
		Type:    "user_joined",
		RoomID:  client.RoomID,
		UserID:  client.UserID,
		Payload: map[string]interface{}{
			"user_id":   client.UserID,
			"username":  client.Username,
			"timestamp": time.Now().Unix(),
		},
		Timestamp: time.Now(),
	}
	
	// Send to all clients in room except the joining client
	for c := range room.clients {
		if c != client {
			select {
			case c.send <- joinMsg:
			default:
				// Client's send channel is full, skip
			}
		}
	}
}

// unregisterClient removes a client from a room
func (h *Hub) unregisterClient(client *Client) {
	h.mu.RLock()
	room, exists := h.rooms[client.RoomID]
	h.mu.RUnlock()
	
	if !exists {
		return
	}

	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
		close(client.send)
		
		log.Printf("❌ Client %s left room %s (remaining: %d)", client.UserID, client.RoomID, len(room.clients))
		
		// Broadcast leave notification
		leaveMsg := &Message{
			Type:    "user_left",
			RoomID:  client.RoomID,
			UserID:  client.UserID,
			Payload: map[string]interface{}{
				"user_id":   client.UserID,
				"username":  client.Username,
				"timestamp": time.Now().Unix(),
			},
			Timestamp: time.Now(),
		}
		
		// Send to remaining clients
		for c := range room.clients {
			select {
			case c.send <- leaveMsg:
			default:
				// Client's send channel is full, skip
			}
		}
		
		// If room is empty, mark for cleanup
		if len(room.clients) == 0 {
			log.Printf("🗑️  Room %s is now empty", client.RoomID)
		}
	}
}

// broadcastMessage sends a message to all clients in a room
// Optimized: Non-blocking sends to avoid slow clients blocking others
func (h *Hub) broadcastMessage(message *Message) {
	h.mu.RLock()
	room, exists := h.rooms[message.RoomID]
	h.mu.RUnlock()
	
	if !exists {
		return
	}

	// Broadcast to all clients in room concurrently
	var wg sync.WaitGroup
	for client := range room.clients {
		// Skip sender if message has sender info
		if message.UserID != "" && client.UserID == message.UserID {
			continue
		}

		wg.Add(1)
		go func(c *Client) {
			defer wg.Done()
			
			select {
			case c.send <- message:
				// Message sent successfully
			case <-time.After(100 * time.Millisecond):
				// Client is too slow, skip
				log.Printf("⚠️  Skipped slow client %s in room %s", c.UserID, message.RoomID)
			}
		}(client)
	}
	
	// Wait for all sends to complete (with timeout)
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		// All sends completed
	case <-time.After(500 * time.Millisecond):
		// Timeout after 500ms
		log.Printf("⚠️  Broadcast timeout for room %s", message.RoomID)
	}
}

// cleanupInactiveRooms removes empty rooms
func (h *Hub) cleanupInactiveRooms() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for roomID, room := range h.rooms {
		if len(room.clients) == 0 {
			// Stop room
			close(room.stopCh)
			delete(h.rooms, roomID)
			log.Printf("🗑️  Cleaned up empty room: %s", roomID)
		}
	}
}

// Stop stops the hub
func (h *Hub) Stop() {
	close(h.stopCh)
	
	// Stop all rooms
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

	users := make([]map[string]string, 0, len(room.clients))
	for client := range room.clients {
		users = append(users, map[string]string{
			"user_id":  client.UserID,
			"username": client.Username,
		})
	}

	return map[string]interface{}{
		"room_id":     roomID,
		"client_count": len(room.clients),
		"users":       users,
	}
}

// GetStats returns hub statistics
func (h *Hub) GetStats() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	totalClients := 0
	for _, room := range h.rooms {
		totalClients += len(room.clients)
	}

	return map[string]interface{}{
		"total_rooms":   len(h.rooms),
		"total_clients": totalClients,
	}
}
