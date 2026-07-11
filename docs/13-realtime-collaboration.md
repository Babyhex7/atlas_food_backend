# 13 — Real-Time Collaboration Feature Plan

> **Status:** Planning  
> **Target:** Multi-user real-time collaboration untuk dietary recall survey  
> **Tanggal:** 2026-07-05

---

## 1. Vision & Scope

### 🎯 Apa yang mau dicapai?

Kolaborasi real-time di Atlas Food memungkinkan **multiple users** (researcher + respondent, atau sesama admin) bekerja bersama dalam satu sesi survey / food database management:

| Fitur | Deskripsi |
|---|---|
| **User Presence** | Lihat siapa aja yang online di sesi yang sama (avatar, nama, role) |
| **Live Cursors** | Lihat posisi cursor & aktivitas user lain secara real-time (warna per user) |
| **Collaborative Food Search** | Satu user search, user lain bisa lihat hasil search & apa yang diklik |
| **Shared Survey Session** | Researcher & respondent isi survey barengan, researcher bisa bantu pilih makanan / portion |
| **Real-Time Food DB Editing** | Multiple admin edit food database barengan tanpa konflik |
| **Activity Feed** | Log aktivitas: "Bagas added 'Nasi Goreng' to Lunch", "Rina changed portion to 250g" |
| **Conflict Resolution** | Deteksi & resolve edit conflict (optimistic locking / CRDT-based) |

### 🧑‍🤝‍🧑 User Stories

```
RESEARCHER: "Saya mau bantu responden isi survey 24h recall, 
             lihat apa yang mereka search, dan bantu pilih makanan 
             yang tepat — semuanya real-time."

ADMIN TEAM: "Saya dan tim admin mau edit food database barengan 
             tanpa saling timpa perubahan satu sama lain."

RESPONDENT: "Saya mau researcher bisa bantu saya isi survey 
            tanpa harus screenshare atau tatap muka langsung."
```

---

## 2. High-Level Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                        FRONTEND (Next.js)                         │
│                                                                    │
│  ┌──────────┐  ┌─────────────┐  ┌──────────────┐  ┌───────────┐ │
│  │ Presence │  │ Live Cursor │  │ Collaborative│  │ Activity  │ │
│  │ Indicator│  │   Overlay   │  │   Search     │  │   Feed    │ │
│  └────┬─────┘  └──────┬──────┘  └──────┬───────┘  └─────┬─────┘ │
│       │               │               │                 │        │
│  ┌────┴───────────────┴───────────────┴─────────────────┴──────┐ │
│  │              useWebSocket Hook + Zustand Store              │ │
│  │    (auto-reconnect, heartbeat, message queue, presence)     │ │
│  └─────────────────────────────┬───────────────────────────────┘ │
└────────────────────────────────┼──────────────────────────────────┘
                                 │ WSS:// (TLS)
                                 │
