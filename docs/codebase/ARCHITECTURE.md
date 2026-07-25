# Muse — Architecture

## High-Level System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        CDN (CloudFront / ArvanCloud)         │
└──────────────────┬──────────────────────┬───────────────────┘
                   │                      │
                   ▼                      ▼
┌──────────────────────────────┐ ┌────────────────────────────┐
│     Nginx (Frontend SPA)     │ │  Go/Gin API Server (:8080) │
│  Vue 3 App served on :80     │ │  JWT-auth, rate-limited     │
└──────────────────────────────┘ └──────────┬─────────────────┘
                   │                         │
                   ▼                         ▼
       ┌──────────────────┐       ┌──────────────────────┐
       │   PostgreSQL 17  │       │   Redis 8.8           │
       │   (Primary DB)   │       │   (Cache + Session)   │
       └──────────────────┘       └──────────────────────┘
                   │                         │
                   ▼                         ▼
       ┌──────────────────┐       ┌──────────────────────┐
       │  OpenSearch 3.2  │       │   MinIO (S3 Storage)  │
       │  (Full-text)     │       │   (Media files)       │
       └──────────────────┘       └──────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│               Python ML Service (:8000)                      │
│  FastAPI + Celery + Whisper + OpenRouter + OpenAI           │
└─────────────────────────────────────────────────────────────┘
```

## Current State: Go Monolith

Despite ambitious microservices docs (`docs/architecture/SERVICES.md`), the **current codebase is a Go monolith** with:

- **1 API binary** (`cmd/api`) — serves all HTTP endpoints
- **1 Worker binary** (`cmd/worker`) — background jobs
- **30 domain modules** in a single `internal/modules/` tree
- **In-process event bus** for decoupling (not Kafka/RabbitMQ)
- **Transactional outbox** pattern for at-least-once delivery

## Key Patterns

### 1. Manual DI Container (`internal/app/container.go`)

No DI framework (no wire, no dig). All dependencies wired manually in `NewContainer()`:

```go
func NewContainer(a *App) *Container {
    c := &Container{
        SQLX:  sqlx.NewDb(a.DB, "pgx"),
        Bus:   events.NewBus(),
        RDB:   a.Redis,
        WSHub: ws.NewHub(),       // starts goroutine
    }
    c.buildAuth(a)
    c.buildCatalog(a)
    c.buildPlaylist(a)
    // ... ~25 more build calls
    c.subscribeEvents()           // cross-module event wiring
    return c
}
```

### 2. Layered Module Pattern

Every module follows:
```
model.go → repository.go → service.go → handler.go → routes.go
```

- **model.go**: DB structs, domain types
- **repository.go**: SQL queries (sqlx)
- **service.go**: Business logic
- **handler.go**: Gin HTTP handlers
- **routes.go**: `RegisterRoutes(router, handler, mw...)`

### 3. Event-Driven Decoupling

Modules communicate via an in-process event bus with `events.Bus`:

```go
// Publisher
bus.Publish(ctx, events.EventTrackPlayed{TrackID: id, UserID: uid})

// Subscriber (wired in container)
bus.Subscribe(events.EventTrackPlayed{}, c.HistoryService.HandleTrackPlayed)
```

Transactional outbox ensures at-least-once delivery for critical events.

### 4. Route Organization (`internal/app/routes.go`)

All routes under `/api/v1`:

| Group | Auth | Modules |
|-------|------|---------|
| Public | None | health, player (public), analytics |
| Authenticated | JWT | auth, lyrics, playlist, library, queue, follow, recommendation, history, notification, moderation, contribution, subscription, reactions, gamification, creator, social, ai, video, importcmd |
| Optional Auth | `OptionalAuthMW` | catalog, search |
| Admin | JWT + admin check | media, ingestion, dashboard, importcmd, catalog admin |
| Internal | HMAC-signed | lyrics callbacks, covers callbacks, video callbacks |
| Special | None | Swagger, Prometheus `/metrics`, WebSocket `/ws` |

### 5. Frontend Architecture Patterns

| Pattern | Description |
|---------|-------------|
| **Dynamic Layout Switching** | `App.vue` reads `route.meta.layout` → renders LayoutAuth, LayoutMusicApp, LayoutAdmin, or LayoutEmpty |
| **Middleware Chain Routing** | Guards composed as `checkMaintenanceGuard → checkLoginGuard → checkAuthGuard` |
| **API Barrel Pattern** | 25 domain API modules, each exporting `useXxxApi()` composable factory |
| **Request Wrapper** | `useRequest()` wraps Axios with token injection, 401 refresh queue, Zod validation, toast hooks |
| **Store Pattern** | Pinia setup stores (`defineStore('name', () => { ... })`) |
| **Token Refresh Queue** | 401 responses queue requests during refresh, replay after success |

### 6. Auth Flow

```
Login → JWT Access Token (60min) + Refresh Token (30 days)
  ├── Access: stored in Pinia + localStorage (encrypted)
  ├── Refresh: stored in httpOnly cookie
  ├── 401 → interceptor queues request → refreshes token → replays
  └── Refresh fails → force logout → redirect to login
```

## Performance Targets

| Metric | Target |
|--------|--------|
| FCP | < 1.5s |
| LCP | < 2.5s |
| TTI | < 3.5s |
| Track playback | < 500ms |
| Search-to-result | < 300ms |
| Page transitions | < 100ms |
| API P99 | < 200ms |
| Concurrent streams | 1M+ |

## Scaling Strategy

Current: Single Go binary + PostgreSQL + Redis.
Aspirational (from docs): 14 microservices on Kubernetes with HPA, Kafka event bus, CQRS, read replicas, CDN edge caching.

The architecture docs describe a future state targeting 50M users, 10M tracks, 1B streams/month.
