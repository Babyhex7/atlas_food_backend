# 14 — Admin Portal & Food Annotation CMS (Brief Lengkap)

> **Status:** Brief / Spec dokumen (belum implementasi penuh)  
> **Tanggal:** 2026-07-19  
> **Scope:** Panel admin CRUD end-to-end yang nyambung ke **Dietary Recall** + **Find Your Food**, plus **Food Annotation CMS** (alat anotasi polygon terpisah dari app user, ala Intake24).  
> **Stack proyek aktual:** Next.js (FE) + Go/Gin (BE) + MySQL/GORM + JWT RBAC — **bukan** Prisma. CMS memakai stack yang sama agar satu repo / satu login.

---

## 1. Visi: Dua permukaan aplikasi

Seperti Intake24, proses anotasi konten **terpisah** dari aplikasi yang dipakai responden.

```
┌─────────────────────────────────────┐
│  A. Food Annotation + Admin CMS     │
│     (Internal — role: admin)        │
│  - CRUD survey, food, kategori      │
│  - As-served / portion methods      │
│  - Annotation Tool (polygon)        │
│  - Submissions / export             │
└──────────────────┬──────────────────┘
                   │ Save metadata + polygon JSON
                   ▼
            MySQL + /uploads
                   ▲
                   │ Fetch published data
┌──────────────────┴──────────────────┐
│  B. User Application                │
│     (Public / respondent)           │
│  - Find Your Food                   │
│  - Dietary Recall wizard            │
│  - Portion picker (SVG dari JSON)   │
└─────────────────────────────────────┘
```

**Prinsip:** Admin/tim gizi mengunggah & menganotasi. User **tidak pernah** tahu ada polygon — mereka hanya klik area / pilih porsi.

---

## 2. Status saat ini vs target

| Modul | Sekarang (repo) | Target brief ini |
|-------|-----------------|------------------|
| Login admin (JWT + `AdminOnly`) | ✅ Ada | Tetap sama — **satu login**, RBAC ketat |
| Survey CRUD | ✅ Sebagian | Lengkap + status draft/active/closed |
| Food / Category CRUD | ⚠️ FE stub / BE ada | FE–BE penuh, sync Find Food + Recall |
| As-served sets & images | ⚠️ Sebagian | CRUD penuh + link ke food |
| Portion methods | ⚠️ Sebagian | CRUD penuh per food |
| Submissions list/export | ✅ Ada | Tetap + filter |
| Collab admin food room | ✅ Ada | Opsional di layar edit |
| **Food Annotation CMS** | ❌ Belum | **Baru** — upload, polygon SVG, draft, publish |
| Public API baca anotasi | ❌ Belum | `GET` published untuk FE user |

---

## 3. Auth & RBAC (jangan miss)

### 3.1 Login — sama dengan sistem yang ada

| Item | Spesifikasi |
|------|-------------|
| Endpoint | `POST /api/v1/auth/login` |
| Token | JWT access + refresh (cookie FE: `atlas_token`) |
| Role | `admin` \| `respondent` |
| Admin gate FE | Route `/admin/**` → cek role `admin`, else redirect `/login` atau `/profile` |
| Admin gate BE | Middleware `JWTAuth()` + `AdminOnly()` di semua `/api/v1/admin/**` |
| Annotation CMS | **Hanya `admin`** (boleh nanti extend role `nutritionist` — out of scope v1) |

### 3.2 Matriks permission (v1)

| Aksi | admin | respondent |
|------|-------|------------|
| Login / profile | ✅ | ✅ |
| Panel `/admin/**` | ✅ | ❌ |
| CRUD survey/food/category/as-served/portion | ✅ | ❌ |
| Annotation CMS (upload, polygon, publish) | ✅ | ❌ |
| Find Food (public read) | ✅ | ✅ (tanpa login juga OK untuk public endpoints) |
| Recall (isi survey) | ✅ (preview) | ✅ |
| Lihat submissions / export | ✅ | ❌ |
| WebSocket collab admin food | ✅ | ❌ |

### 3.3 Flow login admin (end-to-end)

```
1. GET /login
2. POST /api/v1/auth/login { email, password }
3. Response: user.role === "admin"
4. FE setSession + cookies
5. Redirect → /admin/surveys (atau last redirect aman)
6. Setiap request admin: Authorization: Bearer <access_token>
7. 401 → refresh; 403 → toast + redirect keluar admin
```

