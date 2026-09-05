# PRD — Sistem Gamifikasi Atlas Food
> **Feature Branch**: `feature/gamification-ai-chat`
> **Repo**: `atlas_food_backend`
> **Status**: Rancangan Awal (Draft)
> **Author**: Tim Atlas Food
> **Tanggal**: September 2026

---

## 1. Latar Belakang & Tujuan

### 1.1 Masalah yang Diselesaikan

Aplikasi Atlas Food berperan sebagai platform pemetaan dan pengumpulan data gizi pangan terjangkau
berbasis crowdsourcing partisipatif. Tantangan utama platform crowdsourcing adalah menjaga **retensi
pengguna** agar mereka secara konsisten berkontribusi data survei, scan makanan, dan kolaborasi sesi.

Tanpa mekanisme penghargaan yang terukur, pengguna cenderung berhenti berkontribusi setelah sesi
pertama karena kurangnya **motivasi ekstrinsik** yang berkelanjutan.

### 1.2 Solusi

Menerapkan **Sistem Gamifikasi** berbasis kerangka kerja **PBL (Points, Badges, Leaderboard)** dari
Deterding et al. (2011) dan Hamari et al. (2014) — tanpa memerlukan API eksternal berbayar — yang
secara otomatis memberi penghargaan kepada pengguna yang aktif berkontribusi.

### 1.3 Tujuan Bisnis

1. Meningkatkan **retensi harian (Daily Active User)** pengguna Atlas Food.
2. Mendorong **kontribusi data gizi** yang lebih banyak dan berkualitas.
3. Menambahkan **dimensi kompetisi sosial** yang sehat antar kontributor gizi.
4. Memperkuat bobot akademis skripsi dengan penerapan teori motivasi pengguna (SDT — Self-Determination Theory).

---

## 2. Landasan Teori (Relevansi Skripsi)

| Teori / Konsep | Referensi | Kaitan dengan Fitur |
| :--- | :--- | :--- |
| **Gamifikasi (Gamification)** | Deterding et al., 2011 — *ACM MindTrek* | Dasar penggunaan elemen game (XP, Badge, Leaderboard) di konteks non-game |
| **Efektivitas PBL (Points, Badges, Leaderboard)** | Hamari et al., 2014 — *IEEE HICSS* | Justifikasi empiris bahwa PBL meningkatkan partisipasi crowdsourcing |
| **Gamifikasi pada Kesehatan & Gizi** | Johnson et al., 2016 — *Internet Interventions, Elsevier* | Kontekstualisasi gamifikasi pada domain nutrisi & pencatatan gizi |
| **Self-Determination Theory (SDT)** | Ryan & Deci, 2000 — *Contemporary Educational Psychology* | Pemenuhan kebutuhan Competence (level), Relatedness (leaderboard), Autonomy (pilihan aktivitas) |
| **Octalysis Framework** | Chou, Y.K., 2015 — *Actionable Gamification* | Core Drive 2 (Development), Core Drive 5 (Social Influence), Core Drive 8 (Loss Avoidance = Streak) |

---

## 3. Ruang Lingkup Fitur

### 3.1 Yang Termasuk (In-Scope)

- ✅ Atlas XP (Sistem Poin) — penghargaan per aktivitas
- ✅ Daily Streak — konsistensi harian
- ✅ Level & Rank System (Leveling berbasis XP)
- ✅ Achievement Badges (koleksi lencana digital)
- ✅ Leaderboard Global (Papan Peringkat)
- ✅ Riwayat XP per pengguna (XP History)
- ✅ API Endpoint gamifikasi (RESTful, JWT-Protected)
- ✅ Auto-trigger XP saat aktivitas survei, AI, dan kolaborasi

### 3.2 Yang Tidak Termasuk (Out-of-Scope) — untuk versi ini

- ❌ Voucher fisik / reward uang tunai
- ❌ Integrasi payment gateway
- ❌ Push notification (bisa dikembangkan di sprint berikutnya)

---

## 4. Detail Fitur & Mekanisme

### 4.1 🔥 Daily Streak

