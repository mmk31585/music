package storage

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"go.uber.org/zap"
)

type S3Config struct {
	Bucket          string
	Region          string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	RoleARN         string
	AccountID       string
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

// s3Step labels each S3 operation for structured logs.
const (
	s3StepPut     = "s3_put_object"
	s3StepHead    = "s3_head_object"
	s3StepPresign = "s3_presign_put_object"
	s3StepDelete  = "s3_delete_object"
	s3StepGet     = "s3_get_object"
	s3StepCopy    = "s3_copy_object"
)

func NewS3Storage(ctx context.Context, cfg S3Config) (*S3Storage, error) {
	if cfg.Bucket == "" {
		return nil, errors.New("s3 bucket is required")
	}

	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	// Disable EC2 IMDS credential resolution so a missing credential fails fast
	// instead of hanging on the IMDS endpoint. Our OIDC provider (or static keys)
	// fully replaces the default chain.
	_ = os.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	defer os.Unsetenv("AWS_EC2_METADATA_DISABLED")

	loadOptions := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.Region),
	}

	// Precedence: static credentials (S3_ACCESS_KEY_ID/S3_SECRET_ACCESS_KEY) >
	// OIDC web identity via Vercel (AWS_ROLE_ARN) > SDK default chain.
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

	// OIDC web identity: assume role using a per-request OIDC token.
	// The token is supplied via the x-vercel-oidc-token header (Gin middleware
	// stores it in the request context) or the VERCEL_OIDC_TOKEN env var.
	if cfg.RoleARN != "" && (cfg.AccessKeyID == "" || cfg.SecretAccessKey == "") {
		stsClient := sts.NewFromConfig(awsCfg)
		awsCfg.Credentials = &oidcCredentialProvider{
			roleARN:   cfg.RoleARN,
			stsClient: stsClient,
		}
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}

		o.UsePathStyle = cfg.UsePathStyle

		// Explicit HTTP client timeouts so a missing credential fails fast.
		o.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
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
		Bucket:   aws.String(s.bucket),
		Key:      aws.String(key),
		Body:     reader,
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}
	if size >= 0 {
		// S3 requires Content-Length for stream bodies; the caller (media
		// service) passes header.Size, which is known at upload time.
		input.ContentLength = aws.Int64(size)
	}

	log := zap.L().With(
		zap.String("step", s3StepPut),
		zap.String("bucket", s.bucket),
		zap.String("key", key),
	)

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		log.Error("s3 put object failed", zap.Error(err))
	}
	return err
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	key = strings.TrimLeft(key, "/")

	log := zap.L().With(
		zap.String("step", s3StepDelete),
		zap.String("bucket", s.bucket),
		zap.String("key", key),
	)

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Error("s3 delete object failed", zap.Error(err))
	}
	return err
}

func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	key = strings.TrimLeft(key, "/")

	log := zap.L().With(
		zap.String("step", s3StepHead),
		zap.String("bucket", s.bucket),
		zap.String("key", key),
	)

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

		log.Error("s3 head object failed", zap.Error(err))
		return false, err
	}

	return true, nil
}

func (s *S3Storage) Copy(ctx context.Context, srcKey, dstKey string) error {
	srcKey = strings.TrimLeft(srcKey, "/")
	dstKey = strings.TrimLeft(dstKey, "/")

	_, err := s.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s.bucket),
		CopySource: aws.String(s.bucket + "/" + srcKey),
		Key:        aws.String(dstKey),
	})
	return err
}

func (s *S3Storage) Get(ctx context.Context, key string) ([]byte, error) {
	key = strings.TrimLeft(key, "/")

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
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

// PresignPut returns a time-limited presigned PUT URL for direct client-to-S3 upload.
func (s *S3Storage) PresignPut(ctx context.Context, key string, contentType string, ttl time.Duration) (string, error) {
	key = strings.TrimLeft(key, "/")

	log := zap.L().With(
		zap.String("step", s3StepPresign),
		zap.String("bucket", s.bucket),
		zap.String("key", key),
	)

	result, err := s.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = ttl
	})

	if err != nil {
		log.Error("s3 presign put object failed", zap.Error(err))
		return "", err
	}

	return result.URL, nil
}
