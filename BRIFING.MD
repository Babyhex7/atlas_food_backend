# 🎯 Atlas Food — MVP SUPER MINIMAL

> Food recall survey platform (inspired by Intake24). Versi minimal yang bisa jalan.
> Fitur non-esensial ditandai "NANTI AJA".

---

## 📚 Penjelasan Konsep (Q&A)

### Q: Apa itu `user_survey_aliases`?

**A:** Kode/token unik buat respondent masuk ke survey setelah login.

| Contoh                          | Penjelasan                             |
| ------------------------------- | -------------------------------------- |
| `PART-A7X9K2`                   | Alias code yang di-share ke respondent |
| URL: `/survey/abc123?token=xyz` | Link unik dengan token                 |

**Gunanya:** Respondent login dulu → klik link → isi survey.

---

### Q: Apa itu `language_translations`? Pakai i18n/next-intl?

**A:** Tabel ini untuk translate teks di admin/dashboard.

**Tapi buat MVP:**

- ✅ **Pakai hardcode JSON** dulu di frontend
- ✅ atau **next-intl** langsung tanpa tabel DB
- ❌ Skip tabel `language_translations` dulu (NANTI AJA)

```typescript
// i18n/messages/id.json (hardcode dulu)
{
  "survey": {
    "welcome": "Selamat datang di survey gizi",
    "searchFood": "Cari makanan...",
    "portionSize": "Perkiraan porsi (gram)"
  }
}
```

---

### Q: Bedanya `nutrient_unit` vs `nutrient_types`?

**A:**

| Tabel            | Fungsi                | Contoh                            |
| ---------------- | --------------------- | --------------------------------- |
| `nutrient_types` | **APA** yang diukur   | Protein, Lemak, Energi, Vitamin C |
| `nutrient_units` | **SATUAN** pengukuran | gram (g), miligram (mg), kcal     |

**Relasi:**

```
nutrient_types: "Protein" ──unit_id──▶ nutrient_units: "gram (g)"
```

**✅ MVP V2:** Pakai tabel terpisah (`nutrient_units`, `nutrient_types`, `food_nutrients`) untuk perhitungan akurat.

---

### Q: Bedanya `food_category` vs `category`?

**A:**

| Tabel             | Fungsi                                                  |
| ----------------- | ------------------------------------------------------- |
| `categories`      | Master kategori ("Nasi", "Sayur", "Daging")             |
| `food_categories` | Many-to-many linker (Nasi Goreng = "Nasi" + "Gorengan") |

**Kenapa many-to-many?** 1 makanan bisa masuk multiple kategori.

---

### Q: Penjelasan tabel-tabel kompleks?

#### 1. `nutrient_table_records` vs `nutrient_table_record_nutrients`

```
food: "Nasi Putih"
  └──▶ nutrient_table_record: "USDA-20450" (kode dari USDA DB)
          └──▶ nutrient_table_record_nutrients:
                ├── protein: 2.6g
                ├── carbs: 28g
                └── energy: 130kcal
```

**MVP V2:** ✅ **PAKAI** tabel `food_nutrients` + `nutrient_types` + `nutrient_units` (perhitungan akurat & scalable).

#### 2. `food_nutrients`

Simpan nilai gizi per makanan per 100g.
**MVP V2:** ✅ **PAKAI** - simpan di `food_nutrients` dengan struktur: `(food_id, nutrient_type_id, value_per_100g)`.

#### 3. `associated_foods` (sereal → susu)

Fitur "makanan yang sering dipakai bersama" (contoh: sereal → susu).
**MVP V2:** ✅ **PAKAI** - tabel `associated_foods` untuk UX lebih baik.

#### 4. `food_attributes` vs `food_attribute_values`

Custom properties kayak "brand", "is_organic", "halal", etc.
**MVP:** NANTI AJA, pakai JSON di tabel `foods`.

---

## ✅ Alur MVP dengan Portion Images (Lengkap & Benar!)

### ADMIN FLOW

```
Login ──▶ Create Survey ──▶ Upload Food DB ──▶ Set Portion Images ──▶ Share Link/Token
```

### RESPONDENT FLOW (4 Step Saja - Final Version!)

```
Link ──▶ Step 1 ──▶ Step 2 ──▶ Step 3 ──▶ Step 4 ──▶ Done
         Pilih      Add        Portion   Review
         Waktu      Food       (Langsung  & Submit
         Makan      & Drink    Gambar!)

┌────────────────────────────────────────────────────────────────────────┐
│ Step 1: PILIH/EDIT WAKTU MAKAN                                         │
│ ┌────────────────────────────────────────────────────────────┐        │
│ │  Sarapan (07:00) ▶   [Edit Meal]        [? HELP & FAQ]       │        │
│ └────────────────────────────────────────────────────────────┘        │
│                                                                         │
│ Step 2: ADD MAKANAN & MINUMAN (List dulu, portion nanti!)             │
│ ┌────────────────────────────────────────────────────────────┐        │
│ │  What did you have for Sarapan?                            │        │
│ │                                                           │        │
│ │  🍔 Foods                        🥤 Drinks                │        │
│ │  ┌──────────────┐ ┌────┐        ┌──────────────┐ ┌────┐   │        │
│ │  │ [Nasi Goreng]│ │ ADD│        │ [Es Teh     ]│ │ ADD│   │        │
│ │  └──────────────┘ └────┘        └──────────────┘ └────┘   │        │
│ │                                                           │        │
│ │  ✅ Nasi Goreng (akan di-portion-kan nanti)              │        │
│ │  ✅ Telur Goreng (akan di-portion-kan nanti)             │        │
│ │  ✅ Es Teh (akan di-portion-kan nanti)                   │        │
│ │                                                           │        │
│ │  [🕐 Change Meal Time] [🗑️ Delete Meal] [▶ Continue]      │        │
│ └────────────────────────────────────────────────────────────┘        │
└────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│ Step 3: PORTION SELECTION (2 Tipe UI - Sesuai Jenis Makanan!)           │
│                                                                         │
│  "How much did you have?" (1 of 3 foods)                               │
│                                                                         │
│  🔥 TIPE A: Simple Grid (Nasi, Nuggets, dll)                           │
│  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐  ┌─────────────────┐         │
│  │ 🖼️  │ │ 🖼️  │ │ 🖼️  │ │ 🖼️  │ │ 🖼️  │  │  Atau ketik:   │         │
│  │ ½   │ │ 1   │ │ 1½  │ │ 2   │ │ 3   │  │  [___150___] g │         │
│  │porsi│ │porsi│ │porsi│ │porsi│ │porsi│  │                 │         │
│  │~75g │ │~150g│ │~225g│ │~300g│ │~450g│  │ [  ✓ CONFIRM  ] │         │
│  └─────┘ └─────┘ └─────┘ └─────┘ └─────┘  └─────────────────┘         │
│                                                                         │
│  🔥 TIPE B: As Served + Quantity (Pisang, Kentang, dll)                │
│  ┌──────────────────────────────────────────────────────────────────┐   │
│  │                                                                  │   │
│  │     [🍌 Gambar besar pisang di piring (selected)]              │   │
│  │                                                                  │   │
│  │        ┌──────┐           ┌──────┐                              │   │
│  │        │  +   │     and   │  +   │   ← overlay di gambar       │   │
│  │        │  3   │           │  ¼   │                              │   │
│  │        │  -   │           │  -   │                              │   │
│  │        └──────┘           └──────┘                              │   │
│  │       WHOLE              FRACTION                                │   │
│  │                                                                  │   │
│  │  "I had 3 and ¼ of the largest portion (617.50g)"              │   │
│  │                                                                  │   │
│  └──────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  Thumbnail row: [🖼️] [🖼️] [🖼️] [🖼️] [🖼️] [🖼️] [🖼️] [🖼️👆selected] │
│                20g   40g   60g   95g  130g  160g  175g  190g            │
│                                                                         │
│  [    ✓ I HAD THAT MUCH    ]                                            │
│  [◀ Previous Food]  [Next Food ▶]  [Skip This Food]                     │
└────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│ Step 4: REVIEW & SUBMIT                                                 │
│                                                                         │
│  📋 Summary for damn (08:00)                                            │
│  ├─ 🍔 Nasi Goreng - 1.5 porsi (~225g)                                  │
│  ├─ 🍔 Telur Goreng - 1 porsi (~60g)                                    │
│  └─ 🥤 Es Teh - 1 gelas besar (~250ml)                                  │
│                                                                         │
│  [✏️ Edit Portions]  [➕ Add More Food]  [✅ SUBMIT MEAL]                │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 🗄️ TABEL MVP (16 TABEL)

### 1. Auth (Minimal)

```sql
-- users (Admin + Respondent 1 tabel)
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role ENUM('admin', 'respondent') DEFAULT 'respondent',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email)
);

