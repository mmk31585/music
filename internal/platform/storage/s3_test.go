package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewS3Storage_RequiresBucket(t *testing.T) {
	_, err := NewS3Storage(context.Background(), S3Config{})
	require.Error(t, err)
	require.EqualError(t, err, "s3 bucket is required")
}

func TestNewS3Storage_StaticCredsBuildsClient(t *testing.T) {
	s, err := NewS3Storage(context.Background(), S3Config{
		Bucket:          "music-media",
		Region:          "us-east-1",
		AccessKeyID:     "AKIA_TEST",
		SecretAccessKey: "secret",
	})
	require.NoError(t, err)
	require.NotNil(t, s)
	require.Equal(t, "music-media", s.bucket)
	require.False(t, s.presign)
	require.Equal(t, "", s.publicURL)
}

func TestNewS3Storage_RoleARNBuildsClient(t *testing.T) {
	s, err := NewS3Storage(context.Background(), S3Config{
		Bucket:  "music-media",
		Region:  "eu-west-1",
		RoleARN: "arn:aws:iam::123456789012:role/muse-uploads",
	})
	require.NoError(t, err)
	require.NotNil(t, s)
	require.Equal(t, "music-media", s.bucket)
}

func TestS3Storage_GetURL_PublicBaseURL(t *testing.T) {
	s := &S3Storage{
		bucket:    "music-media",
		publicURL: "https://cdn.example.com/music-media",
		presign:   false,
	}
	url, err := s.GetURL(context.Background(), "track-audio/abc.mp3")
	require.NoError(t, err)
	require.Equal(t, "https://cdn.example.com/music-media/track-audio/abc.mp3", url)
}

func TestS3Storage_GetURL_PresignNilPresigner(t *testing.T) {
	// presign with a nil presigner should return an error, not panic
	s := &S3Storage{
		bucket:  "music-media",
		presign: true,
	}
	_, err := s.GetURL(context.Background(), "track-audio/abc.mp3")
	require.Error(t, err)
}

func TestS3Storage_Upload_NilClientReturnsError(t *testing.T) {
	s := &S3Storage{client: nil}
	err := s.Upload(context.Background(), "", nil, 0, "")
	require.Error(t, err)
}
