package auth

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"music/internal/common/response"
	"music/internal/common/validator"
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

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinError(c, err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.GinError(c, err)
		return
	}

	result, err := h.service.Register(c.Request.Context(), req, c.Request.UserAgent(), clientIP(c))
	if err != nil {
		response.GinError(c, err)
		return
	}

	response.GinSuccess(c, http.StatusCreated, "registered successfully", result)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		// Log the actual error before responding
		fmt.Printf("Login bind error: %v\n", err) // or use zap
		response.GinError(c, err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.GinError(c, err)
		return
	}

	result, err := h.service.Login(c.Request.Context(), req, c.Request.UserAgent(), clientIP(c))
	if err != nil {
		response.GinError(c, err)
		return
	}

	response.GinSuccess(c, http.StatusOK, "logged in successfully", result)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinError(c, err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.GinError(c, err)
		return
	}

	result, err := h.service.Refresh(c.Request.Context(), req, c.Request.UserAgent(), clientIP(c))
	if err != nil {
		response.GinError(c, err)
		return
	}

	response.GinSuccess(c, http.StatusOK, "token refreshed successfully", result)
}

func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinError(c, err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.GinError(c, err)
		return
	}

	if err := h.service.Logout(c.Request.Context(), req); err != nil {
		response.GinError(c, err)
		return
	}

	response.GinSuccess(c, http.StatusOK, "logged out successfully", nil)
}

func (h *Handler) Me(c *gin.Context) {
	userID := UserIDFromContext(c) // uses the Gin context getter from middleware.go
	if userID == "" {
		response.GinError(c, ErrUnauthorized())
		return
	}

	result, err := h.service.Me(c.Request.Context(), userID)
	if err != nil {
		response.GinError(c, err)
		return
	}

	response.GinSuccess(c, http.StatusOK, "current user", result)
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