┌────────────────────────────────┼──────────────────────────────────┐
│                        BACKEND (Go/Gin)                            │
│                                                                    │
│  ┌─────────────────────────────┴───────────────────────────────┐ │
│  │                  WebSocket Hub (gorilla/websocket)           │ │
│  │    ┌──────────┐  ┌───────────┐  ┌───────────┐  ┌────────┐  │ │
│  │    │  Room    │  │Broadcast  │  │ Heartbeat │  │  Auth  │  │ │
│  │    │ Manager  │  │  Engine   │  │  Manager  │  │  Guard │  │ │
│  │    └──────────┘  └───────────┘  └───────────┘  └────────┘  │ │
│  └─────────────────────────────┬───────────────────────────────┘ │
│                                                                    │
│  ┌─────────────────────────────┴───────────────────────────────┐ │
│  │              Message Router & Handler                         │ │
│  │    cursor_move | food_search | food_select | presence_join    │ │
│  │    meal_add   | portion_set | review_submit | activity_log    │ │
│  └─────────────────────────────┬───────────────────────────────┘ │
│                                                                    │
│  ┌──────┴──────┐  ┌──────────┐  ┌──────────┐                     │
│  │   Redis     │  │   MySQL  │  │   GORM   │                     │
│  │ (Pub/Sub,   │  │ (Session │  │ (Persist │                     │
│  │  Session)   │  │  History)│  │  Data)   │                     │
│  └─────────────┘  └──────────┘  └──────────┘                     │
└──────────────────────────────────────────────────────────────────┘
```

### Kenapa Redis?

- **Pub/Sub** — Kalau nanti scale horizontal (multiple Go instances), message dari satu instance bisa di-broadcast ke instance lain via Redis pub/sub
- **Session Store** — Simpan state room, daftar user per room, cursor position (ephemeral, TTL-based)
- **Rate Limiting** — Batasi message per detik per user

### Kenapa gorilla/websocket?

- Library WebSocket paling mature untuk Go
- Support auto-ping/pong, read/write deadline, concurrent safe
- Community besar, dokumentasi lengkap

---

## 3. WebSocket Protocol Design

### 3.1 Message Format

Semua message dalam format JSON:

```json
{
  "type": "string",
  "room_id": "uuid",
  "payload": {},
  "sender_id": "uuid",
  "timestamp": "iso8601",
  "sequence": 12345
}
```

### 3.2 Message Types

#### CLIENT → SERVER

| Type | Payload | Deskripsi |
|---|---|---|
| `presence_join` | `{ user_id, display_name, role, avatar_url }` | User masuk room |
| `presence_leave` | `{}` | User keluar room |
| `cursor_move` | `{ x, y, page, element_id }` | Posisi cursor user |
| `food_search` | `{ query, filters }` | User search makanan |
| `food_select` | `{ food_id, food_name }` | User klik/select makanan |
| `meal_add` | `{ meal_type, food_id, food_name }` | User tambah makanan ke meal |
| `portion_set` | `{ meal_id, food_id, portion_gram }` | User set porsi |
| `review_submit` | `{ survey_id }` | User submit review |
| `db_edit_start` | `{ entity_type, entity_id }` | User mulai edit food DB |
| `db_edit_field` | `{ entity_type, entity_id, field, value, version }` | Perubahan field |
| `db_edit_save` | `{ entity_type, entity_id, changes, version }` | Simpan perubahan |
| `db_edit_cancel` | `{ entity_type, entity_id }` | Batal edit |
| `ping` | `{}` | Heartbeat client |

#### SERVER → CLIENT

| Type | Payload | Deskripsi |
|---|---|---|
| `presence_list` | `{ users: [...] }` | Daftar user di room |
| `presence_joined` | `{ user }` | User baru join |
| `presence_left` | `{ user_id }` | User leave |
| `cursor_update` | `{ user_id, x, y, page, color }` | Cursor user lain pindah |
| `food_search_shared` | `{ user_id, query, results }` | Broadcast hasil search |
| `food_selected` | `{ user_id, food_id, food_name }` | Broadcast food select |
| `meal_updated` | `{ user_id, meal }` | Broadcast meal update |
| `portion_updated` | `{ user_id, meal_id, food_id, portion }` | Broadcast portion update |
| `review_submitted` | `{ user_id, survey_id }` | Notifikasi review submit |
| `db_locked` | `{ entity_type, entity_id, locked_by }` | Ada yang lagi edit |
| `db_field_updated` | `{ entity_type, entity_id, field, value }` | Perubahan field real-time |
| `db_edit_saved` | `{ entity_type, entity_id, version }` | Edit selesai disimpan |
| `activity_log` | `{ user_id, action, details }` | Entry activity feed |
| `error` | `{ code, message }` | Error dari server |
| `pong` | `{}` | Heartbeat response |

---

## 4. Backend Implementation Plan

### 4.1 Struktur File (Domain `collab`)

```
atlas_food_backend/internal/domain/collab/
├── handler.go          // WebSocket upgrade handler + Gin route
├── hub.go              // Central Hub: room management, broadcast
├── room.go             // Room struct: client list, message routing
├── client.go           // Client struct: read/write pump, connection
├── message.go          // Message types & parsing
├── auth.go             // WebSocket auth middleware (JWT validation)
├── presence.go         // Presence tracking per room
├── dto.go              // Request/Response DTOs
└── collab_test.go      // Unit & integration tests
```

### 4.2 Central Hub Pattern

```go
// hub.go — Central Hub (singleton)
type Hub struct {
    rooms      map[string]*Room
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

// room.go — Room management
type Room struct {
    ID          string
    clients     map[*Client]bool
    broadcast   chan Message
    presence    *PresenceManager
    history     []Message               // Ring buffer
    mu          sync.RWMutex
}

// client.go — Per-connection handler
type Client struct {
    hub      *Hub
    conn     *websocket.Conn
    room     *Room
    userID   uuid.UUID
    userName string
    userRole string
    send     chan Message
}
```

### 4.3 Connection Lifecycle

```
1. Client connects to WSS://host/ws/collab?token={jwt}&room={uuid}
                    ↓
2.  AuthGuard validates JWT token (from query param)
    - Invalid → Close connection (4001)
    - Valid → Extract userID, role, name
                    ↓
3.  RoomManager.GetOrCreate(roomID)
    - Check if room exists
    - Create if not exists
                    ↓
4.  Client joins Room
    - Add to room.clients map
    - Store presence: ROOM:{roomID}:USERS
    - Broadcast presence_joined to all room members
    - Send presence_list to new client
                    ↓
5.  Start readPump + writePump goroutines
    - readPump:  Read messages → validate → route → handle
    - writePump: Send channel → write to WebSocket
                    ↓
6.  Heartbeat: Server pings every 15s, client must pong within 10s
    - No pong → Connection timeout → cleanup
                    ↓
7.  Disconnect / Leave
    - Remove from room
    - Remove from presence
    - Broadcast presence_left
    - Cleanup goroutines
```

### 4.4 Room Types & Authorization

| Room Type | Format | Who Can Join | Persistent? |
|---|---|---|---|
| **Survey Session** | `survey:{surveyID}:{alias}` | Researcher + Respondent assigned to survey | No (ephemeral, auto-cleanup) |
| **Admin Food DB** | `admin:food-db` | Admin role only | Semi-persistent |
| **Admin Survey Edit** | `admin:survey:{surveyID}` | Admin role only | No |

### 4.5 Message Routing

```go
func (r *Room) routeMessage(msg Message, sender *Client) {
    switch msg.Type {
    case "cursor_move":
        r.broadcastExcept(msg, sender)

    case "food_search":
        r.broadcastExcept(msg, sender)

    case "meal_add", "portion_set", "food_select":
        r.updateMealState(msg)
        r.broadcastExcept(msg, sender)
        r.logActivity(msg, sender)

    case "db_edit_start":
        lock := r.acquireLock(msg.EntityType, msg.EntityID, sender)
        if lock == nil {
            sender.sendError("ENTITY_LOCKED", "Another user is editing this")
            return
        }
        r.broadcast(Message{Type: "db_locked", ...})

    case "db_edit_save":
        r.broadcast(msg)

    case "ping":
        sender.send <- Message{Type: "pong"}
    }
}
```

### 4.6 Database Changes

```sql
-- Tabel untuk persist collaborative session history
CREATE TABLE collaboration_sessions (
    id          CHAR(36) PRIMARY KEY,
    room_id     VARCHAR(255) NOT NULL,
    user_id     CHAR(36) NOT NULL,
    joined_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at     DATETIME,
    INDEX idx_room (room_id)
);

-- Tabel activity log
CREATE TABLE activity_logs (
    id          CHAR(36) PRIMARY KEY,
    room_id     VARCHAR(255) NOT NULL,
    user_id     CHAR(36) NOT NULL,
    action      VARCHAR(100) NOT NULL,
    details     JSON,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_room_time (room_id, created_at DESC)
);

-- Optimistic locking column
ALTER TABLE foods ADD COLUMN version INT NOT NULL DEFAULT 1;
ALTER TABLE surveys ADD COLUMN version INT NOT NULL DEFAULT 1;
```

### 4.7 Router Integration

```go
func SetupRoutes(r *gin.Engine, h *collab.Handler) {
    // ... existing routes ...

    // WebSocket endpoint
    r.GET("/ws/collab", h.UpgradeConnection)

    // REST endpoints untuk collaboration
    collabGroup := r.Group("/api/v1/collab")
    collabGroup.Use(middleware.JWTAuth())
    {
        collabGroup.GET("/rooms/:roomID/history", h.GetRoomHistory)
        collabGroup.GET("/rooms/:roomID/presence", h.GetRoomPresence)
        collabGroup.POST("/rooms/:roomID/invite", h.InviteToRoom)
    }
}
```

---

## 5. Frontend Implementation Plan

### 5.1 Struktur File

```
atlas_food_frontend/src/internal/collab/
├── hooks/
│   ├── useWebSocket.ts
│   ├── usePresence.ts
│   ├── useLiveCursor.ts
│   ├── useCollaborativeSearch.ts
│   ├── useCollaborativeMeal.ts
│   └── useActivityFeed.ts
├── store/
│   └── collabStore.ts
├── components/
│   ├── PresenceAvatars.tsx
│   ├── LiveCursorOverlay.tsx
│   ├── ActivityFeed.tsx
│   ├── CollaborationBar.tsx
│   └── LockIndicator.tsx
├── types/
│   └── collab.ts
└── lib/
    ├── connection.ts
    └── cursor.ts
```

### 5.2 Zustand Collab Store

```typescript
interface CollabUser {
  userId: string;
  displayName: string;
  role: 'admin' | 'researcher' | 'respondent';
  avatarUrl?: string;
  color: string;
  cursor?: { x: number; y: number; page: string };
  lastActive: number;
}

interface CollabState {
  isConnected: boolean;
  roomId: string | null;
  users: CollabUser[];
  activities: ActivityEntry[];
  lockedEntities: Map<string, { lockedBy: string; lockedAt: number }>;

  connect: (roomId: string) => void;
  disconnect: () => void;
  setUsers: (users: CollabUser[]) => void;
  addUser: (user: CollabUser) => void;
  removeUser: (userId: string) => void;
  updateCursor: (userId: string, cursor: { x: number; y: number; page: string }) => void;
  addActivity: (activity: ActivityEntry) => void;
  setLock: (entityType: string, entityId: string, lockedBy: string) => void;
  releaseLock: (entityType: string, entityId: string) => void;
}
```

### 5.3 useWebSocket Hook

```typescript
function useWebSocket(roomId: string | null) {
  const store = useCollabStore();
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout>();
  const reconnectAttemptRef = useRef(0);

  useEffect(() => {
    if (!roomId) return;

    const connect = () => {
      const token = getCookie('atlas_token');
      const wsUrl = `${process.env.NEXT_PUBLIC_WS_URL}/ws/collab?token=${token}&room=${roomId}`;
      const ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        reconnectAttemptRef.current = 0;
        store.setConnected(true);
        startHeartbeat(ws);
      };

      ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        routeIncomingMessage(msg, store);
      };

      ws.onclose = (event) => {
        store.setConnected(false);
        if (event.code !== 1000) {
          const delay = Math.min(1000 * (2 ** reconnectAttemptRef.current), 30000);
          reconnectTimeoutRef.current = setTimeout(connect, delay);
          reconnectAttemptRef.current++;
        }
      };

      wsRef.current = ws;
    };

    connect();

    return () => {
      clearTimeout(reconnectTimeoutRef.current);
      wsRef.current?.close(1000, 'Component unmount');
    };
  }, [roomId]);

  const send = useCallback((type: string, payload: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type, payload }));
    }
  }, []);

  return { send, isConnected: store.isConnected };
}
```

### 5.4 Live Cursor System

```typescript
function useLiveCursor(send: MessageSender) {
  const users = useCollabStore(s => s.users);

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      send('cursor_move', {
        x: e.clientX,
        y: e.clientY,
        page: window.location.pathname
      });
    };

    // Throttle cursor updates: max 15 fps
    const throttled = throttle(handler, 66);
    window.addEventListener('mousemove', throttled);
    return () => window.removeEventListener('mousemove', throttled);
  }, [send]);

  const remoteCursors = users
    .filter(u => u.cursor && u.cursor.page === window.location.pathname)
    .map(u => ({
      userId: u.userId,
      name: u.displayName,
      x: u.cursor!.x,
      y: u.cursor!.y,
      color: u.color,
    }));

  return { remoteCursors };
}
```

```tsx
function LiveCursorOverlay({ cursors }) {
  return (
    <div className="pointer-events-none fixed inset-0 z-50">
      {cursors.map(c => (
        <div
          key={c.userId}
          className="absolute transition-all duration-75 ease-linear"
          style={{ left: c.x, top: c.y }}
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill={c.color}>
            <path d="M3 1l15 12-6 2-3 5-3-1 3-6-6-2z" />
          </svg>
          <span
            className="ml-4 rounded px-2 py-0.5 text-xs text-white whitespace-nowrap"
            style={{ backgroundColor: c.color }}
          >
            {c.name}
          </span>
        </div>
      ))}
    </div>
  );
}
```

### 5.5 Collaborative Search Integration

```tsx
function CollaborativeSearch() {
  const { send } = useWebSocket(roomId);
  const { searchQuery, searchResults, searchingUserId } = useCollabStore();

  const handleSearch = async (query: string) => {
    // Search via REST API (React Query)
    const results = await searchFoods(query);

    // Broadcast ke room
    send('food_search', { query, filters: {} });
  };

  return (
    <div>
      <SearchInput onSearch={handleSearch} />
      {searchingUserId && (
        <div className="text-sm text-muted">
          {searchingUserId} is searching: "{searchQuery}"
        </div>
      )}
    </div>
  );
}
```

---

## 6. Conflict Resolution Strategy

### 6.1 Food DB Editing — Optimistic Locking

```
Scenario: Dua admin edit food yang sama bersamaan

Admin A                          Server                         Admin B
   │                               │                               │
   │── db_edit_start(food:123) ──→│                               │
   │←── db_locked(by: A) ─────────│── db_locked(by: A) ────────→ │
   │                               │                               │
   │── db_edit_field(name:"Nasi")─→│── db_field_updated ────────→│
   │                               │                               │
   │                               │      Admin B coba edit juga   │
   │                               │←── db_edit_start(food:123) ──│
   │                               │── error(LOCKED, by: A) ────→│
   │                               │      ❌ Ditolak!              │
   │                               │                               │
   │── db_edit_save(food:123) ───→│                               │
   │                               │── save to MySQL + version++   │
   │←── db_edit_saved ────────────│── db_edit_saved ────────────→│
   │                               │── db_unlocked ──────────────→│
```

### 6.2 Survey Meal Editing — Last-Write-Wins + Activity Log

Karena meal editing di survey lebih kolaboratif (researcher bantu respondent), konflik jarang terjadi:

1. **Last-write-wins** untuk field yang sama
2. **Merge** untuk meal list (additive, jarang ada delete)
3. **Activity log** mencatat semua perubahan → bisa di-revert manual

### 6.3 CRDT (Untuk Masa Depan)

Untuk collaborative text editing kompleks (food descriptions, survey prompts):

- **Yjs** (CRDT library) + **y-websocket** di frontend
- **y-redis** atau custom Yjs backend di Go
- Ini overkill untuk fase awal — optimistic locking sudah cukup

---

## 7. Error Handling & Edge Cases

### 7.1 Koneksi

| Kondisi | Handling |
|---|---|
| **Disconnect (network issue)** | Exponential backoff reconnect (1s, 2s, 4s, 8s... max 30s). Tampilkan banner "Reconnecting..." |
| **JWT expired saat reconnect** | Trigger token refresh via REST, lalu reconnect |
| **403 Forbidden (no access to room)** | Close connection, redirect ke home, tampilkan toast error |
| **Server crash / restart** | Semua client disconnect → auto-reconnect. Server rebuild room state dari DB + Redis |
| **Rate limit exceeded** | Server kirim `error` type: `RATE_LIMITED`. Client pause sending 5 detik |

### 7.2 State Inconsistency

| Kondisi | Handling |
|---|---|
| **Client state stale** | Server kirim `presence_list` + `state_sync` penuh saat reconnect |
| **Message out-of-order** | Pakai `sequence` number. Client ignore message dengan sequence < lastSeen |
| **Duplicate message** | Client dedup berdasarkan `message_id` (UUID per message dari client) |
| **Conflict detected** | Server return error `VERSION_CONFLICT` → client fetch latest data via REST, lalu retry |

### 7.3 Edge Cases

| Kondisi | Handling |
|---|---|
| **User close tab tanpa leave** | Server heartbeat timeout (15s no ping) → auto cleanup setelah 30s |
| **Multiple tabs dari user yang sama** | Allow, beri suffix `(Tab 2)` di display name. Cursor hanya dari tab terakhir yang aktif |
| **Room untuk survey yang sudah selesai** | Room closed (read-only). User baru tidak bisa join |
| **User dengan koneksi lambat** | Batasi max message size (64KB per message). Drop messages > 5 detik terlambat |
| **100+ user di satu room** | Batch cursor updates (kirim setiap 200ms, bukan setiap move) |

---

## 8. Performance Optimizations

### 8.1 Message Throttling & Batching

#### Client-Side Throttling

```typescript
// Cursor moves - max 15 fps (66ms interval)
const throttledCursorMove = throttle((x: number, y: number) => {
  ws.send(JSON.stringify({
    type: 'cursor_move',
    payload: { x, y, page: window.location.pathname }
  }));
}, 66);

