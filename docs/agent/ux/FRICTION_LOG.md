# Muse UX Friction Log

> Master list of ALL UX findings. Updated as journeys are traced.
> Severity: 🚨 UX-BLOCKER | ⚠️ UX-MAJOR | 💡 UX-IMPROVE | ✨ UX-DELIGHT

## How To Read
Each entry: `# | Severity | Journey | Location | Problem | User Impact | Fix Category`

---

## Phase 1: App Shell & Routing

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-001 | ⚠️ MAJOR | J01 | `router/index.ts:30-36` | Guard chain runs synchronously but auth restore is async | Flash of login screen for authenticated users before redirect completes | Async-aware guard or restore-before-navigation |
| F-002 | 💡 IMPROVE | J01 | `App.vue:24` | Layout switch via dynamic component — no transition between layouts | Abrupt visual context switch when going auth→app or app→admin | Layout transition animation |
| F-003 | 💡 IMPROVE | J01 | `router/index.ts` | No global error boundary for route-level errors | Code-split chunk failures may cause blank screen | Global error handler in router |
| F-004 | 💡 IMPROVE | J01 | `router/routes/app.ts` | `/discover`, `/playlists`, `/recently-played` are hard redirects | User confusion when landing on unexpected page | Toast or animated redirect |
| F-005 | 💡 IMPROVE | J01 | All layouts (no `<html dir>`) | RTL `dir` only on inner divs, not on document `<html>` | Screen readers may misdetect page language | Set `dir` on `<html>` via route watcher |
| F-006 | 💡 IMPROVE | J01 | `LayoutAuth.vue` | No skip link, no offline banner on auth pages | Keyboard users must tab through entire hero panel; no offline notification | Add skip link + offline banner |

## Phase 2: Auth & Login

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix Category |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-007 | 🚨 BLOCKER | J02 | `internal/modules/auth/handler.go:79-83` | Login crashes with 500 on empty body `{}` | User gets server error instead of validation message | Backend input validation |
| F-008 | 🚨 BLOCKER | J02 | `frontend/src/services/api/auth/refresh.ts` | Refresh endpoint: frontend sends empty body, backend expects `{ refreshToken }` | Token refresh fails → user logged out unexpectedly | Align FE/BE contract |
| F-009 | 🚨 BLOCKER | J02 | `composables/auth/useAuth.ts:14` | `?redirect=` query param preserved but always navigates to home | Users always land on home after login, not their intended page | Use redirect query param |
| F-010 | 🚨 BLOCKER | J02 | `useLoginForm.ts` + `useAuth.ts` | Error toast fires TWICE — both in-form banner AND global toast | Redundant, noisy error display | Deduplicate error handling |
| F-011 | ⚠️ MAJOR | J02 | `stores/user-auth.ts:208` | Auth restore failure silently resets state, no user notification | User loses session without explanation | Graceful degradation |
| F-012 | ⚠️ MAJOR | J02 | `LoginForm.vue` | No "forgot password" link | Users can't recover their password | Add forgot password flow |
| F-013 | ⚠️ MAJOR | J02 | `LoginForm.vue` | No offline detection on auth pages | Users can't distinguish network vs. credential error | Add offline banner |
| F-014 | ⚠️ MAJOR | J02 | `LoginForm.vue` | Rate limit (429) not handled — generic error shown | Users may keep trying while locked out | Detect 429, show message |
| F-015 | ⚠️ MAJOR | J02 | `LoginForm.vue` | Label color `text-white/60` on glass card likely fails 4.5:1 contrast | Low-contrast labels on glass surface | Increase opacity or add background |
| F-016 | 💡 IMPROVE | J02 | `AuthHeroPanel.vue` | Brand says "Musicify" not "Muse" | Brand inconsistency | Update to "Muse" |
| F-017 | 💡 IMPROVE | J02 | Genre onboarding | After "skip", user gets no recommendations | Poor first-listening experience | Default genre selection |
| F-018 | 💡 IMPROVE | J02 | `LoginForm.vue`, `RegisterForm.vue` | Auth forms are English-only | Persian users see English labels | Localize to Persian |

## Phase 3: Playback (Core)

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix Category |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-019 | ⚠️ MAJOR | J03 | All pages building PlaybackTrack | No shared factory — 5+ different code paths build PlaybackTrack objects | Inconsistencies, missing fields, duplication | Create shared factory |
| F-020 | ⚠️ MAJOR | J03 | `NowPlayingBar.vue` (871 lines) | Single component handles desktop + mobile + popup menus + error + PiP | Maintainability risk, slow re-renders | Decompose into sub-components |
| F-021 | ⚠️ MAJOR | J03 | `FullscreenPlayer.vue` (1016 lines) | Extremely large component — cover + controls + queue + lyrics + PiP | Same as F-020 | Decompose into sub-components |
| F-022 | ⚠️ MAJOR | J03 | `stores/player.ts:289-297` | Anonymous arrow event listeners in initialize() | Double-binding on hot reloads | Named listeners + proper cleanup |
| F-023 | ⚠️ MAJOR | J03 | `services/player/player-engine.ts` | No retry on transient playback error (I-007) | Users hear glitch on temporary issues | Add retry with exponential backoff |
| F-024 | 💡 IMPROVE | J03 | `stores/player.ts:349-359` | Queue persists only IDs, re-fetches all on restore | Slow restoration for large queues | Cache full track metadata |
| F-025 | 💡 IMPROVE | J03 | `NowPlayingBar.vue` | Volume slider dir="ltr" hardcoded | RTL friction for Persian users | Use logical CSS props |
| F-026 | 💡 IMPROVE | J03 | `FullscreenPlayer.vue` | Lyrics tab loads synchronously | Stutter on lyric-rich tracks | Virtual scroll or lazy render |
| F-027 | ✨ DELIGHT | J03 | Player engine | Crossfade duration is configurable in store but not exposed in UI | Users can't enable crossfade | Add crossfade toggle to settings/overflow |

## Phase 4: Genre Onboarding

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-028 | ⚠️ MAJOR | J02 | `PageGenreOnboarding.vue` | No way to re-run onboarding from Settings/Profile | Users who skipped or want to update preferences are stuck forever | Add "Music preferences" section to settings |
| F-029 | ⚠️ MAJOR | J02 | `PageGenreOnboarding.vue:129-130` | No success feedback after save — immediate redirect to home | User may wonder if preferences were saved | Add confirmation toast or animation |
| F-030 | 💡 IMPROVE | J02 | `PageGenreOnboarding.vue:138-139` | Skip also has no feedback — immediate redirect | Same as F-029 | Brief confirmation before redirect |
| F-031 | 💡 IMPROVE | J02 | `PageGenreOnboarding.vue` | Uses raw `client.post()` directly — not a typed API service | No API contract validation, inconsistent pattern | Create typed onboarding service |

## Phase 5: Home Page & Empty State

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-032 | ⚠️ MAJOR | J03 | `PageHome.vue` | Error from `fetchHomeFeed()` stored in `error` ref but **never displayed in template** | Silent failure — user sees stale skeleton or empty state with no indication | Add error banner with retry button |
| F-033 | ⚠️ MAJOR | J03 | `PageHome.vue:22-29` | Loading skeleton renders, then transitions straight to empty state if APIs fail | Flicker from skeleton → empty even if transient | Smooth transition or automatic retry |
| F-034 | 💡 IMPROVE | J03 | `useHomeFeed.ts:52-69` | Personalized APIs called only for authenticated users — guests see only albums/artists | Guest home page is significantly less engaging | Add curated global recs for guests |
| F-035 | 💡 IMPROVE | J03 | `PageHome.vue` | Hero carousel is derived from `popular` (global trending) — empty on first load for new platform | No hero shown = less visual impact on first visit | Fallback hero content for empty catalog |
| F-036 | 💡 IMPROVE | J03 | `PageHome.vue` | After genre onboarding complete, user returns to home with no toast/feedback | Abrupt transition with no "welcome" moment | Show welcome toast or celebration |
| F-037 | 💡 IMPROVE | J03 | `PageHome.vue` | "Discover" CTA links to `/discover` which **redirects to `/search`** | User clicks "Discover" expecting new content, ends up on search page | Fix CTA destination or remove redirect |

