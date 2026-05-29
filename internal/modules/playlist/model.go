package playlist

import "time"

type Playlist struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CoverURL    *string   `json:"cover_url,omitempty"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PlaylistTrack struct {
	ID         int64     `json:"id"`
	PlaylistID int64     `json:"playlist_id"`
	TrackID    int64     `json:"track_id"`
	Position   int       `json:"position"`
	AddedAt    time.Time `json:"added_at"`
}

type PlaylistTrackItem struct {
	PlaylistTrackID int64   `json:"playlist_track_id"`
	TrackID         int64   `json:"track_id"`
	Position        int     `json:"position"`
	Title           string  `json:"title"`
	ArtistName      *string `json:"artist_name,omitempty"`
	AlbumTitle      *string `json:"album_title,omitempty"`
	CoverURL        *string `json:"cover_url,omitempty"`
	AudioURL        *string `json:"audio_url,omitempty"`
	DurationSeconds *int    `json:"duration_seconds,omitempty"`
}
