package submission

import (
	fooddomain "atlas_food/internal/domain/food"
	"atlas_food/internal/domain/survey"
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
	SubmitSurvey(req SubmitSurveyRequest, userID string) (*SubmissionResponse, error)
	ListSubmissions(surveyID string, page, limit int) ([]ListSubmissionResponse, int64, error)
	GetMySubmissions(userID, userEmail string, page, limit int) ([]ListSubmissionResponse, int64, error)
	GetSubmissionDetail(id string) (*SubmissionDetailResponse, error)
	GetMySubmissionDetail(id, userID, userEmail string) (*SubmissionDetailResponse, error)
	ExportSubmissionsCSV(surveyID string) ([]byte, string, error)
}

type submissionService struct {
	repo       Repository
	surveyRepo survey.Repository
	foodRepo   fooddomain.Repository
}

// NewService - buat instance service submission
func NewService(repo Repository, surveyRepo survey.Repository, foodRepo fooddomain.Repository) Service {
	return &submissionService{repo: repo, surveyRepo: surveyRepo, foodRepo: foodRepo}
}

// SubmitSurvey - simpan hasil recall dari respondent
// userID diambil dari JWT (bukan dari body request) agar participant tidak bisa dipalsukan (anti-IDOR).
func (s *submissionService) SubmitSurvey(req SubmitSurveyRequest, userID string) (*SubmissionResponse, error) {
	if req.SurveyID == "" {
		return nil, errors.New("survey_id wajib diisi")
	}
	if len(req.MealsData) == 0 {
		return nil, errors.New("minimal 1 waktu makan harus diisi")
	}
	if userID == "" {
		return nil, utils.NewAppError(401, "UNAUTHORIZED", "Login diperlukan")
	}

	// Verifikasi survey benar-benar ada dan masih aktif (mencegah submit ke survey
	// yang belum pernah di-access, sudah closed, atau di luar rentang tanggal).
	surv, err := s.surveyRepo.GetSurveyByID(req.SurveyID)
	if err != nil {
		return nil, utils.NewAppError(404, "NOT_FOUND", "Survey tidak ditemukan")
	}
	if surv.Status != "active" {
		return nil, utils.NewAppError(403, "SURVEY_NOT_ACTIVE", "Survey tidak aktif")
	}
	now := time.Now()
	if surv.StartDate != nil && now.Before(*surv.StartDate) {
		return nil, utils.NewAppError(403, "SURVEY_NOT_STARTED", "Survey belum dimulai")
	}
	if surv.EndDate != nil && now.After(*surv.EndDate) {
		return nil, utils.NewAppError(403, "SURVEY_ENDED", "Survey sudah berakhir")
	}

	// Participant selalu diresolusi dari (surveyID, userID) di server, bukan dari
	// participant_id yang dikirim client, supaya tidak bisa submit atas nama
	// participant lain.
	participant, err := s.surveyRepo.GetParticipantBySurveyAndUser(req.SurveyID, userID)
	if err != nil {
		return nil, utils.NewAppError(403, "NOT_JOINED", "Anda belum bergabung ke survey ini")
	}
	req.ParticipantID = participant.ID
	if req.RespondentName == "" {
		req.RespondentName = participant.Alias
	}

	hasFood := false
	for _, meal := range req.MealsData {
		if len(meal.Foods) == 0 {
			continue
		}
		hasFood = true
		for _, food := range meal.Foods {
			if food.PortionGram <= 0 {
				return nil, errors.New("semua makanan harus memiliki porsi valid")
			}
		}
	}
	if !hasFood {
		return nil, errors.New("minimal 1 makanan harus diisi sebelum submit")
	}

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
	if err := s.calculateTotals(&req); err != nil {
		return nil, err
	}

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
		var meals []MealData
		json.Unmarshal([]byte(sub.MealsData), &meals)

		foodCount := 0
		for _, m := range meals {
			foodCount += len(m.Foods)
		}

		resp[i] = ListSubmissionResponse{
			ID:              sub.ID,
			RespondentName:  sub.RespondentName,
			RespondentEmail: sub.RespondentEmail,
			SubmittedAt:     sub.SubmittedAt.Format("2006-01-02 15:04:05"),
			MealCount:       len(meals),
			TotalFoods:      foodCount,
			TotalEnergy:     sub.TotalEnergy,
			TotalProtein:    sub.TotalProtein,
			TotalCarbs:      sub.TotalCarbs,
			TotalFat:        sub.TotalFat,
			MealsData:       json.RawMessage(sub.MealsData),
		}
	}

	return resp, total, nil
}