## Phase 4b: Auth — Additional Findings

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-045 | 🚨 BLOCKER | J01 | Router: no `/maintenance` route | Maintenance guard redirects to `{ name: 'maintenance' }` but the route doesn't exist | Blank screen / Vue-router error if maintenance mode is ever enabled | Add `/maintenance` route |
| F-046 | ⚠️ MAJOR | J01 | `RegisterForm.vue` | Validation fires on submit only — no inline or blur validation | Users fill entire form, submit, then see all errors at once | Add blur/touch validation |
| F-047 | ⚠️ MAJOR | J01 | Router/Auth flow | No "remember me" checkbox — session persistence is all-or-nothing | Users must re-login every time if browser clears storage | Add "remember me" with persistent/ session-only token |
| F-048 | 💡 IMPROVE | J01 | `RegisterForm.vue` | Password rules not explicitly shown — only password strength meter | Users guess at password requirements | Show explicit rules list |

## Phase 5: Home Page & Empty State

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-038 | ⚠️ MAJOR | J04 | `GuestPlayGate.vue` | Play counter is lifetime per-device (localStorage), never resets | Guest who played 3 tracks months ago can never play again without registering | Periodic reset or expiration (e.g., 24h) |
| F-039 | ⚠️ MAJOR | J04 | `GuestPlayGate.vue` | No visible indicator of remaining free plays | Guest surprised by sudden block on 4th play attempt | Show remaining plays badge in player |
| F-040 | ⚠️ MAJOR | J04 | `GuestPlayGate.vue:32-33` | Like/follow actions **immediately blocked** — no free action allowance | No chance to "try" social features before committing | Allow 1 free like/follow as taste preview |
| F-041 | ⚠️ MAJOR | J04 | Router guard chain | Auth-required redirect has no explanatory intermediate page | Confusing: "Why am I being sent to login?" | Show brief overlay: "Log in to access this feature" |
| F-042 | 💡 IMPROVE | J04 | `auth-guard.ts:24-26` | Non-admin redirected to home silently — no toast or message | User thinks the URL is broken or access denied | Add "Admin access required" toast |
| F-043 | 💡 IMPROVE | J04 | `GuestPlayGate.vue` | Guest play count stored in localStorage only — cleared on browser data wipe | Users may lose remaining plays by clearing cache | Server-side tracking with device fingerprint |
| F-044 | ✨ DELIGHT | J04 | Guest→Auth flow | Guest-to-authenticated transition is abrupt — features just appear | Missed "wow" moment showing newly-unlocked features | Animated feature reveal after registration |

## Phase 7: Find & Play a Track (Search → Play)

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-501 | ⚠️ MAJOR | J05 | `PageSearch.vue:588-591` | Search API error silently swallowed — shows "No results found" | User thinks their search has no results, not a network error | Add error banner with retry |
| F-502 | ⚠️ MAJOR | J05 | `PageSearch.vue:650-663` | Search play creates single-track queue — track ends, nothing next | After one song, player stops | Queue more from same artist/album |
| F-503 | ⚠️ MAJOR | J05 | `PageAlbum.vue:402-405` | Album `buildQueue()` re-maps all tracks on every play click | Needless re-computation when clicking multiple tracks | Cache the queue when tracks unchanged |
| F-504 | 💡 IMPROVE | J05 | All play initiators | No shared `buildPlaybackTrack()` factory — each page constructs PlaybackTrack differently (F-019 continuation) | Inconsistencies in fallback fields | Create shared factory |
| F-505 | 💡 IMPROVE | J05 | `PageSearch.vue` | No "Play All" button on search results | Must click each track individually | Add "Play All" to search |
| F-506 | 💡 IMPROVE | J05 | `PageSearch.vue:36-38` | Search placeholder is English-only | Persian users see English | Locale-aware placeholder |
| F-507 | 💡 IMPROVE | J05 | All play buttons | No haptic/visual "playback started" confirmation on the page | User wonders if click registered | Ripple or "Now Playing" badge |
| F-508 | 💡 IMPROVE | J05 | Player engine `doPlay()` | No analytics event on playback start | Can't measure time-to-play | Fire analytics event |
| F-509 | 💡 IMPROVE | J05 | `PageTrack.vue` | `togglePlay()` silently replaces previous queue when playing a new track from track page | User's carefully ordered queue is replaced | Queue insertion vs replacement option |
| F-510 | 💡 IMPROVE | J05 | `PageAlbum.vue` `shuffleAll()` | Client-side shuffle — different order on each device for same album | Inconsistent shuffle experience | Server-seeded shuffle |

## Phase 8: Player Deep Dive (All Controls, Surfaces, Media Session, PiP)

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-601 | ⚠️ MAJOR | J06 | `FullscreenPlayer.vue` (1016 lines) + `ExpandedPlayer.vue` (960 lines) | Both extremely large components (F-021 continuation) | Maintainability risk, slow re-renders | Decompose into sub-components |
| F-602 | ⚠️ MAJOR | J06 | `player-engine.ts` error handler | Consecutive failure guard auto-skips — no user choice to retry current track | User may want to retry, not skip | Show retry option before auto-skip |
| F-603 | ⚠️ MAJOR | J06 | `ExpandedPlayer.vue` | Crossfade toggle exists but crossfade is NOT implemented in audio engine | Setting has no effect — misleading UX | Implement crossfade or remove toggle |
| F-604 | ⚠️ MAJOR | J06 | `NowPlayingBar.vue` | Volume slider hardcoded `dir="ltr"` (F-025, F-604) | RTL users see reversed slider | Use logical CSS properties |
| F-605 | 💡 IMPROVE | J06 | `NowPlayingBar.vue` | Shuffle/repeat popups close on click — can't compare options | Must open/close/open between modes | Show both simultaneously |
| F-606 | 💡 IMPROVE | J06 | `PlayerOverflowMenu.vue` | Crossfade slider present but feature not implemented (same as F-603) | User can set crossfade, nothing happens | Implement or hide |
| F-607 | 💡 IMPROVE | J06 | `FloatingMiniPlayerContent.vue` | Floating mini player has no close/minimize button | Once opened, must toggle from overflow menu | Add minimize button |
| F-608 | 💡 IMPROVE | J06 | `media-session.ts` | Media Session `stop` handler not registered | Lock screen may not clear player UI on stop | Register stop handler |
| F-609 | 💡 IMPROVE | J06 | `FullscreenPlayer.vue` | Auto-PiP on `visibilitychange` is intrusive — no opt-in | Surprise when switching tabs | Opt-in toggle |
| F-610 | 💡 IMPROVE | J06 | `useKeyboardShortcuts.ts` | No way to customize or disable keyboard shortcuts | Power users may trigger accidentally | Add shortcut disable in settings |
| F-611 | 💡 IMPROVE | J06 | `PiPPlayerContent.vue` | PiP uses Chrome-only API — no fallback for Firefox/Safari | Non-Chrome users can't use PiP | Mini-player fallback |
| F-612 | ✨ DELIGHT | J06 | `VisualizerSystem.vue` | 6 visualizer modes but no mode-cycling button in player UI | Users can't try different modes | Add visualizer mode toggle |

## Phase 9: Queue Management

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-701 | 🚨 BLOCKER | J07 | `queue-manager.ts` + UI | **No "add to queue" or "play next" UI anywhere** | Users can't build a listening queue from individual tracks | Add to track context menus and player controls |
| F-702 | ⚠️ MAJOR | J07 | Queue persistence | Queue persists only track IDs, not metadata — re-fetches all on restore (F-024) | Slow restoration; failed fetches drop tracks silently | Cache full metadata with TTL |
| F-703 | ⚠️ MAJOR | J07 | Queue persistence | Current queue position not persisted — restores to index 0 | User loses place in long queue on refresh | Persist queue index |
| F-704 | ⚠️ MAJOR | J07 | `setQueueAndPlay()` | Every new play action REPLACES the entire queue — no append mode | User's curated queue silently replaced | Confirm on replace or allow append |
| F-705 | 💡 IMPROVE | J07 | `QueuePanel.vue` | No indicator that removing current track skips to next | Surprising behavior | Show confirmation dialog |
| F-706 | 💡 IMPROVE | J07 | `QueuePanel.vue` | No "Clear Queue" button | Must remove tracks one by one | Add "Clear Queue" action |
| F-707 | 💡 IMPROVE | J07 | Queue end state | No "queue ended" message — player just stops | User thinks playback is broken | Show "End of queue" state |
| F-708 | 💡 IMPROVE | J07 | `stores/player.ts` | Shuffle/repeat modes reset on every new `setQueueAndPlay()` | Preferred mode lost each time | Persist shuffle/repeat preference |
| F-709 | 💡 IMPROVE | J07 | `QueuePanel.vue` | Drag-reorder has no drop target indicator | User doesn't know where track lands | Show drop indicator line |
| F-710 | 💡 IMPROVE | J07 | `stores/player.ts` | Engine resets shuffle to `off` on every `setQueueAndPlay()` | User with shuffle on gets sequential | Preserve shuffle mode |

