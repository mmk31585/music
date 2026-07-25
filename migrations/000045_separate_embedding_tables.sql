-- +goose Up
-- +goose StatementBegin

-- Step 1: Create track_embeddings_text for OpenAI text embeddings (Go module)
CREATE TABLE IF NOT EXISTS track_embeddings_text (
    track_id UUID PRIMARY KEY REFERENCES tracks(id) ON DELETE CASCADE,
    embedding double precision[] NOT NULL,
    model_version VARCHAR(100) NOT NULL DEFAULT 'v1',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_track_embeddings_text_track_id ON track_embeddings_text(track_id);

-- Step 2: Rename existing track_embeddings to track_embeddings_audio (for openl3)
ALTER TABLE IF EXISTS track_embeddings RENAME TO track_embeddings_audio;

-- Step 3: Add audio-specific columns if missing (they exist in Python model)
ALTER TABLE IF EXISTS track_embeddings_audio ADD COLUMN IF NOT EXISTS embedding_model VARCHAR(100);
ALTER TABLE IF EXISTS track_embeddings_audio ADD COLUMN IF NOT EXISTS tempo_bpm double precision;
ALTER TABLE IF EXISTS track_embeddings_audio ADD COLUMN IF NOT EXISTS genre_ids TEXT[] DEFAULT '{}';
ALTER TABLE IF EXISTS track_embeddings_audio ADD COLUMN IF NOT EXISTS release_year INT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE IF EXISTS track_embeddings_audio RENAME TO track_embeddings;
DROP TABLE IF EXISTS track_embeddings_text CASCADE;

-- +goose StatementEnd
