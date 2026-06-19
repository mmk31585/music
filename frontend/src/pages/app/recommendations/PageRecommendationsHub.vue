<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <!-- Hero -->
    <section
      class="rounded-[2rem] bg-gradient-to-br from-[#1db954] via-[#137d39] to-[#101010] p-8 text-white md:p-12"
    >
      <p class="text-sm font-bold tracking-[0.35em] text-white/70 uppercase">Made for you</p>
      <h1 class="mt-3 text-4xl font-black md:text-6xl">Recommendations</h1>
      <p class="mt-4 max-w-2xl text-white/80">
        Discover popular, recent, and personalized tracks based on your listening experience.
      </p>
    </section>

    <!-- Category Cards -->
    <section class="mt-10 grid gap-5 md:grid-cols-4">
      <RouterLink
        v-for="card in cards"
        :key="card.to"
        :to="card.to"
        class="group rounded-3xl border border-white/10 bg-white/[0.05] p-6 transition hover:-translate-y-1 hover:bg-white/[0.08]"
      >
        <div
          class="flex h-14 w-14 items-center justify-center rounded-2xl text-2xl"
          :class="card.iconClass"
        >
          <i aria-hidden="true" :class="card.icon" />
        </div>
        <h2 class="mt-5 text-xl font-black text-white">{{ card.title }}</h2>
        <p class="mt-2 text-sm leading-6 text-slate-400">{{ card.description }}</p>
      </RouterLink>
    </section>

    <!-- Popular Tracks -->
    <section class="mt-12">
      <div class="mb-5 flex items-end justify-between">
        <div>
          <p class="text-sm font-semibold text-[#1db954]">Trending</p>
          <h2 class="mt-1 text-2xl font-black text-white">Popular tracks</h2>
        </div>
        <RouterLink
          to="/recommendations/popular"
          class="text-sm text-slate-400 transition hover:text-white"
        >
          View all
        </RouterLink>
      </div>

      <div v-if="popularLoading" class="space-y-2">
        <div v-for="i in 5" :key="i" class="h-[68px] animate-pulse rounded-2xl bg-white/[0.06]" />
      </div>
      <div
        v-else-if="popularItems.length === 0"
        class="rounded-3xl border border-white/10 bg-black/20 px-6 py-12 text-center"
      >
        <p class="text-sm text-slate-400">No popular tracks available yet.</p>
      </div>
      <div
        v-else
        class="overflow-hidden rounded-3xl border border-white/10 bg-black/20 p-2 backdrop-blur"
      >
        <div
          v-for="(track, index) in popularItems"
          :key="track.id"
          class="group flex cursor-pointer items-center gap-4 rounded-2xl px-3 py-2.5 transition hover:bg-white/[0.04]"
          @click="playFromPopular(track, index)"
        >
          <span class="flex w-6 items-center justify-center">
            <span class="text-sm text-slate-500 group-hover:hidden">{{ index + 1 }}</span>
            <i aria-hidden="true" class="pi pi-play-fill hidden text-sm text-white group-hover:block" />
          </span>

          <div
            class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10 ring-1 ring-white/10"
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
              <i aria-hidden="true" class="pi pi-music text-xs text-slate-500" />
            </div>
          </div>

          <div class="min-w-0 flex-1">
            <p class="truncate font-semibold text-white">{{ track.title }}</p>
            <p class="truncate text-sm text-slate-400">
              {{ track.artist_name || 'Unknown artist' }}
            </p>
          </div>

          <span class="hidden text-xs text-slate-500 md:block">{{ track.genre || '' }}</span>
          <span class="shrink-0 text-xs text-slate-500 tabular-nums">
            {{ formatTime(track.duration_seconds) }}
          </span>
        </div>
      </div>
    </section>

    <!-- Made For You (personalized) -->
    <section v-if="personalizedItems.length > 0" class="mt-12">
      <div class="mb-5 flex items-end justify-between">
        <div>
          <p class="text-sm font-semibold text-[#1db954]">Just for you</p>
          <h2 class="mt-1 text-2xl font-black text-white">Recommended tracks</h2>
        </div>
        <RouterLink
          to="/recommendations/for-you"
          class="text-sm text-slate-400 transition hover:text-white"
        >
          View all
        </RouterLink>
      </div>
      <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
        <div
          v-for="(track, idx) in personalizedItems"
          :key="track.id"
          class="group cursor-pointer"
          @click="playFromPopular(track, idx)"
        >
          <div
            class="relative mb-3 aspect-square overflow-hidden rounded-2xl bg-white/10 shadow-lg ring-1 ring-white/10 transition group-hover:ring-[#1db954]/50"
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
                class="flex h-12 w-12 items-center justify-center rounded-full bg-[#1db954]/90 text-black shadow-xl"
              >
                <i aria-hidden="true" class="pi pi-play-fill text-lg" />
              </div>
            </div>
          </div>
          <p class="truncate text-sm font-semibold text-white">{{ track.title }}</p>
          <p class="truncate text-xs text-slate-500">{{ track.artist_name || 'Unknown' }}</p>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRecommendationsApi } from '@/services/api/recommendation'