## Phase 10: Continuity (Refresh, Session, State Restore)

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-801 | 🚨 BLOCKER | J08 | `stores/player.ts` — `initialize()` | **Playback does NOT auto-resume after page refresh** | User loses listening session on accidental refresh | Add "Resume playback" option on restore |
| F-802 | 🚨 BLOCKER | J08 | `stores/player.ts` | **Current time not persisted** — track restarts from 0 after refresh | User loses place in long song/podcast | Persist currentTime with timestamp |
| F-803 | ⚠️ MAJOR | J08 | Queue persistence | Queue position not persisted — restores to index 0 (F-703 repeat) | User loses place in long queue | Persist queue index |
| F-804 | ⚠️ MAJOR | J08 | `stores/player.ts` | Shuffle/repeat modes not persisted — reset on page refresh | Preferred mode forgotten | Persist to localStorage |
| F-805 | ⚠️ MAJOR | J08 | Queue persistence | Sequential `getPlaybackTrack()` calls for queue restore — slow | 15 tracks = 15 API calls | Batch fetch or cache metadata |
| F-806 | ⚠️ MAJOR | J08 | `player-engine.ts` | Audio error auto-skips without retry — consecutive failures quickly block | Brief glitch → 3 skips → blocked | Retry with backoff before skip |
| F-807 | 💡 IMPROVE | J08 | Player store | `hydratedAt` timestamp stored but never checked for expiry | Stale queue data persists forever | Add TTL-based expiry |
| F-808 | 💡 IMPROVE | J08 | Player store | No server-side queue sync | Zero cross-device continuity | Sync queue to server |
| F-809 | 💡 IMPROVE | J08 | `FullscreenPlayer.vue` | Auto-PiP has no opt-in | Surprising for new users | Add "Enable Auto-PiP" setting |
| F-810 | 💡 IMPROVE | J08 | `registerBeforeUnload` | `sendBeacon` may not fire on abrupt tab close | Some session data lost | Also fire on `visibilitychange` |
| F-811 | 💡 IMPROVE | J08 | Player engine | No "resume playback" banner on page load after crash | User returns to empty player | Smart resume offer |
| F-812 | 💡 IMPROVE | J08 | `player.initialize()` | `restorePersistedQueue()` fails silently if localStorage unavailable | Queue silently disappears | Log warning, show toast |

---

## Phase 11: Library — Save, Organize, Remove

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-901 | ⚠️ MAJOR | J09 | `PageLibrary.vue` | All 5 library sections fetched in one `Promise.all` — slowest API blocks all | Page loading until slowest response | Stream sections independently |
| F-902 | ⚠️ MAJOR | J09 | `PageLibrary.vue` | No way to remove a liked track FROM the library page | Can't clean up library without navigating away | Add remove button to TrackRow in library |
| F-903 | ⚠️ MAJOR | J09 | `PageLibrary.vue` | No way to unsave album/unfollow artist from library page | Same as F-902 for albums/artists | Add inline remove buttons |
| F-904 | ⚠️ MAJOR | J09 | `PageLibrary.vue` | API error silently swallowed — section shows empty | User thinks they have 0 items | Per-tab error state |
| F-905 | 💡 IMPROVE | J09 | `useTrackLike.ts` | Like/unlike has no undo toast | Accidental unlike requires re-finding track | Add "Undo" snackbar |
| F-906 | 💡 IMPROVE | J09 | `PageLibrary.vue` | No sort options (alpha, date added, artist) | Large libraries hard to navigate | Add sort dropdown per tab |
| F-907 | 💡 IMPROVE | J09 | `PageLibrary.vue` | No pagination — all items fetched at once | Slow for >500 items | Add pagination or virtual scroll |
| F-908 | 💡 IMPROVE | J09 | `PageLibrary.vue` | Track filter is client-side only | Only filters loaded items | Move filter to API |
| F-909 | 💡 IMPROVE | J09 | Library API | No batch operations (unlike multiple tracks) | Cleaning up is one-by-one | Multi-select + batch unlike |
| F-910 | 💡 IMPROVE | J09 | `PageLibrary.vue` | No stale-while-revalidate — library not refreshed on return | Adding track and going back doesn't show update | Refresh via KeepAlive `onActivated` |

## Phase 12: Playlist Lifecycle

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-1001 | ⚠️ MAJOR | J10 | `AddToPlaylistDialog.vue` | No duplicate check — silently adds duplicate tracks | Users accidentally duplicate tracks | Check playlist before add |
| F-1002 | ⚠️ MAJOR | J10 | Playlist detail page | No undo snackbar after add/remove/reorder | Mistakes are permanent | Add undo toast with 5s timeout |
| F-1003 | 💡 IMPROVE | J10 | Create flow | Create → redirect → add tracks is 3+ steps | High friction for first playlist | Allow adding tracks in create flow |
| F-1004 | 💡 IMPROVE | J10 | Playlist detail page | No filter/search within playlist tracks | Long playlists hard to navigate | Add inline search |
| F-1005 | 💡 IMPROVE | J10 | `PagePlaylistDetail.vue` | No batch select for remove/move | One-by-one for large playlists | Multi-select mode |
| F-1006 | 💡 IMPROVE | J10 | `PagePlaylistDetail.vue` | Cover upload has no crop/position | Cover may display poorly | Add crop step |
| F-1007 | 💡 IMPROVE | J10 | Playlist detail page | No public/private toggle after creation | Can't change visibility | Add visibility toggle |
| F-1008 | 💡 IMPROVE | J10 | Playlist detail page | No total duration shown | Users don't know playlist length | Add "X hr Y min" |
| F-1009 | 💡 IMPROVE | J10 | `PagePlaylistDetail.vue` | Delete confirmation doesn't show playlist name | Double-check is harder | Show name in confirmation |

## Phase 13: Collaborative Playlists

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-1101 | 🚨 BLOCKER | J11 | `useCollaborativePlaylist.ts` | No event replay on WebSocket reconnect | Missing changes lost permanently | Replay or re-fetch on reconnect |
| F-1102 | ⚠️ MAJOR | J11 | `useCollaborativePlaylist.ts` | No presence indicators | Users don't know if others are editing | Add "N people editing" badge |
| F-1103 | ⚠️ MAJOR | J11 | WebSocket events | `user_id` in events but not shown in UI | Can't tell who added/removed tracks | Show "Added by Alice" |
| F-1104 | ⚠️ MAJOR | J11 | Collaborative | No conflict handling for simultaneous reorder | Last write silently overwrites | Operational transform or lock |
| F-1105 | 💡 IMPROVE | J11 | Collaborate toggle | No confirmation when toggling collaborative mode | Accidental toggle makes playlist public | Add confirmation dialog |
| F-1106 | 💡 IMPROVE | J11 | Collaborator removal | No undo or confirmation | Accidentally remove someone | Add confirmation + undo |
| F-1107 | 💡 IMPROVE | J11 | Invite flow | Only share-by-link — no in-app notification | Friends don't know they're added | Send in-app notification |
| F-1108 | 💡 IMPROVE | J11 | Collaborate toggle | No granular permissions (add-only, view-only) | All collaborators can delete anything | Add "can edit" vs "can add" roles |
| F-1109 | 💡 IMPROVE | J11 | UI | No "Collaborative" badge on playlist grid cards | Users don't know which are shared | Add badge to grid items |

## Phase 14: Listening History

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-1201 | 🚨 BLOCKER | J12 | Player engine + stores | **Two separate history APIs called simultaneously** (`addPlayHistory` + `addHistory`) | Duplicate data, confused data model | Consolidate to one API |
| F-1202 | ⚠️ MAJOR | J12 | `stores/player.ts` | `.catch(() => {})` on history API — silent failure | History silently lost | Retry + warning log |
| F-1203 | ⚠️ MAJOR | J12 | History tab | No way to remove individual history items | Can't clean up history | Add remove button per item |
| F-1204 | ⚠️ MAJOR | J12 | History tab | No "Clear all history" button | History grows forever, no privacy control | Add clear button with confirmation |
| F-1205 | 💡 IMPROVE | J12 | History tab | No date grouping (Today/Yesterday/This Week) | Hard to find recent vs old | Date section headers |
| F-1206 | 💡 IMPROVE | J12 | Player store | `completed: false` always sent — never updated to `true` | Can't distinguish full listens from skips | Fire `completed: true` on track end |
| F-1207 | 💡 IMPROVE | J12 | History tab | No "Play all" from history | Can't replay recent session | Add "Play all" button |
| F-1208 | 💡 IMPROVE | J12 | History tab | No pagination for large history | Slow for 1000+ entries | Add pagination |
| F-1209 | 💡 IMPROVE | J12 | Privacy | No private session or history pause | All listening always recorded | Add "private session" toggle |

