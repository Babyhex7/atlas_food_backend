# 🔌 API Documentation

## Base URL

```
Development: http://localhost:8080/api/v1
Production:  https://api.atlasfood.com/api/v1
```

---

## Auth Endpoints

### Register Respondent

```http
POST /auth/register
Content-Type: application/json

Request:
{
  "email": "user@example.com",
  "password": "securepassword123",
  "name": "Budi Santoso"
}

Response 201:
{
  "status": "success",
  "data": {
    "id": "uuid-user",
    "email": "user@example.com",
    "name": "Budi Santoso",
    "role": "respondent",
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### Login

```http
POST /auth/login
Content-Type: application/json

Request:
{
  "email": "user@example.com",
  "password": "securepassword123"
}

Response 200:
{
  "status": "success",
  "data": {
    "id": "uuid-user",
    "email": "user@example.com",
    "name": "Budi Santoso",
    "role": "respondent",
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

### Refresh Token

```http
POST /auth/refresh
Content-Type: application/json

Request:
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}

Response 200:
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

### Get Profile

```http
GET /auth/me
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": {
    "id": "uuid-user",
    "email": "user@example.com",
    "name": "Budi Santoso",
    "role": "respondent",
    "is_active": true
  }
}
```

---

## Admin Endpoints

### Surveys

#### List Surveys

```http
GET /admin/surveys
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": [
    {
      "id": "uuid-survey",
      "slug": "gizi-sd-2024",
      "name": "Survey Gizi SD Kelas 5",
      "status": "active",
      "start_date": "2024-01-01",
      "end_date": "2024-12-31",
      "participant_count": 150
    }
  ]
}
```

#### Create Survey

```http
POST /admin/surveys
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "slug": "gizi-sd-2024",
  "name": "Survey Gizi SD Kelas 5",
  "description": "Survey gizi untuk siswa kelas 5",
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
  "start_date": "2024-01-01",
  "end_date": "2024-12-31"
}

Response 201:
{
  "status": "success",
  "data": {
    "id": "uuid-survey",
    "slug": "gizi-sd-2024",
    "name": "Survey Gizi SD Kelas 5",
    "access_token": "gizi-sd-2024-abc123",
    "status": "draft"
  }
}
```

#### Get Survey Detail

```http
GET /admin/surveys/:id
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": {
    "id": "uuid-survey",
    "slug": "gizi-sd-2024",
    "name": "Survey Gizi SD Kelas 5",
    "description": "Survey gizi untuk siswa kelas 5",
    "meals_config": [...],
    "prompts": {...},
    "access_token": "gizi-sd-2024-abc123",
    "status": "active",
    "participant_count": 150,
    "submission_count": 89
  }
}
```

#### Update Survey

```http
PUT /admin/surveys/:id
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "name": "Survey Gizi SD Kelas 5 Updated",
  "status": "active"
}

Response 200:
{
  "status": "success",
  "data": {
    "id": "uuid-survey",
    "name": "Survey Gizi SD Kelas 5 Updated",
    "status": "active"
  }
}
```

#### Delete Survey

```http
DELETE /admin/surveys/:id
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "message": "Survey deleted successfully"
}
```

### Submissions

#### List Submissions

```http
GET /admin/surveys/:id/submissions
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": [
    {
      "id": "uuid-submission",
      "respondent_name": "Budi Santoso",
      "submitted_at": "2024-01-15T14:30:00Z",
      "meal_count": 5,
      "total_foods": 12
    }
  ]
}
```

#### Export Submissions

```http
GET /admin/surveys/:id/export?format=csv
Authorization: Bearer {access_token}

Response: CSV file download

Format CSV:
submission_id,respondent_name,meal_name,food_name,portion_gram,energy,protein,carbs,fat
uuid-123,Budi,Sarapan,Nasi Putih,150,195,3.9,42,0.3
uuid-123,Budi,Sarapan,Telur Goreng,60,90,6,1,6
```

### Foods

#### List Foods

```http
GET /admin/foods?page=1&limit=20&category=protein
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": {
    "foods": [...],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "total_pages": 8
    }
  }
}
```

#### Create Food

```http
POST /admin/foods
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "code": "Nasi-001",
  "name": "Nasi Putih",
  "local_name": "White Rice",
  "description": "Nasi putih matang",
  "category_id": "uuid-cat-1",
  "nutrients": [
    { "type_id": 1, "value_per_100g": 130.00 },
    { "type_id": 2, "value_per_100g": 2.70 }
  ]
}

