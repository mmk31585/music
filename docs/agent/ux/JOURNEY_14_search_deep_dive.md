# Journey 14: Search Deep Dive — Full Search Journey & User Discovery

> Full trace: search from query entry → results → selection → play/discover.

---

## Search Entry Points

| # | Entry Point | Navigate To | UX Flow |
|---|-------------|-------------|---------|
| 1 | NavBar search icon | `PageSearch.vue` via route `/search` | Full page search |
| 2 | Keyboard shortcut `Ctrl+K` / `Cmd+K` | Opens CommandPalette overlay | Quick search modal |
| 3 | Keyboard shortcut `/` (any page) | Focuses NavBar search input | Inline search |
| 4 | Mobile tap on search icon | Focuses full-width search bar | Full page search (mobile) |
| 5 | "Explore Music" CTA on home hero | Route to `/search` | Full page |
| 6 | Context menu "Search for..." | Route to `/search?q=` | Prefilled search |

---

## Command Palette (`CommandPalette.vue`)

**Route**: Not a route — rendered as global overlay from `AppShell.vue`.

```typescript
// Keyboard: Cmd+K or Ctrl+K toggles visibility
// Also triggered from search icon in NavBar on desktop
defineProps<{
  visible: boolean
}>()
const emit = defineEmits<{ close: [] }>()
```

### Sections

| Section | Shows | Data Source |
|---------|-------|-------------|
| **Recent Searches** | Last 5, stored in localStorage | `searchHistory` ref (localStorage key: `muse-search-history`) |
| **Quick Actions** | "Search for..." + current query | Static |
| **Top Results** | Single best match (track, artist, or album) | API search with `limit=1, type=top_result` |
| **Tracks** | Up to 3 tracks | API search with `limit=3, type=track` |
| **Artists** | Up to 3 artists | API search with `limit=3, type=artist` |
| **Albums** | Up to 3 albums | API search with `limit=3, type=album` |

### Debounce
- 300ms debounce on search input
- Results appear inline, no page navigation
- On result click → closes palette → navigates to detail page

### Empty State (no query)
Recent searches list with "Clear all" button. If no recent searches, shows "Start typing to search" placeholder.

### Keyboard Navigation
- `↑` `↓` to navigate results
- `Enter` to select
- `Escape` to close
- `Tab` cycles through result groups

---

## Full Search Page (`PageSearch.vue`)

**Route**: `/search?q=:query` or just `/search`

### Search Input
- Auto-focuses on mount
- Prefilled from URL query param `?q=`
- 300ms debounce before API call
- Search icon + clear button (X) when query present
- Placeholder: Persian "جستجوی موسیقی، هنرمند، آلبوم..." / "Search music, artist, album..."

### Results Rendering

The page uses a **section-based layout**:

```typescript
// API: GET /api/v1/search?q=:query&limit=:limit&offset=:offset&type=:type[,type...]
// Returns { tracks, albums, artists, playlists, top_result }
```

| Section | Empty | Has Results | Error |
|---------|-------|-------------|-------|
| **Top Result** | Hidden | Hero card (largest, most prominent) | Hidden |
| **Tracks** | "No tracks found" | Up to 4 rows, "Show all" link | Silence |
| **Albums** | "No albums found" | Up to 4 cards, "Show all" link | Silence |
| **Artists** | "No artists found" | Up to 4 cards, "Show all" link | Silence |
| **Playlists** | "No playlists found" | Up to 4 cards, "Show all" link | Silence |
| **Lyrics** (genius) | Hidden if no match | Single lyric match shown | Hidden |
| **Genres** | Hidden | Tag chips for matching genres | Hidden |

### Top Result Prioritization

```typescript
// If track matches → show track
// Else if artist matches → show artist
// Else if album matches → show album
// Else → hide top result section
```

### "Show All" Behavior

Each section's "Show all" navigates to a dedicated sub-route with full results:
- `/search/tracks?q=:query` → `PageSearchTracks.vue`
- `/search/albums?q=:query` → `PageSearchAlbums.vue`
- `/search/artists?q=:query` → `PageSearchArtists.vue`
- `/search/playlists?q=:query` → `PageSearchPlaylists.vue`

These sub-pages use the same API but with `limit=20` (no client-side pagination beyond this — see F-1410).

### Search Filters (Desktop sidebar or mobile bottom sheet)

On desktop, a left sidebar with filter chips:
- **Type**: All, Tracks, Albums, Artists, Playlists, Lyrics
- **Genre**: List of available genres (from genre API)
- **Duration**: Short (< 3min), Medium, Long (> 7min)
- **Upload Date**: Today, This Week, This Month, This Year

Filters applied as query params: `&genres=pop,rock&duration=short&upload_date=week`

### Recent Searches

