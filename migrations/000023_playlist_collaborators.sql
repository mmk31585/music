-- +goose Up
-- +goose StatementBegin

-- Add is_collaborative column to playlists if not exists
ALTER TABLE playlists ADD COLUMN IF NOT EXISTS is_collaborative BOOLEAN NOT NULL DEFAULT FALSE;

-- Create playlist_collaborators table
CREATE TABLE IF NOT EXISTS playlist_collaborators (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    playlist_id UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    added_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(playlist_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_playlist_collaborators_playlist_id
    ON playlist_collaborators(playlist_id);

CREATE INDEX IF NOT EXISTS idx_playlist_collaborators_user_id
    ON playlist_collaborators(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS playlist_collaborators;
ALTER TABLE playlists DROP COLUMN IF EXISTS is_collaborative;

-- +goose StatementEnd
