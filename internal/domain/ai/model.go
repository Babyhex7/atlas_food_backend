package ai

import "time"

// AIResultLog - model untuk tabel ai_result_logs.
// Satu submission punya paling banyak satu baris; analisis ulang menimpa baris yang sama.
type AIResultLog struct {
	ID           string `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	SubmissionID string `gorm:"type:char(36);not null;uniqueIndex" json:"submission_id"`
	// InputPayload - pesan user yang dikirim ke LLM (tanpa identitas responden)
	InputPayload string `gorm:"type:json;not null" json:"input_payload"`
	// RawResponse - jawaban LLM apa adanya, untuk audit
	RawResponse string `gorm:"type:json;not null" json:"raw_response"`
	// ResultData - hasil akhir yang dikirim ke klien (narasi LLM + perhitungan server).
	// NULL pada baris lama yang dibuat sebelum kolom ini ada.
	ResultData       *string   `gorm:"type:json" json:"result_data,omitempty"`
	OverallStatus    string    `gorm:"type:enum('good','less','excess');not null" json:"overall_status"`
	ModelUsed        string    `gorm:"type:varchar(100);not null;default:''" json:"model_used"`
	PromptVersion    string    `gorm:"type:varchar(20);not null;default:'v1'" json:"prompt_version"`
	TokenUsed        *int      `gorm:"type:int" json:"token_used,omitempty"`
	PromptTokens     *int      `gorm:"type:int" json:"prompt_tokens,omitempty"`
	CompletionTokens *int      `gorm:"type:int" json:"completion_tokens,omitempty"`
	LatencyMs        *int      `gorm:"type:int" json:"latency_ms,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TableName - set nama tabel
func (AIResultLog) TableName() string {
	return "ai_result_logs"
}
