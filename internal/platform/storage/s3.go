package storage

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
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

type S3Storage struct {
	client     *s3.Client
	presigner  *s3.PresignClient
	bucket     string
	publicURL  string
	presign    bool
	presignTTL time.Duration
}

func NewS3Storage(ctx context.Context, cfg S3Config) (*S3Storage, error) {
	if cfg.Bucket == "" {
		return nil, errors.New("s3 bucket is required")
	}

	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	loadOptions := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.Region),
	}

	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		loadOptions = append(loadOptions, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, loadOptions...)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}

		o.UsePathStyle = cfg.UsePathStyle
	})

	if cfg.PresignTTL == 0 {
		cfg.PresignTTL = 15 * time.Minute
	}

	return &S3Storage{
		client:     client,
		presigner:  s3.NewPresignClient(client),
		bucket:     cfg.Bucket,
		publicURL:  strings.TrimRight(cfg.PublicBaseURL, "/"),
		presign:    cfg.PresignURLs,
		presignTTL: cfg.PresignTTL,
	}, nil
}

func (s *S3Storage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	key = strings.TrimLeft(key, "/")

	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   reader,
	}

	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	_, err := s.client.PutObject(ctx, input)
	return err
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	key = strings.TrimLeft(key, "/")

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	return err
}

func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	key = strings.TrimLeft(key, "/")

	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		// Keep this simple.
		// If you want exact 404 handling, inspect smithy.APIError.
		if strings.Contains(err.Error(), "NotFound") ||
			strings.Contains(err.Error(), "404") ||
			strings.Contains(err.Error(), "NoSuchKey") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *S3Storage) GetURL(ctx context.Context, key string) (string, error) {
	key = strings.TrimLeft(key, "/")

	if s.publicURL != "" {
		return s.publicURL + "/" + key, nil
	}

	if s.presign {
		result, err := s.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = s.presignTTL
		})

		if err != nil {
			return "", err
		}

		return result.URL, nil
	}

	return key, nil
}
