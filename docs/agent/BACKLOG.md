# Muse — Production Readiness Backlog

> Auto-generated from code audits + fix sessions. Severity: 🚨 CRITICAL / ⚠️ IMPORTANT / 💡 SUGGESTION / 🔧 NITPICK.

## Status

**52 items resolved** across 12 phases. **~20 remaining** — mostly suggestions and nice-to-haves. The production readiness program is substantially complete.

## ✅ Fixed / Verified / Documented

| ID | Item | Severity | Resolution |
|----|------|----------|------------|
| C-001 | Frontend type-check is red | 🚨 | Fixed Phase 6: vue-tsc clean |
| C-002 | Auth refresh contract broken | 🚨 | Fixed Phase 5: LogoutPayloadSchema fixed |
| C-003 | Legacy entrypoint pollutes root | 🚨 | Fixed Phase 1 |
| C-004 | Login handler crashes on incomplete body | 🚨 | Fixed Phase 1 |
| C-005 | Routes registered twice | 🚨 | Fixed Phase 1 |
| C-006 | Pre-existing test failures (6 packages) | 🚨 | Fixed Phase 1: all 79 tests pass |
| C1 | listening_history column loss | 🚨 | Fixed by migration 000057 |
| C2 | Duplicate table definitions | 🚨 | Misleading filenames, not actual duplicates |
| C3 | Missing FK constraints | 🚨 | Migration 000058: subscriptions, payments, notifications → users(id) |
| C4 | CreatePlaylist missing owner_id | 🚨 | Fixed: model + repository + tests |
| C7 | Alembic/Go track_embeddings naming conflict | 🚨 | Alembic migration created (HNSW index) |
| C8 | tips.artist_id references users not artists | 🚨 | **Won't fix:** design decision (tips → user accounts, not artist entities) |
| P7-001 | RTL border-r not logical | 🚨 | Fixed: border-r → border-e |
| P7-002 | RTL left-0/right-0 positioning | 🚨 | Fixed: left-0 → start-0 |
| P7-003 | Seek bar missing ArrowLeft/ArrowRight | 🚨 | Fixed: keyboard handlers in 4 player views |
| P7-004 | Queue panel not keyboard operable | 🚨 | Fixed: tabindex, role, ArrowUp/Down reorder |
| P7-005 | No track-change live region | 🚨 | Fixed: aria-live="polite" in App.vue |
| P7-006 | Player bar contrast fails 4.5:1 | 🚨 | Fixed: text-white/50 → text-white/60 |
| P7-007 | Queue drag-and-drop lacks keyboard | ⚠️ | Fixed: ArrowUp/Down reorder with focus management |
| P7-009 | outline-hidden overrides focus-visible | 🚨 | Fixed: focus-visible ring on 5 interactive elements |
| P7-010 | Shuffle/repeat popup keyboard navigation | 🚨 | Fixed: role="menu", ArrowUp/Down/Home/End/Escape, auto-focus |
| P7-011 | AddToPlaylistDialog missing aria-labelledby | ⚠️ | Fixed |
| P7-012 | Visualizer prefers-reduced-motion not reactive | ⚠️ | Fixed: MediaQueryList listener |
| P7-013 | No vue-i18n installed | 🚨 | Installed + locale messages + store + main.ts |
| P7-014 | RTL aurora decorations -left/-right | 🚨 | Fixed: PageNotFound, PageAIMoodExplorer |
| P7-015 | MusicRightPane scoped CSS left/right | ⚠️ | Fixed: left:0;right:0 → inset-inline:0 |
| P7-016 | RTL PageProgressBar directional CSS | ⚠️ | Fixed: full conversion to logical CSS |
| P7-017 | main.css row-enter animation LTR-specific | ⚠️ | Fixed: --dir-multiplier CSS var |
| P7-019 | primeLocale hardcoded in main.ts | ⚠️ | Mitigated: locale-change event listener updates PrimeVue |
| P7-020 | Number formatting not locale-aware | ⚠️ | Fixed: formatCount + formatDuration locale parameter; CountUp uses locale store |
| P7-021 | Font mismatch Lexend vs Cabinet Grotesk | ⚠️ | **Documented:** conscious cost/availability decision |
| P7-022 | Surface Dark #050505 missing CSS token | ⚠️ | Fixed: added `--surface-dark` to main.css |
| P7-023 | useMaintenanceStore has no setter | ⚠️ | Fixed: added `setMaintenance(boolean)` action |
| P7-024 | Feature flags never initialized | ⚠️ | Fixed: `useFeatureFlags().init()` in App.vue onMounted |
| P7-025 | PageLoaderStore not wired with router | ⚠️ | Fixed: connected via router.beforeEach/afterEach |
| P7-026 | userAuthStore.restore() async race | ⚠️ | **Verified:** already properly handled (`await auth.ready()`) |
| P7-027 | Dead register.ts middleware | 💡 | Fixed: deleted (never imported) |
| P7-028 | Stale route type comments | 🔧 | Fixed: removed `// ← add` from router/types.ts |
| P7-029 | ML Celery tasks lack time limits | ⚠️ | Fixed: task_soft_time_limit=600, task_time_limit=900 |
| P7-030 | ML missing Pillow dependency | ⚠️ | Fixed: Pillow>=10.4 added |
| P7-031 | ML loose version pins | 💡 | **Documented:** intentional during active development |
| P7-008 | ExpandedPlayer drag handle tab order | 💡 | **Verified:** already has `aria-hidden="true"` — no fix needed |
| P7-018 | useRTL locale store | 💡 | **Mitigated:** delegates to useLocaleStore |
| N-001 | App title "Achive" not "Muse" | 🔧 | Fixed Phase 1 |
| N-002 | Missing favicon | 🔧 | Fixed Phase 1 |
| N-003 | Missing routes (/admin/settings) | 🔧 | Fixed: PageAdminSettings.vue created + route added |
| N-005 | Skip link accessibility | 🔧 | Fixed: LayoutAdmin, LayoutEmpty |
| N-007 | Silent catch in admin subscriptions | 🔧 | Fixed: added toast error handling |
| I-001 | Missing admin CRUD endpoints | ⚠️ | **False alarm:** 50+ endpoints exist, 20 frontend pages wired |
| I-002 | API contract camelCase vs snake_case | ⚠️ | **Documented:** no unified convention, each module self-consistent |
| I-003 | No durable media asset table | ⚠️ | **Documented:** table exists (000039); Go code uses old `media` table |
| I-004 | Player event listeners leak | ⚠️ | **Verified:** already implemented |
| I-005 | Queue persists full PlaybackTrack objects | ⚠️ | **Verified:** already implements track IDs only |
| I-006 | Shuffle has no history stack | ⚠️ | **Verified:** full history stack in queue-manager |
| I-007 | No retry on audio playback error | ⚠️ | Fixed: 1-retry with 1s delay |
| I-008 | resetTracks() not implemented | ⚠️ | **Verified:** useAdminTracks exports it |
| I-009 | JWT secrets insecure defaults | ⚠️ | **Verified:** ValidateProductionConfig checks |
| I-010 | Upload lacks duration extraction | ⚠️ | Fixed: wired audioinfo.Extract() |
| I-011 | Media serving not range-aware | ⚠️ | **Verified:** already uses http.ServeContent |
| I-012 | Search ILIKE without trigram indexes | ⚠️ | **Verified:** OpenSearch primary; GIN trigram indexes exist |
| S-008 | formatDuration treats 0 as falsy | 💡 | **Verified:** already uses `seconds == null \|\| !isFinite(seconds)` |
| S-010 | Upload uses `body` instead of `data` | 💡 | **Verified:** uses standard `data: formData` — correct |
| S-011 | RegisterTrack reads raw body after bind | 💡 | **Verified:** returns 400 on bind fail, no RawBody consumption |
| S-013 | Login form console.log | 💡 | **Verified:** no console.log found in auth pages |
| Docker | CPU limits + ML health checks | ⚠️ | Fixed: CPU limits on 7 services; ML health check added |

