package ai

import "time"

// MissingFoodSummary - ringkasan makanan hilang yang dilaporkan responden
type MissingFoodSummary struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ReportCount    int       `json:"report_count"`
	LastReportedAt time.Time `json:"last_reported_at"`
	IsInCatalog    bool      `json:"is_in_catalog"`
	CatalogFoodID  *string   `json:"catalog_food_id,omitempty"`
}

// GenerateFoodAIRequest - request dari admin untuk generate gizi & porsi AI
type GenerateFoodAIRequest struct {
	FoodName    string  `json:"food_name" binding:"required"`
	Description string  `json:"description"`
	CategoryID  *string `json:"category_id,omitempty"`
}

// GeneratedPortionView - tampilan porsi yang dihasilkan AI
type GeneratedPortionView struct {
	Code         string  `json:"code"` // A, B, C, D
	Label        string  `json:"label"` // Kecil, Sedang, Besar, Jumbo
	WeightGram   float64 `json:"weight_gram"`
	Description  string  `json:"description"`
	VisualPrompt string  `json:"visual_prompt"`
	ImageURL     string  `json:"image_url"`
}

// GenerateFoodAIResponse - response preview untuk admin sebelum dipublish
type GenerateFoodAIResponse struct {
	FoodName        string                 `json:"food_name"`
	LocalName       string                 `json:"local_name"`
	Description     string                 `json:"description"`
	CategoryID      *string                `json:"category_id,omitempty"`
	CategoryName    string                 `json:"category_name"`
	PhotoType       string                 `json:"photo_type"` // series, range, guide
	CaloriesPer100g float64                `json:"calories_per_100g"`
	ProteinPer100g  float64                `json:"protein_per_100g"`
	FatPer100g      float64                `json:"fat_per_100g"`
	CarbsPer100g    float64                `json:"carbs_per_100g"`
	Portions        []GeneratedPortionView `json:"portions"`
}

// PublishFoodAIRequest - request admin setelah mereview/mengedit untuk simpan ke database
type PublishFoodAIRequest struct {
	FoodName        string                 `json:"food_name" binding:"required"`
	LocalName       string                 `json:"local_name"`
	Description     string                 `json:"description"`
	CategoryID      *string                `json:"category_id,omitempty"`
	PhotoType       string                 `json:"photo_type"` // series, range, guide (default: series)
	CaloriesPer100g float64                `json:"calories_per_100g"`
	ProteinPer100g  float64                `json:"protein_per_100g"`
	FatPer100g      float64                `json:"fat_per_100g"`
	CarbsPer100g    float64                `json:"carbs_per_100g"`
	Portions        []GeneratedPortionView `json:"portions" binding:"required"`
}

// PublishFoodAIResponse - response sukses setelah publikasi
type PublishFoodAIResponse struct {
	FoodID     string `json:"food_id"`
	FoodCode   string `json:"food_code"`
	FoodName   string `json:"food_name"`
	SetID      string `json:"set_id"`
	TotalImage int    `json:"total_image"`
	Message    string `json:"message"`
}
