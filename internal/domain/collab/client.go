package collab

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 64 * 1024 // 64KB
	sendBufferSize = 256
)

// Client represents a WebSocket client
type Client struct {
	conn            *websocket.Conn
	hub             *Hub
	RoomID          string
	UserID          string
	Username        string
	Role            string
	send            chan *Message
	lastMessageTime time.Time
	messageCount    int
	stopCh          chan struct{}
}

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn, hub *Hub, roomID, userID, username, role string) *Client {
	return &Client{
		conn:     conn,
		hub:      hub,
		RoomID:   roomID,
		UserID:   userID,
		Username: username,
		Role:     role,
		send:     make(chan *Message, sendBufferSize),
		stopCh:   make(chan struct{}),
	}
}

func (c *Client) sendQuiet(msg *Message) {
	select {
	case c.send <- msg:
	default:
		log.Printf("Failed to send to client %s", c.UserID)
	}
}

func (c *Client) sendError(code, message string) {
	c.sendQuiet(newMessage(MsgError, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"code":    code,
		"message": message,
	}))
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-c.stopCh:
			return
		default:
			_, messageBytes, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}

			if !c.checkRateLimit() {
				c.sendError("RATE_LIMITED", "Terlalu banyak pesan — tunggu sebentar")
				continue
			}

			var rawMsg map[string]interface{}
			if err := json.Unmarshal(messageBytes, &rawMsg); err != nil {
				c.sendError("INVALID_JSON", "Format pesan tidak valid")
				continue
			}

			msgType, _ := rawMsg["type"].(string)
			payload, _ := rawMsg["payload"].(map[string]interface{})
			msg := newMessage(msgType, c.RoomID, c.UserID, c.Username, payload)
			c.handleMessage(msg)
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			messageJSON, err := json.Marshal(message)
			if err != nil {
				continue
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.stopCh:
			return
		}
	}
}

func (c *Client) handleMessage(msg *Message) {
	switch msg.Type {
	case MsgPresenceJoin:
		// Already registered on connect; re-broadcast presence list
		room := c.hub.GetOrCreateRoom(c.RoomID)
		c.sendQuiet(c.hub.buildPresenceList(room))

	case MsgPresenceLeave:
		c.hub.unregister <- c

	case MsgCursorMove:
		c.handleCursorMove(msg)

	case MsgFoodSearch:
		c.handleFoodSearch(msg)

	case MsgFoodSelect:
		c.handleFoodSelect(msg)

	case MsgMealAdd:
		c.handleMealAdd(msg)

	case MsgPortionSet, MsgPortionSelect:
		c.handlePortionSet(msg)

	case MsgReviewSubmit:
		c.handleReviewSubmit(msg)

	case MsgDBEditStart:
		c.handleDBEditStart(msg)

	case MsgDBEditField:
		c.handleDBEditField(msg)

	case MsgDBEditSave:
		c.handleDBEditSave(msg)

	case MsgDBEditCancel:
		c.handleDBEditCancel(msg)

	case MsgChatMessage:
		c.handleChatMessage(msg)

	case MsgGetHistory:
		c.sendRoomHistory()

	case MsgPing:
		c.sendQuiet(newMessage(MsgPong, c.RoomID, c.UserID, c.Username, map[string]interface{}{}))

	default:
		log.Printf("Unknown message type: %s", msg.Type)
		c.sendError("UNKNOWN_TYPE", "Tipe pesan tidak dikenal: "+msg.Type)
	}
}

func (c *Client) handleCursorMove(msg *Message) {
	payload := msg.Payload
	if payload == nil {
		payload = map[string]interface{}{}
	}
	payload["color"] = colorForUser(c.UserID)
	payload["page"] = payloadString(payload, "page")
	update := newMessage(MsgCursorUpdate, c.RoomID, c.UserID, c.Username, payload)
	room := c.hub.GetOrCreateRoom(c.RoomID)
	room.AddMessage(update)
}

func (c *Client) handleFoodSearch(msg *Message) {
	query := payloadString(msg.Payload, "query")
	c.hub.Publish(newMessage(MsgUserSearching, c.RoomID, c.UserID, c.Username, msg.Payload))
	c.hub.Publish(newMessage(MsgFoodSearchShared, c.RoomID, c.UserID, c.Username, msg.Payload))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "food_search",
		"details": c.Username + ` mencari "` + query + `"`,
		"query":   query,
	}))
}

func (c *Client) handleFoodSelect(msg *Message) {
	foodName := payloadString(msg.Payload, "food_name")
	c.hub.Publish(newMessage(MsgFoodSelected, c.RoomID, c.UserID, c.Username, msg.Payload))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "food_select",
		"details": c.Username + " memilih " + foodName,
	}))
}

func (c *Client) handleMealAdd(msg *Message) {
	foodName := payloadString(msg.Payload, "food_name")
	mealType := payloadString(msg.Payload, "meal_type")
	c.hub.Publish(newMessage(MsgMealUpdated, c.RoomID, c.UserID, c.Username, msg.Payload))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "meal_add",
		"details": c.Username + " menambah " + foodName + " ke " + mealType,
	}))
}

func (c *Client) handlePortionSet(msg *Message) {
	c.hub.Publish(newMessage(MsgPortionUpdated, c.RoomID, c.UserID, c.Username, msg.Payload))
	c.hub.Publish(newMessage(MsgPortionSelected, c.RoomID, c.UserID, c.Username, msg.Payload))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "portion_set",
		"details": c.Username + " mengatur porsi",
	}))
}