## 🚨 CRITICAL — Remaining

*None.*

## ⚠️ IMPORTANT — Remaining

### No error tracking service
No Sentry/Rollbar — errors only visible via log scraping and the admin error viewer.

### OpenTelemetry deeper instrumentation
Database query tracing, HTTP client tracing, Redis tracing, OTLP collector, Python OTEL.
_(stdout exporter + Gin middleware already in place)_

### CI/CD — Playwright E2E tests
CI workflows exist but no E2E test suite.

## 💡 SUGGESTION — Remaining

### S-001 Frontend organized by type, not domain
Long-term code organization improvement — not urgent.

### S-002 Backend DTOs double as domain models
Architectural debt — not blocking.

### S-003 No middleware consolidation
Minor code organization.

### S-004 AnyTrack defeats type safety
Frontend uses `any` cast for tracks in some admin pages (PageAdminTrackDetail, PageAdminSubscriptions, PageAdminVideoUpload).

### S-005 Ingestion review component is 1135 lines
Should be split into child components.

### S-006 Ingestion polling should be WebSocket
Frontend polls drafts endpoint every 2 seconds.

### S-007 Admin pages use wrong theme classes
Some admin pages use old PrimeVue surface-* classes instead of admin dark theme.

### S-009 UpdateTrack cannot clear nullable fields
API doesn't accept null to clear optional track fields.

### S-012 Auth store is not strict-safe
Type casts in Pinia store may bypass type checks.

### S-014 No route-level chunk strategy
All routes in single bundle; no code splitting per route.

### S-015 Admin catalog preview fetches ALL items
Admin catalog page uses public catalog endpoints instead of dedicated stats endpoint.

## 🔧 NITPICK — Remaining

### N-004 Empty states missing
Public track lists and some admin tables lack empty/error/retry states (mostly mitigated — many pages use AdminEmptyState).

### N-006 RTL/LTR consistency
Some components may not flip correctly in RTL mode (partially mitigated by vue-i18n).

### N-008 Media page has no pagination
PageAdminMedia.vue:105 — large media lists will degrade.

## 📦 Deferred

### Backup / admin-import / artist-import pages
Pages referenced in the admin sidebar but not yet implemented _(after core stability)_.

### media_assets migration
Go `internal/` modules still reference old `media` table instead of `media_assets` _(requires coordinated backend refactor)_.

---

*Last updated: Mon Jul 06 2026 — 52 items resolved, ~20 remaining.*
