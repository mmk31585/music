-- +goose Up
-- +goose StatementBegin

ALTER TABLE ingestion_draft_assets
DROP CONSTRAINT IF EXISTS ingestion_draft_assets_asset_type_check;

ALTER TABLE ingestion_draft_assets
ADD CONSTRAINT ingestion_draft_assets_asset_type_check
CHECK (asset_type IN ('cover', 'track_cover', 'album_cover', 'artist_image', 'audio'));

ALTER TABLE ingestion_draft_assets
DROP CONSTRAINT IF EXISTS ingestion_draft_assets_source_check;

ALTER TABLE ingestion_draft_assets
ADD CONSTRAINT ingestion_draft_assets_source_check
CHECK (source IN ('embedded', 'uploaded', 'fetched', 'enrichment'));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE ingestion_draft_assets
DROP CONSTRAINT IF EXISTS ingestion_draft_assets_asset_type_check;

ALTER TABLE ingestion_draft_assets
ADD CONSTRAINT ingestion_draft_assets_asset_type_check
CHECK (asset_type IN ('cover', 'audio'));

ALTER TABLE ingestion_draft_assets
DROP CONSTRAINT IF EXISTS ingestion_draft_assets_source_check;

ALTER TABLE ingestion_draft_assets
ADD CONSTRAINT ingestion_draft_assets_source_check
CHECK (source IN ('embedded', 'uploaded', 'fetched'));

-- +goose StatementEnd
