-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       email VARCHAR(255) NOT NULL UNIQUE,
                       username VARCHAR(50) NOT NULL UNIQUE,
                       display_name VARCHAR(100) NOT NULL,
                       password_hash TEXT NOT NULL,
                       avatar_url TEXT,
                       role VARCHAR(30) NOT NULL DEFAULT 'user',
                       is_active BOOLEAN NOT NULL DEFAULT TRUE,
                       email_verified_at TIMESTAMPTZ,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_created_at ON users(created_at);

CREATE TABLE auth_sessions (
                               id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                               user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               refresh_token_hash TEXT NOT NULL UNIQUE,
                               user_agent TEXT,
                               ip_address TEXT,
                               expires_at TIMESTAMPTZ NOT NULL,
                               revoked_at TIMESTAMPTZ,
                               created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auth_sessions_user_id ON auth_sessions(user_id);
CREATE INDEX idx_auth_sessions_refresh_token_hash ON auth_sessions(refresh_token_hash);
CREATE INDEX idx_auth_sessions_expires_at ON auth_sessions(expires_at);

-- +goose Down

DROP TABLE IF EXISTS auth_sessions;
DROP TABLE IF EXISTS users;
