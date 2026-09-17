# 17 — Redis Architecture & Implementation Best Practices PRD

> **Status:** Proposed / Planning  
> **Target:** High-Performance Caching, WebSocket Pub/Sub Scaling, Distributed Locking & Rate Limiting  
> **Tanggal:** 2026-09-16  
> **Repository:** `atlas_food_backend` & `atlas_food_frontend`

---

## 1. Vision & Executive Summary

### 🎯 Objective
Redis akan diimplementasikan sebagai **high-performance memory backbone** untuk platform **Atlas Food**. Redis tidak hanya digunakan sebagai query cache sederhana, melainkan mendukung 5 pilar utama infrastruktur:

1. **Multi-Level Query Caching** (Public Food Catalog, Search, Categories, Nutrients).
2. **WebSocket Pub/Sub Scaling Broker** (Horizontal scaling untuk multi-instance Go WebSocket Hub).
3. **Distributed Rate Limiting** (Sliding window rate limiter per IP/User).
4. **Auth & Session Security** (JWT Blacklist, Revocation, Temporary Session Store).
5. **Distributed Locks & Idempotency** (Redlock / SETNX untuk concurrency safety pada Food DB Editing & Survey Submission).

---

## 2. High-Level Architecture

```
                                 ┌───────────────────────────┐
                                 │   Frontend Clients (Web)  │
                                 └─────────────┬─────────────┘
                                               │
                                               ▼
                                 ┌───────────────────────────┐
                                 │   NGINX / Load Balancer   │
                                 └─────────────┬─────────────┘
                                               │
                       ┌───────────────────────┴───────────────────────┐
                       │                                               │
                       ▼                                               ▼
         ┌───────────────────────────┐                   ┌───────────────────────────┐
         │     Go API Node 1         │                   │     Go API Node 2         │
         │  (REST + WS Hub Instance) │                   │  (REST + WS Hub Instance) │
         └─────────────┬─────────────┘                   └─────────────┬─────────────┘
                       │                                               │
                       ├───────────────────────┬───────────────────────┤
                       │                       │                       │
                       ▼                       ▼                       ▼
         ┌───────────────────────────┐ ┌───────────────┐ ┌───────────────────────────┐
         │     Redis Cache & State   │ │ Redis Pub/Sub │ │   MySQL Primary DB        │
         │  (Key-Value, Hashes, ZSet)│ │ (WS Channels) │ │   (GORM Persistence)      │
         └───────────────────────────┘ └───────────────┘ └───────────────────────────┘
```

---

## 3. Key Namespace Standards & Schema Design

Untuk mencegah pencemaran key (key collisions) dan mempermudah flush/invalidation, seluruh key di Redis **WAJIB** mengikuti format konvensi tersandar:

```
atlas:<environment>:<module>:<entity>:<identifier_or_hash>
```

### 📋 Key Namespace Table

| Module | Purpose | Key Format Pattern | Redis Data Type | Default TTL | Invalidation Trigger |
|---|---|---|---|---|---|
| **Cache** | Public Food Categories | `atlas:prod:cache:categories:all` | `String` (JSON) | 24 Hours | Admin update/delete category |
| **Cache** | Category Foods List | `atlas:prod:cache:categories:<code_or_id>:foods` | `String` (JSON) | 6 Hours | Admin add/edit food in category |
| **Cache** | Food Detail | `atlas:prod:cache:food:detail:<food_id>` | `String` (JSON) | 12 Hours | Admin edit food/nutrients/photos |
| **Cache** | Search Query | `atlas:prod:cache:search:<query_md5_hash>` | `String` (JSON) | 1 Hour | TTL expire or bulk food change |
| **Pub/Sub**| WebSocket Room Broadcast | `atlas:prod:ws:pubsub:room:<room_id>` | `Channel` | N/A (Ephemeral)| Real-time message event |
| **WS State**| Room Active Members | `atlas:prod:ws:room:<room_id>:users` | `Hash` | 2 Hours | User join/leave / WS disconnect |
| **Rate Limit**| Endpoint Rate Limiting | `atlas:prod:ratelimit:<ip_or_user_id>:<endpoint_slug>` | `ZSet` (Sorted Set) | 1 Minute | Rolling window auto-expire |
| **Auth** | Revoked JWT Blacklist | `atlas:prod:auth:blacklist:<jti_token_id>` | `String` ("1") | Token Expiry | Expire automatically |
| **Lock** | Entity Edit Lock | `atlas:prod:lock:edit:<entity_type>:<entity_id>` | `String` (User ID) | 30 Seconds | Manual unlock or auto-expire |
| **Lock** | Submission Idempotency | `atlas:prod:lock:idem:<submission_key>` | `String` (Req Hash)| 10 Seconds | Completion of request |

