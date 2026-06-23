-- +goose Up
-- +goose StatementBegin

-- Convert listening_history to a partitioned table by month
-- First, rename existing table
ALTER TABLE listening_history RENAME TO listening_history_old;

-- Create partitioned table
CREATE TABLE listening_history (
    id          UUID NOT NULL DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id    UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    played_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    duration    INTEGER NOT NULL DEFAULT 0,
    source      VARCHAR(50) DEFAULT 'direct',
    device_id   VARCHAR(255),

    PRIMARY KEY (id, played_at)
) PARTITION BY RANGE (played_at);

-- Create monthly partitions
CREATE TABLE listening_history_2026_01 PARTITION OF listening_history
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE listening_history_2026_02 PARTITION OF listening_history
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
CREATE TABLE listening_history_2026_03 PARTITION OF listening_history
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
CREATE TABLE listening_history_2026_04 PARTITION OF listening_history
    FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');
CREATE TABLE listening_history_2026_05 PARTITION OF listening_history
    FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');
CREATE TABLE listening_history_2026_06 PARTITION OF listening_history
    FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');
CREATE TABLE listening_history_2026_07 PARTITION OF listening_history
    FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');
CREATE TABLE listening_history_2026_08 PARTITION OF listening_history
    FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');
CREATE TABLE listening_history_2026_09 PARTITION OF listening_history
    FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');
CREATE TABLE listening_history_2026_10 PARTITION OF listening_history
    FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');
CREATE TABLE listening_history_2026_11 PARTITION OF listening_history
    FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');
CREATE TABLE listening_history_2026_12 PARTITION OF listening_history
    FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');

-- Default partition for future dates
CREATE TABLE listening_history_future PARTITION OF listening_history
    FOR VALUES FROM ('2027-01-01') TO ('2030-01-01');

-- Migrate old data
INSERT INTO listening_history (id, user_id, track_id, played_at, duration)
SELECT id, user_id, track_id, played_at, 0 FROM listening_history_old
ON CONFLICT DO NOTHING;

-- Drop old table
DROP TABLE IF EXISTS listening_history_old;

-- Create indexes on the partitioned table
CREATE INDEX idx_listening_history_user_id ON listening_history(user_id);
CREATE INDEX idx_listening_history_track_id ON listening_history(track_id);
CREATE INDEX idx_listening_history_played_at ON listening_history(played_at);
CREATE INDEX idx_listening_history_user_played_at ON listening_history(user_id, played_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Rollback: convert back to regular table
CREATE TABLE listening_history_new (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id    UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    played_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    duration    INTEGER NOT NULL DEFAULT 0,
    source      VARCHAR(50) DEFAULT 'direct',
    device_id   VARCHAR(255)
);

INSERT INTO listening_history_new (id, user_id, track_id, played_at, duration)
SELECT id, user_id, track_id, played_at, duration FROM listening_history;

DROP TABLE IF EXISTS listening_history CASCADE;
ALTER TABLE listening_history_new RENAME TO listening_history;

CREATE INDEX idx_listening_history_user_id ON listening_history(user_id);
CREATE INDEX idx_listening_history_track_id ON listening_history(track_id);
CREATE INDEX idx_listening_history_played_at ON listening_history(played_at);

-- +goose StatementEnd
