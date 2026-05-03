package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger - middleware untuk log setiap HTTP request
// Format: [METHOD] path - status - latency - ip
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Waktu mulai request
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Proses request
		c.Next()

		// Hitung latency
		latency := time.Since(start)

		// Ambil status code
		status := c.Writer.Status()

		// Format path dengan query string kalau ada
		if raw != "" {
			path = path + "?" + raw
		}

		// Log format
		clientIP := c.ClientIP()
		method := c.Request.Method

		fmt.Printf("[%s] %s - %d - %v - %s\n", method, path, status, latency, clientIP)
	}
}
