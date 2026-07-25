# User Profile (/profile/:id)

> Complete redesign of `PageUserProfile.vue` + `UserHero.vue`, `ProfileTabs.vue`
> Fixes: flat hero, hidden music status, boring tabs, no stats, weak empty states

---

## Design Philosophy

The profile is **your identity on Muse** — it should feel like a living music diary:
- **Cinematic hero** — colors adapt to your music taste (your top genre drives the aurora)
- **"Now Playing" is front and center** — it's the first thing people see
- **Stats tell your story** — listening hours, top genres, badges earned
- **Visual grid > text lists** — show album art, avatar walls, cover mosaics
- **Persian-first with global polish** — RTL typography, Jalali dates, Tomans for creators

---

## Layout (Desktop)

```
┌───────────────────────────────────────────────────────────────┐
│  ╔═══════════════════════════════════════════════════════════╗ │
│  ║                   CINEMATIC HERO                          ║ │
│  ║  ┌──────────┐                                            ║ │
│  ║  │          │   Sara Ahmadi                               ║ │
│  ║  │  Avatar  │   @sara_ahmadi                              ║ │
│  ║  │          │                                             ║ │
│  ║  └──────────┘   ۱,۲۳۴ followers    ۵۶۷ following     ║ │
│  ║                                                           ║ │
│  ║  [🎵 Now Playing: Gole Yakh — Mohsen Chavoshi   ⬇]    ║ │
│  ║                                                           ║ │
│  ║  [Follow ✦] [Share ↗] [⋯ More]                           ║ │
│  ╚═══════════════════════════════════════════════════════════╝ │
│                                                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐           │
│  │  ۴۲۸ hrs    │  │  Indie ۴۲%  │  │  🏅 ۱۲      │           │
│  │  Listening  │  │  Top Genre  │  │  Badges     │           │
│  └─────────────┘  └─────────────┘  └─────────────┘           │
│                                                               │
│  Tabs: [🎵 Tracks] [💿 Albums] [✏️ Edits] [👥 Followers]     │
│                                                               │
│  ════════════════════ TRACKS TAB ════════════════════════════ │
│                                                               │
│  [Mosaic view toggle] [List view toggle]                     │
│                                                               │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐               │
│  │ ♥    │ │ ♥    │ │ ♥    │ │ ♥    │ │ ♥    │               │
│  │ Cover│ │ Cover│ │ Cover│ │ Cover│ │ Cover│               │
│  │ Title│ │ Title│ │ Title│ │ Title│ │ Title│               │
│  │Artist│ │Artist│ │Artist│ │Artist│ │Artist│               │
│  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘               │
│                                                               │
│  [Show more ▼]                                                │
└───────────────────────────────────────────────────────────────┘
```

---

## Component Specs