**Tidak ada login terpisah** untuk Annotation CMS. Satu akun admin → sidebar menu “Anotasi”.

---

## 4. Arsitektur Admin Portal (modul CRUD)

Sidebar target `/admin`:

| Menu | Path | Tujuan bisnis | Dipakai oleh |
|------|------|---------------|--------------|
| Survey | `/admin/surveys` | Buat/kelola survey recall | Recall |
| Submissions | `/admin/surveys/:id/submissions` | Review & export hasil | Recall |
| Makanan | `/admin/foods` | Master food + gizi | Recall + Find Food |
| Kategori | `/admin/categories` | 13 kategori Atlas | Find Food + search |
| Foto Porsi (As-served) | `/admin/as-served-sets` | Set foto porsi berbobot gram | Recall portion step |
| Metode Porsi | `/admin/portion-methods` | Metode estimasi per food | Recall |
| **Anotasi** | `/admin/annotations` | Polygon area di foto atlas | Find Food detail / portion UI lanjutan |
| Upload | via API | File ke `/uploads/...` | Semua modul media |

### 4.1 Hubungan data ke produk user

```
categories ──┐
             ├── foods ──┬── food_nutrients          → Find Food detail gizi
             │           ├── food_portion_size_methods → Recall step porsi
             │           ├── as_served_sets/images     → Recall pilih foto porsi
             │           └── food_images + food_areas  → Annotation → SVG klik area
surveys ─────┴── submissions                          → Hasil recall
```

**Find Your Food** membaca:

- `GET /api/v1/public/foods/search`
- `GET /api/v1/public/foods/:id`
- `GET /api/v1/public/categories`
- **Baru:** `GET /api/v1/public/food-images/:id` (anotasi published)

**Recall** membaca:

- Public food search/detail (step add food)
- Portion photos / methods (step portion)
- `POST /api/v1/survey/submit`

Admin yang mengisi CRUD di atas **langsung mempengaruhi** data yang dilihat responden.

---

## 5. Food Annotation CMS (detail produk)

### 5.1 Mengapa terpisah dari “hardcode polygon di FE”

- Ratusan–ribuan foto → tidak mungkin edit source code
- Tim konten / ahli gizi mandiri
- FE user hanya **render JSON** (SVG overlay)

### 5.2 Fitur wajib CMS

1. Upload gambar makanan  
2. Zoom in / zoom out  
3. Klik titik-titik untuk buat polygon  
4. Drag titik untuk edit  
5. Hapus titik / hapus polygon  
6. Isi nama area (mis. “Chicken Breast”)  
7. Opsional: link ke `food_id` master (supaya klik area → food Atlas yang sama)  
8. **Auto-save Draft** setiap perubahan  
9. **Publish** → baru boleh dibaca app user  
10. Preview hasil (SVG overlay seperti user)  
11. Export JSON (opsional) / utama = simpan DB  

### 5.3 Flow anotasi (admin)

```
Upload gambar
    ↓
Gambar muncul di canvas (SVG layer di atas <img>)
    ↓
Zoom / pan
    ↓
Klik titik → polygon mengikuti bentuk makanan
    ↓
Isi nama (+ optional food_id)
    ↓
Auto-save Draft (debounce 1–2s)
    ↓
Ulangi untuk area lain (sayap, paha, …)
    ↓
Preview
    ↓
Publish
    ↓
Status = published → tersedia di User App
```

### 5.4 Contoh payload JSON (kontrak data)

```json
{
  "id": "uuid-food-image",
  "image_url": "/uploads/annotations/fried-chicken.jpg",
  "width": 1200,
  "height": 800,
  "status": "published",
  "foods": [
    {
      "id": "uuid-area-1",
      "name": "Chicken Breast",
      "food_id": "uuid-master-food-optional",
      "polygon": [
        [120, 180],
        [160, 170],
        [210, 210],
        [170, 250]
      ]
    },
    {
      "id": "uuid-area-2",
      "name": "Chicken Wing",
      "food_id": null,
      "polygon": [
        [300, 200],
        [340, 190],
        [360, 240]
      ]
    }
  ]
}
```

**Catatan koordinat:** simpan dalam **pixel space gambar asli** (`width`/`height` tersimpan). FE scale ke tampilan dengan `viewBox="0 0 width height"`.

### 5.5 Flow user app (setelah publish)

