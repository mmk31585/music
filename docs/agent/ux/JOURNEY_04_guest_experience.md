# Journey 04: Guest / Unauthenticated Experience

## What Routes Are Accessible Without Login?

| Route | `requiresAuth` | Guest Can Access? | Notes |
|-------|---------------|-------------------|-------|
| `/` (home) | ❌ (not set) | ✅ Full page | But personalized sections are skipped |
| `/search` | ❌ | ✅ Full search & discover |
| `/track/:id` | ❌ | ✅ Track detail page |
| `/album/:id` | ❌ | ✅ Album page |
| `/artist/:id` | ❌ | ✅ Artist page |
| `/playlist/:id` | ❌ | ✅ Playlist detail (view only) |
| `/recommendations/*` | ❌ | ✅ All rec sub-pages |
| `/explore` | ❌ | ✅ Explore page |
| `/ai/mood-explorer` | ❌ | ✅ Mood explorer |
| `/videos`, `/music-video/:id` | ❌ | ✅ Video pages |
| `/user/:id?` | ❌ | ✅ User profiles |
| `/library` | ✅ `requiresAuth` | 🚫 Redirects to `/auth/login` |
| `/playlists` | redirects to `/library` | 🚫 Same |
| `/recently-played` | redirects to `/library` | 🚫 Same |
| `/notifications` | ✅ | 🚫 Redirect |
| `/social` | ✅ | 🚫 Redirect |
| `/settings` | ✅ | 🚫 Redirect |
| `/subscription` | ✅ | 🚫 Redirect |
| `/contributions` | ✅ | 🚫 Redirect |
| `/creator-dashboard` | ✅ | 🚫 Redirect |
| `/gamification` | ✅ | 🚫 Redirect |
| `/stats` | ✅ | 🚫 Redirect |
| `/ai/playlist-generator` | ✅ | 🚫 Redirect |
| `/admin/*` | ✅ + role=admin | 🚫 Redirect |

**Summary**: Guests have access to ~15 of the ~30 app pages. The catalog browsing experience (search, tracks, albums, artists, playlists, recommendations, videos) is fully available. All interactive/social/personal features require login.

## The Guard Chain for Guests

```
Guest navigates to /library
  → checkMaintenanceGuard → null (continue)
  → checkLoginGuard → to.meta.requiresAuth=true && !isAuthenticated
     → redirect to /auth/login?redirect=/library
  → Guest sees login page
```

The `?redirect` param preserves the intended destination... but as documented in F-009, it's ignored on successful login.

## What Happens When a Guest Hits "Play"?

### The GuestPlayGate Mechanism

```typescript
// GuestPlayGate.vue
function handleAction() {
  if (isAuthenticated) {
    emit('proceed')  // Just play
    return
  }

  // Guest: check play limit
  if (props.action === 'play' && !hasReachedLimit.value) {
    incrementGuestPlay()
    emit('proceed')  // Allow play (up to 3)
    return
  }

  // Limit reached, or non-play action
  showPrompt.value = true  // Show upgrade dialog
}
```

**Play limit**: 3 plays total (stored in localStorage `moja_guest_plays`)
**Counter reset**: Never resets — 3 plays is a lifetime limit per-device unless user clears localStorage
**Non-play actions** (like, playlist, follow): Immediately show upgrade dialog — **no free actions**

### Guest Play Trace

```
Guest clicks play button
  → GuestPlayGate.handleAction()
  → Not authenticated
  → Action is 'play', guestPlayCount < 3
  → incrementGuestPlay() → count saved to localStorage
  → emit('proceed')
  → Player starts playing the track
```

After 3 plays:
```
Guest clicks 4th play button
  → GuestPlayGate.handleAction()
  → Not authenticated
  → Action is 'play', guestPlayCount >= 3 (hasReachedLimit)
  → showPrompt = true
  → GuestUpgradePrompt dialog appears (Persian):
     ┌──────────────────────────────────┐
     │ 🎧 موجا رو کشف کن                │
     │ برای ادامه شنیدن، یه حساب        │
     │ رایگان بساز. همین الان، همین جا. │
     │                                  │
     │ [ ثبت‌نام رایگان ]  ← Persian:   │
     │ [    ورود       ]  ← Persian:    │
     │                                  │
     │ تا حالا 3 از 3 آهنگ رایگان شنیدی  │
     └──────────────────────────────────┘
  → User must register or login to continue
```

### What Happens When a Guest Hits "Like" / "Save" / "Follow"?

```typescript
// GuestPlayGate.vue
// For non-'play' actions, the gate IMMEDIATELY shows upgrade dialog
// No free likes, no free follows, no free playlist saves
```

Examples:
- **Heart/like button on track** → `GuestPlayGate action="like"` → shows upgrade dialog
- **Add to playlist** → `GuestPlayGate action="playlist"` → shows upgrade dialog
- **Follow artist** → `GuestPlayGate action="follow"` → shows upgrade dialog

### What Happens If an API Returns 401?

The request wrapper (`useRequest` / Axios interceptor) handles 401s:
1. If there's a refresh token, try to refresh
2. If refresh fails, clear auth state
3. **But the user doesn't see an error** — the API call fails silently, and the UI may show stale data or empty states

## Social Pages for Guests

| Page | Guard | Guest Behavior |
|------|-------|---------------|
| `/social` (hub) | `requiresAuth` | Redirect to login |
| `/social/party/:id` | ❌ (not set) | ✅ Can view party — but can't interact |
| `/social/room/:id` | ❌ (not set) | ✅ Can view room — but can't interact |
| `/social/club/:id` | ❌ (not set) | ✅ Can view club — but can't interact |
| `/social/clubs/browse` | ✅ | Redirect to login |

