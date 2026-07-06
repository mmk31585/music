# Journey 03: Playback (Core User Journey)

## Entry Points
- **Click play on any track** in search results, homepage, artist page, album page, playlist
- **Tap now-playing bar play/pause** to resume/pause current track
- **Keyboard shortcuts** (space = toggle play, arrow keys = seek)
- **Media Session API** notifications / lock screen controls
- **PiP / FloatingMiniPlayer** controls
- **Radio mode**: seed track → auto-generated queue
- **Auto-play**: next track in queue when current ends

## Step-by-Step Walkthrough (as implemented)

### Step 1: Play Initiation
| Entry | Trigger | Component → Store → Engine |
|-------|---------|---------------------------|
| Home carousel | `@play` → `handlePlay(item)` → `player.toggleTrack(buildPlaybackTrack(item))` | `PageHome.vue:442-444` |
| Search results | `@click` → `playTrack(track)` → `player.setQueueAndPlay([track], 0)` | `PageSearch.vue:650-663` |
| Track detail page | "Play" button → `togglePlay()` → `player.playTrack(pb)` or `player.pause()` | `PageTrack.vue:740-757` |
| Album/Artist | "Play All" → `buildQueue()` → `player.setQueueAndPlay(queue, 0)` | `PageAlbum.vue:432-436`, `PageArtist.vue:221-233` |
| NowPlayingBar | Play button → `togglePlayPause()` → `player.resume()` or `player.pause()` | `usePlayerControls.ts:18-25` |

**Key**: Every page builds `PlaybackTrack` objects differently — no common factory function. This leads to duplicate code and potential inconsistencies.

### Step 2: Store → Engine Handoff
| Aspect | Detail |
|--------|--------|
| **Entry point** | `playerStore.playTrack(track)` or `playerStore.setQueueAndPlay(tracks, index)` |
| **Lazy init** | `initialize()` checks if `PlayerEngine` exists — creates on first use |
| **Error guard** | `resetFailureGuard()` resets `consecutiveFailures` counter before each new play |
| **Queue persistence** | Auto-saves track IDs to localStorage (`player-queue-track-ids`) after each change |
| **Files** | `stores/player.ts:233-337` (initialize), `stores/player.ts:394-437` (playTrack/setQueueAndPlay) |

### Step 3: Player Engine (`PlayerEngine`)
| Aspect | Detail |
|--------|--------|
| **`play(track)`** | Sets track on QueueManager → loads audio src → calls `audio.play()` |
| **`setQueueAndPlay(tracks, idx)`** | Replaces queue → loads track at index → plays |
| **Next/Previous** | Delegates to QueueManager (history-aware for previous, shuffle-aware) |
| **Track change** | Emits `trackchange` event → store sets `currentTrack` |
| **Time updates** | Emits `timeupdate` every ~100ms via `requestAnimationFrame` loop |
| **Error handling** | Emits `error` → store catches → skips to next track (up to 3 consecutive failures) |
| **Media Session** | Updates `navigator.mediaSession` metadata + action handlers |
| **Files** | `services/player/player-engine.ts` (585 lines) |

### Step 4: Audio Engine (`AudioEngine`)
| Aspect | Detail |
|--------|--------|
| **Core** | Wraps a single `HTMLAudioElement` with `crossOrigin='anonymous'` |
| **Progress loop** | Uses `requestAnimationFrame` for smooth time updates (every ~100ms) |
| **AudioContext** | Created lazily on first play — enables analyser for visualizer |
| **Events** | Abstracted via `.on()` pattern: play, pause, ended, waiting, playing, canplay, timeupdate, durationchange, error |
| **Files** | `services/player/audio-engine.ts` (312 lines) |

### Step 5: Now Playing UI (NowPlayingBar)
| Aspect | Detail |
|--------|--------|
| **Desktop** | Fixed bottom bar with two modes: collapsed (mini) and expanded (full controls) |
| **Collapsed** | Shows cover art, title, artist, play/pause button, expand chevron — thin progress line |
| **Expanded** | Shows album art, progress slider, all controls (prev/play/next, shuffle/repeat modes, volume slider, queue preview, lyrics toggle, PiP, overflow menu) |
| **Mobile** | Floating bar above bottom nav with cover, title, play/pause, next |
| **Shuffle modes** | Off, Queue, Random Catalog, Similar Tracks — popup selector |
| **Repeat modes** | Off, All, One — popup selector |
| **Queue preview** | Hover popup showing now-playing + next track |
| **Error state** | Red banner below bar when playback error occurs |
| **Files** | `components/music/player/NowPlayingBar.vue` (871 lines) |

**Key**: NowPlayingBar is 871 lines — very large component handling desktop collapsed/expanded, mobile, queue preview, shuffle/repeat popups, overflow menu, guest play gates, and error states. Could benefit from decomposition.

