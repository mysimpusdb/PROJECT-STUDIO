package services

import (
	"context"
	"fmt"
	"strings"

	"project-studio/internal/auth"
	"project-studio/internal/db"
	"project-studio/internal/models"
	"project-studio/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService provides authentication methods for the frontend
type AuthService struct {
	ctx      context.Context
	userRepo *repository.UserRepository
}

// NewAuthService creates a new AuthService
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(db.DB),
	}
}

// Startup is called at application startup
func (s *AuthService) Startup(ctx context.Context) {
	s.ctx = ctx
	// Update userRepo with the initialized DB
	s.userRepo = repository.NewUserRepository(db.DB)
}

// Login attempts to authenticate a user
func (s *AuthService) Login(req models.LoginRequest) models.LoginResponse {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return models.LoginResponse{
			Success: false,
			Message: "Terjadi kesalahan sistem.",
		}
	}

	if user == nil {
		return models.LoginResponse{
			Success: false,
			Message: "Email atau kata sandi salah.",
		}
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return models.LoginResponse{
			Success: false,
			Message: "Email atau kata sandi salah.",
		}
	}

	// Generate a dummy session token (in a real app, use JWT or similar)
	token := uuid.New().String()

	// Save to keyring
	err = auth.SaveToken(user.ID, token)
	if err != nil {
		fmt.Printf("Warning: failed to save token to keyring: %v\n", err)
	}

	return models.LoginResponse{
		Success: true,
		Message: "Login berhasil.",
		User:    user,
		Token:   token,
	}
}

// Register creates a new local user, attaches them to the default workspace
// (creating one if none exists yet), and signs them in.
func (s *AuthService) Register(req models.RegisterRequest) models.LoginResponse {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	name := strings.TrimSpace(req.FullName)
	if name == "" || email == "" || len(req.Password) < 6 {
		return models.LoginResponse{
			Success: false,
			Message: "Nama, email, dan kata sandi (min. 6 karakter) wajib diisi.",
		}
	}

	existing, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return models.LoginResponse{Success: false, Message: "Terjadi kesalahan sistem."}
	}
	if existing != nil {
		return models.LoginResponse{Success: false, Message: "Email sudah terdaftar."}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.LoginResponse{Success: false, Message: "Gagal memproses kata sandi."}
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hash),
		FullName:     name,
		Timezone:     "Asia/Jakarta",
	}
	if err := s.userRepo.Create(user); err != nil {
		return models.LoginResponse{Success: false, Message: "Gagal mendaftarkan pengguna."}
	}

	// Attach to the first workspace, or create one if this is the very first user.
	var workspaceID string
	err = db.DB.QueryRow(`SELECT id FROM workspaces ORDER BY created_at ASC LIMIT 1`).Scan(&workspaceID)
	role := "member"
	if err != nil || workspaceID == "" {
		workspaceID = uuid.New().String()
		_, _ = db.DB.Exec(`INSERT INTO workspaces (id, name, owner_id) VALUES (?, ?, ?)`, workspaceID, "Default Workspace", user.ID)
		role = "admin"
	}
	_, _ = db.DB.Exec(`INSERT OR IGNORE INTO workspace_members (workspace_id, user_id, role) VALUES (?, ?, ?)`, workspaceID, user.ID, role)

	token := uuid.New().String()
	_ = auth.SaveToken(user.ID, token)

	return models.LoginResponse{
		Success: true,
		Message: "Registrasi berhasil.",
		User:    user,
		Token:   token,
	}
}

// CheckAuth verifies if a user is currently logged in via the Keyring
func (s *AuthService) CheckAuth() models.AuthState {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.AuthState{
			IsLoggedIn: false,
		}
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return models.AuthState{
			IsLoggedIn: false,
		}
	}

	return models.AuthState{
		IsLoggedIn: true,
		User:       user,
	}
}

// Logout clears the current session
func (s *AuthService) Logout() bool {
	auth.DeleteToken()
	return true
}

// UpdateProfile updates the current logged-in user's profile
func (s *AuthService) UpdateProfile(req models.UpdateProfileRequest) models.BaseResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.BaseResponse{Success: false, Message: "Unauthenticated"}
	}

	fullName := strings.TrimSpace(req.FullName)
	if fullName == "" {
		return models.BaseResponse{Success: false, Message: "Nama lengkap tidak boleh kosong."}
	}

	tz := strings.TrimSpace(req.Timezone)
	if tz == "" {
		tz = "Asia/Jakarta"
	}

	err = s.userRepo.UpdateProfile(userID, fullName, tz, req.NotificationPreferences, req.AvatarURL)
	if err != nil {
		return models.BaseResponse{Success: false, Message: "Gagal memperbarui profil."}
	}

	return models.BaseResponse{Success: true, Message: "Profil berhasil diperbarui."}
}

// ChangePassword updates the current user's password
func (s *AuthService) ChangePassword(req models.ChangePasswordRequest) models.BaseResponse {
	_, userID, err := auth.GetToken()
	if err != nil || userID == "" {
		return models.BaseResponse{Success: false, Message: "Unauthenticated"}
	}

	if len(req.NewPassword) < 6 {
		return models.BaseResponse{Success: false, Message: "Kata sandi baru minimal 6 karakter."}
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return models.BaseResponse{Success: false, Message: "Pengguna tidak ditemukan."}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword))
	if err != nil {
		return models.BaseResponse{Success: false, Message: "Kata sandi lama salah."}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return models.BaseResponse{Success: false, Message: "Gagal memproses kata sandi baru."}
	}

	err = s.userRepo.UpdatePassword(userID, string(hash))
	if err != nil {
		return models.BaseResponse{Success: false, Message: "Gagal menyimpan kata sandi baru."}
	}

	return models.BaseResponse{Success: true, Message: "Kata sandi berhasil diubah."}
}
