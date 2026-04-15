package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client    *s3.Client
	presigner *s3.PresignClient
	bucket    string
	region    string
	publicURL string // Wichtig für GetUrl
}

type S3Config struct {
	Bucket    string
	Region    string
	Endpoint  string // Intern: http://minio:9000
	PublicURL string // Extern: https://travel-nugget.duckdns.org/minio
	AccessKey string
	SecretKey string
}

func NewS3Storage(cfg S3Config) (*S3Storage, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket name is required")
	}
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	var opts []func(*config.LoadOptions) error
	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey, cfg.SecretKey, "",
		)))
	}
	opts = append(opts, config.WithRegion(cfg.Region))

	awsCfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// 1. Interner Client (für Server-Side Upload/Read)
	internalClient := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = true
	})

	// 2. Presign Client (für Browser-Upload-URLs)
	// Nutzt die PublicURL, damit die Signatur für die Domain + /minio Pfad passt
	presignClient := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.PublicURL)
		o.UsePathStyle = true
	})

	return &S3Storage{
		client:    internalClient,
		presigner: s3.NewPresignClient(presignClient),
		bucket:    cfg.Bucket,
		region:    cfg.Region,
		publicURL: cfg.PublicURL,
	}, nil
}

func (s *S3Storage) Upload(ctx context.Context, fileName string, file io.Reader) (string, error) {
	log.Printf("Uploading file %s to internal endpoint", fileName)

	// manager.Uploader ist der Standard für S3 v2
	uploader := transfermanager.New(s.client)

	_, err := uploader.UploadObject(ctx, &transfermanager.UploadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return s.GetUrl(ctx, fileName)
}

func (s *S3Storage) ReadFile(ctx context.Context, fileName string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return result.Body, nil
}

func (s *S3Storage) GetUrl(ctx context.Context, fileName string) (string, error) {
	// Resultat: https://travel-nugget.duckdns.org/minio/trip-manager/avatars/file.png
	return fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucket, fileName), nil
}

func (s *S3Storage) GeneratePresignedURL(ctx context.Context, fileName string, expiresIn time.Duration) (string, error) {
	req := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
	}

	result, err := s.presigner.PresignPutObject(ctx, req, s3.WithPresignExpires(expiresIn))
	if err != nil {
		return "", err
	}

	log.Printf("Presigned URL ready for browser: %s", result.URL)
	return result.URL, nil
}