// Food search - debounce 300ms
const debouncedSearch = debounce((query: string) => {
  ws.send(JSON.stringify({
    type: 'food_search',
    payload: { query, filters: {} }
  }));
}, 300);

// Portion updates - debounce 500ms (heavy operation)
const debouncedPortion = debounce((mealId: string, portion: number) => {
  ws.send(JSON.stringify({
    type: 'portion_set',
    payload: { meal_id: mealId, portion_gram: portion }
  }));
}, 500);
```

#### Server-Side Batching

```go
// hub.go - Batch broadcast untuk cursor updates
type Hub struct {
    // ... existing fields ...
    cursorBatch      map[string]*CursorUpdate
    cursorBatchMutex sync.Mutex
    batchTicker      *time.Ticker
}

func (h *Hub) startBatchProcessor() {
    h.batchTicker = time.NewTicker(50 * time.Millisecond)
    go func() {
        for range h.batchTicker.C {
            h.flushCursorBatch()
        }
    }()
}

func (h *Hub) flushCursorBatch() {
    h.cursorBatchMutex.Lock()
    defer h.cursorBatchMutex.Unlock()

    if len(h.cursorBatch) == 0 {
        return
    }

    // Aggregate all cursor updates
    updates := make([]CursorUpdate, 0, len(h.cursorBatch))
    for _, update := range h.cursorBatch {
        updates = append(updates, *update)
    }

    // Broadcast as single message
    h.broadcast <- Message{
        Type:    "cursor_batch",
        Payload: updates,
    }

    // Clear batch
    h.cursorBatch = make(map[string]*CursorUpdate)
}
```

### 8.2 Connection & Memory Optimization

#### WebSocket Configuration

```go
// config/websocket.go
type WebSocketConfig struct {
    // Buffer sizes - tuned untuk message typical size (1-4KB)
    ReadBufferSize  int // 4KB - cukup untuk cursor + search query
    WriteBufferSize int // 8KB - cukup untuk broadcast banyak cursor

    // Compression - enable untuk message > 1KB
    EnableCompression bool

    // Timeouts
    WriteTimeout      time.Duration // 10s
    ReadTimeout       time.Duration // 60s (allow idle time)
    PongTimeout       time.Duration // 10s
    PingInterval      time.Duration // 15s

    // Rate limiting
    MaxMessageSize    int64 // 64KB per message
    MaxMessagesPerSec int   // 50 msg/sec per connection
}

var DefaultWSConfig = WebSocketConfig{
    ReadBufferSize:    4096,
    WriteBufferSize:   8192,
    EnableCompression: true,
    WriteTimeout:      10 * time.Second,
    ReadTimeout:       60 * time.Second,
    PongTimeout:       10 * time.Second,
    PingInterval:      15 * time.Second,
    MaxMessageSize:    65536, // 64KB
    MaxMessagesPerSec: 50,
}

// handler.go - Apply config
upgrader := websocket.Upgrader{
    ReadBufferSize:    cfg.ReadBufferSize,
    WriteBufferSize:   cfg.WriteBufferSize,
    EnableCompression: cfg.EnableCompression,
    WriteBufferPool:   &sync.Pool{}, // Reuse buffers
    CheckOrigin: func(r *http.Request) bool {
        // Whitelist allowed origins
        origin := r.Header.Get("Origin")
        return isAllowedOrigin(origin)
    },
}
```

#### Redis Connection Pooling

```go
// config/redis.go
import "github.com/go-redis/redis/v9"

func NewRedisClient() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_URL"),
        Password: os.Getenv("REDIS_PASSWORD"),

        // Connection pool
        PoolSize:     50,  // Max connections
        MinIdleConns: 10,  // Keep warm connections
        MaxIdleConns: 20,
        PoolTimeout:  4 * time.Second,

        // Timeouts
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,

        // Retry strategy
        MaxRetries:      3,
        MinRetryBackoff: 8 * time.Millisecond,
        MaxRetryBackoff: 512 * time.Millisecond,
    })
}
```

#### Memory Management

```go
// room.go - Limit room history to prevent memory leak
type Room struct {
    // ... existing fields ...
    history     *RingBuffer // Max 100 messages
    maxClients  int         // Max 100 concurrent users per room
}

// Simple ring buffer implementation
type RingBuffer struct {
    data  []Message
    size  int
    index int
    mu    sync.RWMutex
}

func NewRingBuffer(size int) *RingBuffer {
    return &RingBuffer{
        data: make([]Message, size),
        size: size,
    }
}

func (rb *RingBuffer) Add(msg Message) {
    rb.mu.Lock()
    defer rb.mu.Unlock()
    rb.data[rb.index] = msg
    rb.index = (rb.index + 1) % rb.size
}

func (rb *RingBuffer) GetLast(n int) []Message {
    rb.mu.RLock()
    defer rb.mu.RUnlock()
    
    if n > rb.size {
        n = rb.size
    }
    
    result := make([]Message, 0, n)
    for i := 0; i < n; i++ {
        idx := (rb.index - i - 1 + rb.size) % rb.size
        if !rb.data[idx].Timestamp.IsZero() {
            result = append(result, rb.data[idx])
        }
    }
    return result
}
```

### 8.3 Smart Message Filtering

```go
// room.go - Filter messages based on relevance
func (r *Room) shouldBroadcast(msg Message, recipient *Client) bool {
    switch msg.Type {
    case MessageCursorMove:
        // Only send cursor updates if recipient is on same page
        payload := msg.Payload.(map[string]interface{})
        recipientPage := recipient.getCurrentPage()
        return payload["page"] == recipientPage

    case MessageFoodSearch:
        // Always broadcast search to all room members
        return true

    case MessageDBEdit:
        // Only broadcast to admins
        return recipient.role == "admin"

    default:
        return true
    }
}

// Optimize broadcast loop
func (r *Room) broadcastFiltered(msg Message, sender *Client) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    // Pre-serialize message once
    msgBytes, err := json.Marshal(msg)
    if err != nil {
        return
    }

    for client := range r.clients {
        if client == sender {
            continue
        }

        // Apply filtering
        if !r.shouldBroadcast(msg, client) {
            continue
        }

        // Send pre-serialized bytes (no re-marshaling)
        select {
        case client.send <- msgBytes:
        default:
            // Client send buffer full - skip (don't block)
            r.metrics.DroppedMessages.Inc()
        }
    }
}
```

### 8.4 Binary Protocol Optimization (Phase 2)

#### MessagePack Implementation

```go
// For high-traffic messages (cursor, presence)
import "github.com/vmihailenco/msgpack/v5"

type BinaryMessage struct {
    Type      uint8  `msgpack:"t"`  // 1 byte type code
    UserID    uint32 `msgpack:"u"`  // 4 bytes (map UUID to uint32)
    X         uint16 `msgpack:"x"`  // 2 bytes (0-65535)
    Y         uint16 `msgpack:"y"`  // 2 bytes
    Timestamp uint32 `msgpack:"ts"` // 4 bytes (unix timestamp)
}

// Size comparison:
// JSON:        ~80 bytes {"type":"cursor_move","user_id":"uuid",...}
// MessagePack: ~15 bytes (80% reduction!)

func (c *Client) sendBinary(msg BinaryMessage) error {
    data, err := msgpack.Marshal(msg)
    if err != nil {
        return err
    }
    return c.conn.WriteMessage(websocket.BinaryMessage, data)
}
```

#### Protocol Negotiation

```go
// client.go - Support both JSON and binary
type Client struct {
    // ... existing fields ...
    supportsBinary bool
    userIDMapping  uint32 // Local ID untuk binary protocol
}

// Detect client capabilities on connection
func (c *Client) negotiateProtocol() {
    // Check Sec-WebSocket-Protocol header
    protocols := c.conn.Subprotocol()
    c.supportsBinary = strings.Contains(protocols, "msgpack")
}

// Adaptive encoding
func (c *Client) sendMessage(msg Message) error {
    if c.supportsBinary && isBinaryOptimizable(msg.Type) {
        return c.sendBinary(toBinaryMessage(msg))
    }
    return c.sendJSON(msg)
}
```

### 8.5 Horizontal Scaling with Redis Pub/Sub

#### Architecture

```
                    ┌─────────────────┐
                    │   Load Balancer │
                    │  (Sticky Session)│
                    └────────┬─────────┘
          ┌──────────────────┼──────────────────┐
     ┌────┴────┐        ┌────┴────┐        ┌────┴────┐
     │ Go #1   │        │ Go #2   │        │ Go #3   │
     │ :8080   │        │ :8081   │        │ :8082   │
     └────┬────┘        └────┬────┘        └────┬────┘
          │                  │                  │
          │     Redis Pub/Sub Channel            │
          └──────────────────┼──────────────────┘
                        ┌────┴─────┐
                        │  Redis   │
                        │ Cluster  │
                        └──────────┘

Room "survey:abc123" has clients across all 3 instances:
- Instance 1: User A, User B
- Instance 2: User C
- Instance 3: User D

When User A sends message → Go #1 publishes to Redis 
→ Go #2 & #3 receive → broadcast to their local clients
```

#### Implementation

```go
// hub.go - Multi-instance support
type Hub struct {
    // ... existing fields ...
    redis      *redis.Client
    instanceID string // Unique instance identifier
}

