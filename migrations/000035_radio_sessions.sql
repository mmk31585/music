-- +goose Up
CREATE TABLE IF NOT EXISTS radio_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seed_track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    seed_type TEXT NOT NULL DEFAULT 'track' CHECK (seed_type IN ('track', 'artist')),
    played_track_ids JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_radio_sessions_user ON radio_sessions(user_id, last_active_at DESC);

-- +goose Down
DROP TABLE IF EXISTS radio_sessions;
