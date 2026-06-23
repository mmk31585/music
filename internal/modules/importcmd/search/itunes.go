package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// iTunesProvider searches Apple's iTunes Store API for tracks.
// It returns preview URLs (30-second clips) that are directly downloadable
// without requiring yt-dlp or access to blocked services like YouTube.
type iTunesProvider struct {
	client  *http.Client
	baseURL string
}

func NewiTunesProvider() *iTunesProvider {
	return &iTunesProvider{
		client:  &http.Client{Timeout: 8 * time.Second},
		baseURL: "https://itunes.apple.com",
	}
}

func (p *iTunesProvider) Name() string { return "itunes" }

type itunesTrack struct {
	TrackName        string `json:"trackName"`
	ArtistName       string `json:"artistName"`
	CollectionName   string `json:"collectionName"`
	TrackID          int64  `json:"trackId"`
	PreviewURL       string `json:"previewUrl"`
	ArtworkURL       string `json:"artworkUrl100"`
	TrackTimeMillis  int    `json:"trackTimeMillis"`
	PrimaryGenreName string `json:"primaryGenreName"`
	TrackNumber      int    `json:"trackNumber"`
	ReleaseDate      string `json:"releaseDate"`
	TrackViewURL     string `json:"trackViewUrl"`
}

type itunesResponse struct {
	ResultCount int           `json:"resultCount"`
	Results     []itunesTrack `json:"results"`
}

func (p *iTunesProvider) Search(ctx context.Context, q SearchQuery) ([]Result, error) {
	query := q.Raw
	if query == "" {
		query = q.Artist + " " + q.Title
	}

	u, _ := url.Parse(p.baseURL + "/search")
	u.RawQuery = url.Values{
		"term":   {query},
		"limit":  {"10"},
		"entity": {"song"},
	}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("itunes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("itunes: status %d", resp.StatusCode)
	}

	var ir itunesResponse
	if err := json.NewDecoder(resp.Body).Decode(&ir); err != nil {
		return nil, fmt.Errorf("itunes decode: %w", err)
	}

	results := make([]Result, 0, ir.ResultCount)
	for _, t := range ir.Results {
		// iTunes returns results from any matching term; optionally filter for close matches
		duration := 0
		if t.TrackTimeMillis > 0 {
			duration = t.TrackTimeMillis / 1000
		}
		results = append(results, Result{
			Title:     t.TrackName,
			Artist:    t.ArtistName,
			Album:     t.CollectionName,
			URL:       t.PreviewURL,
			Duration:  duration,
			Thumbnail: t.ArtworkURL,
			Source:    "itunes",
			ExternalIDs: ExternalIDs{
				ISRC: fmt.Sprintf("itunes:%d", t.TrackID),
			},
		})
	}

	return results, nil
}
