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

func (m *mockRepo) CreatePlaylist(ctx context.Context, req CreatePlaylistRequest, userID int64) (Playlist, error) {
	args := m.Called(ctx, req, userID)
	return args.Get(0).(Playlist), args.Error(1)
}

func (m *mockRepo) UpdatePlaylist(ctx context.Context, playlistID, userID int64, req UpdatePlaylistRequest) (Playlist, error) {
	args := m.Called(ctx, playlistID, userID, req)
	return args.Get(0).(Playlist), args.Error(1)
}

func (m *mockRepo) DeletePlaylist(ctx context.Context, playlistID, userID int64) error {
	return m.Called(ctx, playlistID, userID).Error(0)
}

func (m *mockRepo) GetPlaylistByID(ctx context.Context, playlistID int64) (Playlist, error) {
	args := m.Called(ctx, playlistID)
	return args.Get(0).(Playlist), args.Error(1)
}

func (m *mockRepo) ListPlaylistTracks(ctx context.Context, playlistID int64) ([]PlaylistTrackItem, error) {
	args := m.Called(ctx, playlistID)
	return args.Get(0).([]PlaylistTrackItem), args.Error(1)
}

func (m *mockRepo) ListPublicPlaylists(ctx context.Context) ([]PlaylistListItemResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]PlaylistListItemResponse), args.Error(1)
}

func (m *mockRepo) ListUserPlaylists(ctx context.Context, userID int64) ([]PlaylistListItemResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]PlaylistListItemResponse), args.Error(1)
}

func (m *mockRepo) AddTrack(ctx context.Context, playlistID, trackID int64) error {
	return m.Called(ctx, playlistID, trackID).Error(0)
}

func (m *mockRepo) RemoveTrack(ctx context.Context, playlistID, trackID int64) error {
	return m.Called(ctx, playlistID, trackID).Error(0)
}

func (m *mockRepo) ReorderTrack(ctx context.Context, playlistID, trackID int64, newPosition int) error {
	return m.Called(ctx, playlistID, trackID, newPosition).Error(0)
}

func TestCreatePlaylist_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	now := time.Now()
	req := CreatePlaylistRequest{Name: "My Playlist", IsPublic: true}
	expected := Playlist{ID: 1, UserID: 10, Name: "My Playlist", IsPublic: true, CreatedAt: now, UpdatedAt: now}

	m.On("CreatePlaylist", mock.Anything, req, int64(10)).Return(expected, nil)

	p, err := svc.CreatePlaylist(context.Background(), req, 10)

	assert.NoError(t, err)
	assert.Equal(t, "My Playlist", p.Name)
	m.AssertExpectations(t)
}

func TestCreatePlaylist_EmptyName(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	_, err := svc.CreatePlaylist(context.Background(), CreatePlaylistRequest{Name: "   ", IsPublic: true}, 10)

	assert.ErrorIs(t, err, ErrInvalidPlaylistName)
	m.AssertNotCalled(t, "CreatePlaylist")
}

func TestUpdatePlaylist_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	req := UpdatePlaylistRequest{Name: "Updated Name", IsPublic: false}
	expected := Playlist{ID: 1, UserID: 10, Name: "Updated Name"}

	m.On("UpdatePlaylist", mock.Anything, int64(1), int64(10), req).Return(expected, nil)

	p, err := svc.UpdatePlaylist(context.Background(), 1, 10, req)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", p.Name)
	m.AssertExpectations(t)
}

func TestUpdatePlaylist_EmptyName(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	_, err := svc.UpdatePlaylist(context.Background(), 1, 10, UpdatePlaylistRequest{Name: ""})

	assert.ErrorIs(t, err, ErrInvalidPlaylistName)
	m.AssertNotCalled(t, "UpdatePlaylist")
}

func TestDeletePlaylist_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetPlaylistByID", mock.Anything, int64(1)).Return(Playlist{ID: 1, UserID: 10}, nil)
	m.On("DeletePlaylist", mock.Anything, int64(1), int64(10)).Return(nil)

	err := svc.DeletePlaylist(context.Background(), 1, 10)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestDeletePlaylist_NotOwner(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetPlaylistByID", mock.Anything, int64(1)).Return(Playlist{ID: 1, UserID: 10}, nil)

	err := svc.DeletePlaylist(context.Background(), 1, 20)

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "DeletePlaylist")
}

func TestDeletePlaylist_NotFound(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetPlaylistByID", mock.Anything, int64(1)).Return(Playlist{}, ErrPlaylistNotFound)

	err := svc.DeletePlaylist(context.Background(), 1, 10)

	assert.ErrorIs(t, err, ErrPlaylistNotFound)
}