-- refresh_tokens (JWT)
CREATE TABLE refresh_tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

**SKIP (NANTI AJA):** roles, permissions, role_user, permission_user, RBAC kompleks.

---

### 2. Survey (Minimal)

```sql
-- locales (hardcode saja: id, en)
CREATE TABLE locales (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,  -- 'id', 'en'
    name VARCHAR(50) NOT NULL
);

-- surveys
CREATE TABLE surveys (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    slug VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,

    -- 🔥 Meals sebagai JSON (simplifed, ga perlu tabel terpisah)
    meals_config JSON NOT NULL,  -- [{"name":"Breakfast","time":"07:00"}, ...]

    -- 🔥 Prompts sebagai JSON
    prompts JSON,  -- {"welcome":"...","search":"..."}

    locale_id INT DEFAULT 1,
    start_date DATE,
    end_date DATE,
    status ENUM('draft', 'active', 'closed') DEFAULT 'draft',

    -- 🔥 Token untuk akses survey
    access_token VARCHAR(255),

    created_by CHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- user_survey_aliases (Simple - untuk respondent terdaftar)
CREATE TABLE survey_participants (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    survey_id CHAR(36) NOT NULL,
  user_id CHAR(36) NOT NULL,
  alias VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

**SKIP (NANTI AJA):**

- `survey_schemes` (pakai JSON di surveys)
- `user_survey_sessions` (save di frontend localStorage dulu)
- `survey_scheme_prompts`
- `language_translations` (pakai next-intl file JSON)

---

### 3. Food Database (Minimal)

```sql
-- categories (simple, flat dulu - ga perlu hierarki)
CREATE TABLE categories (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    icon VARCHAR(50),  -- emoji atau icon name
    display_order INT DEFAULT 0
);

-- foods (tanpa nutrisi JSON, pake tabel terpisah)
CREATE TABLE foods (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    local_name VARCHAR(255),
    description TEXT,

    category_id CHAR(36),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    INDEX idx_name (name),
    FULLTEXT INDEX ft_name (name, local_name)  -- untuk search
);

-- food_categories (many-to-many)
CREATE TABLE food_categories (
    food_id CHAR(36) NOT NULL,
    category_id CHAR(36) NOT NULL,
    PRIMARY KEY (food_id, category_id),
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- 🔥 nutrient_units (satuan: g, mg, kcal, IU)
CREATE TABLE nutrient_units (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,      -- 'g', 'mg', 'kcal', 'IU'
    name VARCHAR(50) NOT NULL,             -- 'gram', 'miligram'
    symbol VARCHAR(10) NOT NULL            -- 'g', 'mg'
);

-- 🔥 nutrient_types (jenis nutrisi)
CREATE TABLE nutrient_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(30) NOT NULL UNIQUE,      -- 'energy', 'protein', 'carbs'
    name VARCHAR(100) NOT NULL,            -- 'Energi', 'Protein'
    unit_id INT NOT NULL,
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (unit_id) REFERENCES nutrient_units(id),
    INDEX idx_code (code)
);

-- 🔥 food_nutrients (nilai gizi per makanan)
CREATE TABLE food_nutrients (
    food_id CHAR(36) NOT NULL,
    nutrient_type_id INT NOT NULL,
    value_per_100g DECIMAL(10, 4) NOT NULL,  -- nilai per 100g
    PRIMARY KEY (food_id, nutrient_type_id),
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    FOREIGN KEY (nutrient_type_id) REFERENCES nutrient_types(id) ON DELETE CASCADE
);

-- 🔥 associated_foods (makanan yang sering dipakai bersama)
CREATE TABLE associated_foods (
    id INT AUTO_INCREMENT PRIMARY KEY,
    food_id CHAR(36) NOT NULL,             -- makanan utama
    associated_food_id CHAR(36) NOT NULL,  -- makanan terkait
    priority INT DEFAULT 0,                -- urutan prioritas (1 = paling sering)
    is_default BOOLEAN DEFAULT FALSE,      -- auto-tambah ke recall?
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    FOREIGN KEY (associated_food_id) REFERENCES foods(id) ON DELETE CASCADE,
    UNIQUE KEY unique_association (food_id, associated_food_id)
);
```

**SKIP (NANTI AJA):**

- `nutrient_tables` (multiple sources USDA/BPOM)
- `nutrient_table_records`, `nutrient_table_record_nutrients`
- `synonym_sets` (search pakai FULLTEXT dulu)
- `food_attributes`, `food_attribute_values`

---

### 4. Portion Size (🔥 dengan Gambar!)

```sql
-- 🔥 Portion Size Methods per Food (dengan gambar!)
CREATE TABLE food_portion_size_methods (
    id INT AUTO_INCREMENT PRIMARY KEY,
    food_id CHAR(36) NOT NULL,

    -- Tipe metode
    method_type ENUM('as_served', 'guide_image', 'weight') NOT NULL,  -- drinkware: SKIP MVP

    -- Label yang ditampilkan (contoh: "In a glass", "In a bottle", "1 nugget")
    label VARCHAR(255) NOT NULL,
    description VARCHAR(255),

    -- 🔥 Gambar untuk metode ini
    image_url VARCHAR(500),
    thumbnail_url VARCHAR(500),

    -- Konfigurasi pilihan (JSON)
    /*
    🔥 TIPE SELECTION:

    1. "simple_grid" - Grid gambar ½,1,1½,2,3 (Nasi, Nuggets)
    {
      "selectionType": "simple_grid",
      "options": [
        {"value": 0.5, "label": "½ porsi", "gram": 75, "image": "/nasi-0.5.jpg"},
        {"value": 1, "label": "1 porsi", "gram": 150, "image": "/nasi-1.jpg"},
        {"value": 1.5, "label": "1½ porsi", "gram": 225, "image": "/nasi-1.5.jpg"},
        {"value": 2, "label": "2 porsi", "gram": 300, "image": "/nasi-2.jpg"}
      ],
      "allowFractions": false
    }

    2. "as_served_quantity" - Thumbnail row + overlay counter (Pisang, Kentang!)
    {
      "selectionType": "as_served_quantity",
      "setCode": "banana-slices",
      "thumbnailPosition": "bottom",
      "maxQuantity": 5,               // Maks 5x dari gambar terbesar
      "allowFractions": true,
      "fractionOptions": [0, 0.25, 0.5, 0.75],
      "defaultQuantity": 1,
      "defaultFraction": 0,
      "showCalculation": true         // Tampilkan "3 and ¼ of largest (617.50g)"
    }

    3. "counter" - WHOLE + FRACTION tanpa gambar (Pack, Sandwich)
    {
      "selectionType": "counter",
      "baseWeight": 200,
      "baseUnit": "pack",
      "wholeLabel": "WHOLE",
      "fractionLabel": "FRACTION",
      "allowFractions": true,
      "fractionOptions": [0, 0.25, 0.5, 0.75]
    }

    4. "input" - Ketik manual gram
    {
      "selectionType": "input",
      "unit": "gram",
      "placeholder": "Enter grams..."
    }
    */
    config JSON,

    -- Urutan tampilan
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    INDEX idx_food (food_id)
);

