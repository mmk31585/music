package social

import (
	"context"

	"github.com/google/uuid"
)

// --- Stage & Raise-Hand ---

func (r *repository) GetHandRaise(ctx context.Context, roomID, userID uuid.UUID) (*HandRaise, error) {
	var hr HandRaise
	err := r.db.GetContext(ctx, &hr, `
		SELECT id, room_id, user_id, status, created_at
		FROM room_hand_raises
		WHERE room_id = $1 AND user_id = $2
	`, roomID, userID)
	if err != nil {
		return nil, err
	}
	return &hr, nil
}

func (r *repository) RaiseHand(ctx context.Context, roomID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO room_hand_raises (room_id, user_id, status)
		VALUES ($1, $2, 'pending')
		ON CONFLICT (room_id, user_id) DO NOTHING
	`, roomID, userID)
	return err
}

func (r *repository) LowerHand(ctx context.Context, roomID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM room_hand_raises
		WHERE room_id = $1 AND user_id = $2 AND status = 'pending'
	`, roomID, userID)
	return err
}

func (r *repository) UpdateHandRaiseStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE room_hand_raises SET status = $2 WHERE id = $1
	`, id, status)
	return err
}

func (r *repository) ListPendingHandRaises(ctx context.Context, roomID uuid.UUID) ([]HandRaise, error) {
	var items []HandRaise
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, room_id, user_id, status, created_at
		FROM room_hand_raises
		WHERE room_id = $1 AND status = 'pending'
		ORDER BY created_at ASC
	`, roomID)
	return items, err
}

func (r *repository) GetStageMember(ctx context.Context, roomID, userID uuid.UUID) (*StageMember, error) {
	var m StageMember
	err := r.db.GetContext(ctx, &m, `
		SELECT room_id, user_id, role, joined_at, muted
		FROM room_stage_members
		WHERE room_id = $1 AND user_id = $2
	`, roomID, userID)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) SetStageMember(ctx context.Context, m *StageMember) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO room_stage_members (room_id, user_id, role, joined_at, muted)
		VALUES (:room_id, :user_id, :role, :joined_at, :muted)
		ON CONFLICT (room_id, user_id) DO UPDATE SET
			role = EXCLUDED.role,
			muted = EXCLUDED.muted,
			joined_at = EXCLUDED.joined_at
	`, m)
	return err
}

func (r *repository) RemoveStageMember(ctx context.Context, roomID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM room_stage_members
		WHERE room_id = $1 AND user_id = $2
	`, roomID, userID)
	return err
}

func (r *repository) ListStageSpeakers(ctx context.Context, roomID uuid.UUID) ([]StageMember, error) {
	var items []StageMember
	err := r.db.SelectContext(ctx, &items, `
		SELECT room_id, user_id, role, joined_at, muted
		FROM room_stage_members
		WHERE room_id = $1 AND role = 'speaker'
		ORDER BY joined_at ASC
	`, roomID)
	return items, err
}

func (r *repository) GetStageHost(ctx context.Context, roomID uuid.UUID) (*StageMember, error) {
	var m StageMember
	err := r.db.GetContext(ctx, &m, `
		SELECT room_id, user_id, role, joined_at, muted
		FROM room_stage_members
		WHERE room_id = $1 AND role = 'host'
	`, roomID)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) UpdateStageMemberMuted(ctx context.Context, roomID, userID uuid.UUID, muted bool) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE room_stage_members SET muted = $3
		WHERE room_id = $1 AND user_id = $2
	`, roomID, userID, muted)
	return err
}
