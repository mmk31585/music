# Journey 05: Find & Play a Track

> Maps the full path from search query → results → click → first audio playback.
> Trace: `PageSearch.vue` → `useSearchApi.search()` → results render → `playTrack()` → `player.setQueueAndPlay()` → `engine.play()` → `audio.play()`.

---

## Entry Points

A user can initiate playback from **7+ distinct entry points**:

| Entry Point | Page | Action | Queue Built |
|-------------|------|--------|------------|
| Search result track | `PageSearch.vue` | Click track row | Single track only |
| Search "Trending" track | `PageSearch.vue` (discover) | Click `HomeTrackCard` | Single track only |
| Search "Viral Hits" track | `PageSearch.vue` (discover) | Click track row | Single track only |
| Album track row | `PageAlbum.vue` | Click track row | Full album (starts at clicked) |
| Album "Play All" | `PageAlbum.vue` | Click play all button | Full album (starts at 0) |
| Album "Shuffle" | `PageAlbum.vue` | Click shuffle button | Shuffled album |
| Track detail page | `PageTrack.vue` | Click play button | Single track |
| "You Might Like" track | `PageTrack.vue` | Click similar track | Full similar list (starts at clicked) |
| Playlist track | `PagePlaylistDetail.vue` | Click track row | Full playlist (starts at clicked) |
| Playlist "Play" | `PagePlaylistDetail.vue` | Click play button | Full playlist (starts at 0) |
| Home page card | `PageHome.vue` | Click track card | Single track (via hero) or section queue |

---

## Step-by-Step: Search → Play

### Step 1: User Types in Search Bar

```
Input field focus → Recent searches panel slides in (if available)
User types → 350ms debounce → doSearch() fires
```

| Aspect | Detail |
|--------|--------|
| Debounce | `350ms` after last keystroke |
| Minimum query | Must have non-whitespace characters |
| API call | `searchApi.searchCatalog({ query: q, limit: 20 })` → `GET /api/v1/search?q=...&type=all&limit=20` |
| Loading state | 6 skeleton rows (cover + 2 lines each) |
| Empty result | `AppEmptyState` with "No results found" — search term displayed |
| Error state | Silent `console.error` — results set to empty, `hasNoResults=true` |

### Step 2: Results Render

Results appear in sections:
1. **Songs** — rows with index, cover art, title, artist, duration
2. **Artists** — circular avatar cards linking to artist page
3. **Albums** — rectangular cover art cards linking to album page

Each song row is a `<div role="button" tabindex="0">` with `@click`, `@keydown.enter`, `@keydown.space.prevent` all calling `playTrack(track)`.

**Error/edge**: If the search API fails, `results` is set to empty arrays and `hasNoResults=true` — the user sees "No results found" even though the real issue is a network/server error (no error banner).

### Step 3: User Clicks a Track

```typescript
function playTrack(track: TrackCardItem) {
  void player.setQueueAndPlay(
    [{
      id: String(track.id),
      title: track.title as string,
      artistName: (track.artist_name as string) || 'Unknown',
      albumTitle: (track.album_title as string) || null,
      coverUrl: (track.cover_url as string) || null,
      durationSeconds: (track.duration_seconds as number) ?? null,
      streamUrl: playerApi.getTrackStreamUrl(String(track.id)),
    }],
    0,
  )
}
```

**Key observation**: From search results, a **single-track queue** is built — only the clicked track is in the queue. After it finishes, the engine has nothing to play next (unless repeat-all is on).

### Step 4: Store Processes the Play Request

```typescript
// stores/player.ts — setQueueAndPlay()
async setQueueAndPlay(tracks, startIndex = 0) {
  this.initialize()
  this.resetFailureGuard()
  this.shuffleMode = 'off'       // Always resets shuffle to off for new queue
  this.isLoadingTrack = true
  this.error = null              // Clear previous error
  
  try {
    await this.engine!.setQueueAndPlay(tracks, startIndex)
  } catch (err) {
    this.error = getErrorMessage(err, 'Failed to play track')
  } finally {
    this.isLoadingTrack = false
  }
}
```

