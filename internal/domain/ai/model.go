package ai

import "time"

// AIResultLog - model untuk tabel ai_result_logs
type AIResultLog struct {
	ID            string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	SubmissionID  string    `gorm:"type:char(36);not null;uniqueIndex" json:"submission_id"`
	InputPayload  string    `gorm:"type:json;not null" json:"input_payload"`
	RawResponse   string    `gorm:"type:json;not null" json:"raw_response"`
	OverallStatus string    `gorm:"type:enum('good','less','excess');not null" json:"overall_status"`
	ModelUsed     string    `gorm:"type:varchar(50);not null;default:'llama3-8b-8192'" json:"model_used"`
	TokenUsed     *int      `gorm:"type:int" json:"token_used,omitempty"`
	LatencyMs     *int      `gorm:"type:int" json:"latency_ms,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// TableName - set nama tabel
func (AIResultLog) TableName() string {
	return "ai_result_logs"
}
