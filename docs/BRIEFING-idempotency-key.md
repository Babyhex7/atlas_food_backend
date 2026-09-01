# BRIEFING: Implementasi Idempotency Key — Atlas Food (BE & FE)

> **Status:** Ready to Implement
> **Tanggal:** 2026-09-01
> **Branch:** `feature/offline-pwa`
> **Referensi:** `docs/16-offline-pwa-sync.md`

---

## 1. Apa itu Idempotency Key & Kenapa Penting?

**Idempotency Key** adalah mekanisme keamanan data yang memastikan:

> "Jika request yang **sama** dikirim **lebih dari satu kali** (karena koneksi putus, retry, dll), data hanya tersimpan **satu kali** di database."

### Kenapa Atlas Food Butuh Ini?

| Skenario Nyata | Tanpa Idempotency | Dengan Idempotency |
|---|---|---|
| Surveyor klik "Simpan" → sinyal putus → klik lagi | **2 record duplikat** di DB ❌ | Hanya **1 record** tersimpan ✅ |
| Auto-sync offline queue kirim data saat online | **Duplikat submissions** jika sync dipanggil 2x ❌ | Aman, backend skip jika sudah ada ✅ |
| User double-click tombol submit karena loading lama | **Duplikat** ❌ | 1 record, idempoten ✅ |
| Network flaky → request timeout → client retry | **Duplikat** ❌ | Aman ✅ |

---

## 2. Overview Arsitektur Idempotency

```
Frontend (Next.js)                    Backend (Go / Gin)
      │                                       │
      │  1. Generate localId = UUID()         │
      │                                       │
      │  POST /api/v1/survey/submit           │
      │  Headers:                             │
      │    Idempotency-Key: {localId}   ───►  │  2. Cek header Idempotency-Key
      │  Body:                                │
      │    { ...payload, local_id: localId }  │  3. Cari di DB: local_id = localId?
      │                                       │
      │                                       │  ┌──────────────────┬─────────────────────┐
      │                                       │  │  DITEMUKAN       │  TIDAK DITEMUKAN     │
      │                                       │  │                  │                      │
      │                                       │  │  Return 200 OK   │  Simpan ke DB        │
      │                                       │  │  (data lama)     │  Return 201 Created  │
      │                                       │  └──────────────────┴─────────────────────┘
      │  ◄─── Response (201 atau 200) ────────│
      │                                       │
```

---

## 3. Cek Kondisi Kode Saat Ini

### 3.1 Backend — Apa yang Belum Ada?

