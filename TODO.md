 # Muse — Project TODO

> Last updated: 2026-06-22
> Sources: inline code TODOs, architecture audits, live exploration audit

---

## Priority Legend

| Tag | Meaning |
|-----|---------|
| 🔴 **CRITICAL** | Build broken, security hole, or core feature non-functional |
| 🟠 **HIGH** | Blocks development, causes incorrect behavior |
| 🟡 **MEDIUM** | Significant UX/arch flaw, tech debt with known impact |
| 🔵 **LOW** | Polish, edge cases, nice-to-have |

---

## 🔴 CRITICAL

### 1. Frontend type-check is red

`npm run type-check` fails. Root causes:
- Stale imports referencing removed files
- Missing `ImportMeta.env` type declarations
- Auth store `storage` generic returns `string | T | null` instead of `T | null`
- Stale test config / dev dependency references

**Files:** `frontend/src/stores/user-auth.ts`, `frontend/src/plugins/client/request-factory.ts`, `frontend/src/composables/useLoading.ts`, missing `env.d.ts`

### 2. Auth refresh contract broken

Backend expects `{ refreshToken }` body on `/auth/refresh`, but frontend `refresh()` posts an empty body. Expired access tokens force logout or loop failed refreshes.

**Files:** `frontend/src/services/api/auth/refresh.ts`, `frontend/src/stores/user-auth.ts`

### 3. Legacy entrypoint pollutes root

Root `main.go` is a separate legacy in-memory song server that coexists with `cmd/api/main.go`. Running `go run .` vs `go run ./cmd/api` starts different apps. Legacy song code under `internal/handler/`, `internal/service/`, `internal/repository/`, `internal/domain/` no longer compiles (references removed `UploadDir`, `AllowedAudioExt`, `MaxUploadSize`).

### 4. `POST /api/v1/auth/login` crashes with 500 on empty/incomplete body

Sending `{}` or `{"email":"..."}` (missing fields) crashes the handler with HTTP 500 `INTERNAL_ERROR`. Should return 400/422 `VALIDATION_ERROR` like `/auth/register` does.

**Root cause:** `LoginRequest` uses `binding:"required"` tags (Gin-level validation during `ShouldBindJSON`) instead of `validate:"required"` tags (validator-level). When Gin returns a validation error, `response.Error()` passes it to `ToAppError()`, which doesn't recognize Gin errors and falls back to `500 Internal Server Error`.

**Fix:** Replace `binding` tags with `validate` tags in `LoginRequest` (matching `RegisterRequest` pattern), and fix the error handler in `Login()` to return proper BadRequest like `Refresh()` does.

**Files:**
- `internal/modules/auth/dto.go:12-15` — `LoginRequest` uses `binding` instead of `validate`
- `internal/modules/auth/handler.go:79-83` — raw err passed to `response.Error()` instead of wrapped BadRequest

### 5. Routes registered twice

`Bootstrap()` creates a Gin engine and calls `RegisterRoutes`, then `cmd/api/main.go` creates a *second* Gin engine and calls `RegisterRoutes` again. Runtime uses the second engine, making startup logic confusing and middleware differences hard to track.

**Files:** `internal/app/bootstrap.go`, `cmd/api/main.go`

---

## 🟠 HIGH

### 6. Missing backend admin CRUD endpoints

Live API exploration revealed 404s on admin endpoints that frontend expects:

| Endpoint | Status | Frontend page that depends on it |
|----------|--------|----------------------------------|
| `GET /admin/catalog/artists` | 404 | PageAdminCatalog (artist listing) |
| `GET /admin/ingestion` | 404 | PageAdminIngestion |
| `GET /admin/moderation` | 404 | PageAdminModeration |
| `GET /admin/subscriptions` | 404 | PageAdminSubscriptions |

**Swagger spec** at `GET /api/v1/swagger/doc.json` also omits these paths. Register routes and handlers for each.

### 7. API contract mismatch — camelCase vs snake_case

Backend catalog models emit `artistId`, `audioUrl`, `durationSeconds` (camelCase). Frontend Zod schemas expect `artist_id`, `audio_url`, `duration_seconds` (snake_case). This breaks rendered fields and form payloads everywhere.

**Files:** Backend `internal/modules/catalog/`, Frontend `services/api/catalog/schemas.ts`, admin pages

### 8. No durable media asset table

Uploads create files on disk but have no DB record. No hash, MIME, size, duration, ownership, or lifecycle tracking. Orphan files cannot be cleaned up. Duplicate detection is filesystem-only.

**Needs:** Migration for `media_assets` table + repository + status pipeline (pending → attached → cleanup)

### 9. Player event listeners leak

AudioEngine uses anonymous arrow functions for event handlers. No `dispose()` method. HMR, test resets, or future multi-device sync can double-bind listeners.

**Files:** `frontend/src/services/player/audio-engine.ts`

### 10. Player queue persists full objects in localStorage

Stores complete `PlaybackTrack` objects — stale URLs/titles accumulate, localStorage bloat. Should persist IDs + snapshot, then hydrate from API on app load.

**File:** `frontend/src/stores/player.ts`

