package app

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"music/internal/common/middleware"
	"music/internal/modules/analytics"
	"music/internal/modules/auth"
	"music/internal/modules/cacheadmin"
	"music/internal/modules/catalog"
	"music/internal/modules/cover"
	"music/internal/modules/dashboard"
	"music/internal/modules/features"
	"music/internal/modules/follow"
	"music/internal/modules/health"
	"music/internal/modules/history"
	"music/internal/modules/ingestion"
	"music/internal/modules/library"
	"music/internal/modules/lyrics"
	"music/internal/modules/media"
	"music/internal/modules/moderation"
	"music/internal/modules/notification"
	"music/internal/modules/permissions"
	"music/internal/modules/playlist"
	"music/internal/modules/player"
	"music/internal/modules/queue"
	"music/internal/modules/reactions"
	"music/internal/modules/recommendation"
	"music/internal/modules/reputation"
	"music/internal/modules/social"
	"music/internal/modules/subscription"
	"music/internal/modules/upload"
	"music/internal/modules/video"
)

// RegisterRoutes registers all API routes on the given router group.
func RegisterRoutes(r *gin.RouterGroup, a *App) {
	api := r.Group("/api/v1")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api.GET("/metrics", gin.WrapH(promhttp.Handler()))
	api.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	c := NewContainer(a)

	health.RegisterRoutes(api, c.HealthHandler)
	features.RegisterRoutes(api, c.FeaturesHandler)
	dashboard.RegisterRoutes(api, c.DashboardHandler)
	auth.RegisterRoutes(api, c.AuthHandler, c.AuthMW, c.AdminMW, c.RDB)
	permissions.RegisterRoutes(api, c.PermissionHandler, c.AuthMW, c.AdminMW, c.PermissionService, a.DB)
	reputation.RegisterRoutes(api, c.ReputationHandler, c.AuthMW, c.TokenManager)
	upload.RegisterRoutes(api, c.UploadHandler, c.AuthMW, c.PermissionService)
	media.RegisterAdminRoutes(api, c.MediaHandler, c.AuthMW, c.AdminMW)

	ingestion.RegisterAdminRoutes(api, c.IngestionHandler, c.AuthMW, c.AdminMW)

	catalogHandlers := catalog.Handlers{
		Artist: c.ArtistHandler,
		Album:  c.AlbumHandler,
		Track:  c.TrackHandler,
		Genre:  c.GenreHandler,
		Stats:  c.StatsHandler,
		Enrich: c.EnrichHandler,
	}

	catalog.RegisterPublicRoutes(api, catalogHandlers, c.OptionalAuthMW, middleware.RateLimitOptional(c.RDB, 60))
	catalog.RegisterAdminRoutes(api, catalogHandlers, c.AuthMW, c.AdminMW)

	lyrics.RegisterRoutes(api, c.LyricsHandler, c.AuthMW, c.AdminMW)
	lyrics.RegisterInternalRoutes(api, c.LyricsHandler, c.MLServiceHMACMW)

	covers.RegisterInternalRoutes(api, c.CoverHandler, c.MLServiceHMACMW)

	playlist.RegisterRoutes(api, c.PlaylistHandler, c.OptionalAuthMW)
	library.RegisterRoutes(api, c.LibraryHandler, c.AuthMW)
	queue.RegisterRoutes(api, c.QueueHandler, c.AuthMW)
	follow.RegisterRoutes(api, c.FollowHandler, c.AuthMW)
	recommendation.RegisterRoutes(api, c.RecommendationHandler, c.OnboardingHandler, c.AuthMW)
	history.RegisterRoutes(api, c.HistoryHandler, c.AuthMW)

	analytics.RegisterRoutes(api, c.AnalyticsHandler, c.AuthMW, c.AdminMW)
	notification.RegisterRoutes(api, c.NotificationHandler, c.AuthMW)
	moderation.RegisterRoutes(api, c.ModerationHandler, c.AuthMW, c.ModeratorMW)
	dashboard.RegisterAdminRoutes(api, c.DashboardHandler, c.AuthMW, c.AdminMW)
	contribution.RegisterRoutes(api, c.ContributionHandler, c.AuthMW, c.ModeratorMW, c.AdminMW)
	importcmd.RegisterAdminRoutes(api, c.ImportHandler, c.AuthMW, c.AdminMW)
	subscription.RegisterRoutes(api, c.SubscriptionHandler, c.AuthMW)

	player.RegisterPublicRoutes(api, c.PlayerHandler)
	player.RegisterPrivateRoutes(api, c.PlayerHandler, c.AuthMW)
	player.RegisterAdminRoutes(api, c.PlayerHandler, c.AuthMW, c.AdminMW)

	if c.SearchHandler != nil {
		search.RegisterRoutes(api, c.SearchHandler, c.OptionalAuthMW, middleware.RateLimitOptional(c.RDB, 60))

		// Alias for frontend compatibility: /catalog/search → /search
		catalogSearch := api.Group("/catalog")
		catalogSearch.Use(c.OptionalAuthMW, middleware.RateLimitOptional(c.RDB, 60))
		catalogSearch.GET("/search", c.SearchHandler.Search)
	}

	api.GET("/ws", func(ctx *gin.Context) {
		tokenStr := ctx.Query("token")
		if tokenStr == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := c.TokenManager.ParseAccessToken(tokenStr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
			return
		}

		c.WSHub.HandleWebSocket(ctx.Writer, ctx.Request, userID)
	})

	reactions.RegisterRoutes(api, c.ReactionsHandler, c.AuthMW)
	gamification.RegisterRoutes(api, c.GamificationHandler, c.AuthMW, c.AdminMW)
	creator.RegisterRoutes(api, c.CreatorHandler, c.AuthMW)
	social.RegisterRoutes(api, c.SocialHandler, c.AuthMW)
	ai.RegisterRoutes(api, c.AIHandler, c.AuthMW)

	cacheadmin.RegisterAdminRoutes(api, c.CacheAdminHandler, c.AuthMW, c.AdminMW)

	video.RegisterRoutes(api, c.VideoHandler, c.AuthMW, c.AdminMW)
	video.RegisterInternalRoutes(api, c.VideoHandler, c.MLServiceHMACMW)

	// Start periodic enrichment for all tracks (runs every 6 hours)
	if c.EnrichHandler != nil {
		enrichCtx := c.enrichCtx
		go func() {
			defer c.enrichDone()
			c.startPeriodicEnrichment(enrichCtx)
		}()
	}
}
