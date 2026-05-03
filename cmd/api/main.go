package main

import (
	"atlas_food/internal/config"
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

	// Setup router Gin dengan middleware
	r := router.Setup(db)

	// Jalankan server pada port yang dikonfigurasi
	addr := ":" + cfg.ServerPort
	log.Printf("Server berjalan pada http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Gagal start server: %v", err)
	}
}
