import { renderMedia, selectComposition } from '@remotion/renderer';
import { bundle } from '@remotion/bundler';
import path from 'path';
import { fileURLToPath } from 'url';
import { dirname } from 'path';
import { existsSync, mkdirSync } from 'fs';
import { VideoComposition } from '../remotion/VideoComposition.js';
import type { Caption, VideoScene } from '../database/client.js';
import * as S3 from '../storage/s3.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

export interface RenderVideoParams {
  videoId: number;
  scenes: VideoScene[];
  captions: Caption[];
  audioUrl: string;
}

export async function renderVideo(params: RenderVideoParams): Promise<string> {
  const { videoId, scenes, captions, audioUrl } = params;

  console.log(`Starting render for video_id: ${videoId}`);
  console.log(`Scenes: ${scenes.length}, Captions: ${captions.length}`);

  // Calculate video duration from captions
  let videoDuration = 0;
  if (captions.length > 0) {
    videoDuration = captions[captions.length - 1].endTime;
  } else {
    // Fallback: 5 seconds per scene
    videoDuration = scenes.length * 5;
  }

  const fps = 30;
  const durationInFrames = Math.ceil(videoDuration * fps);

  // Bundle Remotion app
  console.log('Bundling Remotion app...');
  const entryPoint = path.resolve(__dirname, '../remotion/Root.tsx');
  const bundleLocation = await bundle({
    entryPoint,
    webpackOverride: (config) => config,
  });

  console.log('Bundle created at:', bundleLocation);

  // Select composition
  const composition = await selectComposition({
    serveUrl: bundleLocation,
    id: 'VideoComposition',
    inputProps: {
      scenes: scenes.map((s) => ({
        imageUrl: s.image_url || '',
        index: s.index,
      })),
      captions: captions,
      audioUrl: audioUrl,
    },
  });

  // Output path
  const outputDir = process.env.OUTPUT_DIR || './output';
  if (!existsSync(outputDir)) {
    mkdirSync(outputDir, { recursive: true });
  }
  const outputPath = path.join(outputDir, `video_${videoId}.mp4`);

  console.log(`Rendering video to: ${outputPath}`);

  // Render video
  await renderMedia({
    composition,
    serveUrl: bundleLocation,
    codec: 'h264',
    outputLocation: outputPath,
    inputProps: {
      scenes: scenes.map((s) => ({
        imageUrl: s.image_url || '',
        index: s.index,
      })),
      captions: captions,
      audioUrl: audioUrl,
    },
  });

  console.log(`Video rendered successfully: ${outputPath}`);

  // Upload to S3
  console.log('Uploading video to S3...');
  const videoUrl = await S3.uploadVideo(outputPath, videoId);

  console.log(`Video uploaded to S3: ${videoUrl}`);

  return videoUrl;
}
