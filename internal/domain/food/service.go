package food

import (
	"atlas_food/internal/domain/annotation"
	"atlas_food/internal/pkg/utils"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Service interface {
	// Admin Food
	CreateFood(req CreateFoodRequest) (*FoodResponse, error)
	GetFoodDetail(id string) (*FoodResponse, error)
	ListFoods(categoryID string, page, limit int) ([]Food, int64, error)
	ListFoodsAdmin(f ListFoodsFilter) ([]Food, int64, error)
	UpdateFood(id string, req UpdateFoodRequest) (*FoodResponse, error)
	DeleteFood(id string) error

	// Admin Food Photos (anotasi + gram terpadu)
	ListFoodPhotos(foodID string) (*FoodPhotoListResponse, error)
	CreateFoodPhoto(foodID, userID string, req CreateFoodPhotoRequest) (*FoodPhotoResponse, error)
	UpdateFoodPhoto(foodID, photoID string, req UpdateFoodPhotoRequest) (*FoodPhotoResponse, error)
	DeleteFoodPhoto(foodID, photoID string) error
	PublishFoodPhoto(foodID, photoID string) (*FoodPhotoResponse, error)
	UnpublishFoodPhoto(foodID, photoID string) (*FoodPhotoResponse, error)

	// Admin Portion Method
	AddPortionMethod(foodID string, req CreatePortionMethodRequest) (*PortionMethodResponse, error)
	ListPortionMethods(foodID string) ([]PortionMethodResponse, error)
	UpdatePortionMethod(id int, req UpdatePortionMethodRequest) (*PortionMethodResponse, error)
	DeletePortionMethod(id int) error

	// Admin Category
	CreateCategory(req CreateCategoryRequest) (*Category, error)
	GetCategory(id string) (*Category, error)
	UpdateCategory(id string, req UpdateCategoryRequest) (*Category, error)
	DeleteCategory(id string) error

	// Admin As-Served
	ListAsServedSets() ([]AsServedSet, error)
	CreateAsServedSet(req CreateAsServedSetRequest) (*AsServedSet, error)
	GetAsServedSet(id string) (*AsServedSetDetailResponse, error)
	UpdateAsServedSet(id string, req UpdateAsServedSetRequest) (*AsServedSet, error)
	DeleteAsServedSet(id string) error
	AddAsServedImages(setID string, reqs []AsServedImageRequest) ([]AsServedImage, error)
	UpdateAsServedImage(id string, req UpdateAsServedImageRequest) (*AsServedImage, error)
	DeleteAsServedImage(id string) error

	// Public/Respondent
	SearchFoods(query string, categoryID string, foodType string, limit int) ([]SearchFoodResponse, error)
	ListCategories() ([]Category, error)
	ListFoodsByCategoryCode(categoryCode string, page, limit int) ([]Food, int64, error)
}

type foodService struct {
	repo Repository
	ann  *annotation.Service
}

func NewService(repo Repository, ann *annotation.Service) Service {
	return &foodService{repo: repo, ann: ann}
}

// CreateFood - menambahkan makanan baru ke database (Admin Only)
func (s *foodService) CreateFood(req CreateFoodRequest) (*FoodResponse, error) {
	food := &Food{
		ID:          uuid.New().String(),
		Code:        req.Code,
		Name:        req.Name,
		LocalName:   req.LocalName,
		Description: req.Description,
		PhotoType:   req.PhotoType,
	}

	if food.PhotoType == "" {
		food.PhotoType = "series"
	}

	if req.CategoryID != "" {
		food.CategoryID = &req.CategoryID
	}

	if err := s.repo.CreateFood(food); err != nil {
		return nil, err
	}

	// Add nutrients
	if len(req.Nutrients) > 0 {
		nutrients := make([]FoodNutrient, len(req.Nutrients))
		for i, n := range req.Nutrients {
			nutrients[i] = FoodNutrient{
				FoodID:         food.ID,
				NutrientTypeID: n.TypeID,
				ValuePer100g:   n.ValuePer100g,
			}
		}
		if err := s.repo.UpsertFoodNutrients(nutrients); err != nil {
			return nil, err
		}
	}

	return s.GetFoodDetail(food.ID)
}

// GetFoodDetail - mengambil detail informasi makanan beserta porsi dan gizinya
func (s *foodService) GetFoodDetail(id string) (*FoodResponse, error) {
	food, err := s.repo.GetFoodByID(id)
	if err != nil {
		return nil, utils.NewAppError(404, "NOT_FOUND", "Makanan tidak ditemukan")
	}

	nutrients, _ := s.repo.GetNutrientsByFoodID(id)
	nutrientMap := make(map[string]NutrientDetail)
	for _, n := range nutrients {
		nutrientMap[n.NutrientType.Code] = NutrientDetail{
			Value: n.ValuePer100g,
			Unit:  n.NutrientType.Unit.Symbol,
		}
	}

	methods, _ := s.repo.GetPortionMethodsByFoodID(id)
	portionPhotos := []PortionPhoto{}

	for _, m := range methods {
		var configData struct {
			WeightGram    float64 `json:"weight_gram"`
			ThumbnailURL  string  `json:"thumbnail_url"`
			SetCode       string  `json:"setCode"`
			SelectionType string  `json:"selectionType"`
		}
		_ = json.Unmarshal([]byte(m.Config), &configData)

		// If as_served method with setCode, resolve AsServedSet → AsServedImages
		if m.MethodType == "as_served" && configData.SetCode != "" {
			set, setErr := s.repo.GetAsServedSetByCode(configData.SetCode)
			if setErr == nil {
				images, imgErr := s.repo.GetAsServedImagesBySetID(set.ID)
				if imgErr == nil && len(images) > 0 {
					for _, img := range images {
						thumbURL := img.ThumbnailURL
						if thumbURL == "" {
							thumbURL = img.ImageURL
						}
						portionPhotos = append(portionPhotos, PortionPhoto{
							ID:           img.ID,
							Label:        img.Label,
							ImageURL:     img.ImageURL,
							ThumbnailURL: thumbURL,
							WeightGram:   img.WeightGram,
							Description:  img.Description,
						})
					}
					continue
				}
			}
		}

		// Fallback: legacy method (direct image_url, no setCode)
		thumbnailURL := configData.ThumbnailURL
		if thumbnailURL == "" {
			thumbnailURL = m.ThumbnailURL
		}

		portionPhotos = append(portionPhotos, PortionPhoto{
			ID:           strconv.Itoa(m.ID),
			Label:        m.Label,
			ImageURL:     m.ImageURL,
			ThumbnailURL: thumbnailURL,
			WeightGram:   configData.WeightGram,
			Description:  m.Description,
		})
	}

	var categoryInfo *CategoryInfo
	if food.Category != nil {
		categoryInfo = &CategoryInfo{
			ID:   food.Category.ID,
			Code: food.Category.Code,
			Name: food.Category.Name,
			Icon: food.Category.Icon,
		}
	}

	return &FoodResponse{
		ID:            food.ID,
		Code:          food.Code,
		Name:          food.Name,
		LocalName:     food.LocalName,
		Description:   food.Description,
		PhotoType:     food.PhotoType,
		CategoryID:    food.CategoryID,
		Category:      categoryInfo,
		IsActive:      food.IsActive,
		Nutrients:     nutrientMap,
		PortionPhotos: portionPhotos,
	}, nil
}

// ListFoods - mengambil daftar semua makanan dengan paginasi (Admin)
func (s *foodService) ListFoods(categoryID string, page, limit int) ([]Food, int64, error) {
	return s.repo.ListFoods(categoryID, page, limit)
}

// ListFoodsAdmin - daftar admin dengan search + filter
func (s *foodService) ListFoodsAdmin(f ListFoodsFilter) ([]Food, int64, error) {
	return s.repo.ListFoodsFiltered(f)
}

// UpdateFood - memperbarui data makanan yang sudah ada (Admin)
func (s *foodService) UpdateFood(id string, req UpdateFoodRequest) (*FoodResponse, error) {
	food, err := s.repo.GetFoodByID(id)
	if err != nil {
		return nil, utils.NewAppError(404, "NOT_FOUND", "Makanan tidak ditemukan")
	}

	if req.Name != "" {
		food.Name = req.Name
	}
	if req.LocalName != "" {
		food.LocalName = req.LocalName
	}
	if req.Description != "" {
		food.Description = req.Description
	}
	if req.PhotoType != "" {
		food.PhotoType = req.PhotoType
	}
	// Dikirim = ubah. String kosong berarti lepaskan kategori (simpan NULL),
	// bukan "biarkan apa adanya".
	if req.CategoryID != nil {
		if trimmed := strings.TrimSpace(*req.CategoryID); trimmed == "" {
			food.CategoryID = nil
		} else {
			food.CategoryID = &trimmed
		}
	}
	if req.IsActive != nil {
		food.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateFood(food); err != nil {
		return nil, err
	}

	// Update nutrients if provided
	if len(req.Nutrients) > 0 {
		nutrients := make([]FoodNutrient, len(req.Nutrients))
		for i, n := range req.Nutrients {
			nutrients[i] = FoodNutrient{
				FoodID:         food.ID,
				NutrientTypeID: n.TypeID,
				ValuePer100g:   n.ValuePer100g,
			}
		}
		if err := s.repo.UpsertFoodNutrients(nutrients); err != nil {
			return nil, err
		}
	}

	return s.GetFoodDetail(id)
}

// DeleteFood - menghapus data makanan dari database (Admin)
func (s *foodService) DeleteFood(id string) error {
	return s.repo.DeleteFood(id)
}

// AddPortionMethod - menambahkan metode pengukuran porsi baru untuk makanan tertentu
func (s *foodService) AddPortionMethod(foodID string, req CreatePortionMethodRequest) (*PortionMethodResponse, error) {
	method := &PortionSizeMethod{
		FoodID:      foodID,
		MethodType:  req.MethodType,
		Label:       req.Label,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Config:      string(req.Config),
	}

	if err := s.repo.CreatePortionMethod(method); err != nil {
		return nil, err
	}

	return &PortionMethodResponse{
		ID:          method.ID,
		MethodType:  method.MethodType,
		Label:       method.Label,
		Description: method.Description,
		ImageURL:    method.ImageURL,
		Config:      req.Config,
	}, nil
}

// ListPortionMethods - melihat metode pengukuran porsi yang tersedia untuk satu makanan
func (s *foodService) ListPortionMethods(foodID string) ([]PortionMethodResponse, error) {
	methods, err := s.repo.GetPortionMethodsByFoodID(foodID)
	if err != nil {
		return nil, err
	}

	resp := make([]PortionMethodResponse, len(methods))
	for i, m := range methods {
		resp[i] = PortionMethodResponse{
			ID:          m.ID,
			MethodType:  m.MethodType,
			Label:       m.Label,
			Description: m.Description,
			ImageURL:    m.ImageURL,
			Config:      safeRawJSON(m.Config),
		}
	}
	return resp, nil
}

// SearchFoods - mencari makanan berdasarkan nama untuk responden (Public)
// Parameter foodType: "food" | "drink" | "" (semua)
func (s *foodService) SearchFoods(query string, categoryID string, foodType string, limit int) ([]SearchFoodResponse, error) {
	foods, err := s.repo.SearchFoods(query, categoryID, foodType, limit)
	if err != nil {
		return nil, err
	}

	resp := make([]SearchFoodResponse, len(foods))
	for i, f := range foods {
		var categoryInfo *CategoryInfo
		if f.Category != nil {
			categoryInfo = &CategoryInfo{
				ID:   f.Category.ID,
				Code: f.Category.Code,
				Name: f.Category.Name,
				Icon: f.Category.Icon,
			}
		}
		resp[i] = SearchFoodResponse{
			ID:        f.ID,
			Code:      f.Code,
			Name:      f.Name,
			LocalName: f.LocalName,
			PhotoType: f.PhotoType,
			Category:  categoryInfo,
		}
	}
	return resp, nil
}

// ListCategories - mengambil daftar semua kategori makanan (Public)
func (s *foodService) ListCategories() ([]Category, error) {
	return s.repo.ListCategories()
}

// ListFoodsByCategoryCode - mengambil daftar makanan berdasarkan category code (Public)
func (s *foodService) ListFoodsByCategoryCode(categoryCode string, page, limit int) ([]Food, int64, error) {
	category, err := s.repo.GetCategoryByCode(categoryCode)
	if err != nil {
		return nil, 0, utils.NewAppError(404, "NOT_FOUND", "Kategori tidak ditemukan")
	}
	return s.repo.ListActiveFoods(category.ID, page, limit)
}
