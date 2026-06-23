-- +goose Up
-- Enrich listening_history with signal classification columns
ALTER TABLE listening_history ADD COLUMN IF NOT EXISTS session_id UUID;
ALTER TABLE listening_history ADD COLUMN IF NOT EXISTS played_duration_ms BIGINT NOT NULL DEFAULT 0;
ALTER TABLE listening_history ADD COLUMN IF NOT EXISTS track_duration_ms BIGINT NOT NULL DEFAULT 0;
ALTER TABLE listening_history ADD COLUMN IF NOT EXISTS completion_percent NUMERIC(5,4) NOT NULL DEFAULT 0;
ALTER TABLE listening_history ADD COLUMN IF NOT EXISTS signal_type TEXT NOT NULL DEFAULT '';
ALTER TABLE listening_history ADD COLUMN IF NOT EXISTS is_explicit_like BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_listening_history_signal_type
    ON listening_history(signal_type);
CREATE INDEX IF NOT EXISTS idx_listening_history_session
    ON listening_history(session_id);

-- Taste profiles table for aggregated user preference data
CREATE TABLE IF NOT EXISTS taste_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    top_genre_ids UUID[] NOT NULL DEFAULT '{}',
    top_artist_ids UUID[] NOT NULL DEFAULT '{}',
    seed_track_ids UUID[] NOT NULL DEFAULT '{}',
    avg_tempo_preference NUMERIC(6,2),
    onboarding_genre_ids UUID[] NOT NULL DEFAULT '{}',
    last_computed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS taste_profiles;
ALTER TABLE listening_history DROP COLUMN IF EXISTS is_explicit_like;
ALTER TABLE listening_history DROP COLUMN IF EXISTS signal_type;
ALTER TABLE listening_history DROP COLUMN IF EXISTS completion_percent;
ALTER TABLE listening_history DROP COLUMN IF EXISTS track_duration_ms;
ALTER TABLE listening_history DROP COLUMN IF EXISTS played_duration_ms;
ALTER TABLE listening_history DROP COLUMN IF EXISTS session_id;
