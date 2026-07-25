package config

import (
	"os"
	"strconv"
)

type MediaConfig struct {
	BasePath          string
	PublicBase        string
	MaxFileSizeBytes  int64
	MaxImageSizeBytes int64
	MaxAudioSizeBytes int64
}

func loadMediaConfig() MediaConfig {
	basePath := os.Getenv("MEDIA_BASE_PATH")
	if basePath == "" {
		basePath = "./media"
	}
	publicBase := os.Getenv("MEDIA_PUBLIC_BASE")
	if publicBase == "" {
		publicBase = "http://localhost:8080/media"
	}

	maxFileSizeMB := int64(30)
	if v := os.Getenv("MEDIA_MAX_FILE_SIZE_MB"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			maxFileSizeMB = parsed
		}
	}
	maxImageSizeMB := int64(5)
	if v := os.Getenv("MEDIA_MAX_IMAGE_SIZE_MB"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			maxImageSizeMB = parsed
		}
	}
	maxAudioSizeMB := int64(30)
	if v := os.Getenv("MEDIA_MAX_AUDIO_SIZE_MB"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			maxAudioSizeMB = parsed
		}
	}

	return MediaConfig{
		BasePath:          basePath,
		PublicBase:        publicBase,
		MaxFileSizeBytes:  maxFileSizeMB << 20,
		MaxImageSizeBytes: maxImageSizeMB << 20,
		MaxAudioSizeBytes: maxAudioSizeMB << 20,
	}
}
