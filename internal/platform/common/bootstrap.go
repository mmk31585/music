package common

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"music/internal/shared/middleware"
	"music/internal/platform/validation"
	"music/internal/platform/config"
	"music/internal/platform/cache"
	"music/internal/platform/database"
	"music/internal/platform/eventbus"
	platformLogger "music/internal/platform/logger"
)

func Bootstrap(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	log, err := platformLogger.New(cfg.App.Env, cfg.Log.Level)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Info("running database migrations")
	if err := database.RunMigrations(cfg.Postgres.URL, "migrations"); err != nil {
		return nil, fmt.Errorf("migrations failed: %w", err)
	}

	db, err := database.NewPostgres(ctx, cfg.Postgres.URL)
	if err != nil {
		return nil, fmt.Errorf("init postgres: %w", err)
	}

	redisClient, err := cache.NewRedis(ctx, cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return nil, fmt.Errorf("init redis: %w", err)
	}

	v := validator.New()

	eventBus := events.NewBus(log)

	app := &App{
		Config:    cfg,
		Logger:    log,
		DB:        db,
		Redis:     redisClient,
		Validator: v,
		Events:    eventBus,
	}

	router := gin.New()
	router.Use(middleware.CORS(cfg.CORS.AllowedOrigins))
	router.Use(middleware.GinZapLogger(log))
	router.Use(middleware.GinZapRecovery(log))

	app.RegisterRoutes(router)
	registerMediaRoutes(router, cfg.Media.BasePath)

	app.Router = router
	app.HTTPServer = app.NewHTTPServer()
	app.HTTPServer.Handler = app.Router

	return app, nil
}
