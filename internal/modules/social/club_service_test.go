package social

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"music/internal/modules/playlist"
)

// MockPlaylistSvc satisfies PlaylistCollaborator interface
type MockPlaylistSvc struct {
	mock.Mock
}

func (m *MockPlaylistSvc) CreatePlaylist(ctx context.Context, req playlist.CreatePlaylistRequest, userID uuid.UUID) (playlist.Playlist, error) {
	args := m.Called(ctx, req, userID)
	return args.Get(0).(playlist.Playlist), args.Error(1)
}

func (m *MockPlaylistSvc) SetCollaborative(ctx context.Context, playlistID string, collab bool) error {
	return m.Called(ctx, playlistID, collab).Error(0)
}

func (m *MockPlaylistSvc) AddCollaborator(ctx context.Context, playlistID, userID string) error {
	return m.Called(ctx, playlistID, userID).Error(0)
}

func (m *MockPlaylistSvc) RemoveCollaborator(ctx context.Context, playlistID, userID string) error {
	return m.Called(ctx, playlistID, userID).Error(0)
}

func (m *MockPlaylistSvc) ListPlaylistTracks(ctx context.Context, playlistID string) ([]playlist.PlaylistTrackItem, error) {
	args := m.Called(ctx, playlistID)
	return args.Get(0).([]playlist.PlaylistTrackItem), args.Error(1)
}

func newClubSvc(mockRepo *mockRepo, mockPlaylist PlaylistCollaborator) *ClubService {
	svc := &Service{
		repo: mockRepo,
	}
	return &ClubService{
		repo:        mockRepo,
		playlistSvc: mockPlaylist,
		partySvc:    svc,
	}
}

func TestClubService_CreateClub_CreatesBackingPlaylist(t *testing.T) {
	m := new(mockRepo)
	mp := new(MockPlaylistSvc)

	ownerID := uuid.New()
	playlistID := uuid.New()

	mp.On("CreatePlaylist", mock.Anything, mock.Anything, ownerID).Return(playlist.Playlist{ID: playlistID}, nil)
	mp.On("SetCollaborative", mock.Anything, playlistID.String(), true).Return(nil)
	mp.On("AddCollaborator", mock.Anything, playlistID.String(), ownerID.String()).Return(nil)
	m.On("CreateClub", mock.Anything, mock.Anything).Return(nil)

	cs := newClubSvc(m, mp)
	club, err := cs.CreateClub(context.Background(), ownerID.String(), CreateClubRequest{
		Name: "Test Club",
	})

	assert.NoError(t, err)
	assert.NotNil(t, club)
	assert.Equal(t, "Test Club", club.Name)
	assert.Equal(t, playlistID.String(), club.PlaylistID)
	mp.AssertExpectations(t)
	m.AssertExpectations(t)
}

func TestClubService_JoinClub_Idempotent(t *testing.T) {
	m := new(mockRepo)
	mp := new(MockPlaylistSvc)

	clubID := uuid.New()
	userID := uuid.New()

	m.On("IsClubMember", mock.Anything, clubID, userID).Return(true, nil)

	cs := newClubSvc(m, mp)
	err := cs.JoinClub(context.Background(), clubID.String(), userID.String())

	assert.NoError(t, err)
	m.AssertCalled(t, "IsClubMember", mock.Anything, clubID, userID)
	mp.AssertNotCalled(t, "AddCollaborator")
}

func TestClubService_JoinClub_NotMember_Succeeds(t *testing.T) {
	m := new(mockRepo)
	mp := new(MockPlaylistSvc)

	clubID := uuid.New()
	userID := uuid.New()
	playlistID := uuid.New()

	m.On("IsClubMember", mock.Anything, clubID, userID).Return(false, nil)
	m.On("GetClub", mock.Anything, clubID).Return(&MusicClub{ID: clubID, PlaylistID: &playlistID}, nil)
	m.On("JoinClub", mock.Anything, clubID, userID).Return(nil)
	mp.On("AddCollaborator", mock.Anything, playlistID.String(), userID.String()).Return(nil)

	cs := newClubSvc(m, mp)
	err := cs.JoinClub(context.Background(), clubID.String(), userID.String())

	assert.NoError(t, err)
	mp.AssertCalled(t, "AddCollaborator", mock.Anything, playlistID.String(), userID.String())
}

