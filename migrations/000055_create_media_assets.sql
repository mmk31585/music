-- +goose Up
-- +goose StatementBegin

-- The media_assets table was already created in 000039 with a different
-- schema (object_key, file_size, etc. instead of storage_key, bytes, ...).
-- Add the track_id column if missing and create the index.
ALTER TABLE media_assets
    ADD COLUMN IF NOT EXISTS track_id UUID REFERENCES tracks(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_media_assets_track_id ON media_assets(track_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_media_assets_track_id;
ALTER TABLE media_assets DROP COLUMN IF EXISTS track_id;

-- +goose StatementEnd
