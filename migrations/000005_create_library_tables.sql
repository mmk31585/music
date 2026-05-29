-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS liked_tracks (
                                            id BIGSERIAL PRIMARY KEY,
                                            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_liked_tracks_user_track UNIQUE (user_id, track_id)
    );

CREATE TABLE IF NOT EXISTS liked_albums (
                                            id BIGSERIAL PRIMARY KEY,
                                            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    album_id UUID NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_liked_albums_user_album UNIQUE (user_id, album_id)
    );

CREATE TABLE IF NOT EXISTS followed_artists (
                                                id BIGSERIAL PRIMARY KEY,
                                                user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_followed_artists_user_artist UNIQUE (user_id, artist_id)
    );

CREATE TABLE IF NOT EXISTS play_history (
                                            id BIGSERIAL PRIMARY KEY,
                                            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    played_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_liked_tracks_user_id
    ON liked_tracks(user_id);

CREATE INDEX IF NOT EXISTS idx_liked_tracks_track_id
    ON liked_tracks(track_id);

CREATE INDEX IF NOT EXISTS idx_liked_albums_user_id
    ON liked_albums(user_id);

CREATE INDEX IF NOT EXISTS idx_liked_albums_album_id
    ON liked_albums(album_id);

CREATE INDEX IF NOT EXISTS idx_followed_artists_user_id
    ON followed_artists(user_id);

CREATE INDEX IF NOT EXISTS idx_followed_artists_artist_id
    ON followed_artists(artist_id);

CREATE INDEX IF NOT EXISTS idx_play_history_user_id
    ON play_history(user_id);

CREATE INDEX IF NOT EXISTS idx_play_history_track_id
    ON play_history(track_id);

CREATE INDEX IF NOT EXISTS idx_play_history_user_played_at
    ON play_history(user_id, played_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS play_history;
DROP TABLE IF EXISTS followed_artists;
DROP TABLE IF EXISTS liked_albums;
DROP TABLE IF EXISTS liked_tracks;

-- +goose StatementEnd