-- 🔥 As Served Sets (Foto porsi - contoh: nugget 1,2,3 atau gelas 1,2,3,4,5,6)
CREATE TABLE as_served_sets (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,           -- "Chicken Nuggets - Portion Photos"
    description TEXT,
    food_id CHAR(36),                      -- Bisa spesifik ke 1 food atau general
    category VARCHAR(50),                  -- "nuggets", "glasses", "bottles"
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE SET NULL
);

-- 🔥 As Served Images (Foto individual per porsi)
CREATE TABLE as_served_images (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    set_id CHAR(36) NOT NULL,

    -- Urutan/label (1, 2, 3 atau "small", "medium", "large")
    label VARCHAR(50) NOT NULL,            -- "1", "2", "3", "Small", "Medium"

    -- 🔥 Gambar
    image_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),

    -- Berat dalam gram untuk porsi ini
    weight_gram DECIMAL(10, 2) NOT NULL,

    -- Deskripsi tambahan
    description VARCHAR(255),                -- "1 chicken nugget, ~20g"

    display_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (set_id) REFERENCES as_served_sets(id) ON DELETE CASCADE,
    INDEX idx_set (set_id)
);

-- Note: drinkware_sets & drinkware_scales DI-SKIP untuk MVP
-- Drinks pakai as_served_sets aja (gambar gelas/botol + quantity)
```

**SKIP untuk MVP:** `drinkware_sets`, `drinkware_scales` - minuman pakai `as_served` aja!

**Contoh Data Portion Methods:**

````sql
-- 1. Nuggets dengan foto porsi
INSERT INTO food_portion_size_methods (food_id, method_type, label, description, image_url, config) VALUES
('uuid-nuggets', 'as_served', 'Chicken Nuggets', 'Select how many nuggets you had', '/nuggets-preview.jpg', '
{
  "selectionType": "image",
  "setCode": "nuggets-portions",
  "allowFractions": true,
  "fractionOptions": [0.25, 0.5, 0.75, 1, 1.5, 2, 2.5, 3]
}
');

-- 2. Minuman dengan pilihan botol/gelas/carton
INSERT INTO food_portion_size_methods (food_id, method_type, label, description, image_url, config) VALUES
('uuid-soda', 'as_served', 'In a bottle', 'Choose bottle size', '/bottles-preview.jpg', '
{
  "selectionType": "as_served_quantity",
  "setCode": "soft-drink-bottles",
  "maxQuantity": 3,
  "allowFractions": true
}
'),
('uuid-soda', 'as_served', 'In a glass', 'Choose glass size', '/glasses-preview.jpg', '
{
  "selectionType": "as_served_quantity",
  "setCode": "standard-glasses",
  "maxQuantity": 3,
  "allowFractions": false
}
');

-- 3. Counter method (+/-) - untuk item per pack/unit
INSERT INTO food_portion_size_methods (food_id, method_type, label, description, image_url, config) VALUES
('uuid-toddler-meal', 'standard_portion', 'Individual packs', 'How many packs did you have?', '/packs/pack-preview.jpg', '
{
  "selectionType": "counter",
  "baseWeight": 200,
  "baseUnit": "pack",
  "wholeLabel": "WHOLE",
  "fractionLabel": "FRACTION",
  "allowFractions": true,
  "fractionOptions": [0, 0.25, 0.5, 0.75],
  "defaultWhole": 1,
  "defaultFraction": 0
}
'),
('uuid-sandwich', 'standard_portion', 'Sandwiches', 'How many sandwiches?', '/sandwich/sandwich-preview.jpg', '
{
  "selectionType": "counter",
  "baseWeight": 150,
  "baseUnit": "sandwich",
  "wholeLabel": "WHOLE",
  "fractionLabel": "FRACTION",
  "allowFractions": true,
  "fractionOptions": [0, 0.25, 0.5, 0.75, 1],
  "defaultWhole": 1,
  "defaultFraction": 0
}
');

-- 4. As Served with Quantity Multiplier (🔥 Kaya screenshot pisang!)
INSERT INTO food_portion_size_methods (food_id, method_type, label, description, image_url, config) VALUES
('uuid-banana', 'as_served', 'Sliced banana on plate', 'Choose portion and adjust quantity', '/banana/preview.jpg', '
{
  "selectionType": "as_served_quantity",
  "setCode": "banana-slices",
  "thumbnailPosition": "bottom",
  "maxQuantity": 5,
  "allowFractions": true,
  "fractionOptions": [0, 0.25, 0.5, 0.75],
  "defaultQuantity": 1,
  "defaultFraction": 0,
  "showCalculation": true
}
'),
('uuid-rice', 'as_served', 'Rice on plate', 'Choose portion size', '/rice/preview.jpg', '
{
  "selectionType": "as_served_quantity",
  "setCode": "rice-plates",
  "thumbnailPosition": "bottom",
  "maxQuantity": 3,
  "allowFractions": true,
  "fractionOptions": [0, 0.5],
  "defaultQuantity": 1,
  "defaultFraction": 0
}
');

-- 5. Manual weight (fallback)
INSERT INTO food_portion_size_methods (food_id, method_type, label, description, config) VALUES
('uuid-nuggets', 'weight', 'Enter weight manually', 'If you know the exact weight', '
{
  "selectionType": "input",
  "unit": "gram",
  "placeholder": "Enter grams..."
}
');

**Contoh Data Images:**

```sql
-- As Served Set untuk Nuggets
INSERT INTO as_served_sets (id, code, name, description, category) VALUES
('uuid-set-1', 'nuggets-portions', 'Chicken Nuggets Portion Guide', 'Visual guide for nugget portions', 'nuggets');

-- Gambar per porsi
INSERT INTO as_served_images (set_id, label, image_url, thumbnail_url, weight_gram, description) VALUES
('uuid-set-1', '1', '/nuggets/nugget-1.jpg', '/nuggets/nugget-1-thumb.jpg', 20.0, '1 nugget (~20g)'),
('uuid-set-1', '2', '/nuggets/nugget-2.jpg', '/nuggets/nugget-2-thumb.jpg', 40.0, '2 nuggets (~40g)'),

-- 🔥 As Served Set untuk Pisang (Kaya Screenshot!)
INSERT INTO as_served_sets (id, code, name, description, category, food_id) VALUES
('uuid-set-3', 'banana-slices', 'Sliced Banana Portions', 'Visual guide for banana portions', 'fruits', 'uuid-banana');

