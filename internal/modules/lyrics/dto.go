package lyrics

import (
	"time"
)

// These are aliases for backward compatibility and to avoid import conflicts
type (
	CreateLyricsRequestDTO   = CreateLyricsRequest
	UpdateLyricsRequestDTO   = UpdateLyricsRequest
	LyricsResponseDTO        = LyricsResponse
	LyricsByTrackResponseDTO = LyricsByTrackResponse
)

// Admin-specific DTOs
type AdminCreateLyricsRequest struct {
	TrackID  string `json:"track_id" binding:"required"`
	Language string `json:"language" binding:"required,min=2,max=10"`
	Type     string `json:"type" binding:"required,oneof=plain lrc"`
	Content  string `json:"content" binding:"required"`
}

type AdminUpdateLyricsRequest struct {
	Language *string `json:"language"`
	Type     *string `json:"type" binding:"omitempty,oneof=plain lrc"`
	Content  *string `json:"content"`
}

// Public DTOs
type PublicGetLyricsRequest struct {
	TrackID string `json:"track_id"`
	Lang    string `json:"lang,omitempty"`
}

type PublicLyricsResponse struct {
	ID        string               `json:"id"`
	TrackID   string               `json:"track_id"`
	Language  string               `json:"language"`
	Type      string               `json:"type"`
	Content   string               `json:"content"`
	Lines     []LyricsLineResponse `json:"lines,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

// Response wrapper for admin operations
type AdminLyricsResponse struct {
	ID        string    `json:"id"`
	TrackID   string    `json:"track_id"`
	Language  string    `json:"language"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SyncLyricsRequest is the payload for the OpenRouter sync endpoint.
type SyncLyricsRequest struct {
	TrackID    string `json:"track_id" binding:"required"`
	PlainText  string `json:"plain_text" binding:"required"`
	TrackTitle string `json:"track_title,omitempty"`
	ArtistName string `json:"artist_name,omitempty"`
}

// ReviewLyricsRequest is the payload for the OpenRouter review endpoint.
type ReviewLyricsRequest struct {
	TrackID     string `json:"track_id" binding:"required"`
	ExistingLRC string `json:"existing_lrc" binding:"required"`
	TrackTitle  string `json:"track_title,omitempty"`
	ArtistName  string `json:"artist_name,omitempty"`
}

// Pagination DTO
type LyricsPaginationRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

type PaginatedLyricsResponse struct {
	Data       []LyricsResponse `json:"data"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}
