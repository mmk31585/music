package recommendation

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type DiscoverWeeklyRepository interface {
	GetByWeek(ctx context.Context, userID string, weekOf time.Time) (*DiscoverWeeklyPlaylist, error)
	Save(ctx context.Context, playlist *DiscoverWeeklyPlaylist) error
}

type discoverWeeklyRepo struct {
	db *sqlx.DB
}

func NewDiscoverWeeklyRepository(db *sqlx.DB) DiscoverWeeklyRepository {
	return &discoverWeeklyRepo{db: db}
}

func (r *discoverWeeklyRepo) GetByWeek(ctx context.Context, userID string, weekOf time.Time) (*DiscoverWeeklyPlaylist, error) {
	var p DiscoverWeeklyPlaylist
	var trackIDs []string

	weekDate := weekOf.Format("2006-01-02")
	err := r.db.GetContext(ctx, &p, `
		SELECT user_id::text, generated_at, week_of
		FROM discover_weekly_playlists
		WHERE user_id = $1 AND week_of = $2::date
	`, userID, weekDate)
	if err != nil {
		return nil, err
	}

	// Fetch track_ids separately since pgx/pq handle UUID arrays differently
	err = r.db.GetContext(ctx, &trackIDs, `
		SELECT track_ids::text[]
		FROM discover_weekly_playlists
		WHERE user_id = $1 AND week_of = $2::date
	`, userID, weekDate)
	if err != nil {
		return nil, err
	}

	p.TrackIDs = trackIDs
	return &p, nil
}

func (r *discoverWeeklyRepo) Save(ctx context.Context, p *DiscoverWeeklyPlaylist) error {
	weekDate := p.WeekOf.Format("2006-01-02")
	trackIDs := p.TrackIDs
	if trackIDs == nil {
		trackIDs = []string{}
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO discover_weekly_playlists (user_id, track_ids, generated_at, week_of)
		VALUES ($1, $2, $3, $4::date)
		ON CONFLICT (user_id, week_of) DO UPDATE SET
			track_ids = EXCLUDED.track_ids,
			generated_at = EXCLUDED.generated_at
	`, p.UserID, trackIDs, p.GeneratedAt, weekDate)
	if err != nil {
		return fmt.Errorf("save weekly playlist: %w", err)
	}
	return nil
}
