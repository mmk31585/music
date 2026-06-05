package app

import (
	"context"
	"sync"

	"music/internal/workers/analytics"
	"music/internal/workers/cleanup"
	"music/internal/workers/notification"
	"music/internal/workers/recommendation"

	"go.uber.org/zap"
)

type BackgroundWorker interface {
	Run(ctx context.Context) error
	Name() string
}

type WorkerContainer struct {
	App *App

	AnalyticsWorker      BackgroundWorker
	RecommendationWorker BackgroundWorker
	NotificationWorker   BackgroundWorker
	EmailWorker          BackgroundWorker
	CleanupWorker        BackgroundWorker

	workers []BackgroundWorker
}

func NewWorkerContainer(a *App) *WorkerContainer {
	c := &WorkerContainer{
		App: a,
	}

	c.buildWorkers()

	return c
}

func (c *WorkerContainer) buildWorkers() {
	c.AnalyticsWorker = analytics.NewWorker(c.App.DB, c.App.Logger)
	c.RecommendationWorker = recommendation.NewWorker(c.App.DB, c.App.Logger)
	c.NotificationWorker = notification.NewWorker(c.App.DB, c.App.Logger)
	c.CleanupWorker = cleanup.NewWorker(c.App.DB, c.App.Logger)

	c.workers = []BackgroundWorker{
		c.AnalyticsWorker,
		c.RecommendationWorker,
		c.NotificationWorker,
		c.EmailWorker,
		c.CleanupWorker,
	}
}

func (c *WorkerContainer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(c.workers))

	for _, worker := range c.workers {
		if worker == nil {
			continue
		}

		wg.Add(1)

		go func(w BackgroundWorker) {
			defer wg.Done()

			c.App.Logger.Error("worker failed", zap.String("name", w.Name()))
			if err := w.Run(ctx); err != nil && err != context.Canceled {
				errCh <- err
			}
		}(worker)
	}

	select {
	case <-ctx.Done():
		wg.Wait()
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (c *WorkerContainer) Shutdown(ctx context.Context) error {
	return nil
}
