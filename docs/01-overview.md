# 🎯 Atlas Food — MVP Overview

> Food recall survey platform (inspired by Intake24). Versi minimal yang bisa jalan.
> Fitur non-esensial ditandai "NANTI AJA".

---

## 📋 Deskripsi Proyek

Atlas Food adalah platform survey recall makanan yang memungkinkan responden untuk mencatat apa yang mereka makan dalam periode waktu tertentu. Platform ini terinspirasi dari Intake24 dengan fitur yang disederhanakan untuk MVP.

### Fitur Utama MVP

1. **Survey Management** — Admin dapat membuat, mengelola, dan membagikan survey
2. **Anonymous Access** — Responden dapat mengakses survey tanpa registrasi (via token/link)
3. **Food Database** — Database makanan dengan kategori dan informasi gizi
4. **Portion Selection with Images** — Pemilihan porsi menggunakan gambar visual
5. **Nutrition Calculation** — Perhitungan nutrisi otomatis berdasarkan porsi
6. **Submission & Export** — Pengiriman data dan export ke CSV/JSON

---

## 🏗️ Arsitektur Sistem

### Struktur Database (16 Tabel)

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│    users    │      │   surveys   │      │ survey_part │
├─────────────┤      ├─────────────┤      ├─────────────┤
│ PK id       │──┐   │ PK id       │◀─┐   │ PK id       │
│ email       │  │   │ slug        │  │   │ FK survey_id│
│ role        │  │   │ meals(JSON) │  └───┤ FK user_id  │
│ password    │  │   │ prompts(JSON)    │  │ alias       │
└─────────────┘  │   └─────────────┘         │
     │           │            │              │
     │           │            │              │
     └───────────┼────────────┼──────────────┘
                 │            │
                 ▼            ▼
         ┌─────────────┐     ┌─────────────┐
         │ submissions │     │    foods    │
         ├─────────────┤     ├─────────────┤
         │ PK id       │     │ PK id       │
         │ FK survey_id│     │ name        │
         │ FK part_id  │     │ nutrients   │
         │ meals(JSON) │◀────│  (JSON)     │
         └─────────────┘     └─────────────┘
```

### Fitur yang Masuk MVP V2

| #   | Fitur/Tabel           | Penjelasan                                            |
| --- | --------------------- | ----------------------------------------------------- |
| 1   | **Login/Auth**        | JWT + Refresh Token                                   |
| 2   | **Surveys**           | Create, manage, share link                            |
| 3   | **Food Database**     | Makanan + kategori + nutrisi tabel                    |
| 4   | **Portion Images**    | Langsung pilih gambar porsi (tanpa pilih metode dulu) |
| 5   | **As Served Sets**    | Foto porsi ½,1,1½,2,3 + berat (food & drinks)         |
| 6   | **Nutrisi Tabel**     | nutrient_types, nutrient_units, food_nutrients        |
| 7   | **Associated Foods**  | Sereal → susu, Roti → Selai                           |
| 8   | **Input Flow**        | Add makanan → Continue → Portion (langsung)           |
| 9   | **Submission JSON**   | Simpan hasil recall                                   |

**Total: 16 Tabel**

### Yang DI-SKIP (NANTI AJA)

| Fitur                      | Alasan Skip                     |
| -------------------------- | ------------------------------- |
| Session auto-save (DB)     | Pakai `localStorage` dulu, cukup |
| RBAC kompleks              | 1 field `role` enum cukup        |
| Survey Schemes             | Copy-paste survey dulu           |
| User Sessions              | Pakai `localStorage` dulu        |
| Food Synonyms              | MySQL FULLTEXT cukup             |
| Nutrient Tables            | Multiple sources (USDA, BPOM)    |
| Guide Images               | Area-based portions              |
| Image Processing           | Auto-resize, thumbnail           |
| Email Service              | Reset password, notif            |
| Social Login               | Google/Facebook                  |
| Multi-language DB          | next-intl file JSON dulu         |

---

## 📊 Perbandingan: MVP V2 vs Full Intake24

| Aspek               | MVP V2 (16 Tabel)            | Full Intake24 (~60+ Tabel)      |
| ------------------- | ---------------------------- | ------------------------------- |
| **Tabel**           | 16                           | ~60+                           |
| **Portion Methods** | as_served (food + drinks), weight | + drinkware, guide_image, standard_portion |
| **Images**          | Upload manual                | Auto-process, multiple sizes    |
| **Search**          | MySQL FULLTEXT               | Meilisearch/Elasticsearch       |
| **Auth**            | JWT simple                   | JWT + RBAC + 2FA                |
| **Sessions**        | localStorage                 | Redis + DB                      |
| **Nutrition**       | Tabel: nutrient_types, nutrient_units, food_nutrients | Multiple sources (USDA, BPOM) + history |
| **Time Estimate**   | 1 bulan                      | 6+ bulan                        |

---

## ⏱️ Estimasi Waktu Development

| Phase | Deskripsi | Estimasi |
|-------|-----------|----------|
| Phase 1 | Auth & User (Setup + JWT) | 3-4 hari |
| Phase 2 | Survey Admin (CRUD + Token) | 3-4 hari |
| Phase 3 | Food DB + Portion Size | 5-7 hari |
| Phase 4 | Survey Flow + Submit | 5-7 hari |
| Polish | Bugfix & Testing | 3-5 hari |
| **Total** | | **3-4 minggu** |

---

_versi ini realistis untuk 1 developer dalam 1 bulan._
