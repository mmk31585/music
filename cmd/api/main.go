package main

import (
	"context"
	"music/docs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"music/internal/app"
	"music/internal/common/middleware"

	"github.com/gin-gonic/gin"
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
	docs.SwaggerInfo.Title = "Music API"
	docs.SwaggerInfo.Description = "Music streaming API documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// ── Router setup (exactly once, all middleware explicit) ──
	router := gin.New()
	router.Use(
		middleware.SecurityHeaders(),
		middleware.CORS(application.Config.CORS.AllowedOrigins),
		middleware.GinZapLogger(log),
		middleware.GinZapRecovery(log),
	)

	// Serve uploaded media from local storage
	app.RegisterMediaRoutes(router, application.Config.Storage.Local.BaseDir)

	// Register all application routes (exactly one call)
	application.RegisterRoutes(router)

	application.Router = router
	application.HTTPServer = application.NewHTTPServer()
	application.HTTPServer.Handler = application.Router

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
	workerCtx, cancelWorkers := context.WithCancel(context.Background())
	defer cancelWorkers()

	workerContainer := app.NewWorkerContainer(application)
	go func() {
		log.Info("starting background workers",
			zap.String("env", application.Config.App.Env),
		)
		if err := workerContainer.Run(workerCtx); err != nil && err != context.Canceled {
			log.Fatal("background workers failed", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("shutdown signal received")

	// Cancel workers first so they stop accepting jobs
	cancelWorkers()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), application.Config.App.ShutdownTimeout)
	defer cancel()

	if err := workerContainer.Shutdown(shutdownCtx); err != nil {
		log.Error("worker shutdown failed", zap.Error(err))
	}

	if err := application.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
	}

	log.Info("server stopped")
	time.Sleep(100 * time.Millisecond)
}
