# Current Architecture Refactor Audit

Date: 2026-05-25

## Executive Summary

The app has a workable foundation, but it is not yet a scalable Spotify-like architecture. The main issue is not missing features; it is mixed architectural generations. The backend has an active modular Gin app under `internal/modules/*`, but also legacy root `main.go`, legacy song handlers/services/repositories, and duplicated middleware packages. The frontend has a mostly global structure (`components`, `composables`, `services`, `stores`) rather than feature modules, so auth, catalog, uploads, admin, and player behavior cross through shared barrels and global services.

The correct refactor strategy is incremental:

1. Stabilize correctness and build/type safety.
2. Extract feature modules one at a time, starting with the most self-contained domains: `player`, `uploads`, then `tracks`.
3. Add durable upload domain modeling before expanding media features.
4. Add domain routes only after API contracts and route-level data loading are normalized.
5. Keep existing APIs working while moving internal frontend/backend boundaries.

This audit intentionally does not recommend a giant rewrite.

## Current Structure Assessment

### Frontend

- Routes are centralized in `src/router/routes/*`; this is simple but not feature-owned. Domain routes are not yet implemented for `/track/:id`, `/album/:id`, `/artist/:id`, `/playlist/:id`, `/genre/:id`.
- API services live in `src/services/api/*`; they are typed with Zod, but response normalization is inconsistent and some schemas adapt to backend drift rather than enforcing a clean transport contract.
- Components are grouped by broad UI area (`admin`, `music`, `auth`), not feature modules. This works for small apps but will become harder to reason about as track, artist, album, playlist, library, search, and upload behavior grows.
- Stores are global. `user-auth` and `player` are both in `src/stores`, which encourages cross-feature imports through `@/stores`.
- Player state is centralized, which is good, but the store owns DOM audio event listeners without a disposal strategy. It also persists full queue objects in localStorage instead of compact IDs/snapshots.
- Upload flow has improved: create track with audio is now integrated. However there is no durable upload/media asset table, so the DB still stores a URL string rather than a media asset reference.
- `npm run build-only` passes, but `npm run type-check` fails because of stale imports, missing type packages, `ImportMeta.env` typing, and auth store storage typing.
- Persian font files exist and are loaded, but design tokens and typography are not feature-consistent. The app mixes English UI copy, Persian locale, and LTR/RTL assumptions.
- Accessibility is partial. Buttons have some labels, but table action buttons, row buttons, dialogs, empty states, and keyboard flows need systematic review.

### Backend

- `cmd/api/main.go` is the active server entrypoint. Root `main.go` is a separate legacy in-memory song server and should not coexist long-term in the same package root.
- `internal/modules/*` is the intended modular architecture. `internal/handler`, `internal/service`, `internal/repository`, and `internal/domain/song.go` are legacy song code.
- Middleware exists in both `internal/common/middleware` and `internal/middleware`, with different Gin/net-http assumptions. This increases configuration drift.
- `Bootstrap` builds a Gin engine and registers routes, then `cmd/api/main.go` builds another Gin engine and registers routes again. Runtime uses the second engine. This is confusing and can hide middleware differences.
- CORS config is loaded from environment but runtime CORS is hardcoded in `cmd/api/main.go`.
- Catalog repository uses transactions for track genre replacement, but not all domain operations are modeled as domain transactions.
- Upload-to-track now does file write then DB create with cleanup on DB failure. This is a practical local-filesystem transaction pattern, but it is not fully durable because uploads are not represented in a DB table.
- Media serving uses `c.File`, not range-aware `http.ServeContent` with cache headers, ETags, or explicit MIME handling. Audio playback may work, but production streaming behavior is weak.
- Auth persists refresh sessions in DB, which is good, but frontend refresh calls do not send `refreshToken`. Tokens are stored in localStorage/cookies accessible to JS.
- JWT secrets have insecure defaults and are not rejected in production config.
- Migrations are minimal and lack new tables for playlists, likes, listening history, uploads/media assets, playback sessions, or search indexes.

### Player

- Good: one Pinia store and one `Audio` element are intended as source of truth.
- Risk: event listeners are anonymous and cannot be removed; if store lifecycle ever resets during HMR/tests/microfrontends, listeners can leak.
- Risk: persisted queue stores whole track objects. This creates stale URLs/titles and localStorage bloat.
- Risk: `element.src !== track.audio_url` compares an absolute resolved URL with a possibly relative or different string representation.
- Risk: shuffle has no history stack, so previous behavior after shuffle is not predictable.
- Risk: route changes are mostly safe because playback is store-owned, but layout unmount does not pause/dispose audio by design. That is okay for a music app, but it must be explicit.

