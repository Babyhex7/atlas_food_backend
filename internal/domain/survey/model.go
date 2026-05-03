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
import (
	"time"
)

// Locale - model untuk tabel locales
type Locale struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code string `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"`
	Name string `gorm:"type:varchar(50);not null" json:"name"`
}

func (Locale) TableName() string {
	return "locales"
}
// TableName - set nama tabel untuk GORM
func (Locale) TableName() string {
	return "locales"
}

// Survey - model untuk tabel surveys
type Survey struct {
	ID           string     `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	Slug         string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Name         string     `gorm:"type:varchar(255);not null" json:"name"`
	Description  string     `gorm:"type:text" json:"description"`
	MealsConfig  string     `gorm:"type:json;not null" json:"meals_config"` // JSON string
	Prompts      string     `gorm:"type:json" json:"prompts,omitempty"`      // JSON string
	LocaleID     int        `gorm:"default:1" json:"locale_id"`
	Locale       Locale     `gorm:"foreignKey:LocaleID" json:"locale,omitempty"`
	StartDate    *time.Time `gorm:"type:date" json:"start_date,omitempty"`
	EndDate      *time.Time `gorm:"type:date" json:"end_date,omitempty"`
	Status       string     `gorm:"type:enum('draft','active','closed');default:'draft'" json:"status"`
	AccessToken  string     `gorm:"type:varchar(255);uniqueIndex" json:"-"` // json:"-" agar tidak terexpose
	CreatedBy    string     `gorm:"type:char(36);not null" json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TableName - set nama tabel untuk GORM
func (Survey) TableName() string {
	return "surveys"
}

// SurveyParticipant - model untuk tabel survey_participants
type SurveyParticipant struct {
	ID            string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	SurveyID      string    `gorm:"type:char(36);not null;index" json:"survey_id"`
	Survey        Survey    `gorm:"foreignKey:SurveyID" json:"survey,omitempty"`
	UserID        *string   `gorm:"type:char(36);index" json:"user_id,omitempty"`
	Alias         string    `gorm:"type:varchar(50)" json:"alias,omitempty"`
	IsAnonymous   bool      `gorm:"default:true" json:"is_anonymous"`
	CreatedAt     time.Time `json:"created_at"`
}

// TableName - set nama tabel untuk GORM
func (SurveyParticipant) TableName() string {
	return "survey_participants"
}

// MealsConfig - struct untuk parsing meals_config JSON
type MealsConfig struct {
	Meals []MealConfig `json:"meals"`
}

// MealConfig - konfigurasi satu meal
type MealConfig struct {
	Name     string `json:"name"`
	Time     string `json:"time"`
	Required bool   `json:"required"`
	Prompt   string `json:"prompt,omitempty"`
}

// PromptsConfig - struct untuk parsing prompts JSON
type PromptsConfig struct {
	BeforeMeals string `json:"before_meals,omitempty"`
	AfterMeals  string `json:"after_meals,omitempty"`
	MissingFood string `json:"missing_food,omitempty"`
}
