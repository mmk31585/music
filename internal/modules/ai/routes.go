package ai

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
	aiGroup := api.Group("/ai")
	aiGroup.Use(authMW)
	{
		aiGroup.POST("/embeddings", h.GenerateEmbedding)
		aiGroup.POST("/moods", h.AnalyzeMood)
		aiGroup.GET("/moods/:trackId", h.GetMood)
		aiGroup.POST("/playlists/generate", h.GeneratePlaylist)
		aiGroup.GET("/similar/mood/:trackId", h.SimilarByMood)
		aiGroup.GET("/similar/embedding/:trackId", h.SimilarByEmbedding)
	}
}
