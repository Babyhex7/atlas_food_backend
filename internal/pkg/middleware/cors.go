package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS - middleware untuk mengatur Cross-Origin Resource Sharing
//
// Origin dicocokkan ke allowlist (FRONTEND_URL + CORS_ALLOWED_ORIGINS, plus
// localhost di luar produksi). Sebelumnya middleware ini memantulkan Origin
// apa pun sambil menyalakan Allow-Credentials — kombinasi yang membuat situs
// mana pun bisa memanggil API ini atas nama user yang sedang login.
func CORS(allowed []string) gin.HandlerFunc {
	isAllowed := func(origin string) bool {
		trimmed := strings.TrimRight(origin, "/")
		for _, a := range allowed {
			if strings.EqualFold(a, trimmed) {
				return true
			}
		}
		return false
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Permintaan same-origin / non-browser (curl, health check) tidak
		// mengirim Origin sama sekali — biarkan lewat tanpa header CORS.
		if origin != "" {
			if !isAllowed(origin) {
				if c.Request.Method == http.MethodOptions {
					c.AbortWithStatus(http.StatusForbidden)
					return
				}
				// Tanpa header Allow-Origin, browser sendiri yang memblokir responsnya.
				c.Next()
				return
			}

			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			// Respons berbeda per Origin — tanpa Vary, cache bisa menyajikan
			// header Allow-Origin milik origin lain.
			c.Writer.Header().Add("Vary", "Origin")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With, Idempotency-Key")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			c.Writer.Header().Set("Access-Control-Max-Age", "600")
			// Tanpa ini, browser strip header Content-Disposition dari response
			// cross-origin sehingga frontend tidak bisa baca nama file asli saat
			// download (mis. export CSV survey).
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		}

		// Handle preflight request
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
