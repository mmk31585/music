# Symphony Architecture v2 — Incremental Migration Plan

## Key Decision: Migrate, Don't Restart

**Challenge:** You asked to "start from scratch," but the repo already contains 30 migration files and 275+ Go files representing real business logic, schema evolution, and working features. Restarting would discard all of that.

**Recommended path:** Incrementally restructure the existing codebase into the Symphony vertical-module DDD layout while keeping the app running at every step. This is the only path that preserves existing data and velocity.

---

## Current State
- **Layout:** Flat per-module (`internal/modules/<name>/handler.go`, `service.go`, `repo.go`, `model.go`, `routes.go`)
- **Transports:** REST only (Gin)
- **Database:** PostgreSQL + Redis + MinIO + OpenSearch
- **Shared libs:** `internal/platform/` (database, events, cache, storage, tracing, etc.)
- **Entrypoints:** `cmd/api/main.go`, `cmd/worker/main.go`

## Target State
- **Layout:** Vertical modules with DDD layers (`api/`, `application/`, `domain/`, `infrastructure/`, `events/`, `tests/`, `migrations/`)
- **Transports:** REST + gRPC + GraphQL + WebSocket + SSE via Gateway → BFF
- **Module contract:** Each module exposes a clean interface; extraction is a file copy

---

## Phases (Few Steps, One at a Time)

### Phase 0: Foundation (Weeks 1-3)
**Goal:** New directory layout + platform layer + gRPC bootstrap without breaking existing REST.

**Step 0.1 — Scaffold empty target directories**
- Create `modules/` with all 28 module subdirectories (empty)
- Create `platform/`, `proto/`, `bff/`, `gateway/`, `deployments/`, `docs/adr/`
- Add `.gitignore` entries for new build artifacts
- Validation: `go build ./...` still succeeds

**Step 0.2 — Promote `internal/platform/` to `platform/`**
- Move shared libs (database, redis, events, storage, logger, metrics, tracing, cache, opensearch, web, ws, audio, payment, sentry) to `platform/`
- Update all imports
- Validation: `go build ./...` + existing tests pass

**Step 0.3 — Add gRPC server skeleton**
- `cmd/server/main.go` runs both REST (Gin) and gRPC on separate ports
- Add health checks (`/healthz`, `/readyz`) and OTel middleware
- Validation: Both servers start, existing REST endpoints unchanged

**Step 0.4 — Establish module interface contract**
- Create `platform/module.go` interface: `RegisterGRPC`, `RegisterREST`, `Migrations()`, `Events()`
- Validation: Compiles, no behavior change

---

### Phase 1: Identity Module (Weeks 4-6)
**Goal:** First fully-structured module; proof of concept for vertical slice.

**Step 1.1 — Define Identity domain and Proto**
- `modules/identity/domain/` (user, session, mfa, value objects)
- `proto/identity/*.proto` (UserService, SessionService, AuthService)
- Generate Go stubs
- Validation: Proto compiles, domain tests pass

**Step 1.2 — Implement Identity infrastructure**
- `modules/identity/infrastructure/persistence/` (PostgreSQL repo)
- `modules/identity/infrastructure/auth/` (JWT, password hashing)
- `modules/identity/migrations/001_create_users.sql`, `002_create_sessions.sql`
- Validation: Migrations run against test DB

**Step 1.3 — Implement Identity application + API**
- `modules/identity/application/` (use cases: RegisterUser, Login, VerifyToken)
- `modules/identity/api/grpc/` and `api/rest/` handlers
- Validation: gRPC and REST endpoints respond correctly

**Step 1.4 — Wire Identity into server**
- Register Identity module in `cmd/server/main.go`
- Keep old `internal/modules/auth` running in parallel (feature flag)
- Validation: Old auth routes still work; new Identity module also works

---

### Phase 2: Catalog Module (Weeks 7-10)
**Goal:** Second module; validate cross-module communication via events.

**Step 2.1 — Define Catalog domain and Proto**
- `modules/catalog/domain/` (track, album, artist, genre, label, ISRC, UPC)
- `proto/catalog/*.proto`
- Validation: Domain invariants tested

**Step 2.2 — Catalog infrastructure + migrations**
- Persistence, OpenSearch indexing, enrichment pipelines
- Module-local migrations
- Validation: Full catalog schema migrates cleanly

**Step 2.3 — Catalog application + API**
- Use cases (CreateTrack, SearchTracks, GetArtist)
- gRPC + REST handlers
- Validation: End-to-end catalog CRUD works

**Step 2.4 — Event contracts**
- `track.created`, `album.updated`, `artist.verified` events
- Identity publishes `user.registered`; Catalog subscribes (noop for now)
- Validation: Event bus round-trips

---

### Phase 3: Streaming + Player (Weeks 11-14)
**Goal:** Core playback loop in modular form.

**Step 3.1 — Streaming module**
- `modules/streaming/domain/` (StreamURL, HLSManifest, Format, Bitrate)
- Infrastructure: MinIO presigned URLs, transcoder integration
- Validation: Presigned URLs generated and expire correctly

