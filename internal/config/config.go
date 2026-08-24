package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config - struct untuk menyimpan semua konfigurasi aplikasi
type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Groq AI
	GroqAPIKey      string
	GroqModel       string
	GroqBaseURL     string
	GroqTimeoutSecs int
	GroqMaxTokens   int

	// JWT
	JWTSecret              string
	JWTExpiration          time.Duration
	RefreshTokenExpiration time.Duration

	// Server
	ServerPort string
	ServerMode string

	// Frontend (dipakai untuk generate link survey yang dibagikan ke responden)
	FrontendURL string

	// Upload
	UploadPath    string
	MaxUploadSize int64

	// Seed data
	AdminSeedEmail    string
	AdminSeedPassword string
	AdminSeedName     string
}

// Load - membaca konfigurasi dari environment variables
// Mengembalikan pointer ke Config yang sudah terisi
func Load() *Config {
	// Load .env file (ignore error kalau file tidak ada)
	_ = godotenv.Load()

	cfg := &Config{
		// Database config
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "atlas_food"),

		// JWT config
		JWTSecret:              getEnv("JWT_SECRET", "default-secret-key-minimal-32-chars"),
		JWTExpiration:          parseDuration(getEnv("JWT_EXPIRATION", "24h")),
		RefreshTokenExpiration: parseDuration(getEnv("REFRESH_TOKEN_EXPIRATION", "168h")),

		// Server config
		ServerPort: getEnv("SERVER_PORT", "8080"),
		ServerMode: getEnv("SERVER_MODE", "debug"),

		// Frontend config
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),

		// Groq AI config
		GroqAPIKey:      getEnv("GROQ_API_KEY", ""),
		GroqModel:       getEnv("GROQ_MODEL", "llama-3.3-70b-versatile"),
		GroqBaseURL:     getEnv("GROQ_BASE_URL", "https://api.groq.com/openai/v1"),
		GroqTimeoutSecs: parseInt(getEnv("GROQ_TIMEOUT_SECONDS", "45")),
		GroqMaxTokens:   parseInt(getEnv("GROQ_MAX_TOKENS", "2048")),

		// Upload config
		UploadPath:    getEnv("UPLOAD_PATH", "./uploads"),
		MaxUploadSize: parseInt64(getEnv("MAX_UPLOAD_SIZE", "10485760")),

		// Seed data config
		AdminSeedEmail:    getEnv("ADMIN_SEED_EMAIL", "admin@mail.com"),
		AdminSeedPassword: getEnv("ADMIN_SEED_PASSWORD", "Password123!"),
		AdminSeedName:     getEnv("ADMIN_SEED_NAME", "Admin Atlas"),
	}

	log.Println("Konfigurasi berhasil dimuat")
	return cfg
}

// getEnv - ambil value dari env variable, return default kalau kosong
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDuration - parse string ke time.Duration
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour // default 24 jam
	}
	return d
}

// parseInt64 - parse string ke int64
func parseInt64(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 10485760 // default 10MB
	}
	return n
}

// parseInt - parse string ke int
func parseInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
