-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS plans (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,

    price_cents BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',

    interval VARCHAR(20) NOT NULL DEFAULT 'month',
    is_premium BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    features JSONB NOT NULL DEFAULT '[]'::jsonb,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS subscriptions (
                                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL,
    plan_id UUID NOT NULL REFERENCES plans(id),

    status VARCHAR(30) NOT NULL DEFAULT 'active',
    current_period_start TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    current_period_end TIMESTAMPTZ,
    cancel_at_period_end BOOLEAN NOT NULL DEFAULT FALSE,
    canceled_at TIMESTAMPTZ,

    provider VARCHAR(50),
    provider_subscription_id VARCHAR(255),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT subscriptions_status_check CHECK (
                                                    status IN ('active', 'trialing', 'pending', 'past_due', 'canceled', 'expired')
    )
    );

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_plan_id ON subscriptions(plan_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);

CREATE TABLE IF NOT EXISTS payments (
                                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL,
    subscription_id UUID REFERENCES subscriptions(id),
    plan_id UUID REFERENCES plans(id),

    amount_cents BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',

    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    provider VARCHAR(50),
    provider_payment_id VARCHAR(255),

    failure_reason TEXT,
    paid_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT payments_status_check CHECK (
                                               status IN ('pending', 'paid', 'failed', 'refunded', 'canceled')
    )
    );

CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_subscription_id ON payments(subscription_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);

INSERT INTO plans (
    code,
    name,
    description,
    price_cents,
    currency,
    interval,
    is_premium,
    is_active,
    features
)
VALUES
    (
        'free',
        'Free',
        'Free music account',
        0,
        'USD',
        'month',
        FALSE,
        TRUE,
        '[
            "Listen to music",
            "Create playlists",
            "Like tracks",
            "Follow artists"
        ]'::jsonb
    ),
    (
        'premium_monthly',
        'Premium Monthly',
        'Future premium monthly subscription. Not available yet.',
        999,
        'USD',
        'month',
        TRUE,
        FALSE,
        '[
            "Ad-free listening",
            "High quality audio",
            "Offline downloads",
            "Unlimited skips"
        ]'::jsonb
    )
    ON CONFLICT (code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS plans;
