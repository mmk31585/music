# Social Hub (/social)

> Complete redesign of `PageSocial.vue` + `ListeningPartyCard.vue`, `LiveRoomCard.vue`, `MusicClubCard.vue`
> Fixes: emoji icons, flat cards, hidden discussions, no live activity, stale hero

---

## Design Philosophy

The social hub is a **living, breathing space** — not a list of rooms. It should feel like walking into a vibrant Persian music café where:
- The "hot track right now" is visibly pulsing
- Friends' listening activity scrolls by in real time
- Every room/party/club card shows **who** is there (avatars)
- You can jump from "what I'm listening to" → "start a party with this song" in 1 tap

---

## Layout (Desktop — 3-column bento grid)

```
┌──────────────────────────────────────────────────────────────────────┐
│  SOCIAL HUB                                                           │
│  "با دوستات گوش بده" (Listen with friends)                            │
│                                                                       │
│  [What are you listening to? → Start a party] — inline mini-player   │
│                                                                       │
│  ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────────┐  │
│  │  🔴 LIVE NOW      │ │  TRENDING        │ │  YOUR CLUBS          │ │
│  │                   │ │  PARTIES          │ │                       │ │
│  │  Featured Room    │ │                   │ │  Club 1    ◉ 12/50  │ │
│  │  ┌─────────────┐  │ │  Party A   ◉ 12  │ │  Club 2    ◉ 8/30   │ │
│  │  │ Waveform    │  │ │  Party B   ◉ 8   │ │  Club 3    ◉ 24/100 │ │
│  │  │ animation   │  │ │  Party C   ◉ 5   │ │                       │ │
│  │  └─────────────┘  │ │                   │ │  [Browse All Clubs]  │ │
│  │  "Jazz Night"     │ │ [See All →]      │ │                       │ │
│  │  ۲۸ listening     │ │                   │ │                       │ │
│  │  [Listen Live]    │ │                   │ │                       │ │
│  └──────────────────┘ └──────────────────┘ └──────────────────────┘  │
│                                                                       │
│  ┌─── ACTIVITY RIVER ─────────────────────────────────────────────┐  │
│  │                                                                  │  │
│  │  ◉ Ali started listening to "Gole Yakh"                     2m  │  │
│  │  ◉ Sara joined "Friday Jazz Party"                          5m  │  │
│  │  ◉ "Electronic Night" room went live                       12m  │  │
│  │  ◉ Reza created club "موسیقی کلاسیک ایران"                   1h  │  │
│  │  ◉  ...                                                        │  │
│  │                                                                  │  │
│  │  [Show more]                                                     │  │
│  └──────────────────────────────────────────────────────────────────┘  │
│                                                                       │
│  ┌─── EXPLORE ────────────────────────────────────────────────────┐  │
│  │  All Rooms  |  All Parties  |  All Clubs  |  Discussions        │  │
│  │                                                                  │  │
│  │  ┌────┐ ┌────┐ ┌────┐ ┌────┐ ┌────┐                             │  │
│  │  │Card│ │Card│ │Card│ │Card│ │Card│  ← horizontal scroll        │  │
│  │  └────┘ └────┘ └────┘ └────┘ └────┘                             │  │
│  │                                                                  │  │
│  └──────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────┘
```

---

## Component Specs

