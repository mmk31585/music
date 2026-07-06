# Journey 02: Auth & Login

## Entry Points
- **Direct**: `/auth/login`, `/auth/register`
- **Redirected**: Any `requiresAuth` page when unauthenticated → `/auth/login?redirect=/original-path`
- **Guard trigger**: `checkLoginGuard` redirects unauthenticated users to login
- **Navigation**: "Log in" button in sidebar/footer, "Log out" → redirect to login

## Step-by-Step Walkthrough (as implemented)

### Step 1: Guard Redirect → Login Page
| Aspect | Detail |
|--------|--------|
| **User Action** | Clicks "Library" (requiresAuth) while logged out |
| **System Response** | `checkLoginGuard` → redirects to `/auth/login?redirect=/library` |
| **Layout** | `LayoutAuth.vue` renders — minimal shell with `dir` attribute and `bg-surface-raised` |
| **API Calls** | None yet |
| **States Handled** | No loading state during redirect; page appears instantly |
| **Files** | `router/middleware/login-guard.ts:21-25`, `layouts/LayoutAuth.vue` |

### Step 2: Login Page Render (`/auth/login`)
| Aspect | Detail |
|--------|--------|
| **Screen** | Desktop: Split layout — hero panel (left, hidden on mobile) + login form (right) |
| **User Action** | Views the login form with email + password fields, "Continue as guest" link, "Create one" link |
| **System Response** | Form renders immediately — no loading state needed |
| **States Handled** | 🟢 Success: form ready to use. 🔴 No offline detection. 🔴 No "forgot password" link |
| **Files** | `pages/auth/PageLogin.vue`, `components/auth/LoginForm.vue`, `components/auth/AuthHeroPanel.vue` |

### Step 3: Form Validation & Submission
| Aspect | Detail |
|--------|--------|
| **User Action** | Types email & password → clicks "Log in" |
| **System Response** | Client-side validation (email required, password required). On success, calls `useAuth().login()` |
| **States** | 🟢 Loading: button shows spinner (`:loading="loading"`). 🟢 Validation errors: per-field red text with Transition animation. 🟢 API error: banner with red background + icon. |
| **API Call** | `POST /api/v1/auth/login` with `{ email, password }` |
| **Error handling** | 🚨 **Backend bug**: empty body `{}` returns 500 instead of 400 (C-004 from BACKLOG.md) |
| **Files** | `composables/auth/useLoginForm.ts`, `composables/auth/useAuth.ts:12-33` |

**Critical finding — login toast double-fire**: `useLoginForm.ts` catches errors and sets `apiError` for the in-form banner. But `useAuth.ts` ALSO catches the same error and shows a toast (`'Login failed'` + detail). This means on login failure, the user sees **both** an error banner in the form AND a toast notification. This is redundant and visually noisy.

### Step 4: Login Success → Redirect
| Aspect | Detail |
|--------|--------|
| **System Response** | `store.login()` → sets session (tokens + user) in Pinia + localStorage + cookies. Then `router.push({ name: 'app.home' })` |
| **API Call** | After login, `restore()` fires `me()` to validate token |
| **Feedback** | Toast: "Welcome back — Logged in successfully" |
| **States** | 🔴 **Missing**: The `?redirect=/library` query param is preserved but `useAuth().login()` always redirects to `'app.home'`, NOT the original destination |
| **Files** | `stores/user-auth.ts:116-143`, `composables/auth/useAuth.ts:14` |

**🚨 CRITICAL finding — redirect ignored**: The guard sets `?redirect=/library` but `useAuth().login()` calls `router.push({ name: 'app.home' })` unconditionally. Users are always sent to home, not the page they originally wanted.

### Step 5: Registration Flow
| Aspect | Detail |
|--------|--------|
| **User Action** | Clicks "Create one" → `/auth/register` |
| **Form fields** | Display name, Username, Email, Password (with feedback strength meter) |
| **Validation** | Name 2+ chars, Username 3-20 chars alphanumeric, Email regex, Password 8+ chars |
| **Success** | `useAuth().register()` → stores session → redirects to `/onboarding/genres` |
| **Backend errors** | Parsed from `VALIDATION_ERROR` code with per-field details mapped to form fields |
| **Files** | `components/auth/RegisterForm.vue`, `composables/auth/useRegisterForm.ts` |

**Key**: Registration redirects to the genre onboarding page (`/onboarding/genres`) — this is entirely in Persian and requires selecting at least 3 genres.

### Step 6: Session Persistence (`restore()`)
| Aspect | Detail |
|--------|--------|
| **When** | On Pinia store creation (app init) — async, non-blocking |
| **What** | Reads `localStorage` for token + user → validates via `me()` API call |
| **On failure** | Silently resets state via `$reset()` — user becomes guest |
| **Warning** | 🚨 **C-002**: Refresh endpoint sends empty body — refresh token flow is broken |
| **Files** | `stores/user-auth.ts:183-216` |

### Step 7: Logout
| Aspect | Detail |
|--------|--------|
| **Trigger** | Logout button in sidebar or mobile nav |
| **System** | Calls `store.logout()` → API call to invalidate refresh token → `$reset()` → redirect to login |
| **Fallback** | If API call fails, still `$reset()` and redirect |
| **Files** | `composables/auth/useAuth.ts:73-80`, `stores/user-auth.ts:155-168` |

## State Matrix Findings

