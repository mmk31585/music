package search

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	search := rg.Group("/search")
	{
		search.GET("", handler.Search)
	}
}