### Step 5: Engine Play Pipeline

```
engine.setQueueAndPlay(tracks, startIndex)
  → queue.setQueue(tracks, startIndex)    // Replace queue, set index
  → const track = queue.getCurrent()
  → doPlay(track)
    → emit('buffering', true)              // UI shows buffering
    → currentTrack = track → emit('trackchange')
    → queue.setCurrent(track) → emit('queuechange')
    → updateMediaSession(track, handlers)   // Lock screen controls
    → audio.play(track.streamUrl)
      → audio.src = streamUrl
      → audio.load()                        // Start fetching audio
      → await canplay event                 // WAIT for enough data
      → audio.play()                        // Actually start playback
    → providers.onPlayHistory?.()           // Record play
    → preload.nextTrack(tracks[index+1])    // Preload next track
    → emit('buffering', false)
```

### Step 6: User Sees Playback Start

| Moment | UI State | Duration |
|--------|----------|----------|
| Click happens | Track row shows click feedback | Instant |
| `isLoadingTrack=true` | No visible feedback on search page | 10-50ms |
| `isBuffering=true` | NowPlayingBar shows rotating spinner on play button | 200ms-3s |
| `trackchange` emitted | NowPlayingBar shows new track title, artist, cover | ~100ms after buffer |
| `canplay` fires | Spinner → play icon | Variable (depends on stream) |
| `audio.play()` succeeds | Audio plays, progress bar starts moving | At canplay moment |

### Step 7: If Playback Fails

```
audio.error event → engine emits 'error'
  → store.error handler:
    → sets this.error = message
    → clears isBuffering, isPlaying
    → increments consecutiveFailures
    → if failures < 3: engine.next() (skip to next track, if any)
    → if failures >= 3: playbackStopped = true, engine.stop()
```

**Error messages** (from `audio-engine.ts`):
| Code | Message |
|------|---------|
| `MEDIA_ERR_ABORTED` | "Playback was aborted" |
| `MEDIA_ERR_NETWORK` | "Network error — check your connection" |
| `MEDIA_ERR_DECODE` | "Could not decode audio" |
| `MEDIA_ERR_SRC_NOT_SUPPORTED` | "Audio format not supported or stream unavailable" |
| Other | "Audio playback error" |

**For single-track search play**: If track fails and `consecutiveFailures < 3`, `engine.next()` will find nothing (queue of 1) and playback stops. User sees the error state on NowPlayingBar.

---

## State Matrix Findings

| State | Search Page | Album Page | Track Page | Playlist Page |
|-------|-------------|------------|------------|---------------|
| 🟢 Loading (skeleton) | ✅ 6 rows | ✅ Full page skeleton | ✅ Full page skeleton | ✅ Full page skeleton |
| 🟢 Empty (no results/tracks) | ✅ "No results found" | ✅ "No tracks found" | 🔴 N/A | ✅ Empty state |
| 🟢 Error loading | ✅ Silent (shows empty state) | ✅ Error message | ✅ Error message | ✅ Error message |
| 🟢 Buffering spinner | ✅ (in player bar) | ✅ (in player bar) | ✅ (in player bar) | ✅ (in player bar) |
| 🟢 Playback success | ✅ | ✅ | ✅ | ✅ |
| 🔴 Click feedback on play | ✅ (row highlight) | ✅ (row highlight) | ✅ (button style) | ✅ (overlay play icon) |
| 🔴 No queue after single track | ⚠️ — single track, ends after 1 | ✅ Album fills queue | ✅ Similar fills queue | ✅ Playlist fills queue |
| 🔴 Error display on page | ❌ Silent (no error banner) | ❌ Only in player bar | ❌ Only in player bar | ❌ Only in player bar |
| 🔴 Loading → playing transition flash | ⚠️ — buffering=true → canplay, no intermediate | Same | Same | Same |

---

## Timing & Clicks to Value

