package auth

import (
	"atlas_food/internal/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler - struct untuk HTTP handler auth
type Handler struct {
	service Service
}

// NewHandler - factory function untuk membuat handler
// db: koneksi database GORM
func NewHandler(db *gorm.DB) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return &Handler{service: service}
}

// Register - handler untuk endpoint POST /auth/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	response, err := h.service.Register(req)
	if err != nil {
		if err.Error() == "email sudah terdaftar" {
			utils.ErrorResponse(c, http.StatusConflict, "CONFLICT", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.CreatedResponse(c, response)
}

// Login - handler untuk endpoint POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	response, err := h.service.Login(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// RefreshToken - handler untuk endpoint POST /auth/refresh
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Refresh token diperlukan")
		return
	}

	response, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// GetProfile - handler untuk endpoint GET /auth/me (protected)
func (h *Handler) GetProfile(c *gin.Context) {
	// Ambil userID dari context (set oleh middleware JWT)
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User tidak terautentikasi")
		return
	}

	profile, err := h.service.GetProfile(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	utils.SuccessResponse(c, profile)
}
