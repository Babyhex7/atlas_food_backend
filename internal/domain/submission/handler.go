package submission

import (
	"atlas_food/internal/pkg/middleware"
	"atlas_food/internal/pkg/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler - HTTP handler untuk submission
type Handler struct {
	service Service
}

// NewHandler - buat instance handler submission
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ============ RESPONDENT ENDPOINTS ============

// SubmitSurvey - POST /api/v1/survey/submit
// Menyimpan hasil survey dari responden
func (h *Handler) SubmitSurvey(c *gin.Context) {
	var req SubmitSurveyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data survey tidak lengkap: "+err.Error())
		return
	}

	// Extract Idempotency-Key dari Header jika ada
	idempotencyKey := c.GetHeader("Idempotency-Key")
	if idempotencyKey != "" && req.LocalID == "" {
		req.LocalID = idempotencyKey
	}

	if email, exists := c.Get("email"); exists {
		if req.RespondentEmail == "" {
			req.RespondentEmail = email.(string)
		}
	}

	userID, _ := c.Get("userID")
	userIDStr, _ := userID.(string)

	response, err := h.service.SubmitSurvey(req, userIDStr)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		if strings.Contains(err.Error(), "wajib") || strings.Contains(err.Error(), "minimal") || strings.Contains(err.Error(), "valid") {
			status = http.StatusBadRequest
			code = "VALIDATION_ERROR"
		}
		utils.ErrorResponse(c, status, code, err.Error())
		return
	}

	utils.CreatedResponse(c, response)
}

// BatchSubmitSurvey - POST /api/v1/survey/sync/batch
// Menangani pengiriman banyak survey submission sekaligus dari antrean offline
func (h *Handler) BatchSubmitSurvey(c *gin.Context) {
	var req BatchSyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Format data batch sync tidak valid: "+err.Error())
		return
	}

	userID, _ := c.Get("userID")
	userIDStr, _ := userID.(string)

	response, err := h.service.BatchSubmitSurveys(req, userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, response)
}

// ============ ADMIN ENDPOINTS ============

// ListSubmissions - GET /api/v1/admin/surveys/:id/submissions
// Melihat semua hasil survey yang masuk
func (h *Handler) ListSubmissions(c *gin.Context) {
	surveyID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	response, total, err := h.service.ListSubmissions(surveyID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"submissions": response,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetSubmissionDetail - GET /api/v1/admin/submissions/:id
// Melihat detail satu hasil survey
func (h *Handler) GetSubmissionDetail(c *gin.Context) {
	id := c.Param("id")

	response, err := h.service.GetSubmissionDetail(id)
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

// ExportSubmissions - GET /api/v1/admin/surveys/:id/export
// Download hasil survey dalam format CSV
func (h *Handler) ExportSubmissions(c *gin.Context) {
	surveyID := c.Param("id")

	data, filename, err := h.service.ExportSubmissionsCSV(surveyID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "text/csv", data)
}

// GetMySubmissions - GET /api/v1/survey/my-submissions
// Melihat semua riwayat survey yang pernah diisi oleh responden
func (h *Handler) GetMySubmissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	userID, _ := c.Get("userID")
	userIDStr, _ := userID.(string)

	email, _ := c.Get("email")
	emailStr, _ := email.(string)

	response, total, err := h.service.GetMySubmissions(userIDStr, emailStr, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"submissions": response,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetMySubmissionDetail - GET /api/v1/survey/my-submissions/:id
// Melihat detail satu hasil survey milik responden
func (h *Handler) GetMySubmissionDetail(c *gin.Context) {
	id := c.Param("id")

	userID, _ := c.Get("userID")
	userIDStr, _ := userID.(string)

	email, _ := c.Get("email")
	emailStr, _ := email.(string)

	response, err := h.service.GetMySubmissionDetail(id, userIDStr, emailStr)
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

// SetupRoutes - daftarkan semua route submission
func (h *Handler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// Admin routes
	admin := router.Group("/admin", authMiddleware, middleware.AdminOnly())
	{
		admin.GET("/surveys/:id/submissions", h.ListSubmissions)
		admin.GET("/surveys/:id/export", h.ExportSubmissions)
		admin.GET("/submissions/:id", h.GetSubmissionDetail)
	}

	// Public/Respondent routes
	respondent := router.Group("/survey", authMiddleware, middleware.RespondentOnly())
	{
		respondent.POST("/submit", h.SubmitSurvey)
		respondent.POST("/sync/batch", h.BatchSubmitSurvey)
		respondent.GET("/my-submissions", h.GetMySubmissions)
		respondent.GET("/my-submissions/:id", h.GetMySubmissionDetail)
	}
}
