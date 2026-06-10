package social

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Follow(ctx context.Context, followerID, followedID uuid.UUID) error {
	return m.Called(ctx, followerID, followedID).Error(0)
}
func (m *mockRepo) Unfollow(ctx context.Context, followerID, followedID uuid.UUID) error {
	return m.Called(ctx, followerID, followedID).Error(0)
}
func (m *mockRepo) IsFollowing(ctx context.Context, followerID, followedID uuid.UUID) (bool, error) {
	args := m.Called(ctx, followerID, followedID)
	return args.Bool(0), args.Error(1)
}
func (m *mockRepo) GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]UserFollow, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]UserFollow), args.Int(1), args.Error(2)
}
func (m *mockRepo) GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]UserFollow, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]UserFollow), args.Int(1), args.Error(2)
}
func (m *mockRepo) GetFollowerCount(ctx context.Context, userID uuid.UUID) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}
func (m *mockRepo) GetFollowingCount(ctx context.Context, userID uuid.UUID) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}
func (m *mockRepo) InsertActivity(ctx context.Context, activity Activity) error {
	return m.Called(ctx, activity).Error(0)
}
func (m *mockRepo) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int, types []string) ([]ActivityFeedItem, error) {
	args := m.Called(ctx, userID, limit, offset, types)
	return args.Get(0).([]ActivityFeedItem), args.Error(1)
}
func (m *mockRepo) CreateParty(ctx context.Context, p *ListeningParty) error {
	return m.Called(ctx, p).Error(0)
}
func (m *mockRepo) GetParty(ctx context.Context, id uuid.UUID) (*ListeningParty, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*ListeningParty), args.Error(1)
}
func (m *mockRepo) ListActiveParties(ctx context.Context, limit, offset int) ([]ListeningParty, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]ListeningParty), args.Error(1)
}
func (m *mockRepo) UpdatePartyStatus(ctx context.Context, id uuid.UUID, status string, trackID *uuid.UUID) error {
	return m.Called(ctx, id, status, trackID).Error(0)
}
func (m *mockRepo) UpdatePartyPosition(ctx context.Context, id uuid.UUID, trackID *uuid.UUID, positionMs int64) error {
	return m.Called(ctx, id, trackID, positionMs).Error(0)
}
func (m *mockRepo) JoinParty(ctx context.Context, partyID, userID uuid.UUID) error {
	return m.Called(ctx, partyID, userID).Error(0)
}
func (m *mockRepo) LeaveParty(ctx context.Context, partyID, userID uuid.UUID) error {
	return m.Called(ctx, partyID, userID).Error(0)
}
func (m *mockRepo) GetPartyParticipants(ctx context.Context, partyID uuid.UUID) ([]ListeningPartyParticipant, error) {
	args := m.Called(ctx, partyID)
	return args.Get(0).([]ListeningPartyParticipant), args.Error(1)
}
func (m *mockRepo) CreateRoom(ctx context.Context, room *LiveRoom) error {
	return m.Called(ctx, room).Error(0)
}
func (m *mockRepo) GetRoom(ctx context.Context, id uuid.UUID) (*LiveRoom, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*LiveRoom), args.Error(1)
}
func (m *mockRepo) ListActiveRooms(ctx context.Context, limit, offset int) ([]LiveRoom, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]LiveRoom), args.Error(1)
}
func (m *mockRepo) UpdateRoomStatus(ctx context.Context, id uuid.UUID, status string) error {
	return m.Called(ctx, id, status).Error(0)
}
func (m *mockRepo) JoinRoom(ctx context.Context, roomID, userID uuid.UUID, role string) error {
	return m.Called(ctx, roomID, userID, role).Error(0)
}
func (m *mockRepo) LeaveRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	return m.Called(ctx, roomID, userID).Error(0)
}
func (m *mockRepo) GetRoomParticipants(ctx context.Context, roomID uuid.UUID) ([]LiveRoomParticipant, error) {
	args := m.Called(ctx, roomID)
	return args.Get(0).([]LiveRoomParticipant), args.Error(1)
}
func (m *mockRepo) AddToRoomQueue(ctx context.Context, item *LiveRoomQueueItem) error {
	return m.Called(ctx, item).Error(0)
}
func (m *mockRepo) GetRoomQueue(ctx context.Context, roomID uuid.UUID) ([]LiveRoomQueueItem, error) {
	args := m.Called(ctx, roomID)
	return args.Get(0).([]LiveRoomQueueItem), args.Error(1)
}
func (m *mockRepo) CreateClub(ctx context.Context, club *MusicClub) error {
	return m.Called(ctx, club).Error(0)
}
func (m *mockRepo) GetClub(ctx context.Context, id uuid.UUID) (*MusicClub, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*MusicClub), args.Error(1)
}
func (m *mockRepo) ListClubs(ctx context.Context, limit, offset int) ([]MusicClub, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]MusicClub), args.Error(1)
}
func (m *mockRepo) JoinClub(ctx context.Context, clubID, userID uuid.UUID) error {
	return m.Called(ctx, clubID, userID).Error(0)
}
func (m *mockRepo) LeaveClub(ctx context.Context, clubID, userID uuid.UUID) error {
	return m.Called(ctx, clubID, userID).Error(0)
}
func (m *mockRepo) IsClubMember(ctx context.Context, clubID, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, clubID, userID)
	return args.Bool(0), args.Error(1)
}
func (m *mockRepo) IsClubAdmin(ctx context.Context, clubID, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, clubID, userID)
	return args.Bool(0), args.Error(1)
}
func (m *mockRepo) GetClubMembers(ctx context.Context, clubID uuid.UUID) ([]MusicClubMember, error) {
	args := m.Called(ctx, clubID)
	return args.Get(0).([]MusicClubMember), args.Error(1)
}
func (m *mockRepo) CreateClubPost(ctx context.Context, post *MusicClubPost) error {
	return m.Called(ctx, post).Error(0)
}
func (m *mockRepo) GetClubPosts(ctx context.Context, clubID uuid.UUID, limit, offset int) ([]MusicClubPost, error) {
	args := m.Called(ctx, clubID, limit, offset)
	return args.Get(0).([]MusicClubPost), args.Error(1)
}
func (m *mockRepo) CreateDiscussion(ctx context.Context, d *Discussion) error {
	return m.Called(ctx, d).Error(0)
}
func (m *mockRepo) GetDiscussions(ctx context.Context, targetType string, targetID uuid.UUID, limit, offset int) ([]Discussion, error) {
	args := m.Called(ctx, targetType, targetID, limit, offset)
	return args.Get(0).([]Discussion), args.Error(1)
}
func (m *mockRepo) GetDiscussionReplies(ctx context.Context, parentID uuid.UUID) ([]Discussion, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]Discussion), args.Error(1)
}
func (m *mockRepo) CreateRating(ctx context.Context, rating *TrackRating) error {
	return m.Called(ctx, rating).Error(0)
}
func (m *mockRepo) GetTrackRatings(ctx context.Context, trackID uuid.UUID) ([]TrackRating, error) {
	args := m.Called(ctx, trackID)
	return args.Get(0).([]TrackRating), args.Error(1)
}
func (m *mockRepo) GetTrackRatingAverage(ctx context.Context, trackID uuid.UUID) (float64, int, error) {
	args := m.Called(ctx, trackID)
	return args.Get(0).(float64), args.Int(1), args.Error(2)
}

