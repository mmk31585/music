# Journey 01: Arrival & Auth — The New User's First Screen

## What the Brand-New User Sees

### Step 0: Browser → App (cold start)

```
User types URL or clicks link
  → main.ts mounts Vue app
  → App.vue reads route.meta.layout
  → Pinia stores initialize
     → user-auth store fires restore() ASYNC (reads localStorage → API /me)
     → maintenance store reads flag (sync, from Pinia state, no API)
  → Router resolves `/` → matches `app.home`
  → Guard chain fires:
     1. checkMaintenanceGuard → isMaintenance=false → null (continue)
     2. checkLoginGuard → requiresAuth? No → null (continue)
     3. checkAuthGuard → requiresRole? No → null (continue)
  → LayoutMusicApp renders
  → PageHome loads
```

**CRITICAL TIMING ISSUE**: The guard chain is synchronous but `restore()` in the auth store is async (fire-and-forget, called at store creation). So:
- **Milliseconds 0-50**: User sees `LayoutMusicApp` + `PageHome` as a GUEST (no auth yet)
- **Milliseconds ~200-1500**: `restore()` finishes → if user has a stored token, `userAuth` store updates → `isAuthenticated` becomes `true` → components re-render as authenticated

**But routes with `requiresAuth: true` have a problem**: If the user navigates directly to `/library` during those first ~500ms, the guard sees `isAuthenticated=false` and redirects to login. After restore completes, the user is logged in but already on the login page. The `?redirect=/library` is preserved in the URL but ignored on login (as documented in J02).

### Step 0b: What About `/maintenance`?

The guard `checkMaintenanceGuard` redirects to `{ name: 'maintenance' }` if the maintenance store reports `isMaintenance=true`. But **there is no `/maintenance` route defined** in any of the route files. This means:
- If the store ever has `inMaintenance: true`, the guard redirects to a route name that doesn't exist → **blank screen / Vue-router error**.

### Step 1: First Screen — Home Page (Unauthenticated)

```
LayoutMusicApp renders:
  ┌──────────────────────────────────────┐
  │ MusicSidebar (left)                  │
  │ MusicAppHeader (top)                 │
  │ MAIN:                                │
  │   PageHome loads                     │
  │   ├── Aurora gradient background     │
  │   ├── Loading skeleton (4 sections)  │
  │   └── [after ~500ms-2s]             │
  │       ├── Empty state (v-if !hasData)│
  │       │   "Your soundtrack starts    │
  │       │    right here." + CTAs       │
  │       │   (Explore Music / Discover) │
  │       └── OR content sections        │
  │ MusicRightPane (right)               │
  │ [No NowPlayingBar — no track]        │
  │ MobileBottomNav (mobile)             │
  └──────────────────────────────────────┘
```

The user sees **loading skeletons first**, then **empty state** (because a brand-new user has no history). The empty state is a beautifully designed hero section with "Your soundtrack starts right here" text, gradient backgrounds, and two CTAs: "Explore Music" (→ /search) and "Discover" (→ /discover → redirects to /search).

### Step 2: Clicking Around as Guest

| Action | Result |
|--------|--------|
| Click any track play button | `GuestPlayGate` → allows up to 3 plays → on 4th play, shows upgrade dialog |
| Navigate to `/library` | `checkLoginGuard` → redirects to `/auth/login?redirect=/library` |
| Navigate to `/social` | Same redirect to login |
| Navigate to `/settings` | Same redirect |
| Navigate to `/search`, `/track/:id`, `/album/:id`, `/artist/:id` | ✅ Works — no `requiresAuth` |
| Navigate to `/recommendations` | ✅ Works |
| Click "Like" on a track | `GuestPlayGate` → immediately shows upgrade dialog |

### Step 3: Registration Flow (RegisterForm)

The user clicks "Log in" or "Sign up" from sidebar → `/auth/login` or `/auth/register`.

**Registration form fields** (in order):
1. **Display name** — text input — `@blur` validation? No, **all validation fires on submit only** (`validate()` called at top of `onSubmit`)
2. **Username** — text input (alphanumeric, 3-20 chars, dots/dashes/underscores)
3. **Email** — email input (regex: `^[^\s@]+@[^\s@]+\.[^\s@]+$`)
4. **Password** — PrimeVue `<Password>` component with strength feedback meter (`:feedback="true"`), toggle-mask, 8-char minimum

**Validation timing**: All validation fires **on form submit**, not on blur or on input. Users must click "Create account" before seeing any errors.

