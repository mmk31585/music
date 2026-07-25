# Album Detail (/albums/:id)

> Enhanced from existing `PageAlbum.vue`

---

## Layout (Desktop)

```
┌─────────────────────────────────────────────────────────────┐
│  ┌───────────────────────────────────────────────────────┐  │
│  │  ┌─────────┐  ┌────────────────────────────────────┐  │  │
│  │  │         │  │  Album Title (h1)                    │  │  │
│  │  │  Album  │  │  Artist Name (link)                   │  │  │
│  │  │  Cover  │  │                                       │  │  │
│  │  │ 450×450│  │  Year · Tracks count · Duration       │  │  │
│  │  │         │  │                                       │  │  │
│  │  │         │  │  [▶ Play] [♡ Save] [⋯ More]  │  │  │
│  │  └─────────┘  │                                       │  │  │
│  │               │  Genre tags: [Pop] [Electronic]       │  │  │
│  │               └────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                              │
│  Tracklist                                                    │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  #  Title                      Artist    Duration    │  │
│  │  1  Track One               Artist        3:45  [♡] │  │
│  │  2  Track Two               Artist        4:12  [♡] │  │
│  │  3  Track Three  ▶ (playing) Artist       3:30  [♥] │  │
│  │  4  Track Four              Artist        5:01  [♡] │  │
│  │  ...                                                   │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                              │
│  More by this artist                                          │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ [Card] [Card] [Card] [Card] ← horizontal scroll     │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Component Specs

### Album Header (Glass Hero)
```vue
<div class="relative overflow-hidden rounded-[2rem] p-8 md:p-10">
  <!-- Background: album art blurred -->
  <div class="pointer-events-none absolute inset-0" aria-hidden="true">
    <img :src="album.coverUrl" alt=""
         class="h-full w-full scale-110 object-cover opacity-30"
         style="filter: blur(60px) saturate(1.5)" />
    <div class="absolute inset-0 bg-gradient-to-t from-[#08080A] via-[#08080A]/60 to-transparent" />
  </div>

  <div class="relative z-10 flex flex-col gap-6 sm:flex-row sm:items-end">
    <!-- Cover art -->
    <div class="h-48 w-48 shrink-0 overflow-hidden rounded-2xl shadow-2xl ring-1 ring-white/10
                md:h-56 md:w-56">
      <img :src="album.coverUrl" :alt="album.title"
           class="h-full w-full object-cover" loading="lazy" />
    </div>

    <!-- Meta -->
    <div class="flex flex-col gap-3">
      <p class="text-xs font-bold tracking-[0.25em] text-white/30 uppercase">Album</p>
      <h1 class="text-3xl font-black text-white md:text-5xl">{{ album.title }}</h1>
      <RouterLink :to="`/artists/${album.artistId}`"
                  class="text-sm font-semibold text-white/70 hover:text-white transition-colors">
        {{ album.artistName }}
      </RouterLink>
      <div class="flex items-center gap-2 text-xs text-white/40">
        <span>{{ album.year }}</span>
        <span aria-hidden="true">·</span>
        <span>{{ formattedTrackCount }}</span>
        <span aria-hidden="true">·</span>
        <time datetime="PT{{ totalDuration }}S">{{ formattedDuration }}</time>
      </div>
    </div>
  </div>
