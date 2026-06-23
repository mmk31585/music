<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <!-- Discover Weekly Header -->
    <section
      class="rounded-[2rem] bg-gradient-to-br from-pink-600 via-purple-700 to-[#101010] p-8 text-white md:p-12"
    >
      <p class="text-sm font-bold tracking-[0.35em] text-white/70 uppercase">Discover Weekly</p>
      <h1 class="mt-3 text-4xl font-black md:text-6xl">کشف هفتگی شما</h1>
      <p class="mt-2 text-sm text-white/50">
        {{ weekLabel }}
      </p>
      <p class="mt-4 max-w-2xl text-white/80">
        ۳۰ آهنگ که فکر می‌کنیم دوستشون داری، بر اساس سلیقه‌ات. هر هفته تازه می‌شه.
      </p>

      <button
        v-if="tracks.length > 0"
        class="mt-6 inline-flex items-center gap-2 rounded-full bg-[#1db954] px-8 py-3 text-sm font-bold text-black transition hover:scale-105 hover:bg-[#1ed760]"
        @click="playAll"
      >
        <i aria-hidden="true" class="pi pi-play-fill" />
        پخش همه
      </button>
    </section>

    <!-- Loading -->
    <div
      v-if="loading"
      class="mt-10 grid gap-5 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5"
    >
      <div v-for="i in 15" :key="i" class="space-y-3">
        <div class="aspect-square animate-pulse rounded-2xl bg-white/[0.06]" />
        <div class="h-4 w-3/4 animate-pulse rounded bg-white/[0.06]" />
        <div class="h-3 w-1/2 animate-pulse rounded bg-white/[0.06]" />
      </div>
    </div>

    <!-- Empty / no data -->
    <div
      v-else-if="tracks.length === 0 && !error"
      class="mt-10 flex flex-col items-center gap-4 rounded-3xl border border-white/10 px-6 py-20 text-center"
    >
      <div class="flex h-16 w-16 items-center justify-center rounded-full bg-white/10">
        <i aria-hidden="true" class="pi pi-refresh text-2xl text-slate-400" />
      </div>
      <h3 class="text-xl font-bold text-white">Not enough data yet</h3>
      <p class="max-w-sm text-sm text-slate-400">
        Listen to more tracks and follow artists to get your weekly discovery playlist.
      </p>
      <RouterLink
        to="/recommendations/popular"
        class="rounded-full bg-[#1db954] px-5 py-2 text-sm font-bold text-black transition hover:bg-[#1ed760]"
      >
        Explore popular tracks
      </RouterLink>
    </div>

    <!-- Error state -->
    <div
      v-else-if="error"
      class="mt-10 flex flex-col items-center gap-4 rounded-3xl border border-white/10 px-6 py-20 text-center"
    >
      <div class="flex h-16 w-16 items-center justify-center rounded-full bg-white/10">
        <i aria-hidden="true" class="pi pi-exclamation-triangle text-2xl text-red-400" />
      </div>
      <h3 class="text-xl font-bold text-white">Could not load Discover Weekly</h3>
      <p class="text-sm text-slate-400">{{ error }}</p>
      <button
        class="rounded-full bg-white/10 px-5 py-2 text-sm font-bold text-white transition hover:bg-white/20"
        @click="fetchDiscoverWeekly"
      >
        Try again
      </button>
    </div>

    <!-- Track Grid -->
    <div v-else class="mt-10" aria-live="polite">
      <div class="mb-6 flex items-center gap-3 text-white/40">
        <i aria-hidden="true" class="pi pi-info-circle text-sm" />
        <span class="text-xs">Hover a track to see why it was recommended</span>
      </div>

      <div class="grid gap-5 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
        <div
          v-for="(track, idx) in tracks"
          :key="track.id"
          role="button"
          tabindex="0"
          class="group relative cursor-pointer"
          @click="play(track, idx)"
          @keydown.enter="play(track, idx)"
          @keydown.space.prevent="play(track, idx)"
        >
          <div
            class="relative mb-3 aspect-square overflow-hidden rounded-2xl bg-white/10 shadow-lg ring-1 ring-white/10 transition group-hover:ring-pink-500/50"
          >
            <img
              v-if="track.cover_url"
              :src="track.cover_url"
              :alt="track.title"
              loading="lazy"
              class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
              @error="onImgError"
            />
            <div v-else class="flex h-full items-center justify-center">
              <i aria-hidden="true" class="pi pi-music text-2xl text-slate-500" />
            </div>
            <div
              class="absolute inset-0 flex items-center justify-center bg-black/30 opacity-0 transition group-hover:opacity-100"
            >
              <div
                class="flex h-12 w-12 items-center justify-center rounded-full bg-pink-500/90 text-white shadow-xl"
              >
                <i aria-hidden="true" class="pi pi-play-fill text-lg" />
              </div>
            </div>
            <div
              class="absolute top-2 right-2 rounded-full bg-black/60 px-2 py-0.5 text-[10px] font-bold text-pink-400 opacity-0 backdrop-blur transition group-hover:opacity-100"
              :title="track.genre ? `چون شبیه آهنگ‌های ${track.genre} است` : undefined"
            >
              <i aria-hidden="true" class="pi pi-question-circle mr-1" />
              چرا این آهنگ؟
            </div>
          </div>
          <p class="truncate text-sm font-semibold text-white">{{ track.title }}</p>
          <p class="truncate text-xs text-slate-500">{{ track.artist_name || 'Unknown' }}</p>
          <p class="mt-0.5 truncate text-[10px] text-slate-600">
            {{ track.genre ? `چون شبیه آهنگ‌های ${track.genre} است` : 'توصیه هوشمند' }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRecommendationsApi } from '@/services/api/recommendation'
