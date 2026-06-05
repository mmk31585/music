package track

import (
	"time"

	"github.com/google/uuid"
)

type Genre struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type TrackArtist struct {
	ArtistID uuid.UUID `db:"artist_id" json:"artistId"`
	Name     string    `db:"name" json:"name"`
	Slug     string    `db:"slug" json:"slug"`
	Role     string    `db:"role" json:"role"`
	Position int       `db:"position" json:"position"`
}

type TrackArtistRequest struct {
	ArtistID uuid.UUID `json:"artist_id" validate:"required"`
	Role     string    `json:"role"`
	Position int       `json:"position"`
}

type Track struct {
	ID uuid.UUID `db:"id" json:"id"`

	// Deprecated compatibility field.
	// Keep this until frontend/admin payloads are fully migrated to artists[].
	ArtistID uuid.UUID `db:"artist_id" json:"artistId"`

	AlbumID         *uuid.UUID `db:"album_id" json:"albumId,omitempty"`
	Title           string     `db:"title" json:"title"`
	Slug            string     `db:"slug" json:"slug"`
	DurationSeconds int        `db:"duration_seconds" json:"durationSeconds"`
	TrackNumber     *int       `db:"track_number" json:"trackNumber,omitempty"`
	Explicit        bool       `db:"explicit" json:"explicit"`

	AudioURL *string `db:"audio_url" json:"audioUrl,omitempty"`
	CoverURL *string `db:"cover_url" json:"coverUrl,omitempty"`

	AudioMediaID *uuid.UUID `db:"audio_media_id" json:"audioMediaId,omitempty"`
	CoverMediaID *uuid.UUID `db:"cover_media_id" json:"coverMediaId,omitempty"`

	PlayCount int64 `db:"play_count" json:"playCount"`
	IsPublic  bool  `db:"is_public" json:"isPublic"`

	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`

	Artists []TrackArtist `db:"-" json:"artists,omitempty"`
	Genres  []Genre       `db:"-" json:"genres,omitempty"`
}

type CreateRequest struct {
	ArtistID uuid.UUID `db:"artist_id" json:"artist_id"`

	Artists []TrackArtistRequest `json:"artists"`

	// compatibility with current frontend
	ArtistIDs []uuid.UUID          `json:"artist_ids"`
	Credits   []TrackArtistRequest `json:"credits"`

	AlbumID         *uuid.UUID `db:"album_id" json:"album_id"`
	Title           string     `db:"title" json:"title"`
	DurationSeconds int        `db:"duration_seconds" json:"duration_seconds"`
	TrackNumber     *int       `db:"track_number" json:"track_number"`
	Explicit        *bool      `db:"explicit" json:"explicit"`

	AudioURL *string `db:"audio_url" json:"audio_url"`
	CoverURL *string `db:"cover_url" json:"cover_url"`

	AudioMediaID *uuid.UUID `db:"audio_media_id" json:"audio_media_id"`
	CoverMediaID *uuid.UUID `db:"cover_media_id" json:"cover_media_id"`

	IsPublic *bool       `db:"is_public" json:"is_public"`
	GenreIDs []uuid.UUID `db:"genre_ids" json:"genre_ids"`
}

type UpdateRequest struct {
	Artists []TrackArtistRequest `json:"artists"`

	AlbumID         *uuid.UUID `db:"album_id" json:"album_id"`
	Title           *string    `db:"title" json:"title" validate:"omitempty,min=1,max=250"`
	DurationSeconds *int       `db:"duration_seconds" json:"duration_seconds" validate:"omitempty,min=0"`
	TrackNumber     *int       `db:"track_number" json:"track_number" validate:"omitempty,min=1"`
	Explicit        *bool      `db:"explicit" json:"explicit"`

	AudioURL *string `db:"audio_url" json:"audio_url" validate:"omitempty,url"`
	CoverURL *string `db:"cover_url" json:"cover_url" validate:"omitempty,url"`

	AudioMediaID *uuid.UUID `db:"audio_media_id" json:"audio_media_id"`
	CoverMediaID *uuid.UUID `db:"cover_media_id" json:"cover_media_id"`

	IsPublic *bool       `db:"is_public" json:"is_public"`
	GenreIDs []uuid.UUID `db:"genre_ids" json:"genre_ids"`
}
type TrackCredit struct {
	ID         uuid.UUID `db:"id" json:"id"`
	TrackID    uuid.UUID `db:"track_id" json:"trackId"`
	ArtistID   uuid.UUID `db:"artist_id" json:"artistId"`
	ArtistName string    `db:"artist_name" json:"artistName"`
	ArtistSlug string    `db:"artist_slug" json:"artistSlug"`
	CreditType string    `db:"credit_type" json:"creditType"`
	Position   int       `db:"position" json:"position"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
}

type TrackCreditRequest struct {
	ArtistID   uuid.UUID `json:"artist_id" validate:"required"`
	CreditType string    `json:"credit_type" validate:"required"`
	Position   int       `json:"position"`
}
