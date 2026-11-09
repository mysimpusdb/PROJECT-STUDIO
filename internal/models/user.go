package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Don't expose password hash to frontend
	FullName     string    `json:"full_name"`
	AvatarURL    *string   `json:"avatar_url"`
	Timezone                 string    `json:"timezone"`
	NotificationPreferences  string    `json:"notification_preferences"` // JSON string
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
	Token   string `json:"token,omitempty"`
}

type AuthState struct {
	IsLoggedIn bool   `json:"is_logged_in"`
	User       *User  `json:"user,omitempty"`
}

type UserListResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Users   []User `json:"users,omitempty"`
}

type UpdateProfileRequest struct {
	FullName                string  `json:"full_name"`
	AvatarURL               *string `json:"avatar_url,omitempty"`
	Timezone                string  `json:"timezone"`
	NotificationPreferences string  `json:"notification_preferences,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
