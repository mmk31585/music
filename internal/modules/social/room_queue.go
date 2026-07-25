package social

import (
	"time"

	"github.com/google/uuid"
)

// ----- Internal DB types (used for storage queries) -----

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

// ----- Internal flat DTOs for JOIN queries -----

// nowPlayingTrackRow holds raw joined columns for the now-playing query.
type nowPlayingTrackRow struct {
	RoomID          uuid.UUID  `db:"room_id"`
	TrackID         uuid.UUID  `db:"track_id"`
	StartedAt       time.Time  `db:"started_at"`
	SuggestedBy     *uuid.UUID `db:"suggested_by"`
	Source          string     `db:"source"`
	Title           string     `db:"title"`
	DurationSeconds int        `db:"duration_seconds"`
	CoverURL        *string    `db:"cover_url"`
	ArtistName      string     `db:"artist_name"`
	AlbumTitle      *string    `db:"album_title"`
	Username        *string    `db:"username"`
	UserAvatarURL   *string    `db:"user_avatar_url"`
}

// candidateTrackRow holds raw joined columns for the candidates query.
type candidateTrackRow struct {
	ID              uuid.UUID `db:"id"`
	RoomID          uuid.UUID `db:"room_id"`
	TrackID         uuid.UUID `db:"track_id"`
	SuggestedBy     uuid.UUID `db:"suggested_by"`
	VoteCount       int       `db:"vote_count"`
	CreatedAt       time.Time `db:"created_at"`
	HasVoted        bool      `db:"has_voted"`
	Title           string    `db:"title"`
	DurationSeconds int       `db:"duration_seconds"`
	CoverURL        *string   `db:"cover_url"`
	ArtistName      string    `db:"artist_name"`
	AlbumTitle      *string   `db:"album_title"`
	Username        string    `db:"username"`
	UserAvatarURL   *string   `db:"user_avatar_url"`
}

// ----- Public response types (sent to frontend) -----

type TrackSummaryResponse struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	ArtistName      string  `json:"artist_name"`
	AlbumTitle      *string `json:"album_title,omitempty"`
	CoverURL        *string `json:"cover_url,omitempty"`
	DurationSeconds *int    `json:"duration_seconds,omitempty"`
}

type UserSummaryResponse struct {
	ID        string  `json:"id"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type NowPlayingResponse struct {
	Track       TrackSummaryResponse `json:"track"`
	StartedAt   time.Time            `json:"started_at"`
	SuggestedBy *UserSummaryResponse `json:"suggested_by,omitempty"`
	Source      string               `json:"source"`
}

type CandidateResponse struct {
	ID          uuid.UUID            `json:"id"`
	RoomID      uuid.UUID            `json:"room_id"`
	Track       TrackSummaryResponse `json:"track"`
	SuggestedBy UserSummaryResponse  `json:"suggested_by"`
	VoteCount   int                  `json:"vote_count"`
	HasVoted    bool                 `json:"has_voted"`
	CreatedAt   time.Time            `json:"created_at"`
}

type QueueStateResponse struct {
	NowPlaying *NowPlayingResponse `json:"now_playing"`
	Candidates []CandidateResponse `json:"candidates"`
}

// CandidateWithVoteState is retained for backward compat (internal use only).
type CandidateWithVoteState struct {
	QueueCandidate
	HasVoted bool `json:"has_voted" db:"has_voted"`
}
