package auth

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"music/internal/shared/utils"
	"music/internal/platform/validation"
)

type Handler struct {
	service   *Service
	validator *validator.Validator
}

func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
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
		// Log the actual error before responding
		fmt.Printf("Login bind error: %v\n", err) // or use zap
		response.Error(c, err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.Error(c, err)
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
		response.Error(c, err)
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
