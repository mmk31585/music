# Journey 15: Social — Parties, Clubs, Stages, Presence

> Full trace: social route tree → party/club/stage creation → WebSocket sync → real-time interactions.

---

## Social Architecture Overview

The social system has **three distinct feature categories**, each with its own route, store, and WebSocket event space:

| Feature | Route | Description |
|---------|-------|-------------|
| **Listening Parties** | `/social/parties` | Real-time synchronized listening with friends |
| **Music Clubs** | `/social/clubs` | Persistent communities centered on music taste |
| **Live Stages** | `/social/stages` | Creator-hosted live listening sessions |

All three share a common WebSocket infrastructure (`useSocialSocket`) but have independent stores and event namespaces.

---

## Social Routes & Navigation

```
/social                          → SocialHub.vue (feed of all activity)
  /social/parties                → PartyList.vue (active & upcoming parties)
    /social/parties/new          → PartyCreate.vue
    /social/parties/:id          → PartyRoom.vue
      /social/parties/:id/chat   → PartyChat.vue (route override for mobile)
  /social/clubs                  → ClubList.vue
    /social/clubs/new            → ClubCreate.vue
    /social/clubs/:id            → ClubDetail.vue
  /social/stages                 → StageList.vue
    /social/stages/new           → StageCreate.vue
    /social/stages/:id           → StageRoom.vue
```

---

## Social Hub (`SocialHub.vue`)

A dashboard-style page showing:

| Section | Data Source | Empty State |
|---------|-------------|-------------|
| **Friend Activity** | WebSocket feed (last 50 events) | "No recent activity" with "Invite friends" CTA |
| **Active Parties** | `GET /api/v1/social/parties?active=true` | "No active parties — start one!" → create CTA |
| **Your Clubs** | `GET /api/v1/social/clubs?member=true` | "Join a club to get started" → explore CTA |
| **Live Now** | `GET /api/v1/social/stages?live=true` | "No stages live right now" → schedule CTA |
| **Friend Suggestions** | `GET /api/v1/social/suggestions` | Hidden if no suggestions |

### Initial Load

```typescript
const { data: activity } = await socialApi.getActivity({ limit: 50 })
const { data: parties } = await socialApi.getParties({ active: true })
const { data: clubs } = await socialApi.getClubs({ membership: 'member' })
const { data: stages } = await socialApi.getStages({ live: true })
```

All fetches are parallel via `Promise.all`. If any fails, that section shows its empty state silently.

---

## Listening Parties

### Party Lifecycle

```
1. Create → POST /api/v1/social/parties { name, description, is_public, starts_at }
     → Returns { party_id, invite_link }
2. Invite → Share invite_link (public parties) or direct user invite (private)
3. Join → POST /api/v1/social/parties/:id/join
4. Listen → WebSocket syncs playback for all members
5. Leave → POST /api/v1/social/parties/:id/leave (or closes tab → auto-detect after 5min)
6. End → POST /api/v1/social/parties/:id/end (host only)
```

### Party Room (`PartyRoom.vue`)

Layout:
- **Top**: Party name, host avatar, member count, timer (elapsed)
- **Left sidebar** (desktop) / **Bottom sheet** (mobile): Member list with presence indicators
- **Center**: Playback view — shared now-playing with host's queue
- **Right sidebar**: Chat panel

### Synchronization Model

The party has a **host-centric model**:

```
Host controls playback:
  play/pause/next/prev/seek → emit event → all guests mirror

Guest state:
  "I'm listening along" (no transport controls)
  or "Listening on my own" (can temporarily diverge from host)
  or "Paused" (muted in party)
```

### WebSocket Events (Party Namespace)

| Event | Direction | Payload |
|-------|-----------|---------|
| `party:join` | Client→Server | `{ party_id }` |
| `party:leave` | Client→Server | `{ party_id }` |
| `party:track_change` | Host→Server→All | `{ track_id, position, timestamp }` |
| `party:seek` | Host→Server→All | `{ position, timestamp }` |
| `party:play` | Host→Server→All | `{ position, timestamp }` |
| `party:pause` | Host→Server→All | `{ position, timestamp }` |
| `party:chat` | Client→Server→All | `{ user_id, username, message, timestamp }` |
| `party:member_joined` | Server→All | `{ user_id, username, avatar }` |
| `party:member_left` | Server→All | `{ user_id, username }` |
| `party:member_list` | Server→Client (on join) | `{ members: [] }` |
| `party:ended` | Server→All | `{ ended_by: user_id, reason }` |

