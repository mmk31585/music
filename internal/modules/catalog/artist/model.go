package artist

import (
	"time"

	"github.com/google/uuid"
)

// Domain model - db tags only
type Artist struct {
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Slug     string    `db:"slug"`
	Bio      *string   `db:"bio"`
	ImageURL *string   `db:"image_url"`

	AvatarMediaID *uuid.UUID `db:"avatar_media_id"`
	BannerMediaID *uuid.UUID `db:"banner_media_id"`
	Country       *string    `db:"country"`

	IsVerified       bool       `db:"is_verified"`
	MonthlyListeners int64      `db:"monthly_listeners"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        *time.Time `db:"updated_at"`
}

type ArtistTrack struct {
	ID              uuid.UUID  `db:"id"`
	Title           string     `db:"title"`
	Slug            string     `db:"slug"`
	AlbumID         *uuid.UUID `db:"album_id"`
	DurationSeconds int        `db:"duration_seconds"`
	TrackNumber     *int       `db:"track_number"`
	Explicit        bool       `db:"explicit"`
	AudioURL        *string    `db:"audio_url"`
	CoverURL        *string    `db:"cover_url"`
	PlayCount       int64      `db:"play_count"`
	ArtistRole      string     `db:"artist_role"`
	ArtistPosition  int        `db:"artist_position"`
	CreatedAt       time.Time  `db:"created_at"`
}

type ArtistAlbum struct {
	ID          uuid.UUID  `db:"id"`
	Title       string     `db:"title"`
	Slug        string     `db:"slug"`
	CoverURL    *string    `db:"cover_url"`
	ReleaseDate *time.Time `db:"release_date"`
	AlbumType   string     `db:"album_type"`
	ArtistRole  string     `db:"artist_role"`
	CreatedAt   time.Time  `db:"created_at"`
}

type Overview struct {
	Artist    *Artist       `json:"artist"`
	TopTracks []ArtistTrack `json:"top_tracks"`
	Albums    []ArtistAlbum `json:"albums"`
	Singles   []ArtistAlbum `json:"singles"`
	AppearsOn []ArtistTrack `json:"appears_on"`
}

type RelatedArtist struct {
	ID               uuid.UUID  `db:"id"`
	Name             string     `db:"name"`
	Slug             string     `db:"slug"`
	Bio              *string    `db:"bio"`
	ImageURL         *string    `db:"image_url"`
	AvatarMediaID    *uuid.UUID `db:"avatar_media_id"`
	BannerMediaID    *uuid.UUID `db:"banner_media_id"`
	IsVerified       bool       `db:"is_verified"`
	MonthlyListeners int64      `db:"monthly_listeners"`
	RelatedScore     float64    `db:"related_score"`
	RelatedSource    string     `db:"related_source"`
	RelatedAt        time.Time  `db:"related_at"`
}

// Request DTOs with camelCase JSON tags

type CreateRequest struct {
	Name     string  `json:"name"`
	Bio      *string `json:"bio"`
	ImageURL *string `json:"image_url"`

	AvatarMediaID *uuid.UUID `json:"avatar_media_id"`
	BannerMediaID *uuid.UUID `json:"banner_media_id"`
	Country       *string    `json:"country"`

	IsVerified       *bool  `json:"is_verified"`
	MonthlyListeners *int64 `json:"monthly_listeners"`
}

type UpdateRequest struct {
	Name     *string `json:"name"`
	Bio      *string `json:"bio"`
	ImageURL *string `json:"image_url"`

	AvatarMediaID *uuid.UUID `json:"avatar_media_id"`
	BannerMediaID *uuid.UUID `json:"banner_media_id"`
	Country       *string    `json:"country"`

	IsVerified       *bool  `json:"is_verified"`
	MonthlyListeners *int64 `json:"monthly_listeners"`
}

type RelatedArtistRequest struct {
	RelatedArtistID uuid.UUID `json:"related_artist_id" validate:"required"`
	Score           float64   `json:"score"`
	Source          string    `json:"source"`
}

type ArtistTopTrackRequest struct {
	TrackID  uuid.UUID `json:"track_id" validate:"required"`
	Position int       `json:"position"`
}