### Cinematic Hero (replaces UserHero.vue)
```vue
<div class="relative overflow-hidden rounded-[2rem]">
  <!-- Aurora gradient driven by top genre color -->
  <div class="pointer-events-none absolute inset-0" aria-hidden="true">
    <!-- Dynamic gradient from user's top genre (or fallback green) -->
    <div class="absolute inset-0 bg-gradient-to-br"
         :style="heroGradient" />
    <!-- Floating orbs -->
    <div class="absolute -top-20 -right-20 h-80 w-80 rounded-full blur-3xl opacity-20"
         :style="{ background: heroOrbColor }" />
    <div class="absolute -bottom-16 -left-16 h-60 w-60 rounded-full blur-3xl opacity-15"
         :style="{ background: heroOrbColor2 }" />
    <div class="absolute inset-0 bg-gradient-to-t from-[#08080A] via-[#08080A]/20 to-transparent" />
  </div>

  <div class="relative z-10 flex flex-col gap-6 px-6 pt-16 pb-8 md:flex-row md:items-end md:gap-10 md:pt-12 md:pb-10">
    <!-- Avatar with floating ring -->
    <div class="relative shrink-0 self-center md:self-end">
      <!-- Pulse ring if currently playing -->
      <div v-if="musicStatus?.playing"
           class="absolute -inset-2 animate-ping rounded-full border-2 border-[#1db954] opacity-30" />
      <div class="h-36 w-36 overflow-hidden rounded-full border-4 shadow-2xl md:h-44 md:w-44"
           :class="musicStatus?.playing ? 'border-[#1db954]' : 'border-white/10'">
        <img v-if="avatarUrl" :src="avatarUrl" :alt="displayName"
             loading="lazy" class="h-full w-full object-cover"
             @error="onImgError" />
        <div v-else class="flex h-full w-full items-center justify-center
                          bg-gradient-to-br from-[#1db954]/30 to-[#a855f7]/30 text-4xl text-white/40">
          <Icon name="user" class="h-12 w-12" aria-hidden="true" />
        </div>
      </div>
      <!-- Badge indicators -->
      <div v-if="isCreator"
           class="absolute -bottom-1 -right-1 flex h-7 w-7 items-center justify-center rounded-full
                  bg-[#1db954] ring-2 ring-[#08080A] shadow-lg"
           aria-label="Creator">
        <Icon name="check" class="h-3.5 w-3.5 text-black" aria-hidden="true" />
      </div>
    </div>

    <div class="flex flex-col items-center text-center md:items-start md:text-start">
      <!-- Name + handle -->
      <h1 class="text-3xl font-black text-white md:text-5xl">{{ displayName }}</h1>
      <p v-if="handle" class="mt-1 text-sm text-white/40">@{{ handle }}</p>
      <p v-if="bio" class="mt-2 max-w-md text-sm text-white/50 leading-relaxed">{{ bio }}</p>

      <!-- Stats row -->
      <div class="mt-4 flex flex-wrap items-center justify-center gap-5 text-sm md:justify-start">
        <button @click="$emit('showFollowers')"
                class="transition hover:text-white focus-visible:outline-2 focus-visible:outline-[#1db954]"
                aria-label="View followers">
          <span class="font-bold text-white tabular-nums">{{ formatNumber(followerCount) }}</span>
          <span class="text-white/40"> followers</span>
        </button>
        <button @click="$emit('showFollowing')"
                class="transition hover:text-white focus-visible:outline-2 focus-visible:outline-[#1db954]"
                aria-label="View following">
          <span class="font-bold text-white tabular-nums">{{ formatNumber(followingCount) }}</span>
          <span class="text-white/40"> following</span>
        </button>
        <span class="text-white/30" aria-hidden="true">·</span>
        <span class="text-white/40">
          <span class="font-medium text-white/60">Joined</span>
          {{ formatJoinDate(joinDate) }}
        </span>
      </div>

      <!-- NOW PLAYING — prominent, animated -->
      <div v-if="musicStatus?.playing && musicStatus.currentTrack"
           class="mt-4 flex w-full max-w-md items-center gap-3 rounded-2xl
                  bg-white/[0.06] p-3 ring-1 ring-white/[0.10] motion-safe:animate-fade-in-up">
        <div class="h-12 w-12 shrink-0 overflow-hidden rounded-xl shadow-lg">
          <img :src="musicStatus.currentTrack.coverUrl" alt=""
               class="h-full w-full object-cover" />
        </div>
        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-2">
            <span class="flex gap-0.5" aria-hidden="true">
              <span v-for="i in 4" :key="i"
                    class="h-3 w-0.5 rounded-full bg-[#1db954] motion-safe:animate-equalizer"
                    :style="{ animationDelay: `${i * 100}ms` }" />
            </span>
            <p class="truncate text-xs font-semibold text-white">{{ musicStatus.currentTrack.title }}</p>
          </div>
          <p class="truncate text-[10px] text-white/40 ml-4">{{ musicStatus.currentTrack.artist }}</p>
        </div>
        <button aria-label="Listen along"
                class="shrink-0 rounded-full bg-[#1db954] px-4 py-1.5 text-[10px] font-bold text-black
                       transition hover:bg-[#1ed760] hover:scale-105 active:scale-95
                       focus-visible:outline-2 focus-visible:outline-white">
          Listen
        </button>
      </div>

      <!-- No track status -->
      <div v-else-if="musicStatus && !musicStatus.playing"
           class="mt-4 text-xs text-white/30 flex items-center gap-2">
        <Icon name="pause-circle" class="h-4 w-4" aria-hidden="true" />
        Not listening right now
      </div>

      <!-- Action buttons -->
      <div class="mt-5 flex flex-wrap items-center gap-3">
        <button v-if="!isOwnProfile"
                :class="[
                  'inline-flex items-center gap-2 rounded-full px-6 py-3 text-sm font-bold transition',
                  'focus-visible:outline-2 focus-visible:outline-white active:scale-[0.97]',
                  isFollowing
                    ? 'border border-[#1db954]/50 bg-[#1db954]/10 text-[#1db954] hover:bg-[#1db954]/20'
                    : 'bg-[#1db954] text-black hover:bg-[#1ed760] hover:scale-105',
                ]"
                @click="$emit('toggleFollow')">
          <Icon :name="isFollowing ? 'check' : 'plus'" class="h-4 w-4" aria-hidden="true" />
          {{ isFollowing ? 'Following' : 'Follow' }}
        </button>
        <button aria-label="Share profile"
                class="inline-flex items-center justify-center rounded-full border border-white/15
                       bg-white/[0.04] p-3 text-white/60 backdrop-blur transition
                       hover:bg-white/[0.10] hover:text-white active:scale-95
                       focus-visible:outline-2 focus-visible:outline-[#1db954]">
          <Icon name="share" class="h-5 w-5" aria-hidden="true" />
        </button>
        <button v-if="isOwnProfile"
                class="inline-flex items-center gap-2 rounded-full border border-white/15
                       bg-white/[0.04] px-5 py-3 text-sm font-medium text-white/60 backdrop-blur
                       transition hover:bg-white/[0.10] hover:text-white
                       focus-visible:outline-2 focus-visible:outline-[#1db954]">
          <Icon name="settings" class="h-4 w-4" aria-hidden="true" />
          Edit Profile
        </button>
      </div>
    </div>
  </div>
</div>
```