func newSvc(m *mockRepo) *Service {
	return NewService(m, nil, nil)
}

func TestFollow_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	u1, u2 := uuid.New().String(), uuid.New().String()

	m.On("Follow", mock.Anything, mock.MatchedBy(func(uid uuid.UUID) bool { return true }), mock.Anything).Return(nil)

	err := svc.Follow(context.Background(), u1, u2)
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestFollow_SelfFollow(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	uid := uuid.New().String()

	err := svc.Follow(context.Background(), uid, uid)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot follow yourself")
	m.AssertNotCalled(t, "Follow")
}

func TestFollow_InvalidUUID(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	err := svc.Follow(context.Background(), "not-a-uuid", uuid.New().String())
	assert.Error(t, err)
	m.AssertNotCalled(t, "Follow")
}

func TestUnfollow_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	u1, u2 := uuid.New().String(), uuid.New().String()

	m.On("Unfollow", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := svc.Unfollow(context.Background(), u1, u2)
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestGetFollowers_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	uid := uuid.New().String()

	m.On("GetFollowers", mock.Anything, mock.Anything, 10, 0).Return([]UserFollow{}, 0, nil)

	items, total, err := svc.GetFollowers(context.Background(), uid, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestGetFollowing_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	uid := uuid.New().String()

	m.On("GetFollowing", mock.Anything, mock.Anything, 20, 0).Return([]UserFollow{}, 5, nil)

	items, total, err := svc.GetFollowing(context.Background(), uid, 20, 0)
	assert.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestIsFollowing_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	u1, u2 := uuid.New().String(), uuid.New().String()

	m.On("IsFollowing", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)

	ok, err := svc.IsFollowing(context.Background(), u1, u2)
	assert.NoError(t, err)
	assert.True(t, ok)
	m.AssertExpectations(t)
}

func TestGetFeed_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	uid := uuid.New().String()

	m.On("GetFeed", mock.Anything, mock.Anything, 10, 0, []string(nil)).Return([]ActivityFeedItem{}, nil)

	items, err := svc.GetFeed(context.Background(), uid, 10, 0, "")
	assert.NoError(t, err)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestGetFeed_WithTypeFilter(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	uid := uuid.New().String()

	m.On("GetFeed", mock.Anything, mock.Anything, 10, 0, []string{"follow", "rating"}).Return([]ActivityFeedItem{}, nil)

	items, err := svc.GetFeed(context.Background(), uid, 10, 0, "follow,rating")
	assert.NoError(t, err)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestCreateParty_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	hostID := uuid.New().String()

	m.On("CreateParty", mock.Anything, mock.MatchedBy(func(p *ListeningParty) bool {
		return p.Title == "My Party" && p.IsPublic && p.Status == "active"
	})).Return(nil)
	m.On("JoinParty", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	party, err := svc.CreateParty(context.Background(), CreatePartyRequest{
		Title:    "My Party",
		IsPublic: true,
	}, hostID)

	assert.NoError(t, err)
	assert.Equal(t, "My Party", party.Title)
	assert.Equal(t, "active", party.Status)
	m.AssertExpectations(t)
}

func TestCreateParty_WithOptionalTrack(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	hostID := uuid.New().String()
	trackID := uuid.New().String()

	m.On("CreateParty", mock.Anything, mock.MatchedBy(func(p *ListeningParty) bool {
		return p.CurrentTrackID != nil && p.Title == "Party with Track"
	})).Return(nil)
	m.On("JoinParty", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	party, err := svc.CreateParty(context.Background(), CreatePartyRequest{
		Title:    "Party with Track",
		IsPublic: true,
		TrackID:  &trackID,
	}, hostID)

	assert.NoError(t, err)
	assert.NotNil(t, party.CurrentTrackID)
	m.AssertExpectations(t)
}

func TestGetParty_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	pid := uuid.New()
	party := &ListeningParty{ID: pid, Title: "My Party"}

	m.On("GetParty", mock.Anything, pid).Return(party, nil)

	result, err := svc.GetParty(context.Background(), pid.String())
	assert.NoError(t, err)
	assert.Equal(t, "My Party", result.Title)
	m.AssertExpectations(t)
}

func TestListActiveParties_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	m.On("ListActiveParties", mock.Anything, 10, 0).Return([]ListeningParty{}, nil)

	items, err := svc.ListActiveParties(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestUpdatePartyStatus_NotEnded(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	pid := uuid.New()

	m.On("UpdatePartyStatus", mock.Anything, pid, "paused", (*uuid.UUID)(nil)).Return(nil)

	err := svc.UpdatePartyStatus(context.Background(), pid.String(), "paused", nil)
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestUpdatePartyStatus_Ended(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	pid := uuid.New()

	m.On("UpdatePartyStatus", mock.Anything, pid, "ended", (*uuid.UUID)(nil)).Return(nil)

	err := svc.UpdatePartyStatus(context.Background(), pid.String(), "ended", nil)
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestJoinParty_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	pid, uid := uuid.New(), uuid.New()

	m.On("JoinParty", mock.Anything, pid, uid).Return(nil)

	err := svc.JoinParty(context.Background(), pid.String(), uid.String())
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestLeaveParty_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	pid, uid := uuid.New(), uuid.New()

	m.On("LeaveParty", mock.Anything, pid, uid).Return(nil)

	err := svc.LeaveParty(context.Background(), pid.String(), uid.String())
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestGetPartyParticipants_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	pid := uuid.New()

	m.On("GetPartyParticipants", mock.Anything, pid).Return([]ListeningPartyParticipant{}, nil)

	items, err := svc.GetPartyParticipants(context.Background(), pid.String())
	assert.NoError(t, err)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestCreateRoom_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	hostID := uuid.New().String()

	m.On("CreateRoom", mock.Anything, mock.MatchedBy(func(r *LiveRoom) bool {
		return r.Title == "My Room" && r.Status == "live"
	})).Return(nil)
	m.On("JoinRoom", mock.Anything, mock.Anything, mock.Anything, "host").Return(nil)

	room, err := svc.CreateRoom(context.Background(), CreateRoomRequest{
		Title:    "My Room",
		IsPublic: true,
	}, hostID)

	assert.NoError(t, err)
	assert.Equal(t, "My Room", room.Title)
	assert.Equal(t, "live", room.Status)
	m.AssertExpectations(t)
}

func TestGetRoom_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	rid := uuid.New()

	m.On("GetRoom", mock.Anything, rid).Return(&LiveRoom{ID: rid, Title: "Room"}, nil)

	room, err := svc.GetRoom(context.Background(), rid.String())
	assert.NoError(t, err)
	assert.Equal(t, "Room", room.Title)
	m.AssertExpectations(t)
}

func TestListActiveRooms_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	m.On("ListActiveRooms", mock.Anything, 5, 0).Return([]LiveRoom{}, nil)

	items, err := svc.ListActiveRooms(context.Background(), 5, 0)
	assert.NoError(t, err)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestJoinRoom_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	rid, uid := uuid.New(), uuid.New()

	m.On("JoinRoom", mock.Anything, rid, uid, "listener").Return(nil)

	err := svc.JoinRoom(context.Background(), rid.String(), uid.String())
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestLeaveRoom_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	rid, uid := uuid.New(), uuid.New()

	m.On("LeaveRoom", mock.Anything, rid, uid).Return(nil)

	err := svc.LeaveRoom(context.Background(), rid.String(), uid.String())
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestGetRoomParticipants_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	rid := uuid.New()

	m.On("GetRoomParticipants", mock.Anything, rid).Return([]LiveRoomParticipant{}, nil)

	items, err := svc.GetRoomParticipants(context.Background(), rid.String())
	assert.NoError(t, err)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestAddToRoomQueue_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	rid, tid, uid := uuid.New(), uuid.New(), uuid.New()

	m.On("GetRoomQueue", mock.Anything, rid).Return([]LiveRoomQueueItem{}, nil)
	m.On("AddToRoomQueue", mock.Anything, mock.MatchedBy(func(item *LiveRoomQueueItem) bool {
		return item.RoomID == rid && item.TrackID == tid && item.AddedBy == uid && item.Position == 0
	})).Return(nil)

	err := svc.AddToRoomQueue(context.Background(), rid.String(), tid.String(), uid.String())
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestGetRoomQueue_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	rid := uuid.New()

	m.On("GetRoomQueue", mock.Anything, rid).Return([]LiveRoomQueueItem{}, nil)

	items, err := svc.GetRoomQueue(context.Background(), rid.String())
	assert.NoError(t, err)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestCreateClub_DefaultMaxMembers(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	uid := uuid.New().String()

	m.On("CreateClub", mock.Anything, mock.MatchedBy(func(c *MusicClub) bool {
		return c.Name == "Music Lovers" && c.MaxMembers == 1000
	})).Return(nil)

	club, err := svc.CreateClub(context.Background(), CreateClubRequest{
		Name:     "Music Lovers",
		IsPublic: true,
	}, uid)

	assert.NoError(t, err)
	assert.Equal(t, "Music Lovers", club.Name)
	assert.Equal(t, 1000, club.MaxMembers)
	m.AssertExpectations(t)
}

func TestCreateClub_CustomMaxMembers(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	uid := uuid.New().String()

	m.On("CreateClub", mock.Anything, mock.MatchedBy(func(c *MusicClub) bool {
		return c.MaxMembers == 50
	})).Return(nil)

	club, err := svc.CreateClub(context.Background(), CreateClubRequest{
		Name:       "Small Club",
		IsPublic:   true,
		MaxMembers: 50,
	}, uid)

	assert.NoError(t, err)
	assert.Equal(t, 50, club.MaxMembers)
	m.AssertExpectations(t)
}

func TestGetClub_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	cid := uuid.New()

	m.On("GetClub", mock.Anything, cid).Return(&MusicClub{ID: cid, Name: "Club"}, nil)

	club, err := svc.GetClub(context.Background(), cid.String())
	assert.NoError(t, err)
	assert.Equal(t, "Club", club.Name)
	m.AssertExpectations(t)
}

func TestListClubs_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	m.On("ListClubs", mock.Anything, 10, 0).Return([]MusicClub{}, nil)

	items, err := svc.ListClubs(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestJoinClub_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	cid, uid := uuid.New(), uuid.New()

	m.On("JoinClub", mock.Anything, cid, uid).Return(nil)

	err := svc.JoinClub(context.Background(), cid.String(), uid.String())
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestLeaveClub_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	cid, uid := uuid.New(), uuid.New()

	m.On("LeaveClub", mock.Anything, cid, uid).Return(nil)

	err := svc.LeaveClub(context.Background(), cid.String(), uid.String())
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestIsClubMember_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	cid, uid := uuid.New(), uuid.New()

	m.On("IsClubMember", mock.Anything, cid, uid).Return(true, nil)

	ok, err := svc.IsClubMember(context.Background(), cid.String(), uid.String())
	assert.NoError(t, err)
	assert.True(t, ok)
	m.AssertExpectations(t)
}

func TestIsClubAdmin_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	cid, uid := uuid.New(), uuid.New()

	m.On("IsClubAdmin", mock.Anything, cid, uid).Return(true, nil)

	ok, err := svc.IsClubAdmin(context.Background(), cid.String(), uid.String())
	assert.NoError(t, err)
	assert.True(t, ok)
	m.AssertExpectations(t)
}

func TestGetClubMembers_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	cid := uuid.New()

	m.On("GetClubMembers", mock.Anything, cid).Return([]MusicClubMember{}, nil)

	items, err := svc.GetClubMembers(context.Background(), cid.String())
	assert.NoError(t, err)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestCreateClubPost_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	cid, uid := uuid.New(), uuid.New()

	m.On("CreateClubPost", mock.Anything, mock.MatchedBy(func(p *MusicClubPost) bool {
		return p.Content == "Hello!" && p.ClubID == cid && p.UserID == uid
	})).Return(nil)

	post, err := svc.CreateClubPost(context.Background(), cid.String(), uid.String(), "Hello!")
	assert.NoError(t, err)
	assert.Equal(t, "Hello!", post.Content)
	m.AssertExpectations(t)
}

