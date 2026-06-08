package playlist

import "time"

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
	TrackID int64 `json:"track_id" binding:"required"`
}

type ReorderTrackRequest struct {
	TrackID     int64 `json:"track_id" binding:"required"`
	NewPosition int   `json:"new_position" binding:"required"`
}

type PlaylistResponse struct {
	ID          int64               `json:"id"`
	UserID      int64               `json:"user_id"`
	Name        string              `json:"name"`
	Description *string             `json:"description,omitempty"`
	CoverURL    *string             `json:"cover_url,omitempty"`
	IsPublic    bool                `json:"is_public"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Tracks      []PlaylistTrackItem `json:"tracks,omitempty"`
}

type PlaylistListItemResponse struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CoverURL    *string   `json:"cover_url,omitempty"`
	IsPublic    bool      `json:"is_public"`
	TrackCount  int       `json:"track_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToPlaylistResponse(p Playlist, tracks []PlaylistTrackItem) PlaylistResponse {
	return PlaylistResponse{
		ID:          p.ID,
		UserID:      p.UserID,
		Name:        p.Name,
		Description: p.Description,
		CoverURL:    p.CoverURL,
		IsPublic:    p.IsPublic,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Tracks:      tracks,
	}
}
