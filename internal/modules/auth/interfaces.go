package auth

import (
	"context"
	"time"
)

// TODO: RepositoryInterface is defined but never used as an injection type. Service uses concrete *Repository directly.
// Either wire it into Service or remove.
type RepositoryInterface interface {
	CreateUser(ctx context.Context, user User) (User, error)
	FindUserByEmailOrUsername(ctx context.Context, value string) (User, error)
	FindUserByID(ctx context.Context, id string) (User, error)
	FindPublicUser(ctx context.Context, id string) (User, error)
	FindValidSessionByRefreshToken(ctx context.Context, refreshToken string) (AuthSession, error)
	CreateSession(ctx context.Context, userID, refreshToken string, userAgent, ipAddress *string, expiresAt time.Time) error
	RevokeSessionByRefreshToken(ctx context.Context, refreshToken string) error
	RevokeSessionByID(ctx context.Context, sessionID string) error
	UpdateUser(ctx context.Context, id string, updates map[string]any) error
	UpdatePassword(ctx context.Context, id, passwordHash string) error
	ListUsers(ctx context.Context, params ListUsersParams) ([]User, int, error)
	DeleteUser(ctx context.Context, id string) error
}