-- 8 gambar porsi dari kecil ke besar (kaya screenshot)
INSERT INTO as_served_images (set_id, label, image_url, thumbnail_url, weight_gram, description, display_order) VALUES
('uuid-set-3', '1', '/banana/banana-1.jpg', '/banana/banana-1-thumb.jpg', 20.0, 'Few slices (~20g)', 1),
('uuid-set-3', '2', '/banana/banana-2.jpg', '/banana/banana-2-thumb.jpg', 40.0, 'Small portion (~40g)', 2),
('uuid-set-3', '3', '/banana/banana-3.jpg', '/banana/banana-3-thumb.jpg', 60.0, 'Medium-small (~60g)', 3),
('uuid-set-3', '4', '/banana/banana-4.jpg', '/banana/banana-4-thumb.jpg', 95.0, 'Medium (~95g)', 4),
('uuid-set-3', '5', '/banana/banana-5.jpg', '/banana/banana-5-thumb.jpg', 130.0, 'Medium-large (~130g)', 5),
('uuid-set-3', '6', '/banana/banana-6.jpg', '/banana/banana-6-thumb.jpg', 160.0, 'Large (~160g)', 6),
('uuid-set-3', '7', '/banana/banana-7.jpg', '/banana/banana-7-thumb.jpg', 175.0, 'Very large (~175g)', 7),
('uuid-set-3', '8', '/banana/banana-8.jpg', '/banana/banana-8-thumb.jpg', 190.0, 'Full plate (~190g)', 8),
('uuid-set-1', '3', '/nuggets/nugget-3.jpg', '/nuggets/nugget-3-thumb.jpg', 60.0, '3 nuggets (~60g)'),
('uuid-set-1', '4', '/nuggets/nugget-4.jpg', '/nuggets/nugget-4-thumb.jpg', 80.0, '4 nuggets (~80g)'),
('uuid-set-1', '5', '/nuggets/nugget-5.jpg', '/nuggets/nugget-5-thumb.jpg', 100.0, '5 nuggets (~100g)');

-- 🔥 Contoh Data Nutrisi (Baru!)
-- 1. Nutrient Units
INSERT INTO nutrient_units (code, name, symbol) VALUES
('kcal', 'Kilokalori', 'kcal'),
('g', 'Gram', 'g'),
('mg', 'Miligram', 'mg'),
('mcg', 'Mikrogram', 'μg');

-- 2. Nutrient Types
INSERT INTO nutrient_types (code, name, unit_id, display_order) VALUES
('energy', 'Energi', 1, 1),           -- kcal
('protein', 'Protein', 2, 2),         -- g
('carbs', 'Karbohidrat', 2, 3),      -- g
('fat', 'Lemak Total', 2, 4),         -- g
('fiber', 'Serat', 2, 5),             -- g
('sugar', 'Gula', 2, 6),              -- g
('sodium', 'Natrium', 3, 7),          -- mg
('calcium', 'Kalsium', 3, 8),        -- mg
('iron', 'Zat Besi', 3, 9),           -- mg
('vit_c', 'Vitamin C', 3, 10);       -- mg

-- 3. Contoh Data Foods dengan Nutrisi
INSERT INTO categories (id, code, name, icon) VALUES
('uuid-cat-1', 'staples', 'Makanan Pokok', '🍚'),
('uuid-cat-2', 'protein', 'Lauk Pauk', '🍗'),
('uuid-cat-3', 'fruits', 'Buah-buahan', '🍌'),
('uuid-cat-4', 'drinks', 'Minuman', '🥤');

INSERT INTO foods (id, code, name, local_name, description, category_id) VALUES
('uuid-nasi', 'Nasi-001', 'Nasi Putih', 'White Rice', 'Nasi putih matang', 'uuid-cat-1'),
('uuid-ayam', 'Ayam-001', 'Ayam Goreng', 'Fried Chicken', 'Ayam goreng tepung', 'uuid-cat-2'),
('uuid-pisang', 'Pisang-001', 'Pisang', 'Banana', 'Pisang cavendish', 'uuid-cat-3'),
('uuid-susu', 'Susu-001', 'Susu UHT', 'UHT Milk', 'Susu UHT full cream', 'uuid-cat-4');

-- 4. Food Nutrients (per 100g)
INSERT INTO food_nutrients (food_id, nutrient_type_id, value_per_100g) VALUES
-- Nasi Putih
('uuid-nasi', 1, 130.00),  -- energy: 130 kcal
('uuid-nasi', 2, 2.70),    -- protein: 2.7g
('uuid-nasi', 3, 28.00),   -- carbs: 28g
('uuid-nasi', 4, 0.30),    -- fat: 0.3g
('uuid-nasi', 5, 0.40),    -- fiber: 0.4g
-- Ayam Goreng
('uuid-ayam', 1, 290.00),  -- energy: 290 kcal
('uuid-ayam', 2, 18.00),   -- protein: 18g
('uuid-ayam', 4, 19.00),   -- fat: 19g
-- Pisang
('uuid-pisang', 1, 89.00), -- energy: 89 kcal
('uuid-pisang', 2, 1.10),  -- protein: 1.1g
('uuid-pisang', 3, 22.80), -- carbs: 22.8g
('uuid-pisang', 6, 12.20), -- sugar: 12.2g
-- Susu UHT
('uuid-susu', 1, 64.00),   -- energy: 64 kcal
('uuid-susu', 2, 3.40),    -- protein: 3.4g
('uuid-susu', 8, 120.00);  -- calcium: 120mg

-- 5. Food Categories (many-to-many)
INSERT INTO food_categories (food_id, category_id) VALUES
('uuid-nasi', 'uuid-cat-1'),
('uuid-ayam', 'uuid-cat-2'),
('uuid-pisang', 'uuid-cat-3'),
('uuid-susu', 'uuid-cat-4');

-- 🔥 Associated Foods (Sereal → Susu!)
INSERT INTO associated_foods (food_id, associated_food_id, priority, is_default) VALUES
-- Makan sereal, biasanya minum susu
('uuid-sereal', 'uuid-susu', 1, TRUE),
-- Makan roti, biasanya pakai selai atau susu
('uuid-roti', 'uuid-susu', 1, FALSE),
('uuid-roti', 'uuid-selai', 2, FALSE),
-- Nasi goreng → telur
('uuid-nasi-goreng', 'uuid-telur', 1, FALSE);

-- 🔌 API MVP (Minimal)

**1. Image Method (As Served - untuk Food & Drinks):**

```javascript
// User pilih: 1.5 porsi dari image "nugget-3" (60g per porsi)
const portionGram = selectedImage.weight_gram * quantity;
// 60g * 1.5 = 90g
````

**2. Counter Method (+/-) - Sesuai Screenshot!:**

```javascript
// User pilih: 1 WHOLE + 0 FRACTION
// Base: 1 pack = 200g
const totalUnits = whole + fraction; // 1 + 0 = 1
const portionGram = totalUnits * baseWeight; // 1 * 200g = 200g

// User pilih: 1 WHOLE + 0.5 FRACTION (½)
const totalUnits = 1 + 0.5; // = 1.5
const portionGram = 1.5 * 200; // = 300g

// User pilih: 2 WHOLE + 0.75 FRACTION (¾)
const totalUnits = 2 + 0.75; // = 2.75
const portionGram = 2.75 * 200; // = 550g
```

**3. Weight Method (Manual Input):**

```javascript
// User ketik langsung: 150 gram
const portionGram = userInput; // 150g
```

**4. As Served with Quantity (🔥 Kaya Screenshot Pisang!):**

```javascript
// User pilih: gambar ke-8 (full plate = 190g) + quantity 3.25 (3 and ¼)
const selectedImage = {
  label: "8",
  weight_gram: 190.0, // dari as_served_images
  description: "Full plate (~190g)",
};

const quantity = 3; // WHOLE
const fraction = 0.25; // FRACTION (¼)
const totalQuantity = quantity + fraction; // = 3.25

