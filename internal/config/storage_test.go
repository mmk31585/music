package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadStorageConfig_S3Fields(t *testing.T) {
	os.Setenv("POSTGRES_URL", "postgres://localhost:5432/test?sslmode=disable")
	defer os.Unsetenv("POSTGRES_URL")
	os.Setenv("JWT_ACCESS_SECRET", "test-access-secret-that-is-at-least-32-bytes-long")
	os.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret-that-is-at-least-32-byte")
	defer os.Unsetenv("JWT_ACCESS_SECRET")
	defer os.Unsetenv("JWT_REFRESH_SECRET")

	os.Setenv("STORAGE_DRIVER", "s3")
	os.Setenv("S3_BUCKET_NAME", "music-media")
	os.Setenv("AWS_REGION", "eu-west-1")
	os.Setenv("S3_ENDPOINT", "https://s3.eu-west-1.amazonaws.com")
	os.Setenv("S3_ACCESS_KEY_ID", "AKIA...")
	os.Setenv("S3_SECRET_ACCESS_KEY", "secret")
	os.Setenv("AWS_ROLE_ARN", "arn:aws:iam::123456789012:role/muse-uploads")
	os.Setenv("AWS_ACCOUNT_ID", "123456789012")
	os.Setenv("S3_PUBLIC_BASE_URL", "https://cdn.example.com/music-media")
	os.Setenv("S3_USE_PATH_STYLE", "true")
	os.Setenv("S3_PRESIGN_URLS", "true")
	os.Setenv("S3_PRESIGN_TTL_SECONDS", "3600")
	defer func() {
		for _, k := range []string{
			"STORAGE_DRIVER", "S3_BUCKET_NAME", "AWS_REGION", "S3_ENDPOINT",
			"S3_ACCESS_KEY_ID", "S3_SECRET_ACCESS_KEY", "AWS_ROLE_ARN",
			"AWS_ACCOUNT_ID", "S3_PUBLIC_BASE_URL", "S3_USE_PATH_STYLE",
			"S3_PRESIGN_URLS", "S3_PRESIGN_TTL_SECONDS",
		} {
			os.Unsetenv(k)
		}
	}()

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "s3", cfg.Storage.Driver)
	assert.Equal(t, "music-media", cfg.Storage.S3.Bucket)
	assert.Equal(t, "eu-west-1", cfg.Storage.S3.Region)
	assert.Equal(t, "https://s3.eu-west-1.amazonaws.com", cfg.Storage.S3.Endpoint)
	assert.Equal(t, "AKIA...", cfg.Storage.S3.AccessKeyID)
	assert.Equal(t, "secret", cfg.Storage.S3.SecretAccessKey)
	assert.Equal(t, "arn:aws:iam::123456789012:role/muse-uploads", cfg.Storage.S3.RoleARN)
	assert.Equal(t, "123456789012", cfg.Storage.S3.AccountID)
	assert.Equal(t, "https://cdn.example.com/music-media", cfg.Storage.S3.PublicBaseURL)
	assert.True(t, cfg.Storage.S3.UsePathStyle)
	assert.True(t, cfg.Storage.S3.PresignURLs)
	assert.Equal(t, 3600, int(cfg.Storage.S3.PresignTTL.Seconds()))
}