```
GET /api/v1/public/food-images/:id
        ↓
image_url + areas[].polygon
        ↓
Render <img> + <svg> polygons
        ↓
User klik polygon
        ↓
Highlight area + tampilkan nama / link ke food
        ↓
(Opsional) lanjut pilih porsi as-served untuk food_id tersebut
```

User **tidak** melihat tools anotasi.

### 5.6 Tech UI anotasi

| Pilihan | Keputusan brief |
|---------|-----------------|
| Canvas vs SVG | **SVG Overlay** (bukan Canvas) — lebih mudah hit-test, a11y, scale |
| State editor | Zustand atau local state + React Hook Form untuk metadata |
| Zoom/pan | CSS transform / SVG transform pada group |
| Draft | `PATCH` autosave; status `draft` \| `published` |
| Stack | Next.js + TS + Tailwind (sama admin) + API Go |

---

## 6. Database schema (baru + relasi)

### 6.1 Tabel baru

#### `food_images` (aset gambar yang dianotasi)

```sql
CREATE TABLE food_images (
    id CHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    image_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    width INT NOT NULL,
    height INT NOT NULL,
    status ENUM('draft', 'published', 'archived') NOT NULL DEFAULT 'draft',
    -- opsional: gambar “scene” yang terkait satu food master atau atlas set
    primary_food_id CHAR(36) NULL,
    created_by CHAR(36) NOT NULL,
    published_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (primary_food_id) REFERENCES foods(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id),
    INDEX idx_food_images_status (status)
);
```

#### `food_areas` (polygon per area)

```sql
CREATE TABLE food_areas (
    id CHAR(36) PRIMARY KEY,
    food_image_id CHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    food_id CHAR(36) NULL,          -- link ke master foods (opsional)
    polygon JSON NOT NULL,            -- [[x,y], [x,y], ...]
    z_index INT NOT NULL DEFAULT 0,  -- urutan layer klik
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (food_image_id) REFERENCES food_images(id) ON DELETE CASCADE,
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE SET NULL,
    INDEX idx_food_areas_image (food_image_id)
);
```

#### `food_image_drafts` (opsional — atau kolom di food_images cukup)

Untuk autosave, **v1 cukup** update row `food_images` + `food_areas` dengan `status='draft'`.  
Tabel draft terpisah **tidak wajib** kecuali butuh version history.

### 6.2 Relasi dengan modul yang sudah ada

| Tabel existing | Hubungan ke anotasi |
|----------------|---------------------|
| `foods` | `food_areas.food_id` / `food_images.primary_food_id` |
| `as_served_images` | **Berbeda:** foto porsi berbobot gram untuk estimasi; anotasi = region di satu foto scene |
| `categories` | Lewat `foods.category_id` |

**Jangan campur** as-served weight photos dengan polygon scene annotation di satu tabel — beda use case.

### 6.3 ERD ringkas

```
users ||--o{ food_images : creates
foods ||--o{ food_images : primary_optional
food_images ||--o{ food_areas : contains
foods ||--o{ food_areas : linked_optional
```

---

## 7. API end-to-end

Base: `http://localhost:8080/api/v1`  
Auth admin: `Authorization: Bearer <token>` + role admin.

### 7.1 Auth (existing — dipakai ulang)

| Method | Path | Role |
|--------|------|------|
| POST | `/auth/login` | public |
| POST | `/auth/refresh` | public |
| GET | `/auth/me` | JWT |

### 7.2 Admin — Survey (existing, pastikan FE lengkap)

| Method | Path | Keterangan |
|--------|------|------------|
| GET/POST | `/admin/surveys` | list / create |
| GET/PUT/DELETE | `/admin/surveys/:id` | detail / update / delete |
| POST | `/admin/surveys/:id/clone` | clone |
| POST | `/admin/surveys/:id/regenerate-token` | token baru |
| GET | `/admin/surveys/:id/submissions` | list hasil recall |
| GET | `/admin/surveys/:id/export` | export |
| GET | `/admin/submissions/:id` | detail |

### 7.3 Admin — Food & Category (existing BE, FE harus penuh)

| Method | Path | Keterangan |
|--------|------|------------|
| CRUD | `/admin/foods` | master makanan |
| CRUD | `/admin/categories` | kategori |
| POST | `/admin/foods/:id/portion-methods` | metode porsi |
| CRUD | as-served sets/images | foto porsi + gram |
| POST | `/upload` | upload file |

