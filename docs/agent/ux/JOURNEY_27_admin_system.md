# Journey 27: Admin — System Administration

> Full trace: system health → configuration management → maintenance mode → feature flags → cache.

---

## System Health Infrastructure

### Current State

```go
// internal/modules/health/routes.go
// GET /health         → { status: "ok" }
// GET /health/live    → liveness probe (always 200)
// GET /health/ready   → readiness probe (checks DB + Redis connectivity)
```

The health endpoints are basic Kubernetes-style probes. There is **no admin UI for system health** and **no detailed health dashboard**.

### What Exists

| Endpoint | Purpose | Admin UI? |
|----------|---------|-----------|
| `GET /health` | Simple OK status | ❌ No |
| `GET /health/live` | Liveness (always OK) | ❌ No |
| `GET /health/ready` | Readiness (DB + Redis) | ❌ No |
| `GET /api/v1/metrics` | Prometheus metrics endpoint | ❌ No (Prometheus only) |

### What's Missing

| Feature | Status | Impact |
|---------|--------|--------|
| **Detailed health dashboard** | ❌ Missing | No way to see DB, Redis, storage, worker status at a glance |
| **Admin health page** | ❌ Missing | No `/admin/health` route or page |
| **Service dependency status** | ❌ Missing | Can't see if ML service, WebSocket, or storage is healthy |
| **Worker status monitoring** | ❌ Missing | Can't see if workers are running, their last run time, or error rate |
| **Uptime / version info** | ❌ Missing | No build version, uptime, or deployment info visible |
| **Alert configuration** | ❌ Missing | No way to configure health check thresholds |

---

## Maintenance Mode

### Frontend Store

```typescript
// stores/maintenance.ts (22 lines)
const inMaintenance = ref(false)
const checked = ref(false)

// Composed into:
const isMaintenance = computed(() => inMaintenance.value)
const isChecked = computed(() => checked.value)
```

### Router Guard

```typescript
// router/middleware/maintenance-guard.ts
export function checkMaintenanceGuard(to: Route): Route | undefined {
  const { isMaintenance, isChecked } = useMaintenanceStore()

  if (!isChecked) return undefined // Not yet checked — allow
  if (isMaintenance && to.name !== 'maintenance') {
    return { name: 'maintenance' }
  }
  if (!isMaintenance && to.name === 'maintenance') {
    return { name: 'home' }
  }
}
```

### Gaps

| Issue | Severity | Detail |
|-------|----------|--------|
| No `/maintenance` route | 🚨 F-045/F-2104 | Guard redirects to non-existent route → blank screen |
| No API to toggle maintenance mode | ❌ | Store value is hardcoded/static — no `POST /admin/maintenance` |
| No admin UI toggle | ❌ | Can't enable/disable maintenance from the admin panel |
| No scheduled maintenance | ❌ | Can't set "maintenance starts at midnight" |

---

## Feature Flags

### Frontend Store

```typescript
// stores/feature-flags.ts (36 lines)
const flags = ref<Record<string, boolean>>({})

async function loadFlags() {
  const { data } = await api.get('/features')
  flags.value = data.flags
}

function isEnabled(key: string): boolean {
  return flags.value[key] ?? false
}
```

### Backend

```go
// config/features.go
// Feature flags defined as environment variables / config
// Loaded at startup, no runtime API to toggle
```

### Gaps

| Issue | Severity | Detail |
|-------|----------|--------|
| **No admin UI for flags** | ⚠️ | Must restart server with env changes |
| **No runtime toggle API** | ⚠️ | Flags are read-only once loaded |
| **No gradual rollout** | ⚠️ | No percentage-based rollouts |
| **No A/B testing framework** | ❌ | Can't run experiments |
| **No flag audit log** | ❌ | No record of who changed what flag |

---

## Platform Configuration

### Current State

There is **no admin settings page** (`/admin/settings`) — confirmed by `TODO.md #45`:

> "TODO #45 (LOW): `/admin/settings` route is missing + `/admin/artists/import` path mismatch"

### What an Admin Settings Page Should Include

