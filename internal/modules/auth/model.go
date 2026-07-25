package auth

import "time"

type User struct {
	ID              string
	Email           string
	Username        string
	DisplayName     string
	PasswordHash    string
	AvatarURL       *string
	Role            string
	IsActive        bool
	EmailVerifiedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type AuthSession struct {
	ID               string
	UserID           string
	RefreshTokenHash string
	UserAgent        *string
	IPAddress        *string
	ExpiresAt        time.Time
	RevokedAt        *time.Time
	CreatedAt        time.Time
}

type ListUsersParams struct {
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
	Search    string
}

type AuthUser struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	Username    string  `json:"username"`
	DisplayName string  `json:"displayName"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
	Role        string  `json:"role"`
}