### 11. Shuffle has no history stack

Previous track after shuffle is unpredictable. Needs a proper play history stack.

**File:** `frontend/src/services/player/queue-manager.ts`

### 12. No retry on audio playback error

Error handler only sets `error` message. Transient network errors kill playback permanently. Should retry source once, then optionally skip to next track.

**File:** `frontend/src/stores/player.ts`

### 13. `resetTracks()` called but not implemented

`useAdminCatalog.ts` calls `tracks.resetTracks()` but `useAdminTracks.ts` doesn't export it. Runtime error on admin catalog page.

**Files:**
- `frontend/src/composables/admin/useAdminCatalog.ts:49` — `// TODO HIGH: tracks.resetTracks() does not exist`
- `frontend/src/composables/admin/useAdminTracks.ts:293` — `// TODO MEDIUM: Add resetTracks() method`

### 14. JWT secrets have insecure defaults

Development config has default JWT secrets. No validation rejects them in production. Deploy can be trivially compromised.

**Files:** `internal/config/`, `.env.example`

---

## 🟡 MEDIUM

### 15. Upload lacks duration extraction

Upload service detects MIME/hash/size but not duration. Admins manually enter duration. Use `dhowden/tag` or ffprobe behind an interface.

**File:** `internal/pkg/audioinfo/`

### 16. Media serving is not range-aware

Uses `c.File` instead of `http.ServeContent`. No `Cache-Control`, `ETag`, `Accept-Ranges`, or verified content type. Poor streaming behavior.

**File:** `internal/app/media_routes.go`

### 17. Search uses `ILIKE '%query%'`

No trigram indexes. Slow at scale. Needs PostgreSQL `pg_trgm` extension + GIN indexes.

**File:** `internal/modules/search/repository.go`, add migration

### 18. Frontend organized by type, not domain

Global `components/`, `composables/`, `services/`, `stores/` instead of feature modules. Leads to coupling through shared barrels.

**Strategy:** Incrementally migrate to `features/player` → `features/uploads` → `features/tracks` → `features/auth` (per refactor audit)

### 19. Backend DTOs double as domain models

Catalog models have both `db` and `json` struct tags. DB schema changes leak directly to API contracts. Add transport DTO mapping.

**File:** `internal/modules/catalog/model.go`

### 20. No middleware consolidation

Middleware exists in both `internal/common/middleware/` and `internal/middleware/` with different Gin/net-http assumptions. CORS is hardcoded in `cmd/api/main.go` despite being in config.

### 21. `AnyTrack` defeats type safety

Admin track page uses `AnyTrack` type instead of proper discriminated union. See `PageAdminTracks.vue:369`.

### 22. Ingestion review component is 1135 lines

`PageAdminIngestionReview.vue` is too large. Split into `ArtistStep`, `AlbumStep`, `TrackStep`, `ConfirmStep` sub-components.

### 23. Ingestion polling should be WebSocket

`pollEnrichment` uses 2s interval × 30 = 60s of polling. Replace with WebSocket or SSE.

### 24. Admin pages use wrong theme classes

`PageAdminIngestion.vue` and `PageAdminIngestionReview.vue` use `surface-*`/`primary-*` PrimeVue classes instead of `slate-*`/`white/*` admin dark theme.

### 25. `formatDuration` treats 0 as falsy

0-second tracks show `—` instead of `0:00`. Affects `PageAdminTracks.vue:733` and `PageAdminCatalog.vue:323`.

### 26. Missing database tables

Core tables not yet migrated:
- `media_assets` (upload tracking)
- `playlists` / `playlist_tracks`
- `liked_tracks`
- `listening_history`
- `playback_sessions`

### 27. `UpdateTrack` cannot clear nullable fields

Nil `albumId`, `audioUrl`, `coverUrl` means "not provided" — no way to explicitly set null.

**File:** `internal/modules/catalog/handler.go`

### 28. Upload uses `body` instead of `data`

`UploadMediaDialog.vue` sends multipart as `body` — Axios ignores it in the custom wrapper. Need to use the correct Axios key or bypass the wrapper for uploads.

### 29. `RegisterTrack` reads raw body after `ShouldBindJSON` failure

Consumes the request body, making error logs misleading.

**File:** `internal/modules/catalog/handler.go`

### 30. Auth store is not strict-safe

Storage helper returns `string | T | null` instead of properly typed `T | null`. Causes TypeScript errors and potential runtime shape issues.

**File:** `frontend/src/stores/user-auth.ts`

### 31. Frontend login form has debug `console.log` calls

Remove before production.

**File:** `frontend/src/composables/auth/useLoginForm.ts`

### 32. No route-level chunk strategy

Vite build has only a large chunk warning. Needs manual route-level code splitting for admin, player, and catalog bundles.

### 33. Admin catalog preview fetches ALL items

`PageAdminCatalog.vue:303` fetches entire catalog just for preview counts — wasteful as data grows. Add count-only endpoints or limit queries.

---

## 🔵 LOW

### 34. `VersionedEvent` / `VersionedBaseEvent` defined but unused