---

## Phase 15: Discovery — Home Feed, Radio, Artist/Album Pages

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-1301 | ⚠️ MAJOR | J13 | `useHomeFeed.ts` | Guest home feed skips ALL recommendation APIs — shows only albums + artists + static cards | Guests have a dramatically worse home page with zero personalization | Add generic "popular across all users" sections for guests |
| F-1302 | ⚠️ MAJOR | J13 | `PageHome.vue` | No pull-to-refresh on mobile | Users can't refresh the feed to see new content | Add pull-to-refresh |
| F-1303 | ⚠️ MAJOR | J13 | `useHomeFeed.ts` | Feed errors silently caught — sections go missing with no explanation | User may not know their feed failed to load | Show per-section error state + retry |
| F-1304 | ⚠️ MAJOR | J13 | `PageHome.vue` | Section ordering is hardcoded — no way to reorder or hide sections | Power users can't customize their home | Allow section reorder in settings |
| F-1305 | ⚠️ MAJOR | J13 | `useHomeFeed.ts` | No periodic refresh — feed is static once loaded on mount | Users never see new content without page reload | Add stale-while-revalidate or periodic refresh |
| F-1306 | 💡 IMPROVE | J13 | Radio | Radio start requires a seed track — can't start "genre radio" directly | Users must first find a track of the genre they want | Add "Genre Radio" entry point |
| F-1307 | 💡 IMPROVE | J13 | `PageHome.vue` | Hero carousel uses `popular` (global trending) — may not match user's taste | Hero often irrelevant | Use personalized for hero when available |
| F-1308 | 💡 IMPROVE | J13 | `PageArtist.vue` | "Show all" toggle for top tracks loads same page inline — no "view full discography" | Artists with many tracks get truncated at 5 | Add link to full track list page |
| F-1309 | 💡 IMPROVE | J13 | All catalog pages | No breadcrumb navigation | Users exploring deep links can't see where they are | Add breadcrumb trail |
| F-1310 | 💡 IMPROVE | J13 | Radio | Auto-refill doesn't show countdown/fetch indicator | User surprised when queue suddenly refills | Show "Loading more tracks..." indicator |

## Phase 16: Search Deep Dive

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-1401 | ⚠️ MAJOR | J14 | `useSearch.ts` + `api/search` | Results have **no Persian typo tolerance** — مشق vs موسیق vs موسیقی produce wildly different results | Persian users with common typos find nothing | Add Persian phoneme matching or suggest corrections |
| F-1402 | ⚠️ MAJOR | J14 | `PageSearch.vue` | Reacts to `?q=` URL param changes, but **no debounce** on initial URL load | If URL changes rapidly, fires every keystroke | Debounce URL param watcher |
| F-1403 | ⚠️ MAJOR | J14 | `PageSearch.vue` | Error state is **silent** — API failures just leave previous results showing | User may think search is working but results are stale | Show error banner with retry |
| F-1404 | 💡 IMPROVE | J14 | `CommandPalette.vue` | Recent searches are **localStorage only** — not synced across devices | User who searched on phone can't see on desktop | Sync recent searches to server |
| F-1405 | 💡 IMPROVE | J14 | `PageSearch.vue` | "Show all" navigates to new page with `limit=20` but **no client pagination** — just a "Show more" button | Large result sets (100+ items) require many clicks | Add infinite scroll or pagination |
| F-1406 | 💡 IMPROVE | J14 | `CommandPalette.vue` | Palette shows at most 3 items per result type — **can't expand inline** | User must go to full search page for more results | Add "Show all results" link per section |
| F-1407 | 💡 IMPROVE | J14 | `useSearch.ts` | Search history has **no deduplication** — same query can appear multiple times | Cluttered history | Deduplicate by query string |
| F-1408 | 💡 IMPROVE | J14 | `PageSearch.vue` | Filters are **not persisted** across page reloads | User must re-apply filters after navigation | Persist in URL query params |
| F-1409 | 💡 IMPROVE | J14 | `useSearch.ts` | No **search suggestions** as user types | User must complete their thought before seeing results | Add autocomplete suggestions |
| F-1410 | 💡 IMPROVE | J14 | `PageSearch.vue` | **No "search in my library" toggle** — search always searches all content | Users can't scope search to their liked tracks | Add library-scoped search toggle |
| F-1411 | 💡 IMPROVE | J14 | All search | No **voice search** button on mobile | Typing Persian with virtual keyboard is slow | Add voice search using Web Speech API |
| F-1412 | 💡 IMPROVE | J14 | `CommandPalette.vue` | Palette loads **all result types simultaneously** — no prioritization | Slow on large queries | Load top result first, then sections lazily |

## Phase 17: Social — Parties, Clubs, Stages

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-1501 | ⚠️ MAJOR | J15 | `useSocialSocket.ts` | No **connection status indicator** anywhere in UI | User has no way to know if they're connected to social features | Add connection dot in social nav |
| F-1502 | ⚠️ MAJOR | J15 | `PartyRoom.vue` | Party chat uses **local-only scroll state** — messages before joining not shown | Late joiners see a blank chat | Fetch last 50 messages on join |
| F-1503 | ⚠️ MAJOR | J15 | `useSocialSocket.ts` | No **heartbeat/ping** on socket — disconnection detected only on TCP close | Can take 30+ seconds to detect a dropped connection | Add 30s ping interval |
| F-1504 | 💡 IMPROVE | J15 | `PartyRoom.vue` | No **"listening independently" mode** — guests must follow host | Users who want to explore artist can't | Add "listen on my own" toggle |
| F-1505 | 💡 IMPROVE | J15 | `ClubDetail.vue` | Club **activity feed not paginated** — only shows latest 20 events | Can't see club history | Add "Load more" pagination |
| F-1506 | 💡 IMPROVE | J15 | `StageRoom.vue` | Reactions are **ephemeral — not persisted** | No history of who reacted | Persist top reactions |
| F-1507 | 💡 IMPROVE | J15 | All social pages | **No RTL consideration** in chat layout — messages are always LTR | Persian messages in chat align left | Use `dir="auto"` on chat messages |
| F-1508 | 💡 IMPROVE | J15 | `PartyRoom.vue` | No **party timer / elapsed time** | Users don't know how long they've been listening | Add elapsed time counter |
| F-1509 | 💡 IMPROVE | J15 | `StageRoom.vue` | No **scheduled stages calendar view** | Users can't see upcoming stages in a calendar | Add calendar/list toggle |
| F-1510 | 💡 IMPROVE | J15 | `useSocialSocket.ts` | **No event replay buffer** on reconnect (same as F-1101) | Missed events during reconnection | Buffer last 50 events server-side, replay on reconnect |

## Phase 18: Engagement — Notifications, Gamification, Reactions

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-1601 | ⚠️ MAJOR | J16 | `notificationStore.ts` | Notifications **no pagination on initial load** — fetches only latest 20 | Users can't browse their notification history | Add "Load more" with cursor pagination |
| F-1602 | ⚠️ MAJOR | J16 | `notificationStore.ts` | **No notification grouping** — 10 followers triggers 10 separate notifications | Notification spam when popular | Group "X and Y others followed you" |
| F-1603 | ⚠️ MAJOR | J16 | `PageLeaderboard.vue` | **User's rank not shown** if not in top 100 | Users below top 100 have no reference point | Add "Your rank: #1,234" below top 100 |
| F-1604 | ⚠️ MAJOR | J16 | Leaderboard API | **Friends tab shows no data** if user follows no one | Empty state doesn't explain how to fix it | Add "Find friends via contacts" action |
| F-1605 | 💡 IMPROVE | J16 | `gamificationStore.ts` | XP is **client-side stale** — only fetched once on profile page mount | Level-up animation plays on page reload, not when XP is earned | Push XP updates via WebSocket in real-time |
| F-1606 | 💡 IMPROVE | J16 | `PageBadges.vue` | Badge progress not shown for locked badges — just grayscale icon | User doesn't know how close they are to earning it | Add progress bar "423/1000 tracks" |
| F-1607 | 💡 IMPROVE | J16 | `NotificationPanel.vue` | No **notification sound preferences** — either on or off | Users who want sound for tips but not social can't configure | Per-category sound settings |
| F-1608 | 💡 IMPROVE | J16 | `ReactionPicker.vue` | No **animation** on reaction — just instant state change | Feels flat compared to other platforms | Add micro-animation on reaction toggle |
| F-1609 | 💡 IMPROVE | J16 | `gamificationStore.ts` | **No daily challenge UI** — XP for dailies exists but no "Daily Tasks" panel | Users don't know what to do for XP | Add "Daily Tasks" checklist with progress |
| F-1610 | 💡 IMPROVE | J16 | `PageLeaderboard.vue` | Leaderboard **doesn't highlight user's rank** within the list | Hard to find yourself in the list | Highlight user's row with accent color |
| F-1611 | 💡 IMPROVE | J16 | All gamification | **No notifications for badge/level earned** — user has to visit profile to see | Missed dopamine hits for engagement | Push notification on level-up and badge |

