# 🚀 Quick Start Guide - Find Your Food + WebSocket Collaboration

## ✅ Checklist Sebelum Mulai

- [ ] Go 1.21+ installed (`go version`)
- [ ] MySQL 8.0+ running
- [ ] Database `atlas_food` sudah dibuat
- [ ] Migrations sudah dijalankan
- [ ] File `Atlas_Makananku_FINAL.json` ada di project root
- [ ] `.env` file sudah dikonfigurasi

---

## 📦 Step 1: Seed Data Atlas Makananku

### Windows (Command Prompt):

```cmd
cd F:\Magang\BRIN\atlas_food_backend

REM Jalankan seed
go run cmd\seed\main.go
```

### Expected Output:
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
```

**Troubleshooting:**
- Jika error "JSON file not found", pastikan file ada di: `F:\Magang\BRIN\atlas_food_backend\Atlas_Makananku_FINAL.json`
- Jika error database connection, cek `.env` file

---

## 🔥 Step 2: Start API Server

```cmd
REM Start server
go run cmd\api\main.go
```

Server akan jalan di: `http://localhost:8080`

**Check if running:**
```cmd
curl http://localhost:8080/health
```

Expected: `{"status":"ok","service":"atlas-food-api"}`

---

## 🧪 Step 3: Test Public API Endpoints

### Option A: Using test script (Recommended)

```cmd
REM Run automated tests
test-find-food.bat
```

### Option B: Manual testing dengan curl

#### 1. Get All Categories (13 categories)

```cmd
curl http://localhost:8080/api/v1/public/categories
```

Expected:
```json
{
  "status": "success",
  "data": [
    {"id":"...","code":"MP","name":"Makanan Pokok","icon":"🍚"},
    {"id":"...","code":"LH","name":"Lauk Hewani","icon":"🍗"},
    ...
  ]
}
```

#### 2. Search Foods

```cmd
curl "http://localhost:8080/api/v1/public/foods/search?q=nasi&limit=5"
```

Expected:
```json
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
  "count": 5
}
```

#### 3. Get Foods by Category

```cmd
curl "http://localhost:8080/api/v1/public/categories/MP/foods?limit=10"
```

#### 4. Get Food Detail with Portion Photos

```cmd
REM Replace {food_id} with actual ID from search result
curl http://localhost:8080/api/v1/public/foods/{food_id}
```

Expected:
```json
{
  "status": "success",
  "data": {
    "id": "uuid",
    "code": "MP-01",
    "name": "Nasi",
    "local_name": "Rice",
    "description": "...",
    "photo_type": "series",
    "category": {...},
    "nutrients": {
      "energy": {"value": 130.0, "unit": "kcal"},
      "protein": {"value": 2.7, "unit": "g"}
    },
    "portion_photos": [
      {
        "id": "uuid",
        "label": "A",
        "image_url": "/uploads/atlas/mp/mp-01-a.jpg",
        "thumbnail_url": "/uploads/atlas/mp/mp-01-a-thumb.jpg",
        "weight_gram": 50.0,
        "description": "Porsi A - 50.0g"
      },
      // ... more portion photos
    ]
  }
}
```

---

## 🔌 Step 4: Test WebSocket Real-Time Collaboration

### 4.1 Install wscat (WebSocket client)

```cmd
npm install -g wscat
```

### 4.2 Create 2 Test Users

**Terminal 1:**
```cmd
REM Register User 1
curl -X POST http://localhost:8080/api/v1/auth/register ^
  -H "Content-Type: application/json" ^
  -d "{\"name\":\"Researcher One\",\"email\":\"researcher1@test.com\",\"password\":\"test123\",\"role\":\"researcher\"}"

REM Login User 1
curl -X POST http://localhost:8080/api/v1/auth/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"researcher1@test.com\",\"password\":\"test123\"}"

REM Copy the access_token from response
```

**Terminal 2:**
```cmd
REM Register User 2
curl -X POST http://localhost:8080/api/v1/auth/register ^
  -H "Content-Type: application/json" ^
  -d "{\"name\":\"Researcher Two\",\"email\":\"researcher2@test.com\",\"password\":\"test123\",\"role\":\"researcher\"}"

REM Login User 2
curl -X POST http://localhost:8080/api/v1/auth/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"researcher2@test.com\",\"password\":\"test123\"}"

REM Copy the access_token from response
```

### 4.3 Connect to WebSocket (2 Terminals)

**Terminal 1 - User 1:**
```cmd
set USER1_TOKEN=eyJhbGciOiJIUzI1NiIs...
wscat -c "ws://localhost:8080/ws/collab?token=%USER1_TOKEN%&room=find-food-test-001"
```

**Terminal 2 - User 2:**
```cmd
set USER2_TOKEN=eyJhbGciOiJIUzI1NiIs...
wscat -c "ws://localhost:8080/ws/collab?token=%USER2_TOKEN%&room=find-food-test-001"
```

### 4.4 Test Real-Time Features

#### Test 1: Presence Detection

**Window 1 (User 1) akan melihat:**
```json
< {"type":"presence_joined","user":{"user_id":"...","display_name":"Researcher Two"}}
```

