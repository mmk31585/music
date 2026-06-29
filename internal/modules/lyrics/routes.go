package lyrics

import (
	"github.com/gin-gonic/gin"

	"music/internal/modules/auth"
)

func RegisterRoutes(api *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
	// Public routes
	api.GET("/tracks/:trackId/lyrics", h.GetTrackLyrics)
	api.GET("/tracks/:trackId/lyrics/all", h.GetLyricsByTrackID)

	// Admin routes (protected)
	admin := api.Group("/admin")
	admin.Use(authMW, auth.RequireRole("admin"))
	{
		admin.POST("/lyrics", h.CreateLyrics)
		admin.POST("/lyrics/fetch/:trackId", h.FetchFromLRC)
		admin.POST("/lyrics/fetch-or-generate/:trackId", h.FetchOrGenerateLyrics)
		admin.POST("/lyrics/sync", h.SyncWithAI)
		admin.POST("/lyrics/review", h.ReviewWithAI)
		admin.GET("/lyrics/ai-status/:trackId", h.AILyricsStatus)
		admin.PUT("/lyrics/:id", h.UpdateLyrics)
		admin.DELETE("/lyrics/:id", h.DeleteLyrics)
	}
}

// RegisterInternalRoutes registers internal service-to-service endpoints
// protected by HMAC signature verification (not user JWT auth).
// These are mounted under the /api/v1 prefix so the Python ML service
// can reach them via “{callback_base_url}/internal/v1/lyrics/callback“
// where callback_base_url typically includes /api/v1.
func RegisterInternalRoutes(api *gin.RouterGroup, h *Handler, hmacMW gin.HandlerFunc) {
	internal := api.Group("/internal/v1/lyrics")
	internal.Use(hmacMW)
	{
		internal.POST("/callback", h.HandleLyricsCallback)
	}
}
