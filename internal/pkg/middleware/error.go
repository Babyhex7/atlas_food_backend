package middleware

import (
	"atlas_food/internal/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler - middleware untuk handle error secara global
// Menangkap error dari handler dan format response yang konsisten
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Cek apakah ada error yang dikirim lewat c.Error(err)
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Default values untuk internal server error
			statusCode := http.StatusInternalServerError
			errorCode := "INTERNAL_SERVER_ERROR"
			message := "Terjadi kesalahan pada server"

			// Cek apakah ini adalah custom AppError
			if appErr, ok := err.Err.(*utils.AppError); ok {
				statusCode = appErr.StatusCode
				errorCode = appErr.Code
				message = appErr.Message
			} else {
				// Jika bukan AppError dan dalam mode debug, tampilkan error aslinya
				if gin.Mode() == gin.DebugMode {
					message = err.Error()
				}
			}

			// Kirim response menggunakan utility agar konsisten
			utils.ErrorResponse(c, statusCode, errorCode, message)
			
			// Hentikan eksekusi lebih lanjut jika perlu (opsional karena ini di akhir)
			c.Abort()
		}
	}
}
