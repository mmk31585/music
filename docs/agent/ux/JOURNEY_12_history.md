# Journey 12: History — Listening History, Replay, Privacy

> Full trace: how tracks get into history → where they're displayed → what a user can do with them.

---

## How Tracks Enter History

### Automatic Recording

Every time a track plays, the player engine records history:

```typescript
// player-engine.ts — doPlay()
providers.onPlayHistory?.(track.id, duration)
```

The provider (set up in `stores/player.ts`) calls:

```typescript
// stores/player.ts — initialize()
engine.setProviders({
  onPlayHistory: (trackId, duration) => {
    libraryApi.addPlayHistory({ track_id: trackId, duration, completed: false })
      .catch(() => {})  // Silent failure
    historyApi.addHistory({ track_id: trackId, duration, completed: false })
      .catch(() => {})  // Silent failure
  },
  // ...
})
```

**Two separate history endpoints are called simultaneously**:
1. `libraryApi.addPlayHistory()` → `POST /api/v1/library/history` (Library API)
2. `historyApi.addHistory()` → `POST /api/v1/history` (History API)

This means history is tracked in TWO places — likely a legacy migration artifact.

### Track Completion

The `completed: boolean` field is sent as `false` always — it's never updated to `true` when a track finishes. There is no "listened to completion" tracking in the frontend.

### What Triggers a History Entry

| Trigger | Recorded? | Notes |
|---------|-----------|-------|
| Track starts playing | ✅ Via `onPlayHistory` provider | Fires for every `doPlay()` call |
| Track finishes naturally | ❌ `completed` never set to `true` | |
| Seek into middle of track | ✅ Starts new history entry | Multiple entries for same track |
| Resuming paused track | ❌ Only first play recorded | |
| Skipping to next track | ✅ New entry for the skipped-to track | |
| Radio/catalog/similar mode | ✅ Every track played is recorded | |

### Silent Failure

Both history API calls have `.catch(() => {})` — failures are completely silent. No retry, no queue, no toast.

---

## Where History Appears

### Library → History Tab

```
Library page → History tab
  → Shows TrackRow list of recently played tracks
  → Sorted by most recent first
  → Source: libraryApi.getRecentlyPlayed()
  → Each row: cover, title, artist, timestamp, play button
```

**No pagination**: All history is fetched at once. For users with extensive history, this will be slow.

**No date grouping**: Tracks are just listed chronologically with no "Today", "Yesterday", "This Week" separators.

### History Tab Empty State

```
┌──────────────────────────────────┐
│        🕐 (history icon)         │
│                                  │
│        No history yet            │
│                                  │
│   Start playing tracks and       │
│   your history will appear       │
│   here.                          │
│                                  │
│   [  Discover music  ]           │
└──────────────────────────────────┘
```

CTA button links to `/discover` (which redirects to `/search`).

### Dedicated Recently Played Page

There is a `/recently-played` route that redirects to `/library` — no separate page exists.

### Search Page → Recently Played Section

On the search/discover page, a "Listen Again" section shows recently played items from `historyApi.getHistory()`.

---

## Actions Available from History

| Action | Available? | How |
|--------|-----------|-----|
| Play a history track | ✅ Click TrackRow → calls `player.setQueueAndPlay([track], 0)` |
| Play all from history | ❌ No "Play all" button |
| Add to playlist | ✅ Via TrackRow context menu → "Add to Playlist" |
| Like/Unlike | ✅ Via TrackRow heart icon |
| Remove from history | ❌ No remove button |
| Clear all history | ❌ No clear button |
| See when it was played | ✅ Timestamp shown (relative: "2h ago") |

---

## Privacy Considerations

| Privacy Feature | Status | Detail |
|----------------|--------|--------|
| Clear listening history | ❌ Not available in UI | No way to delete history from frontend |
| Pause history recording | ❌ Not available | Every play is recorded |
| Private session mode | ❌ Not available | No incognito mode within the app |
| "Now Playing" visibility | ✅ Partial | `music-status-privacy` setting in localStorage — but this only controls the "Now Playing" broadcast, not history recording |

---

## State Matrix Findings

| State | Library History Tab | Dedicated History Page |
|-------|-------------------|----------------------|
| 🟢 Loading | ✅ Shared skeleton | N/A (redirects to library) |
| 🟢 Items | ✅ TrackRow list | Same |
| 🟢 Empty | ✅ Distinct CTA | Same |
| 🟢 Error | ❌ Silent | Same |
| 🔴 Remove single item | ❌ Not possible | Same |
| 🔴 Clear all | ❌ Not possible | Same |
| 🔴 Date grouping | ❌ Flat list | Same |
| 🔴 Pagination | ❌ All at once | Same |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-1201 | 🚨 BLOCKER | Player engine + stores | **Two separate history APIs called simultaneously** — `libraryApi.addPlayHistory()` AND `historyApi.addHistory()` | Duplicate data, confusion about which is source of truth | Consolidate to one history API |
| F-1202 | ⚠️ MAJOR | `stores/player.ts` | History recording has `.catch(() => {})` — silent failure | User history silently not saved | Retry on failure + warning log |
| F-1203 | ⚠️ MAJOR | `PageLibrary.vue` history tab | **No way to remove individual tracks from history** | Can't clean up history | Add remove button per track |
| F-1204 | ⚠️ MAJOR | `PageLibrary.vue` history tab | **No "Clear all history" button** | History grows forever, no privacy control | Add clear button with confirmation |
| F-1205 | 💡 IMPROVE | `PageLibrary.vue` history tab | No date grouping in history list | Can't contextually find what was played "yesterday" vs "last week" | Add date section headers |
| F-1206 | 💡 IMPROVE | Player store | `completed: false` is always sent — never updated when track finishes | No way to distinguish "fully listened" from "skipped after 3 seconds" | Fire second API call on track end with `completed: true` |
| F-1207 | 💡 IMPROVE | History | No "Play all" from history | Can't replay recent listening session | Add "Play all" button |
| F-1208 | 💡 IMPROVE | `PageLibrary.vue` history tab | No pagination for large history | Users with 1000+ tracks see slow page load | Add pagination or virtual scroll |
| F-1209 | 💡 IMPROVE | Privacy | No private session or history pause toggle | Users may not want certain tracks recorded | Add "private session" toggle in settings |

## RTL / A11y / Mobile Notes

- ✅ `TrackRow` component used consistently with history tab
- ✅ CTA button in empty state navigates to discover
- ❌ No `aria-live` region when new tracks are added to history (but history is only updated on page load, not live)
- ✅ Timestamps use relative format ("2h ago")

## Delight Opportunities

- ✨ **"Replay your morning"**: Group history by session (e.g., "Your morning listening session: 5 tracks")
- ✨ **Most played stats**: Show a "Most Played This Week" mini-section in history
- ✨ **History export**: Allow exporting listening history as CSV/JSON for data-loving users
