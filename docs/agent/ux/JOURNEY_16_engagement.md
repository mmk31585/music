# Journey 16: Engagement — Notifications, Gamification, Reactions

> Full trace: notification system → XPG/leveling → badges → reactions → leaderboard.

---

## Notification System

### Architecture

```
Notifications managed by: stores/notificationStore.ts (Pinia)
API base: GET    /api/v1/notifications (paginated)
          POST   /api/v1/notifications/:id/read
          POST   /api/v1/notifications/read-all
          DELETE /api/v1/notifications/:id
          GET    /api/v1/notifications/unread-count
WebSocket: notification:* events pushed from server
```

### Notification Types

| Type | Icon | Description | Actionable? |
|------|------|-------------|-------------|
| `new_follower` | 👤 | "{user} followed you" | Visit profile |
| `track_liked` | ❤️ | "{user} liked your track {track}" | View track |
| `comment` | 💬 | "{user} commented on {track}" | View track |
| `party_invite` | 🎉 | "{user} invited you to a party" | Join party |
| `club_activity` | 🏠 | "New activity in {club}" | View club |
| `badge_earned` | 🏆 | "You earned {badge} badge!" | View badges |
| `level_up` | ⬆️ | "You reached level {n}!" | View profile |
| `milestone` | 🎯 | "Your track {track} reached {plays} plays" | View track |
| `tip_received` | 💎 | "{user} tipped you {amount}" | View profile |
| `subscription_gift` | 🎁 | "{user} subscribed to you" | View profile |
| `system` | 🔔 | Platform announcements | Varies |

### Store Structure

```typescript
// stores/notificationStore.ts
const notifications = ref<Notification[]>([])
const unreadCount = ref(0)
const isLoading = ref(false)

function fetchNotifications(params: { page, limit }) { /* ... */ }
async function markAsRead(id: string) { /* optimistic update */ }
async function markAllAsRead() { /* optimistic */ }
async function deleteNotification(id: string) { /* optimistic */ }
function handleWebSocketEvent(event) { /* prepend to list */ }
function getUnreadCount() { /* ... */ }
```

### Notification Bell

```
NavBar notification bell:
  → Shows unread count badge (red dot with number, max 99+)
  → On click → opens NotificationPanel dropdown
  → On mobile → navigates to /notifications page
  → Unread count fetched on mount, updated via WebSocket push
```

### Notification Panel (`NotificationPanel.vue`)

| Section | Content |
|---------|---------|
| **Header** | "Notifications" title, "Mark all read" button |
| **Today** | Notifications from today |
| **Yesterday** | Notifications from yesterday |
| **This Week** | Notifications from the past 7 days |
| **Earlier** | All older notifications |
| **Empty** | "No notifications yet" with "Explore music" CTA |
| **Loading** | Skeleton list |

### WebSocket Integration

```typescript
socket.on('notification:new', (notification) => {
  // Prepend to notifications list
  notifications.value.unshift(notification)
  // Increment unread count
  unreadCount.value++
  // Play notification sound (optional, user-preference)
  // Show browser notification if permission granted
})
```

### Browser Notification Support

```typescript
// Check Notification API permission
async function requestNotificationPermission() {
  if ('Notification' in window && Notification.permission === 'default') {
    const result = await Notification.requestPermission()
    if (result === 'granted') {
      // Register for push notifications
    }
  }
}

// Show browser notification
function showBrowserNotification(notification: Notification) {
  if (Notification.permission === 'granted') {
    new Notification(notification.title, {
      body: notification.body,
      icon: notification.icon,
      tag: notification.id,
    })
  }
}
```

Browser notifications are **opt-in** — requested on first notification received if not yet decided.

---

## Gamification System

### XP & Leveling