func TestGetPlaylist_Public(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: 1, UserID: 10, Name: "Public", IsPublic: true}
	tracks := []PlaylistTrackItem{{TrackID: 1, Title: "Track 1"}}

	m.On("GetPlaylistByID", mock.Anything, int64(1)).Return(p, nil)
	m.On("ListPlaylistTracks", mock.Anything, int64(1)).Return(tracks, nil)

	result, resultTracks, err := svc.GetPlaylist(context.Background(), 1, nil)

	assert.NoError(t, err)
	assert.Equal(t, "Public", result.Name)
	assert.Len(t, resultTracks, 1)
	m.AssertExpectations(t)
}

func TestGetPlaylist_PrivateByOwner(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	userID := int64(10)
	p := Playlist{ID: 1, UserID: userID, Name: "Private", IsPublic: false}

	m.On("GetPlaylistByID", mock.Anything, int64(1)).Return(p, nil)
	m.On("ListPlaylistTracks", mock.Anything, int64(1)).Return([]PlaylistTrackItem{}, nil)

	_, _, err := svc.GetPlaylist(context.Background(), 1, &userID)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestGetPlaylist_PrivateUnauthorized(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: 1, UserID: 10, Name: "Private", IsPublic: false}

	m.On("GetPlaylistByID", mock.Anything, int64(1)).Return(p, nil)

	other := int64(20)
	_, _, err := svc.GetPlaylist(context.Background(), 1, &other)

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "ListPlaylistTracks")
}

func TestGetPlaylist_PrivateNoRequester(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: 1, UserID: 10, Name: "Private", IsPublic: false}

	m.On("GetPlaylistByID", mock.Anything, int64(1)).Return(p, nil)

	_, _, err := svc.GetPlaylist(context.Background(), 1, nil)

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "ListPlaylistTracks")
}

func TestListPublicPlaylists_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	expected := []PlaylistListItemResponse{
		{ID: 1, Name: "Public 1", TrackCount: 5},
		{ID: 2, Name: "Public 2", TrackCount: 3},
	}

	m.On("ListPublicPlaylists", mock.Anything).Return(expected, nil)

	items, err := svc.ListPublicPlaylists(context.Background())

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	m.AssertExpectations(t)
}

func TestListMyPlaylists_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	expected := []PlaylistListItemResponse{
		{ID: 1, Name: "My List", TrackCount: 10},
	}

	m.On("ListUserPlaylists", mock.Anything, int64(10)).Return(expected, nil)

	items, err := svc.ListMyPlaylists(context.Background(), 10)

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	m.AssertExpectations(t)
}

func TestAddTrack_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: 1, UserID: 10}

	m.On("GetPlaylistByID", mock.Anything, int64(1)).Return(p, nil)
	m.On("AddTrack", mock.Anything, int64(1), int64(100)).Return(nil)

	err := svc.AddTrack(context.Background(), 1, 10, 100)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestAddTrack_Forbidden(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetPlaylistByID", mock.Anything, int64(1)).Return(Playlist{ID: 1, UserID: 10}, nil)

	err := svc.AddTrack(context.Background(), 1, 20, 100)

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "AddTrack")
}

func TestAddTrack_RepoError(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: 1, UserID: 10}

	m.On("GetPlaylistByID", mock.Anything, int64(1)).Return(p, nil)
	m.On("AddTrack", mock.Anything, int64(1), int64(100)).Return(errors.New("db error"))

	err := svc.AddTrack(context.Background(), 1, 10, 100)

	assert.Error(t, err)
	m.AssertExpectations(t)
}

func TestRemoveTrack_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: 1, UserID: 10}

	m.On("GetPlaylistByID", mock.Anything, int64(1)).Return(p, nil)
	m.On("RemoveTrack", mock.Anything, int64(1), int64(100)).Return(nil)

	err := svc.RemoveTrack(context.Background(), 1, 10, 100)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestRemoveTrack_NotFound(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: 1, UserID: 10}

	m.On("GetPlaylistByID", mock.Anything, int64(1)).Return(p, nil)
	m.On("RemoveTrack", mock.Anything, int64(1), int64(100)).Return(ErrPlaylistTrackNotFound)

	err := svc.RemoveTrack(context.Background(), 1, 10, 100)

	assert.ErrorIs(t, err, ErrPlaylistTrackNotFound)
	m.AssertExpectations(t)
}

func TestReorderTrack_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: 1, UserID: 10}

	m.On("GetPlaylistByID", mock.Anything, int64(1)).Return(p, nil)
	m.On("ReorderTrack", mock.Anything, int64(1), int64(100), 2).Return(nil)

	err := svc.ReorderTrack(context.Background(), 1, 10, 100, 2)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestReorderTrack_InvalidPosition(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: 1, UserID: 10}

	m.On("GetPlaylistByID", mock.Anything, int64(1)).Return(p, nil)
	m.On("ReorderTrack", mock.Anything, int64(1), int64(100), 0).Return(ErrInvalidTrackPosition)

	err := svc.ReorderTrack(context.Background(), 1, 10, 100, 0)

	assert.ErrorIs(t, err, ErrInvalidTrackPosition)
	m.AssertExpectations(t)
}
