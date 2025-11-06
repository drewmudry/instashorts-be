package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"instashorts-be/is-worker/internal/ai"
	"instashorts-be/is-worker/internal/ai/gemini"
	"instashorts-be/is-worker/internal/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// Task types
const (
	TypeProcessVideo        = "video:process"
	TypeSendEmail           = "email:send"
	TypeGenerateVideoScript = "video:generate_script"
	TypeGenerateAudio       = "video:generate_audio"
	TypeGenerateCaptions    = "video:generate_captions"
	TypeGenerateScenes      = "video:generate_scenes"
	TypeGenerateSceneImage  = "video:generate_scene_image"
	TypeRenderVideo         = "video:render"
	TypeVideoComplete       = "video:complete"
	// Add more task types as needed
)

// ProcessVideoPayload represents the payload for video processing tasks
type ProcessVideoPayload struct {
	VideoID string `json:"video_id"`
	UserID  string `json:"user_id"`
}

// GenerateVideoScriptPayload represents the payload for video script generation tasks
type GenerateVideoScriptPayload struct {
	VideoID int `json:"video_id"`
}

// GenerateAudioPayload represents the payload for audio generation tasks
type GenerateAudioPayload struct {
	VideoID int `json:"video_id"`
}

// GenerateCaptionsPayload represents the payload for caption generation tasks
type GenerateCaptionsPayload struct {
	VideoID int `json:"video_id"`
}

// GenerateScenesPayload represents the payload for scene generation tasks
type GenerateScenesPayload struct {
	VideoID int `json:"video_id"`
}

// GenerateSceneImagePayload represents the payload for scene image generation tasks
type GenerateSceneImagePayload struct {
	SceneID int `json:"scene_id"`
}

// RenderVideoPayload represents the payload for video rendering tasks
type RenderVideoPayload struct {
	VideoID int `json:"video_id"`
}

// VideoCompletePayload represents the payload for video completion tasks
type VideoCompletePayload struct {
	VideoID  int    `json:"video_id"`
	VideoURL string `json:"video_url"`
}

// SendEmailPayload represents the payload for email tasks
type SendEmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// EnqueueProcessVideo enqueues a video processing task
func (c *Client) EnqueueProcessVideo(payload ProcessVideoPayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeProcessVideo, jsonPayload)
	info, err := c.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// EnqueueSendEmail enqueues an email sending task
func (c *Client) EnqueueSendEmail(payload SendEmailPayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeSendEmail, jsonPayload)
	info, err := c.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// EnqueueGenerateVideoScript enqueues a video script generation task
func (c *Client) EnqueueGenerateVideoScript(payload GenerateVideoScriptPayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeGenerateVideoScript, jsonPayload)
	info, err := c.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// EnqueueGenerateAudio enqueues an audio generation task
func (c *Client) EnqueueGenerateAudio(payload GenerateAudioPayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeGenerateAudio, jsonPayload)
	info, err := c.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// EnqueueGenerateCaptions enqueues a caption generation task
func (c *Client) EnqueueGenerateCaptions(payload GenerateCaptionsPayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeGenerateCaptions, jsonPayload)
	info, err := c.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// EnqueueGenerateScenes enqueues a scene generation task
func (c *Client) EnqueueGenerateScenes(payload GenerateScenesPayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeGenerateScenes, jsonPayload)
	info, err := c.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// EnqueueGenerateSceneImage enqueues a scene image generation task
func (c *Client) EnqueueGenerateSceneImage(payload GenerateSceneImagePayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeGenerateSceneImage, jsonPayload)
	info, err := c.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// EnqueueRenderVideo enqueues a video rendering task
func (c *Client) EnqueueRenderVideo(payload RenderVideoPayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeRenderVideo, jsonPayload)
	info, err := c.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// EnqueueVideoComplete enqueues a video completion task
func (c *Client) EnqueueVideoComplete(payload VideoCompletePayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeVideoComplete, jsonPayload)
	info, err := c.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

// Task handlers

