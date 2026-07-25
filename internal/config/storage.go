package config

import (
	"os"
	"strconv"
	"time"
)

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
	driver := os.Getenv("STORAGE_DRIVER")
	if driver == "" {
		driver = "local"
	}
	localBaseDir := os.Getenv("STORAGE_LOCAL_BASE_DIR")
	if localBaseDir == "" {
		localBaseDir = "uploads"
	}
	localBaseURL := os.Getenv("STORAGE_LOCAL_BASE_URL")
	if localBaseURL == "" {
		localBaseURL = "/uploads"
	}
	region := os.Getenv("S3_REGION")
	if region == "" {
		region = "us-east-1"
	}

	usePathStyle := false
	if v := os.Getenv("S3_USE_PATH_STYLE"); v != "" {
		switch v {
		case "true", "TRUE", "1", "yes", "YES", "y", "Y":
			usePathStyle = true
		}
	}
	presignURLs := false
	if v := os.Getenv("S3_PRESIGN_URLS"); v != "" {
		switch v {
		case "true", "TRUE", "1", "yes", "YES", "y", "Y":
			presignURLs = true
		}
	}
	presignTTL := 900
	if v := os.Getenv("S3_PRESIGN_TTL_SECONDS"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			presignTTL = parsed
		}
	}

	return StorageConfig{
		Driver: driver,

		Local: LocalStorageConfig{
			BaseDir: localBaseDir,
			BaseURL: localBaseURL,
		},

		S3: S3StorageConfig{
			Bucket:          os.Getenv("S3_BUCKET"),
			Region:          region,
			Endpoint:        os.Getenv("S3_ENDPOINT"),
			AccessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
			PublicBaseURL:   os.Getenv("S3_PUBLIC_BASE_URL"),
			UsePathStyle:    usePathStyle,
			PresignURLs:     presignURLs,
			PresignTTL:      time.Duration(presignTTL) * time.Second,
		},
	}
}
