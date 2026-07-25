<template>
  <div class="relative mx-auto w-full max-w-7xl px-4 pb-36 pt-6 md:px-6 lg:px-8">
    <!-- Aurora gradient background -->
    <div class="aurora-bg pointer-events-none fixed inset-0" aria-hidden="true">
      <div class="aurora-spot-1 -top-40 -left-40 bg-purple-600/15" />
      <div class="aurora-spot-2 -top-60 -right-40 bg-blue-500/10" />
      <div
        class="aurora-spot-1 top-20 left-1/3 bg-violet-500/10"
        style="animation-delay: -8s; width: 350px; height: 350px"
      />
      <div class="aurora-flow absolute inset-0 opacity-30" />
    </div>

    <div class="relative z-10">
    <!-- ─────────────────────── Search Header ─────────────────────── -->
    <div class="mb-6">
      <p class="text-[10px] font-bold tracking-[0.3em] text-white/40 uppercase">
        {{ hasSearched ? 'Search results' : 'Discover' }}
      </p>
      <h1 class="mt-2 text-4xl font-black text-white md:text-5xl">
        {{ pageTitle }}
      </h1>
    </div>

    <!-- ─────────────────────── Search Bar ─────────────────────── -->
    <div class="relative mb-6">
      <IconField
        class="spring w-full rounded-full shadow-md transition-all duration-300 focus-within:ring-4 focus-within:ring-spotify/15"
      >
        <InputIcon class="text-white/40 text-base">
          <Search aria-hidden="true" class=""  />
        </InputIcon>
        <InputText
          ref="inputRef"
          v-model="query"
          dir="ltr"
          type="text"
          placeholder="What do you want to listen to?"
          aria-label="Search tracks, artists, and albums"
          class="w-full rounded-full border-0 bg-white/4 py-3.5 pl-11 pr-14 text-sm text-white placeholder:text-white/30 focus:bg-white/6"
          style="box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.08)"
          @input="onInput"
          @keydown.enter="doSearch"
          @focus="showRecent = true"
          @blur="hideRecentDelayed"
        />
        <div v-if="query" class="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-1.5">
          <Button
            v-tooltip.top="'Clear'"
            rounded
            text
            severity="secondary"
            aria-label="Clear search"
            class="h-7! w-7! p-0! min-w-0! text-white/40 hover:text-white hover:bg-white/10"
            @click="clearSearch"
          >
            <template #icon>
              <X aria-hidden="true" class="text-xs"  />
            </template>
          </Button>
        </div>
      </IconField>
    </div>

    <!-- ─────────────────── Recent Searches ─────────────────── -->
    <Transition name="fade">
      <div
        v-if="showRecent && !query && recentSearches.length && !hasSearched && !searching"
        class="mb-6 rounded-2xl border border-white/6 bg-white/2 p-4 backdrop-blur-xs"
        @mouseenter="showRecent = true"
        @mouseleave="showRecent = false"
      >
        <div class="mb-3 flex items-center justify-between">
          <h3 class="text-sm font-bold text-white">Recent searches</h3>
          <Button
            text
            severity="secondary"
            size="small"
            aria-label="Clear all recent searches"
            class="text-[11px]! px-2! py-0.5! text-white/40 hover:text-white"
            @click="clearRecentSearches()"
          >
            Clear
          </Button>
        </div>
        <div class="flex flex-wrap gap-2">
          <Chip
            v-for="entry in recentSearches"
            :key="entry.query + entry.timestamp"
            :label="entry.query"
            removable
            class="spring cursor-pointer border border-white/8 bg-white/4 text-xs text-white/70 hover:border-white/20 hover:bg-white/8 hover:text-white active:scale-95"
            @click="selectRecent(entry.query)"
            @remove="removeRecentSearch(entry.query)"
          >
            <template #icon>
              <History aria-hidden="true" class="text-[10px] text-white/30"  />
            </template>
          </Chip>
        </div>
      </div>
    </Transition>

    <!-- ─────────────────── Filter Tabs ─────────────────── -->
    <div v-if="hasSearched && !searching" class="mb-6 flex flex-wrap gap-2">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        class="spring rounded-full px-4 py-1.5 text-xs font-bold transition-all"
        :class="searchType === tab.value
          ? 'bg-white text-black shadow-md'
          : 'border border-white/8 bg-white/4 text-white/50 hover:border-white/20 hover:text-white'"
        @click="searchType = tab.value"
      >
        {{ tab.label }}
        <span v-if="getCount(tab.value)" class="ml-1.5 text-[10px] opacity-60">({{ getCount(tab.value) }})</span>
      </button>
    </div>

    <!-- ─────────────────── Search Loading Skeleton ─────────────────── -->
    <div v-if="searching" class="space-y-2" aria-live="polite">
      <div v-for="i in 6" :key="i" class="flex items-center gap-4 rounded-2xl bg-white/2 p-3">
        <Skeleton shape="rect" size="3.5rem" class="rounded-xl! bg-white/6!" />
        <div class="flex-1 space-y-2.5">
          <Skeleton width="60%" height="1rem" class="rounded! bg-white/6!" />
          <Skeleton width="35%" height="0.75rem" class="rounded! bg-white/4!" />
        </div>
      </div>
    </div>

    <!-- ═══════════════════════ SEARCH RESULTS ═══════════════════════ -->
    <template v-else-if="hasSearched">
      <AppEmptyState
        v-if="searchError"
        icon="pi pi-exclamation-circle"
        title="Search failed"
        description="Something went wrong. Check your connection and try again."
      />
      <AppEmptyState
        v-else-if="hasNoResults"
        icon="pi pi-search"
        title="No results found"
        :description="noResultsText"
      />

      <div v-else class="space-y-10" aria-live="polite">
        <!-- ── Songs ── -->
        <section v-if="filteredResults.tracks.length">
          <div class="mb-4 flex items-center gap-3">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-500/10">
              <PlayCircle aria-hidden="true" class="text-xs text-emerald-400"  />
            </div>
            <h2 class="text-lg font-bold text-white">Songs</h2>
            <Badge :value="filteredResults.tracks.length" severity="success" class="bg-emerald-500/15! text-emerald-400! text-[10px]! font-bold! min-w-5! h-5!" />
            <button
              class="ml-auto flex items-center gap-1.5 rounded-full bg-emerald-500/10 px-3 py-1 text-xs font-bold text-emerald-400 transition hover:bg-emerald-500/20"
              @click.stop="playAllSearchResults"
            >
              <Play aria-hidden="true" class="text-[10px]"  />
              Play All
            </button>
          </div>
          <div class="overflow-hidden rounded-2xl border border-white/6 bg-white/2">
            <button
              v-for="(track, i) in filteredResults.tracks"
              :key="track.id"
              type="button"
              class="group flex w-full items-center gap-3 px-4 py-2.5 transition hover:bg-white/3 active:bg-white/5"
              @click="playTrack(track)"
            >
              <span class="relative flex w-6 items-center justify-center text-xs text-white/30 tabular-nums">
                <span class="group-hover:hidden">{{ i + 1 }}</span>
                <Play aria-hidden="true" class="absolute hidden text-sm text-white group-hover:block"  />
              </span>
              <div class="h-11 w-11 shrink-0 overflow-hidden rounded-lg bg-white/4">
                <img
                  v-if="track.coverUrl"
                  :src="track.coverUrl"
                  :alt="track.title"
                  class="h-full w-full object-cover"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <Music aria-hidden="true" class="text-xs text-white/30"  />
                </div>
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
                <p class="truncate text-xs text-white/40">{{ (track as any).artist_name || (track.artists?.[0]?.name) || 'Unknown artist' }}</p>
              </div>
              <span class="text-xs text-white/40 tabular-nums">{{ formatDuration(track.durationSeconds) }}</span>
            </button>
          </div>
        </section>

        <!-- ── Artists ── -->
        <section v-if="filteredResults.artists.length">
          <div class="mb-4 flex items-center gap-3">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-500/10">
              <Users aria-hidden="true" class="text-xs text-blue-400"  />
            </div>
            <h2 class="text-lg font-bold text-white">Artists</h2>
            <Badge :value="filteredResults.artists.length" severity="info" class="bg-blue-500/15! text-blue-400! text-[10px]! font-bold! min-w-5! h-5!" />
          </div>
          <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
            <RouterLink
              v-for="artist in filteredResults.artists"
              :key="artist.id"
              :to="`/artist/${artist.id}`"
              class="group spring rounded-2xl border border-white/6 bg-white/2 p-4 text-center transition-all hover:border-white/12 hover:bg-white/5 active:scale-[0.97]"
            >
              <div class="mx-auto h-20 w-20 overflow-hidden rounded-full bg-white/6 ring-1 ring-white/10 transition-all duration-300 group-hover:ring-spotify/50">
                <img
                  v-if="artist.imageUrl"
                  :src="artist.imageUrl"
                  :alt="artist.name"
                  class="h-full w-full object-cover"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <User aria-hidden="true" class="text-xl text-white/30"  />
                </div>
              </div>
              <p class="mt-3 truncate text-sm font-medium text-white">{{ artist.name }}</p>
              <p class="text-xs text-white/40">Artist</p>
            </RouterLink>
          </div>
        </section>

        <!-- ── Albums ── -->
        <section v-if="filteredResults.albums.length">
          <div class="mb-4 flex items-center gap-3">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-purple-500/10">
              <Book aria-hidden="true" class="text-xs text-purple-400"  />
            </div>
            <h2 class="text-lg font-bold text-white">Albums</h2>
            <Badge :value="filteredResults.albums.length" severity="warn" class="bg-purple-500/15! text-purple-400! text-[10px]! font-bold! min-w-5! h-5!" />
          </div>
          <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
            <RouterLink
              v-for="album in filteredResults.albums"
              :key="album.id"
              :to="`/album/${album.id}`"
              class="group active:scale-[0.98]"
            >
              <div class="spring relative aspect-square overflow-hidden rounded-2xl bg-white/4 ring-1 ring-white/10 transition-all group-hover:shadow-lg group-hover:ring-spotify/50">
                <img
                  v-if="album.coverUrl"
                  :src="album.coverUrl"
                  :alt="album.title"
                  class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <Disc3 aria-hidden="true" class="text-3xl text-white/30"  />
                </div>
                <div class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition group-hover:opacity-100">
                  <div class="flex h-12 w-12 items-center justify-center rounded-full bg-spotify text-black shadow-xl transition-transform group-hover:scale-110">
                    <Play aria-hidden="true" class="text-lg"  />
                  </div>
                </div>
              </div>
              <p class="mt-2 truncate text-sm font-medium text-white">{{ album.title }}</p>
              <p class="truncate text-xs text-white/40">{{ (album as any).artist_name || (album.artists?.[0]?.name) || '\u2014' }}</p>
            </RouterLink>
          </div>
        </section>
      </div>
    </template>

    <!-- ═══════════════════════ DISCOVER CONTENT ═══════════════════════ -->
    <template v-else>
      <!-- Loading skeleton for discover -->
      <div v-if="discoverLoading" class="space-y-8" aria-live="polite">
        <div v-for="s in 4" :key="s">
          <Skeleton width="30%" height="1.5rem" class="rounded! bg-white/6! mb-2!" />
          <div class="flex gap-4">
            <Skeleton v-for="i in 5" :key="i" shape="rect" size="10rem" class="rounded-2xl! bg-white/4! shrink-0" />
          </div>
        </div>
      </div>

      <template v-else>
        <!-- ── Trending Now ── -->
        <section class="mt-2">
          <HomeSectionHeader
            title="Trending Now"
            eyebrow="Popular"
            see-all-route="/recommendations/popular"
          />
          <HomeCarousel>
            <HomeTrackCard
              v-for="(track, i) in popular"
              :key="track.id"
              :item="track"
              :delay="i * 30"
              badge="Trending"
              @play="playTrack"
            />
          </HomeCarousel>
        </section>

        <!-- ── Made For You ── -->
        <section v-if="forYou.length" class="mt-12">
          <HomeSectionHeader
            title="Made For You"
            eyebrow="Personalized"
          />
          <HomeCarousel>
            <HomeTrackCard
              v-for="(track, i) in forYou"
              :key="track.id"
              :item="track"
              :delay="i * 30"
              @play="playTrack"
            />
          </HomeCarousel>
        </section>

        <!-- ── New Releases ── -->
        <section v-if="recent.length" class="mt-12">
          <HomeSectionHeader
            title="New Releases"
            eyebrow="Latest"
            see-all-route="/recommendations/recent"
          />
          <HomeCarousel>
            <HomeTrackCard
              v-for="(track, i) in recent"
              :key="track.id"
              :item="track"
              :delay="i * 30"
              @play="playTrack"
            />
          </HomeCarousel>
        </section>

        <!-- ── Recently Played ── -->
        <section v-if="recentlyPlayed.length" class="mt-12">
          <HomeSectionHeader
            title="Listen Again"
            eyebrow="Recently played"
            see-all-route="/recently-played"
          />
          <HomeCarousel>
            <HomeTrackCard
              v-for="(item, i) in recentlyPlayed"
              :key="item.track_id"
              :item="item"
              :delay="i * 30"
              @play="playHistoryItem"
            />
          </HomeCarousel>
        </section>

        <!-- ── Browse by Mood ── -->
        <section class="mt-12">
          <div class="mb-5 flex items-end justify-between gap-4">
            <div>
              <p class="text-[10px] font-bold tracking-[0.3em] text-white/40 uppercase">Feel something</p>
              <h2 class="text-2xl font-black text-white md:text-3xl">Browse by Mood</h2>
            </div>
            <Button
              text
              severity="secondary"
              size="small"
              as="router-link"
              to="/ai/mood-explorer"
              class="text-xs! px-3! py-1! text-white/30 hover:text-white"
              icon-pos="right"
            >
              Explore moods
              <template #icon>
                <ChevronLeft aria-hidden="true" class="text-[10px] ml-0.5"  />
              </template>
            </Button>
          </div>
          <div class="grid grid-cols-3 gap-3 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-8">
            <RouterLink
              v-for="mood in moods"
              :key="mood.value"
              :to="`/ai/mood-explorer?mood=${mood.value}`"
              class="group spring flex flex-col items-center gap-2 rounded-2xl border border-white/6 bg-white/2 px-3 py-4 text-center transition-all duration-300 hover:-translate-y-1 hover:border-white/15 hover:bg-white/6 active:scale-[0.95]"
            >
              <div
                class="flex h-10 w-10 items-center justify-center rounded-xl text-lg transition duration-300 group-hover:scale-110"
                :style="{ backgroundColor: mood.bg }"
              >
                <i aria-hidden="true" :class="mood.icon" :style="{ color: mood.fg }" />
              </div>
              <span class="max-w-full truncate text-[11px] font-bold text-white/60 group-hover:text-white/90">{{ mood.label }}</span>
            </RouterLink>
          </div>
        </section>

        <!-- ── Viral Hits ── -->
        <section v-if="popular.length" class="mt-12">
          <div class="mb-5 flex items-end justify-between gap-4">
            <div>
              <p class="text-[10px] font-bold tracking-[0.3em] text-white/40 uppercase">Trending fast</p>
              <h2 class="text-2xl font-black text-white md:text-3xl">Viral Hits</h2>
            </div>
          </div>
          <div class="grid gap-3 md:grid-cols-2">
            <button
              v-for="(track, idx) in popular.slice(0, 4)"
              :key="track.id"
              type="button"
              class="group spring flex w-full items-center gap-4 rounded-2xl border border-white/6 bg-white/2 px-4 py-3 transition-all hover:border-white/12 hover:bg-white/6 active:scale-[0.99]"
              @click="playTrack(track)"
            >
              <div class="flex w-8 items-center justify-center">
                <span class="text-lg font-black text-white/30 tabular-nums">{{ idx + 1 }}</span>
              </div>
              <div class="h-12 w-12 shrink-0 overflow-hidden rounded-xl bg-white/10">
                <img
                  v-if="track.cover_url"
                  :src="track.cover_url"
                  :alt="track.title"
                  loading="lazy"
                  class="h-full w-full object-cover"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <Music aria-hidden="true" class="text-white/30"  />
                </div>
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-bold text-white">{{ track.title }}</p>
                <p class="truncate text-xs text-white/40">{{ track.artist_name }}</p>
              </div>
              <Chip class="bg-spotify/15! text-spotify! text-[10px]! font-medium! px-3! py-1! rounded-full! border-0! h-auto! gap-1!">
                <span class="glow-spread inline-block h-1.5 w-1.5 rounded-full bg-spotify" />
                Trending
              </Chip>
            </button>
          </div>
        </section>

        <!-- ── Genre Worlds ── -->
        <section class="mt-12">
          <div class="mb-5 flex items-end justify-between gap-4">
            <div>
              <p class="text-[10px] font-bold tracking-[0.3em] text-white/40 uppercase">Explore</p>
              <h2 class="text-2xl font-black text-white md:text-3xl">Genre Worlds</h2>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-5">
            <RouterLink
              v-for="genre in genreWorlds"
              :key="genre.name"
              :to="`/search?genre=${encodeURIComponent(genre.name)}`"
              class="group spring relative flex h-28 items-end overflow-hidden rounded-2xl p-5 transition-all duration-300 hover:scale-[1.02] hover:shadow-lg active:scale-[1.01]"
              :style="{ background: genre.gradient }"
            >
              <div class="absolute top-3 right-3 text-2xl opacity-40 transition-all group-hover:scale-125 group-hover:opacity-70">
                {{ genre.icon }}
              </div>
              <p class="relative z-10 text-base font-bold text-white drop-shadow-xl">
                {{ genre.name }}
              </p>
            </RouterLink>
          </div>
        </section>

        <!-- ── From the Community ── -->
        <section v-if="feed.length" class="mt-12">
          <div class="mb-5 flex items-end justify-between gap-4">
            <div>
              <p class="text-[10px] font-bold tracking-[0.3em] text-white/40 uppercase">Activity</p>
              <h2 class="text-2xl font-black text-white md:text-3xl">From the Community</h2>
            </div>
          </div>
          <ActivityItem v-for="item in feed" :key="item.id" :item="item" />
        </section>
      </template>
    </template>

    <!-- ── Scroll to top ── -->
    <ScrollTop
      target="parent"
      :threshold="400"
      class="bottom-8! right-8! bg-spotify! text-black! shadow-lg! hover:bg-spotify-hover!"
      icon="pi pi-arrow-up"
    />
  </div>
  </div>