// NewHandleGenerateVideoScript creates a handler for video script generation
func NewHandleGenerateVideoScript(db *gorm.DB) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload GenerateVideoScriptPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		log.Printf("Generating script for video: video_id=%d", payload.VideoID)

		// Fetch video from database to get theme
		var video struct {
			ID     int
			Theme  string
			Status string
			Script *string
		}
		if err := db.WithContext(ctx).
			Table("videos").
			Where("id = ?", payload.VideoID).
			First(&video).Error; err != nil {
			return fmt.Errorf("failed to fetch video: %w", err)
		}

		log.Printf("Video theme: %s", video.Theme)

		// Update status to "generating_script"
		if err := db.WithContext(ctx).
			Model(&struct {
				ID     int
				Status string
			}{}).
			Table("videos").
			Where("id = ?", payload.VideoID).
			Update("status", "generating_script").Error; err != nil {
			return fmt.Errorf("failed to update video status: %w", err)
		}

		log.Printf("Status updated to 'generating_script' for video_id=%d", payload.VideoID)

		// Create AI service
		aiService, err := gemini.NewService(ctx)
		if err != nil {
			// Update status to "failed" if we can't create the AI service
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("videos").
				Where("id = ?", payload.VideoID).
				Update("status", "failed")
			return fmt.Errorf("failed to create AI service: %w", err)
		}

		// Generate script using Google AI (250-300 words for ~70 seconds)
		script, err := aiService.GenerateVideoScript(ctx, video.Theme)
		if err != nil {
			// Update status to "failed" if script generation fails
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("videos").
				Where("id = ?", payload.VideoID).
				Update("status", "failed")
			return fmt.Errorf("failed to generate script: %w", err)
		}

		log.Printf("Script generated for video_id=%d (length: %d characters)", payload.VideoID, len(script))

		// Update script
		if err := db.WithContext(ctx).
			Model(&struct {
				ID     int
				Script string
			}{}).
			Table("videos").
			Where("id = ?", payload.VideoID).
			Update("script", script).Error; err != nil {
			return fmt.Errorf("failed to update video script: %w", err)
		}

		// Update status to "completed"
		if err := db.WithContext(ctx).
			Model(&struct {
				ID     int
				Status string
			}{}).
			Table("videos").
			Where("id = ?", payload.VideoID).
			Update("status", "completed").Error; err != nil {
			return fmt.Errorf("failed to update video status: %w", err)
		}

		log.Printf("Video script generation completed: video_id=%d", payload.VideoID)

		// Enqueue the generate audio task
		queueClient := GetClient()
		if queueClient != nil {
			log.Printf("Enqueueing generate audio task for video_id=%d", payload.VideoID)
			err := queueClient.EnqueueGenerateAudio(GenerateAudioPayload(payload))
			if err != nil {
				log.Printf("Failed to enqueue generate audio task: %v", err)
				// Don't return error as script generation was successful
			}
		} else {
			log.Printf("Warning: Queue client not initialized, skipping audio generation")
		}

		return nil
	}
}

