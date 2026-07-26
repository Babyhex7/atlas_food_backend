package collab

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// InviteToken — undangan room berbatas waktu (in-memory, best-effort).
type InviteToken struct {
	Token     string
	RoomID    string
	Role      string // editor | viewer
	CreatedBy string
	ExpiresAt time.Time
}

// InviteStore - penyimpanan token undangan di memori, aman diakses banyak goroutine
type InviteStore struct {
	mu      sync.RWMutex
	byToken map[string]*InviteToken
}

// NewInviteStore - buat InviteStore kosong
func NewInviteStore() *InviteStore {
	return &InviteStore{byToken: make(map[string]*InviteToken)}
}

// Create - buat token undangan baru untuk sebuah room (role default editor, TTL default 24 jam)
func (s *InviteStore) Create(roomID, role, createdBy string, ttl time.Duration) *InviteToken {
	if role != RoomRoleViewer && role != RoomRoleEditor {
		role = RoomRoleEditor
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	token := randomToken(16)
	inv := &InviteToken{
		Token:     token,
		RoomID:    roomID,
		Role:      role,
		CreatedBy: createdBy,
		ExpiresAt: time.Now().Add(ttl),
	}
	s.mu.Lock()
	s.byToken[token] = inv
	s.mu.Unlock()
	return inv
}

// Get - ambil token undangan; return false kalau tidak ada atau sudah kedaluwarsa
func (s *InviteStore) Get(token string) (*InviteToken, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	inv, ok := s.byToken[token]
	if !ok {
		return nil, false
	}
	if time.Now().After(inv.ExpiresAt) {
		return nil, false
	}
	return inv, true
}

// Revoke - hapus token undangan sehingga tidak bisa dipakai lagi
func (s *InviteStore) Revoke(token string) {
	s.mu.Lock()
	delete(s.byToken, token)
	s.mu.Unlock()
}

// randomToken - buat string hex acak sepanjang nBytes; fallback ke timestamp kalau crypto/rand gagal
func randomToken(nBytes int) string {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format("150405.000000000")))
	}
	return hex.EncodeToString(b)
}
