package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB - inisialisasi koneksi ke database MySQL
// Mengembalikan pointer ke gorm.DB untuk digunakan di seluruh aplikasi
func InitDB(cfg *Config) *gorm.DB {
	// Format DSN (Data Source Name) untuk MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	// Konfigurasi GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log query SQL
	}

	// Buka koneksi database
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}

	// Test koneksi dengan mengambil SQL DB instance
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal ambil SQL DB: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("Database berhasil terkoneksi")
	return db
}
