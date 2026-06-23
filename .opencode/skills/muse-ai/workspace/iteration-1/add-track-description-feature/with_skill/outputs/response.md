# Track Description Generator — Implementation Guide

This follows the standard 8-step AI capability pattern from the skill. Here's exactly what to touch in the Go backend, with code references.

## 1. `internal/modules/ai/client.go` — Add to `AIClient` interface + implement

**Add the interface method** (line 13-17):
```go
type AIClient interface {
    GenerateEmbedding(...)
    AnalyzeMood(...)
    GeneratePlaylist(...)
    GenerateTrackDescription(ctx context.Context, title, artist, genre, album string, moodTags json.RawMessage) (string, error)
}
```

**Implement on `openAIClient`** — follow the `AnalyzeMood` pattern (lines 120-171). Prompt should ask for a fluent paragraph covering genre, mood, instrumentation, and vibe. Use `response_format: { type: "json_object" }` with a `{ "description": "..." }` schema, temperature ~0.5. The key difference from `AnalyzeMood`: this returns freeform text, not structured floats.

**Implement on `fallbackClient`** (lines 239-343) — return a templated string like `"A {genre} track by {artist}: {mood_tags} vibes with a {genre}-typical instrumentation."` to avoid breaking when no API key is configured.

## 2. `internal/modules/ai/model.go` — New model type (or reuse)

Add a `TrackDescription` struct and a `TrackDescriptionRequest`:

```go
type TrackDescription struct {
    TrackID     uuid.UUID `db:"track_id" json:"track_id"`
    Description string    `db:"description" json:"description"`
    ModelUsed   string    `db:"model_used" json:"model_used"`
    UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
```

You could also store it as a new column `ai_description TEXT` on the `tracks` table, or in a separate `track_descriptions` table like `track_moods`. A separate table avoids bloating the tracks table and keeps the AI module self-contained — mirroring the `track_moods` pattern.

## 3. `internal/modules/ai/repository.go` — New DB methods

Add `UpsertTrackDescription` and `GetTrackDescription` following the `UpsertMood` (line 66) / `GetMood` (line 86) pattern:

```go
func (r *Repository) UpsertTrackDescription(ctx context.Context, td TrackDescription) error
func (r *Repository) GetTrackDescription(ctx context.Context, trackID uuid.UUID) (*TrackDescription, error)
```

## 4. `internal/modules/ai/service.go` — Business logic

Add a `GenerateTrackDescription` method following the `AnalyzeMood` pattern (lines 59-85):

```go
func (s *Service) GenerateTrackDescription(ctx context.Context, trackID string) (*TrackDescription, error) {
    // 1. Check s.enabled
    // 2. Parse UUID, fetch track metadata via s.repo.GetTrackMetadata
    // 3. Optionally fetch existing mood data for richer context
    // 4. Call s.ai.GenerateTrackDescription(ctx, meta.Title, meta.Artist, meta.Genre, meta.Album, moodTags)
    // 5. Upsert via s.repo.UpsertTrackDescription
    // 6. Return *TrackDescription
}
```

Add a `GetTrackDescription` read method (like `GetMood` at line 87), and optionally a `ProcessBatchDescriptions` (like `ProcessBatchMoods` at line 342) for the worker.

## 5. `internal/modules/ai/dto.go` — Request/response DTOs

Add:
```go
type TrackDescriptionRequest struct {
    TrackID string `json:"track_id" binding:"required"`
}

type TrackDescriptionResponse struct {
    TrackID     string `json:"track_id"`
    Description string `json:"description"`
    ModelUsed   string `json:"model_used"`
}
```

## 6. `internal/modules/ai/handler.go` — HTTP handlers

Add two handlers following the `AnalyzeMood` (line 50) / `GetMood` (line 84) pattern:

```go
func (h *Handler) GenerateTrackDescription(c *gin.Context) { ... }
func (h *Handler) GetTrackDescription(c *gin.Context) { ... }
```

`GenerateTrackDescription` reads `TrackDescriptionRequest` from JSON body, calls `service.GenerateTrackDescription`, returns `TrackDescriptionResponse`. `GetTrackDescription` reads `:trackId` from URL param, calls `service.GetTrackDescription`.

## 7. `internal/modules/ai/routes.go` — Route registration

Add under the `aiGroup` block (line 10-17):
```go
aiGroup.POST("/tracks/description", h.GenerateTrackDescription)
aiGroup.GET("/tracks/description/:trackId", h.GetTrackDescription)
```

## 8. `migrations/000035_track_descriptions.sql` — New migration

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS track_descriptions (
    track_id UUID PRIMARY KEY REFERENCES tracks(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    model_used VARCHAR(100) NOT NULL DEFAULT 'v1',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose Down
DROP TABLE IF EXISTS track_descriptions;
```

This follows the naming convention at `migrations/000034_add_lyrics_ai_metadata.sql`.

## 9. (Optional) `internal/workers/ai/worker.go` — Batch processing

Add a batch description call after the existing embedding/mood batches (lines 40-53), using `ProcessBatchDescriptions` with a limit of 20.

## 10. (Optional) `frontend/src/services/api/ai/` — Frontend integration

- **`enums.ts`**: Add `GENERATE_TRACK_DESCRIPTION = '/ai/tracks/description'` and `GET_TRACK_DESCRIPTION = '/ai/tracks/description/:trackId'`
- **`types.ts`**: Add `TrackDescriptionResponseSchema` (Zod), `TrackDescriptionPayload` interface
- **`routes.ts`**: Add `generateTrackDescription(payload)` and `getTrackDescription(trackId)` methods to `useAIApi()` composable

---

## Key design decisions

| Decision | Recommendation | Rationale |
|----------|---------------|-----------|
| **Storage** | Separate `track_descriptions` table | Matches existing `track_embeddings` / `track_moods` pattern; keeps AI concerns self-contained in the module |
| **Prompt location** | Inside `openAIClient.GenerateTrackDescription` in `client.go` | Following the skill's rule: "Keep prompts in the client layer (not service)" — the prompt template and JSON schema live alongside the API call |
| **Fallback** | Template-based string in `fallbackClient` | Ensures the feature doesn't 500 when no API key is set, same pattern as `fallbackClient.AnalyzeMood` |
| **Mood context** | Optionally fetch `track_moods` in the service and pass to the AI client | Richer descriptions — the AI can reference energy, valence, and mood tags in its narrative |
| **Batch worker** | Add `ProcessBatchDescriptions` to service, call from worker | Auto-generates descriptions for tracks that lack them during the 10-minute batch cycle |