### Hero Section — Personalized, Dynamic
```vue
<div class="relative overflow-hidden rounded-[2rem] p-8 md:p-10">
  <!-- Aurora background based on current track or top genre -->
  <div class="pointer-events-none absolute inset-0" aria-hidden="true">
    <div class="absolute -top-40 -right-40 h-[500px] w-[500px] rounded-full bg-[#1db954]/8 blur-3xl" />
    <div class="absolute -bottom-20 -left-20 h-[300px] w-[300px] rounded-full bg-[#a855f7]/6 blur-3xl" />
    <div class="absolute inset-0 bg-gradient-to-br from-black/40 via-transparent to-black/60" />
  </div>

  <div class="relative z-10 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
    <div>
      <p class="text-[10px] font-bold tracking-[0.35em] text-white/30 uppercase">Social Hub</p>
      <h1 class="mt-1 text-3xl font-black text-white md:text-4xl">
        {{ userStore.isLoggedIn ? `👋 خوش اومدی، ${userStore.displayName}` : 'Community' }}
      </h1>
      <p class="mt-1 text-sm text-white/50">با دوستات گوش بده — Listen with friends</p>
    </div>

    <!-- Inline "Start a Party from current track" -->
    <div v-if="currentTrack"
         class="flex items-center gap-3 rounded-2xl bg-white/[0.06] p-3 ring-1 ring-white/[0.10]">
      <div class="h-10 w-10 overflow-hidden rounded-xl">
        <img :src="currentTrack.coverUrl" alt="" class="h-full w-full object-cover" />
      </div>
      <div class="min-w-0 max-w-[180px]">
        <p class="truncate text-xs font-medium text-white">{{ currentTrack.title }}</p>
        <p class="truncate text-[10px] text-white/40">Now Playing</p>
      </div>
      <button aria-label="Start a listening party with this track"
              class="rounded-full bg-[#1db954] px-4 py-2 text-[11px] font-bold text-black
                     transition hover:bg-[#1ed760] hover:scale-105 active:scale-95
                     focus-visible:outline-2 focus-visible:outline-white">
        <span class="flex items-center gap-1.5">
          <Icon name="users" class="h-3 w-3" aria-hidden="true" />
          Party
        </span>
      </button>
    </div>
  </div>
</div>
```

### Live Now Card (Hero Feature — replaces flat tab grid)
```vue
<div class="rounded-2xl bg-white/[0.06] p-5 ring-1 ring-white/[0.10] transition
            hover:bg-white/[0.08] group">
  <div class="flex items-center gap-2 mb-3">
    <span class="flex h-2 w-2 rounded-full bg-red-500 motion-safe:animate-pulse" />
    <span class="text-[10px] font-bold tracking-wider text-red-400 uppercase">Live Now</span>
  </div>

  <!-- Live waveform visualization stub -->
  <div class="flex items-end gap-0.5 h-12 mb-3" aria-hidden="true">
    <div v-for="i in 20" :key="i"
         class="w-1.5 rounded-full bg-[#1db954]/60 transition-all"
         :style="{ height: `${15 + Math.random() * 35}px`, animationDelay: `${i * 80}ms` }"
         :class="motionSafe ? 'motion-safe:animate-waveform' : ''" />
  </div>

  <h3 class="text-lg font-bold text-white">{{ featuredRoom.title }}</h3>
  <p class="mt-1 text-xs text-white/40">{{ featuredRoom.description }}</p>

  <!-- Listener avatars -->
  <div class="mt-3 flex items-center gap-2">
    <div class="flex -space-x-2" aria-label="{{ featuredRoom.listenerCount }} listeners">
      <div v-for="avatar in featuredRoom.avatars.slice(0, 5)" :key="avatar"
           class="h-7 w-7 overflow-hidden rounded-full border-2 border-[#0A0A0A]">
        <img :src="avatar" alt="" class="h-full w-full object-cover" />
      </div>
      <div v-if="featuredRoom.listenerCount > 5"
           class="flex h-7 w-7 items-center justify-center rounded-full border-2
                  border-[#0A0A0A] bg-white/10 text-[9px] font-bold text-white/60">
        +{{ featuredRoom.listenerCount - 5 }}
      </div>
    </div>
    <span class="text-[11px] text-white/40 tabular-nums">{{ featuredRoom.listenerCount }} listening</span>
  </div>

  <button aria-label="Join live room"
          class="mt-4 w-full rounded-xl bg-red-500/10 py-2.5 text-xs font-bold text-red-400
                 transition hover:bg-red-500/20 active:scale-[0.98]
                 focus-visible:outline-2 focus-visible:outline-red-400">
    Listen Live
  </button>
</div>
```

