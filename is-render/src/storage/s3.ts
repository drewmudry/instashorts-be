import { S3Client, PutObjectCommand } from '@aws-sdk/client-s3';
import { readFileSync } from 'fs';
import path from 'path';
import dotenv from 'dotenv';

dotenv.config();

let s3Client: S3Client | null = null;

function getS3Client(): S3Client {
  if (!s3Client) {
    s3Client = new S3Client({
      region: process.env.AWS_REGION || 'us-east-1',
      credentials: {
        accessKeyId: process.env.AWS_ACCESS_KEY_ID || '',
        secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY || '',
      },
    });
  }
  return s3Client;
}

export async function uploadVideo(localPath: string, videoId: number): Promise<string> {
  const client = getS3Client();
  const bucket = process.env.S3_BUCKET_NAME || '';

  if (!bucket) {
    throw new Error('S3_BUCKET_NAME environment variable is not set');
  }

  const key = `videos/${videoId}/rendered_video.mp4`;
  
  const fileContent = readFileSync(localPath);
  
  const command = new PutObjectCommand({
    Bucket: bucket,
    Key: key,
    Body: fileContent,
    ContentType: 'video/mp4',
  });

  await client.send(command);

  // Construct the S3 URL
  const region = process.env.AWS_REGION || 'us-east-1';
  const videoUrl = `https://${bucket}.s3.${region}.amazonaws.com/${key}`;

  return videoUrl;
}
