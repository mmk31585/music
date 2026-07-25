//go:build playlist_integration

package playlist

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockPlaylistService struct {
	mock.Mock
}

func (m *mockPlaylistService) CreatePlaylist(ctx interface{}, req CreatePlaylistRequest, userID uuid.UUID) (*PlaylistResponse, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PlaylistResponse), args.Error(1)
}

func (m *mockPlaylistService) GetPlaylist(ctx interface{}, id uuid.UUID) (*PlaylistResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PlaylistResponse), args.Error(1)
}

func (m *mockPlaylistService) UpdatePlaylist(ctx interface{}, id uuid.UUID, userID uuid.UUID, req UpdatePlaylistRequest) (*PlaylistResponse, error) {
	args := m.Called(ctx, id, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PlaylistResponse), args.Error(1)
}

func (m *mockPlaylistService) DeletePlaylist(ctx interface{}, id uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *mockPlaylistService) ListPublicPlaylists(ctx interface{}, limit, offset int) ([]PlaylistResponse, int, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]PlaylistResponse), args.Int(1), args.Error(2)
}

func (m *mockPlaylistService) ListMyPlaylists(ctx interface{}, userID uuid.UUID, limit, offset int) ([]PlaylistResponse, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]PlaylistResponse), args.Int(1), args.Error(2)
}

func (m *mockPlaylistService) AddTrack(ctx interface{}, playlistID, trackID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, playlistID, trackID, userID)
	return args.Error(0)
}

func (m *mockPlaylistService) RemoveTrack(ctx interface{}, playlistID, trackID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, playlistID, trackID, userID)
	return args.Error(0)
}

func (m *mockPlaylistService) ReorderTrack(ctx interface{}, playlistID uuid.UUID, userID uuid.UUID, req ReorderTrackRequest) error {
	args := m.Called(ctx, playlistID, userID, req)
	return args.Error(0)
}

func (m *mockPlaylistService) ListCollaborators(ctx interface{}, playlistID uuid.UUID) ([]CollaboratorResponse, error) {
	args := m.Called(ctx, playlistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]CollaboratorResponse), args.Error(1)
}

func setupPlaylistTest() (*Handler, *mockPlaylistService, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockPlaylistService)
	h := NewHandler(mockSvc)

	r := gin.New()
	rg := r.Group("/api/v1")
	RegisterRoutes(rg, h, func(c *gin.Context) {
		c.Set("auth_user_id", uuid.New().String())
		c.Next()
	})
	return h, mockSvc, r
}

func TestCreatePlaylist_Success(t *testing.T) {
	_, mockSvc, r := setupPlaylistTest()

	userID := uuid.New()
	reqBody := CreatePlaylistRequest{
		Name: "My Playlist",
	}

	expected := &PlaylistResponse{
		ID:   uuid.New(),
		Name: "My Playlist",
	}

	mockSvc.On("CreatePlaylist", mock.Anything, reqBody, mock.AnythingOfType("uuid.UUID")).
		Return(expected, nil)

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/playlists", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp["success"].(bool))
	mockSvc.AssertNotCalled(t, "CreatePlaylist", mock.Anything, reqBody, userID)
	mockSvc.AssertExpectations(t)
}

func TestCreatePlaylist_InvalidBody(t *testing.T) {
	_, _, r := setupPlaylistTest()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/playlists", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPlaylist_Success(t *testing.T) {
	_, mockSvc, r := setupPlaylistTest()

	playlistID := uuid.New()
	expected := &PlaylistResponse{
		ID:   playlistID,
		Name: "My Playlist",
	}

	mockSvc.On("GetPlaylist", mock.Anything, playlistID).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/playlists/"+playlistID.String(), nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetPlaylist_NotFound(t *testing.T) {
	_, mockSvc, r := setupPlaylistTest()

	playlistID := uuid.New()
	mockSvc.On("GetPlaylist", mock.Anything, playlistID).Return(nil, ErrPlaylistNotFound)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/playlists/"+playlistID.String(), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeletePlaylist_Success(t *testing.T) {
	_, mockSvc, r := setupPlaylistTest()

	playlistID := uuid.New()
	mockSvc.On("DeletePlaylist", mock.Anything, playlistID, mock.AnythingOfType("uuid.UUID")).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/playlists/"+playlistID.String(), nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}
