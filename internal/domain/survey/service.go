package survey

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service - business logic survey
type Service interface {
	List() ([]SurveyResponse, error)
	Create(req CreateSurveyRequest, createdBy string) (*SurveyResponse, error)
	GetByID(id string) (*SurveyResponse, error)
	GetByAccessToken(accessToken string) (*SurveyResponse, error)
	Update(id string, req UpdateSurveyRequest) (*SurveyResponse, error)
	Delete(id string) error
	Clone(id, createdBy string) (*SurveyResponse, error)
	SeedLocales() error
}

type surveyService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &surveyService{repo: repo}
}

func (s *surveyService) SeedLocales() error {
	if err := s.repo.UpsertLocale("id", "Indonesia"); err != nil {
		return err
	}
	if err := s.repo.UpsertLocale("en", "English"); err != nil {
		return err
	}
	return nil
}

func (s *surveyService) List() ([]SurveyResponse, error) {
	surveys, err := s.repo.ListSurveys()
	if err != nil {
		return nil, err
	}

	result := make([]SurveyResponse, 0, len(surveys))
	for i := range surveys {
		result = append(result, mapSurveyToResponse(&surveys[i]))
	}
	return result, nil
}

func (s *surveyService) Create(req CreateSurveyRequest, createdBy string) (*SurveyResponse, error) {
	if len(req.MealsConfig) == 0 {
		return nil, errors.New("meals_config wajib diisi")
	}

	startDate, endDate, err := parseDateRange(req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	mealsJSON, err := json.Marshal(req.MealsConfig)
	if err != nil {
		return nil, errors.New("gagal memproses meals_config")
	}

	promptsJSON, err := json.Marshal(req.Prompts)
	if err != nil {
		return nil, errors.New("gagal memproses prompts")
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = slugify(req.Name)
	}
	slug, err = s.ensureUniqueSlug(slug)
	if err != nil {
		return nil, errors.New("gagal membuat slug unik")
	}

	status := req.Status
	if status == "" {
		status = "draft"
	}

	localeID := req.LocaleID
	if localeID <= 0 {
		localeID = 1
	}

	survey := &Survey{
		ID:          uuid.New().String(),
		Slug:        slug,
		Name:        req.Name,
		Description: req.Description,
		MealsConfig: mealsJSON,
		Prompts:     promptsJSON,
		LocaleID:    localeID,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      status,
		AccessToken: strings.ReplaceAll(uuid.New().String(), "-", ""),
		CreatedBy:   createdBy,
	}

	if err := s.repo.CreateSurvey(survey); err != nil {
		return nil, err
	}

	resp := mapSurveyToResponse(survey)
	return &resp, nil
}

func (s *surveyService) GetByID(id string) (*SurveyResponse, error) {
	survey, err := s.repo.GetSurveyByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("survey tidak ditemukan")
		}
		return nil, err
	}
	resp := mapSurveyToResponse(survey)
	return &resp, nil
}

func (s *surveyService) GetByAccessToken(accessToken string) (*SurveyResponse, error) {
	survey, err := s.repo.GetSurveyByAccessToken(accessToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("survey tidak ditemukan atau belum aktif")
		}
		return nil, err
	}
	resp := mapSurveyToResponse(survey)
	return &resp, nil
}

func (s *surveyService) Update(id string, req UpdateSurveyRequest) (*SurveyResponse, error) {
	survey, err := s.repo.GetSurveyByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("survey tidak ditemukan")
		}
		return nil, err
	}

	if req.Name != nil {
		survey.Name = *req.Name
	}
	if req.Description != nil {
		survey.Description = *req.Description
	}
	if len(req.MealsConfig) > 0 {
		mealsJSON, err := json.Marshal(req.MealsConfig)
		if err != nil {
			return nil, errors.New("gagal memproses meals_config")
		}
		survey.MealsConfig = mealsJSON
	}
	if req.Prompts != nil {
		promptsJSON, err := json.Marshal(req.Prompts)
		if err != nil {
			return nil, errors.New("gagal memproses prompts")
		}
		survey.Prompts = promptsJSON
	}
	if req.LocaleID != nil {
		survey.LocaleID = *req.LocaleID
	}
	if req.Status != nil {
		survey.Status = *req.Status
	}

	startDateStr := ""
	if req.StartDate != nil {
		startDateStr = *req.StartDate
	}
	endDateStr := ""
	if req.EndDate != nil {
		endDateStr = *req.EndDate
	}
	if req.StartDate != nil || req.EndDate != nil {
		currentStart := ""
		currentEnd := ""
		if survey.StartDate != nil {
			currentStart = survey.StartDate.Format("2006-01-02")
		}
		if survey.EndDate != nil {
			currentEnd = survey.EndDate.Format("2006-01-02")
		}
		if startDateStr == "" {
			startDateStr = currentStart
		}
		if endDateStr == "" {
			endDateStr = currentEnd
		}
		startDate, endDate, err := parseDateRange(startDateStr, endDateStr)
		if err != nil {
			return nil, err
		}
		survey.StartDate = startDate
		survey.EndDate = endDate
	}

	if err := s.repo.UpdateSurvey(survey); err != nil {
		return nil, err
	}

	resp := mapSurveyToResponse(survey)
	return &resp, nil
}

