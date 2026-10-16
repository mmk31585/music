package video

import (
	"github.com/gin-gonic/gin"
)

// RegisterInternalRoutes registers internal service-to-service endpoints
// protected by HMAC signature verification (not user JWT auth).
// Mounted under /api/v1 so the Python ML service can reach them via
// {callback_base_url}/internal/v1/video/callback.
func RegisterInternalRoutes(api *gin.RouterGroup, handler *Handler, hmacMW gin.HandlerFunc) {
	internal := api.Group("/internal/v1/video")
	internal.Use(hmacMW)
	{
		internal.POST("/callback", handler.HandleVideoCallback)
	}
}

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc, adminMW gin.HandlerFunc) {
	// --- Admin routes (auth + admin role required) ---
	admin := rg.Group("/admin")
	admin.Use(authMW, adminMW)
	{
		admin.GET("/videos", handler.AdminListVideos)
		admin.POST("/tracks/:trackId/mv", handler.UploadOfficialMV)
		admin.PUT("/videos/:id", handler.AdminUpdateVideo)
		admin.DELETE("/videos/:id", handler.DeleteVideo)
		admin.POST("/videos/:id/approve", handler.ApproveVideo)
		admin.GET("/videos/:id/comments", handler.AdminListVideoComments)
		admin.PUT("/comments/:id", handler.AdminUpdateComment)
		admin.DELETE("/comments/:id", handler.AdminDeleteComment)
	}

	// --- Authenticated user routes ---
	authGroup := rg.Group("/videos")
	authGroup.Use(authMW)
	{
		authGroup.POST("", handler.CreateUserEdit)
		authGroup.GET("", handler.GetUserVideos)
		authGroup.DELETE("/:id", handler.DeleteOwnVideo)
	}

	// --- Public routes ---
	rg.GET("/music-video/:id", handler.GetVideo)
	rg.GET("/music-video/:id/comments", handler.ListComments)
	rg.GET("/music-videos/explore", handler.ListExplore)
	rg.GET("/music-videos/by-track/:trackId", handler.ListByTrack)
}