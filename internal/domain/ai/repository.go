package ai

import (
	"atlas_food/internal/domain/submission"
	"atlas_food/internal/domain/survey"

	"gorm.io/gorm"
)

// Repository - operasi database untuk AI logs dan ownership check
type Repository interface {
	GetSubmissionByID(id string) (*submission.SurveySubmission, error)
	GetParticipantByID(id string) (*survey.SurveyParticipant, error)
	FindBySubmissionID(submissionID string) (*AIResultLog, error)
	Save(log *AIResultLog) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository - buat repository AI
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetSubmissionByID(id string) (*submission.SurveySubmission, error) {
	var item submission.SurveySubmission
	if err := r.db.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) GetParticipantByID(id string) (*survey.SurveyParticipant, error) {
	var item survey.SurveyParticipant
	if err := r.db.Preload("Survey").Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) FindBySubmissionID(submissionID string) (*AIResultLog, error) {
	var log AIResultLog
	if err := r.db.Where("submission_id = ?", submissionID).First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *repository) Save(log *AIResultLog) error {
	return r.db.Create(log).Error
}
