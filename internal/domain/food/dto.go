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
	// omitempty wajib: tanpa itu `oneof` menolak string kosong, sehingga
	// create gagal total kalau klien tidak mengirim photo_type.
	// Service sudah memberi default "series".
	PhotoType  string                `json:"photo_type" binding:"omitempty,oneof=series range"`
	CategoryID string                `json:"category_id"`
	Nutrients  []FoodNutrientRequest `json:"nutrients"`
}

// UpdateFoodRequest - DTO untuk update food
type UpdateFoodRequest struct {
	Name        string `json:"name"`
	LocalName   string `json:"local_name"`
	Description string `json:"description"`
	PhotoType   string `json:"photo_type" binding:"omitempty,oneof=series range"`
	// Pointer agar admin bisa membedakan "tidak diubah" (field tidak dikirim)
	// dari "lepaskan kategori" (dikirim null / string kosong). Dengan tipe
	// string biasa, memilih "Tanpa kategori" di form diam-diam tidak berefek.
	CategoryID *string               `json:"category_id"`
	Nutrients  []FoodNutrientRequest `json:"nutrients"`
	IsActive   *bool                 `json:"is_active"`
}

// CategoryInfo - DTO untuk category dalam food responses
type CategoryInfo struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

// FoodResponse - DTO untuk response detail makanan
//
// CategoryID dan IsActive wajib ikut: form edit admin memuat nilainya dari
// sini. Tanpa keduanya, dropdown kategori selalu tampil kosong dan checkbox
// "aktif" selalu tak tercentang, meski datanya ada di database.
type FoodResponse struct {
	ID            string                    `json:"id"`
	Code          string                    `json:"code"`
	Name          string                    `json:"name"`
	LocalName     string                    `json:"local_name"`
	Description   string                    `json:"description"`
	PhotoType     string                    `json:"photo_type"`
	CategoryID    *string                   `json:"category_id"`
	Category      *CategoryInfo             `json:"category,omitempty"`
	IsActive      bool                      `json:"is_active"`
	Nutrients     map[string]NutrientDetail `json:"nutrients"`
	PortionPhotos []PortionPhoto            `json:"portion_photos,omitempty"`
}

// PortionPhoto - detail foto porsi untuk response public
type PortionPhoto struct {
	ID           string  `json:"id"`
	Label        string  `json:"label"`
	ImageURL     string  `json:"image_url"`
	ThumbnailURL string  `json:"thumbnail_url"`
	WeightGram   float64 `json:"weight_gram"`
	Description  string  `json:"description"`
	// FoodImageID - tautan ke food_images (anotasi published) bila ada
	FoodImageID string `json:"food_image_id,omitempty"`
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
	ID        string        `json:"id"`
	Code      string        `json:"code"`
	Name      string        `json:"name"`
	LocalName string        `json:"local_name"`
	PhotoType string        `json:"photo_type"`
	Category  *CategoryInfo `json:"category,omitempty"`
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
