---
name: muse-ai
description: Use for all AI/ML work in the Muse Persian music ecosystem. Covers the Go backend AI module (OpenAI client for embeddings, mood analysis, playlist generation), the Python ML microservice (Whisper speech-to-text for lyrics, cover image optimization, Celery task queue), the frontend AI pages (playlist generator, mood explorer), the AI worker (batch embedding/mood processing), the contribution AI verifier, and all AI config/database concerns. Trigger whenever the task involves AI features, LLM integration, audio ML, lyrics generation, embeddings, mood analysis, playlist generation, cover optimization, contribution moderation, or adding new AI capabilities to the platform.
---

# Muse AI — AI/ML Platform Guide

Muse uses a **layered AI architecture**: a Go backend module handles LLM interactions (OpenAI-compatible), a separate Python microservice handles heavy ML (Whisper ASR, image processing), and the frontend provides user-facing AI features.

```
Frontend (Vue 3) ─── HTTP ──→ Go Backend (Gin) ─── HTTP ──→ OpenAI API
                                        │
                                   HMAC Callback ←── Python ML Service (FastAPI + Celery)
                                        │
                                   PostgreSQL + Redis
```

---

## 1. Go Backend AI Module (`internal/modules/ai/`)

Standard **handler → service → repository** pattern. The central interface is `AIClient`:

```go
type AIClient interface {
    GenerateEmbedding(ctx context.Context, input string) ([]float64, error)
    AnalyzeMood(ctx context.Context, title, artist, genre string) (*TrackMood, error)
    GeneratePlaylist(ctx context.Context, prompt string, tracks []TrackMeta) ([]string, error)
}
```

Two implementations:
- **`openAIClient`** — Real OpenAI-compatible API calls (configured via env vars)
- **`fallbackClient`** — Genre-based heuristics when no API key is set

### File Layout

```
internal/modules/ai/
├── model.go          # TrackEmbedding, TrackMood, MoodTag, SmartPlaylist, GenerationLog, TrackMeta
├── client.go         # AIClient interface + openAIClient + fallbackClient
├── service.go        # Business logic: GenerateEmbedding, AnalyzeMood, GeneratePlaylist, ProcessBatch*
├── handler.go        # HTTP handlers
├── repository.go     # PostgreSQL queries (track_embeddings, track_moods, ai_generation_log, smart_playlists)
├── dto.go            # Request/response DTOs
└── routes.go         # Route registration
```

