-- +goose Up
-- +goose StatementBegin

-- Enable pg_trgm extension for fuzzy text search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- GIN indexes for full-text search on catalog tables
CREATE INDEX IF NOT EXISTS idx_tracks_title_trgm ON tracks USING GIN (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_artists_name_trgm ON artists USING GIN (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_albums_title_trgm ON albums USING GIN (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_playlists_name_trgm ON playlists USING GIN (name gin_trgm_ops);

-- Full-text search vector columns (PostgreSQL tsvector)
ALTER TABLE tracks ADD COLUMN IF NOT EXISTS search_vector tsvector
    GENERATED ALWAYS AS (
        to_tsvector('simple', coalesce(title, ''))
    ) STORED;

ALTER TABLE artists ADD COLUMN IF NOT EXISTS search_vector tsvector
    GENERATED ALWAYS AS (
        to_tsvector('simple', coalesce(name, '')) ||
        to_tsvector('simple', coalesce(bio, ''))
    ) STORED;

ALTER TABLE albums ADD COLUMN IF NOT EXISTS search_vector tsvector
    GENERATED ALWAYS AS (
        to_tsvector('simple', coalesce(title, ''))
    ) STORED;

-- GIN indexes on search vectors
CREATE INDEX IF NOT EXISTS idx_tracks_search_vector ON tracks USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS idx_artists_search_vector ON artists USING GIN (search_vector);
CREATE INDEX IF NOT EXISTS idx_albums_search_vector ON albums USING GIN (search_vector);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_tracks_title_trgm;
DROP INDEX IF EXISTS idx_artists_name_trgm;
DROP INDEX IF EXISTS idx_albums_title_trgm;
DROP INDEX IF EXISTS idx_playlists_name_trgm;

ALTER TABLE tracks DROP COLUMN IF EXISTS search_vector;
ALTER TABLE artists DROP COLUMN IF EXISTS search_vector;
ALTER TABLE albums DROP COLUMN IF EXISTS search_vector;

DROP EXTENSION IF EXISTS pg_trgm;

-- +goose StatementEnd
