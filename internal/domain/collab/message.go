package collab

import "time"

// Message represents a WebSocket message
type Message struct {
	Type      string                 `json:"type"`
	RoomID    string                 `json:"room_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Username  string                 `json:"username,omitempty"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// Client → Server message types
const (
	MsgPresenceJoin   = "presence_join"
	MsgPresenceLeave  = "presence_leave"
	MsgCursorMove     = "cursor_move"
	MsgViewportUpdate = "viewport_update" // Figma-like: page + scroll + step
	MsgFollowUser     = "follow_user"
	MsgUnfollowUser   = "unfollow_user"
	MsgFoodSearch     = "food_search"
	MsgFoodSelect     = "food_select"
	MsgMealAdd        = "meal_add"
	MsgPortionSet     = "portion_set"
	MsgPortionSelect  = "portion_select" // alias
	MsgReviewSubmit   = "review_submit"
	MsgDBEditStart    = "db_edit_start"
	MsgDBEditField    = "db_edit_field"
	MsgDBEditSave     = "db_edit_save"
	MsgDBEditCancel   = "db_edit_cancel"
	MsgChatMessage    = "chat_message"
	MsgGetHistory     = "get_history"
	MsgPing           = "ping"

	// Cursor chat (ala Figma "/"): bubble teks ephemeral yang nempel di kursor,
	// terpisah dari MsgChatMessage (panel chat biasa) — lihat
	// docs/superpowers/specs/2026-08-16-cursor-chat-prd.md
	MsgCursorChatOpen   = "cursor_chat_open"
	MsgCursorChatUpdate = "cursor_chat_update"
	MsgCursorChatClose  = "cursor_chat_close"

	// Live Canvas Annotation (Coret-coret layar real-time)
	MsgCanvasDrawStart = "canvas_draw_start"
	MsgCanvasDrawMove  = "canvas_draw_move"
	MsgCanvasDrawEnd   = "canvas_draw_end"
	MsgCanvasLaserMove = "canvas_laser_move"
	MsgCanvasClear     = "canvas_clear"

	// Live Dynamic Role Control (Owner ubah role peserta)
	MsgUpdateUserRole = "update_user_role"
)

// cursorChatMaxTextLen - batas panjang teks bubble; cegah payload nakal membanjiri broadcast
const cursorChatMaxTextLen = 200

// Room collaboration roles (per-room, bukan JWT role aplikasi)
const (
	RoomRoleOwner  = "owner"
	RoomRoleEditor = "editor"
	RoomRoleViewer = "viewer"
)

// Server → Client message types
const (
	MsgPresenceList      = "presence_list"
	MsgPresenceJoined    = "presence_joined"
	MsgPresenceLeft      = "presence_left"
	MsgUserJoined        = "user_joined" // legacy alias
	MsgUserLeft          = "user_left"   // legacy alias
	MsgCursorUpdate      = "cursor_update"
	MsgViewportSync      = "viewport_sync" // mirror ke follower
	MsgFollowStarted     = "follow_started"
	MsgFollowStopped     = "follow_stopped"
	MsgFollowState       = "follow_state" // snapshot follow graph
	MsgUserSearching     = "user_searching"
	MsgFoodSearchShared  = "food_search_shared"
	MsgFoodSelected      = "food_selected"
	MsgMealUpdated       = "meal_updated"
	MsgPortionUpdated    = "portion_updated"
	MsgPortionSelected   = "portion_selected" // legacy alias
	MsgReviewSubmitted   = "review_submitted"
	MsgDBLocked          = "db_locked"
	MsgDBFieldUpdated    = "db_field_updated"
	MsgDBEditSaved       = "db_edit_saved"
	MsgDBUnlocked        = "db_unlocked"
	MsgActivityLog       = "activity_log"
	MsgError             = "error"
	MsgPong              = "pong"
	MsgHistory           = "history"
	MsgStateSync         = "state_sync"
	MsgCursorChatUpdated = "cursor_chat_updated"
	MsgCursorChatClosed  = "cursor_chat_closed"
	MsgUserRoleUpdated   = "user_role_updated"

	// Server -> Client Canvas Annotation
	MsgCanvasStrokeStarted = "canvas_stroke_started"
	MsgCanvasStrokeUpdated = "canvas_stroke_updated"
	MsgCanvasStrokeEnded   = "canvas_stroke_ended"
	MsgCanvasLaserUpdated = "canvas_laser_updated"
	MsgCanvasCleared      = "canvas_cleared"
	MsgCanvasStateSync    = "canvas_state_sync"
)

// newMessage - helper pembuat Message dengan timestamp sekarang dan payload non-nil
func newMessage(msgType, roomID, userID, username string, payload map[string]interface{}) *Message {
	if payload == nil {
		payload = map[string]interface{}{}
	}
	return &Message{
		Type:      msgType,
		RoomID:    roomID,
		UserID:    userID,
		Username:  username,
		Payload:   payload,
		Timestamp: time.Now(),
	}
}

// payloadString - ambil nilai string dari payload secara aman; "" kalau tidak ada / bukan string
func payloadString(payload map[string]interface{}, key string) string {
	if payload == nil {
		return ""
	}
	if v, ok := payload[key].(string); ok {
		return v
	}
	return ""
}

// payloadInt - ambil nilai int dari payload secara aman (JSON number masuk sebagai float64); 0 kalau gagal
func payloadInt(payload map[string]interface{}, key string) int {
	if payload == nil {
		return 0
	}
	switch v := payload[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	default:
		return 0
	}
}

// payloadFloat64 - ambil nilai float64 dari payload secara aman; defaultVal kalau gagal
func payloadFloat64(payload map[string]interface{}, key string, defaultVal float64) float64 {
	if payload == nil {
		return defaultVal
	}
	switch v := payload[key].(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return defaultVal
	}
}
