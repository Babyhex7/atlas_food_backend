package survey

import (
	"gorm.io/gorm"
)

// Repository - interface untuk operasi database survey
type Repository interface {
	// Survey operations
	CreateSurvey(survey *Survey) error
	GetSurveyByID(id string) (*Survey, error)
	GetSurveyBySlug(slug string) (*Survey, error)
	GetSurveyByAccessToken(token string) (*Survey, error)
	ListSurveys(createdBy string, page, limit int) ([]Survey, int64, error)
	UpdateSurvey(survey *Survey) error
	DeleteSurvey(id string) error
	CountSurveysBySlug(slug string) (int64, error)

	// Participant operations
	CreateParticipant(participant *SurveyParticipant) error
	GetParticipantByID(id string) (*SurveyParticipant, error)
	GetParticipantBySurveyAndUser(surveyID, userID string) (*SurveyParticipant, error)
	CountParticipantsBySurvey(surveyID string) (int64, error)

	// Locale operations
	GetAllLocales() ([]Locale, error)
	GetLocaleByID(id int) (*Locale, error)
}

// surveyRepository - implementasi Repository
type surveyRepository struct {
	db *gorm.DB
}

// NewRepository - factory function
func NewRepository(db *gorm.DB) Repository {
	return &surveyRepository{db: db}
}

// CreateSurvey - insert survey baru
func (r *surveyRepository) CreateSurvey(survey *Survey) error {
	return r.db.Create(survey).Error
}

// GetSurveyByID - cari survey berdasarkan ID
func (r *surveyRepository) GetSurveyByID(id string) (*Survey, error) {
	var survey Survey
	err := r.db.Preload("Locale").Where("id = ?", id).First(&survey).Error
	if err != nil {
		return nil, err
	}
	return &survey, nil
}

// GetSurveyBySlug - cari survey berdasarkan slug
func (r *surveyRepository) GetSurveyBySlug(slug string) (*Survey, error) {
	var survey Survey
	err := r.db.Preload("Locale").Where("slug = ?", slug).First(&survey).Error
	if err != nil {
		return nil, err
	}
	return &survey, nil
}

// GetSurveyByAccessToken - cari survey berdasarkan access token
func (r *surveyRepository) GetSurveyByAccessToken(token string) (*Survey, error) {
	var survey Survey
	err := r.db.Preload("Locale").Where("access_token = ?", token).First(&survey).Error
	if err != nil {
		return nil, err
	}
	return &survey, nil
}

// ListSurveys - list surveys dengan pagination
func (r *surveyRepository) ListSurveys(createdBy string, page, limit int) ([]Survey, int64, error) {
	var surveys []Survey
	var total int64

	query := r.db.Model(&Survey{})
	if createdBy != "" {
		query = query.Where("created_by = ?", createdBy)
	}

	// Hitung total
	query.Count(&total)

	// Ambil data dengan pagination
	offset := (page - 1) * limit
	err := query.Preload("Locale").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&surveys).Error

	return surveys, total, err
}

// UpdateSurvey - update data survey
func (r *surveyRepository) UpdateSurvey(survey *Survey) error {
	return r.db.Save(survey).Error
}

// DeleteSurvey - hapus survey
func (r *surveyRepository) DeleteSurvey(id string) error {
	return r.db.Where("id = ?", id).Delete(&Survey{}).Error
}

// CountSurveysBySlug - cek apakah slug sudah dipakai
func (r *surveyRepository) CountSurveysBySlug(slug string) (int64, error) {
	var count int64
	err := r.db.Model(&Survey{}).Where("slug = ?", slug).Count(&count).Error
	return count, err
}

// CreateParticipant - tambah participant baru
func (r *surveyRepository) CreateParticipant(participant *SurveyParticipant) error {
	return r.db.Create(participant).Error
}

// GetParticipantByID - cari participant berdasarkan ID
func (r *surveyRepository) GetParticipantByID(id string) (*SurveyParticipant, error) {
	var participant SurveyParticipant
	err := r.db.Preload("Survey").Where("id = ?", id).First(&participant).Error
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

// GetParticipantBySurveyAndUser - cari participant berdasarkan survey + user
func (r *surveyRepository) GetParticipantBySurveyAndUser(surveyID, userID string) (*SurveyParticipant, error) {
	var participant SurveyParticipant
	err := r.db.Where("survey_id = ? AND user_id = ?", surveyID, userID).First(&participant).Error
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

// CountParticipantsBySurvey - hitung jumlah participant di survey
func (r *surveyRepository) CountParticipantsBySurvey(surveyID string) (int64, error) {
	var count int64
	err := r.db.Model(&SurveyParticipant{}).Where("survey_id = ?", surveyID).Count(&count).Error
	return count, err
}

// GetAllLocales - ambil semua locales
func (r *surveyRepository) GetAllLocales() ([]Locale, error) {
	var locales []Locale
	err := r.db.Find(&locales).Error
	return locales, err
}

// GetLocaleByID - ambil locale berdasarkan ID
func (r *surveyRepository) GetLocaleByID(id int) (*Locale, error) {
	var locale Locale
	err := r.db.Where("id = ?", id).First(&locale).Error
	if err != nil {
		return nil, err
	}
	return &locale, nil
}
