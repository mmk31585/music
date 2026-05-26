-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS artists (
                                       id UUID PRIMARY KEY,
                                       name VARCHAR(200) NOT NULL,
    slug VARCHAR(220) NOT NULL UNIQUE,
    bio TEXT,
    image_url TEXT,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    monthly_listeners BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
    );

CREATE INDEX IF NOT EXISTS idx_artists_name ON artists (name);
CREATE INDEX IF NOT EXISTS idx_artists_monthly_listeners ON artists (monthly_listeners DESC);

CREATE TABLE IF NOT EXISTS albums (
                                      id UUID PRIMARY KEY,
                                      artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    title VARCHAR(250) NOT NULL,
    slug VARCHAR(270) NOT NULL,
    cover_url TEXT,
    release_date DATE,
    album_type VARCHAR(50) NOT NULL DEFAULT 'album',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT albums_artist_id_slug_key UNIQUE (artist_id, slug)
    );

CREATE INDEX IF NOT EXISTS idx_albums_artist_id ON albums (artist_id);
CREATE INDEX IF NOT EXISTS idx_albums_release_date ON albums (release_date DESC);

CREATE TABLE IF NOT EXISTS tracks (
                                      id UUID PRIMARY KEY,
                                      artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    album_id UUID REFERENCES albums(id) ON DELETE SET NULL,
    title VARCHAR(250) NOT NULL,
    slug VARCHAR(270) NOT NULL,
    duration_seconds INT NOT NULL DEFAULT 0,
    track_number INT,
    explicit BOOLEAN NOT NULL DEFAULT FALSE,
    audio_url TEXT,
    cover_url TEXT,
    play_count BIGINT NOT NULL DEFAULT 0,
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT tracks_artist_id_slug_key UNIQUE (artist_id, slug)
    );

CREATE INDEX IF NOT EXISTS idx_tracks_artist_id ON tracks (artist_id);
CREATE INDEX IF NOT EXISTS idx_tracks_album_id ON tracks (album_id);
CREATE INDEX IF NOT EXISTS idx_tracks_is_public ON tracks (is_public);
CREATE INDEX IF NOT EXISTS idx_tracks_play_count ON tracks (play_count DESC);

CREATE TABLE IF NOT EXISTS genres (
                                      id UUID PRIMARY KEY,
                                      name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(120) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_genres_name ON genres (name);

CREATE TABLE IF NOT EXISTS track_genres (
                                            track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    genre_id UUID NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (track_id, genre_id)
    );

CREATE INDEX IF NOT EXISTS idx_track_genres_track_id ON track_genres (track_id);
CREATE INDEX IF NOT EXISTS idx_track_genres_genre_id ON track_genres (genre_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS track_genres;
DROP TABLE IF EXISTS genres;
DROP TABLE IF EXISTS tracks;
DROP TABLE IF EXISTS albums;
DROP TABLE IF EXISTS artists;

-- +goose StatementEnd