func (h *Hub) subscribeToRedisPubSub() {
    pubsub := h.redis.Subscribe(context.Background(), "collab:*")
    defer pubsub.Close()

    for msg := range pubsub.Channel() {
        // Parse message
        var collabMsg Message
        if err := json.Unmarshal([]byte(msg.Payload), &collabMsg); err != nil {
            continue
        }

        // Find room and broadcast to local clients
        h.mu.RLock()
        room, exists := h.rooms[collabMsg.RoomID]
        h.mu.RUnlock()

        if exists {
            room.broadcastLocal(collabMsg)
        }
    }
}

// room.go - Publish to Redis when broadcasting
func (r *Room) broadcast(msg Message) {
    // Broadcast to local clients
    r.broadcastLocal(msg)

    // Publish to Redis for other instances
    if r.hub.redis != nil {
        msgBytes, _ := json.Marshal(msg)
        r.hub.redis.Publish(
            context.Background(),
            "collab:"+r.ID,
            msgBytes,
        )
    }
}

func (r *Room) broadcastLocal(msg Message) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    for client := range r.clients {
        select {
        case client.send <- msg:
        default:
            // Buffer full, skip
        }
    }
}
```

#### Sticky Session Configuration (Nginx)

```nginx
upstream websocket_backend {
    # IP hash untuk sticky session
    ip_hash;

    server go1:8080 max_fails=3 fail_timeout=30s;
    server go2:8081 max_fails=3 fail_timeout=30s;
    server go3:8082 max_fails=3 fail_timeout=30s;
}

server {
    listen 443 ssl http2;
    server_name api.atlas-food.com;

    # WebSocket upgrade
    location /ws/collab {
        proxy_pass http://websocket_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

        # Timeouts
        proxy_read_timeout 3600s;
        proxy_send_timeout 3600s;

        # Disable buffering untuk WebSocket
        proxy_buffering off;
    }
}
```

### 8.6 Client-Side Optimizations

#### Connection Management

```typescript
// lib/websocket-manager.ts
class WebSocketManager {
    private ws: WebSocket | null = null;
    private reconnectAttempts = 0;
    private messageQueue: Message[] = [];
    private isReconnecting = false;

    connect(roomId: string) {
        const token = getCookie('atlas_token');
        const url = `${WS_URL}/ws/collab?token=${token}&room=${roomId}`;
        
        this.ws = new WebSocket(url);
        this.setupHandlers();
    }

    private setupHandlers() {
        this.ws!.onopen = () => {
            console.log('✅ WebSocket connected');
            this.reconnectAttempts = 0;
            this.flushMessageQueue(); // Send queued messages
        };

        this.ws!.onclose = (event) => {
            if (event.code !== 1000) { // Not normal closure
                this.reconnect();
            }
        };

        this.ws!.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        this.ws!.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            this.handleMessage(msg);
        };
    }

    private reconnect() {
        if (this.isReconnecting) return;
        this.isReconnecting = true;

        // Exponential backoff: 1s, 2s, 4s, 8s, 16s, 30s (max)
        const delay = Math.min(
            1000 * Math.pow(2, this.reconnectAttempts),
            30000
        );

        console.log(`🔄 Reconnecting in ${delay}ms...`);
        
        setTimeout(() => {
            this.reconnectAttempts++;
            this.isReconnecting = false;
            this.connect(this.currentRoomId);
        }, delay);
    }

    send(type: string, payload: any) {
        const msg = { type, payload, timestamp: Date.now() };

        if (this.ws?.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(msg));
        } else {
            // Queue message for retry
            this.messageQueue.push(msg);
            
            // Limit queue size (drop oldest)
            if (this.messageQueue.length > 100) {
                this.messageQueue.shift();
            }
        }
    }

    private flushMessageQueue() {
        while (this.messageQueue.length > 0) {
            const msg = this.messageQueue.shift()!;
            this.ws!.send(JSON.stringify(msg));
        }
    }
}
```

#### Efficient State Updates

```typescript
// store/collab-store.ts - Zustand with immer
import create from 'zustand';
import { immer } from 'zustand/middleware/immer';

interface CollabState {
    users: Map<string, CollabUser>;
    cursors: Map<string, CursorPosition>;
    activities: ActivityEntry[];
    
    // Actions
    updateCursor: (userId: string, cursor: CursorPosition) => void;
    addActivity: (activity: ActivityEntry) => void;
}

export const useCollabStore = create<CollabState>()(
    immer((set) => ({
        users: new Map(),
        cursors: new Map(),
        activities: [],

        updateCursor: (userId, cursor) => set((state) => {
            // Immer allows mutation-like syntax (efficient)
            state.cursors.set(userId, cursor);
        }),

        addActivity: (activity) => set((state) => {
            state.activities.unshift(activity);
            // Keep last 50 activities
            if (state.activities.length > 50) {
                state.activities.pop();
            }
        }),
    }))
);
```

#### Virtual Rendering for Many Cursors

```tsx
// components/LiveCursorOverlay.tsx
import { useVirtualizer } from '@tanstack/react-virtual';

function LiveCursorOverlay() {
    const cursors = useCollabStore(s => s.cursors);
    const visibleCursors = useMemo(() => {
        // Only render cursors in viewport + 100px margin
        const viewport = {
            left: window.scrollX - 100,
            top: window.scrollY - 100,
            right: window.scrollX + window.innerWidth + 100,
            bottom: window.scrollY + window.innerHeight + 100,
        };

        return Array.from(cursors.values()).filter(c => 
            c.x >= viewport.left && c.x <= viewport.right &&
            c.y >= viewport.top && c.y <= viewport.bottom
        );
    }, [cursors]);

    return (
        <div className="fixed inset-0 pointer-events-none z-50">
            {visibleCursors.map(cursor => (
                <Cursor key={cursor.userId} {...cursor} />
            ))}
        </div>
    );
}

// Memoized cursor component
const Cursor = memo(({ userId, x, y, name, color }: CursorProps) => (
    <div
        className="absolute transition-transform duration-75"
        style={{
            transform: `translate(${x}px, ${y}px)`,
            willChange: 'transform', // GPU acceleration
        }}
    >
        <svg width="20" height="20" viewBox="0 0 20 20" fill={color}>
            <path d="M3 1l15 12-6 2-3 5-3-1 3-6-6-2z" />
        </svg>
        <span
            className="ml-4 rounded px-2 py-0.5 text-xs text-white"
            style={{ backgroundColor: color }}
        >
            {name}
        </span>
    </div>
));
```

### 8.7 Monitoring & Metrics

```go
// metrics/websocket.go
import "github.com/prometheus/client_golang/prometheus"

var (
    ActiveConnections = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "websocket_active_connections",
        Help: "Current number of active WebSocket connections",
    })

    MessagesReceived = prometheus.NewCounterVec(prometheus.CounterOpts{
        Name: "websocket_messages_received_total",
        Help: "Total messages received by type",
    }, []string{"type"})

    MessagesSent = prometheus.NewCounterVec(prometheus.CounterOpts{
        Name: "websocket_messages_sent_total",
        Help: "Total messages sent by type",
    }, []string{"type"})

    MessageLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "websocket_message_latency_seconds",
        Help:    "Message processing latency",
        Buckets: prometheus.DefBuckets,
    }, []string{"type"})

    RoomSize = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "websocket_room_size",
        Help:    "Number of clients per room",
        Buckets: []float64{1, 2, 5, 10, 20, 50, 100},
    }, []string{"room_type"})
)

// Track in client.go
func (c *Client) writePump() {
    defer func() {
        ActiveConnections.Dec()
        c.conn.Close()
    }()

    ActiveConnections.Inc()
    
    for msg := range c.send {
        MessagesSent.WithLabelValues(msg.Type).Inc()
        // ... write logic
    }
}
```

---

## 9. Security

| Aspek | Implementasi |
|---|---|
| **Auth** | JWT token di query param saat WebSocket handshake. Token divalidasi di `AuthGuard` sebelum connection diterima |
| **Room Access Control** | Cek apakah user punya hak akses ke room (assigned ke survey / admin role) |
| **Input Validation** | Semua incoming message divalidasi — type check, size limit (max 64KB), field whitelist |
| **Rate Limiting** | Max 50 msg/sec per connection. Lebih dari itu → disconnect + temp ban 1 menit |
| **TLS** | WSS (WebSocket Secure) via TLS, sama seperti HTTPS |
| **CORS / Origin Check** | WebSocket origin check: hanya allow dari domain frontend yang terdaftar |
| **SQL Injection** | Semua via GORM parameterized queries — tidak ada raw SQL |
| **XSS** | Display name & message content di-escape di frontend. Tidak render raw HTML dari WebSocket messages |
| **Token Exposure** | JWT token di query param hanya saat handshake awal — tidak dikirim di setiap message body |

---

## 10. Frontend Component Integration Points

### 10.1 Survey Wizard Flow

```
app/surveys/[accessToken]/
├── layout.tsx          ← CollaborationBar + LiveCursorOverlay
├── recall/
│   └── page.tsx        ← Init useWebSocket(roomId)
├── add-food/
│   └── page.tsx        ← CollaborativeSearch + ActivityFeed
├── portion/
│   └── page.tsx        ← Broadcast portion_set + LiveCursorOverlay
├── review/
│   └── page.tsx        ← Activity feed + broadcast review_submit
└── done/
    └── page.tsx        ← Disconnect WebSocket, room closed
