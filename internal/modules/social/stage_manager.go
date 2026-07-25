package social

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"music/internal/platform/events"
)

var (
	ErrNotHost             = errors.New("فقط میزبان می‌تواند این کار را انجام دهد")
	ErrHandRaiseNotFound   = errors.New("درخواست دست‌بلندکردن یافت نشد")
	ErrAlreadySpeaker      = errors.New("کاربر در حال حاضر سخنران است")
	ErrNotSpeaker          = errors.New("کاربر سخنران نیست")
	ErrStageMemberNotFound = errors.New("عضو استیج یافت نشد")
)

type StageManager struct {
	repo            Repository
	roomBroadcaster *RoomBroadcaster
	eventBus        events.Publisher
	logger          *zap.Logger
}

func NewStageManager(repo Repository, roomBroadcaster *RoomBroadcaster, eventBus events.Publisher, logger *zap.Logger) *StageManager {
	return &StageManager{
		repo:            repo,
		roomBroadcaster: roomBroadcaster,
		eventBus:        eventBus,
		logger:          logger,
	}
}

func (s *StageManager) verifyHost(ctx context.Context, roomID, userID uuid.UUID) error {
	room, err := s.repo.GetRoom(ctx, roomID)
	if err != nil {
		return fmt.Errorf("room not found: %w", err)
	}
	if room.HostID != userID {
		return ErrNotHost
	}
	return nil
}

