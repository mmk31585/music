---
name: muse
description: Use when working on the Muse Persian music ecosystem — a full-stack music streaming platform with Vue 3 frontend and Go/Gin backend. Includes Persian identity (RTL, Dastgah modal system, Jalali calendar), glassmorphism design system, audio engine, social features, creator economy, and gamification. Trigger on files matching frontend/**, internal/**, cmd/**, docs/**, Makefile, Dockerfile, docker-compose*.yml, migrations/*.sql.
---

# Muse — Persian Music Ecosystem

**Product:** Muse — Premium Persian Music Ecosystem
**Tagline:** "Feel the music. Share the moment."
**Target:** 50M users, 10M tracks, 1B streams/month across Iran and MENA

---

## Core Pillars

- **Streaming:** Premium audio (320kbps AAC / FLAC), gapless playback, crossfade, 8 player modes
- **Discovery:** AI-driven recommendations, mood/genre exploration, viral charts
- **Social:** Activity feeds, listening parties, live rooms, music clubs
- **Community:** Wiki-style contributions (metadata, lyrics, translations), multi-tier moderation
- **Creator Economy:** Artist profiles, release scheduling, analytics, fan tipping
- **Gamification:** XP system, badges, leaderboards, daily challenges
- **Persian Identity:** RTL support, Dastgah modal system, Jalali calendar, IRANYekanWeb font

---

## Tech Stack

### Frontend
Vue 3 (Composition API + `<script setup>`) · TypeScript strict · Vite 7 · Pinia 3 · PrimeVue 4 (auto-imported, RTL enabled) · Tailwind CSS v4 + tailwind-merge + clsx · PrimeIcons + Lucide Vue Next · Vue Router 4 · Axios (JWT refresh interceptor) · HTML5 Audio + Web Audio API (AnalyserNode) · Zod 4 · @vueuse/core

### Backend
Go 1.25 · Gin 1.12 · PostgreSQL 16 (pgx + sqlx) · Redis 7 (go-redis) · OpenSearch 2.x · MinIO / S3 · JWT (golang-jwt v5) · goose 3 · Zap · gorilla/websocket · dhowden/tag + ffmpeg

### Infrastructure
Docker Compose · Dockerfile · Nginx · monitoring (deploy/monitoring/)

---

## Directory Structure

```
music/
├── cmd/api/main.go              # API server entrypoint
├── cmd/worker/main.go           # Background worker entrypoint
├── internal/
│   ├── app/                     # DI container, routes, HTTP server
│   ├── config/                  # Env-based config
│   ├── common/                  # Middleware, errors, pagination, validation
│   ├── modules/                 # 28 domain modules (auth, catalog, playlist, etc.)
│   ├── pkg/audioinfo/           # Audio metadata extraction
│   ├── platform/                # Infrastructure (cache, database, events, etc.)
│   └── workers/                 # 9 background worker types
├── frontend/
│   └── src/
│       ├── assets/              # CSS (fonts.css, main.css), images, fonts
│       ├── components/          # admin/, auth/, common/, contribution/, creator/,
│       │                       # forms/, gamification/, layouts/, music/, social/, tips/
│       ├── composables/         # 27 composables (useAuth, usePlayer, useCatalog, etc.)
│       ├── layouts/             # LayoutMusicApp, LayoutAuth, LayoutAdmin, LayoutEmpty
│       ├── pages/               # app/, admin/, auth/, errors/
│       ├── plugins/             # Axios client, request factory
│       ├── router/              # Vue Router config (app.ts, auth.ts, admin.ts)
│       ├── services/            # API modules, audio engine, WebSocket, storage
│       ├── stores/              # Pinia stores (player, user-auth, feature-flags, etc.)
│       ├── types/               # TypeScript type definitions
│       └── utils/               # Utility functions, PrimeVue preset
├── migrations/                  # 27 SQL migration files
├── deployments/                 # Docker Compose, Dockerfiles
├── docs/                        # DESIGN_SYSTEM.md, SWAGGER.md, architecture/
├── scripts/                     # migrate.sh
├── Dockerfile / Dockerfile.worker / Makefile
└── .env.example
```

---

## Backend Architecture

### Module Pattern (Handler → Service → Repository)

```
internal/modules/{module}/
├── handler.go      # HTTP handlers, request parsing, response
├── service.go      # Business logic
├── repository.go   # Database access
├── model.go        # Domain models
└── dto.go          # Request/Response DTOs
```

### 28 Domain Modules

auth, catalog, playlist, library, queue, player, history, recommendation, search, social, follow, reactions, lyrics, media, analytics, notification, moderation, subscription, ai, creator, gamification, contribution, tips, ingestion, importcmd, dashboard, features, health

### API Base: `/api/v1/`

Roles: Guest → Listener → Creator → Moderator → Admin → Superadmin
WebSocket: `/api/v1/ws` (authenticated)
Event bus: EventTrackPlayed, EventPlaylistCreated, EventSubscriptionPurchased, EventUserRegistered

### Background Workers (9)
ai, analytics, cleanup, eventbus, indexer, notification, recommendation, recommender_v2, transcoder

---

## Frontend Architecture

### App Shell
`App.vue` selects layout via `route.meta.layout`:
- `LayoutMusicApp` — Main app (sidebar + topbar + player + content)
- `LayoutAuth` — `/auth/login`, `/auth/register`
- `LayoutAdmin` — Admin panel
- `LayoutEmpty` — Error pages

### Data Flow
```
Component → Composable → Service API → Axios Client → Backend API
                              ↑
                        Pinia Store (cached state)
```

### Pinia Stores
`useUserAuthStore` (auth tokens, profile, admin detection) · `usePlayerStore` (track, queue, playback state, volume, shuffle, repeat) · `useFeatureFlagsStore` · `useMaintenance` · `usePageLoader`

### Audio Engine
Singleton at `services/player/audio-engine.ts` wrapping HTML5 `<audio>` + Web Audio API AnalyserNode.
Events: play, pause, waiting, playing, canplay, timeupdate, ended, error
Progress via requestAnimationFrame (throttled 250ms)

---

## Design System

- **Dark theme only** with glassmorphism
- **Aurora system:** Dynamic radial gradient blobs extracted from album art
- **Colors:** Primary green `#1DB954`, electric purple `#B646FF`, aurora palette (green, blue, pink, purple)
- **Surfaces:** Dark `#050505`, Base `#0A0A0A`, Raised `#121212`, Overlay `#1A1A1A`
- **Typography:** IRANYekanWeb (Persian), Cabinet Grotesk / Inter / Satoshi (Latin)
- **PrimeVue:** RTL enabled, custom Persian locale, Indigo preset customized
- **Glass tokens:** `.glass`, `.glass-strong`, `.glass-darker`
- **Tailwind CSS v4** `@theme` custom properties in `assets/css/main.css`

---

## Key Conventions

- **Vue:** Composition API with `<script setup lang="ts">`
- **Components:** PascalCase filenames, self-closing tags
- **Composables:** `use` prefix, camelCase files
- **Stores:** `useXxxStore`, Pinia setup stores syntax
- **Pages:** `PageXxx.vue`
- **CSS:** Tailwind utilities + custom `@theme` tokens
- **Routes:** kebab-case paths, dotted notation names (`recommendations.for-you`), domain-separated files
- **Route files:** `router/routes/app.ts`, `auth.ts`, `admin.ts`
- **Git:** Conventional commits preferred, no direct commits without request

---

## Agent Team Architecture

The project has 8 parallel sub-agents defined in `frontend/AGENTS.md`:

| Agent | Focus | Key Exports |
|-------|-------|-------------|
| `@auth` | Auth & Security | `useAuth()`, `useUserAuthStore`, `requireAuth` |
| `@player` | Player & Audio | `usePlayer()`, `audioEngine`, `useQueueManager()` |
| `@catalog` | Catalog & Search | `useCatalogSearch()`, `useTrack()`, `useAlbum()` |
| `@social` | Social & Community | Social API modules, notification composable |
| `@creator` | Creator & Gamification | Creator/gamification API modules |
| `@admin` | Admin & Moderation | Admin composables, admin API modules |
| `@ui` | UI/UX & Layout | `AppLoader`, `AppPageContainer`, layout components |
| `@infra` | Infrastructure | API modules, router instance, common stores |

Each agent owns specific file sets and exports contracts. A coordinator decomposes tasks and resolves cross-agent conflicts. See `frontend/AGENTS.md` for full specification.

---

## How to Run

```bash
# Frontend
cd frontend && cp .env.example .env && npm install && npm run dev

# Backend (Go 1.25+, PG 16+, Redis 7+)
cp .env.example .env && go run cmd/api/main.go

# DB Migrations
./scripts/migrate.sh up

# Lint & Typecheck
cd frontend && npm run lint && npm run type-check
```

---

## Key Docs

- `docs/DESIGN_SYSTEM.md` — Design system & UX specification
- `docs/SWAGGER.md` — Swagger/OpenAPI handler annotations guide
- `docs/architecture/` — API.md, DATABASE.md, OPS.md, OVERVIEW.md, SERVICES.md
- `docs/README.md` — Documentation index
- `frontend/AGENTS.md` — Agent team architecture with file ownership
- `frontend/PERFORMANCE.md` — Performance-sensitive code guidelines
- `FRONTEND_ARCHITECTURE.md` — Frontend architecture & product plan
- `TODO.md` — Active project TODO list
