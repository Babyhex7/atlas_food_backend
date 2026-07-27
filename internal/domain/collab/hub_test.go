package collab

import (
	"testing"
	"time"
)

// newTestClient - Client minimal tanpa koneksi WebSocket asli.
// registerClient tidak menyentuh field conn, jadi nil aman untuk test hub.
func newTestClient(hub *Hub, roomID, userID, username string) *Client {
	return &Client{
		hub:      hub,
		RoomID:   roomID,
		UserID:   userID,
		Username: username,
		Viewport: map[string]interface{}{},
		send:     make(chan *Message, sendBufferSize),
		stopCh:   make(chan struct{}),
	}
}

// isClosed - cek apakah channel send sudah ditutup hub (tanda client ditendang)
func isClosed(ch chan *Message) bool {
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return true
			}
			// buang pesan presence/state lalu cek lagi
		default:
			return false
		}
	}
}

// TestSameUserTwoTabsStayConnected - satu user membuka dua tab harus tetap
// tersambung dua-duanya. Kalau socket lama ditendang, tab lama akan reconnect,
// menendang tab baru, lalu saling tendang tanpa henti (reconnect loop).
func TestSameUserTwoTabsStayConnected(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Stop()

	tabA := newTestClient(hub, "room-1", "user-1", "Budi")
	hub.registerClient(tabA)

	tabB := newTestClient(hub, "room-1", "user-1", "Budi")
	hub.registerClient(tabB)

	room := hub.GetOrCreateRoom("room-1")
	if got := room.GetClientCount(); got != 2 {
		t.Fatalf("dua tab user yang sama harus tetap tersambung, jumlah client = %d, mau 2", got)
	}
	if isClosed(tabA.send) {
		t.Fatal("socket tab pertama ditutup hub — ini yang memicu reconnect loop antar tab")
	}
}

// TestSameUserAppearsOnceInPresence - meski dua socket, presence list hanya
// menampilkan satu avatar per user (dedupe by user_id).
func TestSameUserAppearsOnceInPresence(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Stop()

	hub.registerClient(newTestClient(hub, "room-2", "user-1", "Budi"))
	hub.registerClient(newTestClient(hub, "room-2", "user-1", "Budi"))
	hub.registerClient(newTestClient(hub, "room-2", "user-2", "Sari"))

	room := hub.GetOrCreateRoom("room-2")
	msg := hub.buildPresenceList(room)
	users, _ := msg.Payload["users"].([]map[string]interface{})
	if len(users) != 2 {
		t.Fatalf("presence harus 2 user unik (Budi, Sari), dapat %d", len(users))
	}
}

// TestFirstClientBecomesOwner - client pertama di room jadi owner, sisanya editor.
func TestFirstClientBecomesOwner(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Stop()

	first := newTestClient(hub, "room-3", "user-1", "Budi")
	hub.registerClient(first)
	second := newTestClient(hub, "room-3", "user-2", "Sari")
	hub.registerClient(second)

	if first.RoomRole != RoomRoleOwner {
		t.Fatalf("client pertama harus owner, dapat %q", first.RoomRole)
	}
	if second.RoomRole != RoomRoleEditor {
		t.Fatalf("client kedua harus editor, dapat %q", second.RoomRole)
	}
}

// TestSecondTabInheritsRoomRole - tab kedua dari user yang sama mewarisi role
// room tab pertama (owner tidak turun jadi editor hanya karena buka tab baru).
func TestSecondTabInheritsRoomRole(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Stop()

	tabA := newTestClient(hub, "room-4", "user-1", "Budi")
	hub.registerClient(tabA)
	if tabA.RoomRole != RoomRoleOwner {
		t.Fatalf("tab pertama harus owner, dapat %q", tabA.RoomRole)
	}

	tabB := newTestClient(hub, "room-4", "user-1", "Budi")
	tabB.RoomRole = RoomRoleEditor // hasil ResolveRoomRole untuk room yang sudah ada isinya
	hub.registerClient(tabB)

	if tabB.RoomRole != RoomRoleOwner {
		t.Fatalf("tab kedua user yang sama harus mewarisi role owner, dapat %q", tabB.RoomRole)
	}
}

// TestViewerTetapViewerSaatPindahHalaman - viewer yang pindah halaman kehilangan
// query ?invite= di URL. Role-nya harus tetap viewer, bukan naik jadi editor.
func TestViewerTetapViewerSaatPindahHalaman(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Stop()

	// Owner sudah di dalam room
	hub.registerClient(newTestClient(hub, "room-6", "user-owner", "Budi"))

	// Viewer join lewat invite
	inv := hub.Invites().Create("room-6", RoomRoleViewer, "user-owner", time.Hour)
	viewer := newTestClient(hub, "room-6", "user-viewer", "Sari")
	viewer.RoomRole = hub.ResolveRoomRole("room-6", "user-viewer", inv.Token)
	hub.registerClient(viewer)
	if viewer.RoomRole != RoomRoleViewer {
		t.Fatalf("viewer harus masuk sebagai viewer, dapat %q", viewer.RoomRole)
	}

	// Pindah halaman: socket lama tutup, socket baru konek TANPA invite token
	hub.unregisterClient(viewer)
	rejoin := newTestClient(hub, "room-6", "user-viewer", "Sari")
	rejoin.RoomRole = hub.ResolveRoomRole("room-6", "user-viewer", "")
	hub.registerClient(rejoin)

	if rejoin.RoomRole != RoomRoleViewer {
		t.Fatalf("viewer naik jadi %q setelah pindah halaman — role harus tetap viewer", rejoin.RoomRole)
	}
	if rejoin.canEdit() {
		t.Fatal("viewer lolos pengecekan canEdit setelah pindah halaman")
	}
}

// TestOwnerTetapOwnerSaatPindahHalaman - owner yang reconnect tidak boleh turun
// jadi editor hanya karena room sudah berisi orang lain.
func TestOwnerTetapOwnerSaatPindahHalaman(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Stop()

	owner := newTestClient(hub, "room-7", "user-owner", "Budi")
	hub.registerClient(owner)
	hub.registerClient(newTestClient(hub, "room-7", "user-lain", "Sari"))

	hub.unregisterClient(owner)
	rejoin := newTestClient(hub, "room-7", "user-owner", "Budi")
	rejoin.RoomRole = hub.ResolveRoomRole("room-7", "user-owner", "")
	hub.registerClient(rejoin)

	if rejoin.RoomRole != RoomRoleOwner {
		t.Fatalf("owner turun jadi %q setelah pindah halaman", rejoin.RoomRole)
	}
}

// TestUnregisterKeepsOtherTab - menutup satu tab tidak boleh memutus tab lain
// milik user yang sama.
func TestUnregisterKeepsOtherTab(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Stop()

	tabA := newTestClient(hub, "room-5", "user-1", "Budi")
	hub.registerClient(tabA)
	tabB := newTestClient(hub, "room-5", "user-1", "Budi")
	hub.registerClient(tabB)

	hub.unregisterClient(tabA)
	time.Sleep(10 * time.Millisecond)

	room := hub.GetOrCreateRoom("room-5")
	if got := room.GetClientCount(); got != 1 {
		t.Fatalf("harus tersisa 1 client setelah satu tab ditutup, dapat %d", got)
	}
	if isClosed(tabB.send) {
		t.Fatal("tab yang masih terbuka ikut terputus")
	}
}