```

### 10.2 Admin Panel Flow

```
app/admin/
├── layout.tsx           ← CollaborationBar (auto-join admin:food-db room)
├── foods/
│   ├── page.tsx         ← LiveCursorOverlay + LockIndicator
│   └── [id]/
│       └── page.tsx     ← db_edit_start → db_edit_field → db_edit_save
├── surveys/
│   └── [id]/
│       └── page.tsx     ← Collaborative survey editing
└── ...
```

### 10.3 Collaboration Invite Flow

```
┌─────────────────────────────────────────────────────┐
│  CollaborationBar                                     │
│  ┌──────────┬──────────┬──────────┬───────────────┐ │
│  │ 👤 Avatars│ 📋 Share │ 🔔 Feed  │ ⏺ Recording  │ │
│  └──────────┴──────────┴──────────┴───────────────┘ │
│                                                       │
│  "Share" → Generate invite link:                      │
│  /surveys/{token}/collab?join={sessionId}             │
│                                                       │
│  Avatars → Who's online + cursor colors               │
│  Feed → Slide-out activity panel                      │
└─────────────────────────────────────────────────────┘
```

---

## 11. Implementation Phases

### Phase 1 — Foundation (Core WebSocket)
- [ ] Setup WebSocket handler di Go (gorilla/websocket)
- [ ] Implement Hub + Room + Client struct
- [ ] JWT auth guard untuk WebSocket
- [ ] `useWebSocket` hook di frontend
- [ ] Zustand collab store
- [ ] Basic connect/disconnect + heartbeat
- [ ] Presence: join/leave/list

### Phase 2 — Live Cursors
- [ ] `cursor_move` message handling
- [ ] Cursor throttling (15fps)
- [ ] `LiveCursorOverlay` component
- [ ] `PresenceAvatars` component
- [ ] Warna unik per user

### Phase 3 — Collaborative Survey
- [ ] `food_search` + `food_select` shared
- [ ] `meal_add` + `portion_set` broadcast
- [ ] `ActivityFeed` component
- [ ] `CollaborationBar` component
- [ ] Integrate ke survey wizard pages

### Phase 4 — Collaborative Food DB
- [ ] Optimistic locking (Redis-based)
- [ ] `db_edit_*` message flow
- [ ] `LockIndicator` component
- [ ] Version conflict resolution
- [ ] Integrate ke admin panel

### Phase 5 — Polish & Scale
- [ ] Redis pub/sub untuk multi-instance
- [ ] Message batching & compression
- [ ] Rate limiting
- [ ] Activity log persistence
- [ ] Reconnection with state sync
- [ ] Load testing (k6 WebSocket / artillery)

---

## 12. Dependencies to Add

### Backend (Go)
```
go get github.com/gorilla/websocket
go get github.com/go-redis/redis/v9
go get github.com/google/uuid
```

### Frontend (Next.js)
```bash
# WebSocket API built-in di browser, Zustand sudah ada

npm install uuid
npm install @types/uuid -D

# Optional: auto-reconnect wrapper
npm install reconnecting-websocket
```

---

## 13. Best Practices Checklist

### Connection Management
- [✅] Exponential backoff reconnection (1s → 2s → 4s → ... max 30s)
- [✅] Heartbeat / ping-pong setiap 15 detik
- [✅] Graceful close dengan kode 1000 (normal) atau 1001 (going away)
- [✅] Cleanup semua goroutine & event listener saat disconnect
- [✅] Connection pooling di Redis

### State Management
- [✅] Single source of truth: server adalah authority
- [✅] Sequence number untuk deteksi out-of-order messages
- [✅] Full state sync saat reconnect (bukan delta replay)
- [✅] Client optimistic update + server confirmation (rollback on error)

### Performance
- [✅] Throttle cursor updates (max 15 fps)
- [✅] Batch broadcast (kirim per 50ms, bukan per message)
- [✅] Limit message size (64KB max)
- [✅] Ring buffer untuk room history (last 100 messages)
- [✅] Redis TTL untuk ephemeral data (presence, locks, cursor)

### Security
- [✅] JWT validation on every connection
- [✅] Room-level authorization
- [✅] Input validation (size, type, field whitelist)
- [✅] Rate limiting per connection
- [✅] WSS (TLS) for production
- [✅] No sensitive data in WebSocket messages (passwords, full tokens)

### User Experience
- [✅] "Reconnecting..." overlay saat koneksi terputus
- [✅] Warna unik per user untuk cursor
- [✅] Smooth cursor animation (CSS transition)
- [✅] Activity feed yang informatif tapi tidak mengganggu
- [✅] Lock indicator yang jelas (siapa yang lagi edit apa)
- [✅] Invite link yang mudah di-copy & share

### Observability
- [ ] Metrics: active connections, messages/sec, room count
- [ ] Logging: setiap connect/disconnect/error dengan structured logging
- [ ] Alerting: spike di disconnect rate, latency > threshold
- [ ] Health check endpoint untuk WebSocket hub

---

## 14. Testing & Debugging WebSocket

### 14.1 Unit Testing (Go)

```go
// collab/hub_test.go
package collab_test

import (
    "testing"
    "time"
    "github.com/gorilla/websocket"
    "github.com/stretchr/testify/assert"
)

func TestRoomBroadcast(t *testing.T) {
    hub := NewHub()
    go hub.Run()

    // Create mock clients
    room := hub.GetOrCreateRoom("test-room-1")
    
    client1 := &Client{
        hub:  hub,
        room: room,
        send: make(chan Message, 256),
    }
    
    client2 := &Client{
        hub:  hub,
        room: room,
        send: make(chan Message, 256),
    }

    room.register <- client1
    room.register <- client2

    // Broadcast message
    msg := Message{
        Type:    MessageBroadcast,
        Payload: map[string]string{"test": "data"},
    }
    
    room.broadcast <- msg

    // Verify both clients received
    select {
    case received := <-client1.send:
        assert.Equal(t, msg.Type, received.Type)
    case <-time.After(1 * time.Second):
        t.Fatal("client1 didn't receive message")
    }

    select {
    case received := <-client2.send:
        assert.Equal(t, msg.Type, received.Type)
    case <-time.After(1 * time.Second):
        t.Fatal("client2 didn't receive message")
    }
}

func TestPresenceTracking(t *testing.T) {
    hub := NewHub()
    room := hub.GetOrCreateRoom("test-room-2")

    client := &Client{
        userID:   "user-123",
        userName: "Test User",
        userRole: "admin",
    }

    // Join
    room.addClient(client)
    assert.Equal(t, 1, len(room.clients))
    assert.True(t, room.presence.IsUserInRoom("user-123"))

    // Leave
    room.removeClient(client)
    assert.Equal(t, 0, len(room.clients))
    assert.False(t, room.presence.IsUserInRoom("user-123"))
}

func TestMessageRateLimiting(t *testing.T) {
    client := &Client{
        rateLimiter: NewRateLimiter(10, time.Second), // 10 msg/sec
    }

    // Send 10 messages - should succeed
    for i := 0; i < 10; i++ {
        assert.True(t, client.rateLimiter.Allow())
    }

    // 11th message - should fail
    assert.False(t, client.rateLimiter.Allow())

    // Wait 1 second, should work again
    time.Sleep(1 * time.Second)
    assert.True(t, client.rateLimiter.Allow())
}
```

### 14.2 Integration Testing (Go + WebSocket Client)

```go
// collab/integration_test.go
func TestWebSocketHandshake(t *testing.T) {
    // Start test server
    server := setupTestServer()
    defer server.Close()

    // Generate JWT token
    token := generateTestJWT("user-123", "researcher")

    // Connect as WebSocket client
    wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + 
            "/ws/collab?token=" + token + "&room=test-room"

    ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    assert.NoError(t, err)
    defer ws.Close()

    // Send join message
    joinMsg := Message{
        Type:   MessageJoinRoom,
        RoomID: "test-room",
    }
    
    err = ws.WriteJSON(joinMsg)
    assert.NoError(t, err)

    // Expect join_success response
    var response Message
    err = ws.ReadJSON(&response)
    assert.NoError(t, err)
    assert.Equal(t, MessageJoinSuccess, response.Type)
}

func TestConcurrentConnections(t *testing.T) {
    server := setupTestServer()
    defer server.Close()

    const numClients = 100
    var wg sync.WaitGroup

    for i := 0; i < numClients; i++ {
        wg.Add(1)
        go func(clientID int) {
            defer wg.Done()

            token := generateTestJWT(fmt.Sprintf("user-%d", clientID), "researcher")
            wsURL := fmt.Sprintf("%s/ws/collab?token=%s&room=stress-test",
                server.URL, token)

            ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
            if err != nil {
                t.Errorf("Client %d failed to connect: %v", clientID, err)
                return
            }
            defer ws.Close()

            // Send cursor move
            msg := Message{
                Type: MessageCursorMove,
                Payload: map[string]int{
                    "x": clientID * 10,
                    "y": clientID * 10,
                },
            }
            
            ws.WriteJSON(msg)
            time.Sleep(100 * time.Millisecond)
        }(i)
    }

    wg.Wait()
}
```

### 14.3 End-to-End Testing (Playwright)

```typescript
// e2e/collaboration.spec.ts
import { test, expect } from '@playwright/test';

test('real-time cursor sync between users', async ({ browser }) => {
    // Create two browser contexts (2 users)
    const context1 = await browser.newContext();
    const context2 = await browser.newContext();

    const page1 = await context1.newPage();
    const page2 = await context2.newPage();

    // Login both users
    await page1.goto('/login');
    await page1.fill('[name="email"]', 'researcher@test.com');
    await page1.fill('[name="password"]', 'password123');
    await page1.click('button[type="submit"]');

    await page2.goto('/login');
    await page2.fill('[name="email"]', 'respondent@test.com');
    await page2.fill('[name="password"]', 'password123');
    await page2.click('button[type="submit"]');

    // Join same survey room
    const surveyURL = '/surveys/test-survey-123/recall';
    await page1.goto(surveyURL);
    await page2.goto(surveyURL);

    // Wait for WebSocket connection
    await page1.waitForSelector('[data-ws-status="connected"]');
    await page2.waitForSelector('[data-ws-status="connected"]');

    // Move cursor on page1
    await page1.mouse.move(500, 300);

    // Check if page2 sees the cursor
    await expect(page2.locator('[data-cursor-user="researcher@test.com"]'))
        .toBeVisible({ timeout: 3000 });

    // Verify cursor position
    const cursorElement = page2.locator('[data-cursor-user="researcher@test.com"]');
    const boundingBox = await cursorElement.boundingBox();
    
    expect(boundingBox?.x).toBeCloseTo(500, 50); // Within 50px tolerance
    expect(boundingBox?.y).toBeCloseTo(300, 50);

    await context1.close();
    await context2.close();
});

