# Muse — Backend API (AGENTS.md)

> Start every session by reading this file, then `docs/agent/STATE.md`, then the last 20 lines of `docs/agent/BACKLOG.md` (or `docs/LEARNINGS.md` if present). Run `git status` and `make test` before editing.

## 1. What this is

Muse is a full-stack Persian music streaming platform. This repository is the **Go/Gin backend API** (module `music`) plus the Vue 3 SPA under `frontend/` (git submodule) and the Python ML microservice under `moja-ml-service/`.

- **Goal right now:** it works and is deployed. Polish comes later; do not add features outside an active phase.
- **Audience:** Persian + English speakers (Iranian market). Persian text must keep working everywhere: search, sorting, display.
- **Closed beta:** registration gated by `INVITE_CODE` where configured; track uploads are user-owned content (no scraping, no public catalog).
- **Non-goals for now:** payments not in scope for backend core polish, social/realtime already implemented in modules but not expanded, mobile apps out of scope.

## 2. Repo layout (monorepo shape)

```
music/
├── cmd/api/                 # API server entrypoint (main.go)
├── cmd/worker/              # Background worker entrypoint
├── internal/                # All business logic (not importable outside module)
│   ├── app/                 # bootstrap, DI container, routes, lifecycle
│   ├── common/              # middleware, errors, pagination, response, validation
│   ├── config/              # env/config loading (godotenv + typed structs)
│   ├── modules/             # ~30 domain modules (catalog, auth, playlist, ...)
│   ├── pkg/audioinfo/       # audio metadata helpers
│   ├── platform/            # db, redis, opensearch, storage, ws, events, etc.
│   └── workers/             # async workers (ai, indexer, notification, ...)
├── migrations/              # Goose SQL migrations
├── frontend/                # Vue 3 SPA (git submodule; separate AGENTS.md)
├── moja-ml-service/         # FastAPI + Celery + Whisper
├── scripts/                 # migrate.sh, dev helpers
├── docs/                    # architecture, ADRs-ish notes, agent state
├── Makefile                 # canonical commands
├── Dockerfile               # multi-stage API build
├── fly.toml / railway.json / vercel.json  # deploy targets
└── AGENTS.md                # this file
```

Anything under `frontend/` is owned by the frontend repo; do not edit it from here unless the task explicitly crosses the seam (keep the contract stable).

## 3. Stack

| Concern | Choice |
|---|---|
| Language | Go 1.25.9 |
| HTTP framework | Gin 1.12 |
| DB | PostgreSQL 17 (pgx + sqlx) |
| Migrations | Goose |
| Auth | JWT HS256 (golang-jwt v5), scrypt-style password hashing per module |
| Cache/PubSub | Redis 8 (go-redis v9) |
| Search | OpenSearch 3.2 (Go client) |
| Storage | Local or S3 (MinIO/AWS SDK v2), switchable via config |
| Events | In-process bus + transactional outbox |
| WebSocket | gorilla/websocket |
| Logging | Zap (structured, global singleton) |
| Metrics | Prometheus `/metrics` |
| Docs | Swaggo Swagger at `/swagger/*` |
| ML | Python FastAPI + Celery + Whisper (separate service) |
| Frontend | Vue 3 + Vite + Pinia + Tailwind v4 + PrimeVue |

Do not swap a stack choice without a decision note in `docs/`.

## 4. Entry points and boot

- `cmd/api/main.go`: signal-aware server bootstrap (config, tracing, Sentry, routes, workers, graceful shutdown). Also contains Swagger glue.
- `cmd/worker/main.go`: standalone worker binary.
- `internal/app/bootstrap.go`: assembles `App` (config, DB, Redis, logger, etc.).
- `internal/app/container.go`: manual DI (no framework); ~25 `build*` methods + `subscribeEvents()`.
- `internal/app/routes.go`: all route registration grouped by auth level.

Rules:
- Routes are registered **once** in `internal/app/routes.go`. Do not register in `main.go`.
- `main.go` is for process lifecycle, not business wiring.

## 5. Key patterns

### Module shape
Every domain module should follow:

```
module/
├── model.go
├── repository.go
├── service.go
├── handler.go
├── dto.go
├── routes.go
└── *_test.go
```

Use PascalCase exports, `-er` interfaces (`Reader`, `StorageProvider`), idiomatic Go.

### Manual DI
`NewContainer()` wires concrete implementations and subscribes cross-module events. Add new modules by adding a `buildXxx(a)` call and (if needed) event subscriptions.

### Events
Modules publish/subscribe through `events.Bus`; critical events use the transactional outbox. Never import another module's repository directly when an event fits.

### Response / errors
Standardized via `internal/common/response` and `internal/common/errors`. Keep the JSON envelope consistent.

### Config
`internal/config/*` loads from env (godotenv in dev). Validate at startup; fail fast with a readable error.

## 6. Commands (use these, not ad-hoc)

```bash
make dev            # infra-up + go run ./cmd/api
make dev-worker     # infra-up + go run ./cmd/worker
make test           # go test ./... -race -count=1
make test-coverage  # coverage profile
make lint-go        # golangci-lint
make fmt            # go fmt ./...
make build          # go build -o bin/api ./cmd/api
make build-all      # api + worker
make migrate-up     # applies goose migrations
make migrate-status
make infra-up / infra-down
make docker-up      # compose build + run
```

