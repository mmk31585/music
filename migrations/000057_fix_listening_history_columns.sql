-- +goose Up
-- listening_history lost columns during migration 000043 (partitioning rewrite)
-- and migration 000053 (conflicting schema). Re-add the missing columns.
-- +goose StatementBegin

ALTER TABLE listening_history ADD COLUMN IF NOT EXISTS completed BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE listening_history ADD COLUMN IF NOT EXISTS session_id UUID;
ALTER TABLE listening_history ADD COLUMN IF NOT EXISTS played_duration_ms BIGINT NOT NULL DEFAULT 0;
ALTER TABLE listening_history ADD COLUMN IF NOT EXISTS track_duration_ms BIGINT NOT NULL DEFAULT 0;
ALTER TABLE listening_history ADD COLUMN IF NOT EXISTS completion_percent NUMERIC(5,4) NOT NULL DEFAULT 0;
ALTER TABLE listening_history ADD COLUMN IF NOT EXISTS signal_type TEXT NOT NULL DEFAULT '';
ALTER TABLE listening_history ADD COLUMN IF NOT EXISTS is_explicit_like BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_listening_history_signal_type ON listening_history(signal_type);
CREATE INDEX IF NOT EXISTS idx_listening_history_session ON listening_history(session_id);
CREATE INDEX IF NOT EXISTS idx_listening_history_completed ON listening_history(completed);
CREATE INDEX IF NOT EXISTS idx_listening_history_user_completed_played_at ON listening_history(user_id, completed, played_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_listening_history_user_completed_played_at;
DROP INDEX IF EXISTS idx_listening_history_completed;
DROP INDEX IF EXISTS idx_listening_history_session;
DROP INDEX IF EXISTS idx_listening_history_signal_type;

ALTER TABLE listening_history DROP COLUMN IF EXISTS is_explicit_like;
ALTER TABLE listening_history DROP COLUMN IF EXISTS signal_type;
ALTER TABLE listening_history DROP COLUMN IF EXISTS completion_percent;
ALTER TABLE listening_history DROP COLUMN IF EXISTS track_duration_ms;
ALTER TABLE listening_history DROP COLUMN IF EXISTS played_duration_ms;
ALTER TABLE listening_history DROP COLUMN IF EXISTS session_id;
ALTER TABLE listening_history DROP COLUMN IF EXISTS completed;

-- +goose StatementEnd