**Error handling**:
- **Client-side validation**: Per-field red border + red error text below field, animated with Transition
- **API errors**: Red banner at top of form with `role="alert"`. Backend `VALIDATION_ERROR` code with `details` object is parsed and mapped to individual fields.
- **Duplicate email**: Would return as `VALIDATION_ERROR` with email-specific detail → shown inline
- **Server down**: Caught as generic error → shows "Registration failed" + error message
- **Password feedback**: Strength meter shown live (PrimeVue built-in), but strength rules are not explicitly communicated

**Language**: Entire form is in **English**. Labels: "Display name", "Username", "Email", "Password". No Persian localization.

**After successful form submission**:
1. `useAuth().register()` → calls `POST /api/v1/auth/register`
2. Receives `{ user, access_token, refresh_token }`
3. `store.setSession(...)` → saves to Pinia + localStorage + cookies
4. Router pushes to **`/onboarding/genres`** (genre selection page)
5. Toast appears: "Account created — Welcome to your music app"
6. **User is now auto-logged in** — no separate email verification step

### Step 4: Login Flow (LoginForm)

**Fields**: Email + Password (no "remember me" checkbox)
**Validation**: On submit only — email required, password required
**Error handling**: Per-field inline errors + API error banner (same pattern as register)
**Password toggle**: PrimeVue `toggle-mask` — eye icon to show/hide password

**Wrong password feedback**: Backend returns 401 → `apiError` banner shows the error message. Also, `useAuth().login()` catches the error and shows a **second toast**: "Login failed — [error message]". This is the **double-error bug** (F-010).

**Rate limit (429)**: Not specifically handled — would appear as a generic API error. No cooldown timer shown.

**Forgot password**: **There is no "forgot password" link** on the login form. No route, no component exists.

**"Continue as guest"**: Link at bottom of login form → goes to `/` (home page). Guest session is tracked via localStorage (`moja_guest_plays`) with a limit of 3 track plays.

**"Create one"**: Link to `/auth/register`.

**Login success**:
1. `store.login()` → `POST /api/v1/auth/login`
2. `setSession()` saves credentials
3. `router.push({ name: 'app.home' })` — **always redirects to home, ignoring `?redirect=` query param**
4. Toast: "Welcome back — Logged in successfully"

### Step 5: Session Persistence

**On refresh**: `restore()` is called during Pinia store creation:
1. Reads `localStorage` for token + user
2. If token exists, calls `GET /api/v1/auth/me` to validate
3. If `/me` fails (token expired), calls `$reset()` → user becomes guest
4. **Token refresh is broken** (C-002): when `/me` returns 401, the app should try to refresh the token, but the refresh endpoint contract is broken — frontend sends empty body, backend expects `{ refreshToken }`

**Key finding**: After registration, the user has a valid session. On next visit (days later), if the access token expired but the refresh token is still valid:
- `restore()` finds stored token → calls `/me`
- `/me` returns 401 → should refresh → **refresh fails** → `$reset()` → **user is unexpectedly logged out**

## State Matrix Findings

### Registration Form
| State | Present? | Notes |
|-------|----------|-------|
| 🟢 Empty | ✅ | Clean form with placeholders |
| 🟢 Loading | ✅ | Button spinner + disabled |
| 🔴 Per-field live validation | ❌ | Only on submit — users submit blind, then see all errors at once |
| 🟢 Validation errors | ✅ | Red text, animated, per-field |
| 🟢 API errors | ✅ | Banner + per-field mapping for VALIDATION_ERROR |
| 🔴 Duplicate email UX | Partial | Caught as VALIDATION_ERROR — works but UX could show inline sooner |
| 🔴 Server-down UX | Partial | Generic error — no retry suggestion |
| 🔴 Password rules visibility | ❌ | No visible list of password rules (min length, complexity). Only strength meter. |
| 🔴 Offline | ❌ | No offline detection |
| 🔴 Persian localization | ❌ | Entirely English |

### Login Form
| State | Present? | Notes |
|-------|----------|-------|
| 🟢 Empty | ✅ | |
| 🟢 Loading | ✅ | |
| 🔴 Validation on blur | ❌ | Submit-only |
| 🟢 Wrong password | ✅ | Error message shown (though double-bug causes 2 displays) |
| 🔴 Rate limit UX | ❌ | Not detected or communicated |
| 🔴 Forgot password | ❌ | No link, route, or component |
| 🔴 Remember me | ❌ | Checkbox not present |
| 🔴 Offline | ❌ | No detection |

### Guest Experience
| State | Present? | Notes |
|-------|----------|-------|
| 🟢 Track playing (up to 3) | ✅ | Limited to 3 plays via `GuestPlayGate` + `useGuestSession` |
| 🟢 Upgrade prompt | ✅ | Persian-language dialog with "ثبت‌نام رایگان" and "ورود" |
| 🟢 Guest play counter | ✅ | Persisted to localStorage `moja_guest_plays` |
| 🔴 Play limit visibility | Partial | Only shown in upgrade dialog ("تا حالا X از ۳ آهنگ رایگان شنیدی") — not visible before user hits limit |
| 🔴 Like/save gate | ✅ | Immediately shows upgrade prompt (no free likes) |

