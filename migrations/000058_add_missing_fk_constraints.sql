-- +goose Up
-- +goose StatementBegin

-- ── Missing FK: subscriptions.user_id → users(id) ──────────────────────
-- Migration 000006 created subscriptions.user_id as UUID NOT NULL but
-- omitted the FK reference to users(id). Add it with NOT VALID so we
-- can validate existing data separately (zero-downtime approach).
ALTER TABLE subscriptions
    ADD CONSTRAINT fk_subscriptions_user_id
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        NOT VALID;

-- Throws if any existing rows reference non-existent users.
-- Run this during a maintenance window if the table has orphaned rows.
ALTER TABLE subscriptions VALIDATE CONSTRAINT fk_subscriptions_user_id;

-- ── Missing FK: payments.user_id → users(id) ───────────────────────────
ALTER TABLE payments
    ADD CONSTRAINT fk_payments_user_id
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        NOT VALID;

ALTER TABLE payments VALIDATE CONSTRAINT fk_payments_user_id;

-- ── Missing FK: notifications.user_id → users(id) ──────────────────────
ALTER TABLE notifications
    ADD CONSTRAINT fk_notifications_user_id
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        NOT VALID;

ALTER TABLE notifications VALIDATE CONSTRAINT fk_notifications_user_id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE subscriptions DROP CONSTRAINT IF EXISTS fk_subscriptions_user_id;
ALTER TABLE payments DROP CONSTRAINT IF EXISTS fk_payments_user_id;
ALTER TABLE notifications DROP CONSTRAINT IF EXISTS fk_notifications_user_id;

-- +goose StatementEnd