### 7.4 Admin — Annotation CMS (**baru**)

| Method | Path | Body / notes |
|--------|------|----------------|
| GET | `/admin/food-images` | `?status=draft\|published` list |
| POST | `/admin/food-images` | create setelah upload: `{ title, image_url, width, height }` → status `draft` |
| GET | `/admin/food-images/:id` | detail + areas |
| PATCH | `/admin/food-images/:id` | update title/meta; **autosave** |
| POST | `/admin/food-images/:id/publish` | set `published`, set `published_at` |
| POST | `/admin/food-images/:id/unpublish` | kembali draft |
| DELETE | `/admin/food-images/:id` | soft/hard delete |
| PUT | `/admin/food-images/:id/areas` | **replace-all areas** (cocok autosave editor): `{ areas: [{ id?, name, food_id?, polygon, z_index }] }` |
| POST | `/admin/food-images/:id/areas` | add one area (opsional) |
| PATCH | `/admin/food-areas/:id` | update one |
| DELETE | `/admin/food-areas/:id` | delete one |
| GET | `/admin/food-images/:id/export` | download JSON |

**Validasi publish:**

- `image_url` ada  
- `width/height` > 0  
- minimal 1 area  
- setiap polygon ≥ 3 titik  
- koordinat dalam `[0..width]` × `[0..height]`

### 7.5 Public — konsumsi anotasi (**baru**)

| Method | Path | Auth | Keterangan |
|--------|------|------|------------|
| GET | `/public/food-images` | none | list **published** only (pagination) |
| GET | `/public/food-images/:id` | none | detail + areas (published only) |
| GET | `/public/foods/:id/images` | none | published images linked ke food (opsional) |

404 jika draft diakses publik.

### 7.6 Contoh request autosave

```http
PUT /api/v1/admin/food-images/{id}/areas
Authorization: Bearer <admin_jwt>
Content-Type: application/json

{
  "areas": [
    {
      "name": "Chicken Breast",
      "food_id": "uuid-or-null",
      "polygon": [[120,180],[160,170],[210,210],[170,250]],
      "z_index": 0
    }
  ]
}
```

Response 200:

```json
{
  "status": "success",
  "data": {
    "food_image_id": "...",
    "status": "draft",
    "areas_count": 1,
    "updated_at": "2026-07-19T14:00:00Z"
  }
}
```

### 7.7 Sequence end-to-end (anotasi → user klik)

```
Admin FE          Admin API           DB            User FE         Public API
   |                  |                |               |                |
   |-- upload img --->|                |               |                |
   |-- POST image --->|-- insert ----->|               |                |
   |-- draw poly ---->|                |               |                |
   |-- PUT areas ---->|-- upsert ----->|               |                |
   |-- POST publish ->|-- status=pub ->|               |                |
   |                  |                |               |-- GET image -->|
   |                  |                |<-- published -|                |
   |                  |                |               |-- render SVG   |
   |                  |                |               |-- click area   |
```

---

## 8. FE Admin — struktur halaman (target)

```
app/admin/
├── layout.tsx                 ← sidebar + AdminOnly guard
├── surveys/
├── foods/
├── categories/
├── as-served-sets/
├── portion-methods/
└── annotations/
    ├── page.tsx               ← list draft/published
    ├── new/page.tsx           ← upload → create
    └── [id]/page.tsx          ← Annotation Editor
```

### 8.1 Annotation Editor UI (wajib)

| Area UI | Fungsi |
|---------|--------|
| Toolbar | Zoom ±, undo/redo (nice), hapus titik, hapus polygon, pan |
| Canvas | `<img>` + SVG polygons; mode draw / edit |
| Side panel | Daftar area, nama, food picker, warna |
| Status chip | Draft / Published |
| Autosave indicator | “Menyimpan…” / “Tersimpan HH:mm:ss” |
| Actions | Preview, Publish, Unpublish, Export JSON |

### 8.2 Autosave draft (wajib)

```
Setiap perubahan polygon/nama
    → debounce 1500ms
    → PUT /admin/food-images/:id/areas
    → tampilkan "Tersimpan"
Browser crash / refresh
    → GET detail → restore draft
Publish
    → POST .../publish
    → data masuk public API
```

Tanpa autosave = risiko hilang 30 polygon — **non-negotiable**.

