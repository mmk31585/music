package auth

import (
	"context"
	"testing"
	"time"

	apperrors "music/internal/common/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) CreateUser(ctx context.Context, user User) (User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(User), args.Error(1)
}

func (m *mockRepo) FindUserByEmailOrUsername(ctx context.Context, value string) (User, error) {
	args := m.Called(ctx, value)
	return args.Get(0).(User), args.Error(1)
}

func (m *mockRepo) FindUserByID(ctx context.Context, id string) (User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(User), args.Error(1)
}

func (m *mockRepo) FindPublicUser(ctx context.Context, id string) (User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(User), args.Error(1)
}

func (m *mockRepo) CreateSession(ctx context.Context, userID, refreshToken string, userAgent, ipAddress *string, expiresAt time.Time) error {
	return m.Called(ctx, userID, refreshToken, userAgent, ipAddress, expiresAt).Error(0)
}

func (m *mockRepo) UpdateUser(ctx context.Context, id string, updates map[string]any) error {
	return m.Called(ctx, id, updates).Error(0)
}

func (m *mockRepo) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	return m.Called(ctx, id, passwordHash).Error(0)
}

func (m *mockRepo) ListUsers(ctx context.Context, params ListUsersParams) ([]User, int, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]User), args.Int(1), args.Error(2)
}

func (m *mockRepo) FindValidSessionByRefreshToken(ctx context.Context, refreshToken string) (AuthSession, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(AuthSession), args.Error(1)
}

func (m *mockRepo) RevokeSessionByRefreshToken(ctx context.Context, refreshToken string) error {
	return m.Called(ctx, refreshToken).Error(0)
}

func (m *mockRepo) RevokeSessionByID(ctx context.Context, sessionID string) error {
	return m.Called(ctx, sessionID).Error(0)
}

func (m *mockRepo) DeleteUser(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func TestService_Register_Success(t *testing.T) {
	mockRepo := new(mockRepo)
	tokens := NewTokenManager("access-secret", "refresh-secret", 15*time.Minute, 7*24*time.Hour)
	svc := NewService(mockRepo, tokens)

	now := time.Now()
	expectedUser := User{
		ID:             "user-1",
		Email:          "test@example.com",
		Username:       "testuser",
		DisplayName:    "Test User",
		PasswordHash:   "hashed",
		Role:           "user",
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u User) bool {
		return u.Email == "test@example.com" && u.Username == "testuser"
	})).Return(expectedUser, nil)

	mockRepo.On("CreateSession", mock.Anything, "user-1", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	resp, err := svc.Register(context.Background(), RegisterRequest{
		Email:       "test@example.com",
		Username:    "testuser",
		DisplayName: "Test User",
		Password:    "password123",
	}, "test-agent", "127.0.0.1")

	assert.NoError(t, err)
	assert.Equal(t, "user-1", resp.User.ID)
	assert.Equal(t, "test@example.com", resp.User.Email)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	mockRepo.AssertExpectations(t)
}

func TestService_Register_DuplicateEmail(t *testing.T) {
	mockRepo := new(mockRepo)
	tokens := NewTokenManager("secret", "secret2", 15*time.Minute, 7*24*time.Hour)
	svc := NewService(mockRepo, tokens)

	mockRepo.On("CreateUser", mock.Anything, mock.Anything).
		Return(User{}, apperrors.Conflict("email or username already exists", nil))

	_, err := svc.Register(context.Background(), RegisterRequest{
		Email:       "dup@example.com",
		Username:    "dupuser",
		DisplayName: "Dup User",
		Password:    "password123",
	}, "", "")

	assert.Error(t, err)
	assert.True(t, apperrors.IsConflict(err))
	mockRepo.AssertExpectations(t)
}

