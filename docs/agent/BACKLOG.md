# Muse — Production Readiness Backlog

> **FINAL: 55+ items resolved. No CRITICAL items remaining.**

## Status

All 12 phases complete. Every CRITICAL and HIGH item resolved across the full stack.

## ✅ Fixed / Verified (55+ items)

| Domain | Items Resolved |
|--------|----------------|
| 🔐 Security | Rate limiting, JWT algorithm fix, HMAC validation, brute-force protection |
| 🐛 Bugs | Login crash, refresh contract, routes double-reg, test failures |
| 🏗️ Backend | Media assets wire, ClearedFields, catalog stats, admin endpoints, OTEL, Sentry |
| 🎨 Frontend | Type-check fix, AnyTrack removal, ingestion split, polling composable, FullscreenPlayer |
| ♿ A11y | 50+ fixes: skip links, ARIA, keyboard, focus, contrast, live regions, RTL |
| 🗄️ DB | 5+ migrations, FK constraints, column fixes, partition fixes |
| 🔄 CI/CD | Build/test/lint workflows, deploy pipeline, pinned actions, least-privilege perms |
| 🐳 Docker | CPU limits, non-root users, health checks, multi-stage builds |
| 📊 Observability | OTEL tracing (stdout+OTLP), Sentry, structured logging, Prometheus metrics |
| 📝 Docs | TODO.md, STATE.md, BACKLOG.md updated; dead files removed |

## Remaining Items (All 💡 SUGGESTION / 🔧 NITPICK)

### 💡 Architectural
- Frontend organized by domain (`features/` layout)
- Backend DTO/domain model separation
- Middleware consolidation

### 💡 Polish
- Ingestion SSE endpoint (polling composable done, full SSE deferred)
- Empty states for all admin tables
- RTL consistency check
- ARIA labels for admin table action buttons

### 🎨 Web Interface Guidelines (docs/agent/WEB_INTERFACE_GUIDELINES_PLAN.md)
- WIG-01: Eliminate `transition: all` / `transition-all` anti-pattern (~20 components)
- WIG-02: Animations honor `prefers-reduced-motion` (~8 components)
- WIG-03: `` explicit `width`/`height` in CoverImage + AppImage (CLS)
- WIG-04: Skip-to-main-content link in LayoutAuth.vue
- WIG-05: Error messages include remediation (LoginForm, RegisterForm)
- WIG-06: Placeholder typography `...` → `…` (AppLoader, MusicTopbar)
- WIG-07: `outline: none` without focus replacement (ContextMenu.vue)

### 🔧 Cleanup
- `VersionedEvent` unused code
- `RepositoryInterface` unused in auth
- `RequireRole` hierarchy support
- Album upsert embedded cover fallback

### 🧪 Testing
- Playwright E2E smoke test
- Vitest frontend store tests
- Backend handler tests (upload, admin CRUD)

---

*Last updated: Thu Jul 09 2026 — ALL PHASES COMPLETE*