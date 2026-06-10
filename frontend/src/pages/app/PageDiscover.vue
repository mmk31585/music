<template>
  <div class="pb-36">
    <!-- Aurora hero -->
    <div class="relative overflow-hidden">
      <div class="aurora-flow absolute inset-0 opacity-40" />
      <div class="aurora-spot-1 -top-60 -left-40 bg-[#1db954]/20" style="animation-duration: 18s" />
      <div class="aurora-spot-2 -top-40 right-0 bg-[#60a5fa]/15" style="animation-duration: 22s" />
      <div
        class="absolute right-0 bottom-0 left-0 h-32 bg-gradient-to-t from-[#050505] to-transparent"
      />
      <div class="relative mx-auto w-full max-w-7xl px-4 pt-20 pb-16 md:px-6 lg:px-8">
        <p class="text-xs font-bold tracking-[0.25em] text-[#1db954] uppercase">Discover</p>
        <h1 class="mt-2 text-4xl font-black text-white md:text-6xl">
          Find Your Next<br class="sm:hidden" />
          <span class="text-gradient">Favorite Track</span>
        </h1>
        <p class="mt-3 max-w-lg text-base leading-relaxed text-white/40">
          Explore trending tracks, viral hits, new releases, and personalized picks — all in one
          place.
        </p>
      </div>
    </div>

    <div class="mx-auto w-full max-w-7xl px-4 md:px-6 lg:px-8">
      <div v-if="loading" class="mt-6 space-y-8">
        <div v-for="s in 4" :key="s">
          <SkeletonLoader variant="lines" :lines="1" class="w-40" />
          <div class="mt-4 flex gap-4">
            <SkeletonLoader v-for="i in 5" :key="i" variant="card" class="w-40 shrink-0" />
          </div>
        </div>
      </div>

      <template v-else>
        <!-- Trending -->
        <section class="mt-8">
          <template v-if="popular.length">
          <div class="mb-4 flex items-end justify-between">
            <div>
              <p class="text-xs font-bold tracking-[0.2em] text-[#1db954] uppercase">Popular</p>
              <h2 class="mt-1 text-2xl font-black text-white md:text-3xl">Trending Now</h2>
            </div>
            <RouterLink
              to="/recommendations/popular"
              class="spring text-xs font-bold text-slate-400 transition hover:text-white"
              >See all</RouterLink
            >
          </div>
          <HomeCarousel>
            <div
              v-for="track in popular"
              :key="track.id"
              class="group w-44 shrink-0 cursor-pointer"
              @click="playTrack(track)"
            >
              <div
                class="spring relative aspect-square overflow-hidden rounded-2xl bg-white/[0.04] ring-1 ring-white/10 transition-all group-hover:shadow-[0_8px_32px_rgba(0,0,0,0.4)] group-hover:ring-[#1db954]/50"
              >
                <img
                  v-if="track.cover_url"
                  :src="track.cover_url"
                  :alt="track.title"
                  loading="lazy"
                  class="h-full w-full object-cover transition duration-500 group-hover:scale-110"
                  @error="onImgError"
                />
                <div v-else class="flex h-full w-full items-center justify-center text-white/20">
                  <i class="pi pi-music text-3xl" />
                </div>
                <div class="absolute top-2 left-2">
                  <span
                    class="rounded-full bg-[#1db954]/20 px-2 py-0.5 text-[9px] font-bold text-[#1db954] backdrop-blur-sm"
                    >Trending</span
                  >
                </div>
                <div
                  class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition group-hover:opacity-100"
                >
                  <div
                    class="flex h-12 w-12 items-center justify-center rounded-full bg-[#1db954] text-black shadow-xl transition-transform group-hover:scale-110"
                  >
                    <i class="pi pi-play-fill text-lg" />
                  </div>
                </div>
              </div>
              <p class="mt-2 truncate text-sm font-medium text-white">{{ track.title }}</p>
              <p class="truncate text-xs text-white/40">{{ track.artist_name || 'Unknown' }}</p>
            </div>
          </HomeCarousel>
        </template>
        <template v-else>
          <div class="flex flex-col items-center gap-3 py-16 text-center">
            <i class="pi pi-inbox text-4xl text-slate-500" />
            <p class="text-sm text-slate-400">No trending tracks</p>
          </div>
        </template>
      </section>

        <!-- Recently Played -->
        <section class="mt-10">
          <template v-if="recentlyPlayed.length">
          <div class="mb-4 flex items-end justify-between">
            <div>
              <p class="text-xs font-bold tracking-[0.2em] text-[#1db954] uppercase">
                Recently played
              </p>
              <h2 class="mt-1 text-2xl font-black text-white md:text-3xl">Listen Again</h2>
            </div>
            <RouterLink
              to="/recently-played"
              class="spring text-xs font-bold text-slate-400 transition hover:text-white"
              >See all</RouterLink
            >
          </div>
          <HomeCarousel>
            <div
              v-for="item in recentlyPlayed"
              :key="item.track_id"
              class="group w-44 shrink-0 cursor-pointer"
              @click="playHistoryItem(item)"
            >
              <div
                class="spring relative aspect-square overflow-hidden rounded-2xl bg-white/[0.04] ring-1 ring-white/10 transition-all group-hover:ring-[#1db954]/50"
              >
                <img
                  v-if="item.track_cover_url"
                  :src="item.track_cover_url"
                  :alt="item.track_title"
                  loading="lazy"
                  class="h-full w-full object-cover transition duration-500 group-hover:scale-110"
                  @error="onImgError"
                />
                <div v-else class="flex h-full w-full items-center justify-center text-white/20">
                  <i class="pi pi-music text-3xl" />
                </div>
                <div
                  class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition group-hover:opacity-100"
                >
                  <div
                    class="flex h-12 w-12 items-center justify-center rounded-full bg-[#1db954] text-black shadow-xl transition-transform group-hover:scale-110"
                  >
                    <i class="pi pi-play-fill text-lg" />
                  </div>
                </div>
              </div>
              <p class="mt-2 truncate text-sm font-medium text-white">{{ item.track_title }}</p>
              <p class="truncate text-xs text-white/40">{{ item.artist_name || 'Unknown' }}</p>
            </div>
          </HomeCarousel>
        </template>
        <template v-else>
          <div class="flex flex-col items-center gap-3 py-16 text-center">
            <i class="pi pi-inbox text-4xl text-slate-500" />
            <p class="text-sm text-slate-400">No recently played tracks</p>
          </div>
        </template>
      </section>

        <!-- Mood Grid -->
        <section class="mt-10">
          <div class="mb-4 flex items-end justify-between">
            <div>
              <p class="text-xs font-bold tracking-[0.2em] text-[#1db954] uppercase">
                Feel something
              </p>
              <h2 class="mt-1 text-2xl font-black text-white md:text-3xl">Browse by Mood</h2>
            </div>
            <RouterLink
              to="/ai/mood-explorer"
              class="spring text-xs font-bold text-slate-400 transition hover:text-white"
              >Explore moods</RouterLink
            >
          </div>
          <div class="grid grid-cols-3 gap-3 sm:grid-cols-4 md:grid-cols-5 lg:grid-cols-8">
            <RouterLink
              v-for="mood in moods"
              :key="mood.value"
              :to="`/ai/mood-explorer?mood=${mood.value}`"
              class="group spring flex flex-col items-center gap-2 rounded-2xl border border-white/[0.06] bg-white/[0.02] px-3 py-4 text-center transition-all duration-300 hover:-translate-y-1 hover:border-white/[0.15] hover:bg-white/[0.06]"
            >
              <div
                class="flex h-10 w-10 items-center justify-center rounded-xl text-lg transition duration-300 group-hover:scale-110"
                :style="{ backgroundColor: mood.bg }"
              >
                <i :class="mood.icon" :style="{ color: mood.fg }" />
              </div>
              <span
                class="max-w-full truncate text-[11px] font-bold text-white/60 group-hover:text-white/90"
                >{{ mood.label }}</span
              >
            </RouterLink>
          </div>
        </section>

        <!-- For You -->
        <section class="mt-10">
          <template v-if="forYou.length">
          <div class="mb-4 flex items-end justify-between">
            <div>
              <p class="text-xs font-bold tracking-[0.2em] text-[#1db954] uppercase">
                Personalized
              </p>
              <h2 class="mt-1 text-2xl font-black text-white md:text-3xl">Made For You</h2>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
            <div
              v-for="track in forYou"
              :key="track.id"
              class="group cursor-pointer"
              @click="playTrack(track)"
            >
              <div
                class="spring relative aspect-square overflow-hidden rounded-2xl bg-white/[0.04] ring-1 ring-white/10 transition-all group-hover:shadow-[0_8px_32px_rgba(0,0,0,0.4)] group-hover:ring-[#1db954]/50"
              >
                <img
                  v-if="track.cover_url"
                  :src="track.cover_url"
                  :alt="track.title"
                  loading="lazy"
                  class="h-full w-full object-cover transition duration-500 group-hover:scale-110"
                  @error="onImgError"
                />
                <div v-else class="flex h-full w-full items-center justify-center text-white/20">
                  <i class="pi pi-music text-3xl" />
                </div>
                <div
                  class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition group-hover:opacity-100"
                >
                  <div
                    class="flex h-12 w-12 items-center justify-center rounded-full bg-[#1db954] text-black shadow-xl transition-transform group-hover:scale-110"
                  >
                    <i class="pi pi-play-fill text-lg" />
                  </div>
                </div>
              </div>
              <p class="mt-2 truncate text-sm font-medium text-white">{{ track.title }}</p>
              <p class="truncate text-xs text-white/40">{{ track.artist_name || 'Unknown' }}</p>
            </div>
          </div>
        </template>
        <div v-else class="flex flex-col items-center gap-3 py-16 text-center">
          <i class="pi pi-inbox text-4xl text-slate-500" />
          <p class="text-sm text-slate-400">No personalized recommendations yet</p>
        </div>
      </section>

        <!-- New Releases -->
        <section class="mt-10">
          <template v-if="recent.length">
          <div class="mb-4 flex items-end justify-between">
            <div>
              <p class="text-xs font-bold tracking-[0.2em] text-[#1db954] uppercase">Latest</p>
              <h2 class="mt-1 text-2xl font-black text-white md:text-3xl">New Releases</h2>
            </div>
            <RouterLink
              to="/recommendations/recent"
              class="spring text-xs font-bold text-slate-400 transition hover:text-white"
              >See all</RouterLink
            >
          </div>
          <HomeCarousel>
            <div
              v-for="track in recent"
              :key="track.id"
              class="group w-44 shrink-0 cursor-pointer"
              @click="playTrack(track)"
            >
              <div
                class="spring relative aspect-square overflow-hidden rounded-2xl bg-white/[0.04] ring-1 ring-white/10 transition-all group-hover:ring-[#1db954]/50"
              >
                <img
                  v-if="track.cover_url"
                  :src="track.cover_url"
                  :alt="track.title"
                  loading="lazy"
                  class="h-full w-full object-cover transition duration-500 group-hover:scale-110"
                  @error="onImgError"
                />
                <div v-else class="flex h-full w-full items-center justify-center text-white/20">
                  <i class="pi pi-music text-3xl" />
                </div>
                <div class="absolute top-2 left-2">
                  <span
                    class="rounded-full bg-[#1db954]/20 px-2 py-0.5 text-[9px] font-bold text-[#1db954] backdrop-blur-sm"
                    >NEW</span
                  >
                </div>
                <div
                  class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition group-hover:opacity-100"
                >
                  <div
                    class="flex h-12 w-12 items-center justify-center rounded-full bg-[#1db954] text-black shadow-xl transition-transform group-hover:scale-110"
                  >
                    <i class="pi pi-play-fill text-lg" />
                  </div>
                </div>
              </div>
              <p class="mt-2 truncate text-sm font-medium text-white">{{ track.title }}</p>
              <p class="truncate text-xs text-white/40">{{ track.artist_name || 'Unknown' }}</p>
            </div>
          </HomeCarousel>
        </template>
        <div v-else class="flex flex-col items-center gap-3 py-16 text-center">
          <i class="pi pi-inbox text-4xl text-slate-500" />
          <p class="text-sm text-slate-400">No new releases right now</p>
        </div>
      </section>

        <!-- Viral Hits -->
        <section class="mt-10">
          <template v-if="popular.length">
          <div class="mb-4 flex items-end justify-between">
            <div>
              <p class="text-xs font-bold tracking-[0.2em] text-[#1db954] uppercase">
                Trending fast
              </p>
              <h2 class="mt-1 text-2xl font-black text-white md:text-3xl">Viral Hits</h2>
            </div>
          </div>
          <div class="grid gap-3 md:grid-cols-2">
            <div
              v-for="(track, idx) in popular.slice(0, 4)"
              :key="track.id"
              class="group spring flex cursor-pointer items-center gap-4 rounded-2xl border border-white/[0.06] bg-white/[0.02] px-4 py-3 transition-all hover:border-white/[0.12] hover:bg-white/[0.06]"
              @click="playTrack(track)"
            >
              <div class="flex w-8 items-center justify-center">
                <span class="text-lg font-black text-white/30">{{ idx + 1 }}</span>
              </div>
              <div
                class="h-12 w-12 shrink-0 overflow-hidden rounded-xl bg-white/10 ring-1 ring-white/10"
              >
                <img
                  v-if="track.cover_url"
                  :src="track.cover_url"
                  :alt="track.title"
                  loading="lazy"
                  class="h-full w-full object-cover"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i class="pi pi-music text-white/30" />
                </div>
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-bold text-white">{{ track.title }}</p>
                <p class="truncate text-xs text-white/40">{{ track.artist_name }}</p>
              </div>
              <div class="flex items-center gap-1.5 text-[10px] font-medium text-[#1db954]">
                <span class="glow-spread inline-block h-1.5 w-1.5 rounded-full bg-[#1db954]" />
                Trending
              </div>
            </div>
          </div>
        </template>
        <div v-else class="flex flex-col items-center gap-3 py-16 text-center">
          <i class="pi pi-inbox text-4xl text-slate-500" />
          <p class="text-sm text-slate-400">No viral hits right now</p>
        </div>
      </section>

        <!-- Genre Worlds -->
        <section class="mt-10">
          <div class="mb-4 flex items-end justify-between">
            <div>
              <p class="text-xs font-bold tracking-[0.2em] text-[#1db954] uppercase">Explore</p>
              <h2 class="mt-1 text-2xl font-black text-white md:text-3xl">Genre Worlds</h2>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-5">
            <RouterLink
              v-for="genre in genreWorlds"
              :key="genre.name"
              :to="`/search?genre=${encodeURIComponent(genre.name)}`"
              class="group spring relative flex h-28 items-end overflow-hidden rounded-2xl p-5 transition-all duration-300 hover:scale-[1.02] hover:shadow-[0_8px_32px_rgba(0,0,0,0.4)]"
              :style="{ background: genre.gradient }"
            >
              <div
                class="absolute top-3 right-3 text-2xl opacity-40 transition-all group-hover:scale-125 group-hover:opacity-70"
              >
                {{ genre.icon }}
              </div>
              <p class="relative z-10 text-base font-bold text-white drop-shadow-xl">
                {{ genre.name }}
              </p>
            </RouterLink>
          </div>
        </section>

        <!-- Community Feed -->
        <section class="mt-10">
          <template v-if="feed.length">
          <div class="mb-4 flex items-end justify-between">
            <div>
              <p class="text-xs font-bold tracking-[0.2em] text-[#1db954] uppercase">Activity</p>
              <h2 class="mt-1 text-2xl font-black text-white md:text-3xl">From the Community</h2>
            </div>
          </div>
          <div class="space-y-2">
            <ActivityItem v-for="item in feed" :key="item.id" :item="item" />
          </div>
        </template>
        <div v-else class="flex flex-col items-center gap-3 py-16 text-center">
          <i class="pi pi-inbox text-4xl text-slate-500" />
          <p class="text-sm text-slate-400">No community activity yet</p>
        </div>
      </section>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { SkeletonLoader } from '@/components/common'
