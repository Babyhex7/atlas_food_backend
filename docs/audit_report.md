# 🔍 Audit & Gap Analysis Report: Atlas Food Backend (MVP V2)

Laporan ini menyajikan hasil audit mendalam terhadap basis kode **Atlas Food Backend** untuk menilai kesesuaiannya dengan spesifikasi fungsional **Intake24**, kelengkapan dokumentasi yang diimplementasikan, serta kepatuhan terhadap praktik terbaik (*best practices*) pengembangan perangkat lunak modern menggunakan Go (Golang).

---

## 🎯 1. Apakah Sudah Seperti Intake24?

**Ya, secara konseptual.** Basis kode ini sangat baik dalam mengadopsi esensi inti dari **Intake24**, tetapi dengan pendekatan yang disederhanakan (*pragmatic MVP*) agar realistis dikerjakan dalam waktu singkat.

### 📊 Perbandingan Konseptual:
*   **Intake24 Asli (~60+ Tabel):** Memiliki skema database yang sangat masif, mendukung banyak negara (*multi-locale/region*), integrasi eksternal database nutrisi (seperti USDA, BPOM, dll.), serta metode estimasi porsi yang sangat kompleks (skala visual gelas/piring linier, berat standar kategori, dll.).
*   **Atlas Food MVP V2 (17 Tabel):** Memangkas fitur non-esensial (misalnya *Guide Images* berbasis area koordinat piksel atau pemrosesan gambar otomatis) dan merampingkannya menjadi **17 tabel** terpadu. Fitur inti Intake24 berikut **sukses diadopsi**:
    1.  **Multi-Pass Recall Flow:** Melalui pencatatan makanan per waktu makan (*Meals Data*) yang dikirim responden.
    2.  **Portion Estimation:** Mendukung metode *As-Served* (pilihan gambar visual porsi beserta berat gramnya), *Counter* (whole + fraction porsi), serta input berat manual (*Weight*).
    3.  **Food & Nutrient Database:** Struktur relasi yang rapi antara `foods`, `categories`, `nutrient_types`, `nutrient_units`, dan `food_nutrients` (per 100g).
    4.  **Associated Foods:** Fitur cerdas untuk menyarankan makanan pendamping (contoh: makan Sereal disarankan menambah Susu).

---

## 🛠️ 2. Status Implementasi Fitur (Cek Dokumen vs Code)

Berdasarkan analisis file sumber di Go dan perbandingan dengan `docs/10-implementation-steps.md`, berikut adalah status riil implementasi di dalam basis kode:

| Phase | Fitur | Status di Dokumen | Status di Code | Detail & Catatan |
|---|---|---|---|---|
| **Phase 1** | Foundation | ✅ DONE | **100% DONE** | Project setup, folder structure, logger, CORS, & error middlewares sudah sangat rapi. |
| **Phase 2** | Authentication | ✅ DONE | **100% DONE** | JWT Auth + Refresh Token Rotation (menggunakan SHA256 hash di DB untuk keamanan token) sudah terimplementasi secara kokoh. |
| **Phase 3** | Survey Management | ✅ DONE | **100% DONE** | CRUD Survey, manajemen partisipan, slug generation, dan public survey access token berjalan dengan baik. |
| **Phase 4** | Food Database | ✅ DONE | **100% DONE** | CRUD makanan, kategori, dan pencarian makanan (*Search Foods*) sudah siap pakai. |
| **Phase 5** | Portion Size | ✅ DONE | **100% DONE** | Model data `PortionSizeMethod`, `AsServedSet`, `AsServedImage`, serta `AssociatedFood` sudah terbentuk dan terintegrasi di tingkat GORM. |
| **Phase 6** | Survey Submission | ✅ DONE | **90% DONE** | Handler untuk `/submit` sudah ada, agregasi nutrisi harian dan ekspor data CSV sudah terimplementasi. *(Ada catatan best practice di bawah)*. |
| **Phase 7** | File Upload | ✅ DONE | **100% DONE** | Upload gambar statis ke direktori `./uploads` dengan validasi tipe file dan batas ukuran file berjalan mulus. |
| **Phase 8** | **AI Nutrition (Groq)** | ⏳ PENDING | **❌ 0% (BELUM ADA)** | **Sama sekali belum ada di kode.** Domain `internal/domain/ai` belum dibuat, helper Groq client belum diimplementasikan, tabel `ai_result_logs` belum didaftarkan di migrasi, dan rute `/ai/nutrition-analysis` belum terdaftar. |

