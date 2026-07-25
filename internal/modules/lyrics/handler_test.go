package lyrics

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testHMACSecret = "test-ml-secret-for-testing"

func computeHMACSignature(t *testing.T, ts, body, secret string) string {
	t.Helper()
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ts + "." + body))
	return hex.EncodeToString(mac.Sum(nil))
}

func setupCallbackTest(t *testing.T, payload map[string]any) (*gin.Context, *httptest.ResponseRecorder, *Handler, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewRepository(sqlxDB)
	svc := NewService(repo)
	handler := NewHandler(svc)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	bodyBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sig := computeHMACSignature(t, ts, string(bodyBytes), testHMACSecret)

	req := httptest.NewRequest("POST", "/internal/v1/lyrics/callback", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Moja-Timestamp", ts)
	req.Header.Set("X-Moja-Signature", sig)
	ctx.Request = req

	return ctx, w, handler, mock
}

func TestHandleLyricsCallback_SavesLyricsCorrectly(t *testing.T) {
	payload := map[string]any{
		"track_id":              "550e8400-e29b-41d4-a716-446655440000",
		"lrc_content":           "[00:01.00]Hello world\n[00:05.00]Test line",
		"plain_text":            "",
		"confidence":            0.95,
		"detected_language":     "en",
		"whisper_model_version": "medium-int8",
	}

	ctx, w, handler, mock := setupCallbackTest(t, payload)

	// Expect track existence check
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM tracks WHERE id = \$1\)`).
		WithArgs("550e8400-e29b-41d4-a716-446655440000").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Expect lyrics uniqueness check
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM lyrics WHERE track_id = \$1 AND language = \$2\)`).
		WithArgs("550e8400-e29b-41d4-a716-446655440000", "en").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Expect INSERT with source = "ai_generated" and confidence score
	confidence := 0.95
	mock.ExpectQuery(`INSERT INTO lyrics`).
		WithArgs(
			"550e8400-e29b-41d4-a716-446655440000",
			"en",
			"lrc",
			"[00:01.00]Hello world\n[00:05.00]Test line",
			"ai_generated",
			&confidence,
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "track_id", "language", "type", "content", "source", "confidence_score", "created_at", "updated_at",
		}).AddRow(
			"660e8400-e29b-41d4-a716-446655440001",
			"550e8400-e29b-41d4-a716-446655440000",
			"en",
			"lrc",
			"[00:01.00]Hello world\n[00:05.00]Test line",
			"ai_generated",
			&confidence,
			time.Now(),
			time.Now(),
		))

	handler.HandleLyricsCallback(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "received", resp["status"])
	assert.Equal(t, "660e8400-e29b-41d4-a716-446655440001", resp["lyrics_id"])

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHandleLyricsCallback_RejectsUnknownTrackID(t *testing.T) {
	payload := map[string]any{
		"track_id":    "00000000-0000-0000-0000-000000000000",
		"lrc_content": "[00:01.00]Test",
		"confidence":  0.8,
	}

	ctx, w, handler, mock := setupCallbackTest(t, payload)

	// Track does not exist
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM tracks WHERE id = \$1\)`).
		WithArgs("00000000-0000-0000-0000-000000000000").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	handler.HandleLyricsCallback(ctx)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "track not found")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHandleLyricsCallback_RejectsDuplicateLyrics(t *testing.T) {
	payload := map[string]any{
		"track_id":          "550e8400-e29b-41d4-a716-446655440000",
		"lrc_content":       "[00:01.00]Already exists",
		"confidence":        0.8,
		"detected_language": "fa",
	}

	ctx, w, handler, mock := setupCallbackTest(t, payload)

	// Track exists
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM tracks WHERE id = \$1\)`).
		WithArgs("550e8400-e29b-41d4-a716-446655440000").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// But lyrics already exist for this track + language
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM lyrics WHERE track_id = \$1 AND language = \$2\)`).
		WithArgs("550e8400-e29b-41d4-a716-446655440000", "fa").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	handler.HandleLyricsCallback(ctx)

	assert.Equal(t, http.StatusConflict, w.Code)
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "already exist")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHandleLyricsCallback_RejectsMissingContent(t *testing.T) {
	payload := map[string]any{
		"track_id":    "550e8400-e29b-41d4-a716-446655440000",
		"lrc_content": "",
		"plain_text":  "",
		"confidence":  0.8,
	}

	ctx, w, handler, mock := setupCallbackTest(t, payload)

	handler.HandleLyricsCallback(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "no lyrics content")

	// No DB queries should be made
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHandleLyricsCallback_RejectsInvalidPayload(t *testing.T) {
	// Missing track_id (required)
	payload := map[string]any{
		"lrc_content": "[00:01.00]Test",
	}

	ctx, w, handler, _ := setupCallbackTest(t, payload)

	handler.HandleLyricsCallback(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "invalid payload")
}

func TestHandleLyricsCallback_FallsBackToPlainText(t *testing.T) {
	payload := map[string]any{
		"track_id":    "550e8400-e29b-41d4-a716-446655440000",
		"lrc_content": "",
		"plain_text":  "Hello world\nTest line",
		"confidence":  0.85,
	}

	ctx, w, handler, mock := setupCallbackTest(t, payload)

	// Track exists
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM tracks WHERE id = \$1\)`).
		WithArgs("550e8400-e29b-41d4-a716-446655440000").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Lyrics don't exist
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM lyrics WHERE track_id = \$1 AND language = \$2\)`).
		WithArgs("550e8400-e29b-41d4-a716-446655440000", "fa").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Expect INSERT with plain type
	confidence := 0.85
	mock.ExpectQuery(`INSERT INTO lyrics`).
		WithArgs(
			"550e8400-e29b-41d4-a716-446655440000",
			"fa",
			"plain",
			"Hello world\nTest line",
			"ai_generated",
			&confidence,
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "track_id", "language", "type", "content", "source", "confidence_score", "created_at", "updated_at",
		}).AddRow(
			"660e8400-e29b-41d4-a716-446655440002",
			"550e8400-e29b-41d4-a716-446655440000",
			"fa",
			"plain",
			"Hello world\nTest line",
			"ai_generated",
			&confidence,
			time.Now(),
			time.Now(),
		))

	handler.HandleLyricsCallback(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "received", resp["status"])

	require.NoError(t, mock.ExpectationsWereMet())
}
