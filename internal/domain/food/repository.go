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
	ListFoodsFiltered(f ListFoodsFilter) ([]Food, int64, error)
	ListActiveFoods(categoryID string, page, limit int) ([]Food, int64, error)
	UpdateFood(food *Food) error
	DeleteFood(id string) error
	SearchFoods(query string, categoryID string, foodType string, limit int) ([]Food, error)

	// Public Food operations (Find Your Food)
	SearchFoodsPublic(query string, foodType string, limit int) ([]Food, error)
	GetFoodWithPortionPhotos(foodID string) (*Food, []AsServedImage, error)
	GetFoodsByCategory(categoryCode string, limit int) ([]Food, error)
	GetFoodNutrients(foodID string) ([]FoodNutrient, error)
	GetAllCategories() ([]Category, error)

	// Nutrient operations
	GetNutrientsByFoodID(foodID string) ([]FoodNutrient, error)
	UpsertFoodNutrients(nutrients []FoodNutrient) error
	GetNutrientTypeByID(id int) (*NutrientType, error)

	// Category operations
	ListCategories() ([]Category, error)
	GetCategoryByID(id string) (*Category, error)
	GetCategoryByCode(code string) (*Category, error)
	CreateCategory(category *Category) error
	UpdateCategory(category *Category) error
	DeleteCategory(id string) error
	CountFoodsByCategoryID(id string) (int64, error)

	// Portion Method operations
	GetPortionMethodsByFoodID(foodID string) ([]PortionSizeMethod, error)
	CreatePortionMethod(method *PortionSizeMethod) error
	UpdatePortionMethod(method *PortionSizeMethod) error
	DeletePortionMethod(id int) error

	// As Served operations
	ListAsServedSets() ([]AsServedSet, error)
	CreateAsServedSet(set *AsServedSet) error
	GetAsServedSetByCode(code string) (*AsServedSet, error)
	GetAsServedSetsByFoodID(foodID string) ([]AsServedSet, error)
	GetAsServedImagesBySetID(setID string) ([]AsServedImage, error)
	CreateAsServedImages(images []AsServedImage) error
	GetAsServedSetByID(id string) (*AsServedSet, error)
	UpdateAsServedSet(set *AsServedSet) error
	DeleteAsServedSet(id string) error
	GetAsServedImageByID(id string) (*AsServedImage, error)
	GetAsServedImageByDescription(setID, description string) (*AsServedImage, error)
	UpdateAsServedImage(image *AsServedImage) error
	DeleteAsServedImage(id string) error

	// Portion method lookup (untuk update/delete per-id)
	GetPortionMethodByID(id int) (*PortionSizeMethod, error)
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
	return r.ListFoodsFiltered(ListFoodsFilter{CategoryID: categoryID, Page: page, Limit: limit})
}

func (r *foodRepository) ListActiveFoods(categoryID string, page, limit int) ([]Food, int64, error) {
	return r.ListFoodsFiltered(ListFoodsFilter{CategoryID: categoryID, Page: page, Limit: limit, ActiveOnly: true})
}

// ListFoodsFilter - filter daftar admin
type ListFoodsFilter struct {
	CategoryID string
	Search     string
	PhotoType  string
	IsActive   *bool
	ActiveOnly bool
	Page       int
	Limit      int
}

