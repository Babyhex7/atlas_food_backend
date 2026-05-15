package upload

import (
	"atlas_food/internal/pkg/middleware"
	"atlas_food/internal/pkg/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler - HTTP handler untuk upload file
type Handler struct {
	UploadPath string
}

// NewHandler - buat instance handler upload
func NewHandler(uploadPath string) *Handler {
	// Pastikan folder upload ada
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.MkdirAll(uploadPath, 0755)
	}
	return &Handler{UploadPath: uploadPath}
}

// UploadImage - POST /api/v1/upload
// Upload gambar untuk makanan atau porsi
func (h *Handler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		utils.ValidationErrorResponse(c, "File tidak ditemukan")
		return
	}

	// Validasi tipe file
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		utils.ValidationErrorResponse(c, "Format file tidak didukung (gunakan jpg, png, webp)")
		return
	}

	// Validasi ukuran (max 10MB)
	if file.Size > 10*1024*1024 {
		utils.ValidationErrorResponse(c, "Ukuran file terlalu besar (max 10MB)")
		return
	}

	// Generate nama file unik
	folder := c.DefaultPostForm("folder", "others") // as-served, foods, others
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	
	// Simpan ke disk
	savePath := filepath.Join(h.UploadPath, folder)
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		os.MkdirAll(savePath, 0755)
	}

	fullPath := filepath.Join(savePath, filename)
	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "UPLOAD_FAILED", "Gagal menyimpan file")
		return
	}

	// Response dengan URL file
	// Note: In real production, this should be a full URL or relative path handled by static server
	fileURL := fmt.Sprintf("/uploads/%s/%s", folder, filename)

	utils.SuccessResponse(c, gin.H{
		"url":       fileURL,
		"filename":  filename,
		"size":      file.Size,
		"uploaded_at": time.Now().Format(time.RFC3339),
	})
}

// SetupRoutes - daftarkan route upload
func (h *Handler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// Upload hanya boleh dilakukan oleh admin
	router.POST("/upload", authMiddleware, middleware.AdminOnly(), h.UploadImage)
}
