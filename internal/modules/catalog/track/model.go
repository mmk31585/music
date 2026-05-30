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

type Track struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	ArtistID        uuid.UUID  `db:"artist_id" json:"artistId"`
	AlbumID         *uuid.UUID `db:"album_id" json:"albumId,omitempty"`
	Title           string     `db:"title" json:"title"`
	Slug            string     `db:"slug" json:"slug"`
	DurationSeconds int        `db:"duration_seconds" json:"durationSeconds"`
	TrackNumber     *int       `db:"track_number" json:"trackNumber,omitempty"`
	Explicit        bool       `db:"explicit" json:"explicit"`
	AudioURL        *string    `db:"audio_url" json:"audioUrl,omitempty"`
	CoverURL        *string    `db:"cover_url" json:"coverUrl,omitempty"`
	PlayCount       int64      `db:"play_count" json:"playCount"`
	IsPublic        bool       `db:"is_public" json:"isPublic"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
	Genres          []Genre    `db:"-" json:"genres,omitempty"`
}

type CreateRequest struct {
	ArtistID        uuid.UUID   `db:"artist_id" json:"artist_id" validate:"required"`
	AlbumID         *uuid.UUID  `db:"album_id" json:"album_id"`
	Title           string      `db:"title" json:"title" validate:"required,min=1,max=250"`
	DurationSeconds int         `db:"duration_seconds" json:"duration_seconds" validate:"required,min=0"`
	TrackNumber     *int        `db:"track_number" json:"track_number" validate:"omitempty,min=1"`
	Explicit        *bool       `db:"explicit" json:"explicit"`
	AudioURL        *string     `db:"audio_url" json:"audio_url" validate:"omitempty,url"`
	CoverURL        *string     `db:"cover_url" json:"cover_url" validate:"omitempty,url"`
	IsPublic        *bool       `db:"is_public" json:"is_public"`
	GenreIDs        []uuid.UUID `db:"genre_ids" json:"genre_ids"`
}

type UpdateRequest struct {
	AlbumID         *uuid.UUID  `db:"album_id" json:"album_id"`
	Title           *string     `db:"title" json:"title" validate:"omitempty,min=1,max=250"`
	DurationSeconds *int        `db:"duration_seconds" json:"duration_seconds" validate:"omitempty,min=0"`
	TrackNumber     *int        `db:"track_number" json:"track_number" validate:"omitempty,min=1"`
	Explicit        *bool       `db:"explicit" json:"explicit"`
	AudioURL        *string     `db:"audio_url" json:"audio_url" validate:"omitempty,url"`
	CoverURL        *string     `db:"cover_url" json:"cover_url" validate:"omitempty,url"`
	IsPublic        *bool       `db:"is_public" json:"is_public"`
	GenreIDs        []uuid.UUID `db:"genre_ids" json:"genre_ids"`
}
