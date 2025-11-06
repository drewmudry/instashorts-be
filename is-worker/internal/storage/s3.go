package storage

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Service provides file upload functionality to AWS S3
type S3Service struct {
	client     *s3.Client
	bucketName string
	region     string
}

// NewS3Service creates a new S3 service
func NewS3Service(ctx context.Context) (*S3Service, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	bucketName := os.Getenv("AWS_S3_BUCKET")

	if accessKey == "" {
		return nil, fmt.Errorf("AWS_ACCESS_KEY_ID environment variable not set")
	}
	if secretKey == "" {
		return nil, fmt.Errorf("AWS_SECRET_ACCESS_KEY environment variable not set")
	}
	if region == "" {
		region = "us-east-1" // Default region
	}
	if bucketName == "" {
		bucketName = "instashorts-content" // Default bucket name
	}

	// Create credentials
	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	return &S3Service{
		client:     client,
		bucketName: bucketName,
		region:     region,
	}, nil
}

// UploadAudio uploads an audio file to S3 and returns the public URL
func (s *S3Service) UploadAudio(ctx context.Context, data []byte, videoID int) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("data cannot be empty")
	}

	// Generate a unique file name with timestamp
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("audio/%d/%d.mp3", videoID, timestamp)

	// Upload to S3
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("audio/mpeg"),
		ACL:         "public-read", // Make file publicly accessible
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Construct the public URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, key)

	return url, nil
}

// UploadImage uploads an image file to S3 and returns the public URL
func (s *S3Service) UploadImage(ctx context.Context, data []byte, videoID int, sceneIndex int) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("data cannot be empty")
	}

	// Generate a unique file name with timestamp
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("images/%d/scene_%d_%d.png", videoID, sceneIndex, timestamp)

	// Upload to S3
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("image/png"),
		ACL:         "public-read", // Make file publicly accessible
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Construct the public URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, key)

	return url, nil
}

// UploadVideo uploads a video file to S3 and returns the public URL
func (s *S3Service) UploadVideo(ctx context.Context, data []byte, videoID int) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("data cannot be empty")
	}

	// Generate a unique file name with timestamp
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("videos/%d/%d.mp4", videoID, timestamp)

	// Upload to S3
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("video/mp4"),
		ACL:         "public-read", // Make file publicly accessible
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Construct the public URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, key)

	return url, nil
}

// UploadFile is a generic method to upload any file to S3
func (s *S3Service) UploadFile(ctx context.Context, data []byte, path string, contentType string) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("data cannot be empty")
	}

	// Clean the path
	path = filepath.Clean(path)

	// Upload to S3
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(path),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
		ACL:         "public-read",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Construct the public URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, path)

	return url, nil
}

