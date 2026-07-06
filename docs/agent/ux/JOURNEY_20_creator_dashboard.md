# Journey 20: Creator Dashboard — Analytics, Uploads, Earnings

> Full trace: creator gate → overview → earnings → audience → content management.

---

## Creator Dashboard Architecture

### Route

```typescript
// Route: /creator-dashboard → PageCreatorDashboard.vue (name: 'creator.dashboard')
// Meta: { title: 'Creator Dashboard', requiresAuth: true }
```

**Auth-guarded** — guests cannot access. **Not creator-guarded** at router level — gate is component-level (the page checks `isCreator` and shows an upgrade CTA if not).

### Is-Creator Check

```typescript
// api/creator/routes.ts
async function isCreator(): Promise<{ is_creator: boolean }> {
  return api.get('/creator/check')
  // Returns: { is_creator: true/false }
}

// PageCreatorDashboard.vue — onMounted
const { is_creator } = await creatorApi.isCreator()
if (!is_creator) {
  isNotCreator.value = true
  return // stops render, shows upgrade CTA
}
// Otherwise, fetch all dashboard data
```

### Not-Creator State

```
When is_creator === false:
  → Megaphone icon graphic
  → "Creator Studio" heading
  → "Upload tracks and build your audience"
  → "Start Creating" CTA → links to /admin/media (upload page)
  → No data, no tabs, no stats
```

---

## Dashboard Layout

```
PageCreatorDashboard.vue (689 lines)
  ├── Header: "Studio" label + "Creator Dashboard" title
  │     + Upload link (/admin/media) + Refresh button
  │
  ├── Metric cards (CreatorMetricCard.vue)
  │     ├── Total Plays
  │     ├── Listeners (unique)
  │     ├── Followers
  │     └── Revenue (currency formatted)
  │
  ├── Secondary metrics row
  │     ├── Tracks count
  │     ├── Albums count
  │     └── Playlists count
  │
  └── Tab panel
        ├── Overview tab
        │     ├── Daily Plays bar chart (14 days)
        │     └── Top Tracks list (up to 10)
        ├── Earnings tab
        │     ├── Revenue metrics: Total, Stream, Tips, Subscriptions
        │     ├── Pending payout, Last payout
        │     ├── Revenue breakdown bar chart
        │     └── Payout history list (paid/pending badges)
        ├── Audience tab
        │     ├── Total Listeners, New (7d), Repeat Rate
        │     ├── Top Listeners list (avatar, username, play count)
        │     └── Geographic stats (country flags)
        └── Content tab
              ├── Track list (with edit button per track)
              ├── Album grid (with cover art)
              └── Edit Track modal (Teleported)

Loading state: 4 skeleton cards + skeleton chart
Error state: ❌ Silent — API errors caught but no user-facing feedback
```

---

## Overview Tab

### Metric Cards (`CreatorMetricCard.vue`)

| Metric | Source | Format |
|--------|--------|--------|
| Total Plays | `creatorApi.getOverview()` → `total_plays` | `1,234,567` with K/M suffix |
| Listeners | `creatorApi.getOverview()` → `listeners` | `98,765` |
| Followers | `creatorApi.getOverview()` → `followers` | `5,432` |
| Revenue | `creatorApi.getOverview()` → `revenue` | `$12,345.67` |

Each card: glass surface, icon, formatted number, label, optional trend arrow (up/down/flat).

### Daily Plays Chart

```typescript
const dailyStats = ref<CreatorDailyStat[]>([])

// GET /creator/daily?from=2026-06-22&to=2026-07-06&limit=14
// Returns: [{ date, plays, listeners }, ...]
```

Simple bar chart:
- X-axis: dates (last 14 days)
- Y-axis: play count
- No interactive tooltip — just bars

### Top Tracks

```typescript
// GET /creator/tracks
// Returns: [{ track_id, title, plays, likes, comments }, ...]
// Top 10 by plays, rendered as list rows
// Each row: position, cover, title, play count, like count
// Click → play the track
```

---

## Earnings Tab

### Revenue Breakdown

| Metric | Source | Notes |
|--------|--------|-------|
| Total Revenue | `getEarnings().total` | Sum of all sources |
| Stream Revenue | `getEarnings().stream` | Per-play payouts |
| Tips | `getEarnings().tips` | Listener tips |
| Subscriptions | `getEarnings().subscriptions` | Fan subscriptions |
| Pending Payout | `getEarnings().pending_payout` | Amount awaiting payout |
| Last Payout | `getEarnings().last_payout` | Amount of last successful payout |

