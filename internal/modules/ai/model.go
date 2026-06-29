package ai

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Float64Array implements sql.Scanner and driver.Valuer for PostgreSQL
// double precision[] arrays when using the pgx stdlib driver.
type Float64Array []float64

func (a *Float64Array) Scan(src any) error {
	if src == nil {
		*a = nil
		return nil
	}
	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("Float64Array.Scan: unexpected type %T", src)
	}
	// PostgreSQL array format: {val1,val2,val3}
	s = strings.Trim(s, "{}")
	if s == "" {
		*a = Float64Array{}
		return nil
	}
	parts := strings.Split(s, ",")
	result := make(Float64Array, len(parts))
	for i, p := range parts {
		f, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return fmt.Errorf("Float64Array.Scan: parse %q: %w", p, err)
		}
		result[i] = f
	}
	*a = result
	return nil
}

func (a Float64Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	parts := make([]string, len(a))
	for i, f := range a {
		parts[i] = strconv.FormatFloat(f, 'f', -1, 64)
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

type TrackEmbedding struct {
	TrackID      uuid.UUID    `db:"track_id" json:"track_id"`
	Embedding    Float64Array `db:"embedding" json:"embedding"`
	ModelVersion string       `db:"model_version" json:"model_version"`
	UpdatedAt    time.Time    `db:"updated_at" json:"updated_at"`
}

type TrackMood struct {
	TrackID          uuid.UUID       `db:"track_id" json:"track_id"`
	MoodTags         json.RawMessage `db:"mood_tags" json:"mood_tags"`
	Energy           float64         `db:"energy" json:"energy"`
	Valence          float64         `db:"valence" json:"valence"`
	Tempo            float64         `db:"tempo" json:"tempo"`
	Danceability     float64         `db:"danceability" json:"danceability"`
	Acousticness     float64         `db:"acousticness" json:"acousticness"`
	Instrumentalness float64         `db:"instrumentalness" json:"instrumentalness"`
	Liveness         float64         `db:"liveness" json:"liveness"`
	Speechiness      float64         `db:"speechiness" json:"speechiness"`
	UpdatedAt        time.Time       `db:"updated_at" json:"updated_at"`
}

type MoodTag struct {
	TrackID    string  `json:"track_id"`
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

type SmartPlaylist struct {
	ID          int64     `db:"id" json:"id"`
	UserID      uuid.UUID `db:"user_id" json:"user_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	QueryConfig string    `db:"query_config" json:"query_config"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type GenerationLog struct {
	ID         int64      `db:"id" json:"id"`
	UserID     uuid.UUID  `db:"user_id" json:"user_id"`
	PlaylistID *uuid.UUID `db:"playlist_id" json:"playlist_id"`
	Prompt     string     `db:"prompt" json:"prompt"`
	TrackCount int        `db:"track_count" json:"track_count"`
	ModelUsed  string     `db:"model_used" json:"model_used"`
	LatencyMs  int        `db:"latency_ms" json:"latency_ms"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}

type TrackMeta struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	ArtistID     string  `json:"artist_id,omitempty"`
	Artist       string  `json:"artist,omitempty"`
	Album        string  `json:"album,omitempty"`
	Genre        string  `json:"genre,omitempty"`
	Year         int     `json:"year,omitempty"`
	Duration     int     `json:"duration,omitempty"`
	CoverURL     string  `db:"cover_url" json:"cover_url,omitempty"`
	Energy       float64 `json:"energy,omitempty"`
	Valence      float64 `json:"valence,omitempty"`
	Tempo        float64 `json:"tempo,omitempty"`
	Danceability float64 `json:"danceability,omitempty"`
}
