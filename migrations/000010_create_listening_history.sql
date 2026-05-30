-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS listening_history (
                                                 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    played_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- duration listened in seconds
    duration INTEGER NOT NULL DEFAULT 0 CHECK (duration >= 0),

    completed BOOLEAN NOT NULL DEFAULT FALSE
    );

CREATE INDEX IF NOT EXISTS idx_listening_history_user_played_at
    ON listening_history(user_id, played_at DESC);

CREATE INDEX IF NOT EXISTS idx_listening_history_user_track_played_at
    ON listening_history(user_id, track_id, played_at DESC);

CREATE INDEX IF NOT EXISTS idx_listening_history_track_id
    ON listening_history(track_id);

CREATE INDEX IF NOT EXISTS idx_listening_history_completed
    ON listening_history(completed);

-- Useful later for recommendation queries:
-- "tracks completed by user recently"
CREATE INDEX IF NOT EXISTS idx_listening_history_user_completed_played_at
    ON listening_history(user_id, completed, played_at DESC);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS listening_history;

-- +goose StatementEnd
