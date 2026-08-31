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
	// Buat database otomatis jika belum ada (berguna saat first deploy di Railway/Cloud)
	dsnRoot := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
	)
	if dbRoot, errRoot := gorm.Open(mysql.Open(dsnRoot), &gorm.Config{}); errRoot == nil {
		_ = dbRoot.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", cfg.DBName)).Error
		if sqlRoot, errSql := dbRoot.DB(); errSql == nil {
			_ = sqlRoot.Close()
		}
	}

	// Format DSN (Data Source Name) untuk MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
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
