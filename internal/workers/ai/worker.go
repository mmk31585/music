package aiworker

import (
	"context"
	"time"

	"music/internal/modules/ai"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Worker struct {
	repo   *ai.Repository
	client ai.AIClient
	logger *zap.Logger
}

func NewWorker(db *sqlx.DB, logger *zap.Logger) *Worker {
	repo := ai.NewRepository(db)
	client := ai.NewFallbackClient()
	return &Worker{repo: repo, client: client, logger: logger}
}

func (w *Worker) Name() string { return "ai" }

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	svc := ai.NewService(w.repo, w.client, w.logger, true)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.logger.Info("running AI batch processing")

			count, err := svc.ProcessBatchEmbeddings(ctx, 20)
			if err != nil {
				w.logger.Error("batch embedding processing failed", zap.Error(err))
			} else if count > 0 {
				w.logger.Info("generated embeddings", zap.Int("count", count))
			}

			moodCount, err := svc.ProcessBatchMoods(ctx, 20)
			if err != nil {
				w.logger.Error("batch mood processing failed", zap.Error(err))
			} else if moodCount > 0 {
				w.logger.Info("analyzed moods", zap.Int("count", moodCount))
			}

			if count == 0 && moodCount == 0 {
				w.logger.Debug("AI batch: no pending tracks")
			}
		}
	}
}
