package ai

import (
	"atlas_food/internal/pkg/middleware"
	"atlas_food/internal/pkg/utils"
	"net/http"

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
	ai := router.Group("/ai", authMiddleware, middleware.RespondentOnly())
	{
		ai.POST("/nutrition-analysis", h.AnalyzeNutrition)
	}
}