### Step 6: Fullscreen Player
| Aspect | Detail |
|--------|--------|
| **Trigger** | Click expand button on NowPlayingBar, or tap mobile bar |
| **Layout** | Teleported to body, fullscreen with blurred cover art background |
| **Tabs** | Now Playing, Queue, Lyrics — toggle at top |
| **Cover** | Large album art with pulsing glow when playing |
| **Controls** | Play/pause, next/prev, progress slider, volume, shuffle/repeat |
| **Lyrics panel** | Karaoke-style synced lyrics or plain text — synced to playback position |
| **PiP support** | Pop-out to Picture-in-Picture via PiPPlayerContent |
| **Files** | `components/music/player/FullscreenPlayer.vue` (1016 lines) |

### Step 7: Queue Management
| Aspect | Detail |
|--------|--------|
| **Source** | `QueuePanel.vue` — slide-out panel showing current queue |
| **Features** | Drag-to-reorder tracks, remove from queue, play now, play next |
| **Persistence** | Track IDs saved to `localStorage`, hydrated on app start via `restorePersistedQueue()` |
| **Shuffle order** | `PlayerEngine` maintains internal `shuffleOrder` array for queue shuffle |
| **Files** | `components/music/player/QueuePanel.vue`, `stores/player.ts:348-386` |

### Step 8: Track End / Auto-Advance
| Aspect | Detail |
|--------|--------|
| **Natural end** | `audio.on('ended')` → `PlayerEngine.next()` |
| **Repeat one** | If repeatMode === 'one' → `seek(0)` and play again |
| **Error skip** | Engine emits 'error' → store increments failure counter → skips to next |
| **Stop condition** | 3 consecutive failures → `playbackStopped = true` → no more auto-advance |
| **Play history** | `onPlayHistory` provider called with track ID + duration |
| **Music status** | Debounced (3s) POST to `/api/v1/users/me/music-status` for social presence |

## State Matrix Findings

### NowPlayingBar
| State | Present? | Notes |
|-------|----------|-------|
| 🟢 Empty (no track) | ✅ | `v-if="currentTrack"` — bar is hidden entirely |
| 🟢 Playing | ✅ | Animated progress, play icon → pause, glow effects |
| 🟢 Paused | ✅ | Progress stops, pause icon → play, dimmed cover |
| 🟢 Buffering/loading | ✅ | Spinner on play button, `isLoadingTrack` flag |
| 🟢 Error | ✅ | Red banner below bar with error message |
| 🔴 Offline | Missing | No offline indicator in player bar |
| 🟢 Guest gate | ✅ | `GuestPlayGate` wraps play buttons — redirects to login for guests |
| 🟢 Collapsed/expanded | ✅ | User preference persisted to localStorage |
| 🟢 Queue empty | ✅ | Shows "No upcoming tracks" message |
| 🟢 Shuffle random mode | ✅ | Shows special icon for catalog/similar modes |

### Player Engine
| State | Present? | Notes |
|-------|----------|-------|
| 🟢 First play (lazy init) | ✅ | Engine created on first `initialize()` call |
| 🟢 Track change | ✅ | Emits `trackchange`, updates Media Session |
| 🟢 Queue empty | ✅ | Returns gracefully with no-ops |
| 🟢 Consecutive failure guard | ✅ | Stops after 3 failures |
| 🟢 Repeat modes | ✅ | Off, One, All |
| 🟢 Shuffle modes | ✅ | Off, Queue, Catalog, Similar |
| 🔴 Retry on transient error | Missing | No retry logic — first error skips track (I-007 from BACKLOG) |
| 🔴 Event listener leak | ⚠️ | Anonymous arrow functions in bindAudioEvents — `dispose()` may not clean all (I-004) |
| 🟢 Queue persistence | ✅ | Track IDs saved to localStorage, hydrated on init |

### FullscreenPlayer
| State | Present? | Notes |
|-------|----------|-------|
| 🟢 Open/close | ✅ | Teleport + transition animation |
| 🟢 Playing/paused | ✅ | Visual state reflected in cover scale/opacity |
| 🟢 Tab switching | ✅ | Now Playing, Queue, Lyrics |
| 🟢 Empty queue | ✅ | Shows appropriate state |
| 🔴 Error state in fullscreen | 🟡 | May rely on NowPlayingBar error banner — not independently handled |
| 🟢 PiP toggle | ✅ | Button with active state |
| 🟢 Lyrics karaoke | ✅ | Synced lyrics with highlight (KaraokeLyrics component) |

## Friction Points