## Phase 19: Monetization — Free vs Premium, Subscriptions, Tipping

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-1701 | ⚠️ MAJOR | J17 | `FeatureGate.vue` | Premium gates use **component-level checks** scattered across 15+ files — no centralized feature registry | Hard to maintain, hard to audit, easy to miss gates | Create centralized `FeatureRegistry` with all premium features listed |
| F-1702 | ⚠️ MAJOR | J17 | `PageSubscription.vue` | **No "Restore Purchases"** for Apple/Google in-app purchases | Mobile web users who subscribed via app can't restore on web | Add restore purchases flow |
| F-1703 | ⚠️ MAJOR | J17 | `PageSubscription.vue` | **Cancel flow has no retention offers** — user cancels and that's it | No chance to retain user | Add "1 month free" or "discount" retention offer |
| F-1704 | 💡 IMPROVE | J17 | `TipModal.vue` | Tips are **currency-agnostic** — always USD, no conversion | Persian users want to tip in Tomans or Rials | Add local currency support |
| F-1705 | 💡 IMPROVE | J17 | `PageSubscription.vue` | No **family plan** tier | Competing platforms offer family sharing | Add family plan (up to 6 accounts) |
| F-1706 | 💡 IMPROVE | J17 | `FeatureGate.vue` | **Gate type inconsistency** — some use overlay, some use disabled, some use toast | Inconsistent UX across premium features | Standardize gate types per feature category |
| F-1707 | 💡 IMPROVE | J17 | Subscription | **No trial-to-paid reminder** | Users who sign up for trial and forget lose access with no notice | Send email + push notification 3 days before trial ends |
| F-1708 | 💡 IMPROVE | J17 | `TipModal.vue` | **No tip receipt/invoice** | Users who tip for business can't expense it | Email receipt on successful tip |
| F-1709 | 💡 IMPROVE | J17 | `TipModal.vue` | **No suggested tip amount** based on listening time | User unsure how much to tip | Show "You've listened to X hours of this artist's music" |
| F-1710 | 💡 IMPROVE | J17 | `PageSubscription.vue` | **No subscription gifting** — "Gift Premium to a friend" | Users can't gift subscriptions | Add gift subscription flow |

---

## Phase 20: User Profile — Public Profile, Stats, Activity

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-1801 | ⚠️ MAJOR | J18 | `PageUserProfile.vue` | Profile fetch uses `Promise.all` — if one API fails, the entire profile fails | Partial profile data is better than none | Use `Promise.allSettled` per section |
| F-1802 | ⚠️ MAJOR | J18 | `PageUserProfile.vue` | **No edit profile modal** — "Edit Profile" navigates to `/settings` | User leaves profile context to make changes | Inline edit or modal for basic fields |
| F-1803 | ⚠️ MAJOR | J18 | `PageUserProfile.vue` | Follow/unfollow has **no toast/feedback** on success or failure | User may not know if action registered | Add success toast, error toast on failure |
| F-1804 | 💡 IMPROVE | J18 | `PageUserProfile.vue` | **No profile header cover image** — only avatar + gradient | Profile feels sparse compared to competitors | Add optional cover image upload |
| F-1805 | 💡 IMPROVE | J18 | `PageUserProfile.vue` | **No "tracks in common"** section for other users | No shared-taste discovery | Add mutual liked tracks section |
| F-1806 | 💡 IMPROVE | J18 | `UserHero.vue` | Bio character limit **not shown** when editing | User may exceed limit without warning | Show "X/200" counter |
| F-1807 | 💡 IMPROVE | J18 | `PageUserProfile.vue` | **No pagination** on followers/following tabs | Large follow lists (500+) unfetchable | Add API pagination or virtual scroll |
| F-1808 | 💡 IMPROVE | J18 | `PageUserProfile.vue` | **No "sort by"** on tracks tab — always by date added | Can't find tracks alphabetically | Add sort options |
| F-1809 | 💡 IMPROVE | J18 | `PageUserProfile.vue` | **No "View as guest"** toggle for own profile | Can't preview what others see | Add preview mode toggle |
| F-1810 | 💡 IMPROVE | J18 | `UserHero.vue` | Now Playing section for other users **not updated in real-time** | Stale "now playing" after user changes track | Poll every 30s or use WebSocket push |
| F-1811 | 💡 IMPROVE | J18 | `PageUserProfile.vue` | **No activity feed** tab — user's recent public listens | Can't see what user has been listening to | Add "Activity" tab |
| F-1812 | 💡 IMPROVE | J18 | `PageUserProfile.vue` | **No "Report user"** action | Users can't report inappropriate profiles | Add report button with reason picker |
| F-1813 | 💡 IMPROVE | J18 | `PageUserProfile.vue` | **No "Block user"** action | Users can't block unwanted interactions | Add block button with confirmation |

## Phase 21: Settings — Audio, Appearance, Privacy, Notifications

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-1901 | ⚠️ MAJOR | J19 | `PageUserSettings.vue` | **No audio quality setting** — not present in UI at all | Premium users can't select 320kbps | Add audio quality selector |
| F-1902 | ⚠️ MAJOR | J19 | `PageUserSettings.vue` | **Theme (dark/light) toggle** completely absent | Users locked into dark-only | Add light mode with proper contrast |
| F-1903 | ⚠️ MAJOR | J19 | `PageUserSettings.vue` | **No language/locale picker** — locale inherited from browser | Persian users may see English UI | Add language selector (fa/en) |
| F-1904 | 💡 IMPROVE | J19 | `PageUserSettings.vue` | Save button has **no debounce** — double-click fires two saves | Duplicate saves, possibly overwriting | Add debounce + loading state |
| F-1905 | 💡 IMPROVE | J19 | `PageUserSettings.vue` | Avatar upload has **no progress indicator** | User doesn't know if upload is happening | Add progress bar or spinner |
| F-1906 | 💡 IMPROVE | J19 | `PageUserSettings.vue` | **No "Reset to defaults"** button for preferences | Accidental changes are permanent | Add "Reset to defaults" per section |
| F-1907 | 💡 IMPROVE | J19 | `PageUserSettings.vue` | **No preview** when changing RTL — must save and navigate away | Can't A/B test layout preference | Add live preview area |
| F-1908 | 💡 IMPROVE | J19 | `PageUserSettings.vue` | Password change has **no strength meter** — only min 8 char check | Users may set weak passwords | Add password strength indicator |
| F-1909 | 💡 IMPROVE | J19 | `PageUserSettings.vue` | Notification preferences **only 2 toggles** — no per-category | Users can't granularly control notifications | Add per-category notification settings |
| F-1910 | 💡 IMPROVE | J19 | `PageUserSettings.vue` | **No "Export my data"** GDPR/Privacy option | Users can't download their data | Add data export request button |
| F-1911 | 💡 IMPROVE | J19 | `PageUserSettings.vue` | Delete account labeled "Coming soon" with **no timeline or feedback** | Users who want to leave can't | Add request deletion flow |
| F-1912 | 💡 IMPROVE | J19 | `PageUserSettings.vue` | **Preferences tab saves independently** — each toggle fires its own API call | Multiple rapid toggles spam the server | Debounce + batch preference saves |

