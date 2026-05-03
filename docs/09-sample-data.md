# 📝 Sample Data

## Contoh Data Portion Methods

### 1. Nuggets dengan foto porsi
```sql
INSERT INTO food_portion_size_methods (food_id, method_type, label, description, image_url, config) VALUES
('uuid-nuggets', 'as_served', 'Chicken Nuggets', 'Select how many nuggets you had', '/nuggets-preview.jpg', '
{
  "selectionType": "image",
  "setCode": "nuggets-portions",
  "allowFractions": true,
  "fractionOptions": [0.25, 0.5, 0.75, 1, 1.5, 2, 2.5, 3]
}
');
```

### 2. Minuman dengan pilihan botol/gelas/carton
```sql
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
```

### 3. Counter method (+/-) - untuk item per pack/unit
```sql
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
```

### 4. As Served with Quantity Multiplier (Kaya screenshot pisang!)
```sql
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
```

### 5. Manual weight (fallback)
```sql
INSERT INTO food_portion_size_methods (food_id, method_type, label, description, config) VALUES
('uuid-nuggets', 'weight', 'Enter weight manually', 'If you know the exact weight', '
{
  "selectionType": "input",
  "unit": "gram",
  "placeholder": "Enter grams..."
}
');
```

---

## Contoh Data Images

### As Served Set untuk Nuggets
```sql
INSERT INTO as_served_sets (id, code, name, description, category) VALUES
('uuid-set-1', 'nuggets-portions', 'Chicken Nuggets Portion Guide', 'Visual guide for nugget portions', 'nuggets');

-- Gambar per porsi
INSERT INTO as_served_images (set_id, label, image_url, thumbnail_url, weight_gram, description) VALUES
('uuid-set-1', '1', '/nuggets/nugget-1.jpg', '/nuggets/nugget-1-thumb.jpg', 20.0, '1 nugget (~20g)'),
('uuid-set-1', '2', '/nuggets/nugget-2.jpg', '/nuggets/nugget-2-thumb.jpg', 40.0, '2 nuggets (~40g)'),
('uuid-set-1', '3', '/nuggets/nugget-3.jpg', '/nuggets/nugget-3-thumb.jpg', 60.0, '3 nuggets (~60g)'),
('uuid-set-1', '4', '/nuggets/nugget-4.jpg', '/nuggets/nugget-4-thumb.jpg', 80.0, '4 nuggets (~80g)'),
('uuid-set-1', '5', '/nuggets/nugget-5.jpg', '/nuggets/nugget-5-thumb.jpg', 100.0, '5 nuggets (~100g)');
```

### As Served Set untuk Pisang (Kaya Screenshot!)
```sql
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
('uuid-set-3', '8', '/banana/banana-8.jpg', '/banana/banana-8-thumb.jpg', 190.0, 'Full plate (~190g)', 8);
```

---

## Contoh Data Nutrisi

### 1. Nutrient Units
```sql
INSERT INTO nutrient_units (code, name, symbol) VALUES
('kcal', 'Kilokalori', 'kcal'),
('g', 'Gram', 'g'),
('mg', 'Miligram', 'mg'),
('mcg', 'Mikrogram', 'μg');
```

### 2. Nutrient Types
```sql
INSERT INTO nutrient_types (code, name, unit_id, display_order) VALUES
('energy', 'Energi', 1, 1),           -- kcal
('protein', 'Protein', 2, 2),         -- g
('carbs', 'Karbohidrat', 2, 3),       -- g
('fat', 'Lemak Total', 2, 4),         -- g
('fiber', 'Serat', 2, 5),             -- g
('sugar', 'Gula', 2, 6),              -- g
('sodium', 'Natrium', 3, 7),          -- mg
('calcium', 'Kalsium', 3, 8),         -- mg
('iron', 'Zat Besi', 3, 9),           -- mg
('vit_c', 'Vitamin C', 3, 10);        -- mg
```