### Activity River — Real-time social proof
```vue
<div class="rounded-2xl bg-white/[0.06] p-5 ring-1 ring-white/[0.10]" role="feed"
     aria-label="Live activity feed" aria-live="polite">
  <div class="flex items-center justify-between mb-4">
    <h3 class="text-sm font-bold text-white">Activity</h3>
    <span class="flex items-center gap-1.5 text-[10px] text-white/30">
      <span class="h-1.5 w-1.5 rounded-full bg-[#1db954] motion-safe:animate-pulse" aria-hidden="true" />
      Live
    </span>
  </div>

  <TransitionGroup name="activity-slide" class="space-y-2">
    <ActivityItem v-for="item in visibleActivities" :key="item.id" :item="item" />
  </TransitionGroup>

  <button v-if="hasMoreActivities"
          class="mt-3 w-full rounded-lg py-2 text-xs text-white/40 transition
                 hover:bg-white/[0.04] hover:text-white/60"
          @click="showAllActivities">
    Show {{ remainingCount }} more
  </button>
</div>
```

```vue
<!-- ActivityItem — used in both profile feed and social hub -->
<div class="flex items-start gap-3 rounded-xl p-3 transition hover:bg-white/[0.04]">
  <RouterLink :to="`/profile/${item.userId}`"
              class="shrink-0 h-8 w-8 overflow-hidden rounded-full ring-1 ring-white/10
                     focus-visible:outline-2 focus-visible:outline-[#1db954]">
    <img v-if="item.avatarUrl" :src="item.avatarUrl" :alt="item.userName"
         class="h-full w-full object-cover" loading="lazy" />
    <div v-else class="flex h-full w-full items-center justify-center bg-white/10 text-[10px] font-bold text-white">
      {{ (item.userName || '?')[0] }}
    </div>
  </RouterLink>

  <div class="min-w-0 flex-1">
    <p class="text-xs text-white/70 leading-relaxed">
      <RouterLink :to="`/profile/${item.userId}`"
                  class="font-semibold text-white hover:underline">
        {{ item.userName }}
      </RouterLink>
      {{ item.action }} <!-- e.g. "started listening to", "joined", "created" -->
      <template v-if="item.targetName">
        <RouterLink v-if="item.targetUrl" :to="item.targetUrl"
                    class="font-medium text-[#1db954] hover:underline">
          {{ item.targetName }}
        </RouterLink>
        <span v-else class="font-medium text-white/80">{{ item.targetName }}</span>
      </template>
    </p>
    <p class="mt-0.5 text-[10px] text-white/30">{{ item.timeAgo }}</p>
  </div>

  <!-- Contextual action button -->
  <button v-if="item.actionType === 'party'"
          aria-label="Join party"
          class="shrink-0 rounded-lg bg-[#1db954]/10 px-3 py-1.5 text-[10px] font-semibold
                 text-[#1db954] transition hover:bg-[#1db954]/20
                 focus-visible:outline-2 focus-visible:outline-[#1db954]">
    Join
  </button>
  <button v-else-if="item.actionType === 'room'"
          aria-label="Listen live"
          class="shrink-0 rounded-lg bg-red-500/10 px-3 py-1.5 text-[10px] font-semibold
                 text-red-400 transition hover:bg-red-500/20">
    Listen
  </button>
</div>
```

### Redesigned Cards — No more emoji icons, avatar grids, proper glass

