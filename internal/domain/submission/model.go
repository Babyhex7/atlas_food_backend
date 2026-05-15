package submission

import (
	"time"
)

// SurveySubmission - model untuk tabel survey_submissions
// Menyimpan hasil recall makanan dari responden
type SurveySubmission struct {
	ID              string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	SurveyID        string    `gorm:"type:char(36);not null;index" json:"survey_id"`
	ParticipantID   *string   `gorm:"type:char(36);index" json:"participant_id"`
	RespondentName  string    `gorm:"type:varchar(255)" json:"respondent_name"`
	RespondentEmail string    `gorm:"type:varchar(255)" json:"respondent_email"`
	MealsData       string    `gorm:"type:json;not null" json:"meals_data"`    // JSON string
	MissingFoods    string    `gorm:"type:json" json:"missing_foods"`         // JSON string
	TotalEnergy     float64   `gorm:"type:decimal(10,2);default:0" json:"total_energy"`
	TotalProtein    float64   `gorm:"type:decimal(10,2);default:0" json:"total_protein"`
	TotalCarbs      float64   `gorm:"type:decimal(10,2);default:0" json:"total_carbs"`
	TotalFat        float64   `gorm:"type:decimal(10,2);default:0" json:"total_fat"`
	SubmittedAt     time.Time `gorm:"not null" json:"submitted_at"`
	CreatedAt       time.Time `json:"created_at"`
}

func (SurveySubmission) TableName() string {
	return "survey_submissions"
}
