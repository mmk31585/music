package request

import (
	"net/http"

	"music/internal/common/middleware"
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
