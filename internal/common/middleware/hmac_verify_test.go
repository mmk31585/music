package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-hmac-secret-12345"

func setupGinTest(method, path string, body []byte, headers map[string]string, hmacSecret string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	ctx.Request = req

	handler := VerifyMLServiceWebhook(hmacSecret)
	handler(ctx)

	return w
}

func TestVerifyMLServiceWebhook_RejectsMissingHeaders(t *testing.T) {
	body := []byte(`{"track_id":"abc"}`)

	t.Run("no signature header", func(t *testing.T) {
		w := setupGinTest("POST", "/internal/v1/lyrics/callback", body,
			map[string]string{"X-Moja-Timestamp": "1712345678"}, testSecret)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Contains(t, resp["error"], "missing")
	})

	t.Run("no timestamp header", func(t *testing.T) {
		w := setupGinTest("POST", "/internal/v1/lyrics/callback", body,
			map[string]string{"X-Moja-Signature": "abc123"}, testSecret)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("no headers at all", func(t *testing.T) {
		w := setupGinTest("POST", "/internal/v1/lyrics/callback", body, nil, testSecret)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestVerifyMLServiceWebhook_RejectsStaleTimestamp(t *testing.T) {
	body := []byte(`{"track_id":"abc"}`)
	oldTS := strconv.FormatInt(time.Now().Add(-10*time.Minute).Unix(), 10)
	oldSig := computeSignature(oldTS, string(body), testSecret)

	w := setupGinTest("POST", "/internal/v1/lyrics/callback", body,
		map[string]string{
			"X-Moja-Timestamp": oldTS,
			"X-Moja-Signature": oldSig,
		}, testSecret)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "stale")
}

func computeSignature(ts, rawBody, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ts + "." + rawBody))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyMLServiceWebhook_RejectsInvalidSignature(t *testing.T) {
	body := []byte(`{"track_id":"abc"}`)
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	w := setupGinTest("POST", "/internal/v1/lyrics/callback", body,
		map[string]string{
			"X-Moja-Timestamp": ts,
			"X-Moja-Signature": "this-is-not-a-valid-signature",
		}, testSecret)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "invalid")
}

func TestVerifyMLServiceWebhook_AcceptsValidSignature(t *testing.T) {
	body := []byte(`{"track_id":"abc","lrc_content":"[00:01.00]test"}`)
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sig := computeSignature(ts, string(body), testSecret)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("POST", "/internal/v1/lyrics/callback", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Moja-Timestamp", ts)
	req.Header.Set("X-Moja-Signature", sig)
	ctx.Request = req

	// Next handler captures whether middleware passed
	var passed bool
	handler := VerifyMLServiceWebhook(testSecret)
	handler(ctx)
	if !ctx.IsAborted() {
		passed = true
		// Simulate the next handler reading the body
		restoredBody, _ := io.ReadAll(ctx.Request.Body)
		assert.JSONEq(t, string(body), string(restoredBody), "body should be restored")
	}

	assert.True(t, passed, "middleware should not abort for valid signature")
	assert.Equal(t, http.StatusOK, w.Code) // no handler wrote, so default 200
}

func TestVerifyMLServiceWebhook_EmptySecretRejects(t *testing.T) {
	body := []byte(`{"track_id":"abc"}`)
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sig := computeSignature(ts, string(body), "")

	w := setupGinTest("POST", "/internal/v1/lyrics/callback", body,
		map[string]string{
			"X-Moja-Timestamp": ts,
			"X-Moja-Signature": sig,
		}, "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "not configured")
}

func TestVerifyMLServiceWebhook_RestoresBodyForNextHandler(t *testing.T) {
	// Verify that the middleware reads and restores the body so the
	// downstream handler can still bind JSON.
	body := []byte(`{"hello":"world"}`)
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sig := computeSignature(ts, string(body), testSecret)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("POST", "/internal/v1/lyrics/callback", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Moja-Timestamp", ts)
	req.Header.Set("X-Moja-Signature", sig)
	ctx.Request = req

	mw := VerifyMLServiceWebhook(testSecret)
	mw(ctx)

	// After middleware, the next handler should be able to bind JSON
	var payload struct {
		Hello string `json:"hello"`
	}
	err := ctx.ShouldBindJSON(&payload)
	require.NoError(t, err)
	assert.Equal(t, "world", payload.Hello)
}
