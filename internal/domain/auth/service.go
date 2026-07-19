package auth

import (
	"atlas_food/internal/pkg/utils"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// Service - interface untuk business logic auth
type Service interface {
	Register(req RegisterRequest) (*AuthResponse, error)
	Login(req LoginRequest) (*AuthResponse, error)
	RefreshToken(refreshToken string) (*AuthResponse, error)
	GetProfile(userID string) (*ProfileResponse, error)
	UpdateProfile(userID string, req UpdateProfileRequest) (*ProfileResponse, error)
	ChangePassword(userID string, req ChangePasswordRequest) error
	UpdatePhoto(userID string, file *multipart.FileHeader) (*ProfileResponse, error)
}

// authService - implementasi Service
type authService struct {
	repo Repository
}

// NewService - factory function untuk membuat service
func NewService(repo Repository) Service {
	return &authService{repo: repo}
}

// Register - daftarkan user baru
func (s *authService) Register(req RegisterRequest) (*AuthResponse, error) {
	// Cek apakah email sudah terdaftar
	existingUser, _ := s.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal hash password")
	}

	// Buat user baru
	user := &User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.Name,
		Role:         "respondent", // default role
		IsActive:     true,
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, errors.New("gagal membuat user")
	}

	// Generate tokens
	accessToken, refreshToken, expiresIn, err := s.generateTokens(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User: UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// Login - autentikasi user
func (s *authService) Login(req LoginRequest) (*AuthResponse, error) {
	// Cari user berdasarkan email
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	// Cek password
	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, errors.New("email atau password salah")
	}

	// Cek apakah user aktif
	if !user.IsActive {
		return nil, errors.New("akun tidak aktif")
	}

	// Generate tokens
	accessToken, refreshToken, expiresIn, err := s.generateTokens(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User: UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// RefreshToken - generate access token baru dari refresh token
func (s *authService) RefreshToken(tokenString string) (*AuthResponse, error) {
	// Hash token untuk cari di database (simpan hash, bukan plain)
	tokenHash := utils.HashSHA256(tokenString)

	// Cari refresh token di database
	refreshToken, err := s.repo.GetRefreshToken(tokenHash)
	if err != nil {
		return nil, errors.New("refresh token tidak valid")
	}

	// Cek apakah token sudah expired
	if time.Now().After(refreshToken.ExpiresAt) {
		// Hapus token yang sudah expired
		s.repo.DeleteRefreshToken(tokenHash)
		return nil, errors.New("refresh token sudah expired")
	}

	// Ambil data user
	user, err := s.repo.GetUserByID(refreshToken.UserID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// Generate tokens baru
	accessToken, newRefreshToken, expiresIn, err := s.generateTokens(user)
	if err != nil {
		return nil, err
	}

	// Hapus refresh token lama (one-time use)
	s.repo.DeleteRefreshToken(tokenHash)

	return &AuthResponse{
		User: UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// GetProfile - ambil profil user
func (s *authService) GetProfile(userID string) (*ProfileResponse, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	return toProfileResponse(user), nil
}

// UpdateProfile - update data profil user (name, phone, gender, birth_date)
func (s *authService) UpdateProfile(userID string, req UpdateProfileRequest) (*ProfileResponse, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if req.Gender != nil && *req.Gender != "" && *req.Gender != "male" && *req.Gender != "female" {
		return nil, errors.New("gender tidak valid")
	}

	var birthDate *time.Time
	if req.BirthDate != nil && *req.BirthDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return nil, errors.New("format tanggal lahir tidak valid")
		}
		if parsed.After(time.Now()) {
			return nil, errors.New("tanggal lahir tidak boleh di masa depan")
		}
		birthDate = &parsed
	}

	user.Name = req.Name
	user.Phone = req.Phone
	if req.Gender != nil && *req.Gender == "" {
		user.Gender = nil
	} else {
		user.Gender = req.Gender
	}
	user.BirthDate = birthDate

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, errors.New("gagal menyimpan profil")
	}

	return toProfileResponse(user), nil
}

// ChangePassword - ganti password user, memaksa re-login di device lain
func (s *authService) ChangePassword(userID string, req ChangePasswordRequest) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if err := utils.CheckPassword(req.CurrentPassword, user.PasswordHash); err != nil {
		return errors.New("password saat ini salah")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("gagal hash password")
	}
	user.PasswordHash = hashedPassword

	if err := s.repo.UpdateUser(user); err != nil {
		return errors.New("gagal menyimpan password baru")
	}

	// Paksa re-login di semua device lain
	_ = s.repo.DeleteUserRefreshTokens(userID)

	return nil
}

// UpdatePhoto - upload/ganti foto profil user
func (s *authService) UpdatePhoto(userID string, file *multipart.FileHeader) (*ProfileResponse, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return nil, errors.New("format file tidak didukung (gunakan jpg, png, webp)")
	}
	if file.Size > 5*1024*1024 {
		return nil, errors.New("ukuran file terlalu besar (max 5MB)")
	}

	uploadDir := filepath.Join("uploads", "profile")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	fullPath := filepath.Join(uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return nil, errors.New("gagal membaca file")
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, errors.New("gagal menyimpan file")
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return nil, errors.New("gagal menyimpan file")
	}

	// Hapus foto lama kalau ada
	if user.PhotoURL != nil && *user.PhotoURL != "" {
		oldPath := filepath.Join(".", *user.PhotoURL)
		if _, statErr := os.Stat(oldPath); statErr == nil {
			os.Remove(oldPath)
		}
	}

	photoURL := fmt.Sprintf("/uploads/profile/%s", filename)
	user.PhotoURL = &photoURL

	if err := s.repo.UpdateUser(user); err != nil {
		return nil, errors.New("gagal menyimpan foto profil")
	}

	return toProfileResponse(user), nil
}

// toProfileResponse - konversi User model ke ProfileResponse DTO
func toProfileResponse(user *User) *ProfileResponse {
	resp := &ProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Phone:     user.Phone,
		Gender:    user.Gender,
		PhotoURL:  user.PhotoURL,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02"),
	}
	if user.BirthDate != nil {
		formatted := user.BirthDate.Format("2006-01-02")
		resp.BirthDate = &formatted
	}
	return resp
}

// generateTokens - helper untuk generate access & refresh token
func (s *authService) generateTokens(user *User) (accessToken, refreshToken string, expiresIn int64, err error) {
	if user.Role != "admin" && user.Role != "respondent" {
		return "", "", 0, errors.New("role tidak valid")
	}

	// Generate access token
	accessToken, err = utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return "", "", 0, errors.New("gagal generate access token")
	}

	// Generate refresh token (random string)
	refreshToken = uuid.New().String()
	tokenHash := utils.HashSHA256(refreshToken)

	// Simpan refresh token ke database
	rt := &RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 hari
	}

	if err := s.repo.CreateRefreshToken(rt); err != nil {
		return "", "", 0, errors.New("gagal simpan refresh token")
	}

	return accessToken, refreshToken, 86400, nil // 24 jam dalam detik
}