</template>

<script setup lang="ts">
import { Book, ChevronLeft, Disc3, History, Music, Play, PlayCircle, Search, User, Users, X } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSearchApi } from '@/services/api/catalog/search'
import type { SearchResult } from '@/services/api/catalog/search'
import { useSocialApi } from '@/services/api/social'
import { useRecommendationsApi } from '@/services/api/recommendation'
import { useHistoryApi } from '@/services/api/history'
import { usePlayer } from '@/composables/player'
import { buildPlaybackTrack } from '@/factories/playbackTrack'
import { onImgError } from '@/utils/helpers'
import { AppEmptyState } from '@/components/common'
import { HomeCarousel, HomeTrackCard, ActivityItem } from '@/components/music'
import HomeSectionHeader from '@/components/music/home/HomeSectionHeader.vue'
import { MOOD_OPTIONS } from '@/services/api/ai/types'
import { useRecentSearches } from '@/composables/search/useRecentSearches'
import { formatDuration } from '@/utils/format'
import type { RecommendationTrack } from '@/services/api/recommendation/types'
interface TrackCardItem {
  [key: string]: unknown
  id?: string | number
  cover_url?: string | null; coverUrl?: string | null; track_cover_url?: string | null
  title?: string | null; track_title?: string | null; artist_name?: string | null; artistName?: string | null
  album_title?: string | null; duration_seconds?: number | null
}
import type { HistoryItem } from '@/services/api/history/types'
import type { ActivityItemData } from '@/components/music/ActivityItem.vue'
import type { ActivityFeedItem } from '@/services/api/social/types'

