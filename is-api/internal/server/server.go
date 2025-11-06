package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"instashorts-be/is-api/internal/auth"
	"instashorts-be/pkg/database"
	"instashorts-be/pkg/queue"
	"instashorts-be/is-api/internal/video"
)

type Server struct {
	port int

	db           database.Service
	queueClient  *queue.Client
	authHandler  *auth.Handler
	authRepo     *auth.Repository
	videoHandler *video.Handler
	videoRepo    *video.Repository
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()

	// Initialize auth module with GORM DB
	authRepo := auth.NewRepository(db.GetDB())
	oauthConfig := auth.NewOAuthConfig(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_REDIRECT_URL"),
		os.Getenv("DISCORD_CLIENT_ID"),
		os.Getenv("DISCORD_CLIENT_SECRET"),
		os.Getenv("DISCORD_REDIRECT_URL"),
	)
	oauthService := auth.NewOAuthService(oauthConfig, authRepo)
	authHandler := auth.NewHandler(oauthConfig, oauthService, authRepo)

	// Initialize queue client
	queueClient := queue.NewClient()

	// Initialize video module with GORM DB
	videoRepo := video.NewRepository(db.GetDB())
	videoHandler := video.NewHandler(videoRepo, queueClient)

	NewServer := &Server{
		port:         port,
		db:           db,
		queueClient:  queueClient,
		authHandler:  authHandler,
		authRepo:     authRepo,
		videoHandler: videoHandler,
		videoRepo:    videoRepo,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
