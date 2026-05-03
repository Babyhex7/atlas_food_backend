package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler - middleware untuk handle error secara global
// Menangkap error dari handler dan format response yang konsisten
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Cek error dari context
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": err.Error(),
				},
			})
		}
	}
}
