package social

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type UserFollow struct {
	FollowerID uuid.UUID `db:"follower_id" json:"follower_id"`
	FolloweeID uuid.UUID `db:"followee_id" json:"followed_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type Activity struct {
	ID         uuid.UUID       `db:"id" json:"id"`
	UserID     uuid.UUID       `db:"user_id" json:"user_id"`
	Type       string          `db:"type" json:"type"`
	TargetID   string          `db:"target_id" json:"target_id"`
	TargetType string          `db:"target_type" json:"target_type"`
	Metadata   json.RawMessage `db:"metadata" json:"metadata,omitempty"`
	CreatedAt  time.Time       `db:"created_at" json:"created_at"`
}

type ActivityFeedItem struct {
	Activity
	UserDisplayName string `db:"user_display_name" json:"user_display_name"`
	UserAvatarURL   string `db:"user_avatar_url" json:"user_avatar_url,omitempty"`
	TargetName      string `db:"target_name" json:"target_name,omitempty"`
	TargetImageURL  string `db:"target_image_url" json:"target_image_url,omitempty"`
}

// ListeningParty
type ListeningParty struct {
	ID                uuid.UUID  `db:"id" json:"id"`
	HostID            uuid.UUID  `db:"host_id" json:"host_id"`
	Title             string     `db:"title" json:"title"`
	Description       *string    `db:"description" json:"description,omitempty"`
	CoverURL          *string    `db:"cover_url" json:"cover_url,omitempty"`
	IsPublic          bool       `db:"is_public" json:"is_public"`
	Status            string     `db:"status" json:"status"`
	CurrentTrackID    *uuid.UUID `db:"current_track_id" json:"current_track_id,omitempty"`
	CurrentPositionMs int64      `db:"current_position_ms" json:"current_position_ms"`
	StartedAt         time.Time  `db:"started_at" json:"started_at"`
	EndedAt           *time.Time `db:"ended_at" json:"ended_at,omitempty"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	ParticipantCount  int        `db:"participant_count" json:"participant_count,omitempty"`
}

type ListeningPartyParticipant struct {
	ID       uuid.UUID  `db:"id" json:"id"`
	PartyID  uuid.UUID  `db:"party_id" json:"party_id"`
	UserID   uuid.UUID  `db:"user_id" json:"user_id"`
	JoinedAt time.Time  `db:"joined_at" json:"joined_at"`
	LeftAt   *time.Time `db:"left_at" json:"left_at,omitempty"`
	IsActive bool       `db:"is_active" json:"is_active"`
}

// LiveRoom
type LiveRoom struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	HostID         uuid.UUID  `db:"host_id" json:"host_id"`
	Title          string     `db:"title" json:"title"`
	Description    *string    `db:"description" json:"description,omitempty"`
	CoverURL       *string    `db:"cover_url" json:"cover_url,omitempty"`
	IsPublic       bool       `db:"is_public" json:"is_public"`
	Status         string     `db:"status" json:"status"`
	CurrentTrackID *uuid.UUID `db:"current_track_id" json:"current_track_id,omitempty"`
	ListenerCount  int        `db:"listener_count" json:"listener_count"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	EndedAt        *time.Time `db:"ended_at" json:"ended_at,omitempty"`
}

type LiveRoomParticipant struct {
	ID       uuid.UUID  `db:"id" json:"id"`
	RoomID   uuid.UUID  `db:"room_id" json:"room_id"`
	UserID   uuid.UUID  `db:"user_id" json:"user_id"`
	Role     string     `db:"role" json:"role"`
	JoinedAt time.Time  `db:"joined_at" json:"joined_at"`
	LeftAt   *time.Time `db:"left_at" json:"left_at,omitempty"`
	IsActive bool       `db:"is_active" json:"is_active"`
}

type LiveRoomQueueItem struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	RoomID    uuid.UUID  `db:"room_id" json:"room_id"`
	TrackID   uuid.UUID  `db:"track_id" json:"track_id"`
	AddedBy   uuid.UUID  `db:"added_by" json:"added_by"`
	Position  int        `db:"position" json:"position"`
	Status    string     `db:"status" json:"status"`
	PlayedAt  *time.Time `db:"played_at" json:"played_at,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

// MusicClub
type MusicClub struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Slug        string     `db:"slug" json:"slug"`
	Description *string    `db:"description" json:"description,omitempty"`
	CoverURL    *string    `db:"cover_url" json:"cover_url,omitempty"`
	Genre       *string    `db:"genre" json:"genre,omitempty"`
	PlaylistID  *uuid.UUID `db:"playlist_id" json:"playlist_id,omitempty"`
	CreatedBy   uuid.UUID  `db:"created_by" json:"created_by"`
	IsPublic    bool       `db:"is_public" json:"is_public"`
	MaxMembers  int        `db:"max_members" json:"max_members"`
	MemberCount int        `db:"member_count" json:"member_count"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

type MusicClubMember struct {
	ID       uuid.UUID `db:"id" json:"id"`
	ClubID   uuid.UUID `db:"club_id" json:"club_id"`
	UserID   uuid.UUID `db:"user_id" json:"user_id"`
	Role     string    `db:"role" json:"role"`
	JoinedAt time.Time `db:"joined_at" json:"joined_at"`
}

type MusicClubPost struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ClubID    uuid.UUID `db:"club_id" json:"club_id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Discussion
type Discussion struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	UserID     uuid.UUID  `db:"user_id" json:"user_id"`
	TargetType string     `db:"target_type" json:"target_type"`
	TargetID   uuid.UUID  `db:"target_id" json:"target_id"`
	Content    string     `db:"content" json:"content"`
	ParentID   *uuid.UUID `db:"parent_id" json:"parent_id,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}

// TrackRating
type TrackRating struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	TrackID   uuid.UUID `db:"track_id" json:"track_id"`
	Rating    int       `db:"rating" json:"rating"`
	Review    *string   `db:"review" json:"review,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Request/response types
type CreatePartyRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	IsPublic    bool    `json:"is_public"`
	TrackID     *string `json:"track_id"`
}

type CreateRoomRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type CreateClubRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	IsPublic    bool   `json:"is_public"`
	MaxMembers  int    `json:"max_members"`
}

type ClubDetailResponse struct {
	Club       MusicClub         `json:"club"`
	IsMember   bool              `json:"is_member"`
	MemberRole string            `json:"member_role"`
	Members    []MusicClubMember `json:"members"`
	Posts      []MusicClubPost   `json:"posts"`
	PostCount  int               `json:"post_count"`
	TrackCount int               `json:"track_count"`
}

type LaunchPartyRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type CreateDiscussionRequest struct {
	TargetType string `json:"target_type" binding:"required"`
	TargetID   string `json:"target_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
	ParentID   string `json:"parent_id"`
}

type CreateRatingRequest struct {
	TrackID string `json:"track_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=10"`
	Review  string `json:"review"`
}

type ClubPostRequest struct {
	Content string `json:"content" binding:"required"`
}

type AddToQueueRequest struct {
	TrackID string `json:"track_id" binding:"required"`
}
