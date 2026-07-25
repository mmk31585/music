package upload

import (
	"github.com/gin-gonic/gin"

	"music/internal/common/response"
	"music/internal/modules/auth"
	"music/internal/modules/permissions"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, tokenManager *auth.TokenManager, permSvc permissions.ServiceInterface) {
	api := rg.Group("")
	api.Use(auth.AuthMiddleware(tokenManager))

	// Listener uploads — mine must come before :draftId to avoid matching "mine" as a draft ID
	api.POST("/uploads/drafts", response.Wrap(handler.CreateDraft))
	api.POST("/uploads/drafts/:draftId/file", response.Wrap(handler.UploadDraftFile))
	api.GET("/uploads/drafts/mine", response.Wrap(handler.ListMyDrafts))
	api.GET("/uploads/drafts/:draftId", response.Wrap(handler.GetDraft))
	api.GET("/uploads/slots", response.Wrap(handler.GetUploadSlots))

	// Co-uploaders (any authenticated user can view; manage via track ownership)
	api.GET("/tracks/:trackId/co-uploaders", response.Wrap(handler.GetCoUploaders))

	// Admin/Moderator: review pending drafts
	reviewGroup := api.Group("/uploads/review")
	reviewGroup.Use(permissions.RequireAnyPermission(permSvc, "review_uploads", "manage_users"))
	{
		reviewGroup.GET("", response.Wrap(handler.ListPendingDrafts))
		reviewGroup.POST("/:draftId", response.Wrap(handler.ReviewDraft))
	}

	// Admin: manage co-uploaders
	adminGroup := api.Group("/uploads/co-uploaders")
	adminGroup.Use(permissions.RequireAnyPermission(permSvc, "manage_users"))
	{
		adminGroup.POST("/:trackId", response.Wrap(handler.AddCoUploader))
		adminGroup.DELETE("/:trackId/:userId", response.Wrap(handler.RemoveCoUploader))
	}
}
