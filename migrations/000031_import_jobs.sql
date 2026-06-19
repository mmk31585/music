-- +goose Up
-- +goose StatementBegin

CREATE TABLE import_jobs (
    id              TEXT PRIMARY KEY,
    url             TEXT NOT NULL,
    source          TEXT NOT NULL DEFAULT '',
    query           TEXT NOT NULL DEFAULT '',
    uploaded_by     UUID NOT NULL REFERENCES users(id),
    status          TEXT NOT NULL DEFAULT 'queued'
                    CHECK (status IN ('queued','resolving','downloading','extracting','uploading','complete','failed')),
    progress        INTEGER DEFAULT 0,
    stage           TEXT DEFAULT '',
    draft_id        TEXT REFERENCES ingestion_drafts(id),
    error_message   TEXT DEFAULT '',
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_import_jobs_status ON import_jobs(status);
CREATE INDEX idx_import_jobs_user ON import_jobs(uploaded_by);
CREATE INDEX idx_import_jobs_created ON import_jobs(created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS import_jobs;

-- +goose StatementEnd
