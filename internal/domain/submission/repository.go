package submission

import (
	"gorm.io/gorm"
)

// Repository - interface untuk operasi database submission
type Repository interface {
	CreateSubmission(submission *SurveySubmission) error
	GetSubmissionByID(id string) (*SurveySubmission, error)
	ListSubmissionsBySurvey(surveyID string, page, limit int) ([]SurveySubmission, int64, error)
	DeleteSubmission(id string) error
}

type submissionRepository struct {
	db *gorm.DB
}

// NewRepository - buat instance repository submission
func NewRepository(db *gorm.DB) Repository {
	return &submissionRepository{db: db}
}

// CreateSubmission - simpan hasil survey ke database
func (r *submissionRepository) CreateSubmission(submission *SurveySubmission) error {
	return r.db.Create(submission).Error
}

// GetSubmissionByID - ambil detail submission berdasarkan ID
func (r *submissionRepository) GetSubmissionByID(id string) (*SurveySubmission, error) {
	var submission SurveySubmission
	err := r.db.Where("id = ?", id).First(&submission).Error
	return &submission, err
}

// ListSubmissionsBySurvey - ambil semua submission untuk satu survey (admin)
func (r *submissionRepository) ListSubmissionsBySurvey(surveyID string, page, limit int) ([]SurveySubmission, int64, error) {
	var submissions []SurveySubmission
	var total int64

	query := r.db.Model(&SurveySubmission{}).Where("survey_id = ?", surveyID)
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("submitted_at DESC").Find(&submissions).Error

	return submissions, total, err
}

// DeleteSubmission - hapus submission
func (r *submissionRepository) DeleteSubmission(id string) error {
	return r.db.Where("id = ?", id).Delete(&SurveySubmission{}).Error
}
