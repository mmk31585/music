# Muse Persian Music Platform — Copilot Instructions

Muse is a full-stack Persian music streaming platform with:
- **Backend**: Go 1.22+ with Gin framework, PostgreSQL, Redis
- **Frontend**: Vue 3 (Composition API, `<script setup lang="ts">`) with Pinia, Vue Router, Vite
- **ML/AI**: Python microservice with Celery, Whisper for speech-to-text
- **Infrastructure**: Docker, Docker Compose, GitHub Actions CI/CD
- **Design**: Glassmorphism UI system with Persian (RTL) support

## Project Structure
```
cmd/           — Go entrypoints (api, worker, etc.)
internal/      — Go internal packages (auth, catalog, player, etc.)
frontend/      — Vue 3 SPA
moja-ml-service/ — Python ML microservice
migrations/    — Database migrations
deployments/   — Deployment configs
docs/          — Architecture docs and design system
```

## Coding Standards

### Backend (Go)
- Follow idiomatic Go practices (see `.github/instructions/go.instructions.md`)
- Package layout: `cmd/` for main packages, `internal/` for business logic
- Use Gin framework for HTTP routing and middleware
- Use GORM or sqlx for database access
- All API handlers in `internal/api/handlers/`
- Error handling: wrap errors with context using `fmt.Errorf("%w")`
- Write tests using table-driven tests in `_test.go` files

### Frontend (Vue 3)
- Use Composition API with `<script setup lang="ts">` exclusively
- TypeScript everywhere: components, composables, stores
- Pinia for state management, Vue Router for routing
- Vite for build tooling, Vitest for testing
- See `.github/instructions/vue.instructions.md` for detailed guidance

### Persian/RTL Support
- All user-facing UI must support Persian (Farsi) and RTL layout
- Use logical CSS properties (`margin-inline-start`, `padding-inline-end`, etc.)
- Use `vue-i18n` for translations, `Intl` for date/number formatting
- Set `dir="rtl"` reactively based on locale
- See `.github/instructions/localization.instructions.md`

### Accessibility
- Target WCAG 2.2 AA compliance
- Glassmorphism elements must maintain 4.5:1 text contrast
- Music player controls must have `aria-label` and keyboard support
- Visualizer must respect `prefers-reduced-motion`
- See `.github/instructions/a11y.instructions.md`

### Docker
- Multi-stage builds for all services
- Alpine-based images for Go services
- Non-root users in production containers
- See `.github/instructions/containerization.instructions.md`

## Architecture
- See `docs/DESIGN_SYSTEM.md` for UI/UX design tokens and brand identity
- See `FRONTEND_ARCHITECTURE.md` for frontend architecture
- See `docs/architecture/` for system architecture diagrams
- See `CONTEXT.md` for domain glossary and core concepts

## Agents
This project uses a multi-agent system. See `AGENTS.md` for details.
