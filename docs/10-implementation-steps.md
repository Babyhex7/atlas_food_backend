# 🚀 Tahapan Implementasi Fitur

Dokumen ini menjelaskan langkah-langkah implementasi fitur-fitur Atlas Food API secara bertahap.

---

## ✅ Phase 1: Foundation (DONE)

### 1.1 Project Setup

- [x] Inisialisasi Go module (`go.mod`)
- [x] Setup folder structure (Clean Architecture)
- [x] Install dependencies (Gin, GORM, JWT, bcrypt, godotenv)
- [x] Buat `.env.example` dan `.gitignore`

### 1.2 Configuration & Database

- [x] Buat `config.go` - load env variables
- [x] Buat `database.go` - koneksi MySQL dengan GORM
- [x] Setup connection pool (MaxIdleConns, MaxOpenConns)

### 1.3 Middleware Dasar

- [x] Logger middleware - log request
- [x] CORS middleware - handle cross-origin
- [x] Error Handler middleware - global error handling
- [x] JWT Auth middleware - validasi token

---

## ✅ Phase 2: Authentication (DONE)

### 2.1 Model & Repository Layer

- [x] Buat model `User` dan `RefreshToken`
- [x] Definisi interface `Repository`
- [x] Implementasi `authRepository` dengan GORM

### 2.2 Service Layer

- [x] Interface `Service` - business logic contract
- [x] Implementasi `authService`:
  - [x] `Register()` - hash password, create user
  - [x] `Login()` - validasi credential, generate token
  - [x] `RefreshToken()` - rotate access token
  - [x] `GetProfile()` - ambil data user

### 2.3 Handler & Router

- [x] Buat `Handler` dengan dependency injection
- [x] Endpoint `POST /auth/register`
- [x] Endpoint `POST /auth/login`
- [x] Endpoint `POST /auth/refresh`
- [x] Endpoint `GET /auth/me` (protected)

### 2.4 Utils

- [x] `HashPassword()` - bcrypt
- [x] `CheckPassword()` - compare password
- [x] `GenerateJWT()` - buat token
- [x] `ValidateJWT()` - validasi token
- [x] `HashSHA256()` - hash refresh token
- [x] Response helpers (Success, Error, Validation)

### 2.5 Database Migration

- [x] Tabel `users`
- [x] Tabel `refresh_tokens`
- [x] Auto-migrate di `main.go`

**Status: ✅ SELESAI**

---

## ✅ Phase 3: Survey Management (DONE)

### 3.1 Model & Migration

- [x] Model `Survey`
- [x] Model `SurveyParticipant`
- [x] Model `Locale`
- [x] Auto-migrate di `main.go`

### 3.2 Repository Layer

- [x] Interface `Repository` - semua operasi DB
- [x] CRUD Survey operations
- [x] Participant operations
- [x] Locale operations

### 3.3 Service Layer

- [x] `CreateSurvey()` - buat survey baru
- [x] `GetSurveyByID()` - ambil detail survey
- [x] `ListSurveys()` - list dengan pagination
- [x] `UpdateSurvey()` - update data survey
- [x] `DeleteSurvey()` - hapus survey
- [x] `CloneSurvey()` - duplikat survey
- [x] `GenerateAccessToken()` - regenerate token
- [x] `GetPublicSurveyByToken()` - public access
- [x] `JoinSurvey()` - respondent join

### 3.4 Handler & Endpoints

- [x] `POST /api/v1/admin/surveys` - create survey (admin)
- [x] `GET /api/v1/admin/surveys` - list surveys (admin)
- [x] `GET /api/v1/admin/surveys/:id` - detail survey (admin)
- [x] `PUT /api/v1/admin/surveys/:id` - update survey (admin)
- [x] `DELETE /api/v1/admin/surveys/:id` - delete survey (admin)
- [x] `POST /api/v1/admin/surveys/:id/clone` - clone survey (admin)
- [x] `POST /api/v1/admin/surveys/:id/regenerate-token` - new token (admin)
- [x] `GET /api/v1/s/:token` - public survey access
- [x] `POST /api/v1/s/:token/join` - join survey
- [x] `GET /api/v1/locales` - list locales

