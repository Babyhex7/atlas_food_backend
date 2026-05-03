package utils

import (
	"atlas_food/internal/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims - struct untuk menyimpan data dalam JWT token
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT - membuat token JWT baru untuk user
// userID, email, role: data user yang disimpan dalam token
// Mengembalikan token string atau error
func GenerateJWT(userID, email, role string) (string, error) {
	cfg := config.Load()

	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// ValidateJWT - memvalidasi token JWT dan mengembalikan claims
// tokenString: JWT token yang akan divalidasi
// Mengembalikan JWTClaims jika valid, atau error jika tidak
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	cfg := config.Load()

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Pastikan signing method sesuai
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}
