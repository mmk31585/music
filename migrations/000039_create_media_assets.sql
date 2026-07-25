-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS media_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    media_type       VARCHAR(50) NOT NULL DEFAULT 'audio',
    storage_provider VARCHAR(50) NOT NULL DEFAULT 'local',
    bucket           VARCHAR(255),
    object_key       TEXT NOT NULL,
    public_url       TEXT,

    original_filename VARCHAR(512),
    mime_type         VARCHAR(127),
    file_size         BIGINT,
    checksum_sha256   CHAR(64),

    duration_seconds INTEGER,
    width            INTEGER,
    height           INTEGER,

    metadata         JSONB NOT NULL DEFAULT '{}',

    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_media_assets_media_type ON media_assets(media_type);
CREATE INDEX idx_media_assets_created_by ON media_assets(created_by);
CREATE INDEX idx_media_assets_deleted_at ON media_assets(deleted_at);
CREATE INDEX idx_media_assets_storage_provider ON media_assets(storage_provider);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS media_assets;

-- +goose StatementEnd
