-- +goose Up
-- +goose StatementBegin

ALTER TABLE lyrics
    ADD COLUMN IF NOT EXISTS source VARCHAR(20) NOT NULL DEFAULT 'manual',
    ADD COLUMN IF NOT EXISTS confidence_score REAL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE lyrics
    DROP COLUMN IF EXISTS source,
    DROP COLUMN IF EXISTS confidence_score;

-- +goose StatementEnd
