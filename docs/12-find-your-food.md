# 🔍 Fitur Tambahan: Find Your Food (Atlas Makananku)

> Fitur ini **BERDIRI SENDIRI** dan **TIDAK BERHUBUNGAN** dengan alur food recall survey.
> Tidak butuh login. Semua endpoint yang dipakai bersifat **public/open**.

---

## 1. Gambaran Umum

**Find Your Food** adalah modul katalog makanan publik yang terinspirasi dari buku **Atlas Makananku** (BRIN × UPI). Fungsinya sebagai referensi visual estimasi porsi makanan — peneliti, tenaga kesehatan, dan masyarakat umum bisa mencari makanan dan melihat foto porsinya **tanpa perlu login**.

Fitur ini muncul di **Landing Page** sebagai salah satu dari dua pilihan utama:

```
┌──────────────────────────────────────────────────────────────┐
│           Platform Survey Food Recall & Informasi Gizi       │
│                                                              │
│   ┌─────────────────────┐   ┌──────────────────────────┐     │
│   │  📋  Food Recall    │   │  🔍  Find Your Food       │     │
│   │      Survey         │   │                          │     │
│   │  Isi survey recall  │   │  Cari makanan & lihat    │     │
│   │  makananmu          │   │  estimasi porsi visual   │     │
│   │  [Login Required]   │   │  [No Login Needed]       │     │
│   └─────────────────────┘   └──────────────────────────┘     │
└──────────────────────────────────────────────────────────────┘
```

- **Tombol kiri (Food Recall Survey):** Redirect ke `/login` untuk memulai survey.
- **Tombol kanan (Find Your Food):** Langsung masuk ke halaman `/find-food` tanpa login.

---

## 2. Tipe Foto & Perilaku UI

Sesuai Atlas Makananku, ada **2 tipe foto porsi** yang relevan untuk MVP:

| Tipe Foto | Deskripsi | Perilaku UI |
|-----------|-----------|-------------|
| **`series`** | 4–8 foto bertahap dari kecil ke besar (misal: Nasi 50g → 350g) | Tampilkan slider/carousel horizontal. Ketika satu foto dipilih → menjadi **main focus (tampil besar)**. |
| **`range`** | Variasi bentuk/ukuran alami yang tidak bertahap ketat (misal: variasi ukuran Ayam Goreng) | Tampilkan grid thumbnail. Ketika salah satu dipilih → menjadi **main focus (tampil besar)**. |

> Tipe `guide` (multi-item dalam 1 foto) di-skip untuk MVP karena kompleksitas overlay.

---

## 3. Alur User di Find Your Food

### 3.1 Alur Utama

```
Landing Page
      │
      ▼
Pilih "Find Your Food" (tanpa login)
      │
      ▼
Halaman /find-food
- Search bar besar
- Grid 13 kategori (MP, LH, LN, AS, AB, AP, AMK, KK, ABK, AK, MDL, GK, AH)
      │
      ├── Ketik di search → autocomplete hasil makanan (real-time)
      │         │
      │         ▼
      │     Halaman Hasil Pencarian (/find-food/search?q=...)
      │         │
      │         ▼
      │     Grid card makanan (foto thumbnail, kode, nama)
      │         │
      │         ▼
      │     Klik card → Halaman Detail Makanan (/find-food/:foodId)
      │
      └── Klik kategori → Halaman Kategori (/find-food/category/:categoryCode)
                │
                ▼
            Grid makanan per kategori
                │
                ▼
            Klik card → Halaman Detail Makanan
```

### 3.2 Halaman Detail Makanan (MAIN FOCUS: Foto Porsi)

Ini adalah inti dari fitur ini. Foto porsi adalah elemen paling besar.