Response 201:
{
  "status": "success",
  "data": {
    "id": "uuid-food",
    "code": "Nasi-001",
    "name": "Nasi Putih"
  }
}
```

#### Update Food

```http
PUT /admin/foods/:id
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "name": "Nasi Putih Updated",
  "nutrients": [...]
}

Response 200:
{
  "status": "success",
  "data": { ... }
}
```

#### Delete Food

```http
DELETE /admin/foods/:id
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "message": "Food deleted successfully"
}
```

### Portion Size Management

#### List Portion Methods per Food

```http
GET /admin/foods/:id/portion-methods
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "method_type": "as_served",
      "label": "Chicken Nuggets",
      "image_url": "/nuggets-preview.jpg",
      "config": {
        "selectionType": "image",
        "setCode": "nuggets-portions",
        "allowFractions": true
      }
    }
  ]
}
```

#### Add Portion Method

```http
POST /admin/foods/:id/portion-methods
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "method_type": "as_served",
  "label": "Chicken Nuggets",
  "description": "Select how many nuggets you had",
  "image_url": "/nuggets-preview.jpg",
  "config": {
    "selectionType": "as_served_quantity",
    "setCode": "nuggets-portions",
    "maxQuantity": 5,
    "allowFractions": true
  }
}

Response 201:
{
  "status": "success",
  "data": { ... }
}
```

#### Update Portion Method

```http
PUT /admin/portion-methods/:id
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "label": "Updated Label",
  "config": { ... }
}

Response 200:
{
  "status": "success",
  "data": { ... }
}
```

#### Delete Portion Method

```http
DELETE /admin/portion-methods/:id
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "message": "Portion method deleted"
}
```

### As Served Sets

#### List As Served Sets

```http
GET /admin/as-served-sets
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": [
    {
      "id": "uuid-set",
      "code": "nuggets-portions",
      "name": "Chicken Nuggets Portion Guide",
      "category": "nuggets",
      "image_count": 5
    }
  ]
}
```

#### Create As Served Set

```http
POST /admin/as-served-sets
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "code": "banana-slices",
  "name": "Sliced Banana Portions",
  "description": "Visual guide for banana portions",
  "category": "fruits",
  "food_id": "uuid-banana"
}

Response 201:
{
  "status": "success",
  "data": { ... }
}
```

#### Upload Portion Images

```http
POST /admin/as-served-sets/:id/images
Authorization: Bearer {access_token}
Content-Type: multipart/form-data

Request:
- images[]: [File]
- labels[]: ["1", "2", "3"]
- weights[]: [20.0, 40.0, 60.0]
- descriptions[]: ["1 nugget (~20g)", "2 nuggets (~40g)", "3 nuggets (~60g)"]

Response 201:
{
  "status": "success",
  "data": {
    "uploaded": 3,
    "images": [...]
  }
}
```

---

## Public/Respondent Endpoints

### Survey Access

#### Get Active Surveys

```http
GET /survey/active
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": {
    "surveys": [
      {
        "id": "uuid-survey",
        "slug": "gizi-sd-2024",
        "name": "Survey Gizi SD Kelas 5",
        "status": "active",
        "start_date": "2024-01-01",
        "end_date": "2024-12-31",
        "participant_count": 150,
        "created_at": "2024-01-01"
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10
  }
}
```

#### Access Survey with Token

```http
POST /survey/access
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "token": "gizi-sd-2024-abc123",
  "alias": "PART-A7X9K2",
  "respondent_name": "Budi Santoso"
}

