-- +goose Up
-- Add slug, genre, playlist_id to music_clubs
ALTER TABLE music_clubs ADD COLUMN IF NOT EXISTS slug TEXT NOT NULL DEFAULT '';
ALTER TABLE music_clubs ADD COLUMN IF NOT EXISTS genre TEXT;
ALTER TABLE music_clubs ADD COLUMN IF NOT EXISTS playlist_id UUID REFERENCES playlists(id) ON DELETE SET NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_music_clubs_slug ON music_clubs(slug) WHERE slug != '';
CREATE INDEX IF NOT EXISTS idx_music_clubs_genre ON music_clubs(genre);

-- +goose Down
DROP INDEX IF EXISTS idx_music_clubs_genre;
DROP INDEX IF EXISTS idx_music_clubs_slug;
ALTER TABLE music_clubs DROP COLUMN IF EXISTS playlist_id;
ALTER TABLE music_clubs DROP COLUMN IF EXISTS genre;
ALTER TABLE music_clubs DROP COLUMN IF EXISTS slug;
