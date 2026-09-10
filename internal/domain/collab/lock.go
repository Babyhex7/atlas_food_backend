package collab

import (
	"sync"
	"time"
)

// LockTTL - umur maksimal sebuah lock tanpa aktivitas.
//
// Tanpa batas ini, client yang putus mendadak (tutup tab, jaringan mati)
// meninggalkan entity terkunci selamanya sampai server di-restart.
const LockTTL = 5 * time.Minute

// EntityLock represents an optimistic lock on a food/survey entity during collaborative editing.
type EntityLock struct {
	RoomID     string    `json:"room_id"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	LockedBy   string    `json:"locked_by"`
	Username   string    `json:"username"`
	Version    int       `json:"version"`
	LockedAt   time.Time `json:"locked_at"`
}

// LockManager holds in-memory locks (Redis deferred).
//
// Lock di-scope PER ROOM: dua room yang kebetulan mengedit food yang sama
// tidak boleh saling memblokir, dan snapshot state satu room tidak boleh
// membocorkan lock milik room lain.
type LockManager struct {
	mu    sync.RWMutex
	locks map[string]*EntityLock // key: roomID|entityType:entityID
}

// NewLockManager - buat LockManager kosong
func NewLockManager() *LockManager {
	return &LockManager{
		locks: make(map[string]*EntityLock),
	}
}

// lockKey - susun key map lock dengan format "roomID|entityType:entityID"
func lockKey(roomID, entityType, entityID string) string {
	return roomID + "|" + entityType + ":" + entityID
}

// expired - true kalau lock sudah melewati TTL tanpa aktivitas
func (l *EntityLock) expired(now time.Time) bool {
	return now.Sub(l.LockedAt) > LockTTL
}

// TryLock attempts to acquire a lock. Returns the lock and whether acquisition succeeded.
// Lock milik orang lain yang sudah lewat TTL dianggap tidak ada dan boleh diambil alih.
func (lm *LockManager) TryLock(roomID, entityType, entityID, userID, username string, version int) (*EntityLock, bool, *EntityLock) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	now := time.Now()
	key := lockKey(roomID, entityType, entityID)
	if existing, ok := lm.locks[key]; ok && !existing.expired(now) {
		if existing.LockedBy == userID {
			existing.LockedAt = now
			if version > 0 {
				existing.Version = version
			}
			cp := *existing
			return &cp, true, nil
		}
		cp := *existing
		return nil, false, &cp
	}

	lock := &EntityLock{
		RoomID:     roomID,
		EntityType: entityType,
		EntityID:   entityID,
		LockedBy:   userID,
		Username:   username,
		Version:    version,
		LockedAt:   now,
	}
	if lock.Version <= 0 {
		lock.Version = 1
	}
	lm.locks[key] = lock
	cp := *lock
	return &cp, true, nil
}

// Release removes a lock if owned by userID (or force if userID empty).
func (lm *LockManager) Release(roomID, entityType, entityID, userID string) bool {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	key := lockKey(roomID, entityType, entityID)
	existing, ok := lm.locks[key]
	if !ok {
		return false
	}
	if userID != "" && existing.LockedBy != userID {
		return false
	}
	delete(lm.locks, key)
	return true
}

// ReleaseAllByUser - lepas seluruh lock milik satu user di sebuah room.
//
// Dipanggil saat socket terakhir user itu terputus, supaya entity tidak
// tertinggal dalam keadaan terkunci. Mengembalikan lock yang dilepas agar
// hub bisa menyiarkan db_unlocked ke peserta lain.
func (lm *LockManager) ReleaseAllByUser(roomID, userID string) []EntityLock {
	if roomID == "" || userID == "" {
		return nil
	}
	lm.mu.Lock()
	defer lm.mu.Unlock()

	released := make([]EntityLock, 0)
	for key, lock := range lm.locks {
		if lock.RoomID == roomID && lock.LockedBy == userID {
			released = append(released, *lock)
			delete(lm.locks, key)
		}
	}
	return released
}

// ReleaseRoom - buang seluruh lock milik sebuah room (dipakai saat room dibersihkan)
func (lm *LockManager) ReleaseRoom(roomID string) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	for key, lock := range lm.locks {
		if lock.RoomID == roomID {
			delete(lm.locks, key)
		}
	}
}

// SweepExpired - buang lock yang sudah lewat TTL; dipanggil ticker hub.
// Mengembalikan lock yang dibuang supaya hub bisa memberi tahu room terkait.
func (lm *LockManager) SweepExpired() []EntityLock {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	now := time.Now()
	expired := make([]EntityLock, 0)
	for key, lock := range lm.locks {
		if lock.expired(now) {
			expired = append(expired, *lock)
			delete(lm.locks, key)
		}
	}
	return expired
}

// Get returns current lock if any (nil kalau tidak ada atau sudah kedaluwarsa).
func (lm *LockManager) Get(roomID, entityType, entityID string) *EntityLock {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	if lock, ok := lm.locks[lockKey(roomID, entityType, entityID)]; ok && !lock.expired(time.Now()) {
		cp := *lock
		return &cp
	}
	return nil
}

// BumpVersion increments version after save.
func (lm *LockManager) BumpVersion(roomID, entityType, entityID, userID string) (*EntityLock, bool) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	key := lockKey(roomID, entityType, entityID)
	existing, ok := lm.locks[key]
	if !ok || existing.LockedBy != userID {
		return nil, false
	}
	existing.Version++
	existing.LockedAt = time.Now()
	cp := *existing
	return &cp, true
}

// Snapshot returns locks milik satu room saja (untuk presence/state sync).
func (lm *LockManager) Snapshot(roomID string) []EntityLock {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	now := time.Now()
	out := make([]EntityLock, 0, len(lm.locks))
	for _, l := range lm.locks {
		if l.RoomID == roomID && !l.expired(now) {
			out = append(out, *l)
		}
	}
	return out
}

// Count - jumlah lock aktif di seluruh hub (untuk endpoint stats)
func (lm *LockManager) Count() int {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	now := time.Now()
	n := 0
	for _, l := range lm.locks {
		if !l.expired(now) {
			n++
		}
	}
	return n
}
