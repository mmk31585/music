package recommendation

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
	return "recommendation"
}

func (w *Worker) Run(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			w.logger.Error("recommendation worker panicked", zap.Any("recover", r))
		}
	}()

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
	if err := w.computeContentSimilarity(ctx); err != nil {
		return fmt.Errorf("content similarity: %w", err)
	}
	if err := w.refreshUserAffinities(ctx); err != nil {
		return fmt.Errorf("user affinities: %w", err)
	}
	if err := w.updateTrending(ctx); err != nil {
		return fmt.Errorf("trending: %w", err)
	}
	return nil
}

func (w *Worker) computeContentSimilarity(ctx context.Context) error {
	_, err := w.db.Exec(ctx, `
		INSERT INTO track_similarities (track_id, similar_track_id, similarity, method, updated_at)
		SELECT
			t1.id,
			t2.id,
			CASE
				WHEN t1.artist_id IS NOT NULL AND t1.artist_id = t2.artist_id THEN 0.8
				WHEN t1.album_id IS NOT NULL AND t1.album_id = t2.album_id THEN 0.6
				ELSE 0.0
			END + CASE WHEN tg1.genre_id = tg2.genre_id THEN 0.3 ELSE 0.0 END AS similarity,
			'content',
			NOW()
		FROM tracks t1
		JOIN tracks t2 ON t2.id <> t1.id
		LEFT JOIN track_genres tg1 ON tg1.track_id = t1.id
		LEFT JOIN track_genres tg2 ON tg2.track_id = t2.id AND tg2.genre_id = tg1.genre_id
		WHERE (t1.artist_id IS NOT NULL AND t1.artist_id = t2.artist_id)
		   OR (t1.album_id IS NOT NULL AND t1.album_id = t2.album_id)
		   OR tg1.genre_id = tg2.genre_id
		GROUP BY t1.id, t2.id, t1.artist_id, t2.artist_id, t1.album_id, t2.album_id
		HAVING CASE
			WHEN t1.artist_id IS NOT NULL AND t1.artist_id = t2.artist_id THEN 0.8
			WHEN t1.album_id IS NOT NULL AND t1.album_id = t2.album_id THEN 0.6
			ELSE 0.0
		END + MAX(CASE WHEN tg1.genre_id = tg2.genre_id THEN 0.3 ELSE 0.0 END) > 0
		ON CONFLICT (track_id, similar_track_id)
		DO UPDATE SET similarity = EXCLUDED.similarity, method = 'content', updated_at = NOW()
	`)
	return err
}

func (w *Worker) refreshUserAffinities(ctx context.Context) error {
	_, err := w.db.Exec(ctx, `
		INSERT INTO user_affinities (user_id, target_id, target_type, score, decay, updated_at)
		SELECT
			lh.user_id,
			lh.track_id::text,
			'track',
			COUNT(*) * CASE WHEN lh.completed THEN 1.0 ELSE 0.3 END,
			POW(0.95, EXTRACT(EPOCH FROM (NOW() - MAX(lh.played_at))) / 86400.0),
			NOW()
		FROM listening_history lh
		WHERE lh.played_at > NOW() - INTERVAL '90 days'
		GROUP BY lh.user_id, lh.track_id
		ON CONFLICT (user_id, target_id, target_type)
		DO UPDATE SET
			score = GREATEST(user_affinities.score * 0.7 + EXCLUDED.score * 0.3, 0),
			decay = EXCLUDED.decay,
			updated_at = NOW()
	`)
	return err
}

func (w *Worker) updateTrending(ctx context.Context) error {
	var count int
	err := w.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM (
			SELECT t.id, COUNT(*) * (
				1.0 / (1.0 + EXTRACT(EPOCH FROM (NOW() - MIN(lh.played_at))) / 86400.0)
			) AS velocity
			FROM tracks t
			JOIN listening_history lh ON lh.track_id = t.id
			WHERE lh.played_at > NOW() - INTERVAL '24 hours'
			GROUP BY t.id
			HAVING COUNT(*) > 1
			ORDER BY velocity DESC
			LIMIT 100
		) trending
	`).Scan(&count)
	if err != nil {
		return err
	}
	w.logger.Info("trending tracks", zap.Int64("count", int64(count)))
	return nil
}