Response 200:
{
  "status": "success",
  "data": {
    "survey": {
      "id": "uuid-survey",
      "name": "Survey Gizi SD Kelas 5",
      "meals_config": [...],
      "prompts": {...}
    },
    "participant": {
      "id": "uuid-participant",
      "alias": "PART-A7X9K2"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs..."  // survey session token
  }
}
```

#### Get Survey Info

```http
GET /survey/:id/info
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": {
    "id": "uuid-survey",
    "name": "Survey Gizi SD Kelas 5",
    "meals_config": [...],
    "prompts": {...}
  }
}
```

### Foods

#### Search Foods

```http
GET /foods/search?q=nasi&category=staples&limit=20
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": [
    {
      "id": "uuid-food",
      "code": "Nasi-001",
      "name": "Nasi Putih",
      "local_name": "White Rice",
      "category": "Makanan Pokok",
      "icon": "🍚"
    },
    {
      "id": "uuid-food-2",
      "name": "Nasi Goreng",
      "local_name": "Fried Rice",
      "category": "Makanan Pokok",
      "icon": "🍚"
    }
  ]
}
```

#### Get Food Detail

```http
GET /foods/:id
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": {
    "id": "uuid-food",
    "code": "Nasi-001",
    "name": "Nasi Putih",
    "local_name": "White Rice",
    "description": "Nasi putih matang",
    "category": "Makanan Pokok",
    "nutrients": {
      "energy": { "value": 130.00, "unit": "kcal" },
      "protein": { "value": 2.70, "unit": "g" },
      "carbs": { "value": 28.00, "unit": "g" }
    },
    "portion_methods": [...]
  }
}
```

#### Get Portion Methods per Food

```http
GET /foods/:id/portion-methods
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "method_type": "as_served",
      "label": "In a plate",
      "description": "Choose portion size",
      "image_url": "/rice/preview.jpg",
      "config": {
        "selectionType": "as_served_quantity",
        "setCode": "rice-plates",
        "maxQuantity": 3,
        "allowFractions": true,
        "fractionOptions": [0, 0.5]
      }
    }
  ]
}
```

### Portion Selection

#### Get Portion Options (As Served Images)

```http
GET /portion-methods/:id/options
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": {
    "method": {
      "id": 1,
      "label": "Sliced banana on plate",
      "selection_type": "as_served_quantity"
    },
    "images": [
      {
        "id": "uuid-img-1",
        "label": "1",
        "image_url": "/banana/banana-1.jpg",
        "thumbnail_url": "/banana/banana-1-thumb.jpg",
        "weight_gram": 20.0,
        "description": "Few slices (~20g)"
      },
      {
        "id": "uuid-img-8",
        "label": "8",
        "image_url": "/banana/banana-8.jpg",
        "thumbnail_url": "/banana/banana-8-thumb.jpg",
        "weight_gram": 190.0,
        "description": "Full plate (~190g)"
      }
    ]
  }
}
```

#### Get As Served Set Images

```http
GET /as-served-sets/:code/images
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": {
    "set": {
      "code": "banana-slices",
      "name": "Sliced Banana Portions",
      "category": "fruits"
    },
    "images": [...]
  }
}
```

### Categories

#### List Categories

```http
GET /categories
Authorization: Bearer {access_token}

