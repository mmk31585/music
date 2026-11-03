package upload

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	apperrors "music/internal/common/errors"
)

const testUserID = "auth_user_id"

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

func setupHandlerTest() (*Handler, *mockService, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockService)
	h := NewHandler(mockSvc, nil)

	r := gin.New()
	r.Use(setTestUser)
	return h, mockSvc, r
}

// setTestUser injects an authenticated user into Gin context so handlers
// (which read auth.UserIDFromContext) work without JWT middleware.
func setTestUser(c *gin.Context) {
	c.Set("auth_user_id", testUserID)
	c.Next()
}

func TestCreateDraft_Success(t *testing.T) {
	h, mockSvc, r := setupHandlerTest()

	expected := Draft{
		ID:                "draft_1",
		UploadedBy:        testUserID,
		OriginalFilename:  "song.mp3",
		Status:            "pending",
		UploadSource:      string(SourceListener),
		NeedsReview:       true,
	}
	mockSvc.On("CreateDraft", mock.Anything, testUserID, mock.Anything).Return(expected, nil)

	r.POST("/uploads/drafts", responseWrap(h.CreateDraft))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/uploads/drafts", strings.NewReader(`{"filename":"song.mp3"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCreateDraft_Unauthenticated(t *testing.T) {
	h, _, _ := setupHandlerTest()

	// No user in Gin context: handler should return 401
	r2 := gin.New()
	r2.POST("/uploads/drafts", responseWrap(h.CreateDraft))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/uploads/drafts", strings.NewReader(`{"filename":"song.mp3"}`))
	req.Header.Set("Content-Type", "application/json")
	r2.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateDraft_InvalidBody(t *testing.T) {
	h, mockSvc, r := setupHandlerTest()

	r.POST("/uploads/drafts", responseWrap(h.CreateDraft))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/uploads/drafts", strings.NewReader(`not-json`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "CreateDraft", mock.Anything, testUserID, mock.Anything)
}

func TestGetDraft_Success(t *testing.T) {
	h, mockSvc, r := setupHandlerTest()

	draftID := "draft_1"
	expected := Draft{ID: draftID, UploadedBy: testUserID, Status: "pending"}
	mockSvc.On("GetDraft", mock.Anything, draftID).Return(expected, nil)

	r.GET("/uploads/drafts/:draftId", responseWrap(h.GetDraft))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/uploads/drafts/"+draftID, nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetDraft_NotFound(t *testing.T) {
	h, mockSvc, r := setupHandlerTest()

	draftID := "draft_missing"
	mockSvc.On("GetDraft", mock.Anything, draftID).Return(Draft{}, apperrors.New(http.StatusNotFound, apperrors.CodeNotFound, "draft not found", nil))

	r.GET("/uploads/drafts/:draftId", responseWrap(h.GetDraft))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/uploads/drafts/"+draftID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestListMyDrafts_Success(t *testing.T) {
	h, mockSvc, r := setupHandlerTest()

	expected := []Draft{{ID: "draft_1", UploadedBy: testUserID}}
	mockSvc.On("ListUserDrafts", mock.Anything, testUserID, 20, 0).Return(expected, 1, nil)

	r.GET("/uploads/drafts/mine", responseWrap(h.ListMyDrafts))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/uploads/drafts/mine", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetUploadSlots_Success(t *testing.T) {
	h, mockSvc, r := setupHandlerTest()

	expected := UploadSlotStatus{Available: true, Used: 1, Max: 5}
	mockSvc.On("GetUploadSlots", mock.Anything, testUserID).Return(expected, nil)

	r.GET("/uploads/slots", responseWrap(h.GetUploadSlots))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/uploads/slots", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestReviewDraft_Success(t *testing.T) {
	h, mockSvc, r := setupHandlerTest()

	draftID := "draft_1"
	mockSvc.On("ReviewDraft", mock.Anything, draftID, testUserID, mock.Anything).Return(nil)

	r.POST("/uploads/review/:draftId", responseWrap(h.ReviewDraft))

	w := httptest.NewRecorder()
	body := `{"action":"accept","notes":"looks good"}`
	req, _ := http.NewRequest("POST", "/uploads/review/"+draftID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetCoUploaders_Success(t *testing.T) {
	h, mockSvc, r := setupHandlerTest()

	trackID := int64(42)
	expected := []CoUploader{{TrackID: trackID, UserID: testUserID}}
	mockSvc.On("GetCoUploaders", mock.Anything, trackID).Return(expected, nil)

	r.GET("/tracks/:trackId/co-uploaders", responseWrap(h.GetCoUploaders))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tracks/42/co-uploaders", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

// responseWrap adapts a Handler method (returning error) into a gin.HandlerFunc.
func responseWrap(fn func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c); err != nil {
			apperr := apperrors.ToAppError(err)
			c.AbortWithStatusJSON(apperr.StatusCode, gin.H{"success": false, "code": apperr.Code, "message": apperr.Message})
		}
	}
}