| Setting Category | Examples | Currently Configurable? |
|-----------------|---------|------------------------|
| **Platform** | Site name, description, contact email | ❌ No admin UI |
| **Upload limits** | Max file sizes, allowed MIME types | Only via config files |
| **Authentication** | OAuth providers, registration open/closed | ❌ No admin UI |
| **Email** | SMTP config, from address, templates | ❌ No admin UI |
| **Storage** | Provider, bucket, CDN URL | Only via env variables |
| **Cache** | Redis TTL, invalidation | Only via env variables |
| **Rate limits** | API rate limits per endpoint | Only via config files |
| **Content moderation** | Auto-flag rules, profanity filter | ❌ No admin UI |
| **Regional settings** | Default language, currency, timezone | ❌ No admin UI |

---

## Cache Management

### Current State

Redis is used for caching but there is **no admin UI** to manage it:
- No cache invalidation button
- No cache hit/miss rate display
- No key browser
- No TTL management

---

## State Matrix

| State | Health Dashboard | Maintenance Mode | Feature Flags | Admin Settings |
|-------|-----------------|-----------------|---------------|----------------|
| 🟢 Page exists | ❌ Nonexistent | ❌ Route nonexistent | ❌ No page | ❌ Route nonexistent |
| 🔴 All healthy | N/A | N/A | N/A | N/A |
| 🔴 DB down | ❌ No indicator | ❌ Can't trigger | N/A | N/A |
| 🔴 Worker failure | ❌ No indicator | N/A | N/A | N/A |
| 🔴 Maintenance needed | N/A | ❌ Can't toggle | N/A | N/A |
| 🔴 Flag rollout needed | N/A | N/A | ❌ Can't change | N/A |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-2701 | 🚨 BLOCKER | `stores/maintenance.ts` | **Maintenance mode cannot be toggled via API** — store value is static | Operations can't enable maintenance mode remotely | Add `POST /admin/maintenance` API endpoint |
| F-2702 | ⚠️ MAJOR | `stores/feature-flags.ts` | **Feature flags have no admin UI** — must restart server to toggle | Can't do gradual rollouts or quick kill-switches | Add admin feature flag management page |
| F-2703 | ⚠️ MAJOR | All admin | **No `/admin/settings` page** exists (TODO #45) | No centralized platform configuration | Add admin settings page with categorized sections |
| F-2704 | ⚠️ MAJOR | `internal/modules/health/` | **No detailed health dashboard** — only basic liveness probes | Can't see DB/Redis/storage/worker status at a glance | Add admin health page with per-service status indicators |
| F-2705 | 💡 IMPROVE | System health | **No worker monitoring** — can't see if background jobs are running | Silent worker failures go unnoticed | Add worker status page with last-run, duration, error count |
| F-2706 | 💡 IMPROVE | System health | **No uptime / version display** anywhere in admin | Can't verify which build is deployed | Add build version, git hash, deploy time to admin footer |
| F-2707 | 💡 IMPROVE | Cache management | **No cache management UI** — can't invalidate or inspect cache | Stale data persists until TTL expires | Add admin cache page (invalidation, hit rate, keys) |
| F-2708 | 💡 IMPROVE | Feature flags | **No audit log for flag changes** | Can't track who changed what flag | Add flag change audit trail |
| F-2709 | 💡 IMPROVE | Maintenance mode | **No scheduled maintenance** — can't set "starts at 2 AM" | Maintenance always requires immediate action | Add scheduled maintenance with countdown banner |
| F-2710 | 💡 IMPROVE | Maintenance mode | **No maintenance page UI** (F-045/F-2104 continuation) | Users see blank screen during maintenance | Build proper maintenance page with countdown + status |
| F-2711 | 💡 IMPROVE | System health | **No log viewer** — can't browse server logs from admin | Must SSH into server to debug | Add log viewer with level filter and tail mode |
| F-2712 | 💡 IMPROVE | System health | **No alert configuration** — can't set thresholds for health checks | Silent degradation until outage | Add alert threshold config (e.g., warn when queue > 1000) |

## RTL / A11y / Mobile Notes

- ❌ Maintenance mode page will need full RTL support for Persian users
- ✅ Feature flag store is simple and well-structured (despite missing UI)
- ❌ Admin settings page would need proper form accessibility with `aria-describedby` for each setting
