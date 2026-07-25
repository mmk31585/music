package recommendation

import "time"

type RadioSession struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"user_id" db:"user_id"`
	SeedTrackID    string    `json:"seed_track_id" db:"seed_track_id"`
	SeedType       string    `json:"seed_type" db:"seed_type"`
	PlayedTrackIDs []string  `json:"played_track_ids" db:"played_track_ids"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	LastActiveAt   time.Time `json:"last_active_at" db:"last_active_at"`
}

type RadioTrackItem struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	ArtistName      string  `json:"artist_name"`
	ArtistID        *string `json:"artist_id,omitempty"`
	AlbumTitle      *string `json:"album_title,omitempty"`
	CoverURL        *string `json:"cover_url,omitempty"`
	AudioURL        *string `json:"audio_url,omitempty"`
	DurationSeconds *int    `json:"duration_seconds,omitempty"`
	SimilarityScore float64 `json:"similarity_score,omitempty"`
}

type StartRadioResponse struct {
	SessionID string           `json:"session_id"`
	Tracks    []RadioTrackItem `json:"tracks"`
}

type NextBatchResponse struct {
	Tracks  []RadioTrackItem `json:"tracks"`
	HasMore bool             `json:"has_more"`
}