---

## 4. Best Practice Code Patterns & Best Practices

### 4.1 Query Caching with Single-Flight (Anti Cache-Stampede)
Saat data cache expire dan 1.000 user bersamaan meminta data yang sama, **Cache Stampede (Thundering Herd)** terjadi jika semua query langsung menyerbu MySQL.

**Best Practice Implementation (Go):**
Menggunakan `golang.org/x/sync/singleflight` + Cache-Aside Pattern.

```go
// Code Pattern: Safe Caching with SingleFlight
type FoodCacheService struct {
    redisClient *redis.Client
    db          *gorm.DB
    sfGroup     singleflight.Group
}

func (s *FoodCacheService) GetFoodByID(ctx context.Context, id string) (*dto.FoodDetailResponse, error) {
    cacheKey := fmt.Sprintf("atlas:prod:cache:food:detail:%s", id)

    // 1. Try Read from Redis
    val, err := s.redisClient.Get(ctx, cacheKey).Result()
    if err == nil {
        var resp dto.FoodDetailResponse
        if err := json.Unmarshal([]byte(val), &resp); err == nil {
            return &resp, nil
        }
    }

    // 2. Singleflight suppress duplicate DB calls
    v, err, _ := s.sfGroup.Do(cacheKey, func() (interface{}, error) {
        // Double check cache in singleflight thread
        val, err := s.redisClient.Get(ctx, cacheKey).Result()
        if err == nil {
            var resp dto.FoodDetailResponse
            if json.Unmarshal([]byte(val), &resp) == nil {
                return &resp, nil
            }
        }

        // Fetch from MySQL
        food, err := s.fetchFoodFromDB(ctx, id)
        if err != nil {
            return nil, err
        }

        // Serialize & Store to Redis (TTL 12h + Random Jitter 0-15 min to prevent simultaneous expiry)
        jitter := time.Duration(rand.Intn(900)) * time.Second
        data, _ := json.Marshal(food)
        s.redisClient.Set(ctx, cacheKey, data, 12*time.Hour+jitter)

        return food, nil
    })

    if err != nil {
        return nil, err
    }
    return v.(*dto.FoodDetailResponse), nil
}
```

---

### 4.2 Cache Invalidation Event Hook
Jangan pernah membiarkan data di cache berberda dengan Database. Setiap kali Admin melakukan perubahan pada Food/Category via CMS/API:

```go
func (s *FoodService) UpdateFood(ctx context.Context, food *model.Food) error {
    // 1. Update MySQL DB Transaction
    if err := s.db.WithContext(ctx).Save(food).Error; err != nil {
        return err
    }

    // 2. Invalidate Related Caches (Pipeline for performance)
    pipe := s.redisClient.Pipeline()
    pipe.Del(ctx, fmt.Sprintf("atlas:prod:cache:food:detail:%s", food.ID))
    pipe.Del(ctx, fmt.Sprintf("atlas:prod:cache:categories:%s:foods", food.CategoryCode))
    pipe.Del(ctx, "atlas:prod:cache:categories:all")
    
    // Pattern search keys scan & del if needed
    _, err := pipe.Exec(ctx)
    return err
}
```

---

### 4.3 WebSocket Horizontal Scaling via Redis Pub/Sub

Setiap kali instance Go backend menerima WebSocket message yang perlu dibroadcast ke room (misal: `cursor_move` atau `food_select`), instance tersebut mempublikasikan message ke Redis Channel Room. Semua instance Go backend yang memegang client WebSocket pada room tersebut akan menerima event dari Redis dan meneruskannya ke browser client masing-masing.

```
Client A (Browser) ──> Go Node 1 ──> Redis PubSub Channel ("atlas:prod:ws:pubsub:room:XYZ")
                                                 │
                                                 ├──> Go Node 1 ──> Client B (Browser)
                                                 └──> Go Node 2 ──> Client C (Browser)
```

**Subscriber Loop Implementation:**
```go
func (r *Room) ListenRedisPubSub(ctx context.Context, redisClient *redis.Client) {
    channelName := fmt.Sprintf("atlas:prod:ws:pubsub:room:%s", r.ID)
    pubsub := redisClient.Subscribe(ctx, channelName)
    ch := pubsub.Channel()

    for msg := range ch {
        var wsMsg Message
        if err := json.Unmarshal([]byte(msg.Payload), &wsMsg); err != nil {
            continue
        }
        // Broadcast locally to websocket connections on this node
        r.broadcastLocal(wsMsg)
    }
}
```

---

