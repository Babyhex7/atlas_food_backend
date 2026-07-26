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
