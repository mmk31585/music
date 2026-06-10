package playlist

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) CreatePlaylist(ctx context.Context, req CreatePlaylistRequest, userID string) (Playlist, error) {
	args := m.Called(ctx, req, userID)
	return args.Get(0).(Playlist), args.Error(1)
}

func (m *mockRepo) UpdatePlaylist(ctx context.Context, playlistID string, req UpdatePlaylistRequest) (Playlist, error) {
	args := m.Called(ctx, playlistID, req)
	return args.Get(0).(Playlist), args.Error(1)
}

func (m *mockRepo) DeletePlaylist(ctx context.Context, playlistID, userID string) error {
	return m.Called(ctx, playlistID, userID).Error(0)
}

func (m *mockRepo) GetPlaylistByID(ctx context.Context, playlistID string) (Playlist, error) {
	args := m.Called(ctx, playlistID)
	return args.Get(0).(Playlist), args.Error(1)
}

func (m *mockRepo) ListPlaylistTracks(ctx context.Context, playlistID string) ([]PlaylistTrackItem, error) {
	args := m.Called(ctx, playlistID)
	return args.Get(0).([]PlaylistTrackItem), args.Error(1)
}

func (m *mockRepo) ListPublicPlaylists(ctx context.Context) ([]PlaylistListItemResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]PlaylistListItemResponse), args.Error(1)
}

func (m *mockRepo) ListUserPlaylists(ctx context.Context, userID string) ([]PlaylistListItemResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]PlaylistListItemResponse), args.Error(1)
}

func (m *mockRepo) AddTrack(ctx context.Context, playlistID, trackID string) error {
	return m.Called(ctx, playlistID, trackID).Error(0)
}

func (m *mockRepo) RemoveTrack(ctx context.Context, playlistID, trackID string) error {
	return m.Called(ctx, playlistID, trackID).Error(0)
}

func (m *mockRepo) ReorderTrack(ctx context.Context, playlistID, trackID string, newPosition int) error {
	return m.Called(ctx, playlistID, trackID, newPosition).Error(0)
}

func (m *mockRepo) IsCollaborator(ctx context.Context, playlistID, userID string) (bool, error) {
	args := m.Called(ctx, playlistID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepo) IsCollaborativePlaylist(ctx context.Context, playlistID string) (bool, error) {
	args := m.Called(ctx, playlistID)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepo) SetCollaborative(ctx context.Context, playlistID string, collab bool) error {
	return m.Called(ctx, playlistID, collab).Error(0)
}

func (m *mockRepo) AddCollaborator(ctx context.Context, playlistID, userID string) error {
	return m.Called(ctx, playlistID, userID).Error(0)
}

func (m *mockRepo) RemoveCollaborator(ctx context.Context, playlistID, userID string) error {
	return m.Called(ctx, playlistID, userID).Error(0)
}

func (m *mockRepo) ListCollaborators(ctx context.Context, playlistID string) ([]CollaboratorResponse, error) {
	args := m.Called(ctx, playlistID)
	return args.Get(0).([]CollaboratorResponse), args.Error(1)
}

func ptr(s string) *string { return &s }

func newSvc(mock *mockRepo) *Service {
	return NewService(mock, NewPlaylistBroadcaster(nil))
}

func TestCreatePlaylist_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	now := time.Now()
	req := CreatePlaylistRequest{Name: "My Playlist", IsPublic: true}
	expected := Playlist{ID: "pl-1", UserID: "u1", Name: "My Playlist", IsPublic: true, CreatedAt: now, UpdatedAt: now}

	m.On("CreatePlaylist", mock.Anything, req, "u1").Return(expected, nil)

	p, err := svc.CreatePlaylist(context.Background(), req, "u1")

	assert.NoError(t, err)
	assert.Equal(t, "My Playlist", p.Name)
	m.AssertExpectations(t)
}

func TestCreatePlaylist_EmptyName(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	_, err := svc.CreatePlaylist(context.Background(), CreatePlaylistRequest{Name: "   ", IsPublic: true}, "u1")

	assert.ErrorIs(t, err, ErrInvalidPlaylistName)
	m.AssertNotCalled(t, "CreatePlaylist")
}

