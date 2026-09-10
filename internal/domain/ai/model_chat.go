package ai

import "time"

// ChatSession - model untuk tabel ai_chat_sessions
type ChatSession struct {
	ID        string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	UserID    string    `gorm:"type:char(36);not null" json:"user_id"`
	Title     string    `gorm:"type:varchar(200)" json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ChatSession) TableName() string {
	return "ai_chat_sessions"
}

// ChatMessage - model untuk tabel ai_chat_messages
type ChatMessage struct {
	ID         string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	SessionID  string    `gorm:"type:char(36);not null" json:"session_id"`
	Role       string    `gorm:"type:enum('user','assistant','system');not null" json:"role"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	TokenUsed  int       `gorm:"type:int;default:0" json:"token_used"`
	ModelUsed  string    `gorm:"type:varchar(100)" json:"model_used"`
	LatencyMs  int       `gorm:"type:int;default:0" json:"latency_ms"`
	CreatedAt  time.Time `json:"created_at"`
}

func (ChatMessage) TableName() string {
	return "ai_chat_messages"
}
