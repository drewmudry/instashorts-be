package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"google.golang.org/api/option"
)

// CaptionWord represents a word with timing information
type CaptionWord struct {
	Word      string  `json:"word"`
	StartTime float64 `json:"start_time"` // in seconds
	EndTime   float64 `json:"end_time"`   // in seconds
}

// SpeechToTextService handles speech-to-text transcription with word-level timestamps
type SpeechToTextService struct {
	client *speech.Client
}

// NewSpeechToTextService creates a new Speech-to-Text service
func NewSpeechToTextService(ctx context.Context) (*SpeechToTextService, error) {
	// Check if we have GCP credentials
	credsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsPath != "" {
		// Use explicit credentials if provided
		client, err := speech.NewClient(ctx, option.WithCredentialsFile(credsPath))
		if err != nil {
			return nil, fmt.Errorf("failed to create speech client with credentials: %w", err)
		}
		return &SpeechToTextService{client: client}, nil
	}

	// Otherwise use Application Default Credentials
	client, err := speech.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create speech client: %w", err)
	}

	return &SpeechToTextService{client: client}, nil
}

// Close closes the Speech-to-Text client
func (s *SpeechToTextService) Close() error {
	return s.client.Close()
}

// GenerateCaptionsFromURL generates word-level captions from an audio file URL
// The audio file should be accessible via HTTP/HTTPS (e.g., from S3)
func (s *SpeechToTextService) GenerateCaptionsFromURL(ctx context.Context, audioURL string) ([]CaptionWord, error) {
	// Download the audio file
	audioData, err := downloadAudioFile(ctx, audioURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download audio file: %w", err)
	}

	return s.GenerateCaptions(ctx, audioData)
}

// GenerateCaptions generates word-level captions from audio data
func (s *SpeechToTextService) GenerateCaptions(ctx context.Context, audioData []byte) ([]CaptionWord, error) {
	// Create recognition request
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:                   speechpb.RecognitionConfig_MP3,
			SampleRateHertz:            0, // Auto-detect
			LanguageCode:               "en-US",
			EnableWordTimeOffsets:      true, // This is key for word-level timestamps
			EnableAutomaticPunctuation: true,
			Model:                      "latest_long", // Best for longer audio
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{
				Content: audioData,
			},
		},
	}

	// Perform recognition
	resp, err := s.client.Recognize(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to recognize speech: %w", err)
	}

	// Extract word-level timestamps
	var captions []CaptionWord
	for _, result := range resp.Results {
		// Use the first (most confident) alternative
		if len(result.Alternatives) == 0 {
			continue
		}

		alternative := result.Alternatives[0]
		for _, wordInfo := range alternative.Words {
			// Convert from protobuf Duration to float64 seconds
			startTime := float64(wordInfo.StartTime.Seconds) + float64(wordInfo.StartTime.Nanos)/1e9
			endTime := float64(wordInfo.EndTime.Seconds) + float64(wordInfo.EndTime.Nanos)/1e9

			captions = append(captions, CaptionWord{
				Word:      wordInfo.Word,
				StartTime: startTime,
				EndTime:   endTime,
			})
		}
	}

	if len(captions) == 0 {
		return nil, fmt.Errorf("no captions generated from audio")
	}

	return captions, nil
}

// CaptionsToJSON converts captions to JSON string
func CaptionsToJSON(captions []CaptionWord) (string, error) {
	jsonData, err := json.Marshal(captions)
	if err != nil {
		return "", fmt.Errorf("failed to marshal captions to JSON: %w", err)
	}
	return string(jsonData), nil
}

// downloadAudioFile downloads an audio file from a URL
func downloadAudioFile(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}
