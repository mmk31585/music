package playlist

import (
	"time"

	"github.com/google/uuid"
)

type CreatePlaylistRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description *string `json:"description"`
	CoverURL    *string `json:"cover_url"`
	IsPublic    bool    `json:"is_public"`
}

type UpdatePlaylistRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description *string `json:"description"`
	CoverURL    *string `json:"cover_url"`
	IsPublic    bool    `json:"is_public"`
}

type AddTrackRequest struct {
	TrackID uuid.UUID `json:"track_id" binding:"required"`
}

type ReorderTrackRequest struct {
	TrackID     uuid.UUID `json:"track_id" binding:"required"`
	NewPosition int       `json:"new_position" binding:"required"`
}

type PlaylistResponse struct {
	ID              uuid.UUID           `json:"id"`
	UserID          uuid.UUID           `json:"user_id"`
	Name            string              `json:"name"`
	Description     *string             `json:"description,omitempty"`
	CoverURL        *string             `json:"cover_url,omitempty"`
	IsPublic        bool                `json:"is_public"`
	IsCollaborative bool                `json:"is_collaborative"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	Tracks          []PlaylistTrackItem `json:"tracks,omitempty"`
}

type SetCollaborativeRequest struct {
	Collaborative bool `json:"collaborative"`
}

type PlaylistListItemResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	Name            string    `json:"name"`
	Description     *string   `json:"description,omitempty"`
	CoverURL        *string   `json:"cover_url,omitempty"`
	IsPublic        bool      `json:"is_public"`
	IsCollaborative bool      `json:"is_collaborative"`
	TrackCount      int       `json:"track_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func ToPlaylistResponse(p Playlist, tracks []PlaylistTrackItem) PlaylistResponse {
	return PlaylistResponse{
		ID:              p.ID,
		UserID:          p.UserID,
		Name:            p.Name,
		Description:     p.Description,
		CoverURL:        p.CoverURL,
		IsPublic:        p.IsPublic,
		IsCollaborative: p.IsCollaborative,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
		Tracks:          tracks,
	}
}
