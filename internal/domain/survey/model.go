package survey

import "time"

// Survey - model untuk tabel surveys
type Survey struct {
	ID          string     `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	Slug        string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	MealsConfig []byte     `gorm:"type:json;not null" json:"meals_config"`
	Prompts     []byte     `gorm:"type:json" json:"prompts"`
	LocaleID    int        `gorm:"default:1" json:"locale_id"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Status      string     `gorm:"type:enum('draft','active','closed');default:'draft'" json:"status"`
	AccessToken string     `gorm:"type:varchar(255)" json:"access_token"`
	CreatedBy   string     `gorm:"type:char(36);not null;index" json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (Survey) TableName() string {
	return "surveys"
}

// SurveyParticipant - model untuk tabel survey_participants
type SurveyParticipant struct {
	ID          string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	SurveyID    string    `gorm:"type:char(36);not null;index" json:"survey_id"`
	UserID      *string   `gorm:"type:char(36);index" json:"user_id"`
	Alias       *string   `gorm:"type:varchar(50)" json:"alias"`
	IsAnonymous bool      `gorm:"default:true" json:"is_anonymous"`
	CreatedAt   time.Time `json:"created_at"`
}

func (SurveyParticipant) TableName() string {
	return "survey_participants"
}

// Locale - model untuk tabel locales
type Locale struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code string `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"`
	Name string `gorm:"type:varchar(50);not null" json:"name"`
}

func (Locale) TableName() string {
	return "locales"
}
