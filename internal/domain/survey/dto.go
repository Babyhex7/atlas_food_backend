package survey

// MealConfig - struktur meal untuk surveys.meals_config
type MealConfig struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Time     string `json:"time" binding:"required"`
	Required bool   `json:"required"`
	Prompt   string `json:"prompt,omitempty"`
}

// MealsConfig - wrapper untuk list meals
type MealsConfig struct {
	Meals []MealConfig `json:"meals"`
}

// PromptsConfig - konfigurasi prompt survey
type PromptsConfig struct {
	BeforeMeals string `json:"before_meals,omitempty"`
	AfterMeals  string `json:"after_meals,omitempty"`
	MissingFood string `json:"missing_food,omitempty"`
}

// CreateSurveyRequest - DTO request create survey
type CreateSurveyRequest struct {
	Name        string        `json:"name" binding:"required,min=3,max=255"`
	Slug        string        `json:"slug" binding:"required,min=3,max=100"`
	Description string        `json:"description"`
	MealsConfig MealsConfig   `json:"meals_config" binding:"required"`
	Prompts     PromptsConfig `json:"prompts,omitempty"`
	LocaleID    int           `json:"locale_id" binding:"omitempty,min=1"`
	StartDate   *string       `json:"start_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	EndDate     *string       `json:"end_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	Status      string        `json:"status,omitempty" binding:"omitempty,oneof=draft active closed"`
}

// UpdateSurveyRequest - DTO request update survey
type UpdateSurveyRequest struct {
	Name        string         `json:"name" binding:"omitempty,min=3,max=255"`
	Description string         `json:"description"`
	MealsConfig *MealsConfig   `json:"meals_config,omitempty"`
	Prompts     *PromptsConfig `json:"prompts,omitempty"`
	LocaleID    int            `json:"locale_id" binding:"omitempty,min=1"`
	StartDate   *string        `json:"start_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	EndDate     *string        `json:"end_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	Status      string         `json:"status,omitempty" binding:"omitempty,oneof=draft active closed"`
}

// SurveyResponse - DTO response survey
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
	AccessURL   string        `json:"access_url,omitempty"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
}

// LocaleInfo - info locale untuk response
type LocaleInfo struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// ListSurveysResponse - DTO untuk list surveys
type ListSurveysResponse struct {
	ID               string  `json:"id"`
	Slug             string  `json:"slug"`
	Name             string  `json:"name"`
	Status           string  `json:"status"`
	StartDate        *string `json:"start_date,omitempty"`
	EndDate          *string `json:"end_date,omitempty"`
	ParticipantCount int     `json:"participant_count"`
	CreatedAt        string  `json:"created_at"`
}

// SurveyListResponse - wrapper untuk list pagination
type SurveyListResponse struct {
	Surveys []ListSurveysResponse `json:"surveys"`
	Total   int64                 `json:"total"`
	Page    int                   `json:"page"`
	Limit   int                   `json:"limit"`
}

// AccessSurveyRequest - DTO untuk respondent access survey
type AccessSurveyRequest struct {
	Token          string `json:"token" binding:"required"`
	Alias          string `json:"alias,omitempty"`
	RespondentName string `json:"respondent_name,omitempty"`
}

// AccessSurveyResponse - response setelah access
type AccessSurveyResponse struct {
	Survey      PublicSurveyResponse `json:"survey"`
	Participant ParticipantResponse  `json:"participant"`
	AccessToken string               `json:"access_token"`
}

// ParticipantResponse - detail participant untuk response
type ParticipantResponse struct {
	ID    string `json:"id"`
	Alias string `json:"alias"`
}

// PublicSurveyResponse - response untuk public access
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
