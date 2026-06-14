package search

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, mw ...gin.HandlerFunc) {
	search := rg.Group("/search")
	if len(mw) > 0 {
		search.Use(mw...)
	}
	{
		search.GET("", handler.Search)
	}
}
