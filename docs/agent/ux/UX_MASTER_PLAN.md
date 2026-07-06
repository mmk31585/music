# Muse UX Master Plan — Phase 8 of 8

> **Comprehensive UX Optimization Plan** derived from 30 journey audits (343 findings).
> Last updated: 2026-07-06
> Next: **Approval required before any code changes.**

---

## Table of Contents

1. [Executive Summary — Health Scores](#1-executive-summary--health-scores)
2. [Top 10 UX-Blockers (Ranked by User Impact)](#2-top-10-ux-blockers-ranked-by-user-impact)
3. [Recurring Themes — Fix Once, Fix Many](#3-recurring-themes--fix-once-fix-many)
4. [Prioritized Roadmap — Phase A through E](#4-prioritized-roadmap--phase-a-through-e)
5. [Quick Wins — Top 15](#5-quick-wins--top-15)
6. [Metrics to Track](#6-metrics-to-track)
7. [Open Product Decisions](#7-open-product-decisions)

---

## 1. Executive Summary — Health Scores

Each area rated **1–10** (10 = perfect). Rating reflects: flow completeness, error/empty/loading states, Persian/RTL quality, accessibility, and mobile support.

| Area | Score | Justification |
|------|-------|---------------|
| **🔐 Auth & Onboarding** | **4/10** | Login crashes on empty body (500). Token refresh contract broken — users logged out mid-session. `?redirect=` ignored so users always land on home. No forgot-password, no remember-me, no offline detection on auth pages. Genre onboarding has no re-entry path and no success feedback. Registration validates on submit only. |
| **🎵 Listening** | **5/10** | Core playback works, but: no "add to queue" UI anywhere, no auto-resume after refresh, current time not persisted (tracks restart at 0). Crossfade toggle exists but not implemented. NowPlayingBar (871 lines) and FullscreenPlayer (1016 lines) dangerously large. Volume slider hardcoded LTR. |
| **📚 Library** | **4/10** | All 5 sections in one `Promise.all` (slowest blocks all). No inline remove from library — must navigate to original source. No undo on unlike. No sort, no pagination, no batch operations. Silent errors. |
| **🔍 Discovery & Search** | **4/10** | Guest home skips ALL recommendation APIs. Search has no Persian typo tolerance. Error states silent everywhere. No pull-to-refresh. No periodic refresh. Empty hero on new platforms. Section ordering hardcoded. |
| **👥 Social** | **5/10** | Three solid features (parties, clubs, stages) with WebSocket real-time sync. But: no event replay on reconnect (missed events), no heartbeat (30s+ to detect drop), no presence indicators, chat not RTL-aware. |
| **💰 Monetization** | **5/10** | Two-tier system works end-to-end with 13+ gate touchpoints. But: no centralized FeatureRegistry (scattered across 15+ files), no Restore Purchases for mobile web, no retention offers on cancel, tips are USD-only. |
| **⚙️ Settings & Profile** | **3/10** | Missing major settings: audio quality, theme (dark/light), language picker. Profile uses Promise.all — one failure kills all. No edit profile modal. Follow/unfollow has no feedback. No GDPR data export. Delete account "Coming soon" indefinitely. |
| **📊 Creator Dashboard** | **3/10** | All API errors silently caught. No auto-refresh — stale data for hours. No upload flow (links away to admin). No date range picker, no export, no track-level earnings, no collaboration. |
| **🛠️ Admin — Content** | **5/10** | Sophisticated ingestion pipeline (5 stages, 5 enrichment sources). But: no pagination on media/track lists (slow at 500+). TrackFormDialog (1728 lines) dangerously large, AnyTrack type defeats TS safety. No real-time updates on moderation reports. |
| **🏛️ Admin — Ops & Governance** | **2/10** | Maintenance mode cannot be toggled (store value static). Analytics ingests rich event data but has ZERO admin UI. No general audit log. No GDPR features. No backup UI. No worker monitoring. No bulk operations. No health dashboard. |

**Overall score: 4.0/10** — Functional MVP with systemic gaps across every area. The platform works for happy-path scenarios but crumbles on errors, edge cases, and real-world usage patterns.

---

## 2. Top 10 UX-Blockers (Ranked by User Impact)

Ranked by: **breadth of users affected** × **severity of consequence** × **frequency of encounter**.

| Rank | F-ID | Problem | Journeys | Users Affected | Consequence | Fix |
|------|------|---------|----------|----------------|-------------|-----|
| **1** | F-2102 | **No offline detection** anywhere — no `online`/`offline` listeners | J21 | All users (100%) | Complete silence when network drops — no banner, no retry, no indication | Add global `window.addEventListener('online'/'offline')` + persistent banner |
| **2** | F-2101 | **Invalid route IDs produce blank page** — not 404 | J21 | All users clicking broken links | User sees empty page with no error, no navigation hint, no retry | Add data existence check to route guards before rendering |
| **3** | F-801 | **Playback does not auto-resume after refresh** | J08 | All users who listen (est. 80%) | User must manually find and tap play after every refresh/return | Persist playing state + resume on `initialize()` |
| **4** | F-701 | **No "Add to Queue" / "Play Next" UI anywhere** | J07 | All users (100%) | Users cannot curate their listening session — every play replaces queue | Add context menus + queue controls on every track |
| **5** | F-802 | **Current time not persisted** — track restarts from 0 | J08 | All listeners (80%+) | Long tracks (podcasts, classical) restart from beginning on page refresh | Persist `currentTime` + `updatedAt` timestamp |
| **6** | F-1201 | **Two separate history APIs called simultaneously** per play | J12 | All listeners (100%) | Duplicate data, confused data model, wasted API calls, potential DB duplicates | Consolidate to one API call |
| **7** | F-1101 | **No event replay on WebSocket reconnect** — missed events lost | J11, J15 | Collaborative playlist + social users | Tracks added/deleted during disconnect are lost forever | Buffer last N events server-side, replay on reconnect |
| **8** | F-2701 | **Maintenance mode cannot be toggled** — static store value | J27, J01 | Admins + all users during maintenance | Can't enable maintenance mode remotely; if somehow enabled, causes blank screen | Add `POST /admin/maintenance` endpoint |
| **9** | F-2103 | **ErrorBoundary used inconsistently** — many components not wrapped | J21 | All users on error (est. 5-10% of sessions) | Uncaught errors cause blank UI where they could show a fallback | Wrap all top-level page components + async components |
| **10** | F-008 | **Token refresh contract broken** — FE sends empty body, BE expects `{refreshToken}` | J01, J02 | All returning users (100% hit on token expiry) | Refresh fails → user logged out mid-session with no warning | Align FE request body with BE expectation |

### Honorable Mentions (would be #11-14)

| F-ID | Problem | Why Not Top 10 |
|------|---------|----------------|
| F-007 | Login 500 on empty body | Only affects users who submit empty form (rare manual trigger) |
| F-009 | `?redirect=` ignored after login | Annoying but user reaches app eventually |
| F-045 | No `/maintenance` route | Only triggered if maintenance mode is set (currently impossible) |
| F-010 | Double error toast on login | Annoying but not blocking |

---

## 3. Recurring Themes — Fix Once, Fix Many

These are **systemic patterns** that manifest in multiple findings. Fixing the root cause eliminates many findings at once.

### Theme 1: Silent Errors Everywhere 🔇
**Root cause**: API errors caught with `.catch(() => null)` or stored in `error` refs that are never displayed.
**Affects**: F-032, F-033, F-301, F-302, F-501, F-904, F-1303, F-1403, F-2001, F-2103, F-2111, F-2401
**Count**: ~30 findings
**Fix**: Global API interceptor that shows toast for all non-401 server errors + per-component error state pattern.

### Theme 2: Promise.all Fragility 💥
**Root cause**: Multiple independent API calls bundled in `Promise.all` — one failure kills all sections.
**Affects**: F-901, F-1801 (library tabs, profile sections)
**Fix**: Replace with `Promise.allSettled` pattern + per-section error/empty states. Adopt in code review checklist.

### Theme 3: Missing Empty / Loading / Error States 📋
**Root cause**: Many components handle only the "data present" case. Empty, loading, error, and offline states missing.
**Affects**: F-032, F-033, F-039, F-2101, F-501, F-904, F-1303, F-1403, F-1603, F-1604, F-2001, F-2105
**Count**: ~40 findings
**Fix**: Create `<StateWrapper>` component that renders loading/empty/error/offline variants automatically based on props.

### Theme 4: Missing Undo / Optimistic Updates ↩️
**Root cause**: Destructive actions (unlike, remove, reorder) are immediate with no undo option.
**Affects**: F-705, F-905, F-1002, F-1106, F-2210, F-2306
**Fix**: Add undo snackbar pattern (3s timeout + Undo button) for all non-critical destructive actions.

### Theme 5: No Shared Data Factories 🏭
**Root cause**: No shared factory for building `PlaybackTrack` from raw API data — 5+ duplicate code paths.
**Affects**: F-019, F-504, F-503, F-2302
**Fix**: Create `buildPlaybackTrack()` and `adminTrackToPlaybackTrack()` factories. Enforce usage via code review.

### Theme 6: WebSocket Disconnect/Reconnect Gaps 🔌
**Root cause**: No event replay buffer, no heartbeat, no connection status indicator.
**Affects**: F-1101, F-1501, F-1503, F-1510 (~8 findings)
**Fix**: Standard WebSocket wrapper with: 30s ping, event replay buffer (last 50), connection status store, exponential backoff.

### Theme 7: Session State Not Persisted 💾
**Root cause**: Queue position, current time, shuffle/repeat modes, playing state not persisted to localStorage.
**Affects**: F-702, F-703, F-801, F-802, F-803, F-804, F-807, F-808 (~12 findings)
**Fix**: Single `persistSessionState()` function in player store that saves all relevant state on `beforeunload` + periodic interval.

### Theme 8: Admin Tables Lack Core Features 📊
**Root cause**: No pagination, no sort, no search, no bulk actions on most admin list pages.
**Affects**: F-2201, F-2301, F-2205, F-2405, F-2407, F-2504, F-3001 (~15 findings)
**Fix**: Build reusable `<AdminTable>` component with built-in server-side pagination, sort, search, multi-select, and export.

### Theme 9: Missing Guest Experience Personalization 👤
**Root cause**: Guests get significantly reduced content with no explanation, no expiry on play counter.
**Affects**: F-034, F-038, F-039, F-040, F-041, F-043, F-1301 (~10 findings)
**Fix**: Add curated global recommendations for guests, show remaining play count badge, allow 1 free social action.

### Theme 10: RTL / Persian Gaps 🇮🇷
**Root cause**: Hardcoded LTR directions, no Persian typo tolerance, English-only UI in many places.
**Affects**: F-005, F-018, F-025, F-604, F-506, F-1401, F-1507, F-2112 (~15 findings)
**Fix**: Add RTL lint rule, create Persian search plugin, localize remaining English strings.

---

## 4. Prioritized Roadmap — Phase A through E

Each phase has **checkboxes** → each item can become a session prompt in the fix program.
Legend: `[S]` Small (1-4h) `[M]` Medium (1-3d) `[L]` Large (1+ week)

---

### Phase A: Blockers — Fix Broken/Dead-End Flows 🔴

> **Goal**: Eliminate all scenarios where users hit a dead end, blank screen, or data loss.
> **When**: Sprint 1 — DO NOTHING ELSE until these ship.
> **Dependencies**: Most are independent; some (F-2101, F-2103) depend on each other.

| # | F-ID | Finding | Files | Effort | Depends On |
|---|------|---------|-------|--------|------------|
| A-01 | F-2102 | **No offline detection** — add global `online`/`offline` listeners + persistent banner | Global: `App.vue`, `useOnlineStatus.ts` | [S] | None |
| A-02 | F-2101 | **Invalid dynamic route IDs** produce blank page — add data existence check | Router guards + all dynamic route pages | [M] | None |
| A-03 | F-801 | **Playback doesn't auto-resume** after refresh — persist playing state | `stores/player.ts` — `initialize()` | [M] | None |
| A-04 | F-802 | **Current time not persisted** — track restarts at 0 | `stores/player.ts` — add `currentTime` + timestamp to localStorage | [M] | None |
| A-05 | F-701 | **No "Add to Queue" / "Play Next" UI** — add context menus on all tracks | `queue-manager.ts` + all track row components | [M] | None |
| A-06 | F-1201 | **Two history APIs called simultaneously** — consolidate to one | Player engine + `stores/player.ts` | [M] | A-12 (need single data contract) |
| A-07 | F-1101 | **No WebSocket event replay** on reconnect — buffer last N events | `useCollaborativePlaylist.ts`, `useSocialSocket.ts` | [M] | None |
| A-08 | F-2701 | **Maintenance mode cannot be toggled** — add API endpoint | `stores/maintenance.ts` + backend `POST /admin/maintenance` | [M] | None |
| A-09 | F-2103 | **ErrorBoundary inconsistent** — wrap all top-level components | `ErrorBoundary.vue` + all page `App.vue` boundaries | [M] | None |
| A-10 | F-008 | **Token refresh contract broken** — align FE/BE | `refresh.ts` FE ↔ handler.go BE | [S] | None |
| A-11 | F-007 | **Login 500 on empty body** — add input validation | `handler.go:79-83` | [S] | None |
| A-12 | F-009 | **`?redirect=` ignored** on login — use stored redirect | `composables/auth/useAuth.ts:14` | [S] | None |
| A-13 | F-010 | **Double error toast** on login — deduplicate | `useLoginForm.ts` + `useAuth.ts` | [S] | None |
| A-14 | F-045 | **No `/maintenance` route** — add route + page | Router + `PageMaintenance.vue` | [S] | A-08 |

---

### Phase B: Core Loop Polish — Listen, Search, Library, Queue 🟡

> **Goal**: Make the core listening loop (search → play → queue → library) frictionless.
> **When**: Sprints 2-4
> **Dependencies**: Phase A must be done first (core paths must not be broken).

#### B1. Auth & Onboarding Quality

| # | F-ID | Finding | Files | Effort | Depends On |
|---|------|---------|-------|--------|------------|
| B-01 | F-001 | **Flash of login screen** — async auth guard | `router/index.ts:30-36` | [M] | A-10, A-12 |
| B-02 | F-011 | **Auth restore silent failure** — add toast | `stores/user-auth.ts:208` | [S] | None |
| B-03 | F-012 | **No forgot password** — add flow | `LoginForm.vue` + backend route | [M] | None |
| B-04 | F-013 | **No offline banner** on auth pages | `LoginForm.vue`, `RegisterForm.vue` | [S] | A-01 |
| B-05 | F-014 | **Rate limit (429) not handled** | `LoginForm.vue` | [S] | None |
| B-06 | F-046 | **Validation only on submit** — add blur validation | `RegisterForm.vue` | [M] | None |
| B-07 | F-047 | **No "remember me" checkbox** | `LoginForm.vue` + token storage | [S] | A-10 |
| B-08 | F-028 | **No re-run onboarding** — add to settings | `PageGenreOnboarding.vue` + settings tab | [M] | None |
| B-09 | F-029 | **Onboarding save has no feedback** | `PageGenreOnboarding.vue:129-130` | [S] | None |
| B-10 | F-017 | **Skip onboarding = no recommendations** — default genres | Genre onboarding flow | [M] | None |

#### B2. Player & Queue Quality

| # | F-ID | Finding | Files | Effort | Depends On |
|---|------|---------|-------|--------|------------|
| B-11 | F-019 | **No shared PlaybackTrack factory** | New `factories/playbackTrack.ts` | [M] | None |
| B-12 | F-022 | **Anonymous event listeners** — double-binding on hot reload | `stores/player.ts:289-297` | [S] | None |
| B-13 | F-023 | **No retry on transient error** — auto-skip is too harsh | `player-engine.ts` error handler | [M] | None |
| B-14 | F-702 | **Queue persists IDs only** — cache full metadata | Queue persistence layer | [M] | None |
| B-15 | F-703 | **Queue position not persisted** | Queue persistence | [S] | None |
| B-16 | F-704 | **Every play replaces queue** — add append mode | `setQueueAndPlay()` | [M] | A-05 |
| B-17 | F-804 | **Shuffle/repeat modes not persisted** | `stores/player.ts` | [S] | None |
| B-18 | F-806 | **No retry with backoff** on audio error | `player-engine.ts` | [M] | None |
| B-19 | F-602 | **Auto-skip without user choice** — add retry button | Player surfaces + engine | [M] | None |
| B-20 | F-603 | **Crossfade toggle exists but not implemented** | Implement or remove | [M] | None |
| B-21 | F-604 | **Volume slider hardcoded LTR** — use logical props | `NowPlayingBar.vue` | [S] | None |

#### B3. Library & History Quality

| # | F-ID | Finding | Files | Effort | Depends On |
|---|------|---------|-------|--------|------------|
| B-22 | F-901 | **All sections in one Promise.all** — stream sections | `PageLibrary.vue` | [M] | Theme 2 fix |
| B-23 | F-902 | **No remove from library** — add inline actions | `PageLibrary.vue` | [S] | None |
| B-24 | F-903 | **No unsave/unfollow from library** | `PageLibrary.vue` | [S] | None |
| B-25 | F-904 | **API error shows empty section** — add error state | `PageLibrary.vue` | [S] | Theme 1 fix |
| B-26 | F-905 | **No undo on unlike** — add snackbar | `useTrackLike.ts` | [S] | Theme 4 fix |
| B-27 | F-1001 | **No duplicate check on playlist add** | `AddToPlaylistDialog.vue` | [S] | None |
| B-28 | F-1002 | **No undo on playlist edit** | Playlist detail page | [M] | Theme 4 fix |
| B-29 | F-1202 | **History save has silent `.catch()`** | `stores/player.ts` | [S] | None |
| B-30 | F-1203 | **No remove individual history item** | History tab | [S] | None |
| B-31 | F-1204 | **No "Clear all history"** | History tab | [S] | None |

#### B4. Search & Discovery Quality

| # | F-ID | Finding | Files | Effort | Depends On |
|---|------|---------|-------|--------|------------|
| B-32 | F-1401 | **No Persian typo tolerance** — add phoneme matching | `useSearch.ts` + API | [M] | None |
| B-33 | F-1402 | **No debounce on URL param** | `PageSearch.vue` | [S] | None |
| B-34 | F-1403 | **Search error silent** — add error banner | `PageSearch.vue` | [S] | Theme 1 |
| B-35 | F-1301 | **Guest home skips ALL recommendation APIs** | `useHomeFeed.ts` | [M] | None |
| B-36 | F-1302 | **No pull-to-refresh on mobile** | `PageHome.vue` | [S] | None |
| B-37 | F-1303 | **Feed errors silently caught** | `useHomeFeed.ts` | [S] | Theme 1 |
| B-38 | F-501 | **Search API error shows "No results"** — add error banner | `PageSearch.vue:588-591` | [M] | Theme 1 |
| B-39 | F-502 | **Search play creates single-track queue** — queue related | `PageSearch.vue:650-663` | [S] | A-05 |
| B-40 | F-503 | **Album buildQueue re-maps on every click** | `PageAlbum.vue:402-405` | [S] | None |
| B-41 | F-1305 | **No periodic refresh** — stale-while-revalidate | `useHomeFeed.ts` | [M] | None |

---

### Phase C: State-Matrix Completeness — Empty/Loading/Error/Offline Everywhere 🟢

> **Goal**: Every screen handles all 6 UX states (first-load, loading, empty, error, offline, success) gracefully.
> **When**: Sprints 5-7
> **Dependencies**: Theme 1 (global error interceptor) should be done in Phase B first.

| # | F-ID | Area | What's Missing | Effort | Depends On |
|---|------|------|----------------|--------|------------|
| C-01 | F-032 | Home — error not visible | Error ref stored but never displayed | [S] | Theme 1 |
| C-02 | F-033 | Home — skeleton → empty flicker | No smooth transition on API failure | [S] | Theme 1 |
| C-03 | F-034 | Home — guest only albums/artists | Missing global recs for guests | [M] | B-35 |
| C-04 | F-035 | Home — hero empty on new platform | No fallback hero content | [S] | None |
| C-05 | F-036 | Home — no welcome after onboarding | Abrupt transition with no feedback | [S] | None |
| C-06 | F-037 | Home — "Discover" CTA → /search | Misleading navigation | [S] | None |
| C-07 | F-038 | Guest — play counter never resets | Lifetime limit with no expiry | [M] | None |
| C-08 | F-039 | Guest — no play count badge | User surprised at sudden block | [S] | C-07 |
| C-09 | F-040 | Guest — no free social action | Like/follow immediately blocked | [S] | None |
| C-10 | F-041 | Guest — no redirect explanation | "Why am I being sent to login?" | [S] | None |
| C-11 | F-042 | Guard — non-admin silent redirect | No toast when redirected from admin | [S] | None |
| C-12 | F-505 | Search — no "Play All" button | Must click each track individually | [S] | None |
| C-13 | F-706 | Queue — no "Clear Queue" button | Must remove one by one | [S] | None |
| C-14 | F-707 | Queue — no "end of queue" message | Player just stops | [S] | None |
| C-15 | F-2105 | 404 — uses empty layout | No nav/player on 404 page | [S] | None |
| C-16 | F-2106 | Admin redirect — no toast (F-042 repeat) | Already listed | [S] | C-11 |
| C-17 | F-1603 | Leaderboard — user rank not shown | No reference for users outside top 100 | [S] | None |
| C-18 | F-1604 | Leaderboard — friends tab empty | No guidance for users with 0 friends | [S] | None |
| C-19 | F-1801 | Profile — Promise.all killed by 1 failure | Use Promise.allSettled | [M] | Theme 2 |
| C-20 | F-1802 | Profile — no edit modal | Navigates away to /settings | [M] | None |
| C-21 | F-1803 | Profile — follow/unfollow no feedback | No toast on success/failure | [S] | None |
| C-22 | F-2001 | Creator dashboard — silent errors | All API errors caught with no feedback | [M] | Theme 1 |
| C-23 | F-2002 | Creator dashboard — no auto-refresh | Stale data for hours | [M] | None |
| C-24 | F-2003 | Creator dashboard — no upload flow | Links away to admin | [M] | None |
| C-25 | F-2107 | Global — no `<noscript>` tag | JS-disabled users see blank screen | [S] | None |
| C-26 | F-2108 | Global — no connection quality indicator | Users don't know why content loads slowly | [S] | None |
| C-27 | F-2111 | Global — no API error toast | Users may not notice failed actions | [S] | Theme 1 |
| C-28 | F-2112 | 404 — not localized | Persian users see English 404 | [S] | None |
| C-29 | F-2113 | Router — no route change error handler | Failed lazy imports → blank screen | [M] | None |
| C-30 | F-2114 | All pages — no skip link on error pages | Keyboard users tab through full error display | [S] | None |
| C-31 | F-610 | Keyboard shortcuts — no disable setting | Accidental triggers | [S] | None |
| C-32 | F-611 | PiP — Chrome-only, no fallback | Non-Chrome users can't use | [M] | None |
| C-33 | F-809 | Auto-PiP — no opt-in setting | Intrusive behavior | [S] | None |
| C-34 | F-609 | Same as F-809 | Duplicate | [S] | C-33 |

---

### Phase D: Admin Efficiency — Content Pipeline & Operations 🟣

> **Goal**: Make admins productive — eliminate slow workflows, missing features, and blind spots.
> **When**: Sprints 8-11
> **Dependencies**: Phase A-C must be stable. Admin builds on platform foundations.

#### D1. Admin Content Pipeline (Upload, CRUD, Moderation)

| # | F-ID | Finding | Files | Effort | Depends On |
|---|------|---------|-------|--------|------------|
| D-01 | F-2201 | Media list — no pagination | `PageAdminMedia.vue` | [M] | None |
| D-02 | F-2202 | Upload — no retry button | `UploadMediaDialog.vue` | [S] | None |
| D-03 | F-2203 | Ingestion — old theme classes | `PageAdminIngestion.vue` | [S] | None |
| D-04 | F-2301 | Track list — no pagination | `PageAdminTracks.vue` | [M] | None |
| D-05 | F-2302 | AnyTrack type defeats TypeScript | `PageAdminTracks.vue` | [M] | B-11 |
| D-06 | F-2303 | TrackFormDialog (1728 lines) — decomposing | `TrackFormDialog.vue` | [L] | None |
| D-07 | F-2401 | Moderation — no real-time updates | `PageAdminModeration.vue` | [M] | None |
| D-08 | F-2402 | Moderation — no admin notification | `PageAdminModeration.vue` | [M] | None |
| D-09 | F-2403 | Moderation — bulk action no progress | `PageAdminModeration.vue` | [S] | None |
| D-10 | F-2501 | Admin — no notifications anywhere | All admin pages | [M] | None |
| D-11 | F-2502 | Dashboard — only current counts | Admin dashboard | [M] | None |
| D-12 | F-2503 | Admin — no playlist curation | Admin sidebar + page | [M] | None |
| D-13 | F-2204 | Upload — no batch progress summary | `UploadMediaDialog.vue` | [S] | D-02 |
| D-14 | F-2205 | Media — no search/filter | `PageAdminMedia.vue` | [S] | None |
| D-15 | F-2206 | Ingestion — no enrichment notification | `PageAdminIngestion.vue` | [M] | None |
| D-16 | F-2207 | Ingestion review — no diff view | `PageAdminIngestionReview.vue` | [S] | None |
| D-17 | F-2304 | CRUD — "Enrich All" no progress | All CRUD pages | [S] | None |
| D-18 | F-2305 | CRUD — fragile event-based coupling | `TrackFormDialog.vue` | [S] | None |
| D-19 | F-2306 | CRUD — no unsaved changes warning | All form dialogs | [S] | None |
| D-20 | F-2404 | Report — no "already reported" check | Report flow (user) | [S] | None |
| D-21 | F-2405 | Moderation — no sort options | `PageAdminModeration.vue` | [S] | None |
| D-22 | F-2406 | Moderation — no assignee | `PageAdminModeration.vue` | [S] | None |
| D-23 | F-2407 | Moderation — no search within reports | `PageAdminModeration.vue` | [S] | None |

#### D2. Admin Operations & Governance

| # | F-ID | Finding | Files | Effort | Depends On |
|---|------|---------|-------|--------|------------|
| D-24 | F-2601 | **Analytics has ZERO admin UI** — build dashboard | `internal/modules/analytics/` + frontend | [L] | None |
| D-25 | F-2602 | Dashboard — only 4 counts | `internal/modules/dashboard/` | [M] | None |
| D-26 | F-2603 | Admin — no CSV/JSON export on any page | All admin | [M] | None |
| D-27 | F-2702 | Feature flags — no admin UI | `stores/feature-flags.ts` | [M] | None |
| D-28 | F-2703 | No `/admin/settings` page | All admin (TODO #45) | [M] | None |
| D-29 | F-2704 | No detailed health dashboard | `internal/modules/health/` | [M] | None |
| D-30 | F-2801 | No general-purpose audit log | All admin | [M] | None |
| D-31 | F-2802 | No GDPR features | Platform (export, deletion, consent) | [L] | None |
| D-32 | F-2803 | No data retention policy | Platform | [M] | None |
| D-33 | F-2901 | No on-demand backup trigger | Platform | [M] | None |
| D-34 | F-2902 | No storage usage dashboard | Platform | [M] | None |
| D-35 | F-2903 | No orphaned file detection | `internal/modules/media/` | [M] | None |
| D-36 | F-3001 | No bulk operations on tracks/users/albums | All admin lists | [M] | Theme 8 |
| D-37 | F-3002 | No worker monitoring UI | `internal/app/worker.go` | [M] | None |
| D-38 | F-3003 | Import has no queue or progress | `PageAdminImport.vue` | [M] | None |
| D-39 | F-2604 | Admin — no date range filter | All admin reports | [S] | None |
| D-40 | F-2605 | Dashboard — no top content lists | Admin dashboard | [S] | None |
| D-41 | F-2606 | Dashboard — no growth rate indicators | Admin analytics | [S] | None |
| D-42 | F-2504 | Admin — no user activity log | `PageAdminUsers.vue` | [M] | None |
| D-43 | F-2505 | Admin — no subscription analytics | `PageAdminSubscriptions.vue` | [M] | None |
| D-44 | F-2508 | Admin — no global search (Cmd+K) | Admin sidebar | [M] | None |
| D-45 | F-2705 | System — no worker monitoring | System health | [M] | D-37 |
| D-46 | F-2706 | System — no version/uptime display | Admin footer | [S] | None |
| D-47 | F-2707 | System — no cache management UI | Cache management | [M] | None |
| D-48 | F-3004 | Admin — no enrichment status dashboard | Enrichment | [S] | None |
| D-49 | F-3005 | Admin — no bulk operation progress bar | All admin | [S] | None |
| D-50 | F-3006 | Admin — no keyboard shortcuts | All admin | [M] | None |
| D-51 | F-3008 | Workers — no "Run Now" button | Workers | [S] | D-37 |
| D-52 | F-3010 | Dashboard — no recent actions widget | Admin dashboard | [S] | D-30 |

---

### Phase E: Delight, Retention & Personalization 🟠

> **Goal**: Build brand love through celebrations, smart defaults, and personalization.
> **When**: Sprints 12-15 (only after A-D are stable)
> **Dependencies**: Core must be solid before adding delight.

| # | F-ID | Finding | Files | Effort | Depends On |
|---|------|---------|-------|--------|------------|
| E-01 | F-044 | Guest→Auth celebration — animated feature reveal | Guest→Auth flow | [M] | None |
| E-02 | F-044 | Same as above — tracked once | — | — | — |
| E-03 | F-027 | Crossfade not exposed in UI | Player engine + settings | [S] | B-20 |
| E-04 | F-612 | No visualizer mode cycling button | `VisualizerSystem.vue` | [S] | None |
| E-05 | F-036 | Welcome toast after onboarding | `PageHome.vue` | [S] | None |
| E-06 | F-035 | Fallback hero content for empty platform | `PageHome.vue` | [S] | None |
| E-07 | F-1306 | Genre radio entry point | Radio | [S] | None |
| E-08 | F-1307 | Personalized hero (vs global trending) | `PageHome.vue` | [S] | None |
| E-09 | F-1308 | Full discography link on artist pages | `PageArtist.vue` | [S] | None |
| E-10 | F-1309 | Breadcrumb navigation on catalog pages | All catalog pages | [S] | None |
| E-11 | F-1605 | Real-time XP via WebSocket | `gamificationStore.ts` | [M] | None |
| E-12 | F-1606 | Badge progress bars | `PageBadges.vue` | [S] | None |
| E-13 | F-1608 | Reaction animation (micro-animation) | `ReactionPicker.vue` | [S] | None |
| E-14 | F-1609 | Daily challenge/tasks UI | `gamificationStore.ts` | [M] | None |
| E-15 | F-1610 | Leaderboard row highlight (own rank) | `PageLeaderboard.vue` | [S] | None |
| E-16 | F-1611 | Push notification on badge/level-up | All gamification | [S] | None |
| E-17 | F-1304 | Home section reorder in settings | `PageHome.vue` + settings | [M] | B-41 |
| E-18 | F-704 | Append vs replace queue preference | Queue behavior | [S] | B-16 |
| E-19 | F-705 | Confirm before removing current track | `QueuePanel.vue` | [S] | None |
| E-20 | F-709 | Drop target indicator on drag | `QueuePanel.vue` | [S] | None |
| E-21 | F-1003 | Streamline playlist create flow (1-step) | Playlist create | [M] | None |
| E-22 | F-1009 | Delete confirmation shows playlist name | `PagePlaylistDetail.vue` | [S] | None |
| E-23 | F-1105 | Confirmation on collab toggle | Collaborate toggle | [S] | None |
| E-24 | F-1107 | In-app invite for collaborators | Invite flow | [M] | None |
| E-25 | F-1109 | Collaborative badge on playlist cards | UI | [S] | None |
| E-26 | F-1404 | Sync recent searches across devices | `CommandPalette.vue` | [M] | None |
| E-27 | F-1409 | Search suggestions/autocomplete | `useSearch.ts` | [M] | B-32 |
| E-28 | F-1411 | Voice search (Web Speech API) | All search | [M] | None |
| E-29 | F-1504 | Independent listening mode in parties | `PartyRoom.vue` | [S] | None |
| E-30 | F-1508 | Party timer / elapsed time | `PartyRoom.vue` | [S] | None |
| E-31 | F-1509 | Stage calendar view | `StageRoom.vue` | [S] | None |
| E-32 | F-1704 | Local currency support for tips | `TipModal.vue` | [M] | None |
| E-33 | F-1705 | Family plan tier | `PageSubscription.vue` | [M] | None |
| E-34 | F-1706 | Standardize premium gate types | `FeatureGate.vue` | [M] | None |
| E-35 | F-1707 | Trial expiration reminder | Subscription | [M] | None |
| E-36 | F-1708 | Tip receipt/invoice | `TipModal.vue` | [S] | None |
| E-37 | F-1709 | Suggested tip from listening stats | `TipModal.vue` | [S] | None |
| E-38 | F-1710 | Subscription gifting | `PageSubscription.vue` | [M] | None |
| E-39 | F-1906 | Settings — "Reset to defaults" | `PageUserSettings.vue` | [S] | None |
| E-40 | F-1907 | Settings — live RTL preview | `PageUserSettings.vue` | [S] | None |
| E-41 | F-1908 | Settings — password strength meter | `PageUserSettings.vue` | [S] | None |
| E-42 | F-2004 | Creator — date range picker | `PageCreatorDashboard.vue` | [S] | C-22 |
| E-43 | F-2005 | Creator — export stats | `PageCreatorDashboard.vue` | [M] | C-22 |
| E-44 | F-2006 | Creator — track-level earnings | `PageCreatorDashboard.vue` | [M] | None |
| E-45 | F-2008 | Creator — share/promote actions | `PageCreatorDashboard.vue` | [S] | None |
| E-46 | F-2009 | Creator — scheduled releases | `PageCreatorDashboard.vue` | [M] | None |
| E-47 | F-2010 | Creator — interactive chart tooltips | `PageCreatorDashboard.vue` | [S] | None |
| E-48 | F-1804 | Profile — cover image | `PageUserProfile.vue` | [M] | None |
| E-49 | F-1805 | Profile — tracks in common with other users | `PageUserProfile.vue` | [S] | None |
| E-50 | F-1810 | Profile — real-time Now Playing | `UserHero.vue` | [M] | None |
| E-51 | F-1811 | Profile — activity feed tab | `PageUserProfile.vue` | [M] | None |
| E-52 | F-1205 | History — date grouping | History tab | [S] | None |
| E-53 | F-1206 | History — `completed: true` on track end | Player store | [S] | None |
| E-54 | F-1207 | History — "Play all" button | History tab | [S] | None |
| E-55 | F-1209 | History — private session toggle | Privacy | [M] | None |
| E-56 | F-1607 | Notifications — per-category sound prefs | `NotificationPanel.vue` | [S] | None |
| E-57 | F-1909 | Settings — per-category notification settings | `PageUserSettings.vue` | [M] | None |

---

## 5. Quick Wins — Top 15

**S-effort (1-4h), High Impact**. Can be done in parallel. Listed in priority order.

| # | F-ID | Finding | Files | Why High Impact |
|---|------|---------|-------|-----------------|
| **1** | F-2102 | **Add global offline detection** — `online`/`offline` listeners + banner | `App.vue`, `useOnlineStatus.ts` | Affects 100% of users; offline silence is the #1 blocker |
| **2** | F-008 | **Fix token refresh contract** — align FE sends `{refreshToken}` | `refresh.ts` FE + handler.go BE | All returning users hit this on token expiry |
| **3** | F-009 | **Fix `?redirect=` after login** — use stored param | `useAuth.ts:14` | Every login redirects to wrong page |
| **4** | F-010 | **Deduplicate error toast on login** — prevent double-fire | `useLoginForm.ts` + `useAuth.ts` | Every login error is doubled |
| **5** | F-007 | **Login input validation** — prevent 500 on empty body | `handler.go:79-83` | New user gets server error instead of validation message |
| **6** | F-045 | **Add `/maintenance` route** — prevent blank screen | Router + `PageMaintenance.vue` | Maintenance mode currently causes blank screen |
| **7** | F-2112 | **Localize 404 page** to Persian | `PageNotFound.vue` | Persian users see English-only error page |
| **8** | F-005 | **Set `dir` on `<html>`** element via route watcher | All layouts | Screen reader RTL detection for all Persian users |
| **9** | F-604 | **Fix volume slider RTL** — use logical CSS props | `NowPlayingBar.vue` | Every Persian user sees inverted slider |
| **10** | F-706 | **Add "Clear Queue" button** | `QueuePanel.vue` | Must remove one by one currently |
| **11** | F-707 | **Add "End of queue" state** | Queue end state | Queue just stops with no message |
| **12** | F-1203 | **Add remove button per history item** | History tab | Can't clean individual history items |
| **13** | F-1204 | **Add "Clear all history" button** | History tab | No privacy control for history |
| **14** | F-1001 | **Add duplicate check on playlist add** | `AddToPlaylistDialog.vue` | Silently duplicates tracks |
| **15** | F-505 | **Add "Play All" button to search results** | `PageSearch.vue` | Must click each track individually |

---

## 6. Metrics to Track

Instrument these BEFORE and AFTER each phase to measure improvement.

### North Star Metrics

| Metric | Definition | Current | Target | Instrumentation |
|--------|-----------|---------|--------|-----------------|
| **Time to First Play** (TTFP) | Seconds from cold start to first audio | ~4-5s (est.) | <3s | `player-engine.ts` — emit `playback:start` with timestamp since app mount |
| **Registration Completion Rate** | % of users who complete registration after starting | Unknown | >80% | Add funnel: visit /register → submit → success/error |
| **Guest→Auth Conversion** | % of guest users who register | Unknown | >15% | Track guest play → reg start → completion |
| **Session Continuity Rate** | % of sessions where playback resumes after refresh | 0% | >90% | Count `initialize()` → auto-resume success/failure |
| **Search Success Rate** | % of searches that result in a play | Unknown | >95% | Send `search:resultClicked` event with query + result count |

### Diagnostic Metrics

| Metric | Definition | Current | Target | Instrumentation |
|--------|-----------|---------|--------|-----------------|
| **Login Failure Rate** (non-password) | % of login attempts with 4xx/5xx | ~5% (est.) | <2% | Count login API errors by type |
| **Blank Page Rate** | % of route navigations that render no content | Unknown | 0% | ErrorBoundary catch count + route error handler |
| **Offline Encounter Rate** | % of sessions with offline periods | Unknown | Track | `online`/`offline` event counters |
| **Queue Build Rate** | Avg queue depth per user session | Unknown | >5 tracks | `queue-manager.ts` — emit on add/remove |
| **Playlist Create Rate** | Playlists created per active user | Unknown | >1/month | Track `playlist:created` |
| **Admin Task Completion Time** | Time to upload + finalize a track | Unknown | <5 min | Ingestion pipeline — time per stage |
| **API Error Rate** | % of API calls returning 4xx/5xx | Unknown | <1% | Global interceptor counter |

### Engagement Metrics (for Phase E)

| Metric | Definition | Current | Target |
|--------|-----------|---------|--------|
| **Daily Active Users** (DAU) | Unique users per day | Unknown | +30% |
| **Monthly Active Users** (MAU) | Unique users per month | Unknown | +20% |
| **Avg Session Duration** | Minutes per user per session | Unknown | >25 min |
| **Reaction Rate** | % of plays with a reaction | Unknown | >10% |
| **Social Feature Usage** | % of users in parties/clubs/stages per week | Unknown | >5% |
| **Tip Conversion** | % of artist page visits with a tip | Unknown | >2% |

---

## 7. Open Product Decisions

> Questions surfaced across all 30 journeys. **Mark your decision in each [ ]** before implementation begins.

### Auth & Session

1. **Guest play limit**: Should it be daily (resets every 24h), weekly, or total lifetime (current)?
   - Decision: [x] **Daily** — resets every 24h

2. **Guest play limit count**: What should the limit be? (Current: 3)
   - Decision: [x] **5** plays per day

3. **Forgot password**: Should we implement self-service (email reset link) or require admin contact?
   - Decision: [x] **Self-service** (email reset link) — standard pattern, scales without admin bottleneck

4. **Remember me**: Should "Remember Me" persist indefinitely or for a fixed period (e.g., 30 days)?
   - Decision: [x] **30 days** — balances convenience with security

5. **Email verification**: Is email verification required before the user can use the app? (Current: none)
   - Decision: [x] **Optional** — badge user as unverified, don't block access

6. **Device tracking for guests**: Are we comfortable using device fingerprinting to track guest plays across cache clears?
   - Decision: [x] **localStorage only** — avoids legal complexity of fingerprinting; acceptable trade-off for guest feature

### Playback & Queue

7. **Crossfade**: Should we implement the crossfade feature or remove the toggle? (Current: toggle exists but not implemented)
   - Decision: [x] **Implement** — with 0-12s slider in settings
   
8. **Auto-resume after crash**: Should we show a "Resume playback?" banner (user choice) or auto-resume silently?
   - Decision: [x] **Both** — auto-resume silently + toast notification

9. **Queue append behavior**: When playing a track from search/results, should it:
   - Decision: [x] **Always replace queue** (keep current behavior)

10. **Persist queue across devices**: Should queue sync to server for cross-device continuity?
    - Decision: [x] **No** — localStorage only for now. Revisit post-launch if cross-device usage is high.

### Content & Library

11. **Default genre for skippers**: When a user skips genre onboarding, what genres should be the default?
    - Decision: [x] **Balanced mix** — popular + Persian traditional

12. **Playlist visibility**: Should all playlists default to Public or Private?
    - Decision: [x] **Private** (default)

13. **History retention**: How long should listening history be retained?
    - Decision: [x] **30 days**

14. **Library pagination**: What page size for library lists?
    - Decision: [x] **50** items per page

### Social & Engagement

15. **Collaborative playlist permissions**: What roles should collaborators have?
    - Decision: [x] **Admin-only delete + all can add** — good balance of collaboration and safety. (Noting your feature request for lyric sync collaboration as a separate feature — see FR-02)

16. **Party independence**: Should party guests be able to listen independently (hear their own queue while staying in chat)?
    - Decision: [x] **Yes** — let guests listen independently while staying in chat. More flexible UX.

17. **Notification grouping**: Should similar notifications be grouped (e.g., "5 people followed you")?
    - Decision: [x] **Yes** — group similar notifications to reduce noise

### Monetization

18. **Tip currency**: Should tips support Iranian Rial/Toman in addition to USD?
    - Decision: [x] **Both** — USD + IRR/Toman

19. **Trial reminders**: Should we send trial-expiration reminders?
    - Decision: [x] **Yes** — email + push 3 days before *(recommendation accepted)*

20. **Family plan**: Should we offer a family plan?
    - Decision: [x] **Investigate pricing first** — research competitor pricing and feasibility before committing. Target: if viable, add to Phase E.

### Admin & Governance

21. **Analytics dashboard**: What are the top 5 metrics admins need to see first?
    - Decision: [x] **My recommendation** — here's what I'd build:
      1. **Daily Active Users (DAU)** with 7-day trend line (sparkline or mini-chart)
      2. **Total Plays** with top 10 most-played tracks (table + bar chart)
      3. **New Registrations** (daily count + 7-day trend)
      4. **Revenue / MRR** with churn rate and active subscribers count
      5. **Content Catalog Growth** — total tracks, artists, albums with week-over-week Δ%

22. **Moderation workflow**: Should we require moderator assignment (claim-based) or anyone can act on any report?
    - Decision: [x] **Any admin can act** — simpler workflow for small teams. Revisit with "Assign to me" button if team grows (logged as F-2406).

23. **GDPR priority**: Should GDPR compliance be blocking (must ship before launch) or post-launch?
    - Decision: [x] **Post-launch** — Sprint 14

24. **Backup schedule**: Should backups be manual-only, scheduled, or both?
    - Decision: [x] **Both** — manual trigger + scheduled (configurable)

25. **Admin dark/light mode**: Should admin theme be independent of user theme?
    - Decision: [x] **Better admin theme** (independent, purpose-built for admin workflows)

### Design & Platform

26. **Light mode**: Should we implement a light mode as an alternative to dark-only?
    - Decision: [x] **Yes** — both light and dark mode, plan for Sprint 4

27. **Language priority**: Should we localize admin pages to Persian before or after English admin features are complete?
    - Decision: [x] **English first** — then Persian for Persian users

28. **Offline play**: Should we support offline playback (downloaded tracks)?
    - Decision: [x] **No, streaming only** — offline playback is a large engineering investment. Focus on streaming quality and reliable playback first. Revisit post-launch.

---

## Appendix: Finding Count by Phase

| Severity | Phase A (Blockers) | Phase B (Core) | Phase C (States) | Phase D (Admin) | Phase E (Delight) | Total |
|----------|:-----------------:|:--------------:|:----------------:|:---------------:|:-----------------:|:-----:|
| 🚨 BLOCKER | 14 | 0 | 0 | 0 | 0 | 14 |
| ⚠️ MAJOR | 0 | 41 | 23 | 38 | 2 | 104 |
| 💡 IMPROVE | 0 | 3 | 11 | 14 | 54 | 222 |
| ✨ DELIGHT | 0 | 0 | 0 | 0 | 3 | 3 |
| **Total** | **14** | **44** | **34** | **52** | **59** | **343** |

> **Note**: Phase E contains the most findings (59) but they're all low-priority improvements and delight items.

---

## Appendix: Dependency Chain

```
Phase A (14 blockers)
  ├─→ Phase B (44 core polish) — depends on A being stable
  │     └─→ Phase C (34 state-matrix) — depends on Theme 1 from B
  │           └─→ Phase E (59 delight) — depends on core + states being solid
  └─→ Phase D (52 admin) — depends on A being stable; otherwise independent
```

---

---

## Appendix B: New Feature Requests (Captured from Product Decisions)

These emerged during decision-making and are **not yet in the friction log** — they're new feature asks.

### FR-01: Admin Review → Main App Integration
**Request**: Music that admin reviews/approves from the admin page should integrate into the main app listening flow.
**Implication**: Currently the admin ingestion pipeline and the public catalog are separate. When an admin finalizes a track via the ingestion pipeline, it appears in the catalog. But the "review" step (F-2206, F-2207) is purely internal. The request suggests a **curation/promotion** flow: tracks flagged as "editor's pick" or "approved" could appear in a hero carousel, featured section, or "New from Muse" section on the home page.
**Phase Placement**: Phase D (Admin) — add a "Promote to Home" flag during review → Phase B (Core) — surface promoted content in home feed.

### FR-02: Collaborative Lyric Sync Editor
**Request**: An editor for collaborators to sync lyrics to music — real-time collaborative lyric timing.
**Implication**: This is a **major feature** — real-time collaborative editing of lyric timestamps, similar to Google Docs for lyrics. Requires:
- WebSocket-based collaborative editing (extends existing `useCollaborativePlaylist.ts` pattern)
- Lyric line + timestamp editing interface
- Version history / rollback
- Role-based access (who can edit vs view)
**Phase Placement**: Phase E (Delight) — this is a differentiator, not a blocker. Could be a signature Muse feature.
**Effort**: Large (2-4 weeks)

### FR-03: Better Admin UI/UX (General)
**Request**: General improvement to admin UI/UX.
**Implication**: Covered by Phase D (Admin Efficiency) items — theme consistency (D-03), pagination (D-01, D-04), TypeScript safety (D-05), keyboard shortcuts (D-50), global search (D-44), and admin notifications (D-10, D-08).

---

*This document is the deliverable. Once approved, each Phase A-E item becomes a session prompt for the fix program.*
