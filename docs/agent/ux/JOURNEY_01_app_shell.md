# Journey 01: App Shell & Routing

## Entry Points
- **Any URL** → `router/index.ts` → `beforeEach` guard chain → Layout + Page rendering
- **Deep links** via route params (`/track/:id`, `/album/:id`, `/artist/:id`, `/playlist/:id`, etc.)

## Step-by-Step Walkthrough (as implemented)

### Step 1: App Initialization (`main.ts` → `App.vue`)
| Aspect | Detail |
|--------|--------|
| **User Action** | Navigate to any URL |
| **System Response** | `main.ts` mounts app → `App.vue` reads `route.meta.layout` → renders dynamic layout component |
| **API Calls** | None at mount (auth restore happens via Pinia store async init) |
| **States Handled** | No explicit first-load spinner at App level. Auth store `restore()` fires async on store creation but doesn't block rendering |
| **Files** | `frontend/src/main.ts`, `frontend/src/App.vue` (lines 10-27) |

**Key finding**: `App.vue` mounts a `Toast` component globally and a `PageProgressBar`. The layout is resolved as `route.meta.layout ?? 'layout-empty'`. Auth restore is async — user may see unauthenticated content briefly before `restore()` completes.

### Step 2: Guard Chain (`router/beforeEach`)
| Guard | Logic | What Happens |
|-------|-------|-------------|
| `checkMaintenanceGuard` | Reads `useMaintenanceStore().isMaintenance` | Redirects to `/maintenance` if true (guards against loop) |
| `checkLoginGuard` | Reads `useUserAuthStore()` | Redirects authenticated users *away* from guest pages; redirects unauthenticated users *to* login from auth-required pages |
| `checkAuthGuard` | Reads `useUserAuthStore()` role | Admin-only routes redirect non-admins to home; "redirectIfAdmin" routes redirect admins to `/admin` |

**Files**: `router/middleware/maintenance-guard.ts`, `login-guard.ts`, `auth-guard.ts`

**Critical finding**: The guard chain runs **synchronously** — it reads Pinia stores directly. If `restore()` hasn't completed (async), `isAuthenticated` may be `false` temporarily, causing:
- Authenticated users landing on a `requiresAuth` page flash-redirected to login before restore completes
- Admin users briefly seeing non-admin UI before redirect

### Step 3: Layout Resolution

| Layout | Route Scope | Structure | Mobile? |
|--------|-------------|-----------|---------|
| `layout-empty` (no component, `RouterView` only) | `/onboarding`, 404, `/maintenance` | Raw `router-view` | N/A |
| `LayoutAuth` | `/auth/login`, `/auth/register` | Single centered form | Responsive (hidden hero panel on mobile) |
| `LayoutMusicApp` | All `/` (app) routes | Sidebar + Header + Main + RightPane + PlayerRegion + MobileBottomNav | Responsive with mobile nav overlay |
| `LayoutAdmin` | All `/admin` routes | AdminSidebar + AdminTopbar + Main + NowPlayingBar + Fullscreen/ExpandedPlayer + QueuePanel | Responsive with sliding mobile drawer |

**Key UX properties**:
- `LayoutMusicApp`: Has offline banner, search overlay, radio mode, player region, mobile bottom nav. Uses `KeepAlive` (max 3) for page transitions.
- `LayoutAdmin`: Uses feature flag `redesignedPlayer` to switch between `FullscreenPlayer` and `ExpandedPlayer`.
- Both music layouts use `provide('openRadio', ...)` for cross-component radio launching.

### Step 4: Route Groups

**Auth Routes** (guestOnly):
- `/auth/login` → `PageLogin` (LoginForm + AuthHeroPanel)
- `/auth/register` → `PageRegister`

**App Routes** (under LayoutMusicApp):
- `/` (home), `/search`, `/discover` (→ redirect to search)
- `/recommendations/*` (hub + 4 sub-pages: for-you, popular, best, recent)
- `/library`, `/playlists` (→ redirect to library), `/playlist/:id`
- `/track/:id`, `/album/:id`, `/artist/:id`
- `/explore`, `/create-edit`, `/notifications`
- `/social/*` (hub, parties/:id, rooms/:id, clubs/browse, clubs/:id)
- `/profile`, `/user/:id?`, `/settings`
- `/subscription`, `/contributions`, `/creator-dashboard`, `/gamification`
- `/stats`, `/ai/mood-explorer`, `/ai/playlist-generator`
- `/videos`, `/music-video/:videoId`

**Admin Routes** (requiresAuth + requiresRole:admin):
- `/` (dashboard), `/catalog`, `/tracks`, `/tracks/:id`
- `/artists`, `/artists/:id`, `/albums`, `/albums/:id`, `/genres`
- `/users`, `/media`, `/videos`, `/video-upload`
- `/import`, `/import/artist`, `/ingestion`, `/ingestion/review/:id`
- `/moderation`, `/subscriptions`, `/contributions`

**Missing routes** as identified in BACKLOG.md:
- `/admin/settings` — 404
- `/admin/artists/import` — should be `/admin/import/artist` (mismatch)