**Deskripsi**: Penghitung konsistensi harian pengguna. Streak bertambah +1 jika pengguna melakukan
minimal 1 aktivitas yang diakui dalam satu hari kalender (berdasarkan timezone WIB / UTC+7).

**Aktivitas yang Diakui**:
- Submit survei pangan (Submit 24HR Food Recall)
- Analisis makanan via AI (AI Nutrition Analysis)
- Bertanya via AI Chat (NutriBot AI)
- Bergabung di Room Kolaborasi (Collab Session)

**Aturan**:
- Jika pengguna tidak melakukan aktivitas apapun dalam 24 jam sejak aktivitas terakhir, **streak direset ke 0**.
- `max_streak` menyimpan rekor streak tertinggi yang pernah dicapai.
- Streak dihitung per hari kalender, bukan per 24 jam rolling.

**Tabel DB**:
```sql
-- Di tabel user_gamification
last_activity_date DATE -- tanggal aktivitas terakhir (WIB)
current_streak     INT  -- streak saat ini
max_streak         INT  -- rekor streak tertinggi
```

---

### 4.2 ⭐ Atlas XP & Level Rank

**Deskripsi**: Sistem poin kumulatif yang merepresentasikan total kontribusi pengguna.

**Tabel XP Per Aktivitas**:

| Aktivitas | XP Diberikan | Endpoint Trigger |
| :--- | :---: | :--- |
| Submit Survei Pangan (24HR Recall) | +50 XP | `POST /api/v1/submissions` |
| Analisis AI Nutrition Analysis | +30 XP | `POST /api/v1/ai/analyze` |
| Bertanya ke AI Chat (NutriBot) | +10 XP | `POST /api/v1/ai/chat` |
| Bergabung di Room Kolaborasi | +20 XP | WebSocket `JOIN_ROOM` event |
| Isi Survei (sebagai responden baru) | +25 XP | `POST /api/v1/surveys/:id/fill` |

**Tabel Level & Rank**:

| Level | Nama Rank | XP Minimum | XP Maksimum |
| :---: | :--- | :---: | :---: |
| 1 | 🥉 Nutri Novice | 0 | 150 |
| 2 | 🥈 Gizi Explorer | 151 | 400 |
| 3 | 🥇 Nutri Champion | 401 | 800 |
| 4 | 👑 Pahlawan Gizi Atlas | 801 | ∞ |

**Logic Level-Up**:
```go
func GetLevel(points int) (level int, rankName string) {
    switch {
    case points >= 801: return 4, "Pahlawan Gizi Atlas"
    case points >= 401: return 3, "Nutri Champion"
    case points >= 151: return 2, "Gizi Explorer"
    default:            return 1, "Nutri Novice"
    }
}
```

---

### 4.3 🏆 Leaderboard (Papan Peringkat)

**Deskripsi**: Klasemen pengguna berdasarkan total XP. Tersedia 3 mode filter.

**Filter Mode**:
- `weekly` — Top 10 pengguna dengan XP tertinggi 7 hari terakhir
- `monthly` — Top 10 pengguna dengan XP tertinggi 30 hari terakhir
- `all_time` — Top 10 pengguna dengan total XP tertinggi sepanjang masa

**Endpoint**:
```
GET /api/v1/gamification/leaderboard?period=weekly
GET /api/v1/gamification/leaderboard?period=monthly
GET /api/v1/gamification/leaderboard?period=all_time
```

**Response Schema**:
```json
{
  "period": "weekly",
  "leaderboard": [
    {
      "rank": 1,
      "user_id": "uuid",
      "name": "Budi Santoso",
      "avatar_url": "/uploads/avatars/xxx.jpg",
      "level": 3,
      "rank_name": "Nutri Champion",
      "total_xp": 780,
      "current_streak": 12
    }
  ],
  "my_position": {
    "rank": 5,
    "total_xp": 320
  }
}
```

---

### 4.4 🎖️ Achievement Badges

**Deskripsi**: Lencana digital yang dapat dikumpulkan pengguna sebagai bukti pencapaian.

**Katalog Badge**:

