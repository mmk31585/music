package contribution

import (
	"context"
	"encoding/json"
	"math"
	"strings"
)

type ruleBasedVerifier struct{}

func NewAIVerifier() AIVerifier {
	return &ruleBasedVerifier{}
}

func (v *ruleBasedVerifier) Verify(ctx context.Context, contributionType string, data json.RawMessage) (string, float64, string, error) {
	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		return "reject", 0.8, "invalid JSON data", nil
	}

	if len(body) == 0 {
		return "reject", 0.9, "empty contribution data", nil
	}

	switch contributionType {
	case "lyrics":
		return v.verifyLyrics(body)
	case "translation":
		return v.verifyTranslation(body)
	case "credits":
		return v.verifyCredits(body)
	case "metadata":
		return v.verifyMetadata(body)
	case "album_art":
		return v.verifyAlbumArt(body)
	case "bio":
		return v.verifyBio(body)
	default:
		return "reject", 0.5, "unknown contribution type", nil
	}
}

func (v *ruleBasedVerifier) verifyLyrics(body map[string]any) (string, float64, string, error) {
	text, ok := body["text"].(string)
	if !ok || strings.TrimSpace(text) == "" {
		return "reject", 0.9, "lyrics text is required", nil
	}

	wordCount := len(strings.Fields(text))
	if wordCount < 5 {
		return "reject", 0.8, "lyrics too short", nil
	}

	lang, hasLang := body["language"].(string)
	if hasLang && lang != "" {
		confidence := math.Min(0.5+float64(wordCount)/500.0*0.4, 0.9)
		return "approve", confidence, "valid lyrics content", nil
	}

	return "approve", 0.5, "lyrics content ok, language not specified", nil
}

func (v *ruleBasedVerifier) verifyTranslation(body map[string]any) (string, float64, string, error) {
	text, ok := body["text"].(string)
	if !ok || strings.TrimSpace(text) == "" {
		return "reject", 0.9, "translation text is required", nil
	}

	lang, ok := body["language"].(string)
	if !ok || lang == "" {
		return "reject", 0.7, "translation language is required", nil
	}

	wordCount := len(strings.Fields(text))
	if wordCount < 3 {
		return "reject", 0.8, "translation too short", nil
	}

	return "approve", 0.6, "translation content ok", nil
}

func (v *ruleBasedVerifier) verifyCredits(body map[string]any) (string, float64, string, error) {
	roles, ok := body["roles"].([]any)
	if !ok || len(roles) == 0 {
		arr, hasArr := body["credits"].([]any)
		if !hasArr || len(arr) == 0 {
			return "reject", 0.7, "credits must include at least one role", nil
		}
		return "approve", 0.6, "credits provided", nil
	}

	for _, role := range roles {
		if r, ok := role.(map[string]any); ok {
			if name, has := r["name"].(string); !has || name == "" {
				return "reject", 0.6, "credit entry missing name", nil
			}
		}
	}

	return "approve", 0.6, "credits ok", nil
}

func (v *ruleBasedVerifier) verifyMetadata(body map[string]any) (string, float64, string, error) {
	hasField := false
	for _, key := range []string{"title", "artist", "album", "genre", "year", "track_number", "label", "upc", "isrc"} {
		if val, ok := body[key]; ok && val != nil {
			if s, ok := val.(string); ok && strings.TrimSpace(s) != "" {
				hasField = true
				break
			}
			if _, ok := val.(float64); ok {
				hasField = true
				break
			}
		}
	}

	if !hasField {
		return "reject", 0.6, "metadata must include at least one known field", nil
	}
	return "approve", 0.5, "metadata ok", nil
}

func (v *ruleBasedVerifier) verifyAlbumArt(body map[string]any) (string, float64, string, error) {
	url, ok := body["url"].(string)
	if !ok || url == "" {
		url, ok = body["image_url"].(string)
		if !ok || url == "" {
			return "reject", 0.8, "album art URL is required", nil
		}
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "reject", 0.7, "album art URL must start with http:// or https://", nil
	}

	return "approve", 0.5, "album art URL ok", nil
}

func (v *ruleBasedVerifier) verifyBio(body map[string]any) (string, float64, string, error) {
	text, ok := body["text"].(string)
	if !ok || strings.TrimSpace(text) == "" {
		text, ok = body["bio"].(string)
		if !ok || strings.TrimSpace(text) == "" {
			return "reject", 0.8, "bio text is required", nil
		}
	}

	wordCount := len(strings.Fields(text))
	if wordCount < 10 {
		return "reject", 0.6, "bio too short (minimum 10 words)", nil
	}

	return "approve", 0.5, "bio content ok", nil
}
