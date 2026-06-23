package album

import (
	"time"

	"github.com/google/uuid"
)

// Domain model - db tags only
type Album struct {
	ID uuid.UUID `db:"id"`

	// Deprecated compatibility field.
	ArtistID uuid.UUID `db:"artist_id"`

	Title        string     `db:"title"`
	Slug         string     `db:"slug"`
	CoverURL     *string    `db:"cover_url"`
	CoverMediaID *uuid.UUID `db:"cover_media_id"`

	ReleaseDate *time.Time `db:"release_date"`
	AlbumType   string     `db:"album_type"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	Artists []AlbumArtist `db:"-"`
}

type AlbumArtist struct {
	ArtistID uuid.UUID `db:"artist_id"`
	Name     string    `db:"name"`
	Slug     string    `db:"slug"`
	Role     string    `db:"role"`
	Position int       `db:"position"`
}

type AlbumTrack struct {
	ID              uuid.UUID  `db:"id"`
	AlbumID         *uuid.UUID `db:"album_id"`
	Title           string     `db:"title"`
	Slug            string     `db:"slug"`
	DurationSeconds int        `db:"duration_seconds"`
	TrackNumber     *int       `db:"track_number"`
	Explicit        bool       `db:"explicit"`
	AudioURL        *string    `db:"audio_url"`
	CoverURL        *string    `db:"cover_url"`
	PlayCount       int64      `db:"play_count"`
	IsPublic        bool       `db:"is_public"`
	CreatedAt       time.Time  `db:"created_at"`
}

// CreateRequest is the payload for creating an album.
// JSON tags use camelCase for modern API standards.
type CreateRequest struct {
	// Deprecated compatibility field.
	ArtistID uuid.UUID `json:"artistId"`

	Artists []AlbumArtistRequest `json:"artists"`

	Title        string     `json:"title"`
	CoverURL     *string    `json:"coverUrl"`
	CoverMediaID *uuid.UUID `json:"coverMediaId"`

	ReleaseDate *string `json:"releaseDate"`
	AlbumType   *string `json:"albumType"`
}

// UpdateRequest is the payload for updating an album.
// JSON tags use camelCase for modern API standards.
type UpdateRequest struct {
	Artists []AlbumArtistRequest `json:"artists"`

	Title        *string    `json:"title"`
	CoverURL     *string    `json:"coverUrl"`
	CoverMediaID *uuid.UUID `json:"coverMediaId"`

	ReleaseDate *string `json:"releaseDate"`
	AlbumType   *string `json:"albumType"`
}

type AlbumArtistRequest struct {
	ArtistID uuid.UUID `json:"artistId" validate:"required"`
	Role     string    `json:"role"`
	Position int       `json:"position"`
}

// ListOptions provides optional filters for List queries.
type ListOptions struct {
	ArtistID *uuid.UUID
}
