# 16 — Offline-First PWA dengan Auto-Sync

> **Status:** Planning
> **Target:** Atlas Food dapat digunakan sepenuhnya tanpa koneksi internet (input recall, browse makanan, simpan survey), dan otomatis sinkronisasi saat koneksi kembali aktif.
> **Tanggal:** 2026-09-01
> **Branch:** `feature/offline-pwa`

---

## 1. Vision & Scope

### 🎯 Problem Statement

Surveyor dan responden Atlas Food sering beroperasi di lapangan — di pedesaan, rumah responden, atau fasilitas kesehatan yang minim sinyal internet. Kondisi ini menyebabkan:

- **Data loss** saat koneksi putus di tengah pengisian survey recall
- **Survey tidak bisa dibuka** sama sekali karena halaman tidak ter-cache
- **Frustrasi pengguna** dan **data penelitian tidak akurat** karena responden tidak bisa menyelesaikan recall

### 🏆 Tujuan Fitur

| Target | Deskripsi |
|---|---|
| **Installable App** | Atlas Food bisa di-install di HP Android/iOS/Desktop seperti native app |
| **Full Offline Survey Input** | Semua halaman recall (pilih makanan, atur porsi, review) bisa diisi 100% offline |
| **Local Data Persistence** | Data tersimpan aman di device meskipun browser ditutup/HP restart |
| **Transparent Auto-Sync** | Saat koneksi online, data terkirim otomatis ke server tanpa aksi user |
| **Offline Food Database** | Subset database makanan populer di-cache lokal untuk pencarian offline |

### 🧑‍🤝‍🧑 User Stories

```
SURVEYOR LAPANGAN:
"Saya lagi wawancara responden di desa yang nggak ada sinyal.
 Saya tetap bisa isi data recall dengan lengkap, dan nanti
 begitu balik ke kota data langsung ke-sync ke server."

RESPONDEN MANDIRI:
"Saya mengisi survey recall sendiri di rumah, sinyal lemah.
 Aplikasinya tetap bisa dibuka dan saya bisa isi sampai selesai.
 Pas Wi-Fi nyala lagi, data langsung terkirim sendiri."

ADMIN:
"Data dari lapangan masuk semua meski surveyor offline tadi.
 Tidak ada data yang hilang, semuanya ter-record dengan benar."
```

---

## 2. Arsitektur High-Level

```
┌─────────────────────────────────────────────────────────────────────┐
│                        BROWSER / DEVICE                              │
│                                                                       │
│  ┌──────────────────────┐     ┌──────────────────────────────────┐  │
│  │    Next.js App        │     │          Service Worker           │  │
│  │  (React Components)   │     │    (offline-first cache layer)   │  │
│  │                        │     │                                   │  │
│  │  ┌──────────────────┐ │     │  ┌────────────────────────────┐  │  │
│  │  │  Survey Recall    │ │     │  │  Cache Storage (Workbox)   │  │  │
│  │  │  Pages & UI       │ │     │  │  ● Static Assets (JS/CSS)  │  │  │
│  │  └────────┬──────────┘ │     │  │  ● Next.js Pages (HTML)   │  │  │
│  │           │             │     │  │  ● Food DB Subset (JSON)  │  │  │
│  │  ┌────────▼──────────┐ │     │  │  ● Food Photos (img)      │  │  │
│  │  │  Offline Service   │ │     │  └────────────────────────────┘  │  │
│  │  │  (offlineService)  │◄──────►                                   │  │
│  │  └────────┬──────────┘ │     │  ┌────────────────────────────┐  │  │
│  │           │             │     │  │  Background Sync API       │  │  │
│  │  ┌────────▼──────────┐ │     │  │  (browser-native sync)     │  │  │
│  │  │   IndexedDB        │ │     │  └────────────────────────────┘  │  │
│  │  │ (Dexie.js ORM)     │ │     └──────────────────────────────────┘  │
│  │  │                    │ │                                             │
│  │  │ ● offline_queue    │ │                                             │
│  │  │ ● cached_foods     │ │                                             │
│  │  │ ● survey_drafts    │ │                                             │
│  │  └────────────────────┘ │                                             │
│  └──────────────────────────┘                                             │
└──────────────────┬────────────────────────────────────────────────────────┘
                   │ HTTPS / REST API (saat online)
                   │
┌──────────────────▼────────────────────────────────────────────────────────┐
│                        BACKEND (Go / Gin)                                  │
│                                                                             │
│  ┌────────────────────────────────────────────────────────────────────┐   │
│  │              Existing REST API                                       │   │
│  │  POST /api/v1/submissions    ← Menerima batch sync dari client      │   │
│  │  GET  /api/v1/foods          ← Seed awal food database offline      │   │
│  │  POST /api/v1/sync/batch     ← NEW: Batch submission endpoint       │   │
│  │  GET  /api/v1/sync/status    ← NEW: Cek status sync per device      │   │
│  └────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌────────────────────────────────────────────────────────────────────┐   │
│  │              MySQL + GORM                                            │   │
│  │  submissions table  ← target data sync dari offline queue           │   │
│  │  sync_logs table    ← NEW: audit trail sync history                │   │
│  └────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Komponen Utama

### 3.1 PWA Manifest & Service Worker (Frontend)

```
atlas_food_frontend/
├── public/
│   ├── manifest.json              ← PWA identity (name, icons, colors)
│   ├── sw.js                      ← Service Worker (via @serwist/next)
│   └── icons/
│       ├── icon-192x192.png
│       ├── icon-512x512.png
│       └── icon-maskable-512x512.png
├── next.config.ts                 ← withSerwist() wrapper
└── app/
    └── layout.tsx                 ← <head> manifest link + OfflineToast
