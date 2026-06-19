package social

import (
	"time"

	"github.com/google/uuid"
)

type QueueCandidate struct {
	ID          uuid.UUID `json:"id" db:"id"`
	RoomID      uuid.UUID `json:"room_id" db:"room_id"`
	TrackID     uuid.UUID `json:"track_id" db:"track_id"`
	SuggestedBy uuid.UUID `json:"suggested_by" db:"suggested_by"`
	VoteCount   int       `json:"vote_count" db:"vote_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type QueueVote struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CandidateID uuid.UUID `json:"candidate_id" db:"candidate_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type RoomNowPlaying struct {
	RoomID      uuid.UUID  `json:"room_id" db:"room_id"`
	TrackID     uuid.UUID  `json:"track_id" db:"track_id"`
	StartedAt   time.Time  `json:"started_at" db:"started_at"`
	SuggestedBy *uuid.UUID `json:"suggested_by,omitempty" db:"suggested_by"`
	Source      string     `json:"source" db:"source"`
}

type QueueStateResponse struct {
	NowPlaying *RoomNowPlaying          `json:"now_playing"`
	Candidates []CandidateWithVoteState `json:"candidates"`
}

type CandidateWithVoteState struct {
	QueueCandidate
	HasVoted bool `json:"has_voted" db:"has_voted"`
}
