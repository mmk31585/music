package catalog

import "github.com/google/uuid"

type CreateArtistRequest struct {
	Name             string  `json:"name" validate:"required,min=1,max=200"`
	Bio              *string `json:"bio" validate:"omitempty,max=5000"`
	ImageURL         *string `json:"imageUrl" validate:"omitempty,url"`
	IsVerified       *bool   `json:"isVerified"`
	MonthlyListeners *int64  `json:"monthlyListeners" validate:"omitempty,min=0"`
}

type UpdateArtistRequest struct {
	Name             *string `json:"name" validate:"omitempty,min=1,max=200"`
	Bio              *string `json:"bio" validate:"omitempty,max=5000"`
	ImageURL         *string `json:"imageUrl" validate:"omitempty,url"`
	IsVerified       *bool   `json:"isVerified"`
	MonthlyListeners *int64  `json:"monthlyListeners" validate:"omitempty,min=0"`
}

type CreateAlbumRequest struct {
	ArtistID    uuid.UUID `json:"artistId" validate:"required"`
	Title       string    `json:"title" validate:"required,min=1,max=250"`
	CoverURL    *string   `json:"coverUrl" validate:"omitempty,url"`
	ReleaseDate *string   `json:"releaseDate" validate:"omitempty"`
	AlbumType   *string   `json:"albumType" validate:"omitempty,oneof=album single ep compilation"`
}

type UpdateAlbumRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=1,max=250"`
	CoverURL    *string `json:"coverUrl" validate:"omitempty,url"`
	ReleaseDate *string `json:"releaseDate" validate:"omitempty"`
	AlbumType   *string `json:"albumType" validate:"omitempty,oneof=album single ep compilation"`
}

type CreateTrackRequest struct {
	ArtistID        uuid.UUID   `json:"artistId" validate:"required"`
	AlbumID         *uuid.UUID  `json:"albumId"`
	Title           string      `json:"title" validate:"required,min=1,max=250"`
	DurationSeconds int         `json:"durationSeconds" validate:"required,min=0"`
	TrackNumber     *int        `json:"trackNumber" validate:"omitempty,min=1"`
	Explicit        *bool       `json:"explicit"`
	AudioURL        *string     `json:"audioUrl" validate:"omitempty,url"`
	CoverURL        *string     `json:"coverUrl" validate:"omitempty,url"`
	IsPublic        *bool       `json:"isPublic"`
	GenreIDs        []uuid.UUID `json:"genreIds"`
}

type UpdateTrackRequest struct {
	AlbumID         *uuid.UUID  `json:"albumId"`
	Title           *string     `json:"title" validate:"omitempty,min=1,max=250"`
	DurationSeconds *int        `json:"durationSeconds" validate:"omitempty,min=0"`
	TrackNumber     *int        `json:"trackNumber" validate:"omitempty,min=1"`
	Explicit        *bool       `json:"explicit"`
	AudioURL        *string     `json:"audioUrl" validate:"omitempty,url"`
	CoverURL        *string     `json:"coverUrl" validate:"omitempty,url"`
	IsPublic        *bool       `json:"isPublic"`
	GenreIDs        []uuid.UUID `json:"genreIds"`
}

type CreateGenreRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type UpdateGenreRequest struct {
	Name *string `json:"name" validate:"omitempty,min=2,max=100"`
}

type SearchResponse struct {
	Artists []Artist `json:"artists"`
	Albums  []Album  `json:"albums"`
	Tracks  []Track  `json:"tracks"`
}
