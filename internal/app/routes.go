package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"music/internal/modules/auth"
	"music/internal/modules/catalog"
	"music/internal/modules/media"
)

func (a *App) RegisterRoutes(r *gin.Engine) {
	// API v1 group
	api := r.Group("/api/v1")

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
	authMW := auth.AuthMiddleware(tokenManager) // returns gin.HandlerFunc

	// Register auth routes (pass the group and middleware)
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

	// ---------- Catalog module ----------
	sqlDB := stdlib.OpenDBFromPool(a.DB)
	sqlxDB := sqlx.NewDb(sqlDB, "pgx")
	catalogRepo := catalog.NewRepository(sqlxDB)
	catalogService := catalog.NewService(catalogRepo)
	catalogHandler := catalog.NewHandler(catalogService, a.Validator, mediaService)

	// Public catalog routes (no auth)
	catalog.RegisterPublicRoutes(api, catalogHandler)
	// Admin catalog routes (with auth middleware)
	catalog.RegisterAdminRoutes(api, catalogHandler, authMW)

	mediaHandler := media.NewHandler(mediaService)
	media.RegisterAdminRoutes(api, mediaHandler, authMW)
}
