package recommendation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"music/internal/config"
)

type SimilarTrackResult struct {
	TrackID         string   `json:"track_id"`
	SimilarityScore float64  `json:"similarity_score"`
	MatchedOn       []string `json:"matched_on"`
}

type similarTracksResponse struct {
	Tracks          []SimilarTrackResult `json:"tracks"`
	SeedTrackID     string               `json:"seed_track_id"`
	TotalCandidates int                  `json:"total_candidates"`
}

type SimilarityMLClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewSimilarityMLClient(cfg config.MLServiceConfig) *SimilarityMLClient {
	return &SimilarityMLClient{
		baseURL: cfg.BaseURL,
		httpClient: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
	}
}

func (c *SimilarityMLClient) GetSimilarTracks(
	ctx context.Context,
	trackID string,
	limit int,
	excludeTrackIDs []string,
) ([]SimilarTrackResult, error) {
	fullURL := fmt.Sprintf("%s/api/v1/similarity/tracks/%s?limit=%d", c.baseURL, trackID, limit)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result similarTracksResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if len(excludeTrackIDs) == 0 {
		return result.Tracks, nil
	}

	exclude := make(map[string]struct{}, len(excludeTrackIDs))
	for _, id := range excludeTrackIDs {
		exclude[id] = struct{}{}
	}

	filtered := make([]SimilarTrackResult, 0, len(result.Tracks))
	for _, t := range result.Tracks {
		if _, ok := exclude[t.TrackID]; !ok {
			filtered = append(filtered, t)
		}
	}
	return filtered, nil
}

func (c *SimilarityMLClient) GetSimilarToTrackSet(
	ctx context.Context,
	trackIDs []string,
	limit int,
	excludeTrackIDs []string,
) ([]SimilarTrackResult, error) {
	fullURL := fmt.Sprintf("%s/api/v1/similarity/from-set?limit=%d", c.baseURL, limit)

	body := map[string]interface{}{
		"track_ids": trackIDs,
	}
	if len(excludeTrackIDs) > 0 {
		body["exclude_ids"] = excludeTrackIDs
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result similarTracksResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if len(excludeTrackIDs) == 0 {
		return result.Tracks, nil
	}

	exclude := make(map[string]struct{}, len(excludeTrackIDs))
	for _, id := range excludeTrackIDs {
		exclude[id] = struct{}{}
	}

	filtered := make([]SimilarTrackResult, 0, len(result.Tracks))
	for _, t := range result.Tracks {
		if _, ok := exclude[t.TrackID]; !ok {
			filtered = append(filtered, t)
		}
	}
	return filtered, nil
}

// GetBulkTempo returns tempo_bpm for a batch of track IDs from the
// ml-service's track_embeddings table.  Tracks without embeddings are
// returned with value nil.  Called during taste profile recomputation
// to compute avg_tempo_preference.
func (c *SimilarityMLClient) GetBulkTempo(ctx context.Context, trackIDs []string) (map[string]*float64, error) {
	if len(trackIDs) == 0 {
		return map[string]*float64{}, nil
	}

	joined := strings.Join(trackIDs, ",")
	fullURL := fmt.Sprintf("%s/api/v1/similarity/tempo/bulk?track_ids=%s", c.baseURL, url.QueryEscape(joined))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Tempos map[string]*float64 `json:"tempos"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return result.Tempos, nil
}
