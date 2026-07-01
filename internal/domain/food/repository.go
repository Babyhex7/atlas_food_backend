package food

import (
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	// Food operations
	CreateFood(food *Food) error
	GetFoodByID(id string) (*Food, error)
	GetFoodByCode(code string) (*Food, error)
	ListFoods(categoryID string, page, limit int) ([]Food, int64, error)
	ListActiveFoods(categoryID string, page, limit int) ([]Food, int64, error)
	UpdateFood(food *Food) error
	DeleteFood(id string) error
	SearchFoods(query string, categoryID string, foodType string, limit int) ([]Food, error)

	// Nutrient operations
	GetNutrientsByFoodID(foodID string) ([]FoodNutrient, error)
	UpsertFoodNutrients(nutrients []FoodNutrient) error
	GetNutrientTypeByID(id int) (*NutrientType, error)

	// Category operations
	ListCategories() ([]Category, error)
	GetCategoryByID(id string) (*Category, error)
	GetCategoryByCode(code string) (*Category, error)

	// Portion Method operations
	GetPortionMethodsByFoodID(foodID string) ([]PortionSizeMethod, error)
	CreatePortionMethod(method *PortionSizeMethod) error
	UpdatePortionMethod(method *PortionSizeMethod) error
	DeletePortionMethod(id int) error

	// As Served operations
	ListAsServedSets() ([]AsServedSet, error)
	CreateAsServedSet(set *AsServedSet) error
	GetAsServedSetByCode(code string) (*AsServedSet, error)
	GetAsServedImagesBySetID(setID string) ([]AsServedImage, error)
	CreateAsServedImages(images []AsServedImage) error
}

type foodRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &foodRepository{db: db}
}

func (r *foodRepository) CreateFood(food *Food) error {
	return r.db.Create(food).Error
}

func (r *foodRepository) GetFoodByID(id string) (*Food, error) {
	var food Food
	err := r.db.Preload("Category").Where("id = ?", id).First(&food).Error
	if err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *foodRepository) GetFoodByCode(code string) (*Food, error) {
	var food Food
	err := r.db.Where("code = ?", code).First(&food).Error
	return &food, err
}

func (r *foodRepository) ListFoods(categoryID string, page, limit int) ([]Food, int64, error) {
	return r.listFoods(categoryID, page, limit, false)
}

func (r *foodRepository) ListActiveFoods(categoryID string, page, limit int) ([]Food, int64, error) {
	return r.listFoods(categoryID, page, limit, true)
}

func (r *foodRepository) listFoods(categoryID string, page, limit int, activeOnly bool) ([]Food, int64, error) {
	var foods []Food
	var total int64
	query := r.db.Model(&Food{}).Preload("Category")
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&foods).Error
	return foods, total, err
}

func (r *foodRepository) UpdateFood(food *Food) error {
	return r.db.Save(food).Error
}

func (r *foodRepository) DeleteFood(id string) error {
	return r.db.Where("id = ?", id).Delete(&Food{}).Error
}

func (r *foodRepository) SearchFoods(query string, categoryID string, foodType string, limit int) ([]Food, error) {
	var foods []Food

	// Join kategori untuk filter type
	q := r.db.Preload("Category").
		Joins("LEFT JOIN categories c ON c.id = foods.category_id").
		Where("foods.is_active = ?", true)

	// Filter berdasarkan categoryID jika disediakan (lebih spesifik)
	if categoryID != "" {
		q = q.Where("foods.category_id = ?", categoryID)
	} else if foodType != "" {
		// Filter berdasarkan food type: drink = category code 'drinks', food = bukan 'drinks'
		switch strings.ToLower(foodType) {
		case "drink":
			q = q.Where("c.code = ?", "drinks")
		case "food":
			q = q.Where("(c.code IS NULL OR c.code != ?)", "drinks")
		}
	}

	trimmed := strings.TrimSpace(query)
	if trimmed != "" {
		likeQuery := "%" + trimmed + "%"
		// FULLTEXT wildcard hanya di akhir kata (nasi*), plus fallback LIKE & pencarian kode
		matchQuery := trimmed + "*"
		q = q.Where(
			"(MATCH(foods.name, foods.local_name) AGAINST(? IN BOOLEAN MODE) OR foods.name LIKE ? OR foods.local_name LIKE ? OR foods.code LIKE ?)",
			matchQuery, likeQuery, likeQuery, likeQuery,
		)
	}

	err := q.Limit(limit).Find(&foods).Error
	return foods, err
}

func (r *foodRepository) GetNutrientsByFoodID(foodID string) ([]FoodNutrient, error) {
	var nutrients []FoodNutrient
	err := r.db.Preload("NutrientType.Unit").Where("food_id = ?", foodID).Find(&nutrients).Error
	return nutrients, err
}

func (r *foodRepository) UpsertFoodNutrients(nutrients []FoodNutrient) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, n := range nutrients {
			if err := tx.Save(&n).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *foodRepository) GetNutrientTypeByID(id int) (*NutrientType, error) {
	var nt NutrientType
	err := r.db.Preload("Unit").Where("id = ?", id).First(&nt).Error
	return &nt, err
}

func (r *foodRepository) ListCategories() ([]Category, error) {
	var categories []Category
	err := r.db.Order("display_order ASC").Find(&categories).Error
	return categories, err
}

func (r *foodRepository) GetCategoryByID(id string) (*Category, error) {
	var category Category
	err := r.db.Where("id = ?", id).First(&category).Error
	return &category, err
}

func (r *foodRepository) GetCategoryByCode(code string) (*Category, error) {
	var category Category
	err := r.db.Where("code = ?", code).First(&category).Error
	return &category, err
}

func (r *foodRepository) GetPortionMethodsByFoodID(foodID string) ([]PortionSizeMethod, error) {
	var methods []PortionSizeMethod
	err := r.db.Where("food_id = ?", foodID).Order("display_order ASC").Find(&methods).Error
	return methods, err
}

func (r *foodRepository) CreatePortionMethod(method *PortionSizeMethod) error {
	return r.db.Create(method).Error
}

func (r *foodRepository) UpdatePortionMethod(method *PortionSizeMethod) error {
	return r.db.Save(method).Error
}

func (r *foodRepository) DeletePortionMethod(id int) error {
	return r.db.Where("id = ?", id).Delete(&PortionSizeMethod{}).Error
}

func (r *foodRepository) ListAsServedSets() ([]AsServedSet, error) {
	var sets []AsServedSet
	err := r.db.Find(&sets).Error
	return sets, err
}

func (r *foodRepository) CreateAsServedSet(set *AsServedSet) error {
	return r.db.Create(set).Error
}

func (r *foodRepository) GetAsServedSetByCode(code string) (*AsServedSet, error) {
	var set AsServedSet
	err := r.db.Where("code = ?", code).First(&set).Error
	return &set, err
}

func (r *foodRepository) GetAsServedImagesBySetID(setID string) ([]AsServedImage, error) {
	var images []AsServedImage
	err := r.db.Where("set_id = ?", setID).Order("display_order ASC").Find(&images).Error
	return images, err
}

func (r *foodRepository) CreateAsServedImages(images []AsServedImage) error {
	return r.db.Create(&images).Error
}
