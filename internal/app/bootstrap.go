package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"music/internal/common/validator"
	"music/internal/config"
	"music/internal/platform/cache"
	"music/internal/platform/database"
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

	app := &App{
		Config:    cfg,
		Logger:    log,
		DB:        db,
		Redis:     redisClient,
		Validator: v,
	}

	// Build the Gin engine and register all routes
	engine := gin.New()        // or gin.Default() if you want logger/recovery; we'll add middlewares in main
	app.RegisterRoutes(engine) // RegisterRoutes now accepts *gin.Engine
	app.Router = engine        // engine implements http.Handler

	// Create HTTP server with the Gin engine as handler
	app.HTTPServer = app.NewHTTPServer()

	return app, nil
}