func TestUpdatePlaylist_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	req := UpdatePlaylistRequest{Name: "Updated Name", IsPublic: false}
	playlist := Playlist{ID: "pl-1", UserID: "u1", Name: "Original"}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(playlist, nil)
	m.On("UpdatePlaylist", mock.Anything, "pl-1", req).Return(Playlist{ID: "pl-1", UserID: "u1", Name: "Updated Name"}, nil)

	p, err := svc.UpdatePlaylist(context.Background(), "pl-1", "u1", req)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", p.Name)
	m.AssertExpectations(t)
}

func TestUpdatePlaylist_EmptyName(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	_, err := svc.UpdatePlaylist(context.Background(), "pl-1", "u1", UpdatePlaylistRequest{Name: ""})

	assert.ErrorIs(t, err, ErrInvalidPlaylistName)
	m.AssertNotCalled(t, "GetPlaylistByID")
	m.AssertNotCalled(t, "UpdatePlaylist")
}

func TestUpdatePlaylist_NotOwner(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	playlist := Playlist{ID: "pl-1", UserID: "u1"}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(playlist, nil)
	m.On("IsCollaborativePlaylist", mock.Anything, "pl-1").Return(false, nil)

	_, err := svc.UpdatePlaylist(context.Background(), "pl-1", "u2", UpdatePlaylistRequest{Name: "Test"})

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertExpectations(t)
}

func TestUpdatePlaylist_AsCollaborator(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	playlist := Playlist{ID: "pl-1", UserID: "u1"}
	req := UpdatePlaylistRequest{Name: "Collab Update", IsPublic: true}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(playlist, nil)
	m.On("IsCollaborativePlaylist", mock.Anything, "pl-1").Return(true, nil)
	m.On("IsCollaborator", mock.Anything, "pl-1", "u2").Return(true, nil)
	m.On("UpdatePlaylist", mock.Anything, "pl-1", req).Return(Playlist{ID: "pl-1", Name: "Collab Update"}, nil)

	p, err := svc.UpdatePlaylist(context.Background(), "pl-1", "u2", req)

	assert.NoError(t, err)
	assert.Equal(t, "Collab Update", p.Name)
	m.AssertExpectations(t)
}

func TestUpdatePlaylist_NonCollabNotOwner(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	playlist := Playlist{ID: "pl-1", UserID: "u1"}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(playlist, nil)
	m.On("IsCollaborativePlaylist", mock.Anything, "pl-1").Return(false, nil)

	_, err := svc.UpdatePlaylist(context.Background(), "pl-1", "u2", UpdatePlaylistRequest{Name: "Test"})

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertExpectations(t)
}

func TestDeletePlaylist_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(Playlist{ID: "pl-1", UserID: "u1"}, nil)
	m.On("DeletePlaylist", mock.Anything, "pl-1", "u1").Return(nil)

	err := svc.DeletePlaylist(context.Background(), "pl-1", "u1")

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestDeletePlaylist_NotOwner(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(Playlist{ID: "pl-1", UserID: "u1"}, nil)

	err := svc.DeletePlaylist(context.Background(), "pl-1", "u2")

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "DeletePlaylist")
}

func TestDeletePlaylist_NotFound(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(Playlist{}, ErrPlaylistNotFound)

	err := svc.DeletePlaylist(context.Background(), "pl-1", "u1")

	assert.ErrorIs(t, err, ErrPlaylistNotFound)
}

func TestGetPlaylist_Public(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	p := Playlist{ID: "pl-1", UserID: "u1", Name: "Public", IsPublic: true, IsCollaborative: true}
	tracks := []PlaylistTrackItem{{TrackID: "t1", Title: "Track 1"}}
	collabs := []CollaboratorResponse{{UserID: "u2", IsCreator: false}}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)
	m.On("ListPlaylistTracks", mock.Anything, "pl-1").Return(tracks, nil)
	m.On("ListCollaborators", mock.Anything, "pl-1").Return(collabs, nil)

	result, resultTracks, resultCollabs, err := svc.GetPlaylist(context.Background(), "pl-1", nil)

	assert.NoError(t, err)
	assert.Equal(t, "Public", result.Name)
	assert.Len(t, resultTracks, 1)
	assert.Contains(t, resultCollabs, "u1")
	assert.Contains(t, resultCollabs, "u2")
	m.AssertExpectations(t)
}