### API Endpoints (`/api/v1/ai/`)

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/embeddings` | Generate embeddings for one or more tracks |
| POST | `/moods` | Analyze mood for a track (energy, valence, tempo, etc.) |
| GET | `/moods/:trackId` | Retrieve stored mood for a track |
| POST | `/playlists/generate` | AI playlist from text prompt + mood/activity/genre/seed |
| GET | `/similar/mood/:trackId` | Find similar tracks by mood profile |
| GET | `/similar/embedding/:trackId` | Find similar tracks by embedding cosine similarity |

### Key Patterns

**Adding a new AI capability** follows this flow:
1. Add a method to `AIClient` interface in `client.go`
2. Implement it on `openAIClient` (and `fallbackClient` if a fallback makes sense)
3. Add domain logic in `service.go`
4. Add repository methods in `repository.go` if DB access is needed
5. Add DTOs in `dto.go`
6. Wire up the handler in `handler.go`
7. Register the route in `routes.go`
8. Add the frontend API call in `frontend/src/services/api/ai/routes.ts`

**Prompt engineering:** The system prompt and user prompt template live inside the client method. Keep prompts in the client layer (not service) so the OpenAI interaction is self-contained. Use `JSON mode` (`response_format: { type: "json_object" }`) for structured outputs.

**Database query patterns:** All repository queries joining `track_artists` use `ta.role = 'primary'` (not `ta.position = 0`). This is the canonical way to resolve the primary artist for a track across all AI repository methods.

### Service Methods

- **`GenerateEmbedding`** — Fetches track metadata (title, artist, genre), concatenates as input, calls OpenAI embedding API, stores in `track_embeddings`
- **`AnalyzeMood`** — Sends title/artist/genre to OpenAI chat, gets structured mood analysis (8 audio features), stores in `track_moods`
- **`GetSimilarByMood`** — Queries `track_moods` by mood tag or energy/valence range (±0.2)
- **`GetSimilarByEmbedding`** — Uses PostgreSQL `cosine_distance` operator on `track_embeddings`
- **`GeneratePlaylist`** — Multi-tier candidate selection (seed embedding → mood tag → mood range → plain tracks → random), sends to OpenAI for ranking, logs to `ai_generation_log`. The 4th tier (`GetTracks`) is a last-resort fallback that fetches any available tracks from the catalog when no mood/embedding data exists.
- **`ProcessBatchEmbeddings`** / **`ProcessBatchMoods`** — Process tracks without embeddings/moods (used by background worker)

### Repository Methods

| Method | Purpose |
|--------|---------|
| `GetTracks(ctx, limit)` | Generic track fetch (fallback for playlist generation when no mood/embedding data) — joins across `tracks`, `track_artists` (with `role = 'primary'`), `artists`, `albums`, `track_genres`, `genres` |

### TrackMeta Model

```go
type TrackMeta struct {
    ID       string  `json:"id"`
    Title    string  `json:"title"`
    ArtistID string  `json:"artist_id,omitempty"`
    Artist   string  `json:"artist,omitempty"`
    Album    string  `json:"album,omitempty"`
    Genre    string  `json:"genre,omitempty"`
    Year     int     `json:"year,omitempty"`
    Duration int     `json:"duration,omitempty"`
    CoverURL string  `json:"cover_url,omitempty"`
    Energy   float64 `json:"energy,omitempty"`
    Valence  float64 `json:"valence,omitempty"`
}
```

`CoverURL` is used by frontend for track art display. `Energy`/`Valence` are populated when mood data exists and used as context for AI playlist generation.

### Database Tables

| Table | Key Columns | Purpose |
|-------|------------|---------|
| `track_embeddings` | `track_id`, `embedding` (float[]), `model_version` | Vector embeddings for similarity search |
| `track_moods` | `track_id`, `mood_tags` (jsonb), `energy`, `valence`, `tempo`, `danceability`, `acousticness`, `instrumentalness`, `liveness`, `speechiness` | Mood/audio feature predictions |
| `ai_generation_log` | `user_id`, `playlist_id`, `prompt`, `track_count`, `model_used`, `latency_ms` | Audit log for AI playlist generations |
| `smart_playlists` | `user_id`, `name`, `description`, `query_config` (jsonb), `is_active` | Saved AI playlist configurations |

---

## 2. Python ML Microservice (`moja-ml-service/`)

A standalone FastAPI + Celery service for heavy ML workloads. Communicates with the Go backend via HTTP (job submission) and HMAC-signed callbacks (result delivery).

### Architecture

```
POST /api/v1/lyrics/jobs  ──→  FastAPI  ──→  Celery  ──→  faster-whisper
                                                      └─→  Pillow (covers)
                                                             │
                                   HMAC POST ←──────────────┘
                                   → Go /api/v1/webhooks/lyrics
