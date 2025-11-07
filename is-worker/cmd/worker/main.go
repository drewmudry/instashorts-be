package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	_ "github.com/joho/godotenv/autoload"

	// Updated imports for monorepo
	"instashorts-be/is-worker/internal/handlers"
	"instashorts-be/pkg/database"
	"instashorts-be/pkg/queue"
)

func main() {
	log.Println("Starting worker process...")

	redisAddr := fmt.Sprintf("%s:%s",
		getEnvOrDefault("REDIS_HOST", "localhost"),
		getEnvOrDefault("REDIS_PORT", "6379"),
	)

	// Initialize database (using new package path)
	db := database.New()
	gormDB := db.GetDB()

	// Create asynq server
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority levels
			Queues: map[string]int{
				"critical": 6, // processed 60% of the time
				"default":  3, // processed 30% of the time
				"low":      1, // processed 10% of the time
			},
			// See the godoc for other configuration options
		},
	)

	// Create mux to map task types to handlers
	mux := asynq.NewServeMux()

	// Register task handlers (using new 'handlers' package)
	mux.HandleFunc(queue.TypeGenerateVideoScript, handlers.NewHandleGenerateVideoScript(gormDB))
	mux.HandleFunc(queue.TypeGenerateAudio, handlers.NewHandleGenerateAudio(gormDB))
	mux.HandleFunc(queue.TypeGenerateCaptions, handlers.NewHandleGenerateCaptions(gormDB))
	mux.HandleFunc(queue.TypeGenerateScenes, handlers.NewHandleGenerateScenes(gormDB))
	mux.HandleFunc(queue.TypeGenerateSceneImage, handlers.NewHandleGenerateSceneImage(gormDB))
	// Note: TypeRenderVideo is now handled by the TypeScript renderer service
	mux.HandleFunc(queue.TypeVideoComplete, handlers.NewHandleVideoComplete(gormDB))

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start health check server on port 5100
	healthPort := getEnvOrDefault("WORKER_PORT", "5100")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"worker"}`))
	})

	go func() {
		log.Printf("Worker health check server listening on :%s", healthPort)
		if err := http.ListenAndServe(":"+healthPort, nil); err != nil {
			log.Printf("Health check server error: %v", err)
		}
	}()

	// Run worker in a goroutine
	go func() {
		log.Printf("Worker connected to Redis at %s", redisAddr)
		log.Println("Worker is now listening for tasks...")
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	// Wait for interrupt signal
	sig := <-sigChan
	log.Printf("Received signal: %v", sig)
	log.Println("Shutting down worker gracefully...")

	// Shutdown the server gracefully
	srv.Shutdown()

	log.Println("Worker shutdown complete")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