// NewHandleGenerateAudio creates a handler for audio generation
func NewHandleGenerateAudio(db *gorm.DB) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload GenerateAudioPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		log.Printf("Generating audio for video: video_id=%d", payload.VideoID)

		// Fetch video from database to get script and voice_id
		var video struct {
			ID       int
			Script   *string
			VoiceID  string
			AudioURL *string
			Status   string
		}
		if err := db.WithContext(ctx).
			Table("videos").
			Where("id = ?", payload.VideoID).
			First(&video).Error; err != nil {
			return fmt.Errorf("failed to fetch video: %w", err)
		}

		// Check if script exists
		if video.Script == nil || *video.Script == "" {
			return fmt.Errorf("video has no script to generate audio from")
		}

		log.Printf("Video script length: %d characters", len(*video.Script))

		// Update status to "generating_audio"
		if err := db.WithContext(ctx).
			Model(&struct {
				ID     int
				Status string
			}{}).
			Table("videos").
			Where("id = ?", payload.VideoID).
			Update("status", "generating_audio").Error; err != nil {
			return fmt.Errorf("failed to update video status: %w", err)
		}

		log.Printf("Status updated to 'generating_audio' for video_id=%d", payload.VideoID)

		// Create ElevenLabs service
		elevenLabsService, err := ai.NewElevenLabsService()
		if err != nil {
			log.Printf("ERROR: Failed to create ElevenLabs service for video_id=%d: %v", payload.VideoID, err)
			// Update status to "failed" if we can't create the service
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("videos").
				Where("id = ?", payload.VideoID).
				Update("status", "failed")
			return fmt.Errorf("failed to create ElevenLabs service: %w", err)
		}

		// Generate audio using ElevenLabs
		// Hardcoded voice ID for now
		voiceID := "NNl6r8mD7vthiJatiJt1"
		log.Printf("Calling ElevenLabs API for video_id=%d with voice_id=%s", payload.VideoID, voiceID)
		audioData, err := elevenLabsService.GenerateAudio(ctx, *video.Script, voiceID)
		if err != nil {
			log.Printf("ERROR: Failed to generate audio for video_id=%d: %v", payload.VideoID, err)
			// Update status to "failed" if audio generation fails
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("videos").
				Where("id = ?", payload.VideoID).
				Update("status", "failed")
			return fmt.Errorf("failed to generate audio: %w", err)
		}

		log.Printf("Audio generated for video_id=%d (size: %d bytes)", payload.VideoID, len(audioData))

		// Create S3 service
		s3Service, err := storage.NewS3Service(ctx)
		if err != nil {
			log.Printf("ERROR: Failed to create S3 service for video_id=%d: %v", payload.VideoID, err)
			// Update status to "failed" if we can't create the S3 service
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("videos").
				Where("id = ?", payload.VideoID).
				Update("status", "failed")
			return fmt.Errorf("failed to create S3 service: %w", err)
		}

		// Upload audio to S3
		audioURL, err := s3Service.UploadAudio(ctx, audioData, payload.VideoID)
		if err != nil {
			log.Printf("ERROR: Failed to upload audio to S3 for video_id=%d: %v", payload.VideoID, err)
			// Update status to "failed" if upload fails
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("videos").
				Where("id = ?", payload.VideoID).
				Update("status", "failed")
			return fmt.Errorf("failed to upload audio to S3: %w", err)
		}

		log.Printf("Audio uploaded to S3: %s", audioURL)

		// Update audio_url in the database
		if err := db.WithContext(ctx).
			Model(&struct {
				ID       int
				AudioURL string
			}{}).
			Table("videos").
			Where("id = ?", payload.VideoID).
			Update("audio_url", audioURL).Error; err != nil {
			return fmt.Errorf("failed to update video audio_url: %w", err)
		}

		// Update status to "completed"
		if err := db.WithContext(ctx).
			Model(&struct {
				ID     int
				Status string
			}{}).
			Table("videos").
			Where("id = ?", payload.VideoID).
			Update("status", "completed").Error; err != nil {
			return fmt.Errorf("failed to update video status: %w", err)
		}

		log.Printf("Audio generation completed: video_id=%d, audio_url=%s", payload.VideoID, audioURL)

		// Enqueue the generate captions and scenes tasks in parallel
		queueClient := GetClient()
		if queueClient != nil {
			// Enqueue caption generation
			log.Printf("Enqueueing generate captions task for video_id=%d", payload.VideoID)
			err := queueClient.EnqueueGenerateCaptions(GenerateCaptionsPayload(payload))
			if err != nil {
				log.Printf("Failed to enqueue generate captions task: %v", err)
				// Don't return error as audio generation was successful
			}

			// Enqueue scene generation
			log.Printf("Enqueueing generate scenes task for video_id=%d", payload.VideoID)
			err = queueClient.EnqueueGenerateScenes(GenerateScenesPayload(payload))
			if err != nil {
				log.Printf("Failed to enqueue generate scenes task: %v", err)
				// Don't return error as audio generation was successful
			}
		} else {
			log.Printf("Warning: Queue client not initialized, skipping caption and scene generation")
		}

		return nil
	}
}