import { useSocialApi } from '@/services/api/social'
import { useRecommendationsApi } from '@/services/api/recommendation'
import { useHistoryApi } from '@/services/api/history'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'
import { MOOD_OPTIONS } from '@/services/api/ai/types'
import { onImgError } from '@/utils/helpers'
import { HomeSection, HomeCarousel, ActivityItem } from '@/components/music'

const socialApi = useSocialApi()
const recsApi = useRecommendationsApi()
const historyApi = useHistoryApi()
const player = usePlayer()
const playerApi = usePlayerApi()

const popular = ref<any[]>([])
const forYou = ref<any[]>([])
const recent = ref<any[]>([])
const feed = ref<any[]>([])
const recentlyPlayed = ref<any[]>([])
const loading = ref(true)

const moods: any[] = (MOOD_OPTIONS as unknown as any[]).slice(0, 8).map((m: any) => {
  const colors: Record<string, { bg: string; fg: string }> = {
    energetic: { bg: 'rgba(34,197,94,0.12)', fg: '#22c55e' },
    happy: { bg: 'rgba(250,204,21,0.12)', fg: '#facc15' },
    calm: { bg: 'rgba(96,165,250,0.12)', fg: '#60a5fa' },
    sad: { bg: 'rgba(148,163,184,0.12)', fg: '#94a3b8' },
    focused: { bg: 'rgba(168,85,247,0.12)', fg: '#a855f7' },
    romantic: { bg: 'rgba(244,114,182,0.12)', fg: '#f472b6' },
    dark: { bg: 'rgba(100,116,139,0.12)', fg: '#64748b' },
    party: { bg: 'rgba(239,68,68,0.12)', fg: '#ef4444' },
  }
  const c = colors[m.value] || { bg: 'rgba(255,255,255,0.06)', fg: '#fff' }
  return { ...m, bg: c.bg, fg: c.fg }
})

