package middleware

import (
	"atlas_food/internal/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth - middleware untuk validasi JWT token pada protected routes
// Token diambil dari:
//  1. Header Authorization: Bearer <token> (REST)
//  2. Query ?token=<token> (WebSocket handshake — browser tidak bisa set header)
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Token tidak ditemukan")
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "TOKEN_INVALID", "Token tidak valid atau sudah kadaluarsa")
			c.Abort()
			return
		}

		if claims.Role != "admin" && claims.Role != "respondent" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "ROLE_INVALID", "Role tidak valid")
			c.Abort()
			return
		}

		username := claims.Email
		if at := strings.Index(claims.Email, "@"); at > 0 {
			username = claims.Email[:at]
		}

		c.Set("userID", claims.UserID)
		c.Set("user_id", claims.UserID) // compat WebSocket handler
		c.Set("email", claims.Email)
		c.Set("username", username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// extractToken - ambil JWT dari header "Authorization: Bearer <token>", fallback ke query ?token=
// (dipakai koneksi WebSocket yang tidak bisa kirim custom header)
func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
			return parts[1]
		}
	}
	if q := c.Query("token"); q != "" {
		return q
	}
	return ""
}

// AdminOnly - middleware untuk membatasi akses hanya untuk admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			utils.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "Akses ditolak, hanya admin")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RespondentOnly - middleware untuk membatasi akses hanya untuk respondent
// Admin juga diizinkan agar bisa test/preview survey flow
func RespondentOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
			c.Abort()
			return
		}
		// Admin boleh akses endpoint respondent (untuk test flow)
		if role != "respondent" && role != "admin" {
			utils.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "Akses ditolak")
			c.Abort()
			return
		}
		c.Next()
	}
}
