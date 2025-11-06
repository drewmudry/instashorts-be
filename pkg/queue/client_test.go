package queue

import (
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	// Set test environment variables
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	
	client := NewClient()
	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}
	
	if client.client == nil {
		t.Fatal("Expected underlying asynq client to be initialized, got nil")
	}
	
	// Clean up
	err := client.Close()
	if err != nil {
		t.Logf("Warning: Error closing client: %v", err)
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "returns env value when set",
			envKey:       "TEST_KEY",
			envValue:     "test_value",
			defaultValue: "default",
			expected:     "test_value",
		},
		{
			name:         "returns default when env not set",
			envKey:       "NONEXISTENT_KEY",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}
			
			result := getEnvOrDefault(tt.envKey, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