---

## 9. Struktur folder modular (lengkap)

Prinsip: **domain-driven**, mengikuti pola repo yang sudah ada.

| Layer | Aturan |
|-------|--------|
| `app/` (FE) | Route saja — tipis, compose domain components |
| `internal/domain/<nama>/` (FE) | components, hooks, services, types, schemas, store, constants |
| `internal/domain/<nama>/` (BE) | `handler.go`, `service.go`, `repository.go`, `model.go`, `dto.go` |
| User app vs Admin | Domain **sama** (food, portion); UI beda path (`app/find-food` vs `app/admin/...`) |
| Annotation | Domain baru `annotation` — CMS di admin; consumer SVG di user app |

---

### 9.1 Backend Go — tree target

```
atlas_food_backend/
├── cmd/
│   ├── api/
│   │   └── main.go                 # bootstrap hub, router, migrate
│   └── seed/
│       └── main.go
├── migrations/                     # opsional SQL; AutoMigrate di main tetap
│   └── 0XX_create_food_images.sql
├── uploads/
│   ├── atlas/                      # existing Find Food photos
│   ├── as-served/                  # portion weight photos
│   └── annotations/                # scene images untuk polygon CMS
├── internal/
│   ├── bootstrap/
│   ├── config/
│   ├── router/
│   │   └── router.go               # wire semua domain + middleware
│   ├── pkg/
│   │   ├── middleware/
│   │   │   ├── auth.go             # JWTAuth, AdminOnly, RespondentOnly
│   │   │   ├── cors.go
│   │   │   └── ...
│   │   ├── utils/                  # jwt, response
│   │   └── groq/
│   └── domain/
│       ├── auth/
│       │   ├── handler.go
│       │   ├── service.go
│       │   ├── repository.go
│       │   ├── model.go
│       │   └── dto.go
│       ├── survey/
│       │   ├── handler.go
│       │   ├── service.go
│       │   ├── repository.go
│       │   ├── model.go
│       │   └── dto.go
│       ├── food/                   # master food + public Find Food API
│       │   ├── handler.go          # admin CRUD
│       │   ├── public_handler.go   # /public/foods, categories
│       │   ├── service.go
│       │   ├── repository.go
│       │   ├── model.go
│       │   └── dto.go
│       ├── submission/
│       │   ├── handler.go
│       │   ├── service.go
│       │   ├── repository.go
│       │   ├── model.go
│       │   └── dto.go
│       ├── upload/
│       │   └── handler.go
│       ├── ai/
│       │   ├── handler.go
│       │   ├── service.go
│       │   ├── repository.go
│       │   ├── model.go
│       │   └── dto.go
│       ├── collab/                 # WebSocket (existing)
│       │   ├── handler.go
│       │   ├── hub.go
│       │   ├── room.go
│       │   ├── client.go
│       │   ├── message.go
│       │   ├── lock.go
│       │   └── service.go
│       └── annotation/             # ★ BARU — Food Annotation CMS
│           ├── handler.go          # admin routes
│           ├── public_handler.go   # GET /public/food-images*
│           ├── service.go          # publish rules, validate polygon
│           ├── repository.go
│           ├── model.go            # FoodImage, FoodArea
│           ├── dto.go
│           └── validate.go         # ≥3 titik, bounds width/height
├── docs/
│   └── 14-admin-portal-and-food-annotation-cms.md
└── go.mod
```

**Wiring router (sketch):**

```go
// internal/router/router.go
annotationHandler := annotation.NewHandler(annotationService)
adminAnn := v1.Group("/admin/food-images", middleware.JWTAuth(), middleware.AdminOnly())
{
    adminAnn.GET("", annotationHandler.List)
    adminAnn.POST("", annotationHandler.Create)
    adminAnn.GET("/:id", annotationHandler.Get)
    adminAnn.PATCH("/:id", annotationHandler.Update)
    adminAnn.PUT("/:id/areas", annotationHandler.ReplaceAreas) // autosave
    adminAnn.POST("/:id/publish", annotationHandler.Publish)
    adminAnn.POST("/:id/unpublish", annotationHandler.Unpublish)
    adminAnn.DELETE("/:id", annotationHandler.Delete)
    adminAnn.GET("/:id/export", annotationHandler.ExportJSON)
}
publicAnn := v1.Group("/public/food-images")
{
    publicAnn.GET("", annotationHandler.ListPublished)
    publicAnn.GET("/:id", annotationHandler.GetPublished)
}
```

