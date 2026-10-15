package ingestion

import (
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(rg *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc, adminMW gin.HandlerFunc) {
	admin := rg.Group("/admin/ingestion")
	admin.Use(authMW, adminMW)
	{
		admin.POST("/upload", h.Upload)
		admin.GET("/drafts", h.ListDrafts)
		admin.GET("/drafts/:id", h.GetDraft)
		admin.PATCH("/drafts/:id/metadata", h.UpdateDraftMetadata)
		admin.DELETE("/drafts/:id", h.DeleteDraft)
		admin.POST("/drafts/:id/enrich", h.EnrichDraft)
		admin.GET("/drafts/:id/suggestions", h.GetDraftSuggestions)
		admin.PATCH("/drafts/:id/final-metadata", h.SaveFinalMetadata)
		admin.POST("/drafts/:id/reject", h.RejectDraft)
		admin.POST("/drafts/:id/finalize", h.FinalizeDraft)
		admin.POST("/drafts/:id/images/:entity", h.UploadDraftImage)
		admin.GET("/stats", h.GetStats)
		admin.POST("/cleanup", h.TriggerCleanup)
		admin.GET("/config", h.GetConfig)
	}

	search := rg.Group("/admin/catalog")
	search.Use(authMW, adminMW)
	{
		search.GET("/artists/search", h.SearchArtists)
		search.GET("/albums/search", h.SearchAlbums)
	}
}