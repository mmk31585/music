# Journey 13: Discovery — Home Feed, Radio, Artist/Album Pages

> Full trace: home feed section ordering → radio flow → artist/album discoverability.

---

## Home Feed Section Architecture

PageHome.vue renders **11 possible sections** in a strict order, each conditionally shown:

| # | Section | Shows When | Data Source | Eyebrow |
|---|---------|-----------|-------------|---------|
| 1 | Hero Carousel | `popular.length > 0` (top 5) | `recsApi.getPopular()` | — |
| 2 | Based on Your Taste | `personalized.length > 0` | `recsApi.getPersonalized()` | "Personalized picks" |
| 3 | Recently Played | `recentPlays.length > 0` | `recsApi.getHomeFeed()` → `recently_played` section | "Jump back in" |
| 4 | Trending | `popular.length > 0` | `recsApi.getHomeFeed()` → `trending` section | "Popular right now" |
| 5 | For You | `forYou.length > 0` | `recsApi.getHomeFeed()` → `for_you` section | "Personalized picks" |
| 6 | From Your Artists | `fromYourArtists.length > 0` | `recsApi.getHomeFeed()` → `from_your_artists` | "Based on your listening" |
| 7 | Because Of [Genre] | `becauseOfSections.length > 0` | Dynamic per genre section | "Similar tracks" |
| 8 | Popular in Your Genres | `yourGenres.length > 0` | `recsApi.getHomeFeed()` → `your_genres` | "From genres you listen to" |
| 9 | Albums | `albums.length > 0` | `albumsApi.getAlbums()` | "Collections" |
| 10 | Featured Artists | `artists.length > 0` | `artistsApi.getArtists()` | "Meet the creators" |
| 11 | Quick Links (always) | Always shown | Static 3-card grid | Recommendations / Discover / Mood Explorer |

### Data Fetching

```typescript
const results = await Promise.allSettled([
  auth.isAuthenticated ? recsApi.getHomeFeed().catch(() => null) : Promise.resolve(null),
  auth.isAuthenticated ? recsApi.getPersonalized({ limit: 10 }).catch(() => null) : Promise.resolve(null),
  albumsApi.getAlbums().catch(() => []),
  artistsApi.getArtists().catch(() => []),
])
```

**Guest vs Authenticated**: Guests skip recommendation APIs entirely — they only see Albums, Featured Artists, and Quick Links. Very sparse experience.

### Section Computed Properties

```typescript
// sectionMap = key-by-section.id
recentPlays = sectionMap.recently_played?.items
popular = sectionMap.trending?.items
forYou = sectionMap.for_you?.items
fromYourArtists = sectionMap.from_your_artists?.items
becauseOfSections = sections.filter(s => s.id.startsWith('because_of_'))
yourGenres = sectionMap.your_genres?.items
```

### Empty State (New User, No Data)

A beautiful hero state with gradient orbs and CTAs:
- "Ready to Listen" / "Your soundtrack starts right here"
- Buttons: "Explore Music" → `/search`, "Discover" → `/search`
- 3 utility cards: Recommendations, Discover, Mood Explorer

### Refresh Behavior

- **On mount**: Full fetch every time (no caching beyond component lifecycle)
- **No pull-to-refresh**: Mobile users must navigate away and back
- **No periodic refresh**: Feed is static once loaded
- **No stale-while-revalidate**: Fresh data on every page mount (via `onMounted`)

---

## Discover Weekly

```typescript
// Route: /recommendations/discover-weekly
recsApi.getDiscoverWeekly() → GET /api/v1/recommendations/discover-weekly
```

The `DiscoverWeeklyResponse` is not used in the home page — it's only accessible via the `/recommendations/discover-weekly` route. If the user navigates there:

```
PageRecommendations sub-page
  → Shows personalized weekly playlist
  → Tracks can be played, saved, or added to playlist
  → If not authenticated → empty/redirect
```

---

## Radio Flow

### Starting Radio

Radio can be started from:
| Entry Point | Action | UX |
|-------------|--------|-----|
| Track context menu | Right-click → "Start Radio" | Opens RadioMode component |
| NowPlayingBar overflow | Menu → "Start Radio" | Same |
| Track detail page | "Start Radio" button | Same |
| Artist page | "Artist Radio" | Same |

### What Radio Does

```
Radio starts → POST /api/v1/radio/start { seed_track_id }
  → API returns sessionId
  → RadioMode component renders full-screen
  → First batch of tracks fetched: GET /radio/:sessionId/next?count=10
  → Tracks loaded into queue
  → User sees:
      • Full-screen immersive UI with album art + vinyl spin
      • Seed indicator: pulsing dot + "Radio · From [track/artist]"
      • Transport controls: prev, play/pause, next
      • Progress bar
      • Up Next / History sidebar
  → Auto-refill: when queue ≤ 3 tracks remaining, fetches next batch
  → On close: POST /radio/:sessionId/end
```