- Stored in `localStorage` under key `muse-search-history`
- Max 20 entries (FIFO eviction)
- Each entry: `{ query, type, timestamp }`
- Shows on search page when query is empty
- "Clear all" button at bottom
- Per-item dismiss (X) button
- Not synced across devices

### Search Analytics

```typescript
// POST /api/v1/search-analytics { query, results_count, clicked_result_id, clicked_result_type }
// Sent on result click
// Not sent on page load or mount — only on explicit result click
```

---

## Search API Contract

```typescript
interface SearchResponse {
  top_result?: {
    id: string
    type: 'track' | 'artist' | 'album'
    title: string
    subtitle: string
    image_url?: string
  }
  tracks: PaginatedResult<Track>
  albums: PaginatedResult<Album>
  artists: PaginatedResult<Artist>
  playlists: PaginatedResult<Playlist>
  lyrics?: LyricMatch[]
  genres?: string[]
}

interface PaginatedResult<T> {
  items: T[]
  total: number
  limit: number
  offset: number
}
```

---

## Search State Matrix

| State | Command Palette | Full Search Page |
|-------|----------------|------------------|
| 🟢 Empty (no query) | Recent searches or "Start typing" | Recent searches or "Start typing" |
| 🟢 Loading (debounce) | Previous results + spinner | Previous results + skeleton |
| 🟢 Has results | Sectioned results list | Sectioned results with top result |
| 🟢 No results | "No results for [query]" | "No results for [query]" |
| 🟢 Error | ❌ Silent fallback (last results) | ❌ Silent fallback |
| 🔴 Typo in Persian | ❌ No Persian typo tolerance | ❌ No Persian typo tolerance |
| 🔴 Network error | ❌ Search just stops working | ❌ Search just stops working |
| 🔴 Very long query | ❌ No truncation | ❌ No truncation |

---

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-1401 | ⚠️ MAJOR | `useSearch.ts` + `api/search` | Results have **no Persian typo tolerance** — مشق vs موسیق vs موسیقی produce wildly different results | Persian users with common typos find nothing | Add Persian phoneme matching or suggest corrections |
| F-1402 | ⚠️ MAJOR | `PageSearch.vue` | Reacts to `?q=` URL param changes, but **no debounce** on initial URL load — if URL changes rapidly (e.g., while typing) it fires every keystroke | Excessive API calls | Debounce URL param watcher |
| F-1403 | ⚠️ MAJOR | `PageSearch.vue` | Error state is **silent** — API failures just leave previous results showing | User may think search is working but results are stale | Show error banner with retry |
| F-1404 | 💡 IMPROVE | `CommandPalette.vue` | Recent searches are **localStorage only** — not synced across devices | User who searched on phone can't see on desktop | Sync recent searches to server |
| F-1405 | 💡 IMPROVE | `PageSearch.vue` | "Show all" navigates to new page with `limit=20` but **no client pagination** — just a "Show more" button that re-fetches | Large result sets (100+ items) require many clicks | Add infinite scroll or pagination |
| F-1406 | 💡 IMPROVE | `CommandPalette.vue` | Palette shows at most 3 items per result type — **can't expand inline** | User must go to full search page for more results | Add "Show all results" link per section |
| F-1407 | 💡 IMPROVE | `useSearch.ts` | Search history has **no deduplication** — same query can appear multiple times in history | Cluttered history | Deduplicate by query string (case-insensitive) |
| F-1408 | 💡 IMPROVE | `PageSearch.vue` | Filters are **not persisted** across page reloads | User must re-apply filters after navigation | Persist in URL query params |
| F-1409 | 💡 IMPROVE | `useSearch.ts` | No **search suggestions** as user types | User must complete their thought before seeing results | Add autocomplete suggestions |
| F-1410 | 💡 IMPROVE | `PageSearch.vue` | **No "search in my library" toggle** — search always searches all content | Users can't scope search to their liked tracks | Add library-scoped search toggle |
| F-1411 | 💡 IMPROVE | All search | No **voice search** button on mobile | Typing Persian with virtual keyboard is slow | Add voice search using Web Speech API |
| F-1412 | 💡 IMPROVE | `CommandPalette.vue` | Palette loads **all result types simultaneously** — no prioritization | Slow on large queries | Load top result first, then sections lazily |

## RTL / A11y / Mobile Notes

- ✅ Search input has proper `dir="auto"` — Arabic/Persian text flows RTL, English LTR
- ✅ Command palette handles arrow key navigation
- ❌ Command palette has **no `aria-label`** on close button
- ❌ Full search page filter chips not keyboard accessible on initial render (focus management issue)
- ✅ Touch targets for search results are adequate (>44px)
- ❌ No `aria-live` region for search results count — screen reader doesn't announce "Found X results"
