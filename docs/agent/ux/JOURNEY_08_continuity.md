# Journey 08: Continuity — Session, Refresh, State Restore

> What happens when the user refreshes, closes/reopens the browser, switches devices,
> encounters a session expiry mid-listening, or the token refresh fails during playback.

---

## Player State Persistence Map

| State | Persisted? | Restored? | On What Trigger |
|-------|-----------|-----------|-----------------|
| Queue (track IDs) | ✅ localStorage | ✅ Re-fetched metadata | `player.initialize()` |
| Queue position | ❌ | ❌ Always index 0 | — |
| Current time | ❌ | ❌ Always 0:00 | — |
| Volume | ✅ localStorage | ✅ | `player.initialize()` |
| Muted | ✅ localStorage | ✅ | `player.initialize()` |
| Shuffle mode | ❌ | ❌ Resets to `off` | — |
| Repeat mode | ❌ | ❌ Resets to `off` | — |
| Playback speed | ❌ | ❌ Resets to 1x | — |
| Audio quality | ✅ Pinia | ✅ | `player.initialize()` |
| Sleep timer | ❌ | ❌ | — |
| Crossfade duration | ❌ | ❌ Resets to 0 | — |
| `isPlaying` state | ❌ | ❌ | — |
| Current track | ❌ | ❌ | — |
| Liked tracks cache | ⚠️ In-memory `Set` | ❌ Re-fetched | On first `useTrackLike()` call |
| Player bar collapsed | ✅ localStorage | ✅ | Component initialize |

---

## Step-by-Step: Page Refresh Mid-Playback

### Before Refresh
```
User is listening to track 7 of a 15-track album
  → player.currentTime = 134s (2:14 into the song)
  → player.queue = [15 tracks]
  → player.isPlaying = true
  → player.volume = 0.7
  → player.shuffleMode = 'queue'
  → player.repeatMode = 'all'
```

### On Refresh (`window.location.reload()`)
```
1. Browser navigates away
   → pagehide event fires
   → registerBeforeUnload handler fires:
     sendBeacon('/api/v1/history/music-status', { track_id, ... })
     → Best-effort, fire-and-forget

2. All in-memory Pinia state is destroyed
3. localStorage persists: queue-track-ids, volume, muted
```

### After Refresh (Page Loads)
```
4. Vue app boots
5. Pinia stores initialize with defaults:
     currentTrack: null
     queue: []
     isPlaying: false
     currentTime: 0
     duration: 0
     shuffleMode: 'off'
     repeatMode: 'off'
     volume: 0.85 (default — will be overridden)
     muted: false

6. Player engine not initialized yet

7. User navigates to any page that uses usePlayer()
   → usePlayer() calls player.initialize()

8. player.initialize():
   a. Creates PlayerEngine singleton (if first time)
   b. Calls restorePersistedQueue()
      → Reads 'player-queue-track-ids' from localStorage
      → Makes N sequential API calls to getPlaybackTrack(id)
      → Returns restored tracks
   c. Sets engine.queue with restored tracks
   d. Restores volume from localStorage
   e. Restores muted from localStorage
```

### What User Sees

| Time | State | UI |
|------|-------|-----|
| 0ms (page load) | Everything at default | No player bar (no current track) |
| ~500ms (app boot) | Volume restored (invisible) | No change yet |
| ~1-3s (user navigates to home/album/etc.) | `usePlayer()` called, queue restored | Queue is back but nothing playing |
| User clicks play | Normal play flow | Player bar appears, audio starts |

**The user's session is NOT restored to its pre-refresh state.** The queue is back (if it was persisted), but:
- Playback does NOT auto-resume
- Current track position is lost
- Shuffle/repeat modes are reset
- User must find a play button and click it again

---

## Step-by-Step: Token Expiry Mid-Playback

### Scenario: User has been listening for 70+ minutes (access token expires)

```
1. Audio continues playing (no API calls needed for playback)
2. Background API call happens (e.g., music status update, like, etc.)
3. Request returns 401
4. Axios interceptor catches 401:
   a. Tries to refresh token: POST /api/v1/auth/refresh
   b. If refresh succeeds → retries original request
   c. If refresh fails → clears auth state
5. Audio continues playing uninterrupted
```

