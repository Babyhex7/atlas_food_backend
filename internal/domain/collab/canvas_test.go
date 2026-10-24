package collab

import (
	"fmt"
	"testing"
	"time"
)

// waitForType - kuras channel send client sampai pesan bertipe msgType muncul.
// Pesan lain (presence, state_sync, dll.) dibuang. Gagal jika timeout.
func waitForType(t *testing.T, ch chan *Message, msgType string, timeout time.Duration) *Message {
	t.Helper()
	deadline := time.After(timeout)
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				t.Fatalf("channel ditutup sebelum menerima %q", msgType)
			}
			if msg.Type == msgType {
				return msg
			}
		case <-deadline:
			t.Fatalf("timeout menunggu pesan %q", msgType)
			return nil
		}
	}
}

// noMessageOfType - pastikan pesan bertipe msgType TIDAK diterima dalam jangka waktu tertentu.
func noMessageOfType(t *testing.T, ch chan *Message, msgType string, within time.Duration) {
	t.Helper()
	deadline := time.After(within)
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			if msg.Type == msgType {
				t.Fatalf("tidak boleh menerima %q, tapi diterima: %+v", msgType, msg.Payload)
			}
		case <-deadline:
			return
		}
	}
}

// setupCanvasRoom - room berisi owner (penggambar) dan viewer (penonton).
func setupCanvasRoom(t *testing.T, roomID string) (*Hub, *Client, *Client) {
	t.Helper()
	hub := NewHub()
	go hub.Run()

	owner := newTestClient(hub, roomID, "user-owner", "Enumerator")
	hub.registerClient(owner)

	viewer := newTestClient(hub, roomID, "user-viewer", "Responden")
	viewer.SetRoomRole(hub.ResolveRoomRole(roomID, "user-viewer", ""))
	hub.registerClient(viewer)

	if owner.RoomRole() != RoomRoleOwner || viewer.RoomRole() != RoomRoleViewer {
		t.Fatalf("setup role salah: owner=%q viewer=%q", owner.RoomRole(), viewer.RoomRole())
	}
	return hub, owner, viewer
}

// TestCanvasStrokeBroadcastToPeer - coretan owner (start + move) sampai ke peserta lain.
func TestCanvasStrokeBroadcastToPeer(t *testing.T) {
	hub, owner, viewer := setupCanvasRoom(t, "canvas-1")
	defer hub.Stop()

	owner.handleMessage(newMessage(MsgCanvasDrawStart, "canvas-1", owner.UserID, owner.Username, map[string]interface{}{
		"stroke_id": "str_1", "tool": "pencil", "color": "#EF4444", "width": 3.0, "x": 0.1, "y": 0.2,
	}))
	started := waitForType(t, viewer.send, MsgCanvasStrokeStarted, 2*time.Second)
	if started.Payload["stroke_id"] != "str_1" || started.Payload["color"] != "#EF4444" {
		t.Fatalf("payload stroke_started salah: %+v", started.Payload)
	}

	owner.handleMessage(newMessage(MsgCanvasDrawMove, "canvas-1", owner.UserID, owner.Username, map[string]interface{}{
		"stroke_id": "str_1",
		"points":    []interface{}{[]interface{}{0.3, 0.4}, []interface{}{0.5, 0.6}},
	}))
	updated := waitForType(t, viewer.send, MsgCanvasStrokeUpdated, 2*time.Second)
	if pts, _ := updated.Payload["points"].([][]float64); len(pts) != 2 {
		t.Fatalf("stroke_updated harus membawa 2 titik, dapat %+v", updated.Payload["points"])
	}

	// Server juga menyimpan titik untuk late-joiner.
	strokes := hub.GetOrCreateRoom("canvas-1").GetCanvasStrokes()
	if len(strokes) != 1 || len(strokes[0].Points) != 3 {
		t.Fatalf("ring buffer harus 1 stroke dengan 3 titik, dapat %+v", strokes)
	}
}

// TestCanvasSenderDoesNotReceiveEcho - pengirim tidak menerima siaran coretannya sendiri
// (render lokal optimistis di FE). Konsekuensi: E2E multi-browser wajib pakai 2 akun berbeda.
func TestCanvasSenderDoesNotReceiveEcho(t *testing.T) {
	hub, owner, _ := setupCanvasRoom(t, "canvas-2")
	defer hub.Stop()

	owner.handleMessage(newMessage(MsgCanvasDrawStart, "canvas-2", owner.UserID, owner.Username, map[string]interface{}{
		"stroke_id": "str_echo", "x": 0.5, "y": 0.5,
	}))
	noMessageOfType(t, owner.send, MsgCanvasStrokeStarted, 300*time.Millisecond)
}