| File | Kondisi Sekarang | Yang Perlu Ditambah |
|---|---|---|
| [`submission/model.go`](file:///c:/Users/mybook_bagas/Projek_Agency/atlas_food/atlas_food_backend/internal/domain/submission/model.go) | ❌ Tidak ada field `LocalID` | ✅ Tambah field `LocalID *string` |
| [`submission/repository.go`](file:///c:/Users/mybook_bagas/Projek_Agency/atlas_food/atlas_food_backend/internal/domain/submission/repository.go) | ❌ Tidak ada `FindByLocalID()` | ✅ Tambah method `FindByLocalID(localID string)` |
| [`submission/service.go`](file:///c:/Users/mybook_bagas/Projek_Agency/atlas_food/atlas_food_backend/internal/domain/submission/service.go) | ❌ Tidak ada idempotency check | ✅ Cek `localID` sebelum `CreateSubmission` |
| [`submission/handler.go`](file:///c:/Users/mybook_bagas/Projek_Agency/atlas_food/atlas_food_backend/internal/domain/submission/handler.go) | ❌ Tidak baca header `Idempotency-Key` | ✅ Extract header → pass ke service |
| [`submission/dto.go`](file:///c:/Users/mybook_bagas/Projek_Agency/atlas_food/atlas_food_backend/internal/domain/submission/dto.go) | ❌ Tidak ada field `LocalID` di request | ✅ Tambah `LocalID string` di `SubmitSurveyRequest` |
| Database `survey_submissions` | ❌ Tidak ada kolom `local_id` | ✅ Migration: tambah kolom `local_id VARCHAR(36) UNIQUE` |

### 3.2 Frontend — Apa yang Belum Ada?

| File | Kondisi Sekarang | Yang Perlu Ditambah |
|---|---|---|
| [`submissionService.ts`](file:///c:/Users/mybook_bagas/Projek_Agency/atlas_food/atlas_food_frontend/atlas_food_frontend/internal/domain/submission/services/submissionService.ts) | ❌ Tidak generate UUID, tidak ada header `Idempotency-Key` | ✅ Generate `localId`, kirim header + field `local_id` |
| [`axios.ts`](file:///c:/Users/mybook_bagas/Projek_Agency/atlas_food/atlas_food_frontend/atlas_food_frontend/internal/lib/axios.ts) | ❌ Tidak ada global interceptor idempotency | ✅ (opsional) Tambah interceptor untuk POST/PUT requests |
| `submission/types/submission.ts` | ❌ Type `CreateSubmissionRequest` tidak ada `local_id` | ✅ Tambah `local_id?: string` ke interface |

---

## 4. Perubahan yang Diperlukan

---

### 4.1 DATABASE MIGRATION (Wajib Pertama)

```sql
-- Migration: Tambah kolom local_id untuk idempotency
ALTER TABLE survey_submissions
  ADD COLUMN local_id VARCHAR(36) NULL UNIQUE COMMENT 'UUID dari client, untuk idempotency offline sync'
  AFTER id;

-- Index sudah otomatis dari UNIQUE constraint
-- Verifikasi:
SHOW COLUMNS FROM survey_submissions;
```

**Note:** Kolom `NULL` karena submission lama (sebelum fitur ini) tidak punya `local_id`.

---

### 4.2 BACKEND — File per File

---

#### 📄 `internal/domain/submission/model.go`

**MODIFY:** Tambah field `LocalID` ke struct `SurveySubmission`.

```go
// SEBELUM:
type SurveySubmission struct {
    ID              string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
    SurveyID        string    `gorm:"type:char(36);not null;index" json:"survey_id"`
    // ...
}

// SESUDAH:
type SurveySubmission struct {
    ID              string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
    LocalID         *string   `gorm:"type:varchar(36);uniqueIndex;default:null" json:"local_id,omitempty"` // ← TAMBAH INI
    SurveyID        string    `gorm:"type:char(36);not null;index" json:"survey_id"`
    // ... (field lain tidak berubah)
}
```

---

#### 📄 `internal/domain/submission/dto.go`

**MODIFY:** Tambah field `LocalID` di `SubmitSurveyRequest`.

```go
// SEBELUM:
type SubmitSurveyRequest struct {
    SurveyID        string            `json:"survey_id" binding:"required"`
    ParticipantID   string            `json:"participant_id"`
    // ...
}

// SESUDAH:
type SubmitSurveyRequest struct {
    LocalID         string            `json:"local_id"`             // ← TAMBAH INI (opsional, dari client)
    SurveyID        string            `json:"survey_id" binding:"required"`
    ParticipantID   string            `json:"participant_id"`
    // ... (field lain tidak berubah)
}
```

**TAMBAH** DTO baru untuk response idempotent:

```go
// SubmitSurveyIdempotentResponse — digunakan ketika localId sudah pernah diproses
// Backend tidak proses ulang, langsung return data lama
type SubmitSurveyIdempotentResponse struct {
    SubmissionID string `json:"submission_id"`
    Message      string `json:"message"`
    Idempotent   bool   `json:"idempotent"` // true = data ini hasil cache, bukan baru disimpan
}
```

---

#### 📄 `internal/domain/submission/repository.go`

**MODIFY:** Tambah method `FindByLocalID` ke interface dan implementasi.

```go
// Interface — MODIFY:
type Repository interface {
    CreateSubmission(submission *SurveySubmission) error
    GetSubmissionByID(id string) (*SurveySubmission, error)
    FindByLocalID(localID string) (*SurveySubmission, error) // ← TAMBAH INI
    ListSubmissionsBySurvey(surveyID string, page, limit int) ([]SurveySubmission, int64, error)
    ListSubmissionsByUserID(userID, userEmail string, page, limit int) ([]SurveySubmission, int64, error)
    GetSubmissionByIDAndUser(id, userID, userEmail string) (*SurveySubmission, error)
    DeleteSubmission(id string) error
}

// Implementasi — TAMBAH:
// FindByLocalID - cari submission berdasarkan local_id dari client (untuk idempotency check)
func (r *submissionRepository) FindByLocalID(localID string) (*SurveySubmission, error) {
    var submission SurveySubmission
    err := r.db.Where("local_id = ?", localID).First(&submission).Error
    if err != nil {
        return nil, err
    }
    return &submission, nil
}
```

---

#### 📄 `internal/domain/submission/service.go`

**MODIFY:** Tambah idempotency check di awal fungsi `SubmitSurvey`.

```go
// Service interface — MODIFY:
type Service interface {
    SubmitSurvey(req SubmitSurveyRequest, userID string, idempotencyKey string) (*SubmissionResponse, error) // ← tambah param idempotencyKey
    // ... (method lain tidak berubah)
}

// Implementasi SubmitSurvey — MODIFY:
func (s *submissionService) SubmitSurvey(req SubmitSurveyRequest, userID string, idempotencyKey string) (*SubmissionResponse, error) {
    
    // ==============================
    // IDEMPOTENCY CHECK (TAMBAH INI)
    // ==============================
    if idempotencyKey != "" {
        existing, err := s.repo.FindByLocalID(idempotencyKey)
        if err == nil && existing != nil {
            // Submission dengan local_id ini sudah pernah diproses!
            // Return data lama tanpa menyimpan ulang.
            return &SubmissionResponse{
                SubmissionID: existing.ID,
                Message:      "Survey sudah berhasil dikirim sebelumnya (idempotent).",
            }, nil
        }
    }
    // ==============================

    // ... (sisa logic SubmitSurvey tidak berubah)

    // MODIFY: set LocalID sebelum CreateSubmission
    if idempotencyKey != "" {
        submission.LocalID = &idempotencyKey
    }

    // Simpan ke database
    if err := s.repo.CreateSubmission(submission); err != nil {
        return nil, errors.New("gagal menyimpan hasil survey")
    }

    return &SubmissionResponse{
        SubmissionID: submission.ID,
        Message:      "Survey berhasil dikirim, terima kasih!",
    }, nil
}
```

---

#### 📄 `internal/domain/submission/handler.go`

**MODIFY:** Extract header `Idempotency-Key` dan pass ke service.

```go
// SubmitSurvey handler — MODIFY:
func (h *Handler) SubmitSurvey(c *gin.Context) {
    var req SubmitSurveyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ValidationErrorResponse(c, "Data survey tidak lengkap: "+err.Error())
        return
    }

    // ==============================
    // EXTRACT IDEMPOTENCY KEY (TAMBAH INI)
    // ==============================
    // Prioritas: Header > Body field local_id
    idempotencyKey := c.GetHeader("Idempotency-Key")
    if idempotencyKey == "" {
        idempotencyKey = req.LocalID // Fallback: dari body JSON
    }
    // ==============================

    if email, exists := c.Get("email"); exists {
        if req.RespondentEmail == "" {
            req.RespondentEmail = email.(string)
        }
    }

    userID, _ := c.Get("userID")
    userIDStr, _ := userID.(string)

    // MODIFY: tambah parameter idempotencyKey
    response, err := h.service.SubmitSurvey(req, userIDStr, idempotencyKey)
    if err != nil {
        // ... (error handling tidak berubah)
    }

    utils.CreatedResponse(c, response)
}
```

---

### 4.3 FRONTEND — File per File

---

#### 📄 `internal/domain/submission/types/submission.ts`

**MODIFY:** Tambah `local_id` ke interface request.

```typescript
// SEBELUM:
export interface CreateSubmissionRequest {
  survey_id: string;
  participant_id?: string;
  respondent_name?: string;
  respondent_email?: string;
  meals_data: MealData[];
  daily_total: DailyTotal;
  missing_foods?: MissingFoodData[];
}

// SESUDAH:
export interface CreateSubmissionRequest {
  local_id?: string;              // ← TAMBAH INI — UUID dari client untuk idempotency
  survey_id: string;
  participant_id?: string;
  respondent_name?: string;
  respondent_email?: string;
  meals_data: MealData[];
  daily_total: DailyTotal;
  missing_foods?: MissingFoodData[];
}
```

---

#### 📄 `internal/domain/submission/services/submissionService.ts`

**MODIFY:** Generate `localId` dan sertakan header + field ke request.

```typescript
import { v4 as uuidv4 } from "uuid"; // npm install uuid

// SEBELUM:
export async function submitSurvey(
  payload: CreateSubmissionRequest
): Promise<{ submission_id: string; message: string }> {
  const response = await axiosClient.post("/survey/submit", payload);
  return response.data.data;
}

// SESUDAH:
export async function submitSurvey(
  payload: CreateSubmissionRequest,
  localId?: string  // ← Opsional: jika dari offline queue, pakai localId yang sudah ada
): Promise<{ submission_id: string; message: string }> {
  
  // Generate localId baru jika tidak disediakan (submit online biasa)
  const idempotencyKey = localId ?? uuidv4();

  const response = await axiosClient.post("/survey/submit", {
    ...payload,
    local_id: idempotencyKey,       // ← Sertakan di body
  }, {
    headers: {
      "Idempotency-Key": idempotencyKey,  // ← Sertakan di header
    },
  });

  return response.data.data;
}
```

---

## 5. Alur Lengkap dengan Idempotency

### Skenario A: Submit Online Biasa (Sekali Klik)

```
User klik "Simpan Survey"
        │
        ▼
submitSurvey(payload)
  → Generate localId = uuidv4()      // "abc-123-xyz"
  → POST /survey/submit
      Header: Idempotency-Key: "abc-123-xyz"
      Body: { ...payload, local_id: "abc-123-xyz" }
        │
        ▼ Backend
  Cek: SELECT * FROM survey_submissions WHERE local_id = "abc-123-xyz"
  → Tidak ditemukan
  → Simpan record baru, local_id = "abc-123-xyz"
  → Return 201 Created { submission_id: "server-uuid", idempotent: false }
        │
        ▼
  Toast: "✅ Survey berhasil dikirim!"
```

### Skenario B: Double-Click / Retry (Duplikat dikirim)

```
User double-klik "Simpan" (atau sinyal putus → retry)
        │
        ▼
Request 1 berhasil disimpan (local_id = "abc-123-xyz")
        │
        ▼
Request 2 dikirim dengan Idempotency-Key yang SAMA
  → Backend cek: local_id = "abc-123-xyz" → DITEMUKAN!
  → Return 200 OK { submission_id: "server-uuid", idempotent: true }
        │
        ▼
  Frontend: Proses seperti berhasil biasa
  → Tidak ada duplikat di database! ✅
```

### Skenario C: Offline Mode → Auto-Sync

```
User offline → submit survey
  → offlineService.enqueue(payload)
  → localId = uuidv4()  ← Dibuat saat pertama kali simpan offline
  → Simpan ke IndexedDB: { localId: "abc-123-xyz", payload, status: "PENDING" }

Saat online kembali:
  → syncEngine.sync() jalan
  → Ambil item PENDING dari IndexedDB
  → submitSurvey(payload, "abc-123-xyz")  ← localId yang sudah ada dipakai!
      Header: Idempotency-Key: "abc-123-xyz"

Jika sync dipanggil 2x (bug / race condition):
  → Request kedua di-reject oleh backend (idempotent)
  → Tidak ada duplikat! ✅
```

---

## 6. Checklist Implementasi

### Backend

- [ ] **Migration:** Jalankan SQL alter table tambah kolom `local_id`
- [ ] **`model.go`:** Tambah field `LocalID *string` ke struct `SurveySubmission`
- [ ] **`dto.go`:** Tambah field `LocalID string` ke `SubmitSurveyRequest`
- [ ] **`repository.go`:** Tambah method `FindByLocalID(localID string)` ke interface + implementasi
- [ ] **`service.go`:** Tambah parameter `idempotencyKey` + idempotency check di awal `SubmitSurvey`
- [ ] **`handler.go`:** Extract header `Idempotency-Key` → pass ke service
- [ ] **Test:** Pastikan kirim request yang sama 2x → hanya 1 record di DB

### Frontend

- [ ] **Install dependency:** `npm install uuid && npm install -D @types/uuid`
- [ ] **`submission.ts`:** Tambah field `local_id?: string` ke `CreateSubmissionRequest`
- [ ] **`submissionService.ts`:** Generate UUID, kirim header + field `local_id` di `submitSurvey()`
- [ ] **Test:** Submit survey → cek request di DevTools → verifikasi header `Idempotency-Key` ada

---

## 7. Testing Plan

### Manual Test

| Skenario | Langkah | Expected Result |
|---|---|---|
| Submit normal | Isi survey → klik simpan 1x | 1 record di DB, HTTP 201 |
| Double-click submit | Isi survey → klik simpan 2x cepat | 1 record di DB, request ke-2 HTTP 200 |
| Retry setelah timeout | Paksa timeout (throttle network) → retry | 1 record di DB |
| Submit offline → sync | Offline → submit → online | 1 record di DB setelah sync |
| Sync dipanggil 2x | Jalankan syncEngine.sync() 2x berturut-turut | 1 record di DB, tidak ada duplikat |

### Backend Unit Test

```go
// internal/domain/submission/service_test.go

func TestSubmitSurvey_IdempotencyKey(t *testing.T) {
    // Setup: buat mock repo
    mockRepo := &MockRepository{}
    
    // Simulasi: local_id sudah ada di DB
    existingSubmission := &SurveySubmission{ID: "existing-uuid", LocalID: ptr("test-local-id")}
    mockRepo.On("FindByLocalID", "test-local-id").Return(existingSubmission, nil)
    
    // Act
    result, err := service.SubmitSurvey(req, "user-id", "test-local-id")
    
    // Assert: return data lama, tidak panggil CreateSubmission
    assert.NoError(t, err)
    assert.Equal(t, "existing-uuid", result.SubmissionID)
    mockRepo.AssertNotCalled(t, "CreateSubmission")
}
```

---

## 8. Catatan Penting

> [!WARNING]
> **Jangan** generate `localId` baru setiap kali retry. `localId` harus **sama** untuk request yang sama (sekali generate, pakai terus untuk retry-nya).

> [!IMPORTANT]
> **Urutan implementasi yang benar:**
> 1. Jalankan migration SQL dulu (tambah kolom `local_id`)
> 2. Deploy backend (model + repo + service + handler)
> 3. Deploy frontend (generate UUID + kirim header)
>
> Jika frontend di-deploy duluan sebelum backend siap, tidak ada masalah — backend lama hanya akan mengabaikan field `local_id` yang tidak dikenal.

> [!TIP]
> Untuk debug, cek di DevTools → Network → Request Headers. Harus ada:
> ```
> Idempotency-Key: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
> ```
