package album

import (
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	ArtistID    uuid.UUID  `db:"artist_id" json:"artistId"`
	Title       string     `db:"title" json:"title"`
	Slug        string     `db:"slug" json:"slug"`
	CoverURL    *string    `db:"cover_url" json:"coverUrl,omitempty"`
	ReleaseDate *time.Time `db:"release_date" json:"releaseDate,omitempty"`
	AlbumType   string     `db:"album_type" json:"albumType"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

type CreateRequest struct {
	ArtistID    uuid.UUID `db:"artist_id" json:"artist_id" validate:"required"`
	Title       string    `db:"title" json:"title" validate:"required,min=1,max=250"`
	CoverURL    *string   `db:"cover_url" json:"cover_url" validate:"omitempty,url"`
	ReleaseDate *string   `db:"release_date" json:"release_date"`
	AlbumType   *string   `db:"album_type" json:"album_type" validate:"omitempty,oneof=album single ep compilation"`
}

type UpdateRequest struct {
	Title       *string `db:"title" json:"title" validate:"omitempty,min=1,max=250"`
	CoverURL    *string `db:"cover_url" json:"cover_url" validate:"omitempty,url"`
	ReleaseDate *string `db:"release_date" json:"release_date"`
	AlbumType   *string `db:"album_type" json:"album_type" validate:"omitempty,oneof=album single ep compilation"`
}
