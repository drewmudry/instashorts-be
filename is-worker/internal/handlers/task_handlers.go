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
			}

			log.Printf("Enqueueing generate scenes task for video_id=%d", payload.VideoID)
			err = queueClient.EnqueueGenerateScenes(queue.GenerateScenesPayload(payload))
			if err != nil {
				log.Printf("Failed to enqueue generate scenes task: %v", err)
			}

		} else {
			log.Printf("Warning: Queue client not initialized, skipping audio and scene generation")
		}

		return nil
	}
}

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

		// Create S3 service (using GCS from your new file)
		// Note: The old file used storage.NewS3Service, the new one uses storage.NewGCSClient
		// I'll use the one from your new structure: storage.NewGCSClient
		s3Service, err := storage.NewGCSClient(ctx) // Using NewGCSClient as per your new `task_handlers.go`
		if err != nil {
			log.Printf("ERROR: Failed to create GCS service for video_id=%d: %v", payload.VideoID, err)
			// Update status to "failed" if we can't create the GCS service
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("videos").
				Where("id = ?", payload.VideoID).
				Update("status", "failed")
			return fmt.Errorf("failed to create GCS service: %w", err)
		}

		// Upload audio to S3/GCS
		audioURL, err := s3Service.UploadAudio(ctx, audioData, payload.VideoID)
		if err != nil {
			log.Printf("ERROR: Failed to upload audio to GCS for video_id=%d: %v", payload.VideoID, err)
			// Update status to "failed" if upload fails
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("videos").
				Where("id = ?", payload.VideoID).
				Update("status", "failed")
			return fmt.Errorf("failed to upload audio to GCS: %w", err)
		}

		log.Printf("Audio uploaded to GCS: %s", audioURL)

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
		var payload queue.GenerateCaptionsPayload
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
		log.Printf("Captions JSON: %s", captionsJSON)

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
		var payload queue.GenerateScenesPayload
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
		queueClient := queue.GetClient()
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
				err := queueClient.EnqueueGenerateSceneImage(queue.GenerateSceneImagePayload{
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
		var payload queue.GenerateSceneImagePayload
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

		// Create S3/GCS service
		// Using NewGCSClient as per your new `task_handlers.go`
		s3Service, err := storage.NewGCSClient(ctx)
		if err != nil {
			log.Printf("ERROR: Failed to create GCS service for scene_id=%d: %v", payload.SceneID, err)
			// Update status to "failed" if we can't create the S3 service
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("video_scenes").
				Where("id = ?", payload.SceneID).
				Update("status", "failed")
			return fmt.Errorf("failed to create GCS service: %w", err)
		}

		// Upload image to S3/GCS
		imageURL, err := s3Service.UploadImage(ctx, imageData, scene.VideoID, scene.Index)
		if err != nil {
			log.Printf("ERROR: Failed to upload image to GCS for scene_id=%d: %v", payload.SceneID, err)
			// Update status to "failed" if upload fails
			db.WithContext(ctx).
				Model(&struct {
					ID     int
					Status string
				}{}).
				Table("video_scenes").
				Where("id = ?", payload.SceneID).
				Update("status", "failed")
			return fmt.Errorf("failed to upload image to GCS: %w", err)
		}

		log.Printf("Image uploaded to GCS: %s", imageURL)

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
	// --- MODIFICATION: Fetch video with captions AND status ---
	var video struct {
		ID       int
		Captions *string
		Status   string // <-- Added this field
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

	// --- MODIFICATION: Added check to prevent race condition ---
	// If captions are done, but the scenes task hasn't finished (i.e., status isn't 'generating_images' or 'ready_to_render'),
	// then it's too early for us to check for pending scenes.
	// The scene_image handlers will call this check again later.
	if video.Status != "generating_images" && video.Status != "ready_to_render" {
		log.Printf("Captions are done, but scene generation is not yet complete (status: %s). Skipping render check.", video.Status)
		return
	}
	// --- END MODIFICATION ---

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
		// This check is now safe, as we know the scene gen step has run
		log.Printf("No scenes found for video_id=%d, skipping render", videoID)
		return
	}

	// --- MODIFICATION: Update status to 'ready_to_render' BEFORE enqueueing ---
	log.Printf("All prerequisites met for video_id=%d, updating status and enqueueing render task", videoID)

	// Update status to "ready_to_render"
	if err := db.WithContext(ctx).
		Model(&struct {
			ID     int
			Status string
		}{}).
		Table("videos").
		Where("id = ?", videoID).
		Update("status", "ready_to_render").Error; err != nil {
		log.Printf("ERROR: Failed to update video status to ready_to_render: %v", err)
		// Don't return, still try to enqueue
	}
	// --- END MODIFICATION ---

	// Enqueue render video task
	queueClient := queue.GetClient()
	if queueClient != nil {
		err := queueClient.EnqueueRenderVideo(queue.RenderVideoPayload{VideoID: videoID})
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
		var payload queue.RenderVideoPayload
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

		if err := db.WithContext(ctx).
			Model(&struct {
				ID       int
				VideoURL string
				Status   string
			}{}).
			Table("videos").
			Where("id = ?", payload.VideoID).
			Updates(map[string]interface{}{
				"status": "completed",
			}).Error; err != nil {
			return fmt.Errorf("failed to update video with final URL: %w", err)
		}

		log.Printf("Video render completed successfully: video_id=%d", payload.VideoID)
		return nil
	}
}

// NewHandleVideoComplete creates a handler for video completion tasks
func NewHandleVideoComplete(db *gorm.DB) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var payload queue.VideoCompletePayload
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
