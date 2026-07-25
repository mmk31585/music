-- +goose Up
-- Payment gateways, tips, and invoice support
-- Migration 000022

-- ============================================================
-- TIPS (creator tipping / donations)
-- ============================================================
CREATE TABLE IF NOT EXISTS tips (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    artist_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id        UUID REFERENCES tracks(id) ON DELETE SET NULL,
    amount_cents    BIGINT NOT NULL CHECK (amount_cents > 0),
    currency        VARCHAR(10) NOT NULL DEFAULT 'IRR',
    message         TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    provider        VARCHAR(50),
    provider_pay_id VARCHAR(255),
    paid_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tips_sender   ON tips(sender_id);
CREATE INDEX idx_tips_artist   ON tips(artist_id);
CREATE INDEX idx_tips_status   ON tips(status);

-- ============================================================
-- PAYMENT ATTEMPTS (gateway transaction log)
-- ============================================================
CREATE TABLE IF NOT EXISTS payment_attempts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount_cents    BIGINT NOT NULL,
    currency        VARCHAR(10) NOT NULL DEFAULT 'IRR',
    provider        VARCHAR(50) NOT NULL,
    provider_pay_id VARCHAR(255),
    authority       VARCHAR(255),
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    purpose         VARCHAR(50) NOT NULL,
    target_id       UUID,
    target_type     VARCHAR(50),
    description     TEXT,
    failure_reason  TEXT,
    redirect_url    TEXT,
    paid_at         TIMESTAMPTZ,
    expires_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pay_attempts_user   ON payment_attempts(user_id);
CREATE INDEX idx_pay_attempts_status ON payment_attempts(status);
CREATE INDEX idx_pay_attempts_provider ON payment_attempts(provider);

-- ============================================================
-- Seed payment gateway config (loaded as reference data)
-- ============================================================
INSERT INTO plans (code, name, description, price_cents, currency, interval, is_premium, is_active, features)
VALUES
    ('premium_monthly', 'Premium Monthly', 'Ad-free, high quality audio, offline downloads, unlimited skips',
     99900, 'IRR', 'monthly', TRUE, TRUE,
     '["ad_free","high_quality_audio","offline_downloads","unlimited_skips","early_access"]'::jsonb)
ON CONFLICT (code) DO UPDATE SET
    is_active = TRUE,
    price_cents = 99900,
    description = 'Ad-free, high quality audio, offline downloads, unlimited skips',
    features = '["ad_free","high_quality_audio","offline_downloads","unlimited_skips","early_access"]'::jsonb,
    updated_at = NOW();

INSERT INTO plans (code, name, description, price_cents, currency, interval, is_premium, is_active, features)
VALUES
    ('premium_yearly', 'Premium Yearly', 'All premium features at a discounted annual rate',
     599000, 'IRR', 'yearly', TRUE, TRUE,
     '["ad_free","high_quality_audio","offline_downloads","unlimited_skips","early_access","yearly_discount"]'::jsonb)
ON CONFLICT (code) DO UPDATE SET
    is_active = TRUE,
    price_cents = 599000,
    description = 'All premium features at a discounted annual rate',
    features = '["ad_free","high_quality_audio","offline_downloads","unlimited_skips","early_access","yearly_discount"]'::jsonb,
    updated_at = NOW();

-- +goose Down
DROP TABLE IF EXISTS tips CASCADE;
DROP TABLE IF EXISTS payment_attempts CASCADE;
