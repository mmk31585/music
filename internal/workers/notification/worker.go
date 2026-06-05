package notification

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Worker struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewWorker(db *pgxpool.Pool, logger *zap.Logger) *Worker {
	return &Worker{
		db:     db,
		logger: logger,
	}
}

func (w *Worker) Name() string {
	return "notification"
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.logger.Info("processing notifications")
			if err := w.process(ctx); err != nil {
				w.logger.Error("notification processing failed", zap.Error(err))
			}
		}
	}
}

func (w *Worker) process(ctx context.Context) error {
	// TODO:
	// - deliver queued notifications
	// - mark sent/failed
	// - retry transient errors
	return nil
}
