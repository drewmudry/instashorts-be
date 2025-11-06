import pg from 'pg';
import dotenv from 'dotenv';

dotenv.config();

const { Pool } = pg;

export interface VideoData {
  id: number;
  user_id: number;
  audio_url: string | null;
  captions: string | null; // JSON string
  video_url: string | null;
  status: string;
}

export interface VideoScene {
  id: number;
  video_id: number;
  image_url: string | null;
  index: number;
  prompt: string;
  status: string;
}

export interface Caption {
  word: string;
  start_time: number;
  end_time: number;
}

let pool: pg.Pool | null = null;

export function getDatabasePool(): pg.Pool {
  if (!pool) {
    pool = new Pool({
      host: process.env.BLUEPRINT_DB_HOST || 'localhost',
      port: parseInt(process.env.BLUEPRINT_DB_PORT || '5432'),
      database: process.env.BLUEPRINT_DB_DATABASE || '',
      user: process.env.BLUEPRINT_DB_USERNAME || '',
      password: process.env.BLUEPRINT_DB_PASSWORD || '',
      // Note: schema is set via search_path in query execution
    });
  }
  return pool;
}

export async function fetchVideoData(videoId: number): Promise<VideoData | null> {
  const db = getDatabasePool();
  const result = await db.query<VideoData>(
    'SELECT id, user_id, audio_url, captions, video_url, status FROM videos WHERE id = $1',
    [videoId]
  );
  return result.rows[0] || null;
}

export async function fetchVideoScenes(videoId: number): Promise<VideoScene[]> {
  const db = getDatabasePool();
  const result = await db.query<VideoScene>(
    'SELECT id, video_id, image_url, index, prompt, status FROM video_scenes WHERE video_id = $1 ORDER BY index ASC',
    [videoId]
  );
  return result.rows;
}

export async function updateVideoStatus(videoId: number, status: string): Promise<void> {
  const db = getDatabasePool();
  await db.query('UPDATE videos SET status = $1, updated_at = NOW() WHERE id = $2', [
    status,
    videoId,
  ]);
}

export async function updateVideoUrl(videoId: number, videoUrl: string): Promise<void> {
  const db = getDatabasePool();
  await db.query(
    'UPDATE videos SET video_url = $1, status = $2, updated_at = NOW() WHERE id = $3',
    [videoUrl, 'completed', videoId]
  );
}

export async function closeDatabase(): Promise<void> {
  if (pool) {
    await pool.end();
    pool = null;
  }
}
