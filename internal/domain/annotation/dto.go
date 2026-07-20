package annotation

import "time"

// CreateFoodImageRequest - body POST /admin/food-images (setelah file di-upload)
type CreateFoodImageRequest struct {
	Title         string  `json:"title" binding:"required,max=255"`
	ImageURL      string  `json:"image_url" binding:"required,max=500"`
	ThumbnailURL  string  `json:"thumbnail_url" binding:"omitempty,max=500"`
	Width         int     `json:"width" binding:"required,gt=0"`
	Height        int     `json:"height" binding:"required,gt=0"`
	PrimaryFoodID *string `json:"primary_food_id"`
}

// UpdateFoodImageRequest - body PATCH /admin/food-images/:id
// Semua field pointer agar bisa membedakan "tidak dikirim" vs "dikosongkan".
type UpdateFoodImageRequest struct {
	Title         *string `json:"title" binding:"omitempty,max=255"`
	ThumbnailURL  *string `json:"thumbnail_url" binding:"omitempty,max=500"`
	PrimaryFoodID *string `json:"primary_food_id"`
}

// AreaInput - satu area dalam payload replace-all
type AreaInput struct {
	ID      *string `json:"id"`
	Name    string  `json:"name" binding:"required,max=255"`
	FoodID  *string `json:"food_id"`
	Polygon Polygon `json:"polygon" binding:"required"`
	ZIndex  int     `json:"z_index"`
}

// ReplaceAreasRequest - body PUT /admin/food-images/:id/areas (dipakai autosave editor)
type ReplaceAreasRequest struct {
	Areas []AreaInput `json:"areas"`
}

// ReplaceAreasResponse - ringkasan hasil autosave
type ReplaceAreasResponse struct {
	FoodImageID string    `json:"food_image_id"`
	Status      string    `json:"status"`
	AreasCount  int       `json:"areas_count"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ListFoodImagesQuery - query string GET /admin/food-images
type ListFoodImagesQuery struct {
	Status string `form:"status" binding:"omitempty,oneof=draft published archived"`
	Search string `form:"search"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

// ListFoodImagesResponse - hasil list dengan pagination
type ListFoodImagesResponse struct {
	Items []FoodImageSummary `json:"items"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}

// FoodImageSummary - bentuk ringkas untuk tabel list (tanpa polygon)
type FoodImageSummary struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	ImageURL     string     `json:"image_url"`
	ThumbnailURL string     `json:"thumbnail_url"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Status       string     `json:"status"`
	AreasCount   int        `json:"areas_count"`
	PublishedAt  *time.Time `json:"published_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
