package provider

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/option"
)

type GCStorage struct {
	client    *storage.Client
	iamClient *iamcredentials.Service
	bucket    string
	signerSA  string
	ttl       time.Duration
}

// Exists checks if the file exists in the bucket.
func (s *GCStorage) Exists(ctx context.Context, fileName string) (bool, error) {
	_, err := s.client.Bucket(s.bucket).Object(fileName).Attrs(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("gcs attrs %q: %w", fileName, err)
	}
	return true, nil
}

type GCSConfig struct {
	Bucket   string        // Bucket-Name, z.B. "trip-manager-uploads"
	SignerSA string        // Service-Account-Email für Impersonation-Signing
	TTL      time.Duration // Lifetime für Upload/Download URLs, default 15min
}

func NewGCStorage(ctx context.Context, cfg GCSConfig) (*GCStorage, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket name is required")
	}
	if cfg.SignerSA == "" {
		return nil, fmt.Errorf("signer service account is required")
	}
	if cfg.TTL == 0 {
		cfg.TTL = 15 * time.Minute
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage client: %w", err)
	}

	iamClient, err := iamcredentials.NewService(ctx,
		option.WithScopes(iamcredentials.CloudPlatformScope))
	if err != nil {
		return nil, errors.Join(fmt.Errorf("iamcredentials client: %w", err), client.Close())
	}

	return &GCStorage{
		client:    client,
		iamClient: iamClient,
		bucket:    cfg.Bucket,
		signerSA:  cfg.SignerSA,
		ttl:       cfg.TTL,
	}, nil
}

// GetUploadURL return short-lived upload url.
func (s *GCStorage) GetUploadURL(ctx context.Context, fileName string) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "PUT",
		GoogleAccessID: s.signerSA,
		Expires:        time.Now().Add(s.ttl),
		SignBytes:      s.signBytes(ctx),
	}

	url, err := storage.SignedURL(s.bucket, fileName, opts)
	if err != nil {
		return "", fmt.Errorf("signed put url %q: %w", fileName, err)
	}
	return url, nil
}

// GetDownloadURL return short-lived download url.
func (s *GCStorage) GetDownloadURL(ctx context.Context, fileName string) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: s.signerSA,
		Expires:        time.Now().Add(s.ttl),
		SignBytes:      s.signBytes(ctx),
	}

	url, err := storage.SignedURL(s.bucket, fileName, opts)
	if err != nil {
		return "", fmt.Errorf("signed get url %q: %w", fileName, err)
	}
	return url, nil
}

// ReadFile read file from backend
func (s *GCStorage) ReadFile(ctx context.Context, fileName string) (io.ReadCloser, error) {
	r, err := s.client.Bucket(s.bucket).Object(fileName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("gcs read %q: %w", fileName, err)
	}
	return r, nil
}

// WriteFile write a file from backend
func (s *GCStorage) WriteFile(ctx context.Context, fileName string, file io.Reader) error {
	obj := s.client.Bucket(s.bucket).Object(fileName)
	w := obj.NewWriter(ctx)

	if _, err := io.Copy(w, file); err != nil {
		_ = w.Close()
		return fmt.Errorf("gcs write copy %q: %w", fileName, err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("gcs write close %q: %w", fileName, err)
	}
	return nil
}

// Delete removes object
func (s *GCStorage) Delete(ctx context.Context, fileName string) error {
	err := s.client.Bucket(s.bucket).Object(fileName).Delete(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil
		}
		return fmt.Errorf("gcs delete %q: %w", fileName, err)
	}
	return nil
}

// signBytes creates SignBytes-Callback for Signed-URL-Gen via IAM Credentials API.
// (requires roles/iam.serviceAccountTokenCreator on Signer-SA).
func (s *GCStorage) signBytes(ctx context.Context) func([]byte) ([]byte, error) {
	name := fmt.Sprintf("projects/-/serviceAccounts/%s", s.signerSA)
	return func(payload []byte) ([]byte, error) {
		resp, err := s.iamClient.Projects.ServiceAccounts.SignBlob(
			name,
			&iamcredentials.SignBlobRequest{
				Payload: base64.StdEncoding.EncodeToString(payload),
			},
		).Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("signBlob: %w", err)
		}
		return base64.StdEncoding.DecodeString(resp.SignedBlob)
	}
}
