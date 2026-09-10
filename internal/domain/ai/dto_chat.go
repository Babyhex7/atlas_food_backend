package ai

import "time"

// ChatRequest - DTO untuk user mengirim pesan
type ChatRequest struct {
	SessionID *string `json:"session_id"`
	Message   string  `json:"message" binding:"required"`
}

// ChatResponse - DTO untuk response bot
type ChatResponse struct {
	SessionID string `json:"session_id"`
	Reply     string `json:"reply"`
	XPEarned  int    `json:"xp_earned"`
	ModelUsed string `json:"model_used"`
	LatencyMs int    `json:"latency_ms"`
}

// SessionListResponse - DTO untuk list of sessions
type SessionListResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

// MessageListResponse - DTO untuk list of messages in a session
type MessageListResponse struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
