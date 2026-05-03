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

	"atlas_food/internal/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
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
	JoinSurvey(token string, userID *string, req JoinSurveyRequest) (*JoinSurveyResponse, error)

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
	if req.StartDate != nil {
		sd, _ := time.Parse("2006-01-02", *req.StartDate)
		startDate = &sd
	}
	if req.EndDate != nil {
		ed, _ := time.Parse("2006-01-02", *req.EndDate)
		endDate = &ed
	}

	// Generate access token
	accessToken := uuid.New().String()

	// Default locale
	localeID := req.LocaleID
	if localeID == 0 {
		localeID = 1 // default Indonesia
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
		Slug:        req.Slug,
		Name:        req.Name,
		Description: req.Description,
		MealsConfig: string(mealsJSON),
		Prompts:     string(promptsJSON),
		LocaleID:    localeID,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      "draft",
		AccessToken: accessToken,
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
		return nil, errors.New("gagal membuat survey")
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
	promptsJSON, _ := json.Marshal(req.Prompts)
	survey.Prompts = string(promptsJSON)
	if req.LocaleID > 0 {
		survey.LocaleID = req.LocaleID
	}
	if req.Status != "" {
		survey.Status = req.Status
	}
	if req.StartDate != nil {
		sd, _ := time.Parse("2006-01-02", *req.StartDate)
		survey.StartDate = &sd
	}
	if req.EndDate != nil {
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
		AccessToken: uuid.New().String(),
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
	survey.AccessToken = uuid.New().String()
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

// JoinSurvey - respondent join survey
func (s *surveyService) JoinSurvey(token string, userID *string, req JoinSurveyRequest) (*JoinSurveyResponse, error) {
	survey, err := s.repo.GetSurveyByAccessToken(token)
	if err != nil {
		return nil, utils.NewAppError(404, "NOT_FOUND", "Survey tidak ditemukan")
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

	// Cek apakah user sudah pernah join
	if userID != nil {
		existing, _ := s.repo.GetParticipantBySurveyAndUser(survey.ID, *userID)
		if existing != nil {
			return nil, utils.NewAppError(409, "ALREADY_JOINED", "Anda sudah bergabung di survey ini")
		}
	}

	// Buat participant
	participant := &SurveyParticipant{
		ID:          uuid.New().String(),
		SurveyID:    survey.ID,
		UserID:      userID,
		Alias:       req.Alias,
		IsAnonymous: userID == nil,
	}

	if err := s.repo.CreateParticipant(participant); err != nil {
		return nil, errors.New("gagal join survey")
	}

	// Parse meals_config & prompts untuk response
	var mealsConfig MealsConfig
	json.Unmarshal([]byte(survey.MealsConfig), &mealsConfig)
	var prompts PromptsConfig
	json.Unmarshal([]byte(survey.Prompts), &prompts)

	return &JoinSurveyResponse{
		ParticipantID: participant.ID,
		SurveyID:      survey.ID,
		SurveyName:    survey.Name,
		MealsConfig:   mealsConfig,
		Prompts:       prompts,
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