**Window 2 (User 2) akan melihat:**
```json
< {"type":"presence_list","users":[{"user_id":"...","display_name":"Researcher One"}]}
```

#### Test 2: Food Search Collaboration

**Window 1 - User 1 searches:**
```json
> {"type":"food_search","payload":{"query":"nasi goreng","filters":{}}}
```

**Window 2 - User 2 sees the search:**
```json
< {"type":"food_search_shared","user_id":"...","query":"nasi goreng"}
```

#### Test 3: Food Selection

**Window 1 - User 1 selects food:**
```json
> {"type":"food_select","payload":{"food_id":"uuid-123","food_name":"Nasi Goreng"}}
```

**Window 2 - User 2 sees selection:**
```json
< {"type":"food_selected","user_id":"...","food_id":"uuid-123","food_name":"Nasi Goreng"}
```

#### Test 4: Cursor Movement

**Window 1 - User 1 moves cursor:**
```json
> {"type":"cursor_move","payload":{"x":500,"y":300,"page":"/find-food"}}
```

**Window 2 - User 2 sees cursor:**
```json
< {"type":"cursor_update","user_id":"...","x":500,"y":300,"page":"/find-food","color":"#3B82F6"}
```

#### Test 5: Ping/Pong (Heartbeat)

**Any Window:**
```json
> {"type":"ping"}
< {"type":"pong","timestamp":"2026-07-12T..."}
```

---

## 📊 Verify Data in Database

```sql
-- Connect to MySQL
mysql -u root -p atlas_food

-- Check seeded data
SELECT COUNT(*) FROM categories;  -- Expected: 13
SELECT COUNT(*) FROM foods;       -- Expected: XXX (depends on JSON)
SELECT COUNT(*) FROM as_served_images; -- Expected: XXXX

-- Sample data
SELECT code, name, photo_type FROM foods LIMIT 10;

-- Check portion photos
SELECT f.code, f.name, COUNT(asi.id) as photo_count
FROM foods f
LEFT JOIN as_served_sets ass ON ass.food_id = f.id
LEFT JOIN as_served_images asi ON asi.set_id = ass.id
GROUP BY f.id
LIMIT 10;
```

---

## 🐛 Common Issues

### Issue 1: Seed fails with "JSON file not found"

**Solution:**
```cmd
REM Check if file exists
dir Atlas_Makananku_FINAL.json

REM Set path explicitly
set ATLAS_JSON_PATH=F:\Magang\BRIN\atlas_food_backend\Atlas_Makananku_FINAL.json
go run cmd\seed\main.go
```

### Issue 2: WebSocket connection fails

**Solution:**
```cmd
REM 1. Check if token is valid
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/v1/auth/me

REM 2. Check server logs for errors
REM Look for "WebSocket" related errors

REM 3. Try with fresh login
curl -X POST http://localhost:8080/api/v1/auth/login ...
```

### Issue 3: FULLTEXT search not working

**Solution:**
```sql
-- Check if index exists
SHOW INDEX FROM foods WHERE Key_name = 'ft_name';

-- If not exists, create it:
ALTER TABLE foods ADD FULLTEXT INDEX ft_name (name, local_name);
```

### Issue 4: Empty portion_photos in response

**Solution:**
```sql
-- Check if as_served_sets exist
SELECT f.code, f.name, ass.id as set_id
FROM foods f
LEFT JOIN as_served_sets ass ON ass.food_id = f.id
WHERE f.code = 'MP-01';

-- Check if as_served_images exist
SELECT * FROM as_served_images WHERE set_id = 'YOUR_SET_ID';

-- If empty, re-run seed:
go run cmd\seed\main.go
```

---

## ✅ Success Checklist

After completing all steps, you should have:

- [ ] ✅ Database seeded with categories, foods, and portion photos
- [ ] ✅ Public API endpoints responding correctly
- [ ] ✅ Search working (returns foods when query provided)
- [ ] ✅ Food detail shows portion photos
- [ ] ✅ Categories API returns 13 categories
- [ ] ✅ 2 users can connect to WebSocket
- [ ] ✅ Users can see each other's presence
- [ ] ✅ Food search is shared in real-time
- [ ] ✅ Cursor movements synced (if tested)
- [ ] ✅ Ping/pong heartbeat working

---

## 🎯 Next Steps

### For Backend:
1. ✅ Implement optimistic locking for collaborative editing
2. ✅ Add Redis for horizontal scaling (multi-instance)
3. ✅ Add message batching for performance
4. ✅ Add Prometheus metrics for monitoring

### For Frontend:
1. Create Landing Page dengan 2 cards (Food Recall + Find Your Food)
2. Implement Find Your Food search page
3. Implement food detail page with portion photo viewer
4. Integrate WebSocket untuk real-time collaboration
5. Add live cursor overlay
6. Add activity feed component

---

## 📞 Need Help?

Check these files for detailed documentation:
- `FIND_YOUR_FOOD_README.md` - Complete feature documentation
- `docs/12-find-your-food.md` - Spec & requirements
- `docs/13-realtime-collaboration.md` - WebSocket collaboration spec

Run automated tests:
```cmd
test-find-food.bat
```

**Happy Coding! 🚀**
