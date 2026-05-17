# 📊 Entity Relationship Diagram (ERD)

## ERD MVP (Super Simple)

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
         │ FK part_id  │     │ code        │
         │ meals(JSON) │◀────│ category_id │
         └─────────────┘     └──────┬──────┘
                                    │
                                    ▼
                            ┌─────────────┐
                            │ categories  │
                            ├─────────────┤
                            │ PK id       │
                            │ code        │
                            │ name        │
                            └─────────────┘
```

---

## Detail Relasi Antar Tabel

### Auth Domain

```
users ||--o{ refresh_tokens : has
users ||--o{ surveys : creates
users ||--o| survey_participants : participates
```

| Relasi | Tipe | Deskripsi |
|--------|------|-----------|
| users → refresh_tokens | 1:N | 1 user bisa punya multiple refresh tokens |
| users → surveys | 1:N | 1 user (admin) bisa create multiple surveys |
| users → survey_participants | 1:0/1 | 1 user bisa jadi participant (nullable untuk anonymous) |

---

### Survey Domain

```
surveys ||--o{ survey_participants : has
surveys ||--o{ survey_submissions : receives
surveys }o--|| locales : uses
```

| Relasi | Tipe | Deskripsi |
|--------|------|-----------|
| surveys → survey_participants | 1:N | 1 survey bisa punya multiple participants |
| surveys → survey_submissions | 1:N | 1 survey bisa punya multiple submissions |
| surveys → locales | N:1 | Survey menggunakan 1 locale (default: id) |

---

### Food Domain

```
categories ||--o{ foods : contains
categories ||--o{ food_categories : linked
foods ||--o{ food_categories : belongs_to
foods ||--o{ food_nutrients : has
foods ||--o{ food_portion_size_methods : has
foods ||--o{ as_served_sets : optionally_has
foods ||--o{ associated_foods : has_associations
foods ||--o{ associated_foods : is_associated_to
nutrient_units ||--o{ nutrient_types : measures
nutrient_types ||--o{ food_nutrients : measured_in
```

| Relasi | Tipe | Deskripsi |
|--------|------|-----------|
| categories → foods | 1:N | 1 kategori punya banyak makanan |
| foods ↔ categories | M:N | via food_categories (many-to-many) |
| foods → food_nutrients | 1:N | 1 makanan punya multiple nilai gizi |
| foods → food_portion_size_methods | 1:N | 1 makanan punya multiple portion methods |
| foods → as_served_sets | 1:0/1 | 1 makanan bisa link ke as_served_set (opsional) |
| foods → associated_foods | 1:N | 1 makanan bisa punya multiple associated foods |
| nutrient_units → nutrient_types | 1:N | 1 unit bisa dipakai multiple nutrient types |
| nutrient_types → food_nutrients | 1:N | 1 nutrient type bisa ada di multiple foods |

---

### Portion Domain

```
as_served_sets ||--o{ as_served_images : contains
as_served_sets }o--|| foods : optional_for
food_portion_size_methods }o--|| foods : for
```

| Relasi | Tipe | Deskripsi |
|--------|------|-----------|
| as_served_sets → as_served_images | 1:N | 1 set punya multiple images |
| as_served_sets → foods | N:0/1 | Set bisa spesifik untuk 1 food atau general |
| food_portion_size_methods → foods | N:1 | Method untuk 1 makanan |

---

### Submission Domain

```
survey_submissions }o--|| surveys : belongs_to
survey_submissions }o--o| survey_participants : from
survey_submissions ||--o| ai_result_logs : has_analysis
```

| Relasi | Tipe | Deskripsi |
|--------|------|-----------|
| survey_submissions → surveys | N:1 | Submission untuk 1 survey |
| survey_submissions → survey_participants | N:0/1 | Submission dari 1 participant (nullable) |
| survey_submissions → ai_result_logs | 1:0/1 | Submission bisa punya 1 AI analysis result (optional) |

---

## Complete ERD dengan Semua Tabel

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                    AUTH                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌───────────────┐              ┌─────────────────┐                      │
│   │     users     │──┐         ┌─│  refresh_tokens   │                      │
│   ├───────────────┤  │         │ ├─────────────────┤                      │
│   │ PK id (uuid)  │  │         │ │ PK id           │                      │
│   │ email         │  │         │ │ FK user_id      │                      │
│   │ password_hash │  │         │ │ token_hash      │                      │
│   │ name          │  │         │ │ expires_at      │                      │
│   │ role          │  │         │ └─────────────────┘                      │
│   │ is_active     │  │                                                  │
│   └───────────────┘  │                                                  │
│           │          │                                                  │
│           │ creates  │                                                  │
│           ▼          │                                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                                   SURVEY                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌───────────────┐              ┌─────────────────┐      ┌─────────────┐   │
│   │    surveys    │◀─────────────│survey_participants│◀────│    users    │   │
│   ├───────────────┤              ├─────────────────┤      └─────────────┘   │
│   │ PK id (uuid)  │    has       │ PK id (uuid)    │                      │
│   │ slug          │              │ FK survey_id    │                      │
│   │ name          │              │ FK user_id      │                      │
│   │ meals_config  │              │ alias           │                      │
│   │ prompts       │              │ is_anonymous    │                      │
│   │ access_token  │              └─────────────────┘                      │
│   │ FK locale_id  │◀─────────┐                                            │
│   └───────────────┘          │                                            │
│           │                  │                                            │
│           │ receives         │                                            │
│           ▼                  │                                            │
│   ┌───────────────┐        │                                            │
│   │  submissions  │        │                                            │
│   ├───────────────┤        │                                            │
│   │ PK id (uuid)  │        │                                            │
│   │ FK survey_id  │        │                                            │
│   │ FK part_id    │────────┘                                            │
│   │ meals_data    │                                                     │
│   └───────────────┘                                                     │
│           │                                                              │
│           │ has_analysis                                                 │
│           ▼                                                              │
│   ┌───────────────┐                                                      │
│   │ai_result_logs │                                                      │
│   ├───────────────┤                                                      │
│   │ PK id (uuid)  │                                                      │
│   │ FK submiss_id │                                                      │
│   │ input_payload │                                                      │
│   │ raw_response  │                                                      │
│   │ overall_status│                                                      │
│   │ model_used    │                                                      │
│   └───────────────┘                                                      │
│                                                                             │
│   ┌───────────────┐                                                      │
│   │    locales    │                                                      │
│   ├───────────────┤                                                      │
│   │ PK id         │                                                      │
│   │ code          │                                                      │
│   │ name          │                                                      │
│   └───────────────┘                                                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                                   FOOD                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌───────────────┐              ┌─────────────────┐                      │
│   │  categories   │◀─────────────│  food_categories  │◀────┐              │
│   ├───────────────┤              ├─────────────────┤      │              │
│   │ PK id (uuid)  │    M:N       │ PK (food_id,    │      │              │
│   │ code          │              │    category_id) │      │              │
│   │ name          │              └─────────────────┘      │              │
│   │ icon          │                                         │              │
│   └───────────────┘                                         │              │
│                                                             │ belongs_to   │
│                                                             │              │
│                              ┌─────────────────┐◀──────────┘              │
│                              │      foods      │                           │
│                              ├─────────────────┤                           │
│                              │ PK id (uuid)    │                           │
│                              │ code            │                           │
│                              │ name            │                           │
│                              │ local_name      │                           │
│                              │ category_id     │                           │
│                              └─────────────────┘                           │
│                                    │  │  │                                  │
│            ┌───────────────────────┘  │  └───────────────────────┐          │
│            │ has                    │ has                     │ has       │
│            ▼                        ▼                         ▼          │
│   ┌─────────────────┐    ┌─────────────────────┐    ┌─────────────────┐   │
│   │  food_nutrients │    │food_portion_size_meth│    │ as_served_sets  │   │
│   ├─────────────────┤    ├─────────────────────┤    ├─────────────────┤   │
│   │ PK (food_id,    │    │ PK id               │    │ PK id (uuid)    │   │
│   │     nutrient_id)│    │ FK food_id          │    │ code            │   │
│   │ value_per_100g  │    │ method_type         │    │ name            │   │
│   └─────────────────┘    │ config              │    │ category        │   │
│            ▲             └─────────────────────┘    │ FK food_id      │   │
│            │                        │                 └─────────────────┘   │
│   measured_in                        │                         │          │
│            │                         │ contains                │          │
│   ┌─────────────────┐                ▼                         │          │
│   │  nutrient_types │    ┌─────────────────────┐              │          │
│   ├─────────────────┤    │   as_served_images  │◀─────────────┘          │
│   │ PK id           │    ├─────────────────────┤                         │
│   │ code            │    │ PK id (uuid)        │                         │
│   │ name            │    │ FK set_id           │                         │
│   │ unit_id         │◀───│ label               │                         │
│   └─────────────────┘    │ image_url           │                         │
│            ▲             │ weight_gram           │                         │
│            │             └─────────────────────┘                         │
│   ┌─────────────────┐                                                      │
│   │  nutrient_units │                                                      │
│   ├─────────────────┤                                                      │
│   │ PK id           │                                                      │
│   │ code            │                                                      │
│   │ symbol          │                                                      │
│   └─────────────────┘                                                      │
│                                                                             │
│   ┌─────────────────┐                                                      │
│   │ associated_foods│                                                      │
│   ├─────────────────┤                                                      │
│   │ PK id           │                                                      │
│   │ FK food_id      │                                                      │
│   │ FK assoc_food_id│                                                      │
│   │ priority        │                                                      │
│   │ is_default      │                                                      │
│   └─────────────────┘                                                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Legend

| Symbol | Meaning |
|--------|---------|
| `PK` | Primary Key |
| `FK` | Foreign Key |
| `||--o{` | One-to-Many relationship |
| `}o--o{` | Many-to-Many relationship |
| `}o--\|\|` | Many-to-One relationship |
| `\|\|--o\|` | One-to-Zero-or-One relationship |
