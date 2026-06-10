package social

import (
	"time"

	"github.com/google/uuid"
	"music/internal/platform/ws"
)

// PartyBroadcaster sends real-time listening party events
type PartyBroadcaster struct {
	hub *ws.Hub
}

func NewPartyBroadcaster(hub *ws.Hub) *PartyBroadcaster {
	return &PartyBroadcaster{hub: hub}
}

func (b *PartyBroadcaster) PositionUpdate(partyID uuid.UUID, userID uuid.UUID, trackID *uuid.UUID, positionMs int64, paused bool) {
	b.hub.BroadcastToChannel("party:"+partyID.String(), ws.Message{
		Type: "party.position",
		Payload: map[string]interface{}{
			"user_id":     userID.String(),
			"track_id":    trackID,
			"position_ms": positionMs,
			"paused":      paused,
			"timestamp":   time.Now().UTC(),
		},
		Timestamp: time.Now().UTC(),
	}, uuid.Nil)
}

func (b *PartyBroadcaster) UserJoined(partyID uuid.UUID, userID uuid.UUID) {
	b.hub.BroadcastToChannel("party:"+partyID.String(), ws.Message{
		Type: "party.user_joined",
		Payload: map[string]interface{}{
			"user_id": userID.String(),
		},
		Timestamp: time.Now().UTC(),
	}, uuid.Nil)
}

func (b *PartyBroadcaster) UserLeft(partyID uuid.UUID, userID uuid.UUID) {
	b.hub.BroadcastToChannel("party:"+partyID.String(), ws.Message{
		Type: "party.user_left",
		Payload: map[string]interface{}{
			"user_id": userID.String(),
		},
		Timestamp: time.Now().UTC(),
	}, uuid.Nil)
}

func (b *PartyBroadcaster) TrackChanged(partyID uuid.UUID, trackID uuid.UUID, positionMs int64) {
	b.hub.BroadcastToChannel("party:"+partyID.String(), ws.Message{
		Type: "party.track_changed",
		Payload: map[string]interface{}{
			"track_id":    trackID.String(),
			"position_ms": positionMs,
		},
		Timestamp: time.Now().UTC(),
	}, uuid.Nil)
}

func (b *PartyBroadcaster) PartyEnded(partyID uuid.UUID) {
	b.hub.BroadcastToChannel("party:"+partyID.String(), ws.Message{
		Type:      "party.ended",
		Payload:   map[string]interface{}{},
		Timestamp: time.Now().UTC(),
	}, uuid.Nil)
}

// RoomBroadcaster sends real-time live room events
type RoomBroadcaster struct {
	hub *ws.Hub
}

func NewRoomBroadcaster(hub *ws.Hub) *RoomBroadcaster {
	return &RoomBroadcaster{hub: hub}
}

func (b *RoomBroadcaster) ParticipantJoined(roomID uuid.UUID, userID uuid.UUID, role string) {
	b.hub.BroadcastToChannel("room:"+roomID.String(), ws.Message{
		Type: "room.participant_joined",
		Payload: map[string]interface{}{
			"user_id": userID.String(),
			"role":    role,
		},
		Timestamp: time.Now().UTC(),
	}, uuid.Nil)
}

func (b *RoomBroadcaster) ParticipantLeft(roomID uuid.UUID, userID uuid.UUID) {
	b.hub.BroadcastToChannel("room:"+roomID.String(), ws.Message{
		Type: "room.participant_left",
		Payload: map[string]interface{}{
			"user_id": userID.String(),
		},
		Timestamp: time.Now().UTC(),
	}, uuid.Nil)
}

func (b *RoomBroadcaster) QueueUpdated(roomID uuid.UUID) {
	b.hub.BroadcastToChannel("room:"+roomID.String(), ws.Message{
		Type:      "room.queue_updated",
		Payload:   map[string]interface{}{},
		Timestamp: time.Now().UTC(),
	}, uuid.Nil)
}

func (b *RoomBroadcaster) TrackChanged(roomID uuid.UUID, trackID uuid.UUID) {
	b.hub.BroadcastToChannel("room:"+roomID.String(), ws.Message{
		Type: "room.track_changed",
		Payload: map[string]interface{}{
			"track_id": trackID.String(),
		},
		Timestamp: time.Now().UTC(),
	}, uuid.Nil)
}

func (b *RoomBroadcaster) RoomEnded(roomID uuid.UUID) {
	b.hub.BroadcastToChannel("room:"+roomID.String(), ws.Message{
		Type:      "room.ended",
		Payload:   map[string]interface{}{},
		Timestamp: time.Now().UTC(),
	}, uuid.Nil)
}

func (b *RoomBroadcaster) MessageSent(roomID uuid.UUID, userID uuid.UUID, content string) {
	b.hub.BroadcastToChannel("room:"+roomID.String(), ws.Message{
		Type: "room.message",
		Payload: map[string]interface{}{
			"user_id": userID.String(),
			"content": content,
		},
		Timestamp: time.Now().UTC(),
	}, uuid.Nil)
}