// RUMUS: image_weight × totalQuantity
const portionGram = selectedImage.weight_gram * totalQuantity;
// 190g × 3.25 = 617.5g  ✓ Sama kayak screenshot!

// Tampilkan ke user:
// "I had 3 and ¼ of the largest portion (617.50g)"
```

**Data yang disimpan di submission:**

```json
{
  "foodId": "uuid-banana",
  "foodName": "Pisang Iris",
  "portion": {
    "method": "as_served_quantity",
    "imageId": "uuid-img-8",
    "imageLabel": "8",
    "baseWeight": 190,
    "quantity": 3,
    "fraction": 0.25,
    "totalQuantity": 3.25,
    "portionGram": 617.5
  },
  "nutrients": {
    "energy": 617.5,
    "carbs": 158.3,
    "protein": 7.5
  }
}
```

---

### 5. Submission (Minimal)

```sql
-- survey_submissions (🔥 1 tabel simpan semua)
CREATE TABLE survey_submissions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    survey_id CHAR(36) NOT NULL,
    participant_id CHAR(36),  -- dari survey_participants

    -- 🔥 Data respondent (snapshot, ga perlu join)
    respondent_name VARCHAR(255),
    respondent_email VARCHAR(255),

    -- 🔥 Semua meals & foods simpan di JSON
    meals_data JSON NOT NULL,
    /*
    [
      {
        "name": "Breakfast",
        "time": "08:30",
        "foods": [
          {
            "foodId": "uuid",
            "foodName": "Nasi Putih",
            "portionGram": 150,
            "nutrients": {"energy":195,"protein":3.9}
          }
        ]
      }
    ]
    */

    -- Missing foods (kalau user ketik makanan ga ketemu di DB)
    missing_foods JSON,

    submitted_at TIMESTAMP NOT PNULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE,
    FOREIGN KEY (participant_id) REFERENCES survey_participants(id) ON DELETE SET NULL
);
```

**SKIP (NANTI AJA):**

- `survey_submission_meals` (tabel terpisah)
- `survey_submission_foods` (tabel terpisah)
- `survey_submission_missing_foods` (tabel terpisah)
- `user_survey_sessions` (save di frontend localStorage)

---

## 🔌 API MVP (Minimal)

### Auth

```
POST /auth/register          # Respondent daftar
POST /auth/login             # Login
POST /auth/refresh           # Refresh token
GET  /auth/me                # Profile
```

### Admin

```
GET    /admin/surveys              # List surveys
POST   /admin/surveys              # Create survey
GET    /admin/surveys/:id          # Detail survey
PUT    /admin/surveys/:id          # Update survey
DELETE /admin/surveys/:id          # Delete survey

GET    /admin/surveys/:id/submissions   # List submissions
GET    /admin/surveys/:id/export        # Export CSV/JSON

GET    /admin/foods                # List foods
POST   /admin/foods                # Create food
PUT    /admin/foods/:id            # Update food
DELETE /admin/foods/:id            # Delete food

# 🔥 Portion Size Management
GET    /admin/foods/:id/portion-methods          # List portion methods per food
POST   /admin/foods/:id/portion-methods          # Add portion method
PUT    /admin/portion-methods/:id                # Update portion method
DELETE /admin/portion-methods/:id                # Delete portion method

GET    /admin/as-served-sets                     # List as served sets
POST   /admin/as-served-sets                     # Create as served set
POST   /admin/as-served-sets/:id/images          # Upload portion images

# Note: Drinkware di-skip untuk MVP - drinks pakai as_served_sets aja
```

### Public/Respondent

```
POST /survey/access           # Akses dengan token/alias
GET  /survey/:id/info         # Info survey

GET  /foods/search?q=nasi     # Search makanan
GET  /foods/:id               # Detail makanan

# 🔥 Portion Selection APIs
GET  /foods/:id/portion-methods                  # Get portion methods (dengan gambar!)
GET  /portion-methods/:id/options                # Get options (as_served images)
GET  /as-served-sets/:code/images                # Get portion images (food & drinks)

GET  /categories              # List kategori

POST /survey/submit         # Submit final (meals di JSON)
```

**SKIP (NANTI AJA):**

- Session save/load (pakai localStorage frontend)
- Social login
- Password reset email

---

## 🖼️ UI Flow dengan Portion Selection (Baru!)

```
RESPONDENT JOURNEY dengan Gambar Porsi:
═══════════════════════════════════════════════════════════════

┌─────────────────────────────────────────────────────────────┐
│ Step 1: PILIH WAKTU MAKAN                                    │
│ [Sarapan ▼]  Jam: [08:30 ▼]                                  │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 2: CARI MAKANAN                                        │
│ 🔍 [nasi goreng... ]                                       │
│   ┌─────────┐ ┌─────────┐ ┌─────────┐                      │
│   │Nasi     │ │Nasi     │ │Nasi     │                      │
│   │Goreng   │ │Goreng   │ │Goreng   │                      │
│   │Ayam     │ │Ikan     │ │Spesial  │                      │
│   └────┬────┘ └─────────┘ └─────────┘                      │
│        │                                                     │
│        ▼ (User klik "Nasi Goreng Ayam")                     │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 3: PILIH METODE PORSI                                  │
│ "How do you want to estimate your portion?"                  │
│                                                              │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│   │ [FOTO]      │  │ [FOTO]      │  │ [FOTO]      │         │
│   │ Piring      │  │ Gelas       │  │ Sendok      │         │
│   │ Penuh       │  │ Besar       │  │ Nasi        │         │
│   │             │  │             │  │             │         │
│   │ "In a       │  │ "In a       │  │ "Enter      │         │
│   │  plate"     │  │  glass"     │  │  weight"    │         │
│   └─────────────┘  └─────────────┘  └─────────────┘         │
│        (as_served)      (as_served)      (weight)          │  <-- drinks juga pakai as_served!
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 4: PILIH JUMLAH/PORSI (dengan Gambar!)                 │
│ "How many nuggets did you have?"                           │
│                                                              │
│   ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐                 │
│   │ 🖼️  │ │ 🖼️  │ │ 🖼️  │ │ 🖼️  │ │ 🖼️  │                 │
│   │ 1   │ │ 2   │ │ 3   │ │ 4   │ │ 5   │                 │
│   │~20g │ │~40g │ │~60g │ │~80g │ │~100g│                 │
│   └─────┘ └─────┘ └─────┘ └─────┘ └─────┘                 │
│    [Selected ✓]                                             │
│                                                              │
│   Atau + Fraksi:                                             │
│   [1] whole + [⅓] fraction = 1⅓ nuggets (~27g)             │
│                                                              │
│   [       I HAD THAT MANY        ]                         │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 5: TAMBAH MAKANAN LAIN / NEXT MEAL                      │
│ [+ Tambah Makanan]  atau  [Lanjut ke Makan Siang ▶]         │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 6: REVIEW & SUBMIT                                      │
│                                                              │
│  Sarapan (08:30)                                            │
│  ├─ Nasi Goreng Ayam - 1.3 porsi (~27g)                    │
│  ├─ Telur Goreng - 1 porsi (~60g)                          │
│  └─ Susu - 1 gelas besar (~250ml)                          │
│                                                              │
│  Makan Siang (12:15)                                        │
│  ├─ Nasi Putih - 2 porsi (~200g)                           │
│  └─ ...                                                     │
│                                                              │
│  [✏️ Edit]  [🗑️ Hapus]  [✅ SUBMIT SURVEY]                  │
└─────────────────────────────────────────────────────────────┘
```

---

## 📊 ERD MVP (Super Simple)

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│    users    │      │   surveys   │      │ survey_part │
├─────────────┤      ├─────────────┤      ├─────────────┤
│ PK id       │──┐   │ PK id       │◀─┐   │ PK id       │
│ email       │  │   │ slug        │  │   │ FK survey_id│
│ role        │  │   │ meals(JSON) │  └───┤ FK user_id  │
│ password    │  │   │ prompts(JSON)    │  │ alias       │
└─────────────┘  │   │ access_token     │  └──────┬──────┘
     │           │   └─────────────┘         │
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
         │ missing_foods      │ category_id │──▶┌─────────────┐
         └─────────────┘     └─────────────┘   │ categories  │
                                               ├─────────────┤
                                               │ PK id       │
                                               │ name        │
                                               └─────────────┘
```

