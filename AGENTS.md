# Muse — AI Agents Guide

This file defines the AI agent system for the Muse Persian music platform. It helps AI coding assistants understand how to work with this project effectively.

## Architecture: Two Operating Modes

Muse supports **two modes** depending on the nature of the task:

### Mode 1: Coordinator + Sub-agents (`@muse`)
**Best for:** Multi-domain but loosely coupled tasks, independent fixes, research.

The `@muse` coordinator receives a task, breaks it down, assigns work to isolated sub-agents, collects results, and merges. Sub-agents work independently and report back.

### Mode 2: Team Supervisor + Peers (`@muse-team`)
**Best for:** Tightly coupled, cross-domain work — end-to-end features, shared contracts, multi-layer bugs.

The `@muse-team` supervisor creates a shared task list and deploys **teammates who talk directly** (peer-to-peer). No isolation, no roundtrips through the supervisor. Peers unblock each other instantly.

```
┌──────────────────────────────────────────┐
│           Team Supervisor                 │
│   - Creates task list                     │
│   - Monitors progress                     │
│   - Resolves conflicts                    │
│   - Approves final result                 │
└─┬──────────┬──────────┬───────────────────┘
  │          │          │
  ▼          ▼          ▼
┌────────┐ ┌────────┐ ┌────────┐
│Backend │◄┤Frontend│►│ Infra  │
│Team    │─┤Team    │─┤Team    │
└────────┘ └────────┘ └────────┘
   │           │           │
   └───────────┴───────────┘
    Peer-to-peer (direct comms)
    "Hey @muse-catalog, need /api/v1/foo shape"
    "Hey @muse-ui, here's the schema, unblocked"
```

### When to use which

| Use `@muse` (Coordinator) | Use `@muse-team` (Supervisor) |
|---|---|
| Fix a UI bug | Build an end-to-end feature (API + frontend + routing) |
| Update a DB query | Refactor shared types across frontend/backend |
| Add a single component | Debug a request spanning frontend → API → DB |
| Research/audit/diagnose | Any task needing frontend/backend contract agreement |

## Domain Agents

| Agent | Mode | Responsibility | Handles |
|-------|------|---------------|---------|
| `@muse-auth` | subagent | Auth & Security | Login, JWT, tokens, auth guards, session restoration |
| `@muse-player` | subagent | Player & Audio | Audio engine, queue, radio, visualizer, PiP, lyrics |
| `@muse-catalog` | subagent | Catalog & Search | Tracks, albums, artists, genres, search, playlists, library, recommendations |
| `@muse-social` | subagent | Social & Community | Activity feeds, parties, rooms, clubs, notifications |
| `@muse-creator` | subagent | Creator & Gamification | Dashboard, badges, XP, tips, subscriptions, challenges |
| `@muse-admin` | subagent | Admin & Moderation | Admin panel, CRUD, media upload, moderation |
| `@muse-ui` | subagent | UI/UX & Layout | App shell, components, styling, global stores, design system |
| `@muse-infra` | subagent | Infrastructure | API barrel, router, stores, services, composables, WebSocket |
| `@muse-sre` | primary | Site Reliability | Monitor health, detect crashes, restart servers, watchdog |
| `@muse-team` | primary | Team Supervisor | Shared task list, peer-to-peer coordination |
| `@muse-ai` | skill | AI/ML | Embeddings, mood analysis, playlist generation, Whisper, Celery |

## How the Team Mode Works

1. **`@muse-team`** receives a cross-domain task (e.g., "build artist stats")
2. **Creates a shared task list**: Backend endpoint → Frontend panel → Infra registration
3. **Deploys teammates** with direct peer-to-peer briefing:
   - `@muse-catalog` (backend), `@muse-ui` (frontend), `@muse-infra` (infra)
   - Explicit instruction: *"Message each other directly when blocked — don't wait for me"*
4. **Peers communicate directly**:
   - `@muse-ui` → `@muse-catalog`: "I need `GET /api/v1/artists/:id/stats` — response shape?"
   - `@muse-catalog` → `@muse-ui`: "Done. Schema: `{ listens, plays, rank, trend }`. Registered in router."
   - `@muse-ui`: "Proceeding with panel build."
5. **Supervisor monitors** the task list, resolves blocking, and approves final result.

## Design Principles

### Model-Agnostic
All agents and skills are **model-agnostic** — no agent requires a specific LLM model. No `model:` field binds an agent to Claude, GPT, Gemini, or any other provider. This ensures portability across AI coding hosts.

### Self-Diagnosing
Each Muse domain agent has built-in self-diagnosis: health checks, file existence checks, git log scanning, and type-check verification steps they run when invoked.

### Contract-First
Every domain agent exports a clear contract (`useAuth()`, `usePlayerStore()`, `audioEngine`, etc.) and declares its dependencies on other agents. This enables parallel work and clear boundaries.

## For AI Contributors

### Before Making Changes
1. Read `CONTEXT.md` for domain glossary
2. Read `FRONTEND_ARCHITECTURE.md` for frontend architecture
3. Read `docs/DESIGN_SYSTEM.md` for design tokens
4. Check the existing `internal/` and `frontend/src/` structure for patterns

### Coding Conventions
- **Go**: Follow `.github/instructions/go.instructions.md`
- **Vue 3**: Follow `.github/instructions/vue.instructions.md`
- **Localization**: Follow `.github/instructions/localization.instructions.md`
- **Accessibility**: Follow `.github/instructions/a11y.instructions.md`
- **Docker**: Follow `.github/instructions/containerization.instructions.md`

### Commit Messages
Follow Conventional Commits format:
```
feat: add artist stats endpoint
fix: resolve queue shuffle ordering
refactor: extract player engine interface
```

### Pull Request Process
1. Create feature branch from `main`
2. Write tests for new functionality
3. Ensure all existing tests pass
4. Update documentation if needed
5. Create PR with descriptive title and scope

## Running the Project

### Backend
```bash
# Run API server
go run ./cmd/api

# Run tests
go test ./...
```

### Frontend
```bash
cd frontend
npm run dev
npm run test
npm run build
```

### Docker
```bash
docker compose up -d
```

## Key Files

| File | Purpose |
|------|---------|
| `CONTEXT.md` | Domain glossary and core concepts |
| `FRONTEND_ARCHITECTURE.md` | Frontend architecture docs |
| `docs/DESIGN_SYSTEM.md` | Design tokens and brand identity |
| `.github/copilot-instructions.md` | Shared coding standards |
| `.github/instructions/` | Per-domain instruction files |
| `.opencode/opencode.json` | Agent configuration |
