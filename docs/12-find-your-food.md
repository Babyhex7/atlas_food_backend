# 🔍 Doc #12 — Find Your Food (Fitur Publik Terpisah)

> Fitur ini **BERDIRI SENDIRI** dan **TIDAK BERHUBUNGAN** dengan alur food recall survey.
> Dapat diakses dari **Landing Page** tanpa perlu login.
> Tujuannya adalah menjadi "food encyclopedia" interaktif yang memberikan nilai tambah kepada pengguna umum.

---

## 🎯 Konsep & Tujuan

**Find Your Food** adalah fitur pencarian makanan publik yang memungkinkan siapa saja (tanpa login) untuk:

1. **Mencari makanan** dari database Atlas Food
2. **Melihat detail lengkap** suatu makanan: nama, kategori, deskripsi
3. **Melihat informasi nilai nutrisi** per 100g dan per estimasi porsi umum
4. **Explore makanan terkait** (associated foods: nasi goreng → telur, sereal → susu)
5. **Filter berdasarkan kategori** (Makanan Pokok, Lauk Pauk, Buah, Minuman, dll)

Fitur ini memanfaatkan data yang **sudah ada** di database (`foods`, `categories`, `food_nutrients`, `nutrient_types`, `associated_foods`) — **tidak butuh tabel baru**!

---

## 📍 Entry Point dari Landing Page

Di halaman utama (Landing Page) akan ada **2 tombol/card utama**:

```text
┌─────────────────────────────────────────────────────────────────────┐
│                       🍽️  Atlas Food                                 │
│           Platform Survey Food Recall & Informasi Gizi              │
│                                                                     │
│   ┌──────────────────────────┐  ┌──────────────────────────────┐   │
│   │  📋  Food Recall Survey  │  │  🔍  Find Your Food          │   │
│   │                          │  │                              │   │
│   │  Isi survey recall       │  │  Cari makanan & cek nilai    │   │
│   │  makanan harian kamu     │  │  nutrisinya. Gratis!         │   │
│   │                          │  │                              │   │
│   │  [  Mulai Survey  ]      │  │  [  Cari Sekarang  ]         │   │
│   │   → Ke halaman login     │  │   → Tanpa login!             │   │
│   └──────────────────────────┘  └──────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

- **Tombol kiri (Food Recall Survey):** Redirect ke `/login` untuk memulai survey.
- **Tombol kanan (Find Your Food):** Redirect ke `/find-food` — publik, tanpa perlu login.

---

## 🖼️ UI Flow

### Halaman Utama: `/find-food`

```text
┌─────────────────────────────────────────────────────────────────────┐
│  🔍 Find Your Food                                                  │
│                                                                     │
│  ┌────────────────────────────────────────────────────────────┐    │
│  │  🔍  Cari makanan... (contoh: nasi goreng, pisang, susu)   │    │
│  └────────────────────────────────────────────────────────────┘    │
│                                                                     │
│  Filter Kategori:                                                   │
│  [🍚 Semua] [🍚 Makanan Pokok] [🍗 Lauk Pauk] [🥦 Sayuran]        │
│  [🍌 Buah]  [🥤 Minuman]       [🍰 Camilan]                        │
│                                                                     │
│  Makanan Populer:                                                   │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐             │
│  │ 🍚       │ │ 🍗       │ │ 🍌       │ │ 🥤       │             │
│  │ Nasi     │ │ Ayam     │ │ Pisang   │ │ Susu UHT │             │
│  │ Putih    │ │ Goreng   │ │          │ │          │             │
│  │ 130 kcal │ │ 290 kcal │ │ 89 kcal  │ │ 64 kcal  │             │
│  │ /100g    │ │ /100g    │ │ /100g    │ │ /100g    │             │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘             │
└─────────────────────────────────────────────────────────────────────┘
```

### Halaman Detail Makanan: `/find-food/[foodId]`

```text
┌─────────────────────────────────────────────────────────────────────┐
│  ← Kembali          🍚 Nasi Putih                                   │
│  ─────────────────────────────────────────────────────────────────  │
│                                                                     │
│  📝 Deskripsi:  Nasi putih matang, siap konsumsi.                  │
│  🏷️  Kategori:  Makanan Pokok                                       │
│  🔖 Kode:       NAS-001                                             │
│                                                                     │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  📊 NILAI NUTRISI                                                   │
│                                                                     │
│  Sajian: [  100g  ▼ ]  (pilih: 100g / 150g / 200g / custom)        │
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  ⚡ Energi          130 kcal    ████████░░░░░  43%           │   │
│  │  💪 Protein         2.7 g      ██░░░░░░░░░░░   9%           │   │
│  │  🌾 Karbohidrat     28.0 g     █████████░░░░  30%           │   │
│  │  🧈 Lemak Total     0.3 g      █░░░░░░░░░░░░   1%           │   │
│  │  🌿 Serat           0.4 g      █░░░░░░░░░░░░   2%           │   │
│  │  🍬 Gula            0.0 g      ░░░░░░░░░░░░░   0%           │   │
│  │  🧂 Natrium         1.0 mg     ░░░░░░░░░░░░░   0%           │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  🔗 Sering Dikonsumsi Bersama:                                      │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                            │
│  │ Telur    │ │ Ayam     │ │ Tempe    │                            │
│  │ Goreng   │ │ Goreng   │ │ Goreng   │                            │
│  │ Lihat ▶  │ │ Lihat ▶  │ │ Lihat ▶  │                            │
│  └──────────┘ └──────────┘ └──────────┘                            │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 🔌 API Backend yang Digunakan

