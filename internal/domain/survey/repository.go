package survey

import "gorm.io/gorm"

// Repository - interface operasi database survey
type Repository interface {
	CreateSurvey(s *Survey) error
	GetSurveyByID(id string) (*Survey, error)
	GetSurveyBySlug(slug string) (*Survey, error)
	GetSurveyByAccessToken(accessToken string) (*Survey, error)
	ListSurveys() ([]Survey, error)
	UpdateSurvey(s *Survey) error
	DeleteSurvey(id string) error
	UpsertLocale(code, name string) error
}

type surveyRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &surveyRepository{db: db}
}

func (r *surveyRepository) CreateSurvey(s *Survey) error {
	return r.db.Create(s).Error
}

func (r *surveyRepository) GetSurveyByID(id string) (*Survey, error) {
	var survey Survey
	err := r.db.Where("id = ?", id).First(&survey).Error
	if err != nil {
		return nil, err
	}
	return &survey, nil
}

func (r *surveyRepository) GetSurveyBySlug(slug string) (*Survey, error) {
	var survey Survey
	err := r.db.Where("slug = ?", slug).First(&survey).Error
	if err != nil {
		return nil, err
	}
	return &survey, nil
}

func (r *surveyRepository) GetSurveyByAccessToken(accessToken string) (*Survey, error) {
	var survey Survey
	err := r.db.Where("access_token = ? AND status = ?", accessToken, "active").First(&survey).Error
	if err != nil {
		return nil, err
	}
	return &survey, nil
}

func (r *surveyRepository) ListSurveys() ([]Survey, error) {
	var surveys []Survey
	err := r.db.Order("created_at DESC").Find(&surveys).Error
	return surveys, err
}

func (r *surveyRepository) UpdateSurvey(s *Survey) error {
	return r.db.Save(s).Error
}

func (r *surveyRepository) DeleteSurvey(id string) error {
	return r.db.Where("id = ?", id).Delete(&Survey{}).Error
}

func (r *surveyRepository) UpsertLocale(code, name string) error {
	var locale Locale
	err := r.db.Where("code = ?", code).First(&locale).Error
	if err == nil {
		locale.Name = name
		return r.db.Save(&locale).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	return r.db.Create(&Locale{Code: code, Name: name}).Error
}
