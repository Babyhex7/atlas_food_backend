package bootstrap

import (
	"time"

	"atlas_food/internal/config"
	"atlas_food/internal/domain/auth"
	"atlas_food/internal/domain/food"
	"atlas_food/internal/domain/survey"
	"atlas_food/internal/pkg/utils"

	"gorm.io/gorm"
)

// SeedInitialData - memasukkan data awal (admin, locales, dll.)
func SeedInitialData(db *gorm.DB, cfg *config.Config) error {
	// Seed admin user jika belum ada
	var cnt int64
	if err := db.Model(&auth.User{}).Where("email = ?", cfg.AdminSeedEmail).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt == 0 {
		hash, err := utils.HashPassword(cfg.AdminSeedPassword)
		if err != nil {
			return err
		}
		admin := &auth.User{
			Email:        cfg.AdminSeedEmail,
			PasswordHash: hash,
			Name:         cfg.AdminSeedName,
			Role:         "admin",
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := db.Create(admin).Error; err != nil {
			return err
		}
	}

	// Seed default locales jika tabel kosong
	var localeCnt int64
	if err := db.Model(&survey.Locale{}).Count(&localeCnt).Error; err != nil {
		return err
	}
	if localeCnt == 0 {
		locales := []survey.Locale{
			{Code: "id", Name: "Bahasa Indonesia"},
			{Code: "en", Name: "English"},
		}
		if err := db.Create(&locales).Error; err != nil {
			return err
		}
	}

	// Seed sample data foods jika kosong
	var foodCnt int64
	if err := db.Model(&food.Food{}).Count(&foodCnt).Error; err != nil {
		return err
	}
	if foodCnt == 0 {
		importSampleData(db)
	}
	// Seed default survey
	var surveyCnt int64
	if err := db.Model(&survey.Survey{}).Where("slug = ?", "test-survey").Count(&surveyCnt).Error; err != nil {
		return err
	}
	if surveyCnt == 0 {
		var admin auth.User
		db.Where("email = ?", cfg.AdminSeedEmail).First(&admin)

		testSurvey := &survey.Survey{
			Slug:        "test-survey",
			Name:        "Survey Konsumsi Harian (Test)",
			Description: "Survey untuk mencoba alur MVP.",
			MealsConfig: `{"meals":[{"name":"Sarapan","time":"07:00"},{"name":"Makan Siang","time":"12:00"},{"name":"Makan Malam","time":"19:00"}]}`,
			Prompts:     `{"before_meals":"Selamat datang di survey gizi","after_meals":"Terima kasih"}`,
			LocaleID:    1,
			Status:      "active",
			AccessToken: "test-survey-access-token-2025",
			CreatedBy:   admin.ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := db.Create(testSurvey).Error; err != nil {
			return err
		}

		// Create a participant
		participant := &survey.SurveyParticipant{
			SurveyID:  testSurvey.ID,
			UserID:    admin.ID,
			Alias:     "RESPONDENT-001",
			CreatedAt: time.Now(),
		}
		db.Create(participant)
	}

	return nil
}

func importSampleData(db *gorm.DB) {
	importSQL := `
INSERT INTO nutrient_units (id, code, name, symbol) VALUES
(1, 'kcal', 'Kilokalori', 'kcal'),
(2, 'g', 'Gram', 'g'),
(3, 'mg', 'Miligram', 'mg'),
(4, 'mcg', 'Mikrogram', 'μg');

INSERT INTO nutrient_types (id, code, name, unit_id, display_order) VALUES
(1, 'energy', 'Energi', 1, 1),
(2, 'protein', 'Protein', 2, 2),
(3, 'carbs', 'Karbohidrat', 2, 3),
(4, 'fat', 'Lemak Total', 2, 4),
(5, 'fiber', 'Serat', 2, 5),
(6, 'sugar', 'Gula', 2, 6),
(7, 'sodium', 'Natrium', 3, 7),
(8, 'calcium', 'Kalsium', 3, 8),
(9, 'iron', 'Zat Besi', 3, 9),
(10, 'vit_c', 'Vitamin C', 3, 10);

INSERT INTO categories (id, code, name, icon, display_order) VALUES
('uuid-cat-1', 'MP', 'Makanan Pokok', '🍚', 1),
('uuid-cat-2', 'LH', 'Lauk Hewani', '🍗', 2),
('uuid-cat-3', 'AB', 'Aneka Buah', '🍌', 3),
('uuid-cat-4', 'AMK', 'Makanan & Minuman Kemasan', '🥤', 4);

INSERT INTO foods (id, code, name, local_name, description, photo_type, category_id) VALUES
('uuid-nasi', 'MP-01', 'Nasi', 'Rice', 'Nasi putih matang', 'series', 'uuid-cat-1'),
('uuid-ayam', 'LH-12', 'Ayam Goreng Dada', 'Fried Chicken Breast', 'Ayam goreng tepung bagian dada', 'range', 'uuid-cat-2'),
('uuid-pisang', 'AB-01', 'Pisang', 'Banana', 'Pisang cavendish', 'series', 'uuid-cat-3'),
('uuid-susu', 'AMK-01', 'Susu UHT', 'UHT Milk', 'Susu UHT full cream', 'series', 'uuid-cat-4');

INSERT INTO food_nutrients (food_id, nutrient_type_id, value_per_100g) VALUES
('uuid-nasi', 1, 130.00),
('uuid-nasi', 2, 2.70),
('uuid-nasi', 3, 28.00),
('uuid-nasi', 4, 0.30),
('uuid-nasi', 5, 0.40),
('uuid-ayam', 1, 290.00),
('uuid-ayam', 2, 18.00),
('uuid-ayam', 4, 19.00),
('uuid-pisang', 1, 89.00),
('uuid-pisang', 2, 1.10),
('uuid-pisang', 3, 22.80),
('uuid-pisang', 6, 12.20),
('uuid-susu', 1, 64.00),
('uuid-susu', 2, 3.40),
('uuid-susu', 8, 120.00);

INSERT INTO as_served_sets (id, code, name, description, category, food_id) VALUES
('uuid-set-3', 'banana-slices', 'Sliced Banana Portions', 'Visual guide for banana portions', 'fruits', 'uuid-pisang');

INSERT INTO as_served_images (id, set_id, label, image_url, thumbnail_url, weight_gram, description, display_order) VALUES
(UUID(), 'uuid-set-3', '1', '/banana/banana-1.jpg', '/banana/banana-1-thumb.jpg', 20.0, 'Few slices (~20g)', 1),
(UUID(), 'uuid-set-3', '2', '/banana/banana-2.jpg', '/banana/banana-2-thumb.jpg', 40.0, 'Small portion (~40g)', 2),
(UUID(), 'uuid-set-3', '3', '/banana/banana-3.jpg', '/banana/banana-3-thumb.jpg', 60.0, 'Medium-small (~60g)', 3),
(UUID(), 'uuid-set-3', '4', '/banana/banana-4.jpg', '/banana/banana-4-thumb.jpg', 95.0, 'Medium (~95g)', 4),
(UUID(), 'uuid-set-3', '5', '/banana/banana-5.jpg', '/banana/banana-5-thumb.jpg', 130.0, 'Medium-large (~130g)', 5),
(UUID(), 'uuid-set-3', '6', '/banana/banana-6.jpg', '/banana/banana-6-thumb.jpg', 160.0, 'Large (~160g)', 6),
(UUID(), 'uuid-set-3', '7', '/banana/banana-7.jpg', '/banana/banana-7-thumb.jpg', 175.0, 'Very large (~175g)', 7),
(UUID(), 'uuid-set-3', '8', '/banana/banana-8.jpg', '/banana/banana-8-thumb.jpg', 190.0, 'Full plate (~190g)', 8);

INSERT INTO food_portion_size_methods (food_id, method_type, label, description, image_url, thumbnail_url, config, display_order) VALUES
('uuid-nasi', 'as_served', 'A', 'Porsi sangat kecil', '/uploads/nasi/nasi-A.jpg', '/uploads/nasi/nasi-A-thumb.jpg', '{"weight_gram": 50}', 1),
('uuid-nasi', 'as_served', 'B', 'Porsi kecil', '/uploads/nasi/nasi-B.jpg', '/uploads/nasi/nasi-B-thumb.jpg', '{"weight_gram": 90}', 2),
('uuid-nasi', 'as_served', 'C', 'Porsi sedang-kecil', '/uploads/nasi/nasi-C.jpg', '/uploads/nasi/nasi-C-thumb.jpg', '{"weight_gram": 130}', 3),
('uuid-nasi', 'as_served', 'D', 'Porsi sedang', '/uploads/nasi/nasi-D.jpg', '/uploads/nasi/nasi-D-thumb.jpg', '{"weight_gram": 150}', 4),
('uuid-nasi', 'as_served', 'E', 'Porsi sedang-besar', '/uploads/nasi/nasi-E.jpg', '/uploads/nasi/nasi-E-thumb.jpg', '{"weight_gram": 210}', 5),
('uuid-nasi', 'as_served', 'F', 'Porsi besar', '/uploads/nasi/nasi-F.jpg', '/uploads/nasi/nasi-F-thumb.jpg', '{"weight_gram": 250}', 6),
('uuid-nasi', 'as_served', 'G', 'Porsi sangat besar', '/uploads/nasi/nasi-G.jpg', '/uploads/nasi/nasi-G-thumb.jpg', '{"weight_gram": 270}', 7),
('uuid-nasi', 'as_served', 'H', 'Porsi ekstra besar', '/uploads/nasi/nasi-H.jpg', '/uploads/nasi/nasi-H-thumb.jpg', '{"weight_gram": 350}', 8),
('uuid-pisang', 'as_served', 'Sliced banana on plate', 'Choose portion and adjust quantity', '/banana/preview.jpg', NULL, '{"selectionType":"as_served_quantity","setCode":"banana-slices","thumbnailPosition":"bottom","maxQuantity":5,"allowFractions":true,"fractionOptions":[0,0.25,0.5,0.75],"defaultQuantity":1,"defaultFraction":0,"showCalculation":true}', 1);
`
	// Split by double newline to get individual statements or just execute as one block if driver allows.
	// GORM's Exec usually only does one statement at a time unless multiStatements=true. 
	// We'll execute them manually by splitting on ";"
	importSQL = "SET FOREIGN_KEY_CHECKS=0;" + importSQL + "SET FOREIGN_KEY_CHECKS=1;"
	db.Exec(importSQL) // Note: mysql driver needs parseTime=true&multiStatements=true in DSN to run multiple queries in one Exec.
}
