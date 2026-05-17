package survey

import (
	"atlas_food/internal/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service - interface untuk business logic survey
type Service interface {
	CreateSurvey(req CreateSurveyRequest, createdBy string) (*SurveyResponse, error)
	GetSurveyByID(id string) (*SurveyResponse, error)
	ListSurveys(createdBy string, page, limit int) (*SurveyListResponse, error)
	UpdateSurvey(id string, req UpdateSurveyRequest) (*SurveyResponse, error)
	DeleteSurvey(id string) error
	CloneSurvey(id string, req CloneSurveyRequest, createdBy string) (*SurveyResponse, error)
	GenerateAccessToken(surveyID string) (*AccessTokenResponse, error)

	// Public/Respondent operations
	GetPublicSurveyByToken(token string) (*PublicSurveyResponse, error)
	AccessSurvey(req AccessSurveyRequest, userID *string) (*AccessSurveyResponse, error)

	// Locale operations
	GetAllLocales() ([]Locale, error)
}

// surveyService - implementasi Service
type surveyService struct {
	repo Repository
}

// NewService - factory function
func NewService(repo Repository) Service {
	return &surveyService{repo: repo}
}

// CreateSurvey - buat survey baru
func (s *surveyService) CreateSurvey(req CreateSurveyRequest, createdBy string) (*SurveyResponse, error) {
	// Cek slug unik
	exists, err := s.repo.CountSurveysBySlug(req.Slug)
	if err != nil {
		return nil, errors.New("gagal cek slug")
	}
	if exists > 0 {
		return nil, utils.NewAppError(409, "SLUG_EXISTS", "Slug sudah digunakan")
	}

	// Parse meals_config ke JSON string
	mealsJSON, err := json.Marshal(req.MealsConfig)
	if err != nil {
		return nil, errors.New("gagal parse meals_config")
	}

	// Parse prompts ke JSON string
	promptsJSON, _ := json.Marshal(req.Prompts)

	// Parse dates
	var startDate, endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		sd, _ := time.Parse("2006-01-02", *req.StartDate)
		startDate = &sd
	}
	if req.EndDate != nil && *req.EndDate != "" {
		ed, _ := time.Parse("2006-01-02", *req.EndDate)
		endDate = &ed
	}

	// Default locale
	localeID := req.LocaleID
	if localeID == 0 {
		localeID = 1 // default Indonesia
	}

	survey := &Survey{
		ID:          uuid.New().String(),
		Slug:        req.Slug,
		Name:        req.Name,
		Description: req.Description,
		MealsConfig: string(mealsJSON),
		Prompts:     string(promptsJSON),
		LocaleID:    localeID,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      "draft",
		AccessToken: strings.ReplaceAll(uuid.New().String(), "-", ""),
		CreatedBy:   createdBy,
	}

	if err := s.repo.CreateSurvey(survey); err != nil {
		return nil, err
	}

	return s.mapToResponse(survey), nil
}

// GetSurveyByID - ambil detail survey
func (s *surveyService) GetSurveyByID(id string) (*SurveyResponse, error) {
	survey, err := s.repo.GetSurveyByID(id)
	if err != nil {
		return nil, utils.NewAppError(404, "NOT_FOUND", "Survey tidak ditemukan")
	}
	return s.mapToResponse(survey), nil
}

// ListSurveys - list survey dengan pagination
func (s *surveyService) ListSurveys(createdBy string, page, limit int) (*SurveyListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	surveys, total, err := s.repo.ListSurveys(createdBy, page, limit)
	if err != nil {
		return nil, errors.New("gagal mengambil data survey")
	}

	result := make([]ListSurveysResponse, len(surveys))
	for i, sv := range surveys {
		var startDate, endDate *string
		if sv.StartDate != nil {
			sd := sv.StartDate.Format("2006-01-02")
			startDate = &sd
		}
		if sv.EndDate != nil {
			ed := sv.EndDate.Format("2006-01-02")
			endDate = &ed
		}

		// Hitung participant count
		count, _ := s.repo.CountParticipantsBySurvey(sv.ID)

		result[i] = ListSurveysResponse{
			ID:               sv.ID,
			Slug:             sv.Slug,
			Name:             sv.Name,
			Status:           sv.Status,
			StartDate:        startDate,
			EndDate:          endDate,
			ParticipantCount: int(count),
			CreatedAt:        sv.CreatedAt.Format("2006-01-02"),
		}
	}

	return &SurveyListResponse{
		Surveys: result,
		Total:   total,
		Page:    page,
		Limit:   limit,
	}, nil
}

