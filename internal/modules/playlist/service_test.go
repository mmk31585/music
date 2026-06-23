package playlist

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) CreatePlaylist(ctx context.Context, req CreatePlaylistRequest, userID uuid.UUID) (Playlist, error) {
	args := m.Called(ctx, req, userID)
	return args.Get(0).(Playlist), args.Error(1)
}

func (m *mockRepo) UpdatePlaylist(ctx context.Context, playlistID, userID uuid.UUID, req UpdatePlaylistRequest) (Playlist, error) {
	args := m.Called(ctx, playlistID, userID, req)
	return args.Get(0).(Playlist), args.Error(1)
}

func (m *mockRepo) DeletePlaylist(ctx context.Context, playlistID, userID uuid.UUID) error {
	return m.Called(ctx, playlistID, userID).Error(0)
}

func (m *mockRepo) GetPlaylistByID(ctx context.Context, playlistID uuid.UUID) (Playlist, error) {
	args := m.Called(ctx, playlistID)
	return args.Get(0).(Playlist), args.Error(1)
}

func (m *mockRepo) ListPlaylistTracks(ctx context.Context, playlistID uuid.UUID) ([]PlaylistTrackItem, error) {
	args := m.Called(ctx, playlistID)
	return args.Get(0).([]PlaylistTrackItem), args.Error(1)
}

func (m *mockRepo) ListPublicPlaylists(ctx context.Context) ([]PlaylistListItemResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]PlaylistListItemResponse), args.Error(1)
}

func (m *mockRepo) ListUserPlaylists(ctx context.Context, userID uuid.UUID) ([]PlaylistListItemResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]PlaylistListItemResponse), args.Error(1)
}

func (m *mockRepo) AddTrack(ctx context.Context, playlistID, trackID uuid.UUID) error {
	return m.Called(ctx, playlistID, trackID).Error(0)
}

func (m *mockRepo) RemoveTrack(ctx context.Context, playlistID, trackID uuid.UUID) error {
	return m.Called(ctx, playlistID, trackID).Error(0)
}

func (m *mockRepo) ReorderTrack(ctx context.Context, playlistID, trackID uuid.UUID, newPosition int) error {
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]CollaboratorResponse), args.Error(1)
}

var (
	uid1 = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	uid2 = uuid.MustParse("00000000-0000-0000-0000-000000000002")
	uid3 = uuid.MustParse("00000000-0000-0000-0000-000000000003")
	uid4 = uuid.MustParse("00000000-0000-0000-0000-000000000004")
)

func TestCreatePlaylist_ByName_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	now := time.Now()
	req := CreatePlaylistRequest{Name: "My Playlist", IsPublic: true}
	expected := Playlist{ID: uid1, UserID: uid2, Name: "My Playlist", IsPublic: true, CreatedAt: now, UpdatedAt: now}

	m.On("CreatePlaylist", mock.Anything, req, uid2).Return(expected, nil)

	p, err := svc.CreatePlaylist(context.Background(), req, uid2)

	assert.NoError(t, err)
	assert.Equal(t, "My Playlist", p.Name)
	m.AssertExpectations(t)
}

func TestCreatePlaylist_EmptyName(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	_, err := svc.CreatePlaylist(context.Background(), CreatePlaylistRequest{Name: "   ", IsPublic: true}, uid1)

	assert.ErrorIs(t, err, ErrInvalidPlaylistName)
	m.AssertNotCalled(t, "CreatePlaylist")
}

func TestUpdatePlaylist_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	req := UpdatePlaylistRequest{Name: "Updated Name", IsPublic: false}
	expected := Playlist{ID: uid1, UserID: uid2, Name: "Updated Name"}

	m.On("UpdatePlaylist", mock.Anything, uid1, uid2, req).Return(expected, nil)

	p, err := svc.UpdatePlaylist(context.Background(), uid1, uid2, req)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", p.Name)
	m.AssertExpectations(t)
}

func TestUpdatePlaylist_EmptyName(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	_, err := svc.UpdatePlaylist(context.Background(), uid1, uid2, UpdatePlaylistRequest{Name: ""})

	assert.ErrorIs(t, err, ErrInvalidPlaylistName)
	m.AssertNotCalled(t, "UpdatePlaylist")
}

func TestDeletePlaylist_ByOwner_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetPlaylistByID", mock.Anything, uid1).Return(Playlist{ID: uid1, UserID: uid2}, nil)
	m.On("DeletePlaylist", mock.Anything, uid1, uid2).Return(nil)

	err := svc.DeletePlaylist(context.Background(), uid1, uid2)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestDeletePlaylist_NotOwner(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetPlaylistByID", mock.Anything, uid1).Return(Playlist{ID: uid1, UserID: uid2}, nil)

	err := svc.DeletePlaylist(context.Background(), uid1, uid3)

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "DeletePlaylist")
}

func TestDeletePlaylist_NotFound(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetPlaylistByID", mock.Anything, uid1).Return(Playlist{}, ErrPlaylistNotFound)

	err := svc.DeletePlaylist(context.Background(), uid1, uid2)

	assert.ErrorIs(t, err, ErrPlaylistNotFound)
}