import type { RecommendationTrack } from '@/services/api/recommendation'
import type { DiscoverWeeklyPlaylistMeta } from '@/services/api/recommendation'
import { usePlayer } from '@/composables/player'
import { onImgError } from '@/utils/helpers'
import { usePlayerApi } from '@/services/api/player'

const api = useRecommendationsApi()
const player = usePlayer()
const playerApi = usePlayerApi()

const loading = ref(true)
const error = ref<string | null>(null)
const tracks = ref<RecommendationTrack[]>([])
const playlistMeta = ref<DiscoverWeeklyPlaylistMeta | null>(null)

const weekLabel = computed(() => {
  if (!playlistMeta.value?.week_of) return ''
  const d = new Date(playlistMeta.value.week_of)
  const end = new Date(d)
  end.setDate(end.getDate() + 6)
  const fmt = (date: Date) =>
    date.toLocaleDateString('fa-IR', { month: 'long', day: 'numeric' })
  return `${fmt(d)} — ${fmt(end)}`
})

async function fetchDiscoverWeekly() {
  loading.value = true
  error.value = null
  try {
    const response = await api.getDiscoverWeekly()
    tracks.value = response.tracks ?? []
    playlistMeta.value = response.playlist
  } catch (err) {
    const msg = err instanceof Error ? err.message : 'Failed to load Discover Weekly'
    error.value = msg
  } finally {
    loading.value = false
  }
}

function play(track: RecommendationTrack, index: number) {
  const allTracks = tracks.value.map((t) => ({
    id: String(t.id),
    title: t.title,
    artistName: t.artist_name || 'Unknown',
    albumTitle: t.album_title || null,
    coverUrl: t.cover_url || null,
    durationSeconds: t.duration_seconds ?? null,
    streamUrl: playerApi.getTrackStreamUrl(String(t.id)),
  }))
  player.setQueueAndPlay(allTracks, index)
}

function playAll() {
  const first = tracks.value[0]
  if (first) {
    play(first, 0)
  }
}

onMounted(fetchDiscoverWeekly)
</script>