async function fetchDiscover() {
  loading.value = true
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
    if (feedData?.items) feed.value = feedData.items
    if (historyData?.items) recentlyPlayed.value = historyData.items
  } catch {
    /* silent */
  } finally {
    loading.value = false
  }
}

function playTrack(track: any) {
  player.setQueueAndPlay(
    [
      {
        id: String(track.id),
        title: track.title,
        artistName: track.artist_name || 'Unknown',
        albumTitle: track.album_title || null,
        coverUrl: track.cover_url || null,
        durationSeconds: track.duration_seconds ?? null,
        streamUrl: playerApi.getTrackStreamUrl(String(track.id)),
      },
    ],
    0,
  )
}

function playHistoryItem(item: any) {
  player.setQueueAndPlay(
    [
      {
        id: String(item.track_id),
        title: item.track_title || 'Unknown',
        artistName: item.artist_name || 'Unknown',
        albumTitle: item.album_title || null,
        coverUrl: item.track_cover_url || null,
        durationSeconds: item.track_duration ?? null,
        streamUrl: playerApi.getTrackStreamUrl(String(item.track_id)),
      },
    ],
    0,
  )
}

const genreWorlds = [
  { name: 'Pop', icon: '🌟', gradient: 'linear-gradient(135deg, #831843, #9d174d)' },
  { name: 'Rock', icon: '🎸', gradient: 'linear-gradient(135deg, #7f1d1d, #991b1b)' },
  { name: 'Hip-Hop', icon: '🎤', gradient: 'linear-gradient(135deg, #713f12, #854d0e)' },
  { name: 'Electronic', icon: '🎹', gradient: 'linear-gradient(135deg, #0c4a6e, #075985)' },
  { name: 'Jazz', icon: '🎷', gradient: 'linear-gradient(135deg, #1e1b4b, #312e81)' },
  { name: 'Classical', icon: '🎻', gradient: 'linear-gradient(135deg, #3b0764, #581c87)' },
  { name: 'R&B', icon: '🎙️', gradient: 'linear-gradient(135deg, #831843, #9d174d)' },
  { name: 'Folk', icon: '🪕', gradient: 'linear-gradient(135deg, #422006, #713f12)' },
  { name: 'Ambient', icon: '🌌', gradient: 'linear-gradient(135deg, #0f172a, #1e293b)' },
  { name: 'Traditional', icon: '🏮', gradient: 'linear-gradient(135deg, #78350f, #92400e)' },
]

onMounted(fetchDiscover)
</script>