// NewHandleGenerateCaptions creates a handler for caption generation
func NewHandleGenerateCaptions(db *gorm.DB) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload GenerateCaptionsPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		log.Printf("Generating captions for video: video_id=%d", payload.VideoID)

		// Fetch video from database to get audio_url
		var video struct {
			ID       int
			AudioURL *string
			Status   string
		}
		if err := db.WithContext(ctx).
			Table("videos").
			Where("id = ?", payload.VideoID).
			First(&video).Error; err != nil {
			return fmt.Errorf("failed to fetch video: %w", err)
		}

		// Check if audio URL exists
		if video.AudioURL == nil || *video.AudioURL == "" {
			return fmt.Errorf("video has no audio URL to generate captions from")
		}

		log.Printf("Generating captions from audio URL: %s", *video.AudioURL)

		// Create Speech-to-Text service
		sttService, err := ai.NewSpeechToTextService(ctx)
		if err != nil {
			log.Printf("ERROR: Failed to create Speech-to-Text service for video_id=%d: %v", payload.VideoID, err)
			return fmt.Errorf("failed to create Speech-to-Text service: %w", err)
		}
		defer sttService.Close()

		// Generate captions from audio URL
		captions, err := sttService.GenerateCaptionsFromURL(ctx, *video.AudioURL)
		if err != nil {
			log.Printf("ERROR: Failed to generate captions for video_id=%d: %v", payload.VideoID, err)
			return fmt.Errorf("failed to generate captions: %w", err)
		}

		log.Printf("Generated %d caption words for video_id=%d", len(captions), payload.VideoID)

		// Convert captions to JSON
		captionsJSON, err := ai.CaptionsToJSON(captions)
		if err != nil {
			log.Printf("ERROR: Failed to convert captions to JSON for video_id=%d: %v", payload.VideoID, err)
			return fmt.Errorf("failed to convert captions to JSON: %w", err)
		}

		// Update captions in the database
		if err := db.WithContext(ctx).
			Model(&struct {
				ID       int
				Captions string
			}{}).
			Table("videos").
			Where("id = ?", payload.VideoID).
			Update("captions", captionsJSON).Error; err != nil {
			return fmt.Errorf("failed to update video captions: %w", err)
		}

		log.Printf("Caption generation completed: video_id=%d, caption_count=%d", payload.VideoID, len(captions))

		// Check if all scenes are completed and trigger render if ready
		checkAndEnqueueRender(ctx, db, payload.VideoID)

		return nil
	}
}

// NewHandleGenerateScenes creates a handler for scene generation
func NewHandleGenerateScenes(db *gorm.DB) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload GenerateScenesPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		log.Printf("Generating scenes for video: video_id=%d", payload.VideoID)

		// Fetch video from database to get script
		var video struct {
			ID     int
			Script *string
			Status string
		}
		if err := db.WithContext(ctx).
			Table("videos").
			Where("id = ?", payload.VideoID).
			First(&video).Error; err != nil {
			return fmt.Errorf("failed to fetch video: %w", err)
		}

		// Check if script exists
		if video.Script == nil || *video.Script == "" {
			return fmt.Errorf("video has no script to generate scenes from")
		}

		log.Printf("Generating scenes based on script (length: %d characters)", len(*video.Script))

		// Update status to "generating_scenes"
		if err := db.WithContext(ctx).
			Model(&struct {
				ID     int
				Status string
			}{}).
			Table("videos").
			Where("id = ?", payload.VideoID).
			Update("status", "generating_scenes").Error; err != nil {
			return fmt.Errorf("failed to update video status: %w", err)
		}

		// Create AI service
		aiService, err := gemini.NewService(ctx)
		if err != nil {
			log.Printf("ERROR: Failed to create AI service for video_id=%d: %v", payload.VideoID, err)
			// Update status to "failed" if we can't create the AI service
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("videos").
				Where("id = ?", payload.VideoID).
				Update("status", "failed")
			return fmt.Errorf("failed to create AI service: %w", err)
		}

		// Generate scenes using Gemini
		scenes, err := aiService.GenerateScenes(ctx, *video.Script)
		if err != nil {
			log.Printf("ERROR: Failed to generate scenes for video_id=%d: %v", payload.VideoID, err)
			// Update status to "failed" if scene generation fails
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("videos").
				Where("id = ?", payload.VideoID).
				Update("status", "failed")
			return fmt.Errorf("failed to generate scenes: %w", err)
		}

		log.Printf("Generated %d scenes for video_id=%d", len(scenes), payload.VideoID)

		// Create VideoScene records in the database
		queueClient := GetClient()
		for _, scene := range scenes {
			videoScene := struct {
				VideoID int
				Prompt  string
				Index   int
				Status  string
			}{
				VideoID: payload.VideoID,
				Prompt:  scene.ImagePrompt,
				Index:   scene.Index,
				Status:  "pending",
			}

			// Insert the scene into the database
			var sceneID int
			err := db.WithContext(ctx).
				Table("video_scenes").
				Create(&videoScene).
				Error
			if err != nil {
				log.Printf("ERROR: Failed to create video scene for video_id=%d: %v", payload.VideoID, err)
				continue
			}

			// Get the ID of the created scene
			err = db.WithContext(ctx).
				Table("video_scenes").
				Where("video_id = ? AND index = ?", payload.VideoID, scene.Index).
				Select("id").
				Scan(&sceneID).
				Error
			if err != nil {
				log.Printf("ERROR: Failed to get scene ID for video_id=%d, index=%d: %v", payload.VideoID, scene.Index, err)
				continue
			}

			log.Printf("Created scene %d for video_id=%d (scene_id=%d)", scene.Index, payload.VideoID, sceneID)

			// Enqueue image generation task for this scene
			if queueClient != nil {
				err := queueClient.EnqueueGenerateSceneImage(GenerateSceneImagePayload{
					SceneID: sceneID,
				})
				if err != nil {
					log.Printf("Failed to enqueue generate scene image task for scene_id=%d: %v", sceneID, err)
				} else {
					log.Printf("Enqueued image generation task for scene_id=%d", sceneID)
				}
			}
		}

		// Update video status to "generating_images"
		if err := db.WithContext(ctx).
			Model(&struct {
				ID     int
				Status string
			}{}).
			Table("videos").
			Where("id = ?", payload.VideoID).
			Update("status", "generating_images").Error; err != nil {
			return fmt.Errorf("failed to update video status: %w", err)
		}

		log.Printf("Scene generation completed: video_id=%d, scenes_count=%d", payload.VideoID, len(scenes))
		return nil
	}
}

