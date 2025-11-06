package video

import (
	"time"

	"gorm.io/gorm"
)

// VideoStatus represents the status of a video
type VideoStatus string

const (
	VideoStatusPending          VideoStatus = "pending"
	VideoStatusGeneratingScript VideoStatus = "generating_script"
	VideoStatusGeneratingAudio  VideoStatus = "generating_audio"
	VideoStatusGeneratingScenes VideoStatus = "generating_scenes"
	VideoStatusGeneratingImages VideoStatus = "generating_images"
	VideoStatusRendering        VideoStatus = "rendering"
	VideoStatusProcessing       VideoStatus = "processing"
	VideoStatusCompleted        VideoStatus = "completed"
	VideoStatusFailed           VideoStatus = "failed"
)

// Series represents a collection of videos
type Series struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	UserID    int            `json:"user_id" gorm:"not null;index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Caption represents a word with timing information
type Caption struct {
	Word      string  `json:"word"`
	StartTime float64 `json:"start_time"` // in seconds
	EndTime   float64 `json:"end_time"`   // in seconds
}

// Video represents a video in the system
type Video struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	UserID    int            `json:"user_id" gorm:"not null;index"`
	SeriesID  *int           `json:"series_id,omitempty" gorm:"index"`
	Title     *string        `json:"title,omitempty"`
	Theme     string         `json:"theme" gorm:"not null"`
	VoiceID   string         `json:"voice_id" gorm:"not null"`
	Script    *string        `json:"script,omitempty" gorm:"type:text"`
	AudioURL  *string        `json:"audio_url,omitempty" gorm:"type:text"`
	VideoURL  *string        `json:"video_url,omitempty" gorm:"type:text"` // Final rendered video URL
	Captions  *string        `json:"captions,omitempty" gorm:"type:jsonb"` // JSON array of Caption objects
	Status    VideoStatus    `json:"status" gorm:"type:varchar(50);not null;default:'pending';index"`
	Scenes    []VideoScene   `json:"scenes,omitempty" gorm:"foreignKey:VideoID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// VideoScene represents a scene in a video with its image
type VideoScene struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	VideoID   int            `json:"video_id" gorm:"not null;index"`
	Prompt    string         `json:"prompt" gorm:"type:text;not null"`
	ImageURL  *string        `json:"image_url,omitempty" gorm:"type:text"`
	Index     int            `json:"index" gorm:"not null"`
	Status    string         `json:"status" gorm:"type:varchar(50);not null;default:'pending';index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// CreateVideoRequest represents the request to create a new video
type CreateVideoRequest struct {
	Title   *string `json:"title"`
	Theme   string  `json:"theme" binding:"required"`
	VoiceID string  `json:"voice_id" binding:"required"`
}
