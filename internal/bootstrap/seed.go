package bootstrap

import (
	"time"

	"atlas_food/internal/config"
	"atlas_food/internal/domain/auth"
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

	return nil
}
