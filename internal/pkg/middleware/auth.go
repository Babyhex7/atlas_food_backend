package middleware

import (
	"atlas_food/internal/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth - middleware untuk validasi JWT token pada protected routes
// Token diambil dari header Authorization: Bearer <token>
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			c.Abort()
			return
		}

		// Extract Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token salah"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validasi token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}

		// Simpan user info di context untuk handler
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminOnly - middleware untuk membatasi akses hanya untuk admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak, hanya admin"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// SurveyAccessToken - middleware untuk validasi survey accessToken dari URL param
// accessToken diambil dari param URL /surveys/:accessToken
func SurveyAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil accessToken dari param URL
		accessToken := c.Param("accessToken")
		if accessToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "accessToken tidak ditemukan"})
			c.Abort()
			return
		}

		// Validasi accessToken (cukup check format dasar, validasi DB akan dilakukan di handler)
		if len(accessToken) < 20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "accessToken tidak valid"})
			c.Abort()
			return
		}

		// Simpan ke context untuk digunakan handler
		c.Set("accessToken", accessToken)
		c.Next()
	}
}
