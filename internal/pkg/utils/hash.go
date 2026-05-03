package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword - hash password menggunakan bcrypt
// password: plaintext password dari user
// Mengembalikan hash string atau error
func HashPassword(password string) (string, error) {
	// Cost factor 10 (balance antara security dan performance)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword - compare plaintext password dengan hash
// password: plaintext password yang diinput user
// hash: hash password dari database
// Mengembalikan nil jika cocok, error jika tidak
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// HashSHA256 - hash string menggunakan SHA256
// Digunakan untuk hash refresh token sebelum disimpan ke DB
func HashSHA256(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}
