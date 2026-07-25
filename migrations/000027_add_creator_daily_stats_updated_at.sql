-- +goose Up
ALTER TABLE creator_daily_stats ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- +goose Down
ALTER TABLE creator_daily_stats DROP COLUMN IF EXISTS updated_at;
