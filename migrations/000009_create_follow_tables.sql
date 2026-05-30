-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS user_follows (
                                            follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (follower_id, followee_id),
    CONSTRAINT user_follows_no_self_follow CHECK (follower_id <> followee_id)
    );

CREATE INDEX IF NOT EXISTS idx_user_follows_follower_id
    ON user_follows(follower_id);

CREATE INDEX IF NOT EXISTS idx_user_follows_followee_id
    ON user_follows(followee_id);

CREATE TABLE IF NOT EXISTS artist_follows (
                                              user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, artist_id)
    );

CREATE INDEX IF NOT EXISTS idx_artist_follows_user_id
    ON artist_follows(user_id);

CREATE INDEX IF NOT EXISTS idx_artist_follows_artist_id
    ON artist_follows(artist_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS artist_follows;
DROP TABLE IF EXISTS user_follows;

-- +goose StatementEnd
