package upload

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"music/internal/common/response"
	apperrors "music/internal/common/errors"
	platformstorage "music/internal/platform/storage"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) CreateDraft(ctx context.Context, userID string, req CreateUploadRequest) (Draft, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).(Draft), args.Error(1)
}

func (m *mockService) GetDraft(ctx context.Context, draftID string) (Draft, error) {
	args := m.Called(ctx, draftID)
	return args.Get(0).(Draft), args.Error(1)
}

func (m *mockService) ListPendingDrafts(ctx context.Context, source *string, limit, offset int) ([]Draft, int, error) {
	args := m.Called(ctx, source, limit, offset)
	return args.Get(0).([]Draft), args.Int(1), args.Error(2)
}

func (m *mockService) ListUserDrafts(ctx context.Context, userID string, limit, offset int) ([]Draft, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]Draft), args.Int(1), args.Error(2)
}

func (m *mockService) ReviewDraft(ctx context.Context, draftID string, reviewerID string, req ReviewDraftRequest) error {
	args := m.Called(ctx, draftID, reviewerID, req)
	return args.Error(0)
}

func (m *mockService) AttachFile(ctx context.Context, draftID, userID, storageKey, fileURL string, fileSize int64) error {
	args := m.Called(ctx, draftID, userID, storageKey, fileURL, fileSize)
	return args.Error(0)
}

func (m *mockService) AddCoUploader(ctx context.Context, trackID int64, req CoUploaderRequest) (CoUploader, error) {
	args := m.Called(ctx, trackID, req)
	return args.Get(0).(CoUploader), args.Error(1)
}

func (m *mockService) RemoveCoUploader(ctx context.Context, trackID int64, userID string) error {
	args := m.Called(ctx, trackID, userID)
	return args.Error(0)
}

func (m *mockService) GetCoUploaders(ctx context.Context, trackID int64) ([]CoUploader, error) {
	args := m.Called(ctx, trackID)
	return args.Get(0).([]CoUploader), args.Error(1)
}

func (m *mockService) GetUploadSlots(ctx context.Context, userID string) (UploadSlotStatus, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(UploadSlotStatus), args.Error(1)
}

type mockStorage struct {
	mock.Mock
	uploadErr   error
	getURLErr   error
	getURLValue string
}

func (m *mockStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	args := m.Called(ctx, key, reader, size, contentType)
	return args.Error(0)
}

func (m *mockStorage) Get(ctx context.Context, key string) ([]byte, error) {
	args := m.Called(ctx, key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockStorage) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *mockStorage) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *mockStorage) GetURL(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *mockStorage) Copy(ctx context.Context, srcKey, dstKey string) error {
	args := m.Called(ctx, srcKey, dstKey)
	return args.Error(0)
}

var _ platformstorage.Storage = (*mockStorage)(nil)

func setupHandlerTest(svc ServiceInterface, storage platformstorage.Storage) (*Handler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	h := NewHandler(svc, storage)

	r := gin.New()
	r.Use(setTestUser)
	r.POST("/uploads/drafts", response.Wrap(h.CreateDraft))
	r.POST("/uploads/drafts/:draftId/file", response.Wrap(h.UploadDraftFile))
	r.GET("/uploads/drafts/mine", response.Wrap(h.ListMyDrafts))
	r.GET("/uploads/drafts/:draftId", response.Wrap(h.GetDraft))
	r.GET("/uploads/slots", response.Wrap(h.GetUploadSlots))
	r.GET("/tracks/:trackId/co-uploaders", response.Wrap(h.GetCoUploaders))
	r.POST("/uploads/review/:draftId", response.Wrap(h.ReviewDraft))
	r.GET("/uploads/review", response.Wrap(h.ListPendingDrafts))
	r.POST("/uploads/co-uploaders/:trackId", response.Wrap(h.AddCoUploader))
	r.DELETE("/uploads/co-uploaders/:trackId/:userId", response.Wrap(h.RemoveCoUploader))

	return h, r
}

// setTestUser injects an authenticated user into Gin context so handlers
// (which read auth.UserIDFromContext) work without JWT middleware.
func setTestUser(c *gin.Context) {
	c.Set("auth_user_id", testUserID)
	c.Next()
}

const testUserID = "test-user-id"

func TestCreateDraft_Success(t *testing.T) {
	mockSvc := new(mockService)
	_, r := setupHandlerTest(mockSvc, nil)

	expected := Draft{
		ID:                "draft_1",
		UploadedBy:        testUserID,
		OriginalFilename:  "song.mp3",
		Status:            "pending",
		UploadSource:      string(SourceListener),
		NeedsReview:       true,
	}
	mockSvc.On("CreateDraft", mock.Anything, testUserID, mock.Anything).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/uploads/drafts", strings.NewReader(`{"filename":"song.mp3"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCreateDraft_Unauthenticated(t *testing.T) {
	mockSvc := new(mockService)
	// Engine without setTestUser -> no user in context
	r := gin.New()
	h := NewHandler(mockSvc, nil)
	r.POST("/uploads/drafts", response.Wrap(h.CreateDraft))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/uploads/drafts", strings.NewReader(`{"filename":"song.mp3"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockSvc.AssertNotCalled(t, "CreateDraft", mock.Anything, testUserID, mock.Anything)
}

