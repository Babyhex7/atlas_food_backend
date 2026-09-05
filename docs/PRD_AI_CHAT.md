# PRD — AI Chat NutriBot Atlas Food
> **Feature Branch**: `feature/gamification-ai-chat`
> **Repo**: `atlas_food_backend` & `atlas_food_frontend`
> **Status**: Rancangan Awal (Draft)
> **Author**: Tim Atlas Food
> **Tanggal**: September 2026

---

## 1. Latar Belakang & Tujuan

### 1.1 Masalah yang Diselesaikan

Pengguna Atlas Food — terutama mahasiswa yang mengisi food recall — seringkali memiliki pertanyaan
sederhana seputar gizi dan makanan sehari-hari yang tidak terjawab langsung di dalam aplikasi:
- *"Berapa kalori nasi padang porsi kecil?"*
- *"Apa pengganti protein susu yang murah?"*
- *"Apakah makan mi instan setiap hari itu bahaya?"*

Saat ini, aplikasi hanya menyediakan hasil analisis gizi pasca-submit (AI Nutrition Analysis).
Tidak ada saluran percakapan interaktif yang kontekstual.

### 1.2 Solusi

Membangun **NutriBot AI** — asisten gizi berbasis Large Language Model (LLM) yang di-prompt secara
spesifik untuk konteks **gizi pangan terjangkau** dan **pola makan mahasiswa Indonesia** —
yang bisa diakses kapan saja via Floating Chat Widget atau halaman penuh `/ai-chat`.

### 1.3 Backend AI Engine

AI Chat memanfaatkan **Groq API dengan model LLM** yang sudah terintegrasi di domain `ai` yang
ada (`internal/pkg/groq`). **Tidak ada tambahan dependensi / API luar baru**.

---

## 2. Landasan Teori (Relevansi Skripsi)

| Konsep | Referensi | Kaitan dengan Fitur |
| :--- | :--- | :--- |
| **Conversational AI / Chatbot** | Adamopoulou & Moussiades, 2020 — *IFIP AICT* | Dasar teori sistem chatbot berbasis NLP |
| **AI-Assisted Nutrition Counseling** | Yang et al., 2021 — *JMIR mHealth* | Chatbot gizi terbukti meningkatkan literasi nutrisi pengguna |
| **LLM Prompt Engineering** | Brown et al., 2020 — *NeurIPS (GPT-3 Paper)* | Dasar teknik system prompt untuk mengontrol domain jawaban LLM |
| **Human-AI Interaction** | Luger & Sellen, 2016 — *ACM CHI* | Ekspektasi pengguna terhadap asisten AI percakapan |

---

## 3. Ruang Lingkup Fitur

### 3.1 Yang Termasuk (In-Scope)

- ✅ Floating Chat Widget (mengambang di pojok kanan bawah, global di semua halaman)
- ✅ Halaman penuh `/ai-chat` dengan riwayat percakapan persisten
- ✅ Tombol "Expand ke Layar Penuh" dari widget floating
- ✅ Quick Prompt Suggestions (4-6 pertanyaan populer sekali klik)
- ✅ Riwayat percakapan disimpan ke DB (`ai_chat_histories`)
- ✅ Integrasi Gamifikasi: setiap chat memberi pengguna +10 XP
- ✅ System prompt yang dikustomisasi untuk konteks gizi Atlas Food
- ✅ UI typing indicator (animasi dots saat AI sedang membalas)

### 3.2 Yang Tidak Termasuk (Out-of-Scope)

- ❌ Multi-turn memory di luar satu sesi percakapan (tidak diingat antar sesi)
- ❌ Analisis gambar/foto makanan via chat (sudah di fitur AI Nutrition Analysis)
- ❌ Voice input/output

---

## 4. Detail Fitur

### 4.1 💬 Floating Chat Widget (UI)

**Posisi**: Fixed di pojok kanan bawah layar (`bottom: 24px`, `right: 24px`).

**Tampilan**:
```
┌──────────────────────┐
│   🤖 NutriBot AI     │ ← Header dengan tombol minimize & expand
│──────────────────────│
│  Halo! Saya Atlas    │
│  NutriBot. Tanyakan  │
│  apapun soal gizi!   │
│──────────────────────│
│  💡 Quick Prompts:   │
│  ┌──────────────────┐│
│  │ Makanan murah... ││
│  └──────────────────┘│
│  ┌──────────────────┐│
│  │ Kalori nasi...   ││
│  └──────────────────┘│
│──────────────────────│
│  Ketik pertanyaan... │ ← Input bar + Send button
└──────────────────────┘
```

