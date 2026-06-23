# Muse — AI Agents Guide

This file defines the AI agent system for the Muse Persian music platform. It helps AI coding assistants understand how to work with this project effectively.

## Agent Team

Muse uses a **supervised team model** with specialized agents:

### Coordinator Agent (`@muse`)
- Orchestrates work across all sub-agents
- Knows the full stack and delegates tasks
- Used for multi-domain tasks and architecture decisions

### Domain Agents

| Agent | Responsibility | Handles |
|-------|---------------|---------|
| `@muse-auth` | Auth & Security | Login, JWT, tokens, auth guards, session restoration |
| `@muse-player` | Player & Audio | Audio engine, queue, radio, visualizer, PiP, lyrics |
| `@muse-catalog` | Catalog & Search | Tracks, albums, artists, genres, search, playlists, library, recommendations |
| `@muse-social` | Social & Community | Activity feeds, parties, rooms, clubs, notifications |
| `@muse-creator` | Creator & Gamification | Dashboard, badges, XP, tips, subscriptions, challenges |
| `@muse-admin` | Admin & Moderation | Admin panel, CRUD, media upload, moderation |
| `@muse-ui` | UI/UX & Layout | App shell, components, styling, global stores, design system |
| `@muse-infra` | Infrastructure | API barrel, router, stores, services, composables, WebSocket |
| `@muse-sre` | Site Reliability | Monitor health, detect crashes, restart servers, watchdog |
| `@muse-team` | Team Supervisor | Shared task list, peer-to-peer coordination |
| `@muse-ai` | AI/ML | Embeddings, mood analysis, playlist generation, Whisper, Celery |

## How Agents Work Together

1. **Coordinator** (`@muse`) receives a high-level task
2. **Breaks it down** into sub-tasks and assigns to domain agents
3. **Domain agents** may communicate directly with each other (peer-to-peer):
   - `@muse-ui` asks `@muse-catalog` for API endpoint shapes
   - `@muse-catalog` asks `@muse-infra` for barrel registration
4. **Coordinator** monitors progress, resolves conflicts, approves final result

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