**PartyCard (replaces ListeningPartyCard.vue)**
```vue
<div class="group rounded-2xl bg-white/[0.06] p-5 ring-1 ring-white/[0.10]
            transition-all duration-300 hover:bg-white/[0.08] hover:ring-white/[0.15]
            focus-within:ring-[#1db954] focus-within:ring-2">
  <div class="flex items-start justify-between gap-3">
    <div class="min-w-0 flex-1">
      <div class="flex items-center gap-2">
        <h3 class="truncate text-sm font-bold text-white">{{ party.title }}</h3>
        <span v-if="party.status === 'active'"
              class="inline-flex items-center gap-1 rounded-full bg-green-500/10 px-2 py-0.5
                     text-[9px] font-semibold text-green-400 uppercase">
          <span class="h-1.5 w-1.5 rounded-full bg-green-500 motion-safe:animate-pulse" aria-hidden="true" />
          Live
        </span>
      </div>
      <p v-if="party.description" class="mt-1 line-clamp-2 text-xs text-white/40 leading-relaxed">
        {{ party.description }}
      </p>
    </div>

    <!-- Track art thumbnail if there's a current track -->
    <div v-if="party.currentTrackCover"
         class="h-14 w-14 shrink-0 overflow-hidden rounded-xl ring-1 ring-white/10">
      <img :src="party.currentTrackCover" alt="" class="h-full w-full object-cover" />
    </div>
  </div>

  <!-- Participant avatars row -->
  <div class="mt-4 flex items-center gap-2">
    <div class="flex -space-x-1.5" :aria-label="`${party.participantCount} participants`">
      <div v-for="(avatar, i) in party.avatars.slice(0, 4)" :key="i"
           class="h-6 w-6 overflow-hidden rounded-full border-2 border-[#0A0A0A] transition
                  group-hover:border-white/20">
        <img :src="avatar" alt="" class="h-full w-full object-cover" />
      </div>
      <div v-if="party.participantCount > 4"
           class="flex h-6 w-6 items-center justify-center rounded-full border-2
                  border-[#0A0A0A] bg-white/10 text-[8px] font-bold text-white/50">
        +{{ party.participantCount - 4 }}
      </div>
    </div>
    <span class="text-[10px] text-white/30 tabular-nums">{{ party.participantCount }} listening</span>
  </div>

  <!-- Now playing mini info -->
  <div v-if="party.currentTrackName" class="mt-3 flex items-center gap-1.5 text-[10px] text-white/30">
    <Icon name="music" class="h-3 w-3 shrink-0" aria-hidden="true" />
    <span class="truncate">{{ party.currentTrackName }}</span>
  </div>

  <div class="mt-3 flex items-center gap-2">
    <button aria-label="Join party"
            class="flex-1 rounded-lg bg-[#1db954]/10 py-2.5 text-xs font-bold text-[#1db954]
                   transition hover:bg-[#1db954]/20 active:scale-[0.98]
                   focus-visible:outline-2 focus-visible:outline-[#1db954]">
      <span class="flex items-center justify-center gap-1.5">
        <Icon name="headphones" class="h-3 w-3" aria-hidden="true" />
        Join Party
      </span>
    </button>
    <button aria-label="Share party"
            class="flex h-9 w-9 items-center justify-center rounded-lg text-white/30
                   transition hover:bg-white/[0.06] hover:text-white/60
                   focus-visible:outline-2 focus-visible:outline-[#1db954]">
      <Icon name="share" class="h-4 w-4" aria-hidden="true" />
    </button>
  </div>
</div>
```

### Discussion Tab — No more type+ID filters

Instead of requiring the user to type a `target_type` and `target_id`, show **recent discussions** across the platform and let users filter by category tabs:

```vue
<div class="space-y-4">
  <!-- Category pills -->
  <div class="flex gap-2 overflow-x-auto pb-2 scrollbar-none" role="tablist"
       aria-label="Discussion categories">
    <button v-for="cat in discussionCategories" :key="cat.key"
            role="tab"
            :aria-selected="activeDiscussionTab === cat.key"
            class="shrink-0 rounded-full px-4 py-2 text-xs font-medium transition
                   focus-visible:outline-2 focus-visible:outline-[#1db954]"
            :class="activeDiscussionTab === cat.key
              ? 'bg-white/15 text-white shadow-lg'
              : 'bg-white/[0.04] text-white/40 hover:bg-white/[0.08] hover:text-white/60'"
            @click="activeDiscussionTab = cat.key">
      <span class="flex items-center gap-2">
        <Icon :name="cat.icon" class="h-3.5 w-3.5" aria-hidden="true" />
        {{ cat.label }}
      </span>
    </button>
  </div>

  <!-- Discussion feed -->
  <div class="space-y-2" role="feed" aria-label="Discussions">
    <DiscussionCard v-for="discussion in filteredDiscussions" :key="discussion.id"
                    :discussion="discussion" @reply="openReply" />
    <div v-if="!filteredDiscussions.length"
         class="flex flex-col items-center gap-3 py-12 text-center">
      <Icon name="message-circle" class="h-10 w-10 text-white/10" aria-hidden="true" />
      <p class="text-sm text-white/30">No discussions yet in this category</p>
      <button class="rounded-full bg-[#1db954]/10 px-5 py-2 text-xs font-semibold text-[#1db954]
                     transition hover:bg-[#1db954]/20">
        Start one
      </button>
    </div>
  </div>
</div>
```

