package bootstrap

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gorm.io/gorm"
)

// AtlasFoodItem represents the structure from Atlas_Makananku_FINAL.json
type AtlasFoodItem struct {
	Kode      string `json:"kode"`
	Kategori  string `json:"kategori"`
	NamaID    string `json:"nama_id"`
	NamaEN    string `json:"nama_en"`
	TipeFoto  string `json:"tipe_foto"`
	Keterangan *string `json:"keterangan"`
	Porsi     []struct {
		LabelUkuran string  `json:"label_ukuran"`
		Nilai       float64 `json:"nilai"`
		Satuan      string  `json:"satuan"`
		URT         *string `json:"urt"`
		PerUnit     bool    `json:"per_unit"`
	} `json:"porsi"`
}

// CategoryMapping maps Indonesian category names to category codes
var CategoryMapping = map[string]string{
	"Aneka Buah":           "AB",
	"Lauk Hewani":          "LH",
	"Lauk Nabati":          "LN",
	"Aneka Sayur":          "AS",
	"Makanan Pokok":        "MP",
	"Aneka Produk":         "AP",
	"Aneka Minuman Keras":  "AMK",
	"Kue Kering":           "KK",
	"Aneka Bumbu Kering":   "ABK",
	"Aneka Kue":            "AK",
	"Makanan Daerah Lokal": "MDL",
	"Gula & Kembang":       "GK",
	"Aneka Hasil":          "AH",
}

// SeedFindYourFoodData seeds the database with data from Atlas_Makananku_FINAL.json
func SeedFindYourFoodData(db *gorm.DB, jsonFilePath string) error {
	// Read JSON file
	file, err := os.Open(jsonFilePath)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %v", err)
	}

	var atlasData []AtlasFoodItem
	if err := json.Unmarshal(bytes, &atlasData); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Start transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Seed Categories
	if err := seedCategories(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed categories: %v", err)
	}

	// 2. Seed Foods and Portion Photos
	if err := seedFoodsFromAtlas(tx, atlasData); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed foods: %v", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	fmt.Printf("✅ Successfully seeded %d foods from Atlas Makananku\n", len(atlasData))
	return nil
}

