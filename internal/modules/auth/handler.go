package auth

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	apperrors "music/internal/common/errors"
	"music/internal/common/response"
	"music/internal/common/validator"
)

// ServiceInterface defines the service methods needed by the HTTP handler.
type ServiceInterface interface {
	Register(ctx context.Context, req RegisterRequest, userAgent, ipAddress string) (AuthResponse, error)
	Login(ctx context.Context, req LoginRequest, userAgent, ipAddress string) (AuthResponse, error)
	Refresh(ctx context.Context, req RefreshRequest, userAgent, ipAddress string) (AuthResponse, error)
	Logout(ctx context.Context, req LogoutRequest) error
	Me(ctx context.Context, userID string) (MeResponse, error)
	AdminListUsers(ctx context.Context, params ListUsersParams) (AdminListUsersResponse, error)
	AdminGetUser(ctx context.Context, id string) (AdminUserItem, error)
	AdminUpdateUser(ctx context.Context, id string, req AdminUpdateUserRequest) (AdminUserItem, error)
	AdminDeleteUser(ctx context.Context, id string) error
	UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) error
	ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error
	GetPublicProfile(ctx context.Context, id string) (*User, error)
}

type Handler struct {
	service   ServiceInterface
	validator *validator.Validator
	logger    *zap.Logger
}

func NewHandler(service ServiceInterface, validator *validator.Validator) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
		logger:    zap.L(),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account and returns authentication/session information.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration payload"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.service.Register(c.Request.Context(), req, c.Request.UserAgent(), clientIP(c))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "registered successfully", result)
}

// Login godoc
// @Summary Log in user
// @Description Authenticates a user and returns access and refresh token information.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login payload"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("login bind error", zap.Error(err))
		// Return 400 BadRequest for malformed/empty bodies, not 500
		response.Error(c, apperrors.BadRequest("invalid request body: "+err.Error(), nil))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.Error(c, apperrors.BadRequest("validation failed: "+err.Error(), nil))
		return
	}

	result, err := h.service.Login(c.Request.Context(), req, c.Request.UserAgent(), clientIP(c))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "logged in successfully", result)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Refreshes authentication tokens using a valid refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token payload"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", nil))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.service.Refresh(c.Request.Context(), req, c.Request.UserAgent(), clientIP(c))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "token refreshed successfully", result)
}

// Logout godoc
// @Summary Log out user
// @Description Invalidates the provided session or refresh token and logs the user out.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LogoutRequest true "Logout payload"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.Error(c, err)
		return
	}

	if err := h.service.Logout(c.Request.Context(), req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, http.StatusOK, "logged out successfully", nil)
}

// Me godoc
// @Summary Get current user
// @Description Returns the currently authenticated user's profile.
// @Tags auth
// @Produce json
// @Security Bearer
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	userID := UserIDFromContext(c) // uses the Gin context getter from middleware.go
	if userID == "" {
		response.Error(c, ErrUnauthorized())
		return
	}

	result, err := h.service.Me(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "current user", result)
}

// AdminListUsers godoc
// @Summary List all users (admin)
// @Description Returns a paginated list of all users.
// @Tags admin
// @Produce json
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(20)
// @Param sort_by query string false "Sort field" Enums(created_at,email,username,display_name,role,is_active)
// @Param sort_order query string false "Sort order" Enums(asc,desc)
// @Param search query string false "Search query"
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/users [get]
func (h *Handler) AdminListUsers(c *gin.Context) {
	params := ListUsersParams{
		Page:      1,
		PageSize:  20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Page = p
		}
	}
	if pageSize := c.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			params.PageSize = ps
		}
	}
	if sortBy := c.Query("sort_by"); sortBy != "" {
		params.SortBy = sortBy
	}
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		params.SortOrder = sortOrder
	}
	if search := c.Query("search"); search != "" {
		params.Search = search
	}

	result, err := h.service.AdminListUsers(c.Request.Context(), params)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "users retrieved", result)
}

// AdminGetUser godoc
// @Summary Get user by ID (admin)
// @Description Returns a single user's details including sensitive fields.
// @Tags admin
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/users/{id} [get]
func (h *Handler) AdminGetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, apperrors.BadRequest("user id is required", nil))
		return
	}

	result, err := h.service.AdminGetUser(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "user retrieved", result)
}

// AdminUpdateUser godoc
// @Summary Update user (admin)
// @Description Updates a user's role, active status, or email verification status.
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Param request body AdminUpdateUserRequest true "Update payload"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/users/{id} [put]
func (h *Handler) AdminUpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, apperrors.BadRequest("user id is required", nil))
		return
	}

	var req AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.service.AdminUpdateUser(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "user updated", result)
}

// AdminDeleteUser godoc
// @Summary Delete user (admin)
// @Description Permanently deletes a user account.
// @Tags admin
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /admin/users/{id} [delete]
func (h *Handler) AdminDeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, apperrors.BadRequest("user id is required", nil))
		return
	}

	if err := h.service.AdminDeleteUser(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, http.StatusOK, "user deleted", nil)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID := UserIDFromContext(c)
	if userID == "" {
		response.Error(c, apperrors.Unauthorized("authentication required", nil))
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("update profile bind error", zap.Error(err), zap.Any("body", c.Request.Body))
		response.Error(c, apperrors.BadRequest("invalid request body: "+err.Error(), nil))
		return
	}

	if err := h.service.UpdateProfile(c.Request.Context(), userID, req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.service.Me(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "profile updated", result)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	userID := UserIDFromContext(c)
	if userID == "" {
		response.Error(c, apperrors.Unauthorized("authentication required", nil))
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("invalid request body", nil))
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userID, req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, http.StatusOK, "password changed", nil)
}

func (h *Handler) GetPublicProfile(c *gin.Context) {
	user, err := h.service.GetPublicProfile(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "user not found"})
		return
	}

	response.Success(c, http.StatusOK, "ok", gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"full_name":  user.DisplayName,
		"avatar_url": user.AvatarURL,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// clientIP extracts the client IP from a Gin context
func clientIP(c *gin.Context) string {
	// Try X-Forwarded-For
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	// Try X-Real-IP
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return realIP
	}
	// Fallback to RemoteAddr
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return host
}
