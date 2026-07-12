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

```
CURSOR MOVES:
  Client: throttle → max 15 fps (66ms interval)
  Server: batch → kirim ke room setiap 50ms (akumulasi semua cursor moves)

FOOD SEARCH:
  Client: debounce → 300ms setelah user berhenti ngetik
  Server: cache search results per query (Redis, TTL 30s)

PRESENCE:
  Server: broadcast setiap 10s (bukan setiap join/leave langsung)
  Atau: broadcast langsung untuk join/leave, heartbeat setiap 30s
```

### 8.2 Connection Pooling

```go
// Redis connection pool
RedisPool: &redis.Pool{
    MaxIdle:     20,
    MaxActive:   100,
    IdleTimeout: 240 * time.Second,
}

// WebSocket config
Upgrader: websocket.Upgrader{
    ReadBufferSize:    4096,
    WriteBufferSize:   4096,
    WriteBufferPool:   &sync.Pool{},
    EnableCompression: true,
}
```

### 8.3 Binary Messages (Optimasi Masa Depan)

Untuk bulk state sync: MessagePack atau Protocol Buffers — hemat bandwidth 30-50% vs JSON. MVP cukup JSON dulu.

### 8.4 Horizontal Scaling

```
                    ┌──────────┐
                    │   NLB    │
                    └────┬─────┘
          ┌──────────────┼──────────────┐
     ┌────┴────┐    ┌────┴────┐    ┌────┴────┐
     │  Go #1  │    │  Go #2  │    │  Go #3  │
     │  Hub A  │    │  Hub B  │    │  Hub C  │
     └────┬────┘    └────┬────┘    └────┬────┘
          │              │              │
          └──────────────┼──────────────┘
                    ┌────┴─────┐
                    │  Redis   │  ← Pub/Sub cross-instance broadcast
                    │ Cluster  │
                    └──────────┘
```

Sticky sessions via load balancer (cookie-based). Fallback: Redis pub/sub untuk cross-instance message routing.

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

## 14. Notes & Caveats

1. **Next.js App Router + WebSocket:** Karena App Router pakai React Server Components, WebSocket hanya bisa dijalankan di client component (gunakan `'use client'` directive).

2. **SSR Consideration:** Komponen `LiveCursorOverlay` dan `PresenceAvatars` harus di-render client-side only. Gunakan `dynamic(() => import('./Component'), { ssr: false })`.

3. **Go Goroutine Leak:** Setiap koneksi WebSocket spawn 2 goroutine (read + write pump). Pastikan semua goroutine di-cleanup saat koneksi close. Gunakan `context.Context` dengan `defer cancel()`.

4. **Redis sebagai Dependency:** Phase 1-4 bisa jalan tanpa Redis (single instance Go). Redis baru wajib di Phase 5 untuk horizontal scaling.

5. **Testing WebSocket:** Gunakan `gorilla/websocket` client di Go test, atau Playwright untuk E2E test di browser.

6. **Mobile Consideration:** WebSocket protocol sama persis untuk mobile app. Tidak perlu perubahan di backend.

7. **Backward Compatibility:** Semua REST API tetap jalan. WebSocket adalah tambahan, bukan pengganti. Fitur tanpa kolaborasi tetap berfungsi normal.

---

> **Next Step:** Review planning ini dengan tim, tentukan scope Phase 1, lalu mulai implementasi `internal/domain/collab/handler.go` di backend Go.
