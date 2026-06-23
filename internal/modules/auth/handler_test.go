package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"music/internal/common/validator"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Register(ctx context.Context, req RegisterRequest, userAgent, ip string) (AuthResponse, error) {
	args := m.Called(ctx, req, userAgent, ip)
	return args.Get(0).(AuthResponse), args.Error(1)
}

func (m *mockAuthService) Login(ctx context.Context, req LoginRequest, userAgent, ip string) (AuthResponse, error) {
	args := m.Called(ctx, req, userAgent, ip)
	return args.Get(0).(AuthResponse), args.Error(1)
}

func (m *mockAuthService) Refresh(ctx context.Context, req RefreshRequest, userAgent, ip string) (AuthResponse, error) {
	args := m.Called(ctx, req, userAgent, ip)
	return args.Get(0).(AuthResponse), args.Error(1)
}

func (m *mockAuthService) Logout(ctx context.Context, req LogoutRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockAuthService) Me(ctx context.Context, userID string) (MeResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(MeResponse), args.Error(1)
}

// Stub methods for unused interface methods
func (m *mockAuthService) AdminListUsers(ctx context.Context, params ListUsersParams) (AdminListUsersResponse, error) {
	return AdminListUsersResponse{}, nil
}
func (m *mockAuthService) AdminGetUser(ctx context.Context, id string) (AdminUserItem, error) {
	return AdminUserItem{}, nil
}
func (m *mockAuthService) AdminUpdateUser(ctx context.Context, id string, req AdminUpdateUserRequest) (AdminUserItem, error) {
	return AdminUserItem{}, nil
}
func (m *mockAuthService) AdminDeleteUser(ctx context.Context, id string) error { return nil }
func (m *mockAuthService) UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) error {
	return nil
}
func (m *mockAuthService) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	return nil
}
func (m *mockAuthService) GetPublicProfile(ctx context.Context, id string) (*User, error) {
	return nil, nil
}

func setupAuthTest() (*Handler, *mockAuthService, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockAuthService)
	v := validator.New()
	h := NewHandler(mockSvc, v)
	h.logger = zap.NewNop()

	r := gin.New()
	rg := r.Group("/api/v1/auth")
	RegisterRoutes(rg, h, func(c *gin.Context) { c.Next() })
	return h, mockSvc, r
}

func TestRegister_Success(t *testing.T) {
	_, mockSvc, r := setupAuthTest()

	reqBody := RegisterRequest{
		Email:       "test@example.com",
		Username:    "testuser",
		DisplayName: "Test User",
		Password:    "SecurePass123!",
	}

	mockSvc.On("Register", mock.Anything, reqBody, mock.Anything, mock.Anything).
		Return(AuthResponse{
			AccessToken:  "access_token_123",
			RefreshToken: "refresh_token_456",
			User: AuthUser{
				ID:       "user_1",
				Email:    "test@example.com",
				Username: "testuser",
				Role:     "user",
			},
		}, nil)

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp["success"].(bool))
	mockSvc.AssertExpectations(t)
}

func TestRegister_ValidationError(t *testing.T) {
	_, _, r := setupAuthTest()

	reqBody := RegisterRequest{
		Email:    "invalid",
		Password: "123",
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_Success(t *testing.T) {
	_, mockSvc, r := setupAuthTest()

	reqBody := LoginRequest{
		Email:    "test@example.com",
		Password: "SecurePass123!",
	}

	mockSvc.On("Login", mock.Anything, reqBody, mock.Anything, mock.Anything).
		Return(AuthResponse{
			AccessToken:  "access_token_123",
			RefreshToken: "refresh_token_456",
			User: AuthUser{
				ID:       "user_1",
				Email:    "test@example.com",
				Username: "testuser",
				Role:     "user",
			},
		}, nil)

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	_, mockSvc, r := setupAuthTest()

	reqBody := LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpass",
	}

	mockSvc.On("Login", mock.Anything, reqBody, mock.Anything, mock.Anything).
		Return(AuthResponse{}, errors.New("invalid credentials"))

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestRefresh_Success(t *testing.T) {
	_, mockSvc, r := setupAuthTest()

	reqBody := RefreshRequest{RefreshToken: "valid_refresh_token"}

	mockSvc.On("Refresh", mock.Anything, RefreshRequest{RefreshToken: "valid_refresh_token"}, mock.Anything, mock.Anything).
		Return(AuthResponse{
			AccessToken:  "new_access_token",
			RefreshToken: "new_refresh_token",
		}, nil)

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestRefresh_InvalidToken(t *testing.T) {
	_, mockSvc, r := setupAuthTest()

	reqBody := RefreshRequest{RefreshToken: "expired_token"}

	mockSvc.On("Refresh", mock.Anything, RefreshRequest{RefreshToken: "expired_token"}, mock.Anything, mock.Anything).
		Return(AuthResponse{}, errors.New("invalid or expired refresh token"))

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockSvc.AssertExpectations(t)
}
