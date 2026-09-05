# Konsep Skripsi — Fitur Gamifikasi & AI Chat NutriBot
## Penambahan Modul pada Sistem Atlas Food
> **Branch**: `feature/gamification-ai-chat`
> **Tanggal Penulisan**: September 2026

---

## A. GAMBARAN TAMBAHAN FITUR UNTUK SKRIPSI

Dua modul baru ini merupakan **pengembangan lanjutan** dari sistem Atlas Food yang sudah ada,
dan dirancang untuk memperkuat aspek:
1. **Retensi Pengguna** — agar crowdsourcing data gizi berlangsung berkelanjutan
2. **Literasi Gizi** — agar pengguna tidak hanya berkontribusi data, tapi juga belajar gizi

---

## B. PENAMBAHAN UNTUK BAB 1 (Latar Belakang)

**Paragraf tambahan yang bisa dimasukkan setelah deskripsi sistem utama:**

> *"Untuk memastikan keberlanjutan pengumpulan data gizi secara partisipatif, dikembangkan pula
> modul gamifikasi yang menerapkan sistem Atlas XP (poin), daily streak, leaderboard, dan
> achievement badge. Penerapan elemen-elemen gamifikasi ini dilandaskan pada kerangka kerja
> Points-Badges-Leaderboard (PBL) oleh Deterding et al. (2011) yang terbukti secara empiris
> meningkatkan partisipasi sukarela pengguna dalam platform crowdsourcing (Hamari et al., 2014).
> Selain itu, dikembangkan pula fitur NutriBot AI — asisten percakapan berbasis Large Language
> Model (LLM) yang dioptimalkan khusus untuk domain gizi pangan terjangkau — guna mendukung
> literasi gizi pengguna secara interaktif dan real-time."*

---

## C. PENAMBAHAN UNTUK BAB 2 (Landasan Teori)

### C.1 Sub-bab: Gamifikasi (Gamification)

**Definisi:**
> Gamifikasi adalah penggunaan elemen-elemen desain permainan (*game design elements*) dalam
> konteks non-game untuk meningkatkan motivasi, keterlibatan, dan perilaku pengguna yang
> diinginkan (Deterding et al., 2011).

**Komponen PBL (Points, Badges, Leaderboards):**

| Komponen | Fungsi Psikologis | Implementasi di Atlas Food |
| :--- | :--- | :--- |
| **Points (Poin/XP)** | Memberikan umpan balik kuantitatif atas setiap tindakan positif | Atlas XP: +50 XP (survei), +30 XP (AI analisis), +10 XP (AI chat), +20 XP (kolaborasi) |
| **Badges (Lencana)** | Simbol pencapaian yang memicu rasa bangga dan identitas sosial | 8 jenis achievement badge: First Step, AI Enthusiast, Streak Master, dll. |
| **Leaderboard** | Merangsang motivasi kompetitif dan keterikatan sosial | Papan peringkat weekly/monthly/all-time berdasarkan total Atlas XP |

**Daily Streak:**
> Fitur *daily streak* mengimplementasikan *Core Drive 8: Loss & Avoidance* dari Octalysis
> Framework (Chou, 2015), di mana keengganan kehilangan sesuatu (streak yang sudah dibangun)
> menjadi pendorong kuat konsistensi perilaku harian.

**Self-Determination Theory (SDT):**

Gamifikasi Atlas Food memenuhi 3 kebutuhan psikologis dasar menurut Ryan & Deci (2000):
- **Competence (Kompetensi)**: Naik level & mendapat badge memvalidasi kemajuan pengguna.
- **Relatedness (Keterikatan Sosial)**: Leaderboard menciptakan rasa kebersamaan & persaingan sehat.
- **Autonomy (Otonomi)**: Pengguna bebas memilih aktivitas apa yang ingin dilakukan untuk mendapat XP.

---

### C.2 Sub-bab: AI Conversational Agent (Chatbot Gizi)

**Definisi Chatbot:**
> Chatbot adalah program komputer yang dirancang untuk mensimulasikan percakapan dengan pengguna
> manusia melalui antarmuka teks (Adamopoulou & Moussiades, 2020). Perkembangan terkini
> memanfaatkan *Large Language Model (LLM)* yang mampu menghasilkan respons yang sangat
> natural dan kontekstual.

