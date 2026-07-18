package collab

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development — restrict in production
		return true
	},
}

// Handler handles WebSocket connections
type Handler struct {
	hub *Hub
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

// HandleWebSocket upgrades HTTP to WebSocket for a collaboration room.
// @Summary WebSocket connection for real-time collaboration
// @Tags collaboration
// @Param room_id path string true "Room ID"
// @Param token query string false "JWT access token (required for browser WS)"
// @Security BearerAuth
// @Router /collab/rooms/{room_id}/ws [get]
func (h *Handler) HandleWebSocket(c *gin.Context) {
	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "room_id is required",
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		userID, exists = c.Get("user_id")
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Unauthorized",
		})
		return
	}

	username, _ := c.Get("username")
	if username == nil || username.(string) == "" {
		if email, ok := c.Get("email"); ok && email != nil {
			username = email
		} else {
			username = "Anonymous"
		}
	}

	role, _ := c.Get("role")
	roleStr := ""
	if role != nil {
		roleStr, _ = role.(string)
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := NewClient(conn, h.hub, roomID, userID.(string), username.(string), roleStr)
	h.hub.register <- client

	go client.WritePump()
	go client.ReadPump()

	log.Printf("✅ WebSocket connected: user=%s, room=%s", userID, roomID)
}

// GetRoomInfo returns information about a room
func (h *Handler) GetRoomInfo(c *gin.Context) {
	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "room_id is required",
		})
		return
	}

	info := h.hub.GetRoomInfo(roomID)
	if info == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Room not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   info,
	})
}

// GetHubStats returns hub statistics
func (h *Handler) GetHubStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   h.hub.GetStats(),
	})
}

// InviteToRoom returns a shareable room join hint (in-memory; no Redis).
func (h *Handler) InviteToRoom(c *gin.Context) {
	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "room_id is required",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"room_id":    roomID,
			"join_path":  "?room=" + roomID,
			"expires_at": time.Now().Add(24 * time.Hour).UTC(),
			"note":       "Bagikan URL halaman dengan query ?room=" + roomID + " ke kolaborator (login required).",
		},
	})
}
