-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS track_credits (
                                             id UUID PRIMARY KEY DEFAULT gen_random_uuid()
    );

ALTER TABLE track_credits
    ADD COLUMN IF NOT EXISTS track_id UUID;

ALTER TABLE track_credits
    ADD COLUMN IF NOT EXISTS artist_id UUID;

ALTER TABLE track_credits
    ADD COLUMN IF NOT EXISTS credit_type TEXT;

ALTER TABLE track_credits
    ADD COLUMN IF NOT EXISTS position INT NOT NULL DEFAULT 1;

ALTER TABLE track_credits
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

UPDATE track_credits
SET credit_type = 'unknown'
WHERE credit_type IS NULL;

ALTER TABLE track_credits
    ALTER COLUMN track_id SET NOT NULL;

ALTER TABLE track_credits
    ALTER COLUMN artist_id SET NOT NULL;

ALTER TABLE track_credits
    ALTER COLUMN credit_type SET NOT NULL;

-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'track_credits_track_id_fkey'
    ) THEN
ALTER TABLE track_credits
    ADD CONSTRAINT track_credits_track_id_fkey
        FOREIGN KEY (track_id)
            REFERENCES tracks(id)
            ON DELETE CASCADE;
END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'track_credits_artist_id_fkey'
    ) THEN
ALTER TABLE track_credits
    ADD CONSTRAINT track_credits_artist_id_fkey
        FOREIGN KEY (artist_id)
            REFERENCES artists(id)
            ON DELETE CASCADE;
END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'track_credits_credit_type_not_empty'
    ) THEN
ALTER TABLE track_credits
    ADD CONSTRAINT track_credits_credit_type_not_empty
        CHECK (length(trim(credit_type)) > 0);
END IF;

    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'track_credits_position_positive'
    ) THEN
ALTER TABLE track_credits
    ADD CONSTRAINT track_credits_position_positive
        CHECK (position > 0);
END IF;
END $$;
-- +goose StatementEnd

CREATE INDEX IF NOT EXISTS idx_track_credits_track_id
    ON track_credits(track_id);

CREATE INDEX IF NOT EXISTS idx_track_credits_artist_id
    ON track_credits(artist_id);

CREATE INDEX IF NOT EXISTS idx_track_credits_credit_type
    ON track_credits(credit_type);


CREATE TABLE IF NOT EXISTS related_artists (
                                               artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,
    related_artist_id UUID NOT NULL REFERENCES artists(id) ON DELETE CASCADE,

    score NUMERIC(8, 3) NOT NULL DEFAULT 0,
    source TEXT NOT NULL DEFAULT 'manual',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (artist_id, related_artist_id),

    CONSTRAINT related_artists_not_self
    CHECK (artist_id <> related_artist_id)
    );

CREATE INDEX IF NOT EXISTS idx_related_artists_artist_id
    ON related_artists(artist_id);

CREATE INDEX IF NOT EXISTS idx_related_artists_related_artist_id
    ON related_artists(related_artist_id);

CREATE INDEX IF NOT EXISTS idx_related_artists_score
    ON related_artists(score DESC);


-- +goose Down
DROP INDEX IF EXISTS idx_related_artists_score;
DROP INDEX IF EXISTS idx_related_artists_related_artist_id;
DROP INDEX IF EXISTS idx_related_artists_artist_id;
DROP TABLE IF EXISTS related_artists;

DROP INDEX IF EXISTS idx_track_credits_credit_type;
DROP INDEX IF EXISTS idx_track_credits_artist_id;
DROP INDEX IF EXISTS idx_track_credits_track_id;
DROP TABLE IF EXISTS track_credits;