### Stats Gallery (new)
```vue
<div class="grid grid-cols-2 gap-3 md:grid-cols-4">
  <!-- Listening hours -->
  <StatCard icon="clock" :value="formatHours(stats.listeningHours)" label="Listening"
            :trend="stats.hoursTrend" />

  <!-- Top genre -->
  <StatCard icon="music" :value="stats.topGenre.name" :sub="`${stats.topGenre.percent}%`"
            label="Top Genre" :color="stats.topGenre.color" />

  <!-- Badge count -->
  <RouterLink :to="isOwnProfile ? '/gamification' : ''"
              class="rounded-2xl bg-white/[0.06] p-4 ring-1 ring-white/[0.10]
                     transition hover:bg-white/[0.08] group
                     focus-visible:outline-2 focus-visible:outline-[#1db954]">
    <div class="flex items-center justify-between">
      <Icon name="award" class="h-5 w-5 text-amber-400" aria-hidden="true" />
      <Icon name="chevron-right" class="h-4 w-4 text-white/20 group-hover:text-white/40
                                        transition" aria-hidden="true" />
    </div>
    <p class="mt-3 text-2xl font-black text-white tabular-nums">{{ stats.badgeCount }}</p>
    <p class="text-xs text-white/40">Badges earned</p>
  </RouterLink>

  <!-- Top track -->
  <StatCard v-if="stats.topTrack" icon="award" :value="stats.topTrack.title"
            :sub="stats.topTrack.plays + ' plays'" label="Most Played" />
  <div v-else class="rounded-2xl bg-white/[0.06] p-4 ring-1 ring-white/[0.10]">
    <Icon name="bar-chart" class="h-5 w-5 text-white/20" aria-hidden="true" />
    <p class="mt-3 text-sm text-white/30">No stats yet</p>
  </div>
</div>
```