Fitur ini **100% reuse** endpoint yang sudah ada atau sudah di-scaffold. **Tidak butuh endpoint baru!**

```http
# Sudah aktif di backend:
GET  /categories                          # Untuk filter kategori
GET  /foods/search?q=nasi&category_id=x   # Pencarian makanan (FULLTEXT)
GET  /foods/:id                           # Detail makanan

# Sudah di-scaffold (tinggal implementasi backend):
GET  /foods/:id/portion-methods           # (Opsional, bisa dipakai untuk contoh porsi)
```

> ⚠️ **Catatan penting untuk Backend:** Endpoint `GET /foods/:id` perlu di-update agar juga mengembalikan data **nutrisi** (`food_nutrients` di-join dengan `nutrient_types` & `nutrient_units`) dan **associated foods** di dalam satu response.

> ⚠️ **Auth:** Endpoint `GET /categories`, `GET /foods/search`, dan `GET /foods/:id` harus bisa diakses **tanpa token** (public route, tanpa middleware auth).

### Contoh Response `GET /foods/:id` (Target)

```json
{
  "status": "success",
  "data": {
    "id": "uuid-nasi",
    "code": "NAS-001",
    "name": "Nasi Putih",
    "local_name": "White Rice",
    "description": "Nasi putih matang, siap konsumsi.",
    "category": {
      "id": "uuid-cat-1",
      "name": "Makanan Pokok",
      "icon": "🍚"
    },
    "nutrients": [
      { "code": "energy",  "name": "Energi",       "unit": "kcal", "value_per_100g": 130.0 },
      { "code": "protein", "name": "Protein",      "unit": "g",    "value_per_100g": 2.7   },
      { "code": "carbs",   "name": "Karbohidrat",  "unit": "g",    "value_per_100g": 28.0  },
      { "code": "fat",     "name": "Lemak Total",  "unit": "g",    "value_per_100g": 0.3   },
      { "code": "fiber",   "name": "Serat",        "unit": "g",    "value_per_100g": 0.4   }
    ],
    "associated_foods": [
      { "id": "uuid-telur", "name": "Telur Goreng", "priority": 1 },
      { "id": "uuid-ayam",  "name": "Ayam Goreng",  "priority": 2 }
    ]
  }
}
```

---

## 📦 Struktur Frontend (FE)

Fitur ini extend domain `food/` yang sudah ada. Tidak perlu domain baru.

### Komponen Baru

```text
internal/domain/food/
├── components/
│   ├── FoodSearch.tsx          # (sudah ada — bisa di-reuse)
│   ├── FoodList.tsx            # (sudah ada — bisa di-reuse)
│   ├── FoodForm.tsx            # (sudah ada — untuk admin)
│   ├── FoodFinderPage.tsx      # 🆕 Halaman utama /find-food
│   ├── FoodFinderCard.tsx      # 🆕 Card makanan di grid/list
│   ├── FoodFinderDetail.tsx    # 🆕 Layout halaman detail makanan
│   └── NutritionPanel.tsx      # 🆕 Panel nutrisi dengan progress bar & selector porsi
├── hooks/
│   ├── useFoodQueries.ts       # (sudah ada — extend dengan hook baru)
│   ├── useSearchFoods.ts       # 🆕 Hook search dengan debounce 300ms
│   └── useFoodDetail.ts        # 🆕 Hook get detail + nutrisi + associated foods
├── services/
│   └── foodService.ts          # (sudah ada — tambah fungsi getFoodDetail)
└── types/
    └── food.ts                 # (sudah ada — tambah type NutritionDisplay)
```

### Routing App Router (Next.js)

```text
app/
├── page.tsx                    # Landing Page — tambah 2 tombol entry point
├── find-food/
│   ├── page.tsx                # 🆕 Halaman utama: /find-food
│   └── [foodId]/
│       └── page.tsx            # 🆕 Halaman detail: /find-food/:foodId
```