func TestService_Login_Success(t *testing.T) {
	mockRepo := new(mockRepo)
	tokens := NewTokenManager("secret-a", "secret-r", 15*time.Minute, 7*24*time.Hour)
	svc := NewService(mockRepo, tokens)

	hash, _ := HashPassword("correct-password")
	now := time.Now()
	expectedUser := User{
		ID:           "user-2",
		Email:        "login@example.com",
		Username:     "loginuser",
		DisplayName:  "Login User",
		PasswordHash: hash,
		Role:         "user",
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	mockRepo.On("FindUserByEmailOrUsername", mock.Anything, "login@example.com").Return(expectedUser, nil)
	mockRepo.On("CreateSession", mock.Anything, "user-2", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	resp, err := svc.Login(context.Background(), LoginRequest{
		Email:    "login@example.com",
		Password: "correct-password",
	}, "", "")

	assert.NoError(t, err)
	assert.Equal(t, "login@example.com", resp.User.Email)
	assert.NotEmpty(t, resp.AccessToken)
	mockRepo.AssertExpectations(t)
}

func TestService_Login_WrongPassword(t *testing.T) {
	mockRepo := new(mockRepo)
	tokens := NewTokenManager("secret-a", "secret-r", 15*time.Minute, 7*24*time.Hour)
	svc := NewService(mockRepo, tokens)

	hash, _ := HashPassword("correct-password")
	now := time.Now()
	expectedUser := User{
		ID:           "user-2",
		Email:        "login@example.com",
		Username:     "loginuser",
		DisplayName:  "Login User",
		PasswordHash: hash,
		Role:         "user",
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	mockRepo.On("FindUserByEmailOrUsername", mock.Anything, "login@example.com").Return(expectedUser, nil)

	_, err := svc.Login(context.Background(), LoginRequest{
		Email:    "login@example.com",
		Password: "wrong-password",
	}, "", "")

	assert.Error(t, err)
	assert.True(t, apperrors.IsUnauthorized(err))
	mockRepo.AssertExpectations(t)
}

func TestService_Login_InactiveUser(t *testing.T) {
	mockRepo := new(mockRepo)
	tokens := NewTokenManager("secret-a", "secret-r", 15*time.Minute, 7*24*time.Hour)
	svc := NewService(mockRepo, tokens)

	now := time.Now()
	expectedUser := User{
		ID:           "user-3",
		Email:        "inactive@example.com",
		Username:     "inactive",
		DisplayName:  "Inactive User",
		PasswordHash: "hash",
		Role:         "user",
		IsActive:     false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	mockRepo.On("FindUserByEmailOrUsername", mock.Anything, "inactive@example.com").Return(expectedUser, nil)

	_, err := svc.Login(context.Background(), LoginRequest{
		Email:    "inactive@example.com",
		Password: "any-password",
	}, "", "")

	assert.Error(t, err)
	assert.True(t, apperrors.IsForbidden(err))
	mockRepo.AssertExpectations(t)
}

func TestService_Me_Success(t *testing.T) {
	mockRepo := new(mockRepo)
	tokens := NewTokenManager("secret-a", "secret-r", 15*time.Minute, 7*24*time.Hour)
	svc := NewService(mockRepo, tokens)

	now := time.Now()
	expectedUser := User{
		ID:          "user-1",
		Email:       "me@example.com",
		Username:    "meuser",
		DisplayName: "Me User",
		Role:        "user",
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockRepo.On("FindUserByID", mock.Anything, "user-1").Return(expectedUser, nil)

	resp, err := svc.Me(context.Background(), "user-1")

	assert.NoError(t, err)
	assert.Equal(t, "me@example.com", resp.User.Email)
	mockRepo.AssertExpectations(t)
}

func TestService_AdminListUsers_Success(t *testing.T) {
	mockRepo := new(mockRepo)
	tokens := NewTokenManager("secret-a", "secret-r", 15*time.Minute, 7*24*time.Hour)
	svc := NewService(mockRepo, tokens)

	now := time.Now()
	users := []User{
		{ID: "u1", Email: "a@test.com", Username: "user1", DisplayName: "User One", Role: "user", IsActive: true, CreatedAt: now, UpdatedAt: now},
		{ID: "u2", Email: "b@test.com", Username: "user2", DisplayName: "User Two", Role: "admin", IsActive: true, CreatedAt: now, UpdatedAt: now},
	}

	params := ListUsersParams{Page: 1, PageSize: 20, SortBy: "created_at", SortOrder: "desc"}
	mockRepo.On("ListUsers", mock.Anything, params).Return(users, 2, nil)

	resp, err := svc.AdminListUsers(context.Background(), params)

	assert.NoError(t, err)
	assert.Equal(t, 2, resp.Total)
	assert.Len(t, resp.Items, 2)
	assert.Equal(t, "User One", resp.Items[0].DisplayName)
	assert.Equal(t, "admin", resp.Items[1].Role)
	mockRepo.AssertExpectations(t)
}

func TestService_AdminUpdateUser_Success(t *testing.T) {
	mockRepo := new(mockRepo)
	tokens := NewTokenManager("secret-a", "secret-r", 15*time.Minute, 7*24*time.Hour)
	svc := NewService(mockRepo, tokens)

	role := "admin"
	active := true

	now := time.Now()
	updatedUser := User{
		ID: "u1", Email: "a@test.com", Username: "user1", DisplayName: "User One",
		Role: "admin", IsActive: true, CreatedAt: now, UpdatedAt: now,
	}

	mockRepo.On("UpdateUser", mock.Anything, "u1", mock.MatchedBy(func(m map[string]any) bool {
		return m["role"] == "admin" && m["is_active"] == true
	})).Return(nil)
	mockRepo.On("FindUserByID", mock.Anything, "u1").Return(updatedUser, nil)

	item, err := svc.AdminUpdateUser(context.Background(), "u1", AdminUpdateUserRequest{
		Role:     &role,
		IsActive: &active,
	})

	assert.NoError(t, err)
	assert.Equal(t, "admin", item.Role)
	assert.True(t, item.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestService_AdminDeleteUser_Success(t *testing.T) {
	mockRepo := new(mockRepo)
	tokens := NewTokenManager("secret-a", "secret-r", 15*time.Minute, 7*24*time.Hour)
	svc := NewService(mockRepo, tokens)

	mockRepo.On("DeleteUser", mock.Anything, "u1").Return(nil)

	err := svc.AdminDeleteUser(context.Background(), "u1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