test('collaborative food search', async ({ browser }) => {
    const context1 = await browser.newContext();
    const context2 = await browser.newContext();

    const page1 = await context1.newPage();
    const page2 = await context2.newPage();

    // Setup (login + join room)
    // ... similar to previous test ...

    // User 1 searches for food
    await page1.fill('[data-testid="food-search"]', 'nasi goreng');

    // User 2 should see the search query
    await expect(page2.locator('[data-testid="collab-search-query"]'))
        .toHaveText('nasi goreng', { timeout: 3000 });

    // User 1 selects a food
    await page1.click('[data-food-id="12345"]');

    // User 2 should see activity feed
    await expect(page2.locator('[data-testid="activity-feed"]'))
        .toContainText('researcher@test.com added Nasi Goreng');

    await context1.close();
    await context2.close();
});
```

### 14.4 Load Testing (k6)

```javascript
// loadtest/websocket-stress.js
import ws from 'k6/ws';
import { check } from 'k6';

export let options = {
    stages: [
        { duration: '1m', target: 100 },  // Ramp up to 100 users
        { duration: '3m', target: 100 },  // Stay at 100 users
        { duration: '1m', target: 500 },  // Spike to 500 users
        { duration: '2m', target: 500 },  // Hold spike
        { duration: '1m', target: 0 },    // Ramp down
    ],
};

export default function () {
    const url = 'ws://localhost:8080/ws/collab';
    const token = getJWTToken(); // Helper function
    const roomId = `room-${Math.floor(Math.random() * 10)}`; // 10 rooms

    const res = ws.connect(
        `${url}?token=${token}&room=${roomId}`,
        function (socket) {
            socket.on('open', () => {
                console.log('Connected');

                // Send cursor moves at 15fps
                socket.setInterval(() => {
                    socket.send(JSON.stringify({
                        type: 'cursor_move',
                        payload: {
                            x: Math.random() * 1920,
                            y: Math.random() * 1080,
                            page: '/survey',
                        },
                    }));
                }, 66); // 15fps = 66ms
            });

            socket.on('message', (data) => {
                const msg = JSON.parse(data);
                check(msg, {
                    'message has type': (m) => m.type !== undefined,
                    'message has payload': (m) => m.payload !== undefined,
                });
            });

            socket.on('close', () => console.log('Disconnected'));
            socket.on('error', (e) => console.error('Error:', e));

            // Keep connection for 60 seconds
            socket.setTimeout(() => {
                socket.close();
            }, 60000);
        }
    );

    check(res, { 'status is 101': (r) => r && r.status === 101 });
}
```

### 14.5 Debugging Tools

#### Browser DevTools

```typescript
// Add debug logger in development
if (process.env.NODE_ENV === 'development') {
    window.wsDebug = {
        messageLog: [],
        
        logMessage(direction: 'sent' | 'received', msg: Message) {
            const entry = {
                direction,
                type: msg.type,
                payload: msg.payload,
                timestamp: Date.now(),
            };
            
            this.messageLog.push(entry);
            
            // Keep last 100 messages
            if (this.messageLog.length > 100) {
                this.messageLog.shift();
            }
            
            console.log(`[WS ${direction}]`, msg.type, msg.payload);
        },
        
        getStats() {
            const stats = {
                total: this.messageLog.length,
                byType: {},
                avgInterval: 0,
            };
            
            this.messageLog.forEach(entry => {
                stats.byType[entry.type] = (stats.byType[entry.type] || 0) + 1;
            });
            
            return stats;
        },
        
        clear() {
            this.messageLog = [];
        },
    };
}
```

#### Server-Side Logging

```go
// middleware/websocket_logger.go
type WSLogger struct {
    logger *zap.Logger
}

func (l *WSLogger) LogConnection(clientID, roomID, userID string) {
    l.logger.Info("WebSocket connected",
        zap.String("client_id", clientID),
        zap.String("room_id", roomID),
        zap.String("user_id", userID),
        zap.Time("connected_at", time.Now()),
    )
}

func (l *WSLogger) LogMessage(direction string, msg Message) {
    l.logger.Debug("WebSocket message",
        zap.String("direction", direction),
        zap.String("type", string(msg.Type)),
        zap.String("room_id", msg.RoomID),
        zap.String("user_id", msg.UserID),
        zap.Any("payload", msg.Payload),
    )
}

func (l *WSLogger) LogDisconnection(clientID string, reason string) {
    l.logger.Info("WebSocket disconnected",
        zap.String("client_id", clientID),
        zap.String("reason", reason),
        zap.Duration("session_duration", time.Since(client.connectedAt)),
    )
}
```

#### Monitoring Dashboard

```go
// handler.go - Add debug endpoints
func (h *Handler) GetStats(c *gin.Context) {
    h.hub.mu.RLock()
    defer h.hub.mu.RUnlock()

    stats := map[string]interface{}{
        "active_connections": len(h.hub.clients),
        "active_rooms":       len(h.hub.rooms),
        "rooms": make([]RoomInfo, 0),
    }

    for roomID, room := range h.hub.rooms {
        room.mu.RLock()
        stats["rooms"] = append(stats["rooms"].([]RoomInfo), RoomInfo{
            ID:          roomID,
            ClientCount: len(room.clients),
            CreatedAt:   room.createdAt,
        })
        room.mu.RUnlock()
    }

    c.JSON(http.StatusOK, stats)
}

// Route
router.GET("/api/v1/debug/websocket/stats", handler.GetStats)
```

#### Chrome Extension for WebSocket Inspection

```json
// manifest.json for Chrome extension
{
  "name": "Atlas Food WS Inspector",
  "version": "1.0",
  "manifest_version": 3,
  "devtools_page": "devtools.html",
  "permissions": ["webRequest"],
  "host_permissions": ["ws://localhost:8080/*", "wss://*.atlas-food.com/*"]
}
```

```javascript
// devtools.js
chrome.devtools.network.onRequestFinished.addListener(request => {
    if (request.request.url.includes('/ws/collab')) {
        // Log WebSocket handshake details
        console.log('WebSocket connection:', {
            url: request.request.url,
            status: request.response.status,
            headers: request.response.headers,
        });
    }
});
```

---

## 16. Common Issues & Troubleshooting

### 16.1 Connection Issues

#### Problem: WebSocket fails to connect (status 400/403)

**Symptoms:**
- Console error: `WebSocket connection failed`
- Network tab shows 400 Bad Request or 403 Forbidden

**Causes & Solutions:**

```typescript
// ❌ WRONG: Token expired or missing
const ws = new WebSocket('ws://localhost/ws/collab?room=123');

// ✅ CORRECT: Always include valid JWT token
const token = getCookie('atlas_token');
if (!token || isTokenExpired(token)) {
    await refreshToken(); // Refresh first
}
const ws = new WebSocket(`ws://localhost/ws/collab?token=${token}&room=123`);
```

**Backend validation:**
```go
// handler.go - Proper error responses
func (h *Handler) UpgradeConnection(c *gin.Context) {
    token := c.Query("token")
    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
        return
    }

    claims, err := utils.ValidateJWT(token)
    if err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
        return
    }

    // ... continue with upgrade
}
```

#### Problem: Connection drops after 60 seconds

**Symptoms:**
- Connection works initially but disconnects after ~1 minute
- No error message, just silent disconnect

**Cause:** Missing heartbeat/ping-pong mechanism

**Solution:**
```go
// client.go - Implement ping/pong
func (c *Client) writePump() {
    ticker := time.NewTicker(15 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case msg := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.conn.WriteJSON(msg); err != nil {
                return
            }

        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func (c *Client) readPump() {
    c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })

    // ... read messages
}
```

#### Problem: CORS error on WebSocket connection

**Symptoms:**
- Browser console: `WebSocket connection to 'ws://...' failed: Error during WebSocket handshake: Unexpected response code: 403`

**Solution:**
```go
// handler.go - Configure CORS for WebSocket
upgrader := websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")
        allowedOrigins := []string{
            "http://localhost:3000",
            "https://atlas-food.com",
            "https://app.atlas-food.com",
        }
        
        for _, allowed := range allowedOrigins {
            if origin == allowed {
                return true
            }
        }
        
        return false // Reject unknown origins
    },
}
```

### 16.2 Performance Issues

#### Problem: High CPU usage on server with many connections

**Symptoms:**
- CPU usage spikes when >100 users connected
- Slow message delivery

**Diagnosis:**
```go
// Add profiling endpoint
import _ "net/http/pprof"

router.GET("/debug/pprof/*action", gin.WrapH(http.DefaultServeMux))

// Run profiling:
// go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30
```

**Common causes:**

1. **Inefficient broadcasting:**
```go
// ❌ BAD: Serialize message for each client
for client := range room.clients {
    msg := Message{...}
    msgBytes, _ := json.Marshal(msg) // Repeated!
    client.send <- msgBytes
}

// ✅ GOOD: Serialize once, send to all
msgBytes, _ := json.Marshal(msg)
for client := range room.clients {
    select {
    case client.send <- msgBytes:
    default:
        // Skip slow clients
    }
}
```

2. **Too many goroutines:**
```go
// ❌ BAD: Spawn goroutine per message
for client := range room.clients {
    go func(c *Client) {
        c.send <- msg
    }(client)
}

// ✅ GOOD: Use buffered channel, no goroutine spawn
for client := range room.clients {
    select {
    case client.send <- msg:
    default:
        // Buffer full, skip
    }
}
```

#### Problem: Memory leak - RAM keeps growing

**Diagnosis:**
```bash
# Memory profiling
go tool pprof http://localhost:8080/debug/pprof/heap

# Inside pprof:
(pprof) top10
(pprof) list functionName
```

**Common causes:**

1. **Goroutine leak:**
```go
// ❌ BAD: No cleanup on disconnect
func (c *Client) start() {
    go c.readPump()
    go c.writePump()
}

// ✅ GOOD: Track and cleanup
func (c *Client) start(ctx context.Context) {
    var wg sync.WaitGroup
    
    wg.Add(2)
    go func() {
        defer wg.Done()
        c.readPump(ctx)
    }()
    
    go func() {
        defer wg.Done()
        c.writePump(ctx)
    }()
    
    // Wait for both to finish
    go func() {
        wg.Wait()
        c.cleanup()
    }()
}
```

2. **Unbounded message history:**
```go
// ❌ BAD: History grows forever
room.history = append(room.history, msg)

