package annotation

import (
	"fmt"
	"net/http"

	"atlas_food/internal/pkg/middleware"
	"atlas_food/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Handler - endpoint admin untuk Food Annotation CMS
type Handler struct {
	service *Service
}

// NewHandler - factory Handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// respondError - kirim AppError apa adanya, sisanya jadi 500
func respondError(c *gin.Context, err error) {
	if appErr, ok := err.(*utils.AppError); ok {
		utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
		return
	}
	utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
}

// List - GET /admin/food-images
func (h *Handler) List(c *gin.Context) {
	var q ListFoodImagesQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.ValidationErrorResponse(c, "Query tidak valid: "+err.Error())
		return
	}

	result, err := h.service.List(q)
	if err != nil {
		respondError(c, err)
		return
	}

	utils.SuccessResponse(c, result)
}

// Create - POST /admin/food-images
func (h *Handler) Create(c *gin.Context) {
	var req CreateFoodImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Sesi tidak valid")
		return
	}

	image, err := h.service.Create(req, fmt.Sprintf("%v", userID))
	if err != nil {
		respondError(c, err)
		return
	}

	utils.CreatedResponse(c, image)
}

// Get - GET /admin/food-images/:id
func (h *Handler) Get(c *gin.Context) {
	image, err := h.service.Get(c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}

	utils.SuccessResponse(c, image)
}

// Update - PATCH /admin/food-images/:id
func (h *Handler) Update(c *gin.Context) {
	var req UpdateFoodImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	image, err := h.service.Update(c.Param("id"), req)
	if err != nil {
		respondError(c, err)
		return
	}

	utils.SuccessResponse(c, image)
}

// ReplaceAreas - PUT /admin/food-images/:id/areas (autosave editor)
func (h *Handler) ReplaceAreas(c *gin.Context) {
	var req ReplaceAreasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	result, err := h.service.ReplaceAreas(c.Param("id"), req)
	if err != nil {
		respondError(c, err)
		return
	}

	utils.SuccessResponse(c, result)
}

// Publish - POST /admin/food-images/:id/publish
func (h *Handler) Publish(c *gin.Context) {
	image, err := h.service.Publish(c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}

	utils.SuccessResponse(c, image)
}

// Unpublish - POST /admin/food-images/:id/unpublish
func (h *Handler) Unpublish(c *gin.Context) {
	image, err := h.service.Unpublish(c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}

	utils.SuccessResponse(c, image)
}

// Delete - DELETE /admin/food-images/:id
func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		respondError(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Gambar anotasi dihapus"})
}

// ExportJSON - GET /admin/food-images/:id/export
func (h *Handler) ExportJSON(c *gin.Context) {
	id := c.Param("id")

	image, err := h.service.Get(id)
	if err != nil {
		respondError(c, err)
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="annotation-%s.json"`, id))
	c.JSON(http.StatusOK, image)
}

// SetupRoutes - daftarkan route admin anotasi.
// Seluruh grup dilindungi JWTAuth + AdminOnly sesuai brief §3.
func (h *Handler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	admin := router.Group("/admin/food-images", authMiddleware, middleware.AdminOnly())
	{
		admin.GET("", h.List)
		admin.POST("", h.Create)
		admin.GET("/:id", h.Get)
		admin.PATCH("/:id", h.Update)
		admin.DELETE("/:id", h.Delete)
		admin.PUT("/:id/areas", h.ReplaceAreas)
		admin.POST("/:id/publish", h.Publish)
		admin.POST("/:id/unpublish", h.Unpublish)
		admin.GET("/:id/export", h.ExportJSON)
	}
}
