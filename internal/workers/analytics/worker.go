package analytics

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
	return "analytics"
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.logger.Info("running analytics aggregation")
			if err := w.aggregate(ctx); err != nil {
				w.logger.Error("analytics aggregation failed", zap.Error(err))
			}
		}
	}
}

func (w *Worker) aggregate(ctx context.Context) error {
	// TODO:
	// - aggregate track plays
	// - aggregate playlist creations
	// - write daily/hourly rollups
	return nil
}
