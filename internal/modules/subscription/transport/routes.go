package subscription

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, authMW gin.HandlerFunc) {
	subscriptions := rg.Group("/subscription")
	{
		// Public
		subscriptions.GET("/plans", handler.ListPlans)

		// Auth
		protected := subscriptions.Group("")
		protected.Use(authMW)
		{
			protected.GET("/me", handler.CurrentSubscription)
			protected.GET("/history", handler.ListSubscriptions)
			protected.POST("/checkout", handler.Checkout)
			protected.POST("/cancel", handler.CancelCurrentSubscription)
			protected.GET("/payments", handler.ListPayments)
		}
	}
}
