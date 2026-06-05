package album

import (
	"time"

	"github.com/google/uuid"
)

type AlbumArtist struct {
	ArtistID uuid.UUID `db:"artist_id" json:"artistId"`
	Name     string    `db:"name" json:"name"`
	Slug     string    `db:"slug" json:"slug"`
	Role     string    `db:"role" json:"role"`
	Position int       `db:"position" json:"position"`
}

type AlbumArtistRequest struct {
	ArtistID uuid.UUID `json:"artist_id" validate:"required"`
	Role     string    `json:"role"`
	Position int       `json:"position"`
}

type Album struct {
	ID uuid.UUID `db:"id" json:"id"`

	// Deprecated compatibility field.
	// Keep this while frontend/admin migrates to artists[].
	ArtistID uuid.UUID `db:"artist_id" json:"artistId"`

	Title        string     `db:"title" json:"title"`
	Slug         string     `db:"slug" json:"slug"`
	CoverURL     *string    `db:"cover_url" json:"coverUrl,omitempty"`
	CoverMediaID *uuid.UUID `db:"cover_media_id" json:"coverMediaId,omitempty"`

	ReleaseDate *time.Time `db:"release_date" json:"releaseDate,omitempty"`
	AlbumType   string     `db:"album_type" json:"albumType"`

	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`

	Artists []AlbumArtist `db:"-" json:"artists,omitempty"`
}

type CreateRequest struct {
	// Deprecated compatibility field.
	// If artists[] is empty, this becomes the primary album artist.
	ArtistID uuid.UUID `db:"artist_id" json:"artist_id"`

	Artists []AlbumArtistRequest `json:"artists"`

	Title        string     `db:"title" json:"title" validate:"required,min=1,max=250"`
	CoverURL     *string    `db:"cover_url" json:"cover_url" validate:"omitempty,url"`
	CoverMediaID *uuid.UUID `db:"cover_media_id" json:"cover_media_id"`

	ReleaseDate *string `db:"release_date" json:"release_date"`
	AlbumType   *string `db:"album_type" json:"album_type" validate:"omitempty,oneof=album single ep compilation live"`
}

type UpdateRequest struct {
	Artists []AlbumArtistRequest `json:"artists"`

	Title        *string    `db:"title" json:"title" validate:"omitempty,min=1,max=250"`
	CoverURL     *string    `db:"cover_url" json:"cover_url" validate:"omitempty,url"`
	CoverMediaID *uuid.UUID `db:"cover_media_id" json:"cover_media_id"`

	ReleaseDate *string `db:"release_date" json:"release_date"`
	AlbumType   *string `db:"album_type" json:"album_type" validate:"omitempty,oneof=album single ep compilation live"`
}
type AlbumTrack struct {
	ID              uuid.UUID  `db:"id" json:"id"`
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
}
