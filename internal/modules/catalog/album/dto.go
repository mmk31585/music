package album

import (
	"time"

	"github.com/google/uuid"
)

// ---------- Response DTOs ----------

type AlbumResponse struct {
	ID uuid.UUID `json:"id"`

	// Deprecated compatibility field.
	ArtistID uuid.UUID `json:"artist_id"`

	ArtistName   *string    `json:"artist_name"`
	Title        string     `json:"title"`
	Slug         string     `json:"slug"`
	CoverURL     *string    `json:"cover_url"`
	CoverMediaID *uuid.UUID `json:"cover_media_id"`

	ReleaseDate *time.Time `json:"release_date"`
	AlbumType   string     `json:"album_type"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	TrackCount int                   `json:"track_count"`
	Artists    []AlbumArtistResponse `json:"artists,omitempty"`
}

type AlbumArtistResponse struct {
	ArtistID uuid.UUID `json:"artist_id"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Role     string    `json:"role"`
	Position int       `json:"position"`
}

type AlbumTrackResponse struct {
	ID              uuid.UUID  `json:"id"`
	AlbumID         *uuid.UUID `json:"album_id,omitempty"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	DurationSeconds int        `json:"duration_seconds"`
	TrackNumber     *int       `json:"track_number,omitempty"`
	Explicit        bool       `json:"explicit"`
	AudioURL        *string    `json:"audio_url,omitempty"`
	CoverURL        *string    `json:"cover_url,omitempty"`
	PlayCount       int64      `json:"play_count"`
	IsPublic        bool       `json:"is_public"`
	CreatedAt       time.Time  `json:"created_at"`
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

	var artistName *string
	if len(a.Artists) > 0 {
		name := a.Artists[0].Name
		artistName = &name
	}

	return &AlbumResponse{
		ID:           a.ID,
		ArtistID:     a.ArtistID,
		ArtistName:   artistName,
		Title:        a.Title,
		Slug:         a.Slug,
		CoverURL:     a.CoverURL,
		CoverMediaID: a.CoverMediaID,
		ReleaseDate:  a.ReleaseDate,
		AlbumType:    a.AlbumType,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
		TrackCount:   0,
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