| # | Severity | Location | Problem | User Impact |
|---|----------|----------|---------|-------------|
| 03.01 | ⚠️ MAJOR | All pages building PlaybackTrack | No shared factory — every page builds `PlaybackTrack` differently | Inconsistencies in track metadata, missing fields, duplication |
| 03.02 | ⚠️ MAJOR | `NowPlayingBar.vue` (871 lines) | Single component is too large — handles desktop collapsed/expanded, mobile, queue preview, 2 popup menus, overflow, error, PiP | Maintainability risk, slow re-renders on state changes |
| 03.03 | ⚠️ MAJOR | `FullscreenPlayer.vue` (1016 lines) | Also extremely large — combines cover view, controls, queue tab, lyrics tab, PiP | Same maintainability concerns |
| 03.04 | ⚠️ MAJOR | `stores/player.ts:289-297` | Anonymous arrow event listeners in `initialize()` — `unsubs` array may leak | Double-binding on hot reloads, stale listeners |
| 03.05 | ⚠️ MAJOR | `services/player/player-engine.ts` | No retry on transient playback error — immediate skip to next track (I-007) | Users hear skip-glitch on temporary network issues |
| 03.06 | 💡 IMPROVE | `stores/player.ts:349-359` | Queue persists full ID list but not metadata — re-fetches on restore | Slow queue restoration for large queues |
| 03.07 | 💡 IMPROVE | `NowPlayingBar.vue` | Volume slider is `dir="ltr"` hardcoded — doesn't respect RTL | RTL users see inverted slider behavior |
| 03.08 | 💡 IMPROVE | `FullscreenPlayer.vue` | Lyrics tab loads synchronously — large lyric content blocks UI thread | Stuttering on lyrics-rich tracks |
| 03.09 | 💡 IMPROVE | `media-session.ts` | Media Session action handlers use `navigator.mediaSession` — may not work in all browsers (Safari PiP) | Inconsistent cross-browser behavior |
| 03.10 | 💡 IMPROVE | `NowPlayingBar.vue` | Shuffle/repeat popup menus close on outside click — but click-outside handler is a document-level listener | Performance overhead, potential event conflicts |
| 03.11 | ✨ DELIGHT | `stores/player.ts:169-198` | Music status sharing (debounced, privacy-aware) | Currently best-effort — could be promoted with UI indicator |

## Clicks & Time-to-Value

| Metric | Current | Ideal |
|--------|---------|-------|
| Find a track → hear audio | ~2-5s (page load + API + buffer) | <1s for cached tracks |
| Pause/Resume | ~100ms (state toggle) | Instant |
| Next track | ~200-800ms (engine.next + buffer) | <300ms |
| Fullscreen player open | ~400ms (animation) | <200ms |
| Queue reorder | ~100ms (drag) | Instant |
| **Playback start** (cold) | ~1-3s (API fetch + buffer) | <500ms target |
| **Playback start** (preloaded) | ~100-300ms | <100ms |

## RTL / A11y / Mobile Notes

- **RTL**: Volume slider in NowPlayingBar hardcodes `dir="ltr"` — correct for volume direction, but the progress seek bar also has `dir="ltr"`. For RTL users, seeking should logically reverse (right-to-left = forward).
- **ARIA labels**: Play/pause button has dynamic `aria-label` ('Loading'/'Pause'/'Play'). Previous/next/shuffle/repeat all have labels. Queue items have appropriate roles.
- **Touch targets**: All controls use min 36px (h-9 = 36px, h-11 = 44px). Mobile bar buttons are 36px which is below recommended 44px for touch.
- **Keyboard**: Space toggles play/pause. Arrow keys for seek. Enter activates buttons. Focus order follows visual order.
- **Mobile bottom nav**: `MobileBottomNav` sits at 4.5rem height with safe-area aware positioning. Player bar above it uses `bottom: calc(4.5rem + env(safe-area-inset-bottom))`.
- **prefers-reduced-motion**: Handled in NowPlayingBar — disables all animations and transitions.
- **Focus visible**: Interactive elements have hover/focus styles. Skip link present in LayoutMusicApp.
- **Contrast**: Error banner uses `text-red-400` on `bg-red-500/10` — need to verify 4.5:1. Buttons with `text-white/40` on dark backgrounds may fail contrast.

## Delight Opportunities

- ✨ **Ambient mode**: Color-matched background from cover art palette (already partially done in player engine)
- ✨ **Preload next track intelligently**: Already exists as `PreloadManager` — could be more aggressive
- ✨ **Crossfade between tracks**: `crossfadeDuration` already in store — enable by default
- ✨ **Mini-player animations**: Vinyl spin effect or equalizer animation on now-playing bar
- ✨ **Sleep timer UI**: Already in store — show in player overflow menu
- ✨ **Social listening status**: Let users see what friends are listening to (music-status endpoint exists)

## Open Questions

1. Should the player maintain its own queue independent of page navigation? (e.g., playing a track from search should not create a queue of all search results)
2. Is the catalog/similar shuffle mode working end-to-end? It requires backend endpoints for random/similar tracks — need to verify they exist.
3. What is the intended UX for the floating mini player (PiP) — is it a separate window or in-page component?
4. Should playback continue when navigating to admin pages? Currently LayoutAdmin includes NowPlayingBar and FullscreenPlayer, suggesting yes.
