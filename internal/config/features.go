package config

type FeaturesConfig struct {
	Analytics      bool
	Recommendation bool
	Search         bool
	Social         bool
	Reactions      bool
	Creator        bool
	Moderation     bool
	AI             bool
	Contribution   bool
	Gamification   bool
	Tips           bool
	Subscription   bool
	Notification   bool
}

func loadFeaturesConfig() FeaturesConfig {
	aiEnabled := getEnvAsBool("FEATURE_AI_ENABLED", true)
	if v := getEnv("FEATURE_AI_ENABLED", ""); v == "" {
		aiEnabled = getEnvAsBool("AI_ENABLED", true)
	}

	return FeaturesConfig{
		Analytics:      getEnvAsBool("FEATURE_ANALYTICS_ENABLED", true),
		Recommendation: getEnvAsBool("FEATURE_RECOMMENDATION_ENABLED", true),
		Search:         getEnvAsBool("FEATURE_SEARCH_ENABLED", true),
		Social:         getEnvAsBool("FEATURE_SOCIAL_ENABLED", true),
		Reactions:      getEnvAsBool("FEATURE_REACTIONS_ENABLED", true),
		Creator:        getEnvAsBool("FEATURE_CREATOR_ENABLED", true),
		Moderation:     getEnvAsBool("FEATURE_MODERATION_ENABLED", true),
		AI:             aiEnabled,
		Contribution:   getEnvAsBool("FEATURE_CONTRIBUTION_ENABLED", true),
		Gamification:   getEnvAsBool("FEATURE_GAMIFICATION_ENABLED", true),
		Tips:           getEnvAsBool("FEATURE_TIPS_ENABLED", true),
		Subscription:   getEnvAsBool("FEATURE_SUBSCRIPTION_ENABLED", true),
		Notification:   getEnvAsBool("FEATURE_NOTIFICATION_ENABLED", true),
	}
}