// NewHandleGenerateSceneImage creates a handler for scene image generation
func NewHandleGenerateSceneImage(db *gorm.DB) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload GenerateSceneImagePayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		log.Printf("Generating image for scene: scene_id=%d", payload.SceneID)

		// Fetch scene from database
		var scene struct {
			ID       int
			VideoID  int
			Prompt   string
			ImageURL *string
			Index    int
			Status   string
		}
		if err := db.WithContext(ctx).
			Table("video_scenes").
			Where("id = ?", payload.SceneID).
			First(&scene).Error; err != nil {
			return fmt.Errorf("failed to fetch scene: %w", err)
		}

		log.Printf("Scene prompt: %s", scene.Prompt)

		// Update status to "generating"
		if err := db.WithContext(ctx).
			Model(&struct {
				ID     int
				Status string
			}{}).
			Table("video_scenes").
			Where("id = ?", payload.SceneID).
			Update("status", "generating").Error; err != nil {
			return fmt.Errorf("failed to update scene status: %w", err)
		}

		// Create AI service
		aiService, err := gemini.NewService(ctx)
		if err != nil {
			log.Printf("ERROR: Failed to create AI service for scene_id=%d: %v", payload.SceneID, err)
			// Update status to "failed" if we can't create the AI service
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("video_scenes").
				Where("id = ?", payload.SceneID).
				Update("status", "failed")
			return fmt.Errorf("failed to create AI service: %w", err)
		}

		// Generate image using Imagen
		log.Printf("Generating image with Imagen for scene_id=%d with prompt: %s", payload.SceneID, scene.Prompt)
		imageData, err := aiService.GenerateImage(ctx, scene.Prompt)
		if err != nil {
			log.Printf("ERROR: Failed to generate image for scene_id=%d: %v", payload.SceneID, err)
			// Update status to "failed" if image generation fails
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("video_scenes").
				Where("id = ?", payload.SceneID).
				Update("status", "failed")
			return fmt.Errorf("failed to generate image: %w", err)
		}

		log.Printf("Image generated for scene_id=%d (size: %d bytes)", payload.SceneID, len(imageData))

		// Create S3 service
		s3Service, err := storage.NewS3Service(ctx)
		if err != nil {
			log.Printf("ERROR: Failed to create S3 service for scene_id=%d: %v", payload.SceneID, err)
			// Update status to "failed" if we can't create the S3 service
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("video_scenes").
				Where("id = ?", payload.SceneID).
				Update("status", "failed")
			return fmt.Errorf("failed to create S3 service: %w", err)
		}

		// Upload image to S3
		imageURL, err := s3Service.UploadImage(ctx, imageData, scene.VideoID, scene.Index)
		if err != nil {
			log.Printf("ERROR: Failed to upload image to S3 for scene_id=%d: %v", payload.SceneID, err)
			// Update status to "failed" if upload fails
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("video_scenes").
				Where("id = ?", payload.SceneID).
				Update("status", "failed")
			return fmt.Errorf("failed to upload image to S3: %w", err)
		}

		log.Printf("Image uploaded to S3: %s", imageURL)

		// Update scene with image URL
		if err := db.WithContext(ctx).
			Model(&struct {
				ID       int
				ImageURL string
			}{}).
			Table("video_scenes").
			Where("id = ?", payload.SceneID).
			Update("image_url", imageURL).Error; err != nil {
			return fmt.Errorf("failed to update scene image_url: %w", err)
		}

		// Update status to "completed"
		if err := db.WithContext(ctx).
			Model(&struct {
				ID     int
				Status string
			}{}).
			Table("video_scenes").
			Where("id = ?", payload.SceneID).
			Update("status", "completed").Error; err != nil {
			return fmt.Errorf("failed to update scene status: %w", err)
		}

		log.Printf("Scene image generation completed: scene_id=%d, image_url=%s", payload.SceneID, imageURL)

		// Check if all scenes are completed and captions are ready, then trigger render
		checkAndEnqueueRender(ctx, db, scene.VideoID)

		return nil
	}
}

