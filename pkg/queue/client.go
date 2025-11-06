package queue

import (
	"fmt"
	"os"
	"sync"

	"github.com/hibiken/asynq"
)

// Client wraps the asynq client for enqueuing tasks
type Client struct {
	client *asynq.Client
}

var (
	globalClient *Client
	once         sync.Once
)

// NewClient creates a new queue client
func NewClient() *Client {
	redisAddr := fmt.Sprintf("%s:%s",
		getEnvOrDefault("REDIS_HOST", "localhost"),
		getEnvOrDefault("REDIS_PORT", "6379"),
	)

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	queueClient := &Client{
		client: client,
	}

	// Set the global client
	once.Do(func() {
		globalClient = queueClient
	})

	return queueClient
}

// GetClient returns the global queue client instance
func GetClient() *Client {
	if globalClient == nil {
		globalClient = NewClient()
	}
	return globalClient
}

// GetClient returns the underlying asynq client
func (c *Client) GetClient() *asynq.Client {
	return c.client
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.client.Close()
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