</div>
```

### Track Row (With "Now Playing" State)
```vue
<div role="list" aria-label="Track list">
  <button v-for="(track, index) in album.tracks" :key="track.id"
          @click="playTrack(index)"
          @keydown.enter="playTrack(index)"
          @keydown.space.prevent="playTrack(index)"
          class="group flex w-full items-center gap-4 rounded-xl px-4 py-3
                 transition-colors hover:bg-white/[0.06]
                 focus-visible:outline-2 focus-visible:outline-[#1db954]
                 focus-visible:-outline-offset-2"
          :class="currentTrackId === track.id ? 'bg-[#1db954]/5' : ''"
          :aria-current="currentTrackId === track.id ? 'true' : undefined"
          role="listitem">

    <!-- Index or playing indicator -->
    <span class="flex h-8 w-8 shrink-0 items-center justify-center text-sm"
          aria-hidden="true">
      <template v-if="currentTrackId === track.id">
        <EqualizerIcon class="h-4 w-4 text-[#1db954]" />
      </template>
      <template v-else>
        <span class="text-white/30 tabular-nums group-hover:hidden">{{ index + 1 }}</span>
        <Icon name="play" class="hidden h-3 w-3 text-white group-hover:block ml-0.5" />
      </template>
    </span>

    <!-- Title + artist -->
    <div class="flex-1 truncate text-start">
      <p class="truncate text-sm font-medium"
         :class="currentTrackId === track.id ? 'text-[#1db954]' : 'text-white'">
        {{ track.title }}
      </p>
      <p class="truncate text-xs text-white/50">{{ track.artist }}</p>
    </div>

    <!-- Like button -->
    <button @click.stop="toggleLike(track.id)"
            aria-label="Like track"
            class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full
                   text-white/30 hover:text-white/70 hover:bg-white/[0.06]
                   transition-colors opacity-0 group-hover:opacity-100
                   group-focus-within:opacity-100 focus-visible:opacity-100
                   focus-visible:outline-2 focus-visible:outline-[#1db954]">
      <Icon :name="isLiked(track.id) ? 'heart-solid' : 'heart'"
            class="h-4 w-4" aria-hidden="true" />
    </button>

    <!-- Duration -->
    <time class="w-12 text-end text-xs tabular-nums text-white/40"
          :datetime="`PT${track.duration}S`">
      {{ formatDuration(track.duration) }}
    </time>
  </button>
</div>
```

### Action Buttons
```vue
<div class="flex items-center gap-4">
  <!-- Primary: Play -->
  <button aria-label="Play album"
          class="inline-flex items-center gap-2 rounded-full bg-[#1db954] px-8 py-3
                 text-sm font-bold text-black transition hover:scale-105
                 hover:bg-[#1ed760] active:scale-[1.02]
                 focus-visible:outline-2 focus-visible:outline-white">
    <Icon name="play-fill" class="h-4 w-4" aria-hidden="true" />
    Play
  </button>

  <!-- Secondary: Save -->
  <button aria-label="Save to library"
          class="inline-flex items-center justify-center rounded-full border
                 border-white/15 bg-white/[0.04] p-3 text-white backdrop-blur
                 transition hover:bg-white/[0.10] hover:border-white/25
                 active:scale-95 focus-visible:outline-2 focus-visible:outline-[#1db954]">
    <Icon name="heart" class="h-5 w-5" aria-hidden="true" />
  </button>

  <!-- More menu -->
  <button aria-label="More options"
          class="inline-flex items-center justify-center rounded-full
                 text-white/50 p-3 transition hover:bg-white/[0.06] hover:text-white
                 focus-visible:outline-2 focus-visible:outline-[#1db954]"
          @click="showMenu = true"
          @keydown.escape="showMenu = false">
    <Icon name="more-horizontal" class="h-5 w-5" aria-hidden="true" />
  </button>
</div>
```

## States

| State | Cover Area | Tracklist |
|-------|-----------|-----------|
| **Loading** | Skeleton rect 224×224 + text lines | 8 skeleton row placeholders |
| **Playing** | Equalizer pulse overlay on cover | Green highlight on active row |
| **Saved** | Heart icon filled (♥) | — |
| **Not Saved** | Heart icon outline (♡) | — |
| **Empty** (error/fetch fail) | "Album not found" + Go back | — |
| **RTL** | Cover on right side, meta on left | Index numbers → Persian digits |

## Persian/Arabic Specific

| Element | RTL Adjustment |
|---------|---------------|
| Text alignment | `text-start` not `text-left` |
| Play icon offset | `ml-0.5` → `ms-0.5` (logical) |
| Duration alignment | `text-end` auto-adjusts |
| Track numbers | `font-variant-numeric: tabular-nums` + `Intl` formatting |
| Genre tag order | Right-to-left reading: ["پاپ", "الکترونیک"] |
| Year display | ۱۴۰۴ (Jalali) instead of 2025 |

## Data Flow
```
PageAlbum.vue
├── on beforeRouteEnter → fetchAlbum(route.params.id)
├── albumStore.currentAlbum → album, tracks, metadata
├── playerStore → currentTrackId, isPlaying
├── libraryStore → isLiked(trackId)
├── historyStore → track play events for reshuffle
└── playerStore.playAlbum(album.id) → load all tracks into queue
```