---

### 9.2 Frontend Next.js — tree target

```
atlas_food_frontend/
├── app/
│   ├── layout.tsx
│   ├── page.tsx                          # landing
│   ├── login/
│   ├── register/
│   ├── profile/
│   ├── find-food/                        # USER APP
│   │   ├── page.tsx
│   │   ├── FindFoodContent.tsx
│   │   ├── [id]/page.tsx
│   │   └── category/[code]/page.tsx
│   ├── surveys/
│   │   └── [accessToken]/
│   │       ├── join/
│   │       ├── recall/                   # USER — dietary recall
│   │       └── done/
│   └── admin/                            # ADMIN CMS (AdminOnly)
│       ├── layout.tsx                    # sidebar + guard role=admin
│       ├── surveys/
│       │   ├── page.tsx
│       │   ├── new/page.tsx
│       │   └── [id]/
│       │       ├── page.tsx
│       │       └── submissions/page.tsx
│       ├── foods/
│       │   ├── page.tsx
│       │   ├── new/page.tsx
│       │   └── [id]/page.tsx
│       ├── categories/
│       │   ├── page.tsx
│       │   ├── new/page.tsx
│       │   └── [id]/page.tsx
│       ├── as-served-sets/
│       │   ├── page.tsx
│       │   ├── new/page.tsx
│       │   └── [id]/
│       │       ├── page.tsx
│       │       └── images/page.tsx
│       ├── portion-methods/
│       │   ├── page.tsx
│       │   └── [id]/page.tsx
│       └── annotations/                 # ★ BARU — Annotation CMS routes
│           ├── page.tsx                  # list draft | published
│           ├── new/page.tsx              # upload → create draft
│           └── [id]/
│               ├── page.tsx              # Annotation Editor
│               └── preview/page.tsx      # optional full preview
├── internal/
│   ├── components/                       # shared UI (AppHeader, EmptyState, …)
│   ├── providers/
│   │   └── auth-provider.tsx
│   ├── lib/                              # axios, cookies, cn, layout
│   ├── pkg/                              # Button, Input generik
│   ├── hooks/
│   ├── services/                         # legacy thin API helpers (opsional)
│   └── domain/
│       ├── auth/
│       │   ├── components/
│       │   ├── hooks/
│       │   ├── services/
│       │   ├── store/
│       │   ├── schemas/
│       │   ├── types/
│       │   ├── constants/
│       │   └── index.ts
│       ├── survey/
│       │   ├── components/               # SurveyList, SurveyForm, …
│       │   ├── hooks/
│       │   ├── services/
│       │   ├── store/
│       │   ├── schemas/
│       │   ├── types/
│       │   ├── constants/
│       │   └── index.ts
│       ├── food/
│       │   ├── components/               # FoodList, FoodForm, FoodSearch
│       │   ├── hooks/
│       │   ├── services/
│       │   ├── store/
│       │   ├── schemas/
│       │   ├── types/
│       │   └── index.ts
│       ├── category/
│       │   ├── components/
│       │   ├── hooks/
│       │   ├── services/
│       │   ├── store/
│       │   └── index.ts
│       ├── portion/
│       │   ├── components/               # AsServedSetForm, PortionSelector, …
│       │   ├── hooks/
│       │   ├── services/
│       │   ├── store/
│       │   ├── types/
│       │   └── index.ts
│       ├── submission/
│       │   ├── components/
│       │   ├── hooks/
│       │   ├── services/
│       │   └── index.ts
│       ├── recall/                       # USER wizard (bukan admin)
│       │   ├── components/
│       │   ├── hooks/
│       │   ├── services/
│       │   ├── types/
│       │   └── index.ts
│       ├── collab/
│       │   ├── components/
│       │   ├── hooks/
│       │   ├── store/
│       │   └── index.ts
│       ├── nutrition/
│       └── annotation/                   # ★ BARU — domain Annotation CMS
│           ├── components/
│           │   ├── AnnotationList.tsx          # tabel draft/published
│           │   ├── AnnotationUploader.tsx      # upload + create
│           │   ├── AnnotationEditor.tsx        # shell editor
│           │   ├── AnnotationCanvas.tsx        # img + SVG overlay
│           │   ├── PolygonLayer.tsx            # draw/edit polygons
│           │   ├── VertexHandle.tsx            # drag titik
│           │   ├── ZoomToolbar.tsx
│           │   ├── AreaSidePanel.tsx           # list area + nama + food picker
│           │   ├── AutosaveIndicator.tsx
│           │   ├── PublishBar.tsx              # Publish / Unpublish / Export
│           │   └── AnnotationPreview.tsx       # preview seperti user
│           ├── hooks/
│           │   ├── useAnnotationQueries.ts     # React Query list/detail
│           │   ├── useAnnotationMutations.ts   # create, save areas, publish
│           │   ├── useAnnotationAutosave.ts    # debounce PUT areas
│           │   ├── usePolygonEditor.ts         # state titik, mode draw/edit
│           │   └── useCanvasTransform.ts       # zoom/pan
│           ├── services/
│           │   └── annotationService.ts        # API client admin + public
│           ├── store/
│           │   └── annotationEditorStore.ts    # Zustand: areas, selection, dirty
│           ├── schemas/
│           │   └── annotationSchema.ts         # zod: name, polygon ≥ 3
│           ├── types/
│           │   └── annotation.ts               # FoodImage, FoodArea, Point
│           ├── constants/
│           │   └── annotationStatus.ts         # draft | published | archived
│           ├── utils/
│           │   ├── polygonMath.ts              # hit-test, bounds check
│           │   └── scaleCoordinates.ts         # pixel ↔ display
│           └── index.ts
│       └── find-food-ui/ (opsional)            # atau taruh di components/
│           └── AnnotatedFoodViewer.tsx         # USER: render SVG dari public API
├── docs/
│   └── 14-admin-portal-and-food-annotation-cms.md
├── styles/
└── package.json
```