**Step 3.2 — Player module (new aggregate)**
- `modules/player/domain/` (Player, Queue, PlaybackSession, Crossfade, RepeatShuffle)
- Application: StartPlayback, Pause, Seek, Next, Previous
- Validation: State transitions are valid

**Step 3.3 — Session module**
- `modules/session/domain/` (Session state machine: playing/paused/stopped)
- Redis-backed session state
- Validation: Session survives process restart

**Step 3.4 — Wire playback flow**
- Identity → Catalog → Streaming → Player → Session
- WebSocket for real-time position updates
- Validation: End-to-end playback works in test client

---

### Phase 4: Curation + Discovery (Weeks 15-20)
**Goal:** Playlists, library, search, recommendations.

**Step 4.1 — Playlist + Library modules**
- Vertical slices for playlist and user library
- Validation: CRUD + save/unsave tracks

**Step 4.2 — Search module**
- `modules/search/domain/` (SearchQuery, SearchResult, RankingPipeline)
- Infrastructure: OpenSearch full-text + pgvector semantic + ranker
- Validation: Hybrid search returns relevant results

**Step 4.3 — Recommendations module**
- Plugin-based engine architecture
- First engine: trending
- Validation: Recommendations endpoint returns tracks

---

### Phase 5: Social + Content Types (Weeks 21-26)
**Goal:** Social features and expanded content types.

**Step 5.1 — Social + Community modules**
- Follow, share, comment, feed, listening parties/rooms
- Validation: Social graph operations work

**Step 5.2 — Lyrics + Podcast + Audiobook modules**
- Each as a vertical slice
- Validation: Lyrics fetch, podcast RSS ingestion, audiobook progress tracking

**Step 5.3 — Download + Offline module**
- Download queue, offline storage, license cache, sync state
- Validation: Download + play offline flow

---

### Phase 6: Platform Services (Weeks 27-32)
**Goal:** Monetization, AI, creator tools.

**Step 6.1 — Billing + Notification + Moderation + Analytics**
- Stripe subscriptions, push notifications, content moderation, ClickHouse analytics
- Validation: Billing webhooks, notification delivery

**Step 6.2 — Creator Studio module**
- Upload, schedule releases, manage royalties, audience insights
- Validation: Release pipeline end-to-end

**Step 6.3 — AI Platform module**
- LLM provider abstraction, prompt templates, first agent (playlist generator)
- Validation: AI playlist generation returns coherent track list

**Step 6.4 — Plugin Manager module**
- Hook registry, WASM runtime, built-in plugins
- Validation: Plugin install/enable/disable lifecycle

---

### Phase 7: Gateway + BFF + Observability (Weeks 33-36)
**Goal:** Multi-protocol gateway, client-specific BFFs, full observability.

**Step 7.1 — API Gateway**
- Envoy config: REST → BFF, gRPC-web → BFF, GraphQL federation, WebSocket upgrade, SSE streaming
- Validation: All protocols reachable via gateway

**Step 7.2 — BFF Layer**
- Web, Mobile, TV, Desktop BFFs
- GraphQL schema for web/mobile
- Validation: Each client gets optimized payloads

**Step 7.3 — Observability maturity**
- Tempo + Loki + Mimir + Grafana + Pyroscope
- SLO dashboards, error budgets, service map
- Validation: Traces, metrics, logs visible in Grafana

---

### Phase 8: Extraction Readiness (Weeks 37-40)
**Goal:** Any module can be extracted to an independent service.

**Step 8.1 — Service mesh readiness**
- Linkerd or similar injected; all modules have `/healthz`, `/readyz`, `/metrics`
- Validation: Service map auto-discovers dependencies

**Step 8.2 — Extract first service candidate**
- Copy `modules/streaming` to `services/streaming` (file copy, not refactor)
- Run as separate pod with its own DB schema
- Validation: Streaming works via gRPC across service boundary

**Step 8.3 — Document extraction runbook**
- ADR for each extracted module
- Validation: New team member can extract a module in <1 day

---

## Validation Principles at Every Step
1. `go build ./...` must pass
2. Existing tests must pass (or be marked skip with ticket)
3. Old and new code run in parallel during migration (feature flags)
4. No migration step may break database schema (additive migrations only)
5. Each module must have domain tests before infrastructure tests

---

## Risks
| Risk | Mitigation |
|------|-----------|
| Scope creep — too many modules at once | One module per phase, parallel old/new via feature flag |
| Import cycle after restructuring | Enforce `platform/` only dependency direction |
| gRPC + REST duplication | Use `grpc-gateway` to generate REST from Proto |
| Event schema drift | Centralized event registry + contract tests |
| Performance regression | Benchmark each migrated module against old baseline |

---

## Open Questions for You
1. Do you want to preserve the existing `internal/` layout during migration, or rename to top-level `modules/` + `platform/` immediately?
2. Do you want Proto-first (define `.proto` before domain) or domain-first (Go types before Proto)?
3. Should the worker (`cmd/worker/main.go`) be merged into the monolith or kept separate?
