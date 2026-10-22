package storage

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Service defines standard object storage operations
type Service interface {
	// Upload saves an object to storage and returns its relative path or public URL
	Upload(ctx context.Context, folder, filename string, reader io.Reader, size int64, contentType string) (string, error)
	// Delete removes an object from storage
	Delete(ctx context.Context, key string) error
	// GetPublicURL returns the full accessible URL for an object key
	GetPublicURL(key string) string
	// Driver returns the storage driver name ("local" or "r2")
	Driver() string
}

// Config holds configuration parameters for Storage Service
type Config struct {
	Driver       string // "local" or "r2" / "s3"
	UploadPath   string // local upload path (default "./uploads")
	Bucket       string // R2 / S3 bucket name
	Endpoint     string // R2 / S3 endpoint (e.g. "1942e38345273d01fb206624b19ddd52.r2.cloudflarestorage.com")
	AccessKeyID  string // R2 / S3 access key ID
	SecretKey    string // R2 / S3 secret access key
	PublicURL    string // Public CDN base URL (e.g. "https://pub-b2f72286c4da4fe1b67e2d892478a459.r2.dev")
	UseSSL       bool
}

// LocalStorage implements Service using local disk
type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	if basePath == "" {
		basePath = "./uploads"
	}
	_ = os.MkdirAll(basePath, 0755)
	return &LocalStorage{basePath: basePath}
}

func (l *LocalStorage) Driver() string {
	return "local"
}

func (l *LocalStorage) Upload(ctx context.Context, folder, filename string, reader io.Reader, size int64, contentType string) (string, error) {
	dir := filepath.Join(l.basePath, folder)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	fullPath := filepath.Join(dir, filename)
	out, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, reader); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("/uploads/%s/%s", folder, filename), nil
}

func (l *LocalStorage) Delete(ctx context.Context, key string) error {
	relPath := strings.TrimPrefix(key, "/uploads/")
	relPath = strings.TrimPrefix(relPath, "/")
	fullPath := filepath.Join(l.basePath, relPath)
	if _, err := os.Stat(fullPath); err == nil {
		return os.Remove(fullPath)
	}
	return nil
}

func (l *LocalStorage) GetPublicURL(key string) string {
	if strings.HasPrefix(key, "http://") || strings.HasPrefix(key, "https://") {
		return key
	}
	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}
	return key
}

// R2Storage implements Service using Cloudflare R2 (S3 compatible)
type R2Storage struct {
	client    *minio.Client
	bucket    string
	publicURL string
}

func NewR2Storage(cfg Config) (*R2Storage, error) {
	endpoint := cfg.Endpoint
	// Strip scheme if provided (minio client expects host or host:port)
	if strings.HasPrefix(endpoint, "https://") {
		endpoint = strings.TrimPrefix(endpoint, "https://")
	} else if strings.HasPrefix(endpoint, "http://") {
		endpoint = strings.TrimPrefix(endpoint, "http://")
	}
	// Strip any trailing slash or bucket suffix if passed
	endpoint = strings.TrimSuffix(endpoint, "/")
	if idx := strings.Index(endpoint, "/"); idx != -1 {
		endpoint = endpoint[:idx]
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init R2 client: %w", err)
	}

	publicURL := strings.TrimRight(cfg.PublicURL, "/")

	return &R2Storage{
		client:    client,
		bucket:    cfg.Bucket,
		publicURL: publicURL,
	}, nil
}

func (r *R2Storage) Driver() string {
	return "r2"
}

func (r *R2Storage) Upload(ctx context.Context, folder, filename string, reader io.Reader, size int64, contentType string) (string, error) {
	if contentType == "" {
		ext := filepath.Ext(filename)
		contentType = mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}
	}

	key := fmt.Sprintf("%s/%s", strings.Trim(folder, "/"), strings.Trim(filename, "/"))

	_, err := r.client.PutObject(ctx, r.bucket, key, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object to R2: %w", err)
	}

	if r.publicURL != "" {
		return fmt.Sprintf("%s/%s", r.publicURL, key), nil
	}
	return fmt.Sprintf("/%s", key), nil
}

func (r *R2Storage) Delete(ctx context.Context, key string) error {
	// Strip publicURL prefix if full URL was given
	if r.publicURL != "" && strings.HasPrefix(key, r.publicURL) {
		key = strings.TrimPrefix(key, r.publicURL)
	}
	if parsed, err := url.Parse(key); err == nil && parsed.Path != "" {
		key = parsed.Path
	}
	key = strings.TrimPrefix(key, "/")

	return r.client.RemoveObject(ctx, r.bucket, key, minio.RemoveObjectOptions{})
}

func (r *R2Storage) GetPublicURL(key string) string {
	if strings.HasPrefix(key, "http://") || strings.HasPrefix(key, "https://") {
		return key
	}
	key = strings.TrimPrefix(key, "/")
	if r.publicURL != "" {
		return fmt.Sprintf("%s/%s", r.publicURL, key)
	}
	return "/" + key
}

// NewService creates a storage service based on configuration
func NewService(cfg Config) (Service, error) {
	if strings.ToLower(cfg.Driver) == "r2" || strings.ToLower(cfg.Driver) == "s3" {
		if cfg.AccessKeyID != "" && cfg.SecretKey != "" && cfg.Endpoint != "" {
			return NewR2Storage(cfg)
		}
	}
	return NewLocalStorage(cfg.UploadPath), nil
}
