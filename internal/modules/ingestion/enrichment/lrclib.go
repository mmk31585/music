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
	if len(params) == 0 {
		return nil, nil
	}

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
