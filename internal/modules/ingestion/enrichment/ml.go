package enrichment

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// MLEnrichmentClient queries the ML microservice for enrichment data
// (lyrics and covers) during the enrichment pipeline.
type MLEnrichmentClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewMLEnrichmentClient(baseURL string) *MLEnrichmentClient {
	return &MLEnrichmentClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// lyricsJobStatus mirrors the ML service's job detail response for lyrics.
type lyricsJobStatus struct {
	JobID            string   `json:"job_id"`
	TrackID          string   `json:"track_id"`
	Status           string   `json:"status"`
	ResultLRCContent string   `json:"result_lrc_content,omitempty"`
	ResultPlainText  string   `json:"result_plain_text,omitempty"`
	ResultConfidence *float64 `json:"result_confidence,omitempty"`
	ErrorMessage     string   `json:"error_message,omitempty"`
}

// coverJobStatus mirrors the ML service's job detail response for covers.
type coverJobStatus struct {
	JobID        string `json:"job_id"`
	AlbumID      string `json:"album_id"`
	Status       string `json:"status"`
	OptimizedURL string `json:"optimized_url,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// FetchLyrics queries the ML service for existing lyrics for a track.
// Returns nil if no lyrics are available or the service is unreachable.
func (c *MLEnrichmentClient) FetchLyrics(ctx context.Context, _ TrackQuery) (*MLResult, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/v1/lyrics/jobs", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var jobs []lyricsJobStatus
	if err := json.Unmarshal(body, &jobs); err != nil {
		return nil, nil
	}

	// Return the first completed job with lyrics
	for _, j := range jobs {
		if j.Status == "completed" && (j.ResultPlainText != "" || j.ResultLRCContent != "") {
			lyrics := j.ResultPlainText
			lyricsType := "plain"
			if j.ResultLRCContent != "" {
				lyrics = j.ResultLRCContent
				lyricsType = "lrc"
			}
			confidence := 0.0
			if j.ResultConfidence != nil {
				confidence = *j.ResultConfidence
			}
			return &MLResult{
				Lyrics:     lyrics,
				LyricsType: lyricsType,
				Confidence: confidence,
			}, nil
		}
	}

	return nil, nil
}

// FetchCover queries the ML service for existing optimized cover for a track/album.
// Returns nil if no cover is available or the service is unreachable.
func (c *MLEnrichmentClient) FetchCover(ctx context.Context, _ TrackQuery) (*MLResult, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/v1/covers/jobs", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var jobs []coverJobStatus
	if err := json.Unmarshal(body, &jobs); err != nil {
		return nil, nil
	}

	// Return the first completed job with an optimized URL
	for _, j := range jobs {
		if j.Status == "completed" && j.OptimizedURL != "" {
			return &MLResult{
				AlbumCoverURL: j.OptimizedURL,
			}, nil
		}
	}

	return nil, nil
}
