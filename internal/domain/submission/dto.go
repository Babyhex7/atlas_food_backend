package submission

import "encoding/json"

// SubmitSurveyRequest - DTO untuk input hasil survey
type SubmitSurveyRequest struct {
	SurveyID        string            `json:"survey_id" binding:"required"`
	ParticipantID   string            `json:"participant_id"`
	RespondentName  string            `json:"respondent_name"`
	RespondentEmail string            `json:"respondent_email"`
	MealsData       []MealData        `json:"meals_data" binding:"required"`
	DailyTotal      DailyTotal        `json:"daily_total"`
	MissingFoods    []MissingFoodData `json:"missing_foods"`
}

// MealData - detail makanan per waktu makan
type MealData struct {
	Name      string     `json:"name" binding:"required"`
	Time      string     `json:"time"`
	Foods     []FoodData `json:"foods" binding:"required"`
	MealTotal DailyTotal `json:"meal_total"`
}

// FoodData - detail makanan yang dikonsumsi
type FoodData struct {
	FoodID      string         `json:"food_id"`
	FoodName    string         `json:"food_name" binding:"required"`
	PortionGram float64        `json:"portion_gram" binding:"required"`
	Portion     PortionDetails `json:"portion"`
	Nutrients   NutrientValues `json:"nutrients"`
}

// PortionDetails - detail cara ukur porsi
type PortionDetails struct {
	Method        string  `json:"method"`
	ImageID       string  `json:"image_id"`
	ImageLabel    string  `json:"image_label"`
	BaseWeight    float64 `json:"base_weight"`
	Quantity      float64 `json:"quantity"`
	Fraction      float64 `json:"fraction"`
	TotalQuantity float64 `json:"total_quantity"`
}

// NutrientValues - nilai gizi makanan
type NutrientValues struct {
	Energy  float64 `json:"energy"`
	Protein float64 `json:"protein"`
	Carbs   float64 `json:"carbs"`
	Fat     float64 `json:"fat"`
}

// DailyTotal - total gizi harian
type DailyTotal struct {
	Energy  float64 `json:"energy"`
	Protein float64 `json:"protein"`
	Carbs   float64 `json:"carbs"`
	Fat     float64 `json:"fat"`
}

// MissingFoodData - data makanan yang tidak ada di DB
type MissingFoodData struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// SubmissionResponse - response setelah submit
type SubmissionResponse struct {
	SubmissionID string `json:"submission_id"`
	Message      string `json:"message"`
}

// ListSubmissionResponse - DTO untuk list submission di admin
type ListSubmissionResponse struct {
	ID             string  `json:"id"`
	RespondentName string  `json:"respondent_name"`
	SubmittedAt    string  `json:"submitted_at"`
	MealCount      int     `json:"meal_count"`
	TotalFoods     int     `json:"total_foods"`
	TotalEnergy    float64 `json:"total_energy"`
}

// SubmissionDetailResponse - detail submission untuk admin
type SubmissionDetailResponse struct {
	ID              string          `json:"id"`
	SurveyID        string          `json:"survey_id"`
	RespondentName  string          `json:"respondent_name"`
	RespondentEmail string          `json:"respondent_email"`
	MealsData       json.RawMessage `json:"meals_data"`
	MissingFoods    json.RawMessage `json:"missing_foods"`
	DailyTotal      DailyTotal      `json:"daily_total"`
	SubmittedAt     string          `json:"submitted_at"`
}