```

### File Layout

```
moja-ml-service/
├── app/
│   ├── config.py                     # Pydantic settings (DB, Redis, Whisper, storage, callback)
│   ├── main.py                       # FastAPI app entrypoint
│   ├── api/v1/
│   │   ├── router.py                 # Route aggregation
│   │   ├── health.py                 # Liveness + readiness probes
│   │   ├── lyrics.py                 # POST/GET lyrics jobs
│   │   └── covers.py                 # POST/GET cover jobs
│   ├── domain/
│   │   ├── lyrics/
│   │   │   ├── service.py            # Lyrics generation pipeline orchestrator
│   │   │   ├── transcriber.py        # faster-whisper wrapper (singleton)
│   │   │   ├── lrc_formatter.py      # Segment → LRC format
│   │   │   ├── persian_normalizer.py # Persian text normalization (hazm)
│   │   │   └── confidence.py         # Duration-weighted confidence scoring
│   │   ├── covers/
│   │   │   └── optimizer.py          # Pillow → WebP (3 size variants)
│   │   └── recommendation/           # Placeholder
│   ├── workers/
│   │   ├── celery_app.py             # Celery config (2 queues)
│   │   └── tasks/
│   │       ├── lyrics_tasks.py       # generate_lyrics_task
│   │       └── cover_tasks.py        # optimize_cover_task, refresh_stale_covers
│   └── infra/
│       ├── callback/
│       │   ├── go_client.py          # HMAC-SHA256 callback delivery
│       │   └── cover_client.py       # Same for cover results
│       ├── db/
│       │   ├── models.py             # SQLAlchemy: LyricsJob, CoverJob
│       │   ├── session.py            # Async DB session
│       │   └── sync_session.py       # Sync DB session (Celery)
│       └── storage/
│           ├── factory.py            # shared_volume vs presigned_url
│           └── base.py               # Abstract storage resolver
├── Dockerfile
├── docker-compose.yml
└── pyproject.toml
```

### Adding a New ML Capability

1. **Domain logic** — Create `app/domain/<feature>/` with pure functions/classes (no HTTP or Celery knowledge)
2. **API endpoints** — Add routes in `app/api/v1/<feature>.py` following the lyrics/covers patterns
3. **Celery task** — Add `app/workers/tasks/<feature>_tasks.py` with the async worker
4. **Callback** — Use `HMACCoverClient` or similar for result delivery back to Go
5. **DB model** — Add SQLAlchemy model in `app/infra/db/models.py` if you need job tracking
6. **Alembic migration** — Run `alembic revision --autogenerate` for new tables
7. **Register routes** — Wire up in `app/api/v1/router.py`
8. **Config** — Add new env vars to `app/config.py`

### Lyrics Pipeline

```
audio file → Transcriber.transcribe() → format_as_lrc() → normalize_persian_text()
                                                                    → compute_overall_confidence()
                                                                    → LyricsGenerationResult
