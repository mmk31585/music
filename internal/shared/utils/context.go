package utils

import (
	"net/http"

	"music/internal/shared/middleware"
)

func RequestID(r *http.Request) string {
	value := r.Context().Value(middleware.RequestIDKey)
	if value == nil {
		return ""
	}

	requestID, ok := value.(string)
	if !ok {
		return ""
	}

	return requestID
}
