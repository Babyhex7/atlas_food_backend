package collab

import (
	"sync"
	"time"
)

// Room represents a WebSocket room where users collaborate
type Room struct {
	ID             string
	clients        map[*Client]bool
	hub            *Hub
	messageHistory []*Message
	historyIndex   int
	historyMu      sync.RWMutex
	batchQueue     []*Message
	batchMu        sync.Mutex
	batchTicker    *time.Ticker
	lastBatchSent  time.Time
	stopCh         chan struct{}
	mu             sync.RWMutex

	// roles - room role yang sudah pernah diberikan ke sebuah user (key: userID).
	// Wajib diingat server: kalau role hanya mengandalkan ?invite= di URL, seorang
	// viewer yang pindah halaman (query invite ikut hilang) akan naik jadi editor.
	// Dilindungi mu yang sama dengan clients.
	roles map[string]string

	// canvasStrokes - riwayat coretan aktif di room ini (maksimal 50 stroke terakhir).
	// Memungkinkan late-joiner (peserta yang baru bergabung/rekonek) langsung
	// menerima gambar yang ada di layar tanpa minta pengirim ulang.
	canvasStrokes []*CanvasStrokeItem
	canvasMu      sync.RWMutex
}

// CanvasStrokeItem - struct penyimpan riwayat coretan canvas di server
type CanvasStrokeItem struct {
	StrokeID      string      `json:"stroke_id"`
	UserID        string      `json:"user_id"`
	Username      string      `json:"username"`
	Tool          string      `json:"tool"`
	Color         string      `json:"color"`
	Width         float64     `json:"width"`
	TargetImageID string      `json:"target_image_id,omitempty"`
	Points        [][]float64 `json:"points"`
	Timestamp     int64       `json:"timestamp"`
}

// NewRoom creates a new Room
func NewRoom(id string, hub *Hub) *Room {
	return &Room{
		ID:             id,
		clients:        make(map[*Client]bool),
		hub:            hub,
		messageHistory: make([]*Message, 100),
		batchQueue:     make([]*Message, 0, 50),
		batchTicker:    time.NewTicker(50 * time.Millisecond),
		stopCh:         make(chan struct{}),
		roles:          make(map[string]string),
	}
}

// RememberedRole - room role yang sudah tercatat untuk user ini, "" kalau belum ada
func (r *Room) RememberedRole(userID string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.roles[userID]
}

// Run starts the room's message batching loop (cursor moves)
func (r *Room) Run() {
	defer r.batchTicker.Stop()
	for {
		select {
		case <-r.batchTicker.C:
			r.flushBatchedMessages()
		case <-r.stopCh:
			return
		}
	}
}

// AddMessage adds a message (batch cursor updates; immediate for others)
func (r *Room) AddMessage(msg *Message) {
	if r.shouldBatchMessage(msg) {
		r.batchMu.Lock()
		r.batchQueue = append(r.batchQueue, msg)
		queueSize := len(r.batchQueue)
		r.batchMu.Unlock()
		if queueSize >= 50 {
			r.flushBatchedMessages()
		}
		return
	}
	r.hub.Publish(msg)
}

// shouldBatchMessage - true untuk pesan high-frequency (kursor/viewport/canvas laser/draw move) yang perlu di-batch dulu
func (r *Room) shouldBatchMessage(msg *Message) bool {
	return msg.Type == MsgCursorMove || msg.Type == MsgCursorUpdate ||
		msg.Type == MsgViewportSync ||
		msg.Type == MsgCanvasLaserMove || msg.Type == MsgCanvasLaserUpdated ||
		msg.Type == MsgCanvasDrawMove || msg.Type == MsgCanvasStrokeUpdated
}