---

### 9.3 Pemetaan domain → fitur produk

```
                    ┌──────────── admin/* ────────────┐
                    │                                 │
 survey ────────────┤ CRUD survey / submissions       ├──→ Recall (join/submit)
 food + category ───┤ CRUD master + gizi              ├──→ Find Food + Recall search
 portion ───────────┤ as-served + methods             ├──→ Recall Step Portion
 annotation ────────┤ polygon CMS draft/publish       ├──→ Find Food SVG / area click
 upload ────────────┤ /uploads/*                      ├──→ semua media
 collab ────────────┤ WS room admin:food-db (opsional)│
                    └─────────────────────────────────┘
```

---

### 9.4 Kontrak modul `annotation` (tanggung jawab file)

| File / folder | Tanggung jawab tunggal |
|---------------|------------------------|
| BE `model.go` | Struct GORM `FoodImage`, `FoodArea` |
| BE `repository.go` | SQL/GORM CRUD + replace areas transaction |
| BE `service.go` | Validasi polygon, publish rules, tidak expose draft ke public |
| BE `handler.go` | HTTP admin |
| BE `public_handler.go` | HTTP public published-only |
| FE `annotationService.ts` | Axios calls only |
| FE `usePolygonEditor.ts` | Logika titik/polygon, bukan UI fetch |
| FE `useAnnotationAutosave.ts` | Debounce + mutation, bukan canvas |
| FE `AnnotationCanvas.tsx` | Render + pointer events |
| FE `AnnotatedFoodViewer.tsx` | User read-only SVG (domain annotation atau food) |
| `app/admin/annotations/**` | Routing + compose components — **tanpa business logic** |

---

### 9.5 Dependency antar domain (jangan circular)

```
annotation ──depends──► food          (food_id picker, link master)
annotation ──depends──► upload        (image_url)
portion    ──depends──► food
recall     ──depends──► food, portion, survey, submission
find-food pages ──depends──► food (+ annotation viewer published)
admin layout ──depends──► auth (role guard)
```

`food` **tidak** boleh import `annotation` di BE model inti; relasi cukup FK di sisi annotation.

---

### 9.6 Naming convention

| Tempat | Pola |
|--------|------|
| BE package | `annotation` (singular domain name) |
| FE folder | `internal/domain/annotation/` |
| FE routes | `app/admin/annotations/` (plural resource URL) |
| API | `/admin/food-images`, `/public/food-images` |
| Tabel DB | `food_images`, `food_areas` |
| Komponen | PascalCase: `AnnotationEditor.tsx` |
| Hooks | `use` + Nama: `usePolygonEditor.ts` |

---

