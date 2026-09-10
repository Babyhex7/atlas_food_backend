package gamification

import (
	"atlas_food/internal/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler - HTTP handler untuk gamification
type Handler struct {
	service Service
}

// NewHandler creates a new gamification handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetProfile - GET /api/v1/gamification/profile
func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User tidak terautentikasi")
		return
	}

	profile, err := h.service.GetProfile(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengambil profil gamifikasi")
		return
	}

	// Inject info dari JWT kalau mau
	if name, ok := c.Get("name"); ok {
		profile.Name = name.(string)
	}

	utils.SuccessResponse(c, profile)
}

// GetLeaderboard - GET /api/v1/gamification/leaderboard
func (h *Handler) GetLeaderboard(c *gin.Context) {
	userID, _ := c.Get("userID")
	period := c.DefaultQuery("period", "all_time")

	leaderboard, err := h.service.GetLeaderboard(period, 10, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengambil leaderboard")
		return
	}

	utils.SuccessResponse(c, leaderboard)
}

// GetBadges - GET /api/v1/gamification/badges
func (h *Handler) GetBadges(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User tidak terautentikasi")
		return
	}

	badges, err := h.service.GetBadges(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengambil badge")
		return
	}

	utils.SuccessResponse(c, gin.H{"badges": badges})
}

// GetXPHistory - GET /api/v1/gamification/xp-history
func (h *Handler) GetXPHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User tidak terautentikasi")
		return
	}

	history, err := h.service.GetXPHistory(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengambil riwayat XP")
		return
	}

	utils.SuccessResponse(c, gin.H{"history": history})
}

// SetupRoutes registers gamification endpoints
func (h *Handler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	gamification := router.Group("/gamification", authMiddleware)
	{
		gamification.GET("/profile", h.GetProfile)
		gamification.GET("/leaderboard", h.GetLeaderboard)
		gamification.GET("/badges", h.GetBadges)
		gamification.GET("/xp-history", h.GetXPHistory)
	}
}
