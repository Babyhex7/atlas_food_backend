package collab

import "time"

// Message represents a WebSocket message
type Message struct {
	Type      string                 `json:"type"`
	RoomID    string                 `json:"room_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Username  string                 `json:"username,omitempty"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// Message types:
// - user_joined: User joined the room
// - user_left: User left the room
// - user_searching: User is searching for food
// - food_selected: User selected a food
// - portion_selected: User selected a portion
// - chat_message: Chat message
// - history: Room history
// - error: Error message
