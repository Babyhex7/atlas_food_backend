# 🤖 Briefing Implementasi AI — Atlas Food × Groq

> **Dokumen ini menjelaskan satu fitur AI saja:** Analisis nutrisi & rekomendasi yang ditampilkan **ketika user menekan tombol "AI Recommendation"** di halaman hasil survey. AI **tidak dipanggil saat submission**, melainkan lewat endpoint terpisah yang dipanggil on-demand.

---

## 📌 Gambaran Besar

Setelah respondent menekan tombol **Submit Survey**, backend menghitung total nutrisi dan menyimpan submission. Halaman hasil langsung tampil dengan data nutrisi (kalori, protein, dll) **tanpa perlu menunggu AI**.

Tombol **"AI Recommendation"** tersedia di bagian Actions. Saat dipencet, frontend memanggil endpoint baru `POST /ai/nutrition-analysis` dengan menyertakan `submission_id`. Backend mengambil data nutrisi dari submission tersebut, memanggil Groq, lalu mengembalikan hasil AI.

### Tampilan UI (Referensi Screenshot)

```
┌─────────────────────────────────────────────┐
│  Nutrition Result                           │
│  ─────────────────────────────────────────  │
│  Meal Items             [Edit List]         │
│  ┌──────────────────────────────────────┐   │
│  │ Nasi Goreng    380 kcal  Carb-Heavy  │   │
│  │ Telur Goreng   145 kcal  Protein Src │   │
│  │ Es Teh          20 kcal  Hydration   │   │
│  └──────────────────────────────────────┘   │
├─────────────────────────────────────────────┤
│  Estimated Intake Total         ⚡           │
│  545 TOTAL CALORIES (KCAL)                  │
│  Protein  24g (18%)                         │
│  Carbs    68g (50%)                         │
│  Fats     18g (32%)                         │
├─────────────────────────────────────────────┤
│  Actions                                    │
│  [ ✏  Edit Portions        ]                │
│  [ 🤖 AI Recommendation    ]  ← tombol ini  │
│  [ ↗  Export Data          ]                │
└─────────────────────────────────────────────┘
```

Saat tombol **AI Recommendation** diklik:
1. Button berubah jadi loading state
2. Frontend kirim `POST /ai/nutrition-analysis`
3. Backend ambil data submission → panggil Groq → return JSON
4. Halaman render hasil AI di bawah section Actions

---

## 1. Alur Data (End-to-End)

### Flow 1 — Submit Survey (tidak ada AI di sini)

```
Frontend              Backend (Submission)          DB
   │                          │                      │
   │  POST /submissions        │                      │
   │─────────────────────────▶│                      │
   │                          │  1. Simpan submission│
   │                          │  2. Hitung daily_tot │
   │                          │─────────────────────▶│
   │                          │                      │
   │  Response (cepat):       │                      │
   │  { submission_id, daily  │                      │
   │    nutrition summary }   │                      │
   │◀─────────────────────────│                      │
   │                          │                      │
   │  Tampilkan Nutrition     │                      │
   │  Result page             │                      │
```

### Flow 2 — User Klik "AI Recommendation"

```
Frontend              Backend (AI Endpoint)         Groq API          DB
   │                          │                         │              │
   │  POST /ai/nutrition-     │                         │              │
   │  analysis                │                         │              │
   │  { submission_id }       │                         │              │
   │─────────────────────────▶│                         │              │
   │                          │  1. Fetch submission    │              │
   │                          │  & daily_total from DB  │              │
   │                          │─────────────────────────────────────▶│
   │                          │                         │              │
   │                          │  2. Kirim ke Groq       │              │
   │                          │─────────────────────────▶             │
   │                          │                         │  3. Proses  │
   │                          │                         │  LLM        │
   │                          │  4. JSON result         │              │
   │                          │◀────────────────────────│              │
   │                          │                         │              │
   │                          │  5. Simpan ke           │              │
   │                          │  ai_result_logs         │              │
   │                          │─────────────────────────────────────▶│
   │                          │                         │              │
   │  6. Return ai_result     │                         │              │
   │◀─────────────────────────│                         │              │
   │                          │                         │              │
   │  Render hasil AI         │                         │              │
   │  di halaman              │                         │              │
```

