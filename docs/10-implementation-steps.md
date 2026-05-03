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

## 🔄 Phase 3: Survey Management (NEXT)

### 3.1 Model & Migration
- [ ] Model `Survey`
- [ ] Model `SurveyParticipant`
- [ ] Migration `002_create_surveys.sql`
- [ ] Seeder data `locales`

### 3.2 Admin Endpoints (Protected)
- [ ] `GET /api/v1/admin/surveys` - list surveys
- [ ] `POST /api/v1/admin/surveys` - create survey
- [ ] `GET /api/v1/admin/surveys/:id` - detail survey
- [ ] `PUT /api/v1/admin/surveys/:id` - update survey
- [ ] `DELETE /api/v1/admin/surveys/:id` - delete survey
- [ ] `POST /api/v1/admin/surveys/:id/clone` - clone survey

### 3.3 Survey Access Token
- [ ] Generate access token unik per survey
- [ ] Middleware validasi access token

**Status: 🚧 BELUM MULAI**

---

## 🔄 Phase 4: Food Database (PENDING)

### 4.1 Model & Migration
- [ ] Model `Category`
- [ ] Model `Food`
- [ ] Model `FoodCategory` (junction)
- [ ] Model `NutrientUnit`
- [ ] Model `NutrientType`
- [ ] Model `FoodNutrient`
- [ ] Migration `003_create_foods.sql`

### 4.2 Admin Endpoints
- [ ] CRUD Categories
- [ ] CRUD Foods dengan search
- [ ] Bulk import foods (CSV/JSON)
- [ ] Manage food nutrients

### 4.3 Public Endpoints
- [ ] `GET /api/v1/foods` - list dengan pagination & search
- [ ] `GET /api/v1/foods/:id` - detail food
- [ ] `GET /api/v1/categories` - list categories
- [ ] `GET /api/v1/categories/:id/foods` - foods by category

**Status: 🚧 BELUM MULAI**

---

## 🔄 Phase 5: Portion Size (PENDING)

### 5.1 Model & Migration
- [ ] Model `AssociatedFood`
- [ ] Model `FoodPortionSizeMethod`
- [ ] Model `AsServedSet`
- [ ] Model `AsServedImage`
- [ ] Migration `004_create_portion.sql`

### 5.2 Admin Endpoints
- [ ] Manage portion methods per food
- [ ] Upload as-served images
- [ ] Create as-served sets
- [ ] Manage associated foods

### 5.3 Public Endpoints
- [ ] `GET /api/v1/foods/:id/portion-methods` - list methods
- [ ] `GET /api/v1/as-served-sets/:id/images` - list images

**Status: 🚧 BELUM MULAI**

---

## 🔄 Phase 6: Survey Submission (PENDING)

### 6.1 Model & Migration
- [ ] Model `SurveySubmission`
- [ ] Migration `005_create_submissions.sql`

### 6.2 Respondent Flow
- [ ] `POST /api/v1/surveys/:accessToken/join` - join survey
- [ ] `GET /api/v1/surveys/:accessToken` - get survey detail
- [ ] `POST /api/v1/submissions` - submit survey

### 6.3 Admin Endpoints
- [ ] `GET /api/v1/admin/surveys/:id/submissions` - list submissions
- [ ] `GET /api/v1/admin/submissions/:id` - detail submission
- [ ] `GET /api/v1/admin/submissions/export` - export CSV/Excel

### 6.4 Nutrition Calculation
- [ ] Service hitung nutrisi per meal
- [ ] Agregasi nutrisi per hari
- [ ] Store calculated results

**Status: 🚧 BELUM MULAI**

---

## 🔄 Phase 7: File Upload & Storage (PENDING)

### 7.1 Setup
- [ ] Config upload path
- [ ] Validasi file type (image only)
- [ ] Validasi file size (max 10MB)

### 7.2 Endpoints
- [ ] `POST /api/v1/upload` - upload image
- [ ] Serve static files `/uploads/*`

### 7.3 As-Served Images
- [ ] Upload ke `/uploads/as-served/`
- [ ] Generate thumbnail
- [ ] Store path di database

**Status: 🚧 BELUM MULAI**

---

## 📋 Ringkasan Progress

| Phase | Fitur | Status | Progress |
|-------|-------|--------|----------|
| 1 | Foundation | ✅ Done | 100% |
| 2 | Auth | ✅ Done | 100% |
| 3 | Survey | 🚧 Next | 0% |
| 4 | Food DB | 🚧 Pending | 0% |
| 5 | Portion | 🚧 Pending | 0% |
| 6 | Submission | 🚧 Pending | 0% |
| 7 | Upload | 🚧 Pending | 0% |

---

## 🎯 Prioritas Selanjutnya

1. **Survey Management** - Admin bisa create/manage survey
2. **Food Database** - CRUD foods & categories
3. **Portion Size** - Setup portion methods
4. **Submission Flow** - Respondent bisa submit survey
5. **Export Data** - Admin bisa export hasil survey

---

## 📝 Notes

- Setiap phase menggunakan Clean Architecture (Handler → Service → Repository)
- Semua endpoint admin dilindungi middleware `JWTAuth` + `AdminOnly`
- Semua endpoint public/respondent menggunakan `accessToken` survey
- Gunakan transaction untuk operasi yang melibatkan multiple tables
- Validasi input menggunakan `go-playground/validator`
