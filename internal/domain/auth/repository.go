package auth

import (
	"gorm.io/gorm"
)

// Repository - interface untuk operasi database auth
type Repository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	CreateRefreshToken(token *RefreshToken) error
	GetRefreshToken(tokenHash string) (*RefreshToken, error)
	DeleteRefreshToken(tokenHash string) error
	DeleteUserRefreshTokens(userID string) error
}

// authRepository - implementasi Repository
type authRepository struct {
	db *gorm.DB
}

// NewRepository - factory function untuk membuat repository
func NewRepository(db *gorm.DB) Repository {
	return &authRepository{db: db}
}

// CreateUser - insert user baru ke database
func (r *authRepository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

// GetUserByEmail - cari user berdasarkan email
func (r *authRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID - cari user berdasarkan ID
func (r *authRepository) GetUserByID(id string) (*User, error) {
	var user User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateRefreshToken - simpan refresh token ke database
func (r *authRepository) CreateRefreshToken(token *RefreshToken) error {
	return r.db.Create(token).Error
}

// GetRefreshToken - cari refresh token berdasarkan hash
func (r *authRepository) GetRefreshToken(tokenHash string) (*RefreshToken, error) {
	var token RefreshToken
	err := r.db.Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// DeleteRefreshToken - hapus refresh token (logout)
func (r *authRepository) DeleteRefreshToken(tokenHash string) error {
	return r.db.Where("token_hash = ?", tokenHash).Delete(&RefreshToken{}).Error
}

// DeleteUserRefreshTokens - hapus semua refresh token user (force logout semua device)
func (r *authRepository) DeleteUserRefreshTokens(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&RefreshToken{}).Error
}
