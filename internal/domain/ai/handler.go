package ai

import (
	"atlas_food/internal/pkg/middleware"
	"atlas_food/internal/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler - HTTP handler AI
type Handler struct {
	service Service
}

// NewHandler - factory handler AI
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// AskAssistant - POST /api/v1/ai/assistant
// Menjawab pertanyaan porsi / URT / bahan pelengkap dari responden
func (h *Handler) AskAssistant(c *gin.Context) {
	var req AssistantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	res, err := h.service.AskAssistant(c.Request.Context(), req)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   res,
	})
}

// AnalyzeNutrition - POST /api/v1/ai/nutrition-analysis
func (h *Handler) AnalyzeNutrition(c *gin.Context) {
	var req NutritionAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User tidak terautentikasi")
		return
	}

	result, err := h.service.AnalyzeNutrition(userID.(string), req)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"source": result.Source,
		"data":   result.Data,
	})
}

// SetupRoutes - daftarkan route AI
func (h *Handler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	aiGroup := router.Group("/ai")
	{
		// Asisten interaktif per step recall (rate limit 60 req/min)
		aiGroup.POST("/assistant", middleware.RateLimiter(60, time.Minute), h.AskAssistant)

		// Analisis gizi pasca-recall (harus login respondent)
		aiGroup.POST("/nutrition-analysis", authMiddleware, middleware.RespondentOnly(), h.AnalyzeNutrition)
	}
}
