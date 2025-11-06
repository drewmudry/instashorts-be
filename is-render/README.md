# Instashorts Renderer

TypeScript-based video renderer using Remotion for the Instashorts API. This service listens for `render_video` tasks from Redis (asynq queue) and renders videos using Remotion.

## Features

- Listens for `render_video` tasks from Redis (asynq queue)
- Fetches video data (scenes, captions, audio) from PostgreSQL database
- Renders videos using Remotion
- Uploads rendered videos to S3
- Enqueues `video_complete` tasks back to Redis for the worker to handle

## Setup

### Prerequisites

- Node.js 18+ 
- PostgreSQL database (shared with instashorts-api)
- Redis (shared with instashorts-api)
- AWS S3 access (for uploading rendered videos)

### Installation

```bash
npm install
```

### Environment Variables

Create a `.env` file based on `.env.example`:

```env
# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379

# PostgreSQL Database Configuration
BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=your_database
BLUEPRINT_DB_USERNAME=your_username
BLUEPRINT_DB_PASSWORD=your_password
BLUEPRINT_DB_SCHEMA=public

# S3 Configuration
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=us-east-1
S3_BUCKET_NAME=your_bucket_name

# Renderer Configuration
OUTPUT_DIR=./output
```

### Build

```bash
npm run build
```

### Run

```bash
npm start
```

### Development

```bash
npm run dev
```

## Architecture

### Task Flow

1. **Render Video Task Received**: Listens for `video:render` tasks from Redis
2. **Fetch Video Data**: Retrieves video data from PostgreSQL:
   - Audio URL (S3 link)
   - Captions (JSON)
   - Scenes with image URLs (S3 links)
3. **Render Video**: Uses Remotion to render the video with:
   - Scene images displayed sequentially
   - Captions overlaid on the video
   - Audio synchronized
4. **Upload to S3**: Uploads the rendered video to S3
5. **Update Database**: Updates the video record with the final video URL
6. **Enqueue Complete Task**: Adds a `video:complete` task to Redis for the worker to handle

### Remotion Composition

The video composition (`VideoComposition.tsx`) renders:
- Scene images that transition smoothly
- Captions that appear synchronized with the audio
- Audio track playing in the background

Video specs:
- Resolution: 1080x1920 (vertical/portrait)
- FPS: 30
- Duration: Calculated from captions or scene count

## Integration with Go API

This renderer works alongside the Go-based `instashorts-api`:

- **Task Consumption**: The renderer consumes `video:render` tasks that are enqueued by the Go API
- **Task Production**: The renderer enqueues `video:complete` tasks that the Go worker handles
- **Database**: Both services share the same PostgreSQL database
- **Queue**: Both services use the same Redis instance with asynq

## Notes on asynq Compatibility

The renderer implements compatibility with asynq's Redis task format:
- Tasks are stored in msgpack format in Redis
- The renderer decodes msgpack to extract task data
- Tasks are enqueued in a format compatible with asynq

**Important**: The renderer listens for `render_video` tasks on the `default` queue. To avoid conflicts with the Go worker, you may want to:
1. Use a dedicated queue (e.g., `renderer`) for render_video tasks
2. Update the Go API to enqueue render_video tasks to the dedicated queue
3. Update the queue key in the renderer's `queue/client.ts` to match

Alternatively, you can disable the render_video handler in the Go worker (as shown in `cmd/worker/main.go`).

## Troubleshooting

### Tasks not being received

- Check Redis connection settings
- Verify the queue key format matches asynq's expected format
- Check that tasks are being enqueued with the correct type (`video:render`)

### Rendering fails

- Verify all scene images have valid S3 URLs
- Check that captions are valid JSON
- Ensure audio URL is accessible
- Check Remotion bundle creation logs

### Database connection issues

- Verify PostgreSQL connection settings
- Check that the database schema matches the expected structure
- Ensure database user has proper permissions

## License

ISC
