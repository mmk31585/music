package auth

import (
	"context"
	"strings"
	"time"

	apperrors "music/internal/common/errors"
)

type Service struct {
	repo   RepositoryInterface
	tokens *TokenManager
}

func NewService(repo RepositoryInterface, tokens *TokenManager) *Service {
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

func (s *Service) AdminListUsers(ctx context.Context, params ListUsersParams) (AdminListUsersResponse, error) {
	users, total, err := s.repo.ListUsers(ctx, params)
	if err != nil {
		return AdminListUsersResponse{}, err
	}

	items := make([]AdminUserItem, len(users))
	for i, u := range users {
		items[i] = toAdminUserItem(u)
	}

	return AdminListUsersResponse{
		Items:    items,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

func (s *Service) AdminGetUser(ctx context.Context, id string) (AdminUserItem, error) {
	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return AdminUserItem{}, err
	}
	return toAdminUserItem(user), nil
}

func (s *Service) AdminUpdateUser(ctx context.Context, id string, req AdminUpdateUserRequest) (AdminUserItem, error) {
	updates := make(map[string]any)
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.EmailVerified != nil {
		updates["email_verified"] = *req.EmailVerified
	}

	if len(updates) > 0 {
		if err := s.repo.UpdateUser(ctx, id, updates); err != nil {
			return AdminUserItem{}, err
		}
	}

	return s.AdminGetUser(ctx, id)
}

func (s *Service) UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) error {
	updates := make(map[string]any)
	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Username != nil {
		updates["username"] = strings.ToLower(strings.TrimSpace(*req.Username))
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.Website != nil {
		updates["website"] = *req.Website
	}
	if req.Preferences != nil {
		updates["preferences"] = *req.Preferences
	}

	if len(updates) == 0 {
		return nil
	}

	return s.repo.UpdateUser(ctx, userID, updates)
}

func (s *Service) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if !CheckPassword(req.CurrentPassword, user.PasswordHash) {
		return apperrors.Unauthorized("invalid credentials", nil)
	}

	hash, err := HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, userID, hash)
}

func (s *Service) GetPublicProfile(ctx context.Context, id string) (*User, error) {
	user, err := s.repo.FindPublicUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) AdminDeleteUser(ctx context.Context, id string) error {
	return s.repo.DeleteUser(ctx, id)
}

func toAdminUserItem(user User) AdminUserItem {
	verified := user.EmailVerifiedAt != nil
	var updatedAt *time.Time
	if !user.UpdatedAt.IsZero() {
		updatedAt = &user.UpdatedAt
	}
	return AdminUserItem{
		ID:            user.ID,
		Email:         user.Email,
		Username:      user.Username,
		DisplayName:   user.DisplayName,
		Role:          user.Role,
		IsActive:      user.IsActive,
		EmailVerified: verified,
		AvatarURL:     user.AvatarURL,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     updatedAt,
	}
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
