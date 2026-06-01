package app

import (
	"music/internal/modules/analytics"
	"music/internal/modules/identity"
	"music/internal/modules/catalog"
	"music/internal/modules/catalog/album"
	artist "music/internal/modules/catalog/artist"
	"music/internal/modules/catalog/genre"
	"music/internal/modules/catalog/track"
	"music/internal/modules/follow"
	"music/internal/modules/health"
	"music/internal/modules/history"
	"music/internal/modules/library"
	"music/internal/modules/media"
	"music/internal/modules/notification"
	"music/internal/modules/player"
	"music/internal/modules/playlist"
	"music/internal/modules/queue"
	"music/internal/modules/recommendation"
	"music/internal/modules/search"
	"music/internal/modules/subscription"
	"music/internal/platform/eventbus"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	opensearch "github.com/opensearch-project/opensearch-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func (a *App) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	sqlDB := stdlib.OpenDBFromPool(a.DB)
	sqlxDB := sqlx.NewDb(sqlDB, "pgx")

	logger := zap.L()
	bus := events.NewBus(logger)

	healthHandler := health.NewHandler(a.DB, a.Redis)
	health.RegisterRoutes(api, healthHandler)

	tokenManager := auth.NewTokenManager(
		a.Config.Auth.JWTAccessSecret,
		a.Config.Auth.JWTRefreshSecret,
		a.Config.Auth.AccessTTL,
		a.Config.Auth.RefreshTTL,
	)

	authRepo := auth.NewRepository(a.DB)
	authService := auth.NewService(authRepo, tokenManager)
	authHandler := auth.NewHandler(authService, a.Validator)
	authMW := auth.AuthMiddleware(tokenManager)
	auth.RegisterRoutes(api, authHandler, authMW)

	mediaStorage := media.NewLocalStorage(media.LocalStorageConfig{
		BasePath:   a.Config.Media.BasePath,
		PublicBase: a.Config.Media.PublicBase,
	})

	mediaService := media.NewService(mediaStorage, media.Config{
		MaxFileSizeBytes:  a.Config.Media.MaxFileSizeBytes,
		MaxImageSizeBytes: a.Config.Media.MaxImageSizeBytes,
		MaxAudioSizeBytes: a.Config.Media.MaxAudioSizeBytes,
		AllowedImageMime: []string{
			"image/jpeg",
			"image/png",
			"image/webp",
		},
		AllowedAudioMime: []string{
			"audio/mpeg",
			"audio/ogg",
			"audio/flac",
			"audio/wav",
			"audio/x-wav",
			"audio/wave",
			"audio/mp4",
			"audio/aac",
		},
	})

	mediaHandler := media.NewHandler(mediaService)
	media.RegisterAdminRoutes(api, mediaHandler, authMW)

	artistRepo := artist.NewRepository(sqlxDB)
	artistService := artist.NewService(artistRepo)
	artistHandler := artist.NewHandler(artistService)

	albumRepo := album.NewRepository(sqlxDB)
	albumService := album.NewService(albumRepo)
	albumHandler := album.NewHandler(albumService)

	genreRepo := genre.NewRepository(sqlxDB)
	genreService := genre.NewService(genreRepo)
	genreHandler := genre.NewHandler(genreService)

	trackRepo := track.NewRepository(sqlxDB)
	trackService := track.NewService(trackRepo)
	trackHandler := track.NewHandler(trackService)

	catalogHandlers := catalog.Handlers{
		Artist: artistHandler,
		Album:  albumHandler,
		Genre:  genreHandler,
		Track:  trackHandler,
	}

	catalog.RegisterPublicRoutes(api, catalogHandlers)
	catalog.RegisterAdminRoutes(api, catalogHandlers, authMW)

	playlistRepo := playlist.NewRepository(sqlxDB)
	playlistService := playlist.NewService(playlistRepo)
	playlistHandler := playlist.NewHandler(playlistService)
	playlist.RegisterRoutes(api, playlistHandler, authMW)

	libraryRepo := library.NewRepository(sqlxDB)
	libraryService := library.NewService(libraryRepo)
	libraryHandler := library.NewHandler(libraryService)
	library.RegisterRoutes(api, libraryHandler, authMW)

	queueRepo := queue.NewRepository(sqlxDB)
	queueService := queue.NewService(queueRepo)
	queueHandler := queue.NewHandler(queueService)
	queue.RegisterRoutes(api, queueHandler, authMW)

	followRepo := follow.NewRepository(sqlxDB)
	followService := follow.NewService(followRepo)
	followHandler := follow.NewHandler(followService)
	follow.RegisterRoutes(api, followHandler, authMW)

	recommendationRepo := recommendation.NewRepository(sqlxDB)
	recommendationService := recommendation.NewService(recommendationRepo)
	recommendationHandler := recommendation.NewHandler(recommendationService)
	recommendation.RegisterRoutes(api, recommendationHandler, authMW)

	historyRepo := history.NewRepository(sqlxDB)
	historyService := history.NewService(historyRepo)
	historyHandler := history.NewHandler(historyService)
	history.RegisterRoutes(api, historyHandler, authMW)

	analyticsRepo := analytics.NewRepository(sqlxDB)
	analyticsService := analytics.NewService(analyticsRepo)
	analyticsHandler := analytics.NewHandler(analyticsService)
	analytics.RegisterRoutes(api, analyticsHandler)

	notificationRepo := notification.NewRepository(sqlxDB)
	notificationService := notification.NewService(notificationRepo)
	notificationHandler := notification.NewHandler(notificationService)
	notification.RegisterRoutes(api, notificationHandler, authMW)

	subscriptionRepo := subscription.NewRepository(sqlxDB)
	subscriptionService := subscription.NewService(subscriptionRepo, bus)
	subscriptionHandler := subscription.NewHandler(subscriptionService)
	subscription.RegisterRoutes(api, subscriptionHandler, authMW)

	playerRepo := player.NewRepository(sqlxDB)
	playerService := player.NewService(playerRepo, a.Config.Media.BasePath, bus)
	playerHandler := player.NewHandler(playerService)
	player.RegisterPublicRoutes(api, playerHandler)

	if a.Config.OpenSearch.URL != "" {
		osClient, err := opensearch.NewClient(opensearch.Config{
			Addresses: []string{a.Config.OpenSearch.URL},
		})
		if err == nil {
			searchService := search.NewService(osClient)
			searchHandler := search.NewHandler(searchService)
			search.RegisterRoutes(api, searchHandler)
		}
	}

	analyticsEvents := analytics.NewEventHandler(analyticsService)
	notificationEvents := notification.NewEventHandler(notificationService)
	recommendationEvents := recommendation.NewEventHandler(recommendationService)
	historyEvents := history.NewEventHandler(historyService)

	bus.Subscribe(events.EventTrackPlayed, analyticsEvents.OnTrackPlayed)
	bus.Subscribe(events.EventTrackPlayed, historyEvents.OnTrackPlayed)
	bus.Subscribe(events.EventTrackPlayed, recommendationEvents.OnTrackPlayed)

	bus.Subscribe(events.EventPlaylistCreated, analyticsEvents.OnPlaylistCreated)
	bus.Subscribe(events.EventPlaylistCreated, notificationEvents.OnPlaylistCreated)

	bus.Subscribe(events.EventSubscriptionPurchased, analyticsEvents.OnSubscriptionPurchased)
	bus.Subscribe(events.EventSubscriptionPurchased, notificationEvents.OnSubscriptionPurchased)

	bus.Subscribe(events.EventUserRegistered, analyticsEvents.OnUserRegistered)
	bus.Subscribe(events.EventUserRegistered, notificationEvents.OnUserRegistered)
}