// ── Composables ──
const route = useRoute()
const router = useRouter()
const searchApi = useSearchApi()
const socialApi = useSocialApi()
const recsApi = useRecommendationsApi()
const historyApi = useHistoryApi()
const player = usePlayer()
const { recentSearches, addRecentSearch, removeRecentSearch, clearRecentSearches }
  = useRecentSearches()

// ── Search state ──
const inputRef = ref<HTMLInputElement | null>(null)
const query = ref('')
const lastQuery = ref('')
const results = ref<SearchResult>({ tracks: [], artists: [], albums: [], playlists: [] })
const searching = ref(false)
const hasSearched = ref(false)
const hasNoResults = ref(false)
const searchError = ref(false)
type SearchTab = 'all' | 'tracks' | 'artists' | 'albums'
const searchType = ref<SearchTab>('all')
const showRecent = ref(false)
let debounceTimer: ReturnType<typeof setTimeout> | null = null
// eslint-disable-next-line @typescript-eslint/no-unused-vars
let hideRecentTimer: ReturnType<typeof setTimeout> | null = null

// ── Discover state ──
const popular = ref<RecommendationTrack[]>([])
const forYou = ref<RecommendationTrack[]>([])
const recent = ref<RecommendationTrack[]>([])
const feed = ref<ActivityItemData[]>([])
const recentlyPlayed = ref<HistoryItem[]>([])
const discoverLoading = ref(true)