func TestCreateDraft_InvalidBody(t *testing.T) {
	mockSvc := new(mockService)
	_, r := setupHandlerTest(mockSvc, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/uploads/drafts", strings.NewReader(`not-json`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "CreateDraft", mock.Anything, testUserID, mock.Anything)
}

func TestCreateDraft_ServiceError_Propagates(t *testing.T) {
	mockSvc := new(mockService)
	_, r := setupHandlerTest(mockSvc, nil)

	mockSvc.On("CreateDraft", mock.Anything, testUserID, mock.Anything).
		Return(Draft{}, apperrors.New(http.StatusTooManyRequests, apperrors.CodeTooManyRequests, "slot limit reached", nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/uploads/drafts", strings.NewReader(`{"filename":"song.mp3"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestGetDraft_Success(t *testing.T) {
	mockSvc := new(mockService)
	_, r := setupHandlerTest(mockSvc, nil)

	draftID := "draft_1"
	expected := Draft{ID: draftID, UploadedBy: testUserID, Status: "pending"}
	mockSvc.On("GetDraft", mock.Anything, draftID).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/uploads/drafts/"+draftID, nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetDraft_NotFound(t *testing.T) {
	mockSvc := new(mockService)
	_, r := setupHandlerTest(mockSvc, nil)

	draftID := "draft_missing"
	mockSvc.On("GetDraft", mock.Anything, draftID).Return(Draft{}, apperrors.New(http.StatusNotFound, apperrors.CodeNotFound, "draft not found", nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/uploads/drafts/"+draftID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestListMyDrafts_Success(t *testing.T) {
	mockSvc := new(mockService)
	_, r := setupHandlerTest(mockSvc, nil)

	expected := []Draft{{ID: "draft_1", UploadedBy: testUserID}}
	mockSvc.On("ListUserDrafts", mock.Anything, testUserID, 20, 0).Return(expected, 1, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/uploads/drafts/mine", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetUploadSlots_Success(t *testing.T) {
	mockSvc := new(mockService)
	_, r := setupHandlerTest(mockSvc, nil)

	expected := UploadSlotStatus{Available: true, Used: 1, Max: 5}
	mockSvc.On("GetUploadSlots", mock.Anything, testUserID).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/uploads/slots", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestReviewDraft_Success(t *testing.T) {
	mockSvc := new(mockService)
	_, r := setupHandlerTest(mockSvc, nil)

	draftID := "draft_1"
	mockSvc.On("ReviewDraft", mock.Anything, draftID, testUserID, mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	body := `{"action":"accept","notes":"looks good"}`
	req, _ := http.NewRequest("POST", "/uploads/review/"+draftID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetCoUploaders_Success(t *testing.T) {
	mockSvc := new(mockService)
	_, r := setupHandlerTest(mockSvc, nil)

	trackID := int64(42)
	expected := []CoUploader{{TrackID: trackID, UserID: testUserID}}
	mockSvc.On("GetCoUploaders", mock.Anything, trackID).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tracks/42/co-uploaders", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestAddCoUploader_Success(t *testing.T) {
	mockSvc := new(mockService)
	_, r := setupHandlerTest(mockSvc, nil)

	trackID := int64(42)
	expected := CoUploader{TrackID: trackID, UserID: testUserID, Role: "uploader"}
	mockSvc.On("AddCoUploader", mock.Anything, trackID, mock.Anything).Return(expected, nil)

	w := httptest.NewRecorder()
	body := `{"userId":"user-1","role":"uploader","xpSharePercent":100}`
	req, _ := http.NewRequest("POST", "/uploads/co-uploaders/42", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestRemoveCoUploader_Success(t *testing.T) {
	mockSvc := new(mockService)
	_, r := setupHandlerTest(mockSvc, nil)

	trackID := int64(42)
	userID := "user-1"
	mockSvc.On("RemoveCoUploader", mock.Anything, trackID, userID).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/uploads/co-uploaders/42/user-1", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUploadDraftFile_Success(t *testing.T) {
	mockSvc := new(mockService)
	mockStorage := new(mockStorage)
	_, r := setupHandlerTest(mockSvc, mockStorage)

	draftID := "draft_1"
	mockSvc.On("AttachFile", mock.Anything, draftID, testUserID, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockStorage.On("Upload", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockStorage.On("GetURL", mock.Anything, mock.Anything).Return("https://example.com/audio.mp3", nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Disposition", `form-data; name="audio"; filename="song.mp3"`)
	hdr.Set("Content-Type", "audio/mpeg")
	part, _ := writer.CreatePart(hdr)
	// Use real MP3 magic bytes (ID3v2 tag) so http.DetectContentType
	// returns audio/mpeg instead of application/octet-stream.
	io.Copy(part, bytes.NewReader([]byte("ID3\x03\x00\x00\x00\x00\x00\x00\x00\x00")))
	writer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/uploads/drafts/"+draftID+"/file", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestUploadDraftFile_NoFile(t *testing.T) {
	mockSvc := new(mockService)
	mockStorage := new(mockStorage)
	_, r := setupHandlerTest(mockSvc, mockStorage)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/uploads/drafts/draft_1/file", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/octet-stream")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockStorage.AssertNotCalled(t, "Upload", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestListPendingDrafts_Success(t *testing.T) {
	mockSvc := new(mockService)
	_, r := setupHandlerTest(mockSvc, nil)

	expected := []Draft{{ID: "draft_1", UploadedBy: testUserID}}
	mockSvc.On("ListPendingDrafts", mock.Anything, (*string)(nil), 20, 0).Return(expected, 1, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/uploads/review", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

