package tips

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	r := rg.Group("/tips")
	{
		r.POST("", authMW, handler.Create)
		r.GET("/callback/:id", handler.VerifyCallback)
		r.GET("/sent", authMW, handler.ListSent)
		r.GET("/received/:artist_id", authMW, handler.ListReceived)
	}
}
