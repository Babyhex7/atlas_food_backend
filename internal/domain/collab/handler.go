package collab

import (
	"log"
	"net/http"
	"strings"
	"time"

	"atlas_food/internal/pkg/utils"

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
	fromInvite := false
	if inviteTok != "" {
		inv, ok := h.hub.Invites().Get(inviteTok)
		if !ok || inv.RoomID != roomID {
			utils.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN",
				"Invite tidak valid atau sudah kedaluwarsa")
			return
		}
		fromInvite = true
	}

	roomRole := h.hub.ResolveRoomRole(roomID, userID.(string), inviteTok)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := NewClient(conn, h.hub, roomID, userID.(string), username.(string), roleStr, roomRole)
	client.RoleFromInvite = fromInvite
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
//
// Izin membagikan mengikuti aturan yang sama dengan Figma: hanya anggota room
// dengan hak ubah (owner/editor) yang boleh mengundang, dan tidak seorang pun
// boleh memberi role di atas miliknya sendiri. Tanpa penjagaan ini, viewer bisa
// mencetak undangan editor — menaikkan hak orang lain melebihi haknya sendiri —
// dan orang di luar room bisa membuat undangan untuk room mana pun asal tahu
// id-nya.
func (h *Handler) InviteToRoom(c *gin.Context) {
	roomID := c.Param("room_id")
	if roomID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "room_id wajib diisi")
		return
	}

	userID, _ := c.Get("userID")
	createdBy, _ := userID.(string)

	callerRole := h.hub.RoomRoleOf(roomID, createdBy)
	if callerRole == "" {
		utils.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN",
			"Anda bukan anggota sesi ini")
		return
	}
	if callerRole != RoomRoleOwner && callerRole != RoomRoleEditor {
		utils.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN",
			"Mode hanya lihat tidak bisa membagikan sesi. Minta owner untuk mengundang.")
		return
	}

	var req inviteRequest
	_ = c.ShouldBindJSON(&req)
	role := strings.TrimSpace(req.Role)
	if role == "" {
		role = RoomRoleEditor
	}
	// Owner tidak bisa dibagikan lewat link: pemindahan kepemilikan adalah
	// tindakan tersendiri, bukan efek samping menyalin URL.
	if role != RoomRoleEditor && role != RoomRoleViewer {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR",
			"role harus editor atau viewer")
		return
	}

	inv := h.hub.Invites().Create(roomID, role, createdBy, 24*time.Hour)

	utils.SuccessResponse(c, gin.H{
		"room_id":      roomID,
		"invite_token": inv.Token,
		"role":         inv.Role,
		"join_path":    "?room=" + roomID + "&invite=" + inv.Token,
		"expires_at":   inv.ExpiresAt.UTC(),
		"note":         "Bagikan URL dengan ?room= & ?invite=. Penerima harus login. Role: " + inv.Role,
	})
}

// RevokeInvite membatalkan invite token.
//
// Hanya pembuat undangan atau owner room yang boleh mencabut; sebelumnya siapa
// pun yang tahu tokennya bisa membatalkan undangan orang lain.
func (h *Handler) RevokeInvite(c *gin.Context) {
	token := strings.TrimSpace(c.Param("token"))
	if token == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "token wajib diisi")
		return
	}

	inv, ok := h.hub.Invites().Get(token)
	if !ok {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND",
			"Undangan tidak ditemukan atau sudah kedaluwarsa")
		return
	}

	userID, _ := c.Get("userID")
	callerID, _ := userID.(string)
	if callerID != inv.CreatedBy && h.hub.RoomRoleOf(inv.RoomID, callerID) != RoomRoleOwner {
		utils.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN",
			"Hanya pembuat undangan atau owner sesi yang bisa mencabutnya")
		return
	}

	h.hub.Invites().Revoke(token)
	utils.SuccessResponse(c, gin.H{"revoked": true, "token": token})
}
