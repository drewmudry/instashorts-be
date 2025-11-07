-- Create series table
CREATE TABLE IF NOT EXISTS series (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes on series
CREATE INDEX idx_series_user_id ON series(user_id);
CREATE INDEX idx_series_deleted_at ON series(deleted_at);

-- Create videos table
CREATE TABLE IF NOT EXISTS videos (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    series_id INTEGER REFERENCES series(id) ON DELETE SET NULL,
    title TEXT,
    theme VARCHAR(255) NOT NULL,
    voice_id VARCHAR(255) NOT NULL,
    script TEXT,
    audio_url TEXT,
    video_url TEXT,
    captions JSONB,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for videos
CREATE INDEX idx_videos_user_id ON videos(user_id);
CREATE INDEX idx_videos_series_id ON videos(series_id);
CREATE INDEX idx_videos_status ON videos(status);
CREATE INDEX idx_videos_deleted_at ON videos(deleted_at);

-- Create video_scenes table
CREATE TABLE IF NOT EXISTS video_scenes (
    id SERIAL PRIMARY KEY,
    video_id INTEGER NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    prompt TEXT NOT NULL,
    image_url TEXT,
    index INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for video_scenes
CREATE INDEX idx_video_scenes_video_id ON video_scenes(video_id);
CREATE INDEX idx_video_scenes_status ON video_scenes(status);
CREATE INDEX idx_video_scenes_deleted_at ON video_scenes(deleted_at);

