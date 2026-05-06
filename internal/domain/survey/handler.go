package survey

import (
	"atlas_food/internal/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler - HTTP handler survey
type Handler struct {
	service Service
}

func NewHandler(db *gorm.DB) *Handler {
	repo := NewRepository(db)
	service := NewService(repo)
	return &Handler{service: service}
}

// SeedLocales - inisialisasi locale default
func (h *Handler) SeedLocales() error {
	return h.service.SeedLocales()
}

func (h *Handler) List(c *gin.Context) {
	items, err := h.service.List()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	utils.SuccessResponse(c, items)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User tidak terautentikasi")
		return
	}

	item, err := h.service.Create(req, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	utils.CreatedResponse(c, item)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	item, err := h.service.GetByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}
	utils.SuccessResponse(c, item)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	item, err := h.service.Update(id, req)
	if err != nil {
		if err.Error() == "survey tidak ditemukan" {
			utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	utils.SuccessResponse(c, item)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		if err.Error() == "survey tidak ditemukan" {
			utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Survey deleted successfully"})
}

func (h *Handler) Clone(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User tidak terautentikasi")
		return
	}

	item, err := h.service.Clone(id, userID.(string))
	if err != nil {
		if err.Error() == "survey tidak ditemukan" {
			utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	utils.CreatedResponse(c, item)
}

// GetPublic - get survey detail menggunakan accessToken (untuk respondent)
func (h *Handler) GetPublic(c *gin.Context) {
	accessToken, exists := c.Get("accessToken")
	if !exists {
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "accessToken tidak valid")
		return
	}

	item, err := h.service.GetByAccessToken(accessToken.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	utils.SuccessResponse(c, item)
}

// Join - respondent join survey menggunakan accessToken
func (h *Handler) Join(c *gin.Context) {
	accessToken, exists := c.Get("accessToken")
	if !exists {
		utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "accessToken tidak valid")
		return
	}

	// Validasi bahwa survey aktif
	item, err := h.service.GetByAccessToken(accessToken.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "Survey tidak ditemukan atau belum aktif")
		return
	}

	// Return survey info untuk respondent
	response := gin.H{
		"message": "Berhasil join survey",
		"survey":  item,
	}
	utils.SuccessResponse(c, response)
}
