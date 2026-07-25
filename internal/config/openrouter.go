package config

import "os"

type OpenRouterConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

func loadOpenRouterConfig() OpenRouterConfig {
	model := os.Getenv("OPENROUTER_MODEL")
	if model == "" {
		model = "openrouter/free"
	}
	baseURL := os.Getenv("OPENROUTER_BASE_URL")
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1"
	}
	return OpenRouterConfig{
		APIKey:  os.Getenv("OPENROUTER_API_KEY"),
		Model:   model,
		BaseURL: baseURL,
	}
}
