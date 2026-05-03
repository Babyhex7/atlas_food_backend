# 🛠️ Tech Stack

## Backend Stack

| Layer             | Technology                       | Purpose             |
| ----------------- | -------------------------------- | ------------------- |
| **Language**      | Go (Golang)                      | Backend language    |
| **Framework**     | Gin                              | HTTP web framework  |
| **Database**      | MySQL                            | Primary datastore   |
| **ORM**           | GORM                             | Database operations |
| **Auth**          | JWT (stateless) + Refresh Tokens | Authentication      |
| **Validation**    | go-playground/validator          | Request validation  |
| **Password Hash** | bcrypt                           | Password encryption |

---

## Project Structure (Go - Clean Architecture)

```
atlas_food/
├── cmd/
│   └── api/                     # Entry point API server
│       └── main.go
│   └── worker/                  # (NANTI) Background worker
│       └── main.go
├── internal/
│   ├── config/                  # Config & DB connection
│   │   ├── config.go
│   │   └── database.go
│   │
│   ├── domain/                  # Core business logic (modular per fitur)
│   │   ├── auth/
│   │   │   ├── handler.go       # HTTP handler
│   │   │   ├── service.go       # Business logic
│   │   │   ├── repository.go    # DB access
│   │   │   ├── dto.go           # Request/response
│   │   │   └── model.go         # Domain model
│   │   │
│   │   ├── survey/
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── dto.go
│   │   │   └── model.go
│   │   │
│   │   ├── food/
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── dto.go
│   │   │   └── model.go
│   │   │
│   │   ├── portion/             # Portion size methods (as_served - untuk food & drinks)
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── dto.go
│   │   │   └── model.go
│   │   │
│   │   └── submission/
│   │       ├── handler.go
│   │       ├── service.go
│   │       ├── repository.go
│   │       ├── dto.go
│   │       └── model.go
│   │
│   ├── pkg/                     # Internal shared packages
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── logger.go
│   │   │
│   │   ├── utils/
│   │   │   ├── hash.go
│   │   │   ├── jwt.go
│   │   │   ├── response.go
│   │   │   └── validator.go
│   │   │
│   │   └── errors/              # Custom errors
│   │       └── errors.go
│   │
│   └── router/                  # Route setup
│       └── router.go
│
├── pkg/                         # External shared packages (bisa di-import project lain)
│   └── (kosong dulu untuk MVP)
│
├── migrations/                  # DB migrations
│   ├── 001_create_users.sql
│   ├── 002_create_surveys.sql
│   └── ...
│
├── uploads/                     # Local image storage
│   └── as-served/               # Portion images (food & drinks)
│
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

---

## Kenapa Structure Ini Lebih Bagus?

| Aspek | Old (Layer-based) | New (Feature-based) |
|-------|-------------------|---------------------|
| **Scalability** | Berantakan kalau 50+ tabel | Tiap domain terisolasi |
| **Maintainability** | Harus buka 5 folder beda | 1 folder = 1 fitur lengkap |
| **Testing** | Susah mock layer lain | Tiap domain bisa unit test sendiri |
| **Team Work** | Conflict di folder sama | Beda fitur = beda folder |
| **Microservices** | Susah split | Gampang, tinggal pindahin folder `domain/xxx` |

### Contoh: Semua File Portion di 1 Folder

```
internal/domain/portion/
├── model.go           # AsServedSet, AsServedImage (drinkware: SKIP untuk MVP)
├── repository.go      # Query ke DB
├── service.go         # Business logic (kalkulasi gram, dll)
├── handler.go         # API endpoints
└── dto.go             # Request/response structs
```

**Mau edit fitur portion?** Cukup buka 1 folder ✅

---

## Dependency Requirements

```go
// go.mod
module atlas_food

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/google/uuid v1.3.0
    golang.org/x/crypto v0.14.0
    gorm.io/driver/mysql v1.5.2
    gorm.io/gorm v1.25.5
    github.com/go-playground/validator/v10 v10.15.5
    github.com/joho/godotenv v1.5.1
)
```

---

## Environment Variables

```bash
# .env.example
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=atlas_food

# JWT
JWT_SECRET=your-super-secret-key
JWT_EXPIRATION=24h
REFRESH_TOKEN_EXPIRATION=168h

# Server
SERVER_PORT=8080
SERVER_MODE=debug

# Upload
UPLOAD_PATH=./uploads
MAX_UPLOAD_SIZE=10485760  # 10MB
```
