package video

import (
	"fmt"
	"net/http"
	"strconv"

	"instashorts-be/is-api/internal/auth"
	"instashorts-be/pkg/queue"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo        *Repository
	queueClient *queue.Client
}

func NewHandler(repo *Repository, queueClient *queue.Client) *Handler {
	return &Handler{
		repo:        repo,
		queueClient: queueClient,
	}
}

// CreateVideo handles the creation of a new video
func (h *Handler) CreateVideo(c *gin.Context) {
	// Get authenticated user
	user, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Parse request body
	var req CreateVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	// Create video
	video := &Video{
		UserID:  user.ID,
		Title:   req.Title,
		Theme:   req.Theme,
		VoiceID: req.VoiceID,
		Status:  VideoStatusPending,
	}

	if err := h.repo.CreateVideo(c.Request.Context(), video); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create video"})
		return
	}

	// Enqueue video script generation task
	if err := h.queueClient.EnqueueGenerateVideoScript(queue.GenerateVideoScriptPayload{
		VideoID: video.ID,
	}); err != nil {
		// Log the error but don't fail the request - video is already created
		fmt.Printf("Failed to enqueue video script generation task: %v\n", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"video":   video,
		"message": "Video created successfully",
	})
}

// GetVideo retrieves a single video by ID
func (h *Handler) GetVideo(c *gin.Context) {
	// Get authenticated user
	user, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Parse video ID
	videoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// Get video
	video, err := h.repo.GetVideoByID(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Check if user owns the video
	if video.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"video": video})
}

// GetVideoStatus retrieves the status of a video by ID
func (h *Handler) GetVideoStatus(c *gin.Context) {
	// Get authenticated user
	user, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Parse video ID
	videoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// Get video
	video, err := h.repo.GetVideoByID(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Check if user owns the video
	if video.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"video_id": video.ID,
		"status":   video.Status,
		"theme":    video.Theme,
		"script":   video.Script,
	})
}

// GetMyVideos retrieves all videos for the authenticated user
func (h *Handler) GetMyVideos(c *gin.Context) {
	// Get authenticated user
	user, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Get videos
	videos, err := h.repo.GetVideosByUserID(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve videos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"videos": videos})
}

// DeleteVideo deletes a video
func (h *Handler) DeleteVideo(c *gin.Context) {
	// Get authenticated user
	user, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Parse video ID
	videoID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	// Get video to check ownership
	video, err := h.repo.GetVideoByID(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Check if user owns the video
	if video.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this video"})
		return
	}

	// Delete video
	if err := h.repo.DeleteVideo(c.Request.Context(), videoID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete video"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video deleted successfully"})
}
