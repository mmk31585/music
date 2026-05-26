# Production Audit

Date: 2026-05-25

## Architecture Audit

- The backend has two overlapping generations of code: the active Gin modular stack under `internal/modules/*`, and legacy song upload code under `internal/handler`, `internal/service`, `internal/repository`, and `internal/domain`. The legacy song service no longer compiles because it references removed config fields (`AllowedAudioExt`, `MaxUploadSize`, `UploadDir`).
- Routes are registered twice. `Bootstrap` creates a Gin engine and calls `RegisterRoutes`, then `cmd/api/main.go` creates a second Gin engine and calls `RegisterRoutes` again. This makes startup harder to reason about and hides middleware differences between the bootstrap router and runtime router.
- Configuration is partially centralized but not consistently used. CORS origins are loaded into config, while `cmd/api/main.go` hardcodes `http://localhost:3000`.
- The media module is isolated and simple, but it is not integrated with catalog persistence. Uploading audio creates a file only; creating/updating a track later stores a URL string. There is no DB-side media identity, hash, ownership, lifecycle, or cleanup.
- Catalog repositories use transactions for track and genre relationship writes, which is good. Slug uniqueness is still prechecked in application code, so concurrent duplicate creates rely on DB uniqueness errors after the fact.
- Frontend API schemas do not match backend JSON naming. Backend catalog models emit `artistId`, `audioUrl`, `durationSeconds`; frontend schemas expect `artist_id`, `audio_url`, `duration_seconds`. This breaks rendered fields and form payloads.
- Frontend has a reusable Axios wrapper, but some components bypass the typed media API and use a wrong Axios key (`body` instead of `data`).
- The player is currently static UI. There is no audio element/store, queue, persistence, cleanup, or route-stable playback state.
- Authentication persists tokens in localStorage and cookies. The cookie is not `httpOnly`, so it does not materially reduce XSS token exposure. Refresh-token rotation exists server-side, but frontend refresh response normalization is inconsistent with backend snake_case response fields.

## Bug List

- `go test ./...` fails due legacy song service references to removed config fields.
- `npm run type-check` fails due missing dev type packages, stale imports, unsafe storage generics, and upload dialog typing errors.
- `UploadMediaDialog.vue` sends multipart data as `body`; Axios ignores it for this wrapper.
- `TrackFormDialog.vue` renders cover URL controls twice and declares an audio upload dialog state that is never rendered.
- Admin catalog routes expose `PATCH`, while frontend update code uses `PUT`.
- Backend catalog JSON uses camelCase, while frontend schemas and forms use snake_case.
- Track list queries do not join artists/albums, but frontend expects `artist_name` and `album_title`.
- `RegisterTrack` tries to read raw request body after `ShouldBindJSON` failure, which usually consumes the body and produces misleading logs.
- `UpdateTrack` cannot clear nullable `albumId`, `audioUrl`, or `coverUrl` because nil means "not provided" in the current DTO.
- Media serving uses `c.File`, which lacks explicit cache/range strategy for audio streaming.
- Frontend login form has debug `console.log` calls.

## UX Issues

- Upload UI lacks drag-and-drop states, progress, preview, retry, and actionable validation messages.
- Track creation requires manually entering duration, even when an uploaded audio file is available.
- Empty states are missing in public track lists and admin tables.
- Loading states are basic skeleton blocks, not tied to table/form state or retry behavior.
- Player controls are decorative; clicking tracks only logs to console.
- Mobile player and queue/sidebar do not exist yet.
- Admin forms use URL text fields for media, which is error-prone and makes uploaded media feel detached from tracks.

## Security Issues

- Upload validation depends on `http.DetectContentType` and original extension. It does not persist a trusted hash, does not detect duplicates, and does not record content length after copy.
- Uploaded files are public immediately after upload, even before being attached to a catalog record.
- No upload rate limiting or per-user quota.
- No malware scanning hook or quarantine state.
- JWT default secrets are insecure in development config and are not rejected for production.
- CORS is hardcoded at runtime and credentials are allowed.
- Token storage is XSS-sensitive because access and refresh tokens are stored client-side.
- Media route protects basic path traversal, but it should also use an absolute media root and verify resolved paths remain under it.

## Performance Issues

- Catalog list endpoints do a list query plus count query on every request; acceptable now, but needs pagination strategy and index review as data grows.
- `ILIKE '%query%'` cannot use normal btree indexes. Search needs trigram indexes or a dedicated search strategy.
- No ETag/cache-control for media or catalog responses.
- No response compression middleware.
- Vite has only a large chunk warning; no route-level manual chunk strategy yet.
- Public player has no preload strategy or buffered state, so perceived playback startup will be weak.

## Prioritized Roadmap

1. Restore green backend compile and frontend type-check baseline.
2. Normalize frontend/backend catalog contracts to camelCase, or add explicit transport DTOs on backend.
3. Harden media upload as a first slice: single typed frontend path, progress callback, content hash, safe extension from detected MIME, duplicate detection, structured metadata, and cleanup on failed writes.
4. Add a `media_assets` table and repository so uploads have durable IDs, hashes, category, MIME, size, duration, and lifecycle state.
5. Add a transactional "create track with uploaded media asset" service path so DB and filesystem references cannot drift.
6. Replace legacy song module with either build exclusion or a migration into the active media/catalog modules after compile is green.
7. Add a Pinia playback store with one owned `HTMLAudioElement`, event cleanup, queue, repeat/shuffle, persisted volume and progress.
8. Fix auth hydration and token refresh contracts, then move toward httpOnly refresh cookies.
9. Add production middleware: request ID, structured access logs, CORS from config, secure headers, rate limiting, compression, health/readiness.
10. Add migrations for missing indexes, especially trigram search indexes and partial indexes for public track browsing.
