# Instashorts Backend Monorepo

A monorepo containing all backend services for Instashorts using Go Workspaces.

## Structure

```
instashorts-be/
├── pkg/              # Shared Go packages (database, queue)
├── is-api/           # API service (Go)
├── is-worker/        # Worker service (Go)
├── is-render/        # Video renderer service (TypeScript/Remotion)
└── docker-compose.yml # Unified docker compose for all services
```

## Quick Start

### Prerequisites

- Go 1.24+
- Node.js 20+
- Docker & Docker Compose
- Make

### Setup

1. Clone and navigate to the repo:
```bash
cd instashorts-be
```

2. Install dependencies:
```bash
make install
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Start all services with Docker Compose:
```bash
make up
# or with build
make up-build
```

5. Run database migrations:
```bash
make migrate-up
```

## Services

### API (`is-api`)
- REST API for video management
- Authentication (Google/Discord OAuth)
- Runs on port 8000

### Worker (`is-worker`)
- Background job processor
- Handles video generation tasks
- Health check on port 5100

### Renderer (`is-render`)
- TypeScript service using Remotion
- Renders videos from queue tasks
- Uploads to S3

### Infrastructure
- PostgreSQL database
- Redis queue

## Development

### Running Services Locally

```bash
# API
make dev-api

# Worker
make dev-worker

# Renderer
make dev-renderer
```

### Building

```bash
# Build all
make build

# Build specific service
make build-api
make build-worker
make build-renderer
```

### Testing

```bash
# Run all tests
make test

# Run specific service tests
make test-api
make test-worker
```

## Docker Compose

### Start all services
```bash
make up
```

### View logs
```bash
make logs           # All services
make logs-api       # API only
make logs-worker    # Worker only
make logs-renderer  # Renderer only
```

### Stop services
```bash
make down
```

## Go Workspace

This monorepo uses Go Workspaces (`go.work`) to manage multiple Go modules:

- `pkg/` - Shared packages
- `is-api/` - API service
- `is-worker/` - Worker service

To sync the workspace:
```bash
make work-sync
```

## Database Migrations

Migrations are stored in `is-api/migrations/`.

```bash
make migrate-up      # Run migrations
make migrate-down    # Rollback last migration
```

## Environment Variables

Required environment variables (see `.env.example`):

- Database: `BLUEPRINT_DB_*`
- Redis: `REDIS_HOST`, `REDIS_PORT`
- OAuth: `GOOGLE_*`, `DISCORD_*`
- AWS: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `S3_BUCKET_NAME`
- AI APIs: `GOOGLE_API_KEY`, `ELEVENLABS_API_KEY`

## Architecture

### Shared Code (`pkg/`)
- `database/` - Database connection and service
- `queue/` - Redis queue client and task definitions

### Service-Specific Code
- `is-api/internal/` - API-specific handlers, routes, auth
- `is-worker/internal/` - Worker-specific AI services, storage

### Task Flow
1. API enqueues tasks → Redis
2. Worker processes generation tasks
3. Renderer processes render tasks
4. Worker handles completion tasks

## Make Targets

Run `make help` to see all available targets.

## License

ISC

