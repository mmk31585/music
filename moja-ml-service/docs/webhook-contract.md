# Webhook Contract: Python ML Service → Go Backend

After the ML service finishes processing a lyrics job, it delivers the
result back to the Go backend via this HTTP callback.

---

## Endpoint

```
POST {GO_BACKEND_CALLBACK_URL}/internal/v1/lyrics/callback
```

`GO_BACKEND_CALLBACK_URL` is configured in both services — typically
`http://go-backend:8080/api/v1` inside Docker Compose.

---

## Headers

| Header              | Value                                                        |
|---------------------|--------------------------------------------------------------|
| `Content-Type`      | `application/json`                                           |
| `X-Moja-Signature`  | `hex(HMAC-SHA256(secret, "{timestamp}.{raw_body}"))`         |
| `X-Moja-Timestamp`  | Unix epoch seconds (integer, as string), e.g. `"1712345678"` |

---

## Request Body

```json
{
  "track_id": "string",
  "lrc_content": "string",
  "plain_text": "string",
  "confidence": 0.0,
  "detected_language": "string",
  "whisper_model_version": "string"
}
```

| Field                  | Type          | Description                                        |
|------------------------|---------------|----------------------------------------------------|
| `track_id`             | string        | Track identifier matching Go's catalog              |
| `lrc_content`          | string        | LRC-format lyrics with timestamps (null if failed)  |
| `plain_text`           | string        | Plain-text transcription                           |
| `confidence`           | float (0–1)   | Duration-weighted average confidence score          |
| `detected_language`    | string        | ISO 639-1 language code detected by Whisper         |
| `whisper_model_version`| string        | Model size, e.g. `"medium-int8"`                    |

---

## Signing Scheme

This is the **critical security contract** — both sides must implement
it identically.

### Python side (sending)

```python
import hmac, hashlib, json

timestamp = str(int(time.time()))
body_bytes = json.dumps(payload, ensure_ascii=False, separators=(",", ":")).encode()
message = f"{timestamp}.{body_bytes.decode()}".encode()
signature = hmac.new(secret.encode(), message, hashlib.sha256).hexdigest()
```

Key details:

1. **`separators=(",", ":")`** — compact JSON, no extra whitespace.
2. **`ensure_ascii=False`** — Persian/Arabic text preserved as UTF-8.
3. The signed string is `"{timestamp}.{raw_body}"` (period-separated).
4. The *raw_body* is the **exact byte sequence** sent over the wire.
5. The timestamp is included in the signed payload to prevent replay
   attacks.

### Go side (verifying)

```go
import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "strconv"
    "strings"
    "time"
)

func verifyCallback(rawBody []byte, sigHeader, tsHeader, secret string) bool {
    // 1. Parse and validate timestamp (replay protection)
    ts, err := strconv.ParseInt(tsHeader, 10, 64)
    if err != nil || time.Now().Unix() - ts > 300 {
        return false
    }

    // 2. Recompute signature over the RAW body bytes
    //    (NOT a re-marshalled version of the parsed JSON — whitespace
    //     or key ordering differences would break the match)
    message := fmt.Sprintf("%s.%s", tsHeader, string(rawBody))
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(message))
    expected := hex.EncodeToString(mac.Sum(nil))

    // 3. Constant-time comparison — NEVER use == for HMAC comparison
    return hmac.Equal([]byte(sigHeader), []byte(expected))
}
```

**⚠️ Critical Go-side requirement:** Read `rawBody` *before* parsing JSON.
The HMAC is verified against the raw bytes. If Go parses JSON first and
then re-marshals it, the byte order and whitespace may differ, breaking
the signature check.

---

## Response Handling

| Go Response             | Python Behaviour                                    |
|-------------------------|-----------------------------------------------------|
| `200` or `204`          | ✅ Delivery confirmed. Job marked `completed`.       |
| `400`–`499`             | ❌ Permanent rejection. Python does **NOT** retry.   |
| `500`–`599`             | 🔁 Transient error. Python **retries** via Celery.   |
| Timeout / connection error | 🔁 Transient. Python **retries** via Celery.       |

---

## Replay Protection

- `X-Moja-Timestamp` is included in the HMAC payload.
- Go **must** reject requests where `now - timestamp > 300 seconds`.
- This limits the window for replay attacks to 5 minutes.

---

---

## Video Processing Callback

After the ML service finishes processing a video job (audio replacement), it
delivers the result back to the Go backend via this HTTP callback.

### Endpoint

```
POST {GO_BACKEND_CALLBACK_URL}/internal/v1/video/callback
```

### Headers

Same HMAC signing scheme as the lyrics callback — see above.

### Request Body

```json
{
  "video_id": "string",
  "final_video_path": "string",
  "thumbnail_path": "string",
  "duration_ms": 0,
  "aspect_ratio": "string",
  "processing_status": "completed",
  "error_message": null
}
```

| Field                  | Type          | Description                                        |
|------------------------|---------------|----------------------------------------------------|
| `video_id`             | string        | Video UUID matching Go's video module               |
| `final_video_path`     | string        | Path/URL to the processed video (audio replaced)    |
| `thumbnail_path`       | string        | Path/URL to the extracted thumbnail (JPEG)          |
| `duration_ms`          | integer       | Video duration in milliseconds                      |
| `aspect_ratio`         | string        | Detected aspect ratio, e.g. `"16:9"` or `"9:16"`   |
| `processing_status`    | string        | `"completed"` or `"failed"`                         |
| `error_message`        | string\|null  | Error description if failed                         |

### Response Handling

Same as lyrics callback: 2xx = confirmed, 4xx = permanent rejection,
5xx/timeout = transient (retry).

---

## Implementation Files

| Service   | File                                                                    |
|-----------|-------------------------------------------------------------------------|
| Python    | `app/infra/callback/go_client.py` — HMAC signing + delivery (shared)    |
| Python    | `app/infra/callback/video_client.py` — Video-specific callback client   |
| Python    | `app/infra/callback/cover_client.py` — Cover-specific callback client   |
| Go        | (Phase 3/5 — to be implemented in Go backend)                           |