func TestGetClubPosts_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	cid := uuid.New()

	m.On("GetClubPosts", mock.Anything, cid, 10, 0).Return([]MusicClubPost{}, nil)

	items, err := svc.GetClubPosts(context.Background(), cid.String(), 10, 0)
	assert.NoError(t, err)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestCreateDiscussion_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	uid, tid := uuid.New(), uuid.New()

	m.On("CreateDiscussion", mock.Anything, mock.MatchedBy(func(d *Discussion) bool {
		return d.Content == "Great track!" && d.TargetType == "track" && d.TargetID == tid && d.ParentID == nil
	})).Return(nil)

	d, err := svc.CreateDiscussion(context.Background(), CreateDiscussionRequest{
		TargetType: "track",
		TargetID:   tid.String(),
		Content:    "Great track!",
	}, uid.String())

	assert.NoError(t, err)
	assert.Equal(t, "Great track!", d.Content)
	m.AssertExpectations(t)
}

func TestCreateDiscussion_WithParent(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	uid, tid, parentID := uuid.New(), uuid.New(), uuid.New()

	m.On("CreateDiscussion", mock.Anything, mock.MatchedBy(func(d *Discussion) bool {
		return d.ParentID != nil && *d.ParentID == parentID
	})).Return(nil)

	d, err := svc.CreateDiscussion(context.Background(), CreateDiscussionRequest{
		TargetType: "track",
		TargetID:   tid.String(),
		Content:    "Reply!",
		ParentID:   parentID.String(),
	}, uid.String())

	assert.NoError(t, err)
	assert.NotNil(t, d.ParentID)
	assert.Equal(t, parentID, *d.ParentID)
	m.AssertExpectations(t)
}

