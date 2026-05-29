package app

import (
	"music/internal/modules/auth"
	"music/internal/modules/catalog"
	"music/internal/modules/library"
	"music/internal/modules/media"
	"music/internal/modules/notification"
	"music/internal/modules/player"
	"music/internal/modules/playlist"
	"music/internal/modules/recommendation"
	"music/internal/modules/search"
	"music/internal/modules/subscription"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	opensearch "github.com/opensearch-project/opensearch-go"
)

func (a *App) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	// ---------- Shared DB adapter ----------
	sqlDB := stdlib.OpenDBFromPool(a.DB)
	sqlxDB := sqlx.NewDb(sqlDB, "pgx")

	// ---------- Auth module ----------
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

	// ---------- Media module ----------
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

	// ---------- Catalog module ----------
	catalogRepo := catalog.NewRepository(sqlxDB)
	catalogService := catalog.NewService(catalogRepo)
	catalogHandler := catalog.NewHandler(catalogService, a.Validator, mediaService)

	catalog.RegisterPublicRoutes(api, catalogHandler)
	catalog.RegisterAdminRoutes(api, catalogHandler, authMW)

	// ---------- Player module ----------
	playerRepo := player.NewRepository(sqlxDB)
	playerService := player.NewService(playerRepo, a.Config.Media.BasePath)
	playerHandler := player.NewHandler(playerService)

	player.RegisterPublicRoutes(api, playerHandler)

	// ---------- Playlist module ----------
	playlistRepo := playlist.NewRepository(sqlxDB)
	playlistService := playlist.NewService(playlistRepo)
	playlistHandler := playlist.NewHandler(playlistService)

	playlist.RegisterRoutes(api, playlistHandler, authMW)

	// ---------- Library module ----------
	libraryRepo := library.NewRepository(sqlxDB)
	libraryService := library.NewService(libraryRepo)
	libraryHandler := library.NewHandler(libraryService)

	library.RegisterRoutes(api, libraryHandler, authMW)

	// ---------- Recommendation module ----------
	recommendationRepo := recommendation.NewRepository(sqlxDB)
	recommendationService := recommendation.NewService(recommendationRepo)
	recommendationHandler := recommendation.NewHandler(recommendationService)

	recommendation.RegisterRoutes(api, recommendationHandler, authMW)

	// ---------- Search module ----------
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
	notificationRepo := notification.NewRepository(sqlxDB)
	notificationService := notification.NewService(notificationRepo)
	notificationHandler := notification.NewHandler(notificationService)

	notification.RegisterRoutes(api, notificationHandler, authMW)

	subscriptionRepo := subscription.NewRepository(sqlxDB)
	subscriptionService := subscription.NewService(subscriptionRepo)
	subscriptionHandler := subscription.NewHandler(subscriptionService)

	subscription.RegisterRoutes(api, subscriptionHandler, authMW)

}
