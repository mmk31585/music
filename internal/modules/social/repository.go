package social

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Follow(ctx context.Context, followerID, followedID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followedID uuid.UUID) error
	IsFollowing(ctx context.Context, followerID, followedID uuid.UUID) (bool, error)
	GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]UserFollow, int, error)
	GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]UserFollow, int, error)
	GetFollowerCount(ctx context.Context, userID uuid.UUID) (int, error)
	GetFollowingCount(ctx context.Context, userID uuid.UUID) (int, error)
	InsertActivity(ctx context.Context, activity Activity) error
	GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int, types []string) ([]ActivityFeedItem, error)
	CreateParty(ctx context.Context, p *ListeningParty) error
	GetParty(ctx context.Context, id uuid.UUID) (*ListeningParty, error)
	ListActiveParties(ctx context.Context, limit, offset int) ([]ListeningParty, error)
	UpdatePartyStatus(ctx context.Context, id uuid.UUID, status string, trackID *uuid.UUID) error
	UpdatePartyPosition(ctx context.Context, id uuid.UUID, trackID *uuid.UUID, positionMs int64) error
	JoinParty(ctx context.Context, partyID, userID uuid.UUID) error
	LeaveParty(ctx context.Context, partyID, userID uuid.UUID) error
	GetPartyParticipants(ctx context.Context, partyID uuid.UUID) ([]ListeningPartyParticipant, error)
	CreateRoom(ctx context.Context, room *LiveRoom) error
	GetRoom(ctx context.Context, id uuid.UUID) (*LiveRoom, error)
	ListActiveRooms(ctx context.Context, limit, offset int) ([]LiveRoom, error)
	UpdateRoomStatus(ctx context.Context, id uuid.UUID, status string) error
	JoinRoom(ctx context.Context, roomID, userID uuid.UUID, role string) error
	LeaveRoom(ctx context.Context, roomID, userID uuid.UUID) error
	GetRoomParticipants(ctx context.Context, roomID uuid.UUID) ([]LiveRoomParticipant, error)
	AddToRoomQueue(ctx context.Context, item *LiveRoomQueueItem) error
	GetRoomQueue(ctx context.Context, roomID uuid.UUID) ([]LiveRoomQueueItem, error)
	CreateClub(ctx context.Context, club *MusicClub) error
	GetClub(ctx context.Context, id uuid.UUID) (*MusicClub, error)
	ListClubs(ctx context.Context, limit, offset int) ([]MusicClub, error)
	ListClubsByGenre(ctx context.Context, genre string, limit, offset int) ([]MusicClub, error)
	UpdateClubPlaylistID(ctx context.Context, clubID, playlistID uuid.UUID) error
	JoinClub(ctx context.Context, clubID, userID uuid.UUID) error
	LeaveClub(ctx context.Context, clubID, userID uuid.UUID) error
	IsClubMember(ctx context.Context, clubID, userID uuid.UUID) (bool, error)
	IsClubAdmin(ctx context.Context, clubID, userID uuid.UUID) (bool, error)
	GetClubMembers(ctx context.Context, clubID uuid.UUID) ([]MusicClubMember, error)
	CreateClubPost(ctx context.Context, post *MusicClubPost) error
	GetClubPosts(ctx context.Context, clubID uuid.UUID, limit, offset int) ([]MusicClubPost, error)
	CreateDiscussion(ctx context.Context, d *Discussion) error
	GetDiscussions(ctx context.Context, targetType string, targetID uuid.UUID, limit, offset int) ([]Discussion, error)
	GetDiscussionReplies(ctx context.Context, parentID uuid.UUID) ([]Discussion, error)
	CreateRating(ctx context.Context, rating *TrackRating) error
	GetTrackRatings(ctx context.Context, trackID uuid.UUID) ([]TrackRating, error)
	GetTrackRatingAverage(ctx context.Context, trackID uuid.UUID) (float64, int, error)

	// Room queue voting
	SuggestTrack(ctx context.Context, c *QueueCandidate) error
	GetCandidates(ctx context.Context, roomID uuid.UUID) ([]QueueCandidate, error)
	GetCandidateByID(ctx context.Context, candidateID uuid.UUID) (*QueueCandidate, error)
	CastVoteTx(ctx context.Context, candidateID, userID uuid.UUID) error
	RemoveVoteTx(ctx context.Context, candidateID, userID uuid.UUID) error
	HasVoted(ctx context.Context, candidateID, userID uuid.UUID) (bool, error)
	GetNowPlaying(ctx context.Context, roomID uuid.UUID) (*RoomNowPlaying, error)
	SetNowPlaying(ctx context.Context, np *RoomNowPlaying) error
	RemoveCandidate(ctx context.Context, candidateID uuid.UUID) error
	GetCandidatesWithVoteState(ctx context.Context, roomID, userID uuid.UUID) ([]CandidateWithVoteState, error)
	LockRoomQueue(ctx context.Context, roomID uuid.UUID) (func(), error)
	PickRandomTrack(ctx context.Context) (*uuid.UUID, error)
	GetRoomMembers(ctx context.Context, roomID uuid.UUID) (int, error)

	// Enriched queue queries (JOIN with tracks/artists/albums/users)
	GetNowPlayingWithTrack(ctx context.Context, roomID uuid.UUID) (*NowPlayingResponse, error)
	GetCandidatesWithTrackAndUser(ctx context.Context, roomID, userID uuid.UUID) ([]CandidateResponse, error)

	// Ownership checks
	GetPartyHost(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	GetRoomHost(ctx context.Context, id uuid.UUID) (uuid.UUID, error)

	// Stage & Raise-Hand
	GetHandRaise(ctx context.Context, roomID, userID uuid.UUID) (*HandRaise, error)
	RaiseHand(ctx context.Context, roomID, userID uuid.UUID) error
	LowerHand(ctx context.Context, roomID, userID uuid.UUID) error
	UpdateHandRaiseStatus(ctx context.Context, id uuid.UUID, status string) error
	ListPendingHandRaises(ctx context.Context, roomID uuid.UUID) ([]HandRaise, error)
	GetStageMember(ctx context.Context, roomID, userID uuid.UUID) (*StageMember, error)
	SetStageMember(ctx context.Context, m *StageMember) error
	RemoveStageMember(ctx context.Context, roomID, userID uuid.UUID) error
	ListStageSpeakers(ctx context.Context, roomID uuid.UUID) ([]StageMember, error)
	GetStageHost(ctx context.Context, roomID uuid.UUID) (*StageMember, error)
	UpdateStageMemberMuted(ctx context.Context, roomID, userID uuid.UUID, muted bool) error

	// Club Discussions (Phase 6)
	CreateClubDiscussion(ctx context.Context, d *ClubDiscussion) error
	ListClubDiscussions(ctx context.Context, clubID uuid.UUID, limit, offset int) ([]ClubDiscussion, error)
	GetClubDiscussion(ctx context.Context, id uuid.UUID) (*ClubDiscussion, error)
	DeleteClubDiscussion(ctx context.Context, id uuid.UUID) error
	IncrementClubDiscussionReplyCount(ctx context.Context, id uuid.UUID) error
	DecrementClubDiscussionReplyCount(ctx context.Context, id uuid.UUID) error
	CreateClubDiscussionReply(ctx context.Context, r *ClubDiscussionReply) error
	GetClubDiscussionReplies(ctx context.Context, discussionID uuid.UUID) ([]ClubDiscussionReply, error)
	GetClubDiscussionReply(ctx context.Context, id uuid.UUID) (*ClubDiscussionReply, error)
	DeleteClubDiscussionReply(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}