func TestGetPlaylist_PrivateByOwner(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	userID := "u1"
	p := Playlist{ID: "pl-1", UserID: userID, Name: "Private", IsPublic: false}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)
	m.On("ListPlaylistTracks", mock.Anything, "pl-1").Return([]PlaylistTrackItem{}, nil)
	m.On("ListCollaborators", mock.Anything, "pl-1").Return([]CollaboratorResponse{}, nil)

	_, _, _, err := svc.GetPlaylist(context.Background(), "pl-1", &userID)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestGetPlaylist_PrivateUnauthorized(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	p := Playlist{ID: "pl-1", UserID: "u1", Name: "Private", IsPublic: false}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)

	other := "u2"
	_, _, _, err := svc.GetPlaylist(context.Background(), "pl-1", &other)

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "ListPlaylistTracks")
}

func TestGetPlaylist_PrivateNoRequester(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	p := Playlist{ID: "pl-1", UserID: "u1", Name: "Private", IsPublic: false}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)

	_, _, _, err := svc.GetPlaylist(context.Background(), "pl-1", nil)

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "ListPlaylistTracks")
}

func TestListPublicPlaylists_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	expected := []PlaylistListItemResponse{
		{ID: "pl-1", Name: "Public 1", TrackCount: 5},
		{ID: "pl-2", Name: "Public 2", TrackCount: 3},
	}

	m.On("ListPublicPlaylists", mock.Anything).Return(expected, nil)

	items, err := svc.ListPublicPlaylists(context.Background())

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	m.AssertExpectations(t)
}

func TestListMyPlaylists_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	expected := []PlaylistListItemResponse{
		{ID: "pl-1", Name: "My List", TrackCount: 10},
	}

	m.On("ListUserPlaylists", mock.Anything, "u1").Return(expected, nil)

	items, err := svc.ListMyPlaylists(context.Background(), "u1")

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	m.AssertExpectations(t)
}

func TestAddTrack_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	p := Playlist{ID: "pl-1", UserID: "u1"}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)
	m.On("AddTrack", mock.Anything, "pl-1", "t1").Return(nil)

	err := svc.AddTrack(context.Background(), "pl-1", "u1", "t1")

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestAddTrack_Forbidden(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(Playlist{ID: "pl-1", UserID: "u1"}, nil)
	m.On("IsCollaborativePlaylist", mock.Anything, "pl-1").Return(false, nil)

	err := svc.AddTrack(context.Background(), "pl-1", "u2", "t1")

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "AddTrack")
}

func TestAddTrack_RepoError(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	p := Playlist{ID: "pl-1", UserID: "u1"}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)
	m.On("AddTrack", mock.Anything, "pl-1", "t1").Return(errors.New("db error"))

	err := svc.AddTrack(context.Background(), "pl-1", "u1", "t1")

	assert.Error(t, err)
	m.AssertExpectations(t)
}

func TestRemoveTrack_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	p := Playlist{ID: "pl-1", UserID: "u1"}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)
	m.On("RemoveTrack", mock.Anything, "pl-1", "t1").Return(nil)

	err := svc.RemoveTrack(context.Background(), "pl-1", "u1", "t1")

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestRemoveTrack_NotFound(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	p := Playlist{ID: "pl-1", UserID: "u1"}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)
	m.On("RemoveTrack", mock.Anything, "pl-1", "t1").Return(ErrPlaylistTrackNotFound)

	err := svc.RemoveTrack(context.Background(), "pl-1", "u1", "t1")

	assert.ErrorIs(t, err, ErrPlaylistTrackNotFound)
	m.AssertExpectations(t)
}

func TestReorderTrack_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	p := Playlist{ID: "pl-1", UserID: "u1"}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)
	m.On("ReorderTrack", mock.Anything, "pl-1", "t1", 2).Return(nil)

	err := svc.ReorderTrack(context.Background(), "pl-1", "u1", "t1", 2)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestReorderTrack_InvalidPosition(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	p := Playlist{ID: "pl-1", UserID: "u1"}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)
	m.On("ReorderTrack", mock.Anything, "pl-1", "t1", 0).Return(ErrInvalidTrackPosition)

	err := svc.ReorderTrack(context.Background(), "pl-1", "u1", "t1", 0)

	assert.ErrorIs(t, err, ErrInvalidTrackPosition)
	m.AssertExpectations(t)
}

func TestSetCollaborative_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	p := Playlist{ID: "pl-1", UserID: "u1"}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)
	m.On("SetCollaborative", mock.Anything, "pl-1", true).Return(nil)

	err := svc.SetCollaborative(context.Background(), "pl-1", "u1", true)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestSetCollaborative_NotOwner(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(Playlist{ID: "pl-1", UserID: "u1"}, nil)

	err := svc.SetCollaborative(context.Background(), "pl-1", "u2", true)

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "SetCollaborative")
}

