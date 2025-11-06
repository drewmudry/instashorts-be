import dotenv from 'dotenv';
import { listenForRenderVideoTasks, enqueueVideoComplete, closeRedis } from './queue/client.js';
import {
  fetchVideoData,
  fetchVideoScenes,
  updateVideoStatus,
  updateVideoUrl,
  closeDatabase,
  type Caption,
} from './database/client.js';
import { renderVideo } from './renderer/video-renderer.js';

dotenv.config();

async function handleRenderVideoTask(payload: { video_id: number }): Promise<void> {
  const { video_id } = payload;
  
  console.log(`\n=== Processing render_video task for video_id: ${video_id} ===`);
  
  try {
    // Update status to rendering
    await updateVideoStatus(video_id, 'rendering');
    
    // Fetch video data from database
    const video = await fetchVideoData(video_id);
    if (!video) {
      throw new Error(`Video not found: ${video_id}`);
    }
    
    if (!video.audio_url) {
      throw new Error(`Video has no audio_url: ${video_id}`);
    }
    
    if (!video.captions) {
      throw new Error(`Video has no captions: ${video_id}`);
    }
    
    // Parse captions
    let captions: Caption[];
    try {
      captions = JSON.parse(video.captions);
    } catch (err) {
      throw new Error(`Failed to parse captions: ${err}`);
    }
    
    // Fetch scenes
    const scenes = await fetchVideoScenes(video_id);
    if (scenes.length === 0) {
      throw new Error(`No scenes found for video: ${video_id}`);
    }
    
    // Validate all scenes have image URLs
    const invalidScenes = scenes.filter((s) => !s.image_url);
    if (invalidScenes.length > 0) {
      throw new Error(`Some scenes are missing image URLs: ${invalidScenes.map((s) => s.id).join(', ')}`);
    }
    
    console.log(`Video data fetched: ${scenes.length} scenes, ${captions.length} captions`);
    
    // Render video
    const videoUrl = await renderVideo({
      videoId: video_id,
      scenes,
      captions,
      audioUrl: video.audio_url,
    });
    
    console.log(`Video rendered successfully: ${videoUrl}`);
    
    // Update database with video URL
    await updateVideoUrl(video_id, videoUrl);
    
    // Enqueue video_complete task
    await enqueueVideoComplete({
      video_id,
      video_url: videoUrl,
    });
    
    console.log(`✓ Successfully completed render for video_id: ${video_id}`);
  } catch (error) {
    console.error(`✗ Error processing render_video task for video_id ${video_id}:`, error);
    
    // Update status to failed
    try {
      await updateVideoStatus(video_id, 'failed');
    } catch (updateErr) {
      console.error('Failed to update video status to failed:', updateErr);
    }
  }
}

async function main() {
  console.log('Starting Instashorts Renderer...');
  console.log('Listening for render_video tasks...');
  
  // Set up graceful shutdown
  const shutdown = async () => {
    console.log('\nShutting down gracefully...');
    await closeDatabase();
    await closeRedis();
    process.exit(0);
  };
  
  process.on('SIGINT', shutdown);
  process.on('SIGTERM', shutdown);
  
  // Start listening for tasks
  try {
    await listenForRenderVideoTasks(handleRenderVideoTask);
  } catch (error) {
    console.error('Fatal error:', error);
    await shutdown();
  }
}

main().catch((error) => {
  console.error('Unhandled error:', error);
  process.exit(1);
});
