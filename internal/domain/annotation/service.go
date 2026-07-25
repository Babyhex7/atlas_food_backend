package annotation

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"atlas_food/internal/pkg/utils"

	"github.com/google/uuid"
)

// Batas pagination default untuk list
const (
	defaultPage  = 1
	defaultLimit = 20
	maxLimit     = 100
)

// Service - business rules anotasi: validasi polygon, aturan publish,
// dan menjaga draft tidak bocor ke endpoint publik.
type Service struct {
	repo *Repository
}

// NewService - factory Service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// toAppError - terjemahkan error repository ke AppError yang dipahami handler
func toAppError(err error) error {
	if errors.Is(err, ErrNotFound) {
		return utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Gambar anotasi tidak ditemukan")
	}
	return err
}

// normalizePaging - jaga page/limit selalu masuk akal
func normalizePaging(page, limit int) (int, int) {
	if page < 1 {
		page = defaultPage
	}
	if limit < 1 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return page, limit
}

// Create - buat FoodImage baru; selalu mulai sebagai draft
func (s *Service) Create(req CreateFoodImageRequest, userID string) (*FoodImage, error) {
	if req.PrimaryFoodID != nil && strings.TrimSpace(*req.PrimaryFoodID) == "" {
		req.PrimaryFoodID = nil
	}

	image := &FoodImage{
		ID:            uuid.New().String(),
		Title:         strings.TrimSpace(req.Title),
		ImageURL:      req.ImageURL,
		ThumbnailURL:  req.ThumbnailURL,
		Width:         req.Width,
		Height:        req.Height,
		Status:        StatusDraft,
		PrimaryFoodID: req.PrimaryFoodID,
		CreatedBy:     userID,
	}

	if err := s.repo.Create(image); err != nil {
		return nil, err
	}

	image.Areas = []FoodArea{}
	return image, nil
}

// Get - detail satu gambar beserta areas (admin)
func (s *Service) Get(id string) (*FoodImage, error) {
	image, err := s.repo.FindByID(id)
	if err != nil {
		return nil, toAppError(err)
	}
	return image, nil
}

// GetPublished - detail untuk konsumsi publik; draft dianggap 404
func (s *Service) GetPublished(id string) (*FoodImage, error) {
	image, err := s.repo.FindPublishedByID(id)
	if err != nil {
		return nil, toAppError(err)
	}
	return image, nil
}

// List - daftar gambar untuk panel admin
func (s *Service) List(q ListFoodImagesQuery) (*ListFoodImagesResponse, error) {
	page, limit := normalizePaging(q.Page, q.Limit)

	items, total, err := s.repo.List(q.Status, q.Search, q.PrimaryFoodID, page, limit)
	if err != nil {
		return nil, err
	}

	return &ListFoodImagesResponse{Items: items, Total: total, Page: page, Limit: limit}, nil
}

// ListPublished - daftar published untuk aplikasi user
func (s *Service) ListPublished(page, limit int) (*ListFoodImagesResponse, error) {
	page, limit = normalizePaging(page, limit)

	items, total, err := s.repo.ListPublished(page, limit)
	if err != nil {
		return nil, err
	}

	return &ListFoodImagesResponse{Items: items, Total: total, Page: page, Limit: limit}, nil
}

// ListPublishedByFood - gambar published yang terkait sebuah food master
func (s *Service) ListPublishedByFood(foodID string) ([]FoodImageSummary, error) {
	return s.repo.ListPublishedByFoodID(foodID)
}

// Update - PATCH metadata gambar. Hanya field yang dikirim yang diubah.
func (s *Service) Update(id string, req UpdateFoodImageRequest) (*FoodImage, error) {
	fields := map[string]interface{}{}

	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			return nil, utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Judul tidak boleh kosong")
		}
		fields["title"] = title
	}
	if req.ThumbnailURL != nil {
		fields["thumbnail_url"] = *req.ThumbnailURL
	}
	if req.PrimaryFoodID != nil {
		if strings.TrimSpace(*req.PrimaryFoodID) == "" {
			fields["primary_food_id"] = nil
		} else {
			fields["primary_food_id"] = *req.PrimaryFoodID
		}
	}

	if len(fields) == 0 {
		return s.Get(id)
	}

	if err := s.repo.UpdateFields(id, fields); err != nil {
		return nil, toAppError(err)
	}

	return s.Get(id)
}

// ReplaceAreas - simpan seluruh area (jalur autosave editor).
//
// Menyimpan area TIDAK mengubah status. Gambar yang sudah published dan
// diedit lagi akan tetap published — perubahan langsung terlihat user.
// Editor menampilkan peringatan untuk kasus ini.
func (s *Service) ReplaceAreas(id string, req ReplaceAreasRequest) (*ReplaceAreasResponse, error) {
	image, err := s.repo.FindByID(id)
	if err != nil {
		return nil, toAppError(err)
	}

	areas, err := normalizeAreasForDraft(req.Areas, image.Width, image.Height)
	if err != nil {
		return nil, err
	}

	updatedAt, err := s.repo.ReplaceAreas(id, areas)
	if err != nil {
		return nil, toAppError(err)
	}

	return &ReplaceAreasResponse{
		FoodImageID: id,
		Status:      image.Status,
		AreasCount:  len(areas),
		UpdatedAt:   updatedAt,
	}, nil
}

// Publish - validasi ketat lalu buka gambar ke publik
func (s *Service) Publish(id string) (*FoodImage, error) {
	image, err := s.repo.FindByID(id)
	if err != nil {
		return nil, toAppError(err)
	}

	if err := validateForPublish(image); err != nil {
		return nil, err
	}

	now := time.Now()
	fields := map[string]interface{}{
		"status":       StatusPublished,
		"published_at": now,
	}

	if err := s.repo.UpdateFields(id, fields); err != nil {
		return nil, toAppError(err)
	}

	return s.Get(id)
}

// Unpublish - kembalikan ke draft; endpoint publik langsung 404
func (s *Service) Unpublish(id string) (*FoodImage, error) {
	if _, err := s.repo.FindByID(id); err != nil {
		return nil, toAppError(err)
	}

	fields := map[string]interface{}{
		"status":       StatusDraft,
		"published_at": nil,
	}

	if err := s.repo.UpdateFields(id, fields); err != nil {
		return nil, toAppError(err)
	}

	return s.Get(id)
}

// Delete - hapus gambar beserta seluruh area-nya
func (s *Service) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return toAppError(err)
	}
	return nil
}
