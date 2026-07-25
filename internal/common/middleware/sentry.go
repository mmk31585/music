package middleware

import (
	"net/http"
	"time"

	"music/internal/platform/web"

	sentrygo "github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func Sentry() gin.HandlerFunc {
	_ = sentrygin.New(sentrygin.Options{})

	return func(c *gin.Context) {
		hub := sentrygo.GetHubFromContext(c.Request.Context())
		if hub == nil {
			hub = sentrygo.CurrentHub().Clone()
		}

		hub.Scope().SetRequest(c.Request)

		if userID, ok := web.GetUserIDString(c); ok {
			hub.Scope().SetUser(sentrygo.User{ID: userID})
		}

		hub.Scope().AddBreadcrumb(&sentrygo.Breadcrumb{
			Category:  "http",
			Message:   c.Request.Method + " " + c.Request.URL.String(),
			Level:     sentrygo.LevelInfo,
			Timestamp: time.Now(),
		}, 100)

		c.Request = c.Request.WithContext(sentrygo.SetHubOnContext(c.Request.Context(), hub))

		defer func() {
			if err := recover(); err != nil {
				hub.RecoverWithContext(c.Request.Context(), err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "internal server error",
				})
			}
		}()

		c.Next()

		for _, err := range c.Errors {
			hub.Scope().SetContext("http", map[string]interface{}{
				"method": c.Request.Method,
				"url":    c.Request.URL.String(),
			})
			hub.CaptureException(err.Err)
		}
	}
}