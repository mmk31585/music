-- +goose Up
-- Live Room Stage & Raise-Hand system
-- Migration 000030

CREATE TABLE IF NOT EXISTS room_hand_raises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES live_rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','approved','denied')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(room_id, user_id)
);

CREATE TABLE IF NOT EXISTS room_stage_members (
    room_id UUID NOT NULL REFERENCES live_rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK (role IN ('host','speaker','listener')),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    muted BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (room_id, user_id)
);

CREATE INDEX idx_room_hand_raises_room ON room_hand_raises(room_id, status);
CREATE INDEX idx_room_stage_members_room ON room_stage_members(room_id);

-- +goose Down
DROP TABLE IF EXISTS room_stage_members;
DROP TABLE IF EXISTS room_hand_raises;