func TestGetDiscussions_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	tid := uuid.New()

	m.On("GetDiscussions", mock.Anything, "track", tid, 10, 0).Return([]Discussion{}, nil)

	items, err := svc.GetDiscussions(context.Background(), "track", tid.String(), 10, 0)
	assert.NoError(t, err)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestGetDiscussionReplies_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	pid := uuid.New()

	m.On("GetDiscussionReplies", mock.Anything, pid).Return([]Discussion{}, nil)

	items, err := svc.GetDiscussionReplies(context.Background(), pid.String())
	assert.NoError(t, err)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestCreateRating_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	uid, tid := uuid.New(), uuid.New()

	m.On("CreateRating", mock.Anything, mock.MatchedBy(func(r *TrackRating) bool {
		return r.Rating == 8 && r.UserID == uid && r.TrackID == tid && r.Review == nil
	})).Return(nil)

	rating, err := svc.CreateRating(context.Background(), CreateRatingRequest{
		TrackID: tid.String(),
		Rating:  8,
	}, uid.String())

	assert.NoError(t, err)
	assert.Equal(t, 8, rating.Rating)
	assert.Nil(t, rating.Review)
	m.AssertExpectations(t)
}

func TestCreateRating_WithReview(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	uid, tid := uuid.New(), uuid.New()
	review := "Amazing!"

	m.On("CreateRating", mock.Anything, mock.MatchedBy(func(r *TrackRating) bool {
		return r.Rating == 10 && r.Review != nil && *r.Review == "Amazing!"
	})).Return(nil)

	rating, err := svc.CreateRating(context.Background(), CreateRatingRequest{
		TrackID: tid.String(),
		Rating:  10,
		Review:  review,
	}, uid.String())

	assert.NoError(t, err)
	assert.Equal(t, 10, rating.Rating)
	assert.NotNil(t, rating.Review)
	assert.Equal(t, "Amazing!", *rating.Review)
	m.AssertExpectations(t)
}

