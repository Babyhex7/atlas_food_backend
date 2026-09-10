package config

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// Nilai default yang HANYA boleh dipakai di luar produksi. validate() menolak
// boot kalau nilai-nilai ini masih terpasang saat SERVER_MODE=release.
const (
	defaultJWTSecret     = "default-secret-key-minimal-32-chars"
	defaultAdminPassword = "Password123!"
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

	// explicitProduction - true kalau SERVER_MODE/GIN_MODE benar-benar di-set
	// ke release/production lewat environment (bukan sekadar nilai default).
	explicitProduction bool
}

var (
	loadOnce sync.Once
	loaded   *Config
)

// Load - membaca konfigurasi dari environment variables.
//
// Hasilnya di-cache: Load() dipanggil di jalur panas (setiap generate/validate
// JWT, artinya setiap request ter-autentikasi). Tanpa sync.Once setiap request
// membaca file .env dari disk dan menulis satu baris log — pemborosan I/O yang
// tidak terlihat sampai trafik naik.
func Load() *Config {
	loadOnce.Do(func() {
		loaded = load()
	})
	return loaded
}

// load - pembacaan sebenarnya, dijalankan tepat satu kali
func load() *Config {
	// Load .env file (ignore error kalau file tidak ada)
	_ = godotenv.Load()

	rawURL := getEnv("MYSQL_URL", getEnv("DATABASE_URL", ""))
	urlHost, urlPort, urlUser, urlPass, urlDB := parseDBURL(rawURL)

	cfg := &Config{
		// Database config (Mendukung MYSQL_URL, DATABASE_URL, maupun variabel terpisah Railway MySQL)
		DBHost:     getEnv("DB_HOST", getEnv("MYSQLHOST", firstNonEmpty(urlHost, "localhost"))),
		DBPort:     getEnv("DB_PORT", getEnv("MYSQLPORT", firstNonEmpty(urlPort, "3306"))),
		DBUser:     getEnv("DB_USER", getEnv("MYSQLUSER", firstNonEmpty(urlUser, "root"))),
		DBPassword: getEnv("DB_PASSWORD", getEnv("MYSQLPASSWORD", getEnv("MYSQL_ROOT_PASSWORD", urlPass))),
		DBName:     getEnv("DB_NAME", getEnv("MYSQLDATABASE", getEnv("MYSQL_DATABASE", firstNonEmpty(urlDB, "db_atlas_food")))),

		// JWT config
		JWTSecret:              getEnv("JWT_SECRET", defaultJWTSecret),
		JWTExpiration:          parseDuration(getEnv("JWT_EXPIRATION", "24h")),
		RefreshTokenExpiration: parseDuration(getEnv("REFRESH_TOKEN_EXPIRATION", "168h")),

		// Server config (Railway/Render secara otomatis menyuntikkan PORT)
		ServerPort: getEnv("PORT", getEnv("SERVER_PORT", "8080")),
		ServerMode: getEnv("SERVER_MODE", "release"),

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
		AdminSeedPassword: getEnv("ADMIN_SEED_PASSWORD", defaultAdminPassword),
		AdminSeedName:     getEnv("ADMIN_SEED_NAME", "Admin Atlas"),
	}

	cfg.explicitProduction = isProductionMode(os.Getenv("SERVER_MODE")) || isProductionMode(os.Getenv("GIN_MODE"))
	cfg.validate()
	cfg.logOrigins()

	log.Println("Konfigurasi berhasil dimuat")
	return cfg
}

// IsProduction - true kalau server dijalankan di mode release
func (c *Config) IsProduction() bool {
	return isProductionMode(c.ServerMode)
}