```
┌──────────────────────────────────────────────────────────────┐
│  ← Kembali    MP-01 · Makanan Pokok · Series    [🔖 Simpan] │
│  NASI / RICE                                                  │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐   │
│  │                                                        │   │
│  │         [ FOTO UTAMA BESAR — Main Focus ]              │   │← Full-width
│  │              A · 50 gram                               │   │  zoomable
│  │                                                        │   │
│  └────────────────────────────────────────────────────────┘   │
│                                                              │
│  ◄  [A] [B] [C] [D] [E] [F] [G] [H]  ►    ← Pilih foto    │
│     50g  90g 130g 150g 210g 250g 270g 350g                   │
│                                                              │
├──────────────────────────────────────────────────────────────┤
│  Tabel Ukuran Porsi                                           │
│  ┌──────┬───────┬────────────────────────────────────────┐   │
│  │ Kode │ Berat │ Keterangan                             │   │
│  ├──────┼───────┼────────────────────────────────────────┤   │
│  │  A   │  50 g │ Porsi sangat kecil                     │   │
│  │  B   │  90 g │ Porsi kecil                            │   │
│  │  H   │ 350 g │ Porsi sangat besar                     │   │
│  └──────┴───────┴────────────────────────────────────────┘   │
├──────────────────────────────────────────────────────────────┤
│  📊 Informasi Gizi (per 100g)                                │
│  Energi: 130 kcal · Protein: 2.7g · Karbo: 28g · Lemak: 0.3g│
├──────────────────────────────────────────────────────────────┤
│  ← MP-23 Roti Tawar          MP-02 Nasi Goreng →            │
└──────────────────────────────────────────────────────────────┘
```

**Perilaku foto saat dipilih:**
- Klik thumbnail [A] → foto A menjadi **main focus** (tampil besar di atas)
- Label berat & kode tampil langsung di bawah foto besar
- Animasi smooth transition antar foto (slide/fade)

---

## 4. Perubahan Backend yang Dibutuhkan

Fitur ini memanfaatkan data yang **sudah ada** di database (`foods`, `categories`, `food_nutrients`, `nutrient_types`, `as_served_sets`, `as_served_images`) — **tidak butuh tabel baru**.

### 4.1 Endpoint yang Perlu Disesuaikan / Ditambah

| # | Method | Endpoint | Auth | Deskripsi |
|---|--------|----------|------|-----------|
| 1 | `GET` | `/api/v1/public/foods/search` | ❌ None | Search makanan publik (FULLTEXT) |
| 2 | `GET` | `/api/v1/public/foods/:id` | ❌ None | Detail makanan + nutrisi + foto porsi |
| 3 | `GET` | `/api/v1/public/categories` | ❌ None | List 13 kategori |
| 4 | `GET` | `/api/v1/public/categories/:code/foods` | ❌ None | List makanan per kategori |
| 5 | `GET` | `/api/v1/public/foods/:id/portion-photos` | ❌ None | Ambil semua foto porsi beserta beratnya |

> **Catatan:** Endpoint yang sudah ada di `/api/v1/foods` masih perlu login (middleware JWT). Kita akan buat versi **public** di prefix `/api/v1/public/` yang tanpa auth.

### 4.2 Contoh Response `GET /public/foods/:id`

```json
{
  "status": "success",
  "data": {
    "id": "uuid-food",
    "code": "MP-01",
    "name": "Nasi",
    "local_name": "Rice",
    "description": "Nasi putih matang",
    "category": {
      "code": "MP",
      "name": "Makanan Pokok",
      "icon": "🍚"
    },
    "photo_type": "series",
    "nutrients": {
      "energy":  { "value": 130.00, "unit": "kcal" },
      "protein": { "value": 2.70,   "unit": "g" },
      "carbs":   { "value": 28.00,  "unit": "g" },
      "fat":     { "value": 0.30,   "unit": "g" }
    },
    "portion_photos": [
      {
        "id": "uuid-img-1",
        "label": "A",
        "image_url": "/uploads/nasi/nasi-A.jpg",
        "thumbnail_url": "/uploads/nasi/nasi-A-thumb.jpg",
        "weight_gram": 50.0,
        "description": "Porsi sangat kecil"
      },
      {
        "id": "uuid-img-2",
        "label": "B",
        "image_url": "/uploads/nasi/nasi-B.jpg",
        "thumbnail_url": "/uploads/nasi/nasi-B-thumb.jpg",
        "weight_gram": 90.0,
        "description": "Porsi kecil"
      }
    ]
  }
}
```

