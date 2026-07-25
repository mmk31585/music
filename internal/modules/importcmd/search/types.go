package search

import "context"

type ExternalIDs struct {
	SpotifyID string `json:"spotify_id,omitempty"`
	DeezerID  string `json:"deezer_id,omitempty"`
	MBID      string `json:"mbid,omitempty"`
	ISRC      string `json:"isrc,omitempty"`
}

type SearchQuery struct {
	Raw    string `json:"raw"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

type Result struct {
	Title       string      `json:"title"`
	Artist      string      `json:"artist"`
	Album       string      `json:"album"`
	URL         string      `json:"url"`
	Duration    int         `json:"duration"`
	Thumbnail   string      `json:"thumbnail"`
	Source      string      `json:"source"`
	Score       float64     `json:"score"`
	ISRC        string      `json:"isrc,omitempty"`
	ExternalIDs ExternalIDs `json:"external_ids,omitempty"`
}

type Provider interface {
	Name() string
	Search(ctx context.Context, q SearchQuery) ([]Result, error)
}
