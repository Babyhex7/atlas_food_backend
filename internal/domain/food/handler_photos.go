package food

import (
	"fmt"
	"net/http"

	"atlas_food/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ListFoodPhotos - GET /admin/foods/:id/photos — daftar foto porsi sebuah makanan
func (h *Handler) ListFoodPhotos(c *gin.Context) {
	result, err := h.service.ListFoodPhotos(c.Param("id"))
	if err != nil {
		respondPhotoError(c, err)
		return
	}
	utils.SuccessResponse(c, result)
}

// CreateFoodPhoto - POST /admin/foods/:id/photos — tambah foto porsi baru (dicatat siapa pembuatnya)
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

// UpdateFoodPhoto - PATCH /admin/foods/:id/photos/:photoId — ubah data foto porsi (gram, label, anotasi)
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

// DeleteFoodPhoto - DELETE /admin/foods/:id/photos/:photoId — hapus foto porsi
func (h *Handler) DeleteFoodPhoto(c *gin.Context) {
	if err := h.service.DeleteFoodPhoto(c.Param("id"), c.Param("photoId")); err != nil {
		respondPhotoError(c, err)
		return
	}
	utils.SuccessResponse(c, gin.H{"message": "Foto dihapus"})
}

// PublishFoodPhoto - POST .../publish — terbitkan foto agar terlihat di aplikasi responden
func (h *Handler) PublishFoodPhoto(c *gin.Context) {
	result, err := h.service.PublishFoodPhoto(c.Param("id"), c.Param("photoId"))
	if err != nil {
		respondPhotoError(c, err)
		return
	}
	utils.SuccessResponse(c, result)
}

// UnpublishFoodPhoto - POST .../unpublish — tarik foto kembali ke status draft
func (h *Handler) UnpublishFoodPhoto(c *gin.Context) {
	result, err := h.service.UnpublishFoodPhoto(c.Param("id"), c.Param("photoId"))
	if err != nil {
		respondPhotoError(c, err)
		return
	}
	utils.SuccessResponse(c, result)
}

// respondPhotoError - ubah error service jadi response HTTP; AppError pakai status/code aslinya, sisanya 500
func respondPhotoError(c *gin.Context, err error) {
	if appErr, ok := err.(*utils.AppError); ok {
		utils.ErrorResponse(c, appErr.StatusCode, appErr.Code, appErr.Message)
		return
	}
	utils.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
}
