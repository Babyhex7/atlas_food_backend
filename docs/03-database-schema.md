# 🗄️ Database Schema

## Ringkasan Tabel MVP (16 Tabel)

### 1. Auth (Minimal)

#### users
```sql
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
```

#### refresh_tokens
```sql
CREATE TABLE refresh_tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

---

### 2. Survey (Minimal)

#### locales
```sql
CREATE TABLE locales (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,  -- 'id', 'en'
    name VARCHAR(50) NOT NULL
);
```

#### surveys
```sql
CREATE TABLE surveys (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    slug VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,

    -- Meals sebagai JSON (simplified, ga perlu tabel terpisah)
    meals_config JSON NOT NULL,  -- [{"name":"Breakfast","time":"07:00"}, ...]

    -- Prompts sebagai JSON
    prompts JSON,  -- {"welcome":"...","search":"..."}

    locale_id INT DEFAULT 1,
    start_date DATE,
    end_date DATE,
    status ENUM('draft', 'active', 'closed') DEFAULT 'draft',

    -- Simple alias/token untuk anonymous access
    access_token VARCHAR(255),  -- Token untuk akses survey

    created_by CHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id)
);
```

#### survey_participants
```sql
CREATE TABLE survey_participants (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    survey_id CHAR(36) NOT NULL,
    user_id CHAR(36),  -- NULL kalau anonymous
    alias VARCHAR(50), -- NULL kalau respondent terdaftar
    is_anonymous BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

---

### 3. Food Database (Minimal)

#### categories
```sql
CREATE TABLE categories (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    icon VARCHAR(50),  -- emoji atau icon name
    display_order INT DEFAULT 0
);
```

#### foods
```sql
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
```

#### food_categories (many-to-many)
```sql
CREATE TABLE food_categories (
    food_id CHAR(36) NOT NULL,
    category_id CHAR(36) NOT NULL,
    PRIMARY KEY (food_id, category_id),
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);
```

#### nutrient_units
```sql
CREATE TABLE nutrient_units (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,      -- 'g', 'mg', 'kcal', 'IU'
    name VARCHAR(50) NOT NULL,             -- 'gram', 'miligram'
    symbol VARCHAR(10) NOT NULL            -- 'g', 'mg'
);
```

#### nutrient_types
```sql
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
```

#### food_nutrients
```sql
CREATE TABLE food_nutrients (
    food_id CHAR(36) NOT NULL,
    nutrient_type_id INT NOT NULL,
    value_per_100g DECIMAL(10, 4) NOT NULL,  -- nilai per 100g
    PRIMARY KEY (food_id, nutrient_type_id),
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    FOREIGN KEY (nutrient_type_id) REFERENCES nutrient_types(id) ON DELETE CASCADE
);
```

#### associated_foods
```sql
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

---

### 4. Portion Size (dengan Gambar!)

#### food_portion_size_methods
```sql
CREATE TABLE food_portion_size_methods (
    id INT AUTO_INCREMENT PRIMARY KEY,
    food_id CHAR(36) NOT NULL,

    -- Tipe metode
    method_type ENUM('as_served', 'guide_image', 'weight') NOT NULL,  -- drinkware: SKIP MVP

    -- Label yang ditampilkan (contoh: "In a glass", "In a bottle", "1 nugget")
    label VARCHAR(255) NOT NULL,
    description VARCHAR(255),

    -- Gambar untuk metode ini
    image_url VARCHAR(500),
    thumbnail_url VARCHAR(500),

    -- Konfigurasi pilihan (JSON)
    config JSON,

    -- Urutan tampilan
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE CASCADE,
    INDEX idx_food (food_id)
);
```

#### as_served_sets
```sql
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
```

#### as_served_images
```sql
CREATE TABLE as_served_images (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    set_id CHAR(36) NOT NULL,

    -- Urutan/label (1, 2, 3 atau "small", "medium", "large")
    label VARCHAR(50) NOT NULL,            -- "1", "2", "3", "Small", "Medium"

    -- Gambar
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
```

Note: `drinkware_sets` & `drinkware_scales` **DI-SKIP** untuk MVP — drinks pakai `as_served_sets` aja!

---

### 5. Submission (Minimal)

#### survey_submissions
```sql
CREATE TABLE survey_submissions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    survey_id CHAR(36) NOT NULL,
    participant_id CHAR(36),  -- dari survey_participants

    -- Data respondent (snapshot, ga perlu join)
    respondent_name VARCHAR(255),
    respondent_email VARCHAR(255),

    -- Semua meals & foods simpan di JSON
    meals_data JSON NOT NULL,

    -- Missing foods (kalau user ketik makanan ga ketemu di DB)
    missing_foods JSON,

    submitted_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (survey_id) REFERENCES surveys(id) ON DELETE CASCADE,
    FOREIGN KEY (participant_id) REFERENCES survey_participants(id) ON DELETE SET NULL
);
```

---

## Penjelasan Konsep

### Q: Apa itu `user_survey_aliases`?

**A:** Kode/token unik buat respondent masuk ke survey tanpa register/login.

| Contoh                          | Penjelasan                             |
| ------------------------------- | -------------------------------------- |
| `PART-A7X9K2`                   | Alias code yang di-share ke respondent |
| URL: `/survey/abc123?token=xyz` | Link unik dengan token                 |

**Gunanya:** Respondent tinggal klik link → langsung isi survey, tanpa perlu daftar akun.

---

### Q: Bedanya `nutrient_unit` vs `nutrient_types`?

| Tabel            | Fungsi                | Contoh                            |
| ---------------- | --------------------- | --------------------------------- |
| `nutrient_types` | **APA** yang diukur   | Protein, Lemak, Energi, Vitamin C |
| `nutrient_units` | **SATUAN** pengukuran | gram (g), miligram (mg), kcal     |

**Relasi:**
```
nutrient_types: "Protein" ──unit_id──▶ nutrient_units: "gram (g)"
```

---

### Q: Bedanya `food_category` vs `category`?

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

**MVP V2:** **PAKAI** tabel `food_nutrients` + `nutrient_types` + `nutrient_units` (perhitungan akurat & scalable).

#### 2. `food_nutrients`

Simpan nilai gizi per makanan per 100g.
**MVP V2:** **PAKAI** - simpan di `food_nutrients` dengan struktur: `(food_id, nutrient_type_id, value_per_100g)`.

#### 3. `associated_foods` (sereal → susu)

Fitur "makanan yang sering dipakai bersama" (contoh: sereal → susu).
**MVP V2:** **PAKAI** - tabel `associated_foods` untuk UX lebih baik.

#### 4. `food_attributes` vs `food_attribute_values`

Custom properties kayak "brand", "is_organic", "halal", etc.
**MVP:** NANTI AJA, pakai JSON di tabel `foods`.