### 4.3 Perubahan pada Model Food

Perlu tambahkan field `photo_type` di tabel `foods` untuk membedakan tipe foto:

```sql
-- Migration baru: 006_add_photo_type_to_foods.sql
ALTER TABLE foods 
ADD COLUMN photo_type ENUM('series', 'range', 'guide') DEFAULT 'series'
AFTER description;
```

### 4.4 Perbaikan Search: LIKE → FULLTEXT MATCH

Ubah implementasi `SearchFoods` di `food/repository.go` dari `LIKE` menjadi `MATCH ... AGAINST` agar memanfaatkan `FULLTEXT INDEX ft_name`:

```go
// SEKARANG (tidak optimal — tidak pakai FULLTEXT index):
err := q.Where("name LIKE ? OR local_name LIKE ?", "%"+query+"%", "%"+query+"%").Limit(limit).Find(&foods).Error

// SEHARUSNYA (pakai FULLTEXT):
err := q.Where("MATCH(name, local_name) AGAINST(? IN BOOLEAN MODE)", query+"*").Limit(limit).Find(&foods).Error
```

---

## 5. Perubahan File Backend

### 5.1 File yang Perlu Dimodifikasi

| File | Perubahan |
|------|-----------|
| `internal/domain/food/repository.go` | Ubah `SearchFoods` ke FULLTEXT; tambah `GetFoodWithPortionPhotos()` |
| `internal/domain/food/service.go` | Tambah method `GetPublicFoodDetail()` yang return nutrisi + foto porsi |
| `internal/domain/food/handler.go` | Tambah handler untuk public endpoints |
| `internal/domain/food/dto.go` | Tambah `PublicFoodDetailResponse` dengan `portion_photos` |
| `internal/router/router.go` | Register `/public/*` routes tanpa middleware auth |
| `migrations/` | Tambah `006_add_photo_type_to_foods.sql` |

---

## 6. Frontend Checklist (Scope FE)

> Dikerjakan di repo FE di branch `dashboard-recall`.

- [ ] Update `app/page.tsx` (Landing Page) — Tambah 2 card: Food Recall Survey & Find Your Food
- [ ] Buat halaman `/find-food` — Search bar + Grid 13 kategori
- [ ] Buat halaman `/find-food/search` — Halaman hasil pencarian
- [ ] Buat halaman `/find-food/[foodId]` — Halaman detail dengan foto jadi main focus
  - [ ] Komponen `PortionPhotoViewer` — Series: klik thumbnail → jadi main focus + animasi
  - [ ] Komponen `PortionPhotoViewer` — Range: grid gallery, klik → main focus
  - [ ] Tabel ukuran porsi
  - [ ] Kartu nutrisi per 100g
  - [ ] Navigasi prev/next makanan dalam kategori
- [ ] Buat halaman `/find-food/category/[code]` — Grid makanan per kategori
- [ ] Fitur Bookmark / Simpan (localStorage, tanpa login)

---

## 7. Estimasi Pengerjaan

| Task | BE | FE |
|------|----|----|
| Public endpoints + routing tanpa auth | 0.5 hari | — |
| FULLTEXT search fix | 0.5 hari | — |
| Migrasi `photo_type` field | 0.5 hari | — |
| Landing page (2 card) | — | 0.5 hari |
| Halaman Find Your Food (search + kategori) | — | 1 hari |
| Halaman Detail Makanan (foto main focus + tabel) | — | 1.5 hari |
| **Total** | **1.5 hari** | **3 hari** |
