package catalog

import (
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID               uuid.UUID  `db:"id" json:"id"`
	Name             string     `db:"name" json:"name"`
	Slug             string     `db:"slug" json:"slug"`
	Bio              *string    `db:"bio" json:"bio,omitempty"`
	ImageURL         *string    `db:"image_url" json:"imageUrl,omitempty"`
	IsVerified       bool       `db:"is_verified" json:"isVerified"`
	MonthlyListeners int64      `db:"monthly_listeners" json:"monthlyListeners"`
	CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt        *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

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

	Genres []Genre `db:"-" json:"genres,omitempty"`
}

type Genre struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
