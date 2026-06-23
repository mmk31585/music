-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS playback_sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id        UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    device_id       VARCHAR(255),
    device_type     VARCHAR(50),
    ip_address      INET,
    user_agent      TEXT,

    started_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at        TIMESTAMPTZ,
    duration_played INTEGER NOT NULL DEFAULT 0,
    completion_rate REAL NOT NULL DEFAULT 0,

    source          VARCHAR(50) DEFAULT 'direct',
    session_data    JSONB NOT NULL DEFAULT '{}',

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_playback_sessions_user_id ON playback_sessions(user_id);
CREATE INDEX idx_playback_sessions_track_id ON playback_sessions(track_id);
CREATE INDEX idx_playback_sessions_started_at ON playback_sessions(started_at);
CREATE INDEX idx_playback_sessions_user_track ON playback_sessions(user_id, track_id);
CREATE INDEX idx_playback_sessions_source ON playback_sessions(source);

-- Partition by month for performance
-- Requires pg_partman or manual partition management for production
-- CREATE TABLE playback_sessions_y2026m01 PARTITION OF playback_sessions
--     FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS playback_sessions;

-- +goose StatementEnd
