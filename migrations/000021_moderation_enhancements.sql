-- +goose Up
-- Moderation enhancements: audit log, resolution notes, moderation queue metadata
-- Migration 000021

-- ============================================================
-- MODERATION ACTIONS (audit log for all moderation activity)
-- ============================================================
CREATE TABLE IF NOT EXISTS moderation_actions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id       UUID REFERENCES content_reports(id) ON DELETE SET NULL,
    moderator_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action          VARCHAR(50) NOT NULL,
    target_id       VARCHAR(255) NOT NULL,
    target_type     VARCHAR(50) NOT NULL,
    previous_status VARCHAR(20),
    new_status      VARCHAR(20),
    note            TEXT,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_mod_actions_report    ON moderation_actions(report_id);
CREATE INDEX idx_mod_actions_moderator ON moderation_actions(moderator_id);
CREATE INDEX idx_mod_actions_target    ON moderation_actions(target_id, target_type);
CREATE INDEX idx_mod_actions_created   ON moderation_actions(created_at DESC);

-- ============================================================
-- Add resolution_note to content_reports
-- ============================================================
ALTER TABLE content_reports ADD COLUMN IF NOT EXISTS resolution_note TEXT;

-- Populate initial resolution notes for existing resolved reports
UPDATE content_reports SET resolution_note = 'Migrated from legacy moderation' WHERE status != 'pending' AND resolution_note IS NULL;

CREATE INDEX IF NOT EXISTS idx_reports_status_created ON content_reports(status, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS moderation_actions CASCADE;
ALTER TABLE content_reports DROP COLUMN IF EXISTS resolution_note;
DROP INDEX IF EXISTS idx_reports_status_created;
