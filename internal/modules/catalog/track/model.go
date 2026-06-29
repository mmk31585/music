package track

import (
	"time"

	"github.com/google/uuid"
)

// Domain models - db tags only
type Genre struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Slug      string    `db:"slug"`
	CreatedAt time.Time `db:"created_at"`
}

type TrackArtist struct {
	ArtistID uuid.UUID `db:"artist_id"`
	Name     string    `db:"name"`
	Slug     string    `db:"slug"`
	Role     string    `db:"role"`
	Position int       `db:"position"`
}

type Track struct {
	ID uuid.UUID `db:"id"`

	// Deprecated compatibility field.
	ArtistID uuid.UUID `db:"artist_id"`

	AlbumID         *uuid.UUID `db:"album_id"`
	Title           string     `db:"title"`
	Slug            string     `db:"slug"`
	DurationSeconds int        `db:"duration_seconds"`
	TrackNumber     *int       `db:"track_number"`
	Explicit        bool       `db:"explicit"`

	AudioURL *string `db:"audio_url"`
	CoverURL *string `db:"cover_url"`

	AudioMediaID *uuid.UUID `db:"audio_media_id"`
	CoverMediaID *uuid.UUID `db:"cover_media_id"`

	PlayCount int64 `db:"play_count"`
	IsPublic  bool  `db:"is_public"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	Artists []TrackArtist `db:"-"`
	Genres  []Genre       `db:"-"`
}

type TrackCredit struct {
	ID         uuid.UUID `db:"id"`
	TrackID    uuid.UUID `db:"track_id"`
	ArtistID   uuid.UUID `db:"artist_id"`
	ArtistName string    `db:"artist_name"`
	ArtistSlug string    `db:"artist_slug"`
	CreditType string    `db:"credit_type"`
	Position   int       `db:"position"`
	CreatedAt  time.Time `db:"created_at"`
}

type CreateRequest struct {
	ArtistID uuid.UUID `json:"artist_id"`

	Artists []TrackArtistRequest `json:"artists"`

	ArtistIDs []uuid.UUID          `json:"artist_ids"`
	Credits   []TrackArtistRequest `json:"credits"`

	AlbumID         *uuid.UUID `json:"album_id"`
	Title           string     `json:"title"`
	DurationSeconds int        `json:"duration_seconds"`
	TrackNumber     *int       `json:"track_number"`
	Explicit        *bool      `json:"explicit"`

	AudioURL *string `json:"audio_url"`
	CoverURL *string `json:"cover_url"`

	AudioMediaID *uuid.UUID `json:"audio_media_id"`
	CoverMediaID *uuid.UUID `json:"cover_media_id"`

	IsPublic *bool       `json:"is_public"`
	GenreIDs []uuid.UUID `json:"genre_ids"`
}

type UpdateRequest struct {
	Artists []TrackArtistRequest `json:"artists"`

	AlbumID         *uuid.UUID `json:"album_id"`
	Title           *string    `json:"title"`
	DurationSeconds *int       `json:"duration_seconds"`
	TrackNumber     *int       `json:"track_number"`
	Explicit        *bool      `json:"explicit"`

	AudioURL *string `json:"audio_url"`
	CoverURL *string `json:"cover_url"`

	AudioMediaID *uuid.UUID `json:"audio_media_id"`
	CoverMediaID *uuid.UUID `json:"cover_media_id"`

	IsPublic *bool       `json:"is_public"`
	GenreIDs []uuid.UUID `json:"genre_ids"`

	ClearFields []string `json:"clear_fields,omitempty"`
}

// TrackArtistRequest is used in CreateRequest and UpdateRequest.
type TrackArtistRequest struct {
	ArtistID uuid.UUID `json:"artist_id" validate:"required"`
	Role     string    `json:"role"`
	Position int       `json:"position"`
}

// TrackCreditRequest is used for credit replacement.
type TrackCreditRequest struct {
	ArtistID   uuid.UUID `json:"artist_id" validate:"required"`
	CreditType string    `json:"credit_type" validate:"required"`
	Position   int       `json:"position"`
}

// ListOptions provides optional filters for List queries.
type ListOptions struct {
	AlbumID  *uuid.UUID
	ArtistID *uuid.UUID
	Query    string
}
