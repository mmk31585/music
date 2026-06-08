package auth

import (
	"context"
	"strings"

	apperrors "music/internal/common/errors"
)

type Service struct {
	repo   *Repository
	tokens *TokenManager
}

func NewService(repo *Repository, tokens *TokenManager) *Service {
	return &Service{
		repo:   repo,
		tokens: tokens,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest, userAgent, ipAddress string) (AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	username := strings.ToLower(strings.TrimSpace(req.Username))
	displayName := strings.TrimSpace(req.DisplayName)

	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		return AuthResponse{}, err
	}

	user, err := s.repo.CreateUser(ctx, User{
		Email:        email,
		Username:     username,
		DisplayName:  displayName,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return AuthResponse{}, err
	}

	return s.createAuthResponse(ctx, user, userAgent, ipAddress)
}

func (s *Service) Login(ctx context.Context, req LoginRequest, userAgent, ipAddress string) (AuthResponse, error) {
	value := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.repo.FindUserByEmailOrUsername(ctx, value)
	if err != nil {
		return AuthResponse{}, apperrors.Unauthorized("invalid credentials", nil)
	}

	if !user.IsActive {
		return AuthResponse{}, apperrors.Forbidden("user account is disabled", nil)
	}

	if !CheckPassword(req.Password, user.PasswordHash) {
		return AuthResponse{}, apperrors.Unauthorized("invalid credentials", nil)
	}

	return s.createAuthResponse(ctx, user, userAgent, ipAddress)
}

func (s *Service) Refresh(ctx context.Context, req RefreshRequest, userAgent, ipAddress string) (AuthResponse, error) {
	claims, err := s.tokens.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return AuthResponse{}, err
	}

	session, err := s.repo.FindValidSessionByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return AuthResponse{}, err
	}

	if session.UserID != claims.UserID {
		return AuthResponse{}, apperrors.Unauthorized("invalid refresh token", nil)
	}

	user, err := s.repo.FindUserByID(ctx, session.UserID)
	if err != nil {
		return AuthResponse{}, err
	}

	if !user.IsActive {
		return AuthResponse{}, apperrors.Forbidden("user account is disabled", nil)
	}

	// Refresh token rotation:
	// old refresh token becomes invalid after successful refresh.
	if err := s.repo.RevokeSessionByID(ctx, session.ID); err != nil {
		return AuthResponse{}, err
	}

	return s.createAuthResponse(ctx, user, userAgent, ipAddress)
}

func (s *Service) Logout(ctx context.Context, req LogoutRequest) error {
	return s.repo.RevokeSessionByRefreshToken(ctx, req.RefreshToken)
}

func (s *Service) Me(ctx context.Context, userID string) (MeResponse, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return MeResponse{}, err
	}

	return MeResponse{
		User: toAuthUser(user),
	}, nil
}

func (s *Service) createAuthResponse(ctx context.Context, user User, userAgent, ipAddress string) (AuthResponse, error) {
	accessToken, accessExpiresAt, err := s.tokens.GenerateAccessToken(user)
	if err != nil {
		return AuthResponse{}, err
	}

	refreshToken, refreshExpiresAt, err := s.tokens.GenerateRefreshToken(user)
	if err != nil {
		return AuthResponse{}, err
	}

	var userAgentPtr *string
	if userAgent != "" {
		userAgentPtr = &userAgent
	}

	var ipAddressPtr *string
	if ipAddress != "" {
		ipAddressPtr = &ipAddress
	}

	if err := s.repo.CreateSession(ctx, user.ID, refreshToken, userAgentPtr, ipAddressPtr, refreshExpiresAt); err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		User:         toAuthUser(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessExpiresAt.Sub(accessExpiresAt.Add(-s.tokens.accessTTL)).Seconds()),
	}, nil
}

func toAuthUser(user User) AuthUser {
	return AuthUser{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Role:        user.Role,
	}
}
