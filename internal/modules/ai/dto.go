package ai

import "encoding/json"

type EmbeddingRequest struct {
	TrackIDs []string `json:"track_ids" binding:"required"`
}

type MoodAnalysisRequest struct {
	TrackID string `json:"track_id" binding:"required"`
}

type GeneratePlaylistRequest struct {
	Prompt      string   `json:"prompt"`
	Mood        string   `json:"mood,omitempty"`
	Activity    string   `json:"activity,omitempty"`
	SeedTrackID string   `json:"seed_track_id,omitempty"`
	Genre       string   `json:"genre,omitempty"`
	Limit       int      `json:"limit,omitempty"`
	ExcludeIDs  []string `json:"exclude_ids,omitempty"`
}

type EmbeddingResponse struct {
	TrackID      string    `json:"track_id"`
	ModelVersion string    `json:"model_version"`
	Dimensions   int       `json:"dimensions"`
	Vector       []float64 `json:"vector,omitempty"`
}

type MoodResponse struct {
	TrackID          string          `json:"track_id"`
	MoodTags         json.RawMessage `json:"mood_tags"`
	Energy           float64         `json:"energy"`
	Valence          float64         `json:"valence"`
	Tempo            float64         `json:"tempo"`
	Danceability     float64         `json:"danceability"`
	Acousticness     float64         `json:"acousticness"`
	Instrumentalness float64         `json:"instrumentalness"`
	Liveness         float64         `json:"liveness"`
	Speechiness      float64         `json:"speechiness"`
}

type SimilarTracksRequest struct {
	TrackID string `json:"track_id" binding:"required"`
	Limit   int    `json:"limit,omitempty"`
	Mood    string `json:"mood,omitempty"`
}

type AIPlaylistResponse struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Tracks      []TrackItem `json:"tracks"`
	GeneratedAt string      `json:"generated_at"`
}

type TrackItem struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Artist   string  `json:"artist"`
	CoverURL string  `json:"cover_url,omitempty"`
	Duration int     `json:"duration,omitempty"`
	Energy   float64 `json:"energy,omitempty"`
	Valence  float64 `json:"valence,omitempty"`
}
