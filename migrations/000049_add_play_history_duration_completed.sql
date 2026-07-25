-- +goose Up
-- +goose StatementBegin

ALTER TABLE play_history
    ADD COLUMN IF NOT EXISTS duration INTEGER,
    ADD COLUMN IF NOT EXISTS completed BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE play_history
    DROP COLUMN IF EXISTS duration,
    DROP COLUMN IF EXISTS completed;

-- +goose StatementEnd
