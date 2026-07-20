package food

// DTO untuk modul katalog admin: kategori, as-served set/image, portion method.
// Dipisah dari dto.go agar file per-tanggung-jawab tetap mudah dibaca.

// ============ CATEGORY ============

// CreateCategoryRequest - body POST /admin/categories
type CreateCategoryRequest struct {
	Code         string `json:"code" binding:"required,max=50"`
	Name         string `json:"name" binding:"required,max=255"`
	Icon         string `json:"icon" binding:"omitempty,max=50"`
	DisplayOrder int    `json:"display_order"`
}

// UpdateCategoryRequest - body PUT /admin/categories/:id
// Pointer agar field yang tidak dikirim tidak ikut tertimpa.
type UpdateCategoryRequest struct {
	Code         *string `json:"code" binding:"omitempty,max=50"`
	Name         *string `json:"name" binding:"omitempty,max=255"`
	Icon         *string `json:"icon" binding:"omitempty,max=50"`
	DisplayOrder *int    `json:"display_order"`
}

// ============ AS-SERVED SET ============

// UpdateAsServedSetRequest - body PUT /admin/as-served-sets/:id
type UpdateAsServedSetRequest struct {
	Code        *string `json:"code" binding:"omitempty,max=50"`
	Name        *string `json:"name" binding:"omitempty,max=255"`
	Description *string `json:"description"`
	FoodID      *string `json:"food_id"`
	Category    *string `json:"category" binding:"omitempty,max=50"`
}

// AsServedImageRequest - satu foto porsi berbobot gram
type AsServedImageRequest struct {
	Label        string  `json:"label" binding:"required,max=50"`
	ImageURL     string  `json:"image_url" binding:"required,max=500"`
	ThumbnailURL string  `json:"thumbnail_url" binding:"omitempty,max=500"`
	WeightGram   float64 `json:"weight_gram" binding:"required,gt=0"`
	Description  string  `json:"description" binding:"omitempty,max=255"`
	DisplayOrder int     `json:"display_order"`
}

// UpdateAsServedImageRequest - body PUT /admin/as-served-images/:imageId
type UpdateAsServedImageRequest struct {
	Label        *string  `json:"label" binding:"omitempty,max=50"`
	ImageURL     *string  `json:"image_url" binding:"omitempty,max=500"`
	ThumbnailURL *string  `json:"thumbnail_url" binding:"omitempty,max=500"`
	WeightGram   *float64 `json:"weight_gram" binding:"omitempty,gt=0"`
	Description  *string  `json:"description" binding:"omitempty,max=255"`
	DisplayOrder *int     `json:"display_order"`
}

// AsServedSetDetailResponse - set beserta seluruh foto porsinya
type AsServedSetDetailResponse struct {
	AsServedSet
	Images []AsServedImage `json:"images"`
}

// ============ PORTION METHOD ============

// UpdatePortionMethodRequest - body PUT /admin/portion-methods/:methodId
type UpdatePortionMethodRequest struct {
	MethodType   *string `json:"method_type" binding:"omitempty,oneof=as_served guide_image weight"`
	Label        *string `json:"label" binding:"omitempty,max=255"`
	Description  *string `json:"description" binding:"omitempty,max=255"`
	ImageURL     *string `json:"image_url" binding:"omitempty,max=500"`
	ThumbnailURL *string `json:"thumbnail_url" binding:"omitempty,max=500"`
	Config       *string `json:"config"`
	DisplayOrder *int    `json:"display_order"`
	IsActive     *bool   `json:"is_active"`
}
