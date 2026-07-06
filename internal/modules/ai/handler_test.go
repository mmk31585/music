package ai

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func setupTestHandler() *Handler {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.id::text, t.title, COALESCE(a.name, '') as artist`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "album", "genre", "duration", "cover_url"}).
			AddRow("test-id", "Test", "Artist", "Album", "Pop", 200, ""))

	mock.ExpectExec(`INSERT INTO track_embeddings_text`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "v1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewRepository(sqlxDB)
	svc := NewService(repo, &aiClient{embedding: []float64{0.1, 0.2, 0.3}}, zap.NewNop(), true)
	return NewHandler(svc)
}

func TestHandler_GenerateEmbedding(t *testing.T) {
	h := setupTestHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := EmbeddingRequest{TrackIDs: []string{uuid.New().String()}}
	data, _ := json.Marshal(body)
	c.Request = httptest.NewRequest(http.MethodPost, "/embeddings", bytes.NewReader(data))
	c.Request.Header.Set("Content-Type", "application/json")

	h.GenerateEmbedding(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GenerateEmbedding_NoTrackIDs(t *testing.T) {
	h := setupTestHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := EmbeddingRequest{TrackIDs: []string{}}
	data, _ := json.Marshal(body)
	c.Request = httptest.NewRequest(http.MethodPost, "/embeddings", bytes.NewReader(data))
	c.Request.Header.Set("Content-Type", "application/json")

	h.GenerateEmbedding(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandler_AnalyzeMood_InvalidRequest(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := NewRepository(sqlxDB)
	svc := NewService(repo, &aiClient{}, zap.NewNop(), true)
	h := NewHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/moods", bytes.NewReader([]byte(`{}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	h.AnalyzeMood(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandler_GetMood_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT track_id, mood_tags, energy, valence, tempo`)).
		WillReturnError(sqlmock.ErrCancelled)

	repo := NewRepository(sqlxDB)
	svc := NewService(repo, &aiClient{}, zap.NewNop(), true)
	h := NewHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "trackId", Value: uuid.New().String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/moods/track-id", nil)

	h.GetMood(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestHandler_SimilarByMood(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT t.id, t.title FROM tracks t`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("1", "A"))

	repo := NewRepository(sqlxDB)
	svc := NewService(repo, &aiClient{}, zap.NewNop(), true)
	h := NewHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "trackId", Value: uuid.New().String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/similar/mood/track-id?mood=calm", nil)

	h.SimilarByMood(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandler_SimilarByEmbedding(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT track_id, embedding, model_version, updated_at FROM track_embeddings_text`)).
		WillReturnRows(sqlmock.NewRows([]string{"track_id", "embedding", "model_version", "updated_at"}).
			AddRow(uuid.New(), "{0.1,0.2}", "v1", time.Now()))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM track_embeddings_audio te`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "album", "genre", "duration", "cover_url"}).
			AddRow("1", "A", "Artist", "Album", "Pop", 200, ""))

	repo := NewRepository(sqlxDB)
	svc := NewService(repo, &aiClient{}, zap.NewNop(), true)
	h := NewHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "trackId", Value: uuid.New().String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/similar/embedding/track-id", nil)

	h.SimilarByEmbedding(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandler_handleError_NotFound(t *testing.T) {
	h := NewHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.handleError(c, ErrNotFound)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestHandler_handleError_Generic(t *testing.T) {
	h := NewHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.handleError(c, http.ErrNoLocation)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
