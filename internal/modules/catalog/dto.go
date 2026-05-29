package catalog

import "github.com/google/uuid"

type CreateArtistRequest struct {
	Name             string  `db:"name" json:"name" validate:"required,min=1,max=200"`
	Bio              *string `db:"bio" json:"bio" validate:"omitempty,max=5000"`
	ImageURL         *string `db:"image_url" json:"image_url" validate:"omitempty,url"`
	IsVerified       *bool   `db:"is_verified" json:"is_verified"`
	MonthlyListeners *int64  `db:"monthly_listeners" json:"monthly_listeners" validate:"omitempty,min=0"`
}

type UpdateArtistRequest struct {
	Name             *string `db:"name" json:"name" validate:"omitempty,min=1,max=200"`
	Bio              *string `db:"bio" json:"bio" validate:"omitempty,max=5000"`
	ImageURL         *string `db:"image_url" json:"image_url" validate:"omitempty,url"`
	IsVerified       *bool   `db:"is_verified" json:"is_verified"`
	MonthlyListeners *int64  `db:"monthly_listeners" json:"monthly_listeners" validate:"omitempty,min=0"`
}

type CreateAlbumRequest struct {
	ArtistID    uuid.UUID `db:"artist_id" json:"artist_id" validate:"required"`
	Title       string    `db:"title" json:"title" validate:"required,min=1,max=250"`
	CoverURL    *string   `db:"cover_url" json:"cover_url" validate:"omitempty,url"`
	ReleaseDate *string   `db:"release_date" json:"release_date" validate:"omitempty"`
	AlbumType   *string   `db:"album_type" json:"album_type" validate:"omitempty,oneof=album single ep compilation"`
}

type UpdateAlbumRequest struct {
	Title       *string `db:"title" json:"title" validate:"omitempty,min=1,max=250"`
	CoverURL    *string `db:"cover_url" json:"cover_url" validate:"omitempty,url"`
	ReleaseDate *string `db:"release_date" json:"release_date" validate:"omitempty"`
	AlbumType   *string `db:"album_type" json:"album_type" validate:"omitempty,oneof=album single ep compilation"`
}

type CreateTrackRequest struct {
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

type UpdateTrackRequest struct {
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

type CreateGenreRequest struct {
	Name string `db:"name" json:"name" validate:"required,min=2,max=100"`
}

type UpdateGenreRequest struct {
	Name *string `db:"name" json:"name" validate:"omitempty,min=2,max=100"`
}

type SearchResponse struct {
	Artists []Artist `db:"artists" json:"artists"`
	Albums  []Album  `db:"albums" json:"albums"`
	Tracks  []Track  `db:"tracks" json:"tracks"`
}
