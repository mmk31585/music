-- +goose Up
-- +goose StatementBegin

ALTER TABLE ingestion_drafts
    ADD COLUMN IF NOT EXISTS file_hash TEXT;

ALTER TABLE ingestion_drafts
    ADD COLUMN IF NOT EXISTS stale BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE ingestion_drafts
    DROP CONSTRAINT IF EXISTS ingestion_drafts_status_check;

ALTER TABLE ingestion_drafts
    ADD CONSTRAINT ingestion_drafts_status_check
        CHECK (status IN ('pending', 'enriching', 'review', 'accepted', 'rejected', 'published', 'enrichment_failed'));

CREATE INDEX IF NOT EXISTS idx_ingestion_drafts_file_hash ON ingestion_drafts(file_hash);
CREATE INDEX IF NOT EXISTS idx_ingestion_drafts_stale ON ingestion_drafts(stale) WHERE stale = TRUE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE ingestion_drafts
    DROP COLUMN IF EXISTS file_hash,
    DROP COLUMN IF EXISTS stale;

ALTER TABLE ingestion_drafts
    DROP CONSTRAINT IF EXISTS ingestion_drafts_status_check;

ALTER TABLE ingestion_drafts
    ADD CONSTRAINT ingestion_drafts_status_check
        CHECK (status IN ('pending', 'enriching', 'review', 'accepted', 'rejected', 'published'));

-- +goose StatementEnd