### Type Tambahan di `food.ts`

```typescript
// types/food.ts — tambahan untuk fitur Find Your Food

export interface NutrientItem {
  code: string;        // "energy", "protein", "carbs", ...
  name: string;        // "Energi", "Protein", "Karbohidrat", ...
  unit: string;        // "kcal", "g", "mg"
  value_per_100g: number;
}

export interface AssociatedFood {
  id: string;
  name: string;
  priority: number;
}

export interface FoodDetail extends Food {
  nutrients: NutrientItem[];
  associated_foods: AssociatedFood[];
}

// Untuk kalkulasi nutrisi dinamis berdasarkan gram yang dipilih user
export interface NutritionDisplay {
  nutrient: NutrientItem;
  calculated_value: number;  // (value_per_100g / 100) * selectedGram
  percentage: number;        // % dari AKG harian (opsional)
}
```

### Logic Kalkulasi Nutrisi (di Frontend)

```typescript
// Rumus: (value_per_100g / 100) × selectedGram
function calculateNutrients(
  nutrients: NutrientItem[],
  selectedGram: number
): NutritionDisplay[] {
  return nutrients.map((n) => ({
    nutrient: n,
    calculated_value: Math.round((n.value_per_100g / 100) * selectedGram * 10) / 10,
    percentage: 0, // bisa diisi dengan referensi AKG nanti
  }));
}

// Contoh penggunaan:
// User pilih 150g nasi putih
// energy: (130 / 100) × 150 = 195 kcal
// protein: (2.7 / 100) × 150 = 4.1 g
```

---

## 🗓️ Estimasi Pengerjaan

| Task | Estimasi |
| ---- | -------- |
| **BE:** Update `GET /foods/:id` — include `nutrients[]` & `associated_foods[]` | 1–2 hari |
| **BE:** Set public route (tanpa auth) untuk `/categories`, `/foods/search`, `/foods/:id` | 0.5 hari |
| **FE:** Landing Page — tambah 2 tombol entry point | 0.5 hari |
| **FE:** Halaman `/find-food` — search bar + filter kategori + grid card | 1–2 hari |
| **FE:** Halaman `/find-food/[foodId]` — detail + `NutritionPanel` + associated foods | 1–2 hari |
| **Total** | **3–5 hari** |

---

## ✅ Checklist Implementasi

### Backend

- [ ] Update `GET /foods/:id` — include `nutrients[]` & `associated_foods[]` di response
- [ ] Update `GET /foods/search` — support query param `?category_id=` untuk filter
- [ ] Set `GET /categories` sebagai public route (tanpa auth)
- [ ] Set `GET /foods/search` sebagai public route (tanpa auth)
- [ ] Set `GET /foods/:id` sebagai public route (tanpa auth)

### Frontend

- [ ] Update `app/page.tsx` (Landing Page) — tambah 2 card tombol: Food Recall & Find Your Food
- [ ] Buat `app/find-food/page.tsx` — halaman utama search & filter
- [ ] Buat `app/find-food/[foodId]/page.tsx` — halaman detail makanan
- [ ] Buat `FoodFinderCard.tsx` — card makanan dengan nama & energi singkat
- [ ] Buat `NutritionPanel.tsx` — tabel nutrisi + progress bar + dropdown selector gram (100/150/200/custom)
- [ ] Buat `FoodFinderDetail.tsx` — layout halaman detail lengkap
- [ ] Buat `useSearchFoods.ts` — React Query + debounce 300ms
- [ ] Buat `useFoodDetail.ts` — React Query untuk get detail + nutrisi + associated
- [ ] Tambahkan tipe `NutrientItem`, `AssociatedFood`, `FoodDetail` ke `food.ts`
- [ ] Implementasi logika kalkulasi nutrisi di frontend berdasarkan gram yang dipilih
- [ ] Tampilkan `associated_foods` sebagai card clickable (link ke `/find-food/[foodId]`)

---

## ❌ Scope yang Di-Skip (Fase Ini)

| Fitur | Alasan Skip |
| ----- | ----------- |
| Login/auth untuk Find Your Food | Sengaja publik, tidak perlu login |
| Simpan histori pencarian ke DB | Tidak ada akun, pakai localStorage jika perlu |
| Perbandingan makanan (A vs B) | Bisa dikembangkan di fase berikutnya |
| Rekomendasi makanan (AI-based) | Out of scope MVP |
| AKG (Angka Kecukupan Gizi) harian | Data referensi belum tersedia, bisa ditambah nanti |
| Export hasil nutrisi ke PDF | Out of scope MVP |