| ID Badge | Nama Badge | Ikon | Deskripsi Trigger |
| :--- | :--- | :---: | :--- |
| `FIRST_STEP` | First Step | 🔰 | Pertama kali submit survei pangan |
| `AI_ENTHUSIAST` | AI Enthusiast | 🤖 | Gunakan analisis AI sebanyak 5 kali |
| `STREAK_7` | Streak Master | 🔥 | Menjaga daily streak 7 hari berturut-turut |
| `STREAK_30` | Flame Legend | 🔥🔥 | Menjaga daily streak 30 hari berturut-turut |
| `COLLAB_HERO` | Collab Hero | 👥 | Bergabung di 5 sesi kolaborasi live |
| `NUTRI_CHAMPION` | Nutri Champion | 🥗 | Mencapai XP level 3 (401 XP) |
| `TOP_WEEKLY` | Juara Mingguan | 🏅 | Berada di peringkat #1 leaderboard mingguan |
| `NUTRIBOT_USER` | NutriBot User | 💬 | Pertama kali menggunakan AI Chat NutriBot |

**Sistem Trigger (Auto-award)**:
- Service `gamification` mengecek kondisi badge setiap kali XP ditambahkan.
- Badge yang sudah diraih tidak diberikan ulang (idempotent check via tabel `user_badges`).

---

## 5. Skema Database (Migrations)

