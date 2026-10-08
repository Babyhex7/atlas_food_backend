package ai

import (
	"atlas_food/internal/pkg/middleware"
	"atlas_food/internal/pkg/utils"
	"context"
	"errors"
	"log"
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
		utils.ValidationErrorResponse(c, "submission_id wajib diisi")
		return
	}

	// Context sengaja dilepas dari request: bila user menutup halaman di tengah
	// analisis, pemanggilan LLM tetap diselesaikan dan hasilnya tersimpan, jadi
	// kuota yang sudah terpakai tidak terbuang. Batas waktunya dipegang client LLM.
	result, err := h.service.AnalyzeNutrition(context.WithoutCancel(c.Request.Context()), c.GetString("userID"), req)
	if err != nil {
		respondError(c, err)
		return
	}
	respondResult(c, result)
}

// GetNutritionAnalysis - GET /api/v1/ai/nutrition-analysis/:submission_id
// Mengembalikan hasil tersimpan saja; tidak pernah memanggil LLM.
func (h *Handler) GetNutritionAnalysis(c *gin.Context) {
	result, err := h.service.GetNutritionAnalysis(c.GetString("userID"), c.Param("submission_id"))
	if err != nil {
		respondError(c, err)
		return
	}
	respondResult(c, result)
}

// respondResult - bentuk response sukses endpoint AI.
// `source` dan `meta` menjadi saudara `data` (bukan di dalamnya) sesuai kontrak brief AI.
func respondResult(c *gin.Context, result *NutritionAnalysisResult) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"source": result.Source,
		"data":   result.Data,
		"meta":   result.Meta,
	})
}

// respondError - AppError diteruskan apa adanya; error lain dicatat dan disamarkan
func respondError(c *gin.Context, err error) {
	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
		return
	}
	log.Printf("[ai] error tak terduga: %v", err)
	utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Terjadi kesalahan pada server")
}

// SetupRoutes - daftarkan route AI
func (h *Handler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	ai := router.Group("/ai", authMiddleware, middleware.RespondentOnly())
	{
		ai.POST("/nutrition-analysis", h.AnalyzeNutrition)
		ai.GET("/nutrition-analysis/:submission_id", h.GetNutritionAnalysis)
	}
}