## Phase 22: Creator Dashboard — Analytics, Uploads, Earnings

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-2001 | ⚠️ MAJOR | J20 | `PageCreatorDashboard.vue` | **All API errors silently caught** — `.catch(() => null)` with no user feedback | Creator thinks dashboard is working but data is stale | Add per-section error state + retry button |
| F-2002 | ⚠️ MAJOR | J20 | `PageCreatorDashboard.vue` | **No auto-refresh or WebSocket push** — stats only update on manual refresh | Creators see stale data for hours | Add 30s auto-refresh or WebSocket push |
| F-2003 | ⚠️ MAJOR | J20 | `PageCreatorDashboard.vue` | **No upload flow in dashboard** — "Upload" links to `/admin/media` | Disconnected experience, creator leaves context | Add inline upload widget |
| F-2004 | 💡 IMPROVE | J20 | `PageCreatorDashboard.vue` | **No date range picker** — daily stats always show last 14 days | Can't compare month-over-month | Add date range selector with presets |
| F-2005 | 💡 IMPROVE | J20 | `PageCreatorDashboard.vue` | **No export/download** for stats (CSV/PDF) | Creators can't share reports | Add export button per tab |
| F-2006 | 💡 IMPROVE | J20 | `PageCreatorDashboard.vue` | **No track-level earnings** — only aggregate | Can't tell which tracks earn most | Add "Revenue by Track" breakdown |
| F-2007 | 💡 IMPROVE | J20 | `PageCreatorDashboard.vue` | **No notification settings** for creator milestones | Creators don't know when they hit milestones | Add milestone alert preferences |
| F-2008 | 💡 IMPROVE | J20 | `PageCreatorDashboard.vue` | **No "Promote" or "Share" actions** for tracks | Creators can't easily share their work | Add share link/copy per track |
| F-2009 | 💡 IMPROVE | J20 | `PageCreatorDashboard.vue` | **No "Scheduled Releases"** — upload goes live immediately | No pre-release planning | Add scheduled publish date |
| F-2010 | 💡 IMPROVE | J20 | `PageCreatorDashboard.vue` | **Daily chart has no interactivity** — static bars, no tooltips | Hard to read exact numbers | Add tooltip on hover |
| F-2011 | 💡 IMPROVE | J20 | `PageCreatorDashboard.vue` | **No "Top Cities"** in geographic stats — only country level | Less granular audience insight | Add city-level data |
| F-2012 | 💡 IMPROVE | J20 | `PageCreatorDashboard.vue` | **No collaboration/invite** — other creators can't be added | No co-creator workflow | Add "Add collaborator" on track edit |

## Phase 23: Error Pages & Edge Cases — 404, Offline, Maintenance, Degraded

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-2101 | 🚨 BLOCKER | J21 | Dynamic routes (track/album/artist) | **Invalid/non-existent IDs return empty pages** — route matches but data is missing | User sees blank page instead of 404 | Add data existence check to route guards |
| F-2102 | 🚨 BLOCKER | J21 | Global | **No offline detection anywhere** — no `online`/`offline` event listeners | Users have no way to know app is offline | Add global offline detection + banner |
| F-2103 | 🚨 BLOCKER | J21 | `ErrorBoundary.vue` | **ErrorBoundary used inconsistently** — many components not wrapped | Uncaught errors cause blank UI | Wrap all top-level page components |
| F-2104 | ⚠️ MAJOR | J21 | Maintenance store | **No `/maintenance` route** exists (F-045 continuation) | Maintenance mode causes blank screen | Add `/maintenance` route with UI |
| F-2105 | ⚠️ MAJOR | J21 | `PageNotFound.vue` | **404 page uses "layout-empty"** — no nav, no player, no footer | Once on 404, user can't access player | Use default layout for 404 |
| F-2106 | ⚠️ MAJOR | J21 | `router/middleware/auth-guard.ts` | Non-admin redirected to home **silently** (F-042 continuation) | User thinks the URL is broken | Add "Admin access required" toast |
| F-2107 | 💡 IMPROVE | J21 | Global | **No site-wide `<noscript>` tag** | Users with JS disabled see white screen | Add `<noscript>` with message |
| F-2108 | 💡 IMPROVE | J21 | Global | **No connection quality indicator** — no slow network detection | Users don't know why content loads slowly | Add connection quality detection |
| F-2109 | 💡 IMPROVE | J21 | `stores/feature-flags.ts` | Feature flags **not polled** — require page refresh to update | Flag rollouts require user to refresh | Add periodic flag polling |
| F-2110 | 💡 IMPROVE | J21 | `ErrorBoundary.vue` | **No error reporting** — errors caught but not logged to server | Can't diagnose production errors | Add error reporting service |
| F-2111 | 💡 IMPROVE | J21 | Global | **No "Oops, something broke" toast** for failed API mutations | Users may not notice a failed action | Add global API error toast |
| F-2112 | 💡 IMPROVE | J21 | `PageNotFound.vue` | **404 page not localized** — English only | Persian users see English 404 | Localize 404 page to Persian |
| F-2113 | 💡 IMPROVE | J21 | Router | **No route change error handling** — code-split failures, no fallback | White screen on failed lazy import | Add route-level error handler |
| F-2114 | 💡 IMPROVE | J21 | All pages | **No skip link on error pages** — keyboard users can't skip to main | Must tab through entire error state | Add skip link to all layouts |

## Phase 24: Admin Upload & Ingestion Pipeline

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-2201 | ⚠️ MAJOR | J22 | `PageAdminMedia.vue` | **Media list has no pagination** — fetches all items at once | Slow page load with 500+ media items | Add pagination or virtual scroll |
| F-2202 | ⚠️ MAJOR | J22 | `UploadMediaDialog.vue` | Upload error has **no retry button** — must re-select file | Frustrating for large files on flaky connections | Add retry button |
| F-2203 | ⚠️ MAJOR | J22 | `PageAdminIngestion.vue` | Ingestion stats uses **old PrimeVue surface-* theme classes** | Visual inconsistency with admin dark theme | Migrate to admin dark theme tokens |
| F-2204 | 💡 IMPROVE | J22 | `UploadMediaDialog.vue` | No **batch upload progress summary** — only per-file | Can't see overall progress | Add batch progress bar |
| F-2205 | 💡 IMPROVE | J22 | `PageAdminMedia.vue` | No **search/filter** on media list | Can't find specific file | Add search by filename + type filter |
| F-2206 | 💡 IMPROVE | J22 | `PageAdminIngestion.vue` | No **notification when enrichment completes** | Admin must manually refresh | Add WebSocket push or polling badge |
| F-2207 | 💡 IMPROVE | J22 | `PageAdminIngestionReview.vue` | No **version comparison** for enriched fields | Can't see what enrichment changed | Add diff view (original → enriched) |
| F-2208 | 💡 IMPROVE | J22 | Ingestion pipeline | **Media and ingestion are disconnected** — media uploads don't appear in drafts | Admin confused about which tool to use | Unify media and ingestion workflows |
| F-2209 | 💡 IMPROVE | J22 | `PageAdminImport.vue` | Import search has **no rate-limit indicator** | Admin thinks search is broken when rate-limited | Show "API rate limited — retry in X seconds" |
| F-2210 | 💡 IMPROVE | J22 | `PageAdminIngestion.vue` | **No undo** on draft delete | Accidental delete loses work | Add soft-delete or undo toast |
| F-2211 | 💡 IMPROVE | J22 | `PageAdminIngestion.vue` | **formatDuration treats 0 as falsy** — 0-second shows "--" | Confusing for short audio | Fix formatDuration |
| F-2212 | 💡 IMPROVE | J22 | `PageAdminImport.vue` | Import doesn't **navigate to review** after import | Admin must find the draft manually | Auto-navigate to review |
| F-2213 | 💡 IMPROVE | J22 | `UploadMediaDialog.vue` | No **bulk delete** for media items | One-by-one cleanup | Add multi-select + bulk delete |
| F-2214 | 💡 IMPROVE | J22 | `PageAdminIngestionReview.vue` | No **"finalize and add another"** | Admin re-opening drafts one at a time | Add "Finalize & Next" button |

