package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"music/internal/app"

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

	go func() {
		log.Info("starting server",
			zap.String("addr", application.HTTPServer.Addr),
			zap.String("env", application.Config.App.Env),
		)

		if err := application.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), application.Config.App.ShutdownTimeout)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
	}

	log.Info("server stopped")
	time.Sleep(100 * time.Millisecond)
}