### 4.4 Distributed Rate Limiting (Sliding Window counter via Redis Lua)

Untuk proteksi endpoint publik (`/api/v1/public/*`) dan Auth API dari brute-force / DDoS attack, digunakan **Sliding Window Log** berbasis Redis Sorted Set (ZSet) yang dieksekusi secara atomic dengan Script Lua:

```lua
-- Lua Script: sliding_window_rate_limiter.lua
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local clearBefore = now - window

redis.call('ZREMRANGEBYSCORE', key, 0, clearBefore)
local currentRequests = redis.call('ZCARD', key)

if currentRequests < limit then
    redis.call('ZADD', key, now, now)
    redis.call('EXPIRE', key, math.ceil(window / 1000))
    return 1 -- Allowed
else
    return 0 -- Rate Limited
end
```

---

### 4.5 Safe Distributed Lock with Lua Release

Digunakan untuk mencegah dua admin mengubah porsi/makanan yang sama secara bersamaan, atau mencegah *double-submit* survey idempotency:

```go
// Acquire Lock
func AcquireLock(ctx context.Context, rdb *redis.Client, lockKey string, lockValue string, ttl time.Duration) (bool, error) {
    return rdb.SetNX(ctx, lockKey, lockValue, ttl).Result()
}

// Release Lock safely (Lua Script to prevent deleting other user's lock if expired)
var releaseLockScript = redis.NewScript(`
    if redis.call("get", KEYS[1]) == ARGV[1] then
        return redis.call("del", KEYS[1])
    else
        return 0
    end
`)

func ReleaseLock(ctx context.Context, rdb *redis.Client, lockKey string, lockValue string) error {
    return releaseLockScript.Run(ctx, rdb, []string{lockKey}, lockValue).Err()
}
```

---

## 5. Redis Server Configuration & Resilience Best Practices

### 5.1 Eviction Policy Configuration
Di `redis.conf`:
```ini
# Maximum Memory Limit
maxmemory 512mb

# Eviction Policy: Hapus key yang paling jarang dipakai yang memiliki TTL (volatile-lru)
# atau hapus key apapun yang paling jarang dipakai (allkeys-lru).
maxmemory-policy volatile-lru
```

### 5.2 Persistence Settings
```ini
# RDB Snapshot strategy
save 900 1
save 300 10
save 60 10000

# AOF (Append Only File) strategy for zero-data-loss durability
appendonly yes
appendfsync everysec
```

### 5.3 Go Connection Pool Rules (`go-redis/v9`)
```go
rdb := redis.NewClient(&redis.Options{
    Addr:         fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
    Password:     cfg.RedisPassword,
    DB:           cfg.RedisDB,
    PoolSize:     50,                // Default max connections
    MinIdleConns: 10,                // Always keep 10 idle connections warm
    DialTimeout:  5 * time.Second,   // Connection timeout
    ReadTimeout:  3 * time.Second,   // Socket read timeout
    WriteTimeout: 3 * time.Second,   // Socket write timeout
    MaxRetries:   3,                 // Retry logic for transient network failures
})
```

---

## 6. Docker & Infrastructure Integration Plan

### File `docker-compose.yml` (Backend):
```yaml
version: '3.8'

services:
  redis:
    image: redis:7.2-alpine
    container_name: atlas_redis
    restart: always
    ports:
      - "6379:6379"
    command: redis-server --save 60 1 --loglevel notice --requirepass ${REDIS_PASSWORD:-atlasredis123} --maxmemory 512mb --maxmemory-policy volatile-lru
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "${REDIS_PASSWORD:-atlasredis123}", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5
    networks:
      - atlas-network

volumes:
  redis_data:
```

---

## 7. Implementation Roadmap & Phases

- [ ] **Phase 1: Redis Infra Setup**: Tambahkan Redis service di `docker-compose.yml` backend & tambahkan konfigurasi `.env`.
- [ ] **Phase 2: Redis Client Core Package**: Implementasi `internal/infra/redis/client.go` dengan connection pooling, healthcheck, dan fallback wrapper.
- [ ] **Phase 3: Public Query Caching**: Refactor `internal/domain/food/repository.go` dan `public_handler.go` dengan Redis Cache-Aside + Singleflight.
- [ ] **Phase 4: Distributed Rate Limiter Middleware**: Implementasi middleware sliding window Lua script pada Gin Router.
- [ ] **Phase 5: WebSocket Pub/Sub Multi-Node Broker**: Hubungkan WebSocket Hub di `internal/domain/collab` dengan Redis Pub/Sub channel.
- [ ] **Phase 6: Lock & Invalidation Verification**: Buat automated integration tests untuk pengujian invalidasi cache & failover.