## Phase 25: Admin CRUD — Track, Album, Artist, Genre

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-2301 | ⚠️ MAJOR | J23 | `PageAdminTracks.vue` | **No pagination on track list** — fetches all tracks at once | Extremely slow with 1000+ tracks | Add server-side pagination |
| F-2302 | ⚠️ MAJOR | J23 | `PageAdminTracks.vue` | **AnyTrack type defeats TypeScript safety** — 10+ fragile getter functions | Runtime errors from unexpected API shapes | Create proper Track type |
| F-2303 | ⚠️ MAJOR | J23 | `TrackFormDialog.vue` (1728 lines) | **Extremely large form** — audio, metadata, lyrics, enrichment all in one | Maintainability risk, slow render | Decompose into sub-forms |
| F-2304 | 💡 IMPROVE | J23 | All CRUD pages | **No "Enrich All" progress indicator** — batch POST with no feedback | Admin doesn't know if it's working | Add progress bar with item count |
| F-2305 | 💡 IMPROVE | J23 | `TrackFormDialog.vue` | **Artist/Album inline creation dispatches events** — fragile coupling | Race conditions, hard to debug | Use callback props instead of events |
| F-2306 | 💡 IMPROVE | J23 | All form dialogs | **No unsaved changes warning** when closing dialog | Accidental close loses form data | Add `beforeClose` check with confirmation |
| F-2307 | 💡 IMPROVE | J23 | `PageAdminTracks.vue` | **No column customization** in table view | Admin may want different data visible | Add column toggle/visibility |
| F-2308 | 💡 IMPROVE | J23 | `AdminArtistsCard.vue` | **No album/track count** shown in artist grid | Can't gauge artist popularity at a glance | Add track/album count badges |
| F-2309 | 💡 IMPROVE | J23 | `AdminAlbumsCard.vue` | **No release year filter** in album grid | Hard to find albums from a specific year | Add year filter or facet |
| F-2310 | 💡 IMPROVE | J23 | `PageAdminTrackDetail.vue` | **Delete from detail loses scroll position** | Must find track again in list | Preserve scroll position |
| F-2311 | 💡 IMPROVE | J23 | All CRUD pages | **No "duplicate" action** for tracks/albums/artists | Manually re-entering similar metadata | Add "Duplicate" button to pre-fill form |
| F-2312 | 💡 IMPROVE | J23 | All CRUD pages | **No audit log** shown for edits | Can't see what changed or who changed it | Add "Last edited by X on date" metadata |

## Phase 26: Admin Moderation Queue

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-2401 | ⚠️ MAJOR | J24 | `PageAdminModeration.vue` | **No real-time updates** — must refresh to see new reports | Admins miss urgent copyright reports | Add WebSocket push for new reports |
| F-2402 | ⚠️ MAJOR | J24 | `PageAdminModeration.vue` | **No notification** for admins when new report arrives | High-priority reports sit unnoticed | Add admin notification bell badge |
| F-2403 | ⚠️ MAJOR | J24 | `PageAdminModeration.vue` | **Bulk action has no progress bar** — text only | No visual feedback for long batches | Add progress bar with percentage |
| F-2404 | 💡 IMPROVE | J24 | Report flow (user) | **No "already reported" check** — can report same content multiple times | Duplicate reports | Check existing pending report before creating |
| F-2405 | 💡 IMPROVE | J24 | `PageAdminModeration.vue` | **No sort options** in queue — always newest first | Can't prioritize by severity | Add sort by reason, date, reporter |
| F-2406 | 💡 IMPROVE | J24 | `PageAdminModeration.vue` | **No moderator assignment** — any admin can act on any report | No accountability | Add "Assign to me" button |
| F-2407 | 💡 IMPROVE | J24 | `PageAdminModeration.vue` | **No searching within reports** | Hard to investigate repeat offenders | Add search by target/reporter name |
| F-2408 | 💡 IMPROVE | J24 | Report flow (user) | **No edit/cancel** after submitting report | User submits wrong reason, can't correct | Add "Edit" / "Cancel" within 5 min |
| F-2409 | 💡 IMPROVE | J24 | `PageAdminModeration.vue` | **No appeal mechanism** for resolved reports | Can't contest a dismissal | Add "Appeal" button with reason |
| F-2410 | 💡 IMPROVE | J24 | Flag system | **No auto-unflag** notification when expired | Content may still be problematic | Add auto-flag renewal for repeat offenders |
| F-2411 | 💡 IMPROVE | J24 | `PageAdminModeration.vue` | **No keyboard shortcuts** — all actions require mouse | Slow workflow for high-volume moderation | Add keyboard shortcuts (r/d/j/k) |
| F-2412 | 💡 IMPROVE | J24 | `PageAdminModeration.vue` | **Stats tab not auto-refreshing** | Stale data | Add periodic stats refresh (60s) |

## Phase 27: Admin Dashboard & Operations

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-2501 | ⚠️ MAJOR | J25 | All admin pages | **No admin notifications** — no bell, no real-time alerts | Admins must manually check each section | Add admin notification system |
| F-2502 | ⚠️ MAJOR | J25 | Admin dashboard | **Dashboard shows only current counts** — no trends/charts | Can't monitor platform health over time | Add time-series charts (7/30/90 days) |
| F-2503 | ⚠️ MAJOR | J25 | Admin sidebar | **No playlist curation route** — playlists managed only on public side | No way to feature/promote playlists | Add `/admin/playlists` page |
| F-2504 | 💡 IMPROVE | J25 | `PageAdminUsers.vue` | **No user activity log** — can't see user's recent actions | Hard to investigate problem users | Add "Activity" tab per user |
| F-2505 | 💡 IMPROVE | J25 | `PageAdminSubscriptions.vue` | **No subscription analytics** — no MRR, churn, conversion | Can't measure business health | Add subscription metrics dashboard |
| F-2506 | 💡 IMPROVE | J25 | `PageAdminUsers.vue` | **No impersonation mode** | Support can't debug user issues | Add "Log in as user" with audit trail |
| F-2507 | 💡 IMPROVE | J25 | `PageAdminSubscriptions.vue` | **No refund flow** — must process externally | Support friction | Add "Issue refund" button |
| F-2508 | 💡 IMPROVE | J25 | Admin sidebar | **No search across admin** — must navigate to each section | Slow to find content | Add global admin search (Cmd+K) |
| F-2509 | 💡 IMPROVE | J25 | All admin pages | **No bulk user operations** — suspend/ban/verify one at a time | Slow for spam waves | Add multi-select + bulk actions |
| F-2510 | 💡 IMPROVE | J25 | All admin pages | **No "last edited by"** on any catalog item | No accountability | Add audit metadata to catalog items |
| F-2511 | 💡 IMPROVE | J25 | Admin dashboard | **No system health section** — server status, queue depth, errors | Can't monitor infrastructure | Add health panel |
| F-2512 | 💡 IMPROVE | J25 | `PageAdminContributions.vue` | **No bulk contribution approval** — one at a time | Slow for bulk submissions | Add "Approve all pending" with XP presets |
| F-2513 | 💡 IMPROVE | J25 | All admin pages | **No export/download** for any list | Can't do offline analysis | Add CSV export per page |
| F-2514 | 💡 IMPROVE | J25 | `PageAdminVideos.vue` | **No bulk video operations** — approve/reject one at a time | Slow for video moderation | Add multi-select + bulk approve/reject |

## Phase 28: Admin Reporting & Analytics

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-2601 | ⚠️ MAJOR | J26 | `internal/modules/analytics/` | **Analytics module has rich event data but ZERO admin UI** to query it | Admins can't see user growth, top content, engagement | Build admin analytics dashboard with charts |
| F-2602 | ⚠️ MAJOR | J26 | `internal/modules/dashboard/` | **Dashboard stats returns only 4 counts** — no users, plays, DAU/WAU/MAU | Dashboard nearly useless for monitoring | Expand with richer stats |
| F-2603 | ⚠️ MAJOR | J26 | All admin pages | **No CSV/JSON export** for any data | Can't do offline analysis | Add "Export" button to every list page |
| F-2604 | 💡 IMPROVE | J26 | All admin pages | **No date range filter** on any admin report | Can't compare periods | Add date range picker |
| F-2605 | 💡 IMPROVE | J26 | Admin dashboard | **No top content lists** (tracks, artists, albums) | Can't see what's popular | Add "Top 10" sections |
| F-2606 | 💡 IMPROVE | J26 | Admin analytics | **No growth rate indicators** — no week-over-week Δ% | Can't spot trends | Add Δ% badges to metrics |
| F-2607 | 💡 IMPROVE | J26 | `internal/modules/analytics/` | **No custom query endpoint** — can't filter by date or type | Can't answer ad-hoc questions | Add query endpoint |
| F-2608 | 💡 IMPROVE | J26 | Admin analytics | **No search query analytics** — can't see what users search for | Missed content gap opportunities | Add search term analytics page |
| F-2609 | 💡 IMPROVE | J26 | Admin analytics | **No retention/cohort analysis** | Can't measure retention | Add cohort tables |
| F-2610 | 💡 IMPROVE | J26 | Admin analytics | **No geographic analytics** | Can't target regional content | Add geo-map or country list |
| F-2611 | 💡 IMPROVE | J26 | Admin analytics | **No exportable report scheduling** — weekly PDF/CSV | Manual reporting only | Add scheduled report delivery |