```typescript
// stores/gamificationStore.ts
const xp = ref(0)
const level = ref(1)
const xpToNextLevel = ref(1000)
const isLoading = ref(false)

async function fetchGamification() {
  const { data } = await api.get('/api/v1/gamification')
  xp.value = data.xp
  level.value = data.level
  xpToNextLevel.value = data.xp_to_next_level
}

// GET /api/v1/gamification
// Returns: { xp: number, level: number, xp_to_next_level: number, xp_progress: number }
```

### XP Actions & Rewards

| Action | XP | Limit |
|--------|----|-------|
| Listen to a track | 10 XP | Per track, 1x per play |
| Like a track | 5 XP | 50 per day |
| Add to playlist | 10 XP | 30 per day |
| Share a track | 15 XP | 20 per day |
| Complete profile | 50 XP | Once |
| Follow a creator | 10 XP | 30 per day |
| Comment on track | 5 XP | 20 per day |
| Create a playlist | 25 XP | 10 per day |
| Daily login | 20 XP | Once per day |
| Complete all dailies | 100 XP | Once per day |

### Level Thresholds

```typescript
// Level thresholds are server-side
// Frontend receives: level + xp_to_next_level + progress percentage
// Visual: XP bar shown in user profile + bottom of NavBar dropdown

const LEVELS = [0, 100, 300, 600, 1000, 1500, 2100, 2800, 3600, 4500]
// Each level requires: LEVELS[level] XP to reach (cumulative)
```

### Level-Up Animation

```
When xp crosses threshold:
  → LevelUpModal or LevelUpToast appears
  → Animation: glowing particle effect, badge reveal
  → Duration: 3 seconds, auto-dismiss
  → Shows: new level, reward unlocked (if any)
```

### Badges (`PageBadges.vue` → `/profile/badges`)

| Badge Category | Examples | How to Earn |
|----------------|----------|-------------|
| **Listening** | "Early Bird" (100 tracks), "Audiophile" (1000 tracks), "Night Owl" (50 night plays) | Track counts |
| **Social** | "Social Butterfly" (10 friends), "Party Starter" (host 5 parties) | Social actions |
| **Creator** | "First Upload", "Rising Star" (1000 plays), "Hit Maker" (10000 plays) | Creator milestones |
| **Streak** | "Week Streak", "Month Streak", "Year Streak" | Daily login streaks |
| **Special** | "Beta Tester", "Founding Member", seasonal | Time-limited events |

### Badge Display

```
BadgeGrid.vue:
  → Grid of badge cards
  → Each card: icon, name, description
  → If earned: full color + earned date
  → If not earned: grayscale + "How to earn" tooltip
  → Featured badges show first (user's top 3 badges, set in profile)
```

### Leaderboard

```typescript
// PageLeaderboard.vue → /leaderboard
// GET /api/v1/leaderboard?period=weekly&limit=20

interface LeaderboardEntry {
  rank: number
  user_id: string
  username: string
  avatar: string
  xp: number
  level: number
  badges: Badge[]
  is_following: boolean
}
```

| Tab | Period | Shows |
|-----|--------|-------|
| Weekly | This week (Mon-Sun) | Top 100 |
| Monthly | This month | Top 100 |
| All Time | All time | Top 100 |
| Friends | Following only | Top 50 |

### Leaderboard State

| State | Behavior |
|-------|----------|
| 🟢 Loading | Skeleton rows |
| 🟢 Empty (friends) | "Follow friends to see them here" with "Find friends" CTA |
| 🟢 Error | ❌ Silent — shows empty table |
| 🔴 User not in top 100 | ❌ User's rank not shown at all — no "your rank" indicator |

---

## Reactions (Track Level)

### Track Reaction System

```typescript
// POST /api/v1/tracks/:id/reaction { type: 'like' | 'love' | 'fire' | 'sad' | 'chill' }
// DELETE /api/v1/tracks/:id/reaction (remove)
// GET /api/v1/tracks/:id/reactions → { counts: { like: 42, love: 15, ... }, user_reaction }
```

