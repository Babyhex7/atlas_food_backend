package food

import (
	"atlas_food/internal/pkg/utils"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
)

type Service interface {
	// Admin Food
	CreateFood(req CreateFoodRequest) (*FoodResponse, error)
	GetFoodDetail(id string) (*FoodResponse, error)
	ListFoods(categoryID string, page, limit int) ([]Food, int64, error)
	UpdateFood(id string, req UpdateFoodRequest) (*FoodResponse, error)
	DeleteFood(id string) error

	// Admin Portion Method
	AddPortionMethod(foodID string, req CreatePortionMethodRequest) (*PortionMethodResponse, error)
	ListPortionMethods(foodID string) ([]PortionMethodResponse, error)

	// Public/Respondent
	SearchFoods(query string, categoryID string, foodType string, limit int) ([]SearchFoodResponse, error)
	ListCategories() ([]Category, error)
	ListFoodsByCategoryCode(categoryCode string, page, limit int) ([]Food, int64, error)
}

type foodService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &foodService{repo: repo}
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
	portionPhotos := make([]PortionPhoto, len(methods))
	for i, m := range methods {
		var configData struct {
			WeightGram   float64 `json:"weight_gram"`
			ThumbnailURL string  `json:"thumbnail_url"`
		}
		_ = json.Unmarshal([]byte(m.Config), &configData)

		thumbnailURL := configData.ThumbnailURL
		if thumbnailURL == "" {
			thumbnailURL = m.ThumbnailURL
		}

		portionPhotos[i] = PortionPhoto{
			ID:           strconv.Itoa(m.ID),
			Label:        m.Label,
			ImageURL:     m.ImageURL,
			ThumbnailURL: thumbnailURL,
			WeightGram:   configData.WeightGram,
			Description:  m.Description,
		}
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
		ID:             food.ID,
		Code:           food.Code,
		Name:           food.Name,
		LocalName:      food.LocalName,
		Description:    food.Description,
		PhotoType:      food.PhotoType,
		Category:       categoryInfo,
		Nutrients:      nutrientMap,
		PortionPhotos:  portionPhotos,
	}, nil
}

// ListFoods - mengambil daftar semua makanan dengan paginasi (Admin)
func (s *foodService) ListFoods(categoryID string, page, limit int) ([]Food, int64, error) {
	return s.repo.ListFoods(categoryID, page, limit)
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
	if req.CategoryID != "" {
		food.CategoryID = &req.CategoryID
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
			Config:      json.RawMessage(m.Config),
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