### Tabs — Pill-style (replaces ProfileTabs.vue underline)
```vue
<div class="flex gap-1 rounded-xl bg-white/[0.04] p-1 overflow-x-auto scrollbar-none"
     role="tablist" :aria-label="`${displayName}'s profile`">
  <button v-for="tab in tabs" :key="tab.key"
          role="tab"
          :aria-selected="activeTab === tab.key"
          :aria-controls="`tabpanel-${tab.key}`"
          class="flex items-center gap-2 rounded-lg px-4 py-2.5 text-sm font-medium
                 transition-all duration-200 whitespace-nowrap
                 focus-visible:outline-2 focus-visible:outline-[#1db954]"
          :class="activeTab === tab.key
            ? 'bg-white/10 text-white shadow-lg'
            : 'text-white/30 hover:text-white/50'"
          @click="$emit('update:activeTab', tab.key)">
    <Icon :name="tab.icon" class="h-4 w-4" aria-hidden="true" />
    {{ tab.label }}
  </button>
</div>
```

Tabs definition:
```ts
const tabs = [
  { key: 'tracks', label: 'آهنگ‌ها', icon: 'music' },
  { key: 'albums', label: 'آلبوم‌ها', icon: 'disc' },
  { key: 'edits', label: 'ادیت‌ها', icon: 'video' },
  { key: 'followers', label: 'دنبال‌کننده‌ها', icon: 'users' },
  { key: 'following', label: 'دنبال‌شونده‌ها', icon: 'user-plus' },
]
```

### Tracks Tab — Mosaic View (new)
```vue
<div role="tabpanel" id="tabpanel-tracks"
     :aria-labelledby="`tab-tracks`">
  <!-- View toggle -->
  <div class="flex items-center justify-between mb-4">
    <p class="text-xs text-white/30 tabular-nums">{{ likedTracks.length }} tracks</p>
    <div class="flex gap-1 rounded-lg bg-white/[0.04] p-0.5">
      <button aria-label="Mosaic view"
              class="rounded-md p-1.5 transition"
              :class="trackViewMode === 'mosaic' ? 'bg-white/10 text-white' : 'text-white/30 hover:text-white/50'"
              @click="trackViewMode = 'mosaic'">
        <Icon name="grid" class="h-4 w-4" aria-hidden="true" />
      </button>
      <button aria-label="List view"
              class="rounded-md p-1.5 transition"
              :class="trackViewMode === 'list' ? 'bg-white/10 text-white' : 'text-white/30 hover:text-white/50'"
              @click="trackViewMode = 'list'">
        <Icon name="list" class="h-4 w-4" aria-hidden="true" />
      </button>
    </div>
  </div>

  <!-- MOSAIC VIEW -->
  <div v-if="trackViewMode === 'mosaic'"
       class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
    <div v-for="track in likedTracks" :key="track.id"
         class="group relative overflow-hidden rounded-2xl bg-white/[0.04] ring-1 ring-white/[0.06]
                transition-all duration-300 hover:ring-[#1db954]/30 hover:bg-white/[0.06]
                focus-within:ring-[#1db954]">
      <!-- Cover art -->
      <div class="aspect-square overflow-hidden">
        <img v-if="track.cover_url" :src="track.cover_url"
             :alt="track.track_title" loading="lazy"
             class="h-full w-full object-cover transition duration-500
                    group-hover:scale-105" />
        <div v-else class="flex h-full items-center justify-center bg-gradient-to-br from-white/[0.04] to-white/[0.02]">
          <Icon name="music" class="h-8 w-8 text-white/20" aria-hidden="true" />
        </div>
      </div>

      <!-- Overlay on hover -->
      <div class="absolute inset-0 flex items-center justify-center gap-2
                  bg-black/40 opacity-0 group-hover:opacity-100
                  group-focus-within:opacity-100 transition-opacity duration-200">
        <button aria-label="Play"
                class="flex h-10 w-10 items-center justify-center rounded-full bg-[#1db954]
                       text-black shadow-lg transition hover:scale-110 active:scale-95
                       focus-visible:outline-2 focus-visible:outline-white">
          <Icon name="play" class="h-5 w-5 ml-0.5" aria-hidden="true" />
        </button>
      </div>

      <!-- Info -->
      <div class="p-3">
        <p class="truncate text-sm font-medium text-white">{{ track.track_title }}</p>
        <p class="truncate text-xs text-white/40">{{ track.artist_name }}</p>
      </div>

      <!-- Visibility badge (own profile) -->
      <div v-if="isOwnProfile" class="absolute top-2 left-2">
        <span class="flex h-6 w-6 items-center justify-center rounded-full
                     backdrop-blur-sm text-[10px]"
              :class="track.is_public
                ? 'bg-black/30 text-white/60'
                : 'bg-amber-500/30 text-amber-300'"
              :title="track.is_public ? 'Public' : 'Private'">
          <Icon :name="track.is_public ? 'globe' : 'lock'" class="h-3 w-3" aria-hidden="true" />
        </span>
      </div>
    </div>
  </div>

  <!-- LIST VIEW (original, enhanced) -->
  <div v-else class="space-y-1">
    <div v-for="(track, i) in likedTracks" :key="track.id"
         class="group flex items-center gap-3 rounded-xl px-3 py-2.5 transition
                hover:bg-white/[0.04] focus-within:bg-white/[0.04]">
      <!-- Cover thumbnail -->
      <div class="h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10">
        <img v-if="track.cover_url" :src="track.cover_url" alt=""
             class="h-full w-full object-cover" loading="lazy" />
        <div v-else class="flex h-full items-center justify-center">
          <Icon name="music" class="h-4 w-4 text-white/30" aria-hidden="true" />
        </div>
      </div>
      <!-- Info -->
      <div class="min-w-0 flex-1">
        <p class="truncate text-sm font-medium text-white">{{ track.track_title }}</p>
        <p class="truncate text-xs text-white/40">{{ track.artist_name }}</p>
      </div>
      <!-- Visibility toggle (own profile) -->
      <button v-if="isOwnProfile"
              @click="toggleTrackVisibility(track, i)"
              class="flex h-8 w-8 items-center justify-center rounded-lg text-sm
                     opacity-0 group-hover:opacity-100 group-focus-within:opacity-100
                     transition hover:bg-white/10 focus-visible:opacity-100
                     focus-visible:outline-2 focus-visible:outline-[#1db954]"
              :class="track.is_public ? 'text-white/40' : 'text-amber-400/70'"
              :aria-label="track.is_public ? 'Set private' : 'Set public'">
        <Icon :name="track.is_public ? 'globe' : 'lock'" class="h-4 w-4" aria-hidden="true" />
      </button>
      <!-- Play -->
      <button aria-label="Play track"
              class="flex h-8 w-8 items-center justify-center rounded-full text-white/40
                     transition hover:bg-white/10 hover:text-white
                     focus-visible:outline-2 focus-visible:outline-[#1db954]"
              @click="playLikedTrack(track)">
        <Icon name="play" class="h-4 w-4 ml-0.5" aria-hidden="true" />
      </button>
    </div>
  </div>

  <!-- Empty state -->
  <div v-if="likedTracks.length === 0"
       class="flex flex-col items-center gap-4 py-20 text-center">
    <div class="flex h-20 w-20 items-center justify-center rounded-full bg-white/[0.04]">
      <Icon name="heart" class="h-8 w-8 text-white/10" aria-hidden="true" />
    </div>
    <div>
      <p class="text-base font-semibold text-white/40">
        {{ isOwnProfile ? 'No liked tracks yet' : 'Tracks are private' }}
      </p>
      <p class="mt-1 text-sm text-white/30">
        {{ isOwnProfile
          ? 'Heart tracks to save them here'
          : 'This user keeps their likes private' }}
      </p>
    </div>
    <RouterLink v-if="isOwnProfile" to="/"
                class="rounded-full bg-[#1db954] px-6 py-2.5 text-sm font-bold text-black
                       transition hover:bg-[#1ed760]">
      Discover music
    </RouterLink>
  </div>
</div>
```

### Followers/Following Tab — Avatar Wall (replaces flat list)
```vue
<div v-if="followers.length === 0"
     class="flex flex-col items-center gap-4 py-20 text-center">
  <div class="flex h-20 w-20 items-center justify-center rounded-full bg-white/[0.04]">
    <Icon name="users" class="h-8 w-8 text-white/10" aria-hidden="true" />
  </div>
  <p class="text-base font-semibold text-white/40">No followers yet</p>
  <p class="text-sm text-white/30">Share your profile to grow your community</p>
</div>

<div v-else class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
  <RouterLink v-for="f in followers" :key="f.follower_id"
              :to="`/profile/${f.follower_id}`"
              class="flex flex-col items-center gap-3 rounded-2xl bg-white/[0.04] p-5
                     ring-1 ring-white/[0.06] transition-all duration-200
                     hover:bg-white/[0.08] hover:ring-white/[0.12]
                     focus-visible:outline-2 focus-visible:outline-[#1db954]">
    <div class="h-16 w-16 overflow-hidden rounded-full ring-2 ring-white/10">
      <img v-if="f.avatar_url" :src="f.avatar_url" :alt="f.follower_name"
           class="h-full w-full object-cover" loading="lazy" />
      <div class="flex h-full w-full items-center justify-center bg-gradient-to-br
                  from-white/[0.08] to-white/[0.02] text-lg font-bold text-white/50">
        {{ (f.follower_name || f.follower_id).charAt(0).toUpperCase() }}
      </div>
    </div>
    <div class="text-center">
      <p class="truncate text-sm font-semibold text-white max-w-[100px]">{{ f.follower_name }}</p>
      <p class="text-[10px] text-white/30">Follows you</p>
    </div>
  </RouterLink>
</div>
```

---

## States Reference

| State | Hero | Stats | Tracks Tab | Followers Tab |
|-------|------|-------|-----------|---------------|
| **Loading** | Avatar skeleton + 3 text lines | 4 skeleton cards | 6 skeleton mosaic items | 6 skeleton avatar circles |
| **Loaded** | Full hero with music status | Genre, hours, badges | Mosaic or list view | Avatar wall |
| **Empty** (new user) | Show join date + "Start exploring" | "No stats yet" | "Heart tracks to save them" | "Share your profile" |
| **Error** | "Couldn't load profile" + Retry | — | — | — |
| **Another user's profile (private)** | Show public info only | Limited stats | "Tracks are private" | Show public followers |
| **Creator profile** | Creator badge on avatar | Show earnings stat | — | — |
| **RTL** | Persian text, right-aligned | Grid auto-flips | Mosaic grid auto-flips | Avatar wall auto-flips |

---

## Data Flow (Optimized)

```
PageUserProfile.vue
├── on beforeRouteEnter:
│   ├── userStore.fetchProfile(targetUserId) → avatar, name, bio, handle, joinDate
│   ├── socialStore.fetchFollowStats(targetUserId) → followers, following, counts
│   ├── reactionsStore.fetchLikedTracks(targetUserId, { limit: 50 })
│   ├── reactionsStore.fetchLikedAlbums(targetUserId, { limit: 10 })
│   ├── videoStore.fetchUserVideos(targetUserId)
│   ├── musicStatusStore.fetchStatus(targetUserId) → currentTrack, playing state
│   └── gamificationStore.fetchUserBadges(targetUserId) → badgeCount
│
├── Computed:
│   ├── heroGradient → derived from topGenre.color or currentTrack.cover palette
│   ├── stats → aggregated from all stores
│   └── isOwnProfile → computed from auth.user?.id === targetUserId
│
├── WebSocket (if watching profile):
│   └── socialSocket.on('user:status') → update musicStatus in real-time
│
└── Cleanup
```

## Accessibility Improvements

| Current Issue | Fix |
|---------------|-----|
| `"pi pi-user"` text icon | `<Icon name="user">` with `aria-hidden="true"` |
| No `role="tablist"` on tabs | Added `role="tablist"`, `role="tab"`, `aria-selected`, `aria-controls` |
| Flat follower list hard to tap | Avatar wall with larger touch targets (≥64px) |
| No `focus-visible` on hero buttons | Added to follow/share/edit buttons |
| Empty states are just text | Added illustration + contextual CTA |
| Loading has no `role="status"` | Added to skeleton containers |
| Music status has no live region | Added `aria-live="polite"` to status section |