func TestGetPlaylist_Public(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: uid1, UserID: uid2, Name: "Public", IsPublic: true}
	tracks := []PlaylistTrackItem{{PlaylistTrackID: uid3, TrackID: uid4, Title: "Track 1"}}

	m.On("GetPlaylistByID", mock.Anything, uid1).Return(p, nil)
	m.On("ListPlaylistTracks", mock.Anything, uid1).Return(tracks, nil)

	result, resultTracks, err := svc.GetPlaylist(context.Background(), uid1, nil)

	assert.NoError(t, err)
	assert.Equal(t, "Public", result.Name)
	assert.Len(t, resultTracks, 1)
	m.AssertExpectations(t)
}

func TestGetPlaylist_PrivateByOwner(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: uid1, UserID: uid2, Name: "Private", IsPublic: false}

	m.On("GetPlaylistByID", mock.Anything, uid1).Return(p, nil)
	m.On("ListPlaylistTracks", mock.Anything, uid1).Return([]PlaylistTrackItem{}, nil)

	_, _, err := svc.GetPlaylist(context.Background(), uid1, &uid2)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestGetPlaylist_PrivateUnauthorized(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: uid1, UserID: uid2, Name: "Private", IsPublic: false}

	m.On("GetPlaylistByID", mock.Anything, uid1).Return(p, nil)

	_, _, err := svc.GetPlaylist(context.Background(), uid1, &uid3)

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "ListPlaylistTracks")
}

func TestGetPlaylist_PrivateNoRequester(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: uid1, UserID: uid2, Name: "Private", IsPublic: false}

	m.On("GetPlaylistByID", mock.Anything, uid1).Return(p, nil)

	_, _, err := svc.GetPlaylist(context.Background(), uid1, nil)

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "ListPlaylistTracks")
}

func TestListPublicPlaylists_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	expected := []PlaylistListItemResponse{
		{ID: uid1, Name: "Public 1", TrackCount: 5},
		{ID: uid2, Name: "Public 2", TrackCount: 3},
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
		{ID: uid1, Name: "My List", TrackCount: 10},
	}

	m.On("ListUserPlaylists", mock.Anything, uid2).Return(expected, nil)

	items, err := svc.ListMyPlaylists(context.Background(), uid2)

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	m.AssertExpectations(t)
}

func TestAddTrack_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: uid1, UserID: uid2}

	m.On("GetPlaylistByID", mock.Anything, uid1).Return(p, nil)
	m.On("AddTrack", mock.Anything, uid1, uid3).Return(nil)

	err := svc.AddTrack(context.Background(), uid1, uid2, uid3)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestAddTrack_Forbidden(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	m.On("GetPlaylistByID", mock.Anything, uid1).Return(Playlist{ID: uid1, UserID: uid2}, nil)

	err := svc.AddTrack(context.Background(), uid1, uid3, uid4)

	assert.ErrorIs(t, err, ErrForbiddenPlaylistAccess)
	m.AssertNotCalled(t, "AddTrack")
}

func TestAddTrack_RepoError(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: uid1, UserID: uid2}

	m.On("GetPlaylistByID", mock.Anything, uid1).Return(p, nil)
	m.On("AddTrack", mock.Anything, uid1, uid3).Return(errors.New("db error"))

	err := svc.AddTrack(context.Background(), uid1, uid2, uid3)

	assert.Error(t, err)
	m.AssertExpectations(t)
}

func TestRemoveTrack_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: uid1, UserID: uid2}

	m.On("GetPlaylistByID", mock.Anything, uid1).Return(p, nil)
	m.On("RemoveTrack", mock.Anything, uid1, uid3).Return(nil)

	err := svc.RemoveTrack(context.Background(), uid1, uid2, uid3)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestRemoveTrack_NotFound(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: uid1, UserID: uid2}

	m.On("GetPlaylistByID", mock.Anything, uid1).Return(p, nil)
	m.On("RemoveTrack", mock.Anything, uid1, uid3).Return(ErrPlaylistTrackNotFound)

	err := svc.RemoveTrack(context.Background(), uid1, uid2, uid3)

	assert.ErrorIs(t, err, ErrPlaylistTrackNotFound)
	m.AssertExpectations(t)
}

func TestReorderTrack_Success(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: uid1, UserID: uid2}

	m.On("GetPlaylistByID", mock.Anything, uid1).Return(p, nil)
	m.On("ReorderTrack", mock.Anything, uid1, uid3, 2).Return(nil)

	err := svc.ReorderTrack(context.Background(), uid1, uid2, uid3, 2)

	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestReorderTrack_InvalidPosition(t *testing.T) {
	m := new(mockRepo)
	svc := NewService(m)

	p := Playlist{ID: uid1, UserID: uid2}

	m.On("GetPlaylistByID", mock.Anything, uid1).Return(p, nil)
	m.On("ReorderTrack", mock.Anything, uid1, uid3, 0).Return(ErrInvalidTrackPosition)

	err := svc.ReorderTrack(context.Background(), uid1, uid2, uid3, 0)

	assert.ErrorIs(t, err, ErrInvalidTrackPosition)
	m.AssertExpectations(t)
}
