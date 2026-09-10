package middleware

import (
	"fmt"
	"strings"
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

			// Format path dengan query string kalau ada.
			// Token JWT ikut lewat query pada handshake WebSocket (browser tidak
			// bisa mengirim header di sana), jadi harus di-redact sebelum dicetak —
			// kalau tidak, kredensial penuh mendarat di log server.
			if raw != "" {
				path = path + "?" + redactSensitiveQuery(raw)
			}

			// Log format
			clientIP := c.ClientIP()
			method := c.Request.Method

			fmt.Printf("[%s] %s - %d - %v - %s\n", method, path, status, latency, clientIP)
		}
	}


	// sensitiveQueryKeys - parameter query yang isinya tidak boleh masuk log
	var sensitiveQueryKeys = map[string]bool{
		"token":         true,
		"access_token":  true,
		"refresh_token": true,
		"invite":        true,
		"api_key":       true,
	}

	// redactSensitiveQuery - ganti nilai parameter sensitif dengan "REDACTED",
	// tanpa mengubah bentuk query string lainnya.
	func redactSensitiveQuery(raw string) string {
		parts := strings.Split(raw, "&")
		for i, part := range parts {
			key, _, found := strings.Cut(part, "=")
			if !found {
				continue
			}
			if sensitiveQueryKeys[strings.ToLower(key)] {
				parts[i] = key + "=REDACTED"
			}
		}
		return strings.Join(parts, "&")
	}