func TestGetTrackRatings_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	tid := uuid.New()

	m.On("GetTrackRatings", mock.Anything, tid).Return([]TrackRating{}, nil)

	items, err := svc.GetTrackRatings(context.Background(), tid.String())
	assert.NoError(t, err)
	assert.Empty(t, items)
	m.AssertExpectations(t)
}

func TestGetTrackRatingAverage_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	tid := uuid.New()

	m.On("GetTrackRatingAverage", mock.Anything, tid).Return(7.5, 10, nil)

	avg, count, err := svc.GetTrackRatingAverage(context.Background(), tid.String())
	assert.NoError(t, err)
	assert.Equal(t, 7.5, avg)
	assert.Equal(t, 10, count)
	m.AssertExpectations(t)
}

func TestRecordActivity_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	uid := uuid.New()

	m.On("InsertActivity", mock.Anything, mock.MatchedBy(func(a Activity) bool {
		return a.Type == "track_like" && a.TargetType == "track"
	})).Return(nil)

	err := svc.RecordActivity(context.Background(), uid.String(), "track_like", "target-1", "track", nil)
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestUpdatePartyPosition_EmptyTrackID(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	pid := uuid.New()

	m.On("UpdatePartyPosition", mock.Anything, pid, (*uuid.UUID)(nil), int64(50000)).Return(nil)

	err := svc.UpdatePartyPosition(context.Background(), pid.String(), "", 50000)
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestUpdatePartyPosition_WithTrackID(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	pid, tid := uuid.New(), uuid.New()

	m.On("UpdatePartyPosition", mock.Anything, pid, &tid, int64(120000)).Return(nil)

	err := svc.UpdatePartyPosition(context.Background(), pid.String(), tid.String(), 120000)
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestUpdateRoomStatus_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)
	rid := uuid.New()

	m.On("UpdateRoomStatus", mock.Anything, rid, "ended").Return(nil)

	err := svc.UpdateRoomStatus(context.Background(), rid.String(), "ended")
	assert.NoError(t, err)
	m.AssertExpectations(t)
}
