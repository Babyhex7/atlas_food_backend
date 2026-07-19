package auth

// RegisterRequest - DTO untuk request register
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2,max=100"`
}

// LoginRequest - DTO untuk request login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest - DTO untuk request refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UpdateProfileRequest - DTO untuk request update profile
type UpdateProfileRequest struct {
	Name      string  `json:"name" binding:"required,min=2,max=100"`
	Phone     *string `json:"phone"`
	Gender    *string `json:"gender"`
	BirthDate *string `json:"birth_date"` // format "2006-01-02"
}

// ChangePasswordRequest - DTO untuk request ganti password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// AuthResponse - DTO untuk response auth (login/register)
type AuthResponse struct {
	User         UserInfo `json:"user"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
}

// UserInfo - informasi user untuk response
type UserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

// ProfileResponse - DTO untuk response profile
type ProfileResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	Phone     *string `json:"phone,omitempty"`
	Gender    *string `json:"gender,omitempty"`
	BirthDate *string `json:"birth_date,omitempty"`
	PhotoURL  *string `json:"photo_url,omitempty"`
	Role      string  `json:"role"`
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
}