---

## 🛠️ Backend Technical Stack

| Layer             | Technology                       | Purpose             |
| ----------------- | -------------------------------- | ------------------- |
| **Language**      | Go (Golang)                      | Backend language    |
| **Framework**     | Gin                              | HTTP web framework  |
| **Database**      | MySQL                            | Primary datastore   |
| **ORM**           | GORM                             | Database operations |
| **Auth**          | JWT (stateless) + Refresh Tokens | Authentication      |
| **Validation**    | go-playground/validator          | Request validation  |
| **Password Hash** | bcrypt                           | Password encryption |

### Project Structure (Go - Clean Architecture)

```
atlas_food/
├── cmd/
│   └── api/                     # Entry point API server
│       └── main.go
│   └── worker/                  # (NANTI) Background worker
│       └── main.go
├── internal/
│   ├── config/                  # Config & DB connection
│   │   ├── config.go
│   │   └── database.go
│   │
│   ├── domain/                  # 🔥 Core business logic (modular per fitur)
│   │   ├── auth/
│   │   │   ├── handler.go       # HTTP handler
│   │   │   ├── service.go       # Business logic
│   │   │   ├── repository.go    # DB access
│   │   │   ├── dto.go           # Request/response
│   │   │   └── model.go         # Domain model
│   │   │
│   │   ├── survey/
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── dto.go
│   │   │   └── model.go
│   │   │
│   │   ├── food/
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── dto.go
│   │   │   └── model.go
│   │   │
│   │   ├── portion/             # 🔥 Portion size methods (as_served - untuk food & drinks)
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── dto.go
│   │   │   └── model.go
│   │   │
│   │   └── submission/
│   │       ├── handler.go
│   │       ├── service.go
│   │       ├── repository.go
│   │       ├── dto.go
│   │       └── model.go
│   │
│   ├── pkg/                     # Internal shared packages
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── logger.go
│   │   │
│   │   ├── utils/
│   │   │   ├── hash.go
│   │   │   ├── jwt.go
│   │   │   ├── response.go
│   │   │   └── validator.go
│   │   │
│   │   └── errors/              # Custom errors
│   │       └── errors.go
│   │
│   └── router/                  # Route setup
│       └── router.go
│
├── pkg/                         # 🔥 External shared packages (bisa di-import project lain)
│   └── (kosong dulu untuk MVP)
│
├── migrations/                  # DB migrations
│   ├── 001_create_users.sql
│   ├── 002_create_surveys.sql
│   └── ...
│
├── uploads/                     # Local image storage
│   └── as-served/               # Portion images (food & drinks)
│
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

### 🔥 Kenapa Structure Ini Lebih Bagus?

| Aspek               | Old (Layer-based)          | New (Feature-based)                           |
| ------------------- | -------------------------- | --------------------------------------------- |
| **Scalability**     | Berantakan kalau 50+ tabel | Tiap domain terisolasi                        |
| **Maintainability** | Harus buka 5 folder beda   | 1 folder = 1 fitur lengkap                    |
| **Testing**         | Susah mock layer lain      | Tiap domain bisa unit test sendiri            |
| **Team Work**       | Conflict di folder sama    | Beda fitur = beda folder                      |
| **Microservices**   | Susah split                | Gampang, tinggal pindahin folder `domain/xxx` |

### 📁 Contoh: Semua File Portion di 1 Folder

```
internal/domain/portion/
├── model.go           # AsServedSet, AsServedImage (drinkware: SKIP untuk MVP)
├── repository.go      # Query ke DB
├── service.go         # Business logic (kalkulasi gram, dll)
├── handler.go         # API endpoints
└── dto.go             # Request/response structs
```

**Mau edit fitur portion?** Cukup buka 1 folder ✅

---

## 🎯 MVP Checklist (Realistis)

### Phase 1: Auth & User (Week 1)

- [ ] Setup Go + Gin + GORM + MySQL
- [ ] Tabel `users` (admin + respondent)
- [ ] Tabel `refresh_tokens`
- [ ] API Register/Login JWT

### Phase 2: Survey Admin (Week 2)

- [ ] Tabel `surveys` dengan JSON meals & prompts
- [ ] CRUD Survey API
- [ ] Generate access token/link
- [ ] Tabel `survey_participants`

### Phase 3: Food DB + Portion Size (Week 2-3)

- [ ] Tabel `categories`
- [ ] Tabel `foods` (tanpa nutrisi JSON - pake tabel terpisah)
- [ ] Tabel `food_categories`
- [ ] 🔥 Tabel `nutrient_units`, `nutrient_types`, `food_nutrients`
- [ ] 🔥 Tabel `associated_foods` (sereal → susu)
- [ ] 🔥 Tabel `food_portion_size_methods` (dengan gambar!)
- [ ] 🔥 Tabel `as_served_sets` & `as_served_images` (foto porsi - food & drinks)
- [ ] ~~Tabel `drinkware_sets` & `drinkware_scales`~~ **DI-SKIP** - drinks pakai as_served
- [ ] API Search foods (FULLTEXT MySQL)
- [ ] 🔥 API Get Portion Methods per Food
- [ ] 🔥 API Get Portion Images
- [ ] Seed data makanan Indonesia (~100 item)
- [ ] 🔥 Upload gambar portion (Admin)

### Phase 4: Survey Flow dengan Portion (Week 3-4)

- [ ] API Akses survey (token/alias)
- [ ] Frontend: Step 1 - Pilih jam makan
- [ ] Frontend: Step 2 - Search makanan
- [ ] 🔥 Frontend: Step 3 - Pilih metode porsi (dengan gambar)
- [ ] 🔥 Frontend: Step 4 - Pilih jumlah/porsi (dengan gambar + fraksi)
- [ ] Frontend: Step 5 - Review & Submit
- [ ] API Submit dengan meals JSON
- [ ] Admin: Lihat submissions & Export CSV

### SKIP (NANTI AJA setelah MVP jalan):

- [ ] Session auto-save (DB) - pake localStorage dulu
- [ ] Respondent register/login
- [ ] Advanced food search (synonyms, fuzzy)
- [ ] Multi-language (i18n file JSON dulu)
- [ ] RBAC roles & permissions
- [ ] Email notifications
- [ ] Guide images (area-based portions)
- [ ] Image processing automation (resize, thumbnail)
- [ ] Multiple nutrient sources (USDA, BPOM, etc)

---

## 💡 Contoh Data Real

### Survey JSON Structure:

```json
{
  "slug": "gizi-sd-2024",
  "name": "Survey Gizi SD Kelas 5",
  "meals_config": [
    { "name": "Sarapan", "time": "06:00-08:00", "order": 1 },
    { "name": "Snack Pagi", "time": "09:00-10:00", "order": 2 },
    { "name": "Makan Siang", "time": "11:00-13:00", "order": 3 },
    { "name": "Snack Sore", "time": "15:00-16:00", "order": 4 },
    { "name": "Makan Malam", "time": "18:00-20:00", "order": 5 }
  ],
  "prompts": {
    "welcome": "Halo! Ayo ceritakan apa yang kamu makan kemarin.",
    "instructions": "Pilih waktu makan, lalu cari makanan yang kamu konsumsi."
  },
  "access_token": "gizi-sd-2024-abc123"
}
```

### Submission JSON Structure:

```json
{
  "survey_id": "uuid",
  "respondent_name": "Budi Santoso",
  "meals_data": [
    {
      "name": "Sarapan",
      "time": "07:30",
      "foods": [
        {
          "food_id": "uuid-nasi-putih",
          "food_name": "Nasi Putih",
          "portion_gram": 200,
          "nutrients": {
            "energy": 260,
            "protein": 5.2,
            "carbs": 57.2,
            "fat": 0.8
          }
        },
        {
          "food_id": "uuid-telor-goreng",
          "food_name": "Telor Goreng",
          "portion_gram": 60,
          "nutrients": { "energy": 140, "protein": 12, "fat": 10 }
        }
      ]
    },
    {
      "name": "Makan Siang",
      "time": "12:15",
      "foods": [ ... ]
    }
  ],
  "missing_foods": [
    { "name": "Kerupuk Udang", "description": "Kerupuk warna pink dari toko sebelah" }
  ]
}
```

---

## ❌ YANG DI-SKIP (NANTI AJA)

| Fitur                        | Alasan Skip                     |
| ---------------------------- | ------------------------------- |
| **Semua portion gambar**     | Complex, butuh image processing |
| **Tabel nutrisi terpisah**   | JSON cukup buat MVP             |
| **Synonyms & fuzzy search**  | MySQL FULLTEXT cukup dulu       |
| **Associated foods**         | Bisa ditambah di form manual    |
| **Session auto-save DB**     | localStorage frontend cukup     |
| **RBAC kompleks**            | 1 field `role` enum cukup       |
| **Language translations DB** | next-intl file JSON cukup       |
| **Email service**            | Skip, password reset nanti      |
| **Social login**             | Email/pass dulu                 |

---

## 🏁 FINAL: MASUK MVP V2 vs SKIP (Penjelasan Lengkap)

### ✅ MASUK MVP V2 (16 Tabel + Fitur)

| #   | Fitur/Tabel             | Penjelasan                                            | Kenapa Masuk                     |
| --- | ----------------------- | ----------------------------------------------------- | -------------------------------- |
| 1   | **Login/Auth**          | JWT + Refresh Token                                   | Core requirement                 |
| 2   | **Surveys**             | Create, manage, share link                            | Core requirement                 |
| 3   | **Food Database**       | Makanan + kategori + 🔥 nutrisi tabel                 | Core requirement                 |
| 4   | **🔥 Portion Images**   | Langsung pilih gambar porsi (tanpa pilih metode dulu) | **WAJIB** - simple & cepat       |
| 5   | **🔥 As Served Sets**   | Foto porsi ½,1,1½,2,3 + berat (food **& drinks!**)    | **WAJIB** - visual portion guide |
| 6   | **🔥 Nutrisi Tabel**    | nutrient_types, nutrient_units, food_nutrients        | **WAJIB** - perhitungan akurat   |
| 8   | **🔥 Associated Foods** | Sereal → susu, Roti → Selai                           | **WAJIB** - UX lebih baik        |
| 9   | **Input Flow**          | Add makanan → Continue → Portion (langsung)           | **WAJIB** - 4 step saja          |
| 10  | **Submission JSON**     | Simpan hasil recall                                   | Core requirement                 |

**Total: 16 Tabel**

- `users`, `refresh_tokens`
- `surveys`, `survey_participants`
- `locales`
- `categories`, `foods`, `food_categories`
- `food_portion_size_methods` 🔥
- `as_served_sets`, `as_served_images` 🔥 (juga untuk drinks!)
- 🔥 `nutrient_units`, `nutrient_types`, `food_nutrients`
- 🔥 `associated_foods`
- `survey_submissions`

> 💡 **Note:** `drinkware_sets` & `drinkware_scales` **DI-SKIP** - drinks pakai `as_served` aja!

---

### ❌ SKIP / NANTI AJA (Setelah MVP Jalan)

| Fitur                      | Penjelasan                    | Kenapa Skip                      |
| -------------------------- | ----------------------------- | -------------------------------- |
| **Session auto-save (DB)** | Simpan progress ke database   | Pakai `localStorage` dulu, cukup |
| **RBAC kompleks**          | Roles + permissions terpisah  | 1 field `role` enum cukup        |
| **Survey Schemes**         | Template reusable             | Copy-paste survey dulu           |
| **User Sessions**          | State di database             | Pakai `localStorage` dulu        |
| **Food Synonyms**          | Nama alternatif makanan       | MySQL FULLTEXT cukup             |
| **Nutrient Tables**        | Multiple sources (USDA, BPOM) | Single source cukup untuk MVP    |
| **Guide Images**           | Area-based portions           | As served images dulu            |
| **Image Processing**       | Auto-resize, thumbnail        | Upload manual dulu               |
| **Email Service**          | Reset password, notif         | Skip dulu                        |
| **Social Login**           | Google/Facebook               | Email/pass dulu                  |
| **Multi-language DB**      | `language_translations` tabel | next-intl file JSON dulu         |
| **Help & FAQ**             | Help button di UI             | Hardcode dulu                    |
| **Missing Foods DB**       | Tabel terpisah                | Simpan di submission JSON        |

---

### 📊 Perbandingan: MVP V2 vs Full Intake24

| Aspek                | MVP V2 (16 Tabel)                                        | Full Intake24 (~60+ Tabel)                 |
| -------------------- | -------------------------------------------------------- | ------------------------------------------ |
| **Tabel**            | 16 (14 → + nutrient & associated, -drinkware)            | ~60+                                       |
| **Portion Methods**  | as_served (food + drinks), weight                        | + drinkware, guide_image, standard_portion |
| **Images**           | Upload manual                                            | Auto-process, multiple sizes               |
| **Search**           | MySQL FULLTEXT                                           | Meilisearch/Elasticsearch                  |
| **Auth**             | JWT simple                                               | JWT + RBAC + 2FA                           |
| **Sessions**         | localStorage                                             | Redis + DB                                 |
| **Nutrition**        | ✅ Tabel: nutrient_types, nutrient_units, food_nutrients | Multiple sources (USDA, BPOM) + history    |
| **Associated Foods** | ✅ Ada (sereal → susu)                                   | Lebih complex dengan priority              |
| **Time Estimate**    | 1 bulan                                                  | 6+ bulan                                   |

---

### 🎯 Kesimpusan

**MVP V2 Fokus ke:**

1. ✅ Survey create & share
2. ✅ Input makanan (food + drinks dipisah)
3. ✅ **Portion selection dengan gambar** (FITUR UTAMA)
4. ✅ Submit & export data

**Flow Lengkap Respondent (4 Step):**

```
Link ──▶ Pilih Waktu ──▶ Add Foods/Drinks ──▶ Continue ──▶
Portion Selection (langsung gambar!) ──▶ Review ──▶ Submit
```

**Simplify:** Tidak ada step "Pilih Metode" (In a plate/glass/weight) — langsung tampilkan gambar porsi + input manual dalam satu layar!

---

## 🧮 RUMUS PERHITUNGAN GIZI (Nutrition Calculations)

### Formula Dasar (Intake24 Standard)

**1. Query Nutrisi dari Tabel (SQL):**

```sql
-- Ambil semua nutrisi untuk 1 makanan
SELECT
    nt.code AS nutrient_code,
    nt.name AS nutrient_name,
    nu.symbol AS unit,
    fn.value_per_100g