// ✅ GOOD: Use ring buffer (fixed size)
if len(room.history) >= maxHistory {
    room.history = room.history[1:]
}
room.history = append(room.history, msg)
```

3. **Zombie rooms:**
```go
// ❌ BAD: Empty rooms never cleaned up
func (r *Room) removeClient(client *Client) {
    delete(r.clients, client)
    // Room still exists in hub!
}

// ✅ GOOD: Delete empty rooms
func (r *Room) removeClient(client *Client) {
    delete(r.clients, client)
    
    if len(r.clients) == 0 {
        r.hub.deleteRoom(r.ID)
    }
}
```

#### Problem: Slow message delivery (high latency)

**Symptoms:**
- Cursor movements lag by 2-3 seconds
- Activity feed updates delayed

**Diagnosis:**
```go
// Add latency tracking
func (c *Client) writePump() {
    for msg := range c.send {
        start := time.Now()
        err := c.conn.WriteJSON(msg)
        latency := time.Since(start)
        
        if latency > 100*time.Millisecond {
            log.Warnf("Slow write: %v for user %s", latency, c.userID)
        }
    }
}
```

**Solutions:**

1. **Client send buffer full:**
```go
// Increase buffer size
client := &Client{
    send: make(chan Message, 1024), // Increased from 256
}

// Or drop messages for slow clients
select {
case client.send <- msg:
case <-time.After(10 * time.Millisecond):
    // Client too slow, drop message
    metrics.DroppedMessages.Inc()
}
```

2. **Network congestion:**
```go
// Enable WebSocket compression
upgrader := websocket.Upgrader{
    EnableCompression: true,
}

// Or switch to binary protocol (MessagePack)
```

### 16.3 State Synchronization Issues

#### Problem: Clients see different room state

**Symptoms:**
- User A sees 3 people in room, User B sees 5
- Cursor positions out of sync

**Cause:** Missed messages or race conditions

**Solution: Full state sync on reconnect**
```go
// client.go
func (c *Client) onReconnect() {
    // Send full room state
    presence := c.room.getAllUsers()
    cursors := c.room.getAllCursors()
    
    c.send <- Message{
        Type: "state_sync",
        Payload: map[string]interface{}{
            "users":   presence,
            "cursors": cursors,
            "version": c.room.version,
        },
    }
}
```

```typescript
// Frontend
ws.addEventListener('message', (event) => {
    const msg = JSON.parse(event.data);
    
    if (msg.type === 'state_sync') {
        // Replace entire state, don't merge
        collabStore.setState({
            users: new Map(msg.payload.users),
            cursors: new Map(msg.payload.cursors),
            version: msg.payload.version,
        });
    }
});
```

#### Problem: Optimistic locking conflicts

**Symptoms:**
- User edits food data, gets `VERSION_CONFLICT` error
- Changes lost after another user saves

**Solution: Implement proper conflict resolution**
```go
// service.go
func (s *Service) UpdateFood(foodID string, updates FoodUpdate, version int) error {
    var food Food
    result := s.db.Model(&Food{}).
        Where("id = ? AND version = ?", foodID, version).
        Updates(updates).
        Update("version", version+1)
    
    if result.RowsAffected == 0 {
        // Version mismatch - conflict!
        return ErrVersionConflict
    }
    
    return nil
}
```

```typescript
// Frontend - Retry with latest version
async function saveFood(foodId: string, updates: FoodUpdate) {
    let retries = 3;
    
    while (retries > 0) {
        try {
            const food = await getFood(foodId); // Get latest version
            await updateFood(foodId, updates, food.version);
            return; // Success
        } catch (err) {
            if (err.code === 'VERSION_CONFLICT') {
                retries--;
                if (retries === 0) {
                    // Show merge UI to user
                    showMergeDialog(foodId, updates);
                }
            } else {
                throw err;
            }
        }
    }
}
```

### 16.4 Debugging Checklist

When WebSocket not working as expected, check:

- [ ] **Backend logs:** Any errors during connection upgrade?
- [ ] **JWT token:** Is it valid and not expired?
- [ ] **Network tab:** What's the WebSocket handshake status code?
- [ ] **Console logs:** Any JavaScript errors?
- [ ] **Redis connection:** Is Redis running and accessible?
- [ ] **Room authorization:** Does user have permission to join room?
- [ ] **Message format:** Is JSON valid? Check payload structure
- [ ] **Rate limiting:** Is client being throttled?
- [ ] **Firewall/Proxy:** Is WebSocket traffic allowed?
- [ ] **Load balancer:** Is sticky session configured?

**Quick diagnostic:**
```bash
# Test WebSocket endpoint directly
wscat -c "ws://localhost:8080/ws/collab?token=YOUR_JWT&room=test-123"

# Should connect and allow sending JSON:
> {"type": "ping"}
< {"type": "pong", "timestamp": "2026-07-11T..."}
```

## 17. Production Deployment Checklist

### 17.1 Infrastructure Requirements

- [ ] **Load Balancer:** Nginx/HAProxy with WebSocket support + sticky sessions
- [ ] **SSL/TLS:** WSS (secure WebSocket) dengan valid certificate
- [ ] **Redis Cluster:** For horizontal scaling (3+ nodes recommended)
- [ ] **Monitoring:** Prometheus + Grafana for metrics
- [ ] **Logging:** Centralized logging (ELK/Loki)
- [ ] **Auto-scaling:** Based on active connection count

### 17.2 Security Hardening

```go
// Production WebSocket config
var ProdWSConfig = WebSocketConfig{
    // Stricter timeouts
    WriteTimeout:  5 * time.Second,
    ReadTimeout:   30 * time.Second,
    PongTimeout:   5 * time.Second,
    PingInterval:  10 * time.Second,
    
    // Tighter limits
    MaxMessageSize:    32768, // 32KB (reduced from 64KB)
    MaxMessagesPerSec: 30,    // Reduced from 50
    
    // Connection limits
    MaxConnectionsPerUser: 5, // Prevent abuse
    MaxConnectionsPerIP:   20,
    
    // Rate limiting
    RateLimitWindow:     time.Second,
    RateLimitBurst:      10,
    
    // Security
    RequireAuth:         true,
    AllowedOrigins:      []string{"https://app.atlas-food.com"},
    CheckOrigin:         true,
    EnableCompression:   true,
}
```

### 17.3 Environment Variables

```bash
# .env.production
WS_HOST=0.0.0.0
WS_PORT=8080
WS_MAX_CONNECTIONS=10000
WS_ENABLE_COMPRESSION=true

REDIS_URL=redis://redis-cluster:6379
REDIS_PASSWORD=***
REDIS_POOL_SIZE=100

JWT_SECRET=***
JWT_EXPIRY=24h

LOG_LEVEL=info
ENABLE_DEBUG_ENDPOINTS=false

METRICS_PORT=9090
HEALTH_CHECK_PORT=8081
```

### 17.4 Docker Configuration

```dockerfile
# Dockerfile - Multi-stage build
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8080 9090
CMD ["./main"]
```

```yaml
# docker-compose.yml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
      - "9090:9090"
    environment:
      - REDIS_URL=redis://redis:6379
      - WS_MAX_CONNECTIONS=5000
    depends_on:
      - redis
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '2'
          memory: 1G
        reservations:
          cpus: '1'
          memory: 512M

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    command: redis-server --appendonly yes

  nginx:
    image: nginx:alpine
    ports:
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - api

volumes:
  redis-data:
```

### 17.5 Nginx Configuration

```nginx
# nginx.conf
upstream websocket_backend {
    ip_hash; # Sticky sessions
    
    server api1:8080 max_fails=3 fail_timeout=30s;
    server api2:8080 max_fails=3 fail_timeout=30s;
    server api3:8080 max_fails=3 fail_timeout=30s;
}

