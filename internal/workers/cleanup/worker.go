package cleanup

import (
	"context"
	"fmt"
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
	defer func() {
		if r := recover(); r != nil {
			w.logger.Error("cleanup worker panicked", zap.Any("recover", r))
		}
	}()

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
	if err := w.deleteExpiredSessions(ctx); err != nil {
		return fmt.Errorf("expired sessions: %w", err)
	}
	if err := w.markStaleDrafts(ctx); err != nil {
		return fmt.Errorf("stale drafts: %w", err)
	}
	if err := w.pruneOldNotifications(ctx); err != nil {
		return fmt.Errorf("old notifications: %w", err)
	}
	if err := w.cleanupDisconnectedSessions(ctx); err != nil {
		return fmt.Errorf("disconnected sessions: %w", err)
	}
	return nil
}

func (w *Worker) deleteExpiredSessions(ctx context.Context) error {
	tag, err := w.db.Exec(ctx, `
		DELETE FROM auth_sessions
		WHERE expires_at < NOW()
	`)
	if err != nil {
		return err
	}
	if tag.RowsAffected() > 0 {
		w.logger.Info("deleted expired auth sessions", zap.Int64("count", tag.RowsAffected()))
	}
	return nil
}

func (w *Worker) markStaleDrafts(ctx context.Context) error {
	tag, err := w.db.Exec(ctx, `
		UPDATE ingestion_drafts
		SET stale = TRUE
		WHERE status IN ('pending', 'enriching')
		  AND created_at < NOW() - INTERVAL '24 hours'
		  AND stale = FALSE
	`)
	if err != nil {
		return err
	}
	if tag.RowsAffected() > 0 {
		w.logger.Info("marked stale ingestion drafts", zap.Int64("count", tag.RowsAffected()))
	}
	return nil
}

func (w *Worker) pruneOldNotifications(ctx context.Context) error {
	tag, err := w.db.Exec(ctx, `
		DELETE FROM notifications
		WHERE is_read = TRUE
		  AND read_at < NOW() - INTERVAL '30 days'
	`)
	if err != nil {
		return err
	}
	if tag.RowsAffected() > 0 {
		w.logger.Info("pruned old read notifications", zap.Int64("count", tag.RowsAffected()))
	}
	return nil
}

func (w *Worker) cleanupDisconnectedSessions(ctx context.Context) error {
	tag, err := w.db.Exec(ctx, `
		DELETE FROM ws_sessions
		WHERE disconnected_at IS NOT NULL
		  AND disconnected_at < NOW() - INTERVAL '24 hours'
	`)
	if err != nil {
		return err
	}
	if tag.RowsAffected() > 0 {
		w.logger.Info("cleaned up old ws sessions", zap.Int64("count", tag.RowsAffected()))
	}
	return nil
}
