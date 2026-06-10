package features

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"music/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter(cfg *config.FeaturesConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewHandler(cfg)
	rg := r.Group("/api/v1")
	RegisterRoutes(rg, handler)
	return r
}

func TestGetFeatures_AllEnabled(t *testing.T) {
	cfg := &config.FeaturesConfig{
		Analytics:      true,
		Recommendation: true,
		Search:         true,
		Social:         true,
		Reactions:      true,
		Creator:        true,
		Moderation:     true,
		AI:             true,
		Contribution:   true,
		Gamification:   true,
		Tips:           true,
		Subscription:   true,
		Notification:   true,
	}

	r := setupRouter(cfg)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/features", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, true, resp.Data["social"])
	assert.Equal(t, true, resp.Data["ai"])
	assert.Equal(t, true, resp.Data["analytics"])
}

func TestGetFeatures_SomeDisabled(t *testing.T) {
	cfg := &config.FeaturesConfig{
		Analytics:      true,
		Recommendation: true,
		Search:         true,
		Social:         false,
		Reactions:      true,
		Creator:        false,
		Moderation:     true,
		AI:             true,
		Contribution:   true,
		Gamification:   true,
		Tips:           true,
		Subscription:   true,
		Notification:   true,
	}

	r := setupRouter(cfg)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/features", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, false, resp.Data["social"])
	assert.Equal(t, false, resp.Data["creator"])
	assert.Equal(t, true, resp.Data["analytics"])
}

func TestGetFeatures_AllDisabled(t *testing.T) {
	falseVal := false
	cfg := &config.FeaturesConfig{
		Analytics:      falseVal,
		Recommendation: falseVal,
		Search:         falseVal,
		Social:         falseVal,
		Reactions:      falseVal,
		Creator:        falseVal,
		Moderation:     falseVal,
		AI:             falseVal,
		Contribution:   falseVal,
		Gamification:   falseVal,
		Tips:           falseVal,
		Subscription:   falseVal,
		Notification:   falseVal,
	}

	r := setupRouter(cfg)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/features", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, 13, len(resp.Data))
	for k, v := range resp.Data {
		assert.False(t, v.(bool), "feature %s should be false", k)
	}
}

func TestGetFeatures_ResponseKeys(t *testing.T) {
	cfg := &config.FeaturesConfig{
		Analytics: true, Recommendation: true, Search: true,
		Social: true, Reactions: true, Creator: true,
		Moderation: true, AI: true, Contribution: true,
		Gamification: true, Tips: true, Subscription: true, Notification: true,
	}

	r := setupRouter(cfg)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/features", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	expectedKeys := []string{
		"analytics", "recommendation", "search", "social",
		"reactions", "creator", "moderation", "ai",
		"contribution", "gamification", "tips", "subscription", "notification",
	}
	for _, key := range expectedKeys {
		_, exists := resp.Data[key]
		assert.True(t, exists, "response should contain key %s", key)
	}
}
