-- +goose Up
-- +goose StatementBegin

-- Add audio_asset_id and cover_asset_id to tracks, linking to media_assets.
-- This is the new FK reference replacing the old audio_media_id/cover_media_id
-- columns (from migration 000013) that pointed to the legacy `media` table.
-- Using NOT VALID so existing data isn't rechecked (zero-downtime pattern).

ALTER TABLE tracks
    ADD COLUMN IF NOT EXISTS audio_asset_id UUID REFERENCES media_assets(id) ON DELETE SET NULL;

ALTER TABLE tracks
    ADD COLUMN IF NOT EXISTS cover_asset_id UUID REFERENCES media_assets(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_tracks_audio_asset_id ON tracks(audio_asset_id);
CREATE INDEX IF NOT EXISTS idx_tracks_cover_asset_id ON tracks(cover_asset_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_tracks_cover_asset_id;
DROP INDEX IF EXISTS idx_tracks_audio_asset_id;

ALTER TABLE tracks DROP COLUMN IF EXISTS cover_asset_id;
ALTER TABLE tracks DROP COLUMN IF EXISTS audio_asset_id;

-- +goose StatementEnd
