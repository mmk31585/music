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
	if err := w.aggregateTrackPlays(ctx); err != nil {
		return fmt.Errorf("track plays: %w", err)
	}
	if err := w.aggregateArtistPlays(ctx); err != nil {
		return fmt.Errorf("artist plays: %w", err)
	}
	if err := w.aggregateAlbumPlays(ctx); err != nil {
		return fmt.Errorf("album plays: %w", err)
	}
	if err := w.aggregateCreatorDailyStats(ctx); err != nil {
		return fmt.Errorf("creator daily stats: %w", err)
	}
	playlistCount, err := w.countRecentPlaylistCreations(ctx)
	if err != nil {
		return fmt.Errorf("playlist creations: %w", err)
	}
	w.logger.Info("analytics summary", zap.Int64("playlist_creations_24h", playlistCount))
	return nil
}

func (w *Worker) aggregateTrackPlays(ctx context.Context) error {
	_, err := w.db.Exec(ctx, `
		INSERT INTO track_analytics (track_id, play_count, pause_count, skip_count, completion_count, updated_at)
		SELECT
			lh.track_id,
			COUNT(*) AS play_count,
			0, 0,
			COUNT(*) FILTER (WHERE lh.completed = TRUE) AS completion_count,
			NOW()
		FROM listening_history lh
		WHERE lh.played_at > NOW() - INTERVAL '24 hours'
		GROUP BY lh.track_id
		ON CONFLICT (track_id) DO UPDATE SET
			play_count = track_analytics.play_count + EXCLUDED.play_count,
			completion_count = track_analytics.completion_count + EXCLUDED.completion_count,
			updated_at = NOW()
	`)
	return err
}

func (w *Worker) aggregateArtistPlays(ctx context.Context) error {
	_, err := w.db.Exec(ctx, `
		INSERT INTO artist_analytics (artist_id, play_count, pause_count, skip_count, completion_count, updated_at)
		SELECT
			t.artist_id,
			COUNT(*) AS play_count,
			0, 0,
			COUNT(*) FILTER (WHERE lh.completed = TRUE) AS completion_count,
			NOW()
		FROM listening_history lh
		JOIN tracks t ON t.id = lh.track_id
		WHERE lh.played_at > NOW() - INTERVAL '24 hours'
		  AND t.artist_id IS NOT NULL
		GROUP BY t.artist_id
		ON CONFLICT (artist_id) DO UPDATE SET
			play_count = artist_analytics.play_count + EXCLUDED.play_count,
			completion_count = artist_analytics.completion_count + EXCLUDED.completion_count,
			updated_at = NOW()
	`)
	return err
}

func (w *Worker) aggregateAlbumPlays(ctx context.Context) error {
	_, err := w.db.Exec(ctx, `
		INSERT INTO album_analytics (album_id, play_count, pause_count, skip_count, completion_count, updated_at)
		SELECT
			t.album_id,
			COUNT(*) AS play_count,
			0, 0,
			COUNT(*) FILTER (WHERE lh.completed = TRUE) AS completion_count,
			NOW()
		FROM listening_history lh
		JOIN tracks t ON t.id = lh.track_id
		WHERE lh.played_at > NOW() - INTERVAL '24 hours'
		  AND t.album_id IS NOT NULL
		GROUP BY t.album_id
		ON CONFLICT (album_id) DO UPDATE SET
			play_count = album_analytics.play_count + EXCLUDED.play_count,
			completion_count = album_analytics.completion_count + EXCLUDED.completion_count,
			updated_at = NOW()
	`)
	return err
}

func (w *Worker) aggregateCreatorDailyStats(ctx context.Context) error {
	_, err := w.db.Exec(ctx, `
		INSERT INTO creator_daily_stats (user_id, date, plays, listeners, updated_at)
		SELECT
			t.artist_id,
			lh.played_at::date AS date,
			COUNT(*) AS plays,
			COUNT(DISTINCT lh.user_id) AS listeners,
			NOW()
		FROM listening_history lh
		JOIN tracks t ON t.id = lh.track_id
		WHERE lh.played_at > NOW() - INTERVAL '24 hours'
		  AND t.artist_id IS NOT NULL
		GROUP BY t.artist_id, lh.played_at::date
		ON CONFLICT (user_id, date) DO UPDATE SET
			plays = creator_daily_stats.plays + EXCLUDED.plays,
			listeners = creator_daily_stats.listeners + EXCLUDED.listeners,
			updated_at = NOW()
	`)
	return err
}

func (w *Worker) countRecentPlaylistCreations(ctx context.Context) (int64, error) {
	var count int64
	err := w.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM playlists
		WHERE created_at > NOW() - INTERVAL '24 hours'
	`).Scan(&count)
	return count, err
}