### LoginForm
| State | Present? | Details |
|-------|----------|---------|
| 🟢 Empty | ✅ | Clean form with placeholders |
| 🟢 Loading | ✅ | Button spinner, disabled during submit |
| 🟢 Validation error | ✅ | Per-field red text with animated transitions |
| 🟢 API error | ✅ | Red banner with icon, animated |
| 🔴 Offline | Missing | No detection — user would see "network error" as generic API error |
| 🔴 Rate limit (429) | Missing | Would appear as generic error — no cooldown indicator |
| 🔴 Forgot password | Missing | No "forgot password?" link anywhere |
| 🔴 Redirect intent | Missing | `?redirect=` param ignored — always goes to home |

### RegisterForm
| State | Present? | Details |
|-------|----------|---------|
| 🟢 Empty | ✅ | Clean form |
| 🟢 Loading | ✅ | Button spinner |
| 🟢 Validation error | ✅ | Detailed per-field |
| 🟢 API error | ✅ | Per-field backend validation mapping |
| 🔴 Offline | Missing | Same as login |
| 🔴 Rate limit | Missing | Same as login |

### Genre Onboarding
| State | Present? | Details |
|-------|----------|---------|
| 🟢 Empty | ✅ | Shows Persian "No genres available" |
| 🟢 Loading | ✅ | Pulse skeleton chips |
| 🟢 Error | ✅ | API error in Persian |
| 🟢 Success | ✅ | Saves and redirects to home |
| 🔴 Selection minimum | ✅ | Disables button if < 3 selected |

## Friction Points

| # | Severity | Location | Problem | User Impact |
|---|----------|----------|---------|-------------|
| 02.01 | 🚨 BLOCKER | `useAuth.ts:14` | `?redirect=` query param is ignored — always redirects to home | Users who clicked a protected link (e.g., Library) are sent to home after login instead of their intended page |
| 02.02 | 🚨 BLOCKER | `useLoginForm.ts` + `useAuth.ts` | Error toast fires TWICE — once in form (apiError banner) and once in useAuth (toast) | Redundant error display — confusing, visually noisy |
| 02.03 | 🚨 BLOCKER | `stores/user-auth.ts:208` | Auth restore failure silently resets state — no notification | User loses session without understanding why |
| 02.04 | ⚠️ MAJOR | `LoginForm.vue` | No "forgot password" link anywhere | Users who forget password have no recovery path |
| 02.05 | ⚠️ MAJOR | `LoginForm.vue` | No offline detection on auth pages | Users can't tell if login failure is due to network vs. credentials |
| 02.06 | ⚠️ MAJOR | `LoginForm.vue` | Rate limit (429) not detected — shows generic error | Users may keep trying without knowing they're locked out |
| 02.07 | 💡 IMPROVE | `LoginForm.vue` | "Email" label is English-only — not localized for Persian users | UX friction for Persian speakers |
| 02.08 | 💡 IMPROVE | `LoginForm.vue` | Password `toggle-mask` may not be accessible — missing aria labels on toggle | Screen reader users can't identify toggle purpose |
| 02.09 | 💡 IMPROVE | `useLoginForm.ts` | Debug `console.log` calls left in production code | (S-013 from BACKLOG.md) |
| 02.10 | 💡 IMPROVE | `AuthHeroPanel.vue:22` | Brand name says "Musicify" not "Muse" | Brand inconsistency |
| 02.11 | 💡 IMPROVE | `Genres:15` | Title is Persian but `/auth/login` and `/auth/register` are English-only | Inconsistent Persian support across auth |
| 02.12 | 💡 IMPROVE | Genre onboarding | After "skip", user goes to home with no genres selected — homepage recommendations may be poor | New users get generic suggestions |

## Clicks & Time-to-Value

| Metric | Current | Ideal |
|--------|---------|-------|
| Login (from click → app) | ~2-5s (2 pages: login → home) | 1-click if remembered |
| Register → first content | ~3-6s (register → genres → home) | Streamlined with smart defaults |
| Session restore | Async, non-blocking — but flash of wrong state possible | <500ms seamless |
| Password reset | ❌ Not available | Add flow |

## RTL / A11y / Mobile Notes

- **RTL**: Auth pages do not set RTL layout — form labels and inputs are LTR only. Persian users see English labels.
- **Skip link**: Missing in `LayoutAuth.vue` — keyboard users must tab through entire hero panel before reaching the form.
- **Focus management**: No auto-focus on email field on page load.
- **Touch targets**: All form inputs and buttons use adequate sizing (min 44px).
- **Contrast**: Glass card on dark bg — text is white on rgba(0,0,0,0.4) with backdrop-blur. Need to verify 4.5:1 ratio for the "Email" and "Password" labels (white/60 = #ffffff99 = about 3.7:1 on the glass background — **FAILS WCAG AA**).
- **ARIA**: Inputs have `aria-label`, errors have `role="alert"`. Password component toggle mask may lack proper ARIA.

## Delight Opportunities

- ✨ **Smart redirect after login**: Use the `?redirect=` param to send users where they intended to go
- ✨ **Social login**: OAuth buttons (Google, Spotify, etc.) for 1-click auth
- ✨ **Remember me**: Option to persist session longer than default
- ✨ **Persian welcome**: Show a Persian greeting on the login form for users with Persian locale
- ✨ **Animated brand intro**: The login page could show a brief musical animation on first visit

## Open Questions

1. Is there a plan for password reset/forgot password? Currently there's no route or component.
2. Should the auth forms be fully Persian (RTL) when the user's locale is Persian?
3. Is guest access intended to be permanent or temporary? Currently users can browse home/search as guest but hit auth walls on social/library features.
4. What is the expected UX when auth refresh fails mid-session? Current behavior is silent reset — should there be a toast or modal?
