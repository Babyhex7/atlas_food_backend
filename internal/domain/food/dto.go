package food

import "encoding/json"

// FoodNutrientRequest - DTO untuk input nutrisi makanan
type FoodNutrientRequest struct {
	TypeID       int     `json:"type_id" binding:"required"`
	ValuePer100g float64 `json:"value_per_100g" binding:"required"`
}

// CreateFoodRequest - DTO untuk create food
type CreateFoodRequest struct {
	Code        string                `json:"code" binding:"required"`
	Name        string                `json:"name" binding:"required"`
	LocalName   string                `json:"local_name"`
	Description string                `json:"description"`
	CategoryID  string                `json:"category_id"`
	Nutrients   []FoodNutrientRequest `json:"nutrients"`
}

// UpdateFoodRequest - DTO untuk update food
type UpdateFoodRequest struct {
	Name        string                `json:"name"`
	LocalName   string                `json:"local_name"`
	Description string                `json:"description"`
	CategoryID  string                `json:"category_id"`
	Nutrients   []FoodNutrientRequest `json:"nutrients"`
	IsActive    *bool                 `json:"is_active"`
}

// FoodResponse - DTO untuk response detail makanan
type FoodResponse struct {
	ID          string            `json:"id"`
	Code        string            `json:"code"`
	Name        string            `json:"name"`
	LocalName   string            `json:"local_name"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	Icon        string            `json:"icon,omitempty"`
	Nutrients   map[string]NutrientDetail `json:"nutrients"`
	PortionMethods []PortionMethodResponse `json:"portion_methods,omitempty"`
}

// NutrientDetail - detail nutrisi untuk response
type NutrientDetail struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

// PortionMethodResponse - detail portion method untuk response
type PortionMethodResponse struct {
	ID          int             `json:"id"`
	MethodType  string          `json:"method_type"`
	Label       string          `json:"label"`
	Description string          `json:"description"`
	ImageURL    string          `json:"image_url"`
	Config      json.RawMessage `json:"config"`
}

// SearchFoodResponse - DTO untuk hasil pencarian makanan
type SearchFoodResponse struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	LocalName string `json:"local_name"`
	Category  string `json:"category"`
	Icon      string `json:"icon"`
}

// CreatePortionMethodRequest - DTO untuk tambah portion method
type CreatePortionMethodRequest struct {
	MethodType  string          `json:"method_type" binding:"required,oneof=as_served guide_image weight"`
	Label       string          `json:"label" binding:"required"`
	Description string          `json:"description"`
	ImageURL    string          `json:"image_url"`
	Config      json.RawMessage `json:"config" binding:"required"`
}

// CreateAsServedSetRequest - DTO untuk create as served set
type CreateAsServedSetRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	CategoryID  string `json:"category"`
	FoodID      string `json:"food_id"`
}

// AsServedSetResponse - DTO untuk response as served set
type AsServedSetResponse struct {
	ID         string `json:"id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	ImageCount int    `json:"image_count"`
}
