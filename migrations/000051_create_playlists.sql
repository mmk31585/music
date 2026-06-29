-- +goose Up
-- +goose StatementBegin

-- The playlists table was already created in 000004 with `user_id`
-- instead of `owner_id`. Add the column if missing and sync existing rows.
ALTER TABLE playlists
    ADD COLUMN IF NOT EXISTS owner_id UUID REFERENCES users(id);

UPDATE playlists SET owner_id = user_id WHERE owner_id IS NULL AND user_id IS NOT NULL;
ALTER TABLE playlists ALTER COLUMN owner_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_playlists_owner_id ON playlists(owner_id);
CREATE INDEX IF NOT EXISTS idx_playlists_is_public ON playlists(is_public);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_playlists_is_public;
DROP INDEX IF EXISTS idx_playlists_owner_id;
ALTER TABLE playlists DROP COLUMN IF EXISTS owner_id;

-- +goose StatementEnd