FROM food_nutrients fn
JOIN nutrient_types nt ON fn.nutrient_type_id = nt.id
JOIN nutrient_units nu ON nt.unit_id = nu.id
WHERE fn.food_id = 'uuid-nasi';

-- Hasil:
-- | nutrient_code | nutrient_name | unit | value_per_100g |
-- |---------------|---------------|------|----------------|
-- | energy        | Energi        | kcal | 130.0000       |
-- | protein       | Protein       | g    | 2.7000         |
-- | carbs         | Karbohidrat   | g    | 28.0000        |
```

**2. Hitung Nutrisi per Portion (Backend):**

```javascript
// Data dari query tabel food_nutrients (per 100g)
const nutrientsPer100g = [
  { code: "energy", name: "Energi", unit: "kcal", value: 130.0 },
  { code: "protein", name: "Protein", unit: "g", value: 2.7 },
  { code: "carbs", name: "Karbohidrat", unit: "g", value: 28.0 },
  { code: "fat", name: "Lemak Total", unit: "g", value: 0.3 },
  { code: "fiber", name: "Serat", unit: "g", value: 0.4 },
];

// User makan 150 gram nasi
const portionGram = 150;

// RUMUS: (nutrient_value / 100) × portionGram
function calculateNutrients(nutrientsPer100g, portionGram) {
  const result = {};

  for (const nutrient of nutrientsPer100g) {
    const calculated = (nutrient.value / 100) * portionGram;
    // Bulatkan ke 1 desimal
    result[nutrient.code] = {
      name: nutrient.name,
      unit: nutrient.unit,
      value: Math.round(calculated * 10) / 10,
    };
  }

  return result;
}

