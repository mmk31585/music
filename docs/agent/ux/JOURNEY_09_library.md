# Journey 09: Library — Save, Organize, Remove

> Full trace: like/unlike a track/album/artist → where saved → library page → organization → empty states → remove flows.

---

## Entry Points — Saving Content

### Like a Track

| Surface | Trigger | UX | API |
|---------|---------|----|-----|
| NowPlayingBar | Heart icon (♡) click | Filled heart animation + optimistic update | `libraryApi.likeTrack({ track_id, album_id?, artist_name? })` |
| FullscreenPlayer | Heart icon | Same | Same |
| ExpandedPlayer | Heart icon | Same | Same |
| MobileBottomSheet | Heart icon | Same | Same |
| TrackRow component | Heart icon on hover | Inline toggle | Same |
| Track detail page | Heart button in hero | Full button toggle | Same |
| Album page | Heart icon next to album title | Saves the album, not individual tracks | `libraryApi.likeAlbum(...)` |

**Optimistic update**: `useTrackLike()` composable immediately toggles the liked state in a shared in-memory `Set<string>` (unique track IDs). If the API call fails, it reverts.

**No undo toast**: Unlike modern patterns (YouTube's "Added to Library — Undo?" snackbar), Muse provides no undo. The action is immediate with no way to reverse except clicking unlike again.

### Follow/Unfollow an Artist

| Surface | Trigger | UX |
|---------|---------|-----|
| Artist page | Follow/Following button in hero | Text toggle |
| Artist card | Context menu or detail link | Navigates to artist page |
| Search results | Artist card click → artist page → follow | Two-step |

### Save an Album

| Surface | Trigger | UX |
|---------|---------|-----|
| Album page | Save/Unsave button or heart icon | Toggle |

## Library Page (`PageLibrary.vue`)

```
┌─────────────────────────────────────────────────────────────┐
│  ┌─────────────────────────────────────────────────────────┐│
│  │  Your Collection                  Stats bar            ││
│  │  LIBRARY                          ♡ 42  ■ 8  👥 12    ││
│  │                                                         ││
│  └─────────────────────────────────────────────────────────┘│
│                                                             │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐            │
│  │♡ Tracks││■ Albums││👥 Artists││📋 Playlists││🕐 History││
│  │  42   ││   8   ││   12  ││   5   ││   28  │            │
│  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘            │
│                                                             │
│  ── Content for selected tab ──                             │
│                                                             │
│  Tracks: searchable list with TrackRow component            │
│  Albums: grid of covers linking to album pages              │
│  Artists: circular avatars linking to artist pages           │
│  Playlists: grid + Create button                            │
│  History: TrackRow list                                     │
└─────────────────────────────────────────────────────────────┘
```

### Tracks Tab
- **Filter**: Text input filtering by title/artist/album (client-side, case-insensitive)
- **List**: `TrackRow` component (cover, title, artist, duration, actions)
- **Empty (no tracks)**: "No liked tracks yet. Tap the heart icon on any track to save it here."
- **Empty (with filter)**: "No tracks match your filter."

### Albums Tab
- **Grid**: 2-5 columns depending on viewport, cover art cards
- **Hover**: Play overlay on cover
- **Empty**: "No saved albums. Save albums to your library to find them quickly."

### Artists Tab
- **Grid**: Circular avatar cards
- **Empty**: "No followed artists. Follow artists to keep up with their latest releases."

### Playlists Tab
- **Grid**: Playlist cards with cover grid fallback
- **Create button**: Opens inline `Dialog` with Name, Description, Public checkbox
- **Empty**: "No playlists yet. Create your first playlist to start organizing your music." + CTA button

### History Tab
- **List**: `TrackRow` component (same as tracks tab)
- **Empty**: "No history yet. Start playing tracks and your history will appear here." + "Discover music" CTA

## Removing Items from Library

| Item | Remove Method | UX Feedback | Undo? |
|------|--------------|-------------|-------|
| Track | Unlike in any player surface | Heart unfills, track removed from library | No — must re-like |
| Album | Unsave on album page | Toggle off | No |
| Artist | Unfollow on artist page | Toggle off | No |
| Playlist track | Remove button on playlist detail | Track removed from list | No |
| Entire playlist | Delete on playlist detail | Confirmation dialog, redirect to library | No |
| History | No clear-all button | Must play new tracks to push old out | No |

## State Matrix Findings

| State | Tracks Tab | Albums Tab | Artists Tab | Playlists Tab | History Tab |
|-------|-----------|------------|-------------|--------------|-------------|
| 🟢 Loading | ✅ 3 shimmer sections | ✅ Same | ✅ Same | ✅ Same | ✅ Same |
| 🟢 Empty (no data) | ✅ Distinct message | ✅ Distinct message | ✅ Distinct message | ✅ Distinct message + CTA | ✅ Distinct message + CTA |
| 🟢 Empty (filter no match) | ✅ "No tracks match" | N/A | N/A | N/A | N/A |
| 🟢 Error loading | ❌ Silent `console.error` | ❌ Silent | ❌ Silent | ❌ Silent | ❌ Silent |
| 🔴 Remove from library | ❌ No inline remove | ❌ No inline remove | ❌ No inline remove | ❌ Must go to detail page | ❌ No clear history |
| 🔴 Item count caching | ❌ Fresh fetch every visit | Same | Same | Same | Same |
| 🔴 Sorting | ❌ Only filtering | ❌ No sort | ❌ No sort | ❌ No sort | ❌ Chronological only |
| 🔴 Pagination for large libraries | ❌ Fetches all at once | Same | Same | Same | Same |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-901 | ⚠️ MAJOR | `PageLibrary.vue` | All 5 sections fetched in one `Promise.all` — one slow API blocks all from rendering | Page shows loading skeleton until slowest API responds | Stream sections as they arrive |
| F-902 | ⚠️ MAJOR | `PageLibrary.vue` | No way to remove a liked track FROM the library page — must go to the track or use player | User looking at library can't clean it up | Add remove button to TrackRow in library context |
| F-903 | ⚠️ MAJOR | `PageLibrary.vue` | No way to unsave an album/unfollow artist from the library page | Same: must navigate to detail page | Add inline remove buttons |
| F-904 | ⚠️ MAJOR | `PageLibrary.vue` | Error silently swallowed — if an API fails, section is just empty | User thinks they have 0 items in a section | Show error state per tab |
| F-905 | 💡 IMPROVE | `useTrackLike.ts` | Like/unlike has no undo toast | Accidental unlike requires re-finding the track | Add "Undo" snackbar |
| F-906 | 💡 IMPROVE | `PageLibrary.vue` | No sort options (alphabetical, date added, artist) | Large libraries are hard to navigate | Add sort dropdown per tab |
| F-907 | 💡 IMPROVE | `PageLibrary.vue` | No pagination — all items fetched at once | Large libraries (>500 tracks) will be slow | Add pagination or virtual scroll |
| F-908 | 💡 IMPROVE | `PageLibrary.vue` | Tracks tab filter is client-side only | Only filters currently loaded items (pagination issue) | Move filter to API |
| F-909 | 💡 IMPROVE | Library API | No batch operations (unlike multiple tracks) | Cleaning up library is one-by-one | Add multi-select + batch unlike |
| F-910 | 💡 IMPROVE | `PageLibrary.vue` | Library not refreshed on return from track page (no stale-while-revalidate) | Adding a new liked track, going back, doesn't show update unless page is refreshed | Use `onActivated` for KeepAlive refresh |

## RTL / A11y / Mobile Notes

- ✅ Each tab button has `role="tab"` and `aria-selected`
- ✅ Each tab panel has `role="tabpanel"` and `aria-live="polite"`
- ✅ Hero stats use semantic labels
- ✅ Touch targets on grid items are adequate
- ❌ No keyboard shortcut to switch tabs (arrows or numbers)
- ❌ Tab counts are numeric badges — screen readers read them but order could be confusing
- ✅ Create playlist dialog is modal with proper focus trap
- ❌ Filter input has no `aria-label` beyond placeholder