// UpdateSurvey - update data survey
func (s *surveyService) UpdateSurvey(id string, req UpdateSurveyRequest) (*SurveyResponse, error) {
	survey, err := s.repo.GetSurveyByID(id)
	if err != nil {
		return nil, utils.NewAppError(404, "NOT_FOUND", "Survey tidak ditemukan")
	}

	// Update field yang dikirim
	if req.Name != "" {
		survey.Name = req.Name
	}
	if req.Description != "" {
		survey.Description = req.Description
	}
	if req.MealsConfig != nil {
		mealsJSON, _ := json.Marshal(req.MealsConfig)
		survey.MealsConfig = string(mealsJSON)
	}
	if req.Prompts != nil {
		promptsJSON, _ := json.Marshal(req.Prompts)
		survey.Prompts = string(promptsJSON)
	}
	if req.LocaleID > 0 {
		survey.LocaleID = req.LocaleID
	}
	if req.Status != "" {
		survey.Status = req.Status
	}
	if req.StartDate != nil && *req.StartDate != "" {
		sd, _ := time.Parse("2006-01-02", *req.StartDate)
		survey.StartDate = &sd
	}
	if req.EndDate != nil && *req.EndDate != "" {
		ed, _ := time.Parse("2006-01-02", *req.EndDate)
		survey.EndDate = &ed
	}

	if err := s.repo.UpdateSurvey(survey); err != nil {
		return nil, errors.New("gagal update survey")
	}

	return s.mapToResponse(survey), nil
}

// DeleteSurvey - hapus survey
func (s *surveyService) DeleteSurvey(id string) error {
	_, err := s.repo.GetSurveyByID(id)
	if err != nil {
		return utils.NewAppError(404, "NOT_FOUND", "Survey tidak ditemukan")
	}

	return s.repo.DeleteSurvey(id)
}

// CloneSurvey - duplikat survey
func (s *surveyService) CloneSurvey(id string, req CloneSurveyRequest, createdBy string) (*SurveyResponse, error) {
	// Cek slug baru
	exists, _ := s.repo.CountSurveysBySlug(req.NewSlug)
	if exists > 0 {
		return nil, utils.NewAppError(409, "SLUG_EXISTS", "Slug baru sudah digunakan")
	}

	// Ambil survey asli
	original, err := s.repo.GetSurveyByID(id)
	if err != nil {
		return nil, utils.NewAppError(404, "NOT_FOUND", "Survey tidak ditemukan")
	}

	// Buat survey baru dengan data dari original
	newSurvey := &Survey{
		ID:          uuid.New().String(),
		Slug:        req.NewSlug,
		Name:        req.NewName,
		Description: original.Description,
		MealsConfig: original.MealsConfig,
		Prompts:     original.Prompts,
		LocaleID:    original.LocaleID,
		Status:      "draft",
		AccessToken: strings.ReplaceAll(uuid.New().String(), "-", ""),
		CreatedBy:   createdBy,
	}

	if err := s.repo.CreateSurvey(newSurvey); err != nil {
		return nil, errors.New("gagal clone survey")
	}

	return s.mapToResponse(newSurvey), nil
}

// GenerateAccessToken - generate token baru untuk survey
func (s *surveyService) GenerateAccessToken(surveyID string) (*AccessTokenResponse, error) {
	survey, err := s.repo.GetSurveyByID(surveyID)
	if err != nil {
		return nil, utils.NewAppError(404, "NOT_FOUND", "Survey tidak ditemukan")
	}

	// Generate token baru
	survey.AccessToken = strings.ReplaceAll(uuid.New().String(), "-", "")
	if err := s.repo.UpdateSurvey(survey); err != nil {
		return nil, errors.New("gagal update token")
	}

	return &AccessTokenResponse{
		SurveyID:    survey.ID,
		AccessToken: survey.AccessToken,
		AccessURL:   fmt.Sprintf("/s/%s", survey.AccessToken),
	}, nil
}