### Auth Store Session
| State | Present? | Notes |
|-------|----------|-------|
| 🔴 Token refresh (broken) | ❌ | C-002 — empty body sent, backend expects refreshToken |
| 🟢 restore() on init | ✅ | Async, non-blocking |
| 🔴 Silent session loss | ❌ | F-011 — no notification when session fails |
| 🔴 Maintenance route | ❌ | Route name 'maintenance' doesn't exist in router |

## Friction Points (New This Journey)

| # | Severity | Location | Problem | User Impact |
|---|----------|----------|---------|-------------|
| F-101 | 🚨 BLOCKER | Router: no `/maintenance` route | Guard redirects to non-existent route name 'maintenance' | Blank screen / Vue-router error if maintenance mode is ever enabled |
| F-102 | 🚨 BLOCKER | `stores/user-auth.ts:116-143` + refresh.ts | Token refresh contract broken (C-002) | Users logged out mid-session without explanation; session loss on page refresh |
| F-103 | 🚨 BLOCKER | `reg/useAuth.ts:14` | `?redirect=` query ignored — always navigates to home after login | Users who clicked "Library" get sent home after authenticating |
| F-104 | ⚠️ MAJOR | `LoginForm.vue` | No "forgot password" link or flow | Users who forget their password are stranded |
| F-105 | ⚠️ MAJOR | `useLoginForm.ts:41-54` + `useAuth.ts:26-31` | Error toast fires TWICE on login failure | Redundant, confusing error display |
| F-106 | ⚠️ MAJOR | `RegisterForm.vue` | Validation fires on submit only — no inline/blur validation | Users fill entire form, submit, then see all errors at once |
| F-107 | ⚠️ MAJOR | `LoginForm.vue`, `RegisterForm.vue` | No offline detection on auth pages | Users can't tell if failure is network vs. credentials |
| F-108 | ⚠️ MAJOR | `LoginForm.vue` | No rate-limit UX for 429 responses | Users keep trying while locked out |
| F-109 | ⚠️ MAJOR | `router/index.ts:30-36` | Guard chain synchronous, but auth restore async | Authenticated users flash as guests briefly; redirect race on initial load |
| F-110 | ⚠️ MAJOR | `Router/Auth` | No "remember me" — session persistence is all-or-nothing | Users must re-login every time if browser clears storage |
| F-111 | 💡 IMPROVE | `RegisterForm.vue` | Password rules not explicitly shown — only strength meter | Users guess at requirements |
| F-112 | 💡 IMPROVE | `LoginForm.vue:label` | Labels in English, not localized for Persian | Persian users see English labels |
| F-113 | 💡 IMPROVE | `GuestPlayGate.vue` | No counter visible before hitting limit — 3 free plays is a surprise | Users don't know they have limited plays until prompted |

## RTL / A11y / Mobile Notes

- **No Persian localization** on auth forms — all labels, placeholders, and validation messages are English
- **Guest upgrade prompt** is fully Persian (dialog header, body, buttons)
- **Skip link** missing in LayoutAuth — keyboard users tab through entire hero panel before reaching form
- **Contrast risk**: Labels use `text-white/60` on glass card background — likely below 4.5:1 WCAG AA threshold
- **Touch targets**: All inputs and buttons are adequately sized (44px+)
- **ARIA**: Error messages have `role="alert"`. Inputs have `aria-label`. The password toggle (eye icon) may not have accessible labels.

## Delight Opportunities

- ✨ **Smart redirect**: Use `?redirect=` to send user to original destination after login
- ✨ **Social login**: Google/Spotify OAuth for 1-click registration
- ✨ **Persian-first auth**: Show Persian language on auth forms for Persian users
- ✨ **Password rules visible**: Small info icon/expandable showing password requirements before user starts typing
- ✨ **Guest play counter**: Show "X of 3 free plays remaining" somewhere visible (e.g., small badge in player)
- ✨ **Smooth auth transition**: After login, animate the transition from guest to authenticated state (sidebar updates, content refreshes)

## Open Questions

1. Is email verification required (confirm email link)? There's no verification step after registration — user is immediately logged in.
2. What happens on password change — does the refresh token invalidate? Need backend behavior clarification.
3. Is the `/maintenance` route intentionally missing or a known bug?
4. Should guest plays be per-device or per-IP? Currently per-localStorage (device-specific).
