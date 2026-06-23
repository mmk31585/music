package lyrics

import "time"

type Lyrics struct {
	ID              string    `json:"id" db:"id"`
	TrackID         string    `json:"track_id" db:"track_id"`
	Language        string    `json:"language" db:"language"`
	Type            string    `json:"type" db:"type"`
	Content         string    `json:"content" db:"content"`
	Source          string    `json:"source" db:"source"`
	ConfidenceScore *float64  `json:"confidence_score,omitempty" db:"confidence_score"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type LyricsLine struct {
	TimeSeconds float64 `json:"time_seconds"`
	Text        string  `json:"text"`
}

type CreateLyricsRequest struct {
	TrackID         string   `json:"track_id" binding:"required"`
	Language        string   `json:"language" binding:"required,min=2,max=10"`
	Type            string   `json:"type" binding:"required,oneof=plain lrc"`
	Content         string   `json:"content" binding:"required"`
	Source          string   `json:"source"`                     // "manual", "ai_generated", "contributed"
	ConfidenceScore *float64 `json:"confidence_score,omitempty"` // AI confidence 0-1
}

type UpdateLyricsRequest struct {
	Language *string `json:"language"`
	Type     *string `json:"type"`
	Content  *string `json:"content"`
}

type LyricsResponse struct {
	ID              string    `json:"id"`
	TrackID         string    `json:"track_id"`
	Language        string    `json:"language"`
	Type            string    `json:"type"`
	Content         string    `json:"content"`
	Source          string    `json:"source"`
	ConfidenceScore *float64  `json:"confidence_score,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type LyricsLineResponse struct {
	TimeSeconds float64 `json:"time_seconds"`
	Text        string  `json:"text"`
}

type LyricsByTrackResponse struct {
	ID        string               `json:"id"`
	TrackID   string               `json:"track_id"`
	Language  string               `json:"language"`
	Type      string               `json:"type"`
	Content   string               `json:"content"`
	Lines     []LyricsLineResponse `json:"lines,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

func ToLyricsResponse(l Lyrics) LyricsResponse {
	return LyricsResponse{
		ID:              l.ID,
		TrackID:         l.TrackID,
		Language:        l.Language,
		Type:            l.Type,
		Content:         l.Content,
		Source:          l.Source,
		ConfidenceScore: l.ConfidenceScore,
		CreatedAt:       l.CreatedAt,
		UpdatedAt:       l.UpdatedAt,
	}
}

func ToLyricsByTrackResponse(l Lyrics, lines []LyricsLine) LyricsByTrackResponse {
	lineResponses := make([]LyricsLineResponse, len(lines))
	for i, line := range lines {
		lineResponses[i] = LyricsLineResponse{
			TimeSeconds: line.TimeSeconds,
			Text:        line.Text,
		}
	}

	return LyricsByTrackResponse{
		ID:        l.ID,
		TrackID:   l.TrackID,
		Language:  l.Language,
		Type:      l.Type,
		Content:   l.Content,
		Lines:     lineResponses,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
	}
}
