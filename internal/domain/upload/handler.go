package upload

import (
	"atlas_food/internal/infra/storage"
	"atlas_food/internal/pkg/middleware"
	"atlas_food/internal/pkg/utils"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var safeFolderName = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

// Handler - HTTP handler untuk upload file
type Handler struct {
	UploadPath string
	Storage    storage.Service
}

// NewHandler - buat instance handler upload dengan storage service
func NewHandler(uploadPath string, storageSvc storage.Service) *Handler {
	if storageSvc == nil {
		storageSvc = storage.NewLocalStorage(uploadPath)
	}
	return &Handler{
		UploadPath: uploadPath,
		Storage:    storageSvc,
	}
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
	folder := c.DefaultPostForm("folder", "others") // as-served, annotations, others
	if !safeFolderName.MatchString(folder) {
		utils.ValidationErrorResponse(c, "Nama folder tidak valid")
		return
	}
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Buka stream file
	src, err := file.Open()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "UPLOAD_FAILED", "Gagal membaca file")
		return
	}
	defer src.Close()

	// Simpan ke storage (Local disk atau Cloudflare R2)
	fileURL, err := h.Storage.Upload(c.Request.Context(), folder, filename, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "UPLOAD_FAILED", "Gagal menyimpan file: "+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{
		"url":         fileURL,
		"filename":    filename,
		"size":        file.Size,
		"driver":      h.Storage.Driver(),
		"uploaded_at": time.Now().Format(time.RFC3339),
	})
}

// SetupRoutes - daftarkan route upload
func (h *Handler) SetupRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	// Upload hanya boleh dilakukan oleh admin
	router.POST("/upload", authMiddleware, middleware.AdminOnly(), h.UploadImage)
}
