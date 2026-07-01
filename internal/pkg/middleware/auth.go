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
			utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Token tidak ditemukan")
			c.Abort()
			return
		}

		// Extract Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Format token salah")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validasi token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "TOKEN_INVALID", "Token tidak valid atau sudah kadaluarsa")
			c.Abort()
			return
		}

		// Simpan user info di context untuk handler
		if claims.Role != "admin" && claims.Role != "respondent" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "ROLE_INVALID", "Role tidak valid")
			c.Abort()
			return
		}

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

// SurveyAccessToken - middleware untuk validasi survey accessToken dari URL param
// accessToken diambil dari param URL /surveys/:accessToken
func SurveyAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Param("accessToken")
		if accessToken == "" {
			utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "accessToken tidak ditemukan")
			c.Abort()
			return
		}

		if len(accessToken) < 20 {
			utils.ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", "accessToken tidak valid")
			c.Abort()
			return
		}

		c.Set("accessToken", accessToken)
		c.Next()
	}
}
