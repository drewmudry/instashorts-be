package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

type Service struct {
	client *genai.Client
}

// NewService creates a new AI service with Vertex AI
func NewService(ctx context.Context) (*Service, error) {
	// Get GCP project and location from environment
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("GCP_PROJECT_ID environment variable not set")
	}

	location := os.Getenv("GCP_LOCATION")
	if location == "" {
		location = "us-central1" // Default location
	}

	// Create client using Vertex AI with Application Default Credentials
	// This will use the credentials from `gcloud auth application-default login`
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  projectID,
		Location: location,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Vertex AI client: %w", err)
	}

	return &Service{client: client}, nil
}

// GenerateVideoScript generates a video script based on the theme
// The script will be 250-300 words for approximately 70 seconds of narration
func (s *Service) GenerateVideoScript(ctx context.Context, theme string) (string, error) {
	prompt := fmt.Sprintf(`Generate a compelling and engaging video script about: %s

Requirements:
- The script should be between 50-100 words (approximately 10 seconds when narrated)
- Write in a conversational and engaging tone suitable for short-form video content
- The script should be in paragraph format with no headers or subheaders
- Make it informative yet entertaining
- Include a strong hook at the beginning to capture attention
- Structure the content with a clear beginning, middle, and end
- End with a thought-provoking conclusion
- Write ONLY the script text, no additional formatting, labels, or stage directions


Generate the script now:`, theme)
	// - Include voice controls like: [laughs], [laughs harder], [starts laughing], [wheezing], [whispers], [sighs], [exhales],[sarcastic], [curious], [excited], [crying], [snorts], [mischievously]

	result, err := s.client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-exp",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	script := result.Text()
	if script == "" {
		return "", fmt.Errorf("generated script is empty")
	}

	return script, nil
}

// ScenePrompt represents a scene with its image generation prompt
type ScenePrompt struct {
	ImagePrompt string `json:"image_prompt"`
	Index       int    `json:"index"`
}

// GenerateScenes generates 2-3 scene prompts based on the video script
func (s *Service) GenerateScenes(ctx context.Context, script string) ([]ScenePrompt, error) {
	prompt := fmt.Sprintf(`Based on the following video script, generate 2-3 scene descriptions that will be used to create images for the video.

Script:
%s

Requirements:
- Generate between 2-3 scenes that flow with the narration. 
- the script should be in paragraph format with no headers or subheaders.
- Each scene should be a detailed, descriptive image prompt that can be used for image generation. include consistent style and coloring across all scenes
- Use consistent styling and coloring across all scenes (e.g., "cinematic style", "vibrant colors", "minimalist illustration")
- Make each prompt very descriptive (3-4 sentences) to generate high-quality images
- Prompts should include detailed descriptions of the scene, including the characters, objects, and background.
- Order the scenes to match the progression of the script
- Return ONLY a valid JSON array with no additional text, markdown formatting, or code blocks
- Format: [{"image_prompt": "detailed description here", "index": 0}, {"image_prompt": "another detailed description", "index": 1}, ...]

Generate the JSON array now:`, script)

	result, err := s.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-pro",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate scenes: %w", err)
	}

	responseText := result.Text()
	if responseText == "" {
		return nil, fmt.Errorf("generated scenes response is empty")
	}

	// Strip markdown code blocks if present (```json ... ```)
	responseText = stripMarkdownCodeBlocks(responseText)

	// Parse the JSON response
	var scenes []ScenePrompt
	if err := json.Unmarshal([]byte(responseText), &scenes); err != nil {
		return nil, fmt.Errorf("failed to parse scenes JSON: %w (response: %s)", err, responseText)
	}

	if len(scenes) == 0 {
		return nil, fmt.Errorf("no scenes generated")
	}

	return scenes, nil
}

// GenerateImage generates an image using the Imagen model based on a text prompt
// Returns the image data as bytes
func (s *Service) GenerateImage(ctx context.Context, prompt string) ([]byte, error) {
	// Generate image using Imagen 4.0
	config := &genai.GenerateImagesConfig{
		NumberOfImages: 1,
		AspectRatio:    "9:16", // Vertical format for short-form videos
	}

	response, err := s.client.Models.GenerateImages(
		ctx,
		"imagen-4.0-generate-001",
		prompt,
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate image: %w", err)
	}

	if len(response.GeneratedImages) == 0 {
		return nil, fmt.Errorf("no images generated")
	}

	// Return the first (and only) generated image
	imageBytes := response.GeneratedImages[0].Image.ImageBytes
	if len(imageBytes) == 0 {
		return nil, fmt.Errorf("generated image has no data")
	}

	return imageBytes, nil
}

// stripMarkdownCodeBlocks removes markdown code block formatting from a string
// Handles cases like ```json ... ``` or ``` ... ```
func stripMarkdownCodeBlocks(text string) string {
	// Trim whitespace
	text = strings.TrimSpace(text)

	// Check if it starts with ``` and ends with ```
	if strings.HasPrefix(text, "```") && strings.HasSuffix(text, "```") {
		// Remove the opening ```
		text = strings.TrimPrefix(text, "```")

		// Remove language identifier if present (e.g., "json")
		if idx := strings.Index(text, "\n"); idx != -1 {
			text = text[idx+1:]
		}

		// Remove the closing ```
		text = strings.TrimSuffix(text, "```")

		// Trim any remaining whitespace
		text = strings.TrimSpace(text)
	}

	return text
}
