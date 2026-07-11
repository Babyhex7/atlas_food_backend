package ai

import (
	"atlas_food/internal/domain/submission"
	"encoding/json"
)

// NutritionAnalysisRequest - request body endpoint AI
type NutritionAnalysisRequest struct {
	SubmissionID string `json:"submission_id" binding:"required"`
}

// NutritionAnalysisItem - item analisis nutrisi
type NutritionAnalysisItem struct {
	Label       string `json:"label"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// HealthInsight - insight kesehatan
type HealthInsight struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// NutritionAnalysisData - payload response AI
type NutritionAnalysisData struct {
	OverallStatus       string                  `json:"overall_status"`
	OverallMessage      string                  `json:"overall_message"`
	NutritionalAnalysis []NutritionAnalysisItem `json:"nutritional_analysis"`
	AIRecommendation    string                  `json:"ai_recommendation"`
	RecommendedFoods    []string                `json:"recommended_foods"`
	HealthInsight       HealthInsight           `json:"health_insight"`
	SuggestedActivities []string                `json:"suggested_activities"`
}

// NutritionAnalysisResult - response endpoint AI
type NutritionAnalysisResult struct {
	Source string                `json:"source"`
	Data   NutritionAnalysisData `json:"data"`
}

// GroqInput - input untuk Groq
type GroqInput struct {
	SubmissionID   string                `json:"submission_id"`
	SurveyID       string                `json:"survey_id"`
	RespondentName string                `json:"respondent_name"`
	MealsData      json.RawMessage       `json:"meals_data"`
	MissingFoods   json.RawMessage       `json:"missing_foods"`
	DailyTotal     submission.DailyTotal `json:"daily_total"`
}

// GroqResult - hasil parse dari Groq
type GroqResult struct {
	Data        NutritionAnalysisData `json:"data"`
	RawResponse json.RawMessage       `json:"raw_response"`
	ModelUsed   string                `json:"model_used"`
	TokenUsed   *int                  `json:"token_used,omitempty"`
	LatencyMs   *int                  `json:"latency_ms,omitempty"`
}