Event versioning types exist but nothing uses them. Clean up or implement.

**File:** `internal/platform/events/version.go`

### 35. `RepositoryInterface` defined but unused

Auth module has an interface that's never used as injection type. Service uses concrete `*Repository` directly.

**File:** `internal/modules/auth/interfaces.go:8`

### 36. `RequireRole` is exact-match only

No role hierarchy (e.g., super-admin > admin). Extend when multi-role support is needed.

**File:** `internal/modules/auth/middleware.go:53`

### 37. Album upsert doesn't fall back to embedded cover

Ingestion finalization sets album `cover_url` only from user-provided field. If user provides no cover but embedded cover asset exists, album gets NULL.

**File:** `internal/modules/ingestion/finalization/service.go:118`

### 38. Silent catch in admin subscriptions

`PageAdminSubscriptions.vue:141` silently catches errors — should show a toast.

### 39. Media page has no pagination

`PageAdminMedia.vue:105` — large media lists will degrade performance.

### 40. Empty states missing

Public track lists and admin tables lack empty/error/retry states.

### 41. Accessibility gaps

Buttons lack systematic ARIA labels. Table action buttons, dialogs, empty states, and keyboard flows need review. Skip-link exists in LayoutMusicApp but not universally tested.

### 42. RTL/LTR consistency

App mixes English UI copy, Persian locale, and LTR/RTL assumptions. Some components may not flip correctly in RTL mode.

### 43. Missing favicon

HTML `<link rel="icon">` references `/archive.ico` which doesn't exist. Browser tab shows no icon.

**File:** `frontend/index.html`

### 44. App title is "Achive" not "Muse"

HTML `<title>` is "Achive" — appears to be leftover from previous branding.

**File:** `frontend/index.html`

### 45. Frontend missing routes: `/admin/settings`, wrong path for import

`/admin/settings` has no route definition → shows 404 page.  
`/admin/artists/import` should be `/admin/import/artist` (exists but under different path).

**Files:** `frontend/src/router/routes/admin.ts`

---

## ✅ Recently Completed

### 46. Import by Artist feature (search + batch)

Backend `GET /admin/import/artist?name=...` searches MusicBrainz (fallback from Deezer) for artist discography. Returns albums grouped with Cover Art Archive images, tracks with durations/sources. Frontend `PageAdminImportArtist.vue` has album-grid UI with multi-select (per-track, per-album, select-all), batch import trigger, and progress bar.

**Files:** `internal/modules/importcmd/artist_search.go`, `handler.go`, `import_service.go`, `routes.go`, `dto.go`, `frontend/src/pages/admin/PageAdminImportArtist.vue`, `frontend/src/router/routes/admin.ts`, `AdminSidebar.vue`, `frontend/src/services/api/importcmd/routes.ts`, `frontend/src/services/api/importcmd/types.ts`

### 47. Import worker consumer group fix

Worker created Redis consumer group with `$` (start from latest), so jobs enqueued *before* the worker started were never delivered. Changed to `"0"` (start from beginning), ensuring all existing stream messages are available.

**File:** `internal/modules/importcmd/worker/worker.go:186`

### 48. Import worker JSON payload unwrap fix

Worker's `poll()` received double-wrapped JSON `{"payload":{"id":...,"title":...}}` inside a Redis field also named `payload`. Fixed with `json.RawMessage` intermediate unmarshal so `Job` struct fields populate correctly instead of silently zeroing out.

**File:** `internal/modules/importcmd/worker/worker.go:235-240`

---

## 📦 Database Migration Roadmap

Ordered by dependency:

1. `media_assets` — upload/media tracking with hash, MIME, size, duration, status
2. `playlists` + `playlist_tracks` — user playlists with ordering
3. `liked_tracks` — user track likes
4. `listening_history` — time-partitioned playback log
5. `playback_sessions` — cross-device state sync
6. `pg_trgm` extension + GIN indexes for catalog search
7. `tracks.audio_asset_id` — link track to media asset (optional FK)

---

## 🧪 Testing Plan

- Backend unit: media storage path safety, MIME validation, duplicate hashing, cleanup
- Backend handler: `POST /admin/catalog/tracks/upload`
- Repository: track create/update, playlist ordering (test Postgres)
- Frontend store: player repeat/shuffle/seek/persist (Vitest)
- API services: auth refresh, track upload payload mapping
- Router: auth guard, lazy-loaded domain pages
- E2E smoke: login → create artist → create track with audio → see in list → play

---

## 🏗️ Architecture Migration Path

### Frontend feature modules (incremental)

```
src/features/player/     ← first, clear boundary, low API coupling
src/features/uploads/    ← second, isolates media API + progress UI
src/features/tracks/     ← third, track rows, forms, types
src/features/auth/       ← fix token contract + storage while moving
src/features/search/     ← then artists, albums, playlists, library
```

### Backend refactor

```
internal/app/usecase/            ← coordinate cross-module ops
internal/common/httpx/            ← consolidate middleware
internal/platform/storage/        ← clean storage abstraction
Legacy: remove internal/handler/, internal/service/, internal/repository/, internal/domain/, root main.go
```
