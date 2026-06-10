package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_PopulatesFeatures(t *testing.T) {
	os.Setenv("POSTGRES_URL", "postgres://localhost:5432/test?sslmode=disable")
	defer os.Unsetenv("POSTGRES_URL")

	os.Setenv("FEATURE_SOCIAL_ENABLED", "false")
	os.Setenv("FEATURE_GAMIFICATION_ENABLED", "false")
	defer os.Unsetenv("FEATURE_SOCIAL_ENABLED")
	defer os.Unsetenv("FEATURE_GAMIFICATION_ENABLED")

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.False(t, cfg.Features.Social)
	assert.False(t, cfg.Features.Gamification)
	assert.True(t, cfg.Features.Analytics)
	assert.True(t, cfg.Features.Recommendation)
	assert.True(t, cfg.Features.Search)
	assert.True(t, cfg.Features.Reactions)
	assert.True(t, cfg.Features.Creator)
	assert.True(t, cfg.Features.Moderation)
	assert.True(t, cfg.Features.AI)
	assert.True(t, cfg.Features.Contribution)
	assert.True(t, cfg.Features.Tips)
	assert.True(t, cfg.Features.Subscription)
	assert.True(t, cfg.Features.Notification)
}