**Key insight**: Audio playback uses `HTMLAudioElement` with a direct stream URL. It does NOT go through the API client. **Token expiry does NOT interrupt audio playback.** The user can continue listening indefinitely even with an expired token.

The expiry only affects:
- Like/unlike operations (silently fail)
- Library operations
- Music status updates (silently fail)
- Any API call that checks auth

### If Refresh Fails Mid-Session

```
1. Auth state is cleared
2. GuestPlayGate activates on next user interaction
3. User can still finish the current track (audio continues)
4. If user clicks "next track" → stream URL fetch may fail → playback stops
```

---

## Step-by-Step: Browser Tab Switch

### With PiP

```
User is listening in FullscreenPlayer
  → Tab switch triggers 'visibilitychange'
  → FullscreenPlayer detects document.hidden === true
  → If playing, auto-opens PiP window
  → PiP window shows compact player (always-on-top)
  → User can control playback from PiP
  → User returns to tab → PiP may auto-close (depends on browser)
```

### Without PiP

```
User is listening (any mode)
  → Tab switch → audio continues (HTMLAudioElement is not paused)
  → Media Session updates lock screen with current track info
  → User can control via lock screen (play/pause/next/prev/seek)
  → User returns to tab → full experience restored
```

**Note**: `HTMLAudioElement` is NOT paused on tab switch unless the user explicitly pauses or the browser throttles the tab. Most modern browsers do NOT throttle audio-playing tabs.

---

## Step-by-Step: Incognito / Private Browsing

```
1. localStorage is available BUT is isolated to the session
2. All queue/volume/muted persistence works within the session
3. On tab close → localStorage is wiped
4. On new incognito session → everything starts from defaults
5. Queue, volume, preferences are all lost
```

This means: In incognito mode, every page refresh within the session preserves state, but closing the tab loses everything. This is expected browser behavior.

---

## Step-by-Step: Device Switch

```
User listens on Desktop (Chrome)
  → Queue: [album tracks], Volume: 0.7, Shuffle: on
  → Closes browser

User opens Mobile (Safari)
  → No queue restored (different browser, different localStorage)
  → Volume defaults to 0.85
  → Shuffle defaults to 'off'
  → Nothing is playing
```

**Zero cross-device state sharing**. No server-side queue, volume, or playback position is stored. Each device is a completely fresh experience.

---

## Step-by-Step: Network Interruption

### Audio Streaming Interrupted

```
Network drops mid-stream
  → HTMLAudioElement fires 'waiting' event
  → isBuffering = true (UI shows spinner)
  → Browser buffers indefinitely (up to its internal timeout)
  → If network returns within buffer window → 'playing' fires → resume normal
  → If network stays down long enough → 'error' fires (MEDIA_ERR_NETWORK)
  → Consecutive failure guard fires:
    → If < 3 failures: auto-skip to next track → network still down → another error → skip → ...
    → If >= 3 failures: playbackStopped = true, engine.stop()
    → User sees "Playback unavailable — the media server may be offline"
```

### API Request Interrupted (during play)

```
User clicks play
  → store.setQueueAndPlay() called
  → API is down (stream URL not available)
  → engine.play(track) → audio.load() hangs on network
  → Eventually 'error' fires → same chain as above
  → User sees error in player bar
```

**No retry on network error** for the current track — it's immediately skipped. On the **next** play action, `resetFailureGuard()` clears the counter and tries again.

---

## State Matrix: What Survives What