### Key Finding: No Event Replay on Reconnect

```typescript
// socket.ts — reconnection logic
socket.on('connect', () => {
  // Re-joins party room
  // BUT: no request for missed events
  // Party state may be out of sync
})
```

**F-1101 (P0 BLOCKER)**: When a user's WebSocket reconnects after a disconnect, the party does **not** replay missed events — the user misses track changes, chat messages, and member activity that occurred while disconnected.

### Presence Indicators

Participants show:
- Online/green dot (socket connected + party active)
- Away/yellow dot (socket connected, tab not focused for 5+ min)
- Offline/gray dot (disconnected, party tab closed)

### Chat Panel

```
PartyChat.vue
  → Chat messages rendered in reverse scroll (newest at bottom)
  → Auto-scrolls to bottom on new message (if user hasn't scrolled up)
  → "New messages" badge if user has scrolled up
  → Empty state: "No messages yet — start the conversation"
  → Send: Enter key or send button
  → Character limit: 500
  → Rate limit: 5 messages per 10 seconds
```

---

## Music Clubs

### Club Lifecycle

```
1. Create → POST /api/v1/social/clubs { name, description, genres[], is_public, cover_image }
     → Returns { club_id }
2. Discover → Browse public clubs
3. Join → POST /api/v1/social/clubs/:id/join
4. Participate → Club feed, shared playlists, discussions
5. Leave → POST /api/v1/social/clubs/:id/leave
```

### Club Detail (`ClubDetail.vue`)

| Section | Content |
|---------|---------|
| **Header** | Cover image, name, description, member count, genre tags |
| **Action Bar** | Join/Leave button, Share, Notification bell |
| **Feed** | Club activity: member-joined, track-shared, playlist-created |
| **Collaborative Playlist** | Club's shared playlist (or multiple) |
| **Members** | Avatar grid with "View all" link |

### Club Discovery (`ClubList.vue`)

- Search bar to filter by name/genre
- Tabs: "Popular" / "New" / "Your Clubs"
- Empty state: "No clubs match your search" or "Join a club to see it here"

---

## Live Stages

### Stage Lifecycle

```
1. Schedule → POST /api/v1/social/stages { name, scheduled_at, genres[], is_18_plus }
     → Returns { stage_id }
2. Go Live → POST /api/v1/social/stages/:id/go-live (creator only)
     → GET /api/v1/social/stages/:id returns { is_live: true }
3. Host → Plays tracks, talks, engages
4. Listen → Audience listens synchronously
5. React → Emoji reactions, comments, tips (⊂ monetization)
6. End → POST /api/v1/social/stages/:id/end (creator only)
```

### Stage Room (`StageRoom.vue`)

Layout:
- **Top**: Stage name, host info, listener count, live badge
- **Center**: Host's playback — synchronized to all listeners
- **Sidebar**: Reactions (emoji waterfall), comments, tip button
- **Bottom**: No transport controls for listeners — host controls playback

### Reactions

```
Reaction emojis: ❤️ 🔥 🎵 👏 🎉 💯
  → Send: POST /api/v1/social/stages/:id/reactions { emoji }
  → Display: Waterfall animation falling across the screen
  → Rate: 1 reaction per 500ms
```

---

## WebSocket Infrastructure (`useSocialSocket`)

```typescript
// composables/useSocialSocket.ts
export function useSocialSocket() {
  const socket = ref<WebSocket | null>(null)
  const listeners = new Map<string, Set<Function>>()
  
  function connect(room: SocialRoom) { /* ... */ }
  function disconnect() { /* ... */ }
  function on<T>(event: string, cb: (data: T) => void) { /* ... */ }
  function off(event: string, cb: Function) { /* ... */ }
  function emit(event: string, data: unknown) { /* ... */ }
  function reconnect() { /* ... */ }
}
```

