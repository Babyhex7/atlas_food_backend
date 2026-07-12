package collab

import (
	"log"
	"sync"
	"time"
)

// Room represents a WebSocket room where users collaborate
type Room struct {
	// Room identifier
	ID string

	// Registered clients in this room
	clients map[*Client]bool

	// Hub reference
	hub *Hub

	// Message history (ring buffer for last 100 messages)
	messageHistory []*Message
	historyIndex   int
	historyMu      sync.RWMutex

	// Message batching
	batchQueue    []*Message
	batchMu       sync.Mutex
	batchTicker   *time.Ticker
	lastBatchSent time.Time

	// Stop channel
	stopCh chan struct{}
}

// NewRoom creates a new Room
func NewRoom(id string, hub *Hub) *Room {
	return &Room{
		ID:             id,
		clients:        make(map[*Client]bool),
		hub:            hub,
		messageHistory: make([]*Message, 100), // Ring buffer size
		batchQueue:     make([]*Message, 0, 50),
		batchTicker:    time.NewTicker(100 * time.Millisecond), // Batch every 100ms
		stopCh:         make(chan struct{}),
	}
}

// Run starts the room's message batching loop
func (r *Room) Run() {
	defer r.batchTicker.Stop()

	for {
		select {
		case <-r.batchTicker.C:
			r.flushBatchedMessages()
		case <-r.stopCh:
			return
		}
	}
}

// AddMessage adds a message to the room (with batching for optimization)
func (r *Room) AddMessage(msg *Message) {
	// Add to history
	r.addToHistory(msg)

	// Check if we should batch or send immediately
	shouldBatch := r.shouldBatchMessage(msg)
	
	if shouldBatch {
		r.batchMu.Lock()
		r.batchQueue = append(r.batchQueue, msg)
		queueSize := len(r.batchQueue)
		r.batchMu.Unlock()

		// Flush if batch is full (50 messages)
		if queueSize >= 50 {
			r.flushBatchedMessages()
		}
	} else {
		// Send immediately for important messages
		r.hub.broadcast <- msg
	}
}

// shouldBatchMessage determines if a message should be batched
func (r *Room) shouldBatchMessage(msg *Message) bool {
	// Don't batch these important message types
	switch msg.Type {
	case "user_joined", "user_left", "error":
		return false
	default:
		return true
	}
}

// flushBatchedMessages sends all batched messages
func (r *Room) flushBatchedMessages() {
	r.batchMu.Lock()
	
	if len(r.batchQueue) == 0 {
		r.batchMu.Unlock()
		return
	}

	// Copy and clear queue
	messages := make([]*Message, len(r.batchQueue))
	copy(messages, r.batchQueue)
	r.batchQueue = r.batchQueue[:0]
	r.lastBatchSent = time.Now()
	
	r.batchMu.Unlock()

	// Send all messages concurrently
	var wg sync.WaitGroup
	for _, msg := range messages {
		wg.Add(1)
		go func(m *Message) {
			defer wg.Done()
			r.hub.broadcast <- m
		}(msg)
	}
	wg.Wait()

	log.Printf("📤 Flushed %d batched messages for room %s", len(messages), r.ID)
}

// addToHistory adds message to ring buffer history
func (r *Room) addToHistory(msg *Message) {
	r.historyMu.Lock()
	defer r.historyMu.Unlock()

	r.messageHistory[r.historyIndex] = msg
	r.historyIndex = (r.historyIndex + 1) % len(r.messageHistory)
}

// GetHistory returns recent message history
func (r *Room) GetHistory(limit int) []*Message {
	r.historyMu.RLock()
	defer r.historyMu.RUnlock()

	if limit > len(r.messageHistory) {
		limit = len(r.messageHistory)
	}

	history := make([]*Message, 0, limit)
	
	// Get messages in chronological order
	for i := 0; i < len(r.messageHistory); i++ {
		idx := (r.historyIndex - len(r.messageHistory) + i + len(r.messageHistory)) % len(r.messageHistory)
		if r.messageHistory[idx] != nil {
			history = append(history, r.messageHistory[idx])
		}
	}

	// Return last N messages
	if len(history) > limit {
		history = history[len(history)-limit:]
	}

	return history
}

// BroadcastToRoom sends a message to all clients in this room
func (r *Room) BroadcastToRoom(msg *Message) {
	r.AddMessage(msg)
}

// GetClientCount returns number of connected clients
func (r *Room) GetClientCount() int {
	return len(r.clients)
}

// Stop stops the room
func (r *Room) Stop() {
	close(r.stopCh)
}
