import { Composition } from 'remotion';
import { VideoComposition } from './VideoComposition';

export const RemotionRoot: React.FC = () => {
  return (
    <>
      <Composition
        id="VideoComposition"
        component={VideoComposition}
        durationInFrames={2100} // 70 seconds at 30fps
        fps={30}
        width={1080}
        height={1920}
        defaultProps={{
          scenes: [],
          captions: [],
          audioUrl: '',
        }}
      />
    </>
  );
};
