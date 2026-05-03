# 🔄 Workflow & Alur Aplikasi

## User Flow Overview

### Admin Flow

```
Login ──▶ Create Survey ──▶ Upload Food DB ──▶ Set Portion Images ──▶ Share Link/Token
```

### Respondent Flow (4 Step Saja - Final Version!)

```
Link ──▶ Step 1 ──▶ Step 2 ──▶ Step 3 ──▶ Step 4 ──▶ Done
         Pilih      Add        Portion   Review
         Waktu      Food       (Langsung  & Submit
                    & Drink    Gambar!)
```

---

## Detail Alur Respondent

### Step 1: PILIH/EDIT WAKTU MAKAN

```
┌────────────────────────────────────────────────────────────┐
│  Sarapan (07:00) ▶   [Edit Meal]        [? HELP & FAQ]       │
└────────────────────────────────────────────────────────────┘
```

**Actions:**
- Pilih waktu makan dari daftar (Sarapan, Snack, Makan Siang, dll)
- Edit waktu makan jika perlu
- Lanjut ke step berikutnya

---

### Step 2: ADD MAKANAN & MINUMAN (List dulu, portion nanti!)

```
┌────────────────────────────────────────────────────────────┐
│  What did you have for Sarapan?                            │
│                                                           │
│  🍔 Foods                        🥤 Drinks                │
│  ┌──────────────┐ ┌────┐        ┌──────────────┐ ┌────┐   │
│  │ [Nasi Goreng]│ │ ADD│        │ [Es Teh     ]│ │ ADD│   │
│  └──────────────┘ └────┘        └──────────────┘ └────┘   │
│                                                           │
│  ✅ Nasi Goreng (akan di-portion-kan nanti)              │
│  ✅ Telur Goreng (akan di-portion-kan nanti)             │
│  ✅ Es Teh (akan di-portion-kan nanti)                   │
│                                                           │
│  [🕐 Change Meal Time] [🗑️ Delete Meal] [▶ Continue]      │
└────────────────────────────────────────────────────────────┘
```

**Actions:**
- Search makanan/minuman
- Tambahkan ke list
- Setelah selesai, klik **Continue** untuk ke portion selection

---

### Step 3: PORTION SELECTION (2 Tipe UI - Sesuai Jenis Makanan!)

#### 🔥 TIPE A: Simple Grid (Nasi, Nuggets, dll)

```
┌────────────────────────────────────────────────────────────────────────┐
│  "How much did you have?" (1 of 3 foods)                               │
│                                                                       │
│  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐  ┌─────────────────┐         │
│  │ 🖼️  │ │ 🖼️  │ │ 🖼️  │ │ 🖼️  │ │ 🖼️  │  │  Atau ketik:   │         │
│  │ ½   │ │ 1   │ │ 1½  │ │ 2   │ │ 3   │  │  [___150___] g │         │
│  │porsi│ │porsi│ │porsi│ │porsi│ │porsi│  │                 │         │
│  │~75g │ │~150g│ │~225g│ │~300g│ │~450g│  │ [  ✓ CONFIRM  ] │         │
│  └─────┘ └─────┘ └─────┘ └─────┘ └─────┘  └─────────────────┘         │
│                                                                       │
│  [◀ Previous Food]  [Next Food ▶]  [Skip This Food]                   │
└────────────────────────────────────────────────────────────────────────┘
```

#### 🔥 TIPE B: As Served + Quantity (Pisang, Kentang, dll)

```
┌──────────────────────────────────────────────────────────────────────┐
│                                                                    │
│     [🍌 Gambar besar pisang di piring (selected)]                  │
│                                                                    │
│        ┌──────┐           ┌──────┐                                │
│        │  +   │     and   │  +   │   ← overlay di gambar         │
│        │  3   │           │  ¼   │                                │
│        │  -   │           │  -   │                                │
│        └──────┘           └──────┘                                │
│       WHOLE              FRACTION                                  │
│                                                                    │
│  "I had 3 and ¼ of the largest portion (617.50g)"                  │
│                                                                    │
│  Thumbnail row: [🖼️] [🖼️] [🖼️] [🖼️] [🖼️] [🖼️] [🖼️] [🖼️👆selected] │
│                20g   40g   60g   95g  130g  160g  175g  190g          │
│                                                                    │
│  [    ✓ I HAD THAT MUCH    ]                                        │
│  [◀ Previous Food]  [Next Food ▶]  [Skip This Food]               │
└──────────────────────────────────────────────────────────────────────┘
```

