package app

import (
	"context"
	"music/internal/common/middleware"
	ai "music/internal/modules/ai"
	"music/internal/modules/analytics"
	"music/internal/modules/auth"
	"music/internal/modules/catalog"
	"music/internal/modules/contribution"
	"music/internal/modules/covers"
	"music/internal/modules/creator"
	"music/internal/modules/dashboard"
	"music/internal/modules/follow"
	"music/internal/modules/gamification"
	"music/internal/modules/health"
	"music/internal/modules/history"
	"music/internal/modules/importcmd"
	"music/internal/modules/ingestion"
	"music/internal/modules/library"
	"music/internal/modules/lyrics"
	"music/internal/modules/media"
	"music/internal/modules/moderation"
	"music/internal/modules/notification"
	"music/internal/modules/player"
	"music/internal/modules/playlist"
	"music/internal/modules/queue"
	"music/internal/modules/reactions"
	"music/internal/modules/recommendation"
	"music/internal/modules/search"
	"music/internal/modules/social"
	"music/internal/modules/subscription"
	"music/internal/modules/video"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (a *App) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api.GET("/metrics", gin.WrapH(promhttp.Handler()))
	api.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	c := NewContainer(a)

	health.RegisterRoutes(api, c.HealthHandler)
	auth.RegisterRoutes(api, c.AuthHandler, c.AuthMW)

	media.RegisterAdminRoutes(api, c.MediaHandler, c.AuthMW)

	ingestion.RegisterAdminRoutes(api, c.IngestionHandler, c.AuthMW)

	catalogHandlers := catalog.Handlers{
		Artist: c.ArtistHandler,
		Album:  c.AlbumHandler,
		Genre:  c.GenreHandler,
		Track:  c.TrackHandler,
		Enrich: c.EnrichHandler,
	}

	catalog.RegisterPublicRoutes(api, catalogHandlers, c.OptionalAuthMW, middleware.RateLimitOptional(c.RDB, 60))
	catalog.RegisterAdminRoutes(api, catalogHandlers, c.AuthMW)

	lyrics.RegisterRoutes(api, c.LyricsHandler, c.AuthMW)
	lyrics.RegisterInternalRoutes(api, c.LyricsHandler, c.MLServiceHMACMW)

	covers.RegisterInternalRoutes(api, c.CoverHandler, c.MLServiceHMACMW)

	playlist.RegisterRoutes(api, c.PlaylistHandler, c.AuthMW, c.OptionalAuthMW)
	library.RegisterRoutes(api, c.LibraryHandler, c.AuthMW)
	queue.RegisterRoutes(api, c.QueueHandler, c.AuthMW)
	follow.RegisterRoutes(api, c.FollowHandler, c.AuthMW)
	recommendation.RegisterRoutes(api, c.RecommendationHandler, c.OnboardingHandler, c.AuthMW)
	history.RegisterRoutes(api, c.HistoryHandler, c.AuthMW)

	analytics.RegisterRoutes(api, c.AnalyticsHandler)
	notification.RegisterRoutes(api, c.NotificationHandler, c.AuthMW)
	moderation.RegisterRoutes(api, c.ModerationHandler, c.AuthMW)
	dashboard.RegisterAdminRoutes(api, c.DashboardHandler, c.AuthMW)
	contribution.RegisterRoutes(api, c.ContributionHandler, c.AuthMW)
	importcmd.RegisterAdminRoutes(api, c.ImportHandler, c.AuthMW)
	subscription.RegisterRoutes(api, c.SubscriptionHandler, c.AuthMW)

	player.RegisterPublicRoutes(api, c.PlayerHandler)
	player.RegisterPrivateRoutes(api, c.PlayerHandler, c.AuthMW)

	if c.SearchHandler != nil {
		search.RegisterRoutes(api, c.SearchHandler, c.OptionalAuthMW, middleware.RateLimitOptional(c.RDB, 60))

		// Alias for frontend compatibility: /catalog/search → /search
		catalogSearch := api.Group("/catalog")
		catalogSearch.Use(c.OptionalAuthMW, middleware.RateLimitOptional(c.RDB, 60))
		catalogSearch.GET("/search", c.SearchHandler.Search)
	}

	tokenManager := auth.NewTokenManager(
		a.Config.Auth.JWTAccessSecret,
		a.Config.Auth.JWTRefreshSecret,
		a.Config.Auth.AccessTTL,
		a.Config.Auth.RefreshTTL,
	)

	api.GET("/ws", func(ctx *gin.Context) {
		tokenStr := ctx.Query("token")
		if tokenStr == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := tokenManager.ParseAccessToken(tokenStr)
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
	gamification.RegisterRoutes(api, c.GamificationHandler, c.AuthMW)
	creator.RegisterRoutes(api, c.CreatorHandler, c.AuthMW)
	social.RegisterRoutes(api, c.SocialHandler, c.AuthMW)
	ai.RegisterRoutes(api, c.AIHandler, c.AuthMW)

	video.RegisterRoutes(api, c.VideoHandler, c.AuthMW)
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

// startPeriodicEnrichment runs enrichment on a 6-hour loop until the context is cancelled.
func (c *Container) startPeriodicEnrichment(ctx context.Context) {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	// Run once at startup (non-blocking)
	c.EnrichHandler.EnrichAllBackground(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.EnrichHandler.EnrichAllBackground(ctx)
		}
	}
}
