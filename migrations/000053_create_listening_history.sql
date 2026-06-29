-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS listening_history (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    played_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ms_played INT NOT NULL DEFAULT 0,
    source VARCHAR(50) NOT NULL DEFAULT 'direct',
    PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS idx_listening_history_user_played ON listening_history(user_id, played_at DESC);
CREATE INDEX IF NOT EXISTS idx_listening_history_track_id ON listening_history(track_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS listening_history;

-- +goose StatementEnd