// GetMySubmissions - list submission milik user yang sedang login (respondent)
func (s *submissionService) GetMySubmissions(userID, userEmail string, page, limit int) ([]ListSubmissionResponse, int64, error) {
	submissions, total, err := s.repo.ListSubmissionsByUserID(userID, userEmail, page, limit)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]ListSubmissionResponse, len(submissions))
	for i, sub := range submissions {
		var meals []MealData
		json.Unmarshal([]byte(sub.MealsData), &meals)

		foodCount := 0
		for _, m := range meals {
			foodCount += len(m.Foods)
		}

		resp[i] = ListSubmissionResponse{
			ID:              sub.ID,
			RespondentName:  sub.RespondentName,
			RespondentEmail: sub.RespondentEmail,
			SubmittedAt:     sub.SubmittedAt.Format("2006-01-02 15:04:05"),
			MealCount:       len(meals),
			TotalFoods:      foodCount,
			TotalEnergy:     sub.TotalEnergy,
			TotalProtein:    sub.TotalProtein,
			TotalCarbs:      sub.TotalCarbs,
			TotalFat:        sub.TotalFat,
			MealsData:       json.RawMessage(sub.MealsData),
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
		SubmittedAt: sub.SubmittedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetMySubmissionDetail - detail submission milik user yang sedang login
func (s *submissionService) GetMySubmissionDetail(id, userID, userEmail string) (*SubmissionDetailResponse, error) {
	sub, err := s.repo.GetSubmissionByIDAndUser(id, userID, userEmail)
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
		SubmittedAt: sub.SubmittedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// calculateTotals - menghitung total nutrisi per meal dan per hari (server-side validation)
func (s *submissionService) calculateTotals(req *SubmitSurveyRequest) error {
	var dailyEnergy, dailyProtein, dailyCarbs, dailyFat float64

	for i := range req.MealsData {
		var mealEnergy, mealProtein, mealCarbs, mealFat float64
		for j := range req.MealsData[i].Foods {
			energy, protein, carbs, fat, err := s.calculateFoodNutrients(req.MealsData[i].Foods[j])
			if err != nil {
				return err
			}
			req.MealsData[i].Foods[j].Nutrients = NutrientValues{
				Energy:  energy,
				Protein: protein,
				Carbs:   carbs,
				Fat:     fat,
			}
			mealEnergy += energy
			mealProtein += protein
			mealCarbs += carbs
			mealFat += fat
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
	return nil
}

// calculateFoodNutrients - hitung energi/protein/karbo/lemak satu item makanan.
// Makanan "missing-" (tidak ada di database) memakai nilai gizi yang dikirim frontend apa adanya,
// selebihnya diambil dari nilai per 100g di database lalu diskalakan sesuai berat porsi
func (s *submissionService) calculateFoodNutrients(food FoodData) (float64, float64, float64, float64, error) {
	if food.FoodID == "" || strings.HasPrefix(food.FoodID, "missing-") {
		return food.Nutrients.Energy, food.Nutrients.Protein, food.Nutrients.Carbs, food.Nutrients.Fat, nil
	}

	if s.foodRepo == nil {
		return 0, 0, 0, 0, errors.New("food repository tidak tersedia")
	}

	if _, err := s.foodRepo.GetFoodByID(food.FoodID); err != nil {
		return 0, 0, 0, 0, errors.New("makanan tidak ditemukan di database")
	}

	nutrients, err := s.foodRepo.GetNutrientsByFoodID(food.FoodID)
	if err != nil {
		return 0, 0, 0, 0, errors.New("gagal mengambil nutrisi makanan")
	}

	var energyPer100, proteinPer100, carbsPer100, fatPer100 float64
	for _, nutrient := range nutrients {
		switch strings.ToLower(nutrient.NutrientType.Code) {
		case "energy", "calories", "calorie", "energy_kcal":
			energyPer100 = nutrient.ValuePer100g
		case "protein":
			proteinPer100 = nutrient.ValuePer100g
		case "carbs", "carbohydrate", "carbohydrates":
			carbsPer100 = nutrient.ValuePer100g
		case "fat", "total_fat":
			fatPer100 = nutrient.ValuePer100g
		}
	}

	factor := food.PortionGram / 100
	return energyPer100 * factor, proteinPer100 * factor, carbsPer100 * factor, fatPer100 * factor, nil
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