### Auto-Refill Logic

```typescript
watch(() => player.queue.value.length, (len) => {
  if (radioActive && len <= 3 && !isLoadingBatch) {
    fetchNextBatch(10)
  }
})
```

### Radio Modes (Shuffle)

The radio works in conjunction with the player's shuffle modes:
- When radio starts, shuffle mode is set to `'similar'` or `'catalog'`
- `next()` dispatches via shuffle mode → fetches from API when needed

---

## Artist Page

```
PageArtist.vue
  → Params: artistId from route
  → Sections:
    1. ArtistHero: avatar, name, monthly listeners, follow/unfollow CTA, play all + shuffle buttons
    2. Popular: top 5 tracks (with "Show all" toggle)
    3. Discography: album horizontal carousel
    4. You Might Also Like: related artists carousel
    5. Bio: expandable description (>300 chars truncated)
  → Background: dynamic color from artist image (useAlbumColors)
```

### Follow CTA

| State | Button Text | Action |
|-------|-------------|--------|
| Not following | "Follow" | `libraryApi.followArtist({ artist_id })` → optimistic update |
| Following | "Following" | `libraryApi.unfollowArtist(artistId)` → optimistic update |

### Empty States

| Missing Data | What Shows |
|-------------|------------|
| No tracks | Empty section, section hidden |
| No albums | "No albums yet" placeholder |
| No related artists | Section hidden |
| No bio | Bio section hidden |
| Image load error | Fallback icon |

---

## Album Page

```
PageAlbum.vue
  → Params: albumId from route
  → Sections:
    1. AlbumHero: cover art (with color glow), title, artist, year, track count, total duration
    2. Track list: numbered rows with play, title, artist, duration
    3. "Play All" + "Shuffle" buttons in hero
    4. "More from [artist]" section: artist album carousel
    5. "You Might Also Like": related albums section
  → Background: dynamic color from album art
```

---

## State Matrix Findings

| State | Home Feed | Radio | Artist Page | Album Page |
|-------|-----------|-------|-------------|------------|
| 🟢 Loading | ✅ 4 skeleton sections | ✅ Spinner on start | ✅ Full page skeleton | ✅ Full page skeleton |
| 🟢 Empty (no data) | ✅ Hero empty state | N/A (always has seed) | ✅ Hidden sections | ✅ "No tracks" |
| 🟢 Error loading | ⚠️ Silent catch | ❌ Silent | ❌ Silent | ❌ Silent |
| 🟢 Guest experience | 🟡 Albums + artists only | ❌ Requires auth | ✅ Viewable | ✅ Viewable |
| 🔴 Feed refresh | ❌ No pull-to-refresh | N/A | ❌ Manual refresh only | ❌ Manual refresh only |
| 🔴 Section ordering | ❌ No way to customize | N/A | N/A | N/A |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-1301 | ⚠️ MAJOR | `useHomeFeed.ts` | Guest home feed skips ALL recommendation APIs — shows only albums + artists + static cards | Guests have a dramatically worse home page with zero personalization | Add generic "popular across all users" sections for guests |
| F-1302 | ⚠️ MAJOR | `PageHome.vue` | No pull-to-refresh on mobile | Users can't refresh the feed to see new content | Add pull-to-refresh |
| F-1303 | ⚠️ MAJOR | `useHomeFeed.ts` | Feed errors silently caught — sections go missing with no explanation | User may not know their feed failed to load | Show per-section error state + retry |
| F-1304 | 💡 IMPROVE | `PageHome.vue` | Section ordering is hardcoded — no way to reorder or hide sections | Power users can't customize their home | Allow section reorder in settings |
| F-1305 | 💡 IMPROVE | Radio | Radio start requires a seed track — can't start "genre radio" directly | Users must first find a track of the genre they want | Add "Genre Radio" entry point |
| F-1306 | 💡 IMPROVE | `PageHome.vue` | Hero carousel uses `popular` (global trending) — may not match user's taste | Hero often irrelevant | Use personalized for hero when available |
| F-1307 | 💡 IMPROVE | `PageArtist.vue` | "Show all" toggle for top tracks loads same page inline — no "view full discography" | Artists with many tracks get truncated at 5 | Add link to full track list page |
| F-1308 | 💡 IMPROVE | All catalog pages | No breadcrumb navigation | Users exploring deep links can't see where they are | Add breadcrumb trail |
| F-1309 | 💡 IMPROVE | Radio | Auto-refill doesn't show countdown/fetch indicator | User surprised when queue suddenly refills | Show "Loading more tracks..." indicator |

## RTL / A11y / Mobile Notes

- ✅ Home page hero carousel wraps correctly in RTL
- ✅ Artist page uses `dir="auto"` for Persian artist names
- ❌ Album page track rows not labeled for screen readers (no `aria-label` on individual rows)
- ✅ Touch targets on all carousels are adequate
- ❌ Radio mode full-screen has no close button with `aria-label`
