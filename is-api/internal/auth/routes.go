package auth

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all auth routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, authRepo *Repository) {
	auth := router.Group("/auth")
	{
		// Google OAuth
		auth.GET("/google", handler.GoogleLogin)
		auth.GET("/google/callback", handler.GoogleCallback)

		// Discord OAuth (ready for future implementation)
		auth.GET("/discord", handler.DiscordLogin)
		auth.GET("/discord/callback", handler.DiscordCallback)

		// Protected routes (require authentication)
		protected := auth.Group("")
		protected.Use(RequireAuth(authRepo))
		{
			protected.GET("/me", handler.GetCurrentUser)
			protected.POST("/logout", handler.Logout)
		}
	}
}
