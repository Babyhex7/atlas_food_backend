# 🔍 Find Your Food - Implementation Guide

## 📋 Overview

Feature **Find Your Food** adalah katalog makanan publik (tanpa login) yang menggunakan data dari **Atlas Makananku** (BRIN × UPI). Feature ini mencakup:

- ✅ Public API endpoints (no authentication required)
- ✅ Search makanan dengan FULLTEXT search
- ✅ Browse makanan per kategori (13 kategori)
- ✅ Detail makanan dengan foto porsi (series/range)
- ✅ WebSocket real-time collaboration (untuk testing)

---

## 🚀 Quick Start

### 1. Prerequisites

```bash
# Check Go version
go version  # Should be 1.21 or higher

# Check MySQL is running
mysql -u root -p -e "SELECT VERSION();"
```

### 2. Setup Database

```sql
-- Create database if not exists
CREATE DATABASE IF NOT EXISTS atlas_food;
USE atlas_food;

-- Run migrations (jika belum)
-- File migration ada di folder: migrations/
```

### 3. Seed Data from Atlas Makananku JSON

```bash
# Set environment variable (optional, default: ./Atlas_Makananku_FINAL.json)
export ATLAS_JSON_PATH="F:\Magang\BRIN\atlas_food_backend\Atlas_Makananku_FINAL.json"

# Run seed command
go run cmd/seed/main.go
```

**Expected Output:**
```
🌱 Starting Find Your Food data seeding...
============================================================
📄 Reading data from: F:\Magang\BRIN\atlas_food_backend\Atlas_Makananku_FINAL.json
📦 Processed 50/XXX foods...
📦 Processed 100/XXX foods...
✅ Successfully seeded XXX foods from Atlas Makananku
============================================================
✅ Seeding completed successfully!

📊 Summary:
   - Categories: 13
   - Foods: XXX
   - Portion Photos: XXXX

🚀 You can now start the API server:
   go run cmd/api/main.go
```

### 4. Start API Server

```bash
# Run API server
go run cmd/api/main.go
```

Server akan jalan di `http://localhost:8080`

---

## 📡 Public API Endpoints (No Auth Required)

### 1. Search Foods

```bash
# Search by query
curl "http://localhost:8080/api/v1/public/foods/search?q=nasi&limit=10"

# Response
{
  "status": "success",
  "data": [
    {
      "id": "uuid",
      "code": "MP-01",
      "name": "Nasi",
      "local_name": "Rice",
      "photo_type": "series",
      "category": {
        "id": "uuid",
        "code": "MP",
        "name": "Makanan Pokok",
        "icon": "🍚"
      }
    }
  ],
  "count": 1
}
```

### 2. Get Food Detail with Portion Photos

```bash
# Get food detail by ID
curl "http://localhost:8080/api/v1/public/foods/{food_id}"

# Response
{
  "status": "success",
  "data": {
    "id": "uuid",
    "code": "MP-01",
    "name": "Nasi",
    "local_name": "Rice",
    "description": "Nasi putih matang",
    "photo_type": "series",
    "category": {
      "id": "uuid",
      "code": "MP",
      "name": "Makanan Pokok",
      "icon": "🍚"
    },
    "nutrients": {
      "energy": { "value": 130.0, "unit": "kcal" },
      "protein": { "value": 2.7, "unit": "g" }
    },
    "portion_photos": [
      {
        "id": "uuid",
        "label": "A",
        "image_url": "/uploads/atlas/mp/mp-01-a.jpg",
        "thumbnail_url": "/uploads/atlas/mp/mp-01-a-thumb.jpg",
        "weight_gram": 50.0,
        "description": "Porsi A - 50.0g"
      }
    ]
  }
}
```

### 3. Get All Categories

```bash
curl "http://localhost:8080/api/v1/public/categories"

# Response
{
  "status": "success",
  "data": [
    {
      "id": "uuid",
      "code": "MP",
      "name": "Makanan Pokok",
      "icon": "🍚"
    },
    {
      "id": "uuid",
      "code": "LH",
      "name": "Lauk Hewani",
      "icon": "🍗"
    }
    // ... 11 more categories
  ]
}
```

### 4. Get Foods by Category

```bash
# Get all foods in Makanan Pokok category
curl "http://localhost:8080/api/v1/public/categories/MP/foods?limit=20"

# Response
{
  "status": "success",
  "data": [
    {
      "id": "uuid",
      "code": "MP-01",
      "name": "Nasi",
      "local_name": "Rice",
      "photo_type": "series",
      "category": {...}
    }
  ],
  "count": 20
}
```

---

## 🧪 Testing dengan 2 User (WebSocket Collaboration)

### Setup Test Environment

1. **Start API Server:**
```bash
go run cmd/api/main.go
```

2. **Create 2 Test Users:**

```bash
# User 1: Researcher
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Researcher One",
    "email": "researcher1@test.com",
    "password": "password123",
    "role": "researcher"
  }'

# User 2: Researcher
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Researcher Two",
    "email": "researcher2@test.com",
    "password": "password123",
    "role": "researcher"
  }'
```

3. **Login Both Users:**

```bash
# Login User 1
USER1_TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "researcher1@test.com",
    "password": "password123"
  }' | jq -r '.data.access_token')

echo "User 1 Token: $USER1_TOKEN"

# Login User 2
USER2_TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "researcher2@test.com",
    "password": "password123"
  }' | jq -r '.data.access_token')

echo "User 2 Token: $USER2_TOKEN"
```

### WebSocket Testing with 2 Windows

