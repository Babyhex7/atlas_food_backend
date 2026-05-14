package food

import (
	"time"
)

// Category - model untuk tabel categories
type Category struct {
	ID           string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	Code         string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Icon         string    `gorm:"type:varchar(50)" json:"icon"`
	DisplayOrder int       `gorm:"default:0" json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
}

func (Category) TableName() string {
	return "categories"
}

// Food - model untuk tabel foods
type Food struct {
	ID          string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	LocalName   string    `gorm:"type:varchar(255)" json:"local_name"`
	Description string    `gorm:"type:text" json:"description"`
	CategoryID  *string   `gorm:"type:char(36)" json:"category_id"`
	Category    *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Food) TableName() string {
	return "foods"
}

// NutrientUnit - model untuk tabel nutrient_units
type NutrientUnit struct {
	ID     int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code   string `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"`
	Name   string `gorm:"type:varchar(50);not null" json:"name"`
	Symbol string `gorm:"type:varchar(10);not null" json:"symbol"`
}

func (NutrientUnit) TableName() string {
	return "nutrient_units"
}

// NutrientType - model untuk tabel nutrient_types
type NutrientType struct {
	ID           int          `gorm:"primaryKey;autoIncrement" json:"id"`
	Code         string       `gorm:"type:varchar(30);uniqueIndex;not null" json:"code"`
	Name         string       `gorm:"type:varchar(100);not null" json:"name"`
	UnitID       int          `json:"unit_id"`
	Unit         NutrientUnit `gorm:"foreignKey:UnitID" json:"unit"`
	DisplayOrder int          `gorm:"default:0" json:"display_order"`
	IsActive     bool         `gorm:"default:true" json:"is_active"`
}

func (NutrientType) TableName() string {
	return "nutrient_types"
}

// FoodNutrient - model untuk tabel food_nutrients
type FoodNutrient struct {
	FoodID         string       `gorm:"type:char(36);primaryKey" json:"food_id"`
	NutrientTypeID int          `gorm:"primaryKey" json:"nutrient_type_id"`
	NutrientType   NutrientType `gorm:"foreignKey:NutrientTypeID" json:"nutrient_type"`
	ValuePer100g   float64      `gorm:"type:decimal(10,4);not null" json:"value_per_100g"`
}

func (FoodNutrient) TableName() string {
	return "food_nutrients"
}

// AssociatedFood - model untuk tabel associated_foods
type AssociatedFood struct {
	ID               int       `gorm:"primaryKey;autoIncrement" json:"id"`
	FoodID           string    `gorm:"type:char(36);not null;index" json:"food_id"`
	AssociatedFoodID string    `gorm:"type:char(36);not null" json:"associated_food_id"`
	AssociatedFood   Food      `gorm:"foreignKey:AssociatedFoodID" json:"associated_food"`
	Priority         int       `gorm:"default:0" json:"priority"`
	IsDefault        bool      `gorm:"default:false" json:"is_default"`
	CreatedAt        time.Time `json:"created_at"`
}

func (AssociatedFood) TableName() string {
	return "associated_foods"
}

// PortionSizeMethod - model untuk tabel food_portion_size_methods
type PortionSizeMethod struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	FoodID       string    `gorm:"type:char(36);not null;index" json:"food_id"`
	MethodType   string    `gorm:"type:enum('as_served','guide_image','weight');not null" json:"method_type"`
	Label        string    `gorm:"type:varchar(255);not null" json:"label"`
	Description  string    `gorm:"type:varchar(255)" json:"description"`
	ImageURL     string    `gorm:"type:varchar(500)" json:"image_url"`
	ThumbnailURL string    `gorm:"type:varchar(500)" json:"thumbnail_url"`
	Config       string    `gorm:"type:json" json:"config"` // JSON string
	DisplayOrder int       `gorm:"default:0" json:"display_order"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

func (PortionSizeMethod) TableName() string {
	return "food_portion_size_methods"
}

// AsServedSet - model untuk tabel as_served_sets
type AsServedSet struct {
	ID          string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	FoodID      *string   `gorm:"type:char(36)" json:"food_id"`
	Category    string    `gorm:"type:varchar(50)" json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

func (AsServedSet) TableName() string {
	return "as_served_sets"
}

// AsServedImage - model untuk tabel as_served_images
type AsServedImage struct {
	ID           string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	SetID        string    `gorm:"type:char(36);not null;index" json:"set_id"`
	Label        string    `gorm:"type:varchar(50);not null" json:"label"`
	ImageURL     string    `gorm:"type:varchar(500);not null" json:"image_url"`
	ThumbnailURL string    `gorm:"type:varchar(500)" json:"thumbnail_url"`
	WeightGram   float64   `gorm:"type:decimal(10,2);not null" json:"weight_gram"`
	Description  string    `gorm:"type:varchar(255)" json:"description"`
	DisplayOrder int       `gorm:"default:0" json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
}

func (AsServedImage) TableName() string {
	return "as_served_images"
}