**Actions:**
- Pilih gambar porsi yang sesuai
- Atur quantity (jika diperlukan)
- Atur fraksi (jika diperlukan)
- Konfirmasi porsi
- Lanjut ke makanan berikutnya atau ke review

---

### Step 4: REVIEW & SUBMIT

```
┌────────────────────────────────────────────────────────────┐
│  📋 Summary for Sarapan (08:00)                             │
│  ├─ 🍔 Nasi Goreng - 1.5 porsi (~225g)                      │
│  ├─ 🍔 Telur Goreng - 1 porsi (~60g)                        │
│  └─ 🥤 Es Teh - 1 gelas besar (~250ml)                      │
│                                                           │
│  [✏️ Edit Portions]  [➕ Add More Food]  [✅ SUBMIT MEAL]  │
└────────────────────────────────────────────────────────────┘
```

**Actions:**
- Review semua makanan yang diinput
- Edit porsi jika perlu
- Tambah makanan lain jika perlu
- Submit meal atau lanjut ke waktu makan berikutnya

---

## UI Flow Lengkap dengan Portion Selection (Baru!)

### RESPONDENT JOURNEY dengan Gambar Porsi:

```
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
│        (as_served)      (as_served)      (weight)          │
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

## Simplify Final Flow

**Flow Lengkap Respondent (4 Step):**

```
Link ──▶ Pilih Waktu ──▶ Add Foods/Drinks ──▶ Continue ──▶
Portion Selection (langsung gambar!) ──▶ Review ──▶ Submit
```

**Simplify:** Tidak ada step "Pilih Metode" (In a plate/glass/weight) — langsung tampilkan gambar porsi + input manual dalam satu layar!

---

## MVP Checklist (Realistis)

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
- [ ] Tabel `nutrient_units`, `nutrient_types`, `food_nutrients`
- [ ] Tabel `associated_foods` (sereal → susu)
- [ ] Tabel `food_portion_size_methods` (dengan gambar!)
- [ ] Tabel `as_served_sets` & `as_served_images` (foto porsi - food & drinks)
- [ ] ~~Tabel `drinkware_sets` & `drinkware_scales`~~ **DI-SKIP** - drinks pakai as_served
- [ ] API Search foods (FULLTEXT MySQL)
- [ ] API Get Portion Methods per Food
- [ ] API Get Portion Images
- [ ] Seed data makanan Indonesia (~100 item)
- [ ] Upload gambar portion (Admin)

### Phase 4: Survey Flow dengan Portion (Week 3-4)

- [ ] API Akses survey (token/alias)
- [ ] Frontend: Step 1 - Pilih jam makan
- [ ] Frontend: Step 2 - Search makanan
- [ ] Frontend: Step 3 - Pilih metode porsi (dengan gambar)
- [ ] Frontend: Step 4 - Pilih jumlah/porsi (dengan gambar + fraksi)
- [ ] Frontend: Step 5 - Review & Submit
- [ ] API Submit dengan meals JSON
- [ ] Admin: Lihat submissions & Export CSV

### SKIP (NANTI AJA setelah MVP jalan):

- [ ] Session auto-save (DB) - pake localStorage dulu
- [ ] Respondent register/login (anonymous dulu)
- [ ] Advanced food search (synonyms, fuzzy)
- [ ] Multi-language (i18n file JSON dulu)
- [ ] RBAC roles & permissions
- [ ] Email notifications
- [ ] Guide images (area-based portions)
- [ ] Image processing automation (resize, thumbnail)
- [ ] Multiple nutrient sources (USDA, BPOM, etc)