// ── Computed ──
const pageTitle = computed(() =>
  hasSearched.value ? `\u201c${lastQuery.value}\u201d` : 'Find Your Sound',
)
const noResultsText = computed(() =>
  `No results for "${lastQuery.value}". Try a different search term.`,
)

const filteredResults = computed(() => {
  if (searchType.value === 'all') return results.value
  return {
    tracks: searchType.value === 'tracks' ? results.value.tracks : [],
    artists: searchType.value === 'artists' ? results.value.artists : [],
    albums: searchType.value === 'albums' ? results.value.albums : [],
  }
})

const tabs = [
  { value: 'all' as const, label: 'All' },
  { value: 'tracks' as const, label: 'Songs' },
  { value: 'artists' as const, label: 'Artists' },
  { value: 'albums' as const, label: 'Albums' },
]

function getCount(tab: string): number {
  if (tab === 'all') return results.value.tracks.length + results.value.artists.length + results.value.albums.length
  if (tab === 'tracks') return results.value.tracks.length
  if (tab === 'artists') return results.value.artists.length
  if (tab === 'albums') return results.value.albums.length
  return 0
}

// ── Moods ──
const moods = MOOD_OPTIONS.slice(0, 8).map((m) => {
  const colors: Record<string, { bg: string; fg: string }> = {
    energetic: { bg: 'rgba(34,197,94,0.12)', fg: '#22c55e' },
    happy: { bg: 'rgba(250,204,21,0.12)', fg: '#facc15' },
    chill: { bg: 'rgba(96,165,250,0.12)', fg: '#60a5fa' },
    calm: { bg: 'rgba(148,163,184,0.12)', fg: '#94a3b8' },
    sad: { bg: 'rgba(148,163,184,0.12)', fg: '#94a3b8' },
    focus: { bg: 'rgba(168,85,247,0.12)', fg: '#a855f7' },
    romantic: { bg: 'rgba(244,114,182,0.12)', fg: '#f472b6' },
    intense: { bg: 'rgba(239,68,68,0.12)', fg: '#ef4444' },
  }
  const c = colors[m.value] ?? { bg: 'rgba(255,255,255,0.06)', fg: '#fff' }
  return { ...m, bg: c.bg, fg: c.fg }
})

