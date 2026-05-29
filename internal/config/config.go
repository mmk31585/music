package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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

	// Legacy song module compatibility. The active upload path lives in
	// internal/modules/media, but these keep old packages buildable until the
	// legacy module is either migrated or explicitly removed.
	AllowedAudioExt map[string]bool
	MaxUploadSize   int64
	UploadDir       string
}

type MediaConfig struct {
	BasePath          string
	PublicBase        string
	MaxFileSizeBytes  int64
	MaxImageSizeBytes int64
	MaxAudioSizeBytes int64
}

type OpenSearchConfig struct {
	URL      string
	Username string
	Password string
	Index    string
}

type AuthConfig struct {
	JWTAccessSecret  string
	JWTRefreshSecret string
	AccessTTL        time.Duration
	RefreshTTL       time.Duration
}

type AppConfig struct {
	Name            string
	Env             string
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type PostgresConfig struct {
	URL string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type LogConfig struct {
	Level string
}

type CORSConfig struct {
	AllowedOrigins []string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name:            getEnv("APP_NAME", "musicapp"),
			Env:             getEnv("APP_ENV", "development"),
			Host:            getEnv("APP_HOST", "0.0.0.0"),
			Port:            getEnv("APP_PORT", "8080"),
			ReadTimeout:     getEnvAsDurationSeconds("APP_READ_TIMEOUT", 15),
			WriteTimeout:    getEnvAsDurationSeconds("APP_WRITE_TIMEOUT", 15),
			IdleTimeout:     getEnvAsDurationSeconds("APP_IDLE_TIMEOUT", 60),
			ShutdownTimeout: getEnvAsDurationSeconds("APP_SHUTDOWN_TIMEOUT", 10),
		},
		Postgres: PostgresConfig{
			URL: getEnv("POSTGRES_URL", ""),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "debug"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsStringSlice("CORS_ALLOWED_ORIGINS", []string{
				"http://localhost:3000",
				"http://localhost:5173",
			}),
		},
		Auth: AuthConfig{
			JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "change_me_access_secret"),
			JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "change_me_refresh_secret"),
			AccessTTL:        time.Duration(getEnvAsInt("JWT_ACCESS_TTL_MINUTES", 15)) * time.Minute,
			RefreshTTL:       time.Duration(getEnvAsInt("JWT_REFRESH_TTL_DAYS", 30)) * 24 * time.Hour,
		},
		Media: MediaConfig{
			BasePath:          getEnv("MEDIA_BASE_PATH", "./media"),
			PublicBase:        getEnv("MEDIA_PUBLIC_BASE", "http://localhost:8080/media"),
			MaxFileSizeBytes:  mbToBytes(getInt64Env("MEDIA_MAX_FILE_SIZE_MB", 30)),
			MaxImageSizeBytes: mbToBytes(getInt64Env("MEDIA_MAX_IMAGE_SIZE_MB", 5)),
			MaxAudioSizeBytes: mbToBytes(getInt64Env("MEDIA_MAX_AUDIO_SIZE_MB", 30)),
		},
		OpenSearch: OpenSearchConfig{
			URL:      getEnv("OPENSEARCH_URL", ""),
			Username: getEnv("OPENSEARCH_USERNAME", ""),
			Password: getEnv("OPENSEARCH_PASSWORD", ""),
			Index:    getEnv("OPENSEARCH_INDEX", "tracks"),
		},
	}

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

	if cfg.Postgres.URL == "" {
		return nil, fmt.Errorf("POSTGRES_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getInt64Env(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvAsDurationSeconds(key string, fallback int) time.Duration {
	return time.Duration(getEnvAsInt(key, fallback)) * time.Second
}

func getEnvAsStringSlice(key string, fallback []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			result = append(result, item)
		}
	}

	if len(result) == 0 {
		return fallback
	}

	return result
}

func mbToBytes(mb int64) int64 {
	return mb << 20
}
