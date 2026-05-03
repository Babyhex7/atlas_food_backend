package survey

import (
	"atlas_food/internal/pkg/utils"
	"net/http"
	"atlas_food/internal/pkg/middleware"
	"atlas_food/internal/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler - HTTP handler survey
>>>>>>> c72d9d4 (add survey be schema and logic setup)

<<<<<<< HEAD
=======
// NewHandler - factory function
	repo := NewRepository(db)
	return &Handler{service: service}
}

<<<<<<< HEAD
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
=======
// ============ ADMIN ENDPOINTS ============

// CreateSurvey - POST /api/v1/admin/surveys
func (h *Handler) CreateSurvey(c *gin.Context) {
>>>>>>> c72d9d4 (add survey be schema and logic setup)
	var req CreateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

<<<<<<< HEAD
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
=======
	// Ambil userID dari context (JWT)
	userID, _ := c.Get("userID")

	response, err := h.service.CreateSurvey(req, userID.(string))
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.CreatedResponse(c, response)
}

// ListSurveys - GET /api/v1/admin/surveys
func (h *Handler) ListSurveys(c *gin.Context) {
	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Ambil userID dari context
	userID, _ := c.Get("userID")

	response, err := h.service.ListSurveys(userID.(string), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// GetSurvey - GET /api/v1/admin/surveys/:id
func (h *Handler) GetSurvey(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.GetSurveyByID(id)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// UpdateSurvey - PUT /api/v1/admin/surveys/:id
func (h *Handler) UpdateSurvey(c *gin.Context) {
>>>>>>> c72d9d4 (add survey be schema and logic setup)
	id := c.Param("id")

	var req UpdateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

<<<<<<< HEAD
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
=======
	response, err := h.service.UpdateSurvey(id, req)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
>>>>>>> c72d9d4 (add survey be schema and logic setup)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

<<<<<<< HEAD
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
=======
	utils.SuccessResponse(c, response)
}

// DeleteSurvey - DELETE /api/v1/admin/surveys/:id
func (h *Handler) DeleteSurvey(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteSurvey(id); err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Survey berhasil dihapus"})
}

// CloneSurvey - POST /api/v1/admin/surveys/:id/clone
func (h *Handler) CloneSurvey(c *gin.Context) {
	id := c.Param("id")

	var req CloneSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	// Ambil userID dari context
	userID, _ := c.Get("userID")

	response, err := h.service.CloneSurvey(id, req, userID.(string))
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.CreatedResponse(c, response)
}

// RegenerateAccessToken - POST /api/v1/admin/surveys/:id/regenerate-token
func (h *Handler) RegenerateAccessToken(c *gin.Context) {
	surveyID := c.Param("id")

	response, err := h.service.GenerateAccessToken(surveyID)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// ============ PUBLIC/RESPONDENT ENDPOINTS ============

// GetPublicSurvey - GET /api/v1/s/:token (public access)
func (h *Handler) GetPublicSurvey(c *gin.Context) {
	token := c.Param("token")

	response, err := h.service.GetPublicSurveyByToken(token)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// JoinSurvey - POST /api/v1/s/:token/join
func (h *Handler) JoinSurvey(c *gin.Context) {
	token := c.Param("token")

	var req JoinSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	// Cek apakah user sudah login (optional)
	var userID *string
	if id, exists := c.Get("userID"); exists {
		uid := id.(string)
		userID = &uid
	}

	response, err := h.service.JoinSurvey(token, userID, req)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// ============ UTILITY ENDPOINTS ============

// ListLocales - GET /api/v1/locales
func (h *Handler) ListLocales(c *gin.Context) {
	locales, err := h.service.GetAllLocales()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, locales)
}

// SetupRoutes - daftarkan semua route survey ke router
func (h *Handler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// Admin routes (protected)
	admin := router.Group("/admin/surveys", authMiddleware, middleware.AdminOnly())
	{
		admin.POST("", h.CreateSurvey)
		admin.GET("", h.ListSurveys)
		admin.GET("/:id", h.GetSurvey)
		admin.PUT("/:id", h.UpdateSurvey)
		admin.DELETE("/:id", h.DeleteSurvey)
		admin.POST("/:id/clone", h.CloneSurvey)
		admin.POST("/:id/regenerate-token", h.RegenerateAccessToken)
	}

	// Public routes (survey access)
	public := router.Group("/s")
	{
		public.GET("/:token", h.GetPublicSurvey)
		public.POST("/:token/join", h.JoinSurvey)
	}

	// Locales (public)
	router.GET("/locales", h.ListLocales)
}
>>>>>>> c72d9d4 (add survey be schema and logic setup)
