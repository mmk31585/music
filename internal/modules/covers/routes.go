package covers

import "github.com/gin-gonic/gin"

// RegisterInternalRoutes registers internal service-to-service endpoints
// for cover optimization callbacks, protected by HMAC signature verification.
// Mounted under /api/v1 so the Python ML service can reach them via
// {callback_base_url}/internal/v1/covers/callback.
func RegisterInternalRoutes(api *gin.RouterGroup, h *Handler, hmacMW gin.HandlerFunc) {
	internal := api.Group("/internal/v1/covers")
	internal.Use(hmacMW)
	{
		internal.POST("/callback", h.HandleCoverCallback)
	}
}
