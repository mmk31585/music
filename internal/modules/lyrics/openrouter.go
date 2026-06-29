package lyrics

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"music/internal/config"
)

// OpenRouterClient communicates with OpenRouter's OpenAI-compatible API
// to fetch, sync, and review lyrics. Sits between LRCLIB and Whisper in
// the pipeline: LRCLIB → OpenRouter → Whisper.
type OpenRouterClient struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

// NewOpenRouterClient creates a new client from config. Returns nil if no
// API key is configured (OpenRouter is optional).
func NewOpenRouterClient(cfg config.OpenRouterConfig) *OpenRouterClient {
	if cfg.APIKey == "" {
		return nil
	}
	return &OpenRouterClient{
		apiKey:  cfg.APIKey,
		model:   cfg.Model,
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
		httpClient: &http.Client{
			Timeout: 220 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout:   60 * time.Second,
				ResponseHeaderTimeout: 60 * time.Second,
				ExpectContinueTimeout: 10 * time.Second,
				MaxIdleConns:          2,
				IdleConnTimeout:       90 * time.Second,
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
				},
			},
		},
	}
}

// openRouterMessage represents a chat message in the OpenAI format.
type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openRouterRequest is the request body for the OpenAI-compatible API.
type openRouterRequest struct {
	Model       string              `json:"model"`
	Messages    []openRouterMessage `json:"messages"`
	Temperature float64             `json:"temperature"`
	MaxTokens   int                 `json:"max_tokens"`
}

// openRouterResponse is the response from the OpenAI-compatible API.
type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// chat calls the OpenRouter chat completions endpoint with the given messages.
func (c *OpenRouterClient) chat(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	body := openRouterRequest{
		Model: c.model,
		Messages: []openRouterMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.3,
		MaxTokens:   4000,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("openrouter marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("openrouter request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("HTTP-Referer", "https://moja-music.app")
	req.Header.Set("X-Title", "Moja Music")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("openrouter do: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("openrouter read: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openrouter status %d: %s", resp.StatusCode, string(respBody))
	}

	var result openRouterResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("openrouter parse: %w", err)
	}

	if result.Error != nil && result.Error.Message != "" {
		return "", fmt.Errorf("openrouter error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("openrouter: no choices returned")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

// parseResponse splits the AI response at the first newline to extract
// LANGUAGE: <code> prefix. Returns (language, content).
func parseResponse(raw string) (language string, content string) {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "LANGUAGE:") {
		parts := strings.SplitN(raw, "\n", 2)
		langPart := strings.TrimSpace(strings.TrimPrefix(parts[0], "LANGUAGE:"))
		langPart = strings.TrimSpace(langPart)
		if len(parts) > 1 {
			content = strings.TrimSpace(parts[1])
		}
		return langPart, content
	}
	// No LANGUAGE prefix — return as-is with empty language
	return "", raw
}

// FetchLyrics generates lyrics from scratch given track title and artist.
// Returns LRC-formatted content and detected language.
func (c *OpenRouterClient) FetchLyrics(ctx context.Context, trackTitle, artistName string) (lrcContent string, language string, err error) {
	systemPrompt := `You are a professional lyricist who writes song lyrics in the style and language of the given song.
Respond ONLY in this format:
LANGUAGE: <ISO 639-1 code>
<LRC-lyrics with timestamps in [mm:ss.xx] format>

The lyrics should be:
- In the same language as the song title/artist suggests
- Each line timestamped approximately evenly across a 3-4 minute track
- 2-4 lines per verse/chorus section
- Include verse, chorus, and outro sections if appropriate`

	userPrompt := fmt.Sprintf(`Write lyrics for the song titled "%s" by artist "%s".
Return LRC-formatted lyrics with timestamps.`, trackTitle, artistName)

	raw, err := c.chat(ctx, systemPrompt, userPrompt)
	if err != nil {
		return "", "", err
	}

	lang, content := parseResponse(raw)
	return content, lang, nil
}

// SyncLyrics takes plain text lyrics and adds LRC timestamps.
// Returns LRC-formatted content.
func (c *OpenRouterClient) SyncLyrics(ctx context.Context, trackTitle, artistName string, plainLyrics string, durationSeconds int) (lrcContent string, language string, err error) {
	systemPrompt := `You are a professional lyrics synchronizer. Given plain song lyrics, add [mm:ss.xx] timestamps to create LRC format.
Respond ONLY in this format:
LANGUAGE: <ISO 639-1 code>
<LRC-lyrics with timestamps>

Rules:
- Distribute lines approximately evenly across the track duration
- Each line gets a timestamp like [01:23.45]
- If the plain text has section headers (like [Verse], [Chorus]), remove them or convert to LRC comments
- Keep the original words intact — only add timestamps, do not rewrite
- Use the same language as the original lyrics`

	userPrompt := fmt.Sprintf(`Track: "%s" by %s
Duration: %d seconds
Plain lyrics:
%s

Add [mm:ss.xx] timestamps to create LRC format.`, trackTitle, artistName, durationSeconds, plainLyrics)

	raw, err := c.chat(ctx, systemPrompt, userPrompt)
	if err != nil {
		return "", "", err
	}

	lang, content := parseResponse(raw)
	return content, lang, nil
}

// ReviewLyrics reviews and fixes existing LRC lyrics (spelling, timing).
// Returns corrected LRC content.
func (c *OpenRouterClient) ReviewLyrics(ctx context.Context, trackTitle, artistName string, existingLRC string) (lrcContent string, language string, err error) {
	systemPrompt := `You are a professional lyrics editor. Review the given LRC lyrics and fix any issues:
- Correct spelling mistakes
- Adjust timestamps that seem out of order or incorrect
- Ensure consistent formatting
- Do NOT rewrite the lyrics creatively — only fix errors

Respond ONLY in this format:
LANGUAGE: <ISO 639-1 code>
<corrected LRC-lyrics>`

	userPrompt := fmt.Sprintf(`Track: "%s" by %s
Existing LRC lyrics:
%s

Review and fix any issues.`, trackTitle, artistName, existingLRC)

	raw, err := c.chat(ctx, systemPrompt, userPrompt)
	if err != nil {
		return "", "", err
	}

	lang, content := parseResponse(raw)
	return content, lang, nil
}
