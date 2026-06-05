package recommendation

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
	return "recommendation"
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.logger.Info("running recommendation generation")
			if err := w.generate(ctx); err != nil {
				w.logger.Error("recommendation generation failed", zap.Error(err))
			}
		}
	}
}

func (w *Worker) generate(ctx context.Context) error {
	// TODO:
	// - refresh user recommendations
	// - compute similarity graph
	// - rebuild trending lists
	return nil
}
