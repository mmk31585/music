-- +goose Up
-- +goose StatementBegin

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE email = 'admin@example.com') THEN
        INSERT INTO users (email, username, display_name, password_hash, role, is_active, email_verified_at)
        VALUES ('admin@example.com', 'admin', 'Admin', '$2a$10$cXbSheUUueeN1Csk87JSjeSA.N.qZkz2a4Pc3PzAhDYJ2a7LroG6O', 'admin', TRUE, NOW());
    END IF;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM users WHERE email = 'admin@example.com';

-- +goose StatementEnd