### Create Modal — Step-by-step, rich
```vue
<Teleport to="body">
  <Transition name="modal">
    <div v-if="showCreateModal"
         class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
         @click.self="showCreateModal = false"
         role="dialog" aria-modal="true" :aria-label="`Create ${entityType}`">

      <div class="mx-4 w-full max-w-md rounded-2xl bg-[#121212] p-6 shadow-2xl ring-1 ring-white/10
                  motion-safe:animate-modal-in">
        <!-- Step indicator -->
        <div class="flex items-center gap-2 mb-6">
          <div v-for="step in 3" :key="step"
               class="h-1 flex-1 rounded-full transition-colors duration-300"
               :class="step <= createStep ? 'bg-[#1db954]' : 'bg-white/10'" />
        </div>

        <!-- Step content -->
        <template v-if="createStep === 1">
          <h2 class="text-lg font-bold text-white">{{ entityType === 'party' ? 'Start a Party' : entityType === 'room' ? 'Go Live' : 'Create Club' }}</h2>
          <p class="mt-1 text-sm text-white/40">Choose a name for your {{ entityType }}</p>
          <input v-model="createForm.title"
                 type="text"
                 :placeholder="entityType === 'club' ? 'Club name' : 'Party name'"
                 aria-label="Name"
                 class="mt-4 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm
                        text-white placeholder-white/20 outline-none transition
                        focus:border-white/20 focus:bg-white/[0.08]" />
        </template>

        <template v-if="createStep === 2">
          <h2 class="text-lg font-bold text-white">Details</h2>
          <p class="mt-1 text-sm text-white/40">Add a description and set privacy</p>
          <textarea v-model="createForm.description"
                    placeholder="What's this about?"
                    rows="3"
                    aria-label="Description"
                    class="mt-4 w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm
                           text-white placeholder-white/20 outline-none transition
                           focus:border-white/20 focus:bg-white/[0.08]" />
          <label class="mt-4 flex items-center gap-3 cursor-pointer">
            <input v-model="createForm.is_public" type="checkbox"
                   class="h-5 w-5 rounded border-white/10 bg-white/5 accent-[#1db954]" />
            <div>
              <span class="text-sm text-white">Public</span>
              <p class="text-xs text-white/30">Anyone can find and join</p>
            </div>
          </label>
        </template>

        <template v-if="createStep === 3">
          <div class="flex flex-col items-center gap-4 py-4 text-center">
            <div class="flex h-16 w-16 items-center justify-center rounded-full bg-[#1db954]/10">
              <Icon :name="entityIcon" class="h-8 w-8 text-[#1db954]" aria-hidden="true" />
            </div>
            <h2 class="text-lg font-bold text-white">Almost there!</h2>
            <p class="text-sm text-white/40">Review and launch your {{ entityType }}</p>
            <div class="w-full rounded-xl bg-white/[0.04] p-4 text-start">
              <p class="text-xs text-white/30">Name</p>
              <p class="text-sm font-medium text-white">{{ createForm.title }}</p>
              <p v-if="createForm.description" class="mt-2 text-xs text-white/30">Description</p>
              <p v-if="createForm.description" class="text-sm text-white/60">{{ createForm.description }}</p>
              <p class="mt-2 text-xs text-white/30">Visibility</p>
              <p class="text-sm text-white/60">{{ createForm.is_public ? '🌍 Public' : '🔒 Private' }}</p>
            </div>
          </div>
        </template>

        <!-- Navigation -->
        <div class="mt-6 flex gap-3">
          <button v-if="createStep > 1"
                  class="flex-1 rounded-xl bg-white/5 py-3 text-sm font-medium text-white/50
                         transition hover:bg-white/10"
                  @click="createStep--">
            Back
          </button>
          <button v-else
                  class="flex-1 rounded-xl bg-white/5 py-3 text-sm font-medium text-white/50
                         transition hover:bg-white/10"
                  @click="showCreateModal = false">
            Cancel
          </button>
          <button v-if="createStep < 3"
                  :disabled="createStep === 1 && !createForm.title.trim()"
                  class="flex-1 rounded-xl bg-[#1db954] py-3 text-sm font-bold text-black
                         transition hover:bg-[#1ed760] disabled:opacity-40
                         focus-visible:outline-2 focus-visible:outline-white"
                  @click="createStep++">
            Next
          </button>
          <button v-else
                  :disabled="isCreating"
                  class="flex-1 rounded-xl bg-[#1db954] py-3 text-sm font-bold text-black
                         transition hover:bg-[#1ed760] disabled:opacity-40
                         focus-visible:outline-2 focus-visible:outline-white"
                  @click="handleCreate">
            <span v-if="isCreating" class="inline-flex items-center gap-2">
              <span class="inline-block h-4 w-4 animate-spin rounded-full border-2 border-black border-t-transparent" />
              Creating...
            </span>
            <span v-else>Launch 🚀</span>
          </button>
        </div>
      </div>
    </div>
  </Transition>
</Teleport>
```