Response 200:
{
  "status": "success",
  "data": [
    { "id": "uuid-cat-1", "code": "staples", "name": "Makanan Pokok", "icon": "🍚" },
    { "id": "uuid-cat-2", "code": "protein", "name": "Lauk Pauk", "icon": "🍗" },
    { "id": "uuid-cat-3", "code": "fruits", "name": "Buah-buahan", "icon": "🍌" },
    { "id": "uuid-cat-4", "code": "drinks", "name": "Minuman", "icon": "🥤" }
  ]
}
```

### Survey Submission

#### Submit Survey

```http
POST /survey/submit
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "survey_id": "uuid-survey",
  "respondent_name": "Budi Santoso",
  "respondent_email": "budi@example.com",
  "meals_data": [
    {
      "name": "Sarapan",
      "time": "07:30",
      "foods": [
        {
          "food_id": "uuid-nasi-putih",
          "food_name": "Nasi Putih",
          "portion_gram": 150,
          "portion": {
            "method": "as_served_quantity",
            "image_id": "uuid-img-3",
            "image_label": "3",
            "base_weight": 60,
            "quantity": 2,
            "fraction": 0.5,
            "total_quantity": 2.5
          },
          "nutrients": {
            "energy": 195,
            "protein": 3.9,
            "carbs": 42,
            "fat": 0.3
          }
        }
      ],
      "meal_total": {
        "energy": 195,
        "protein": 3.9
      }
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

Response 201:
{
  "status": "success",
  "data": {
    "submission_id": "uuid-submission",
    "message": "Survey submitted successfully"
  }
}
```

---

## AI Nutrition Analysis Endpoints

### Get Nutrition Analysis (On-Demand)

> **PENTING:** AI dipanggil **on-demand** hanya ketika user klik tombol "AI Recommendation" di halaman hasil survey. Tidak dipanggil saat submission.

#### Request AI Analysis

```http
POST /ai/nutrition-analysis
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "submission_id": "uuid-submission"
}

Response 200 (Fresh from Groq):
{
  "status": "success",
  "source": "groq",
  "data": {
    "overall_status": "less",
    "overall_message": "Your current nutrition is still below the recommended daily requirement. Additional balanced nutrients are needed.",

    "nutritional_analysis": [
      {
        "label": "Calories",
        "status": "low",
        "description": "Current calorie level is still relatively low for optimal daily energy needs."
      },
      {
        "label": "Protein",
        "status": "low",
        "description": "Protein source is limited and should be increased to support body recovery and muscle maintenance."
      },
      {
        "label": "Balance",
        "status": "partial",
        "description": "Your meal already contains sufficient carbohydrates, but fiber and micronutrient sources are still lacking."
      }
    ],

    "ai_recommendation": "To improve your nutritional balance, consider adding:\n- Grilled chicken or fish for additional protein\n- Vegetables such as broccoli or spinach for fiber and vitamins\n- Fruits like banana or apple for natural nutrients\n- More water intake to maintain hydration balance",

    "recommended_foods": [
      "Grilled Chicken", "Boiled Egg", "Broccoli",
      "Spinach", "Banana", "Apple", "Greek Yogurt", "Mineral Water"
    ],

    "health_insight": {
      "title": "Mild Nutritional Deficiency",
      "description": "Your current meal composition is considered partially balanced, but additional protein, vegetables, and hydration are recommended to better fulfill daily nutritional needs."
    },

    "suggested_activities": ["Light Walking", "Yoga", "Stretching"]
  }
}

Response 200 (From Cache):
{
  "status": "success",
  "source": "cache",
  "data": { ... }  // Same structure as above
}

Response 404 (Submission Not Found):
{
  "status": "error",
  "message": "Submission not found or access denied"
}

Response 503 (Groq Service Error):
{
  "status": "error",
  "message": "AI service temporarily unavailable, please try again"
}
```

**Notes:**

- `source: "groq"` → Fresh analysis from Groq API
- `source: "cache"` → Retrieved from database (already analyzed before)
- Cache by `submission_id`: Same submission always returns same result
- Ownership check: User can only analyze their own submissions
- Groq API key never exposed to frontend

---

## AI Result (Terintegrasi di Submission) - DEPRECATED

> ⚠️ **DEPRECATED:** Versi lama di mana AI result dikirim otomatis saat submission. Sekarang AI dipanggil on-demand via endpoint terpisah di atas.

### Submit Survey (dengan AI Result) - OLD VERSION

```http
POST /survey/submit
Authorization: Bearer {access_token}
Content-Type: application/json

Request: (sama seperti sebelumnya)
{ ... }

Response 201:
{
  "status": "success",
  "data": {
    "submission_id": "uuid-submission",
    "message": "Survey submitted successfully",

    "ai_result": {
      "overall_status": "less",
      "overall_message": "Your current nutrition is still below the recommended daily requirement. Additional balanced nutrients are needed.",

      "nutritional_analysis": [
        {
          "label": "Calories",
          "status": "low",
          "description": "Current calorie level is still relatively low for optimal daily energy needs."
        },
        {
          "label": "Protein",
          "status": "low",
          "description": "Protein source is limited and should be increased to support body recovery."
        },
        {
          "label": "Balance",
          "status": "partial",
          "description": "Your meal already contains sufficient carbohydrates, but fiber sources are still lacking."
        }
      ],

      "ai_recommendation": "To improve your nutritional balance, consider adding:\n- Grilled chicken or fish for additional protein\n- Vegetables such as broccoli or spinach for fiber and vitamins\n- Fruits like banana or apple for natural nutrients",

      "recommended_foods": [
        "Grilled Chicken", "Boiled Egg", "Broccoli",
        "Spinach", "Banana", "Apple", "Greek Yogurt", "Mineral Water"
      ],

      "health_insight": {
        "title": "Mild Nutritional Deficiency",
        "description": "Your current meal composition is considered partially balanced, but additional protein, vegetables, and hydration are recommended."
      },

      "suggested_activities": ["Light Walking", "Yoga", "Stretching"]
    }
  }
}
```

> **Catatan:** Jika Groq API gagal (timeout/error), `ai_result` akan bernilai `null`. Frontend harus handle kondisi ini dengan graceful fallback.

---

## Error Responses

### Standard Error Format

```json
{
  "status": "error",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      { "field": "email", "message": "Email is required" },
      {
        "field": "password",
        "message": "Password must be at least 8 characters"
      }
    ]
  }
}
```

### Error Codes

| Code               | HTTP Status | Description              |
| ------------------ | ----------- | ------------------------ |
| `UNAUTHORIZED`     | 401         | Invalid or missing token |
| `FORBIDDEN`        | 403         | Insufficient permissions |
| `NOT_FOUND`        | 404         | Resource not found       |
| `VALIDATION_ERROR` | 422         | Invalid input data       |
| `CONFLICT`         | 409         | Resource already exists  |
| `INTERNAL_ERROR`   | 500         | Server error             |
