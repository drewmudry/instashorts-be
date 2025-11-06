package queue

// Example: How to integrate the queue client into your server
//
// 1. Update internal/server/server.go to include the queue client:
//
// type Server struct {
//     port int
//     db           database.Service
//     authHandler  *auth.Handler
//     authRepo     *auth.Repository
//     videoHandler *video.Handler
//     videoRepo    *video.Repository
//     queueClient  *queue.Client  // Add this
// }
//
// func NewServer() *http.Server {
//     // ... existing setup ...
//     
//     // Initialize queue client
//     queueClient := queue.NewClient()
//     
//     NewServer := &Server{
//         port:         port,
//         db:           db,
//         authHandler:  authHandler,
//         authRepo:     authRepo,
//         videoHandler: videoHandler,
//         videoRepo:    videoRepo,
//         queueClient:  queueClient,  // Add this
//     }
//     
//     // ... rest of setup ...
// }
//
// 2. Update your handlers to accept the queue client:
//
// In internal/video/handlers.go (example):
//
// type Handler struct {
//     repo        *Repository
//     queueClient *queue.Client  // Add this
// }
//
// func NewHandler(repo *Repository, queueClient *queue.Client) *Handler {
//     return &Handler{
//         repo:        repo,
//         queueClient: queueClient,
//     }
// }
//
// func (h *Handler) CreateVideo(c *gin.Context) {
//     // ... create video in database ...
//     
//     // Enqueue background processing
//     err := h.queueClient.EnqueueProcessVideo(queue.ProcessVideoPayload{
//         VideoID: video.ID,
//         UserID:  userID,
//     })
//     if err != nil {
//         log.Printf("Failed to enqueue video processing: %v", err)
//         // Still return success to user, task will be retried
//     }
//     
//     c.JSON(202, gin.H{
//         "message": "Video created and queued for processing",
//         "video_id": video.ID,
//     })
// }
//
// 3. Don't forget to close the queue client on server shutdown:
//
// In cmd/api/main.go, in the gracefulShutdown function:
//
// func gracefulShutdown(apiServer *http.Server, queueClient *queue.Client, done chan bool) {
//     ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
//     defer stop()
//     
//     <-ctx.Done()
//     
//     log.Println("shutting down gracefully...")
//     
//     // Close queue client
//     if err := queueClient.Close(); err != nil {
//         log.Printf("Error closing queue client: %v", err)
//     }
//     
//     // Shutdown API server
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()
//     if err := apiServer.Shutdown(ctx); err != nil {
//         log.Printf("Server forced to shutdown with error: %v", err)
//     }
//     
//     done <- true
// }

