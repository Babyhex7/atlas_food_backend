package food

import (
	"atlas_food/internal/pkg/middleware"
	"atlas_food/internal/pkg/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler - HTTP handler domain food, meneruskan request ke Service
type Handler struct {
	service Service
}

// NewHandler - buat handler food dengan service yang dipakai
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ============ ADMIN ENDPOINTS ============

// CreateFood - POST /admin/foods — tambah makanan baru
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

// ListFoods - GET /admin/foods — daftar makanan untuk admin, mendukung filter kategori/pencarian/photo_type/is_active + pagination
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

// GetFood - GET /admin/foods/:id — detail satu makanan (termasuk gizi & metode porsi)
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

// UpdateFood - PUT /admin/foods/:id — ubah data makanan
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

// DeleteFood - DELETE /admin/foods/:id — hapus makanan
func (h *Handler) DeleteFood(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteFood(id); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	utils.SuccessResponse(c, gin.H{"message": "Food deleted successfully"})
}

// AddPortionMethod - POST /admin/foods/:id/portion-methods — tambah metode porsi ke sebuah makanan
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

// GetAssociatedFoods - GET /admin/foods/:id/associated — daftar makanan terkait
func (h *Handler) GetAssociatedFoods(c *gin.Context) {
	foodID := c.Param("id")
	if foodID == "" {
		utils.ValidationErrorResponse(c, "ID makanan wajib diisi")
		return
	}

	items, err := h.service.GetAssociatedFoods(foodID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, items)
}

// SetAssociatedFoods - PUT /admin/foods/:id/associated — atur makanan terkait
func (h *Handler) SetAssociatedFoods(c *gin.Context) {
	foodID := c.Param("id")
	if foodID == "" {
		utils.ValidationErrorResponse(c, "ID makanan wajib diisi")
		return
	}

	var req SetAssociatedFoodsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Format data tidak valid: "+err.Error())
		return
	}

	items, err := h.service.SetAssociatedFoods(foodID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, items)
}

// ============ PUBLIC ENDPOINTS ============

// SearchFoods - GET pencarian makanan untuk responden; query minimal 3 karakter, bisa difilter kategori/tipe
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

// ListCategories - GET daftar kategori makanan untuk ditampilkan di frontend
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

// ListNutrientTypes - GET /nutrient-types atau /admin/nutrient-types — daftar semua jenis nutrisi
func (h *Handler) ListNutrientTypes(c *gin.Context) {
	types, err := h.service.ListNutrientTypes()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	utils.SuccessResponse(c, types)
}

// ListNutrientUnits - GET /nutrient-units atau /admin/nutrient-units — daftar semua satuan nutrisi
func (h *Handler) ListNutrientUnits(c *gin.Context) {
	units, err := h.service.ListNutrientUnits()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	utils.SuccessResponse(c, units)
}

// SetupRoutes - daftarkan seluruh route domain food (admin butuh auth + role admin; route publik diurus PublicHandler)
func (h *Handler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// Endpoint nutrisi dapat diakses melalui /api/v1/nutrient-types & /api/v1/nutrient-units
	router.GET("/nutrient-types", h.ListNutrientTypes)
	router.GET("/nutrient-units", h.ListNutrientUnits)

	// Admin routes
	admin := router.Group("/admin", authMiddleware, middleware.AdminOnly())
	{
		admin.GET("/nutrient-types", h.ListNutrientTypes)
		admin.GET("/nutrient-units", h.ListNutrientUnits)

		foods := admin.Group("/foods")
		{
			foods.POST("", h.CreateFood)
			foods.GET("", h.ListFoods)
			foods.GET("/:id", h.GetFood)
			foods.PUT("/:id", h.UpdateFood)
			foods.DELETE("/:id", h.DeleteFood)
			foods.POST("/:id/portion-methods", h.AddPortionMethod)
			foods.GET("/:id/portion-methods", h.ListPortionMethodsAdmin)

			// Makanan terkait (associated foods)
			foods.GET("/:id/associated", h.GetAssociatedFoods)
			foods.PUT("/:id/associated", h.SetAssociatedFoods)

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