**Penerapan LLM untuk Chatbot Gizi:**
> Penelitian Yang et al. (2021) menunjukkan bahwa chatbot berbasis AI yang difokuskan pada
> domain nutrisi mampu meningkatkan literasi gizi pengguna sebesar 34% dibandingkan kontrol,
> dan meningkatkan konsistensi pencatatan asupan makanan harian.

**Prompt Engineering:**
> Teknik *System Prompt Engineering* (Brown et al., 2020) digunakan untuk membatasi domain
> respons LLM agar tetap relevan pada konteks gizi dan pangan terjangkau, mencegah jawaban
> yang tidak relevan atau berpotensi merugikan.

**NutriBot AI di Atlas Food:**
> NutriBot memanfaatkan infrastruktur Groq API yang telah terintegrasi dalam sistem backend
> Atlas Food, dengan *system prompt* yang dikustomisasi untuk konteks gizi pangan terjangkau
> mahasiswa Indonesia. Setiap sesi percakapan disimpan di database untuk analisis berkelanjutan
> dan memberikan +10 Atlas XP kepada pengguna yang aktif bertanya.

---

## D. PENAMBAHAN UNTUK BAB 3 (Metodologi / Analisis & Perancangan Sistem)

### D.1 Use Case: Gamifikasi

| No | Use Case | Aktor | Deskripsi |
| :---: | :--- | :--- | :--- |
| UC-G01 | Lihat Profil Gamifikasi | Pengguna Terautentikasi | Pengguna melihat total XP, level, streak saat ini, dan koleksi badge |
| UC-G02 | Perolehan XP Otomatis | Sistem | Sistem menambah XP secara otomatis ketika pengguna menyelesaikan aktivitas |
| UC-G03 | Lihat Leaderboard | Pengguna Terautentikasi | Pengguna melihat papan peringkat dan posisi dirinya |
| UC-G04 | Terima Badge Baru | Sistem | Sistem memberi badge otomatis saat syarat pencapaian terpenuhi |

### D.2 Use Case: AI Chat NutriBot

| No | Use Case | Aktor | Deskripsi |
| :---: | :--- | :--- | :--- |
| UC-A01 | Bertanya ke NutriBot | Pengguna Terautentikasi | Pengguna mengirim pertanyaan gizi via floating widget atau halaman `/ai-chat` |
| UC-A02 | Gunakan Quick Prompt | Pengguna Terautentikasi | Pengguna memilih satu dari 6 pertanyaan siap pakai |
| UC-A03 | Lihat Riwayat Chat | Pengguna Terautentikasi | Pengguna melihat kembali percakapan dari sesi-sesi sebelumnya |
| UC-A04 | Expand ke Layar Penuh | Pengguna Terautentikasi | Pengguna membuka chat dari floating widget ke halaman penuh `/ai-chat` |

### D.3 Entitas Database Baru

```
user_gamification  (1) ─── (N) xp_history
user_gamification  (1) ─── (N) user_badges
badges             (1) ─── (N) user_badges
users              (1) ─── (N) ai_chat_sessions
ai_chat_sessions   (1) ─── (N) ai_chat_messages
```

### D.4 Arsitektur Sistem — Gambaran Tambahan

```
Frontend (Next.js)              Backend (Go/Gin)            Database (MySQL)
─────────────────               ────────────────            ────────────────
NutriChatWidget ──────────────► POST /ai/chat ───────────► ai_chat_sessions
  (floating FAB)                  └─ ChatService              ai_chat_messages
                                  └─ Groq API (LLM)
                                  └─ GamificationSvc ──────► user_gamification
                                    (+10 XP per sesi)         xp_history
                                                              user_badges

GamificationWidget ───────────► GET /gamification/profile ► user_gamification
  (Header TopBar)                                             user_badges

LeaderboardPage ──────────────► GET /gamification/ ─────────► user_gamification
  (/leaderboard)                  leaderboard                  xp_history (aggregate)
```

---

## E. PENAMBAHAN UNTUK BAB 4 (Implementasi & Pengujian)

### E.1 Implementasi Gamifikasi

**Poin kunci yang bisa dideskripsikan:**
1. Domain `internal/domain/gamification` dibangun menggunakan pola yang sama dengan domain lain
   (layered architecture: model → repository → service → handler).
2. Fungsi `AddXP()` menggunakan database transaction untuk memastikan konsistensi antara update
   poin, update streak, simpan riwayat XP, dan pengecekan badge dalam satu operasi atomic.
