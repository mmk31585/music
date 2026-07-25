package lyrics

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"music/internal/config"
)

// MLServiceClient communicates with the moja-ml-service (Python) to enqueue
// lyrics generation jobs.
type MLServiceClient struct {
	baseURL    string
	secret     string
	httpClient *http.Client
}

// NewMLServiceClient creates a new client from the ML service config.
func NewMLServiceClient(cfg config.MLServiceConfig) *MLServiceClient {
	return &MLServiceClient{
		baseURL: cfg.BaseURL,
		secret:  cfg.WebhookHMACSecret,
		httpClient: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
	}
}

// CreateLyricsJobRequest is the payload sent to moja-ml-service to enqueue a
// lyrics generation job. Exactly one of AudioFilePath or AudioDownloadURL
// must be set.
type CreateLyricsJobRequest struct {
	TrackID          string `json:"track_id"`
	AudioFilePath    string `json:"audio_file_path,omitempty"`
	AudioDownloadURL string `json:"audio_download_url,omitempty"`
	TrackTitle       string `json:"track_title,omitempty"`
	TrackArtist      string `json:"track_artist,omitempty"`
}

// CreateLyricsJobResponse is the response from moja-ml-service after
// enqueuing a job.
type CreateLyricsJobResponse struct {
	JobID string `json:"job_id"`
}

// LyricsJobStatus mirrors the ML service's LyricsJobDetailResponse.
type LyricsJobStatus struct {
	JobID                  string   `json:"job_id"`
	TrackID                string   `json:"track_id"`
	Status                 string   `json:"status"`
	AudioFilePath          string   `json:"audio_file_path,omitempty"`
	AudioDownloadURL       string   `json:"audio_download_url,omitempty"`
	TrackTitle             string   `json:"track_title,omitempty"`
	TrackArtist            string   `json:"track_artist,omitempty"`
	ResultLRCContent       string   `json:"result_lrc_content,omitempty"`
	ResultPlainText        string   `json:"result_plain_text,omitempty"`
	ResultConfidence       *float64 `json:"result_confidence,omitempty"`
	ResultDetectedLanguage string   `json:"result_detected_language,omitempty"`
	ErrorMessage           string   `json:"error_message,omitempty"`
	RetryCount             int      `json:"retry_count"`
	CallbackDelivered      bool     `json:"callback_delivered"`
	CreatedAt              string   `json:"created_at"`
	UpdatedAt              string   `json:"updated_at"`
}

// EnqueueLyricsJob sends a job creation request to the ML service.
// It returns the job ID on success.
func (c *MLServiceClient) EnqueueLyricsJob(ctx context.Context, req CreateLyricsJobRequest) (string, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	fullURL := c.baseURL + "/api/v1/lyrics/jobs"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("ML service returned %d: %s", resp.StatusCode, string(respBody))
	}

	var jobResp CreateLyricsJobResponse
	if err := json.Unmarshal(respBody, &jobResp); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if jobResp.JobID == "" {
		return "", fmt.Errorf("ML service returned empty job_id")
	}

	return jobResp.JobID, nil
}

// GetJobByTrack queries the ML service for the latest lyrics job for a track.
func (c *MLServiceClient) GetJobByTrack(ctx context.Context, trackID string) (*LyricsJobStatus, error) {
	fullURL := c.baseURL + "/api/v1/lyrics/jobs/by-track/" + trackID
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

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no job found for track %s", trackID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned %d: %s", resp.StatusCode, string(respBody))
	}

	var jobStatus LyricsJobStatus
	if err := json.Unmarshal(respBody, &jobStatus); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &jobStatus, nil
}

// SignCallbackPayload creates the HMAC-SHA256 signature for a webhook
// callback body, matching the Python signing scheme exactly.
//
//	signature = hex(HMAC-SHA256(secret, "{timestamp}.{raw_body}"))
func SignCallbackPayload(body []byte, secret string) (timestamp string, signature string) {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	payload := ts + "." + string(body)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return ts, hex.EncodeToString(mac.Sum(nil))
}
