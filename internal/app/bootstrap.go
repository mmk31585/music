package app

import (
	"context"
	"fmt"

	"music/internal/common/validator"
	"music/internal/config"
	"music/internal/platform/cache"
	"music/internal/platform/database"
	"music/internal/platform/events"
	platformLogger "music/internal/platform/logger"
)

// Bootstrap creates the App with all core dependencies (config, logger, DB, Redis,
// validator, event bus) but does NOT create the Gin engine or register routes.
//
// Route registration and HTTP server creation are handled by main.go via
// SetupRouter() and app.NewHTTPServer(), ensuring routes are registered exactly
// once and middleware choices are explicit at the call site.
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

	eventBus := events.NewBus(log)

	appCtx, appCancel := context.WithCancel(context.Background())

	return &App{
		Config:    cfg,
		Logger:    log,
		DB:        db,
		Redis:     redisClient,
		Validator: v,
		Events:    eventBus,
		ctx:       appCtx,
		cancel:    appCancel,
	}, nil
}
