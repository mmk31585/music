package social

import (
	"context"

	"github.com/google/uuid"
)

// --- Track Ratings ---

func (r *repository) CreateRating(ctx context.Context, rating *TrackRating) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO track_ratings (id, user_id, track_id, rating, review, created_at, updated_at)
		VALUES (:id, :user_id, :track_id, :rating, :review, :created_at, :updated_at)
		ON CONFLICT (user_id, track_id) DO UPDATE SET rating = EXCLUDED.rating, review = EXCLUDED.review, updated_at = NOW()
	`, rating)
	return err
}

func (r *repository) GetTrackRatings(ctx context.Context, trackID uuid.UUID) ([]TrackRating, error) {
	var items []TrackRating
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM track_ratings WHERE track_id = $1 ORDER BY created_at DESC`, trackID)
	return items, err
}

func (r *repository) GetTrackRatingAverage(ctx context.Context, trackID uuid.UUID) (float64, int, error) {
	var avg float64
	var count int
	err := r.db.GetContext(ctx, &avg, `SELECT COALESCE(AVG(rating), 0) FROM track_ratings WHERE track_id = $1`, trackID)
	if err != nil {
		return 0, 0, err
	}
	err = r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM track_ratings WHERE track_id = $1`, trackID)
	return avg, count, err
}