func (s *surveyService) Delete(id string) error {
	_, err := s.repo.GetSurveyByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("survey tidak ditemukan")
		}
		return err
	}
	return s.repo.DeleteSurvey(id)
}

func (s *surveyService) Clone(id, createdBy string) (*SurveyResponse, error) {
	src, err := s.repo.GetSurveyByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("survey tidak ditemukan")
		}
		return nil, err
	}

	baseSlug := fmt.Sprintf("%s-copy", src.Slug)
	cloneSlug, err := s.ensureUniqueSlug(baseSlug)
	if err != nil {
		return nil, errors.New("gagal membuat slug clone")
	}

	clone := &Survey{
		ID:          uuid.New().String(),
		Slug:        cloneSlug,
		Name:        src.Name + " (Copy)",
		Description: src.Description,
		MealsConfig: src.MealsConfig,
		Prompts:     src.Prompts,
		LocaleID:    src.LocaleID,
		StartDate:   src.StartDate,
		EndDate:     src.EndDate,
		Status:      "draft",
		AccessToken: strings.ReplaceAll(uuid.New().String(), "-", ""),
		CreatedBy:   createdBy,
	}

	if err := s.repo.CreateSurvey(clone); err != nil {
		return nil, err
	}

	resp := mapSurveyToResponse(clone)
	return &resp, nil
}

func parseDateRange(startDateStr, endDateStr string) (*time.Time, *time.Time, error) {
	var startDate *time.Time
	var endDate *time.Time

	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return nil, nil, errors.New("format start_date harus YYYY-MM-DD")
		}
		startDate = &t
	}

	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, nil, errors.New("format end_date harus YYYY-MM-DD")
		}
		endDate = &t
	}

	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		return nil, nil, errors.New("end_date tidak boleh lebih kecil dari start_date")
	}

	return startDate, endDate, nil
}

func (s *surveyService) ensureUniqueSlug(base string) (string, error) {
	candidate := slugify(base)
	if candidate == "" {
		candidate = "survey"
	}

	for i := 0; i < 50; i++ {
		check := candidate
		if i > 0 {
			check = fmt.Sprintf("%s-%d", candidate, i+1)
		}
		_, err := s.repo.GetSurveyBySlug(check)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return check, nil
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}
	}

	return "", errors.New("slug tidak tersedia")
}

func slugify(input string) string {
	s := strings.ToLower(strings.TrimSpace(input))
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")
	re := regexp.MustCompile(`[^a-z0-9-]+`)
	s = re.ReplaceAllString(s, "")
	s = strings.Trim(s, "-")
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return s
}

func mapSurveyToResponse(s *Survey) SurveyResponse {
	var startDate *string
	var endDate *string
	if s.StartDate != nil {
		v := s.StartDate.Format("2006-01-02")
		startDate = &v
	}
	if s.EndDate != nil {
		v := s.EndDate.Format("2006-01-02")
		endDate = &v
	}

	prompts := json.RawMessage(s.Prompts)
	if len(prompts) == 0 || string(prompts) == "null" {
		prompts = nil
	}

	return SurveyResponse{
		ID:          s.ID,
		Slug:        s.Slug,
		Name:        s.Name,
		Description: s.Description,
		MealsConfig: json.RawMessage(s.MealsConfig),
		Prompts:     prompts,
		LocaleID:    s.LocaleID,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      s.Status,
		AccessToken: s.AccessToken,
		CreatedBy:   s.CreatedBy,
		CreatedAt:   s.CreatedAt.Format(time.RFC3339),
	}
}
