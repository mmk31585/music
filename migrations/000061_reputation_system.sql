-- +goose Up
-- +goose StatementBegin

-- 1. User reputation profile
CREATE TABLE IF NOT EXISTS user_reputation (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    total_contributions INT NOT NULL DEFAULT 0,
    accepted_contributions INT NOT NULL DEFAULT 0,
    rejected_contributions INT NOT NULL DEFAULT 0,
    pending_contributions INT NOT NULL DEFAULT 0,
    trust_score NUMERIC(5,2) NOT NULL DEFAULT 0.0,
    tier VARCHAR(20) NOT NULL DEFAULT 'newcomer',
    upload_slots INT NOT NULL DEFAULT 0,
    auto_publish BOOLEAN NOT NULL DEFAULT FALSE,
    can_review BOOLEAN NOT NULL DEFAULT FALSE,
    last_contribution_at TIMESTAMPTZ,
    contribution_streak_days INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_reputation_user ON user_reputation(user_id);
CREATE INDEX idx_user_reputation_tier ON user_reputation(tier);
CREATE INDEX idx_user_reputation_score ON user_reputation(trust_score DESC);
COMMENT ON TABLE user_reputation IS 'Contributor reputation profile with trust tiers and scoring';

-- 2. Trust tier definitions
CREATE TABLE IF NOT EXISTS trust_tiers (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(20) NOT NULL UNIQUE,
    label VARCHAR(50) NOT NULL,
    description TEXT,
    min_score NUMERIC(5,2) NOT NULL DEFAULT 0,
    min_accepted INT NOT NULL DEFAULT 0,
    upload_slots INT NOT NULL DEFAULT 1,
    auto_publish BOOLEAN NOT NULL DEFAULT FALSE,
    can_review BOOLEAN NOT NULL DEFAULT FALSE,
    hierarchy_level INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE trust_tiers IS 'Configurable trust tier thresholds';

-- 3. Contribution score history
CREATE TABLE IF NOT EXISTS contribution_scores (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contribution_type VARCHAR(50) NOT NULL,
    contribution_id BIGINT,
    score_delta NUMERIC(5,2) NOT NULL DEFAULT 0,
    reason TEXT,
    reviewed_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_contribution_scores_user ON contribution_scores(user_id);
CREATE INDEX idx_contribution_scores_type ON contribution_scores(contribution_type);
CREATE INDEX idx_contribution_scores_created ON contribution_scores(created_at DESC);
COMMENT ON TABLE contribution_scores IS 'Audit trail for trust score changes';

-- 4. Seed trust tier definitions
INSERT INTO trust_tiers (slug, label, description, min_score, min_accepted, upload_slots, auto_publish, can_review, hierarchy_level) VALUES
    ('newcomer',    'Newcomer',    'Just joined the community',                0,    0, 1, FALSE, FALSE, 0),
    ('contributor', 'Contributor', 'Active contributor with accepted work',    10,   3, 2, FALSE, FALSE, 1),
    ('trusted',     'Trusted',     'Consistently high-quality contributions',  30,  10, 5, FALSE, TRUE,  2),
    ('verified',    'Verified',    'Verified creator with proven track record', 50,  25, 10, TRUE,  TRUE,  3),
    ('elite',       'Elite',       'Elite contributor, top of the community',  75,  50, 20, TRUE,  TRUE,  4),
    ('legend',      'Legend',      'Legendary contributor, community pillar',  90, 100, 50, TRUE,  TRUE,  5)
ON CONFLICT (slug) DO NOTHING;

-- 5. Initialize user_reputation for existing users
INSERT INTO user_reputation (user_id, tier, upload_slots, auto_publish, can_review)
SELECT u.id, 'newcomer', 1, FALSE, FALSE
FROM users u
ON CONFLICT (user_id) DO NOTHING;

-- 6. Updated_at trigger
CREATE OR REPLACE FUNCTION update_reputation_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_user_reputation_updated_at
    BEFORE UPDATE ON user_reputation
    FOR EACH ROW EXECUTE FUNCTION update_reputation_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS trigger_user_reputation_updated_at ON user_reputation;
DROP FUNCTION IF EXISTS update_reputation_updated_at();
DROP TABLE IF EXISTS contribution_scores;
DROP TABLE IF EXISTS trust_tiers;
DROP TABLE IF EXISTS user_reputation;

-- +goose StatementEnd
