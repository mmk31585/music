package ai

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TrackEmbedding struct {
	TrackID      uuid.UUID `db:"track_id" json:"track_id"`
	Embedding    []float64 `db:"embedding" json:"embedding"`
	ModelVersion string    `db:"model_version" json:"model_version"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
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
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	ArtistID string  `json:"artist_id,omitempty"`
	Artist   string  `json:"artist,omitempty"`
	Album    string  `json:"album,omitempty"`
	Genre    string  `json:"genre,omitempty"`
	Year     int     `json:"year,omitempty"`
	Duration int     `json:"duration,omitempty"`
	Energy   float64 `json:"energy,omitempty"`
	Valence  float64 `json:"valence,omitempty"`
}
