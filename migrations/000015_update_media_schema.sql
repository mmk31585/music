-- +goose Up
BEGIN;

ALTER TABLE media
    ADD COLUMN IF NOT EXISTS media_type varchar(50),
    ADD COLUMN IF NOT EXISTS storage_provider varchar(50) DEFAULT 'local',
    ADD COLUMN IF NOT EXISTS bucket varchar(255),
    ADD COLUMN IF NOT EXISTS object_key text,
    ADD COLUMN IF NOT EXISTS public_url text,
    ADD COLUMN IF NOT EXISTS file_size bigint,
    ADD COLUMN IF NOT EXISTS checksum_sha256 varchar(64),
    ADD COLUMN IF NOT EXISTS original_filename text,
    ADD COLUMN IF NOT EXISTS metadata jsonb DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS created_by uuid;

UPDATE media
SET media_type = kind
WHERE media_type IS NULL;

UPDATE media
SET public_url = url
WHERE public_url IS NULL;

UPDATE media
SET object_key = url
WHERE object_key IS NULL;

UPDATE media
SET file_size = size_bytes
WHERE file_size IS NULL;

UPDATE media
SET storage_provider = 'local'
WHERE storage_provider IS NULL;

UPDATE media
SET metadata = '{}'::jsonb
WHERE metadata IS NULL;

ALTER TABLE media
    ALTER COLUMN media_type SET NOT NULL,
ALTER COLUMN storage_provider SET NOT NULL,
    ALTER COLUMN metadata SET NOT NULL;

-- New repository inserts do not populate old legacy columns.
ALTER TABLE media
    ALTER COLUMN kind DROP NOT NULL,
ALTER COLUMN url DROP NOT NULL;

CREATE INDEX IF NOT EXISTS idx_media_media_type
    ON media (media_type);

CREATE INDEX IF NOT EXISTS idx_media_checksum_sha256
    ON media (checksum_sha256);

CREATE INDEX IF NOT EXISTS idx_media_created_by
    ON media (created_by);

COMMIT;

-- +goose Down
BEGIN;

DROP INDEX IF EXISTS idx_media_created_by;
DROP INDEX IF EXISTS idx_media_checksum_sha256;
DROP INDEX IF EXISTS idx_media_media_type;

ALTER TABLE media
    ALTER COLUMN kind SET NOT NULL,
ALTER COLUMN url SET NOT NULL;

ALTER TABLE media
DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS metadata,
    DROP COLUMN IF EXISTS original_filename,
    DROP COLUMN IF EXISTS checksum_sha256,
    DROP COLUMN IF EXISTS file_size,
    DROP COLUMN IF EXISTS public_url,
    DROP COLUMN IF EXISTS object_key,
    DROP COLUMN IF EXISTS bucket,
    DROP COLUMN IF EXISTS storage_provider,
    DROP COLUMN IF EXISTS media_type;

COMMIT;
