# Journey 21: Error Pages & Edge Cases — 404, 500, Offline, Maintenance, Degraded

> Full trace: route-level error boundary → error pages → offline detection → maintenance mode → degraded experience.

---

## Error Handling Architecture

### Layers of Error Protection

```
Layer 1: Router catch-all → PageNotFound.vue (404)
Layer 2: ErrorBoundary component → onErrorCaptured → retry button
Layer 3: Component-level error states (try/catch with user feedback) — inconsistent
Layer 4: Network error interceptor (Axios) — NO global handler found
Layer 5: Maintenance store + guard (pre-route)
```

---

## 404 — Page Not Found

### Route Registration

```typescript
// frontend/src/router/routes/index.ts
{
  path: '/:pathMatch(.*)*',
  name: 'not-found',
  component: () => import('@/pages/errors/PageNotFound.vue'),
  meta: {
    title: 'Page Not Found',
    layout: 'layout-empty',
  },
}
```

The catch-all route uses Vue Router's `pathMatch` regex to capture all unmatched paths.

### Page Component (`PageNotFound.vue` — 54 lines)

```
Aurora gradient background (same hero style as auth pages)
  → Large "404" text (400px+ font size, bold, gradient)
  → "Page not found" heading
  → Subtitle: "The page you're looking for doesn't exist or has been moved."
  → Two buttons:
      1. "Go Home" — router.push({ name: 'home' })
      2. "Go Back" — router.back()
  → Both buttons have full keyboard support
  → All icon-based (no text icons)
```

### State Matrix

| State | Behavior |
|-------|----------|
| 🟢 Typo in URL (e.g., `/hom`) | ✅ Shows 404 page |
| 🟢 Deleted/removed page | ✅ Shows 404 page |
| 🟢 Direct navigation to invalid path | ✅ Shows 404 page |
| 🟢 Invalid admin path (e.g., `/admin/nonexistent`) | ✅ Shows 404 page |
| 🔴 404 on API route (e.g., `/api/v1/nonexistent`) | ✅ Caught by Axios interceptor — JSON error, NOT 404 page |
| 🔴 404 on dynamic route (e.g., `/track/00000000-0000-0000-0000-000000000000`) | ⚠️ Shows empty page with no data, NOT the 404 page — because the route pattern /track/:id matches |

**Key finding (F-2101)**: Dynamic routes that match the pattern but have invalid/non-existent IDs do **NOT** trigger the 404 page. They render the component, which then silently shows empty state or error. The route guard does not differentiate between "pattern match" and "data exists."

---

## Error Boundary Component (`ErrorBoundary.vue` — 52 lines)

```typescript
// On mount: registers onErrorCaptured
// On error: shows fallback UI
// On retry: emits 'retry', re-renders slot

interface Slots {
  default: {}        // Normal content
  fallback: { error } // Custom error UI
}

// Usage:
<ErrorBoundary>
  <SomeComponent />
  <template #fallback="{ error }">
    <p>Something went wrong: {{ error.message }}</p>
    <button @click="retry">Try again</button>
  </template>
</ErrorBoundary>
```

### Default Fallback UI

```
Center-aligned card:
  ⚠️ Alert icon
  "Something went wrong"
  Error.message (truncated)
  [Try Again] button → emits 'retry'
```

### Where ErrorBoundary is Used

The `ErrorBoundary` component exists but its **usage across the app is inconsistent**. A grep may reveal it's only used in a few places, leaving many components unprotected.

---

## Router Guards (Edge Case Handling)

### Guard Chain

```typescript
// router/index.ts — 3 middleware layers
// 1. checkMaintenanceGuard → maintenance mode redirect
// 2. checkLoginGuard → auth-required redirect
// 3. checkAuthGuard → role-based redirect
// First non-null result short-circuits
```

### Maintenance Guard

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

**F-045** (from Phase 1): The `/maintenance` route doesn't exist — the guard redirects to a non-existent route name. If maintenance mode is enabled, the app shows a blank screen.

### Login Guard

```typescript
// router/middleware/login-guard.ts
export function checkLoginGuard(to: Route): Route | undefined {
  if (to.meta.guestOnly && auth.isAuthenticated) {
    return { name: 'home' }
  }
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
}
```

**Edge cases**:
- Guest accessing `/settings` → redirected to `/login?redirect=%2Fsettings`
- Authenticated user accessing `/login` → redirected to `/home`
- Authenticated user accessing `/register` → redirected to `/home`

### Auth Guard