**Status: ✅ SELESAI**

---

## 🔄 Phase 4: Food Database (PENDING)

### 4.1 Model & Migration

- [x] Model `Category`
- [x] Model `Food`
- [x] Model `FoodCategory` (junction)
- [x] Model `NutrientUnit`
- [x] Model `NutrientType`
- [x] Model `FoodNutrient`
- [x] Migration `003_create_foods.sql`

### 4.2 Admin Endpoints

- [x] CRUD Categories
- [x] CRUD Foods dengan search
- [x] Bulk import foods (CSV/JSON)
- [x] Manage food nutrients

### 4.3 Public Endpoints

- [x] `GET /api/v1/foods` - list dengan pagination & search
- [x] `GET /api/v1/foods/:id` - detail food
- [x] `GET /api/v1/categories` - list categories
- [x] `GET /api/v1/categories/:id/foods` - foods by category

**Status: ✅ SELESAI**

---

## 🔄 Phase 5: Portion Size (PENDING)

### 5.1 Model & Migration

- [x] Model `AssociatedFood`
- [x] Model `FoodPortionSizeMethod`
- [x] Model `AsServedSet`
- [x] Model `AsServedImage`
- [x] Migration `004_create_portion.sql`

### 5.2 Admin Endpoints

- [x] Manage portion methods per food
- [x] Upload as-served images
- [x] Create as-served sets
- [x] Manage associated foods

### 5.3 Public Endpoints

- [x] `GET /api/v1/foods/:id/portion-methods` - list methods
- [x] `GET /api/v1/as-served-sets/:id/images` - list images

**Status: ✅ SELESAI**

---

## 🔄 Phase 6: Survey Submission (PENDING)

### 6.1 Model & Migration

- [x] Model `SurveySubmission`
- [x] Migration `005_create_submissions.sql`

### 6.2 Respondent Flow

- [x] `POST /api/v1/surveys/:accessToken/join` - join survey
- [x] `GET /api/v1/surveys/:accessToken` - get survey detail
- [x] `POST /api/v1/submissions` - submit survey

### 6.3 Admin Endpoints

- [x] `GET /api/v1/admin/surveys/:id/submissions` - list submissions
- [x] `GET /api/v1/admin/submissions/:id` - detail submission
- [x] `GET /api/v1/admin/submissions/export` - export CSV/Excel

### 6.4 Nutrition Calculation

- [x] Service hitung nutrisi per meal
- [x] Agregasi nutrisi per hari
- [x] Store calculated results

**Status: ✅ SELESAI**

---

## 🔄 Phase 7: File Upload & Storage (PENDING)

### 7.1 Setup

- [x] Config upload path
- [x] Validasi file type (image only)
- [x] Validasi file size (max 10MB)

### 7.2 Endpoints

- [x] `POST /api/v1/upload` - upload image
- [x] Serve static files `/uploads/*`

### 7.3 As-Served Images

- [x] Upload ke `/uploads/as-served/`
- [x] Generate thumbnail
- [x] Store path di database

**Status: ✅ SELESAI**

---

## 📋 Ringkasan Progress

| Phase | Fitur         | Status      | Progress |
| ----- | ------------- | ----------- | -------- |
| 1     | Foundation    | ✅ Done     | 100%     |
| 2     | Auth          | ✅ Done     | 100%     |
| 3     | Survey        | ✅ Done     | 100%     |
| 4     | Food DB       | ✅ Done     | 100%     |
| 5     | Portion       | ✅ Done     | 100%     |
| 6     | Submission    | ✅ Done     | 100%     |
| 7     | Upload        | ✅ Done     | 100%     |
| **8** | **AI (Groq)** | ⏳ Pending  | 0%       |

