package playlist

import (
	"time"

	"github.com/google/uuid"
)

type Playlist struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	OwnerID         uuid.UUID `json:"owner_id"`
	Name            string    `json:"name"`
	Description     *string   `json:"description,omitempty"`
	CoverURL        *string   `json:"cover_url,omitempty"`
	IsPublic        bool      `json:"is_public"`
	IsCollaborative bool      `json:"is_collaborative"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type PlaylistTrack struct {
	ID         uuid.UUID `json:"id"`
	PlaylistID uuid.UUID `json:"playlist_id"`
	TrackID    uuid.UUID `json:"track_id"`
	Position   int       `json:"position"`
	AddedAt    time.Time `json:"added_at"`
}

type PlaylistTrackItem struct {
	PlaylistTrackID uuid.UUID `json:"playlist_track_id"`
	TrackID         uuid.UUID `json:"track_id"`
	Position        int       `json:"position"`
	Title           string    `json:"title"`
	ArtistName      *string   `json:"artist_name,omitempty"`
	AlbumTitle      *string   `json:"album_title,omitempty"`
	CoverURL        *string   `json:"cover_url,omitempty"`
	AudioURL        *string   `json:"audio_url,omitempty"`
	DurationSeconds *int      `json:"duration_seconds,omitempty"`
}