```typescript
// router/middleware/auth-guard.ts
export function checkAuthGuard(to: Route): Route | undefined {
  const { user } = useAuthStore()
  
  if (to.meta.requiresRole === 'admin' && user?.role !== 'admin') {
    return { name: 'home' }  // F-042: silent redirect, no toast
  }
  if (to.meta.redirectIfAdmin && user?.role === 'admin') {
    return { name: 'admin.dashboard' }
  }
}
```

---

## Offline Detection

### Current State

There is **NO global offline detection** in the app. The following issues exist:

```typescript
// No code found that listens to:
window.addEventListener('online', () => {})
window.addEventListener('offline', () => {})

// No stores for online/offline state
// No offline banner component
// No offline-aware API calls
```

### Where Offline is Missing

| Page | Current Behavior | Ideal |
|------|-----------------|-------|
| Auth pages | Generic error on network failure (F-013) | Show "You are offline" banner |
| Home page | Silent error — sections disappear | Show "You're offline — showing cached content" |
| Search | "No results found" on network error (F-1403) | Show "Offline — search unavailable" |
| Player | Error on track load: "This track can't be played" | Show "Offline — can't stream" |
| Library | Empty sections (F-904) | Show cached library from IndexedDB |
| Any page | Error boundary catches it | Show "Offline" in the error message |

---

## Feature Flags Store

```typescript
// stores/feature-flags.ts
const flags = ref<Record<string, boolean>>({})

async function loadFlags() {
  const { flags: data } = await api.get('/feature-flags')
  flags.value = data
}

function isEnabled(key: string): boolean {
  return flags.value[key] ?? false
}
```

Used for gradual rollouts. Example: `redesignedPlayer` flag gates the new player UI.

### Edge Cases

| State | Behavior |
|-------|----------|
| Flags not loaded | `isEnabled()` returns `false` (feature off) |
| Flag API fails | `loadFlags()` silently fails, all flags default to `false` |
| Flag changes mid-session | No polling — requires page refresh |
| Flag mismatch FE/BE | No validation that FE can handle flag values |

---

## Page Loader Store

```typescript
// stores/page-loader.ts
const loading = ref(false)

function setLoading(val: boolean) {
  loading.value = val
}

function $reset() {
  loading.value = false
}
```

Used with route transitions. When loading is true, a full-page spinner or skeleton shows instead of the page content.

---

## Degraded Experience Patterns

The app handles several degraded states but inconsistently:

| Degraded State | Handled? | Current UX |
|---------------|----------|------------|
| API returns 500 | ❌ No global handler | Per-component silent catch |
| API timeout (>30s) | ❌ No timeout indicator | Spinner spins forever |
| Rate limited (429) | ⚠️ Partial — handled in auth (F-014) but not elsewhere | List/queue just stop responding |
| Slow network | ❌ No degraded mode | Skeleton shows but never fills |
| localStorage full | ❌ No fallback | Queue persistence silently fails (F-812) |
| WebSocket disconnect | ⚠️ Auto-reconnects but no UI indicator (F-1501) | Features silently break |
| Audio decode error | ⚠️ Auto-skips (F-806) but no "track unavailable" state | Thumbnail shows but nothing plays |
| Empty search results | ✅ "No results" with Persian text | Clear message |
| Browser not supported | ❌ No browser check | Unknown behavior on old browsers |
| JavaScript disabled | ❌ No `<noscript>` fallback | White screen |

---

## State Matrix