---

## 🎯 Prioritas Selanjutnya

1. **AI Integration (Phase 8)** — Groq recommendation & nutritional insight
2. Lihat detail di `docs/brif_ai.md` untuk panduan lengkap implementasi AI.

---

## 🤖 Phase 8: AI Integration (PENDING)

> Referensi lengkap: [`docs/brif_ai.md`](./brif_ai.md)
>
> **Catatan penting:** AI **tidak punya endpoint tersendiri**. Groq dipanggil dari dalam `submission/service.go` setelah submission berhasil disimpan, dan hasilnya dikirim dalam response `POST /submissions` sebagai field `ai_result`.

### 8.1 Setup & Konfigurasi

- [ ] Tambahkan env vars Groq ke `.env` dan `.env.example`
  ```bash
  GROQ_API_KEY=gsk_xxxxxxxxxxxxxxxx
  GROQ_MODEL=llama3-8b-8192
  GROQ_BASE_URL=https://api.groq.com/openai/v1
  GROQ_TIMEOUT_SECONDS=15
  GROQ_MAX_TOKENS=512
  ```
- [ ] Tambahkan parsing config di `internal/config/config.go`

### 8.2 Groq Client Package

- [ ] Buat package baru: `internal/pkg/groq/client.go`
  - [ ] Struct `Client` dengan `apiKey`, `model`, `httpClient`
  - [ ] Method `AnalyzeNutrition(daily DailyTotal, foodNames []string) (*AIResult, error)`
  - [ ] Logic build system prompt + user prompt
  - [ ] HTTP POST ke Groq API (`/chat/completions`)
  - [ ] Parse & validasi JSON response

### 8.3 Database Migration

- [ ] Buat `migrations/008_create_ai_result_logs.sql` — tabel `ai_result_logs`
- [ ] Relasi: `1:1` dengan `survey_submissions`
- [ ] Tambahkan `AutoMigrate` di `main.go`

### 8.4 Modifikasi Domain Submission

- [ ] Tambahkan struct baru di `submission/dto.go`:
  - [ ] `AIResult`
  - [ ] `NutritionalAnalysis`
  - [ ] `HealthInsight`
  - [ ] Update `SubmitResponse` dengan field `AIResult`
- [ ] Modifikasi `submission/service.go`:
  - [ ] Inject `groq.Client` ke service struct
  - [ ] Setelah `s.repo.Create()`, panggil `s.groqClient.AnalyzeNutrition()`
  - [ ] Simpan hasil ke `ai_result_logs` via repo baru
  - [ ] Pastikan error Groq **tidak** membatalkan submission (log warning saja)
- [ ] Modifikasi `submission/handler.go`:
  - [ ] Sertakan `ai_result` di response JSON

### 8.5 Testing

- [ ] Test Postman: `POST /submissions` → cek ada field `ai_result` di response
- [ ] Cek isi `ai_result` mengandung: `overall_status`, `nutritional_analysis`, `recommended_foods`, `health_insight`, `suggested_activities`
- [ ] Test graceful degradation: set `GROQ_API_KEY=invalid` → submission tetap `201`, `ai_result: null`
- [ ] Verifikasi tabel `ai_result_logs` terisi setelah submission

**Status: ⏳ PENDING**

---

## 📝 Notes

- Setiap phase menggunakan Clean Architecture (Handler → Service → Repository)
- Semua endpoint admin dilindungi middleware `JWTAuth` + `AdminOnly`
- Semua endpoint public/respondent menggunakan `accessToken` survey
- Gunakan transaction untuk operasi yang melibatkan multiple tables
- Validasi input menggunakan `go-playground/validator`
- **AI Phase:** Groq API key diperoleh gratis di [console.groq.com](https://console.groq.com). AI hanya dipanggil saat submit, tidak ada endpoint AI terpisah.


