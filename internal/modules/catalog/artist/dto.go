package artist

import (
	"time"

	"github.com/google/uuid"
)

// ---------- Response DTOs ----------

type ArtistResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Bio      *string   `json:"bio,omitempty"`
	ImageURL *string   `json:"imageUrl,omitempty"`

	AvatarMediaID *uuid.UUID `json:"avatarMediaId,omitempty"`
	BannerMediaID *uuid.UUID `json:"bannerMediaId,omitempty"`
	Country       *string    `json:"country,omitempty"`

	IsVerified       bool       `json:"isVerified"`
	MonthlyListeners int64      `json:"monthlyListeners"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
}

type ArtistTrackResponse struct {
	ID              uuid.UUID  `json:"id"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	AlbumID         *uuid.UUID `json:"albumId,omitempty"`
	DurationSeconds int        `json:"durationSeconds"`
	TrackNumber     *int       `json:"trackNumber,omitempty"`
	Explicit        bool       `json:"explicit"`
	AudioURL        *string    `json:"audioUrl,omitempty"`
	CoverURL        *string    `json:"coverUrl,omitempty"`
	PlayCount       int64      `json:"playCount"`
	ArtistRole      string     `json:"artistRole"`
	ArtistPosition  int        `json:"artistPosition"`
	CreatedAt       time.Time  `json:"createdAt"`
}

type ArtistAlbumResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	CoverURL    *string    `json:"coverUrl,omitempty"`
	ReleaseDate *time.Time `json:"releaseDate,omitempty"`
	AlbumType   string     `json:"albumType"`
	ArtistRole  string     `json:"artistRole"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type ArtistOverviewResponse struct {
	Artist    *ArtistResponse       `json:"artist"`
	TopTracks []ArtistTrackResponse `json:"topTracks"`
	Albums    []ArtistAlbumResponse `json:"albums"`
	Singles   []ArtistAlbumResponse `json:"singles"`
	AppearsOn []ArtistTrackResponse `json:"appearsOn"`
}

type RelatedArtistResponse struct {
	ID               uuid.UUID  `json:"id"`
	Name             string     `json:"name"`
	Slug             string     `json:"slug"`
	Bio              *string    `json:"bio,omitempty"`
	ImageURL         *string    `json:"imageUrl,omitempty"`
	AvatarMediaID    *uuid.UUID `json:"avatarMediaId,omitempty"`
	BannerMediaID    *uuid.UUID `json:"bannerMediaId,omitempty"`
	IsVerified       bool       `json:"isVerified"`
	MonthlyListeners int64      `json:"monthlyListeners"`
	RelatedScore     float64    `json:"relatedScore"`
	RelatedSource    string     `json:"relatedSource"`
	RelatedAt        time.Time  `json:"relatedAt"`
}

// Mapper functions

func ArtistToResponse(a *Artist) *ArtistResponse {
	if a == nil {
		return nil
	}
	return &ArtistResponse{
		ID:               a.ID,
		Name:             a.Name,
		Slug:             a.Slug,
		Bio:              a.Bio,
		ImageURL:         a.ImageURL,
		AvatarMediaID:    a.AvatarMediaID,
		BannerMediaID:    a.BannerMediaID,
		Country:          a.Country,
		IsVerified:       a.IsVerified,
		MonthlyListeners: a.MonthlyListeners,
		CreatedAt:        a.CreatedAt,
		UpdatedAt:        a.UpdatedAt,
	}
}

func ArtistListToResponse(items []Artist) []ArtistResponse {
	if items == nil {
		return nil
	}
	res := make([]ArtistResponse, len(items))
	for i := range items {
		res[i] = *ArtistToResponse(&items[i])
	}
	return res
}