// TestCanvasViewerForbidden - viewer ditolak (FORBIDDEN) untuk semua event canvas dan tidak ada yang tersiar.
func TestCanvasViewerForbidden(t *testing.T) {
	hub, owner, viewer := setupCanvasRoom(t, "canvas-3")
	defer hub.Stop()

	events := []string{MsgCanvasDrawStart, MsgCanvasDrawMove, MsgCanvasDrawEnd, MsgCanvasLaserMove, MsgCanvasClear}
	for _, ev := range events {
		viewer.handleMessage(newMessage(ev, "canvas-3", viewer.UserID, viewer.Username, map[string]interface{}{
			"stroke_id": "str_hack", "x": 0.1, "y": 0.1,
		}))
		errMsg := waitForType(t, viewer.send, MsgError, time.Second)
		if errMsg.Payload["code"] != "FORBIDDEN" {
			t.Fatalf("%s: kode error harus FORBIDDEN, dapat %+v", ev, errMsg.Payload)
		}
	}

	noMessageOfType(t, owner.send, MsgCanvasStrokeStarted, 300*time.Millisecond)
	if n := len(hub.GetOrCreateRoom("canvas-3").GetCanvasStrokes()); n != 0 {
		t.Fatalf("stroke viewer tidak boleh tersimpan, dapat %d", n)
	}
}

// TestCanvasClearBroadcast - clear menghapus riwayat server & tersiar ke peserta.
func TestCanvasClearBroadcast(t *testing.T) {
	hub, owner, viewer := setupCanvasRoom(t, "canvas-4")
	defer hub.Stop()

	owner.handleMessage(newMessage(MsgCanvasDrawStart, "canvas-4", owner.UserID, owner.Username, map[string]interface{}{
		"stroke_id": "str_c", "x": 0.1, "y": 0.1,
	}))
	waitForType(t, viewer.send, MsgCanvasStrokeStarted, 2*time.Second)

	owner.handleMessage(newMessage(MsgCanvasClear, "canvas-4", owner.UserID, owner.Username, map[string]interface{}{}))
	waitForType(t, viewer.send, MsgCanvasCleared, 2*time.Second)

	if n := len(hub.GetOrCreateRoom("canvas-4").GetCanvasStrokes()); n != 0 {
		t.Fatalf("setelah clear, riwayat harus kosong, dapat %d", n)
	}
}

// TestCanvasRingBufferCap - server hanya menyimpan 50 stroke terakhir.
func TestCanvasRingBufferCap(t *testing.T) {
	hub, owner, _ := setupCanvasRoom(t, "canvas-5")
	defer hub.Stop()

	for i := 0; i < 60; i++ {
		owner.handleMessage(newMessage(MsgCanvasDrawStart, "canvas-5", owner.UserID, owner.Username, map[string]interface{}{
			"stroke_id": fmt.Sprintf("str_%d", i), "x": 0.1, "y": 0.1,
		}))
	}
	strokes := hub.GetOrCreateRoom("canvas-5").GetCanvasStrokes()
	if len(strokes) != 50 {
		t.Fatalf("ring buffer harus 50, dapat %d", len(strokes))
	}
	if strokes[0].StrokeID != "str_10" {
		t.Fatalf("stroke tertua yang tersisa harus str_10, dapat %s", strokes[0].StrokeID)
	}
}

// TestCanvasLateJoinerStateSync - peserta yang baru join menerima coretan yang sudah ada via state_sync.
func TestCanvasLateJoinerStateSync(t *testing.T) {
	hub, owner, _ := setupCanvasRoom(t, "canvas-6")
	defer hub.Stop()

	owner.handleMessage(newMessage(MsgCanvasDrawStart, "canvas-6", owner.UserID, owner.Username, map[string]interface{}{
		"stroke_id": "str_old", "x": 0.2, "y": 0.2,
	}))

	late := newTestClient(hub, "canvas-6", "user-late", "Telat")
	late.SetRoomRole(hub.ResolveRoomRole("canvas-6", "user-late", ""))
	hub.registerClient(late)

	sync := waitForType(t, late.send, MsgStateSync, 2*time.Second)
	strokes, _ := sync.Payload["canvas_strokes"].([]*CanvasStrokeItem)
	if len(strokes) != 1 || strokes[0].StrokeID != "str_old" {
		t.Fatalf("state_sync harus berisi str_old, dapat %+v", sync.Payload["canvas_strokes"])
	}
}
