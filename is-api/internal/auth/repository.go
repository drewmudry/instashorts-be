package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionExpired     = errors.New("session expired")
	ErrOAuthAccountExists = errors.New("oauth account already exists")
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindUserByEmail finds a user by email
func (r *Repository) FindUserByEmail(email string) (*User, error) {
	user := &User{}
	result := r.db.Where("email = ?", email).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return user, nil
}

// FindUserByID finds a user by ID
func (r *Repository) FindUserByID(id int) (*User, error) {
	user := &User{}
	result := r.db.First(user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return user, nil
}

// CreateUser creates a new user
func (r *Repository) CreateUser(email, username string, avatarURL *string) (*User, error) {
	user := &User{
		Email:     email,
		Username:  username,
		AvatarURL: avatarURL,
	}
	result := r.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// FindOAuthAccount finds an OAuth account by provider and provider ID
func (r *Repository) FindOAuthAccount(provider OAuthProvider, providerID string) (*OAuthAccount, error) {
	account := &OAuthAccount{}
	result := r.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(account)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return account, nil
}

// CreateOAuthAccount creates a new OAuth account
func (r *Repository) CreateOAuthAccount(userID int, provider OAuthProvider, providerID, email, accessToken string, refreshToken *string, expiresAt *time.Time) (*OAuthAccount, error) {
	account := &OAuthAccount{
		UserID:       userID,
		Provider:     provider,
		ProviderID:   providerID,
		Email:        email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}
	result := r.db.Create(account)
	if result.Error != nil {
		return nil, result.Error
	}
	return account, nil
}

// UpdateOAuthAccountTokens updates the OAuth account tokens
func (r *Repository) UpdateOAuthAccountTokens(id int, accessToken string, refreshToken *string, expiresAt *time.Time) error {
	result := r.db.Model(&OAuthAccount{}).Where("id = ?", id).Updates(map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_at":    expiresAt,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateSession creates a new session for a user
func (r *Repository) CreateSession(userID int, duration time.Duration) (*Session, error) {
	session := &Session{
		ID:        uuid.New().String(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(duration),
	}
	result := r.db.Create(session)
	if result.Error != nil {
		return nil, result.Error
	}
	return session, nil
}

// FindSession finds a session by ID
func (r *Repository) FindSession(sessionID string) (*Session, error) {
	session := &Session{}
	result := r.db.Where("id = ?", sessionID).First(session)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, result.Error
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	return session, nil
}

// DeleteSession deletes a session
func (r *Repository) DeleteSession(sessionID string) error {
	result := r.db.Delete(&Session{}, "id = ?", sessionID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CleanupExpiredSessions removes expired sessions
func (r *Repository) CleanupExpiredSessions() error {
	result := r.db.Where("expires_at < ?", time.Now()).Delete(&Session{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
