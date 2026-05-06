package main

import (
	"atlas_food/internal/bootstrap"
	"atlas_food/internal/config"
	"atlas_food/internal/domain/auth"
	"atlas_food/internal/domain/survey"
	"atlas_food/internal/router"
	"log"
)

// main - entry point aplikasi Atlas Food API
// Inisialisasi config, database, dan start HTTP server
func main() {
	// Load konfigurasi dari .env file
	cfg := config.Load()

	// Koneksi ke database MySQL
	db := config.InitDB(cfg)

	// Auto Migration - buat tabel otomatis
	log.Println("Menjalankan Auto Migration...")
	err := db.AutoMigrate(
		&auth.User{},
		&auth.RefreshToken{},
		&survey.Locale{},
		&survey.Survey{},
		&survey.SurveyParticipant{},
	)
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	// Seed data awal untuk development
	if err := bootstrap.SeedInitialData(db, cfg); err != nil {
		log.Fatalf("Gagal seed data awal: %v", err)
	}

	// Setup router Gin dengan middleware
	r := router.Setup(db)

	// Jalankan server pada port yang dikonfigurasi
	addr := ":" + cfg.ServerPort
	log.Printf("Server berjalan pada http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Gagal start server: %v", err)
	}
}