```

The `Transcriber` is a singleton (`get_transcriber()`) — loads Whisper model once per process. Uses `faster-whisper` with `medium` model, `int8` quantization, CPU. Auto-detects language per track (don't force Persian — catalog is multi-language). VAD filter enabled.

### Cover Optimizer

Produces 3 WebP variants per image:
- `original.webp` — ≤1200px (full-res fallback)
- `med.webp` — ≤300px (album cards, lists)
- `thumb.webp` — ≤100px (mini player, search)

### HMAC Callback Contract

When a Celery task completes, it delivers results to Go via signed POST.

- **Signature:** `HMAC-SHA256` of `{job_id}.{status}.{timestamp}` with shared secret
- **Header:** `X-Signature-256`
- **Endpoint:** Configured as `go_backend_callback_url` — appends `/webhooks/lyrics` or `/webhooks/covers`
- **Retry:** Transient errors retry; permanent errors mark job as `failed`

### Configuration

All via `.env` in the `moja-ml-service/` directory:

```
DATABASE_URL=postgresql+asyncpg://moja_ml:password@localhost:5433/moja_ml
REDIS_URL=redis://localhost:6379/0
CELERY_BROKER_URL=redis://localhost:6379/1
STORAGE_MODE=shared_volume  # or presigned_url
AUDIO_SHARED_PATH=/mnt/audio-uploads
WHISPER_MODEL_SIZE=medium
GO_BACKEND_CALLBACK_URL=http://localhost:8080/api/v1
WEBHOOK_HMAC_SECRET=dev-secret-change-in-production
```

---

## 3. AI Worker (`internal/workers/ai/worker.go`)

A background worker that runs every **10 minutes** to batch-process tracks without embeddings or moods.

```go
type Worker struct {
    repo   *ai.Repository
    client ai.AIClient
    logger *zap.Logger
}
```

Processes up to 20 tracks per batch for both embeddings and moods. Uses `ai.NewFallbackClient()` by default — in production, inject the real OpenAI client.

To change the batch size or schedule, edit the worker's ticker interval and the `limit` parameter passed to `ProcessBatchEmbeddings`/`ProcessBatchMoods`.

---

## 4. Frontend AI Features

### AI API Composable (`frontend/src/services/api/ai/`)

```
frontend/src/services/api/ai/
├── enums.ts    # Route definitions (6 endpoints)
├── types.ts    # Zod schemas, TS types, mood/activity constants
├── routes.ts   # useAIApi() composable (generatePlaylist, analyzeMood, etc.)
└── index.ts    # Barrel export
```

**Usage in any component:**
```ts
import { useAIApi } from '@/services/api/ai'
const aiApi = useAIApi()
const result = await aiApi.generatePlaylist({ prompt: 'chill morning', mood: 'calm', limit: 15 })
```

**Mood and activity options** are defined as const arrays in `types.ts` — add new moods/activities there to make them available across all AI UIs.

### AI Pages

| Route | Component | Purpose |
|-------|-----------|---------|
| `/ai/playlist-generator` | `PageAIPlaylistGenerator.vue` | Text prompt + mood/activity/genre selectors → AI playlist |
| `/ai/mood-explorer` | `PageAIMoodExplorer.vue` | Browse tracks by mood with energy bar visualization |

Both pages follow the same pattern:
1. Import `useAIApi()` and call the appropriate method
2. Display results in a track list with play-all and per-track play
3. Import `usePlayer()` to queue and play tracks

**Accessibility patterns used across AI pages:**
- `<i aria-hidden="true" ...>` — all icon-only elements get `aria-hidden` to hide from screen readers
- `<button aria-label="Play track">` — icon-only buttons have descriptive labels
- Track lists use `aria-live="polite"` for dynamic content regions
- Reactive refs typed as `Record<string, unknown>[]` instead of `any[]` — maintain this pattern in new AI pages

### Adding a New AI Page

1. Create `frontend/src/pages/app/PageNewAIFeature.vue`
2. Add route in `frontend/src/router/routes/app.ts`
3. Add navigation link in `LayoutMusicApp.vue`
4. Use `useAIApi()` for backend calls, `usePlayer()` for playback

### Feature Flags

AI features are gated by `feature-flags` service:
```ts
// frontend/src/services/api/feature-flags/types.ts
ai: boolean  // Controls whether AI features are available
```

Backend: `FEATURE_AI_ENABLED` / `AI_ENABLED` env vars in `internal/config/features.go`.

---

## 5. Contribution AI Verifier (`internal/modules/contribution/ai_verifier.go`)

Rule-based content moderation for user contributions. The `AIVerifier` interface:

```go
type AIVerifier interface {
    Verify(ctx context.Context, contributionType string, data json.RawMessage) (verdict string, confidence float64, reason string, err error)
}
```

Supports contribution types: `lyrics`, `translation`, `credits`, `metadata`, `album_art`, `bio`.

The current implementation is a `ruleBasedVerifier` with hardcoded heuristics. To upgrade to an LLM-powered verifier, implement the interface with an OpenAI chat completion client similar to `internal/modules/ai/client.go`.

---

## 6. Configuration Reference

### Go Backend (`internal/config/ai.go` + `internal/config/features.go`)

| Env Var | Default | Description |
|---------|---------|-------------|
| `AI_OPENAI_ENDPOINT` | `https://api.openai.com/v1` | OpenAI-compatible API base URL |
| `AI_OPENAI_KEY` | `""` | API key (empty = fallback mode) |
| `AI_EMBEDDING_MODEL` | `text-embedding-3-small` | Embedding model name |
| `AI_ENABLED` | `true` | Master AI toggle |
| `FEATURE_AI_ENABLED` | (inherits `AI_ENABLED`) | Feature-flag override |

