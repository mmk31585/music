-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS xp_transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount INT NOT NULL,
    balance_after INT NOT NULL,
    source VARCHAR(50) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_xp_user ON xp_transactions (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_xp_source ON xp_transactions (source, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_xp_created ON xp_transactions (created_at DESC);

CREATE TABLE IF NOT EXISTS level_definitions (
    level INT PRIMARY KEY,
    xp_required BIGINT NOT NULL,
    title VARCHAR(100) NOT NULL,
    title_persian VARCHAR(100),
    privileges JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS badges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon_url TEXT,
    category VARCHAR(20) NOT NULL,
    rarity VARCHAR(20) NOT NULL DEFAULT 'common',
    criteria JSONB NOT NULL DEFAULT '{}',
    xp_reward INT NOT NULL DEFAULT 0,
    is_hidden BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_badges (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    badge_id UUID NOT NULL REFERENCES badges(id) ON DELETE CASCADE,
    earned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_displayed BOOLEAN NOT NULL DEFAULT true,
    PRIMARY KEY (user_id, badge_id)
);

CREATE TABLE IF NOT EXISTS daily_challenges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    challenge_type VARCHAR(50) NOT NULL,
    target_count INT NOT NULL,
    xp_reward INT NOT NULL,
    badge_reward_id UUID REFERENCES badges(id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    valid_from DATE NOT NULL,
    valid_until DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_challenges_active ON daily_challenges (is_active, valid_from, valid_until);

CREATE TABLE IF NOT EXISTS user_challenges (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    challenge_id UUID NOT NULL REFERENCES daily_challenges(id) ON DELETE CASCADE,
    progress INT NOT NULL DEFAULT 0,
    is_completed BOOLEAN NOT NULL DEFAULT false,
    completed_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, challenge_id)
);

CREATE TABLE IF NOT EXISTS leaderboard_snapshots (
    id BIGSERIAL,
    leaderboard_type VARCHAR(30) NOT NULL,
    rank INT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    score BIGINT NOT NULL,
    snapshot_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (leaderboard_type, snapshot_date, rank)
);

CREATE INDEX IF NOT EXISTS idx_leaderboard_snapshot ON leaderboard_snapshots (leaderboard_type, snapshot_date, rank);

-- Seed level definitions
INSERT INTO level_definitions (level, xp_required, title, title_persian, privileges) VALUES
    (1, 0, 'Novice Listener', 'شنونده تازه‌کار', '{"create_playlist": true}'::jsonb),
    (2, 100, 'Curious Ear', 'گوش کنجکاو', '{"create_playlist": true}'::jsonb),
    (3, 520, 'Melody Seeker', 'جوینده ملودی', '{"create_playlist": true, "upload_avatar": true}'::jsonb),
    (4, 1600, 'Rhythm Catcher', 'ریتم‌گیر', '{"create_playlist": true, "upload_avatar": true}'::jsonb),
    (5, 3900, 'Harmony Lover', 'عاشق هارمونی', '{"create_playlist": true, "upload_avatar": true, "collaborative_playlist": true}'::jsonb),
    (6, 8000, 'Tone Explorer', 'کاوشگر نغمه', '{}'::jsonb),
    (7, 14000, 'Music Wanderer', 'سرگردان موسیقی', '{}'::jsonb),
    (8, 22000, 'Note Collector', 'جمع‌آورنده نت', '{}'::jsonb),
    (9, 33000, 'Groove Master', 'استاد گروو', '{"contribute_lyrics": true}'::jsonb),
    (10, 48000, 'Beat Connoisseur', 'خبره ضرب‌آهنگ', '{"contribute_lyrics": true}'::jsonb),
    (11, 66000, 'Vocal Virtuoso', 'استاد آواز', '{"moderate_content": true}'::jsonb),
    (12, 88000, 'Lyric Sage', 'فرزانه شعر', '{"moderate_content": true}'::jsonb),
    (13, 115000, 'Dastgah Navigator', 'راهبر دستگاه', '{"moderate_content": true, "direct_publish": true}'::jsonb),
    (14, 147000, 'Genre Weaver', 'بافنده سبک', '{"direct_publish": true}'::jsonb),
    (15, 185000, 'Orchestra Commander', 'فرمانده ارکستر', '{"direct_publish": true}'::jsonb),
    (16, 230000, 'Melody Architect', 'معمار ملودی', '{"create_club": true}'::jsonb),
    (17, 283000, 'Music Scholar', 'دانشور موسیقی', '{"create_club": true}'::jsonb),
    (18, 345000, 'Sound Alchemist', 'کیمیاگر صدا', '{"create_club": true, "veto_moderation": true}'::jsonb),
    (19, 417000, 'Legendary Listener', 'شنونده افسانه‌ای', '{"veto_moderation": true}'::jsonb),
    (20, 500000, 'Persian Gem', 'گوهر پارسی', '{"veto_moderation": true, "all": true}'::jsonb);

-- Seed badges
INSERT INTO badges (id, name, description, category, rarity, criteria, xp_reward) VALUES
    (gen_random_uuid(), 'First Stream', 'Stream your first track', 'listener', 'common', '{"type": "stream_count", "target": 1}', 10),
    (gen_random_uuid(), 'Night Owl', 'Stream 100 tracks between midnight and 5am', 'listener', 'rare', '{"type": "night_streams", "target": 100}', 100),
    (gen_random_uuid(), 'Explorer', 'Listen to tracks from 20 different genres', 'listener', 'epic', '{"type": "genre_count", "target": 20}', 250),
    (gen_random_uuid(), 'Marathon', 'Stream 1000 tracks in a single week', 'listener', 'legendary', '{"type": "weekly_streams", "target": 1000}', 500),
    (gen_random_uuid(), 'Wordsmith', 'Submit 10 approved lyrics contributions', 'contributor', 'rare', '{"type": "lyrics_approved", "target": 10}', 150),
    (gen_random_uuid(), 'Polyglot', 'Submit 5 approved translations', 'contributor', 'epic', '{"type": "translations_approved", "target": 5}', 200),
    (gen_random_uuid(), 'Curator', 'Create 10 playlists', 'contributor', 'common', '{"type": "playlists_created", "target": 10}', 50),
    (gen_random_uuid(), 'Networker', 'Gain 100 followers', 'social', 'rare', '{"type": "followers", "target": 100}', 150),
    (gen_random_uuid(), 'Party Host', 'Host 10 listening parties', 'social', 'epic', '{"type": "parties_hosted", "target": 10}', 300),
    (gen_random_uuid(), 'Club President', 'Create a club with 100+ members', 'social', 'legendary', '{"type": "club_members", "target": 100}', 500),
    (gen_random_uuid(), 'Chart Climber', 'Reach top 10 on the weekly leaderboard', 'special', 'epic', '{"type": "leaderboard_rank", "target": 10}', 200),
    (gen_random_uuid(), 'Community Hero', 'Complete 50 moderation reviews', 'special', 'legendary', '{"type": "moderation_reviews", "target": 50}', 1000);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS leaderboard_snapshots;
DROP TABLE IF EXISTS user_challenges;
DROP TABLE IF EXISTS daily_challenges;
DROP TABLE IF EXISTS user_badges;
DROP TABLE IF EXISTS badges;
DROP TABLE IF EXISTS level_definitions;
DROP TABLE IF EXISTS xp_transactions;

-- +goose StatementEnd