### Migration 011: Tabel `user_gamification`
```sql
CREATE TABLE user_gamification (
    user_id           VARCHAR(36)  PRIMARY KEY,
    total_points      INT          NOT NULL DEFAULT 0,
    current_streak    INT          NOT NULL DEFAULT 0,
    max_streak        INT          NOT NULL DEFAULT 0,
    last_activity_date DATE,
    level             INT          NOT NULL DEFAULT 1,
    rank_name         VARCHAR(50)  NOT NULL DEFAULT 'Nutri Novice',
    created_at        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### Migration 012: Tabel `badges`
```sql
CREATE TABLE badges (
    id           VARCHAR(50)  PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    description  TEXT         NOT NULL,
    icon_emoji   VARCHAR(10),
    criteria     TEXT,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Migration 013: Tabel `user_badges`
```sql
CREATE TABLE user_badges (
    id          VARCHAR(36)  PRIMARY KEY DEFAULT (UUID()),
    user_id     VARCHAR(36)  NOT NULL,
    badge_id    VARCHAR(50)  NOT NULL,
    earned_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id)  REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (badge_id) REFERENCES badges(id) ON DELETE CASCADE,
    UNIQUE KEY uq_user_badge (user_id, badge_id)
);
```

### Migration 014: Tabel `xp_history`
```sql
CREATE TABLE xp_history (
    id          VARCHAR(36)  PRIMARY KEY DEFAULT (UUID()),
    user_id     VARCHAR(36)  NOT NULL,
    activity    VARCHAR(100) NOT NULL,
    xp_earned   INT          NOT NULL,
    description TEXT,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_xp_history_user (user_id),
    INDEX idx_xp_history_created_at (created_at)
);
```

---

## 6. API Endpoints Gamifikasi

### 6.1 Endpoint Daftar Lengkap

| Method | Endpoint | Auth | Deskripsi |
| :--- | :--- | :---: | :--- |
| `GET` | `/api/v1/gamification/profile` | ✅ JWT | Profil gamifikasi pengguna saat ini (XP, streak, level, badges) |
| `GET` | `/api/v1/gamification/leaderboard` | ✅ JWT | Papan peringkat (query: `?period=weekly\|monthly\|all_time`) |
| `GET` | `/api/v1/gamification/badges` | ✅ JWT | Semua badge yang diraih pengguna |
| `GET` | `/api/v1/gamification/xp-history` | ✅ JWT | Riwayat perolehan XP |

### 6.2 Response: `GET /api/v1/gamification/profile`
```json
{
  "user_id": "uuid-xxx",
  "name": "Bagas Pratama",
  "avatar_url": "/uploads/avatars/xxx.jpg",
  "total_points": 380,
  "level": 2,
  "rank_name": "Gizi Explorer",
  "next_level_points": 401,
  "progress_to_next_level": 95,
  "current_streak": 5,
  "max_streak": 12,
  "badges_earned": 3,
  "badges": [
    { "id": "FIRST_STEP", "name": "First Step", "icon_emoji": "🔰", "earned_at": "2026-09-01T10:00:00Z" },
    { "id": "AI_ENTHUSIAST", "name": "AI Enthusiast", "icon_emoji": "🤖", "earned_at": "2026-09-03T15:30:00Z" }
  ]
}
```

---

## 7. Arsitektur Internal Domain `gamification`

```
internal/domain/gamification/
├── model.go        -- struct UserGamification, Badge, UserBadge, XPHistory, Leaderboard
├── dto.go          -- request/response DTO untuk API layer
├── repository.go   -- interface Repository + MySQL implementation
├── service.go      -- business logic (AddXP, CheckStreak, AwardBadge, GetLeaderboard)
└── handler.go      -- Gin HTTP handler (GetProfile, GetLeaderboard, GetBadges)
```

### Alur `AddXP(userID, activity, xp)`:
```
1. Ambil / buat user_gamification record
2. Tambah total_points += xp
3. Recalculate level & rank_name
4. Update streak: cek last_activity_date vs today (WIB)
   - Jika sama hari → tidak ubah streak (sudah dihitung hari ini)
   - Jika kemarin → current_streak++, max_streak = max(max_streak, current_streak)
   - Jika > kemarin → current_streak = 1 (reset)
5. Simpan XP ke xp_history
6. Cek & award badge yang relevan (CheckAndAwardBadges)
7. Return hasil terbaru
```

---

## 8. Integration Points (Trigger Auto-XP)

XP ditambahkan secara internal — bukan via endpoint publik tersendiri — melainkan dipanggil dari
service lain yang sudah ada:

```go
// Di submission/service.go, setelah submit berhasil:
s.gamification.AddXP(userID, "submit_survey", 50)

// Di ai/service.go, setelah AnalyzeNutrition berhasil:
s.gamification.AddXP(userID, "ai_analysis", 30)

// Di ai/service.go, setelah Chat berhasil:
s.gamification.AddXP(userID, "ai_chat", 10)

// Di collab/handler.go, setelah JOIN_ROOM WebSocket:
gamificationSvc.AddXP(userID, "collab_join", 20)
```

---

## 9. Acceptance Criteria

- [ ] `GET /api/v1/gamification/profile` mengembalikan data XP, streak, level, dan badges milik user yang login.
- [ ] XP otomatis bertambah saat user submit survei (terintegrasi di domain `submission`).
- [ ] XP otomatis bertambah saat user pakai AI (terintegrasi di domain `ai`).
- [ ] Streak bertambah +1 pada hari yang berbeda dari aktivitas sebelumnya.
- [ ] Streak direset ke 0 jika tidak ada aktivitas lebih dari 1 hari.
- [ ] Badge `FIRST_STEP` diberikan otomatis setelah submit survei pertama.
- [ ] Badge `AI_ENTHUSIAST` diberikan setelah 5x analisis AI.
- [ ] Badge `STREAK_7` diberikan setelah streak mencapai 7 hari.
- [ ] Leaderboard `weekly` hanya menghitung XP yang diperoleh dalam 7 hari terakhir.
- [ ] `my_position` di leaderboard response menampilkan peringkat user yang sedang login.

---

## 10. Rencana Sprint Pengerjaan

| Sprint | Task |
| :--- | :--- |
| Sprint 1 | Buat 4 migration SQL (011-014) + seed data badges |
| Sprint 2 | Buat domain `gamification` (model, repo, service, handler) |
| Sprint 3 | Integrasi AddXP ke domain `submission` dan `ai` |
| Sprint 4 | Integrasi trigger streak & badge checker |
| Sprint 5 | Integrasi AddXP ke domain `collab` (WebSocket JOIN event) |
| Sprint 6 | Pengujian end-to-end + dokumentasi API (Postman Collection) |
