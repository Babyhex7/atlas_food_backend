package ai

import (
	"atlas_food/internal/domain/submission"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository - operasi database untuk AI logs dan data pendukung analisis
type Repository interface {
	// FindOwnedSubmission - ambil submission milik user berdasarkan id server ATAU local_id.
	// Mengembalikan gorm.ErrRecordNotFound bila tidak ada atau bukan milik user.
	FindOwnedSubmission(ref, userID string) (*submission.SurveySubmission, error)
	// GetProfile - jenis kelamin dan usia user pada tanggal tertentu; kosong bila belum diisi
	GetProfile(userID string, at time.Time) Profile
	// FindBySubmissionID - hasil analisis tersimpan; (nil, nil) bila belum ada
	FindBySubmissionID(submissionID string) (*AIResultLog, error)
	// Save - simpan hasil; menimpa baris lama untuk submission yang sama
	Save(log *AIResultLog) error
}

// repository - implementasi Repository AI di atas GORM
type repository struct {
	db *gorm.DB
}

// NewRepository - buat repository AI
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// FindOwnedSubmission - kepemilikan diperiksa di query yang sama dengan pencarian,
// sehingga submission milik orang lain tidak pernah dimuat ke memori.
func (r *repository) FindOwnedSubmission(ref, userID string) (*submission.SurveySubmission, error) {
	var item submission.SurveySubmission
	err := r.db.
		Joins("JOIN survey_participants ON survey_participants.id = survey_submissions.participant_id").
		Where("survey_participants.user_id = ?", userID).
		Where("(survey_submissions.id = ? OR survey_submissions.local_id = ?)", ref, ref).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetProfile - baca hanya kolom yang dibutuhkan untuk memilih rujukan AKG
func (r *repository) GetProfile(userID string, at time.Time) Profile {
	var row struct {
		Gender    *string
		BirthDate *time.Time
	}
	if err := r.db.Table("users").Select("gender, birth_date").Where("id = ?", userID).Take(&row).Error; err != nil {
		return Profile{}
	}

	var profile Profile
	if row.Gender != nil {
		profile.Sex = *row.Gender
	}
	if row.BirthDate != nil {
		profile.Age = ageAt(*row.BirthDate, at)
	}
	return profile
}

// FindBySubmissionID - cari hasil analisis yang sudah pernah dibuat untuk sebuah submission
func (r *repository) FindBySubmissionID(submissionID string) (*AIResultLog, error) {
	var log AIResultLog
	err := r.db.Where("submission_id = ?", submissionID).First(&log).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// Save - upsert berdasarkan submission_id (kolom unik), supaya analisis ulang
// maupun dua permintaan yang nyaris bersamaan tidak gagal karena baris ganda.
func (r *repository) Save(log *AIResultLog) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "submission_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"input_payload", "raw_response", "result_data", "overall_status",
			"model_used", "prompt_version", "token_used", "prompt_tokens",
			"completion_tokens", "latency_ms", "updated_at",
		}),
	}).Create(log).Error
}
