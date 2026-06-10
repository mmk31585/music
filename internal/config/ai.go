package config

type AIConfig struct {
	OpenAIEndpoint string
	OpenAIKey      string
	EmbeddingModel string
	Enabled        bool
}

func loadAIConfig() AIConfig {
	return AIConfig{
		OpenAIEndpoint: getEnv("AI_OPENAI_ENDPOINT", "https://api.openai.com/v1"),
		OpenAIKey:      getEnv("AI_OPENAI_KEY", ""),
		EmbeddingModel: getEnv("AI_EMBEDDING_MODEL", "text-embedding-3-small"),
		Enabled:        getEnv("AI_ENABLED", "true") == "true",
	}
}
