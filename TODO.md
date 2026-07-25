# Muse — Project TODO

> Last updated: 2026-07-09 (FINAL)
> Sources: inline code TODOs, architecture audits, live exploration audit, agent fix sessions

---

## Final Status: ✅ ALL MAJOR ITEMS COMPLETE

Build clean: `go build` ✅ `go vet` ✅ `go test` 24/24 ✅  
All CRITICAL and HIGH items resolved. Only MEDIUM/LOW remain.

---

## 🟡 MEDIUM (Optional)

### M-001 Frontend organized by type, not domain
Long-term migration to `features/player`, `features/uploads`, etc. Architectural.

### M-002 Backend DTOs double as domain models
Add transport DTO mapping layer. Architectural debt.

### M-003 No middleware consolidation
`internal/common/middleware/` and `internal/middleware/` coexist. CORS hardcoded in `cmd/api/main.go`.

### M-006 Ingestion polling → WebSocket/SSE
Now using `useEnrichmentPolling` composable (cleaner polling). Full SSE backend would need pub/sub.

### M-016 Missing database tables (still TODO)
- `playlists` + `playlist_tracks` — user playlists
- `liked_tracks` — user likes  
- `listening_history` — time-partitioned playback log
- `playback_sessions` — cross-device state sync

*(DBA confirmed all these migrations already exist in the repo)*

---

## 🔵 LOW

### L-001 `VersionedEvent` / `VersionedBaseEvent` unused
`internal/platform/events/version.go` — clean up or implement.

### L-002 `RepositoryInterface` unused in auth
`internal/modules/auth/interfaces.go` — never used as injection type.

### L-003 `RequireRole` exact-match only
No role hierarchy. Extend when multi-role support needed.

### L-004 Album upsert doesn't fall back to embedded cover
`internal/modules/ingestion/finalization/service.go:118`

### L-005 Empty states missing
Some admin tables lack empty/error/retry states (partially mitigated).

### L-006 RTL/LTR consistency
Some components may not flip correctly (partially mitigated by vue-i18n + logical CSS).

### L-007 Accessibility — remaining gaps
Buttons lack ARIA labels in some admin tables. Skip-link exists but not universally tested.

---

## 📦 Database Migration Roadmap (Already Migrated)
- ✅ `media_assets` (000039 + 000055)
- ✅ `listening_history` fix (000057)
- ✅ FK constraints (000058)
- ✅ `playlists` + `playlist_tracks` (000004 + 000051)
- ✅ `liked_tracks` (000005)
- ✅ `playback_sessions` (000041 + 000054)

---

## 🧪 Testing Remaining
- E2E smoke: Playwright test for login → create → play flow
- Frontend store tests (Vitest)
- Backend handler tests for upload, admin CRUD

---

## ✅ Completed In This Session
| Item | Status |
|------|--------|
| Frontend type-check fixed | ✅ |
| Auth refresh contract fixed | ✅ |
| Login handler validation fix | ✅ |
| Routes double-registration fixed | ✅ |
| Media assets table wired (`media_assets`) | ✅ |
| FullscreenPlayer built (1164 lines) | ✅ |
| Admin analytics + cache pages | ✅ |
| Admin analytics backend endpoint | ✅ |
| Admin cache backend endpoint | ✅ |
| OpenTelemetry OTLP dep added | ✅ |
| Redis OTEL tracing hook | ✅ |
| Sentry error tracking wired | ✅ (already done) |
| AnyTrack type safety fixed (3 pages) | ✅ |
| Ingestion review split (933→888 lines) | ✅ |
| Ingestion polling composable | ✅ |
| UpdateTrack null clearing (ClearedFields) | ✅ (already implemented) |
| Admin catalog stats endpoint | ✅ (already exists) |
| CI/CD workflows audited & fixed | ✅ |
| DB migrations verified (all exist) | ✅ |
| All docs updated | ✅ |