## Prioritized Problems

### 1. Critical Bugs

**Frontend type-check is red.**
Root cause: stale imports, missing type packages, auth store generic storage returning `string | T`, and missing `ImportMeta.env` declarations.
Production impact: regressions can ship despite TypeScript; refactors become unsafe.
Fix strategy: add env/type declarations, fix auth store storage typing, remove stale captcha/maintenance API coupling, add missing dev dependencies or remove stale test config references.

**Auth refresh contract is broken.**
Root cause: backend expects `{ refreshToken }`; frontend `refresh()` posts no body.
Production impact: expired access tokens force logout or loop failed refresh attempts.
Fix strategy: make auth store own refresh token retrieval and pass it to API; handle refresh rotation atomically.

**Root `main.go` and modular `cmd/api/main.go` coexist.**
Root cause: legacy prototype server was not removed or isolated.
Production impact: `go run .` and `go run ./cmd/api` start different apps with different storage and routes.
Fix strategy: move legacy server behind build tag or delete after migration approval; document `cmd/api` as only entrypoint.

### 2. UX Flaws

**Standalone media upload creates files but not tracks.**
Root cause: media upload and catalog track creation are separate workflows.
Production impact: admins think uploads are lost because they do not appear in track tables.
Fix strategy: keep standalone media upload as asset upload, but label it clearly; default track creation should use integrated upload endpoint.

**Track table lacks rich feedback.**
Root cause: simple DataTable with minimal empty/error states.
Production impact: failed loads or no data look broken.
Fix strategy: add empty, error, retry, and row-level loading states.

**Mobile player is functional but dense.**
Root cause: desktop-first bottom bar.
Production impact: controls may be cramped on phones.
Fix strategy: split mini-player and expanded player states inside `features/player`.

### 3. Architecture Flaws

**Frontend is organized by technical type, not domain feature.**
Root cause: global `components`, `composables`, `services`, `stores`.
Production impact: track/admin/upload/player behavior becomes coupled through shared barrels.
Fix strategy: migrate incrementally to `features/player`, `features/uploads`, `features/tracks`, then route-owned features.

**Backend module boundaries are blurred by media-aware catalog handler.**
Root cause: `catalog.Handler` directly depends on `media.Service` for integrated upload.
Production impact: catalog and media modules become harder to test independently.
Fix strategy: introduce an application/use-case layer for `CreateTrackWithAudio` that coordinates media and catalog services.

**Response DTOs double as domain models.**
Root cause: catalog models have `db` and `json` tags together.
Production impact: DB schema changes leak directly to API contracts.
Fix strategy: add transport DTO mapping for public/admin responses before expanding routes.

### 4. State Management Flaws

**Auth store is not strict-safe.**
Root cause: storage helper returns `string | T | null`.
Production impact: runtime shape errors and TypeScript failures.
Fix strategy: typed storage parser per namespace; no generic “maybe string” for structured auth state.

**Player persistence stores full objects.**
Root cause: simplest persistence strategy.
Production impact: stale metadata and storage bloat.
Fix strategy: persist queue item IDs plus current track snapshot; refresh details from API when route/app loads.

### 5. Security Issues

**Upload has no durable DB media asset record.**
Root cause: filesystem storage returns URL metadata only.
Production impact: orphan files, unclear ownership, hard cleanup, weak auditing.
Fix strategy: add `uploads` or `media_assets` table with hash, path, MIME, size, duration, uploader, status.

**Tokens are JS-readable.**
Root cause: localStorage/cookie token persistence.
Production impact: XSS can steal access and refresh tokens.
Fix strategy: short-lived access token in memory, httpOnly secure refresh cookie, CSRF-aware refresh endpoint.

**Production config allows default JWT secrets.**
Root cause: config defaults are development-friendly.
Production impact: deploys can be trivially compromised.
Fix strategy: validate production env and reject default secrets.

### 6. Upload Pipeline Risks

**Filesystem and DB are not truly transactional.**
Root cause: local filesystem cannot join DB transactions.
Production impact: crash between file write and DB write can still orphan files.
Fix strategy: write media asset row with `pending` status, store file, create track in DB transaction, mark asset `attached`; scheduled cleanup for stale pending assets.

**No duration extraction.**
Root cause: upload service only detects MIME/hash/size.
Production impact: admins manually enter duration; player metadata can be wrong.
Fix strategy: use a real audio metadata library or ffprobe integration behind an interface.

**Duplicate detection is filesystem-only.**
Root cause: hash filename is used without DB uniqueness.
Production impact: no way to query duplicate assets, attach count, or ownership.
Fix strategy: unique index on `(sha256, category)` in media assets.

