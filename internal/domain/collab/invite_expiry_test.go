package collab

import (
	"testing"
	"time"
)

// TestInviteExpiredRejected - token yang melewati TTL tidak boleh lagi diterima (skenario S9 artikel E2).
func TestInviteExpiredRejected(t *testing.T) {
	store := NewInviteStore()
	inv := store.Create("room-exp", RoomRoleEditor, "owner", 30*time.Millisecond)

	if _, ok := store.Get(inv.Token); !ok {
		t.Fatal("token harus valid sebelum TTL habis")
	}
	time.Sleep(60 * time.Millisecond)
	if _, ok := store.Get(inv.Token); ok {
		t.Fatal("token kedaluwarsa tidak boleh diterima")
	}
}

// TestInviteRevokedRejected - token yang dicabut langsung tidak berlaku (skenario S8).
func TestInviteRevokedRejected(t *testing.T) {
	store := NewInviteStore()
	inv := store.Create("room-rev", RoomRoleEditor, "owner", time.Hour)
	store.Revoke(inv.Token)
	if _, ok := store.Get(inv.Token); ok {
		t.Fatal("token yang dicabut tidak boleh diterima")
	}
}
