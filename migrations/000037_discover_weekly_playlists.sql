-- +goose Up
CREATE TABLE IF NOT EXISTS discover_weekly_playlists (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_ids UUID[] NOT NULL DEFAULT '{}',
    generated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    week_of DATE NOT NULL,
    PRIMARY KEY (user_id, week_of)
);

CREATE INDEX IF NOT EXISTS idx_discover_weekly_user_week ON discover_weekly_playlists(user_id, week_of);

-- +goose Down
DROP TABLE IF EXISTS discover_weekly_playlists;