### 7. Playback Bugs

**No cleanup listener API.**
Root cause: anonymous event handlers inside `ensureAudio`.
Production impact: test/HMR leaks; future multi-device sync can double-bind.
Fix strategy: named listeners and `dispose()` action.

**Shuffle previous behavior is weak.**
Root cause: random next without history.
Production impact: previous track does not match user expectation.
Fix strategy: maintain play history stack.

**No retry or fallback on audio error.**
Root cause: error handler only sets message.
Production impact: transient network errors stop playback.
Fix strategy: retry current source once, then optionally skip to next with visible state.

### 8. Performance Bottlenecks

**Large common bundle and admin table chunk.**
Root cause: PrimeVue/DataTable and services are globally imported through pages/barrels.
Production impact: slower initial load.
Fix strategy: route-level feature chunks and narrower imports.

**Search uses `ILIKE '%query%'`.**
Root cause: no search-specific indexing.
Production impact: slow catalog search at scale.
Fix strategy: PostgreSQL trigram indexes first; full-text later if needed.

**Media serving lacks range/cache strategy.**
Root cause: `c.File` convenience route.
Production impact: poor streaming and revalidation behavior.
Fix strategy: use `http.ServeContent`, cache-control, ETag, and verified content type.

### 9. Technical Debt Hotspots

- `frontend/src/stores/user-auth.ts`
- `frontend/src/plugins/client/request-factory.ts`
- `frontend/src/composables/useLoading.ts`
- `internal/modules/catalog/handler.go`
- `internal/modules/catalog/repository.go`
- `cmd/api/main.go`
- root `main.go`
- legacy `internal/handler`, `internal/service`, `internal/repository`, `internal/domain`

## Target Frontend Feature Modules

The target shape is valid, but should be introduced gradually:

```text
src/features/
  auth/
  player/
  tracks/
  albums/
  artists/
  playlists/
  library/
  search/
  uploads/
  admin/
```

Each feature should own its components, composables, stores, services, types, routes, and utils where applicable. Shared UI primitives should remain outside features only when they are truly generic.

Recommended migration order:

1. `features/player`: low API coupling, clear state boundary.
2. `features/uploads`: isolates media API, upload UI, progress/retry.
3. `features/tracks`: owns track rows, track forms, track API types.
4. `features/auth`: fix token contract and storage typing while moving.
5. Route features: search, artists, albums, playlists, library.

## Backend Refactor Direction

Recommended backend module shape:

```text
internal/modules/catalog
internal/modules/media
internal/modules/auth
internal/modules/playlists
internal/modules/library
internal/app/usecase
internal/common/httpx
internal/platform/storage
```

`internal/app/usecase` should coordinate cross-module operations such as `CreateTrackWithAudio`. Catalog should not permanently depend on media internals.

## Database Roadmap

Add migrations in this order:

1. `media_assets`: id, category, original_name, path, public_url, mime_type, size_bytes, sha256, duration_seconds, status, created_by, created_at, attached_at.
2. `track_media`: optional join/reference from track to media asset, or add `tracks.audio_asset_id`.
3. `playlists`: id, owner_id, name, description, visibility, cover_url, created_at, updated_at.
4. `playlist_tracks`: playlist_id, track_id, position, added_by, added_at.
5. `liked_tracks`: user_id, track_id, created_at.
6. `listening_history`: user_id, track_id, played_at, source, duration_ms.
7. `playback_sessions`: user_id, device_id, current_track_id, state, updated_at.

Avoid adding recommendation/AI tables until basic listening history exists.

## Testing Plan

- Backend unit tests for media storage path safety, MIME validation, duplicate hashing, and cleanup.
- Backend handler tests for `POST /admin/catalog/tracks/upload`.
- Repository tests with test Postgres for track create/update and playlist ordering.
- Frontend store tests for player repeat/shuffle/seek/persist.
- API service tests for auth refresh and track upload payload mapping.
- Route tests for auth guard and lazy-loaded domain pages.
- E2E smoke: login, create artist, create track with audio, see it in list, play it.

## First Safe Refactor Slice

Start with `features/player` because it has a clear boundary and low backend coupling.

Scope:

- Move player store and music player components into `src/features/player`.
- Keep compatibility exports from existing `src/stores` and `src/components/music`.
- Add `dispose()` and named audio event cleanup.
- Add store-level tests once Vitest dependencies are fixed.

Tradeoff:

- This introduces a feature module without forcing every existing import to change immediately.
- It avoids a giant migration while creating the pattern for later modules.