**Poin penting:**
- Submit tetap **cepat** karena tidak ada call ke Groq.
- AI dipanggil **on-demand** hanya ketika user klik tombol.
- Jika submission sudah pernah dianalisis sebelumnya (ada di `ai_result_logs`), langsung return dari DB tanpa panggil Groq lagi (**cache by submission_id**).

---

## 2. Endpoint AI

### `POST /ai/nutrition-analysis`

**Akses:** Customer (harus login, submission milik user tersebut)

**Request Body:**
```json
{
  "submission_id": "uuid-submission"
}
```

**Response Sukses (200):**
```json
{
  "status": "success",
  "source": "groq",
  "data": {
    "overall_status": "less",
    "overall_message": "Your current nutrition is still below the recommended daily requirement...",
    "nutritional_analysis": [
      {
        "label": "Calories",
        "status": "low",
        "description": "Current calorie level is still relatively low for optimal daily energy needs."
      },
      {
        "label": "Protein",
        "status": "low",
        "description": "Protein source is limited and should be increased to support body recovery."
      },
      {
        "label": "Balance",
        "status": "partial",
        "description": "Your meal already contains sufficient carbs, but fiber and micronutrients are lacking."
      }
    ],
    "ai_recommendation": "To improve your nutritional balance, consider adding:\n- Grilled chicken or fish\n- Broccoli or spinach for fiber\n- Banana or apple for natural nutrients\n- More water intake",
    "recommended_foods": [
      "Grilled Chicken", "Boiled Egg", "Broccoli",
      "Spinach", "Banana", "Apple", "Greek Yogurt", "Mineral Water"
    ],
    "health_insight": {
      "title": "Mild Nutritional Deficiency",
      "description": "Your current meal composition is partially balanced, but additional protein, vegetables, and hydration are recommended."
    },
    "suggested_activities": ["Light Walking", "Yoga", "Stretching"]
  }
}
```

> `"source": "groq"` → fresh dari Groq API
> `"source": "cache"` → diambil dari `ai_result_logs` DB (sudah pernah dianalisis)

**Response Gagal — Submission tidak ditemukan (404):**
```json
{
  "status": "error",
  "message": "Submission not found or access denied"
}
```

**Response Gagal — Groq error (503):**
```json
{
  "status": "error",
  "message": "AI service temporarily unavailable, please try again"
}
```

---

## 3. Input yang Dikirim ke AI

AI menerima data nutrisi harian yang diambil dari submission tersimpan di DB. Data ini sudah dihitung saat submission dibuat (lihat rumus di `07-portion-calculation.md`).

### Struktur Input ke Groq
```json
{
  "daily_total": {
    "energy_kcal": 545,
    "protein_g": 24,
    "carbs_g": 68,
    "fat_g": 18,
    "fiber_g": 3.2,
    "sugar_g": 12.5,
    "sodium_mg": 840,
    "calcium_mg": 210,
    "iron_mg": 3.5,
    "vitamin_c_mg": 8
  },
  "meal_count": 3,
  "food_names": ["Nasi Goreng", "Telur Goreng", "Es Teh"]
}
```

> `food_names` disertakan agar AI bisa memberi rekomendasi yang **kontekstual** (tidak merekomendasikan makanan yang sudah dikonsumsi).

---

## 4. Output yang Diharapkan dari AI

AI **wajib** mengembalikan JSON murni. Tidak boleh ada teks narasi di luar JSON.

```json
{
  "overall_status": "less",
  "overall_message": "Your current nutrition is still below the recommended daily requirement. Additional balanced nutrients are needed.",

  "nutritional_analysis": [
    {
      "label": "Calories",
      "status": "low",
      "description": "Current calorie level is still relatively low for optimal daily energy needs."
    },
    {
      "label": "Protein",
      "status": "low",
      "description": "Protein source is limited and should be increased to support body recovery and muscle maintenance."
    },
    {
      "label": "Balance",
      "status": "partial",
      "description": "Your meal already contains sufficient carbohydrates, but fiber and micronutrient sources are still lacking."
    }
  ],

  "ai_recommendation": "To improve your nutritional balance, consider adding:\n- Grilled chicken or fish for additional protein\n- Vegetables such as broccoli or spinach for fiber and vitamins\n- Fruits like banana or apple for natural nutrients\n- More water intake to maintain hydration balance",

  "recommended_foods": [
    "Grilled Chicken", "Boiled Egg", "Broccoli",
    "Spinach", "Banana", "Apple", "Greek Yogurt", "Mineral Water"
  ],

  "health_insight": {
    "title": "Mild Nutritional Deficiency",
    "description": "Your current meal composition is considered partially balanced, but additional protein, vegetables, and hydration are recommended to better fulfill daily nutritional needs."
  },

  "suggested_activities": ["Light Walking", "Yoga", "Stretching"]
}
```