### Connection Lifecycle

1. User authenticates → socket connects (JWT token in query param)
2. User joins a room → socket sends `join` event
3. Socket disconnects → exponential backoff reconnect (1s, 2s, 4s, 8s, max 30s)
4. Socket reconnects → re-joins all previously joined rooms
5. User navigates away → socket stays connected (kept alive for background activity)
6. User logs out → socket disconnects

### Event Listener Pattern

Each feature (party, club, stage) registers its own listeners:

```typescript
// PartyRoom.vue
onMounted(() => {
  ws.on('party:member_joined', handleMemberJoined)
  ws.on('party:chat', handleChatMessage)
  ws.on('party:track_change', handleTrackChange)
  // ...
})

onUnmounted(() => {
  // Cleanup all listeners
})
```

---

## Social State Matrix

| State | Social Hub | Party Room | Club Detail | Stage Room |
|-------|-----------|------------|-------------|------------|
| 🟢 Loading | ✅ Skeleton grid | ✅ Full skeleton | ✅ Skeleton sections | ✅ Skeleton |
| 🟢 Empty | ✅ Per-section messages | ✅ "No members" for empty party | ✅ "No activity" | ✅ "No listeners" |
| 🟢 Has data | ✅ Feed + sections | ✅ Full room | ✅ Detail view | ✅ Live stage |
| 🟢 Error | ⚠️ Silent per-section fallback | ❌ Silent fail on connect | ❌ Silent | ❌ Silent |
| 🟢 Reconnecting | ❌ No banner/indicator | ❌ Out of sync after reconnect | ❌ No indicator | ❌ Out of sync |
| 🔴 Rate limited (chat) | ❌ No visible rate limit feedback | ❌ Chat failure silent | ❌ N/A | ❌ Reaction failure silent |

---

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-1501 | ⚠️ MAJOR | `useSocialSocket.ts` | No **connection status indicator** anywhere in UI | User has no way to know if they're connected to social features | Add connection dot in social nav |
| F-1502 | ⚠️ MAJOR | `PartyRoom.vue` | Party chat uses **local-only scroll state** — messages before joining not shown | Late joiners see a blank chat | Fetch last 50 messages on join |
| F-1503 | ⚠️ MAJOR | `useSocialSocket.ts` | No **heartbeat/ping** on socket — disconnection detected only on TCP close | Can take 30+ seconds to detect a dropped connection | Add 30s ping interval |
| F-1504 | 💡 IMPROVE | `PartyRoom.vue` | No **"listening independently" mode** — guests must follow host | Users who want to explore artist can't | Add "listen on my own" toggle that diverges temporarily |
| F-1505 | 💡 IMPROVE | `ClubDetail.vue` | Club **activity feed not paginated** — only shows latest 20 events | Can't see club history | Add "Load more" pagination |
| F-1506 | 💡 IMPROVE | `StageRoom.vue` | Reactions are **ephemeral — not persisted** | No history of who reacted | Persist top reactions as "most reacted" |
| F-1507 | 💡 IMPROVE | All social pages | **No RTL consideration** in chat layout — messages are always LTR | Persian messages in chat align left | Use `dir="auto"` on chat messages |
| F-1508 | 💡 IMPROVE | `PartyRoom.vue` | No **party timer / elapsed time** | Users don't know how long they've been listening | Add elapsed time counter |
| F-1509 | 💡 IMPROVE | `StageRoom.vue` | No **scheduled stages calendar view** | Users can't see upcoming stages in a calendar | Add calendar/list toggle |
| F-1510 | 💡 IMPROVE | `useSocialSocket.ts` | **No event replay buffer** on reconnect (same as F-1101) | Missed events during reconnection | Buffer last 50 events server-side, replay on reconnect |

## RTL / A11y / Mobile Notes

- ❌ Chat bubbles don't respect RTL layout — fixed left-aligned
- ✅ Live region announcements for party track changes (screen reader says "Now playing: [track]")
- ❌ Stage reactions (emoji waterfall) have no text alternative
- ✅ Touch targets for all social controls adequate
- ❌ Club and party creation forms have no `aria-describedby` for errors
