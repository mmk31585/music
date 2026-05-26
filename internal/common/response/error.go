package response

import (
	"encoding/json"
	"net/http"

	apperrors "music/internal/common/errors"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error writes an error response using http.ResponseWriter (for non-Gin handlers)
func Error(w http.ResponseWriter, err error) {
	appErr := apperrors.ToAppError(err)
	writeError(w, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
}

// AppError is an alias for Error (kept for compatibility)
func AppError(w http.ResponseWriter, err error) {
	Error(w, err)
}

// GinError writes an error response using Gin's context
func GinError(c *gin.Context, err error) {
	appErr := apperrors.ToAppError(err)
	writeGinError(c, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
}

// GinAppError is an alias for GinError
func GinAppError(c *gin.Context, err error) {
	GinError(c, err)
}

// writeError is the internal writer for standard http.ResponseWriter
func writeError(w http.ResponseWriter, statusCode int, code string, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Code:    code,
		Message: message,
		Details: details,
	})
}

// writeGinError is the internal writer for gin.Context
func writeGinError(c *gin.Context, statusCode int, code string, message string, details interface{}) {
	c.Header("Content-Type", "application/json")
	c.Status(statusCode)
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Code:    code,
		Message: message,
		Details: details,
	})
}
