package analytics

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
	return "analytics"
}

func (w *Worker) Run(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			w.logger.Error("analytics worker panicked", zap.Any("recover", r))
		}
	}()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.aggregate(ctx); err != nil {
				w.logger.Error("analytics aggregation failed", zap.Error(err))
			}
		}
	}
}

// aggregate refreshes the analytics that are not maintained in real time.
//
// track_analytics / artist_analytics / album_analytics are incremented by the
// event-driven analytics module (internal/modules/analytics) whenever a
// TrackPlayed event arrives. The worker used to re-add a rolling 24h window to
// those same cumulative columns on every run, double-counting each play (and
// growing the counts without bound). It must not touch them.
//
// creator_daily_stats has no event-driven writer, so it is refreshed here.
func (w *Worker) aggregate(ctx context.Context) error {
	if err := w.refreshCreatorDailyStats(ctx); err != nil {
		return fmt.Errorf("creator daily stats: %w", err)
	}
	return nil
}

// refreshCreatorDailyStats recomputes plays/listeners per creator per day for
// today and yesterday and REPLACES the stored values. Replacing (rather than
// `+ EXCLUDED`) makes the job idempotent: the previous accumulation added the
// same rolling window on every run. The window covers all of yesterday and
// today, so those two days are always recomputed from complete data; older days
// keep their last value.
//
// The creator is the user who uploaded the track, not the artist: artists are
// catalog entities with their own uuid and no owner column, so tracks.artist_id
// can never satisfy the creator_daily_stats.user_id -> users.id foreign key.
func (w *Worker) refreshCreatorDailyStats(ctx context.Context) error {
	tag, err := w.db.Exec(ctx, `
		INSERT INTO creator_daily_stats (user_id, date, plays, listeners, updated_at)
		SELECT
			d.uploaded_by,
			lh.played_at::date AS date,
			COUNT(*) AS plays,
			COUNT(DISTINCT lh.user_id) AS listeners,
			NOW()
		FROM listening_history lh
		JOIN (
			SELECT DISTINCT track_id, uploaded_by
			FROM ingestion_drafts
			WHERE track_id IS NOT NULL AND uploaded_by IS NOT NULL
		) d ON d.track_id = lh.track_id
		WHERE lh.played_at >= CURRENT_DATE - INTERVAL '1 day'
		GROUP BY d.uploaded_by, lh.played_at::date
		ON CONFLICT (user_id, date) DO UPDATE SET
			plays = EXCLUDED.plays,
			listeners = EXCLUDED.listeners,
			updated_at = NOW()
	`)
	if err != nil {
		return err
	}
	if tag.RowsAffected() > 0 {
		w.logger.Info("refreshed creator daily stats", zap.Int64("rows", tag.RowsAffected()))
	}
	return nil
}
