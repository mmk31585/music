# Muse — Known Concerns & Risks

## ⚠️ CRITICAL Issues (Blocking)

### 1. Frontend Type-Check Fails
**Status**: Broken
**Impact**: CI blocked, dev experience degraded
**File**: `frontend/tsconfig.app.json`
**Cause**: Stale imports, missing `env.d.ts`, type errors accumulate
**Fix**: Update imports, add `env.d.ts`, run `vue-tsc --build` incrementally

### 2. Auth Refresh Contract Broken
**Status**: Broken
**Impact**: Users logged out unexpectedly after token expiry
**File**: Backend auth handler returns `{}` instead of proper token response
**Fix**: Align backend response with frontend `useAuthApi().refresh()` expectation

### 3. Legacy Root Entrypoint (`main.go`)
**Status**: Present
**Impact**: Build confusion, outdated patterns
**File**: `/main.go` at project root (outside `cmd/`)
**Fix**: Remove legacy file, ensure `cmd/api/main.go` is the sole entrypoint

### 4. Login Endpoint Crashes on Incomplete Body
**Status**: Broken
**Cause**: Gin binding vs. validate tag mismatch — empty body bypasses validation
**Fix**: Ensure `binding:"required"` on all DTO fields or use explicit body size check

### 5. Routes Registered Twice
**Status**: Partially fixed (may still occur)
**Cause**: Bootstrap confusion — routes registered in both `main.go` and `app.go`
**Fix**: Single registration point in `internal/app/routes.go`

## 🔴 HIGH Priority

### 6. Missing Admin CRUD Endpoints
**Status**: Frontend calls return 404
**Files**: `/api/v1/admin/tracks`, `/api/v1/admin/artists`, etc.
**Impact**: Admin panel non-functional for core CRUD operations
**Fix**: Implement handler + route registration for missing admin endpoints

### 7. API Contract Mismatch (camelCase vs snake_case)
**Status**: Inconsistent
**Impact**: Frontend may receive unexpected field names
**Note**: Go uses `json:"snake_case"` tags; frontend Zod schemas expect matching

### 8. Missing `media_assets` DB Table
**Status**: Migration not applied
**Impact**: Media upload features broken
**Fix**: Create and run migration `000039_create_media_assets.sql` (or equivalent)

### 9. Player Event Listener Leaks
**Status**: Anonymous functions in `addEventListener` not tracked
**Impact**: Memory leak, duplicate handler execution on remount
**Fix**: Use named/ref-tracked listeners, clean up in `onUnmounted`

### 10. Queue Persists Full Objects in localStorage
**Status**: Suboptimal
**Impact**: localStorage bloat, stale track data
**Fix**: Persist only track IDs + cache metadata separately

### 11. Shuffle Lacks History Stack
**Status**: Missing
**Impact**: "Previous" after shuffle plays wrong track
**Fix**: Implement LIFO history for shuffle mode (`prevShuffleHistory`)

### 12. No Retry on Audio Error
**Status**: Track fails silently on network error
**Fix**: Add exponential backoff retry (max 3) before advancing queue

### 13. `resetTracks()` Not Implemented
**Status**: Referenced but not defined
**Impact**: Broken reset flow in admin/player
**Fix**: Implement the function

### 14. Insecure Default JWT Secrets
**Status**: Hardcoded in `.env` and `.env.example`
**Impact**: Security risk if `.env` is committed or leaked
**Fix**: Generate per-deployment secrets, document rotation, add warning

## 🟡 MEDIUM Priority

### 15. Upload Lacks Duration Extraction
**Status**: Uploaded audio files don't auto-extract duration
**Impact**: Tracks show 0:00 duration until manually updated

### 16. Media Serving Not Range-Aware
**Status**: No `Accept-Ranges: bytes` / `Content-Range` headers
**Impact**: No seeking in audio streams, no partial content

### 17. Search Uses `ILIKE '%query%'`
**Status**: Full-text search not enabled (OpenSearch may be unused)
**Impact**: Inefficient, no relevance ranking, no prefix matching

### 18. Frontend Organized by Type, Not Domain
**Status**: Files split across `components/`, `composables/`, `pages/`, `services/`, `stores/`
**Impact**: Related functionality scattered, harder to navigate at scale
**Note**: TODO.md mentions planned migration to feature-based modules

### 19. Backend DTOs Double as Domain Models
**Status**: Same structs used for DB, API request, and API response
**Impact**: Tight coupling, cannot evolve API independently of DB schema

### 20. Middleware Consolidation Needed
**Status**: Auth middleware logic split across `common/middleware/` and module-level `middleware.go`
**Impact**: Duplication, inconsistent error handling

## 🟢 LOW Priority

### 21. Unused Interfaces/Types
**Status**: Various interfaces defined but not implemented
**Impact**: Dead code, confusing for newcomers

### 22. Missing Favicon
**Status**: Browser tab shows default Vite icon

### 23. Wrong App Title ("Achive" not "Muse")
**Status**: Title tag says "Achive" instead of "Muse" in some places
**Note**: Low priority but affects brand consistency

### 24. Accessibility Gaps
**Status**: Some interactive elements lack `aria-label`, keyboard support
**Note**: FRONTEND_ARCHITECTURE.md lists a11y as low priority

### 25. RTL/LTR Consistency
**Status**: Some components use hardcoded `left`/`right` instead of logical properties

## 🏗️ Architecture Concerns

### 26. Microservices Prematurity
**Docs** (`docs/architecture/SERVICES.md`) describe 14 microservices, but **current codebase is a monolith**. Risk of premature decomposition vs. actually scaling the monolith.

### 27. In-Process Event Bus
**Current**: Simple in-process bus with transactional outbox
**Risk**: No persistence if process dies before outbox flush. No horizontal scaling.

### 28. Single Database
**Current**: All modules share one PostgreSQL database
**Risk**: No read replicas, no CQRS, no domain isolation

### 29. No Schema Versioning for Frontend-Backend Contracts
**Risk**: Frontend Zod schemas can drift from backend responses without detection

### 30. No E2E Tests
**Risk**: Integration bugs (like the auth refresh contract) not caught before deploy

## 📋 Database Migration Roadmap (from TODO.md)

Pending migrations (in dependency order):
1. `media_assets` — Create missing media_assets table
2. `playlists` — Playlist schema refactor
3. `liked_tracks` — Liked tracks table
4. `listening_history` — Listening history partition fix
5. `playback_sessions` — Session tracking
6. `pg_trgm` — Enable trigram extension for search
7. `audio_asset_id` — Add audio_asset_id to tracks

## 🧪 Testing Debt

| Area | What's Missing | Impact |
|------|---------------|--------|
| API handlers | Table-driven tests for all handlers | Regression risk on every change |
| Frontend components | No component tests | UI breaks undetected |
| Frontend stores | No Pinia store tests | State logic untested |
| Integration | No API E2E tests | Contract drift undetected |
| WebSocket | No socket tests | Real-time feature fragility |

## 📊 Coverage Dashboard

| Metric | Current | Target |
|--------|---------|--------|
| Go line coverage | Unknown (below 30%) | ≥ 30% (CI gate) |
| Frontend coverage | Near 0% | ≥ 50% |
| ML service coverage | Minimal | ≥ 40% |
| E2E coverage | None | Critical paths only |