The `openAIClient` HTTP client has a **30-second timeout** (hardcoded in `client.go`). Server-level timeouts are configured in `internal/config/app.go`:
| Env Var | Default | Description |
|---------|---------|-------------|
| `APP_READ_TIMEOUT` | `30s` | Max duration for reading request body |
| `APP_WRITE_TIMEOUT` | `60s` | Max duration for writing response (AI endpoints with long-running LLM calls may need this) |

### Python ML Service (`moja-ml-service/app/config.py`)

See `moja-ml-service/app/config.py` for the full `Settings` model. Key settings:

| Setting | Default | Description |
|---------|---------|-------------|
| `storage_mode` | `shared_volume` | `shared_volume` or `presigned_url` |
| `whisper_model_size` | `medium` | Whisper model size |
| `whisper_compute_type` | `int8` | Quantization type |
| `go_backend_callback_url` | `http://localhost:8080/api/v1` | Go webhook endpoint |

---

## 7. Common Tasks

### Add a New LLM-Powered Feature

1. Define the prompt + JSON schema in `internal/modules/ai/client.go` (add method to `AIClient`, implement on `openAIClient`)
2. Add business logic in `service.go`
3. Store results in DB via `repository.go`
4. Expose via `handler.go` + `routes.go`
5. Call from frontend via `useAIApi()` in `routes.ts`

### Debug AI Playlist Generation

- Check `ai_generation_log` table for latency, track count, prompt
- Enable debug logging for the AI module
- The `GeneratePlaylist` method has a **4-tier fallback chain**: seed embedding → mood tag → mood range → plain `GetTracks()`. Check logs to see which tier is being used. If you see "mood range query failed, falling back to plain tracks" in logs, it means no mood/embedding data exists for any track.
- Verify track embeddings exist (`track_embeddings`) if using seed track similarity
- If all candidates come from `GetTracks()`, consider running the AI worker to generate embeddings/moods

### Run ML Service Locally

```bash
cd moja-ml-service
cp .env.example .env  # edit as needed
uv sync
uv run alembic upgrade head
uv run uvicorn app.main:app --reload --port 8000
# In another terminal:
uv run celery -A app.workers.celery_app worker -l info -Q lyrics_queue,covers_queue
```

### Add a New ML Pipeline

1. Create domain logic in `moja-ml-service/app/domain/<feature>/`
2. Add API endpoints in `moja-ml-service/app/api/v1/<feature>.py`
3. Add Celery task in `moja-ml-service/app/workers/tasks/<feature>_tasks.py`
4. Add callback delivery in `moja-ml-service/app/infra/callback/`
5. Add DB model in `moja-ml-service/app/infra/db/models.py`
6. Register in router, add config vars, run alembic migration

### Deploy ML Service

- **Docker:** Use `Dockerfile` + `docker-compose.yml` in `moja-ml-service/`
- **Storage mode:** For local/VPS use `shared_volume` (mount Go's uploads dir). For S3 use `presigned_url`
- **HMAC secret:** Must match between ML service and Go backend
- **Celery concurrency:** Set to `1` (RAM constraint for Whisper model)

### Regenerate Embeddings / Moods for All Tracks

The AI background worker (`internal/workers/ai/worker.go`) processes tracks without embeddings/moods every 10 minutes in batches of 20. To force a full re-index:

1. Truncate `track_embeddings` and/or `track_moods` tables
2. Restart the worker
3. It will reprocess all tracks in batches

---

## 8. Testing Patterns

### Go Backend AI Tests

- Unit test the `fallbackClient` directly — it's deterministic
- Test `service.go` methods with a mock `AIClient` and mock repository
- Test handler error paths by injecting different error types from the service

### Python ML Service Tests

Tests are in `moja-ml-service/tests/`. Test patterns:
- Unit test domain logic (transcriber, LRC formatter, confidence, optimizer) with small fixtures
- Mock the Whisper model for test speed
- Test HMAC callback signing end-to-end with known keys

### Frontend AI Tests

- Test Zod schema validation with invalid payloads
- Test page component rendering with mock API responses