// ── Genre Worlds ──
const genreWorlds = [
  { name: 'Pop', icon: '\u{1F31F}', gradient: 'linear-gradient(135deg, #831843, #9d174d)' },
  { name: 'Rock', icon: '\u{1F3B8}', gradient: 'linear-gradient(135deg, #7f1d1d, #991b1b)' },
  { name: 'Hip-Hop', icon: '\u{1F3A4}', gradient: 'linear-gradient(135deg, #713f12, #854d0e)' },
  { name: 'Electronic', icon: '\u{1F3B9}', gradient: 'linear-gradient(135deg, #0c4a6e, #075985)' },
  { name: 'Jazz', icon: '\u{1F3B7}', gradient: 'linear-gradient(135deg, #1e1b4b, #312e81)' },
  { name: 'Classical', icon: '\u{1F3BB}', gradient: 'linear-gradient(135deg, #3b0764, #581c87)' },
  { name: 'R&B', icon: '\u{1F399}\uFE0F', gradient: 'linear-gradient(135deg, #831843, #9d174d)' },
  { name: 'Folk', icon: '\u{1FA95}', gradient: 'linear-gradient(135deg, #422006, #713f12)' },
  { name: 'Ambient', icon: '\u{1F30C}', gradient: 'linear-gradient(135deg, #0f172a, #1e293b)' },
  { name: 'Traditional', icon: '\u{1F3EE}', gradient: 'linear-gradient(135deg, #78350f, #92400e)' },
]

