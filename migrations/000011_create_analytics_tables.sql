-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS analytics_events (
                                                id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NULL REFERENCES users(id) ON DELETE SET NULL,

    event_type TEXT NOT NULL CHECK (
                                       event_type IN (
                                       'play',
                                       'pause',
                                       'skip',
                                       'completion',
                                       'search',
                                       'playlist_open'
                                                     )
    ),

    track_id UUID NULL REFERENCES tracks(id) ON DELETE SET NULL,
    artist_id UUID NULL REFERENCES artists(id) ON DELETE SET NULL,
    album_id UUID NULL REFERENCES albums(id) ON DELETE SET NULL,
    playlist_id UUID NULL REFERENCES playlists(id) ON DELETE SET NULL,

    query TEXT NULL,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_analytics_events_event_type
    ON analytics_events(event_type);

CREATE INDEX IF NOT EXISTS idx_analytics_events_user_id
    ON analytics_events(user_id);

CREATE INDEX IF NOT EXISTS idx_analytics_events_track_id
    ON analytics_events(track_id);

CREATE INDEX IF NOT EXISTS idx_analytics_events_artist_id
    ON analytics_events(artist_id);

CREATE INDEX IF NOT EXISTS idx_analytics_events_album_id
    ON analytics_events(album_id);

CREATE INDEX IF NOT EXISTS idx_analytics_events_playlist_id
    ON analytics_events(playlist_id);

CREATE INDEX IF NOT EXISTS idx_analytics_events_created_at
    ON analytics_events(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_analytics_events_query
    ON analytics_events USING GIN (to_tsvector('simple', COALESCE(query, '')));

CREATE TABLE IF NOT EXISTS track_analytics (
                                               track_id UUID PRIMARY KEY REFERENCES tracks(id) ON DELETE CASCADE,

    play_count BIGINT NOT NULL DEFAULT 0,
    pause_count BIGINT NOT NULL DEFAULT 0,
    skip_count BIGINT NOT NULL DEFAULT 0,
    completion_count BIGINT NOT NULL DEFAULT 0,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS artist_analytics (
                                                artist_id UUID PRIMARY KEY REFERENCES artists(id) ON DELETE CASCADE,

    play_count BIGINT NOT NULL DEFAULT 0,
    pause_count BIGINT NOT NULL DEFAULT 0,
    skip_count BIGINT NOT NULL DEFAULT 0,
    completion_count BIGINT NOT NULL DEFAULT 0,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS album_analytics (
                                               album_id UUID PRIMARY KEY REFERENCES albums(id) ON DELETE CASCADE,

    play_count BIGINT NOT NULL DEFAULT 0,
    pause_count BIGINT NOT NULL DEFAULT 0,
    skip_count BIGINT NOT NULL DEFAULT 0,
    completion_count BIGINT NOT NULL DEFAULT 0,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_track_analytics_play_count
    ON track_analytics(play_count DESC);

CREATE INDEX IF NOT EXISTS idx_artist_analytics_play_count
    ON artist_analytics(play_count DESC);

CREATE INDEX IF NOT EXISTS idx_album_analytics_play_count
    ON album_analytics(play_count DESC);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS album_analytics;
DROP TABLE IF EXISTS artist_analytics;
DROP TABLE IF EXISTS track_analytics;
DROP TABLE IF EXISTS analytics_events;

-- +goose StatementEnd
