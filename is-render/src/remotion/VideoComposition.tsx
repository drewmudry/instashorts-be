import React from 'react';
import {
  AbsoluteFill,
  interpolate,
  useCurrentFrame,
  useVideoConfig,
  Img,
  Audio,
  Sequence,
  staticFile,
} from 'remotion';

interface VideoCompositionProps {
  scenes: Array<{
    imageUrl: string;
    index: number;
  }>;
  captions: Array<{
    word: string;
    startTime: number;
    endTime: number;
  }>;
  audioUrl: string;
}

export const VideoComposition: React.FC<VideoCompositionProps> = ({
  scenes,
  captions,
  audioUrl,
}) => {
  const frame = useCurrentFrame();
  const { fps, durationInFrames } = useVideoConfig();

  // Calculate current time in seconds
  const currentTime = frame / fps;

  // Find current scene based on time
  // Assuming scenes are evenly distributed across the video duration
  const sceneDuration = durationInFrames / scenes.length;
  const currentSceneIndex = Math.floor(frame / sceneDuration);
  const currentScene = scenes[currentSceneIndex] || scenes[scenes.length - 1];

  // Find current caption words
  const currentCaptions = captions.filter(
    (caption) => currentTime >= caption.startTime && currentTime <= caption.endTime
  );

  // Calculate scene transition
  const sceneProgress = (frame % sceneDuration) / sceneDuration;

  return (
    <AbsoluteFill style={{ backgroundColor: '#000' }}>
      {/* Audio */}
      <Audio src={audioUrl} />

      {/* Current Scene Image */}
      <AbsoluteFill
        style={{
          opacity: interpolate(sceneProgress, [0, 0.1, 0.9, 1], [0, 1, 1, 0]),
        }}
      >
        <Img
          src={currentScene.imageUrl}
          style={{
            width: '100%',
            height: '100%',
            objectFit: 'cover',
          }}
        />
      </AbsoluteFill>

      {/* Captions Overlay */}
      <AbsoluteFill
        style={{
          justifyContent: 'center',
          alignItems: 'center',
          padding: '40px',
        }}
      >
        <div
          style={{
            backgroundColor: 'rgba(0, 0, 0, 0.7)',
            padding: '20px 40px',
            borderRadius: '10px',
            maxWidth: '80%',
          }}
        >
          <div
            style={{
              fontSize: '48px',
              fontWeight: 'bold',
              color: '#fff',
              textAlign: 'center',
              lineHeight: '1.2',
            }}
          >
            {currentCaptions.map((caption, idx) => (
              <span
                key={idx}
                style={{
                  marginRight: '8px',
                  opacity: interpolate(
                    currentTime,
                    [caption.startTime, caption.startTime + 0.1, caption.endTime - 0.1, caption.endTime],
                    [0, 1, 1, 0],
                    {
                      extrapolateLeft: 'clamp',
                      extrapolateRight: 'clamp',
                    }
                  ),
                }}
              >
                {caption.word}
              </span>
            ))}
          </div>
        </div>
      </AbsoluteFill>
    </AbsoluteFill>
  );
};
