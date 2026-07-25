-- +goose Up
-- +goose StatementBegin

-- The liked_tracks table was already created in 000005 with `created_at`
-- instead of `liked_at`. Add the column if missing, then create indexes.
ALTER TABLE liked_tracks
    ADD COLUMN IF NOT EXISTS liked_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- Sync liked_at from created_at for existing rows.
UPDATE liked_tracks SET liked_at = created_at WHERE liked_at IS DISTINCT FROM created_at;

CREATE INDEX IF NOT EXISTS idx_liked_tracks_user_id ON liked_tracks(user_id);
CREATE INDEX IF NOT EXISTS idx_liked_tracks_track_id ON liked_tracks(track_id);
CREATE INDEX IF NOT EXISTS idx_liked_tracks_liked_at ON liked_tracks(liked_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_liked_tracks_liked_at;
DROP INDEX IF EXISTS idx_liked_tracks_track_id;
DROP INDEX IF EXISTS idx_liked_tracks_user_id;
ALTER TABLE liked_tracks DROP COLUMN IF EXISTS liked_at;

-- +goose StatementEnd