// Contoh hasil untuk 150g nasi:
// {
//   "energy": { name: "Energi", unit: "kcal", value: 195 },
//   "protein": { name: "Protein", unit: "g", value: 4.1 },
//   "carbs": { name: "Karbohidrat", unit: "g", value: 42 },
//   "fat": { name: "Lemak Total", unit: "g", value: 0.5 },
//   "fiber": { name: "Serat", unit: "g", value: 0.6 }
// }
```

**3. SQL Query Lengkap (Hitung Langsung di DB):**

```sql
-- Hitung nutrisi untuk portion 150g nasi
SELECT
    nt.code,
    nt.name,
    nu.symbol AS unit,
    ROUND((fn.value_per_100g / 100) * 150, 1) AS calculated_value
FROM food_nutrients fn
JOIN nutrient_types nt ON fn.nutrient_type_id = nt.id
JOIN nutrient_units nu ON nt.unit_id = nu.id
WHERE fn.food_id = 'uuid-nasi';
```

**2. Hitung Total per Meal:**

```javascript
// Contoh: Sarapan dengan 3 makanan
const meal = {
  name: "Sarapan",
  foods: [
    {
      name: "Nasi Putih",
      portionGram: 150,
      nutrients: { energy: 195, protein: 3.9 },
    },
    {
      name: "Telur Goreng",
      portionGram: 60,
      nutrients: { energy: 90, protein: 6 },
    },
    {
      name: "Kecap Manis",
      portionGram: 15,
      nutrients: { energy: 45, protein: 0.3 },
    },
  ],
};

// Jumlahkan semua nutrisi
const mealTotals = meal.foods.reduce((acc, food) => {
  for (const [nutrient, value] of Object.entries(food.nutrients)) {
    acc[nutrient] = (acc[nutrient] || 0) + value;
  }
  return acc;
}, {});

// Hasil: { energy: 330, protein: 10.2, ... }
```

**3. Hitung Total per Hari (Submission):**

```javascript
// Submission dengan 5 meals
const submission = {
  meals: [sarapan, snackPagi, makanSiang, snackSore, makanMalam],
};

// Jumlahkan semua meals
const dailyTotals = submission.meals.reduce((acc, meal) => {
  for (const [nutrient, value] of Object.entries(meal.totals)) {
    acc[nutrient] = (acc[nutrient] || 0) + value;
  }
  return acc;
}, {});

// Simpan di submission.meals_data[n].totalNutrients
```

### Standard Nutrients (Sesuai Intake24)

| Nutrient  | Unit | Keterangan        |
| --------- | ---- | ----------------- |
| energy    | kcal | Energi total      |
| protein   | g    | Protein           |
| carbs     | g    | Karbohidrat total |
| fat       | g    | Lemak total       |
| fiber     | g    | Serat             |
| sugar     | g    | Gula              |
| sodium    | mg   | Natrium           |
| calcium   | mg   | Kalsium           |
| iron      | mg   | Zat besi          |
| vitamin_c | mg   | Vitamin C         |

---

## 📝 CATATAN DEVELOPER (Biar Mirip Intake24)

### Struktur JSON Submission (Kaya Intake24)

```json
{
  "surveyId": "uuid",
  "participantId": "uuid",
  "submitTime": "2024-01-15T14:30:00Z",
  "meals": [
    {
      "name": "Breakfast",
      "time": "07:30",
      "foods": [
        {
          "foodId": "uuid-nasi",
          "foodCode": "NAS001",
          "foodName": "Nasi Putih",
          "portion": {
            "method": "as_served",
            "imageId": "uuid-img-3",
            "quantity": 1.5,
            "gram": 225
          },
          "nutrients": {
            "energy": 292.5,
            "protein": 5.85,
            "carbs": 63,
            "fat": 0.45
          }
        }
      ],
      "mealTotal": {
        "energy": 292.5,
        "protein": 5.85
      }
    }
  ],
  "dailyTotal": {
    "energy": 2100,
    "protein": 65
  }
}
```

### Tips UI/UX (Kaya Intake24):

1. **Food Search**: Auto-complete dengan debounce (300ms)
2. **Portion Images**: Grid 3-5 gambar, highlight saat hover
3. **Navigation**: Progress bar (Step 1/4 → 2/4 → 3/4 → 4/4)
4. **Summary**: Table dengan collapsible rows per meal
5. **Mobile First**: Touch-friendly, swipe untuk next/prev food

### Validasi Penting:

```javascript
// Validasi portionGram (max 5000g ~ 5kg per item)
if (portionGram <= 0 || portionGram > 5000) {
  throw new Error("Portion must be between 1-5000g");
}

// Validasi nutrients exists
if (!food.nutrients_per_100g || !food.nutrients_per_100g.energy) {
  throw new Error("Food must have energy data");
}
```

### Export Format (CSV):

```csv
submission_id,respondent_name,meal_name,food_name,portion_gram,energy,protein,carbs,fat
uuid-123,Budi,Sarapan,Nasi Putih,150,195,3.9,42,0.3
uuid-123,Budi,Sarapan,Telur Goreng,60,90,6,1,6
```

---

## 🚀 SIAP DEVELOP?

**Checklist Pra-Develop:**

- [ ] Setup: Node.js + MySQL + Git
- [ ] Design DB: Jalankan SQL schema
- [ ] Seed: Insert data makanan Indonesia
- [ ] API: Postman collection ready
- [ ] Frontend: Next.js + shadcn/ui installed

**Estimasi Waktu:**

- Phase 1 (Auth): 3-4 hari
- Phase 2 (Survey): 3-4 hari
- Phase 3 (Food + Portion): 5-7 hari
- Phase 4 (Flow + Submit): 5-7 hari
- Polish & Bugfix: 3-5 hari

**Total: 3-4 minggu untuk MVP** ✅

---

_Versi ini realistis untuk 1 developer dalam 1 bulan._
