package submission

import (
	"atlas_food/internal/pkg/utils"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service - interface untuk business logic submission
type Service interface {
	SubmitSurvey(req SubmitSurveyRequest) (*SubmissionResponse, error)
	ListSubmissions(surveyID string, page, limit int) ([]ListSubmissionResponse, int64, error)
	GetSubmissionDetail(id string) (*SubmissionDetailResponse, error)
	ExportSubmissionsCSV(surveyID string) ([]byte, string, error)
}

type submissionService struct {
	repo Repository
}

// NewService - buat instance service submission
func NewService(repo Repository) Service {
	return &submissionService{repo: repo}
}

// SubmitSurvey - simpan hasil recall dari respondent
func (s *submissionService) SubmitSurvey(req SubmitSurveyRequest) (*SubmissionResponse, error) {
	// Marshal data ke JSON string untuk disimpan di DB
	mealsJSON, err := json.Marshal(req.MealsData)
	if err != nil {
		return nil, errors.New("gagal memproses data makanan")
	}

	missingJSON, _ := json.Marshal(req.MissingFoods)

	// Buat model submission
	submission := &SurveySubmission{
		ID:              uuid.New().String(),
		SurveyID:        req.SurveyID,
		RespondentName:  req.RespondentName,
		RespondentEmail: req.RespondentEmail,
		MealsData:       string(mealsJSON),
		MissingFoods:    string(missingJSON),
		SubmittedAt:     time.Now(),
	}

	if req.ParticipantID != "" {
		pID := req.ParticipantID
		submission.ParticipantID = &pID
	}

	// Calculate/Verify totals
	s.calculateTotals(&req)
	
	// Update meals JSON with calculated totals
	updatedMealsJSON, _ := json.Marshal(req.MealsData)
	submission.MealsData = string(updatedMealsJSON)
	
	// Set aggregate totals
	submission.TotalEnergy = req.DailyTotal.Energy
	submission.TotalProtein = req.DailyTotal.Protein
	submission.TotalCarbs = req.DailyTotal.Carbs
	submission.TotalFat = req.DailyTotal.Fat

	// Simpan ke database
	if err := s.repo.CreateSubmission(submission); err != nil {
		return nil, errors.New("gagal menyimpan hasil survey")
	}

	return &SubmissionResponse{
		SubmissionID: submission.ID,
		Message:      "Survey berhasil dikirim, terima kasih!",
	}, nil
}

// ListSubmissions - list submission untuk admin
func (s *submissionService) ListSubmissions(surveyID string, page, limit int) ([]ListSubmissionResponse, int64, error) {
	submissions, total, err := s.repo.ListSubmissionsBySurvey(surveyID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]ListSubmissionResponse, len(submissions))
	for i, sub := range submissions {
		// Hitung jumlah meal dan food dari JSON
		var meals []MealData
		json.Unmarshal([]byte(sub.MealsData), &meals)

		foodCount := 0
		for _, m := range meals {
			foodCount += len(m.Foods)
		}

		resp[i] = ListSubmissionResponse{
			ID:             sub.ID,
			RespondentName: sub.RespondentName,
			SubmittedAt:    sub.SubmittedAt.Format("2006-01-02 15:04:05"),
			MealCount:      len(meals),
			TotalFoods:     foodCount,
			TotalEnergy:    sub.TotalEnergy,
		}
	}

	return resp, total, nil
}

// GetSubmissionDetail - detail submission untuk admin
func (s *submissionService) GetSubmissionDetail(id string) (*SubmissionDetailResponse, error) {
	sub, err := s.repo.GetSubmissionByID(id)
	if err != nil {
		return nil, utils.NewAppError(404, "NOT_FOUND", "Submission tidak ditemukan")
	}

	return &SubmissionDetailResponse{
		ID:              sub.ID,
		SurveyID:        sub.SurveyID,
		RespondentName:  sub.RespondentName,
		RespondentEmail: sub.RespondentEmail,
		MealsData:       json.RawMessage(sub.MealsData),
		MissingFoods:    json.RawMessage(sub.MissingFoods),
		DailyTotal: DailyTotal{
			Energy:  sub.TotalEnergy,
			Protein: sub.TotalProtein,
			Carbs:   sub.TotalCarbs,
			Fat:     sub.TotalFat,
		},
		SubmittedAt:     sub.SubmittedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// calculateTotals - menghitung total nutrisi per meal dan per hari (server-side validation)
func (s *submissionService) calculateTotals(req *SubmitSurveyRequest) {
	var dailyEnergy, dailyProtein, dailyCarbs, dailyFat float64

	for i := range req.MealsData {
		var mealEnergy, mealProtein, mealCarbs, mealFat float64
		for _, f := range req.MealsData[i].Foods {
			mealEnergy += f.Nutrients.Energy
			mealProtein += f.Nutrients.Protein
			mealCarbs += f.Nutrients.Carbs
			mealFat += f.Nutrients.Fat
		}
		
		// Update meal totals
		req.MealsData[i].MealTotal = DailyTotal{
			Energy:  mealEnergy,
			Protein: mealProtein,
			Carbs:   mealCarbs,
			Fat:     mealFat,
		}

		dailyEnergy += mealEnergy
		dailyProtein += mealProtein
		dailyCarbs += mealCarbs
		dailyFat += mealFat
	}

	// Update daily totals
	req.DailyTotal = DailyTotal{
		Energy:  dailyEnergy,
		Protein: dailyProtein,
		Carbs:   dailyCarbs,
		Fat:     dailyFat,
	}
}

// ExportSubmissionsCSV - generate CSV untuk export data survey
func (s *submissionService) ExportSubmissionsCSV(surveyID string) ([]byte, string, error) {
	submissions, _, err := s.repo.ListSubmissionsBySurvey(surveyID, 1, 1000) // limit 1000 for export
	if err != nil {
		return nil, "", err
	}

	// Buat buffer untuk CSV
	// Note: In real production, use streaming for large data
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	// Header
	writer.Write([]string{"SubmissionID", "Respondent", "Meal", "Food", "Portion(g)", "Energy", "Protein", "Carbs", "Fat", "SubmittedAt"})

	for _, sub := range submissions {
		var meals []MealData
		json.Unmarshal([]byte(sub.MealsData), &meals)

		for _, m := range meals {
			for _, f := range m.Foods {
				writer.Write([]string{
					sub.ID,
					sub.RespondentName,
					m.Name,
					f.FoodName,
					fmt.Sprintf("%.2f", f.PortionGram),
					fmt.Sprintf("%.2f", f.Nutrients.Energy),
					fmt.Sprintf("%.2f", f.Nutrients.Protein),
					fmt.Sprintf("%.2f", f.Nutrients.Carbs),
					fmt.Sprintf("%.2f", f.Nutrients.Fat),
					sub.SubmittedAt.Format("2006-01-02 15:04:05"),
				})
			}
		}
	}

	writer.Flush()
	filename := fmt.Sprintf("export-survey-%s-%s.csv", surveyID, time.Now().Format("20060102"))

	return []byte(buf.String()), filename, nil
}
