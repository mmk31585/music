package social

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func newStageManager(m *mockRepo) *StageManager {
	return NewStageManager(m, nil, nil, zap.NewNop())
}

func TestStageManager_RaiseHand_Idempotent(t *testing.T) {
	m := new(mockRepo)
	sm := newStageManager(m)
	roomID := uuid.New().String()
	userID := uuid.New().String()

	m.On("RaiseHand", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := sm.RaiseHand(context.Background(), roomID, userID)
	assert.NoError(t, err)

	// Calling again should not error (ON CONFLICT DO NOTHING)
	m.On("RaiseHand", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	err = sm.RaiseHand(context.Background(), roomID, userID)
	assert.NoError(t, err)

	m.AssertExpectations(t)
}

func TestStageManager_ApproveHand_RequiresHost(t *testing.T) {
	m := new(mockRepo)
	sm := newStageManager(m)
	roomID := uuid.New()
	hostID := uuid.New()
	nonHostID := uuid.New()
	targetUserID := uuid.New()

	// Non-host is not the room host
	m.On("GetRoom", mock.Anything, roomID).Return(&LiveRoom{
		ID:     roomID,
		HostID: hostID,
	}, nil)

	err := sm.ApproveHand(context.Background(), roomID.String(), nonHostID.String(), targetUserID.String())
	assert.Error(t, err)
	assert.Equal(t, ErrNotHost, err)
	m.AssertExpectations(t)
}

func TestStageManager_ApproveHand_PromotesToSpeaker(t *testing.T) {
	m := new(mockRepo)
	sm := newStageManager(m)
	roomID := uuid.New()
	hostID := uuid.New()
	targetUserID := uuid.New()
	handRaiseID := uuid.New()

	m.On("GetRoom", mock.Anything, roomID).Return(&LiveRoom{
		ID:     roomID,
		HostID: hostID,
	}, nil)
	m.On("GetHandRaise", mock.Anything, roomID, targetUserID).Return(&HandRaise{
		ID:     handRaiseID,
		RoomID: roomID,
		UserID: targetUserID,
		Status: "pending",
	}, nil)
	m.On("UpdateHandRaiseStatus", mock.Anything, handRaiseID, "approved").Return(nil)
	m.On("SetStageMember", mock.Anything, mock.MatchedBy(func(member *StageMember) bool {
		return member.RoomID == roomID && member.UserID == targetUserID && member.Role == RoleSpeaker
	})).Return(nil)

	err := sm.ApproveHand(context.Background(), roomID.String(), hostID.String(), targetUserID.String())
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestStageManager_RemoveFromStage_RequiresHost(t *testing.T) {
	m := new(mockRepo)
	sm := newStageManager(m)
	roomID := uuid.New()
	hostID := uuid.New()
	nonHostID := uuid.New()
	targetUserID := uuid.New()

	m.On("GetRoom", mock.Anything, roomID).Return(&LiveRoom{
		ID:     roomID,
		HostID: hostID,
	}, nil)

	err := sm.RemoveFromStage(context.Background(), roomID.String(), nonHostID.String(), targetUserID.String())
	assert.Error(t, err)
	assert.Equal(t, ErrNotHost, err)
	m.AssertExpectations(t)
}

func TestStageManager_PendingRaisesOnlyVisibleToHost(t *testing.T) {
	m := new(mockRepo)
	sm := newStageManager(m)
	roomID := uuid.New()
	hostID := uuid.New()
	listenerID := uuid.New()
	now := time.Now().UTC()

	pending := []HandRaise{
		{RoomID: roomID, UserID: listenerID, Status: "pending", CreatedAt: now},
	}

	// Host sees pending raises
	m.On("GetRoom", mock.Anything, roomID).Return(&LiveRoom{
		ID:        roomID,
		HostID:    hostID,
		CreatedAt: now,
	}, nil)
	m.On("ListStageSpeakers", mock.Anything, roomID).Return([]StageMember{}, nil)
	m.On("ListPendingHandRaises", mock.Anything, roomID).Return(pending, nil)

	hostState, err := sm.GetStageState(context.Background(), roomID.String(), hostID.String())
	assert.NoError(t, err)
	assert.NotNil(t, hostState)
	assert.Equal(t, hostID, hostState.Host.UserID)
	assert.Len(t, hostState.PendingRequests, 1)

	// Non-host should not see pending raises
	m.On("GetRoom", mock.Anything, roomID).Return(&LiveRoom{
		ID:        roomID,
		HostID:    hostID,
		CreatedAt: now,
	}, nil)
	m.On("ListStageSpeakers", mock.Anything, roomID).Return([]StageMember{}, nil)

	listenerState, err := sm.GetStageState(context.Background(), roomID.String(), listenerID.String())
	assert.NoError(t, err)
	assert.NotNil(t, listenerState)
	assert.Empty(t, listenerState.PendingRequests)

	m.AssertExpectations(t)
}

func TestStageManager_LeaveStage_Success(t *testing.T) {
	m := new(mockRepo)
	sm := newStageManager(m)
	roomID := uuid.New()
	userID := uuid.New()

	m.On("GetStageMember", mock.Anything, roomID, userID).Return(&StageMember{
		RoomID: roomID,
		UserID: userID,
		Role:   RoleSpeaker,
	}, nil)
	m.On("RemoveStageMember", mock.Anything, roomID, userID).Return(nil)

	err := sm.LeaveStage(context.Background(), roomID.String(), userID.String())
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestStageManager_ToggleMute_SelfAllowed(t *testing.T) {
	m := new(mockRepo)
	sm := newStageManager(m)
	roomID := uuid.New()
	userID := uuid.New()

	m.On("UpdateStageMemberMuted", mock.Anything, roomID, userID, true).Return(nil)

	err := sm.ToggleMute(context.Background(), roomID.String(), userID.String(), userID.String(), true)
	assert.NoError(t, err)
	m.AssertExpectations(t)
}