### Nilai `overall_status`
| Value | Arti | Warna UI |
|---|---|---|
| `"good"` | Nutrisi sudah cukup | Hijau |
| `"less"` | Nutrisi kurang dari kebutuhan | Merah |
| `"excess"` | Nutrisi berlebih | Kuning/Oranye |

---

## 5. Prompt Engineering

### System Prompt
```
Kamu adalah asisten nutrisi Atlas Food. Tugasmu adalah menganalisis data nutrisi harian
seorang responden dan memberikan umpan balik yang jelas dan akurat dalam bahasa Inggris.

ATURAN WAJIB:
1. Kembalikan HANYA JSON valid. DILARANG ada teks narasi di luar JSON.
2. Ikuti struktur JSON yang sudah ditentukan, tidak boleh berbeda.
3. overall_status HANYA boleh salah satu dari: "good", "less", "excess".
4. nutritional_analysis wajib berisi 3 item: Calories, Protein, Balance.
5. recommended_foods berisi 6-8 nama makanan yang belum ada di daftar makanan user.
6. suggested_activities berisi 2-4 aktivitas fisik ringan yang sesuai kondisi.
7. Semua teks dalam bahasa Inggris.
8. Gunakan Angka Kecukupan Gizi (AKG) Indonesia sebagai acuan:
   - Energi: 2150 kcal/hari
   - Protein: 65 g/hari
   - Lemak: 75 g/hari
   - Karbohidrat: 325 g/hari
```

### User Prompt (dibangun dinamis oleh backend)
```
Analisis nutrisi harian berikut:

Total Nutrisi Hari Ini:
- Energi: {energy_kcal} kcal (AKG: 2150 kcal)
- Protein: {protein_g} g (AKG: 65 g)
- Karbohidrat: {carbs_g} g (AKG: 325 g)
- Lemak: {fat_g} g (AKG: 75 g)
- Serat: {fiber_g} g (AKG: 30 g)
- Natrium: {sodium_mg} mg (batas: 2300 mg)
- Kalsium: {calcium_mg} mg (AKG: 1200 mg)

Makanan yang sudah dikonsumsi hari ini: {food_names_joined}

Berikan analisis dan rekomendasi dalam format JSON yang sudah ditentukan.
```

---

## 6. Skema Database — Tabel `ai_result_logs`

> Tabel ini bisa menggunakan **database terpisah** atau database yang sama. Bebas, yang penting terpisah logikanya dari domain submission.

```sql
CREATE TABLE ai_result_logs (
    id              CHAR(36) PRIMARY KEY DEFAULT (UUID()),

    -- Relasi ke submission (UNIQUE untuk cache — 1 submission = 1 hasil AI)
    submission_id   CHAR(36) NOT NULL UNIQUE,

    -- Input yang dikirim ke Groq (untuk audit/debug)
    input_payload   JSON NOT NULL,

    -- Output mentah dari Groq (disimpan as-is)
    raw_response    JSON NOT NULL,

    -- Field yang sudah di-parse (untuk query cepat)
    overall_status  ENUM('good', 'less', 'excess') NOT NULL,

    -- Metadata Groq
    model_used      VARCHAR(50) NOT NULL DEFAULT 'llama3-8b-8192',
    token_used      INT,
    latency_ms      INT,

    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- FK ke survey_submissions di DB yang sama (jika satu DB)
    -- Hapus baris FK ini jika pakai DB terpisah:
    FOREIGN KEY (submission_id) REFERENCES survey_submissions(id) ON DELETE CASCADE,

    INDEX idx_submission (submission_id),
    INDEX idx_status (overall_status)
);
```

