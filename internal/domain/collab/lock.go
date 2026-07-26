package collab

import (
	"sync"
	"time"
)

// EntityLock represents an optimistic lock on a food/survey entity during collaborative editing.
type EntityLock struct {
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	LockedBy   string    `json:"locked_by"`
	Username   string    `json:"username"`
	Version    int       `json:"version"`
	LockedAt   time.Time `json:"locked_at"`
}

// LockManager holds in-memory locks (Redis deferred).
type LockManager struct {
	mu    sync.RWMutex
	locks map[string]*EntityLock // key: entityType:entityID
}

// NewLockManager - buat LockManager kosong
func NewLockManager() *LockManager {
	return &LockManager{
		locks: make(map[string]*EntityLock),
	}
}

// lockKey - susun key map lock dengan format "entityType:entityID"
func lockKey(entityType, entityID string) string {
	return entityType + ":" + entityID
}

// TryLock attempts to acquire a lock. Returns the lock and whether acquisition succeeded.
func (lm *LockManager) TryLock(entityType, entityID, userID, username string, version int) (*EntityLock, bool, *EntityLock) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	key := lockKey(entityType, entityID)
	if existing, ok := lm.locks[key]; ok {
		if existing.LockedBy == userID {
			existing.LockedAt = time.Now()
			if version > 0 {
				existing.Version = version
			}
			return existing, true, nil
		}
		return nil, false, existing
	}

	lock := &EntityLock{
		EntityType: entityType,
		EntityID:   entityID,
		LockedBy:   userID,
		Username:   username,
		Version:    version,
		LockedAt:   time.Now(),
	}
	if lock.Version <= 0 {
		lock.Version = 1
	}
	lm.locks[key] = lock
	return lock, true, nil
}

// Release removes a lock if owned by userID (or force if userID empty).
func (lm *LockManager) Release(entityType, entityID, userID string) bool {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	key := lockKey(entityType, entityID)
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

// Get returns current lock if any.
func (lm *LockManager) Get(entityType, entityID string) *EntityLock {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	if lock, ok := lm.locks[lockKey(entityType, entityID)]; ok {
		cp := *lock
		return &cp
	}
	return nil
}

// BumpVersion increments version after save.
func (lm *LockManager) BumpVersion(entityType, entityID, userID string) (*EntityLock, bool) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	key := lockKey(entityType, entityID)
	existing, ok := lm.locks[key]
	if !ok || existing.LockedBy != userID {
		return nil, false
	}
	existing.Version++
	existing.LockedAt = time.Now()
	cp := *existing
	return &cp, true
}

// Snapshot returns all locks (for presence/state sync).
func (lm *LockManager) Snapshot() []EntityLock {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	out := make([]EntityLock, 0, len(lm.locks))
	for _, l := range lm.locks {
		out = append(out, *l)
	}
	return out
}
