package food

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"atlas_food/internal/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Business rules katalog admin: kategori, as-served set/image, portion method.

// notFound - AppError 404 dengan pesan spesifik entitas
func notFound(entity string) error {
	return utils.NewAppError(http.StatusNotFound, "NOT_FOUND", entity+" tidak ditemukan")
}

// conflict - AppError 409 untuk pelanggaran keunikan / relasi
func conflict(message string) error {
	return utils.NewAppError(http.StatusConflict, "CONFLICT", message)
}

// ============ CATEGORY ============

// CreateCategory - buat kategori baru; kode wajib unik supaya URL Find Food tidak bentrok
func (s *foodService) CreateCategory(req CreateCategoryRequest) (*Category, error) {
	code := strings.TrimSpace(req.Code)

	// Cek duplikat lebih dulu agar admin dapat pesan jelas, bukan error SQL mentah
	if existing, err := s.repo.GetCategoryByCode(code); err == nil && existing != nil {
		return nil, conflict(fmt.Sprintf("Kategori dengan kode %q sudah ada", code))
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	category := &Category{
		Code:         code,
		Name:         strings.TrimSpace(req.Name),
		Icon:         req.Icon,
		DisplayOrder: req.DisplayOrder,
	}

	if err := s.repo.CreateCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategory - ambil kategori berdasarkan ID, balikan error 404 kalau tidak ada
func (s *foodService) GetCategory(id string) (*Category, error) {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("Kategori")
		}
		return nil, err
	}
	return category, nil
}

// UpdateCategory - ubah kategori secara partial (hanya field yang dikirim); kode baru dicek keunikannya
func (s *foodService) UpdateCategory(id string, req UpdateCategoryRequest) (*Category, error) {
	category, err := s.GetCategory(id)
	if err != nil {
		return nil, err
	}

	if req.Code != nil {
		code := strings.TrimSpace(*req.Code)
		if code == "" {
			return nil, utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Kode kategori tidak boleh kosong")
		}
		// Kode dipakai di URL Find Food (/find-food/category/:code), jadi
		// bentroknya harus ditolak eksplisit.
		if code != category.Code {
			if existing, err := s.repo.GetCategoryByCode(code); err == nil && existing != nil {
				return nil, conflict(fmt.Sprintf("Kategori dengan kode %q sudah ada", code))
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		}
		category.Code = code
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Nama kategori tidak boleh kosong")
		}
		category.Name = name
	}
	if req.Icon != nil {
		category.Icon = *req.Icon
	}
	if req.DisplayOrder != nil {
		category.DisplayOrder = *req.DisplayOrder
	}

	if err := s.repo.UpdateCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory - hapus kategori; ditolak (409) kalau masih dipakai makanan
func (s *foodService) DeleteCategory(id string) error {
	if _, err := s.GetCategory(id); err != nil {
		return err
	}

	// Kategori yang masih dipakai food tidak boleh hilang begitu saja —
	// foods.category_id akan menggantung dan Find Food kehilangan filternya.
	count, err := s.repo.CountFoodsByCategoryID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return conflict(fmt.Sprintf("Kategori masih dipakai %d makanan, pindahkan dulu sebelum menghapus", count))
	}

	return s.repo.DeleteCategory(id)
}

// ============ AS-SERVED SET ============

// ListAsServedSets - ambil semua set foto porsi
func (s *foodService) ListAsServedSets() ([]AsServedSet, error) {
	return s.repo.ListAsServedSets()
}

// CreateAsServedSet - buat set foto porsi baru dengan kode unik
func (s *foodService) CreateAsServedSet(req CreateAsServedSetRequest) (*AsServedSet, error) {
	code := strings.TrimSpace(req.Code)

	if existing, err := s.repo.GetAsServedSetByCode(code); err == nil && existing != nil {
		return nil, conflict(fmt.Sprintf("Set foto porsi dengan kode %q sudah ada", code))
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	set := &AsServedSet{
		ID:          uuid.New().String(),
		Code:        code,
		Name:        strings.TrimSpace(req.Name),
		Description: req.Description,
		FoodID:      normalizeOptionalID(&req.FoodID),
		Category:    req.CategoryID,
	}

	if err := s.repo.CreateAsServedSet(set); err != nil {
		return nil, err
	}

	return set, nil
}

// GetAsServedSet - ambil detail set foto porsi beserta seluruh fotonya
func (s *foodService) GetAsServedSet(id string) (*AsServedSetDetailResponse, error) {
	set, err := s.repo.GetAsServedSetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("Set foto porsi")
		}
		return nil, err
	}

	images, err := s.repo.GetAsServedImagesBySetID(id)
	if err != nil {
		return nil, err
	}
	if images == nil {
		images = []AsServedImage{}
	}

	return &AsServedSetDetailResponse{AsServedSet: *set, Images: images}, nil
}

// UpdateAsServedSet - ubah set foto porsi secara partial; kode baru dicek keunikannya
func (s *foodService) UpdateAsServedSet(id string, req UpdateAsServedSetRequest) (*AsServedSet, error) {
	set, err := s.repo.GetAsServedSetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("Set foto porsi")
		}
		return nil, err
	}

	if req.Code != nil {
		code := strings.TrimSpace(*req.Code)
		if code == "" {
			return nil, utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Kode set tidak boleh kosong")
		}
		if code != set.Code {
			if existing, err := s.repo.GetAsServedSetByCode(code); err == nil && existing != nil {
				return nil, conflict(fmt.Sprintf("Set foto porsi dengan kode %q sudah ada", code))
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		}
		set.Code = code
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Nama set tidak boleh kosong")
		}
		set.Name = name
	}
	if req.Description != nil {
		set.Description = *req.Description
	}
	if req.Category != nil {
		set.Category = *req.Category
	}
	if req.FoodID != nil {
		set.FoodID = normalizeOptionalID(req.FoodID)
	}

	if err := s.repo.UpdateAsServedSet(set); err != nil {
		return nil, err
	}

	return set, nil
}

