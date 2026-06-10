-- +goose Up
-- Social features: Listening Parties, Live Rooms, Music Clubs, Discussions
-- Migration 000020

-- ============================================================
-- LISTENING PARTIES
-- ============================================================
CREATE TABLE IF NOT EXISTS listening_parties (
                                                 id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    host_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title          VARCHAR(255) NOT NULL,
    description    TEXT,
    cover_url      TEXT,
    is_public      BOOLEAN NOT NULL DEFAULT true,
    status         VARCHAR(20) NOT NULL DEFAULT 'active'
    CHECK (status IN ('active', 'paused', 'ended')),
    current_track_id   UUID REFERENCES tracks(id) ON DELETE SET NULL,
    current_position_ms BIGINT DEFAULT 0,
    started_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at       TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS listening_party_participants (
                                                            id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    party_id   UUID NOT NULL REFERENCES listening_parties(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    left_at    TIMESTAMPTZ,
    is_active  BOOLEAN NOT NULL DEFAULT true,
    UNIQUE(party_id, user_id)
    );

CREATE INDEX idx_listening_parties_host ON listening_parties(host_id);
CREATE INDEX idx_listening_parties_status ON listening_parties(status);
CREATE INDEX idx_lpp_party ON listening_party_participants(party_id);
CREATE INDEX idx_lpp_user ON listening_party_participants(user_id);
CREATE INDEX idx_lpp_active ON listening_party_participants(party_id, is_active) WHERE is_active = true;

-- ============================================================
-- LIVE ROOMS
-- ============================================================
CREATE TABLE IF NOT EXISTS live_rooms (
                                          id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    host_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    cover_url       TEXT,
    is_public       BOOLEAN NOT NULL DEFAULT true,
    status          VARCHAR(20) NOT NULL DEFAULT 'live'
    CHECK (status IN ('live', 'ended')),
    current_track_id    UUID REFERENCES tracks(id) ON DELETE SET NULL,
    listener_count  INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at        TIMESTAMPTZ
    );

CREATE TABLE IF NOT EXISTS live_room_participants (
                                                      id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id   UUID NOT NULL REFERENCES live_rooms(id) ON DELETE CASCADE,
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role      VARCHAR(20) NOT NULL DEFAULT 'listener'
    CHECK (role IN ('host', 'co_host', 'speaker', 'listener')),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    left_at   TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT true,
    UNIQUE(room_id, user_id)
    );

CREATE TABLE IF NOT EXISTS live_room_queue (
                                               id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id    UUID NOT NULL REFERENCES live_rooms(id) ON DELETE CASCADE,
    track_id   UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    added_by   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    position   INT NOT NULL DEFAULT 0,
    status     VARCHAR(20) NOT NULL DEFAULT 'queued'
    CHECK (status IN ('queued', 'playing', 'played', 'skipped')),
    played_at  TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX idx_live_rooms_host ON live_rooms(host_id);
CREATE INDEX idx_live_rooms_status ON live_rooms(status);
CREATE INDEX idx_lrp_room ON live_room_participants(room_id);
CREATE INDEX idx_lrp_user ON live_room_participants(user_id);
CREATE INDEX idx_lrq_room ON live_room_queue(room_id, position);

-- ============================================================
-- MUSIC CLUBS
-- ============================================================
CREATE TABLE IF NOT EXISTS music_clubs (
                                           id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    cover_url   TEXT,
    created_by  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_public   BOOLEAN NOT NULL DEFAULT true,
    max_members INT NOT NULL DEFAULT 1000,
    member_count INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS music_club_members (
                                                  id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    club_id   UUID NOT NULL REFERENCES music_clubs(id) ON DELETE CASCADE,
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role      VARCHAR(20) NOT NULL DEFAULT 'member'
    CHECK (role IN ('admin', 'moderator', 'member')),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(club_id, user_id)
    );

CREATE TABLE IF NOT EXISTS music_club_posts (
                                                id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    club_id    UUID NOT NULL REFERENCES music_clubs(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX idx_music_clubs_creator ON music_clubs(created_by);
CREATE INDEX idx_mcm_club ON music_club_members(club_id);
CREATE INDEX idx_mcm_user ON music_club_members(user_id);
CREATE INDEX idx_mcp_club ON music_club_posts(club_id);

-- ============================================================
-- DISCUSSIONS (track/album/playlist comments)
-- ============================================================
CREATE TABLE IF NOT EXISTS discussions (
                                           id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_type VARCHAR(50) NOT NULL
    CHECK (target_type IN ('track', 'album', 'playlist', 'artist')),
    target_id   UUID NOT NULL,
    content     TEXT NOT NULL,
    parent_id   UUID REFERENCES discussions(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX idx_discussions_target ON discussions(target_type, target_id);
CREATE INDEX idx_discussions_user ON discussions(user_id);
CREATE INDEX idx_discussions_parent ON discussions(parent_id);

-- ============================================================
-- TRACK RATINGS
-- ============================================================
CREATE TABLE IF NOT EXISTS track_ratings (
                                             id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id   UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    rating     INT NOT NULL CHECK (rating >= 1 AND rating <= 10),
    review     TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, track_id)
    );

CREATE INDEX idx_track_ratings_track ON track_ratings(track_id);
CREATE INDEX idx_track_ratings_user ON track_ratings(user_id);

-- +goose Down
DROP TABLE IF EXISTS track_ratings;
DROP TABLE IF EXISTS discussions;
DROP TABLE IF EXISTS music_club_posts;
DROP TABLE IF EXISTS music_club_members;
DROP TABLE IF EXISTS music_clubs;
DROP TABLE IF EXISTS live_room_queue;
DROP TABLE IF EXISTS live_room_participants;
DROP TABLE IF EXISTS live_rooms;
DROP TABLE IF EXISTS listening_party_participants;
DROP TABLE IF EXISTS listening_parties;