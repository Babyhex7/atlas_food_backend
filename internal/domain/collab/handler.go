package collab

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// TODO: Restrict in production
		return true
	},
}

// Handler handles WebSocket connections
type Handler struct {
	hub *Hub
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

// HandleWebSocket handles WebSocket connection requests
// @Summary WebSocket connection for real-time collaboration
// @Description Connect to a room for real-time food search collaboration
// @Tags collaboration
// @Param room_id path string true "Room ID"
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

	// Get user info from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Unauthorized",
		})
		return
	}

	username, _ := c.Get("username")
	if username == nil {
		username = "Anonymous"
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Create new client
	client := NewClient(conn, h.hub, roomID, userID.(string), username.(string))

	// Register client
	h.hub.register <- client

	// Start read and write pumps in goroutines
	go client.WritePump()
	go client.ReadPump()

	log.Printf("✅ WebSocket connected: user=%s, room=%s", userID, roomID)
}

// GetRoomInfo returns information about a room
// @Summary Get room information
// @Description Get information about active users in a room
// @Tags collaboration
// @Param room_id path string true "Room ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /collab/rooms/{room_id} [get]
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
// @Summary Get hub statistics
// @Description Get statistics about all rooms and clients
// @Tags collaboration
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /collab/stats [get]
func (h *Handler) GetHubStats(c *gin.Context) {
	stats := h.hub.GetStats()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   stats,
	})
}
