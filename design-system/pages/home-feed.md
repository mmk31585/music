# Home Feed (/)

> Enhanced from existing `PageHome.vue` with UI/UX Pro Max recommendations

---

## Layout (Mobile-first)

```
┌─────────────────────────────────────────────┐
│  Header: "عصر بخیر، علی" (Good evening)      │
│  ┌─────────────────────────────────────────┐ │
│  │ Mood pills (horizontal scroll)           │ │
│  │ [Focus 🎯] [Chill 🧘] [Party 🎉]        │ │
│  │ [Workout 💪] [Romance 💕] [Sad 🥲]      │ │
│  └─────────────────────────────────────────┘ │
│                                               │
│  Made for You (Section)                        │
│  ┌─────────────────────────────────────────┐ │
│  │ [Card] [Card] ← horizontal scroll snap   │ │
│  │ [Card] [Card]                            │ │
│  └─────────────────────────────────────────┘ │
│                                               │
│  Continue Listening (Section)                  │
│  ┌─────────────────────────────────────────┐ │
│  │ [Row] Track Title - Artist          3:45│ │
│  │ [Row] Track Title - Artist          4:12│ │
│  │ [Row] Track Title - Artist          3:30│ │
│  └─────────────────────────────────────────┘ │
│                                               │
│  Trending Now (Section with live badges)       │
│  ┌─────────────────────────────────────────┐ │
│  │ [Card] [Card] [Card] [Card] [Card]      │ │
│  └─────────────────────────────────────────┘ │
│                                               │
│  Popular Artists (horizontal scroll circles)   │
│  ┌─────────────────────────────────────────┐ │
│  │ (O) (O) (O) (O) (O)                     │ │
│  └─────────────────────────────────────────┘ │
│                                               │
│  New Releases + Persian Traditional +         │
│  More sections...                              │
└─────────────────────────────────────────────┘
```

## Component Specs

### Mood Pills
```vue
<button class="h-10 shrink-0 rounded-full px-5 text-sm font-medium
               transition-all duration-200
               bg-white/5 text-white/60 hover:bg-white/10 hover:text-white
               active:scale-95
               focus-visible:outline-2 focus-visible:outline-[#1db954]"
        aria-label="Focus mood">
  <span class="flex items-center gap-2">
    <Icon name="focus" class="h-4 w-4" aria-hidden="true" />
    <span>Focus</span>
  </span>
</button>
```
- Touch target: 40px height + 40px padding = exceeded ✅
- RTL: `gap-2` works in both directions ✅
- Active state: `scale-95` for press feedback ✅

### Continue Listening (Row Layout)
```vue
<TrackRow>
  <template #prefix>
    <button aria-label="Continue playing"
            class="flex h-10 w-10 items-center justify-center rounded-full
                   bg-[#1db954] text-black transition-transform
                   hover:scale-105 active:scale-95"
            @click="continueTrack(track)">
      <Icon name="play" class="h-4 w-4 ml-0.5" aria-hidden="true" />
    </button>
  </template>
  <template #info>
    <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
    <p class="truncate text-xs text-white/50">{{ track.artist }}</p>
  </template>
  <template #duration>
    <time class="text-xs text-white/40 tabular-nums">{{ track.duration }}</time>
  </template>
</TrackRow>
```
- Progress bar: subtle 2px green line under the row if partially played
- Touch target: 56px row height ✅
- RTL: `ml-0.5` → use `ms-0.5` for RTL support → **FIX**

### Trending Badges
```vue
<span class="inline-flex items-center gap-1 rounded-full bg-[#1db954]/10
             px-2 py-0.5 text-[10px] font-semibold text-[#1db954]"
      aria-label="Live trend: 1.2M plays this week">
  <span class="h-1.5 w-1.5 rounded-full bg-[#1db954] animate-pulse" aria-hidden="true" />
  1.2M
</span>
```

## States

| State | Visual | Action |
|-------|--------|--------|
| **Loading** | Skeleton grid (6 AlbumCard shimmer placeholders) | — |
| **Empty** | "Welcome! Start by exploring music" + CTA button | Navigate to Discover |
| **Error** | "Couldn't load recommendations" + Retry button | Refetch API |
| **Offline** | Offline banner + cached content | Show downloaded tracks |
| **RTL Active** | Cards slide from right → left | Logical CSS properties |

## Persian/Arabic Specific

| Element | RTL Adjustment |
|---------|---------------|
| Mood pills order | Right-to-left scroll |
| Time display | Use Persian digits (۱۲۳ not 123) via `Intl.NumberFormat('fa-IR')` |
| Greeting | "عصر بخیر، علی" — time-aware Persian greeting |
| Song count | "۱۲ آهنگ" (12 tracks) — number + Persian unit |
| Direction | `dir="rtl"` on <html>, use `text-start` not `text-left` |

## Data Flow
```
PageHome.vue
├── On mount: fetchRecommendations(userId)
├── userStore.greeting → computed: getTimeAwareGreeting(locale)
├── moodStore.moods → from API based on user history
├── trackStore.recommendations → Made for You section
├── historyStore.recent → Continue Listening section
├── catalogStore.trending → Trending Now section
└── catalogStore.newReleases → New Releases section
```