**Cache logic:** Sebelum panggil Groq, backend cek dulu apakah `submission_id` sudah ada di `ai_result_logs`. Kalau ada, langsung return dari DB dengan `"source": "cache"`. Tidak ada panggilan redundan ke Groq.

---

## 7. Struktur Kode — Domain AI Terpisah

AI diimplementasikan sebagai **domain baru** `internal/domain/ai`, bukan diinjek ke domain submission.

```
internal/
├── domain/
│   ├── submission/          ← tidak diubah, submit tetap cepat
│   │   ├── service.go
│   │   ├── handler.go
│   │   └── dto.go
│   │
│   └── ai/                  ← domain baru
│       ├── handler.go       ← handle POST /ai/nutrition-analysis
│       ├── service.go       ← orchestrate: fetch submission → call Groq → save log
│       ├── repository.go    ← CRUD ke tabel ai_result_logs
│       └── dto.go           ← struct request/response AI
│
└── pkg/
    └── groq/                ← package shared (bisa dipakai domain lain nanti)
        └── client.go        ← HTTP wrapper ke Groq API
```

### Alur di `ai/service.go`

```go
func (s *aiService) GetNutritionAnalysis(userID, submissionID string) (*AIAnalysisResponse, error) {
    // 1. Validasi: submission harus milik userID ini
    submission, err := s.submissionRepo.FindByIDAndUser(submissionID, userID)
    if err != nil {
        return nil, ErrSubmissionNotFound
    }

    // 2. Cek cache di ai_result_logs
    cached, err := s.aiRepo.FindBySubmissionID(submissionID)
    if err == nil && cached != nil {
        // Sudah pernah dianalisis, langsung return
        return &AIAnalysisResponse{
            Source: "cache",
            Data:   cached.ParsedResult,
        }, nil
    }

    // 3. Bangun input ke Groq dari data submission
    input := buildGroqInput(submission)

    // 4. Panggil Groq
    aiResult, latencyMs, tokenUsed, err := s.groqClient.AnalyzeNutrition(input)
    if err != nil {
        return nil, ErrGroqUnavailable
    }

    // 5. Simpan ke ai_result_logs
    s.aiRepo.Save(AIResultLog{
        SubmissionID:  submissionID,
        InputPayload:  input,
        RawResponse:   aiResult,
        OverallStatus: aiResult.OverallStatus,
        ModelUsed:     s.groqClient.Model,
        TokenUsed:     tokenUsed,
        LatencyMs:     latencyMs,
    })

    // 6. Return hasil
    return &AIAnalysisResponse{
        Source: "groq",
        Data:   aiResult,
    }, nil
}
```

### Handler `ai/handler.go`

```go
// POST /ai/nutrition-analysis
func (h *AIHandler) GetNutritionAnalysis(c *gin.Context) {
    userID := c.GetString("user_id")  // dari JWT middleware

    var req AIAnalysisRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"status": "error", "message": "Invalid request"})
        return
    }

    result, err := h.service.GetNutritionAnalysis(userID, req.SubmissionID)
    if err != nil {
        switch err {
        case ErrSubmissionNotFound:
            c.JSON(404, gin.H{"status": "error", "message": "Submission not found or access denied"})
        case ErrGroqUnavailable:
            c.JSON(503, gin.H{"status": "error", "message": "AI service temporarily unavailable, please try again"})
        default:
            c.JSON(500, gin.H{"status": "error", "message": "Internal server error"})
        }
        return
    }

    c.JSON(200, gin.H{
        "status": "success",
        "source": result.Source,
        "data":   result.Data,
    })
}
```

---

## 8. Routing

Tambahkan route baru di `router.go`:

```go
// AI Routes — customer only
aiGroup := router.Group("/ai")
aiGroup.Use(middleware.AuthMiddleware())
aiGroup.Use(middleware.RoleMiddleware("customer"))
{
    aiGroup.POST("/nutrition-analysis", aiHandler.GetNutritionAnalysis)
}
```

---

## 9. Package `pkg/groq/client.go`

```go
type Client struct {
    apiKey  string
    Model   string
    baseURL string
    timeout time.Duration
}

func NewClient(apiKey, model string) *Client { ... }

// AnalyzeNutrition kirim ke Groq, return (result, latencyMs, tokenUsed, error)
func (c *Client) AnalyzeNutrition(input GroqInput) (*AIResult, int, int, error) {
    // 1. Build system prompt dan user prompt
    // 2. POST ke https://api.groq.com/openai/v1/chat/completions
    // 3. Parse response JSON
    // 4. Validasi struktur output
    // 5. Return *AIResult, latency, token_used
}
```

