package media

import (
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockMediaService struct {
	mock.Mock
}

func (m *mockMediaService) Upload(ctx context.Context, category UploadCategory, file multipart.File, header *multipart.FileHeader, createdBy *uuid.UUID) (*UploadResponse, error) {
	args := m.Called(ctx, category, file, header, createdBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UploadResponse), args.Error(1)
}

func (m *mockMediaService) ListMedia(ctx context.Context) ([]Media, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Media), args.Error(1)
}

func (m *mockMediaService) DeleteMedia(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupMediaHandlerTest() (*Handler, *mockMediaService, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockMediaService)
	h := NewHandler(mockSvc)

	r := gin.New()
	// Direct handler test: test handler methods individually
	return h, mockSvc, r
}

func TestDeleteAdminMedia_Success(t *testing.T) {
	h, mockSvc, r := setupMediaHandlerTest()

	mediaID := uuid.New().String()
	mockSvc.On("DeleteMedia", mock.Anything, mock.Anything).Return(nil)

	r.DELETE("/admin/media/:id", h.DeleteAdminMedia)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/admin/media/"+mediaID, nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeleteAdminMedia_InvalidID(t *testing.T) {
	h, _, r := setupMediaHandlerTest()

	r.DELETE("/admin/media/:id", h.DeleteAdminMedia)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/admin/media/not-a-uuid", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.False(t, resp["success"].(bool))
}

func TestListAdminMedia_Success(t *testing.T) {
	h, mockSvc, r := setupMediaHandlerTest()

	expected := []Media{
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), MediaType: "image"},
		{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), MediaType: "audio"},
	}
	mockSvc.On("ListMedia", mock.Anything).Return(expected, nil)

	r.GET("/admin/media", h.ListAdminMedia)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/media", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestListAdminMedia_Error(t *testing.T) {
	h, mockSvc, r := setupMediaHandlerTest()

	mockSvc.On("ListMedia", mock.Anything).Return(nil, errors.New("db error"))

	r.GET("/admin/media", h.ListAdminMedia)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/media", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}
