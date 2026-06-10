-- +goose Up

-- Event outbox for reliable event publishing
CREATE TABLE event_outbox (
    id UUID PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ,
    retry_count INT NOT NULL DEFAULT 0,
    last_error TEXT,
    trace_id VARCHAR(64)
);

CREATE INDEX idx_event_outbox_status ON event_outbox(status, created_at);
CREATE INDEX idx_event_outbox_type ON event_outbox(event_type);
CREATE INDEX idx_event_outbox_created ON event_outbox(created_at);

-- Social: unified follow system
CREATE TABLE IF NOT EXISTS user_follows (
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followed_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower_id, followed_id)
);

CREATE INDEX idx_user_follows_follower ON user_follows(follower_id, created_at DESC);
CREATE INDEX idx_user_follows_followed ON user_follows(followed_id, created_at DESC);

-- Social: activity feed
CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    target_id VARCHAR(255) NOT NULL,
    target_type VARCHAR(50) NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activities_user ON activities(user_id, created_at DESC);
CREATE INDEX idx_activities_type ON activities(type, created_at DESC);
CREATE INDEX idx_activities_target ON activities(target_id, target_type);
CREATE INDEX idx_activities_created ON activities(created_at DESC);

-- Reactions: unified like/reaction system (replaces old liked_tracks, liked_albums)
CREATE TABLE reactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_id VARCHAR(255) NOT NULL,
    target_type VARCHAR(50) NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'like',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, target_id, target_type)
);

CREATE INDEX idx_reactions_target ON reactions(target_id, target_type);
CREATE INDEX idx_reactions_user ON reactions(user_id, target_type);
CREATE INDEX idx_reactions_count ON reactions(target_id, target_type, type);

-- Content moderation
CREATE TABLE content_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reporter_id UUID NOT NULL REFERENCES users(id),
    target_id VARCHAR(255) NOT NULL,
    target_type VARCHAR(50) NOT NULL,
    reason VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    moderator_id UUID REFERENCES users(id),
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reports_status ON content_reports(status);
CREATE INDEX idx_reports_target ON content_reports(target_id, target_type);

CREATE TABLE content_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    target_id VARCHAR(255) NOT NULL,
    target_type VARCHAR(50) NOT NULL,
    flag_type VARCHAR(50) NOT NULL,
    flagged_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX idx_flags_target ON content_flags(target_id, target_type);

-- Creator dashboard
CREATE TABLE creator_stats (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_plays BIGINT NOT NULL DEFAULT 0,
    unique_listeners BIGINT NOT NULL DEFAULT 0,
    total_followers INT NOT NULL DEFAULT 0,
    total_tracks INT NOT NULL DEFAULT 0,
    total_albums INT NOT NULL DEFAULT 0,
    total_playlists INT NOT NULL DEFAULT 0,
    estimated_revenue BIGINT NOT NULL DEFAULT 0,
    last_calculated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE creator_daily_stats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    plays INT NOT NULL DEFAULT 0,
    listeners INT NOT NULL DEFAULT 0,
    likes INT NOT NULL DEFAULT 0,
    follows INT NOT NULL DEFAULT 0,
    shares INT NOT NULL DEFAULT 0,
    revenue_cents INT NOT NULL DEFAULT 0,
    UNIQUE (user_id, date)
);

CREATE INDEX idx_creator_daily_user ON creator_daily_stats(user_id, date DESC);

-- User profiles (extensible)
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio TEXT DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS location VARCHAR(255) DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS website VARCHAR(500) DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS social_links JSONB DEFAULT '{}';
ALTER TABLE users ADD COLUMN IF NOT EXISTS preferences JSONB DEFAULT '{}';
ALTER TABLE users ADD COLUMN IF NOT EXISTS header_image_url TEXT DEFAULT '';

-- HLS streaming
ALTER TABLE tracks ADD COLUMN IF NOT EXISTS hls_path TEXT DEFAULT '';
ALTER TABLE tracks ADD COLUMN IF NOT EXISTS has_hls BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE tracks ADD COLUMN IF NOT EXISTS hls_qualities TEXT[] DEFAULT '{}';

CREATE TABLE transcode_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    source_path TEXT NOT NULL,
    output_path TEXT NOT NULL,
    qualities TEXT[] NOT NULL DEFAULT '{ "128k", "192k", "320k" }',
    error_message TEXT,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transcode_jobs_status ON transcode_jobs(status);
CREATE INDEX idx_transcode_jobs_track ON transcode_jobs(track_id);