### Payout History

```typescript
// GET /creator/earnings/payouts?limit=20
// Returns: [{ id, amount, status, date, method }]
// Status: 'paid' (green badge) | 'pending' (yellow badge) | 'failed' (red badge)
// Method: 'bank_transfer', 'paypal', etc.
```

| State | Payout Row |
|-------|-----------|
| Paid | ✅ Green badge + checkmark + date |
| Pending | ⏳ Yellow badge + spinner icon |
| Failed | ❌ Red badge + error icon |
| Empty history | "No payouts yet" |
| Error loading | ❌ Silent — no error state |

### Revenue Chart

Same bar chart format as Daily Plays but showing revenue breakdown by source (stacked bars or grouped).

---

## Audience Tab

### Audience Metrics

| Metric | Source | Description |
|--------|--------|-------------|
| Total Listeners | `getAudience().total_listeners` | All-time unique listeners |
| New (7 days) | `getAudience().new_listeners_7d` | New listeners in last week |
| Repeat Rate | `getAudience().repeat_rate` | % of listeners who played >1 track |

### Top Listeners

```typescript
// GET /creator/audience?top_limit=20
// Returns: { top_listeners: [{ user_id, username, avatar_url, play_count }], geo: [...] }
// Rendered as list: avatar, username, play count
// Empty state: "No listeners yet"
```

### Geographic Stats

```typescript
// geo: [{ country_code, country_name, listener_count, flag_emoji }]
// Rendered as horizontal bar list
// Sorted by listener_count desc
// Empty state: "No geographic data yet"
```

---

## Content Tab

### Track List

```
Track rows with:
  → Cover art thumbnail
  → Title (primary) + Persian title (secondary)
  → Play count
  → Like count
  → Edit button → opens EditTrackModal
  → Delete button → confirmation dialog → DELETE /creator/tracks/:trackId

Empty state: "No tracks uploaded yet — start creating!"
```

### Album Grid

```
Album cards with:
  → Cover art
  → Title
  → Track count
  → Edit button → opens EditAlbumModal (or inline)
  → Delete button → confirmation → DELETE /creator/albums/:albumId

Empty state: "No albums yet"
```

### Edit Track Modal

```typescript
// Teleported modal (<Teleport to="body">)
// Props: trackId, initial data
// Fields:
//   - Title (required)
//   - Persian Title (optional)
//   - Lyrics (large textarea)
//   - Explicit toggle
//   - Genre selector
//   - Save/Cancel buttons

// PUT /creator/tracks/:trackId { title, persian_title, lyrics, explicit, genre }
```

---

## Dashboard Data Refresh

```typescript
async function refreshStats() {
  // Called manually via "Refresh" button
  // Also called on mount
  // No auto-refresh or polling
  
  // Re-fetches:
  await Promise.all([
    creatorApi.getOverview(),
    creatorApi.getDailyStats({ from, to, limit: 14 }),
    creatorApi.getTrackStats(),
    creatorApi.getEarnings(),
    creatorApi.getAudience({ top_limit: 20 }),
    creatorApi.getContent(),
  ])
}
```

**No WebSocket push** for real-time stats — manual refresh required. **No auto-refresh interval**. **No stale-while-revalidate**.

---

## Creator API Contracts (Backend)

```go
// internal/modules/creator/handler.go
type Handler struct {
  service Service
}

func (h *Handler) GetOverview(c *gin.Context)  → CreatorStats
func (h *Handler) GetDailyStats(c *gin.Context) → []CreatorDailyStat
func (h *Handler) GetTrackStats(c *gin.Context) → TrackStats
func (h *Handler) RefreshStats(c *gin.Context)  → message
func (h *Handler) IsCreator(c *gin.Context)     → { is_creator: bool }
func (h *Handler) GetEarnings(c *gin.Context)   → EarningsBreakdown
func (h *Handler) GetPayoutHistory(c *gin.Context) → Payout[]
func (h *Handler) GetPayoutMethods(c *gin.Context) → PayoutMethod[]
func (h *Handler) GetAudience(c *gin.Context)   → AudienceData
func (h *Handler) GetContent(c *gin.Context)    → CreatorContentData
func (h *Handler) UpdateTrack(c *gin.Context)   → Track
func (h *Handler) UpdateAlbum(c *gin.Context)   → Album
func (h *Handler) DeleteTrack(c *gin.Context)   → message
func (h *Handler) DeleteAlbum(c *gin.Context)   → message
```

