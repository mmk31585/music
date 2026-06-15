package auth

import "time"

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email,max=255"`
	Username    string `json:"username" validate:"required,min=3,max=50"`
	DisplayName string `json:"displayName" validate:"required,min=2,max=100"`
	Password    string `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type AuthResponse struct {
	User         AuthUser `json:"user"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	TokenType    string   `json:"token_type"`
	ExpiresIn    int64    `json:"expires_in"`
}

type MeResponse struct {
	User AuthUser `json:"user"`
}

type AdminUserItem struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Username      string     `json:"username"`
	DisplayName   string     `json:"display_name"`
	Role          string     `json:"role"`
	IsActive      bool       `json:"is_active"`
	EmailVerified bool       `json:"email_verified"`
	AvatarURL     *string    `json:"avatar_url,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

type AdminListUsersResponse struct {
	Items    []AdminUserItem `json:"items"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

type AdminUpdateUserRequest struct {
	Role          *string `json:"role"`
	IsActive      *bool   `json:"is_active"`
	EmailVerified *bool   `json:"email_verified"`
}
