-- +goose Up
-- +goose StatementBegin

ALTER TABLE ingestion_drafts
    ADD COLUMN IF NOT EXISTS artist_id UUID REFERENCES artists(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS album_id  UUID REFERENCES albums(id)  ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS track_id  UUID REFERENCES tracks(id)  ON DELETE SET NULL;

ALTER TABLE ingestion_drafts
    DROP CONSTRAINT IF EXISTS ingestion_drafts_status_check;

ALTER TABLE ingestion_drafts
    ADD CONSTRAINT ingestion_drafts_status_check
        CHECK (status IN ('pending', 'enriching', 'review', 'accepted', 'rejected', 'published'));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE ingestion_drafts
    DROP COLUMN IF EXISTS track_id,
    DROP COLUMN IF EXISTS album_id,
    DROP COLUMN IF EXISTS artist_id;

ALTER TABLE ingestion_drafts
    DROP CONSTRAINT IF EXISTS ingestion_drafts_status_check;

ALTER TABLE ingestion_drafts
    ADD CONSTRAINT ingestion_drafts_status_check
        CHECK (status IN ('pending', 'enriching', 'review', 'accepted', 'rejected'));

-- +goose StatementEnd
