package enrichment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type LRCLibResult struct {
	ID           int    `json:"id"`
	TrackName    string `json:"trackName"`
	ArtistName   string `json:"artistName"`
	AlbumName    string `json:"albumName"`
	Duration     int    `json:"duration"`
	Synced       bool   `json:"synced"`
	PlainLyrics  string `json:"plainLyrics"`
	SyncedLyrics string `json:"syncedLyrics"`
}

type LRCLibClient interface {
	SearchLyrics(ctx context.Context, query TrackQuery) (*LRCLibResult, error)
}

type lrclibClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewLRCLibClient() LRCLibClient {
	return &lrclibClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://lrclib.net/api",
	}
}

func (c *lrclibClient) callGet(ctx context.Context, params url.Values) (*LRCLibResult, error) {
	u := fmt.Sprintf("%s/get?%s", c.baseURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "MojaMusic/1.0 (music-ingestion)")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("lrclib returned %d", resp.StatusCode)
	}

	var result LRCLibResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

func (c *lrclibClient) callSearch(ctx context.Context, q string) ([]LRCLibResult, error) {
	u := fmt.Sprintf("%s/search?q=%s", c.baseURL, url.QueryEscape(q))
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "MojaMusic/1.0 (music-ingestion)")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	var results []LRCLibResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}
	return results, nil
}

func (c *lrclibClient) SearchLyrics(ctx context.Context, query TrackQuery) (*LRCLibResult, error) {
	params := url.Values{}
	if query.Artist != "" {
		params.Set("artist_name", query.Artist)
	}
	if query.Title != "" {
		params.Set("track_name", query.Title)
	}
	if query.Album != "" {
		params.Set("album_name", query.Album)
	}
	if query.Duration > 0 {
		params.Set("duration", fmt.Sprintf("%d", query.Duration))
	}
	if len(params) == 0 {
		return nil, nil
	}

	// Step 1: Try exact match via /get
	result, err := c.callGet(ctx, params)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}

	// Step 2: Fallback to /search with just track + artist (fuzzy)
	if query.Title == "" && query.Artist == "" {
		return nil, nil
	}
	searchQ := query.Title
	if query.Artist != "" {
		searchQ = query.Artist + " " + searchQ
	}

	results, err := c.callSearch(ctx, searchQ)
	if err != nil {
		return nil, nil // silent fallback — search is best-effort
	}

	// Prefer a result with synced lyrics, then any with lyrics
	var best *LRCLibResult
	for i := range results {
		r := &results[i]
		if r.Synced && r.SyncedLyrics != "" {
			return r, nil // synced is best, return immediately
		}
		if best == nil && (r.PlainLyrics != "" || r.SyncedLyrics != "") {
			best = r
		}
	}
	return best, nil
}
