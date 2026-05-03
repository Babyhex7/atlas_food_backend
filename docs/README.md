# 📚 Atlas Food Documentation

Selamat datang di dokumentasi Atlas Food! Dokumentasi ini berisi informasi lengkap tentang arsitektur, API, database, dan alur aplikasi Atlas Food.

---

## 📖 Daftar Dokumen

| No | Dokumen | Deskripsi |
|----|---------|-----------|
| 1 | [01-overview.md](./01-overview.md) | Overview proyek, fitur MVP, dan estimasi waktu |
| 2 | [02-tech-stack.md](./02-tech-stack.md) | Tech stack, project structure, dan dependencies |
| 3 | [03-database-schema.md](./03-database-schema.md) | Schema database lengkap dengan 16 tabel |
| 4 | [04-api-documentation.md](./04-api-documentation.md) | Dokumentasi API lengkap dengan request/response |
| 5 | [05-workflow-alur.md](./05-workflow-alur.md) | Alur aplikasi dan user flow |
| 6 | [06-user-story.md](./06-user-story.md) | User stories untuk admin dan respondent |
| 7 | [07-portion-calculation.md](./07-portion-calculation.md) | Rumus perhitungan porsi dan nutrisi |
| 8 | [08-erd.md](./08-erd.md) | Entity Relationship Diagram lengkap |
| 9 | [09-sample-data.md](./09-sample-data.md) | Sample data SQL dan JSON |

---

## 🚀 Quick Start

### 1. Baca Overview
Mulai dari [01-overview.md](./01-overview.md) untuk memahami scope proyek dan fitur MVP.

### 2. Setup Development
Lihat [02-tech-stack.md](./02-tech-stack.md) untuk setup environment dan project structure.

### 3. Database
Lihat [03-database-schema.md](./03-database-schema.md) dan [08-erd.md](./08-erd.md) untuk memahami struktur database.

### 4. API Development
Lihat [04-api-documentation.md](./04-api-documentation.md) untuk endpoint dan format request/response.

### 5. Business Logic
Lihat [05-workflow-alur.md](./05-workflow-alur.md), [06-user-story.md](./06-user-story.md), dan [07-portion-calculation.md](./07-portion-calculation.md) untuk memahami alur bisnis.

---

## 📝 Ringkasan Proyek

**Atlas Food** adalah platform survey recall makanan (inspired by Intake24) dengan fitur:

- **Survey Management** — Buat dan kelola survey dengan mudah
- **Anonymous Access** — Responden akses via link/token tanpa registrasi
- **Food Database** — Database makanan dengan informasi gizi lengkap
- **Portion Selection dengan Gambar** — Pilih porsi menggunakan visual
- **Nutrition Calculation** — Perhitungan gizi otomatis
- **Export Data** — Export submission ke CSV/JSON

### MVP V2 (16 Tabel)

| Fitur | Status |
|-------|--------|
| Auth (JWT + Refresh Token) | ✅ Wajib |
| Survey CRUD | ✅ Wajib |
| Food Database | ✅ Wajib |
| Portion Images | ✅ Wajib |
| Nutrition Calculation | ✅ Wajib |
| Associated Foods | ✅ Wajib |
| Survey Submission | ✅ Wajib |

### Tech Stack

- **Backend:** Go + Gin + GORM + MySQL
- **Auth:** JWT + Refresh Tokens
- **Architecture:** Clean Architecture (Feature-based)

### Estimasi Waktu

**Total: 3-4 minggu untuk MVP**

---

## 🔗 Referensi

- [Intake24](https://intake24.co.uk/) — Platform recall makanan yang jadi inspirasi
- [Gin Framework](https://gin-gonic.com/) — Go web framework
- [GORM](https://gorm.io/) — Go ORM

---

_versi ini realistis untuk 1 developer dalam 1 bulan._