Frontend commands live in `frontend/` (`npm run lint`, `npx vitest run`, `npx vue-tsc --noEmit`).

CI gate: `go test ./... -race -count=1` + lint + build; coverage threshold currently 30%. Fix the first error only, then rerun.

## 7. Database and migrations

- Goose SQL migrations in `migrations/`.
- Run via `scripts/migrate.sh {up|down|status}` (uses `DATABASE_URL`/unpooled URL).
- Never run migrations inside an HTTP request.
- Never edit a migration that has already shipped; add a new one.

## 8. Environment

See `.env.example` (do not commit `.env`). Key vars:

- `APP_ENV`, `APP_HOST`, `APP_PORT`, `PORT` (Vercel injects `PORT`; config falls back to it)
- `DATABASE_SOURCE` / `POSTGRES_*`, `REDIS_ADDR`, `REDIS_PASSWORD`
- `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`
- `CORS_ALLOWED_ORIGINS` (exact origins, no wildcard)
- payment keys (`IDPAY_*`, `ZARINPAL_*`), ML/external service keys, storage keys
- tracing/Sentry (`SENTRY_DSN`, OTLP endpoint)

`internal/config/env`/`app.go` etc. own validation; add new fields there, not scattered.

## 9. Persians and search

`internal/.../search` and normalization helpers should follow the Persian rules: NFKC, Arabic→Persian letters, strip ZWNJ/tatweel/diacritics, digit normalization, lowercase, whitespace collapse. Keep Persian-aware search/sorting in tests.

## 10. Working agreement

For every session:

1. Read this file, `docs/agent/STATE.md`, and scan `docs/agent/BACKLOG.md`.
2. Run `git status` + `make test`.
3. **One task at a time** (one bullet of one phase). Keep commits small (< ~200 changed lines) with Conventional Commit messages.
4. After each task: run `make test` (and `go vet ./...`, `golangci-lint` when relevant). Commit only when green, then update `docs/agent/STATE.md`.
5. On 3 failed attempts on the same problem, record it in `docs/agent/BACKLOG.md`/concerns and move on.
6. Record non-trivial bugs/surprises in `docs/LEARNINGS.md` (if present) append-only: `date | symptom | cause | fix | rule`.
7. Never commit secrets, real URLs/tokens in tests, disable failing tests, or add deps without noting why.
8. Validate inputs with the validator package; use parameterized SQL; avoid SQL injection; include security headers.
9. Do not work on a later concern before earlier blockers are cleared.

Tests: table-driven with `testify`; DB tests via mocks/sqlmock or real Postgres in Docker. Race detector enabled in CI.

## 11. State, docs map

- `docs/agent/STATE.md` — phase completion, build state, current blockers (last updated 2026-07-09).
- `docs/agent/BACKLOG.md` — open issues/tasks.
- `docs/agent/ux/` — UX journey specs and friction log.
- `docs/codebase/*` — generated stack/structure/conventions/concerns/testing.
- `docs/architecture/*` — aspirational microservices, DB, ops, services docs.
- `docs/MASTER_PLAN.md`, `docs/ROADMAP.md`, `docs/ARCHITECTURE_REPOS.md`, `docs/ENGINEERING_OPERATING_SYSTEM.md` — product/strategy spine.
- `docs/swagger.json` / `swagger.yaml` — API contract.
- Frontend has its own `frontend/AGENTS.md` and `frontend/PERFORMANCE.md`.

Current state highlights from `docs/agent/STATE.md`:
- Go build/test/vet reported clean; Vue `npm run build` failed earlier due to environment/binary format (not code) — re-verify before release.
- 12/12 planned phases marked complete; remaining items are MEDIUM/LOW concerns.
- Known concerns live in `docs/codebase/CONCERNS.md` (frontend type-check, auth refresh contract, admin CRUD gaps, etc.).

## 12. Agent skills available in this environment

When working in this repo, prefer these installed skills:
- `muse` — project-wide Muse conventions (RTL, Dastgah, visual system).
- `muse-auth`, `muse-catalog`, `muse-player`, `muse-social`, `muse-creator`, `muse-admin`, `muse-ui`, `muse-infra` — domain specialists by area.
- `golang-*` — idiomatic Go (concurrency, errors, testing, security, etc.).
- `vue`, `nuxt-ui` (if frontend work), `playwright` for E2E.
- `code-review-and-quality`, `security-and-hardening`, `test-driven-development`.

## 13. Human-only tasks

Stop and ask the human for:
- Secret rotation / credential management.
- Running migrations against production Neon.
- Vercel/Fly/Railway project linking and dashboard configuration.
- Production deploy sign-off and rollback decisions.
- Any change that breaks the frontend contract (`docs/swagger.json`).

## 14. Safety notes

- Do not log tokens, passwords, or PII.
- CORS origins must be exact; never wildcard in production.
- Audio bytes never stream through this API in the intended model (local files + MinIO); keep upload/download policies consistent with configured storage.
- Vercel deployment note: Vercel Go server must listen on `PORT` when provided (config falls back to it); provided `APP_PORT` can override in dev.
- External services can be down; ML/lyrics calls must not break core track endpoints.

---

Maintenance: If a subsequent session discovers this file is stale or wrong (build commands changed, new env vars, new module), edit it in the same commit as the change. Keep this file under ~400 lines.