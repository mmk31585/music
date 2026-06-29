package config

type OpenRouterConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

func loadOpenRouterConfig() OpenRouterConfig {
	return OpenRouterConfig{
		APIKey:  getEnv("OPENROUTER_API_KEY", ""),
		Model:   getEnv("OPENROUTER_MODEL", "openrouter/free"),
		BaseURL: getEnv("OPENROUTER_BASE_URL", "https://openrouter.ai/api/v1"),
	}
}
