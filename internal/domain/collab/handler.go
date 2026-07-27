package collab

import (
	"log"
	"net/http"
	"strings"
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
// Query: token (JWT), invite (optional invite token for room role).
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
			"message": "Unauthorized — login diperlukan",
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

	inviteTok := strings.TrimSpace(c.Query("invite"))
	// Validasi invite bila ada: harus match room
	if inviteTok != "" {
		inv, ok := h.hub.Invites().Get(inviteTok)
		if !ok || inv.RoomID != roomID {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Invite tidak valid atau sudah kedaluwarsa",
			})
			return
		}
	}

	roomRole := h.hub.ResolveRoomRole(roomID, userID.(string), inviteTok)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := NewClient(conn, h.hub, roomID, userID.(string), username.(string), roleStr, roomRole)
	h.hub.register <- client

	go client.WritePump()
	go client.ReadPump()

	log.Printf("✅ WebSocket connected: user=%s, room=%s, room_role=%s", userID, roomID, roomRole)
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

type inviteRequest struct {
	Role string `json:"role"` // editor | viewer
}

// InviteToRoom membuat invite token berbatas waktu + URL share.
func (h *Handler) InviteToRoom(c *gin.Context) {
	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "room_id is required",
		})
		return
	}

	var req inviteRequest
	_ = c.ShouldBindJSON(&req)
	role := strings.TrimSpace(req.Role)
	if role == "" {
		role = RoomRoleEditor
	}
	if role != RoomRoleEditor && role != RoomRoleViewer {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "role harus editor atau viewer",
		})
		return
	}

	userID, _ := c.Get("userID")
	createdBy, _ := userID.(string)

	inv := h.hub.Invites().Create(roomID, role, createdBy, 24*time.Hour)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"room_id":      roomID,
			"invite_token": inv.Token,
			"role":         inv.Role,
			"join_path":    "?room=" + roomID + "&invite=" + inv.Token,
			"expires_at":   inv.ExpiresAt.UTC(),
			"note":         "Bagikan URL dengan ?room= & ?invite=. Penerima harus login. Role: " + inv.Role,
		},
	})
}

// RevokeInvite membatalkan invite token.
func (h *Handler) RevokeInvite(c *gin.Context) {
	token := strings.TrimSpace(c.Param("token"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "token wajib"})
		return
	}
	h.hub.Invites().Revoke(token)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   gin.H{"revoked": true, "token": token},
	})
}
