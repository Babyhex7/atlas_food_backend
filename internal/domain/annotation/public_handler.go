package annotation

import (
	"strconv"

	"atlas_food/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// PublicHandler - endpoint tanpa auth untuk aplikasi responden.
// Hanya menyajikan data berstatus published; draft selalu 404.
type PublicHandler struct {
	service *Service
}

// NewPublicHandler - factory PublicHandler
func NewPublicHandler(service *Service) *PublicHandler {
	return &PublicHandler{service: service}
}

// ListPublished - GET /public/food-images
func (h *PublicHandler) ListPublished(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.service.ListPublished(page, limit)
	if err != nil {
		respondError(c, err)
		return
	}

	utils.SuccessResponse(c, result)
}

// GetPublished - GET /public/food-images/:id
func (h *PublicHandler) GetPublished(c *gin.Context) {
	image, err := h.service.GetPublished(c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}

	utils.SuccessResponse(c, image)
}

// ListByFood - GET /public/foods/:id/images
func (h *PublicHandler) ListByFood(c *gin.Context) {
	items, err := h.service.ListPublishedByFood(c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}

	utils.SuccessResponse(c, items)
}

// SetupRoutes - daftarkan route publik anotasi ke grup /public yang sudah ada
func (h *PublicHandler) SetupRoutes(publicGroup *gin.RouterGroup) {
	publicGroup.GET("/food-images", h.ListPublished)
	publicGroup.GET("/food-images/:id", h.GetPublished)
	publicGroup.GET("/foods/:id/images", h.ListByFood)
}
