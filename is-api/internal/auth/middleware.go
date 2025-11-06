package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireAuth is a middleware that requires authentication
func RequireAuth(repo *Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get session ID from cookie
		sessionID, err := c.Cookie(sessionCookieName)
		if err != nil || sessionID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}

		// Find session
		session, err := repo.FindSession(sessionID)
		if err != nil {
			if err == ErrSessionNotFound || err == ErrSessionExpired {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Session invalid or expired"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate session"})
			c.Abort()
			return
		}

		// Get user
		user, err := repo.FindUserByID(session.UserID)
		if err != nil {
			if err == ErrUserNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)

		c.Next()
	}
}

// GetUserFromContext retrieves the authenticated user from the context
func GetUserFromContext(c *gin.Context) (*User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	u, ok := user.(*User)
	return u, ok
}
