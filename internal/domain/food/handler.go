package food

import (
	"atlas_food/internal/pkg/middleware"
	"atlas_food/internal/pkg/utils"
	"net/http"
	"strconv"
	"strings"

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
	search := c.Query("search")
	photoType := c.Query("photo_type")

	filter := ListFoodsFilter{
		CategoryID: categoryID,
		Search:     search,
		PhotoType:  photoType,
		Page:       page,
		Limit:      limit,
	}

	if status := c.Query("is_active"); status == "true" {
		v := true
		filter.IsActive = &v
	} else if status == "false" {
		v := false
		filter.IsActive = &v
	}

	foods, total, err := h.service.ListFoodsAdmin(filter)
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
	if len(strings.TrimSpace(query)) < 3 {
		utils.ValidationErrorResponse(c, "Query pencarian minimal 3 karakter")
		return
	}
	categoryID := c.Query("category")
	// foodType: "food" | "drink" | "" (kosong = semua)
	foodType := c.Query("type")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	response, err := h.service.SearchFoods(query, categoryID, foodType, limit)
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

// GetFoodsByCategory - GET /api/v1/public/categories/:code/foods
// Mengambil daftar makanan berdasarkan kategori code
func (h *Handler) GetFoodsByCategory(c *gin.Context) {
	categoryCode := c.Param("code")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	foods, total, err := h.service.ListFoodsByCategoryCode(categoryCode, page, limit)
	if err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
			return
		}
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
			foods.GET("/:id/portion-methods", h.ListPortionMethodsAdmin)

			// Foto makanan terpadu: anotasi draft + gram porsi (series/range)
			foods.GET("/:id/photos", h.ListFoodPhotos)
			foods.POST("/:id/photos", h.CreateFoodPhoto)
			foods.PATCH("/:id/photos/:photoId", h.UpdateFoodPhoto)
			foods.DELETE("/:id/photos/:photoId", h.DeleteFoodPhoto)
			foods.POST("/:id/photos/:photoId/publish", h.PublishFoodPhoto)
			foods.POST("/:id/photos/:photoId/unpublish", h.UnpublishFoodPhoto)
		}

		// Kategori, as-served set/image, dan metode porsi
		h.setupCatalogRoutes(admin)
	}

	// Respondent routes
	respondent := router.Group("", authMiddleware, middleware.RespondentOnly())
	{
		// Route khusus responden (reserved)
		_ = respondent
	}

	// NOTE: Public routes now handled by PublicHandler in router.go
	// No need to register them here to avoid duplicate routes
}
