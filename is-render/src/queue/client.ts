import Redis from 'ioredis';
import dotenv from 'dotenv';
import { unpack } from 'msgpackr';

dotenv.config();

// Task types matching the Go API
export const TASK_TYPES = {
  RENDER_VIDEO: 'video:render',
  VIDEO_COMPLETE: 'video:complete',
} as const;

export interface RenderVideoPayload {
  video_id: number;
}

export interface VideoCompletePayload {
  video_id: number;
  video_url: string;
}

let redisClient: Redis | null = null;
let redisSubscriber: Redis | null = null;

export function getRedisClient(): Redis {
  if (!redisClient) {
    const host = process.env.REDIS_HOST || 'localhost';
    const port = parseInt(process.env.REDIS_PORT || '6379');
    redisClient = new Redis({
      host,
      port,
      retryStrategy: (times) => {
        const delay = Math.min(times * 50, 2000);
        return delay;
      },
    });

    redisClient.on('error', (err) => {
      console.error('Redis Client Error:', err);
    });

    redisClient.on('connect', () => {
      console.log('Redis Client Connected');
    });
  }
  return redisClient;
}

export function getRedisSubscriber(): Redis {
  if (!redisSubscriber) {
    const host = process.env.REDIS_HOST || 'localhost';
    const port = parseInt(process.env.REDIS_PORT || '6379');
    redisSubscriber = new Redis({
      host,
      port,
      retryStrategy: (times) => {
        const delay = Math.min(times * 50, 2000);
        return delay;
      },
    });

    redisSubscriber.on('error', (err) => {
      console.error('Redis Subscriber Error:', err);
    });

    redisSubscriber.on('connect', () => {
      console.log('Redis Subscriber Connected');
    });
  }
  return redisSubscriber;
}

// asynq stores tasks in Redis with specific keys
// Tasks are stored in sorted sets and lists with patterns like:
// asynq:queues:{queue}:pending, asynq:queues:{queue}:active, etc.
// We'll use BLPOP on the pending queue to get tasks
export async function listenForRenderVideoTasks(
  callback: (payload: RenderVideoPayload) => Promise<void>
): Promise<void> {
  const client = getRedisClient();
  
  // asynq stores pending tasks in a list: asynq:queues:{queue}:pending
  // The default queue is "default"
  const queueKey = 'asynq:queues:default:pending';
  
  console.log(`Listening for render_video tasks on queue: ${queueKey}`);
  
  while (true) {
    try {
      // Use BLPOP to block until a task is available
      // BLPOP returns [key, value] when a task is available
      const result = await client.blpop(queueKey, 0);
      
      if (result && result.length >= 2) {
        const [, taskData] = result;
        try {
          // asynq stores tasks as msgpack (binary format)
          // Try to decode as msgpack first, then fall back to JSON
          let task: any;
          
          try {
            // Try msgpack first (asynq's default format)
            const buffer = Buffer.from(taskData, 'binary');
            task = unpack(buffer);
          } catch {
            try {
              // Fall back to JSON (for backwards compatibility or custom formats)
              task = JSON.parse(taskData);
            } catch {
              console.warn('Task is not in msgpack or JSON format, skipping');
              continue;
            }
          }
          
          // asynq task structure: {Type: string, Payload: []byte}
          // Payload is stored as a byte array that needs to be parsed as JSON
          if (task && task.Type === TASK_TYPES.RENDER_VIDEO && task.Payload) {
            let payload: RenderVideoPayload;
            
            // Payload might be a Buffer or Uint8Array
            if (Buffer.isBuffer(task.Payload) || task.Payload instanceof Uint8Array) {
              const payloadStr = Buffer.from(task.Payload).toString('utf-8');
              payload = JSON.parse(payloadStr);
            } else if (typeof task.Payload === 'string') {
              payload = JSON.parse(task.Payload);
            } else {
              payload = task.Payload;
            }
            
            console.log(`Received render_video task for video_id: ${payload.video_id}`);
            await callback(payload);
          }
        } catch (err) {
          console.error('Error processing task:', err);
        }
      }
    } catch (err) {
      console.error('Error listening for tasks:', err);
      // Wait a bit before retrying
      await new Promise((resolve) => setTimeout(resolve, 1000));
    }
  }
}

// Enqueue a video_complete task
// asynq expects tasks in a specific format in the pending queue
export async function enqueueVideoComplete(payload: VideoCompletePayload): Promise<void> {
  const client = getRedisClient();
  
  // Create task in asynq format
  // asynq stores tasks with: Type, Payload (as JSON string), and metadata
  const task = {
    Type: TASK_TYPES.VIDEO_COMPLETE,
    Payload: JSON.stringify(payload),
    // Add minimal metadata that asynq expects
    Timeout: 0,
    Retry: 0,
    Queue: 'default',
  };
  
  // Push to the pending queue
  const queueKey = 'asynq:queues:default:pending';
  await client.rpush(queueKey, JSON.stringify(task));
  
  console.log(`Enqueued video_complete task for video_id: ${payload.video_id}`);
}

export async function closeRedis(): Promise<void> {
  if (redisClient) {
    await redisClient.quit();
    redisClient = null;
  }
  if (redisSubscriber) {
    await redisSubscriber.quit();
    redisSubscriber = null;
  }
}