```

### 3.2 Local Database Schema (IndexedDB via Dexie.js)

```typescript
// lib/offlineDb.ts

export interface OfflineQueueItem {
  id?: number;                   // Auto-increment (IndexedDB PK)
  localId: string;               // UUID v4 — generated di client sebelum sync
  type: 'SUBMISSION' | 'DRAFT';  // Tipe data
  surveyAccessToken: string;     // Untuk grouping per survey
  payload: RecallSubmissionDTO;  // Data recall lengkap
  createdAt: string;             // ISO 8601
  syncStatus: 'PENDING' | 'SYNCING' | 'SYNCED' | 'FAILED';
  retryCount: number;            // Berapa kali gagal sync
  lastAttemptAt?: string;        // ISO 8601 terakhir coba sync
  errorMessage?: string;         // Pesan error terakhir jika FAILED
}

export interface CachedFood {
  id: string;                    // Backend food UUID
  name: string;
  categoryCode: string;
  photoUrl?: string;
  portionSizeGram: number;
  cachedAt: string;              // Untuk TTL validation (7 days)
}

export interface SurveyDraft {
  accessToken: string;           // PK
  lastStep: string;              // halaman terakhir diisi ('add-food' | 'portion' | 'review')
  draftData: Partial<RecallSubmissionDTO>;
  updatedAt: string;
}
```

### 3.3 Online/Offline Status Hook

```typescript
// hooks/useOnlineStatus.ts

