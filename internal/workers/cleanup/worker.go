package cleanup

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
	return "cleanup"
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.logger.Info("running cleanup jobs")
			if err := w.cleanup(ctx); err != nil {
				w.logger.Error("cleanup jobs failed", zap.Error(err))
			}
		}
	}
}

func (w *Worker) cleanup(ctx context.Context) error {
	// TODO:
	// - delete expired sessions/tokens
	// - remove temp files
	// - prune stale notifications
	// - cleanup orphaned media
	return nil
}