## Phase 29: System Administration (Health, Config, Maintenance)

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-2701 | 🚨 BLOCKER | J27 | `stores/maintenance.ts` | **Maintenance mode cannot be toggled via API** — store value is static | Can't enable maintenance remotely | Add `POST /admin/maintenance` API endpoint |
| F-2702 | ⚠️ MAJOR | J27 | `stores/feature-flags.ts` | **Feature flags have no admin UI** — must restart server to toggle | Can't do gradual rollouts or kill-switches | Add admin feature flag page |
| F-2703 | ⚠️ MAJOR | J27 | All admin | **No `/admin/settings` page** (TODO #45) | No centralized platform configuration | Add admin settings page |
| F-2704 | ⚠️ MAJOR | J27 | `internal/modules/health/` | **No detailed health dashboard** — only basic liveness probes | Can't see DB/Redis/storage/worker status | Add health page with per-service indicators |
| F-2705 | 💡 IMPROVE | J27 | System health | **No worker monitoring** — can't see if background jobs are running | Silent worker failures go unnoticed | Add worker status page |
| F-2706 | 💡 IMPROVE | J27 | System health | **No uptime / version display** in admin | Can't verify build deployment | Add build version + deploy time to footer |
| F-2707 | 💡 IMPROVE | J27 | Cache management | **No cache management UI** — can't invalidate or inspect cache | Stale data persists until TTL | Add cache page with invalidation |
| F-2708 | 💡 IMPROVE | J27 | Feature flags | **No audit log for flag changes** | Can't track who changed what | Add flag change audit trail |
| F-2709 | 💡 IMPROVE | J27 | Maintenance mode | **No scheduled maintenance** — can't set "starts at 2 AM" | Always requires immediate action | Add scheduled maintenance with countdown |
| F-2710 | 💡 IMPROVE | J27 | Maintenance mode | **No maintenance page UI** (F-045/F-2104) | Users see blank screen during maintenance | Build proper maintenance page |
| F-2711 | 💡 IMPROVE | J27 | System health | **No log viewer** — must SSH to debug | Debugging friction | Add log viewer with level filter |
| F-2712 | 💡 IMPROVE | J27 | System health | **No alert configuration** — can't set health thresholds | Silent degradation until outage | Add alert threshold config |

## Phase 30: Audit & Compliance

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-2801 | ⚠️ MAJOR | J28 | All admin | **No general-purpose audit log** — only moderation actions tracked | No accountability for catalog/settings changes | Add AdminAction model for all admin ops |
| F-2802 | ⚠️ MAJOR | J28 | Platform | **No GDPR compliance features** — no data export, account deletion, consent | Legal risk, user trust erosion | Implement export + deletion + consent |
| F-2803 | ⚠️ MAJOR | J28 | Platform | **No data retention policy** — no auto-purge of old data | Storage grows unbounded, privacy risk | Add retention policies with automated cleanup |
| F-2804 | 💡 IMPROVE | J28 | Audit | **No diff view** for changes — can't see before/after | Hard to understand what changed | Add JSON diff viewer |
| F-2805 | 💡 IMPROVE | J28 | Audit | **No admin session audit** — no record of who logged in | Can't detect unauthorized access | Add login/logout event tracking |
| F-2806 | 💡 IMPROVE | J28 | Compliance | **No privacy policy page** — no `/privacy` route | Legal requirement unmet | Add privacy policy page |
| F-2807 | 💡 IMPROVE | J28 | Compliance | **No terms of service page** — no `/terms` route | Legal requirement unmet | Add terms of service page |
| F-2808 | 💡 IMPROVE | J28 | Compliance | **No cookie consent banner** | GDPR requirement for EU users | Add cookie consent banner |
| F-2809 | 💡 IMPROVE | J28 | Audit | **No pagination/search on moderation audit** — per-report only | Can't search across all actions | Add global audit log view |
| F-2810 | 💡 IMPROVE | J28 | Compliance | **No "right to be forgotten" cascade** | User data persists after deletion | Add cascade delete for all user data |
| F-2811 | 💡 IMPROVE | J28 | Audit | **No export of audit log** — can't download for compliance | Must manually screenshot | Add CSV export with date filter |

## Phase 31: Backup & Data Management

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-2901 | ⚠️ MAJOR | J29 | Platform | **No on-demand backup trigger** — no API or UI | Can't take a pre-deployment snapshot | Add `POST /admin/backup` endpoint |
| F-2902 | ⚠️ MAJOR | J29 | Platform | **No storage usage dashboard** | Storage costs grow without visibility | Add storage usage page per category |
| F-2903 | ⚠️ MAJOR | J29 | `internal/modules/media/` | **No orphaned file detection** | Storage waste from orphaned uploads | Add orphan detection + cleanup tool |
| F-2904 | 💡 IMPROVE | J29 | Platform | **No backup listing** — can't see existing backups | Can't verify backup health | Add backup list with metadata |
| F-2905 | 💡 IMPROVE | J29 | Platform | **No backup download** — can't retrieve from admin UI | Must access storage directly | Add one-time download URL |
| F-2906 | 💡 IMPROVE | J29 | Platform | **No scheduled backup configuration** | Backups may be inconsistent | Add backup schedule UI |
| F-2907 | 💡 IMPROVE | J29 | Platform | **No CDN cache invalidation UI** | Delayed content updates | Add "Purge CDN for file" button |
| F-2908 | 💡 IMPROVE | J29 | Platform | **No database maintenance UI** — analyze, vacuum, reindex | DB performance degrades | Add "Run Maintenance" button |
| F-2909 | 💡 IMPROVE | J29 | Platform | **No export/import for catalog data** | Manual data migration | Add catalog export/import wizard |
| F-2910 | 💡 IMPROVE | J29 | Platform | **Soft-delete purge missing** | Storage grows after deletion | Add cleanup job with configurable TTL |

## Phase 32: Advanced Administration

| # | Severity | Journey | Location (file:line) | Problem | User Impact | Fix |
|---|----------|---------|---------------------|---------|-------------|-----|
| F-3001 | ⚠️ MAJOR | J30 | All admin lists | **No bulk operations on tracks/users/albums** — only moderation has bulk | Slow workflows for large catalogs | Add multi-select + bulk toolbar |
| F-3002 | ⚠️ MAJOR | J30 | `internal/app/worker.go` | **No worker monitoring UI** — can't see worker status or errors | Silent worker failures undetected | Add worker status page |
| F-3003 | ⚠️ MAJOR | J30 | `PageAdminImport.vue` | **Import has no queue or progress** | Admin unsure if import is working | Add import queue with progress |
| F-3004 | 💡 IMPROVE | J30 | Enrichment | **No enrichment status dashboard** | Incomplete metadata unnoticed | Add enrichment coverage report |
| F-3005 | 💡 IMPROVE | J30 | All admin | **No bulk operation progress bar** — text-only | No visual feedback during batches | Add progress bar with ETA |
| F-3006 | 💡 IMPROVE | J30 | Admin | **No keyboard shortcuts** in admin pages | Slow for power users | Add admin-wide shortcuts |
| F-3007 | 💡 IMPROVE | J30 | Import | **No import history** | Accidental duplicate imports | Add import log |
| F-3008 | 💡 IMPROVE | J30 | Workers | **No "Run Now" for workers** | Must wait for next interval | Add manual worker trigger |
| F-3009 | 💡 IMPROVE | J30 | Workers | **No worker interval configuration** | Can't adjust frequency without redeploy | Add interval config in settings |
| F-3010 | 💡 IMPROVE | J30 | Dashboard | **No "recent actions" widget** on dashboard | Can't see latest admin activity | Add recent actions feed |
| F-3011 | 💡 IMPROVE | J30 | All admin | **No bulk edit dialog** — can't edit multiple tracks' genres | One-by-one for large catalogs | Add bulk edit modal |
| F-3012 | 💡 IMPROVE | J30 | Admin | **No admin tour / onboarding** | Learning curve for admin tools | Add onboarding overlay |
| F-3013 | 💡 IMPROVE | J30 | Enrichment | **No enrichment rollback** | Data quality risk | Add "Undo enrichment" |
| F-3014 | 💡 IMPROVE | J30 | Admin | **No admin dark/light mode** | Admin may prefer different theme | Add independent admin theme |

---

## Summary Statistics

| Severity | Count |
|----------|-------|
| 🚨 UX-BLOCKER | 14 |
| ⚠️ UX-MAJOR | 104 |
| 💡 UX-IMPROVE | 222 |
| ✨ UX-DELIGHT | 3 |
| **Total** | **343** (was 285 at end of Phase 6) |

> Last updated: 2026-07-06
