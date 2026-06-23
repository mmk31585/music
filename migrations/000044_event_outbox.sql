-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS event_outbox (
    id            TEXT        PRIMARY KEY,
    event_type    TEXT        NOT NULL,
    payload       JSONB       NOT NULL DEFAULT '{}',
    status        TEXT        NOT NULL DEFAULT 'pending'
                              CHECK (status IN ('pending', 'processed', 'failed')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at  TIMESTAMPTZ,
    retry_count   INTEGER     NOT NULL DEFAULT 0,
    last_error    TEXT,
    trace_id      TEXT
);

-- Index for the relay's FetchPending query
CREATE INDEX IF NOT EXISTS idx_event_outbox_status_created
    ON event_outbox (status, created_at ASC)
    WHERE status = 'pending';

-- Periodic cleanup: remove processed events older than 7 days
CREATE INDEX IF NOT EXISTS idx_event_outbox_cleanup
    ON event_outbox (processed_at)
    WHERE status = 'processed';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS event_outbox;

-- +goose StatementEnd
