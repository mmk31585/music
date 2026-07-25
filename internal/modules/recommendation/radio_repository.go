package recommendation

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RadioRepository interface {
	CreateSession(ctx context.Context, session *RadioSession) error
	GetSession(ctx context.Context, sessionID string) (*RadioSession, error)
	GetActiveSession(ctx context.Context, userID string) (*RadioSession, error)
	AppendPlayedTrack(ctx context.Context, sessionID string, trackID string) error
	UpdateLastActive(ctx context.Context, sessionID string) error
	GetPlayedTrackIDs(ctx context.Context, sessionID string) ([]string, error)
	GetTracksByIDs(ctx context.Context, trackIDs []string) ([]RadioTrackItem, error)
	GetPopularTracks(ctx context.Context, genre string, limit int) ([]RadioTrackItem, error)
}

type radioRepo struct {
	db *sqlx.DB
}

func NewRadioRepository(db *sqlx.DB) RadioRepository {
	return &radioRepo{db: db}
}

func marshalIDs(ids []string) string {
	if ids == nil {
		ids = []string{}
	}
	b, _ := json.Marshal(ids)
	return string(b)
}

func unmarshalIDs(data []byte) []string {
	var ids []string
	if len(data) == 0 {
		return []string{}
	}
	_ = json.Unmarshal(data, &ids)
	if ids == nil {
		return []string{}
	}
	return ids
}

func (r *radioRepo) CreateSession(ctx context.Context, session *RadioSession) error {
	session.ID = uuid.New().String()
	session.CreatedAt = time.Now().UTC()
	session.LastActiveAt = session.CreatedAt
	if session.SeedType == "" {
		session.SeedType = "track"
	}
	if session.PlayedTrackIDs == nil {
		session.PlayedTrackIDs = []string{}
	}

	query := `
		INSERT INTO radio_sessions (id, user_id, seed_track_id, seed_type, played_track_ids, created_at, last_active_at)
		VALUES ($1, $2, $3, $4, $5::jsonb, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.UserID, session.SeedTrackID, session.SeedType,
		marshalIDs(session.PlayedTrackIDs), session.CreatedAt, session.LastActiveAt,
	)
	return err
}

func (r *radioRepo) GetSession(ctx context.Context, sessionID string) (*RadioSession, error) {
	query := `
		SELECT id, user_id, seed_track_id, seed_type, played_track_ids, created_at, last_active_at
		FROM radio_sessions
		WHERE id = $1
	`
	var s RadioSession
	var playedRaw []byte
	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&s.ID, &s.UserID, &s.SeedTrackID, &s.SeedType,
		&playedRaw, &s.CreatedAt, &s.LastActiveAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	s.PlayedTrackIDs = unmarshalIDs(playedRaw)
	return &s, nil
}

func (r *radioRepo) GetActiveSession(ctx context.Context, userID string) (*RadioSession, error) {
	query := `
		SELECT id, user_id, seed_track_id, seed_type, played_track_ids, created_at, last_active_at
		FROM radio_sessions
		WHERE user_id = $1
		ORDER BY last_active_at DESC
		LIMIT 1
	`
	var s RadioSession
	var playedRaw []byte
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&s.ID, &s.UserID, &s.SeedTrackID, &s.SeedType,
		&playedRaw, &s.CreatedAt, &s.LastActiveAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	s.PlayedTrackIDs = unmarshalIDs(playedRaw)
	return &s, nil
}

func (r *radioRepo) AppendPlayedTrack(ctx context.Context, sessionID string, trackID string) error {
	// Get current ids, append, save back
	session, err := r.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("session not found")
	}

	// Deduplicate
	for _, id := range session.PlayedTrackIDs {
		if id == trackID {
			return r.UpdateLastActive(ctx, sessionID)
		}
	}

	session.PlayedTrackIDs = append(session.PlayedTrackIDs, trackID)
	query := `UPDATE radio_sessions SET played_track_ids = $1::jsonb, last_active_at = NOW() WHERE id = $2`
	_, err = r.db.ExecContext(ctx, query, marshalIDs(session.PlayedTrackIDs), sessionID)
	return err
}

func (r *radioRepo) UpdateLastActive(ctx context.Context, sessionID string) error {
	query := `UPDATE radio_sessions SET last_active_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, sessionID)
	return err
}

func (r *radioRepo) GetPlayedTrackIDs(ctx context.Context, sessionID string) ([]string, error) {
	query := `SELECT played_track_ids FROM radio_sessions WHERE id = $1`
	var playedRaw []byte
	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(&playedRaw)
	if err != nil {
		if err == sql.ErrNoRows {
			return []string{}, nil
		}
		return nil, err
	}
	return unmarshalIDs(playedRaw), nil
}

func (r *radioRepo) GetTracksByIDs(ctx context.Context, trackIDs []string) ([]RadioTrackItem, error) {
	if len(trackIDs) == 0 {
		return []RadioTrackItem{}, nil
	}

	placeholders := make([]string, len(trackIDs))
	args := make([]any, len(trackIDs))
	for i, id := range trackIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT t.id, t.title,
			COALESCE(a.name, '') AS artist_name,
			t.artist_id,
			al.title AS album_title,
			t.cover_url,
			t.audio_url,
			t.duration_seconds
		FROM tracks t
		LEFT JOIN artists a ON a.id = t.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		WHERE t.id IN (%s)
	`, strings.Join(placeholders, ","))

	items := make([]RadioTrackItem, 0)
	err := r.db.SelectContext(ctx, &items, query, args...)
	return items, err
}

func (r *radioRepo) GetPopularTracks(ctx context.Context, genre string, limit int) ([]RadioTrackItem, error) {
	var query string
	var args []any

	if genre != "" {
		query = `
			SELECT t.id, t.title,
				COALESCE(a.name, '') AS artist_name,
				t.artist_id,
				al.title AS album_title,
				t.cover_url,
				t.audio_url,
				t.duration_seconds
			FROM tracks t
			LEFT JOIN artists a ON a.id = t.artist_id
			LEFT JOIN albums al ON al.id = t.album_id
			WHERE EXISTS (
				SELECT 1 FROM track_genres tg
				JOIN genres g ON g.id = tg.genre_id
				WHERE tg.track_id = t.id AND g.name = $1
			)
			ORDER BY t.play_count DESC NULLS LAST, t.title ASC
			LIMIT $2
		`
		args = []any{genre, limit}
	} else {
		query = `
			SELECT t.id, t.title,
				COALESCE(a.name, '') AS artist_name,
				t.artist_id,
				al.title AS album_title,
				t.cover_url,
				t.audio_url,
				t.duration_seconds
			FROM tracks t
			LEFT JOIN artists a ON a.id = t.artist_id
			LEFT JOIN albums al ON al.id = t.album_id
			ORDER BY t.play_count DESC NULLS LAST, t.title ASC
			LIMIT $1
		`
		args = []any{limit}
	}

	items := make([]RadioTrackItem, 0)
	err := r.db.SelectContext(ctx, &items, query, args...)
	return items, err
}