**State Widget**:
- **Collapsed**: Tombol bulat `🤖` floating (FAB — Floating Action Button)
- **Expanded**: Pop-up chat window (lebar 360px, tinggi 520px, responsive)
- **Full Page**: Redirect ke `/ai-chat` (window tetap ada, history tersinkron)

---

### 4.2 📄 Halaman Penuh `/ai-chat`

**Fitur Halaman**:
- Riwayat percakapan lengkap dari semua sesi sebelumnya (tersimpan di DB).
- Input area yang lebih besar.
- Panel kiri (desktop): Daftar sesi percakapan sebelumnya (sidebar history).
- Indikator XP yang diperoleh dari sesi chat (+10 XP per sesi chat baru).

---

### 4.3 ⚡ Quick Prompt Suggestions

**6 Prompt Siap Pakai**:
1. 💡 *"Rekomendasi makanan sehat & murah di bawah Rp 15.000"*
2. 🥗 *"Berapa kalori dan protein Nasi Goreng Telur satu porsi?"*
3. 🏋️ *"Tips pemenuhan kebutuhan protein harian untuk mahasiswa aktif"*
4. 🩺 *"Apakah konsumsi mi instan setiap hari itu berbahaya?"*
5. 📊 *"Apa itu skor gizi di Atlas Food dan bagaimana cara membacanya?"*
6. 🌱 *"Makanan tinggi zat besi yang cocok untuk sarapan murah"*

---

### 4.4 🤖 System Prompt NutriBot AI

**System prompt yang dikirim ke Groq setiap sesi**:

```
Kamu adalah NutriBot, asisten gizi cerdas dari aplikasi Atlas Food — sebuah platform
pemetaan pangan dan gizi terjangkau untuk mahasiswa dan masyarakat umum Indonesia.

Peranmu:
1. Menjawab pertanyaan seputar gizi, nutrisi, kalori, dan komposisi makanan Indonesia
   secara akurat, ramah, dan mudah dipahami.
2. Memberikan rekomendasi makanan bergizi dengan harga terjangkau (di bawah Rp 20.000/porsi)
   yang umum ditemui di sekitar kampus atau warung makan lokal.
3. Membantu pengguna memahami hasil food recall dan skor gizi dari Atlas Food.

Aturan ketat:
- Jawab HANYA pertanyaan terkait gizi, nutrisi, makanan, dan kesehatan umum.
- Jika ditanya di luar topik gizi/kesehatan, tolak dengan ramah dan arahkan kembali ke topik.
- Jangan memberikan diagnosis medis. Selalu sarankan konsultasi dokter untuk kondisi medis.
- Gunakan Bahasa Indonesia yang santai, ramah, dan mudah dipahami mahasiswa.
- Sertakan angka/estimasi gizi jika relevan.
```

---

## 5. Skema Database

### Migration 015: Tabel `ai_chat_sessions`
```sql
CREATE TABLE ai_chat_sessions (
    id          VARCHAR(36)  PRIMARY KEY DEFAULT (UUID()),
    user_id     VARCHAR(36)  NOT NULL,
    title       VARCHAR(200),
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_chat_sessions_user (user_id)
);
```

### Migration 016: Tabel `ai_chat_messages`
```sql
CREATE TABLE ai_chat_messages (
    id           VARCHAR(36)  PRIMARY KEY DEFAULT (UUID()),
    session_id   VARCHAR(36)  NOT NULL,
    role         ENUM('user', 'assistant') NOT NULL,
    content      TEXT         NOT NULL,
    token_used   INT          DEFAULT 0,
    model_used   VARCHAR(100),
    latency_ms   INT          DEFAULT 0,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES ai_chat_sessions(id) ON DELETE CASCADE,
    INDEX idx_chat_messages_session (session_id)
);
```

---

## 6. API Endpoints AI Chat

| Method | Endpoint | Auth | Deskripsi |
| :--- | :--- | :---: | :--- |
| `POST` | `/api/v1/ai/chat` | ✅ JWT | Kirim pesan ke NutriBot dan terima balasan |
| `GET` | `/api/v1/ai/chat/sessions` | ✅ JWT | Daftar semua sesi percakapan pengguna |
| `GET` | `/api/v1/ai/chat/sessions/:session_id` | ✅ JWT | Riwayat pesan dalam satu sesi |
| `DELETE` | `/api/v1/ai/chat/sessions/:session_id` | ✅ JWT | Hapus satu sesi percakapan |

### 6.1 Request: `POST /api/v1/ai/chat`
```json
{
  "session_id": "uuid-existing-session-or-null",
  "message": "Berapa kalori nasi padang?"
}
```

