package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	ErrInvalidProvider = errors.New("invalid oauth provider")
	ErrInvalidState    = errors.New("invalid state parameter")
)

// OAuthConfig holds OAuth configuration for different providers
type OAuthConfig struct {
	Google  *oauth2.Config
	Discord *oauth2.Config
}

// NewOAuthConfig creates OAuth configurations for all providers
func NewOAuthConfig(googleClientID, googleClientSecret, googleRedirectURL string, discordClientID, discordClientSecret, discordRedirectURL string) *OAuthConfig {
	config := &OAuthConfig{}

	// Google OAuth configuration
	if googleClientID != "" && googleClientSecret != "" {
		config.Google = &oauth2.Config{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  googleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
	}

	// Discord OAuth configuration (for future use)
	if discordClientID != "" && discordClientSecret != "" {
		config.Discord = &oauth2.Config{
			ClientID:     discordClientID,
			ClientSecret: discordClientSecret,
			RedirectURL:  discordRedirectURL,
			Scopes:       []string{"identify", "email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://discord.com/api/oauth2/authorize",
				TokenURL: "https://discord.com/api/oauth2/token",
			},
		}
	}

	return config
}

// GetAuthURL returns the OAuth authorization URL for the specified provider
func (c *OAuthConfig) GetAuthURL(provider OAuthProvider, state string) (string, error) {
	switch provider {
	case OAuthProviderGoogle:
		if c.Google == nil {
			return "", ErrInvalidProvider
		}
		return c.Google.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
	case OAuthProviderDiscord:
		if c.Discord == nil {
			return "", ErrInvalidProvider
		}
		return c.Discord.AuthCodeURL(state), nil
	default:
		return "", ErrInvalidProvider
	}
}

// ExchangeCode exchanges an authorization code for a token
func (c *OAuthConfig) ExchangeCode(ctx context.Context, provider OAuthProvider, code string) (*oauth2.Token, error) {
	switch provider {
	case OAuthProviderGoogle:
		if c.Google == nil {
			return nil, ErrInvalidProvider
		}
		return c.Google.Exchange(ctx, code)
	case OAuthProviderDiscord:
		if c.Discord == nil {
			return nil, ErrInvalidProvider
		}
		return c.Discord.Exchange(ctx, code)
	default:
		return nil, ErrInvalidProvider
	}
}

// GetUserProfile fetches user profile from the OAuth provider
func (c *OAuthConfig) GetUserProfile(ctx context.Context, provider OAuthProvider, token *oauth2.Token) (*UserProfile, error) {
	switch provider {
	case OAuthProviderGoogle:
		return c.getGoogleUserProfile(ctx, token)
	case OAuthProviderDiscord:
		return c.getDiscordUserProfile(ctx, token)
	default:
		return nil, ErrInvalidProvider
	}
}

// getGoogleUserProfile fetches user info from Google
func (c *OAuthConfig) getGoogleUserProfile(ctx context.Context, token *oauth2.Token) (*UserProfile, error) {
	client := c.Google.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: status %d, body: %s", resp.StatusCode, string(body))
	}

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Generate username from email or name
	username := googleUser.Name
	if username == "" {
		username = googleUser.Email
	}

	return &UserProfile{
		ProviderID: googleUser.ID,
		Email:      googleUser.Email,
		Username:   username,
		AvatarURL:  &googleUser.Picture,
	}, nil
}

// getDiscordUserProfile fetches user info from Discord
func (c *OAuthConfig) getDiscordUserProfile(ctx context.Context, token *oauth2.Token) (*UserProfile, error) {
	client := c.Discord.Client(ctx, token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: status %d, body: %s", resp.StatusCode, string(body))
	}

	var discordUser struct {
		ID            string  `json:"id"`
		Username      string  `json:"username"`
		Discriminator string  `json:"discriminator"`
		Email         string  `json:"email"`
		Avatar        *string `json:"avatar"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&discordUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Construct Discord username
	username := discordUser.Username
	if discordUser.Discriminator != "" && discordUser.Discriminator != "0" {
		username = fmt.Sprintf("%s#%s", discordUser.Username, discordUser.Discriminator)
	}

	// Construct avatar URL if available
	var avatarURL *string
	if discordUser.Avatar != nil && *discordUser.Avatar != "" {
		url := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", discordUser.ID, *discordUser.Avatar)
		avatarURL = &url
	}

	return &UserProfile{
		ProviderID: discordUser.ID,
		Email:      discordUser.Email,
		Username:   username,
		AvatarURL:  avatarURL,
	}, nil
}

// OAuthService handles OAuth authentication logic
type OAuthService struct {
	config *OAuthConfig
	repo   *Repository
}

func NewOAuthService(config *OAuthConfig, repo *Repository) *OAuthService {
	return &OAuthService{
		config: config,
		repo:   repo,
	}
}

// HandleOAuthCallback processes the OAuth callback and creates or updates user
func (s *OAuthService) HandleOAuthCallback(ctx context.Context, provider OAuthProvider, code string) (*User, *Session, error) {
	// Exchange code for token
	token, err := s.config.ExchangeCode(ctx, provider, code)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user profile from provider
	profile, err := s.config.GetUserProfile(ctx, provider, token)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Check if OAuth account exists
	oauthAccount, err := s.repo.FindOAuthAccount(provider, profile.ProviderID)
	if err != nil && err != ErrUserNotFound {
		return nil, nil, fmt.Errorf("error checking oauth account: %w", err)
	}

	var user *User
	var expiresAt *time.Time
	if token.Expiry != (time.Time{}) {
		expiresAt = &token.Expiry
	}

	if oauthAccount != nil {
		// Existing OAuth account - get user
		user, err = s.repo.FindUserByID(oauthAccount.UserID)
		if err != nil {
			return nil, nil, fmt.Errorf("error finding user: %w", err)
		}

		// Update OAuth tokens
		refreshToken := &token.RefreshToken
		if token.RefreshToken == "" {
			refreshToken = nil
		}
		err = s.repo.UpdateOAuthAccountTokens(oauthAccount.ID, token.AccessToken, refreshToken, expiresAt)
		if err != nil {
			return nil, nil, fmt.Errorf("error updating tokens: %w", err)
		}
	} else {
		// New user - create user and OAuth account
		user, err = s.repo.CreateUser(profile.Email, profile.Username, profile.AvatarURL)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating user: %w", err)
		}

		refreshToken := &token.RefreshToken
		if token.RefreshToken == "" {
			refreshToken = nil
		}
		_, err = s.repo.CreateOAuthAccount(
			user.ID,
			provider,
			profile.ProviderID,
			profile.Email,
			token.AccessToken,
			refreshToken,
			expiresAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating oauth account: %w", err)
		}
	}

	// Create session (30 days)
	session, err := s.repo.CreateSession(user.ID, 30*24*time.Hour)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating session: %w", err)
	}

	return user, session, nil
}
