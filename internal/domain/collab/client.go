package collab

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 8192 // 8KB

	// Send buffer size
	sendBufferSize = 256
)

// Client represents a WebSocket client
type Client struct {
	// WebSocket connection
	conn *websocket.Conn

	// Hub reference
	hub *Hub

	// Room ID
	RoomID string

	// User information
	UserID   string
	Username string

	// Buffered channel of outbound messages
	send chan *Message

	// Rate limiting
	lastMessageTime time.Time
	messageCount    int

	// Stop channel
	stopCh chan struct{}
}

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn, hub *Hub, roomID, userID, username string) *Client {
	return &Client{
		conn:     conn,
		hub:      hub,
		RoomID:   roomID,
		UserID:   userID,
		Username: username,
		send:     make(chan *Message, sendBufferSize),
		stopCh:   make(chan struct{}),
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub
// Runs in a goroutine per client
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

			// Rate limiting: max 10 messages per second
			if !c.checkRateLimit() {
				log.Printf("⚠️  Rate limit exceeded for user %s", c.UserID)
				continue
			}

			// Parse message
			var rawMsg map[string]interface{}
			if err := json.Unmarshal(messageBytes, &rawMsg); err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}

			// Create message object
			msgType, _ := rawMsg["type"].(string)
			payload, _ := rawMsg["payload"].(map[string]interface{})

			msg := &Message{
				Type:      msgType,
				RoomID:    c.RoomID,
				UserID:    c.UserID,
				Username:  c.Username,
				Payload:   payload,
				Timestamp: time.Now(),
			}

			// Handle message based on type
			c.handleMessage(msg)
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
// Runs in a goroutine per client
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
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Convert message to JSON
			messageJSON, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			// Write message
			if err := c.conn.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
				log.Printf("Error writing message: %v", err)
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

// handleMessage processes incoming messages
func (c *Client) handleMessage(msg *Message) {
	switch msg.Type {
	case "food_search":
		// Handle food search
		c.handleFoodSearch(msg)
	
	case "food_select":
		// Handle food selection
		c.handleFoodSelect(msg)
	
	case "portion_select":
		// Handle portion selection
		c.handlePortionSelect(msg)
	
	case "chat_message":
		// Handle chat message
		c.handleChatMessage(msg)
	
	case "get_history":
		// Send room history to this client
		c.sendRoomHistory()
	
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// handleFoodSearch handles food search requests
func (c *Client) handleFoodSearch(msg *Message) {
	// Broadcast search activity to other users
	broadcastMsg := &Message{
		Type:      "user_searching",
		RoomID:    c.RoomID,
		UserID:    c.UserID,
		Username:  c.Username,
		Payload:   msg.Payload,
		Timestamp: time.Now(),
	}
	
	c.hub.broadcast <- broadcastMsg
}

// handleFoodSelect handles food selection
func (c *Client) handleFoodSelect(msg *Message) {
	// Broadcast food selection to other users
	broadcastMsg := &Message{
		Type:      "food_selected",
		RoomID:    c.RoomID,
		UserID:    c.UserID,
		Username:  c.Username,
		Payload:   msg.Payload,
		Timestamp: time.Now(),
	}
	
	c.hub.broadcast <- broadcastMsg
}

// handlePortionSelect handles portion selection
func (c *Client) handlePortionSelect(msg *Message) {
	// Broadcast portion selection to other users
	broadcastMsg := &Message{
		Type:      "portion_selected",
		RoomID:    c.RoomID,
		UserID:    c.UserID,
		Username:  c.Username,
		Payload:   msg.Payload,
		Timestamp: time.Now(),
	}
	
	c.hub.broadcast <- broadcastMsg
}

// handleChatMessage handles chat messages
func (c *Client) handleChatMessage(msg *Message) {
	// Broadcast chat message to other users
	broadcastMsg := &Message{
		Type:      "chat_message",
		RoomID:    c.RoomID,
		UserID:    c.UserID,
		Username:  c.Username,
		Payload:   msg.Payload,
		Timestamp: time.Now(),
	}
	
	c.hub.broadcast <- broadcastMsg
}

// sendRoomHistory sends room history to this client
func (c *Client) sendRoomHistory() {
	room := c.hub.GetOrCreateRoom(c.RoomID)
	history := room.GetHistory(50) // Last 50 messages

	historyMsg := &Message{
		Type:   "history",
		RoomID: c.RoomID,
		Payload: map[string]interface{}{
			"messages": history,
		},
		Timestamp: time.Now(),
	}

	select {
	case c.send <- historyMsg:
	default:
		log.Printf("Failed to send history to client %s", c.UserID)
	}
}

// checkRateLimit checks if the client is sending messages too fast
func (c *Client) checkRateLimit() bool {
	now := time.Now()
	
	// Reset counter if more than 1 second has passed
	if now.Sub(c.lastMessageTime) > time.Second {
		c.messageCount = 0
		c.lastMessageTime = now
	}

	c.messageCount++
	
	// Max 10 messages per second
	return c.messageCount <= 10
}

// Stop stops the client
func (c *Client) Stop() {
	close(c.stopCh)
}
