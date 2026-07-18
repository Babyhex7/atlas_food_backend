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
	MsgPresenceJoin  = "presence_join"
	MsgPresenceLeave = "presence_leave"
	MsgCursorMove    = "cursor_move"
	MsgFoodSearch    = "food_search"
	MsgFoodSelect    = "food_select"
	MsgMealAdd       = "meal_add"
	MsgPortionSet    = "portion_set"
	MsgPortionSelect = "portion_select" // alias
	MsgReviewSubmit  = "review_submit"
	MsgDBEditStart   = "db_edit_start"
	MsgDBEditField   = "db_edit_field"
	MsgDBEditSave    = "db_edit_save"
	MsgDBEditCancel  = "db_edit_cancel"
	MsgChatMessage   = "chat_message"
	MsgGetHistory    = "get_history"
	MsgPing          = "ping"
)

// Server → Client message types
const (
	MsgPresenceList     = "presence_list"
	MsgPresenceJoined   = "presence_joined"
	MsgPresenceLeft     = "presence_left"
	MsgUserJoined       = "user_joined" // legacy alias
	MsgUserLeft         = "user_left"   // legacy alias
	MsgCursorUpdate     = "cursor_update"
	MsgUserSearching    = "user_searching"
	MsgFoodSearchShared = "food_search_shared"
	MsgFoodSelected     = "food_selected"
	MsgMealUpdated      = "meal_updated"
	MsgPortionUpdated   = "portion_updated"
	MsgPortionSelected  = "portion_selected" // legacy alias
	MsgReviewSubmitted  = "review_submitted"
	MsgDBLocked         = "db_locked"
	MsgDBFieldUpdated   = "db_field_updated"
	MsgDBEditSaved      = "db_edit_saved"
	MsgDBUnlocked       = "db_unlocked"
	MsgActivityLog      = "activity_log"
	MsgError            = "error"
	MsgPong             = "pong"
	MsgHistory          = "history"
	MsgStateSync        = "state_sync"
)

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

func payloadString(payload map[string]interface{}, key string) string {
	if payload == nil {
		return ""
	}
	if v, ok := payload[key].(string); ok {
		return v
	}
	return ""
}

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