### Search → Play

| Step | Approx Time | User Feels |
|------|-------------|------------|
| Type query | 350ms debounce | Responsive |
| API call | 200-800ms | Loading skeleton |
| Results render | ~100ms after API | See results instantly |
| Click track | Instant | Click feedback |
| Buffering | 200ms-3s | Spinner on play button |
| Audio starts | Varies | Hear music |

**Total: ~700ms-4s from first keystroke to audio** (if search API is fast).

### Album/Track/Playlist → Play

| Step | Approx Time | User Feels |
|------|-------------|------------|
| Navigate to page | 500ms-2s (page load) | Skeleton loading |
| Click play/track | Instant | Click feedback |
| Buffering | 200ms-3s | Spinner |
| Audio starts | Varies | Hear music |

**Total: ~700ms-5s from page entry to audio**.

---

## RTL / A11y / Mobile Notes

- ✅ Search input uses `dir="ltr"` — correct for search queries
- ✅ `aria-label` on search input, clear button, play buttons
- ✅ `role="button"` + `tabindex="0"` on all clickable rows
- ✅ Empty state has informative descriptions
- ❌ **No keyboard shortcut to focus search** — user must tab to it
- ❌ **No "play" announcement** — screen readers aren't notified when playback starts (no live region on NowPlayingBar)
- ❌ **Search page error is silent** — if API fails, user sees "No results found" which is misleading
- ✅ Touch targets are adequate (44px+ for buttons, rows are large enough)

## Delight Opportunities

- ✨ **Instant play from search**: Preload the top result so first click is instant (no buffering)
- ✨ **"Play all" in search results**: Allow playing all search results as a queue
- ✨ **Search suggestions**: Typeahead dropdown with popular searches as user types
- ✨ **Visual feedback on play**: Show NowPlayingBar slide up with animation when first track starts playing
- ✨ **Error with retry**: When a track fails, show a "Retry" button in the player bar

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-501 | ⚠️ MAJOR | `PageSearch.vue:588-591` | Search API error silently swallowed — shows "No results found" | User thinks their search has no results, not a network error | Add error banner with retry |
| F-502 | ⚠️ MAJOR | `PageSearch.vue:650-663` | Search play creates single-track queue — track ends, nothing next | After one song, player stops | Queue more from same artist/album |
| F-503 | ⚠️ MAJOR | `PageAlbum.vue:402-405` | Album `buildQueue()` re-maps on every click — wasteful if clicking multiple tracks | Needless re-computation | Cache the queue when tracks haven't changed |
| F-504 | 💡 IMPROVE | All play initiators | No shared `buildPlaybackTrack()` factory — each page constructs PlaybackTrack differently (F-019 continuation) | Inconsistencies in fallback fields, missing albumTitle, etc. | Create shared factory |
| F-505 | 💡 IMPROVE | `PageSearch.vue` | No direct "Play All" on search results | User must click each track individually | Add "Play All" button |
| F-506 | 💡 IMPROVE | `PageSearch.vue:36-38` | Search placeholder is English "What do you want to listen to?" only | Persian users see English | Add locale-aware placeholder |
| F-507 | 💡 IMPROVE | All play buttons | No haptic/visual "playback started" confirmation on the page | User may wonder if click registered | Brief ripple or "Now Playing" badge on clicked row |
| F-508 | 💡 IMPROVE | Player engine `doPlay()` | No transition tracking (playback-start analytics event) | Can't measure time-to-play | Fire analytics event on first audio frame |
| F-509 | 💡 IMPROVE | `PageTrack.vue` | `togglePlay()` — if track is already playing, pauses. But if user clicked from another source, previous queue is replaced silently | User's carefully ordered queue is replaced by single track | Queue insertion vs replacement option |
| F-510 | 💡 IMPROVE | `PageAlbum.vue` `shuffleAll()` | Uses Fisher-Yates on the client | Different shuffle order on different devices for same album | Server-seeded shuffle or saved shuffle order |
