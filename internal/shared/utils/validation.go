package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "music/internal/common/errors"
)

// ValidationError sends a 422 Unprocessable Entity response for validation failures.
func ValidationError(c *gin.Context, details interface{}) {
	Error(c, apperrors.New(
		http.StatusUnprocessableEntity,
		"VALIDATION_ERROR",
		"validation failed",
		details,
	))
}