| Event | Queue | Volume | Playing | Current Time | Shuffle | Repeat |
|-------|-------|--------|---------|-------------|---------|--------|
| Page refresh | ✅ (IDs) | ✅ | ❌ | ❌ | ❌ | ❌ |
| Tab switch (same session) | ✅ (memory) | ✅ (memory) | ✅ | ✅ | ✅ | ✅ |
| Browser close → reopen | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| Device switch | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Incognito close | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Token refresh failure | ✅ (memory) | ✅ | ✅ (may stop on next) | ✅ | ✅ | ✅ |
| Network interruption | ✅ (memory) | ✅ | ❌ (may fail) | ❌ (resets) | ✅ | ✅ |
| Auth logout | ❌ (cleared) | ✅ | ❌ | ❌ | ❌ | ❌ |
| App update (HMR) | ⚠️ (depends) | ❌ | ❌ | ❌ | ❌ | ❌ |

---

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-801 | 🚨 BLOCKER | `stores/player.ts` — `initialize()` | **Playback does not auto-resume after page refresh** — queue is restored but user must manually click play | User loses their listening session on accidental refresh | Add "Resume playback" option on restore |
| F-802 | 🚨 BLOCKER | `stores/player.ts` — persistence | **Current time not persisted** — track restarts from 0 after refresh | User loses their place in a long song (e.g., podcast, classical) | Persist current time with timestamp, resume from position |
| F-803 | ⚠️ MAJOR | Queue persistence | Queue position not persisted — always restores to index 0 (F-703 continuation) | User loses place in long queue | Persist queue index |
| F-804 | ⚠️ MAJOR | `stores/player.ts` — persistence | Shuffle/repeat modes not persisted — reset on every new queue and page refresh | User's preferred listening mode is forgotten | Persist to localStorage |
| F-805 | ⚠️ MAJOR | Queue persistence | Queue persistence uses **sequential** `getPlaybackTrack()` calls — slow for large queues | 15-track queue = 15 sequential API calls; UI blocks restoring | Batch fetch or cache metadata |
| F-806 | ⚠️ MAJOR | `player-engine.ts` | Audio error mid-playback auto-skips without retry — consecutive failures quickly block playback | Brief network glitch → 3 skipped tracks → "Playback unavailable" | Add retry with exponential backoff before skip |
| F-807 | 💡 IMPROVE | Player store | Queue persistence has `hydratedAt` timestamp stored but never checked | Stale queue data never expires | Add TTL-based queue expiry |
| F-808 | 💡 IMPROVE | Player store | No server-side queue sync | Zero cross-device continuity | Sync queue to user account on server |
| F-809 | 💡 IMPROVE | `FullscreenPlayer.vue` | Auto-PiP on tab switch has no opt-in — happens automatically with no user choice | Surprising behavior for new users | Add "Enable Auto-PiP" setting |
| F-810 | 💡 IMPROVE | `registerBeforeUnload` | `pagehide` handler sends `sendBeacon` — but if browser kills the tab quickly, this may not fire | Some session data lost on abrupt close | Use `navigator.sendBeacon` with `keepalive: true` on `visibilitychange` |
| F-811 | 💡 IMPROVE | Player engine | No "resume playback" banner/dialog on page load after crash | User returns to app and sees nothing playing — they may not remember what they were listening to | Offer smart resume: "Continue listening to [album]" |
| F-812 | 💡 IMPROVE | `player.initialize()` | `restorePersistedQueue()` fails silently if `localStorage` is full or unavailable | User's queue silently disappears | Log warning, show toast |

## RTL / A11y / Mobile Notes

- ✅ Audio continues playing in background tabs (native HTMLAudioElement behavior)
- ✅ Media Session provides lock screen controls (play/pause/next/prev)
- ✅ PiP provides always-on-top controls on supporting browsers
- ❌ No "resume on app open" feature — user must navigate back and click play
- ❌ No offline indication — user can't tell if playback failed due to network or other error
- ✅ `sendBeacon` is used for fire-and-forget music status updates (works even during page unload)

## Delight Opportunities

- ✨ **Smart Resume**: On app open, show a toast: "Continue listening to [album title]?" with play button
- ✨ **Cross-Device Queue**: Save queue to user account, sync across devices
- ✨ **Playback History Timeline**: Show a visual timeline of what was listened to, with "resume from where you were"
- ✨ **Offline Playlist Hint**: When network fails mid-playlist, show "You were listening to [album]. Would you like to download it for offline?"
