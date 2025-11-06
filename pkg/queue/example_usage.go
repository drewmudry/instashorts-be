package queue

// Example usage in your API handlers:
//
// 1. Initialize the queue client in your server setup:
//    queueClient := queue.NewClient()
//    defer queueClient.Close()
//
// 2. Use it in your handlers to enqueue tasks:
//
//    // Example: Enqueue a video processing task
//    err := queueClient.EnqueueProcessVideo(queue.ProcessVideoPayload{
//        VideoID: "video-123",
//        UserID:  "user-456",
//    })
//    if err != nil {
//        log.Printf("Failed to enqueue video processing: %v", err)
//    }
//
//    // Example: Enqueue an email task
//    err = queueClient.EnqueueSendEmail(queue.SendEmailPayload{
//        To:      "user@example.com",
//        Subject: "Welcome!",
//        Body:    "Thanks for signing up!",
//    })
//    if err != nil {
//        log.Printf("Failed to enqueue email: %v", err)
//    }
//
// 3. To run the worker:
//    go run cmd/worker/main.go
//
// 4. Add these environment variables to your .env file:
//    REDIS_HOST=localhost
//    REDIS_PORT=6379
