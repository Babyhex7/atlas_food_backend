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
// CreateSurveyRequest - DTO untuk create survey
type CreateSurveyRequest struct {
	Name        string        `json:"name" binding:"required,min=3,max=255"`
	Slug        string        `json:"slug" binding:"required,min=3,max=100"`
	Description string        `json:"description"`
	MealsConfig MealsConfig   `json:"meals_config" binding:"required"`
	Prompts     PromptsConfig `json:"prompts,omitempty"`
	LocaleID    int           `json:"locale_id" binding:"omitempty,min=1"`
	StartDate   *string       `json:"start_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	EndDate     *string       `json:"end_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
}

// UpdateSurveyRequest - DTO untuk update survey
type UpdateSurveyRequest struct {
	Name        string        `json:"name" binding:"omitempty,min=3,max=255"`
	Description string        `json:"description"`
	MealsConfig *MealsConfig  `json:"meals_config,omitempty"`
	Prompts     PromptsConfig `json:"prompts,omitempty"`
	LocaleID    int           `json:"locale_id" binding:"omitempty,min=1"`
	StartDate   *string       `json:"start_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	EndDate     *string       `json:"end_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	Status      string        `json:"status,omitempty" binding:"omitempty,oneof=draft active closed"`
}

// SurveyResponse - DTO untuk response survey
type SurveyResponse struct {
	ID          string        `json:"id"`
	Slug        string        `json:"slug"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	MealsConfig MealsConfig   `json:"meals_config"`
	Prompts     PromptsConfig `json:"prompts,omitempty"`
	Locale      LocaleInfo    `json:"locale"`
	StartDate   *string       `json:"start_date,omitempty"`
	EndDate     *string       `json:"end_date,omitempty"`
	Status      string        `json:"status"`
	CreatedBy   string        `json:"created_by"`
	AccessURL   string        `json:"access_url,omitempty"` // URL publik untuk akses survey
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
}

// LocaleInfo - info locale untuk response
type LocaleInfo struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// ListSurveysResponse - DTO untuk list surveys (tanpa detail lengkap)
type ListSurveysResponse struct {
	ID               string  `json:"id"`
	Slug             string  `json:"slug"`
	Name             string  `json:"name"`
	Status           string  `json:"status"`
	StartDate        *string `json:"start_date,omitempty"`
	EndDate          *string `json:"end_date,omitempty"`
	ParticipantCount int     `json:"participant_count,omitempty"`
	CreatedAt        string  `json:"created_at"`
}

// SurveyListResponse - wrapper untuk list
type SurveyListResponse struct {
	Surveys []ListSurveysResponse `json:"surveys"`
	Total   int64                 `json:"total"`
	Page    int                   `json:"page"`
	Limit   int                   `json:"limit"`
}

// JoinSurveyRequest - DTO untuk respondent join survey
type JoinSurveyRequest struct {
	Alias string `json:"alias,omitempty"` // Nama alias (optional)
}

// JoinSurveyResponse - response setelah join
type JoinSurveyResponse struct {
	ParticipantID string        `json:"participant_id"`
	SurveyID      string        `json:"survey_id"`
	SurveyName    string        `json:"survey_name"`
	MealsConfig   MealsConfig   `json:"meals_config"`
	Prompts       PromptsConfig `json:"prompts"`
	AccessToken   string        `json:"access_token,omitempty"`
}

// PublicSurveyResponse - response untuk public access (tanpa sensitive data)
type PublicSurveyResponse struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	MealsConfig MealsConfig   `json:"meals_config"`
	Prompts     PromptsConfig `json:"prompts"`
	Locale      LocaleInfo    `json:"locale"`
	StartDate   *string       `json:"start_date,omitempty"`
	EndDate     *string       `json:"end_date,omitempty"`
	Status      string        `json:"status"`
}

// GenerateAccessTokenRequest - request regenerate token
type GenerateAccessTokenRequest struct {
	SurveyID string `json:"survey_id" binding:"required"`
}

// AccessTokenResponse - response dengan token baru
type AccessTokenResponse struct {
	SurveyID    string `json:"survey_id"`
	AccessToken string `json:"access_token"`
	AccessURL   string `json:"access_url"`
}

// CloneSurveyRequest - request clone survey
type CloneSurveyRequest struct {
	NewName string `json:"new_name" binding:"required,min=3,max=255"`
	NewSlug string `json:"new_slug" binding:"required,min=3,max=100"`
}
