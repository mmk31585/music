# Journey 07: Queue Management

> Full trace of queue operations: add, remove, reorder, persist, restore, cross-source behavior.

---

## Queue Data Model

```typescript
class QueueManager {
  private queue: PlaybackTrack[] = []      // Linear ordered list
  private index = -1                        // Current position in queue
  private history: HistoryEntry[] = []      // Play history stack
}

interface HistoryEntry {
  trackId: string
  queuePosition: number
}
```

The queue is a **simple linear array** with a current index pointer. The history stack enables correct "previous" navigation even after shuffle or manual jumps.

---

## Queue Operations

### View Queue

| Surface | How to Open | Shows |
|---------|-------------|-------|
| QueuePanel | Click queue icon in NowPlayingBar / press `Q` | Right-side slide-in panel: "Now Playing" + draggable "Up Next" |
| FullscreenPlayer tab | Open fullscreen, click "Queue" tab | Queue list with drag reorder |
| ExpandedPlayer tab | Same | Same |

### Add to Queue

There are **two concepts**: "Play Next" and "Add to Queue". The user's `useTrackContextMenu` composable provides both:

| Action | Position | Effect |
|--------|----------|--------|
| **Play Next** | Inserts at `index + 1` | Track plays immediately after current one |
| **Add to Queue** | Appends to end | Track plays after all other queued items |

**Not currently implemented**: Neither "Play Next" nor "Add to Queue" appears in the UI context menus or any visible control. The queue is populated only by `setQueueAndPlay()` from pages (album, playlist, search, etc.) and by shuffle/catalog/similar modes. Individual track "add to queue" functionality does not exist in the current UI.

### Reorder Queue

| Surface | Mechanism | Status |
|---------|-----------|--------|
| QueuePanel | `vuedraggable` drag-and-drop | ✅ Works |
| FullscreenPlayer queue tab | Same `QueuePanel` component | ✅ Works |
| ExpandedPlayer queue tab | Same | ✅ Works |

**Reorder behavior** (`reorderQueue(oldIndex, newIndex)`):
```typescript
reorderQueue(oldIndex: number, newIndex: number) {
  const [moved] = this.queue.splice(oldIndex, 1)
  this.queue.splice(newIndex, 0, moved)
  // Corrects this.index if the current track's position has changed
  if (oldIndex < this.index && newIndex >= this.index) this.index--
  else if (oldIndex > this.index && newIndex <= this.index) this.index++
  else if (oldIndex === this.index) this.index = newIndex
}
```

### Remove from Queue

Each track in the queue panel has a remove button. Calling remove:
1. Splice the track from `queue` array
2. Adjust `index` if the removed track was before or at the current position
3. If current track was removed, start playing the next track (or stop)

**Note**: Removing the **current track** causes an immediate advance to the next track — no confirmation.

### Clear Queue

```typescript
clear() {
  this.queue = []
  this.index = -1
  this.history = []
}
```
Called when `setQueueAndPlay()` replaces the queue, or on `engine.stop()`.

---

## Queue Persistence

### Storage

```typescript
// Persist (on every queuechange event)
localStorage.setItem('player-queue-track-ids', JSON.stringify({
  trackIds: queue.map(t => t.id),
  hydratedAt: Date.now(),
}))

// Restore (on player.initialize())
async restorePersistedQueue(): Promise<PlaybackTrack[]> {
  const raw = localStorage.getItem('player-queue-track-ids')
  if (!raw) return []
  
  const { trackIds } = JSON.parse(raw)
  const tracks = await Promise.all(
    trackIds.map(id => 
      playerApi.getPlaybackTrack(id).catch(() => null)
    )
  )
  return tracks.filter(Boolean) as PlaybackTrack[]
}
```

### What's Persisted

| Data | Persisted? | Key | When |
|------|-----------|-----|------|
| Track IDs in queue | ✅ | `player-queue-track-ids` | Every `queuechange` event |
| Full track metadata | ❌ | — | Must re-fetch on restore |
| Current position in queue | ❌ | — | Always restores to index 0 |
| Queue play history stack | ❌ | — | Lost on refresh |
| Current time position | ❌ | — | Always starts from beginning |

### Persistence Limitations

1. **Only IDs, not metadata**: On restore, ALL track metadata is re-fetched via individual `getPlaybackTrack(id)` API calls. This means:
   - Slow restoration for large queues (sequential API calls)
   - Failed fetches skip that track silently
   - If the API is down, the queue is empty on restore

2. **No position**: Queue always restores to position 0. If user was on track 15 of a 20-track queue, refresh resets to track 0.

3. **No currentTime**: Position in the current song is lost. Song starts from 0.

4. **Expiry**: `hydratedAt` timestamp is stored but **never checked** — there's no expiration logic.

---

## Cross-Source Queue Behavior