func (s *StageManager) RaiseHand(ctx context.Context, roomID, userID string) error {
	rid, err := uuid.Parse(roomID)
	if err != nil {
		return fmt.Errorf("invalid room id: %w", err)
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	if err := s.repo.RaiseHand(ctx, rid, uid); err != nil {
		return fmt.Errorf("raise hand: %w", err)
	}

	if s.roomBroadcaster != nil {
		s.roomBroadcaster.HandRaised(rid, uid)
	}
	return nil
}

func (s *StageManager) LowerHand(ctx context.Context, roomID, userID string) error {
	rid, err := uuid.Parse(roomID)
	if err != nil {
		return fmt.Errorf("invalid room id: %w", err)
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	if err := s.repo.LowerHand(ctx, rid, uid); err != nil {
		return fmt.Errorf("lower hand: %w", err)
	}

	if s.roomBroadcaster != nil {
		s.roomBroadcaster.HandLowered(rid, uid)
	}
	return nil
}

func (s *StageManager) ApproveHand(ctx context.Context, roomID, hostID, targetUserID string) error {
	rid, err := uuid.Parse(roomID)
	if err != nil {
		return fmt.Errorf("invalid room id: %w", err)
	}
	hid, err := uuid.Parse(hostID)
	if err != nil {
		return fmt.Errorf("invalid host id: %w", err)
	}
	tuid, err := uuid.Parse(targetUserID)
	if err != nil {
		return fmt.Errorf("invalid target user id: %w", err)
	}

	if err := s.verifyHost(ctx, rid, hid); err != nil {
		return err
	}

	hr, err := s.repo.GetHandRaise(ctx, rid, tuid)
	if err != nil {
		return ErrHandRaiseNotFound
	}

	if err := s.repo.UpdateHandRaiseStatus(ctx, hr.ID, "approved"); err != nil {
		return fmt.Errorf("update hand raise: %w", err)
	}

	member := &StageMember{
		RoomID:   rid,
		UserID:   tuid,
		Role:     RoleSpeaker,
		JoinedAt: time.Now().UTC(),
		Muted:    false,
	}
	if err := s.repo.SetStageMember(ctx, member); err != nil {
		return fmt.Errorf("set stage member: %w", err)
	}

	if s.roomBroadcaster != nil {
		s.roomBroadcaster.HandApproved(rid, tuid)
		s.broadcastStageState(ctx, rid)
	}
	return nil
}

func (s *StageManager) DenyHand(ctx context.Context, roomID, hostID, targetUserID string) error {
	rid, err := uuid.Parse(roomID)
	if err != nil {
		return fmt.Errorf("invalid room id: %w", err)
	}
	hid, err := uuid.Parse(hostID)
	if err != nil {
		return fmt.Errorf("invalid host id: %w", err)
	}
	tuid, err := uuid.Parse(targetUserID)
	if err != nil {
		return fmt.Errorf("invalid target user id: %w", err)
	}

	if err := s.verifyHost(ctx, rid, hid); err != nil {
		return err
	}

	hr, err := s.repo.GetHandRaise(ctx, rid, tuid)
	if err != nil {
		return ErrHandRaiseNotFound
	}

	if err := s.repo.UpdateHandRaiseStatus(ctx, hr.ID, "denied"); err != nil {
		return fmt.Errorf("update hand raise: %w", err)
	}

	if s.roomBroadcaster != nil {
		s.roomBroadcaster.HandDenied(rid, tuid)
	}
	return nil
}

func (s *StageManager) RemoveFromStage(ctx context.Context, roomID, hostID, targetUserID string) error {
	rid, err := uuid.Parse(roomID)
	if err != nil {
		return fmt.Errorf("invalid room id: %w", err)
	}
	hid, err := uuid.Parse(hostID)
	if err != nil {
		return fmt.Errorf("invalid host id: %w", err)
	}
	tuid, err := uuid.Parse(targetUserID)
	if err != nil {
		return fmt.Errorf("invalid target user id: %w", err)
	}

	if err := s.verifyHost(ctx, rid, hid); err != nil {
		return err
	}

	member, err := s.repo.GetStageMember(ctx, rid, tuid)
	if err != nil {
		return ErrStageMemberNotFound
	}
	if member.Role != RoleSpeaker {
		return ErrNotSpeaker
	}

	if err := s.repo.RemoveStageMember(ctx, rid, tuid); err != nil {
		return fmt.Errorf("remove stage member: %w", err)
	}

	if s.roomBroadcaster != nil {
		s.broadcastStageState(ctx, rid)
	}
	return nil
}

func (s *StageManager) LeaveStage(ctx context.Context, roomID, userID string) error {
	rid, err := uuid.Parse(roomID)
	if err != nil {
		return fmt.Errorf("invalid room id: %w", err)
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	member, err := s.repo.GetStageMember(ctx, rid, uid)
	if err != nil {
		return ErrStageMemberNotFound
	}
	if member.Role != RoleSpeaker {
		return ErrNotSpeaker
	}

	if err := s.repo.RemoveStageMember(ctx, rid, uid); err != nil {
		return fmt.Errorf("remove stage member: %w", err)
	}

	if s.roomBroadcaster != nil {
		s.broadcastStageState(ctx, rid)
	}
	return nil
}

func (s *StageManager) ToggleMute(ctx context.Context, roomID, actorID, targetUserID string, muted bool) error {
	rid, err := uuid.Parse(roomID)
	if err != nil {
		return fmt.Errorf("invalid room id: %w", err)
	}
	aid, err := uuid.Parse(actorID)
	if err != nil {
		return fmt.Errorf("invalid actor id: %w", err)
	}
	tuid, err := uuid.Parse(targetUserID)
	if err != nil {
		return fmt.Errorf("invalid target user id: %w", err)
	}

	// Allow host to mute anyone, or speaker to mute self
	if aid != tuid {
		if err := s.verifyHost(ctx, rid, aid); err != nil {
			return err
		}
	}

	if err := s.repo.UpdateStageMemberMuted(ctx, rid, tuid, muted); err != nil {
		return fmt.Errorf("update muted: %w", err)
	}

	if s.roomBroadcaster != nil {
		s.broadcastStageState(ctx, rid)
	}
	return nil
}

func (s *StageManager) GetStageState(ctx context.Context, roomID, requestingUserID string) (*StageStateResponse, error) {
	rid, err := uuid.Parse(roomID)
	if err != nil {
		return nil, fmt.Errorf("invalid room id: %w", err)
	}
	ruid, err := uuid.Parse(requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	room, err := s.repo.GetRoom(ctx, rid)
	if err != nil {
		return nil, fmt.Errorf("room not found: %w", err)
	}

	speakers, err := s.repo.ListStageSpeakers(ctx, rid)
	if err != nil {
		return nil, fmt.Errorf("list speakers: %w", err)
	}

	host := &StageMember{
		RoomID:   rid,
		UserID:   room.HostID,
		Role:     RoleHost,
		JoinedAt: room.CreatedAt,
	}

	resp := &StageStateResponse{
		Host:     host,
		Speakers: speakers,
	}

	// Pending raises only visible to host
	if room.HostID == ruid {
		pending, err := s.repo.ListPendingHandRaises(ctx, rid)
		if err != nil {
			return nil, fmt.Errorf("list pending raises: %w", err)
		}
		resp.PendingRequests = pending
	}

	return resp, nil
}

func (s *StageManager) broadcastStageState(ctx context.Context, roomID uuid.UUID) {
	if s.roomBroadcaster == nil {
		return
	}

	room, err := s.repo.GetRoom(ctx, roomID)
	if err != nil {
		s.logger.Error("stage broadcast: get room", zap.Error(err))
		return
	}

	speakers, err := s.repo.ListStageSpeakers(ctx, roomID)
	if err != nil {
		s.logger.Error("stage broadcast: list speakers", zap.Error(err))
		return
	}

	host := &StageMember{
		RoomID:   roomID,
		UserID:   room.HostID,
		Role:     RoleHost,
		JoinedAt: room.CreatedAt,
	}

	s.roomBroadcaster.StageUpdated(roomID, speakers, host)
}
