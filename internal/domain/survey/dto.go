package survey

import "encoding/json"

// MealConfig - struktur meal untuk surveys.meals_config
type MealConfig struct {
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Time  string `json:"time" binding:"required"`
	Order int    `json:"order"`
}

// CreateSurveyRequest - DTO request create survey
type CreateSurveyRequest struct {
	Slug        string            `json:"slug" binding:"omitempty,min=3,max=100"`
	Name        string            `json:"name" binding:"required,min=3,max=255"`
	Description string            `json:"description" binding:"omitempty,max=1000"`
	MealsConfig []MealConfig      `json:"meals_config" binding:"required,min=1,dive"`
	Prompts     map[string]string `json:"prompts"`
	LocaleID    int               `json:"locale_id" binding:"omitempty,min=1"`
	StartDate   string            `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate     string            `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
	Status      string            `json:"status" binding:"omitempty,oneof=draft active closed"`
}

// UpdateSurveyRequest - DTO request update survey
type UpdateSurveyRequest struct {
	Name        *string           `json:"name" binding:"omitempty,min=3,max=255"`
	Description *string           `json:"description" binding:"omitempty,max=1000"`
	MealsConfig []MealConfig      `json:"meals_config" binding:"omitempty,min=1,dive"`
	Prompts     map[string]string `json:"prompts"`
	LocaleID    *int              `json:"locale_id" binding:"omitempty,min=1"`
	StartDate   *string           `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate     *string           `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
	Status      *string           `json:"status" binding:"omitempty,oneof=draft active closed"`
}

// SurveyResponse - DTO response survey
type SurveyResponse struct {
	ID          string          `json:"id"`
	Slug        string          `json:"slug"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	MealsConfig json.RawMessage `json:"meals_config"`
	Prompts     json.RawMessage `json:"prompts,omitempty"`
	LocaleID    int             `json:"locale_id"`
	StartDate   *string         `json:"start_date,omitempty"`
	EndDate     *string         `json:"end_date,omitempty"`
	Status      string          `json:"status"`
	AccessToken string          `json:"access_token"`
	CreatedBy   string          `json:"created_by"`
	CreatedAt   string          `json:"created_at"`
}