func TestAddCollaborator_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	p := Playlist{ID: "pl-1", UserID: "u1"}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)
	m.On("AddCollaborator", mock.Anything, "pl-1", "u2").Return(nil)

	err := svc.AddCollaborator(context.Background(), "pl-1", "u1", "u2")

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestAddCollaborator_NotOwner(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(Playlist{ID: "pl-1", UserID: "u1"}, nil)

	err := svc.AddCollaborator(context.Background(), "pl-1", "u2", "u3")

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "AddCollaborator")
}

func TestAddCollaborator_AlreadyAdded(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	p := Playlist{ID: "pl-1", UserID: "u1"}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)
	m.On("AddCollaborator", mock.Anything, "pl-1", "u2").Return(ErrAlreadyCollaborator)

	err := svc.AddCollaborator(context.Background(), "pl-1", "u1", "u2")

	assert.ErrorIs(t, err, ErrAlreadyCollaborator)
	m.AssertExpectations(t)
}

func TestRemoveCollaborator_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	p := Playlist{ID: "pl-1", UserID: "u1"}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)
	m.On("RemoveCollaborator", mock.Anything, "pl-1", "u2").Return(nil)

	err := svc.RemoveCollaborator(context.Background(), "pl-1", "u1", "u2")

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestRemoveCollaborator_NotOwner(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(Playlist{ID: "pl-1", UserID: "u1"}, nil)

	err := svc.RemoveCollaborator(context.Background(), "pl-1", "u2", "u3")

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "RemoveCollaborator")
}

func TestRemoveCollaborator_NotACollaborator(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	p := Playlist{ID: "pl-1", UserID: "u1"}

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(p, nil)
	m.On("RemoveCollaborator", mock.Anything, "pl-1", "u2").Return(ErrNotCollaborator)

	err := svc.RemoveCollaborator(context.Background(), "pl-1", "u1", "u2")

	assert.ErrorIs(t, err, ErrNotCollaborator)
	m.AssertExpectations(t)
}

func TestListCollaborators_Success(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	expected := []CollaboratorResponse{
		{UserID: "u2", IsCreator: false},
	}

	m.On("ListCollaborators", mock.Anything, "pl-1").Return(expected, nil)

	items, err := svc.ListCollaborators(context.Background(), "pl-1")

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	m.AssertExpectations(t)
}

func TestCanModify_OwnerAlways(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(Playlist{ID: "pl-1", UserID: "u1"}, nil)

	ok, err := svc.canModify(context.Background(), "pl-1", "u1")

	assert.NoError(t, err)
	assert.True(t, ok)
	m.AssertExpectations(t)
	m.AssertNotCalled(t, "IsCollaborativePlaylist")
}

func TestCanModify_Collaborator(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(Playlist{ID: "pl-1", UserID: "u1"}, nil)
	m.On("IsCollaborativePlaylist", mock.Anything, "pl-1").Return(true, nil)
	m.On("IsCollaborator", mock.Anything, "pl-1", "u2").Return(true, nil)

	ok, err := svc.canModify(context.Background(), "pl-1", "u2")

	assert.NoError(t, err)
	assert.True(t, ok)
	m.AssertExpectations(t)
}

func TestCanModify_NonCollaborator(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(Playlist{ID: "pl-1", UserID: "u1"}, nil)
	m.On("IsCollaborativePlaylist", mock.Anything, "pl-1").Return(true, nil)
	m.On("IsCollaborator", mock.Anything, "pl-1", "u2").Return(false, nil)

	ok, err := svc.canModify(context.Background(), "pl-1", "u2")

	assert.NoError(t, err)
	assert.False(t, ok)
	m.AssertExpectations(t)
}

func TestCanModify_NonCollabNotOwner(t *testing.T) {
	m := new(mockRepo)
	svc := newSvc(m)

	m.On("GetPlaylistByID", mock.Anything, "pl-1").Return(Playlist{ID: "pl-1", UserID: "u1"}, nil)
	m.On("IsCollaborativePlaylist", mock.Anything, "pl-1").Return(false, nil)

	ok, err := svc.canModify(context.Background(), "pl-1", "u2")

	assert.NoError(t, err)
	assert.False(t, ok)
	m.AssertExpectations(t)
	m.AssertNotCalled(t, "IsCollaborator")
}
