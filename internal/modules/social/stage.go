package social

import (
	"time"

	"github.com/google/uuid"
)

type StageRole string

const (
	RoleHost     StageRole = "host"
	RoleSpeaker  StageRole = "speaker"
	RoleListener StageRole = "listener"
)

type HandRaise struct {
	ID        uuid.UUID `json:"id" db:"id"`
	RoomID    uuid.UUID `json:"room_id" db:"room_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type StageMember struct {
	RoomID   uuid.UUID `json:"room_id" db:"room_id"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Role     StageRole `json:"role" db:"role"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
	Muted    bool      `json:"muted" db:"muted"`
}

type StageStateResponse struct {
	Host            *StageMember  `json:"host"`
	Speakers        []StageMember `json:"speakers"`
	PendingRequests []HandRaise   `json:"pending_requests,omitempty"`
}