## 10. User stories (admin & anotasi)

1. Sebagai **admin**, saya login dengan akun yang sama, masuk panel admin.  
2. Sebagai **admin**, saya CRUD makanan/kategori agar muncul di Find Food & Recall.  
3. Sebagai **admin**, saya kelola survey + lihat submissions recall.  
4. Sebagai **admin/tim gizi**, saya upload foto scene, gambar polygon, isi nama, autosave.  
5. Sebagai **admin**, saya publish anotasi agar user app bisa fetch.  
6. Sebagai **responden**, saya di Find Food/Recall melihat area klik / porsi tanpa tahu CMS.  
7. Sebagai **admin**, saya tidak ingin coding ulang saat ada 100 gambar baru — cukup upload + anotasi.

---

## 11. Acceptance criteria (ringkas)

### Admin CRUD (nyambung produk)

- [ ] Semua route `/admin/**` hanya `role=admin`  
- [ ] Food CRUD tercermin di `GET /public/foods/search` & detail  
- [ ] Category CRUD tercermin di public categories  
- [ ] As-served images dipakai di Recall portion step  
- [ ] Survey create → access token → respondent join → submit → muncul di submissions  

### Annotation CMS

- [ ] Upload + create `food_images` draft  
- [ ] Buat ≥1 polygon (≥3 titik), edit drag, hapus  
- [ ] Autosave draft tanpa publish  
- [ ] Publish → public GET sukses; unpublish → public 404  
- [ ] User FE render SVG dari JSON tanpa hardcode koordinat  
- [ ] 100 gambar baru = 0 perubahan kode FE user  
- [ ] Folder mengikuti struktur modular §9 (domain `annotation` terpisah)

---

## 12. Fase implementasi (usulan)

| Fase | Isi | Hasil |
|------|-----|--------|
| **A** | Harden Admin CRUD FE–BE (food, category, as-served, portion) | Recall + Find Food data lengkap dari admin |
| **B** | Scaffold `internal/domain/annotation` BE + migration | Backend siap |
| **C** | Scaffold FE `domain/annotation` + `app/admin/annotations` + editor + autosave | CMS internal jalan |
| **D** | Public API + `AnnotatedFoodViewer` di Find Food | User konsumsi anotasi |
| **E** | Polish: undo, keyboard, batch upload, food picker | Produksi konten |

---

## 13. Out of scope v1 (catat agar tidak miss ekspektasi)

- Role `nutritionist` terpisah  
- AI auto-detect polygon  
- Mobile annotation app  
- Multi-language nama area  
- Version history penuh / CRDT collab di editor (boleh pakai collab room admin nanti)  
- Pindah ke Postgres/Prisma (tetap MySQL + Go)  
- Monorepo terpisah “annotation-only app” (cukup path `/admin/annotations` dulu)

---

## 14. Ringkasan keputusan arsitektur

1. **Satu login JWT**, RBAC `admin` vs `respondent`.  
2. **Admin Portal** = CRUD operasional + **Annotation CMS** sebagai menu di dalamnya (bukan repo terpisah wajib; *secara produk* terpisah dari user app).  
3. Polygon disimpan **JSON di DB**, FE user hanya fetch & render SVG.  
4. **Draft + Publish** + **autosave** wajib.  
5. As-served (berat gram) ≠ annotation scene (region klik) — dua model, boleh saling link lewat `food_id`.  
6. Semua perubahan admin harus **traceable** ke Find Food + Recall lewat tabel `foods` / public API yang sama.  
7. **Modular domain folders** — BE `domain/annotation`, FE `internal/domain/annotation` + `app/admin/annotations` (lihat §9).

---

## 15. Referensi dokumen terkait

- [01-overview.md](./01-overview.md) — scope MVP  
- [02-tech-stack.md](./02-tech-stack.md) — struktur proyek  
- [03-database-schema.md](./03-database-schema.md) — schema existing  
- [04-api-documentation.md](./04-api-documentation.md) — API existing  
- [05-workflow-alur.md](./05-workflow-alur.md) — alur recall  
- [12-find-your-food.md](./12-find-your-food.md) — Find Food  
- [13-realtime-collaboration.md](./13-realtime-collaboration.md) — collab (opsional di admin edit)

---

*Dokumen ini adalah brief implementasi. Setelah disetujui, lanjut rencana teknis per fase (migration → API → editor UI → public consumer).*
