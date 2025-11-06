package auth

import (
	"crypto/rand"
	"encoding/base64"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	stateCallbackKey  = "oauth_state"
	sessionCookieName = "session_token"
)

type Handler struct {
	oauthConfig  *OAuthConfig
	oauthService *OAuthService
	repo         *Repository
}

func NewHandler(oauthConfig *OAuthConfig, oauthService *OAuthService, repo *Repository) *Handler {
	return &Handler{
		oauthConfig:  oauthConfig,
		oauthService: oauthService,
		repo:         repo,
	}
}

// generateStateToken generates a random state token for CSRF protection
func generateStateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GoogleLogin initiates Google OAuth flow
func (h *Handler) GoogleLogin(c *gin.Context) {
	state, err := generateStateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state token"})
		return
	}

	// Store state in session cookie for verification
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(stateCallbackKey, state, 600, "/", "", false, true) // 10 minutes

	url, err := h.oauthConfig.GetAuthURL(OAuthProviderGoogle, state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OAuth not configured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

// GoogleCallback handles the OAuth callback from Google
func (h *Handler) GoogleCallback(c *gin.Context) {
	// Verify state token
	state := c.Query("state")
	storedState, err := c.Cookie(stateCallbackKey)
	if err != nil || state != storedState {
		c.Redirect(http.StatusFound, "http://localhost:5173/?error=invalid_state")
		return
	}

	// Clear the state cookie
	c.SetCookie(stateCallbackKey, "", -1, "/", "", false, true)

	// Get authorization code
	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusFound, "http://localhost:5173/?error=missing_code")
		return
	}

	// Handle OAuth callback
	_, session, err := h.oauthService.HandleOAuthCallback(c.Request.Context(), OAuthProviderGoogle, code)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:5173/?error=auth_failed")
		return
	}

	// Set session cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		sessionCookieName,
		session.ID,
		int(time.Until(session.ExpiresAt).Seconds()),
		"/",
		"",
		false, // Set to true in production with HTTPS
		true,  // HttpOnly
	)

	// Redirect to frontend dashboard
	c.Redirect(http.StatusFound, "http://localhost:5173/dashboard")
}

// DiscordLogin initiates Discord OAuth flow
func (h *Handler) DiscordLogin(c *gin.Context) {
	state, err := generateStateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state token"})
		return
	}

	// Store state in session cookie for verification
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(stateCallbackKey, state, 600, "/", "", false, true) // 10 minutes

	url, err := h.oauthConfig.GetAuthURL(OAuthProviderDiscord, state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OAuth not configured"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

// DiscordCallback handles the OAuth callback from Discord
func (h *Handler) DiscordCallback(c *gin.Context) {
	// Verify state token
	state := c.Query("state")
	storedState, err := c.Cookie(stateCallbackKey)
	if err != nil || state != storedState {
		c.Redirect(http.StatusFound, "http://localhost:5173/?error=invalid_state")
		return
	}

	// Clear the state cookie
	c.SetCookie(stateCallbackKey, "", -1, "/", "", false, true)

	// Get authorization code
	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusFound, "http://localhost:5173/?error=missing_code")
		return
	}

	// Handle OAuth callback
	_, session, err := h.oauthService.HandleOAuthCallback(c.Request.Context(), OAuthProviderDiscord, code)
	if err != nil {
		c.Redirect(http.StatusFound, "http://localhost:5173/?error=auth_failed")
		return
	}

	// Set session cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		sessionCookieName,
		session.ID,
		int(time.Until(session.ExpiresAt).Seconds()),
		"/",
		"",
		false, // Set to true in production with HTTPS
		true,  // HttpOnly
	)

	// Redirect to frontend dashboard
	c.Redirect(http.StatusFound, "http://localhost:5173/dashboard")
}

// GetCurrentUser returns the currently authenticated user
func (h *Handler) GetCurrentUser(c *gin.Context) {
	// User is already set in context by RequireAuth middleware
	user, exists := GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Logout logs out the current user
func (h *Handler) Logout(c *gin.Context) {
	// Get session ID from cookie
	sessionID, err := c.Cookie(sessionCookieName)
	if err == nil && sessionID != "" {
		// Delete session from database
		_ = h.repo.DeleteSession(sessionID)
	}

	// Clear session cookie
	c.SetCookie(sessionCookieName, "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
