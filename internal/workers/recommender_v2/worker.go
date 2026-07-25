package recommender_v2

import (
	"context"
	"math"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Worker struct {
	db     *sqlx.DB
	rdb    *redis.Client
	logger *zap.Logger
}

func NewWorker(db *sqlx.DB, rdb *redis.Client, logger *zap.Logger) *Worker {
	return &Worker{db: db, rdb: rdb, logger: logger}
}

func (w *Worker) Name() string { return "recommender_v2" }

func (w *Worker) Run(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			w.logger.Error("recommender_v2 worker panicked", zap.Any("recover", r))
		}
	}()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.logger.Info("running collaborative filtering")
			if err := w.computePlayCooccurrence(ctx); err != nil {
				w.logger.Error("cooccurrence computation failed", zap.Error(err))
			}
			if err := w.computeUserAffinities(ctx); err != nil {
				w.logger.Error("affinity computation failed", zap.Error(err))
			}
			if err := w.updateTrendingCache(ctx); err != nil {
				w.logger.Error("trending update failed", zap.Error(err))
			}
		}
	}
}

func (w *Worker) computePlayCooccurrence(ctx context.Context) error {
	_, err := w.db.ExecContext(ctx, `
		INSERT INTO track_similarities (track_id, similar_track_id, similarity, method)
		SELECT
			a.track_id,
			b.track_id,
			COUNT(*)::float / NULLIF(
				(SELECT COUNT(*) FROM listening_history WHERE track_id = a.track_id AND played_at > NOW() - INTERVAL '7 days'),
			0) AS similarity,
			'cooccurrence'
		FROM listening_history a
		JOIN listening_history b ON a.user_id = b.user_id AND a.track_id != b.track_id
		WHERE a.played_at > NOW() - INTERVAL '7 days'
		AND b.played_at > NOW() - INTERVAL '7 days'
		GROUP BY a.track_id, b.track_id
		HAVING COUNT(*) > 1
		ON CONFLICT (track_id, similar_track_id)
		DO UPDATE SET similarity = EXCLUDED.similarity, updated_at = NOW()
	`)
	return err
}

func (w *Worker) computeUserAffinities(ctx context.Context) error {
	_, err := w.db.ExecContext(ctx, `
		INSERT INTO user_affinities (user_id, target_id, target_type, score, decay, updated_at)
		SELECT
			lh.user_id,
			lh.track_id::text,
			'track',
			COUNT(*) * CASE WHEN lh.completed THEN 1.0 ELSE 0.3 END,
			1.0,
			NOW()
		FROM listening_history lh
		WHERE lh.played_at > NOW() - INTERVAL '30 days'
		GROUP BY lh.user_id, lh.track_id
		ON CONFLICT (user_id, target_id, target_type)
		DO UPDATE SET
			score = user_affinities.score * 0.9 + EXCLUDED.score * 0.1,
			decay = 1.0,
			updated_at = NOW()
	`)
	return err
}

func (w *Worker) updateTrendingCache(ctx context.Context) error {
	if w.rdb == nil {
		return nil
	}

	rows, err := w.db.QueryContext(ctx, `
		SELECT t.id, COUNT(*) * (
			1.0 / (1.0 + EXTRACT(EPOCH FROM (NOW() - MIN(lh.played_at))) / 86400.0)
		) AS velocity
		FROM tracks t
		JOIN listening_history lh ON lh.track_id = t.id
		WHERE lh.played_at > NOW() - INTERVAL '24 hours'
		GROUP BY t.id
		ORDER BY velocity DESC
		LIMIT 100
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	pipe := w.rdb.Pipeline()
	pipe.Del(ctx, "trending:tracks")

	for rows.Next() {
		var trackID string
		var velocity float64
		if err := rows.Scan(&trackID, &velocity); err != nil {
			continue
		}
		pipe.ZAdd(ctx, "trending:tracks", redis.Z{
			Score:  math.Round(velocity*1000) / 1000,
			Member: trackID,
		})
	}

	pipe.Expire(ctx, "trending:tracks", 1*time.Hour)
	_, err = pipe.Exec(ctx)
	return err
}