**Window 1 - User 1 (Researcher One):**
```bash
# Install wscat if not installed
npm install -g wscat

# Connect as User 1
wscat -c "ws://localhost:8080/ws/collab?token=$USER1_TOKEN&room=find-food-test-001"
```

**Window 2 - User 2 (Researcher Two):**
```bash
# Connect as User 2
wscat -c "ws://localhost:8080/ws/collab?token=$USER2_TOKEN&room=find-food-test-001"
```

### Test Scenarios

#### 1. Test Presence Detection

**Window 1:**
```json
// User 1 akan melihat User 2 join
< {"type":"presence_joined","user":{"user_id":"uuid","display_name":"Researcher Two"}}
```

**Window 2:**
```json
// User 2 akan melihat presence list (User 1 sudah di room)
< {"type":"presence_list","users":[{"user_id":"uuid","display_name":"Researcher One"}]}
```

#### 2. Test Food Search Collaboration

**Window 1 - User 1 searches for food:**
```json
> {"type":"food_search","payload":{"query":"nasi goreng","filters":{}}}
```

**Window 2 - User 2 sees the search in real-time:**
```json
< {"type":"food_search_shared","user_id":"user1_uuid","query":"nasi goreng","timestamp":"..."}
```

#### 3. Test Food Selection

**Window 1 - User 1 selects a food:**
```json
> {"type":"food_select","payload":{"food_id":"uuid","food_name":"Nasi Goreng"}}
```

**Window 2 - User 2 sees the selection:**
```json
< {"type":"food_selected","user_id":"user1_uuid","food_id":"uuid","food_name":"Nasi Goreng"}
```

#### 4. Test Cursor Movement (Live Collaboration)

**Window 1 - User 1 moves cursor:**
```json
> {"type":"cursor_move","payload":{"x":500,"y":300,"page":"/find-food"}}
```

**Window 2 - User 2 sees cursor update:**
```json
< {"type":"cursor_update","user_id":"user1_uuid","x":500,"y":300,"page":"/find-food","color":"#3B82F6"}
```

#### 5. Test Activity Feed

**Window 1 - User 1 adds food to meal:**
```json
> {"type":"meal_add","payload":{"meal_type":"breakfast","food_id":"uuid","food_name":"Nasi Goreng"}}
```

**Window 2 - User 2 sees activity:**
```json
< {"type":"activity_log","user_id":"user1_uuid","action":"added","details":{"meal_type":"breakfast","food_name":"Nasi Goreng"}}
```

---

## 🐛 Troubleshooting

### Issue: Seed script fails

```bash
# Check if JSON file exists
ls -la Atlas_Makananku_FINAL.json

# Check database connection
mysql -u root -p atlas_food -e "SHOW TABLES;"

# Check Go modules
go mod tidy
go mod download
```

### Issue: WebSocket connection fails

```bash
# Check if token is valid
curl -H "Authorization: Bearer $USER1_TOKEN" http://localhost:8080/api/v1/auth/me

# Check if WebSocket endpoint is accessible
curl -I http://localhost:8080/ws/collab
```

### Issue: FULLTEXT search not working

```sql
-- Check if FULLTEXT index exists
SHOW INDEX FROM foods WHERE Key_name = 'ft_name';

-- If not exists, create it:
ALTER TABLE foods ADD FULLTEXT INDEX ft_name (name, local_name);
```

---

## 📊 Database Schema Summary

### Tables Created by Seed:

1. **categories** (13 categories)
   - MP: Makanan Pokok
   - LH: Lauk Hewani
   - LN: Lauk Nabati
   - AS: Aneka Sayur
   - AB: Aneka Buah
   - ... (13 total)

2. **foods** (XXX foods from JSON)
   - Each food linked to category
   - Contains: code, name, local_name, photo_type

3. **as_served_sets** (1 per food)
   - Contains portion set metadata

4. **as_served_images** (multiple per food)
   - Contains portion photos (A, B, C, D, etc.)
   - Each with weight_gram, image_url, thumbnail_url

---

## 📝 Next Steps

### For Frontend Development:

1. **Landing Page:**
   - Add 2 cards: "Food Recall Survey" + "Find Your Food"
   - Route: `/` or `/home`

2. **Find Your Food Page:**
   - Search bar + 13 category grid
   - Route: `/find-food`

3. **Search Results:**
   - Grid of food cards
   - Route: `/find-food/search?q={query}`

4. **Category Browse:**
   - Grid of foods in selected category
   - Route: `/find-food/category/{code}`

5. **Food Detail:**
   - Main focus: Portion photos (large display)
   - Series type: Horizontal slider
   - Range type: Grid gallery
   - Route: `/find-food/{food_id}`

### For WebSocket Collaboration Testing:

1. Open 2 browser windows with different users
2. Both join same "room" (e.g., survey session)
3. Test real-time features:
   - Live cursors
   - Food search sharing
   - Portion selection sync
   - Activity feed

---

## 🎯 Success Criteria

- [ ] JSON data successfully seeded to database
- [ ] Public API endpoints working without authentication
- [ ] Search returns relevant results using FULLTEXT
- [ ] Food detail shows all portion photos
- [ ] Categories list shows all 13 categories
- [ ] WebSocket connects successfully with JWT token
- [ ] 2 users can see each other's presence in same room
- [ ] Food search is shared in real-time between users
- [ ] Cursor movements are synced (if implemented)
- [ ] Activity feed updates in real-time

---

## 📞 Support

Jika ada issue, check:
1. Logs di console (go run cmd/api/main.go)
2. Database schema (SHOW TABLES; DESCRIBE foods;)
3. Network requests di browser DevTools
4. WebSocket messages di wscat

**Happy Coding! 🚀**