### 6.2 Response: `POST /api/v1/ai/chat`
```json
{
  "session_id": "uuid-xxx",
  "reply": "Nasi Padang satu porsi standar mengandung sekitar 700–900 kkal tergantung lauk...",
  "xp_earned": 10,
  "model_used": "llama-3.3-70b-versatile",
  "latency_ms": 1250
}
```

---

## 7. Arsitektur Internal — Perluasan Domain `ai`

Domain `ai` yang sudah ada diperluas dengan sub-modul chat:

```
internal/domain/ai/
├── model.go            -- [EXISTING] struct AIResultLog
├── model_chat.go       -- [NEW] struct ChatSession, ChatMessage
├── dto.go              -- [EXISTING] NutritionAnalysisRequest, dll
├── dto_chat.go         -- [NEW] ChatRequest, ChatResponse, SessionListResponse
├── repository.go       -- [EXISTING]
├── repository_chat.go  -- [NEW] interface ChatRepository + MySQL impl
├── service.go          -- [EXISTING] AnalyzeNutrition
├── service_chat.go     -- [NEW] interface ChatService + impl (Chat, GetSessions, GetMessages)
└── handler.go          -- [EXISTING]
└── handler_chat.go     -- [NEW] Gin handlers untuk Chat endpoints
```

### Alur `Chat(userID, sessionID, message)`:
```
1. Jika sessionID kosong → buat sesi baru (ai_chat_sessions)
2. Ambil history pesan sebelumnya di sesi (max. 10 pesan terakhir untuk context window)
3. Susun messages array: [system_prompt, ...history, {role: "user", content: message}]
4. Kirim ke Groq API via pkg/groq
5. Simpan pesan user + balasan AI ke ai_chat_messages
6. Trigger gamification.AddXP(userID, "ai_chat", 10) — hanya sekali per sesi baru
7. Return balasan AI + session_id + xp_earned
```

---

## 8. Frontend Implementation Plan

```
app/ai-chat/
├── page.tsx                   -- Halaman penuh /ai-chat
├── AiChatContent.tsx          -- Client component utama halaman
└── components/
    ├── ChatBubble.tsx         -- Komponen bubble pesan (user & AI)
    ├── QuickPrompts.tsx       -- Panel 6 quick prompt suggestions
    └── SessionSidebar.tsx     -- Daftar sesi percakapan (desktop)

internal/components/ai/
├── NutriChatWidget.tsx        -- Floating widget (FAB + pop-up window)
└── NutriChatWindow.tsx        -- Jendela chat pop-up (360x520px)

internal/domain/ai/
├── hooks/
│   └── useNutriChat.ts        -- React Query hook: sendMessage, getSessions
└── types/
    └── chat.ts                -- TypeScript types: ChatMessage, ChatSession
```

---

## 9. Acceptance Criteria

- [ ] Widget floating 🤖 muncul di semua halaman (kecuali halaman `/login` & `/register`).
- [ ] Klik tombol floating membuka pop-up chat window tanpa pindah halaman.
- [ ] Quick Prompt mengisi input dan langsung mengirim pesan ke NutriBot.
- [ ] Balasan AI muncul dengan typing indicator (animasi dots selama loading).
- [ ] Klik tombol expand memindahkan ke `/ai-chat` dengan riwayat sesi yang sama.
- [ ] Setiap sesi baru memberikan +10 XP kepada pengguna (konfirmasi via toast/notif).
- [ ] NutriBot menolak pertanyaan di luar topik gizi dengan ramah.
- [ ] Riwayat percakapan tersimpan dan bisa diakses kembali di `/ai-chat`.
- [ ] Performa: waktu respons AI < 3 detik (P95).

---

## 10. Rencana Sprint Pengerjaan

| Sprint | Task |
| :--- | :--- |
| Sprint 1 | Migrasi DB: tabel `ai_chat_sessions` & `ai_chat_messages` |
| Sprint 2 | Backend: `model_chat.go`, `dto_chat.go`, `repository_chat.go` |
| Sprint 3 | Backend: `service_chat.go` (logika chat + context window + XP trigger) |
| Sprint 4 | Backend: `handler_chat.go` + register route di router |
| Sprint 5 | Frontend: hooks `useNutriChat.ts` + types |
| Sprint 6 | Frontend: `NutriChatWidget.tsx` (floating FAB + pop-up window) |
| Sprint 7 | Frontend: Halaman `/ai-chat` (full page + sidebar session) |
| Sprint 8 | Testing e2e + UX polish (animasi, responsive mobile) |