// DeleteAsServedSet - hapus set foto porsi beserta seluruh fotonya
func (s *foodService) DeleteAsServedSet(id string) error {
	if _, err := s.repo.GetAsServedSetByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return notFound("Set foto porsi")
		}
		return err
	}
	return s.repo.DeleteAsServedSet(id)
}

// ============ AS-SERVED IMAGE ============

// AddAsServedImages - tambah beberapa foto porsi sekaligus ke dalam sebuah set
func (s *foodService) AddAsServedImages(setID string, reqs []AsServedImageRequest) ([]AsServedImage, error) {
	if _, err := s.repo.GetAsServedSetByID(setID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("Set foto porsi")
		}
		return nil, err
	}

	if len(reqs) == 0 {
		return nil, utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Tidak ada foto yang dikirim")
	}

	images := make([]AsServedImage, 0, len(reqs))
	for _, req := range reqs {
		images = append(images, AsServedImage{
			ID:           uuid.New().String(),
			SetID:        setID,
			Label:        strings.TrimSpace(req.Label),
			ImageURL:     req.ImageURL,
			ThumbnailURL: req.ThumbnailURL,
			WeightGram:   req.WeightGram,
			Description:  req.Description,
			DisplayOrder: req.DisplayOrder,
		})
	}

	if err := s.repo.CreateAsServedImages(images); err != nil {
		return nil, err
	}

	return images, nil
}

// UpdateAsServedImage - ubah data satu foto porsi secara partial (label, URL, berat gram, urutan)
func (s *foodService) UpdateAsServedImage(id string, req UpdateAsServedImageRequest) (*AsServedImage, error) {
	image, err := s.repo.GetAsServedImageByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("Foto porsi")
		}
		return nil, err
	}

	if req.Label != nil {
		label := strings.TrimSpace(*req.Label)
		if label == "" {
			return nil, utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Label foto tidak boleh kosong")
		}
		image.Label = label
	}
	if req.ImageURL != nil {
		image.ImageURL = *req.ImageURL
	}
	if req.ThumbnailURL != nil {
		image.ThumbnailURL = *req.ThumbnailURL
	}
	if req.WeightGram != nil {
		image.WeightGram = *req.WeightGram
	}
	if req.Description != nil {
		image.Description = *req.Description
	}
	if req.DisplayOrder != nil {
		image.DisplayOrder = *req.DisplayOrder
	}

	if err := s.repo.UpdateAsServedImage(image); err != nil {
		return nil, err
	}

	return image, nil
}

// DeleteAsServedImage - hapus satu foto porsi
func (s *foodService) DeleteAsServedImage(id string) error {
	if _, err := s.repo.GetAsServedImageByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return notFound("Foto porsi")
		}
		return err
	}
	return s.repo.DeleteAsServedImage(id)
}

// ============ PORTION METHOD ============

// UpdatePortionMethod - ubah metode porsi secara partial lalu kembalikan bentuk response-nya
func (s *foodService) UpdatePortionMethod(id int, req UpdatePortionMethodRequest) (*PortionMethodResponse, error) {
	method, err := s.repo.GetPortionMethodByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("Metode porsi")
		}
		return nil, err
	}

	if req.MethodType != nil {
		method.MethodType = *req.MethodType
	}
	if req.Label != nil {
		label := strings.TrimSpace(*req.Label)
		if label == "" {
			return nil, utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Label metode tidak boleh kosong")
		}
		method.Label = label
	}
	if req.Description != nil {
		method.Description = *req.Description
	}
	if req.ImageURL != nil {
		method.ImageURL = *req.ImageURL
	}
	if req.ThumbnailURL != nil {
		method.ThumbnailURL = *req.ThumbnailURL
	}
	if req.Config != nil {
		method.Config = *req.Config
	}
	if req.DisplayOrder != nil {
		method.DisplayOrder = *req.DisplayOrder
	}
	if req.IsActive != nil {
		method.IsActive = *req.IsActive
	}

	if err := s.repo.UpdatePortionMethod(method); err != nil {
		return nil, err
	}

	return &PortionMethodResponse{
		ID:          method.ID,
		MethodType:  method.MethodType,
		Label:       method.Label,
		Description: method.Description,
		ImageURL:    method.ImageURL,
		Config:      safeRawJSON(method.Config),
	}, nil
}

// DeletePortionMethod - hapus metode porsi berdasarkan ID
func (s *foodService) DeletePortionMethod(id int) error {
	if _, err := s.repo.GetPortionMethodByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return notFound("Metode porsi")
		}
		return err
	}
	return s.repo.DeletePortionMethod(id)
}

// normalizeOptionalID - ubah pointer string kosong jadi nil agar tersimpan NULL.
// String kosong akan melanggar foreign key ke foods(id).
func normalizeOptionalID(id *string) *string {
	if id == nil || strings.TrimSpace(*id) == "" {
		return nil
	}
	return id
}

// safeRawJSON - bungkus kolom config jadi JSON yang aman di-marshal.
//
// json.RawMessage kosong membuat encoder gagal ("unexpected end of JSON
// input") dan seluruh response berubah jadi 500. Kolom config bisa kosong
// untuk baris hasil seed lama, jadi nilainya dinormalkan ke `null`.
func safeRawJSON(raw string) json.RawMessage {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || !json.Valid([]byte(trimmed)) {
		return json.RawMessage("null")
	}
	return json.RawMessage(trimmed)
}
