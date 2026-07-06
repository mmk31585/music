# Muse — Coding Conventions

## Backend (Go)

### Core Rules
- Idiomatic Go per [Effective Go](https://go.dev/doc/effective_go), [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- Package layout: `cmd/` for main, `internal/` for business logic
- Use `gofmt` + `goimports` for formatting
- Return early, keep happy path left-aligned
- Wrap errors with `fmt.Errorf("%w")`
- Prefer stdlib over custom implementations

### Naming
- **Packages**: lowercase, single word, no underscores (`auth`, not `authentication_service`)
- **Exported**: PascalCase (`NewHandler`, `Authenticate`)
- **Unexported**: camelCase (`newHandler`, `authenticate`)
- **Interfaces**: `-er` suffix (`Reader`, `StorageProvider`)
- No stuttering: `http.Server` not `http.HTTPServer`

### Module Structure
Each domain module follows:
```
module/
├── model.go         # DB structs, domain types
├── repository.go    # SQL queries (sqlx)
├── service.go       # Business logic
├── handler.go       # Gin HTTP handlers
├── dto.go           # Request/response DTOs
├── routes.go        # RegisterRoutes(router, handler, mw)
├── *_test.go        # Table-driven tests
└── event_handler.go # Event bus subscribers (optional)
```

### API Response Format
```json
{
  "data": { ... },       // single object
  "data": [ ... ],       // array
  "meta": {              // pagination
    "page": 1,
    "per_page": 20,
    "total": 100
  },
  "error": {             // on error
    "code": "NOT_FOUND",
    "message": "Track not found"
  }
}
```

### Route Registration
```go
func RegisterRoutes(rg *gin.RouterGroup, h *Handler, mw ...gin.HandlerFunc) {
    r := rg.Group("/resource", mw...)
    r.GET("/", h.List)
    r.GET("/:id", h.Get)
    r.POST("/", h.Create)
    r.PUT("/:id", h.Update)
    r.DELETE("/:id", h.Delete)
}
```

### Error Handling
- Define domain error types in `internal/common/errors/`
- Use `errors.Is()` / `errors.As()` for checking
- Don't log AND return — choose one

### Concurrency
- `go >= 1.25` in `go.mod` → use new `sync.WaitGroup.Go()` method
- Always know how goroutines exit
- Use `sync.Once` for one-time init
- `sync.RWMutex` for read-heavy state

## Frontend (Vue 3 + TypeScript)

### Vue Conventions
- **Always** `<script setup lang="ts">`
- **Always** Composition API (no Options API)
- Use `defineProps`/`defineEmits` with TypeScript generics
- Components in PascalCase files (`TrackRow.vue`)
- Composables in camelCase with `use` prefix (`useAuth()`)
- Pages prefixed with `Page` (`PageHome.vue`)

### TypeScript
- Strict mode enabled
- Types defined next to their usage (not in global `types/`)
- API types in `services/api/*/types.ts`
- Use Zod schemas for API response validation

### State Management (Pinia)
- Setup function pattern: `defineStore('name', () => { ... })`
- Actions over mutations (Pinia setup stores don't have mutations)
- Store references: `useUserAuthStore()`, `usePlayerStore()`

### API Client Pattern
```typescript
// services/api/<domain>/routes.ts
export function useXxxApi() {
  const request = useRequest()
  return {
    list: (params) => request.get('/resource', params),
    get: (id) => request.get(`/resource/${id}`),
    create: (data) => request.post('/resource', data),
  }
}
```

### Composable Pattern
```typescript
// composables/<domain>/useXxx.ts
export function useXxx() {
  const data = ref(null)
  const loading = ref(false)
  async function fetch() { ... }
  return { data, loading, fetch }
}
```

### Component Organization
- **Smart** (pages): fetch data, manage state, compose dumb components
- **Dumb** (components/): receive props, emit events, no direct API calls

## Design System

### Glassmorphism Hierarchy
```css
.glass           { background: rgba(18,18,18,0.6); backdrop-filter: blur(20px); }
.glass-strong    { background: rgba(18,18,18,0.8); backdrop-filter: blur(28px); }
.glass-darker    { background: rgba(10,10,10,0.85); backdrop-filter: blur(32px); }
```

### Colors
- Primary: `#1DB954` (green)
- Accent: `#B646FF` (electric purple)
- Dark surface: `#121212`, `#0A0A0A`
- Slate grays for text hierarchy (50-500)

### Typography
- Persian: `IRANYekanWeb` (light 300, regular 400, bold 700, black 900)
- Latin: `Cabinet Grotesk` / `Inter` / `Satoshi`
- Base: 14px on body

### Spacing
- 4px base unit, scale: 1, 2, 3, 4, 5, 6, 8, 10, 12, 16, 20, 24
- Tailwind utility classes throughout

## Persian/RTL

- Logical CSS properties: `margin-inline-start`, `padding-inline-end`, `border-inline-start`
- Set `dir="rtl"` reactively based on locale
- `Intl.NumberFormat('fa-IR')` for Persian numerals
- `Intl.DateTimeFormat('fa-IR')` for Jalali calendar dates
- Font preloading for IRANYekan font family

## Accessibility (WCAG 2.2 AA)

- **4.5:1** minimum contrast on all text (including glassmorphism surfaces)
- **3:1** minimum contrast for UI components and graphics
- Player controls require `aria-label` (e.g., "Play", "Pause", "Next track")
- Volume slider: `role="slider"` with `aria-valuemin/max/now`
- Keyboard operable: all interactive elements
- `prefers-reduced-motion` respected for animations/visualizer
- Skip link as first focusable element
- Touch targets ≥ 24×24 CSS px

## Docker

- Multi-stage builds (builder + runtime)
- Alpine base images
- Non-root user in production (`appuser:1001`)
- Build cache mounts for `go mod download`

## Git

- Conventional Commits: `feat:`, `fix:`, `refactor:`, `chore:`, `docs:`, `test:`
- Feature branches from `main`
- PRs with descriptive titles and scope

## AI Agent System

- 11 domain agents (auth, player, catalog, social, creator, admin, ui, infra, sre, team, ai)
- Two modes: Coordinator (`@muse`) for loose coupling, Team Supervisor (`@muse-team`) for tight coupling
- Agents own specific file sets (see `AGENTS.md`)
- Contract exports must be flagged for coordinator review
