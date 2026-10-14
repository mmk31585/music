# Muse — Execution Roadmap

> **Territory:** phases · per-phase tests · per-phase CI · exit criteria · week-1 sprint.
> This is the *when* and *how-much* layer.
> - Map (why/what/architecture): [MASTER_PLAN.md](MASTER_PLAN.md)
> - Architecture (repos/contracts/migration): [ARCHITECTURE_REPOS.md](ARCHITECTURE_REPOS.md)
> - Culture/SRE/strategy: [ENGINEERING_OPERATING_SYSTEM.md](ENGINEERING_OPERATING_SYSTEM.md)

---

## 0. Rules of the road

1. **One sentence product:** *A place where someone can find and play Persian music, and it loads fast enough that they stay.* Every phase must serve this sentence.
2. **5 things that must work:** search works · click-to-play < 1s · artist page feels owned · playlists persist · Google indexes you (SSR).
3. **3 things we are NOT building (yet):** social rooms · gamification · subscriptions. Written down because they are tempting.
4. **Suspicion toward own ambition.** A new technology is justified only when you can name the *user-visible* feature it unlocks or the *measured* pain it removes. Deferred-infra triggers live in [ARCHITECTURE_REPOS.md](ARCHITECTURE_REPOS.md#7-deferred-infra-triggers).
5. **Ship broken early.** The first deploy happens in Phase 1, not when it's pretty. Deploying turns architecture questions into measured facts.

---

## 1. The phase map (high level)

| Phase | Name | Duration | Ships (user-visible) | New infra |
|---|---|---|---|---|
| 0 | Foundations | 1 wk | nothing yet — tooling + CI + dev env | Postgres, Redis, MinIO |
| 1 | Vertical slice MVP | 2–4 wk | search → play, end-to-end, on a VM | nginx, outbox (relay off) |
| 2 | Player & Library | 2–3 wk | player, queue, playlists, liked, history | partitioned history, outbox relay real |
| 3 | AI & Recommendations | 2–3 wk | lyrics, similar tracks, daily mix | FastAPI + Celery |
| 4 | Scale on evidence | trigger-driven | charts, HLS, better recs | Kafka, ClickHouse, Rust streamer |
| 5 | Ecosystem | with users | social, gamification, subscriptions | feature flags |

> Repo layout, runtime topology, and what "ships" means in repo terms: [ARCHITECTURE_REPOS.md](ARCHITECTURE_REPOS.md).
> **Exit criteria gate every phase** — a phase is not "done when code compiles"; it's done when its exit criteria pass (measured, not asserted).

---

## 2. Testing & CI standards (apply to every phase)

The pyramid that governs what kind of test goes where:

1. **Unit-test domain logic hard** — player engine, queue/shuffle, search ranking, Persian normalization. This is where bugs are expensive.
2. **Integration-test repositories against *real* Postgres** (testcontainers / CI service containers). Never mock the database — that's how phase-2 surprises happen.
3. **One Playwright e2e per golden path** — the phase's headline feature exercised from UI to data.

**CI gates** grow with the phases. The always-on set:
- `api` — `go vet` + `go test ./...` + `go build`
- `web` — `npm run lint` + `vue-tsc --build` + `vite build`
- `contracts` — OpenAPI drift check (a change to `libs/contracts` must be additive & versioned)
- `integration` — Go integration tests against Postgres/Redis service containers
- `e2e-smoke` — Playwright golden path

Coverage is a **floor, not a target**. More tests never replace the one golden-path e2e.

---

## 3. Phase 0 — Foundations (Week 1)

### Deliverables
- Monorepo skeleton (`apps/{web,api}`, `libs/contracts`, `infra/compose`, `docs/adr`).
- Dev environment: `make infra-up` → Postgres 17, Redis 7, MinIO on localhost.
- `Makefile` targets: `run`, `test`, `lint`, `fmt`, `build`, `infra-up/down`, `migrate-up`.
- Git hygiene: `.gitignore` (no binaries, no `.env`, no `uploads/`), commitlint, pre-commit hooks.
- **ADR-001:** target stack & phased roadmap. **ADR-002:** modular monolith before microservices.
- CI skeleton: every PR runs lint + test + build; every merge to `main` runs the full gate.

### Tests
- None for product code yet. The CI *pipeline itself* is the thing under test: it must fail on a seeded bad commit.

### CI
```yaml
# .github/workflows/ci.yml — runs on every PR
jobs:
  api:        # go vet, go test ./..., go build ./cmd/...
  web:        # npm run lint, vue-tsc --build, vite build
  contracts:  # openapi diff — additive & versioned
```
Postgres/Redis available to integration tests via `service:` containers.

### Exit criteria
- Fresh clone + `make infra-up` + `make run` on a clean machine works.
- CI green on a trivial PR, red on a seeded failure.
- ADR-001/002 written and linked from README.

---

## 4. Phase 1 — Vertical slice MVP (Weeks 2–4)

**Goal: one feature through the whole stack.** `search → open track → it plays`, deployed to one VM. Nothing else is real yet.

### Deliverables
- **Go API module skeleton**: `internal/app` (composition root) + `health`, `auth`, `catalog`, `search`, `upload`, `media`. Each module = `handler/service/repository` + `routes.go`.
- **Contracts**: OpenAPI for `GET /tracks/:id`, `GET /search?q=`, `POST /auth/register|login`, `POST /uploads/audio`.
- **Auth**: register/login/refresh, JWT access + refresh, **httpOnly cookie** for refresh, RBAC from current repo's `000060`.
- **Catalog read path**: `tracks/albums/artists/genres` many-to-many (current `000013`). Search via **Postgres FTS + `pg_trgm`** (not OpenSearch).
- **Upload + storage**: `STORAGE_DRIVER=local|s3`; Phase-1 default = local filesystem + MinIO (Iranian DC: no provider S3 dependency). nginx serves `uploads/` with range + TLS.
- **Nuxt 3 SSR app**: public SEO pages (`/artist/:id`, `/album/:id`, `/track/:id`) server-rendered with full meta; app shell with session restore; search page; track page wired to a real `<audio>`.
- **event_outbox** (pattern from current `000044`) — created, relay logs "would publish". Future-proofing that costs ~nothing now.
- **Deploy**: Docker Compose to **one VM** in the Iranian DC. nginx in front (TLS via CDN or ACME). `/healthz` + `/metrics` exposed.

### Tests
- **Unit (Go):** auth token issue/refresh/expiry; search ranking + Persian normalization (reuse current repo's `persian.go` logic + tests); storage driver path mapping.
- **Unit (frontend):** auth composable (session restore, refresh-on-401); search store.
- **Integration (Go, real Postgres):** auth register→login→refresh; catalog read path with seeded rows. **Repository layer never mocked.**
- **Contract tests:** request/response shapes against the OpenAPI spec.

### CI
- Phase 0 jobs + `integration` (Postgres+Redis service containers) + `typecheck` (vue-tsc) + `e2e-smoke` (Playwright: register → search → play).

### Exit criteria
- Playwright golden path passes: register, search a Persian term, open track, audio plays.
- SEO check: `curl` an artist page returns server-rendered HTML with meta (proves SSR is real).
- Click-to-play on a VM over HTTPS < 1s (measured, not assumed).
- `make deploy` to the VM is a repeatable documented runbook.

---

## 5. Phase 2 — Player & Library (Weeks 5–7)

**Goal: playback becomes a product.** Queue, shuffle/repeat, playlists, liked tracks, listening history — and the event path that will later feed analytics.

### Deliverables
- **Frontend player engine** (deep module, one seam): audio + queue + preload + Media Session + shuffle/repeat + history emit. Unit-testable by substituting the audio backend (reuse the design from the current repo).
- **Playlists CRUD** + `playlist_tracks` (current `000004`/`000040`), **liked tracks** (`000050`), **library** view.
- **Listening history** — time-partitioned (current `000043`) — with a `track_played` event written **through the outbox**, atomically with the DB write.
- **Outbox relay becomes real**: a single Go worker reads `pending`, publishes to `events_log` (still no Kafka — Kafka is a sink swap later).
- **Now-playing bar** + Media Session (mobile + desktop).

### Tests
- **Unit (frontend):** player engine (play/pause/seek, next/prev, shuffle modes, repeat modes, queue mutation) — cover hard, bugs are expensive here.
- **Unit (Go):** playlist service (ownership, collaborators, reorder); history service; outbox relay (idempotency, retry, ordering).
- **Integration:** playlist CRUD vs real Postgres; outbox → relay → `events_log` round trip.
- **E2E:** create playlist, add tracks, play from it, history appears.

### CI
- Frontend unit (vitest) with a coverage gate on `player/`; Go integration job extended; e2e extended to the playlist golden path.

### Exit criteria
- Playback survives: reload mid-track, tab-away (Media Session), low network (buffering UI).
- A played track appears in history *and* in `events_log` within seconds (event path end-to-end).
- No listener-facing state lost on refresh.

---

## 6. Phase 3 — AI & Recommendations (Weeks 8–10)

**Goal: the AI seam from Phase 1 gets its real implementation.** FastAPI behind the HMAC seam; recommendations start simple and honest.

### Deliverables
- **Python/FastAPI service** (`apps/ai`) + Celery worker:
  - `transcribe-lyrics` (Whisper) — Persian normalization + LRC formatting (reuse `moja-ml-service/app/domain/lyrics/` design).
  - `embed-track` / `similar-tracks` — embeddings in **separate tables** (keep current `000045` idea; embeddings never live in main tables).
- **Recommendation v1 = popularity + recency + genre affinity** (a SQL query + Redis cache, **not ML yet**). Daily-mix + popular on Home. ML collaborative filtering is Phase 4 — only with real listening data.
- **HMAC internal routes** on Go (`/ml/*`), like the current repo's `MLServiceHMACMW`.
- Lyrics UI: synced LRC viewer, karaoke highlight, plain-text fallback.

### Tests
- **Python (pytest):** Persian normalizer, LRC formatter, confidence scoring, embedding pipeline (deterministic fixtures; no model download in CI).
- **Contract tests (Go ↔ Python):** `httptest` server acting as the AI service; Go client against the OpenAPI contract — proves the seam without running torch in CI.
- **Integration:** upload track → task queued → transcribed → lyrics appear in API.
- **E2E:** track page shows synced lyrics; Home shows a daily mix.

### CI
- New `ai` job: `uv run pytest` (or pip), lint (ruff), typecheck.
- New `contract` job: spin Go stub + Python stub against `libs/contracts`, assert both conform.
- Docker: Go API + `ai` images built; Celery worker image.

### Exit criteria
- Full lyrics round-trip works from upload to synced display.
- Home loads with a daily mix in < 300 ms (cached).
- `similar tracks` returns sensible results for 5 hand-picked Persian tracks.

---

## 7. Phase 4 — Scale on Evidence (only when metrics justify)

**Nothing in this phase happens on a schedule.** Each item has a trigger. If a trigger doesn't fire, we don't build it — that is a feature, not a failure.

| Item | Trigger (measured) | What lands |
|---|---|---|
| Kafka | outbox `events_log` grows + a second consumer needs events | outbox relay → Kafka topic; consumers subscribe |
| ClickHouse | listen events make Postgres analytics slow + a user-facing chart/leaderboard needs it | events → ClickHouse, rollups, dashboards |
| Rust streamer | nginx range-serving is the bottleneck (huge catalog, HLS/DASH, high concurrency) | media edge behind CDN |
| Search upgrade | Postgres FTS measurably underperforms | OpenSearch or pgvector |
| Recs v2 | enough listening history for collaborative filtering to beat popularity | train offline on ClickHouse aggregates → mixes into Redis |

> **The canonical trigger table lives in [ARCHITECTURE_REPOS.md](ARCHITECTURE_REPOS.md#7-deferred-infra-triggers).** This copy is here for execution context only.

### Tests
- Kafka consumers: at-least-once + idempotent handling (the outbox guarantees the producer side; the consumer must dedupe).
- ClickHouse ingest: schema + rollup idempotency tests.
- Rust: range-request unit + integration tests against a real S3/MinIO.

---

## 8. Phase 5 — Ecosystem (with users & revenue)

**Goal: copy the *designed* feature set from the current repo — social rooms, clubs, gamification (XP/badges/challenges), creator payouts, subscriptions — and build it only for an audience that exists.**

Rules:
- Reuse current repo migration designs: `000020` (social), `000019` (gamification), `000022` (payments).
- WebSockets for rooms/clubs via a single hub (current `internal/platform/ws` pattern).
- Payments in Iran = local-gateway abstraction (current `internal/platform/payment` design).
- Each feature ships with its own unit + integration tests + e2e golden path, behind a **feature flag**.

---

## 9. Definition of Done (apply to every phase)

- [ ] Golden-path e2e passing (Playwright) for the phase's headline feature.
- [ ] Unit + integration tests committed and green in CI.
- [ ] Type checks pass (Go vet, vue-tsc, ruff).
- [ ] `make deploy` runbook executed once on the VM; health endpoint green.
- [ ] Observability wired: the phase's key action produces a log line + a metric + a trace.
- [ ] Exit criteria for the phase (above) actually met — measured, not asserted.
- [ ] If a hard-to-reverse decision was made → an ADR exists.

> Ownership, SLOs, and "is it Landed not just Launched" are judged against [ENGINEERING_OPERATING_SYSTEM.md](ENGINEERING_OPERATING_SYSTEM.md).

---

## 10. First sprint (Week 1)

```text
Day 1  mkdir the monorepo; git init; .gitignore; commitlint; Makefile skeleton; infra/compose (postgres, redis, minio); make infra-up works.
Day 2  Go module apps/api; internal/app composition root; health module; /healthz + /metrics; CI job `api` green (vet+test+build).
Day 3  migrations/ + outbox table (000044 pattern) + a seed migration for genres; migrate-up wired to Makefile.
Day 4  OpenAPI skeleton in libs/contracts (auth + tracks + search); a drift-check job that fails CI if the spec and code diverge.
Day 5  Nuxt app apps/web; dev proxy to API; one SSR page (/artist/:id) rendering from the API; CI job `web` green.
Day 6  ADR-001 (stack+phases+deferred triggers), ADR-002 (monolith); README with run instructions.
Day 7  Demo: fresh clone → make infra-up → make run → the SSR page renders. CI green. VM provisioning notes written.
```

> If Week 1 completes, the seam structure is real, CI is real, and Phase 1's slice is one vertical feature away. That is the whole plan working.