All routes under `/creator` prefix. All require authentication. No admin role required — any user can become a creator by uploading.

---

## State Matrix

| State | Overview | Earnings | Audience | Content |
|-------|----------|----------|----------|---------|
| 🟢 Loading | ✅ 4 skeleton cards + chart skeleton | ✅ Skeleton cards | ✅ Skeleton lists | ✅ Skeleton list |
| 🟢 Loaded (has data) | ✅ Full dashboard | ✅ Metrics + chart + history | ✅ Listeners + geo | ✅ Track list + album grid |
| 🟢 Loaded (no data) | ✅ Zero-state metrics ("0 plays") | ✅ "No earnings yet" | ✅ "No listeners yet" | ✅ "No tracks uploaded yet" |
| 🟢 Not creator | ✅ Upgrade CTA (megaphone) | Hidden | Hidden | Hidden |
| 🟢 Error | ❌ Silent catch — shows empty/zero | ❌ Silent | ❌ Silent | ❌ Silent |
| 🔴 Refresh in progress | ✅ Spinner on refresh button | ✅ Same | ✅ Same | ✅ Same |
| 🔴 Track delete error | N/A | N/A | N/A | ❌ Silent — track stays in list |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-2001 | ⚠️ MAJOR | `PageCreatorDashboard.vue` | **All API errors silently caught** — `.catch(() => null)` with no user feedback | Creator thinks dashboard is working but data is stale/missing | Add per-section error state + "Try again" button |
| F-2002 | ⚠️ MAJOR | `PageCreatorDashboard.vue` | **No auto-refresh or WebSocket push** — stats only update on manual refresh | Creators see stale data for hours | Add 30s auto-refresh or WebSocket push |
| F-2003 | ⚠️ MAJOR | `PageCreatorDashboard.vue` | **No upload flow in dashboard** — "Upload" links to `/admin/media` | Disconnected experience, creator leaves dashboard | Add inline upload widget or upload from content tab |
| F-2004 | 💡 IMPROVE | `PageCreatorDashboard.vue` | **No date range picker** — daily stats always show last 14 days | Can't compare month-over-month | Add date range selector with presets (7d, 30d, 90d, custom) |
| F-2005 | 💡 IMPROVE | `PageCreatorDashboard.vue` | **No export/download** for stats (CSV/PDF) | Creators can't share reports | Add export button per tab |
| F-2006 | 💡 IMPROVE | `PageCreatorDashboard.vue` | **No track-level earnings** — only aggregate | Can't tell which tracks earn most | Add "Revenue by Track" breakdown |
| F-2007 | 💡 IMPROVE | `PageCreatorDashboard.vue` | **No notification settings** for creator milestones | Creators don't know when they hit 1000 plays | Add milestone alert preferences |
| F-2008 | 💡 IMPROVE | `PageCreatorDashboard.vue` | **No "Promote" or "Share" actions** for tracks | Creators can't easily share their work | Add share link/copy per track |
| F-2009 | 💡 IMPROVE | `PageCreatorDashboard.vue` | **No "Scheduled Releases"** — upload goes live immediately | No pre-release planning | Add scheduled publish date |
| F-2010 | 💡 IMPROVE | `PageCreatorDashboard.vue` | **Daily chart has no interactivity** — static bars, no tooltips | Hard to read exact numbers | Add tooltip on hover with exact count |
| F-2011 | 💡 IMPROVE | `PageCreatorDashboard.vue` | **No "Top Cities"** in geographic stats — only country level | Less granular audience insight | Add city-level data |
| F-2012 | 💡 IMPROVE | `PageCreatorDashboard.vue` | **No collaboration/invite** — other creators can't be added to tracks | No co-creator workflow | Add "Add collaborator" on track edit |

## RTL / A11y / Mobile Notes

- ✅ Metric cards use proper `aria-label` for numeric values
- ❌ Bar chart has **no text alternative** for screen readers — just visual bars
- ❌ Geographic stats use flag emojis as only identifier — missing country name in text
- ✅ Edit track modal uses `<Teleport>` for proper stacking
- ❌ Delete track has confirmation but **no undo toast** — permanent action
- ✅ Touch targets for all buttons adequate
- ❌ Daily chart bars are not keyboard navigable — no way to "hover" with keyboard
