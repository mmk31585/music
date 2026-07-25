package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"music/internal/app"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	application, err := app.Bootstrap(ctx)
	if err != nil {
		panic(err)
	}

	log := application.Logger

	workerContainer := app.NewWorkerContainer(application)

	go func() {
		log.Info("starting worker process",
			zap.String("env", application.Config.App.Env),
		)

		if err := workerContainer.Run(ctx); err != nil && err != context.Canceled {
			log.Fatal("worker process failed", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("worker shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), application.Config.App.ShutdownTimeout)
	defer cancel()

	if err := workerContainer.Shutdown(shutdownCtx); err != nil {
		log.Error("worker shutdown failed", zap.Error(err))
	}

	if err := application.Shutdown(shutdownCtx); err != nil {
		log.Error("application shutdown failed", zap.Error(err))
	}

	log.Info("worker stopped")
	time.Sleep(100 * time.Millisecond)
}
