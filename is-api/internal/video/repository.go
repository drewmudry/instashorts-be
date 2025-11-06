package video

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateVideo creates a new video in the database
func (r *Repository) CreateVideo(ctx context.Context, video *Video) error {
	return r.db.WithContext(ctx).Create(video).Error
}

// GetVideoByID retrieves a video by its ID with its scenes
func (r *Repository) GetVideoByID(ctx context.Context, id int) (*Video, error) {
	var video Video
	err := r.db.WithContext(ctx).
		Preload("Scenes", func(db *gorm.DB) *gorm.DB {
			return db.Order("video_scenes.index ASC")
		}).
		First(&video, id).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// GetVideosByUserID retrieves all videos for a user
func (r *Repository) GetVideosByUserID(ctx context.Context, userID int) ([]Video, error) {
	var videos []Video
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&videos).Error
	return videos, err
}

// UpdateVideo updates an existing video
func (r *Repository) UpdateVideo(ctx context.Context, video *Video) error {
	return r.db.WithContext(ctx).Save(video).Error
}

// UpdateVideoScript updates only the script of a video
func (r *Repository) UpdateVideoScript(ctx context.Context, id int, script string) error {
	return r.db.WithContext(ctx).
		Model(&Video{}).
		Where("id = ?", id).
		Update("script", script).Error
}

// DeleteVideo soft deletes a video
func (r *Repository) DeleteVideo(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&Video{}, id).Error
}
