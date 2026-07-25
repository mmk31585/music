-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS contributions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contributor_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contribution_type VARCHAR(50) NOT NULL,
    target_type VARCHAR(20) NOT NULL,
    target_id UUID NOT NULL,
    locale VARCHAR(10),
    data JSONB NOT NULL,
    summary TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    ai_verdict VARCHAR(20),
    ai_confidence DOUBLE PRECISION,
    ai_reason TEXT,
    moderator_id UUID REFERENCES users(id),
    moderator_note TEXT,
    version INT NOT NULL DEFAULT 1,
    is_minor BOOLEAN NOT NULL DEFAULT false,
    xp_awarded INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    decided_at TIMESTAMPTZ,
    applied_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contributions_status ON contributions (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_contributions_type ON contributions (contribution_type, status);
CREATE INDEX IF NOT EXISTS idx_contributions_target ON contributions (target_type, target_id);
CREATE INDEX IF NOT EXISTS idx_contributions_contributor ON contributions (contributor_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_contributions_moderator ON contributions (moderator_id) WHERE moderator_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS contribution_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contribution_id UUID NOT NULL REFERENCES contributions(id) ON DELETE CASCADE,
    data JSONB NOT NULL,
    previous_data JSONB,
    changed_by UUID NOT NULL REFERENCES users(id),
    change_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contribution_history_contribution ON contribution_history (contribution_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_contribution_history_changed_by ON contribution_history (changed_by, created_at DESC);

CREATE TABLE IF NOT EXISTS content_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    target_type VARCHAR(20) NOT NULL,
    target_id UUID NOT NULL,
    version INT NOT NULL,
    data JSONB NOT NULL,
    applied_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (target_type, target_id, version)
);

CREATE INDEX IF NOT EXISTS idx_content_versions_target ON content_versions (target_type, target_id, version DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS content_versions;
DROP TABLE IF EXISTS contribution_history;
DROP TABLE IF EXISTS contributions;

-- +goose StatementEnd
