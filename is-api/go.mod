module instashorts-be/is-api

go 1.24.4

require (
	instashorts-be/pkg v0.0.0
	github.com/gin-contrib/cors v1.7.6
	github.com/gin-gonic/gin v1.11.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/oauth2 v0.32.0
	google.golang.org/api v0.247.0
	google.golang.org/genai v1.33.0
)

replace instashorts-be/pkg => ../pkg

