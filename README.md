# Instashorts Backend Monorepo

AI-powered video generation platform backend with Go API, Worker, and TypeScript renderer.

## 🏗️ Architecture

- **is-api**: REST API service (Go/Gin)
- **is-worker**: Background job processor (Go/Asynq)
- **is-render**: Video renderer (TypeScript/Remotion) - *In Development*
- **pkg**: Shared Go packages (database, queue, etc.)

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL (via Docker)
- Redis (via Docker)
- Google Cloud Project with Vertex AI enabled
- AWS S3 bucket for storage
- ElevenLabs API key for text-to-speech

### Setup Instructions

1. **Clone and navigate to the project**:
   ```bash
   cd instashorts-backend
   ```

2. **Fix Go versions** (Go 1.24 doesn't exist yet, use 1.23):
   ```bash
   chmod +x fix-monorepo.sh
   ./fix-monorepo.sh
   ```

3. **Configure environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your actual credentials
   ```

4. **Add Google Cloud credentials**:
   ```bash
   # Copy your Vertex AI service account key
   cp /path/to/your/vertex-ai-key.json ./vertex-ai-key.json
   ```

5. **Build and run with Docker**:
   ```bash
   # Build all services
   docker-compose build

   # Start all services
   docker-compose up -d

   # View logs
   docker-compose logs -f
   ```

6. **Run database migrations**:
   ```bash
   # Wait for postgres to be ready, then:
   make migrate-up
   ```

## 📝 Development

### Local Development (without Docker)

```bash
# Install dependencies
make install

# Run API locally
make dev-api

# Run Worker locally
make dev-worker

# Run tests
make test
```

### Working with the Monorepo

The project uses Go workspaces for managing multiple modules:

```bash
# Sync workspace dependencies
go work sync

# Add a new module to workspace
go work use ./new-module
```

### Database Migrations

```bash
# Create new migration
migrate create -ext sql -dir is-api/migrations -seq migration_name

# Run migrations
make migrate-up

# Rollback last migration
make migrate-down
```

## 🔧 Configuration

### Required Environment Variables

- **Database**: `BLUEPRINT_DB_*` variables for PostgreSQL
- **Redis**: `REDIS_HOST`, `REDIS_PORT`
- **OAuth**: Google and Discord OAuth credentials
- **Google Cloud**: `GCP_PROJECT_ID`, service account key
- **AWS**: S3 credentials and bucket name
- **ElevenLabs**: API key for text-to-speech

See `.env.example` for all variables.

## 📁 Project Structure

```
.
├── go.work                 # Go workspace configuration
├── pkg/                    # Shared packages
│   ├── database/          # Database connection
│   └── queue/            # Queue client and tasks
├── is-api/                # API service
│   ├── cmd/api/          # API entry point
│   ├── internal/         # API business logic
│   └── migrations/       # Database migrations
├── is-worker/             # Worker service
│   ├── cmd/worker/       # Worker entry point
│   └── internal/         # Worker business logic
└── is-render/             # Renderer service (TypeScript)
```

## 🐛 Troubleshooting

### Import Issues

If you see import errors:
1. Ensure all modules use the same prefix: `instashorts-be/`
2. Run `go work sync`
3. Check that `replace` directives point to `../pkg`

### Docker Build Issues

If Docker builds fail:
1. Ensure go.work is present in the build context
2. Check that Go version is 1.23 (not 1.24)
3. Verify all source directories are copied in Dockerfile

### Database Connection Issues

1. Check PostgreSQL is running: `docker-compose ps postgres`
2. Verify credentials in .env match docker-compose.yml
3. Ensure migrations have run successfully

## 🚦 API Endpoints

- `GET /health` - Health check
- `GET /api/auth/google` - Start Google OAuth
- `POST /api/videos` - Create new video
- `GET /api/videos/:id` - Get video details
- `GET /api/videos/:id/status` - Get video processing status

## 📊 Video Processing Pipeline

1. **Script Generation** - Generate video script using Gemini
2. **Audio Generation** - Convert script to speech with ElevenLabs
3. **Caption Generation** - Extract word-level timestamps
4. **Scene Generation** - Generate scene descriptions with Gemini
5. **Image Generation** - Create images with Imagen 4.0
6. **Video Rendering** - Combine assets with Remotion Lambda

## 🤝 Contributing

1. Create a feature branch
2. Make your changes
3. Run tests: `make test`
4. Submit a pull request

## 📄 License

[Your License Here]