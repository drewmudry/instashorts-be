module instashorts-be/is-worker

go 1.24.4

require (
	instashorts-be/pkg v0.0.0
	cloud.google.com/go/speech v1.28.1
	github.com/aws/aws-sdk-go-v2 v1.39.6
	github.com/aws/aws-sdk-go-v2/config v1.31.17
	github.com/aws/aws-sdk-go-v2/service/lambda v1.81.1
	github.com/hibiken/asynq v0.25.1
	github.com/joho/godotenv v1.5.1
	google.golang.org/api v0.247.0
	google.golang.org/genai v1.33.0
)

replace instashorts-be/pkg => ../pkg

