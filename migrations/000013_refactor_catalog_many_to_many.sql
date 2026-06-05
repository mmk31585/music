-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Fix existing UUID defaults if missing.
ALTER TABLE artists
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE albums
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE tracks
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE genres
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- Central media table.
CREATE TABLE IF NOT EXISTS media (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    kind VARCHAR(50) NOT NULL,
    url TEXT NOT NULL,

    mime_type VARCHAR(120),
    size_bytes BIGINT,

    width INT,
    height INT,

    duration_seconds INT,
    bitrate INT,
    sample_rate INT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
    );

CREATE INDEX IF NOT EXISTS idx_media_kind ON media (kind);
CREATE INDEX IF NOT EXISTS idx_media_created_at ON media (created_at DESC);

-- Artist media/profile expansion.
ALTER TABLE artists
    ADD COLUMN IF NOT EXISTS avatar_media_id UUID REFERENCES media(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS banner_media_id UUID REFERENCES media(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS country VARCHAR(2);

CREATE INDEX IF NOT EXISTS idx_artists_slug ON artists (slug);
CREATE INDEX IF NOT EXISTS idx_artists_avatar_media_id ON artists (avatar_media_id);
CREATE INDEX IF NOT EXISTS idx_artists_banner_media_id ON artists (banner_media_id);

-- Album media expansion.
ALTER TABLE albums
    ADD COLUMN IF NOT EXISTS cover_media_id UUID REFERENCES media(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_albums_cover_media_id ON albums (cover_media_id);

-- Track media expansion.
ALTER TABLE tracks
    ADD COLUMN IF NOT EXISTS audio_media_id UUID REFERENCES media(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS cover_media_id UUID REFERENCES media(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_tracks_audio_media_id ON tracks (audio_media_id);
CREATE INDEX IF NOT EXISTS idx_tracks_cover_media_id ON tracks (cover_media_id);

-- Many-to-many track artists.
CREATE TABLE IF NOT EXISTS track_artists (
                                             track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,

    role VARCHAR(50) NOT NULL DEFAULT 'primary',
    position INT NOT NULL DEFAULT 1,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (track_id, artist_id, role)
    );

CREATE INDEX IF NOT EXISTS idx_track_artists_track_id ON track_artists (track_id);
CREATE INDEX IF NOT EXISTS idx_track_artists_artist_id ON track_artists (artist_id);
CREATE INDEX IF NOT EXISTS idx_track_artists_role ON track_artists (role);
CREATE INDEX IF NOT EXISTS idx_track_artists_position ON track_artists (position);

-- Many-to-many album artists.
CREATE TABLE IF NOT EXISTS album_artists (
                                             album_id UUID NOT NULL REFERENCES albums(id) ON DELETE CASCADE,
    artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,

    role VARCHAR(50) NOT NULL DEFAULT 'primary',
    position INT NOT NULL DEFAULT 1,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (album_id, artist_id, role)
    );

CREATE INDEX IF NOT EXISTS idx_album_artists_album_id ON album_artists (album_id);
CREATE INDEX IF NOT EXISTS idx_album_artists_artist_id ON album_artists (artist_id);
CREATE INDEX IF NOT EXISTS idx_album_artists_role ON album_artists (role);
CREATE INDEX IF NOT EXISTS idx_album_artists_position ON album_artists (position);

-- Optional structured credits beyond display artists.
CREATE TABLE IF NOT EXISTS track_credits (
                                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    artist_id UUID REFERENCES artists(id) ON DELETE SET NULL,

    name VARCHAR(200),
    role VARCHAR(80) NOT NULL,

    position INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT track_credits_artist_or_name_check
    CHECK (artist_id IS NOT NULL OR name IS NOT NULL)
    );

CREATE INDEX IF NOT EXISTS idx_track_credits_track_id ON track_credits (track_id);
CREATE INDEX IF NOT EXISTS idx_track_credits_artist_id ON track_credits (artist_id);
CREATE INDEX IF NOT EXISTS idx_track_credits_role ON track_credits (role);

-- Backfill current one-artist design into new join tables.
INSERT INTO track_artists (track_id, artist_id, role, position)
SELECT id, artist_id, 'primary', 1
FROM tracks
    ON CONFLICT DO NOTHING;

INSERT INTO album_artists (album_id, artist_id, role, position)
SELECT id, artist_id, 'primary', 1
FROM albums
    ON CONFLICT DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS track_credits;
DROP TABLE IF EXISTS album_artists;
DROP TABLE IF EXISTS track_artists;

ALTER TABLE tracks
DROP COLUMN IF EXISTS audio_media_id,
    DROP COLUMN IF EXISTS cover_media_id;

ALTER TABLE albums
DROP COLUMN IF EXISTS cover_media_id;

ALTER TABLE artists
DROP COLUMN IF EXISTS avatar_media_id,
    DROP COLUMN IF EXISTS banner_media_id,
    DROP COLUMN IF EXISTS country;

DROP TABLE IF EXISTS media;

-- Keep UUID defaults on rollback intentionally.
-- They fix the original insert behavior and are safe to preserve.

-- +goose StatementEnd
