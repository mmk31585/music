# Track Description Generator — Implementation Guide

## Overview

Add a `TrackDescription` model, a `GenerateTrackDescription` method on the `AIClient` interface, and wire it through the existing handler → service → repository stack. Since this is a pure text generation (no numeric analysis, no DB persistence required), the footprint is light.

---

## Files to touch

### 1. `internal/modules/ai/client.go` — Add new method to `AIClient` interface

Add `GenerateTrackDescription` to the interface:

```
GenerateTrackDescription(ctx context.Context, title, artist, genre, album string, duration int) (string, error)
```

Implement it on `*openAIClient` using the same chat completion pattern as `AnalyzeMood` and `GeneratePlaylist` — system prompt asking for a 2–3 paragraph description covering genre, mood, instrumentation, and vibe. Return the raw string, not JSON. Add a stub on `*fallbackClient` that returns a heuristic description built from the genre/artist/title (e.g. `"A {genre} track by {artist} — {title}."`).

### 2. `internal/modules/ai/model.go` — New domain model

```
type TrackDescription struct {
    TrackID     string `json:"track_id"`
    Title       string `json:"title"`
    Artist      string `json:"artist,omitempty"`
    Description string `json:"description"`
    ModelUsed   string `json:"model_used"`
    GeneratedAt string `json:"generated_at"`
}
```

No DB struct needed unless you want to cache descriptions. If caching is desired, add a `track_descriptions` table later — start in-memory or on-demand.

### 3. `internal/modules/ai/dto.go` — Request/response types

```
type TrackDescriptionRequest struct {
    TrackID string `json:"track_id" binding:"required"`
}

type TrackDescriptionResponse struct {
    TrackID     string `json:"track_id"`
    Title       string `json:"title"`
    Artist      string `json:"artist,omitempty"`
    Description string `json:"description"`
    GeneratedAt string `json:"generated_at"`
}
```

### 4. `internal/modules/ai/service.go` — Business logic

```
func (s *Service) GenerateTrackDescription(ctx context.Context, trackID string) (*TrackDescription, error)
```

Flow:
1. Parse+validate `trackID` as UUID.
2. Fetch metadata via `s.repo.GetTrackMetadata(ctx, trackID)` (already exists).
3. Call `s.ai.GenerateTrackDescription(ctx, meta.Title, meta.Artist, meta.Genre, meta.Album, meta.Duration)`.
4. Wrap into `*TrackDescription` and return.

The existing `generated_at` pattern is `time.Now().UTC().Format(time.RFC3339)` (see `GeneratePlaylist`).

### 5. `internal/modules/ai/handler.go` — HTTP handler

```
func (h *Handler) GenerateTrackDescription(c *gin.Context) {
    var req TrackDescriptionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
        return
    }
    desc, err := h.service.GenerateTrackDescription(c.Request.Context(), req.TrackID)
    if err != nil {
        if strings.Contains(err.Error(), "AI features are disabled") {
            c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"success": true, "data": desc})
}
```

### 6. `internal/modules/ai/routes.go` — Wire the route

```
aiGroup.POST("/descriptions", h.GenerateTrackDescription)
```

---

## Sequence diagram

```
POST /api/v1/ai/descriptions
  │
  ├─ Handler.GenerateTrackDescription
  │   ├─ Bind JSON → TrackDescriptionRequest
  │   └─ ctx → Service.GenerateTrackDescription
  │       ├─ repo.GetTrackMetadata(trackID)
  │       │   └─ SELECT … FROM tracks LEFT JOIN artists/albums/genres
  │       ├─ ai.GenerateTrackDescription(title, artist, genre, album, duration)
  │       │   └─ POST /v1/chat/completions  (system: "Write a vivid paragraph…")
  │       └─ return *TrackDescription
  └─ JSON 200 {success, data}
```

---

## OpenAI prompt design

```
system: "You are a music description writer. Write 2–3 engaging sentences describing the track's genre, mood, instrumentation, and overall vibe. Be vivid but concise."

user: "Track: '{title}' by {artist}\nGenre: {genre}\nAlbum: {album}\nDuration: {duration}s\n\nDescribe the track."
```

Temperature 0.7, max_tokens 250, no JSON response format (just raw text).

---

## Optional: caching descriptions

If you want to avoid regenerating on every request, add to `repository.go`:

```
func (r *Repository) SaveTrackDescription(ctx context.Context, desc TrackDescription) error
func (r *Repository) GetTrackDescription(ctx context.Context, trackID uuid.UUID) (*TrackDescription, error)
```

Then in the service: check cache first, fall through to AI, save result. This matches the existing pattern in `GetMood`/`AnalyzeMood`.
