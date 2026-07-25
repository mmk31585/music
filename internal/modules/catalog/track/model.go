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

// UpdateRequest for updating a track.
// Nullable fields use double-pointer types (**T) to distinguish three states:
//   - ptr == nil        → field not in JSON → skip (keep current value)
//   - *ptr == nil       → field is explicit JSON null → set to SQL NULL
//   - *ptr != nil        → field has a value → set to *ptr
//
// ClearedFields provides an additional explicit mechanism: list column names
// that should be set to NULL in the database. This is useful when the caller
// cannot send the field at all (e.g. PATCH semantics) but still wants to clear it.
//
// Non-nullable fields use single-pointer (*T) for optional updates (nil = skip, non-nil = set).
type UpdateRequest struct {
	Artists []TrackArtistRequest `json:"artists"`

	// ClearedFields lists nullable column names to set to NULL.
	// Must be processed BEFORE the dynamic builder so that a non-null value
	// in the same request overrides the NULL.
	ClearedFields []string `json:"cleared_fields"`

	AlbumID         **uuid.UUID `json:"album_id"`
	Title           *string     `json:"title"`
	DurationSeconds *int        `json:"duration_seconds"`
	TrackNumber     **int       `json:"track_number"`
	Explicit        *bool       `json:"explicit"`

	AudioURL     **string   `json:"audio_url"`
	CoverURL     **string   `json:"cover_url"`
	AudioMediaID **uuid.UUID `json:"audio_media_id"`
	CoverMediaID **uuid.UUID `json:"cover_media_id"`

	IsPublic *bool       `json:"is_public"`
	GenreIDs []uuid.UUID `json:"genre_ids"`
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