-- AI features
CREATE TABLE track_embeddings (
    track_id UUID PRIMARY KEY REFERENCES tracks(id) ON DELETE CASCADE,
    embedding FLOAT[] NOT NULL,
    model_version VARCHAR(50) NOT NULL DEFAULT 'v1',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE track_moods (
    track_id UUID PRIMARY KEY REFERENCES tracks(id) ON DELETE CASCADE,
    mood_tags JSONB NOT NULL DEFAULT '[]',
    energy FLOAT NOT NULL DEFAULT 0.5,
    valence FLOAT NOT NULL DEFAULT 0.5,
    tempo FLOAT NOT NULL DEFAULT 120,
    danceability FLOAT NOT NULL DEFAULT 0.5,
    acousticness FLOAT NOT NULL DEFAULT 0.5,
    instrumentalness FLOAT NOT NULL DEFAULT 0,
    liveness FLOAT NOT NULL DEFAULT 0.5,
    speechiness FLOAT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Recommendation v2
CREATE TABLE user_affinities (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_id VARCHAR(255) NOT NULL,
    target_type VARCHAR(20) NOT NULL,
    score FLOAT NOT NULL DEFAULT 0,
    decay FLOAT NOT NULL DEFAULT 1.0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, target_id, target_type)
);

CREATE TABLE user_similarities (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    similar_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    similarity FLOAT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, similar_user_id)
);

CREATE TABLE track_similarities (
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    similar_track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    similarity FLOAT NOT NULL,
    method VARCHAR(20) NOT NULL DEFAULT 'content',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (track_id, similar_track_id)
);

CREATE TABLE user_sessions (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    session_vector FLOAT[] NOT NULL DEFAULT '{}',
    recent_tracks UUID[] NOT NULL DEFAULT '{}',
    current_affinity JSONB DEFAULT '{}',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Collaborative playlists
ALTER TABLE playlists ADD COLUMN IF NOT EXISTS is_collaborative BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE playlists ADD COLUMN IF NOT EXISTS snapshot_id VARCHAR(32) DEFAULT '';

CREATE TABLE playlist_collaborators (
    playlist_id UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission VARCHAR(20) NOT NULL DEFAULT 'edit',
    added_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (playlist_id, user_id)
);

CREATE INDEX idx_playlist_collab_user ON playlist_collaborators(user_id);

-- WebSocket sessions
CREATE TABLE ws_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    connected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    disconnected_at TIMESTAMPTZ,
    ip_address VARCHAR(45),
    user_agent TEXT
);

CREATE INDEX idx_ws_sessions_user ON ws_sessions(user_id);
CREATE INDEX idx_ws_sessions_active
    ON ws_sessions(user_id)
    WHERE disconnected_at IS NULL;

-- OpenSearch sync tracking
CREATE TABLE search_sync (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(20) NOT NULL,
    synced BOOLEAN NOT NULL DEFAULT false,
    synced_at TIMESTAMPTZ,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_search_sync_pending ON search_sync(synced, created_at);

-- Audio quality versions for HLS
CREATE TABLE audio_qualities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    quality VARCHAR(20) NOT NULL,
    bitrate INT NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (track_id, quality)
);

CREATE INDEX idx_audio_qualities_track ON audio_qualities(track_id);

-- +goose Down

DROP TABLE IF EXISTS audio_qualities CASCADE;
DROP TABLE IF EXISTS search_sync CASCADE;
DROP TABLE IF EXISTS ws_sessions CASCADE;
DROP TABLE IF EXISTS playlist_collaborators CASCADE;
DROP TABLE IF EXISTS user_sessions CASCADE;
DROP TABLE IF EXISTS track_similarities CASCADE;
DROP TABLE IF EXISTS user_similarities CASCADE;
DROP TABLE IF EXISTS user_affinities CASCADE;
DROP TABLE IF EXISTS track_moods CASCADE;
DROP TABLE IF EXISTS track_embeddings CASCADE;
DROP TABLE IF EXISTS transcode_jobs CASCADE;
DROP TABLE IF EXISTS creator_daily_stats CASCADE;
DROP TABLE IF EXISTS creator_stats CASCADE;
DROP TABLE IF EXISTS content_flags CASCADE;
DROP TABLE IF EXISTS content_reports CASCADE;
DROP TABLE IF EXISTS reactions CASCADE;
DROP TABLE IF EXISTS activities CASCADE;
DROP TABLE IF EXISTS user_follows CASCADE;
DROP TABLE IF EXISTS event_outbox CASCADE;

ALTER TABLE users DROP COLUMN IF EXISTS bio;
ALTER TABLE users DROP COLUMN IF EXISTS location;
ALTER TABLE users DROP COLUMN IF EXISTS website;
ALTER TABLE users DROP COLUMN IF EXISTS social_links;
ALTER TABLE users DROP COLUMN IF EXISTS preferences;
ALTER TABLE users DROP COLUMN IF EXISTS header_image_url;
ALTER TABLE tracks DROP COLUMN IF EXISTS hls_path;
ALTER TABLE tracks DROP COLUMN IF EXISTS has_hls;
ALTER TABLE tracks DROP COLUMN IF EXISTS hls_qualities;
ALTER TABLE playlists DROP COLUMN IF EXISTS is_collaborative;
ALTER TABLE playlists DROP COLUMN IF EXISTS snapshot_id;