// logOrigins - cetak allowlist yang berlaku saat boot.
//
// CORS dan handshake WebSocket sekarang memakai allowlist. Kalau origin
// frontend tidak terdaftar, browser memblokir SEMUA request tanpa pesan yang
// berguna di sisi server — gejalanya terlihat seperti "backend mati", padahal
// hanya FRONTEND_URL yang belum di-set. Mencetaknya di log boot membuat
// penyebabnya terlihat dalam hitungan detik, bukan jam.
func (c *Config) logOrigins() {
	origins := c.AllowedOrigins()
	log.Printf("CORS/WebSocket mengizinkan origin: %v", origins)

	if !c.IsProduction() {
		return
	}
	if os.Getenv("FRONTEND_URL") == "" && os.Getenv("CORS_ALLOWED_ORIGINS") == "" {
		log.Printf("PERINGATAN: FRONTEND_URL dan CORS_ALLOWED_ORIGINS keduanya kosong di mode produksi. " +
			"Hanya %v yang diizinkan, jadi frontend yang ter-deploy AKAN diblokir browser. " +
			"Set FRONTEND_URL ke domain frontend Anda.", origins)
	}
}

// isProductionMode - kenali penamaan mode produksi yang lazim
func isProductionMode(mode string) bool {
	return strings.EqualFold(mode, "release") || strings.EqualFold(mode, "production")
}

// AllowedOrigins - daftar origin yang boleh mengakses API (CORS & WebSocket).
//
// Diambil dari CORS_ALLOWED_ORIGINS (dipisah koma) dan selalu menyertakan
// FRONTEND_URL. Di mode non-produksi localhost ikut diizinkan supaya dev
// tidak perlu mengonfigurasi apa pun.
func (c *Config) AllowedOrigins() []string {
	origins := make([]string, 0, 4)
	seen := map[string]bool{}
	add := func(o string) {
		o = strings.TrimRight(strings.TrimSpace(o), "/")
		if o == "" || seen[o] {
			return
		}
		seen[o] = true
		origins = append(origins, o)
	}

	add(c.FrontendURL)
	for _, o := range strings.Split(getEnv("CORS_ALLOWED_ORIGINS", ""), ",") {
		add(o)
	}
	if !c.IsProduction() {
		add("http://localhost:3000")
		add("http://127.0.0.1:3000")
	}
	return origins
}

// validate - hentikan boot kalau ada konfigurasi berbahaya di produksi.
//
// Rahasia default yang lolos ke produksi berarti siapa pun bisa menandatangani
// token admin sendiri. Lebih baik gagal saat start daripada diam-diam tidak aman.
//
// Pemeriksaan keras HANYA berlaku kalau produksi dinyatakan eksplisit lewat
// SERVER_MODE/GIN_MODE. Default ServerMode memang "release", tapi memakai itu
// sebagai pemicu log.Fatal akan membuat `go run` di laptop tanpa .env langsung
// mati — bukan itu yang kita mau.
func (c *Config) validate() {
	insecureSecret := c.JWTSecret == defaultJWTSecret || len(c.JWTSecret) < 32
	insecurePassword := c.AdminSeedPassword == defaultAdminPassword

	if !c.explicitProduction {
		if insecureSecret {
			log.Println("PERINGATAN: JWT_SECRET lemah/default — WAJIB diganti sebelum deploy")
		}
		if insecurePassword {
			log.Println("PERINGATAN: ADMIN_SEED_PASSWORD masih default — WAJIB diganti sebelum deploy")
		}
		return
	}

	if insecureSecret {
		log.Fatal("JWT_SECRET wajib di-set minimal 32 karakter saat SERVER_MODE=release")
	}
	if insecurePassword {
		log.Fatal("ADMIN_SEED_PASSWORD wajib diganti saat SERVER_MODE=release")
	}
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

// firstNonEmpty - mengembalikan argumen pertama yang tidak kosong
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// parseDBURL - mengekstrak host, port, user, password, dbname dari format URL database (misal MYSQL_URL)
func parseDBURL(rawURL string) (host, port, user, password, dbname string) {
	if rawURL == "" {
		return
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	host = u.Hostname()
	port = u.Port()
	if port == "" {
		port = "3306"
	}
	if u.User != nil {
		user = u.User.Username()
		password, _ = u.User.Password()
	}
	dbname = strings.TrimPrefix(u.Path, "/")
	return
}
