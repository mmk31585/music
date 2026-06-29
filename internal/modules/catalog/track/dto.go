package track

import (
	"time"

	"github.com/google/uuid"
)

// ---------- Response DTOs ----------

type TrackResponse struct {
	ID uuid.UUID `json:"id"`

	// Deprecated compatibility field.
	ArtistID uuid.UUID `json:"artist_id"`

	AlbumID         *uuid.UUID `json:"album_id"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	DurationSeconds int        `json:"duration_seconds"`
	TrackNumber     *int       `json:"track_number"`
	Explicit        bool       `json:"explicit"`

	AudioURL *string `json:"audio_url"`
	CoverURL *string `json:"cover_url"`

	AudioMediaID *uuid.UUID `json:"audio_media_id"`
	CoverMediaID *uuid.UUID `json:"cover_media_id"`

	ArtistName *string `json:"artist_name"`

	PlayCount int64 `json:"play_count"`
	IsPublic  bool  `json:"is_public"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	Artists []TrackArtistResponse `json:"artists,omitempty"`
	Genres  []GenreResponse       `json:"genres,omitempty"`
}

type TrackArtistResponse struct {
	ArtistID uuid.UUID `json:"artist_id"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Role     string    `json:"role"`
	Position int       `json:"position"`
}

type GenreResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type TrackCreditResponse struct {
	ID         uuid.UUID `json:"id"`
	TrackID    uuid.UUID `json:"track_id"`
	ArtistID   uuid.UUID `json:"artist_id"`
	ArtistName string    `json:"artist_name"`
	ArtistSlug string    `json:"artist_slug"`
	CreditType string    `json:"credit_type"`
	Position   int       `json:"position"`
	CreatedAt  time.Time `json:"created_at"`
}

// Mapper functions

func TrackToResponse(t *Track) *TrackResponse {
	if t == nil {
		return nil
	}
	artists := make([]TrackArtistResponse, len(t.Artists))
	for i, a := range t.Artists {
		artists[i] = TrackArtistResponse{
			ArtistID: a.ArtistID,
			Name:     a.Name,
			Slug:     a.Slug,
			Role:     a.Role,
			Position: a.Position,
		}
	}
	genres := make([]GenreResponse, len(t.Genres))
	for i, g := range t.Genres {
		genres[i] = GenreResponse{
			ID:        g.ID,
			Name:      g.Name,
			Slug:      g.Slug,
			CreatedAt: g.CreatedAt,
		}
	}
	var artistName *string
	if len(t.Artists) > 0 {
		name := t.Artists[0].Name
		artistName = &name
	}

	return &TrackResponse{
		ID:              t.ID,
		ArtistID:        t.ArtistID,
		AlbumID:         t.AlbumID,
		Title:           t.Title,
		Slug:            t.Slug,
		DurationSeconds: t.DurationSeconds,
		TrackNumber:     t.TrackNumber,
		Explicit:        t.Explicit,
		AudioURL:        t.AudioURL,
		CoverURL:        t.CoverURL,
		AudioMediaID:    t.AudioMediaID,
		CoverMediaID:    t.CoverMediaID,
		ArtistName:      artistName,
		PlayCount:       t.PlayCount,
		IsPublic:        t.IsPublic,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
		Artists:         artists,
		Genres:          genres,
	}
}

func TrackListToResponse(items []Track) []TrackResponse {
	res := make([]TrackResponse, len(items))
	for i := range items {
		res[i] = *TrackToResponse(&items[i])
	}
	return res
}

func TrackCreditsToResponse(items []TrackCredit) []TrackCreditResponse {
	if items == nil {
		return nil
	}
	res := make([]TrackCreditResponse, len(items))
	for i, c := range items {
		res[i] = TrackCreditResponse{
			ID:         c.ID,
			TrackID:    c.TrackID,
			ArtistID:   c.ArtistID,
			ArtistName: c.ArtistName,
			ArtistSlug: c.ArtistSlug,
			CreditType: c.CreditType,
			Position:   c.Position,
			CreatedAt:  c.CreatedAt,
		}
	}
	return res
}
