package video

import (
	"github.com/gin-gonic/gin"

	"music/internal/modules/auth"
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

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	// --- Admin routes (auth + admin role required) ---
	admin := rg.Group("/admin")
	admin.Use(authMW, auth.RequireRole("admin"))
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
		authGroup.DELETE("/:id", handler.DeleteOwnVideo)
		authGroup.POST("/:id/comments", handler.CreateComment)
		authGroup.POST("/:id/like", handler.LikeVideo)
		authGroup.DELETE("/:id/like", handler.UnlikeVideo)
		authGroup.POST("/:id/view", handler.ViewVideo)
	}

	// --- Public video routes ---
	rg.GET("/videos/:id", handler.GetVideo)
	rg.GET("/videos/:id/stream", handler.StreamVideo)
	rg.GET("/videos/:id/comments", handler.ListComments)
	rg.GET("/tracks/:trackId/videos", handler.ListByTrack)
	rg.GET("/explore/videos", handler.ListExplore)

	// --- User videos ---
	rg.GET("/users/:id/videos", handler.GetUserVideos)

	// --- Track like visibility (authenticated) ---
	rg.PUT("/tracks/:trackId/like/visibility", authMW, handler.SetTrackLikeVisibility)
	rg.GET("/users/:id/liked-tracks", handler.GetPublicLikedTracks)

	// --- Music status ---
	rg.PUT("/users/me/music-status", authMW, handler.SetMusicStatus)
	rg.DELETE("/users/me/music-status", authMW, handler.ClearMusicStatus)
	rg.GET("/users/:id/music-status", handler.GetMusicStatus)
}