// ── Search ──
function onInput() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    doSearch()
  }, 350)
}

async function doSearch() {
  const q = query.value.trim()
  if (!q) return
  searchType.value = 'all'

  lastQuery.value = q
  searching.value = true
  hasSearched.value = true
  hasNoResults.value = false
  searchError.value = false
  showRecent.value = false
  addRecentSearch(q)

  // Sync query to URL param (debounced)
  router.replace({ query: { ...route.query, q } })

  try {
    const res = await searchApi.searchCatalog({ query: q, limit: 20 })
    results.value = res
    const empty = !res.tracks.length && !res.artists.length && !res.albums.length
    hasNoResults.value = empty
    searchError.value = false
  } catch (err) {
    console.error('Failed to search catalog:', err)
    searchError.value = true
    hasNoResults.value = false
    results.value = { tracks: [], artists: [], albums: [], playlists: [] }
  } finally {
    searching.value = false
  }
}

function clearSearch() {
  query.value = ''
  results.value = { tracks: [], artists: [], albums: [], playlists: [] }
  hasSearched.value = false
  hasNoResults.value = false
  searchType.value = 'all'
  inputRef.value?.focus()
}

function selectRecent(q: string) {
  query.value = q
  showRecent.value = false
  doSearch()
}

function hideRecentDelayed() {
  hideRecentTimer = setTimeout(() => {
    showRecent.value = false
  }, 200)
}

