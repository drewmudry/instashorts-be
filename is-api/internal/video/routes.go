package video

import (
	"instashorts-be/is-api/internal/auth"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all video routes
func RegisterRoutes(router *gin.RouterGroup, handler *Handler, authRepo *auth.Repository) {
	videos := router.Group("/videos")
	videos.Use(auth.RequireAuth(authRepo))
	{
		// Video routes
		videos.POST("", handler.CreateVideo)
		videos.GET("", handler.GetMyVideos)
		videos.GET("/:id", handler.GetVideo)
		videos.GET("/:id/status", handler.GetVideoStatus)
		videos.DELETE("/:id", handler.DeleteVideo)
	}
}