// checkAndEnqueueRender checks if all prerequisites are met to render the video
func checkAndEnqueueRender(ctx context.Context, db *gorm.DB, videoID int) {
	// Fetch video with captions
	var video struct {
		ID       int
		Captions *string
	}
	if err := db.WithContext(ctx).
		Table("videos").
		Where("id = ?", videoID).
		First(&video).Error; err != nil {
		log.Printf("ERROR: Failed to fetch video for render check: %v", err)
		return
	}

	// Check if captions exist
	if video.Captions == nil || *video.Captions == "" {
		log.Printf("Captions not ready yet for video_id=%d, skipping render", videoID)
		return
	}

	// Check if all scenes have completed
	var pendingScenes int64
	if err := db.WithContext(ctx).
		Table("video_scenes").
		Where("video_id = ? AND (status != ? OR image_url IS NULL)", videoID, "completed").
		Count(&pendingScenes).Error; err != nil {
		log.Printf("ERROR: Failed to count pending scenes for video_id=%d: %v", videoID, err)
		return
	}

	if pendingScenes > 0 {
		log.Printf("Still have %d pending scenes for video_id=%d, skipping render", pendingScenes, videoID)
		return
	}

	// Check if there are any scenes at all
	var totalScenes int64
	if err := db.WithContext(ctx).
		Table("video_scenes").
		Where("video_id = ?", videoID).
		Count(&totalScenes).Error; err != nil {
		log.Printf("ERROR: Failed to count total scenes for video_id=%d: %v", videoID, err)
		return
	}

	if totalScenes == 0 {
		log.Printf("No scenes found for video_id=%d, skipping render", videoID)
		return
	}

	log.Printf("All prerequisites met for video_id=%d, enqueueing render task", videoID)

	// Enqueue render video task
	queueClient := GetClient()
	if queueClient != nil {
		err := queueClient.EnqueueRenderVideo(RenderVideoPayload{VideoID: videoID})
		if err != nil {
			log.Printf("ERROR: Failed to enqueue render video task for video_id=%d: %v", videoID, err)
		} else {
			log.Printf("Successfully enqueued render video task for video_id=%d", videoID)
		}
	} else {
		log.Printf("ERROR: Queue client not initialized, cannot enqueue render task")
	}
}