func (c *Client) handleReviewSubmit(msg *Message) {
	c.hub.Publish(newMessage(MsgReviewSubmitted, c.RoomID, c.UserID, c.Username, msg.Payload))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "review_submit",
		"details": c.Username + " mengirim review/submit survey",
	}))
}

func (c *Client) handleDBEditStart(msg *Message) {
	entityType := payloadString(msg.Payload, "entity_type")
	entityID := payloadString(msg.Payload, "entity_id")
	version := payloadInt(msg.Payload, "version")
	if entityType == "" || entityID == "" {
		c.sendError("INVALID_PAYLOAD", "entity_type dan entity_id wajib")
		return
	}

	lock, ok, existing := c.hub.Locks().TryLock(entityType, entityID, c.UserID, c.Username, version)
	if !ok {
		c.sendError("LOCKED", "Sedang diedit oleh "+existing.Username)
		c.sendQuiet(newMessage(MsgDBLocked, c.RoomID, existing.LockedBy, existing.Username, map[string]interface{}{
			"entity_type": existing.EntityType,
			"entity_id":   existing.EntityID,
			"locked_by":   existing.LockedBy,
			"username":    existing.Username,
			"version":     existing.Version,
		}))
		return
	}

	payload := map[string]interface{}{
		"entity_type": lock.EntityType,
		"entity_id":   lock.EntityID,
		"locked_by":   lock.LockedBy,
		"username":    lock.Username,
		"version":     lock.Version,
	}
	c.hub.Publish(newMessage(MsgDBLocked, c.RoomID, c.UserID, c.Username, payload))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "db_edit_start",
		"details": c.Username + " mulai edit " + entityType + " " + entityID,
	}))
}

func (c *Client) handleDBEditField(msg *Message) {
	entityType := payloadString(msg.Payload, "entity_type")
	entityID := payloadString(msg.Payload, "entity_id")
	lock := c.hub.Locks().Get(entityType, entityID)
	if lock == nil || lock.LockedBy != c.UserID {
		c.sendError("NOT_LOCK_OWNER", "Anda tidak memegang lock entity ini")
		return
	}
	c.hub.Publish(newMessage(MsgDBFieldUpdated, c.RoomID, c.UserID, c.Username, msg.Payload))
}

func (c *Client) handleDBEditSave(msg *Message) {
	entityType := payloadString(msg.Payload, "entity_type")
	entityID := payloadString(msg.Payload, "entity_id")
	clientVersion := payloadInt(msg.Payload, "version")

	lock := c.hub.Locks().Get(entityType, entityID)
	if lock == nil || lock.LockedBy != c.UserID {
		c.sendError("NOT_LOCK_OWNER", "Anda tidak memegang lock entity ini")
		return
	}
	if clientVersion > 0 && clientVersion < lock.Version {
		c.sendError("VERSION_CONFLICT", "Versi data sudah berubah — refresh lalu coba lagi")
		return
	}

	bumped, ok := c.hub.Locks().BumpVersion(entityType, entityID, c.UserID)
	if !ok {
		c.sendError("VERSION_CONFLICT", "Gagal menyimpan versi")
		return
	}

	c.hub.Locks().Release(entityType, entityID, c.UserID)
	c.hub.Publish(newMessage(MsgDBEditSaved, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"entity_type": entityType,
		"entity_id":   entityID,
		"version":     bumped.Version,
		"changes":     msg.Payload["changes"],
	}))
	c.hub.Publish(newMessage(MsgDBUnlocked, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"entity_type": entityType,
		"entity_id":   entityID,
	}))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "db_edit_save",
		"details": c.Username + " menyimpan " + entityType + " " + entityID,
	}))
}

func (c *Client) handleDBEditCancel(msg *Message) {
	entityType := payloadString(msg.Payload, "entity_type")
	entityID := payloadString(msg.Payload, "entity_id")
	if !c.hub.Locks().Release(entityType, entityID, c.UserID) {
		c.sendError("NOT_LOCK_OWNER", "Tidak bisa melepas lock")
		return
	}
	c.hub.Publish(newMessage(MsgDBUnlocked, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"entity_type": entityType,
		"entity_id":   entityID,
	}))
	c.hub.Publish(newMessage(MsgActivityLog, c.RoomID, c.UserID, c.Username, map[string]interface{}{
		"action":  "db_edit_cancel",
		"details": c.Username + " membatalkan edit",
	}))
}

func (c *Client) handleChatMessage(msg *Message) {
	c.hub.Publish(newMessage(MsgChatMessage, c.RoomID, c.UserID, c.Username, msg.Payload))
}

func (c *Client) sendRoomHistory() {
	room := c.hub.GetOrCreateRoom(c.RoomID)
	c.sendQuiet(newMessage(MsgHistory, c.RoomID, "", "", map[string]interface{}{
		"messages": room.GetHistory(50),
	}))
}

func (c *Client) checkRateLimit() bool {
	now := time.Now()
	if now.Sub(c.lastMessageTime) > time.Second {
		c.messageCount = 0
		c.lastMessageTime = now
	}
	c.messageCount++
	return c.messageCount <= 50
}

// Stop stops the client
func (c *Client) Stop() {
	close(c.stopCh)
}