### Step 5: Page Rendering

Inside `LayoutMusicApp`:
```html
<RouterView v-slot="{ Component }">
  <Transition name="page" mode="out-in">
    <KeepAlive :max="3">
      <component :is="Component" />
    </KeepAlive>
  </Transition>
</RouterView>
```

**Key findings**:
- `KeepAlive` max=3 — only 3 pages cached. Navigation beyond 3 destroys/recreates components
- Page transition uses CSS class `page` (enter/leave animation)
- Each page manages its own data fetching in `onMounted` — no standard data-loading pattern

## State Matrix Findings

### App.vue (Root)
| State | Present? | Notes |
|-------|----------|-------|
| 🟢 First load | Partial | No global skeleton — layout renders immediately, individual pages handle their own loading |
| 🟢 Empty | N/A | App-level doesn't have an empty state |
| 🟡 Loading | Partial | `PageProgressBar` shows top bar progress, but no global page-load spinner |
| 🔴 Error | Missing | No global error boundary — individual pages may or may not handle errors |
| 🔴 Offline | Partial | `LayoutMusicApp` has an offline banner, but `LayoutAdmin` and `LayoutAuth` don't |

### Guard Chain
| State | Present? | Notes |
|-------|----------|-------|
| 🔴 Auth flash | Missing | No guard for the brief window when auth restore hasn't completed. Users may see a flash of login before redirect |
| 🟢 Maintenance | Present | `checkMaintenanceGuard` properly prevents loop |
| 🟢 Admin role | Present | `checkAuthGuard` redirects non-admin users |
| 🔴 Guest access | Partial | `/` has `requiresAuth: false` which means unauthenticated users can browse home/search but hit auth wall on library/social/etc |
| 🔴 Auth-required redirect | Present but fragile | Redirect query `?redirect=/path` exists but no guarantee the target page re-renders correctly after login |

## Friction Points

| # | Severity | Location | Problem | User Impact |
|---|----------|----------|---------|-------------|
| 01.01 | ⚠️ MAJOR | `router/index.ts:30-36` | Guard chain runs synchronously — auth restore is async. Authenticated users may be briefly seen as unauthenticated | Flash of login page → redirect to intended page; disorienting |
| 01.02 | 💡 IMPROVE | `App.vue:24` | Layout switch via `component :is` — no transition between layouts | Abrupt visual switch when going auth → app or app → admin |
| 01.03 | 💡 IMPROVE | `router/index.ts` | No global error boundary for route-level errors | If a page component fails to load (code-split chunk), user gets blank white screen or Vue error overlay in dev |
| 01.04 | 💡 IMPROVE | `app.ts` (routes) | `/discover` and `/playlists` and `/recently-played` are hard redirects to `/search` or `/library` | User sees unexpected page — no toast explaining the redirect |
| 01.05 | 💡 IMPROVE | All layouts | No `dir` attribute set on `<html>` element by default — RTL only set on inner containers | Screen readers may not correctly detect Persian language at document level |
| 01.06 | 💡 IMPROVE | `LayoutAuth.vue` | No skip link, no offline banner | Keyboard users can't skip to content; no offline notification on auth pages |

## Clicks & Time-to-Value

| Metric | Current | Ideal |
|--------|---------|-------|
| Navigate to any protected page (authenticated) | ~200ms (after restore) | <100ms |
| Layout switch (auth → app) | Instant (no transition) | Smooth cross-fade |
| Page transition | ~160ms (CSS transition) | <100ms perceived |

## RTL / A11y / Mobile Notes

- **RTL**: `LayoutMusicApp` uses `:dir="rtlDir"` from `useRTL()`. `LayoutAdmin` uses `:dir="dir"`. Both apply at the top-level div — but document-level `<html>` dir is not set.
- **Skip link**: Present in `LayoutMusicApp` (`#main-content`), absent in `LayoutAuth` and `LayoutAdmin`.
- **Mobile**: `LayoutMusicApp` has a mobile sidebar overlay (w-80 max-w-[85vw]) with backdrop and keyboard dismiss. MobileBottomNav is rendered globally.
- **Touch targets**: Sidebar items use adequate padding (py-3). Close buttons are 40×40.
- **Focus management**: Mobile overlay uses `tabindex="0"` on backdrop — keyboard trap is implemented manually. No auto-focus on sidebar open.

## Delight Opportunities

- ✨ **Smart auth redirect**: After login, restore user scroll position and playing state from before the redirect
- ✨ **Layout transition**: Animate layout switching with a subtle fade/scale to signal context change
- ✨ **Persistent state across routing**: Use `KeepAlive` more strategically (e.g., always keep player state alive regardless of page)

## Open Questions

1. What is the expected behavior when auth restore fails? Currently it calls `$reset()` and the user becomes a guest — should there be a retry or notification?
2. Is the `/maintenance` route page implemented? The guard redirects there but I haven't found the page component.
3. Should admin routes also be accessible to super-admin or only exact "admin" role?