// RemotionLambdaRequest represents the payload sent to Remotion Lambda
type RemotionLambdaRequest struct {
	Scenes []struct {
		ImageURL string `json:"image_url"`
		Index    int    `json:"index"`
	} `json:"scenes"`
	Captions []struct {
		Word      string  `json:"word"`
		StartTime float64 `json:"start_time"`
		EndTime   float64 `json:"end_time"`
	} `json:"captions"`
	AudioURL      string  `json:"audioUrl"`
	VideoDuration float64 `json:"videoDuration"`
	VideoID       int     `json:"videoId"`
}

// RemotionLambdaResponse represents the response from Remotion Lambda
type RemotionLambdaResponse struct {
	Success  bool    `json:"success"`
	VideoURL *string `json:"videoUrl,omitempty"`
	RenderId *string `json:"renderId,omitempty"`
	Error    *string `json:"error,omitempty"`
}

// NewHandleRenderVideo creates a handler for video rendering with Remotion Lambda
func NewHandleRenderVideo(db *gorm.DB) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload RenderVideoPayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		log.Printf("Starting video render for video_id=%d", payload.VideoID)

		// Update status to "rendering"
		if err := db.WithContext(ctx).
			Model(&struct {
				ID     int
				Status string
			}{}).
			Table("videos").
			Where("id = ?", payload.VideoID).
			Update("status", "rendering").Error; err != nil {
			return fmt.Errorf("failed to update video status: %w", err)
		}

		// Fetch video with all required data
		var video struct {
			ID       int
			AudioURL *string
			Captions *string
		}
		if err := db.WithContext(ctx).
			Table("videos").
			Where("id = ?", payload.VideoID).
			First(&video).Error; err != nil {
			updateVideoStatusToFailed(ctx, db, payload.VideoID)
			return fmt.Errorf("failed to fetch video: %w", err)
		}

		// Validate audio URL
		if video.AudioURL == nil || *video.AudioURL == "" {
			updateVideoStatusToFailed(ctx, db, payload.VideoID)
			return fmt.Errorf("video has no audio URL")
		}

		// Validate captions
		if video.Captions == nil || *video.Captions == "" {
			updateVideoStatusToFailed(ctx, db, payload.VideoID)
			return fmt.Errorf("video has no captions")
		}

		// Parse captions JSON
		var captionsData []struct {
			Word      string  `json:"word"`
			StartTime float64 `json:"start_time"`
			EndTime   float64 `json:"end_time"`
		}
		if err := json.Unmarshal([]byte(*video.Captions), &captionsData); err != nil {
			updateVideoStatusToFailed(ctx, db, payload.VideoID)
			return fmt.Errorf("failed to parse captions: %w", err)
		}

		// Fetch all scenes ordered by index
		var scenes []struct {
			ID       int
			VideoID  int
			ImageURL *string
			Index    int
		}
		if err := db.WithContext(ctx).
			Table("video_scenes").
			Where("video_id = ?", payload.VideoID).
			Order("index ASC").
			Find(&scenes).Error; err != nil {
			updateVideoStatusToFailed(ctx, db, payload.VideoID)
			return fmt.Errorf("failed to fetch video scenes: %w", err)
		}

		if len(scenes) == 0 {
			updateVideoStatusToFailed(ctx, db, payload.VideoID)
			return fmt.Errorf("no scenes found for video")
		}

		log.Printf("Found %d scenes for video_id=%d", len(scenes), payload.VideoID)

		// Calculate video duration from captions
		videoDuration := 0.0
		if len(captionsData) > 0 {
			videoDuration = captionsData[len(captionsData)-1].EndTime
		}
		if videoDuration == 0 {
			// Fallback: estimate 5 seconds per scene
			videoDuration = float64(len(scenes)) * 5.0
		}

		// Prepare Lambda request
		lambdaReq := RemotionLambdaRequest{
			AudioURL:      *video.AudioURL,
			VideoDuration: videoDuration,
			VideoID:       payload.VideoID,
		}

		// Add scenes
		for _, scene := range scenes {
			if scene.ImageURL == nil || *scene.ImageURL == "" {
				updateVideoStatusToFailed(ctx, db, payload.VideoID)
				return fmt.Errorf("scene %d has no image URL", scene.Index)
			}
			lambdaReq.Scenes = append(lambdaReq.Scenes, struct {
				ImageURL string `json:"image_url"`
				Index    int    `json:"index"`
			}{
				ImageURL: *scene.ImageURL,
				Index:    scene.Index,
			})
		}

		// Add captions
		for _, caption := range captionsData {
			lambdaReq.Captions = append(lambdaReq.Captions, struct {
				Word      string  `json:"word"`
				StartTime float64 `json:"start_time"`
				EndTime   float64 `json:"end_time"`
			}{
				Word:      caption.Word,
				StartTime: caption.StartTime,
				EndTime:   caption.EndTime,
			})
		}

		// Invoke Lambda
		videoURL, err := invokeLambdaRender(ctx, lambdaReq)
		if err != nil {
			updateVideoStatusToFailed(ctx, db, payload.VideoID)
			return fmt.Errorf("failed to render video on Lambda: %w", err)
		}

		log.Printf("Lambda render completed: video_url=%s", videoURL)

		// Update video with final URL and status
		if err := db.WithContext(ctx).
			Model(&struct {
				ID       int
				VideoURL string
				Status   string
			}{}).
			Table("videos").
			Where("id = ?", payload.VideoID).
			Updates(map[string]interface{}{
				"video_url": videoURL,
				"status":    "completed",
			}).Error; err != nil {
			return fmt.Errorf("failed to update video with final URL: %w", err)
		}

		log.Printf("Video render completed successfully: video_id=%d, video_url=%s", payload.VideoID, videoURL)
		return nil
	}
}