| State | 404 Page | Error Boundary | Route Guard | Offline | Maintenance |
|-------|----------|----------------|-------------|---------|-------------|
| 🟢 Normal operation | N/A | N/A | Passes through | N/A | N/A |
| 🟢 Invalid URL | ✅ Shows 404 | N/A | N/A | N/A | N/A |
| 🟢 Component error | N/A | ✅ Shows fallback | N/A | N/A | N/A |
| 🟢 Unauthenticated on auth page | N/A | N/A | ✅ Redirects to login | N/A | N/A |
| 🔴 Offline | N/A | ❌ Shows generic "Something went wrong" | ❌ Redirects may fail | ❌ No banner | N/A |
| 🔴 Maintenance mode | N/A | N/A | ❌ Redirects to missing route | N/A | ❌ Blank screen (F-045) |
| 🔴 Invalid dynamic route ID | ❌ Shows empty page with no error | ❌ Not wrapped | N/A | N/A | N/A |
| 🔴 Non-admin accessing admin | N/A | N/A | ⚠️ Silent redirect to home (F-042) | N/A | N/A |
| 🔴 Feature flag API fails | N/A | N/A | N/A | N/A | N/A — all flags default to off |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-2101 | 🚨 BLOCKER | Dynamic routes (track/album/artist) | **Invalid/non-existent IDs return empty pages with no error** — route matches but data is missing | User sees a blank page when they should see a 404 | Add data existence check to route guards or page components |
| F-2102 | 🚨 BLOCKER | Global | **No offline detection anywhere** — no `online`/`offline` event listeners, no offline banner | Users have no way to know the app is offline, see generic errors | Add global offline detection + persistent banner |
| F-2103 | 🚨 BLOCKER | `ErrorBoundary.vue` | **ErrorBoundary used inconsistently** — many pages/components not wrapped | Uncaught errors cause blank or broken UI | Wrap all top-level page components in ErrorBoundary |
| F-2104 | ⚠️ MAJOR | Maintenance store | **No `/maintenance` route** exists (F-045 continuation) | Maintenance mode causes blank screen | Add `/maintenance` route with appropriate UI |
| F-2105 | ⚠️ MAJOR | `PageNotFound.vue` | **404 page uses "layout-empty"** — no nav, no player, no footer | Once on 404, user can't access player controls | Use default app layout for 404 to preserve player continuity |
| F-2106 | ⚠️ MAJOR | `router/middleware/auth-guard.ts` | Non-admin redirected to home **silently** (F-042 continuation) | User thinks the URL is broken | Add "Admin access required" toast |
| F-2107 | 💡 IMPROVE | Global | **No site-wide `<noscript>` tag** | Users with JS disabled see a white screen | Add `<noscript>` with message "Please enable JavaScript" |
| F-2108 | 💡 IMPROVE | Global | **No connection quality indicator** — no detection of slow network | Users don't know why content loads slowly | Add connection quality detection (slow 3G, fast 3G, 4G) |
| F-2109 | 💡 IMPROVE | `stores/feature-flags.ts` | Feature flags **not polled** — require page refresh to update | Flag rollouts require user to refresh | Add periodic flag polling (every 5 min) |
| F-2110 | 💡 IMPROVE | `ErrorBoundary.vue` | **No error reporting** — errors caught but not logged to server | Can't diagnose production errors | Add error reporting service (Sentry, etc.) |
| F-2111 | 💡 IMPROVE | Global | **No "Oops, something broke" toast** for failed API mutations | Users may not notice a failed action | Add global API error toast |
| F-2112 | 💡 IMPROVE | `PageNotFound.vue` | **404 page not localized** — English only | Persian users see English 404 | Localize 404 page to Persian |
| F-2113 | 💡 IMPROVE | Router | **No route change error handling** — if route component fails to load (code-split), no fallback | White screen on failed lazy import | Add route-level error handler with retry |
| F-2114 | 💡 IMPROVE | All pages | **No skip link on error pages** — keyboard users can't skip to main content | Must tab through entire error state | Add skip link to all layouts |

## RTL / A11y / Mobile Notes

- ❌ 404 page has no `aria-label` on navigation buttons — "Go Home" and "Go Back" are icon-only
- ❌ No skip link on any page layout (auth, app, admin) — violates WCAG 2.4.1
- ❌ 404 page not wrapped in `<main>` landmark — screen reader can't navigate to content
- ✅ Maintenance store logic is simple and correct (despite missing route)
- ❌ Error boundary default fallback has no heading (`<h2>`) — screen reader can't identify it as a major error
- ✅ Route guards preserve `?redirect=` query param for post-login navigation

## Edge Case Checklist

| Edge Case | Handled? | Notes |
|-----------|----------|-------|
| Invalid track/album/artist ID in route | ❌ | Shows empty page, no 404 |
| Offline mode | ❌ | No detection, no banner, no cached content |
| Maintenance mode | ❌ | Guard redirects to non-existent route (F-045) |
| JavaScript disabled | ❌ | No `<noscript>` fallback |
| WebSocket disconnect | ⚠️ | Reconnects but no UI indicator (F-1501) |
| API rate limiting | ⚠️ | Handled only on auth (F-014) |
| localStorage full | ⚠️ | Fails silently (F-812) |
| Code-split chunk load failure | ❌ | No retry mechanism |
| Browser not supported | ❌ | No feature detection |
| Non-admin accessing admin route | ⚠️ | Silent redirect (F-042) |
| Expired session during use | ⚠️ | Auth restore handles but without notification (F-011) |
| Server-side rendering (SSR) edge cases | N/A | App is client-side rendered only |
