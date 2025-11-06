package auth

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Username  string         `json:"username" gorm:"not null"`
	AvatarURL *string        `json:"avatar_url,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// OAuthProvider represents the type of OAuth provider
type OAuthProvider string

const (
	OAuthProviderGoogle  OAuthProvider = "google"
	OAuthProviderDiscord OAuthProvider = "discord"
)

// OAuthAccount represents a user's OAuth account from a provider
type OAuthAccount struct {
	ID           int            `json:"id" gorm:"primaryKey"`
	UserID       int            `json:"user_id" gorm:"not null;index"`
	Provider     OAuthProvider  `json:"provider" gorm:"type:varchar(50);not null"`
	ProviderID   string         `json:"provider_id" gorm:"not null"`
	Email        string         `json:"email" gorm:"not null"`
	AccessToken  string         `json:"-" gorm:"not null"`
	RefreshToken *string        `json:"-"`
	ExpiresAt    *time.Time     `json:"expires_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName overrides the default table name for GORM
func (OAuthAccount) TableName() string {
	return "oauth_accounts"
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(255)"`
	UserID    int       `json:"user_id" gorm:"not null;index"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// UserProfile represents OAuth user info from providers
type UserProfile struct {
	ProviderID string
	Email      string
	Username   string
	AvatarURL  *string
}