Reaction types:
| Type | Icon | Emotion |
|------|------|---------|
| `like` | 👍 | General appreciation |
| `love` | ❤️ | Deep affection for the track |
| `fire` | 🔥 | Energetic / banger |
| `sad` | 😢 | Melancholic / emotional |
| `chill` | 😎 | Relaxed / vibing |

### UI Integration Points

| Location | Display | Interaction |
|----------|---------|-------------|
| NowPlayingBar | Single reaction button (user's current reaction or "React") | Click → reaction picker dropdown |
| FullscreenPlayer | Full reaction bar with all 5 types + counts | Click type to toggle |
| TrackRow (catalog) | Inline reaction count | Click count → tooltip with breakdown |
| TrackDetail page | Full reaction section with user list per type | Click to view reactors |

### Reaction Picker

```
ReactionPicker.vue:
  → 5 emoji buttons in a row
  → Selected state: filled background, slightly larger
  → Count shown below each emoji
  → Hover: tooltip with first 3 reactors
  → Click toggles reaction on/off
```

---

## Grok / AI Integration

```
/ai/playlist → PageMoodPlaylist.vue (see muse-ai skill)
  → Mood-based playlist generation
  → Separate from core engagement flow
```

---

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-1601 | ⚠️ MAJOR | `notificationStore.ts` | Notifications **no pagination on initial load** — fetches only latest 20 | Users can't browse their notification history | Add "Load more" with cursor pagination |
| F-1602 | ⚠️ MAJOR | `notificationStore.ts` | **No notification grouping** — 10 followers triggers 10 separate notifications | Notification spam when popular | Group "X and Y others followed you" |
| F-1603 | ⚠️ MAJOR | `PageLeaderboard.vue` | **User's rank not shown** if not in top 100 | Users below top 100 have no reference point | Add "Your rank: #1,234" below top 100 |
| F-1604 | ⚠️ MAJOR | Leaderboard API | **Friends tab shows no data** if user follows no one | Empty state doesn't explain how to fix it | Add "Find friends via contacts" action |
| F-1605 | 💡 IMPROVE | `gamificationStore.ts` | XP is **client-side stale** — only fetched once on profile page mount | Level-up animation plays on page reload, not when XP is earned | Push XP updates via WebSocket in real-time |
| F-1606 | 💡 IMPROVE | `PageBadges.vue` | Badge progress not shown for locked badges — just grayscale icon | User doesn't know how close they are to earning it | Show progress bar "423/1000 tracks" |
| F-1607 | 💡 IMPROVE | `NotificationPanel.vue` | No **notification sound preferences** — either on or off | Users who want sound for tips but not social can't configure | Per-category sound settings |
| F-1608 | 💡 IMPROVE | `ReactionPicker.vue` | No **animation** on reaction — just instant state change | Feels flat compared to other platforms | Add micro-animation on reaction toggle |
| F-1609 | 💡 IMPROVE | `gamificationStore.ts` | **No daily challenge UI** — XP for dailies exists but no "Daily Tasks" panel | Users don't know what to do for XP | Add "Daily Tasks" checklist with progress |
| F-1610 | 💡 IMPROVE | `PageLeaderboard.vue` | Leaderboard **doesn't highlight user's rank** within the list | Hard to find yourself in the list | Highlight user's row with accent color |
| F-1611 | 💡 IMPROVE | All gamification | **No notifications for badge/level earned** — user has to visit profile to see | Missed dopamine hits for engagement | Push notification on level-up and badge |

## RTL / A11y / Mobile Notes

- ✅ Notification panel uses `dir="auto"` for mixed-language content
- ❌ Level-up animation has **no `prefers-reduced-motion`** respect — users with motion sensitivity get stuck with particles
- ✅ Reaction emoji picker has proper `aria-label` on each button
- ❌ Leaderboard rows don't announce rank position to screen readers
- ✅ Badge tooltips accessible via keyboard (focus + hover)