func (r *foodRepository) ListFoodsFiltered(f ListFoodsFilter) ([]Food, int64, error) {
	var foods []Food
	var total int64
	query := r.db.Model(&Food{}).Preload("Category")

	if f.CategoryID != "" {
		query = query.Where("category_id = ?", f.CategoryID)
	}
	if f.ActiveOnly {
		query = query.Where("is_active = ?", true)
	} else if f.IsActive != nil {
		query = query.Where("is_active = ?", *f.IsActive)
	}
	if f.PhotoType == "series" || f.PhotoType == "range" {
		query = query.Where("photo_type = ?", f.PhotoType)
	}
	if s := strings.TrimSpace(f.Search); s != "" {
		like := "%" + s + "%"
		query = query.Where(
			"(foods.name LIKE ? OR foods.local_name LIKE ? OR foods.code LIKE ?)",
			like, like, like,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := f.Page
	limit := f.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit
	err := query.Order("foods.name ASC").Offset(offset).Limit(limit).Find(&foods).Error
	return foods, total, err
}

func (r *foodRepository) listFoods(categoryID string, page, limit int, activeOnly bool) ([]Food, int64, error) {
	return r.ListFoodsFiltered(ListFoodsFilter{
		CategoryID: categoryID,
		Page:       page,
		Limit:      limit,
		ActiveOnly: activeOnly,
	})
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

func (r *foodRepository) GetAsServedSetsByFoodID(foodID string) ([]AsServedSet, error) {
	var sets []AsServedSet
	err := r.db.Where("food_id = ?", foodID).Order("created_at ASC").Find(&sets).Error
	return sets, err
}

func (r *foodRepository) GetAsServedImageByDescription(setID, description string) (*AsServedImage, error) {
	var image AsServedImage
	err := r.db.Where("set_id = ? AND description = ?", setID, description).First(&image).Error
	return &image, err
}

func (r *foodRepository) GetAsServedImagesBySetID(setID string) ([]AsServedImage, error) {
	var images []AsServedImage
	err := r.db.Where("set_id = ?", setID).Order("display_order ASC").Find(&images).Error
	return images, err
}

func (r *foodRepository) CreateAsServedImages(images []AsServedImage) error {
	return r.db.Create(&images).Error
}

// GetFoodWithPortionPhotos - Get food detail with portion photos (Public)
// Optimized: No N+1 query, uses JOIN and Preload
func (r *foodRepository) GetFoodWithPortionPhotos(foodID string) (*Food, []AsServedImage, error) {
	// Use channels for concurrent operations
	type result struct {
		food          *Food
		portionPhotos []AsServedImage
		err           error
	}
	
	resultChan := make(chan result, 2)
	
	// Goroutine 1: Fetch food with category (preloaded)
	go func() {
		var f Food
		err := r.db.Preload("Category").
			Where("id = ? AND is_active = ?", foodID, true).
			First(&f).Error
		resultChan <- result{food: &f, err: err}
	}()
	
	// Goroutine 2: Fetch portion photos with JOIN (avoid N+1)
	go func() {
		var photos []AsServedImage
		err := r.db.Table("as_served_images asi").
			Select("asi.*").
			Joins("JOIN as_served_sets ass ON ass.id = asi.set_id").
			Where("ass.food_id = ?", foodID).
			Order("asi.display_order ASC").
			Find(&photos).Error
		resultChan <- result{portionPhotos: photos, err: err}
	}()
	
	// Collect results from both goroutines
	var foodResult, photosResult result
	for i := 0; i < 2; i++ {
		res := <-resultChan
		if res.food != nil {
			foodResult = res
		} else {
			photosResult = res
		}
	}
	
	// Check for errors
	if foodResult.err != nil {
		return nil, nil, foodResult.err
	}
	if photosResult.err != nil {
		return foodResult.food, []AsServedImage{}, nil // Return empty photos if error
	}
	
	return foodResult.food, photosResult.portionPhotos, nil
}

// SearchFoodsPublic - Search foods with FULLTEXT (Public, no auth required)
// Optimized: Single query with Preload, no N+1
//
// foodType: "food" | "drink" | "" (kosong = semua). Filter mengikuti aturan yang
// sama dengan SearchFoods (authenticated): drink = kategori 'drinks'.
func (r *foodRepository) SearchFoodsPublic(query string, foodType string, limit int) ([]Food, error) {
	var foods []Food

	// LEFT JOIN kategori supaya filter food_type bisa memakai categories.code.
	// Select dibuat eksplisit: tanpa itu query menjadi "SELECT *" atas hasil join,
	// dan kolom bernama sama di kedua tabel (id, code, name, created_at) bisa
	// saling menimpa saat di-scan ke struct Food.
	q := r.db.Preload("Category"). // Preload to avoid N+1
					Select("foods.*").
					Joins("LEFT JOIN categories c ON c.id = foods.category_id").
					Where("foods.is_active = ?", true)

	switch strings.ToLower(strings.TrimSpace(foodType)) {
	case "drink":
		q = q.Where("c.code = ?", "drinks")
	case "food":
		q = q.Where("(c.code IS NULL OR c.code != ?)", "drinks")
	}

	trimmed := strings.TrimSpace(query)
	if trimmed != "" {
		// FULLTEXT sendirian gagal untuk pencarian parsial di tengah kata dan untuk
		// token di bawah innodb_ft_min_token_size. Gabungkan dengan LIKE + kode agar
		// hasil tetap muncul, mengikuti perilaku SearchFoods (authenticated).
		likeQuery := "%" + trimmed + "%"
		matchQuery := trimmed + "*"
		q = q.Where(
			"(MATCH(foods.name, foods.local_name) AGAINST(? IN BOOLEAN MODE) OR foods.name LIKE ? OR foods.local_name LIKE ? OR foods.code LIKE ?)",
			matchQuery, likeQuery, likeQuery, likeQuery,
		)
	}

	err := q.Limit(limit).Find(&foods).Error
	return foods, err
}

// GetFoodsByCategory - Get all foods in a category (Public)
// Optimized: Single JOIN query, no N+1
func (r *foodRepository) GetFoodsByCategory(categoryCode string, limit int) ([]Food, error) {
	var foods []Food
	
	err := r.db.Preload("Category").
		Joins("JOIN categories ON categories.id = foods.category_id").
		Where("categories.code = ? AND foods.is_active = ?", categoryCode, true).
		Limit(limit).
		Find(&foods).Error
	
	return foods, err
}

// GetAllCategories - Get all categories (Public)
func (r *foodRepository) GetAllCategories() ([]Category, error) {
	var categories []Category
	err := r.db.Order("display_order ASC").Find(&categories).Error
	return categories, err
}

// GetFoodNutrients - Get nutrients for a food
// Optimized: Uses Preload to avoid N+1
func (r *foodRepository) GetFoodNutrients(foodID string) ([]FoodNutrient, error) {
	var nutrients []FoodNutrient
	err := r.db.
		Preload("NutrientType.Unit"). // Nested preload
		Where("food_id = ?", foodID).
		Find(&nutrients).Error
	return nutrients, err
}
