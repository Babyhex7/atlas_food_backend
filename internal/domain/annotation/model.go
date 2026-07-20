package annotation

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Status anotasi — hanya `published` yang boleh dibaca publik
const (
	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusArchived  = "archived"
)

// Point - satu titik polygon dalam pixel space gambar asli ([x, y])
type Point [2]float64

func (p Point) X() float64 { return p[0] }
func (p Point) Y() float64 { return p[1] }

// Polygon - kumpulan titik yang disimpan sebagai kolom JSON MySQL
type Polygon []Point

// Value - serialize Polygon ke JSON untuk disimpan GORM
func (p Polygon) Value() (driver.Value, error) {
	if p == nil {
		return "[]", nil
	}
	b, err := json.Marshal([]Point(p))
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan - deserialize kolom JSON MySQL ke Polygon
func (p *Polygon) Scan(value interface{}) error {
	if value == nil {
		*p = Polygon{}
		return nil
	}

	var raw []byte
	switch v := value.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return errors.New("annotation: tipe kolom polygon tidak didukung")
	}

	if len(raw) == 0 {
		*p = Polygon{}
		return nil
	}

	var points []Point
	if err := json.Unmarshal(raw, &points); err != nil {
		return err
	}
	*p = Polygon(points)
	return nil
}

// FoodImage - gambar scene yang dianotasi (tabel food_images)
type FoodImage struct {
	ID            string     `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	Title         string     `gorm:"type:varchar(255);not null" json:"title"`
	ImageURL      string     `gorm:"type:varchar(500);not null" json:"image_url"`
	ThumbnailURL  string     `gorm:"type:varchar(500)" json:"thumbnail_url"`
	Width         int        `gorm:"not null" json:"width"`
	Height        int        `gorm:"not null" json:"height"`
	Status        string     `gorm:"type:enum('draft','published','archived');not null;default:'draft';index" json:"status"`
	PrimaryFoodID *string    `gorm:"type:char(36)" json:"primary_food_id"`
	CreatedBy     string     `gorm:"type:char(36);not null" json:"created_by"`
	PublishedAt   *time.Time `json:"published_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Areas dimuat lewat Preload; urut z_index agar hit-test konsisten dengan FE
	Areas []FoodArea `gorm:"foreignKey:FoodImageID" json:"areas"`
}

func (FoodImage) TableName() string {
	return "food_images"
}

// FoodArea - satu region polygon di dalam sebuah FoodImage (tabel food_areas)
type FoodArea struct {
	ID          string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	FoodImageID string    `gorm:"type:char(36);not null;index" json:"food_image_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	FoodID      *string   `gorm:"type:char(36)" json:"food_id"`
	Polygon     Polygon   `gorm:"type:json;not null" json:"polygon"`
	ZIndex      int       `gorm:"not null;default:0" json:"z_index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (FoodArea) TableName() string {
	return "food_areas"
}
