-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS queue_items (
                                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    position INTEGER NOT NULL CHECK (position > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX IF NOT EXISTS idx_queue_items_user_position
    ON queue_items(user_id, position);

CREATE INDEX IF NOT EXISTS idx_queue_items_user_id
    ON queue_items(user_id);

CREATE INDEX IF NOT EXISTS idx_queue_items_track_id
    ON queue_items(track_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS queue_items;

-- +goose StatementEnd
