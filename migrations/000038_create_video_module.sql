-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    uploader_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL CHECK (type IN ('official_mv', 'user_edit')),
    status TEXT NOT NULL DEFAULT 'processing' CHECK (status IN ('processing','ready','failed')),
    title TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    raw_video_path TEXT,
    final_video_path TEXT,
    thumbnail_path TEXT,
    duration_ms BIGINT NOT NULL DEFAULT 0,
    aspect_ratio TEXT NOT NULL DEFAULT '',
    file_size_bytes BIGINT NOT NULL DEFAULT 0,
    track_start_ms BIGINT NOT NULL DEFAULT 0,
    track_end_ms BIGINT NOT NULL DEFAULT 0,
    view_count BIGINT NOT NULL DEFAULT 0,
    like_count BIGINT NOT NULL DEFAULT 0,
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    is_approved BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_videos_track_type ON videos(track_id, type, is_approved, is_public);
CREATE INDEX IF NOT EXISTS idx_videos_uploader ON videos(uploader_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_videos_explore ON videos(is_approved, is_public, created_at DESC);

CREATE TABLE IF NOT EXISTS video_likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(video_id, user_id)
);

CREATE TABLE IF NOT EXISTS track_like_visibility (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    visibility TEXT NOT NULL DEFAULT 'private' CHECK (visibility IN ('public','private')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, track_id)
);

CREATE TABLE IF NOT EXISTS user_music_status (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current_track_id UUID REFERENCES tracks(id) ON DELETE SET NULL,
    visibility TEXT NOT NULL DEFAULT 'private' CHECK (visibility IN ('public','followers','private')),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS user_music_status;
DROP TABLE IF EXISTS track_like_visibility;
DROP TABLE IF EXISTS video_likes;
DROP TABLE IF EXISTS videos;

-- +goose StatementEnd