func TestClubService_LeaveClub_OwnerCannotLeave(t *testing.T) {
	m := new(mockRepo)
	mp := new(MockPlaylistSvc)

	clubID := uuid.New()
	ownerID := uuid.New()

	m.On("IsClubAdmin", mock.Anything, clubID, ownerID).Return(true, nil)

	cs := newClubSvc(m, mp)
	err := cs.LeaveClub(context.Background(), clubID.String(), ownerID.String())

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrOwnerCannotLeave)
}

func TestClubService_LaunchParty_RequiresMembership(t *testing.T) {
	m := new(mockRepo)
	mp := new(MockPlaylistSvc)

	clubID := uuid.New()
	userID := uuid.New()

	m.On("IsClubMember", mock.Anything, clubID, userID).Return(false, nil)

	cs := newClubSvc(m, mp)
	party, err := cs.LaunchListeningParty(context.Background(), clubID.String(), userID.String(), LaunchPartyRequest{
		Title: "Test Party",
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrLaunchNotMember)
	assert.Nil(t, party)
}

func TestClubService_LaunchParty_Member_Succeeds(t *testing.T) {
	m := new(mockRepo)
	mp := new(MockPlaylistSvc)

	clubID := uuid.New()
	userID := uuid.New()
	playlistID := uuid.New()
	trackID := uuid.New()

	m.On("CreateRoom", mock.Anything, mock.Anything).Return(nil)
	m.On("IsClubMember", mock.Anything, clubID, userID).Return(true, nil)
	m.On("GetClub", mock.Anything, clubID).Return(&MusicClub{ID: clubID, PlaylistID: &playlistID}, nil)
	m.On("JoinRoom", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("JoinParty", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	m.On("CreateParty", mock.Anything, mock.Anything).Return(nil)
	m.On("JoinParty", mock.Anything, mock.Anything, userID).Return(nil)
	m.On("UpdatePartyStatus", mock.Anything, mock.Anything, "active", mock.Anything).Return(nil)
	mp.On("ListPlaylistTracks", mock.Anything, playlistID.String()).Return([]playlist.PlaylistTrackItem{
		{TrackID: trackID},
	}, nil)
	m.On("AddToRoomQueue", mock.Anything, mock.Anything).Return(nil)

	cs := newClubSvc(m, mp)
	party, err := cs.LaunchListeningParty(context.Background(), clubID.String(), userID.String(), LaunchPartyRequest{
		Title: "Test Party",
	})

	assert.NoError(t, err)
	assert.NotNil(t, party)
}

func TestClubService_GetClubDetail_Success(t *testing.T) {
	m := new(mockRepo)
	mp := new(MockPlaylistSvc)

	clubID := uuid.New()
	userID := uuid.New()
	playlistID := uuid.New()

	m.On("GetClub", mock.Anything, clubID).Return(&MusicClub{ID: clubID, Name: "Test Club", PlaylistID: &playlistID}, nil)
	m.On("GetClubMembers", mock.Anything, clubID).Return([]MusicClubMember{
		{UserID: userID, Role: "admin"},
	}, nil)
	m.On("GetClubPosts", mock.Anything, clubID, 50, 0).Return([]MusicClubPost{}, nil)
	m.On("IsClubMember", mock.Anything, clubID, userID).Return(true, nil)
	mp.On("ListPlaylistTracks", mock.Anything, playlistID.String()).Return([]playlist.PlaylistTrackItem{}, nil)

	cs := newClubSvc(m, mp)
	detail, err := cs.GetClubDetail(context.Background(), clubID.String(), userID.String())

	assert.NoError(t, err)
	assert.NotNil(t, detail)
	assert.Equal(t, "Test Club", detail.Club.Name)
	assert.True(t, detail.IsMember)
	assert.Equal(t, "admin", detail.MemberRole)
}
