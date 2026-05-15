package food

import (
	"atlas_food/internal/pkg/middleware"
	"atlas_food/internal/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ============ ADMIN ENDPOINTS ============

func (h *Handler) CreateFood(c *gin.Context) {
	var req CreateFoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	response, err := h.service.CreateFood(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.CreatedResponse(c, response)
}

func (h *Handler) ListFoods(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	categoryID := c.Query("category")

	foods, total, err := h.service.ListFoods(categoryID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"foods": foods,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *Handler) GetFood(c *gin.Context) {
	id := c.Param("id")
	response, err := h.service.GetFoodDetail(id)
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

func (h *Handler) UpdateFood(c *gin.Context) {
	id := c.Param("id")
	var req UpdateFoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	response, err := h.service.UpdateFood(id, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	utils.SuccessResponse(c, response)
}

func (h *Handler) DeleteFood(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteFood(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	utils.SuccessResponse(c, gin.H{"message": "Food deleted successfully"})
}

func (h *Handler) AddPortionMethod(c *gin.Context) {
	foodID := c.Param("id")
	var req CreatePortionMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	response, err := h.service.AddPortionMethod(foodID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	utils.CreatedResponse(c, response)
}

// ============ PUBLIC ENDPOINTS ============

func (h *Handler) SearchFoods(c *gin.Context) {
	query := c.Query("q")
	categoryID := c.Query("category")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	response, err := h.service.SearchFoods(query, categoryID, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	utils.SuccessResponse(c, response)
}

func (h *Handler) ListCategories(c *gin.Context) {
	response, err := h.service.ListCategories()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	utils.SuccessResponse(c, response)
}

// GetFoodsByCategory - GET /api/v1/categories/:id/foods
// Mengambil daftar makanan berdasarkan kategori
func (h *Handler) GetFoodsByCategory(c *gin.Context) {
	categoryID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	foods, total, err := h.service.ListFoods(categoryID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"foods": foods,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *Handler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// Admin routes
	admin := router.Group("/admin", authMiddleware, middleware.AdminOnly())
	{
		foods := admin.Group("/foods")
		{
			foods.POST("", h.CreateFood)
			foods.GET("", h.ListFoods)
			foods.GET("/:id", h.GetFood)
			foods.PUT("/:id", h.UpdateFood)
			foods.DELETE("/:id", h.DeleteFood)
			foods.POST("/:id/portion-methods", h.AddPortionMethod)
		}
	}

	// Public routes
	router.GET("/foods/search", h.SearchFoods)
	router.GET("/foods/:id", h.GetFood)
	router.GET("/categories", h.ListCategories)
	router.GET("/categories/:id/foods", h.GetFoodsByCategory)
}
