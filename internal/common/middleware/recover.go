package middleware

import (
	"net/http"

	"go.uber.org/zap"
	apperrors "music/internal/common/errors"
	"music/internal/common/response"
)

func Recover(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic recovered",
						zap.Any("panic", rec),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
					)

					response.AppError(w, apperrors.Internal("internal server error", nil))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