### 3. Categories
```sql
INSERT INTO categories (id, code, name, icon) VALUES
('uuid-cat-1', 'staples', 'Makanan Pokok', '🍚'),
('uuid-cat-2', 'protein', 'Lauk Pauk', '🍗'),
('uuid-cat-3', 'fruits', 'Buah-buahan', '🍌'),
('uuid-cat-4', 'drinks', 'Minuman', '🥤');
```

### 4. Foods
```sql
INSERT INTO foods (id, code, name, local_name, description, category_id) VALUES
('uuid-nasi', 'Nasi-001', 'Nasi Putih', 'White Rice', 'Nasi putih matang', 'uuid-cat-1'),
('uuid-ayam', 'Ayam-001', 'Ayam Goreng', 'Fried Chicken', 'Ayam goreng tepung', 'uuid-cat-2'),
('uuid-pisang', 'Pisang-001', 'Pisang', 'Banana', 'Pisang cavendish', 'uuid-cat-3'),
('uuid-susu', 'Susu-001', 'Susu UHT', 'UHT Milk', 'Susu UHT full cream', 'uuid-cat-4');
```

### 5. Food Nutrients (per 100g)
```sql
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
```

### 6. Food Categories (many-to-many)
```sql
INSERT INTO food_categories (food_id, category_id) VALUES
('uuid-nasi', 'uuid-cat-1'),
('uuid-ayam', 'uuid-cat-2'),
('uuid-pisang', 'uuid-cat-3'),
('uuid-susu', 'uuid-cat-4');
```

### 7. Associated Foods (Sereal → Susu!)
```sql
INSERT INTO associated_foods (food_id, associated_food_id, priority, is_default) VALUES
-- Makan sereal, biasanya minum susu
('uuid-sereal', 'uuid-susu', 1, TRUE),
-- Makan roti, biasanya pakai selai atau susu
('uuid-roti', 'uuid-susu', 1, FALSE),
('uuid-roti', 'uuid-selai', 2, FALSE),
-- Nasi goreng → telur
('uuid-nasi-goreng', 'uuid-telur', 1, FALSE);
```

---

## Contoh Survey JSON Structure

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

---

## Contoh Submission JSON Structure

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
          "portion_gram": 150,
          "nutrients": {
            "energy": 195,
            "protein": 3.9,
            "carbs": 42,
            "fat": 0.3
          }
        },
        {
          "food_id": "uuid-telor-goreng",
          "food_name": "Telor Goreng",
          "portion_gram": 60,
          "nutrients": { "energy": 90, "protein": 6, "fat": 6 }
        }
      ],
      "meal_total": {
        "energy": 285,
        "protein": 9.9
      }
    },
    {
      "name": "Makan Siang",
      "time": "12:15",
      "foods": [ 
        {
          "food_id": "uuid-pisang",
          "food_name": "Pisang",
          "portion": {
            "method": "as_served_quantity",
            "image_id": "uuid-img-8",
            "image_label": "8",
            "base_weight": 190,
            "quantity": 3,
            "fraction": 0.25,
            "total_quantity": 3.25,
            "portion_gram": 617.5
          },
          "nutrients": {
            "energy": 617.5,
            "carbs": 158.3,
            "protein": 7.5
          }
        }
      ]
    }
  ],
  "daily_total": {
    "energy": 2100,
    "protein": 65
  },
  "missing_foods": [
    { "name": "Kerupuk Udang", "description": "Kerupuk warna pink dari toko sebelah" }
  ]
}
```

---

## Export Format (CSV)

```csv
submission_id,respondent_name,meal_name,food_name,portion_gram,energy,protein,carbs,fat
uuid-123,Budi,Sarapan,Nasi Putih,150,195,3.9,42,0.3
uuid-123,Budi,Sarapan,Telur Goreng,60,90,6,1,6
uuid-123,Budi,Makan Siang,Pisang,617.5,617.5,7.5,158.3,0
```