Social viewing pages don't have `requiresAuth`, so guests can see the page, but any interactive action (chat, react, join) will fail or require auth via GuestPlayGate.

## Admin Pages for Non-Admin Users

```typescript
// auth-guard.ts
if (to.meta.requiresRole === 'admin') {
  if (!isAuthenticated) return redirect to login
  if (!isAdmin) return redirect to home
}
```

Non-admin users who somehow know `/admin` URLs:
- If not authenticated → redirected to login
- If authenticated but not admin → redirected to home (no error message)
- **No 403 page** — just a silent redirect

## Guest-to-Authenticated Transition

When a guest registers or logs in:
1. Auth state updates → Pinia store updates
2. `isAuthenticated` becomes true
3. Components re-render:
   - Sidebar shows user info + logout button instead of login/register links
   - Library/social/settings become accessible
   - GuestPlayGate stops showing prompts
   - Player can play unlimited tracks

**However**: The transition is not animated or communicated. There's no "You're now logged in" full-page transition. The toast ("Welcome back" / "Account created") is the only feedback.

## State Matrix Findings

### Guest Access
| Aspect | Status | Notes |
|--------|--------|-------|
| 🟢 Browse catalog | ✅ | All read-only catalog features work |
| 🟢 Search | ✅ | Full search capability |
| 🟢 View recommendations | ✅ | Global (non-personalized) recommendations |
| 🟢 View playlists | ✅ | See playlist contents, can't edit |
| 🟢 Play up to 3 tracks | ✅ | Limited to 3 lifetime plays |
| 🔴 Play counter never resets | ❌ | 3 plays = lifetime limit per device — can't listen to more even months later |
| 🔴 Guest play count visibility | ❌ | No indicator until dialog appears at limit |
| 🔴 Like/save/follow | ❌ | Immediately blocked — no "try before you buy" |
| 🔴 Guest→Auth transition | 🟡 | Works but abrupt — no guided transition |

### Guard Chain
| Aspect | Status | Notes |
|--------|--------|-------|
| 🟢 Auth-required redirect | ✅ | Redirects to login with `?redirect=` param |
| 🟢 Admin role check | ✅ | Redirects non-admins |
| 🔴 No "you need to login" page | ❌ | No intermediate page explaining *why* user is being redirected |
| 🔴 ?redirect ignored after login | ❌ | F-009 — always goes to home |
| 🔴 Non-admin gets silent redirect | ❌ | No toast or message when redirected from admin page |

## Friction Points

| # | Severity | Location | Problem | User Impact |
|---|----------|----------|---------|-------------|
| F-401 | ⚠️ MAJOR | `GuestPlayGate.vue` | Play counter is lifetime per-device, never resets | Guest who played 3 tracks months ago can never play again without registering |
| F-402 | ⚠️ MAJOR | `GuestPlayGate.vue` | No visible indicator of remaining free plays | Guest suddenly blocked on 4th play attempt — unexpected friction |
| F-403 | ⚠️ MAJOR | `GuestPlayGate.vue:32-33` | `action='like'` etc. immediately blocks with no free action allowance | No chance to "try" the like/save feature before committing |
| F-404 | ⚠️ MAJOR | Router guard | Auth-required redirect has no explanatory page — just abrupt redirect to login | Confusing: "Why am I being sent to login?" |
| F-405 | 💡 IMPROVE | `auth-guard.ts:24-26` | Non-admin redirected to home silently — no toast/message | User thinks URL is broken |
| F-406 | 💡 IMPROVE | `GuestPlayGate.vue` | Guest play count stored in localStorage — cleared on browser data wipe | Frustrating: users may lose remaining plays by clearing cache |
| F-407 | ✨ DELIGHT | Guest flow | Guest-to-auth transition is abrupt | Could show a "Your library is now unlocked!" animation after registration |

## RTL / A11y / Mobile Notes

- ✅ **Guest upgrade dialog** is fully Persian — great UX for Persian-speaking guests
- ✅ **Dialog is modal** — forces attention to the upgrade choice
- ✅ **`GuestPlayGate` pattern** is well-designed — wraps any action component, shows dialog when needed
- ❌ **No keyboard shortcut** from upgrade dialog — user must click button to register
- ❌ **Focus trap not verified** — dialog may not trap focus correctly (PrimeVue Dialog handles this, but needs verification)
- ✅ **Dialog closable** — `:closable="true"` allows dismissing the upgrade prompt

## Delight Opportunities

- ✨ **Show play counter somewhere visible**: Small "3 free plays remaining" badge in the player or sidebar
- ✨ **Free feature preview**: Allow 1 like or 1 follow to let guests taste social features
- ✨ **Graduated gating**: After 3 plays, reduce audio quality rather than blocking completely
- ✨ **Explain redirects**: Instead of abrupt redirect to login, show a brief overlay: "Log in to access your Library"
- ✨ **Guest→Auth celebration**: On registration from guest, highlight newly-unlocked features: "Now you can save songs, follow artists, and more!"

## Open Questions

1. Should the guest play counter expire (e.g., reset daily/weekly) rather than being permanent?
2. Should guests be able to see their play count somewhere (like a small badge on the play button)?
3. Is there a plan for a "limited mode" instead of hard block (e.g., audio ads, lower quality)?
4. Should the `?redirect=` flow be fixed so guests who hit a protected page end up at their intended destination after auth?
