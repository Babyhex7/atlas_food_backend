package auth

import (
	"time"
)

// User - model untuk tabel users
type User struct {
	ID           string    `gorm:"type:char(36);primaryKey;default:(UUID())" json:"id"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"` // json:"-" agar tidak ikut di serialize
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Role         string    `gorm:"type:enum('admin','respondent');default:'respondent'" json:"role"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName - set nama tabel untuk GORM
func (User) TableName() string {
	return "users"
}

// RefreshToken - model untuk tabel refresh_tokens
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"type:char(36);not null;index" json:"user_id"`
	TokenHash string    `gorm:"type:varchar(255);not null" json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName - set nama tabel untuk GORM
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