// invokeLambdaRender invokes the Remotion Lambda function to render the video
func invokeLambdaRender(ctx context.Context, req RemotionLambdaRequest) (string, error) {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create Lambda client
	lambdaClient := lambda.NewFromConfig(cfg)

	// Marshal request to JSON
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal lambda request: %w", err)
	}

	// Get Lambda function name from environment
	functionName := os.Getenv("REMOTION_LAMBDA_FUNCTION")
	if functionName == "" {
		return "", fmt.Errorf("REMOTION_LAMBDA_FUNCTION environment variable not set")
	}

	log.Printf("Invoking Lambda function: %s", functionName)

	// Invoke Lambda
	result, err := lambdaClient.Invoke(ctx, &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		Payload:      reqJSON,
	})
	if err != nil {
		return "", fmt.Errorf("failed to invoke lambda: %w", err)
	}

	// Parse response
	var lambdaResp RemotionLambdaResponse
	if err := json.Unmarshal(result.Payload, &lambdaResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal lambda response: %w", err)
	}

	// Check for errors
	if !lambdaResp.Success {
		errorMsg := "unknown error"
		if lambdaResp.Error != nil {
			errorMsg = *lambdaResp.Error
		}
		return "", fmt.Errorf("lambda render failed: %s", errorMsg)
	}

	// Validate video URL
	if lambdaResp.VideoURL == nil || *lambdaResp.VideoURL == "" {
		return "", fmt.Errorf("lambda returned empty video URL")
	}

	return *lambdaResp.VideoURL, nil
}

// NewHandleVideoComplete creates a handler for video completion tasks
func NewHandleVideoComplete(db *gorm.DB) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload VideoCompletePayload
		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		log.Printf("Processing video_complete task for video_id=%d", payload.VideoID)

		// Update video with final URL and status
		if err := db.WithContext(ctx).
			Model(&struct {
				ID       int
				VideoURL string
				Status   string
			}{}).
			Table("videos").
			Where("id = ?", payload.VideoID).
			Updates(map[string]interface{}{
				"video_url": payload.VideoURL,
				"status":    "completed",
			}).Error; err != nil {
			return fmt.Errorf("failed to update video with final URL: %w", err)
		}

		log.Printf("Video completion processed successfully: video_id=%d, video_url=%s", payload.VideoID, payload.VideoURL)
		return nil
	}
}

// Helper functions for video rendering

func updateVideoStatusToFailed(ctx context.Context, db *gorm.DB, videoID int) {
	db.WithContext(ctx).
		Model(&struct {
			ID     int
			Status string
		}{}).
		Table("videos").
		Where("id = ?", videoID).
		Update("status", "failed")
}
