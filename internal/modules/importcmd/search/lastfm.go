package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type LastFMProvider struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewLastFMProvider(apiKey string) *LastFMProvider {
	return &LastFMProvider{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: "https://ws.audioscrobbler.com/2.0",
		apiKey:  apiKey,
	}
}

func (p *LastFMProvider) Name() string { return "lastfm" }

type lfmTrack struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Duration   string `json:"duration"`
	Streamable string `json:"streamable"`
	Artist     string `json:"artist"`
	Image      []struct {
		Text string `json:"#text"`
		Size string `json:"size"`
	} `json:"image"`
	MBID string `json:"mbid"`
}

type lfmResults struct {
	TrackMatches []lfmTrack `json:"trackmatches"`
}

type lfmResponse struct {
	Results struct {
		TrackMatches struct {
			Track []lfmTrack `json:"track"`
		} `json:"trackmatches"`
	} `json:"results"`
}

func (p *LastFMProvider) Search(ctx context.Context, q SearchQuery) ([]Result, error) {
	query := q.Raw
	if query == "" {
		query = q.Artist + " " + q.Title
	}

	u, _ := url.Parse(p.baseURL)
	u.RawQuery = url.Values{
		"method":  {"track.search"},
		"track":   {query},
		"api_key": {p.apiKey},
		"format":  {"json"},
		"limit":   {"10"},
	}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lastfm: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("lastfm: status %d", resp.StatusCode)
	}

	var lr lfmResponse
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		return nil, fmt.Errorf("lastfm decode: %w", err)
	}

	results := make([]Result, 0, len(lr.Results.TrackMatches.Track))
	for _, t := range lr.Results.TrackMatches.Track {
		duration := 0
		if t.Duration != "" {
			fmt.Sscanf(t.Duration, "%d", &duration)
		}

		thumbnail := ""
		for _, img := range t.Image {
			if img.Size == "extralarge" || img.Size == "large" {
				if img.Text != "" {
					thumbnail = img.Text
				}
			}
		}

		results = append(results, Result{
			Title:     t.Name,
			Artist:    t.Artist,
			URL:       t.URL,
			Duration:  duration,
			Thumbnail: thumbnail,
			Source:    "lastfm",
			ExternalIDs: ExternalIDs{
				MBID: t.MBID,
			},
		})
	}

	return results, nil
}
