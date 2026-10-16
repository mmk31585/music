package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// resolveURL resolves ${VAR} references in a connection string.
// If the entire value is ${VAR_NAME}, it is replaced with os.Getenv(VAR_NAME).
// This handles platforms (e.g. Vercel) that inject a variable reference
// instead of the resolved value.
func resolveURL(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		name := value[2 : len(value)-1]
		if v := os.Getenv(name); v != "" {
			return v
		}
	}
	return value
}

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
	Sentry     SentryConfig

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
		Sentry:     loadSentryConfig(),
	}

	// Resolve ${VAR} references in all connection URLs so platforms like Vercel
	// that inject variable references instead of resolved values work correctly.
	cfg.Postgres.URL = resolveURL(cfg.Postgres.URL)
	cfg.Postgres.UnpooledURL = resolveURL(cfg.Postgres.UnpooledURL)
	cfg.Redis.Addr = resolveURL(cfg.Redis.Addr)
	cfg.OpenSearch.URL = resolveURL(cfg.OpenSearch.URL)
	cfg.MLService.BaseURL = resolveURL(cfg.MLService.BaseURL)
	cfg.OpenRouter.BaseURL = resolveURL(cfg.OpenRouter.BaseURL)
	cfg.Storage.S3.Endpoint = resolveURL(cfg.Storage.S3.Endpoint)
	cfg.Storage.S3.PublicBaseURL = resolveURL(cfg.Storage.S3.PublicBaseURL)
	cfg.Media.BasePath = resolveURL(cfg.Media.BasePath)
	cfg.Media.PublicBase = resolveURL(cfg.Media.PublicBase)

	cfg.applyLegacyUploadCompatibility()

	if cfg.Postgres.URL == "" && cfg.Postgres.UnpooledURL == "" {
		return nil, fmt.Errorf("POSTGRES_URL or DATABASE_URL is required")
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