---

## ⚠️ 3. Analisis Gaps & Praktik Terbaik (Best Practices Check)

Meskipun fondasi Clean Architecture-nya sudah sangat bagus dan rapi, ada beberapa **celah keamanan gizi (*nutrient validation*)** dan **gaps** penting yang perlu diperbaiki agar mencapai tingkat *best practice* sesungguhnya:

### 🚨 Gap 1: Kerentanan Manipulasi Data Nutrisi oleh Client (Critical)
*   **Kondisi Saat Ini:** Di file [service.go](file:///c:/Users/mybook_bagas/Projek_Agency/atlas_food/atlas_food_backend/internal/domain/submission/service.go#L137-L171), proses `calculateTotals` hanya menjumlahkan data gizi yang **dikirim langsung oleh client** via request body JSON (`f.Nutrients.Energy`, dll.).
*   **Dampak:** Terjadi celah keamanan data (*data tampering vulnerability*). Respondent nakal bisa memanipulasi request JSON di frontend (misalnya mengirim porsi 500g nasi goreng tetapi mengisi `nutrients.energy = 0` dan `nutrients.protein = 100`). Backend akan langsung menyimpan data palsu tersebut ke database tanpa memvalidasinya ke database makanan riil.
*   **Solusi Best Practice:** Backend **wajib** melakukan verifikasi nilai gizi server-side.
    1. Ambil `food_id` dan `portion_gram` dari request.
    2. Query database makanan riil via `food_nutrients` untuk mendapatkan gizi per 100g.
    3. Lakukan perhitungan secara dinamis di server menggunakan rumus: `(nilai_gizi_per_100g / 100) * portion_gram`.
    4. Simpan nilai hasil kalkulasi server ke database submission untuk menjamin integritas data riset.

### 🤖 Gap 2: Phase 8 (AI Groq) Belum Terimplementasi
Untuk memenuhi spesifikasi di [brif_ai.md](file:///c:/Users/mybook_bagas/Projek_Agency/atlas_food/atlas_food_backend/docs/brif_ai.md), langkah-langkah berikut harus segera diwujudkan:
1.  **GORM Migrations:** Tambahkan model/tabel `ai_result_logs` ke dalam list `AutoMigrate` pada `main.go`.
2.  **Groq Client:** Buat HTTP wrapper di `internal/pkg/groq/client.go` untuk memanggil Groq Chat Completions API (`https://api.groq.com/openai/v1/chat/completions`) menggunakan model `llama3-8b-8192` dan mengembalikan JSON terstruktur.
3.  **Clean Architecture AI Domain:** Buat folder `internal/domain/ai` berisi `dto.go`, `model.go`, `repository.go`, `service.go`, dan `handler.go`.
4.  **Cache By Submission ID:** Pastikan ada pengecekan ke tabel `ai_result_logs` terlebih dahulu sebelum menembak API Groq. Jika data cache ada, langsung return dengan status `source: "cache"`.
5.  **Ownership Check:** Validasi bahwa responden yang meminta analisis AI adalah pemilik sah dari `submission_id` tersebut.

---

## 💡 4. Kesimpulan & Rekomendasi Langkah Selanjutnya

1.  **Arsitektur & Kode Dasar:** **Sangat Baik!** Penggunaan Clean Architecture di Go sudah sangat konsisten dan rapi.
2.  **Kesesuaian Dokumentasi:** Dokumentasi di direktori `/docs` sangat komprehensif dan profesional. Namun, ada *discrepancy* di mana dokumen mengasumsikan Phase 8 sudah berjalan (atau siap), padahal di kode nyata masih belum tersentuh.
3.  **Langkah Selanjutnya:**
    *   **Perbaikan Server-Side Nutrition Aggregation:** Ubah logika di `submission.Service` agar mengambil nilai gizi dari database makanan riil alih-alih mempercayai input mentah dari client.
    *   **Implementasi Phase 8 (AI Groq):** Mulai buat modul AI on-demand ini agar aplikasi Atlas Food MVP V2 menjadi utuh 100% dan siap di-deploy.

---

*Laporan audit ini dibuat pada tanggal **17 Mei 2026**.*
