package ai

import (
	"atlas_food/internal/domain/submission"
	"atlas_food/internal/pkg/groq"
	"atlas_food/internal/pkg/utils"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

// Service - business logic AI analysis
type Service interface {
	AnalyzeNutrition(userID string, req NutritionAnalysisRequest) (*NutritionAnalysisResult, error)
}

type service struct {
	repo Repository
	groq *groq.Client
}

// NewService - buat service AI
func NewService(repo Repository, groqClient *groq.Client) Service {
	return &service{repo: repo, groq: groqClient}
}

// AnalyzeNutrition - analisis gizi sebuah submission memakai Groq.
// Memastikan submission memang milik user yang meminta, lalu memakai hasil lama bila sudah pernah dianalisis
func (s *service) AnalyzeNutrition(userID string, req NutritionAnalysisRequest) (*NutritionAnalysisResult, error) {
	if req.SubmissionID == "" {
		return nil, utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "submission_id wajib diisi")
	}

	sub, err := s.repo.GetSubmissionByID(req.SubmissionID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Submission not found or access denied")
	}

	if sub.ParticipantID == nil || *sub.ParticipantID == "" {
		return nil, utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Submission not found or access denied")
	}

	participant, err := s.repo.GetParticipantByID(*sub.ParticipantID)
	if err != nil || participant.UserID != userID {
		return nil, utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Submission not found or access denied")
	}

	if cached, err := s.repo.FindBySubmissionID(req.SubmissionID); err == nil {
		var cachedData NutritionAnalysisData
		if err := json.Unmarshal([]byte(cached.RawResponse), &cachedData); err == nil {
			return &NutritionAnalysisResult{Source: "cache", Data: cachedData}, nil
		}
	}

	inputPayload := GroqInput{
		SubmissionID:   sub.ID,
		SurveyID:       sub.SurveyID,
		RespondentName: sub.RespondentName,
		MealsData:      json.RawMessage(sub.MealsData),
		MissingFoods:   json.RawMessage(sub.MissingFoods),
		DailyTotal: submission.DailyTotal{
			Energy:  sub.TotalEnergy,
			Protein: sub.TotalProtein,
			Carbs:   sub.TotalCarbs,
			Fat:     sub.TotalFat,
		},
	}

	payloadBytes, _ := json.MarshalIndent(inputPayload, "", "  ")
	systemPrompt := `Analyze the submission and return only valid JSON with keys overall_status, overall_message, nutritional_analysis, ai_recommendation, recommended_foods, health_insight, suggested_activities.`
	result, err := s.groq.Analyze(systemPrompt, string(payloadBytes))
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return nil, appErr
		}
		return nil, utils.NewAppError(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "AI service temporarily unavailable, please try again")
	}

	inputBytes, _ := json.Marshal(inputPayload)
	var analysisData NutritionAnalysisData
	if err := json.Unmarshal([]byte(result.Content), &analysisData); err != nil {
		return nil, utils.NewAppError(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "AI response tidak valid")
	}

	log := &AIResultLog{
		ID:            uuid.New().String(),
		SubmissionID:  req.SubmissionID,
		InputPayload:  string(inputBytes),
		RawResponse:   result.Content,
		OverallStatus: analysisData.OverallStatus,
		ModelUsed:     result.ModelUsed,
		TokenUsed:     result.TokenUsed,
		LatencyMs:     result.LatencyMs,
	}

	if err := s.repo.Save(log); err != nil {
		return nil, errors.New("gagal menyimpan hasil AI")
	}

	return &NutritionAnalysisResult{Source: "groq", Data: analysisData}, nil
}
