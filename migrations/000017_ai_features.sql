-- +goose Up
-- +goose StatementBegin

-- AI generation log for tracking playlist generation requests
CREATE TABLE IF NOT EXISTS ai_generation_log (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    playlist_id UUID REFERENCES playlists(id) ON DELETE SET NULL,
    prompt TEXT NOT NULL DEFAULT '',
    track_count INT NOT NULL DEFAULT 0,
    model_used VARCHAR(100) NOT NULL DEFAULT 'v1',
    latency_ms INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_generation_log_user_id ON ai_generation_log(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_generation_log_created_at ON ai_generation_log(created_at);

-- Smart playlists for saved AI playlist configurations
CREATE TABLE IF NOT EXISTS smart_playlists (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    query_config TEXT NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_smart_playlists_user_id ON smart_playlists(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS smart_playlists CASCADE;
DROP TABLE IF EXISTS ai_generation_log CASCADE;
-- +goose StatementEnd
