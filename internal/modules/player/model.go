package player

import "time"

type PlaybackTrack struct {
	ID              string     `json:"id" db:"id"`
	Title           string     `json:"title" db:"title"`
	ArtistName      string     `json:"artistName" db:"artist_name"`
	AlbumTitle      *string    `json:"albumTitle,omitempty" db:"album_title"`
	AudioURL        string     `json:"audioUrl" db:"audio_url"`
	CoverURL        *string    `json:"coverUrl,omitempty" db:"cover_url"`
	DurationSeconds *int       `json:"durationSeconds,omitempty" db:"duration_seconds"`
	IsPublic        bool       `json:"isPublic" db:"is_public"`
	CreatedAt       *time.Time `json:"createdAt,omitempty" db:"created_at"`
}
