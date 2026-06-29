package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) UpsertEmbedding(ctx context.Context, trackID uuid.UUID, embedding []float64, modelVersion string) error {
	query := `
		INSERT INTO track_embeddings_text (track_id, embedding, model_version, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (track_id)
		DO UPDATE SET embedding = $2, model_version = $3, updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, trackID, embedding, modelVersion)
	return err
}

func (r *Repository) GetEmbedding(ctx context.Context, trackID uuid.UUID) (*TrackEmbedding, error) {
	var e TrackEmbedding
	err := r.db.GetContext(ctx, &e, `SELECT track_id, embedding, model_version, updated_at FROM track_embeddings_text WHERE track_id = $1`, trackID)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *Repository) GetTracksWithoutEmbeddings(ctx context.Context, limit int) ([]TrackMeta, error) {
	query := `
		SELECT t.id, t.title FROM tracks t
		LEFT JOIN track_embeddings_text te ON te.track_id = t.id
		WHERE te.track_id IS NULL
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TrackMeta
	for rows.Next() {
		var m TrackMeta
		var idStr string
		if err := rows.Scan(&idStr, &m.Title); err != nil {
			return nil, err
		}
		m.ID = idStr
		results = append(results, m)
	}
	return results, nil
}

func (r *Repository) UpsertMood(ctx context.Context, mood TrackMood) error {
	tagsJSON, err := json.Marshal(mood.MoodTags)
	if err != nil {
		return fmt.Errorf("marshal mood tags: %w", err)
	}
	query := `
		INSERT INTO track_moods (track_id, mood_tags, energy, valence, tempo, danceability, acousticness, instrumentalness, liveness, speechiness, updated_at)
		VALUES ($1, $2::jsonb, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
		ON CONFLICT (track_id)
		DO UPDATE SET mood_tags = $2::jsonb, energy = $3, valence = $4, tempo = $5, danceability = $6,
			acousticness = $7, instrumentalness = $8, liveness = $9, speechiness = $10, updated_at = NOW()
	`
	_, err = r.db.ExecContext(ctx, query,
		mood.TrackID, string(tagsJSON),
		mood.Energy, mood.Valence, mood.Tempo, mood.Danceability,
		mood.Acousticness, mood.Instrumentalness, mood.Liveness, mood.Speechiness,
	)
	return err
}

func (r *Repository) GetMood(ctx context.Context, trackID uuid.UUID) (*TrackMood, error) {
	var m TrackMood
	err := r.db.GetContext(ctx, &m, `
		SELECT track_id, mood_tags, energy, valence, tempo, danceability, acousticness, instrumentalness, liveness, speechiness, updated_at
		FROM track_moods WHERE track_id = $1
	`, trackID)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repository) GetTracksByMood(ctx context.Context, moodName string, limit int) ([]TrackMeta, error) {
	query := `
		SELECT t.id, t.title FROM tracks t
		JOIN track_moods tm ON tm.track_id = t.id
		WHERE tm.mood_tags @> jsonb_build_array(jsonb_build_object('name', $1))
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, moodName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TrackMeta
	for rows.Next() {
		var m TrackMeta
		var idStr string
		if err := rows.Scan(&idStr, &m.Title); err != nil {
			return nil, err
		}
		m.ID = idStr
		results = append(results, m)
	}
	return results, nil
}

func (r *Repository) GetTracksByMoodRange(ctx context.Context, minEnergy, maxEnergy, minValence, maxValence float64, limit int) ([]TrackMeta, error) {
	query := `
		SELECT t.id, t.title FROM tracks t
		JOIN track_moods tm ON tm.track_id = t.id
		WHERE tm.energy BETWEEN $1 AND $2
		AND tm.valence BETWEEN $3 AND $4
		LIMIT $5
	`
	rows, err := r.db.QueryContext(ctx, query, minEnergy, maxEnergy, minValence, maxValence, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TrackMeta
	for rows.Next() {
		var m TrackMeta
		var idStr string
		if err := rows.Scan(&idStr, &m.Title); err != nil {
			return nil, err
		}
		m.ID = idStr
		results = append(results, m)
	}
	return results, nil
}

func (r *Repository) GetTracksWithoutMoods(ctx context.Context, limit int) ([]TrackMeta, error) {
	query := `
		SELECT t.id, t.title FROM tracks t
		LEFT JOIN track_moods tm ON tm.track_id = t.id
		WHERE tm.track_id IS NULL
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TrackMeta
	for rows.Next() {
		var m TrackMeta
		var idStr string
		if err := rows.Scan(&idStr, &m.Title); err != nil {
			return nil, err
		}
		m.ID = idStr
		results = append(results, m)
	}
	return results, nil
}

