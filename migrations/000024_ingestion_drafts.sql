-- +goose Up
-- +goose StatementBegin

CREATE TABLE ingestion_drafts (
    id TEXT PRIMARY KEY,
    uploaded_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    original_filename TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT NOT NULL DEFAULT 0,
    duration_seconds DOUBLE PRECISION,
    bitrate INTEGER,
    format TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'enriching', 'review', 'accepted', 'rejected')),
    extracted_metadata JSONB NOT NULL DEFAULT '{}',
    enriched_metadata JSONB,
    final_metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ingestion_drafts_uploaded_by ON ingestion_drafts(uploaded_by);
CREATE INDEX idx_ingestion_drafts_status ON ingestion_drafts(status);
CREATE INDEX idx_ingestion_drafts_created_at ON ingestion_drafts(created_at DESC);

CREATE TABLE ingestion_draft_assets (
    id TEXT PRIMARY KEY,
    draft_id TEXT NOT NULL REFERENCES ingestion_drafts(id) ON DELETE CASCADE,
    asset_type TEXT NOT NULL CHECK (asset_type IN ('cover', 'audio')),
    url TEXT NOT NULL,
    source TEXT NOT NULL CHECK (source IN ('embedded', 'uploaded', 'fetched')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ingestion_draft_assets_draft_id ON ingestion_draft_assets(draft_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS ingestion_draft_assets;
DROP TABLE IF EXISTS ingestion_drafts;

-- +goose StatementEnd
