-- +goose Up
-- Democratic voting queue for live rooms and listening parties
-- Migration 000029

CREATE TABLE IF NOT EXISTS room_queue_candidates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES live_rooms(id) ON DELETE CASCADE,
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    suggested_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vote_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(room_id, track_id)
);
CREATE INDEX idx_room_queue_candidates_room_votes ON room_queue_candidates(room_id, vote_count DESC, created_at ASC);

CREATE TABLE IF NOT EXISTS room_queue_votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    candidate_id UUID NOT NULL REFERENCES room_queue_candidates(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(candidate_id, user_id)
);

CREATE TABLE IF NOT EXISTS room_now_playing (
    room_id UUID PRIMARY KEY REFERENCES live_rooms(id) ON DELETE CASCADE,
    track_id UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    suggested_by UUID REFERENCES users(id) ON DELETE SET NULL,
    source TEXT NOT NULL CHECK (source IN ('vote', 'autofill'))
);

-- +goose Down
DROP TABLE IF EXISTS room_now_playing;
DROP TABLE IF EXISTS room_queue_votes;
DROP TABLE IF EXISTS room_queue_candidates;
