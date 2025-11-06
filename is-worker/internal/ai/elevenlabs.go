package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ElevenLabsService provides text-to-speech functionality using ElevenLabs API
type ElevenLabsService struct {
	apiKey     string
	httpClient *http.Client
}

// NewElevenLabsService creates a new ElevenLabs service
func NewElevenLabsService() (*ElevenLabsService, error) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ELEVENLABS_API_KEY environment variable not set")
	}

	return &ElevenLabsService{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}, nil
}

// TextToSpeechRequest represents the request body for ElevenLabs text-to-speech API
type TextToSpeechRequest struct {
	Text          string                 `json:"text"`
	ModelID       string                 `json:"model_id"`
	VoiceSettings map[string]interface{} `json:"voice_settings,omitempty"`
}

// GenerateAudio generates audio from text using the specified voice ID
// Returns the audio data as a byte slice
func (s *ElevenLabsService) GenerateAudio(ctx context.Context, text string, voiceID string) ([]byte, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}
	if voiceID == "" {
		return nil, fmt.Errorf("voiceID cannot be empty")
	}

	// Prepare request body
	requestBody := TextToSpeechRequest{
		Text:    text,
		ModelID: "eleven_multilingual_v2", // or eleven_v3 or eleven_ttv_v3
		VoiceSettings: map[string]interface{}{
			"stability":         0.5,
			"similarity_boost":  0.75,
			"style":             0.0,
			"use_speaker_boost": true,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voiceID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", s.apiKey)

	// Make request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ElevenLabs API returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
