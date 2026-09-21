package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	Endpoint  string `envconfig:"S3_ENDPOINT" required:"true"`
	Bucket    string `envconfig:"S3_BUCKET" required:"true"`
	Region    string `envconfig:"S3_REGION" required:"true"`
	AccessKey string `envconfig:"S3_ACCESS_KEY" required:"true"`
	SecretKey string `envconfig:"S3_SECRET_KEY" required:"true"`
}

type S3 struct {
	cfg           S3Config
	client        *s3.Client
	presignClient *s3.PresignClient
}

func NewS3(cfg S3Config) (*S3, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// create s3 client with custom endpoint and path style
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = true
	})

	return &S3{
		cfg:           cfg,
		client:        s3Client,
		presignClient: s3.NewPresignClient(s3Client),
	}, nil
}

func (s *S3) PresignPut(ctx context.Context, key string, contentType string, ttl time.Duration) (string, error) {
	req, err := s.presignClient.PresignPutObject(ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(s.cfg.Bucket),
			Key:         aws.String(key),
			ContentType: aws.String(contentType),
		},
		func(o *s3.PresignOptions) {
			o.Expires = ttl
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to presign put object: %w", err)
	}

	return req.URL, nil
}

func (s *S3) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

func (s *S3) PublicURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", s.cfg.Endpoint, s.cfg.Bucket, key)
}

func (s *S3) KeyFromURL(url string) string {
	prefix := fmt.Sprintf("%s/%s/", s.cfg.Endpoint, s.cfg.Bucket)
	return strings.TrimPrefix(url, prefix)
}
