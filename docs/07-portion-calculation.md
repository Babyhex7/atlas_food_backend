# 🧮 Portion Calculation & Nutrition

## RUMUS PERHITUNGAN GIZI (Nutrition Calculations)

### Formula Dasar (Intake24 Standard)

#### 1. Query Nutrisi dari Tabel (SQL)

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

#### 2. Hitung Nutrisi per Portion (Backend)

```javascript
// Data dari query tabel food_nutrients (per 100g)
const nutrientsPer100g = [
  { code: "energy", name: "Energi", unit: "kcal", value: 130.00 },
  { code: "protein", name: "Protein", unit: "g", value: 2.70 },
  { code: "carbs", name: "Karbohidrat", unit: "g", value: 28.00 },
  { code: "fat", name: "Lemak Total", unit: "g", value: 0.30 },
  { code: "fiber", name: "Serat", unit: "g", value: 0.40 }
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
      value: Math.round(calculated * 10) / 10
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

#### 3. SQL Query Lengkap (Hitung Langsung di DB)

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

---

## Metode Portion Calculation

### 1. Image Method (As Served - untuk Food & Drinks)

```javascript
// User pilih: 1.5 porsi dari image "nugget-3" (60g per porsi)
const selectedImage = {
  weight_gram: 60.0
};
const quantity = 1.5;

const portionGram = selectedImage.weight_gram * quantity;
// 60g * 1.5 = 90g
```

---

### 2. Counter Method (+/-) - Sesuai Screenshot!

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

---

### 3. Weight Method (Manual Input)

```javascript
// User ketik langsung: 150 gram
const portionGram = userInput; // 150g
```

---

### 4. As Served with Quantity (Kaya Screenshot Pisang!)

```javascript
// User pilih: gambar ke-8 (full plate = 190g) + quantity 3.25 (3 and ¼)
const selectedImage = {
  label: "8",
  weight_gram: 190.0,  // dari as_served_images
  description: "Full plate (~190g)"
};

const quantity = 3;        // WHOLE
const fraction = 0.25;     // FRACTION (¼)
const totalQuantity = quantity + fraction;  // = 3.25

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

## Aggregation Calculations

### 2. Hitung Total per Meal

```javascript
// Contoh: Sarapan dengan 3 makanan
const meal = {
  name: "Sarapan",
  foods: [
    { name: "Nasi Putih", portionGram: 150, nutrients: {energy:195, protein:3.9} },
    { name: "Telur Goreng", portionGram: 60, nutrients: {energy:90, protein:6} },
    { name: "Kecap Manis", portionGram: 15, nutrients: {energy:45, protein:0.3} }
  ]
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

### 3. Hitung Total per Hari (Submission)

```javascript
// Submission dengan 5 meals
const submission = {
  meals: [sarapan, snackPagi, makanSiang, snackSore, makanMalam]
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

---

## Standard Nutrients (Sesuai Intake24)

| Nutrient | Unit | Keterangan |
|----------|------|------------|
| energy | kcal | Energi total |
| protein | g | Protein |
| carbs | g | Karbohidrat total |
| fat | g | Lemak total |
| fiber | g | Serat |
| sugar | g | Gula |
| sodium | mg | Natrium |
| calcium | mg | Kalsium |
| iron | mg | Zat besi |
| vitamin_c | mg | Vitamin C |

---

## Contoh Data Real

### Nutrient Units
```sql
INSERT INTO nutrient_units (code, name, symbol) VALUES
('kcal', 'Kilokalori', 'kcal'),
('g', 'Gram', 'g'),
('mg', 'Miligram', 'mg'),
('mcg', 'Mikrogram', 'μg');
```

### Nutrient Types
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

### Food Nutrients (per 100g)
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

---

## Validasi Penting

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
