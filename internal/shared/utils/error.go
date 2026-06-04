package utils

import (
	"github.com/gin-gonic/gin"
	apperrors "music/internal/contracts/errors"
)

// ErrorResponse represents the standard error payload structure.
type ErrorResponse struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error sends a standardized error response and aborts the Gin context.
// Accepts either *apperrors.AppError or any error (converted via ToAppError).
func Error(c *gin.Context, err error) {
	appErr := apperrors.ToAppError(err)
	c.AbortWithStatusJSON(appErr.StatusCode, ErrorResponse{
		Success: false,
		Code:    appErr.Code,
		Message: appErr.Message,
		Details: appErr.Details,
	})
}

// HandlerFunc defines a typed handler signature that returns an error.
// When used with Wrap, errors are automatically sent via Error().
type HandlerFunc func(c *gin.Context) error

// Wrap converts a HandlerFunc into a standard gin.HandlerFunc with automatic error handling.
func Wrap(fn HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c); err != nil {
			Error(c, err) // automatically aborts context
		}
	}
}
