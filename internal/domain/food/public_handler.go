package food

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// PublicHandler handles public (no-auth) food endpoints
type PublicHandler struct {
	repo Repository
}

// NewPublicHandler creates a new public food handler
func NewPublicHandler(repo Repository) *PublicHandler {
	return &PublicHandler{repo: repo}
}

// SearchFoods godoc
// @Summary Search foods (Public, no auth)
// @Description Search foods by name or local_name using FULLTEXT search
// @Tags public-food
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param limit query int false "Limit results" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /public/foods/search [get]
func (h *PublicHandler) SearchFoods(c *gin.Context) {
	query := c.DefaultQuery("q", "")
	limit := 20
	
	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	foods, err := h.repo.SearchFoodsPublic(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to search foods",
		})
		return
	}

	// Transform to response
	var results []SearchFoodResponse
	for _, food := range foods {
		item := SearchFoodResponse{
			ID:        food.ID,
			Code:      food.Code,
			Name:      food.Name,
			LocalName: food.LocalName,
			PhotoType: food.PhotoType,
		}
		
		if food.Category != nil {
			item.Category = &CategoryInfo{
				ID:   food.Category.ID,
				Code: food.Category.Code,
				Name: food.Category.Name,
				Icon: food.Category.Icon,
			}
		}
		
		results = append(results, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   results,
		"count":  len(results),
	})
}

// GetFoodDetail godoc
// @Summary Get food detail with portion photos (Public)
// @Description Get complete food information including nutrients and portion photos
// @Tags public-food
// @Accept json
// @Produce json
// @Param id path string true "Food ID"
// @Success 200 {object} map[string]interface{}
// @Router /public/foods/{id} [get]
func (h *PublicHandler) GetFoodDetail(c *gin.Context) {
	foodID := c.Param("id")

	// Use channels for concurrent data fetching
	type dataResult struct {
		food          *Food
		portionPhotos []AsServedImage
		nutrients     []FoodNutrient
		err           error
		dataType      string
	}

	resultChan := make(chan dataResult, 2)

	// Goroutine 1: Get food with portion photos
	go func() {
		food, photos, err := h.repo.GetFoodWithPortionPhotos(foodID)
		resultChan <- dataResult{
			food:          food,
			portionPhotos: photos,
			err:           err,
			dataType:      "food",
		}
	}()

	// Goroutine 2: Get nutrients
	go func() {
		nutrients, err := h.repo.GetFoodNutrients(foodID)
		resultChan <- dataResult{
			nutrients: nutrients,
			err:       err,
			dataType:  "nutrients",
		}
	}()

	// Collect results
	var food *Food
	var portionPhotos []AsServedImage
	var nutrients []FoodNutrient

	for i := 0; i < 2; i++ {
		res := <-resultChan
		
		switch res.dataType {
		case "food":
			if res.err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  "error",
					"message": "Food not found",
				})
				return
			}
			food = res.food
			portionPhotos = res.portionPhotos
			
		case "nutrients":
			if res.err == nil {
				nutrients = res.nutrients
			}
			// Don't fail if nutrients not found, just use empty array
		}
	}

	// Transform nutrients to map
	nutrientMap := make(map[string]NutrientDetail)
	for _, n := range nutrients {
		nutrientMap[n.NutrientType.Code] = NutrientDetail{
			Value: n.ValuePer100g,
			Unit:  n.NutrientType.Unit.Symbol,
		}
	}

	// Transform portion photos — food_image_id untuk overlay anotasi responden
	var portionPhotoList []PortionPhoto
	for _, p := range portionPhotos {
		desc := p.Description
		foodImageID := ""
		if strings.HasPrefix(desc, "food_image:") {
			foodImageID = strings.TrimPrefix(desc, "food_image:")
			desc = ""
		}
		portionPhotoList = append(portionPhotoList, PortionPhoto{
			ID:           p.ID,
			Label:        p.Label,
			ImageURL:     p.ImageURL,
			ThumbnailURL: p.ThumbnailURL,
			WeightGram:   p.WeightGram,
			Description:  desc,
			FoodImageID:  foodImageID,
		})
	}

	response := FoodResponse{
		ID:            food.ID,
		Code:          food.Code,
		Name:          food.Name,
		LocalName:     food.LocalName,
		Description:   food.Description,
		PhotoType:     food.PhotoType,
		Nutrients:     nutrientMap,
		PortionPhotos: portionPhotoList,
	}

	if food.Category != nil {
		response.Category = &CategoryInfo{
			ID:   food.Category.ID,
			Code: food.Category.Code,
			Name: food.Category.Name,
			Icon: food.Category.Icon,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}

// GetCategories godoc
// @Summary Get all categories (Public)
// @Description Get list of all food categories
// @Tags public-food
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /public/categories [get]
func (h *PublicHandler) GetCategories(c *gin.Context) {
	categories, err := h.repo.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get categories",
		})
		return
	}

	var results []CategoryInfo
	for _, cat := range categories {
		results = append(results, CategoryInfo{
			ID:   cat.ID,
			Code: cat.Code,
			Name: cat.Name,
			Icon: cat.Icon,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   results,
	})
}

// GetFoodsByCategory godoc
// @Summary Get foods by category (Public)
// @Description Get all foods in a specific category
// @Tags public-food
// @Accept json
// @Produce json
// @Param code path string true "Category code (e.g., MP, LH, AB)"
// @Param limit query int false "Limit results" default(50)
// @Success 200 {object} map[string]interface{}
// @Router /public/categories/{code}/foods [get]
func (h *PublicHandler) GetFoodsByCategory(c *gin.Context) {
	categoryCode := c.Param("code")
	limit := 50
	
	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 200 {
			limit = l
		}
	}

	foods, err := h.repo.GetFoodsByCategory(categoryCode, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get foods",
		})
		return
	}

	// Transform to response
	var results []SearchFoodResponse
	for _, food := range foods {
		item := SearchFoodResponse{
			ID:        food.ID,
			Code:      food.Code,
			Name:      food.Name,
			LocalName: food.LocalName,
			PhotoType: food.PhotoType,
		}
		
		if food.Category != nil {
			item.Category = &CategoryInfo{
				ID:   food.Category.ID,
				Code: food.Category.Code,
				Name: food.Category.Name,
				Icon: food.Category.Icon,
			}
		}
		
		results = append(results, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   results,
		"count":  len(results),
	})
}