// seedCategories - isi 13 kategori Atlas Makananku (MP, LH, LN, ...); kategori yang sudah ada dilewati
func seedCategories(tx *gorm.DB) error {
	categories := []map[string]interface{}{
		{"code": "MP", "name": "Makanan Pokok", "icon": "🍚", "display_order": 1},
		{"code": "LH", "name": "Lauk Hewani", "icon": "🍗", "display_order": 2},
		{"code": "LN", "name": "Lauk Nabati", "icon": "🥜", "display_order": 3},
		{"code": "AS", "name": "Aneka Sayur", "icon": "🥗", "display_order": 4},
		{"code": "AB", "name": "Aneka Buah", "icon": "🍎", "display_order": 5},
		{"code": "AP", "name": "Aneka Produk", "icon": "🥛", "display_order": 6},
		{"code": "AMK", "name": "Aneka Minuman Keras", "icon": "🥤", "display_order": 7},
		{"code": "KK", "name": "Kue Kering", "icon": "🍪", "display_order": 8},
		{"code": "ABK", "name": "Aneka Bumbu Kering", "icon": "🧂", "display_order": 9},
		{"code": "AK", "name": "Aneka Kue", "icon": "🧁", "display_order": 10},
		{"code": "MDL", "name": "Makanan Daerah Lokal", "icon": "🍲", "display_order": 11},
		{"code": "GK", "name": "Gula & Kembang", "icon": "🍬", "display_order": 12},
		{"code": "AH", "name": "Aneka Hasil", "icon": "🥫", "display_order": 13},
	}

	for _, cat := range categories {
		// Check if exists
		var count int64
		tx.Table("categories").Where("code = ?", cat["code"]).Count(&count)
		if count == 0 {
			if err := tx.Table("categories").Create(cat).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// seedFoodsFromAtlas - isi tabel foods dari data JSON Atlas Makananku, dipetakan ke kategori berdasarkan kode
func seedFoodsFromAtlas(tx *gorm.DB, atlasData []AtlasFoodItem) error {
	// Get category IDs
	categoryMap := make(map[string]string) // category_code -> category_id
	var categories []struct {
		ID   string
		Code string
	}
	if err := tx.Table("categories").Select("id, code").Find(&categories).Error; err != nil {
		return err
	}
	for _, cat := range categories {
		categoryMap[cat.Code] = cat.ID
	}

	// Process foods sequentially (no goroutines for DB writes in transaction)
	// But optimize with batch inserts
	for i, item := range atlasData {
		if err := processSingleFood(tx, item, categoryMap); err != nil {
			return err
		}

		if (i+1)%50 == 0 {
			fmt.Printf("📦 Processed %d/%d foods...\n", i+1, len(atlasData))
		}
	}

	return nil
}

// processSingleFood handles single food item processing
func processSingleFood(tx *gorm.DB, item AtlasFoodItem, categoryMap map[string]string) error {
	// Get category code from category name
	categoryCode := CategoryMapping[item.Kategori]
	if categoryCode == "" {
		// Try to extract from kode (e.g., "AB-01" -> "AB")
		parts := strings.Split(item.Kode, "-")
		if len(parts) > 0 {
			categoryCode = parts[0]
		}
	}

	categoryID := categoryMap[categoryCode]

	// Map photo type
	photoType := "series"
	if strings.ToLower(item.TipeFoto) == "range" {
		photoType = "range"
	}

	// Check if food already exists
	var existingFood struct {
		ID string
	}
	err := tx.Table("foods").Select("id").Where("code = ?", item.Kode).First(&existingFood).Error

	var foodID string
	if err == gorm.ErrRecordNotFound {
		// Insert new food
		food := map[string]interface{}{
			"code":        item.Kode,
			"name":        item.NamaID,
			"local_name":  item.NamaEN,
			"description": item.Keterangan,
			"photo_type":  photoType,
			"category_id": categoryID,
			"is_active":   true,
		}

		if err := tx.Table("foods").Create(&food).Error; err != nil {
			return fmt.Errorf("failed to create food %s: %v", item.Kode, err)
		}

		// Get the inserted food ID
		if err := tx.Table("foods").Select("id").Where("code = ?", item.Kode).First(&existingFood).Error; err != nil {
			return fmt.Errorf("failed to get food ID for %s: %v", item.Kode, err)
		}
		foodID = existingFood.ID
	} else if err != nil {
		return fmt.Errorf("failed to check food existence: %v", err)
	} else {
		foodID = existingFood.ID
		// Update existing food
		updates := map[string]interface{}{
			"name":        item.NamaID,
			"local_name":  item.NamaEN,
			"description": item.Keterangan,
			"photo_type":  photoType,
			"category_id": categoryID,
		}
		if err := tx.Table("foods").Where("id = ?", foodID).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update food %s: %v", item.Kode, err)
		}
	}

	// Create as_served_set if not exists
	var asServedSetID string
	var existingSet struct {
		ID string
	}

	// Generate unique code for as_served_set
	setCode := fmt.Sprintf("SET-%s", item.Kode)

	err = tx.Table("as_served_sets").Select("id").Where("code = ?", setCode).First(&existingSet).Error
	if err == gorm.ErrRecordNotFound {
		asServedSet := map[string]interface{}{
			"code":        setCode,
			"name":        fmt.Sprintf("%s Portion Set", item.NamaID),
			"description": fmt.Sprintf("Portion sizes for %s from Atlas Makananku", item.NamaID),
			"food_id":     foodID,
			"category":    categoryCode,
		}
		if err := tx.Table("as_served_sets").Create(&asServedSet).Error; err != nil {
			return fmt.Errorf("failed to create as_served_set: %v", err)
		}
		// Get ID
		tx.Table("as_served_sets").Select("id").Where("code = ?", setCode).First(&existingSet)
		asServedSetID = existingSet.ID
	} else if err != nil {
		return fmt.Errorf("failed to check as_served_set: %v", err)
	} else {
		asServedSetID = existingSet.ID
	}

	// Delete existing as_served_images for this set
	tx.Table("as_served_images").Where("set_id = ?", asServedSetID).Delete(nil)

	// Tentukan apakah food ini bertipe guide
	isGuide := strings.ToLower(item.TipeFoto) == "guide"

	// Batch insert portion photos for better performance
	var portionPhotos []map[string]interface{}
	for j, porsi := range item.Porsi {
		// Skip if nilai is 0 or negative
		if porsi.Nilai <= 0 {
			continue
		}

		// Build description from URT if available
		description := fmt.Sprintf("Porsi %s - %.1f%s", porsi.LabelUkuran, porsi.Nilai, porsi.Satuan)
		if porsi.URT != nil && *porsi.URT != "" {
			description = fmt.Sprintf("%s (%s)", description, *porsi.URT)
		}

		// Bangun URL sesuai nama file asli di folder foto:
		//   Series/Range : /uploads/atlas/[CATCODE]/[KODE]_[LABEL].jpg  (misal: MP-01_A.jpg)
		//   Guide        : /uploads/atlas/[CATCODE]/[KODE]_guide.jpg     (misal: AS-22_guide.jpg)
		// URL ini sekaligus menjadi object key di MinIO: photos/[CATCODE]/[filename]
		imageURL := BuildImageURL(categoryCode, item.Kode, porsi.LabelUkuran, isGuide)

		portionPhotos = append(portionPhotos, map[string]interface{}{
			"set_id":        asServedSetID,
			"label":         porsi.LabelUkuran,
			"weight_gram":   porsi.Nilai,
			"image_url":     imageURL,
			"thumbnail_url": imageURL, // thumbnail = same file, frontend dapat resize via query param
			"description":   description,
			"display_order": j + 1,
		})
	}

	// Batch insert all portion photos at once
	if len(portionPhotos) > 0 {
		if err := tx.Table("as_served_images").Create(&portionPhotos).Error; err != nil {
			return fmt.Errorf("failed to batch create portion photos for %s: %v", item.Kode, err)
		}
	}

	return nil
}