import type { RecommendationTrack } from '@/services/api/recommendation'
import { usePlayer } from '@/composables/player'
import { onImgError } from '@/utils/helpers'
import { usePlayerApi } from '@/services/api/player'

const api = useRecommendationsApi()
const player = usePlayer()
const playerApi = usePlayerApi()

const popularItems = ref<RecommendationTrack[]>([])
const personalizedItems = ref<RecommendationTrack[]>([])
const popularLoading = ref(false)

const cards = [
  {
    title: 'For You',
    description: 'Personalized picks based on your taste and listening history.',
    to: '/recommendations/for-you',
    icon: 'pi pi-heart',
    iconClass: 'bg-pink-500/20 text-pink-400',
  },
  {
    title: 'Popular Tracks',
    description: 'The most played tracks across the catalog.',
    to: '/recommendations/popular',
    icon: 'pi pi-chart-line',
    iconClass: 'bg-[#1db954]/20 text-[#1db954]',
  },
  {
    title: 'Best Tracks',
    description: 'Curated high-quality tracks for the best listening session.',
    to: '/recommendations/best',
    icon: 'pi pi-star',
    iconClass: 'bg-yellow-400/20 text-yellow-300',
  },
  {
    title: 'Recently Played',
    description: 'Tracks you listened to recently.',
    to: '/recommendations/recent',
    icon: 'pi pi-clock',
    iconClass: 'bg-sky-400/20 text-sky-300',
  },
]

async function fetchPopular() {
  popularLoading.value = true
  try {
    const response = await api.getPopular({ limit: 5 })
    popularItems.value = response.items
  } catch {
    // silent
  } finally {
    popularLoading.value = false
  }
}

async function fetchPersonalized() {
  try {
    const response = await api.getForYou({ limit: 5 })
    personalizedItems.value = response.items
  } catch {
    // silent
  }
}

function playFromPopular(track: RecommendationTrack) {
  const source = [...popularItems.value, ...personalizedItems.value]
  const allTracks = source.map((t) => ({
    id: String(t.id),
    title: t.title,
    artistName: t.artist_name || 'Unknown',
    albumTitle: t.album_title || null,
    coverUrl: t.cover_url || null,
    durationSeconds: t.duration_seconds ?? null,
    streamUrl: playerApi.getTrackStreamUrl(String(t.id)),
  }))
  const startIdx = allTracks.findIndex((t) => t.id === String(track.id))
  if (startIdx >= 0) {
    player.setQueueAndPlay(allTracks, startIdx)
  }
}

function formatTime(seconds?: number | null) {
  if (!seconds) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}

onMounted(() => {
  fetchPopular()
  fetchPersonalized()
})
</script>
