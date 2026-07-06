# Journey 03: First Home — The Empty State / Cold-Start Experience

## What Happens on First Visit (Brand-New User, No History)

### The Loading Phase

```
Time 0ms:     PageHome component mounts
              → fetchHomeFeed() called in onMounted
              → loading = true
              → 4 skeleton sections render:
                 ┌── Skeleton 1 ──┐
                 │ ████░░░░ (lines) │
                 │ ██ ██ ██ ██ ██  │  ← 5 card skeletons
                 └────────────────┘
                 (× 4 sections)
```

The skeletons show immediately. There are 4 skeleton sections, each with a line skeleton (title) and 5 card skeletons (album cover placeholders).

### The Data Fetch

`fetchHomeFeed()` makes these API calls:

```typescript
const results = await Promise.allSettled([
  // 1. Personalized home feed — ONLY if authenticated
  auth.isAuthenticated
    ? recsApi.getHomeFeed().catch(() => null)
    : Promise.resolve(null),

  // 2. Personalized recommendations — ONLY if authenticated
  auth.isAuthenticated
    ? recsApi.getPersonalized({ limit: 10 }).catch(() => null)
    : Promise.resolve(null),

  // 3. Albums catalog (always fetched)
  albumsApi.getAlbums().catch(() => [] as Album[]),

  // 4. Artists catalog (always fetched)
  artistsApi.getArtists().catch(() => [] as Artist[]),
])
```

**For a brand-new user (authenticated but zero history):**
- Home feed API → returns `{ sections: [] }` or sections with empty items
- Personalized API → returns `{ items: [] }`
- Albums API → returns whatever albums exist in the catalog (may be empty)
- Artists API → returns whatever artists exist (may be empty)

### The Empty State

After loading completes (~500ms-2s):

```typescript
const hasData = computed(() =>
  sections.value.length > 0 ||
  albums.value.length > 0 ||
  artists.value.length > 0
)
```

If `hasData` is `false` (catalog is empty too), the **empty state hero** is shown:

```
┌──────────────────────────────────────────────────┐
│  ┌──────────────────┐                            │
│  │ 🎧  🔔  ⭐  💗   │  ← Icon cluster          │
│  └──────────────────┘                            │
│                                                  │
│  READY TO LISTEN                                 │
│                                                  │
│  Your soundtrack                                 │
│  starts right here.                              │
│                                                  │
│  Discover new artists, save albums you love,     │
│  and build playlists that match every mood.      │
│                                                  │
│  ─────────────────────────────────               │
│                                                  │
│  [ 🔍 Explore Music ]  [ 🧭 Discover ]           │
│                                                  │
│  ┌──────────────────┬──────────────────┬────────┐ │
│  │ ⭐ Recommendations│ 🧭 Discover     │ ✨ Mood │ │
│  │                  │                  │ Explorer│ │
│  └──────────────────┴──────────────────┴────────┘ │
│                   3-card grid below                │
└──────────────────────────────────────────────────┘
```

The empty state is visually **beautiful** — gradient orbs, grid pattern background, glass cards, well-designed CTA buttons. It's a **deliberate empty state** not an afterthought.

### What If the Catalog HAS Content (Albums/Artists Exist)?

This is the more realistic scenario. In that case:

```
┌──────────────────────────────────────────────────┐
│  Hero carousel (top 5 popular tracks → may       │
│  be empty if no popular tracks exist)             │
│                                                  │
│  Personalized sections (ALL empty for new user):  │
│  ❌ "Based on Your Taste"  ← 0 items → hidden    │
│  ❌ "Recently played"      ← 0 items → hidden    │
│  ❌ "Trending"             ← 0 items → hidden    │
│  ❌ "For You"              ← 0 items → hidden    │
│  ❌ "More from your favorites" ← hidden          │
│  ❌ Because-of sections    ← hidden              │
│  ❌ "Popular in your genres" ← hidden             │
│                                                  │
│  ✅ Albums section (if albums exist)              │
│  ✅ Artists section (if artists exist)            │
│  ✅ 3-card CTA grid (always shown)               │
│     Recommendations / Discover / Mood Explorer   │
└──────────────────────────────────────────────────┘
```

The personalized sections all use `v-if="section.length"` so they're **hidden if empty**. The user sees:
- Albums (horizontal scroll carousel)
- Artists (circular avatar grid)
- The 3 utility cards at the bottom

### What About the Hero Carousel?

```typescript
const heroItems = computed<HeroItem[]>(() => {
  const tracks = popular.value.slice(0, 5)
  return tracks.map(...)
})
```

