package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFeaturesConfig_Defaults(t *testing.T) {
	cfg := loadFeaturesConfig()

	assert.True(t, cfg.Analytics)
	assert.True(t, cfg.Recommendation)
	assert.True(t, cfg.Search)
	assert.True(t, cfg.Social)
	assert.True(t, cfg.Reactions)
	assert.True(t, cfg.Creator)
	assert.True(t, cfg.Moderation)
	assert.True(t, cfg.AI)
	assert.True(t, cfg.Contribution)
	assert.True(t, cfg.Gamification)
	assert.True(t, cfg.Tips)
	assert.True(t, cfg.Subscription)
	assert.True(t, cfg.Notification)
}

func TestLoadFeaturesConfig_DisableSingleFeature(t *testing.T) {
	os.Setenv("FEATURE_SOCIAL_ENABLED", "false")
	defer os.Unsetenv("FEATURE_SOCIAL_ENABLED")

	cfg := loadFeaturesConfig()

	assert.False(t, cfg.Social)
	assert.True(t, cfg.Analytics)
	assert.True(t, cfg.Recommendation)
	assert.True(t, cfg.Search)
	assert.True(t, cfg.Reactions)
	assert.True(t, cfg.Creator)
	assert.True(t, cfg.Moderation)
	assert.True(t, cfg.AI)
	assert.True(t, cfg.Contribution)
	assert.True(t, cfg.Gamification)
	assert.True(t, cfg.Tips)
	assert.True(t, cfg.Subscription)
	assert.True(t, cfg.Notification)
}

func TestLoadFeaturesConfig_DisableMultipleFeatures(t *testing.T) {
	os.Setenv("FEATURE_ANALYTICS_ENABLED", "false")
	os.Setenv("FEATURE_GAMIFICATION_ENABLED", "false")
	os.Setenv("FEATURE_TIPS_ENABLED", "false")
	defer os.Unsetenv("FEATURE_ANALYTICS_ENABLED")
	defer os.Unsetenv("FEATURE_GAMIFICATION_ENABLED")
	defer os.Unsetenv("FEATURE_TIPS_ENABLED")

	cfg := loadFeaturesConfig()

	assert.False(t, cfg.Analytics)
	assert.False(t, cfg.Gamification)
	assert.False(t, cfg.Tips)
	assert.True(t, cfg.Recommendation)
	assert.True(t, cfg.Search)
	assert.True(t, cfg.Social)
}

func TestLoadFeaturesConfig_AIEnabledFallback(t *testing.T) {
	os.Unsetenv("FEATURE_AI_ENABLED")
	os.Setenv("AI_ENABLED", "false")
	defer os.Unsetenv("AI_ENABLED")

	cfg := loadFeaturesConfig()

	assert.False(t, cfg.AI)
}

func TestLoadFeaturesConfig_FEATURE_AI_TakesPrecedence(t *testing.T) {
	os.Setenv("FEATURE_AI_ENABLED", "true")
	os.Setenv("AI_ENABLED", "false")
	defer os.Unsetenv("FEATURE_AI_ENABLED")
	defer os.Unsetenv("AI_ENABLED")

	cfg := loadFeaturesConfig()

	assert.True(t, cfg.AI)
}

func TestLoadFeaturesConfig_CaseInsensitiveFalse(t *testing.T) {
	os.Setenv("FEATURE_SEARCH_ENABLED", "FALSE")
	defer os.Unsetenv("FEATURE_SEARCH_ENABLED")

	cfg := loadFeaturesConfig()

	assert.False(t, cfg.Search)
}

func TestLoadFeaturesConfig_DisableAllNonCore(t *testing.T) {
	keys := []string{
		"FEATURE_ANALYTICS_ENABLED",
		"FEATURE_RECOMMENDATION_ENABLED",
		"FEATURE_SEARCH_ENABLED",
		"FEATURE_SOCIAL_ENABLED",
		"FEATURE_REACTIONS_ENABLED",
		"FEATURE_CREATOR_ENABLED",
		"FEATURE_MODERATION_ENABLED",
		"FEATURE_AI_ENABLED",
		"FEATURE_CONTRIBUTION_ENABLED",
		"FEATURE_GAMIFICATION_ENABLED",
		"FEATURE_TIPS_ENABLED",
		"FEATURE_SUBSCRIPTION_ENABLED",
		"FEATURE_NOTIFICATION_ENABLED",
	}
	for _, k := range keys {
		os.Setenv(k, "false")
	}
	defer func() {
		for _, k := range keys {
			os.Unsetenv(k)
		}
	}()

	cfg := loadFeaturesConfig()

	assert.False(t, cfg.Analytics)
	assert.False(t, cfg.Recommendation)
	assert.False(t, cfg.Search)
	assert.False(t, cfg.Social)
	assert.False(t, cfg.Reactions)
	assert.False(t, cfg.Creator)
	assert.False(t, cfg.Moderation)
	assert.False(t, cfg.AI)
	assert.False(t, cfg.Contribution)
	assert.False(t, cfg.Gamification)
	assert.False(t, cfg.Tips)
	assert.False(t, cfg.Subscription)
	assert.False(t, cfg.Notification)
}

func TestLoadFeaturesConfig_InvalidValueDefaultsToTrue(t *testing.T) {
	os.Setenv("FEATURE_CONTRIBUTION_ENABLED", "maybe")
	defer os.Unsetenv("FEATURE_CONTRIBUTION_ENABLED")

	cfg := loadFeaturesConfig()

	assert.True(t, cfg.Contribution)
}