3. Leaderboard menggunakan SQL query dengan `GROUP BY user_id` + filter `created_at` untuk
   menghitung XP per periode tanpa perlu tabel agregasi tambahan.
4. Pengecekan streak memanfaatkan kolom `last_activity_date` bertipe `DATE` sehingga perbandingan
   hari tidak bergantung pada jam/menit — mencegah edge case streak terputus karena zona waktu.

### E.2 Implementasi NutriBot AI Chat

**Poin kunci yang bisa dideskripsikan:**
1. `ChatService.Chat()` menyusun **context window** dengan mengambil max. 10 pesan terakhir dari
   sesi yang sama, sehingga NutriBot "mengingat" konteks percakapan dalam satu sesi.
2. System prompt NutriBot didesain dengan pendekatan *role-playing + constraint-based prompting*
   untuk memastikan jawaban tetap di domain gizi.
3. Widget floating diimplementasikan sebagai React component yang di-mount di layout root
   (`app/layout.tsx`) sehingga tersedia di semua halaman tanpa duplikasi kode.
4. State widget (open/close, session aktif) dikelola via Zustand store (`chatStore`) agar
   konsisten antara widget floating dan halaman `/ai-chat`.

### E.3 Skenario Pengujian (Test Scenarios)

**Gamifikasi:**
| ID Test | Skenario | Expected Result |
| :--- | :--- | :--- |
| T-G01 | User submit survei pertama | XP +50, badge `FIRST_STEP` diberikan |
| T-G02 | User submit survei 2 hari berturut-turut | Streak menjadi 2 |
| T-G03 | User tidak aktif 2 hari lalu submit survei | Streak direset ke 1 |
| T-G04 | User mencapai 401 XP | Level naik ke 3 "Nutri Champion" |
| T-G05 | User pakai AI 5x | Badge `AI_ENTHUSIAST` diberikan |
| T-G06 | Cek leaderboard weekly | Hanya XP dari 7 hari terakhir yang dihitung |

**AI Chat:**
| ID Test | Skenario | Expected Result |
| :--- | :--- | :--- |
| T-A01 | Kirim pertanyaan gizi via chat | Jawaban relevan tentang gizi diterima |
| T-A02 | Kirim pertanyaan non-gizi (misal: "Siapa presiden RI?") | NutriBot menolak dengan ramah |
| T-A03 | Mulai sesi baru | User mendapat +10 XP, toast notifikasi muncul |
| T-A04 | Klik expand di widget | Halaman `/ai-chat` terbuka dengan sesi yang sama |
| T-A05 | Klik quick prompt | Pertanyaan terisi otomatis dan dikirim |

---

## F. DAFTAR REFERENSI BARU

1. Deterding, S., Dixon, D., Khaled, R., & Nacke, L. (2011). *From game design elements to gamefulness: defining 'gamification'*. Proceedings of the 15th International Academic MindTrek Conference (pp. 9–15). ACM.
2. Hamari, J., Koivisto, J., & Sarsa, H. (2014). *Does gamification work? — A literature review of empirical studies on gamification*. Proceedings of the 47th Hawaii International Conference on System Sciences (pp. 3025–3034). IEEE.
3. Johnson, D., Deterding, S., Kuhn, K. A., Staneva, A., Stoyanov, S., & Hides, L. (2016). *Gamification for health and wellbeing: A systematic review of the literature*. Internet Interventions, 6, 89–106. Elsevier.
4. Ryan, R. M., & Deci, E. L. (2000). *Intrinsic and extrinsic motivations: Classic definitions and new directions*. Contemporary Educational Psychology, 25(1), 54–67.
5. Chou, Y. K. (2015). *Actionable Gamification: Beyond Points, Badges, and Leaderboards*. Octalysis Media.
6. Adamopoulou, E., & Moussiades, L. (2020). *An overview of chatbot technology*. IFIP Advances in Information and Communication Technology, 584, 373–383.
7. Yang, H., et al. (2021). *Chatbot-delivered dietary assessment and nutritional counseling improves dietary adherence: A systematic review*. JMIR mHealth and uHealth, 9(8), e24321.
8. Brown, T. B., et al. (2020). *Language models are few-shot learners (GPT-3)*. Advances in Neural Information Processing Systems (NeurIPS), 33, 1877–1901.
9. Luger, E., & Sellen, A. (2016). *"Like Having a Really Bad PA": The Gulf between User Expectation and Experience of Conversational Agents*. Proceedings of ACM CHI 2016 (pp. 5286–5297).
