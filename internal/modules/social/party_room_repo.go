package social

import (
	"context"

	"github.com/google/uuid"
)

// --- Listening Parties ---

func (r *repository) CreateParty(ctx context.Context, p *ListeningParty) error {
	p.ID = uuid.New()
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO listening_parties (id, host_id, title, description, cover_url, is_public, status, current_track_id, current_position_ms, started_at)
		VALUES (:id, :host_id, :title, :description, :cover_url, :is_public, :status, :current_track_id, :current_position_ms, :started_at)
	`, p)
	return err
}

func (r *repository) GetParty(ctx context.Context, id uuid.UUID) (*ListeningParty, error) {
	var p ListeningParty
	err := r.db.GetContext(ctx, &p, `
		SELECT lp.*, COUNT(lpp.id) AS participant_count
		FROM listening_parties lp
		LEFT JOIN listening_party_participants lpp ON lpp.party_id = lp.id AND lpp.is_active = true
		WHERE lp.id = $1
		GROUP BY lp.id
	`, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) ListActiveParties(ctx context.Context, limit, offset int) ([]ListeningParty, error) {
	var items []ListeningParty
	err := r.db.SelectContext(ctx, &items, `
		SELECT lp.*, COUNT(lpp.id) AS participant_count
		FROM listening_parties lp
		LEFT JOIN listening_party_participants lpp ON lpp.party_id = lp.id AND lpp.is_active = true
		WHERE lp.status IN ('active', 'paused') AND lp.is_public = true
		GROUP BY lp.id
		ORDER BY lp.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	return items, err
}

func (r *repository) GetPartyHost(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	var hostID uuid.UUID
	err := r.db.GetContext(ctx, &hostID, `SELECT host_id FROM listening_parties WHERE id = $1`, id)
	return hostID, err
}

func (r *repository) GetRoomHost(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	var hostID uuid.UUID
	err := r.db.GetContext(ctx, &hostID, `SELECT host_id FROM live_rooms WHERE id = $1`, id)
	return hostID, err
}

func (r *repository) UpdatePartyStatus(ctx context.Context, id uuid.UUID, status string, trackID *uuid.UUID) error {
	if trackID != nil {
		_, err := r.db.ExecContext(ctx, `
			UPDATE listening_parties
			SET status = $2::text,
			    ended_at = CASE WHEN $2::text = 'ended' THEN NOW() ELSE ended_at END,
			    current_track_id = $3
			WHERE id = $1
		`, id, status, *trackID)
		return err
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE listening_parties
		SET status = $2::text,
		    ended_at = CASE WHEN $2::text = 'ended' THEN NOW() ELSE ended_at END
		WHERE id = $1
	`, id, status)
	return err
}

func (r *repository) UpdatePartyPosition(ctx context.Context, id uuid.UUID, trackID *uuid.UUID, positionMs int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE listening_parties SET current_track_id = $2, current_position_ms = $3
		WHERE id = $1
	`, id, trackID, positionMs)
	return err
}

func (r *repository) JoinParty(ctx context.Context, partyID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO listening_party_participants (party_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (party_id, user_id) DO UPDATE SET is_active = true, left_at = NULL
	`, partyID, userID)
	return err
}

func (r *repository) LeaveParty(ctx context.Context, partyID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE listening_party_participants
		SET is_active = false, left_at = NOW()
		WHERE party_id = $1 AND user_id = $2
	`, partyID, userID)
	return err
}

func (r *repository) GetPartyParticipants(ctx context.Context, partyID uuid.UUID) ([]ListeningPartyParticipant, error) {
	var items []ListeningPartyParticipant
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM listening_party_participants
		WHERE party_id = $1 AND is_active = true
		ORDER BY joined_at ASC
	`, partyID)
	return items, err
}

// --- Live Rooms ---

func (r *repository) CreateRoom(ctx context.Context, room *LiveRoom) error {
	if room.ID == uuid.Nil {
		room.ID = uuid.New()
	}
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO live_rooms (id, host_id, title, description, cover_url, is_public, status, created_at)
		VALUES (:id, :host_id, :title, :description, :cover_url, :is_public, :status, :created_at)
	`, room)
	return err
}

func (r *repository) GetRoom(ctx context.Context, id uuid.UUID) (*LiveRoom, error) {
	var room LiveRoom
	err := r.db.GetContext(ctx, &room, `SELECT * FROM live_rooms WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *repository) ListActiveRooms(ctx context.Context, limit, offset int) ([]LiveRoom, error) {
	var items []LiveRoom
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM live_rooms
		WHERE status = 'live' AND is_public = true
		ORDER BY listener_count DESC, created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	return items, err
}

func (r *repository) UpdateRoomStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE live_rooms SET status = $2, ended_at = CASE WHEN $2 = 'ended' THEN NOW() ELSE ended_at END
		WHERE id = $1
	`, id, status)
	return err
}

func (r *repository) JoinRoom(ctx context.Context, roomID, userID uuid.UUID, role string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO live_room_participants (room_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (room_id, user_id) DO UPDATE SET is_active = true, left_at = NULL, role = $3
	`, roomID, userID, role)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `UPDATE live_rooms SET listener_count = (SELECT COUNT(*) FROM live_room_participants WHERE room_id = $1 AND is_active = true) WHERE id = $1`, roomID)
	return err
}

func (r *repository) LeaveRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE live_room_participants
		SET is_active = false, left_at = NOW()
		WHERE room_id = $1 AND user_id = $2
	`, roomID, userID)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `UPDATE live_rooms SET listener_count = (SELECT COUNT(*) FROM live_room_participants WHERE room_id = $1 AND is_active = true) WHERE id = $1`, roomID)
	return err
}

func (r *repository) GetRoomParticipants(ctx context.Context, roomID uuid.UUID) ([]LiveRoomParticipant, error) {
	var items []LiveRoomParticipant
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM live_room_participants
		WHERE room_id = $1 AND is_active = true
		ORDER BY role ASC, joined_at ASC
	`, roomID)
	return items, err
}

func (r *repository) AddToRoomQueue(ctx context.Context, item *LiveRoomQueueItem) error {
	item.ID = uuid.New()
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO live_room_queue (id, room_id, track_id, added_by, position, status, created_at)
		VALUES (:id, :room_id, :track_id, :added_by, :position, :status, :created_at)
	`, item)
	return err
}

func (r *repository) GetRoomQueue(ctx context.Context, roomID uuid.UUID) ([]LiveRoomQueueItem, error) {
	var items []LiveRoomQueueItem
	err := r.db.SelectContext(ctx, &items, `
		SELECT * FROM live_room_queue
		WHERE room_id = $1 AND status = 'queued'
		ORDER BY position ASC
	`, roomID)
	return items, err
}
