package food

import (
	"net/http"
	"strconv"

	"atlas_food/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// HTTP handler untuk katalog admin: kategori, as-served set/image, portion method.

// respondError - kirim AppError apa adanya, sisanya jadi 500
func respondError(c *gin.Context, err error) {
	if appErr, ok := err.(*utils.AppError); ok {
		utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
		return
	}
	utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
}

// ============ CATEGORY ============

// ListCategoriesAdmin - GET /admin/categories
func (h *Handler) ListCategoriesAdmin(c *gin.Context) {
	categories, err := h.service.ListCategories()
	if err != nil {
		respondError(c, err)
		return
	}
	utils.SuccessResponse(c, categories)
}

// CreateCategory - POST /admin/categories
func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	category, err := h.service.CreateCategory(req)
	if err != nil {
		respondError(c, err)
		return
	}
	utils.CreatedResponse(c, category)
}

// GetCategory - GET /admin/categories/:id
func (h *Handler) GetCategory(c *gin.Context) {
	category, err := h.service.GetCategory(c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	utils.SuccessResponse(c, category)
}

// UpdateCategory - PUT /admin/categories/:id
func (h *Handler) UpdateCategory(c *gin.Context) {
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	category, err := h.service.UpdateCategory(c.Param("id"), req)
	if err != nil {
		respondError(c, err)
		return
	}
	utils.SuccessResponse(c, category)
}

// DeleteCategory - DELETE /admin/categories/:id
func (h *Handler) DeleteCategory(c *gin.Context) {
	if err := h.service.DeleteCategory(c.Param("id")); err != nil {
		respondError(c, err)
		return
	}
	utils.SuccessResponse(c, gin.H{"message": "Kategori dihapus"})
}

// ============ AS-SERVED SET ============

// ListAsServedSets - GET /admin/as-served-sets
func (h *Handler) ListAsServedSets(c *gin.Context) {
	sets, err := h.service.ListAsServedSets()
	if err != nil {
		respondError(c, err)
		return
	}
	utils.SuccessResponse(c, sets)
}

// CreateAsServedSet - POST /admin/as-served-sets
func (h *Handler) CreateAsServedSet(c *gin.Context) {
	var req CreateAsServedSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	set, err := h.service.CreateAsServedSet(req)
	if err != nil {
		respondError(c, err)
		return
	}
	utils.CreatedResponse(c, set)
}

// GetAsServedSet - GET /admin/as-served-sets/:id (beserta foto)
func (h *Handler) GetAsServedSet(c *gin.Context) {
	set, err := h.service.GetAsServedSet(c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	utils.SuccessResponse(c, set)
}

// UpdateAsServedSet - PUT /admin/as-served-sets/:id
func (h *Handler) UpdateAsServedSet(c *gin.Context) {
	var req UpdateAsServedSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	set, err := h.service.UpdateAsServedSet(c.Param("id"), req)
	if err != nil {
		respondError(c, err)
		return
	}
	utils.SuccessResponse(c, set)
}

// DeleteAsServedSet - DELETE /admin/as-served-sets/:id
func (h *Handler) DeleteAsServedSet(c *gin.Context) {
	if err := h.service.DeleteAsServedSet(c.Param("id")); err != nil {
		respondError(c, err)
		return
	}
	utils.SuccessResponse(c, gin.H{"message": "Set foto porsi dihapus"})
}

// ============ AS-SERVED IMAGE ============

// AddAsServedImages - POST /admin/as-served-sets/:id/images
// Menerima satu objek atau array agar upload batch tidak butuh endpoint terpisah.
func (h *Handler) AddAsServedImages(c *gin.Context) {
	var reqs []AsServedImageRequest

	if err := c.ShouldBindJSON(&reqs); err != nil {
		// Fallback: body berupa satu objek tunggal
		var single AsServedImageRequest
		if errSingle := c.ShouldBindJSON(&single); errSingle != nil {
			utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
			return
		}
		reqs = []AsServedImageRequest{single}
	}

	images, err := h.service.AddAsServedImages(c.Param("id"), reqs)
	if err != nil {
		respondError(c, err)
		return
	}
	utils.CreatedResponse(c, images)
}

// UpdateAsServedImage - PUT /admin/as-served-images/:imageId
func (h *Handler) UpdateAsServedImage(c *gin.Context) {
	var req UpdateAsServedImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	image, err := h.service.UpdateAsServedImage(c.Param("imageId"), req)
	if err != nil {
		respondError(c, err)
		return
	}
	utils.SuccessResponse(c, image)
}

// DeleteAsServedImage - DELETE /admin/as-served-images/:imageId
func (h *Handler) DeleteAsServedImage(c *gin.Context) {
	if err := h.service.DeleteAsServedImage(c.Param("imageId")); err != nil {
		respondError(c, err)
		return
	}
	utils.SuccessResponse(c, gin.H{"message": "Foto porsi dihapus"})
}

// ============ PORTION METHOD ============

// ListPortionMethodsAdmin - GET /admin/foods/:id/portion-methods
func (h *Handler) ListPortionMethodsAdmin(c *gin.Context) {
	methods, err := h.service.ListPortionMethods(c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	utils.SuccessResponse(c, methods)
}

// parseMethodID - baca :methodId sebagai integer
func parseMethodID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("methodId"))
	if err != nil {
		utils.ValidationErrorResponse(c, "ID metode porsi harus berupa angka")
		return 0, false
	}
	return id, true
}

// UpdatePortionMethod - PUT /admin/portion-methods/:methodId
func (h *Handler) UpdatePortionMethod(c *gin.Context) {
	id, ok := parseMethodID(c)
	if !ok {
		return
	}

	var req UpdatePortionMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	method, err := h.service.UpdatePortionMethod(id, req)
	if err != nil {
		respondError(c, err)
		return
	}
	utils.SuccessResponse(c, method)
}

// DeletePortionMethod - DELETE /admin/portion-methods/:methodId
func (h *Handler) DeletePortionMethod(c *gin.Context) {
	id, ok := parseMethodID(c)
	if !ok {
		return
	}

	if err := h.service.DeletePortionMethod(id); err != nil {
		respondError(c, err)
		return
	}
	utils.SuccessResponse(c, gin.H{"message": "Metode porsi dihapus"})
}

// setupCatalogRoutes - daftarkan route katalog ke grup /admin yang sudah
// dilindungi JWTAuth + AdminOnly di SetupRoutes.
func (h *Handler) setupCatalogRoutes(admin *gin.RouterGroup) {
	categories := admin.Group("/categories")
	{
		categories.GET("", h.ListCategoriesAdmin)
		categories.POST("", h.CreateCategory)
		categories.GET("/:id", h.GetCategory)
		categories.PUT("/:id", h.UpdateCategory)
		categories.DELETE("/:id", h.DeleteCategory)
	}

	sets := admin.Group("/as-served-sets")
	{
		sets.GET("", h.ListAsServedSets)
		sets.POST("", h.CreateAsServedSet)
		sets.GET("/:id", h.GetAsServedSet)
		sets.PUT("/:id", h.UpdateAsServedSet)
		sets.DELETE("/:id", h.DeleteAsServedSet)
		sets.POST("/:id/images", h.AddAsServedImages)
	}

	images := admin.Group("/as-served-images")
	{
		images.PUT("/:imageId", h.UpdateAsServedImage)
		images.DELETE("/:imageId", h.DeleteAsServedImage)
	}

	methods := admin.Group("/portion-methods")
	{
		methods.PUT("/:methodId", h.UpdatePortionMethod)
		methods.DELETE("/:methodId", h.DeletePortionMethod)
	}
}
