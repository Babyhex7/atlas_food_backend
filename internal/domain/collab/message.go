package collab

import "time"

// MessageType merupakan tipe event websocket.
type MessageType string

const (

	// ===========================
	// Connection
	// ===========================

	MessageJoinRoom MessageType = "join_room"

	MessageLeaveRoom MessageType = "leave_room"

	MessageJoinSuccess MessageType = "join_success"

	MessageDisconnect MessageType = "disconnect"

	// ===========================
	// Health Check
	// ===========================

	MessagePing MessageType = "ping"

	MessagePong MessageType = "pong"

	// ===========================
	// Broadcast
	// ===========================

	MessageBroadcast MessageType = "broadcast"

	// ===========================
	// Error
	// ===========================

	MessageError MessageType = "error"
)

// Message merupakan format komunikasi websocket.
type Message struct {

	// Jenis event
	Type MessageType `json:"type"`

	// Room tujuan
	RoomID string `json:"room_id,omitempty"`

	// User pengirim
	UserID string `json:"user_id,omitempty"`

	// Isi data
	Payload any `json:"payload,omitempty"`

	// Waktu event
	Timestamp time.Time `json:"timestamp"`
}
