package config

import "os"

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
	aiEnabled := true
	if v := os.Getenv("FEATURE_AI_ENABLED"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			aiEnabled = false
		}
	}
	if os.Getenv("FEATURE_AI_ENABLED") == "" {
		if v := os.Getenv("AI_ENABLED"); v != "" {
			switch v {
			case "false", "FALSE", "0", "no", "NO", "n", "N":
				aiEnabled = false
			}
		}
	}

	analytics := true
	if v := os.Getenv("FEATURE_ANALYTICS_ENABLED"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			analytics = false
		}
	}
	recommendation := true
	if v := os.Getenv("FEATURE_RECOMMENDATION_ENABLED"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			recommendation = false
		}
	}
	search := true
	if v := os.Getenv("FEATURE_SEARCH_ENABLED"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			search = false
		}
	}
	social := true
	if v := os.Getenv("FEATURE_SOCIAL_ENABLED"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			social = false
		}
	}
	reactions := true
	if v := os.Getenv("FEATURE_REACTIONS_ENABLED"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			reactions = false
		}
	}
	creator := true
	if v := os.Getenv("FEATURE_CREATOR_ENABLED"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			creator = false
		}
	}
	moderation := true
	if v := os.Getenv("FEATURE_MODERATION_ENABLED"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			moderation = false
		}
	}
	contribution := true
	if v := os.Getenv("FEATURE_CONTRIBUTION_ENABLED"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			contribution = false
		}
	}
	gamification := true
	if v := os.Getenv("FEATURE_GAMIFICATION_ENABLED"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			gamification = false
		}
	}
	tips := true
	if v := os.Getenv("FEATURE_TIPS_ENABLED"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			tips = false
		}
	}
	subscription := true
	if v := os.Getenv("FEATURE_SUBSCRIPTION_ENABLED"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			subscription = false
		}
	}
	notification := true
	if v := os.Getenv("FEATURE_NOTIFICATION_ENABLED"); v != "" {
		switch v {
		case "false", "FALSE", "0", "no", "NO", "n", "N":
			notification = false
		}
	}

	return FeaturesConfig{
		Analytics:      analytics,
		Recommendation: recommendation,
		Search:         search,
		Social:         social,
		Reactions:      reactions,
		Creator:        creator,
		Moderation:     moderation,
		AI:             aiEnabled,
		Contribution:   contribution,
		Gamification:   gamification,
		Tips:           tips,
		Subscription:   subscription,
		Notification:   notification,
	}
}
