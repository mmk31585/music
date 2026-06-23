package album

import (
	"time"

	"github.com/google/uuid"
)

// ---------- Response DTOs ----------

type AlbumResponse struct {
	ID uuid.UUID `json:"id"`

	// Deprecated compatibility field.
	ArtistID uuid.UUID `json:"artistId"`

	Title        string     `json:"title"`
	Slug         string     `json:"slug"`
	CoverURL     *string    `json:"coverUrl,omitempty"`
	CoverMediaID *uuid.UUID `json:"coverMediaId,omitempty"`

	ReleaseDate *time.Time `json:"releaseDate,omitempty"`
	AlbumType   string     `json:"albumType"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`

	Artists []AlbumArtistResponse `json:"artists,omitempty"`
}

type AlbumArtistResponse struct {
	ArtistID uuid.UUID `json:"artistId"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Role     string    `json:"role"`
	Position int       `json:"position"`
}

type AlbumTrackResponse struct {
	ID              uuid.UUID  `json:"id"`
	AlbumID         *uuid.UUID `json:"albumId,omitempty"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	DurationSeconds int        `json:"durationSeconds"`
	TrackNumber     *int       `json:"trackNumber,omitempty"`
	Explicit        bool       `json:"explicit"`
	AudioURL        *string    `json:"audioUrl,omitempty"`
	CoverURL        *string    `json:"coverUrl,omitempty"`
	PlayCount       int64      `json:"playCount"`
	IsPublic        bool       `json:"isPublic"`
	CreatedAt       time.Time  `json:"createdAt"`
}

// Mapper functions

func AlbumToResponse(a *Album) *AlbumResponse {
	if a == nil {
		return nil
	}
	artists := make([]AlbumArtistResponse, len(a.Artists))
	for i, ar := range a.Artists {
		artists[i] = AlbumArtistResponse{
			ArtistID: ar.ArtistID,
			Name:     ar.Name,
			Slug:     ar.Slug,
			Role:     ar.Role,
			Position: ar.Position,
		}
	}
	return &AlbumResponse{
		ID:           a.ID,
		ArtistID:     a.ArtistID,
		Title:        a.Title,
		Slug:         a.Slug,
		CoverURL:     a.CoverURL,
		CoverMediaID: a.CoverMediaID,
		ReleaseDate:  a.ReleaseDate,
		AlbumType:    a.AlbumType,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
		Artists:      artists,
	}
}

func AlbumListToResponse(items []Album) []AlbumResponse {
	res := make([]AlbumResponse, len(items))
	for i := range items {
		res[i] = *AlbumToResponse(&items[i])
	}
	return res
}

func AlbumTrackListToResponse(items []AlbumTrack) []AlbumTrackResponse {
	if items == nil {
		return nil
	}
	res := make([]AlbumTrackResponse, len(items))
	for i, t := range items {
		res[i] = AlbumTrackResponse{
			ID:              t.ID,
			AlbumID:         t.AlbumID,
			Title:           t.Title,
			Slug:            t.Slug,
			DurationSeconds: t.DurationSeconds,
			TrackNumber:     t.TrackNumber,
			Explicit:        t.Explicit,
			AudioURL:        t.AudioURL,
			CoverURL:        t.CoverURL,
			PlayCount:       t.PlayCount,
			IsPublic:        t.IsPublic,
			CreatedAt:       t.CreatedAt,
		}
	}
	return res
}