| Action | Old Queue | New Queue | Behavior |
|--------|-----------|-----------|----------|
| Play track from search | Any existing queue | Replaced with single track | `setQueueAndPlay([single], 0)` — old queue dropped |
| Play track from album | Any existing queue | Replaced with full album | `setQueueAndPlay(albumTracks, clickedIndex)` |
| Play from playlist | Any existing queue | Replaced with full playlist | Same pattern |
| Play similar track | Any existing queue | Replaced with similar tracks | Same pattern |
| Play single track (track page) | Any existing queue | Replaced with single track | `playTrack(single)` — `engine.play()` without queue |
| Click current track again | Unchanged | Unchanged | Pauses if playing, resumes if paused |
| Shuffle mode change | Unchanged | Same tracks, new order | Engine rebuilds shuffle order |

**Key insight**: Every new play action **completely replaces the queue**. There is no "add to queue" or "play next" from the UI. This means:
- If a user is listening to an album and clicks a search result, the album queue is lost
- If a user is halfway through a playlist and clicks a track from "You Might Like," the playlist is replaced

---

## Queue at End States

### End of Queue (repeat off)
```
audio 'ended' → engine.next()
  → queue.next() returns null (nothing after current)
  → engine stops playback
  → isPlaying = false
  → MediaSession state = 'none'
  → NowPlayingBar shows current track (still displayed) but paused
  → Progress bar at end (100%)
```

User sees: The track is displayed as if it just finished. No auto-advance, no "queue ended" message. The player remains visible with the last track's info.

### End of Queue (repeat all)
```
audio 'ended' → engine.next()
  → queue.next() returns null
  → repeatMode is 'all'
  → wrap to queue[0], play from start
```

### Single-Track Queue
When playing from search or track page:
```
Queue is [1 track]
Track ends → engine.next()
  → queue.next() returns null
  → repeat off → stop
  → repeat all → replay same track
```

---

## Volume & Audio State Persistence

| Setting | Persisted? | Key | Restored? |
|---------|-----------|-----|-----------|
| Volume level | ✅ | `player-volume` | ✅ On `initialize()` |
| Muted state | ✅ | `player-muted` | ✅ On `initialize()` |
| Shuffle mode | ❌ | — | Resets to `off` on new queue |
| Repeat mode | ❌ | — | Resets to `off` on new queue |
| Playback speed | ❌ | — | Resets to 1x on page reload |
| Audio quality | ✅ | (in player store) | ✅ |
| Sleep timer | ❌ | — | Lost on refresh |

---

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-701 | 🚨 BLOCKER | `queue-manager.ts` | **No "add to queue" or "play next" UI anywhere** — queue is only populated by `setQueueAndPlay()` | Users can't build a listening queue from individual tracks | Add "Add to Queue" and "Play Next" to track context menus and player controls |
| F-702 | ⚠️ MAJOR | Queue persistence | Queue persists only track IDs, not metadata — re-fetches all on restore (F-024) | Slow restoration; failed fetches drop tracks silently | Cache full metadata with TTL |
| F-703 | ⚠️ MAJOR | Queue persistence | Current queue position not persisted — restores to index 0 | User loses their place in a long queue on refresh | Persist current index |
| F-704 | ⚠️ MAJOR | `setQueueAndPlay()` | Every new play action **replaces the entire queue** — no "play next" or "add to end" | User's carefully curated queue is silently replaced | Confirm on replace or allow append mode |
| F-705 | 💡 IMPROVE | `QueuePanel.vue` | No visual indicator that removing current track advances immediately | Surprising behavior — user removes track, playback jumps | Show confirmation: "Removing current track will skip to next" |
| F-706 | 💡 IMPROVE | `QueuePanel.vue` | No "Clear Queue" button | User must remove tracks one by one | Add "Clear Queue" action |
| F-707 | 💡 IMPROVE | Queue at end | No "queue ended" message or visual — player just stops | User may think playback is broken | Show "End of queue — add more tracks" state |
| F-708 | 💡 IMPROVE | `stores/player.ts` | Shuffle and repeat modes are reset on every new `setQueueAndPlay()` | User's preferred listening mode is lost on each new play | Persist shuffle/repeat preference |
| F-709 | 💡 IMPROVE | `QueuePanel.vue` | Drag-reorder works but no visual feedback for drop target position | User doesn't know exactly where the track will land | Show drop indicator line |
| F-710 | 💡 IMPROVE | `setQueueAndPlay()` | Engine always resets shuffle to `off` when `setQueueAndPlay()` is called (line in store) | User who had shuffle on gets sequential playback | Preserve shuffle mode, reshuffle new queue |

## RTL / A11y / Mobile Notes

- ✅ Queue panel is slide-in from right (correct for RTL — left side would be better for RTL)
- ❌ **Drag reorder is not keyboard accessible** — `vuedraggable` uses mouse/touch only
- ✅ Each track in queue has `aria-label` with title and artist
- ✅ Remove button has `aria-label="Remove from queue"`
- ❌ No keyboard shortcut for "remove from queue"
