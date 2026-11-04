package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"music/internal/app"
	"music/internal/common/middleware"
	"music/internal/platform/sentry"
	"music/internal/platform/tracing"

	// Registers the generated OpenAPI spec with swag so the swagger UI at
	// /swagger/* can serve /swagger/doc.json. Without this blank import the
	// registry is empty and the endpoint returns HTTP 500.
	_ "music/docs"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

// @title           Music Application API
// @version         1.0
// @description     This is the API server for the Music Streaming Application.
// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer " followed by your token

func main() {
	ctx := context.Background()

	application, err := app.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	log := application.Logger

	// ── OpenTelemetry Tracing ──────────────────────────────────────
	otelShutdown, err := tracing.Init(ctx, tracing.Config{
		ServiceName: "music-api",
		Environment: application.Config.App.Env,
	})
	if err != nil {
		log.Fatal("failed to init tracing", zap.Error(err))
	}
	defer otelShutdown()

	// ── Sentry Error Tracking ──────────────────────────────────────
	if application.Config.Sentry.DSN != "" {
		if err := sentry.Init(application.Config.Sentry.DSN, application.Config.App.Env); err != nil {
			log.Fatal("failed to init sentry", zap.Error(err))
		}
		defer sentry.Flush()
		log.Info("sentry initialized")
	}

	// ── Router setup (exactly once, all middleware explicit) ──
	router := gin.New()
	router.Use(
		middleware.SecurityHeaders(),
		middleware.CORS(application.Config.CORS.AllowedOrigins),
		otelgin.Middleware("music-api"),
		middleware.GinZapLogger(log),
		middleware.GinZapRecovery(log),
		middleware.Sentry(),
	)

	// Serve uploaded media from local storage
	app.RegisterMediaRoutes(router, application.Config.Storage.Local.BaseDir)

	// Register swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register all application routes (exactly one call)
	application.RegisterRoutes(router)

	application.Router = router
	application.HTTPServer = application.NewHTTPServer()

	go func() {
		log.Info("starting server",
			zap.String("addr", application.HTTPServer.Addr),
			zap.String("env", application.Config.App.Env),
		)

		if err := application.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	// ── Background workers (import, notifications, etc.) ──
	// Only started when this process owns them. When a dedicated worker
	// container is deployed, RUN_BACKGROUND_WORKERS=false keeps the API from
	// running every job a second time.
	workerCtx, cancelWorkers := context.WithCancel(context.Background())
	defer cancelWorkers()

	var workerContainer *app.WorkerContainer
	if application.Config.App.RunBackgroundWorkers {
		workerContainer = app.NewWorkerContainer(application)
		go func() {
			log.Info("starting background workers",
				zap.String("env", application.Config.App.Env),
			)
			if err := workerContainer.Run(workerCtx); err != nil && err != context.Canceled {
				log.Fatal("background workers failed", zap.Error(err))
			}
		}()
	} else {
		log.Info("in-process background workers disabled; expecting a dedicated worker container")
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("shutdown signal received")

	// Cancel workers first so they stop accepting jobs
	cancelWorkers()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), application.Config.App.ShutdownTimeout)
	defer cancel()

	if workerContainer != nil {
		if err := workerContainer.Shutdown(shutdownCtx); err != nil {
			log.Error("worker shutdown failed", zap.Error(err))
		}
	}

	if err := application.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
	}

	log.Info("server stopped")
	time.Sleep(100 * time.Millisecond)
}
