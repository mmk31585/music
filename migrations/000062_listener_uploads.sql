-- +goose Up
-- +goose StatementBegin

ALTER TABLE ingestion_drafts
    ADD COLUMN IF NOT EXISTS upload_source TEXT NOT NULL DEFAULT 'admin'
        CHECK (upload_source IN ('admin', 'listener', 'club')),
    ADD COLUMN IF NOT EXISTS club_id UUID,
    ADD COLUMN IF NOT EXISTS needs_review BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS review_notes TEXT,
    ADD COLUMN IF NOT EXISTS reviewed_by UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMPTZ;

DO $$
BEGIN
    IF EXISTS (SELECT FROM pg_tables WHERE tablename = 'music_clubs') THEN
        ALTER TABLE ingestion_drafts
            ADD CONSTRAINT fk_ingestion_drafts_club
            FOREIGN KEY (club_id) REFERENCES music_clubs(id) ON DELETE SET NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_ingestion_drafts_source ON ingestion_drafts(upload_source);
CREATE INDEX IF NOT EXISTS idx_ingestion_drafts_needs_review ON ingestion_drafts(needs_review) WHERE needs_review = TRUE;
CREATE INDEX IF NOT EXISTS idx_ingestion_drafts_club ON ingestion_drafts(club_id) WHERE club_id IS NOT NULL;

COMMENT ON COLUMN ingestion_drafts.upload_source IS 'Origin: admin (direct), listener (community), club (club submission)';
COMMENT ON COLUMN ingestion_drafts.club_id IS 'If upload_source=club, the club this was submitted to';
COMMENT ON COLUMN ingestion_drafts.needs_review IS 'Whether this draft requires manual review before publishing';
COMMENT ON COLUMN ingestion_drafts.review_notes IS 'Notes from the reviewer on accept/reject';

CREATE TABLE IF NOT EXISTS co_uploaders (
    id BIGSERIAL PRIMARY KEY,
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'contributor'
        CHECK (role IN ('uploader', 'contributor', 'featured')),
    xp_share_percent NUMERIC(5,2) NOT NULL DEFAULT 100.0,
    contribution_type TEXT,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(track_id, user_id)
);

CREATE INDEX idx_co_uploaders_track ON co_uploaders(track_id);
CREATE INDEX idx_co_uploaders_user ON co_uploaders(user_id);
COMMENT ON TABLE co_uploaders IS 'Tracks can have multiple uploaders; XP is split by xp_share_percent';

CREATE TABLE IF NOT EXISTS upload_slots (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    used_slots INT NOT NULL DEFAULT 0,
    max_slots INT NOT NULL DEFAULT 1,
    last_upload_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_upload_slots_user ON upload_slots(user_id);
COMMENT ON TABLE upload_slots IS 'Tracks concurrent upload slots per user based on trust tier';

INSERT INTO upload_slots (user_id, used_slots, max_slots)
SELECT u.id, 0, 1
FROM users u
ON CONFLICT (user_id) DO NOTHING;

CREATE OR REPLACE FUNCTION update_upload_slots_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_upload_slots_updated_at
    BEFORE UPDATE ON upload_slots
    FOR EACH ROW EXECUTE FUNCTION update_upload_slots_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS trigger_upload_slots_updated_at ON upload_slots;
DROP FUNCTION IF EXISTS update_upload_slots_updated_at();
DROP TABLE IF EXISTS upload_slots;
DROP TABLE IF EXISTS co_uploaders;

ALTER TABLE ingestion_drafts DROP CONSTRAINT IF EXISTS fk_ingestion_drafts_club;
ALTER TABLE ingestion_drafts
    DROP COLUMN IF EXISTS reviewed_at,
    DROP COLUMN IF EXISTS reviewed_by,
    DROP COLUMN IF EXISTS review_notes,
    DROP COLUMN IF EXISTS needs_review,
    DROP COLUMN IF EXISTS club_id,
    DROP COLUMN IF EXISTS upload_source;

-- +goose StatementEnd