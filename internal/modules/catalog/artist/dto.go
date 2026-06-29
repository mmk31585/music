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
	ImageURL *string   `json:"image_url,omitempty"`

	AvatarMediaID *uuid.UUID `json:"avatar_media_id,omitempty"`
	BannerMediaID *uuid.UUID `json:"banner_media_id,omitempty"`
	Country       *string    `json:"country,omitempty"`

	IsVerified       bool       `json:"is_verified"`
	MonthlyListeners int64      `json:"monthly_listeners"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

type ArtistTrackResponse struct {
	ID              uuid.UUID  `json:"id"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	AlbumID         *uuid.UUID `json:"album_id,omitempty"`
	DurationSeconds int        `json:"duration_seconds"`
	TrackNumber     *int       `json:"track_number,omitempty"`
	Explicit        bool       `json:"explicit"`
	AudioURL        *string    `json:"audio_url,omitempty"`
	CoverURL        *string    `json:"cover_url,omitempty"`
	PlayCount       int64      `json:"play_count"`
	ArtistRole      string     `json:"artist_role"`
	ArtistPosition  int        `json:"artist_position"`
	CreatedAt       time.Time  `json:"created_at"`
}

type ArtistAlbumResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	CoverURL    *string    `json:"cover_url,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	AlbumType   string     `json:"album_type"`
	ArtistRole  string     `json:"artist_role"`
	CreatedAt   time.Time  `json:"created_at"`
}

type ArtistOverviewResponse struct {
	Artist    *ArtistResponse       `json:"artist"`
	TopTracks []ArtistTrackResponse `json:"top_tracks"`
	Albums    []ArtistAlbumResponse `json:"albums"`
	Singles   []ArtistAlbumResponse `json:"singles"`
	AppearsOn []ArtistTrackResponse `json:"appears_on"`
}

type RelatedArtistResponse struct {
	ID               uuid.UUID  `json:"id"`
	Name             string     `json:"name"`
	Slug             string     `json:"slug"`
	Bio              *string    `json:"bio,omitempty"`
	ImageURL         *string    `json:"image_url,omitempty"`
	AvatarMediaID    *uuid.UUID `json:"avatar_media_id,omitempty"`
	BannerMediaID    *uuid.UUID `json:"banner_media_id,omitempty"`
	IsVerified       bool       `json:"is_verified"`
	MonthlyListeners int64      `json:"monthly_listeners"`
	RelatedScore     float64    `json:"related_score"`
	RelatedSource    string     `json:"related_source"`
	RelatedAt        time.Time  `json:"related_at"`
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
	res := make([]ArtistResponse, len(items))
	for i := range items {
		res[i] = *ArtistToResponse(&items[i])
	}
	return res
}
