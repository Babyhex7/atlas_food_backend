package food

// Batas foto per tipe (admin CRUD terpadu)
const (
	MaxSeriesPhotos = 10
	MaxRangePhotos  = 1
)

// CreateFoodPhotoRequest - unggah satu foto makanan (anotasi draft + porsi berbobot)
type CreateFoodPhotoRequest struct {
	Title        string  `json:"title" binding:"required,max=255"`
	ImageURL     string  `json:"image_url" binding:"required,max=500"`
	ThumbnailURL string  `json:"thumbnail_url" binding:"omitempty,max=500"`
	Width        int     `json:"width" binding:"required,gt=0"`
	Height       int     `json:"height" binding:"required,gt=0"`
	WeightGram   float64 `json:"weight_gram" binding:"required,gt=0"`
	Label        string  `json:"label" binding:"omitempty,max=50"`
}

// UpdateFoodPhotoRequest - ubah judul/label/berat (pointer = optional)
type UpdateFoodPhotoRequest struct {
	Title      *string  `json:"title" binding:"omitempty,max=255"`
	Label      *string  `json:"label" binding:"omitempty,max=50"`
	WeightGram *float64 `json:"weight_gram" binding:"omitempty,gt=0"`
}

// FoodPhotoResponse - satu kartu foto di admin (anotasi + gram porsi)
type FoodPhotoResponse struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Label           string  `json:"label"`
	ImageURL        string  `json:"image_url"`
	ThumbnailURL    string  `json:"thumbnail_url"`
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	Status          string  `json:"status"`
	AreasCount      int     `json:"areas_count"`
	WeightGram      float64 `json:"weight_gram"`
	AsServedImageID string  `json:"as_served_image_id"`
	DisplayOrder    int     `json:"display_order"`
	PublishedAt     *string `json:"published_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// FoodPhotoListResponse - daftar + meta batas tipe
type FoodPhotoListResponse struct {
	Items     []FoodPhotoResponse `json:"items"`
	PhotoType string              `json:"photo_type"`
	MaxPhotos int                 `json:"max_photos"`
	Count     int                 `json:"count"`
}