// flushBatchedMessages - kirim pesan tertahan; cursor & laser di-coalesce ke latest per user
func (r *Room) flushBatchedMessages() {
	r.batchMu.Lock()
	if len(r.batchQueue) == 0 {
		r.batchMu.Unlock()
		return
	}
	messages := make([]*Message, len(r.batchQueue))
	copy(messages, r.batchQueue)
	r.batchQueue = r.batchQueue[:0]
	r.lastBatchSent = time.Now()
	r.batchMu.Unlock()

	latestCursor := make(map[string]*Message)
	latestLaser := make(map[string]*Message)
	var others []*Message
	for _, msg := range messages {
		if msg.Type == MsgCursorMove || msg.Type == MsgCursorUpdate {
			latestCursor[msg.UserID] = msg
			continue
		}
		if msg.Type == MsgCanvasLaserMove || msg.Type == MsgCanvasLaserUpdated {
			latestLaser[msg.UserID] = msg
			continue
		}
		others = append(others, msg)
	}
	for _, msg := range latestCursor {
		r.hub.Publish(msg)
	}
	for _, msg := range latestLaser {
		r.hub.Publish(msg)
	}
	for _, msg := range others {
		r.hub.Publish(msg)
	}
}

// AddCanvasStroke - simpan stroke baru ke riwayat room (maksimal 50 stroke)
func (r *Room) AddCanvasStroke(stroke *CanvasStrokeItem) {
	r.canvasMu.Lock()
	defer r.canvasMu.Unlock()
	if len(r.canvasStrokes) >= 50 {
		r.canvasStrokes = r.canvasStrokes[1:]
	}
	r.canvasStrokes = append(r.canvasStrokes, stroke)
}

// AppendCanvasPoints - tambahkan kumpulan titik koordinat baru ke strokeID yang ada
func (r *Room) AppendCanvasPoints(strokeID string, points [][]float64) {
	r.canvasMu.Lock()
	defer r.canvasMu.Unlock()
	for _, s := range r.canvasStrokes {
		if s.StrokeID == strokeID {
			s.Points = append(s.Points, points...)
			break
		}
	}
}

// ClearCanvasStrokes - hapus seluruh stroke canvas di room (atau per targetImageID tertentu)
func (r *Room) ClearCanvasStrokes(targetImageID string) {
	r.canvasMu.Lock()
	defer r.canvasMu.Unlock()
	if targetImageID == "" {
		r.canvasStrokes = nil
		return
	}
	filtered := make([]*CanvasStrokeItem, 0, len(r.canvasStrokes))
	for _, s := range r.canvasStrokes {
		if s.TargetImageID != targetImageID {
			filtered = append(filtered, s)
		}
	}
	r.canvasStrokes = filtered
}

// GetCanvasStrokes - ambil snapshot seluruh stroke canvas aktif untuk sync late-joiner
func (r *Room) GetCanvasStrokes() []*CanvasStrokeItem {
	r.canvasMu.RLock()
	defer r.canvasMu.RUnlock()
	strokes := make([]*CanvasStrokeItem, len(r.canvasStrokes))
	copy(strokes, r.canvasStrokes)
	return strokes
}

// addToHistory - simpan pesan ke ring buffer history (kursor & viewport diabaikan supaya tidak membanjiri)
func (r *Room) addToHistory(msg *Message) {
	// Don't flood history with cursor / viewport updates
	if msg.Type == MsgCursorUpdate || msg.Type == MsgCursorMove || msg.Type == MsgViewportSync {
		return
	}
	r.historyMu.Lock()
	defer r.historyMu.Unlock()
	r.messageHistory[r.historyIndex] = msg
	r.historyIndex = (r.historyIndex + 1) % len(r.messageHistory)
}

// GetHistory returns recent message history
func (r *Room) GetHistory(limit int) []*Message {
	r.historyMu.RLock()
	defer r.historyMu.RUnlock()

	if limit > len(r.messageHistory) {
		limit = len(r.messageHistory)
	}

	history := make([]*Message, 0, limit)
	for i := 0; i < len(r.messageHistory); i++ {
		idx := (r.historyIndex - len(r.messageHistory) + i + len(r.messageHistory)) % len(r.messageHistory)
		if r.messageHistory[idx] != nil {
			history = append(history, r.messageHistory[idx])
		}
	}
	if len(history) > limit {
		history = history[len(history)-limit:]
	}
	return history
}

// GetClientCount returns number of connected clients
func (r *Room) GetClientCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}

// Stop stops the room
func (r *Room) Stop() {
	close(r.stopCh)
}