server {
    listen 443 ssl http2;
    server_name api.atlas-food.com;

    # SSL config
    ssl_certificate /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # WebSocket location
    location /ws/collab {
        proxy_pass http://websocket_backend;
        
        # WebSocket headers
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Timeouts
        proxy_connect_timeout 7d;
        proxy_send_timeout 7d;
        proxy_read_timeout 7d;

        # Disable buffering
        proxy_buffering off;
        proxy_cache off;

        # Limits
        client_max_body_size 64k;
        limit_conn conn_limit_per_ip 20;
        limit_req zone=req_limit_per_ip burst=10 nodelay;
    }

    # REST API
    location /api/ {
        proxy_pass http://websocket_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Health check
    location /health {
        access_log off;
        return 200 "healthy\n";
    }
}

# Rate limiting zones
http {
    limit_conn_zone $binary_remote_addr zone=conn_limit_per_ip:10m;
    limit_req_zone $binary_remote_addr zone=req_limit_per_ip:10m rate=30r/s;
}
```

### 17.6 Kubernetes Deployment (Optional)

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: atlas-food-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: atlas-food-api
  template:
    metadata:
      labels:
        app: atlas-food-api
    spec:
      containers:
      - name: api
        image: atlas-food-api:latest
        ports:
        - containerPort: 8080
          name: websocket
        - containerPort: 9090
          name: metrics
        env:
        - name: REDIS_URL
          valueFrom:
            secretKeyRef:
              name: redis-secret
              key: url
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "2000m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: atlas-food-api-service
spec:
  type: LoadBalancer
  sessionAffinity: ClientIP # Sticky sessions
  selector:
    app: atlas-food-api
  ports:
  - name: websocket
    port: 443
    targetPort: 8080
  - name: metrics
    port: 9090
    targetPort: 9090
```

### 17.7 Monitoring & Alerts

```yaml
# prometheus-alerts.yaml
groups:
- name: websocket_alerts
  interval: 30s
  rules:
  
  # High connection count
  - alert: HighWebSocketConnections
    expr: websocket_active_connections > 8000
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High WebSocket connection count"
      description: "{{ $value }} active connections (threshold: 8000)"

  # Connection spike
  - alert: ConnectionSpike
    expr: rate(websocket_connections_total[1m]) > 100
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "WebSocket connection spike"
      description: "{{ $value }} new connections per second"

  # High message drop rate
  - alert: HighMessageDropRate
    expr: rate(websocket_messages_dropped_total[5m]) > 10
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "High WebSocket message drop rate"
      description: "{{ $value }} messages dropped per second"

  # High latency
  - alert: HighWebSocketLatency
    expr: histogram_quantile(0.95, rate(websocket_message_latency_seconds_bucket[5m])) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High WebSocket message latency"
      description: "P95 latency: {{ $value }}s (threshold: 1s)"

  # Instance down
  - alert: WebSocketInstanceDown
    expr: up{job="atlas-food-api"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "WebSocket instance down"
      description: "Instance {{ $labels.instance }} is down"
```

### 17.8 Graceful Shutdown

```go
// main.go - Handle shutdown gracefully
func main() {
    // ... setup ...
    
    hub := collab.NewHub()
    go hub.Run()

    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    // Start server in goroutine
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    // 1. Stop accepting new connections
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // 2. Notify all connected clients
    hub.BroadcastShutdown("Server maintenance in progress")

    // 3. Wait for existing connections to close (max 30s)
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("Server forced shutdown: %v", err)
    }

    // 4. Close Redis connection
    hub.redis.Close()

    log.Println("Server exited")
}
```

```go
// hub.go
func (h *Hub) BroadcastShutdown(message string) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    shutdownMsg := Message{
        Type: MessageDisconnect,
        Payload: map[string]string{
            "reason":  message,
            "code":    "1001", // Going away
            "reconnect": "true",
        },
    }

    for _, room := range h.rooms {
        room.broadcast <- shutdownMsg
    }

    // Give clients time to receive message
    time.Sleep(2 * time.Second)

    // Close all connections
    for _, room := range h.rooms {
        for client := range room.clients {
            client.conn.Close()
        }
    }
}
```

### 17.9 Rollback Plan

If WebSocket feature causes issues in production:

1. **Feature flag:** Disable collab via environment variable
```go
if os.Getenv("ENABLE_COLLAB") != "true" {
    log.Println("Collaboration disabled via feature flag")
    return
}
```

2. **Circuit breaker:** Auto-disable if error rate > threshold
```go
if errorRate > 0.05 { // 5% error rate
    enableCollab = false
    log.Error("Collaboration disabled due to high error rate")
}
```

3. **Gradual rollout:** Enable for percentage of users
```go
func shouldEnableCollab(userID string) bool {
    // Enable for 10% of users (A/B test)
    hash := md5.Sum([]byte(userID))
    return hash[0] < 26 // 26/256 ≈ 10%
}
```

---

## 18. Notes & Caveats

1. **Next.js App Router + WebSocket:** Karena App Router pakai React Server Components, WebSocket hanya bisa dijalankan di client component (gunakan `'use client'` directive).

2. **SSR Consideration:** Komponen `LiveCursorOverlay` dan `PresenceAvatars` harus di-render client-side only. Gunakan `dynamic(() => import('./Component'), { ssr: false })`.

3. **Go Goroutine Leak:** Setiap koneksi WebSocket spawn 2 goroutine (read + write pump). Pastikan semua goroutine di-cleanup saat koneksi close. Gunakan `context.Context` dengan `defer cancel()`.

4. **Redis sebagai Dependency:** Phase 1-4 bisa jalan tanpa Redis (single instance Go). Redis baru wajib di Phase 5 untuk horizontal scaling.

5. **Testing WebSocket:** Gunakan `gorilla/websocket` client di Go test, atau Playwright untuk E2E test di browser. Load testing dengan k6 atau artillery.

6. **Mobile Consideration:** WebSocket protocol sama persis untuk mobile app. Tidak perlu perubahan di backend.

7. **Backward Compatibility:** Semua REST API tetap jalan. WebSocket adalah tambahan, bukan pengganti. Fitur tanpa kolaborasi tetap berfungsi normal.

8. **Message Size:** Limit ke 64KB per message untuk mencegah abuse. Untuk data besar, kirim via REST API lalu notifikasi via WebSocket.

9. **Binary Protocol:** MessagePack bisa hemat 50-80% bandwidth tapi tambah kompleksitas. Mulai dengan JSON, optimize nanti kalau perlu.

10. **Cursor Accuracy:** Cursor position relative ke viewport, bukan absolute ke document. Butuh sync scroll position untuk akurasi penuh.

11. **Safari Compatibility:** Safari kadang close WebSocket connection saat tab di-background. Butuh reconnection logic yang robust.

12. **Rate Limiting Strategy:** Per-user limit (50 msg/sec) + per-IP limit (500 msg/sec) untuk mencegah abuse tapi allow multiple tabs.

---

## 19. Updated Implementation Phases

### ✅ Phase 0 — Planning & Design (DONE)
- [x] Architecture planning
- [x] WebSocket protocol design
- [x] Security & performance analysis
- [x] Documentation complete

### 🔄 Phase 1 — Foundation (Week 1-2)
**Backend:**
- [ ] Setup WebSocket handler di Go (gorilla/websocket)
- [ ] Implement Hub + Room + Client struct
- [ ] JWT auth guard untuk WebSocket
- [ ] Basic connect/disconnect + heartbeat
- [ ] Presence: join/leave/list messages
- [ ] Unit tests untuk core functionality

**Frontend:**
- [ ] `useWebSocket` hook di frontend
- [ ] Zustand collab store setup
- [ ] Connection manager dengan auto-reconnect
- [ ] Basic presence indicators

**Success Criteria:**
- WebSocket dapat connect/disconnect dengan JWT auth
- Multiple users dapat join same room
- Presence list updates real-time

### 🔄 Phase 2 — Live Cursors (Week 3)
- [ ] `cursor_move` message handling (backend)
- [ ] Cursor throttling 15fps (frontend)
- [ ] `LiveCursorOverlay` component
- [ ] `PresenceAvatars` component dengan warna unik
- [ ] Viewport filtering untuk performance
- [ ] Integration dengan survey pages

**Success Criteria:**
- 10+ users dapat lihat cursor real-time dengan <200ms latency
- CPU usage <50% dengan 50 concurrent users

### 🔄 Phase 3 — Collaborative Survey (Week 4-5)
- [ ] `food_search` + `food_select` broadcast
- [ ] `meal_add` + `portion_set` messages
- [ ] `ActivityFeed` component
- [ ] `CollaborationBar` component
- [ ] Survey wizard integration
- [ ] Room auto-join untuk survey sessions
- [ ] Activity log persistence (DB)

**Success Criteria:**
- Researcher dapat bantu respondent pilih makanan real-time
- Activity feed shows all user actions
- Survey data konsisten across clients

### 🔄 Phase 4 — Collaborative Food DB (Week 6-7)
- [ ] Optimistic locking dengan Redis
- [ ] `db_edit_*` message flow
- [ ] `LockIndicator` component
- [ ] Version conflict resolution UI
- [ ] Admin panel integration
- [ ] Edit history tracking

**Success Criteria:**
- Multiple admin dapat edit food DB tanpa konflik
- Locks expire setelah 5 menit idle
- No data loss dari concurrent edits

### 🔄 Phase 5 — Polish & Scale (Week 8-9)
- [ ] Redis pub/sub untuk multi-instance
- [ ] Message batching & compression
- [ ] Rate limiting implementation
- [ ] Binary protocol (MessagePack) untuk cursor
- [ ] Full state sync on reconnect
- [ ] Load testing (target: 500 concurrent users)
- [ ] Monitoring & alerting setup
- [ ] Production deployment

**Success Criteria:**
- Handle 500+ concurrent connections per instance
- P95 message latency <500ms
- Horizontal scaling works (3+ instances)
- Zero downtime deployment

### 🔄 Phase 6 — Advanced Features (Future)
- [ ] Voice/video chat integration (WebRTC)
- [ ] Screen sharing untuk researcher-respondent
- [ ] Collaborative text editing (CRDT-based)
- [ ] Replay recorded sessions
- [ ] AI-powered suggestions based on collab data
- [ ] Mobile app support (React Native)

---

## 20. Performance Benchmarks (Target)

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Connection Time** | <500ms | Time from connect() to first message |
| **Message Latency (P50)** | <100ms | Server receive → all clients receive |
| **Message Latency (P95)** | <500ms | 95th percentile |
| **Cursor Update Rate** | 15 fps | Actual updates rendered |
| **Max Concurrent Users/Room** | 100 | Without performance degradation |
| **Max Concurrent Connections/Instance** | 5,000 | Single Go instance |
| **CPU Usage (50 users)** | <30% | Single core |
| **Memory Usage (1000 users)** | <500MB | Total heap allocation |
| **Message Throughput** | 10,000 msg/sec | Single instance broadcast |
| **Reconnection Time** | <2s | Auto-reconnect after disconnect |
| **State Sync Time** | <1s | Full room state on reconnect |

**Load Test Command:**
```bash
# Test with k6
k6 run --vus 500 --duration 5m loadtest/websocket-stress.js

# Expected results:
# ✓ 95% of messages delivered within 500ms
# ✓ 0% connection failures
# ✓ <1% message drop rate
```

---

## 21. Useful Resources

### Documentation
- [gorilla/websocket](https://github.com/gorilla/websocket) - Go WebSocket library
- [WebSocket API (MDN)](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket) - Browser API docs
- [Redis Pub/Sub](https://redis.io/docs/manual/pubsub/) - For multi-instance scaling

### Tools
- [wscat](https://github.com/websockets/wscat) - CLI WebSocket client untuk testing
- [k6](https://k6.io/) - Load testing tool dengan WebSocket support
- [Prometheus](https://prometheus.io/) - Metrics & monitoring
- [Grafana](https://grafana.com/) - Visualization dashboard

### Tutorials & Examples
- [Yjs Demos](https://docs.yjs.dev/) - CRDT collaborative editing examples
- [Socket.io vs WebSocket](https://ably.com/blog/socketio-vs-websocket) - Protocol comparison
- [Scaling WebSockets](https://socket.io/docs/v4/using-multiple-nodes/) - Multi-instance patterns

---

> **Next Steps:**  
> 1. Review planning dengan tim development  
> 2. Setup development environment (Go + Redis lokal)  
> 3. Start Phase 1: Implement `handler.go`, `hub.go`, `client.go`  
> 4. Create feature branch: `feature/realtime-collaboration`  
> 5. Daily sync untuk tracking progress

**Estimasi Total Development Time:** 9-10 minggu (1 developer full-time)  
**Priority:** High (major feature untuk research collaboration use case)  
**Risk Level:** Medium (new technology stack, butuh thorough testing)
