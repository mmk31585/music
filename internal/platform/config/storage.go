package config

import "time"

type StorageConfig struct {
	Driver string

	Local LocalStorageConfig
	S3    S3StorageConfig
}

type LocalStorageConfig struct {
	BaseDir string
	BaseURL string
}

type S3StorageConfig struct {
	Bucket          string
	Region          string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	PublicBaseURL   string
	UsePathStyle    bool
	PresignURLs     bool
	PresignTTL      time.Duration
}

func loadStorageConfig() StorageConfig {
	return StorageConfig{
		Driver: getEnv("STORAGE_DRIVER", "local"),

		Local: LocalStorageConfig{
			BaseDir: getEnv("STORAGE_LOCAL_BASE_DIR", "uploads"),
			BaseURL: getEnv("STORAGE_LOCAL_BASE_URL", "/uploads"),
		},

		S3: S3StorageConfig{
			Bucket:          getEnv("S3_BUCKET", ""),
			Region:          getEnv("S3_REGION", "us-east-1"),
			Endpoint:        getEnv("S3_ENDPOINT", ""),
			AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),
			PublicBaseURL:   getEnv("S3_PUBLIC_BASE_URL", ""),
			UsePathStyle:    getEnvAsBool("S3_USE_PATH_STYLE", false),
			PresignURLs:     getEnvAsBool("S3_PRESIGN_URLS", false),
			PresignTTL:      getEnvAsDurationSeconds("S3_PRESIGN_TTL_SECONDS", 900),
		},
	}
}
