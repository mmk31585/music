-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS lyrics (
                                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Default added
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    language VARCHAR(10) NOT NULL,
    type VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(), -- Default added for consistency

    CONSTRAINT lyrics_track_id_language_key UNIQUE (track_id, language)
    );

-- Index for track_id is good for finding all languages for a track
CREATE INDEX IF NOT EXISTS idx_lyrics_track_id ON lyrics (track_id);

-- Note: No need for idx_lyrics_track_id_language because the UNIQUE constraint
-- creates an index for us automatically.

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS lyrics;

-- +goose StatementEnd
