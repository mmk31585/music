package artist

import (
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID       uuid.UUID `db:"id" json:"id"`
	Name     string    `db:"name" json:"name"`
	Slug     string    `db:"slug" json:"slug"`
	Bio      *string   `db:"bio" json:"bio,omitempty"`
	ImageURL *string   `db:"image_url" json:"imageUrl,omitempty"`

	AvatarMediaID *uuid.UUID `db:"avatar_media_id" json:"avatarMediaId,omitempty"`
	BannerMediaID *uuid.UUID `db:"banner_media_id" json:"bannerMediaId,omitempty"`
	Country       *string    `db:"country" json:"country,omitempty"`

	IsVerified       bool       `db:"is_verified" json:"isVerified"`
	MonthlyListeners int64      `db:"monthly_listeners" json:"monthlyListeners"`
	CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt        *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

type CreateRequest struct {
	Name     string  `db:"name" json:"name" validate:"required,min=1,max=200"`
	Bio      *string `db:"bio" json:"bio" validate:"omitempty,max=5000"`
	ImageURL *string `db:"image_url" json:"image_url" validate:"omitempty,url"`

	AvatarMediaID *uuid.UUID `db:"avatar_media_id" json:"avatar_media_id"`
	BannerMediaID *uuid.UUID `db:"banner_media_id" json:"banner_media_id"`
	Country       *string    `db:"country" json:"country" validate:"omitempty,len=2"`

	IsVerified       *bool  `db:"is_verified" json:"is_verified"`
	MonthlyListeners *int64 `db:"monthly_listeners" json:"monthly_listeners" validate:"omitempty,min=0"`
}

type UpdateRequest struct {
	Name     *string `db:"name" json:"name" validate:"omitempty,min=1,max=200"`
	Bio      *string `db:"bio" json:"bio" validate:"omitempty,max=5000"`
	ImageURL *string `db:"image_url" json:"image_url" validate:"omitempty,url"`

	AvatarMediaID *uuid.UUID `db:"avatar_media_id" json:"avatar_media_id"`
	BannerMediaID *uuid.UUID `db:"banner_media_id" json:"banner_media_id"`
	Country       *string    `db:"country" json:"country" validate:"omitempty,len=2"`

	IsVerified       *bool  `db:"is_verified" json:"is_verified"`
	MonthlyListeners *int64 `db:"monthly_listeners" json:"monthly_listeners" validate:"omitempty,min=0"`
}

type ArtistTrack struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	Title           string     `db:"title" json:"title"`
	Slug            string     `db:"slug" json:"slug"`
	AlbumID         *uuid.UUID `db:"album_id" json:"albumId,omitempty"`
	DurationSeconds int        `db:"duration_seconds" json:"durationSeconds"`
	TrackNumber     *int       `db:"track_number" json:"trackNumber,omitempty"`
	Explicit        bool       `db:"explicit" json:"explicit"`
	AudioURL        *string    `db:"audio_url" json:"audioUrl,omitempty"`
	CoverURL        *string    `db:"cover_url" json:"coverUrl,omitempty"`
	PlayCount       int64      `db:"play_count" json:"playCount"`
	ArtistRole      string     `db:"artist_role" json:"artistRole"`
	ArtistPosition  int        `db:"artist_position" json:"artistPosition"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
}

type ArtistAlbum struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Title       string     `db:"title" json:"title"`
	Slug        string     `db:"slug" json:"slug"`
	CoverURL    *string    `db:"cover_url" json:"coverUrl,omitempty"`
	ReleaseDate *time.Time `db:"release_date" json:"releaseDate,omitempty"`
	AlbumType   string     `db:"album_type" json:"albumType"`
	ArtistRole  string     `db:"artist_role" json:"artistRole"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
}

type Overview struct {
	Artist    *Artist       `json:"artist"`
	TopTracks []ArtistTrack `json:"topTracks"`
	Albums    []ArtistAlbum `json:"albums"`
	Singles   []ArtistAlbum `json:"singles"`
	AppearsOn []ArtistTrack `json:"appearsOn"`
}
type RelatedArtist struct {
	ID               uuid.UUID  `db:"id" json:"id"`
	Name             string     `db:"name" json:"name"`
	Slug             string     `db:"slug" json:"slug"`
	Bio              *string    `db:"bio" json:"bio,omitempty"`
	ImageURL         *string    `db:"image_url" json:"imageUrl,omitempty"`
	AvatarMediaID    *uuid.UUID `db:"avatar_media_id" json:"avatarMediaId,omitempty"`
	BannerMediaID    *uuid.UUID `db:"banner_media_id" json:"bannerMediaId,omitempty"`
	IsVerified       bool       `db:"is_verified" json:"isVerified"`
	MonthlyListeners int64      `db:"monthly_listeners" json:"monthlyListeners"`

	RelatedScore  float64   `db:"related_score" json:"relatedScore"`
	RelatedSource string    `db:"related_source" json:"relatedSource"`
	RelatedAt     time.Time `db:"related_at" json:"relatedAt"`
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
