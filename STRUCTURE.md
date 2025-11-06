# Instashorts Backend - Monorepo Structure

```
instashorts-be/
├── .gitignore
├── .env.example              # Environment variables template
├── docker-compose.yml        # Unified docker compose for all services
├── go.work                   # Go workspace file
├── Makefile                  # Build and development commands
├── README.md                 # Main documentation
│
├── pkg/                      # Shared Go packages
│   ├── go.mod
│   ├── database/            # Database connection and service
│   │   ├── database.go
│   │   └── database_test.go
│   └── queue/               # Redis queue client and tasks
│       ├── client.go
│       ├── client_test.go
│       ├── tasks.go
│       ├── example_usage.go
│       └── integration_example.go
│
├── is-api/                   # API Service (Go)
│   ├── Dockerfile
│   ├── go.mod
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
│   ├── internal/
│   │   ├── auth/            # Authentication handlers
│   │   │   ├── handlers.go
│   │   │   ├── middleware.go
│   │   │   ├── models.go
│   │   │   ├── oauth_service.go
│   │   │   ├── repository.go
│   │   │   └── routes.go
│   │   ├── server/          # HTTP server setup
│   │   │   ├── server.go
│   │   │   ├── routes.go
│   │   │   └── routes_test.go
│   │   └── video/           # Video management
│   │       ├── handlers.go
│   │       ├── models.go
│   │       ├── repository.go
│   │       └── routes.go
│   └── migrations/          # Database migrations
│       ├── 000001_create_auth_tables.up.sql
│       ├── 000001_create_auth_tables.down.sql
│       ├── 000002_create_series_and_videos.up.sql
│       └── 000002_create_series_and_videos.down.sql
│
├── is-worker/                # Worker Service (Go)
│   ├── Dockerfile
│   ├── go.mod
│   ├── cmd/
│   │   └── worker/
│   │       └── main.go
│   └── internal/
│       ├── ai/              # AI services
│       │   ├── elevenlabs.go
│       │   ├── speechtotext.go
│       │   └── gemini/
│       │       └── service.go
│       └── storage/         # S3 storage
│           └── s3.go
│
└── is-render/                # Renderer Service (TypeScript)
    ├── Dockerfile
    ├── package.json
    ├── tsconfig.json
    ├── remotion.config.ts
    ├── README.md
    └── src/
        ├── index.ts         # Main entry point
        ├── database/
        │   └── client.ts    # PostgreSQL client
        ├── queue/
        │   └── client.ts    # Redis queue client
        ├── remotion/
        │   ├── Root.tsx
        │   └── VideoComposition.tsx
        ├── renderer/
        │   └── video-renderer.ts
        └── storage/
            └── s3.ts
```

## Module Dependencies

- **pkg/** - No dependencies on other modules
- **is-api/** - Depends on `pkg/`
- **is-worker/** - Depends on `pkg/`
- **is-render/** - Independent TypeScript module

## Import Paths

- `instashorts-be/pkg/database` - Database service
- `instashorts-be/pkg/queue` - Queue client and tasks
- `instashorts-be/is-api/internal/*` - API-specific code
- `instashorts-be/is-worker/internal/*` - Worker-specific code



