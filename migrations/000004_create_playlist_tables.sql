-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS playlists (
                                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    cover_url TEXT NULL,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS playlist_tracks (
                                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    playlist_id UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    position INT NOT NULL CHECK (position > 0),
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_playlist_track UNIQUE (playlist_id, track_id),
    CONSTRAINT uq_playlist_position UNIQUE (playlist_id, position)
    );

CREATE INDEX IF NOT EXISTS idx_playlists_user_id
    ON playlists(user_id);

CREATE INDEX IF NOT EXISTS idx_playlists_is_public
    ON playlists(is_public);

CREATE INDEX IF NOT EXISTS idx_playlist_tracks_playlist_id
    ON playlist_tracks(playlist_id);

CREATE INDEX IF NOT EXISTS idx_playlist_tracks_track_id
    ON playlist_tracks(track_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS playlist_tracks;
DROP TABLE IF EXISTS playlists;

-- +goose StatementEnd
