package config

type MediaConfig struct {
	BasePath          string
	PublicBase        string
	MaxFileSizeBytes  int64
	MaxImageSizeBytes int64
	MaxAudioSizeBytes int64
}

func loadMediaConfig() MediaConfig {
	return MediaConfig{
		BasePath:          getEnv("MEDIA_BASE_PATH", "./media"),
		PublicBase:        getEnv("MEDIA_PUBLIC_BASE", "http://localhost:8080/media"),
		MaxFileSizeBytes:  mbToBytes(getInt64Env("MEDIA_MAX_FILE_SIZE_MB", 30)),
		MaxImageSizeBytes: mbToBytes(getInt64Env("MEDIA_MAX_IMAGE_SIZE_MB", 5)),
		MaxAudioSizeBytes: mbToBytes(getInt64Env("MEDIA_MAX_AUDIO_SIZE_MB", 30)),
	}
}