export function useOnlineStatus(): boolean {
  const [isOnline, setIsOnline] = useState<boolean>(
    typeof navigator !== 'undefined' ? navigator.onLine : true
  );

  useEffect(() => {
    const handleOnline  = () => setIsOnline(true);
    const handleOffline = () => setIsOnline(false);

    window.addEventListener('online',  handleOnline);
    window.addEventListener('offline', handleOffline);

    return () => {
      window.removeEventListener('online',  handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, []);

  return isOnline;
}
```

### 3.4 Offline Submission Service

```typescript
// domain/survey/services/offlineService.ts

export class OfflineSubmissionService {
  private db = getOfflineDb();

  // Simpan submission ke antrian lokal (offline/online sama-sama lewat sini)
  async enqueue(
    surveyAccessToken: string,
    payload: RecallSubmissionDTO
  ): Promise<string> {
    const localId = crypto.randomUUID();
    await this.db.offlineQueue.add({
      localId,
      type: 'SUBMISSION',
      surveyAccessToken,
      payload,
      createdAt: new Date().toISOString(),
      syncStatus: 'PENDING',
      retryCount: 0,
    });
    return localId;
  }

  // Ambil semua yang belum ter-sync
  async getPending(): Promise<OfflineQueueItem[]> {
    return this.db.offlineQueue
      .where('syncStatus').anyOf(['PENDING', 'FAILED'])
      .and(item => item.retryCount < 5) // Max 5 retry
      .toArray();
  }

  // Tandai sebagai berhasil sync
  async markSynced(localId: string, serverId: string): Promise<void> {
    await this.db.offlineQueue
      .where('localId').equals(localId)
      .modify({ syncStatus: 'SYNCED', serverId });
  }

  // Tandai sebagai gagal, increment retry count
  async markFailed(localId: string, error: string): Promise<void> {
    const item = await this.db.offlineQueue.where('localId').equals(localId).first();
    if (!item) return;
    await this.db.offlineQueue
      .where('localId').equals(localId)
      .modify({
        syncStatus: item.retryCount >= 4 ? 'FAILED' : 'PENDING',
        retryCount: item.retryCount + 1,
        lastAttemptAt: new Date().toISOString(),
        errorMessage: error,
      });
  }
}
```

### 3.5 Sync Engine

```typescript
// lib/syncEngine.ts

export class SyncEngine {
  private offlineService = new OfflineSubmissionService();
  private isSyncing = false;

  // Dipanggil saat event 'online' atau saat app mount
  async sync(): Promise<SyncResult> {
    if (this.isSyncing || !navigator.onLine) return { synced: 0, failed: 0 };
    this.isSyncing = true;

    const pending = await this.offlineService.getPending();
    let synced = 0;
    let failed = 0;

    for (const item of pending) {
      try {
        // Tandai sebagai sedang sync
        await this.offlineService.markSyncing(item.localId);

        // Kirim ke backend dengan header idempotency key
        const res = await apiClient.post('/submissions', {
          ...item.payload,
          localId: item.localId, // Backend pakai ini untuk idempotency
        }, {
          headers: { 'Idempotency-Key': item.localId }
        });

        await this.offlineService.markSynced(item.localId, res.data.id);
        synced++;
      } catch (err: any) {
        await this.offlineService.markFailed(item.localId, err.message);
        failed++;
      }
    }

    this.isSyncing = false;
    return { synced, failed };
  }
}

// Inisialisasi global — listen event 'online'
export function initSyncEngine() {
  const engine = new SyncEngine();

  // Sync saat pertama kali online
  if (navigator.onLine) engine.sync();

  // Sync otomatis saat koneksi kembali
  window.addEventListener('online', () => engine.sync());

  return engine;
}
```

---

## 4. Alur / Flow Mekanisme

### 4.1 Flow: Submit Recall Saat Offline

```
Responden klik "Simpan & Lanjut"
           │
           ▼
   useOnlineStatus() → false?
           │
    ┌──────┴──────┐
    │   OFFLINE   │
    └──────┬──────┘
           │
           ▼
  offlineService.enqueue(payload)
  → Simpan ke IndexedDB
  → localId = UUID.v4()
  → syncStatus = 'PENDING'
           │
           ▼
  Toast: "📥 Data tersimpan lokal.
          Akan dikirim saat ada koneksi."
           │
           ▼
  Navigasi ke halaman berikutnya
  (flow survey tetap jalan normal)
```

### 4.2 Flow: Auto-Sync Saat Online Kembali

```
Browser deteksi koneksi: event 'online'
           │
           ▼
  SyncEngine.sync() dipanggil
           │
           ▼
  Ambil semua item syncStatus == 'PENDING' dari IndexedDB
           │
    ┌──────┴──────────────────────────────┐
    │  Untuk setiap pending item:          │
    │                                      │
    │  1. Tandai syncStatus = 'SYNCING'    │
    │  2. POST /api/v1/submissions         │
    │     + Header: Idempotency-Key: {localId}
    │                                      │
    │  ┌────────────┬─────────────────┐    │
    │  │  Sukses ✅  │  Gagal ❌        │    │
    │  │            │                 │    │
    │  │ markSynced │ markFailed      │    │
    │  │ (SYNCED)   │ (retry < 5 →    │    │
    │  │            │  PENDING lagi)  │    │
    │  └────────────┴─────────────────┘    │
    └─────────────────────────────────────-┘
           │
           ▼
  Toast: "✅ {N} data berhasil disinkronisasi."
```

### 4.3 Flow: Caching Food Database untuk Offline Search

```
App pertama kali dibuka (online):
           │
           ▼
  GET /api/v1/foods?limit=500&popular=true
           │
           ▼
  Simpan ke IndexedDB.cachedFoods
  + TTL metadata: cachedAt = now()
           │
           ▼
  Saat user search makanan offline:
           │
           ▼
  foodSearchService.search(query)
    → isOnline? → false
    → db.cachedFoods.filter(name.includes(query))
    → Return cached results
    → Tampilkan banner "🔌 Menggunakan data offline"
```

### 4.4 Flow: PWA Installation

```
Pengguna buka Atlas Food di mobile browser (Chrome/Safari)
           │
           ▼
  Service Worker ter-install
  manifest.json terbaca
           │
           ▼
  Browser tampilkan: "Tambahkan ke layar utama?" → User tap "Tambah"
           │
           ▼
  App ter-install sebagai standalone (tanpa browser chrome/address bar)
           │
           ▼
  Semua asset statis & halaman survey di-cache via Workbox/Serwist
           │
           ▼
  Offline mode → App tetap bisa dibuka
```

---

## 5. Caching Strategy (Workbox / Serwist)

| Resource Type | Strategy | TTL | Deskripsi |
|---|---|---|---|
| **Static Assets** (JS, CSS, fonts) | `CacheFirst` | Seumur hidup (hash-based) | Next.js asset sudah pake content hash, aman CacheFirst |
| **HTML Pages** (app shell) | `StaleWhileRevalidate` | 1 hari | Tampilkan cache dulu, update di background |
| **API: GET /foods** | `NetworkFirst` fallback CacheStorage | 7 hari | Coba network, jika gagal pakai cache |
| **API: GET /surveys/{token}** | `NetworkFirst` fallback CacheStorage | 1 hari | Tampilkan data survey dari cache jika offline |
| **Food Photos** (as-served images) | `CacheFirst` | 30 hari | Gambar jarang berubah, cache agresif |
| **API: POST /submissions** | ❌ Jangan cache! Pakai IndexedDB Queue | — | Data write → selalu via IndexedDB offline queue |

---

## 6. Backend API Changes

### 6.1 Idempotency di Endpoint Submission

Endpoint `POST /api/v1/submissions` sudah ada. Perlu ditambah:

```go
// internal/domain/submission/handler.go

func (h *Handler) CreateSubmission(c *gin.Context) {
    // Cek Idempotency-Key dari header
    idempotencyKey := c.GetHeader("Idempotency-Key")

    if idempotencyKey != "" {
        // Cek apakah sudah ada submission dengan local_id ini
        existing, err := h.service.FindByLocalID(c.Request.Context(), idempotencyKey)
        if err == nil && existing != nil {
            // Sudah pernah diproses → return hasil lama (idempotent)
            c.JSON(http.StatusOK, existing)
            return
        }
    }

    // Proses submission baru...
    // Simpan idempotencyKey sebagai local_id di DB
}
```

### 6.2 Batch Sync Endpoint (NEW)

Untuk efisiensi: daripada N request terpisah, kirim semua sekaligus:

```
POST /api/v1/sync/batch

Request Body:
{
  "items": [
    {
      "local_id": "uuid-v4",
      "type": "SUBMISSION",
      "survey_access_token": "abc123",
      "payload": { ...RecallSubmissionDTO... },
      "created_at": "2026-09-01T10:00:00Z"
    },
    ...
  ]
}

Response:
{
  "results": [
    { "local_id": "uuid-v4", "status": "SYNCED",  "server_id": "server-uuid" },
    { "local_id": "uuid-v4", "status": "FAILED",  "error": "Survey not found" },
    { "local_id": "uuid-v4", "status": "SKIPPED", "reason": "Already synced" }
  ],
  "synced_count": 1,
  "failed_count": 1,
  "skipped_count": 1
}
```

### 6.3 Database Changes

```sql
-- Tambah kolom local_id untuk idempotency di tabel submissions
ALTER TABLE submissions
  ADD COLUMN local_id VARCHAR(36) NULL UNIQUE AFTER id,
  ADD INDEX idx_local_id (local_id);

-- Tabel audit trail sync (opsional tapi direkomendasikan)
CREATE TABLE sync_logs (
  id           CHAR(36)     PRIMARY KEY DEFAULT (UUID()),
  local_id     VARCHAR(36)  NOT NULL,
  user_id      CHAR(36)     NOT NULL,
  device_info  JSON         NULL,              -- user agent, dll
  sync_type    VARCHAR(50)  NOT NULL,          -- 'SUBMISSION', 'DRAFT'
  status       VARCHAR(20)  NOT NULL,          -- 'SYNCED', 'FAILED'
  error_detail TEXT         NULL,
  created_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_user (user_id),
  INDEX idx_status (status),
  INDEX idx_local_id (local_id)
);
```

### 6.4 File Structure Backend (NEW)

```
atlas_food_backend/internal/domain/submission/
├── handler.go        ← MODIFY: tambah handler BatchSync, idempotency check
├── service.go        ← MODIFY: tambah FindByLocalID(), BatchSync()
├── repository.go     ← MODIFY: tambah FindByLocalID(), CreateWithLocalID()
├── dto.go            ← MODIFY: tambah BatchSyncRequest/Response DTOs
└── ...existing files...

atlas_food_backend/internal/domain/sync/     ← NEW domain
├── handler.go        ← POST /api/v1/sync/batch
├── service.go
├── repository.go
└── dto.go
```

---

## 7. Frontend File Structure (NEW)

```
atlas_food_frontend/
├── public/
│   ├── manifest.json                       ← [NEW] PWA Manifest
│   └── icons/                              ← [NEW] App icons (192, 512, maskable)
│
├── next.config.ts                          ← [MODIFY] withSerwist() wrapper
│
├── internal/
│   ├── lib/
│   │   ├── offlineDb.ts                    ← [NEW] Dexie.js database setup
│   │   ├── syncEngine.ts                   ← [NEW] Auto-sync logic
│   │   └── offlineService.ts              ← [NEW] Queue management
│   │
│   ├── hooks/
│   │   ├── useOnlineStatus.ts             ← [NEW] navigator.onLine watcher
│   │   ├── useSyncStatus.ts               ← [NEW] pending count, sync state
│   │   └── useOfflineFood.ts              ← [NEW] food search w/ offline fallback
│   │
│   ├── components/
│   │   ├── OfflineStatusBar.tsx           ← [NEW] Banner "Kamu offline / N data menunggu sync"
│   │   ├── SyncProgressToast.tsx          ← [NEW] Toast saat sync berjalan
│   │   └── OfflineFoodSearch.tsx          ← [NEW] Search makanan dengan offline fallback
│   │
│   └── domain/
│       └── survey/
│           └── services/
│               └── offlineSubmission.ts   ← [NEW] Offline-first submission service
│
└── app/
    └── surveys/
        └── [accessToken]/
            └── layout.tsx                 ← [MODIFY] Tambah OfflineStatusBar
```

---

## 8. UI / UX Design

### 8.1 Offline Status Bar

```
┌─────────────────────────────────────────────────────────────────┐
│  🔌  Kamu sedang offline. Data akan tersimpan lokal.             │
│  3 data menunggu sinkronisasi.                [Sync Manual?]     │
└─────────────────────────────────────────────────────────────────┘
```

- **Warna bar**: merah/oranye saat offline, hilang saat online
- **Animasi**: fade-in saat pertama offline, shake subtle jika gagal sync

### 8.2 Tombol Submit (Offline State)

```
┌─────────────────────────────────────────────────┐
│  [   💾 Simpan Lokal (Offline Mode)   ]           │
│   Data akan dikirim otomatis saat online         │
└─────────────────────────────────────────────────┘
```

- Tombol tetap bisa diklik (tidak di-disable)
- Text berubah menjadi "Simpan Lokal" saat offline

### 8.3 Sync Progress Toast

```
✅ 3 data recall berhasil disinkronisasi ke server.
⚠️  1 data gagal dikirim — akan dicoba ulang otomatis.
```

### 8.4 Food Search Offline Indicator

```
┌────────────────────────────────────────────────┐
│  🔍 Cari makanan...                            │
│  [Mode Offline: Menampilkan data lokal]        │
│                                                 │
│  📦 Nasi Putih          250g / porsi           │
│  📦 Ayam Goreng         100g / porsi           │
│  ...                                            │
└────────────────────────────────────────────────┘
```

---

## 9. Dependencies Baru

### Frontend

```bash
# PWA / Service Worker
npm install @serwist/next serwist

# IndexedDB ORM (Dexie.js — lightweight, TypeScript-first)
npm install dexie

# UUID untuk localId generation
npm install uuid
npm install -D @types/uuid
```

### Backend

Tidak ada dependency baru. Perubahan hanya pada handler/service/repository layer yang sudah ada.

---

## 10. Best Practices & Prinsip Implementasi

### 10.1 Offline-First Principles

| Prinsip | Implementasi di Atlas Food |
|---|---|
| **Optimistic UI** | Tampilkan UI seolah request berhasil, sync di background |
| **Local-first writes** | Semua input survey → IndexedDB dulu, API menyusul |
| **Idempotent operations** | `localId` (UUID) sebagai Idempotency-Key → sync aman diulang |
| **Graceful degradation** | Online: full features; Offline: core features tetap berjalan |
| **Transparent sync** | User tidak perlu tahu kapan data di-sync — terjadi otomatis |

### 10.2 IndexedDB Best Practices

```typescript
// ✅ BENAR: Selalu pakai Dexie.js transaction untuk operasi multiple tables
await db.transaction('rw', db.offlineQueue, db.surveyDrafts, async () => {
  await db.offlineQueue.add(item);
  await db.surveyDrafts.put(draft);
});

// ❌ SALAH: Operasi terpisah tanpa transaction → bisa terjadi partial write
await db.offlineQueue.add(item);
await db.surveyDrafts.put(draft); // Jika ini gagal, terjadi inkonsistensi!
```

### 10.3 Service Worker Best Practices

```typescript
// ✅ BENAR: Pakai @serwist/next bukan manual sw.js untuk Next.js
// next.config.ts
import withSerwist from '@serwist/next';
export default withSerwist({
  swSrc: 'app/sw.ts',
  swDest: 'public/sw.js',
})(nextConfig);

// ✅ Cache busting: Serwist otomatis handle revision/versioning
// ✅ Skip waiting: aktifkan agar update SW langsung aktif
```

### 10.4 Sync Engine Best Practices

```typescript
// ✅ Exponential backoff untuk retry
const RETRY_DELAYS = [1000, 2000, 5000, 10000, 30000]; // ms
const delay = RETRY_DELAYS[item.retryCount] ?? 30000;

// ✅ Jangan sync jika tab/app sedang di-close
// Gunakan Background Sync API (via Service Worker) untuk sync yang lebih reliable
self.addEventListener('sync', (event) => {
  if (event.tag === 'sync-submissions') {
    event.waitUntil(syncAllPending());
  }
});

// ✅ Limit max retry = 5, setelah itu status = 'PERMANENTLY_FAILED'
// Tampilkan UI "Data gagal sync — hubungi admin" untuk kasus ini
```

### 10.5 Data Security Offline

| Aspek | Implementasi |
|---|---|
| **Data sensitif di IndexedDB** | IndexedDB menggunakan same-origin policy — hanya Atlas Food yang bisa akses |
| **Token di IndexedDB** | Jangan simpan JWT di IndexedDB! Pakai `httpOnly` cookie atau memory store |
| **Enkripsi data offline** | Untuk data sangat sensitif: enkripsi payload dengan AES-256 sebelum simpan ke IndexedDB (future scope) |
| **Bersihkan data setelah sync** | Setelah `SYNCED`, hapus payload dari IndexedDB (simpan hanya metadata untuk audit) |

---

## 11. Testing Plan

### 11.1 Manual Testing

| Skenario | Cara Test | Expected |
|---|---|---|
| Buka app offline | DevTools → Network → Offline | App tetap terbuka, halaman survey ter-cache |
| Submit recall offline | Network Offline → isi survey → klik simpan | Data masuk IndexedDB, toast "tersimpan lokal" |
| Auto-sync saat online kembali | DevTools → Online | Data terkirim ke backend, toast sukses |
| Duplikat sync (kirim 2x) | Paksa sync ulang item yang sudah synced | Backend deduplikasi via `local_id`, tidak ada duplikat data |
| Install PWA | Chrome → address bar → install icon | App ter-install, bisa dibuka standalone |
| Food search offline | Network Offline → search makanan | Tampilkan hasil dari cached foods di IndexedDB |

### 11.2 Automated Tests (Frontend)

```typescript
// __tests__/offlineService.test.ts

describe('OfflineSubmissionService', () => {
  it('should enqueue submission with PENDING status', async () => {
    const service = new OfflineSubmissionService();
    const localId = await service.enqueue('token123', mockPayload);
    const item = await service.findByLocalId(localId);
    expect(item?.syncStatus).toBe('PENDING');
  });

  it('should mark as SYNCED correctly', async () => {
    // ...
  });

  it('should not return items with retryCount >= 5', async () => {
    // ...
  });
});
```

### 11.3 Automated Tests (Backend)

```go
// internal/domain/submission/handler_test.go

func TestCreateSubmission_IdempotencyKey(t *testing.T) {
    // Kirim 2x request dengan Idempotency-Key yang sama
    // Pastikan response kedua identik dengan pertama
    // Pastikan tidak ada duplikat di database
}

func TestBatchSync(t *testing.T) {
    // Kirim batch berisi 3 item: 2 baru, 1 sudah ada (skip)
    // Pastikan hasil: 2 SYNCED, 1 SKIPPED
}
```

---

## 12. Implementation Phases

### Phase 1 — PWA Foundation (Backend: tidak ada perubahan)
- [x] Branch `feature/offline-pwa` dibuat (BE + FE)
- [ ] Tambah `@serwist/next` ke frontend
- [ ] Setup `manifest.json` + app icons
- [ ] Konfigurasi `withSerwist` di `next.config.ts`
- [ ] Verifikasi PWA install prompt muncul di browser

### Phase 2 — IndexedDB & Offline Service
- [ ] Setup Dexie.js schema (`offlineDb.ts`)
- [ ] Implementasi `OfflineSubmissionService`
- [ ] Implementasi `useOnlineStatus` hook
- [ ] Implementasi `SyncEngine` dasar (manual trigger)

### Phase 3 — Offline Submit Flow (FE)
- [ ] Modifikasi survey submission service → offline-first
- [ ] Tambah `OfflineStatusBar` component ke layout survey
- [ ] Modifikasi tombol simpan → "Simpan Lokal" saat offline
- [ ] `SyncProgressToast` component

### Phase 4 — Auto-Sync & Backend Idempotency
- [ ] Tambah `localId` kolom di tabel `submissions`
- [ ] Modifikasi `POST /api/v1/submissions` → support Idempotency-Key
- [ ] Implementasi `POST /api/v1/sync/batch` endpoint
- [ ] Implementasi auto-sync via event `'online'`
- [ ] Background Sync API via Service Worker

### Phase 5 — Offline Food Search
- [ ] Implementasi food caching ke IndexedDB saat online
- [ ] Modifikasi food search service → offline fallback
- [ ] `OfflineFoodSearch` component dengan indicator mode offline
- [ ] TTL management untuk cached foods (7 hari)

### Phase 6 — Polish & Testing
- [ ] Unit tests frontend (Jest)
- [ ] Integration tests backend (Go test)
- [ ] Manual testing seluruh offline scenario
- [ ] Performance audit (Lighthouse PWA score target: 90+)

---

## 13. Lighthouse PWA Target Score

| Kategori | Target |
|---|---|
| **Performance** | ≥ 90 |
| **Accessibility** | ≥ 95 |
| **Best Practices** | ≥ 90 |
| **SEO** | ≥ 90 |
| **PWA** | ✅ All checks pass |

Checklist PWA Lighthouse:
- ✅ Served over HTTPS
- ✅ Has a `<meta name="viewport">` tag
- ✅ Has valid `manifest.json`
- ✅ Service Worker registered
- ✅ App responds with 200 when offline
- ✅ Has icons for various sizes
- ✅ `start_url` responds while offline

---

## 14. References

- [Web.dev — Offline Cookbook](https://web.dev/offline-cookbook/)
- [Serwist — Next.js PWA (fork of next-pwa)](https://serwist.pages.dev/docs/next)
- [Dexie.js — Docs](https://dexie.org/docs/API-Reference)
- [Background Sync API](https://developer.chrome.com/docs/workbox/modules/workbox-background-sync/)
- [IndexedDB Best Practices — MDN](https://developer.mozilla.org/en-US/docs/Web/API/IndexedDB_API/Using_IndexedDB)
- [PWA Manifest — MDN](https://developer.mozilla.org/en-US/docs/Web/Manifest)
