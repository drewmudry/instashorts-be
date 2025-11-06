package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
)

// GCSClient provides file upload functionality to Google Cloud Storage
type GCSClient struct {
	client *storage.Client
	bucket string
}

// NewGCSClient creates a new GCS service
func NewGCSClient(ctx context.Context) (*GCSClient, error) {
	// Auth is automatic!
	// The client will find the GOOGLE_APPLICATION_CREDENTIALS
	// environment variable set by your docker-compose.yml.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		return nil, fmt.Errorf("GCS_BUCKET_NAME environment variable not set")
	}

	return &GCSClient{
		client: client,
		bucket: bucketName,
	}, nil
}

// getSignedURL generates a V4 signed URL for a GCS object.
// This is the secure way to grant temporary access.
func (c *GCSClient) getSignedURL(objectName string) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute), // 15 minutes of access
	}

	url, err := c.client.Bucket(c.bucket).SignedURL(objectName, opts)
	if err != nil {
		return "", fmt.Errorf("Bucket.SignedURL: %w", err)
	}
	return url, nil
}

// upload private a helper to upload data
func (c *GCSClient) upload(ctx context.Context, data []byte, key string, contentType string) error {
	if len(data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	// Get a handle to the GCS object
	obj := c.client.Bucket(c.bucket).Object(key)

	// Create a writer
	wc := obj.NewWriter(ctx)
	wc.ContentType = contentType

	// Copy the data into the writer
	if _, err := io.Copy(wc, bytes.NewReader(data)); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	// Close the writer to finalize the upload
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	return nil
}

// UploadAudio uploads an audio file to GCS and returns a signed URL
func (c *GCSClient) UploadAudio(ctx context.Context, data []byte, videoID int) (string, error) {
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("audio/%d/%d.mp3", videoID, timestamp)

	if err := c.upload(ctx, data, key, "audio/mpeg"); err != nil {
		return "", fmt.Errorf("failed to upload audio to GCS: %w", err)
	}

	return c.getSignedURL(key)
}

// UploadImage uploads an image file to GCS and returns a signed URL
func (c *GCSClient) UploadImage(ctx context.Context, data []byte, videoID int, sceneIndex int) (string, error) {
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("images/%d/scene_%d_%d.png", videoID, sceneIndex, timestamp)

	if err := c.upload(ctx, data, key, "image/png"); err != nil {
		return "", fmt.Errorf("failed to upload image to GCS: %w", err)
	}

	return c.getSignedURL(key)
}

// UploadVideo uploads a video file to GCS and returns a signed URL
func (c *GCSClient) UploadVideo(ctx context.Context, data []byte, videoID int) (string, error) {
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("videos/%d/%d.mp4", videoID, timestamp)

	if err := c.upload(ctx, data, key, "video/mp4"); err != nil {
		return "", fmt.Errorf("failed to upload video to GCS: %w", err)
	}

	return c.getSignedURL(key)
}

// UploadFile is a generic method to upload any file to GCS
func (c *GCSClient) UploadFile(ctx context.Context, data []byte, path string, contentType string) (string, error) {
	key := filepath.Clean(path)

	if err := c.upload(ctx, data, key, contentType); err != nil {
		return "", fmt.Errorf("failed to upload file to GCS: %w", err)
	}

	return c.getSignedURL(key)
}