func (r *Repository) GetTrackMetadata(ctx context.Context, trackID string) (*TrackMeta, error) {
	query := `
		SELECT t.id::text, t.title, COALESCE(a.name, '') as artist,
			COALESCE(al.title, '') as album, COALESCE(g.name, '') as genre,
			COALESCE(t.duration_seconds, 0) as duration,
			COALESCE(t.cover_url, '') as cover_url
		FROM tracks t
		LEFT JOIN track_artists ta ON ta.track_id = t.id AND ta.role = 'primary'
		LEFT JOIN artists a ON a.id = ta.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		LEFT JOIN track_genres tg ON tg.track_id = t.id
		LEFT JOIN genres g ON g.id = tg.genre_id
		WHERE t.id = $1
	`
	var m TrackMeta
	err := r.db.GetContext(ctx, &m, query, trackID)
	if err != nil {
		return nil, err
	}
	m.ID = trackID
	return &m, nil
}

func (r *Repository) GetTracksByIDs(ctx context.Context, ids []string) ([]TrackMeta, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := `
		SELECT t.id::text, t.title, COALESCE(a.name, '') as artist,
			COALESCE(al.title, '') as album, COALESCE(g.name, '') as genre,
			COALESCE(t.duration_seconds, 0) as duration,
			COALESCE(t.cover_url, '') as cover_url
		FROM tracks t
		LEFT JOIN track_artists ta ON ta.track_id = t.id AND ta.role = 'primary'
		LEFT JOIN artists a ON a.id = ta.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		LEFT JOIN track_genres tg ON tg.track_id = t.id
		LEFT JOIN genres g ON g.id = tg.genre_id
		WHERE t.id = ANY($1)
	`
	rows, err := r.db.QueryContext(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TrackMeta
	for rows.Next() {
		var m TrackMeta
		if err := rows.Scan(&m.ID, &m.Title, &m.Artist, &m.Album, &m.Genre, &m.Duration, &m.CoverURL); err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

func (r *Repository) GetSimilarByEmbedding(ctx context.Context, embedding []float64, limit int, embeddingSpace string) ([]TrackMeta, error) {
	tableName := "track_embeddings_audio"
	if embeddingSpace == "text" {
		tableName = "track_embeddings_text"
	}

	query := fmt.Sprintf(`
		SELECT t.id::text, t.title, COALESCE(a.name, '') as artist,
			COALESCE(al.title, '') as album, COALESCE(g.name, '') as genre,
			COALESCE(t.duration_seconds, 0) as duration,
			COALESCE(t.cover_url, '') as cover_url
		FROM %s te
		JOIN tracks t ON t.id = te.track_id
		LEFT JOIN track_artists ta ON ta.track_id = t.id AND ta.role = 'primary'
		LEFT JOIN artists a ON a.id = ta.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		LEFT JOIN track_genres tg ON tg.track_id = t.id
		LEFT JOIN genres g ON g.id = tg.genre_id
		WHERE NOT te.embedding IS NULL
		ORDER BY cosine_distance(te.embedding, $1)
		LIMIT $2
	`, tableName)
	rows, err := r.db.QueryContext(ctx, query, embedding, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TrackMeta
	for rows.Next() {
		var m TrackMeta
		if err := rows.Scan(&m.ID, &m.Title, &m.Artist, &m.Album, &m.Genre, &m.Duration, &m.CoverURL); err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

func (r *Repository) GetTracks(ctx context.Context, limit int) ([]TrackMeta, error) {
	query := `
		SELECT t.id::text, t.title, COALESCE(a.name, '') as artist,
			COALESCE(al.title, '') as album, COALESCE(g.name, '') as genre,
			COALESCE(t.duration_seconds, 0) as duration,
			COALESCE(t.cover_url, '') as cover_url
		FROM tracks t
		LEFT JOIN track_artists ta ON ta.track_id = t.id AND ta.role = 'primary'
		LEFT JOIN artists a ON a.id = ta.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		LEFT JOIN track_genres tg ON tg.track_id = t.id
		LEFT JOIN genres g ON g.id = tg.genre_id
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TrackMeta
	for rows.Next() {
		var m TrackMeta
		if err := rows.Scan(&m.ID, &m.Title, &m.Artist, &m.Album, &m.Genre, &m.Duration, &m.CoverURL); err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

// GetTracksWithMood returns tracks with their mood analysis data (energy, valence, tempo).
// Unlike GetTracks, this LEFT JOINs with track_moods so the returned TrackMeta has mood fields populated.
func (r *Repository) GetTracksWithMood(ctx context.Context, limit int) ([]TrackMeta, error) {
	query := `
		SELECT t.id::text, t.title, COALESCE(a.name, '') as artist,
			COALESCE(al.title, '') as album, COALESCE(g.name, '') as genre,
			COALESCE(t.duration_seconds, 0) as duration,
			COALESCE(t.cover_url, '') as cover_url,
			COALESCE(tm.energy, 0) as energy,
			COALESCE(tm.valence, 0) as valence,
			COALESCE(tm.tempo, 0) as tempo,
			COALESCE(tm.danceability, 0) as danceability
		FROM tracks t
		LEFT JOIN track_artists ta ON ta.track_id = t.id AND ta.role = 'primary'
		LEFT JOIN artists a ON a.id = ta.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		LEFT JOIN track_genres tg ON tg.track_id = t.id
		LEFT JOIN genres g ON g.id = tg.genre_id
		LEFT JOIN track_moods tm ON tm.track_id = t.id
		GROUP BY t.id, a.name, al.title, g.name, tm.energy, tm.valence, tm.tempo, tm.danceability
		ORDER BY t.created_at DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TrackMeta
	for rows.Next() {
		var m TrackMeta
		if err := rows.Scan(&m.ID, &m.Title, &m.Artist, &m.Album, &m.Genre,
			&m.Duration, &m.CoverURL, &m.Energy, &m.Valence, &m.Tempo, &m.Danceability); err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

// GetTracksByGenres returns tracks matching any of the given genre names (case-insensitive).
func (r *Repository) GetTracksByGenres(ctx context.Context, genres []string, limit int) ([]TrackMeta, error) {
	if len(genres) == 0 {
		return nil, nil
	}
	// Build IN clause with parameterized args
	placeholders := make([]string, len(genres))
	args := make([]any, len(genres)+1)
	args[len(genres)] = limit
	for i, g := range genres {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = strings.ToLower(g)
	}
	query := fmt.Sprintf(`
		SELECT t.id::text, t.title, COALESCE(a.name, '') as artist,
			COALESCE(al.title, '') as album, COALESCE(g.name, '') as genre,
			COALESCE(t.duration_seconds, 0) as duration,
			COALESCE(t.cover_url, '') as cover_url
		FROM tracks t
		LEFT JOIN track_artists ta ON ta.track_id = t.id AND ta.role = 'primary'
		LEFT JOIN artists a ON a.id = ta.artist_id
		LEFT JOIN albums al ON al.id = t.album_id
		LEFT JOIN track_genres tg ON tg.track_id = t.id
		LEFT JOIN genres g ON g.id = tg.genre_id
		WHERE LOWER(g.name) IN (%s)
		LIMIT $%d
	`, strings.Join(placeholders, ", "), len(genres)+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TrackMeta
	for rows.Next() {
		var m TrackMeta
		if err := rows.Scan(&m.ID, &m.Title, &m.Artist, &m.Album, &m.Genre, &m.Duration, &m.CoverURL); err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

func (r *Repository) LogGeneration(ctx context.Context, log GenerationLog) error {
	query := `
		INSERT INTO ai_generation_log (user_id, playlist_id, prompt, track_count, model_used, latency_ms)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, log.UserID, log.PlaylistID, log.Prompt, log.TrackCount, log.ModelUsed, log.LatencyMs)
	return err
}

func (r *Repository) SaveSmartPlaylist(ctx context.Context, sp SmartPlaylist) (int64, error) {
	query := `
		INSERT INTO smart_playlists (user_id, name, description, query_config, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id int64
	err := r.db.QueryRowContext(ctx, query, sp.UserID, sp.Name, sp.Description, sp.QueryConfig, sp.IsActive).Scan(&id)
	return id, err
}

func (r *Repository) EmbeddingExists(ctx context.Context, trackID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM track_embeddings_text WHERE track_id = $1)`, trackID).Scan(&exists)
	return exists, err
}
