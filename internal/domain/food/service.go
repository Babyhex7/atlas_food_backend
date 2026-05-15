package food

import (
	"atlas_food/internal/pkg/utils"
	"encoding/json"

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
	SearchFoods(query string, categoryID string, limit int) ([]SearchFoodResponse, error)
	ListCategories() ([]Category, error)
}

type foodService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &foodService{repo: repo}
}

func (s *foodService) CreateFood(req CreateFoodRequest) (*FoodResponse, error) {
	food := &Food{
		ID:          uuid.New().String(),
		Code:        req.Code,
		Name:        req.Name,
		LocalName:   req.LocalName,
		Description: req.Description,
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
	methodResponses := make([]PortionMethodResponse, len(methods))
	for i, m := range methods {
		methodResponses[i] = PortionMethodResponse{
			ID:          m.ID,
			MethodType:  m.MethodType,
			Label:       m.Label,
			Description: m.Description,
			ImageURL:    m.ImageURL,
			Config:      json.RawMessage(m.Config),
		}
	}

	categoryName := ""
	categoryIcon := ""
	if food.Category != nil {
		categoryName = food.Category.Name
		categoryIcon = food.Category.Icon
	}

	return &FoodResponse{
		ID:             food.ID,
		Code:           food.Code,
		Name:           food.Name,
		LocalName:      food.LocalName,
		Description:    food.Description,
		Category:       categoryName,
		Icon:           categoryIcon,
		Nutrients:      nutrientMap,
		PortionMethods: methodResponses,
	}, nil
}

func (s *foodService) ListFoods(categoryID string, page, limit int) ([]Food, int64, error) {
	return s.repo.ListFoods(categoryID, page, limit)
}

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

func (s *foodService) DeleteFood(id string) error {
	return s.repo.DeleteFood(id)
}

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

func (s *foodService) SearchFoods(query string, categoryID string, limit int) ([]SearchFoodResponse, error) {
	foods, err := s.repo.SearchFoods(query, categoryID, limit)
	if err != nil {
		return nil, err
	}

	resp := make([]SearchFoodResponse, len(foods))
	for i, f := range foods {
		catName := ""
		catIcon := ""
		if f.Category != nil {
			catName = f.Category.Name
			catIcon = f.Category.Icon
		}
		resp[i] = SearchFoodResponse{
			ID:        f.ID,
			Code:      f.Code,
			Name:      f.Name,
			LocalName: f.LocalName,
			Category:  catName,
			Icon:      catIcon,
		}
	}
	return resp, nil
}

func (s *foodService) ListCategories() ([]Category, error) {
	return s.repo.ListCategories()
}