// ── Discover ──
async function fetchDiscover() {
  discoverLoading.value = true
  try {
    const [popularData, forYouData, recentData, feedData, historyData] = await Promise.all([
      recsApi.getPopular({ limit: 10 }).catch(() => null),
      recsApi.getForYou({ limit: 10 }).catch(() => null),
      recsApi.getRecent({ limit: 10 }).catch(() => null),
      socialApi.getFeed({ limit: 20, types: 'upload' }).catch(() => null),
      historyApi.getHistory({ limit: 10 }).catch(() => null),
    ])
    if (popularData?.items) popular.value = popularData.items
    if (forYouData?.items) forYou.value = forYouData.items
    if (recentData?.items) recent.value = recentData.items
    if (feedData?.items) {
      feed.value = feedData.items.map((item: ActivityFeedItem) => ({
        id: item.id,
        userId: item.user_id,
        userName: item.user_display_name,
        avatarUrl: item.user_avatar_url,
        action: item.type,
        targetName: item.target_name,
        targetUrl: item.target_id ? `/${item.target_type}/${item.target_id}` : null,
        createdAt: item.created_at,
      }))
    }
    if (historyData?.items) recentlyPlayed.value = historyData.items
  } catch (err) {
    console.error('Failed to fetch discover data:', err)
  } finally {
    discoverLoading.value = false
  }
}

function playTrack(track: TrackCardItem) {
  void player.playTrack(
    buildPlaybackTrack({
      id: String(track.id),
      title: (track.title as string) || undefined,
      artist_name: (track.artist_name as string) || 'Unknown',
      album_title: (track.album_title as string) || null,
      cover_url: (track.cover_url as string) || null,
      duration_seconds: (track.duration_seconds as number) ?? null,
    }),
  )
}

function playHistoryItem(item: TrackCardItem) {
  void player.playTrack(
    buildPlaybackTrack({
      id: String(item.track_id),
      title: (item.track_title as string) || 'Unknown',
      artist_name: (item.artist_name as string) || 'Unknown',
      cover_url: (item.track_cover_url as string) || null,
      duration_seconds: (item.track_duration as number) ?? null,
    }),
  )
}

function playAllSearchResults() {
  const tracks = filteredResults.value.tracks
  if (!tracks.length) return
  const queue = tracks.map((t: any) =>
    buildPlaybackTrack({
      id: String(t.id),
      title: t.title || undefined,
      artist_name: t.artist_name || 'Unknown',
      album_title: t.album_title || null,
      cover_url: t.cover_url || null,
      duration_seconds: t.duration_seconds ?? null,
    })
  )
  void player.setQueueAndPlay(queue, 0)
}

// Restore query from URL on mount
const initialQuery = route.query.q as string | undefined
if (initialQuery) {
  query.value = initialQuery
  doSearch()
}

// Sync query to URL with debounce
watch(query, () => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    router.replace({ query: { q: query.value || undefined } })
  }, 500)
})

onMounted(fetchDiscover)
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
