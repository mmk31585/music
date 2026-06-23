package acquisition

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// iTunesResolver searches the Apple iTunes Search API for a downloadable
// preview URL. This resolver does not use yt-dlp — it makes direct HTTP
// requests to the iTunes API. iTunes preview URLs are always available
// from any network that can reach Apple's CDN.
type iTunesResolver struct{}

func newiTunesResolver() *iTunesResolver {
	return &iTunesResolver{}
}

func (r *iTunesResolver) Name() string { return "itunes" }

// Priority is lowest — tried only after SoundCloud/YouTube fail.
func (r *iTunesResolver) Priority() int { return 100 }

func (r *iTunesResolver) Resolve(ctx context.Context, q ResolveQuery) (*Candidate, error) {
	// Build iTunes search URL
	searchTerm := url.QueryEscape(q.Artist + " " + q.Title)
	apiURL := fmt.Sprintf(
		"https://itunes.apple.com/search?term=%s&media=music&limit=5&entity=song",
		searchTerm,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("itunes request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("itunes api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("itunes api returned status %d", resp.StatusCode)
	}

	var result struct {
		ResultCount int `json:"resultCount"`
		Results     []struct {
			TrackName       string `json:"trackName"`
			ArtistName      string `json:"artistName"`
			PreviewURL      string `json:"previewUrl"`
			TrackTimeMillis int    `json:"trackTimeMillis"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("itunes decode: %w", err)
	}

	if result.ResultCount == 0 {
		return nil, fmt.Errorf("itunes: no results for %s - %s", q.Artist, q.Title)
	}

	// Find the best match — prefer exact artist match
	for _, track := range result.Results {
		if track.PreviewURL == "" {
			continue
		}
		// Accept if artist name roughly matches
		if containsIgnoreCase(track.ArtistName, q.Artist) ||
			containsIgnoreCase(q.Artist, track.ArtistName) {
			return &Candidate{
				URL:         track.PreviewURL,
				Source:      "itunes",
				Quality:     Quality128,
				Confidence:  0.8,
				AudioFormat: "m4a",
			}, nil
		}
	}

	// Fallback: return the first result even if artist doesn't match
	if len(result.Results) > 0 && result.Results[0].PreviewURL != "" {
		return &Candidate{
			URL:         result.Results[0].PreviewURL,
			Source:      "itunes",
			Quality:     Quality128,
			Confidence:  0.5,
			AudioFormat: "m4a",
		}, nil
	}

	return nil, fmt.Errorf("itunes: no preview URL found for %s - %s", q.Artist, q.Title)
}

func containsIgnoreCase(s, substr string) bool {
	sMax := len(s)
	subMax := len(substr)
	if subMax > sMax {
		return false
	}
	for i := 0; i <= sMax-subMax; i++ {
		if caseInsensitiveEqual(s[i:i+subMax], substr) {
			return true
		}
	}
	return false
}

func caseInsensitiveEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i]|0x20 != b[i]|0x20 {
			return false
		}
	}
	return true
}