---

## States Reference

| State | Hero | Live Now Card | Activity River | Card Grid | Create Modal |
|-------|------|---------------|----------------|-----------|-------------|
| **Loading** | Skeleton gradient pulse | Waveform shimmer | 4 skeleton ActivityItems | 6 skeleton cards | — |
| **Loaded** | Greeting + party prompt | Featured room | Live stream | Cards grid | Closed |
| **Empty** (no rooms/parties) | — | "No rooms yet — go live!" + CTA | "No activity yet" | Empty state per tab | — |
| **Error** | Retry link in hero | Error state card | "Could not load" | Retry button | Toast error |
| **Creating** | — | — | — | Optimistic add | Step 3 spinner |
| **Joined** | — | Update listener count | "You joined" appears | Card updates | Close + success toast |
| **RTL** | Greeting in Persian first | All gaps logical | Timeline on right side | Cards flipped | All inputs start-aligned |

---

## Accessibility Fixes (Compared to Current)

| Current Issue | Fix |
|---------------|-----|
| Emoji icons (👥 🎵 👤 🏛) | Replaced with Lucide/PrimeIcons SVGs |
| ⟳ refresh (Unicode symbol) | Replaced with `<Icon name="refresh">` |
| No role/aria on cards | `role="article"` on cards, `aria-label` on all buttons |
| Filter discussion requires manual input | Category pill tabs with `role="tablist"` |
| Generic create form | Step-by-step wizard with clear states |
| No `focus-visible` on card buttons | Added `focus-visible:outline-*` to all interactive elements |
| No keyboard navigation for cards | Cards are `<div>` but inner buttons are `<button>` ✅ |

---

## Data Flow
```
PageSocialHub.vue
├── on mount:
│   ├── socialStore.fetchActivity(limit: 20) → activityRiver
│   ├── socialStore.fetchFeatured() → featuredRoom
│   ├── socialStore.fetchParties({ status: 'active', limit: 20 }) → trendingParties
│   ├── socialStore.fetchRooms({ limit: 10 }) → liveRooms
│   ├── socialStore.fetchMyClubs() → yourClubs
│   └── socialStore.fetchDiscussions({ limit: 30 }) → discussions
├── WebSocket:
│   └── socialSocket.on('activity:new') → prepend to activityRiver
├── playerStore.currentTrack → "Start a Party" inline CTA
├── userStore → greeting personalization
└── Cleanup: disconnect socket on unmount
```