// GetPublicSurveyByToken - untuk respondent akses survey
func (s *surveyService) GetPublicSurveyByToken(token string) (*PublicSurveyResponse, error) {
	survey, err := s.repo.GetSurveyByAccessToken(token)
	if err != nil {
		return nil, utils.NewAppError(404, "NOT_FOUND", "Survey tidak ditemukan")
	}

	// Parse meals_config
	var mealsConfig MealsConfig
	json.Unmarshal([]byte(survey.MealsConfig), &mealsConfig)

	// Parse prompts
	var prompts PromptsConfig
	json.Unmarshal([]byte(survey.Prompts), &prompts)

	var startDate, endDate *string
	if survey.StartDate != nil {
		sd := survey.StartDate.Format("2006-01-02")
		startDate = &sd
	}
	if survey.EndDate != nil {
		ed := survey.EndDate.Format("2006-01-02")
		endDate = &ed
	}

	return &PublicSurveyResponse{
		ID:          survey.ID,
		Name:        survey.Name,
		Description: survey.Description,
		MealsConfig: mealsConfig,
		Prompts:     prompts,
		Locale: LocaleInfo{
			ID:   survey.Locale.ID,
			Code: survey.Locale.Code,
			Name: survey.Locale.Name,
		},
		StartDate: startDate,
		EndDate:   endDate,
		Status:    survey.Status,
	}, nil
}

// AccessSurvey - respondent access survey
func (s *surveyService) AccessSurvey(req AccessSurveyRequest, userID *string) (*AccessSurveyResponse, error) {
	survey, err := s.repo.GetSurveyByAccessToken(req.Token)
	if err != nil {
		return nil, utils.NewAppError(404, "NOT_FOUND", "Survey tidak ditemukan")
	}

	if userID == nil || *userID == "" {
		return nil, utils.NewAppError(401, "UNAUTHORIZED", "Login sebagai respondent diperlukan")
	}

	// Cek status survey
	if survey.Status != "active" {
		return nil, utils.NewAppError(403, "SURVEY_NOT_ACTIVE", "Survey tidak aktif")
	}

	// Cek tanggal
	now := time.Now()
	if survey.StartDate != nil && now.Before(*survey.StartDate) {
		return nil, utils.NewAppError(403, "SURVEY_NOT_STARTED", "Survey belum dimulai")
	}
	if survey.EndDate != nil && now.After(*survey.EndDate) {
		return nil, utils.NewAppError(403, "SURVEY_ENDED", "Survey sudah berakhir")
	}

	// Alias default
	alias := req.Alias
	if alias == "" && req.RespondentName != "" {
		alias = req.RespondentName
	}

	// Cek apakah user sudah pernah join
	var participant *SurveyParticipant
	if userID != nil {
		participant, _ = s.repo.GetParticipantBySurveyAndUser(survey.ID, *userID)
	}

	if participant == nil {
		// Buat participant baru
		participant = &SurveyParticipant{
			ID:       uuid.New().String(),
			SurveyID: survey.ID,
			UserID:   *userID,
			Alias:    alias,
		}

		if err := s.repo.CreateParticipant(participant); err != nil {
			return nil, errors.New("gagal access survey")
		}
	}

	// Map survey to public response
	publicSurvey, _ := s.GetPublicSurveyByToken(req.Token)

	return &AccessSurveyResponse{
		Survey: *publicSurvey,
		Participant: ParticipantResponse{
			ID:    participant.ID,
			Alias: participant.Alias,
		},
		AccessToken: "dummy-session-token", // In real case, you might generate a survey-specific JWT
	}, nil
}

// GetAllLocales - ambil semua locales
func (s *surveyService) GetAllLocales() ([]Locale, error) {
	return s.repo.GetAllLocales()
}

// Helper: map Survey ke SurveyResponse
func (s *surveyService) mapToResponse(survey *Survey) *SurveyResponse {
	// Parse meals_config
	var mealsConfig MealsConfig
	json.Unmarshal([]byte(survey.MealsConfig), &mealsConfig)

	// Parse prompts
	var prompts PromptsConfig
	json.Unmarshal([]byte(survey.Prompts), &prompts)

	var startDate, endDate *string
	if survey.StartDate != nil {
		sd := survey.StartDate.Format("2006-01-02")
		startDate = &sd
	}
	if survey.EndDate != nil {
		ed := survey.EndDate.Format("2006-01-02")
		endDate = &ed
	}

	return &SurveyResponse{
		ID:          survey.ID,
		Slug:        survey.Slug,
		Name:        survey.Name,
		Description: survey.Description,
		MealsConfig: mealsConfig,
		Prompts:     prompts,
		Locale: LocaleInfo{
			ID:   survey.Locale.ID,
			Code: survey.Locale.Code,
			Name: survey.Locale.Name,
		},
		StartDate: startDate,
		EndDate:   endDate,
		Status:    survey.Status,
		CreatedBy: survey.CreatedBy,
		AccessURL: fmt.Sprintf("/s/%s", survey.AccessToken),
		CreatedAt: survey.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: survey.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
