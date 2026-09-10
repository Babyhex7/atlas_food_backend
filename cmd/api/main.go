package main

import (
	"atlas_food/internal/bootstrap"
	"atlas_food/internal/config"
	"atlas_food/internal/domain/ai"
	"atlas_food/internal/domain/annotation"
	"atlas_food/internal/domain/auth"
	"atlas_food/internal/domain/collab"
	"atlas_food/internal/domain/food"
	"atlas_food/internal/domain/submission"
	"atlas_food/internal/domain/survey"
	"atlas_food/internal/router"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
		&food.Category{},
		&food.Food{},
		&food.NutrientUnit{},
		&food.NutrientType{},
		&food.FoodNutrient{},
		&food.AssociatedFood{},
		&food.PortionSizeMethod{},
		&food.AsServedSet{},
		&food.AsServedImage{},
		&submission.SurveySubmission{},
		&ai.AIResultLog{},
		&annotation.FoodImage{},
		&annotation.FoodArea{},
	)
	if err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}

	log.Println("Memastikan FULLTEXT index untuk tabel foods ada...")
	// Index dibutuhkan agar MATCH(name, local_name) AGAINST(...) berfungsi
	errFT := db.Exec("CREATE FULLTEXT INDEX ft_name_local ON foods(name, local_name)").Error
	if errFT != nil {
		log.Printf("Info FULLTEXT index (diabaikan jika sudah ada): %v", errFT)
	}

	// Seed data awal untuk development
	if err := bootstrap.SeedInitialData(db, cfg); err != nil {
		log.Fatalf("Gagal seed data awal: %v", err)
	}

	// Initialize WebSocket Hub for real-time collaboration
	log.Println("Initializing WebSocket Hub...")
	hub := collab.NewHub()
	go hub.Run() // Start hub in goroutine

	// Setup router Gin dengan middleware
	r := router.Setup(db, cfg, hub)

	// Jalankan server pada port yang dikonfigurasi
	addr := ":" + cfg.ServerPort
	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		log.Printf("Server berjalan pada http://localhost%s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Gagal start server: %v", err)
		}
	}()

	// Graceful shutdown: tanpa ini hub tidak pernah dihentikan dan koneksi
	// WebSocket yang sedang aktif diputus mendadak saat proses dimatikan.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Sinyal shutdown diterima, menutup server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server dipaksa berhenti: %v", err)
	}
	hub.Stop()

	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	log.Println("Server berhenti dengan bersih")
}
