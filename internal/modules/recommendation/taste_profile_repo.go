package recommendation

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type TasteProfileRepository interface {
	GetProfile(ctx context.Context, userID string) (*TasteProfile, error)
	UpsertProfile(ctx context.Context, profile *TasteProfile) error
	GetSignalRows(ctx context.Context, userID string, since time.Time) ([]signalRow, error)
	GetRecentPositiveTrackIDs(ctx context.Context, userID string, limit int) ([]string, error)
	SetOnboardingGenres(ctx context.Context, userID string, genreIDs []string) error
}

type tasteProfileRepo struct {
	db *sqlx.DB
}

func NewTasteProfileRepository(db *sqlx.DB) TasteProfileRepository {
	return &tasteProfileRepo{db: db}
}

func (r *tasteProfileRepo) GetProfile(ctx context.Context, userID string) (*TasteProfile, error) {
	var p TasteProfile
	err := r.db.GetContext(ctx, &p, `
		SELECT
			user_id,
			top_genre_ids,
			top_artist_ids,
			seed_track_ids,
			avg_tempo_preference,
			onboarding_genre_ids,
			last_computed_at
		FROM taste_profiles
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *tasteProfileRepo) UpsertProfile(ctx context.Context, profile *TasteProfile) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO taste_profiles (
			user_id, top_genre_ids, top_artist_ids, seed_track_ids,
			avg_tempo_preference, onboarding_genre_ids, last_computed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id) DO UPDATE SET
			top_genre_ids = EXCLUDED.top_genre_ids,
			top_artist_ids = EXCLUDED.top_artist_ids,
			seed_track_ids = EXCLUDED.seed_track_ids,
			avg_tempo_preference = EXCLUDED.avg_tempo_preference,
			last_computed_at = EXCLUDED.last_computed_at
	`,
		profile.UserID,
		profile.TopGenreIDs,
		profile.TopArtistIDs,
		profile.SeedTrackIDs,
		profile.AvgTempoPreference,
		profile.OnboardingGenreIDs,
		profile.LastComputedAt,
	)
	return err
}

func (r *tasteProfileRepo) GetSignalRows(ctx context.Context, userID string, since time.Time) ([]signalRow, error) {
	var rows []signalRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT
			lh.id,
			lh.user_id::text,
			lh.track_id::text,
			lh.played_at,
			lh.played_duration_ms,
			lh.track_duration_ms,
			lh.completion_percent,
			lh.signal_type,
			lh.is_explicit_like,
			t.artist_id::text AS artist_id,
			(SELECT string_agg(DISTINCT tg.genre_id::text, ',') FROM track_genres tg WHERE tg.track_id = lh.track_id) AS genre_id
		FROM listening_history lh
		JOIN tracks t ON t.id = lh.track_id
		WHERE lh.user_id = $1
		  AND lh.played_at >= $2
		ORDER BY lh.played_at DESC
	`, userID, since)
	if err != nil {
		return nil, fmt.Errorf("get signal rows: %w", err)
	}
	return rows, nil
}

func (r *tasteProfileRepo) GetRecentPositiveTrackIDs(ctx context.Context, userID string, limit int) ([]string, error) {
	var ids []string
	err := r.db.SelectContext(ctx, &ids, `
		SELECT lh.track_id::text
		FROM listening_history lh
		WHERE lh.user_id = $1
		  AND lh.signal_type IN ('complete_positive', 'replay_strong')
		ORDER BY lh.played_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *tasteProfileRepo) SetOnboardingGenres(ctx context.Context, userID string, genreIDs []string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO taste_profiles (user_id, onboarding_genre_ids, last_computed_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			onboarding_genre_ids = EXCLUDED.onboarding_genre_ids,
			last_computed_at = NOW()
	`, userID, genreIDs)
	return err
}