If `popular` is empty (no trending data), `heroItems` is `[]` and the hero is hidden entirely. The page jumps straight to albums/artists/grid.

### Recommendations Hub

The `/recommendations` hub has sub-pages: Popular, For You, Best, Recent. For a new user:
- `/recommendations/popular` → may show content (global popularity)
- `/recommendations/for-you` → likely empty (no personalization data)
- `/recommendations/best` → may show content
- `/recommendations/recent` → may show content

## State Matrix Findings

### PageHome
| State | Present? | Notes |
|-------|----------|-------|
| 🟢 Loading | ✅ | 4 skeleton sections (titles + card grids) |
| 🟢 Empty (no data at all) | ✅ | Beautiful hero empty state with CTAs |
| 🟢 Partial data (albums/artists but no recs) | ✅ | Hides empty sections, shows what exists |
| 🟢 Full data | ✅ | All sections with horizontal carousels |
| 🔴 Error state | Partial | `error` ref exists but not displayed in template — errors are silently caught |
| 🔴 Offline | ❌ | No offline-specific handling on home page |
| 🔴 Loading flash re-entry | ❌ | `KeepAlive max=3` — navigating away and back may re-fetch |

### Home Feed Data Flow
| State | Present? | Notes |
|-------|----------|-------|
| 🟢 Authenticated home feed | ✅ | `getHomeFeed()` called for auth users |
| 🔴 Guest home feed | Partial | Guests skip personalized APIs — albums/artists are fetched but no recommendations |
| 🟢 API failure resilience | ✅ | Each API wrapped in `.catch(() => null)` — graceful degradation |
| 🟢 Empty sections hidden | ✅ | `v-if="section.length"` pattern |

## Friction Points

| # | Severity | Location | Problem | User Impact |
|---|----------|----------|---------|-------------|
| F-301 | ⚠️ MAJOR | `PageHome.vue` | Error from `fetchHomeFeed` is stored in `error` ref but **never displayed in template** | Silent failure — user sees stale skeleton or empty state with no indication something went wrong |
| F-302 | ⚠️ MAJOR | `PageHome.vue:22-29` | Loading skeleton renders for ALL APIs, but transitions straight to empty state if APIs fail | Flicker from skeleton → empty even if data would have loaded after retry |
| F-303 | 💡 IMPROVE | `useHomeFeed.ts:52-69` | Personalized APIs are only called for authenticated users — guests see only albums/artists | Guest home page is significantly less engaging |
| F-304 | 💡 IMPROVE | `PageHome.vue:373-465` | Hero is derived from `popular` (global trending) — can be empty on first load for new platforms | No hero shown = less visual impact on first visit |
| F-305 | 💡 IMPROVE | `PageHome.vue` | After genre onboarding complete, user returns to home but no toast/feedback about it | Abrupt transition from onboarding → home with no "welcome" follow-up |
| F-306 | 💡 IMPROVE | `PageHome.vue` | The "Discover" CTA links to `/discover` which **redirects to `/search`** | User clicks "Discover" expecting new content, ends up on search page |

## RTL / A11y / Mobile Notes

- ✅ **Skip link present** in LayoutMusicApp
- ✅ **Image error handling**: `@error="onImgError"` on all images
- ✅ **Loading lazy**: `loading="lazy"` on all non-hero images
- ✅ **Aurora background**: `pointer-events-none` so it doesn't interfere
- ❌ **Error message not displayed**: No `role="alert"` for fetch errors
- ℹ️ **Hero carousel**: Horizontal scroll on mobile, hidden scrollbar (`scrollbar-hidden` class)
- ✅ **Touch targets**: CTA buttons and cards use adequate sizing

## Delight Opportunities

- ✨ **Welcome back after onboarding**: Show a brief "Welcome, [name]! Your music journey begins" toast/animation
- ✨ **Progressive reveal**: Show sections one by one with staggered entrance animation instead of all at once
- ✨ **Empty state with sample content**: If no data exists, offer to load sample/popular Persian tracks to seed the experience
- ✨ **"Mood check" prompt**: If the user has no history, suggest a quick "How are you feeling?" mood picker instead of empty sections
- ✨ **Genre onboarding reminder**: If user skipped onboarding, show a subtle banner on home: "Pick your favorite genres for personalized suggestions"

## Open Questions

1. Should the home page show a "retry" option when all API calls fail? Currently it silently shows empty state.
2. What happens when `albumsApi.getAlbums()` returns empty vs. error? Both result in empty arrays — no distinction.
3. Is there a seed-data script for new deployments? A brand-new deployment with no tracks would show only the empty state.
4. Should guests see curated recommendations (global popular) even without personalization?
