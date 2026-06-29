package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	App        AppConfig
	Postgres   PostgresConfig
	Redis      RedisConfig
	Log        LogConfig
	CORS       CORSConfig
	Auth       AuthConfig
	Media      MediaConfig
	OpenSearch OpenSearchConfig
	Storage    StorageConfig
	Enrichment EnrichmentConfig
	Features   FeaturesConfig
	AI         AIConfig
	Payment    PaymentConfig
	MLService  MLServiceConfig
	OpenRouter OpenRouterConfig

	// Legacy song module compatibility. The active upload path lives in
	// internal/modules/media, but these keep old packages buildable until the
	// legacy module is either migrated or explicitly removed.
	AllowedAudioExt map[string]bool
	MaxUploadSize   int64
	UploadDir       string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		App:        loadAppConfig(),
		Postgres:   loadPostgresConfig(),
		Redis:      loadRedisConfig(),
		Log:        loadLogConfig(),
		CORS:       loadCORSConfig(),
		Auth:       loadAuthConfig(),
		Media:      loadMediaConfig(),
		OpenSearch: loadOpenSearchConfig(),
		Storage:    loadStorageConfig(),
		Enrichment: loadEnrichmentConfig(),
		Features:   loadFeaturesConfig(),
		AI:         loadAIConfig(),
		Payment:    loadPaymentConfig(),
		MLService:  loadMLServiceConfig(),
		OpenRouter: loadOpenRouterConfig(),
	}

	cfg.applyLegacyUploadCompatibility()

	if cfg.Postgres.URL == "" {
		return nil, fmt.Errorf("POSTGRES_URL is required")
	}

	// Validate critical configuration — panics in production if secrets are weak
	ValidateProductionConfig(cfg)

	return cfg, nil
}

func (cfg *Config) applyLegacyUploadCompatibility() {
	cfg.AllowedAudioExt = map[string]bool{
		".mp3":  true,
		".ogg":  true,
		".flac": true,
		".wav":  true,
		".m4a":  true,
		".aac":  true,
	}

	cfg.MaxUploadSize = cfg.Media.MaxAudioSizeBytes
	cfg.UploadDir = cfg.Media.BasePath
}
