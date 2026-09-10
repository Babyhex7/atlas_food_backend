package collab

import (
	"sync"
	"testing"
	"time"
)

// TestSendToUnregisteredClientDoesNotPanic - regresi untuk crash paling berbahaya
// di domain ini.
//
// Dulu peserta lain mengirim langsung ke `client.send` (mis. saat follow/unfollow
// atau update role) sementara hub menutup channel itu di unregister. Kalau
// keduanya berbarengan, prosesnya panic "send on closed channel" — dan karena
// ReadPump berjalan di goroutine telanjang tanpa recover, SELURUH server mati.
func TestSendToUnregisteredClientDoesNotPanic(t *testing.T) {
	hub := NewHub()

	client := newTestClient(hub, "room-panic", "user-1", "Budi")
	hub.registerClient(client)
	hub.unregisterClient(client)

	// Setelah unregister, siapa pun boleh mencoba mengirim — harus no-op, bukan panic.
	if client.trySend(newMessage(MsgPong, "room-panic", "user-1", "Budi", nil)) {
		t.Fatal("trySend ke client yang sudah keluar seharusnya gagal, bukan berhasil")
	}
	client.sendQuiet(newMessage(MsgPong, "room-panic", "user-1", "Budi", nil))

	// closeSend ganda juga tidak boleh panic "close of closed channel".
	client.closeSend()
	client.closeSend()
}

// TestConcurrentSendAndUnregister - jalankan pengirim dan unregister bersamaan
// berulang kali. Tanpa gerbang sendOpen, ini panic dalam beberapa iterasi.
func TestConcurrentSendAndUnregister(t *testing.T) {
	hub := NewHub()

	for i := 0; i < 200; i++ {
		client := newTestClient(hub, "room-race", "user-x", "Budi")
		hub.registerClient(client)

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				client.sendQuiet(newMessage(MsgPong, "room-race", "user-x", "Budi", nil))
			}
		}()
		go func() {
			defer wg.Done()
			hub.unregisterClient(client)
		}()
		wg.Wait()
	}
}

// TestRoleAccessorsAreConcurrencySafe - RoomRole ditulis goroutine HTTP
// (UpdateUserRole) sambil dibaca ReadPump lewat canEdit().
func TestRoleAccessorsAreConcurrencySafe(t *testing.T) {
	hub := NewHub()
	owner := newTestClient(hub, "room-role", "user-owner", "Owner")
	hub.registerClient(owner)
	target := newTestClient(hub, "room-role", "user-target", "Target")
	hub.registerClient(target)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			role := RoomRoleEditor
			if i%2 == 0 {
				role = RoomRoleViewer
			}
			_ = hub.UpdateUserRole("room-role", "user-owner", "user-target", role)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			_ = target.canEdit()
			_ = target.RoomRole()
		}
	}()
	wg.Wait()
}

// TestLocksAreScopedPerRoom - lock room A tidak boleh memblokir room B, dan
// snapshot satu room tidak boleh membocorkan lock room lain.
func TestLocksAreScopedPerRoom(t *testing.T) {
	lm := NewLockManager()

	if _, ok, _ := lm.TryLock("room-a", "food", "F-1", "user-1", "Budi", 1); !ok {
		t.Fatal("lock pertama di room-a harus berhasil")
	}
	if _, ok, _ := lm.TryLock("room-b", "food", "F-1", "user-2", "Siti", 1); !ok {
		t.Fatal("entity yang sama di room berbeda harus bisa dikunci terpisah")
	}
	if _, ok, existing := lm.TryLock("room-a", "food", "F-1", "user-3", "Andi", 1); ok || existing == nil {
		t.Fatal("user lain di room yang sama harus ditolak")
	}

	if got := lm.Snapshot("room-a"); len(got) != 1 || got[0].LockedBy != "user-1" {
		t.Fatalf("snapshot room-a harus berisi 1 lock miliknya sendiri, dapat %+v", got)
	}
}

// TestLockReleasedOnDisconnect - entity tidak boleh tertinggal terkunci saat
// pemegangnya putus koneksi.
func TestLockReleasedOnDisconnect(t *testing.T) {
	hub := NewHub()
	client := newTestClient(hub, "room-lock", "user-1", "Budi")
	hub.registerClient(client)

	if _, ok, _ := hub.Locks().TryLock("room-lock", "food", "F-9", "user-1", "Budi", 1); !ok {
		t.Fatal("lock awal harus berhasil")
	}

	hub.unregisterClient(client)

	if got := hub.Locks().Get("room-lock", "food", "F-9"); got != nil {
		t.Fatalf("lock harus dilepas saat pemiliknya keluar, masih ada: %+v", got)
	}
}

// TestExpiredLockCanBeTakenOver - jaring pengaman terakhir kalau unregister
// tidak pernah jalan (proses dibunuh, socket menggantung).
func TestExpiredLockCanBeTakenOver(t *testing.T) {
	lm := NewLockManager()
	lock, ok, _ := lm.TryLock("room-ttl", "food", "F-2", "user-1", "Budi", 1)
	if !ok {
		t.Fatal("lock awal harus berhasil")
	}

	// Mundurkan waktu lock melewati TTL
	lm.mu.Lock()
	lm.locks[lockKey("room-ttl", "food", "F-2")].LockedAt = time.Now().Add(-LockTTL - time.Minute)
	lm.mu.Unlock()

	if got := lm.Get("room-ttl", "food", "F-2"); got != nil {
		t.Fatal("lock kedaluwarsa harus dianggap tidak ada")
	}
	if _, ok, _ := lm.TryLock("room-ttl", "food", "F-2", "user-2", "Siti", 1); !ok {
		t.Fatalf("user lain harus bisa mengambil alih lock kedaluwarsa milik %s", lock.LockedBy)
	}
	if swept := lm.SweepExpired(); len(swept) != 0 {
		t.Fatalf("lock yang sudah diambil alih tidak boleh ikut tersapu, dapat %+v", swept)
	}
}

// TestEmptyRoomSurvivesGracePeriod - room tidak boleh langsung dibuang begitu
// kosong; kalau dibuang, viewer cukup menunggu sepi lalu masuk lagi jadi owner.
func TestEmptyRoomSurvivesGracePeriod(t *testing.T) {
	hub := NewHub()
	client := newTestClient(hub, "room-grace", "user-viewer", "Viewer")
	client.SetRoomRole(RoomRoleViewer)
	client.RoleFromInvite = true
	hub.registerClient(client)
	hub.unregisterClient(client)

	hub.cleanupInactiveRooms()

	if got := hub.RoomRoleOf("room-grace", "user-viewer"); got != RoomRoleViewer {
		t.Fatalf("role harus tetap diingat selama grace period, dapat %q", got)
	}
	if got := hub.ResolveRoomRole("room-grace", "user-viewer", ""); got != RoomRoleViewer {
		t.Fatalf("viewer tidak boleh naik jadi %q hanya karena room sepi sesaat", got)
	}
}