---

## 10. Environment Variables

Tambahkan ke `.env` dan `.env.example`:

```bash
# =========================================
# GROQ AI CONFIGURATION
# =========================================
GROQ_API_KEY=gsk_xxxxxxxxxxxxxxxxxxxxxxxx   # Dari console.groq.com (GRATIS)
GROQ_MODEL=llama3-8b-8192
GROQ_BASE_URL=https://api.groq.com/openai/v1
GROQ_TIMEOUT_SECONDS=15
GROQ_MAX_TOKENS=512
```

### Cara Dapat API Key (Gratis)
1. Buka [console.groq.com](https://console.groq.com)
2. Login dengan Google → Menu **API Keys** → **Create API Key**
3. Copy key (`gsk_xxx...`) → paste ke `.env`
4. Free tier: ~14.400 request/hari. Lebih dari cukup.

---

## 11. Sinkronisasi dengan Dokumen Lain

| Dokumen | Perubahan yang Perlu Dilakukan |
|---|---|
| `03-database-schema.md` | Tambahkan tabel `ai_result_logs` (tabel baru) |
| `04-api-documentation.md` | Tambahkan endpoint baru `POST /ai/nutrition-analysis` |
| `08-erd.md` | Tambahkan entitas `ai_result_logs` dengan relasi ke `survey_submissions` |
| `10-implementation-steps.md` | Phase 8 — buat domain `ai/` dan package `pkg/groq/` |

---

## 12. Checklist Implementasi

```
Phase 8: AI Nutritional Analysis (Estimasi: 1-2 hari)
```

- [ ] **8.1** Tambahkan env vars Groq ke `.env` dan `.env.example`
- [ ] **8.2** Buat `internal/pkg/groq/client.go` dengan method `AnalyzeNutrition()`
- [ ] **8.3** Buat `internal/domain/ai/dto.go` — struct `AIAnalysisRequest`, `AIAnalysisResponse`, `AIResult`, dll
- [ ] **8.4** Buat `internal/domain/ai/repository.go` — CRUD ke tabel `ai_result_logs`
- [ ] **8.5** Buat `internal/domain/ai/service.go` — logic cache check → Groq → save log
- [ ] **8.6** Buat `internal/domain/ai/handler.go` — handle `POST /ai/nutrition-analysis`
- [ ] **8.7** Buat migration `migrations/008_create_ai_result_logs.sql`
- [ ] **8.8** Tambahkan `AutoMigrate(&AIResultLog{})` di `main.go`
- [ ] **8.9** Daftarkan route `POST /ai/nutrition-analysis` di `router.go`
- [ ] **8.10** Test Postman: submit dulu → copy `submission_id` → hit `/ai/nutrition-analysis`
- [ ] **8.11** Test cache: hit endpoint 2x dengan `submission_id` yang sama → response kedua `"source": "cache"`
- [ ] **8.12** Test graceful error: set `GROQ_API_KEY` salah → endpoint return 503 (submission tidak terdampak)

---

## ⚠️ Catatan Penting

> **Ownership check wajib:** Endpoint `/ai/nutrition-analysis` harus memvalidasi bahwa `submission_id` yang dikirim memang milik user yang sedang login. Jangan sampai user A bisa ambil analisis AI dari submission user B.

> **Cache by submission_id:** Karena data submission tidak berubah setelah disimpan, analisis AI untuk satu `submission_id` akan selalu sama. Cache ini aman dan menghemat quota Groq.

> **Security:** `GROQ_API_KEY` tidak boleh pernah ada di frontend. Semua panggilan ke Groq harus melalui backend.

> **Latency:** Groq sangat cepat (~300-800ms). Tampilkan loading spinner di button "AI Recommendation" saat menunggu response, lalu render hasilnya.

> **Boleh DB terpisah:** Tabel `ai_result_logs` bisa diletakkan di database terpisah jika diinginkan. Hapus baris `FOREIGN KEY` dari SQL schema-nya dan handle relasi di application layer (cukup simpan `submission_id` sebagai string biasa).
