package config

import "os"

type AIConfig struct {
	OpenAIEndpoint string
	OpenAIKey      string
	EmbeddingModel string
	Enabled        bool
}

func loadAIConfig() AIConfig {
	endpoint := os.Getenv("AI_OPENAI_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1"
	}
	model := os.Getenv("AI_EMBEDDING_MODEL")
	if model == "" {
		model = "text-embedding-3-small"
	}
	enabled := os.Getenv("AI_ENABLED")
	if enabled == "" {
		enabled = "true"
	}
	return AIConfig{
		OpenAIEndpoint: endpoint,
		OpenAIKey:      os.Getenv("AI_OPENAI_KEY"),
		EmbeddingModel: model,
		Enabled:        enabled == "true",
	}
}
