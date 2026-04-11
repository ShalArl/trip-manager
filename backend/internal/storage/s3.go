package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Storage implements Storage interface using AWS S3 or S3-compatible services (like MinIO)
type S3Storage struct {
	client    *s3.Client
	bucket    string
	region    string
	publicURL string // Base URL for public file access
}

// S3Config contains configuration for S3 storage
type S3Config struct {
	Bucket    string // S3 bucket name
	Region    string // AWS region (e.g., "us-east-1")
	Endpoint  string // S3 endpoint (optional, for MinIO or custom S3 services)
	PublicURL string // Base URL for public file access (e.g., "https://s3.example.com" or "http://localhost:9000")
	AccessKey string // AWS access key or MinIO access key
	SecretKey string // AWS secret key or MinIO secret key
	UseSSL    bool   // Use SSL for S3 connection
}

// NewS3Storage creates a new S3Storage instance
// Works with AWS S3, MinIO, and other S3-compatible services
func NewS3Storage(cfg S3Config) (*S3Storage, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket name is required")
	}

	if cfg.Region == "" {
		cfg.Region = "us-east-1" // Default region
	}

	if cfg.PublicURL == "" {
		cfg.PublicURL = fmt.Sprintf("https://s3.%s.amazonaws.com/%s", cfg.Region, cfg.Bucket)
	}

	// Build AWS config
	var opts []func(*config.LoadOptions) error

	// Add custom credentials if provided
	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"", // No session token
		)))
	}

	// Add region
	opts = append(opts, config.WithRegion(cfg.Region))

	awsCfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with optional custom endpoint (for MinIO)
	var client *s3.Client
	if cfg.Endpoint != "" {
		// For MinIO or custom S3 services
		client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true // Required for MinIO
		})
	} else {
		// For AWS S3
		client = s3.NewFromConfig(awsCfg)
	}

	return &S3Storage{
		client:    client,
		bucket:    cfg.Bucket,
		region:    cfg.Region,
		publicURL: cfg.PublicURL,
	}, nil
}

// Upload saves a file to S3 and returns the public URL
func (s *S3Storage) Upload(ctx context.Context, fileName string, file io.Reader) (string, error) {
	uploader := manager.NewUploader(s.client)

	// Upload file with public-read ACL (or private, depending on use case)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    types.ObjectCannedACLPublicRead, // Make file publicly readable
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Return public URL
	fileURL := fmt.Sprintf("%s/%s", s.publicURL, fileName)
	return fileURL, nil
}

// ReadFile retrieves a file from S3
func (s *S3Storage) ReadFile(ctx context.Context, fileName string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read file from S3: %w", err)
	}

	return result.Body, nil
}

// GetUrl returns the public URL for a file
func (s *S3Storage) GetUrl(ctx context.Context, fileName string) (string, error) {
	fileURL := fmt.Sprintf("%s/%s", s.publicURL, fileName)
	return fileURL, nil
}


