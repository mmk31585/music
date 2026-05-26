package response

import (
	"net/http"

	apperrors "music/internal/common/errors"
)

func ValidationError(w http.ResponseWriter, details interface{}) {
	Error(w, apperrors.New(
		http.StatusUnprocessableEntity,
		"VALIDATION_ERROR",
		"validation failed",
		details,
	))
}
