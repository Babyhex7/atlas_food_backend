package survey

import (
	"atlas_food/internal/pkg/middleware"
	"atlas_food/internal/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler - HTTP handler survey
type Handler struct {
	service Service
}

// NewHandler - factory function
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ============ ADMIN ENDPOINTS ============

// CreateSurvey - POST /api/v1/admin/surveys
func (h *Handler) CreateSurvey(c *gin.Context) {
	var req CreateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

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
	id := c.Param("id")

	var req UpdateSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	response, err := h.service.UpdateSurvey(id, req)
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

// AccessSurvey - POST /api/v1/survey/access
func (h *Handler) AccessSurvey(c *gin.Context) {
	var req AccessSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	userIDRaw, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User tidak terautentikasi")
		return
	}

	userID := userIDRaw.(string)

	response, err := h.service.AccessSurvey(req, &userID)
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

// GetSurveyInfo - GET /api/v1/survey/:id/info
func (h *Handler) GetSurveyInfo(c *gin.Context) {
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

	// Return limited info for respondent
	utils.SuccessResponse(c, response)
}

// ListActiveSurveys - GET /api/v1/survey/active
// Menampilkan semua kuesioner aktif yang bisa diisi oleh responden
func (h *Handler) ListActiveSurveys(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	response, err := h.service.ListActiveSurveys(page, limit)
	if err != nil {
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

	// Respondent routes
	resp := router.Group("/survey", authMiddleware, middleware.RespondentOnly())
	{
		resp.GET("/active", h.ListActiveSurveys)
		resp.POST("/access", h.AccessSurvey)
		resp.GET("/:id/info", h.GetSurveyInfo)
	}

	// Locales (public)
	router.GET("/locales", h.ListLocales)
}
