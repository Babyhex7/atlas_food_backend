package food

import (
	"fmt"
	"net/http"

	"atlas_food/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListFoodPhotos(c *gin.Context) {
	result, err := h.service.ListFoodPhotos(c.Param("id"))
	if err != nil {
		respondPhotoError(c, err)
		return
	}
	utils.SuccessResponse(c, result)
}

func (h *Handler) CreateFoodPhoto(c *gin.Context) {
	var req CreateFoodPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Sesi tidak valid")
		return
	}

	result, err := h.service.CreateFoodPhoto(c.Param("id"), fmt.Sprintf("%v", userID), req)
	if err != nil {
		respondPhotoError(c, err)
		return
	}
	utils.CreatedResponse(c, result)
}

func (h *Handler) UpdateFoodPhoto(c *gin.Context) {
	var req UpdateFoodPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Data tidak valid: "+err.Error())
		return
	}

	result, err := h.service.UpdateFoodPhoto(c.Param("id"), c.Param("photoId"), req)
	if err != nil {
		respondPhotoError(c, err)
		return
	}
	utils.SuccessResponse(c, result)
}

func (h *Handler) DeleteFoodPhoto(c *gin.Context) {
	if err := h.service.DeleteFoodPhoto(c.Param("id"), c.Param("photoId")); err != nil {
		respondPhotoError(c, err)
		return
	}
	utils.SuccessResponse(c, gin.H{"message": "Foto dihapus"})
}

func (h *Handler) PublishFoodPhoto(c *gin.Context) {
	result, err := h.service.PublishFoodPhoto(c.Param("id"), c.Param("photoId"))
	if err != nil {
		respondPhotoError(c, err)
		return
	}
	utils.SuccessResponse(c, result)
}

func (h *Handler) UnpublishFoodPhoto(c *gin.Context) {
	result, err := h.service.UnpublishFoodPhoto(c.Param("id"), c.Param("photoId"))
	if err != nil {
		respondPhotoError(c, err)
		return
	}
	utils.SuccessResponse(c, result)
}

func respondPhotoError(c *gin.Context, err error) {
	if appErr, ok := err.(*utils.AppError); ok {
		utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
		return
	}
	utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
}
