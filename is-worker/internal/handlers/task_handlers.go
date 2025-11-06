package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"instashorts-be/is-worker/internal/ai"
	"instashorts-be/is-worker/internal/ai/gemini"
	"instashorts-be/is-worker/internal/storage"
	"instashorts-be/pkg/queue"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// NewHandleGenerateVideoScript creates a handler for video script generation
func NewHandleGenerateVideoScript(db *gorm.DB) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload queue.GenerateVideoScriptPayload
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
		queueClient := queue.GetClient()
		if queueClient != nil {
			log.Printf("Enqueueing generate audio task for video_id=%d", payload.VideoID)
			err := queueClient.EnqueueGenerateAudio(queue.GenerateAudioPayload(payload))
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
		var payload queue.GenerateAudioPayload
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
		s3Service, err := storage.NewGCSClient(ctx)
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
		queueClient := queue.GetClient()
		if queueClient != nil {
			// Enqueue caption generation
			log.Printf("Enqueueing generate captions task for video_id=%d", payload.VideoID)
			err := queueClient.EnqueueGenerateCaptions(queue.GenerateCaptionsPayload(payload))
			if err != nil {
				log.Printf("Failed to enqueue generate captions task: %v", err)
				// Don't return error as audio generation was successful
			}

			// Enqueue scene generation
			log.Printf("Enqueueing generate scenes task for video_id=%d", payload.VideoID)
			err = queueClient.EnqueueGenerateScenes(queue.GenerateScenesPayload(payload))
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

// Additional handlers for captions, scenes, etc. would go here...
// For brevity, I'll create placeholder functions

func NewHandleGenerateCaptions(db *gorm.DB) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		// TODO: Implement caption generation
		log.Println("Caption generation handler - not yet implemented")
		return nil
	}
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

func NewHandleVideoComplete(db *gorm.DB) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		// TODO: Implement video completion
		log.Println("Video completion handler - not yet implemented")
		return nil
	}
}
