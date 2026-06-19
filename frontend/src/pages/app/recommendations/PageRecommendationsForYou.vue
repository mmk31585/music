<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <section
      class="rounded-[2rem] bg-gradient-to-br from-pink-600 via-purple-700 to-[#101010] p-8 text-white md:p-12"
    >
      <p class="text-sm font-bold tracking-[0.35em] text-white/70 uppercase">Personalized</p>
      <h1 class="mt-3 text-4xl font-black md:text-6xl">For You</h1>
      <p class="mt-4 max-w-2xl text-white/80">
        Tracks selected for you based on your listening history, followed artists, and what similar
        listeners enjoy.
      </p>
    </section>

    <div
      v-if="loading"
      class="mt-10 grid gap-5 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5"
    >
      <div v-for="i in 20" :key="i" class="space-y-3">
        <div class="aspect-square animate-pulse rounded-2xl bg-white/[0.06]" />
        <div class="h-4 w-3/4 animate-pulse rounded bg-white/[0.06]" />
        <div class="h-3 w-1/2 animate-pulse rounded bg-white/[0.06]" />
      </div>
    </div>

    <div
      v-else-if="items.length === 0"
      class="mt-10 flex flex-col items-center gap-4 rounded-3xl border border-white/10 px-6 py-20 text-center"
    >
      <div class="flex h-16 w-16 items-center justify-center rounded-full bg-white/10">
        <i aria-hidden="true" class="pi pi-heart text-2xl text-slate-400" />
      </div>
      <h3 class="text-xl font-bold text-white">No recommendations yet</h3>
      <p class="max-w-sm text-sm text-slate-400">
        Listen to more tracks and follow artists to get personalized recommendations.
      </p>
      <RouterLink
        to="/recommendations/popular"
        class="rounded-full bg-[#1db954] px-5 py-2 text-sm font-bold text-black transition hover:bg-[#1ed760]"
      >
        Explore popular tracks
      </RouterLink>
    </div>

    <div
      v-else
      class="mt-10 grid gap-5 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5"
    >
      <div
        v-for="(track, idx) in items"
        :key="track.id"
        class="group cursor-pointer"
        @click="play(track, idx)"
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
            v-if="track.score"
            class="absolute top-2 left-2 rounded-full bg-black/60 px-2 py-0.5 text-[10px] font-bold text-pink-400 backdrop-blur"
          >
            {{ track.score.toFixed(1) }}
          </div>
        </div>
        <p class="truncate text-sm font-semibold text-white">{{ track.title }}</p>
        <p class="truncate text-xs text-slate-500">{{ track.artist_name || 'Unknown' }}</p>
      </div>
    </div>
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

const loading = ref(true)
const items = ref<RecommendationTrack[]>([])

async function fetchForYou() {
  loading.value = true
  try {
    const response = await api.getForYou({ limit: 50 })
    items.value = response.items
  } catch {
    // silent
  } finally {
    loading.value = false
  }
}

function play(track: RecommendationTrack, index: number) {
  const allTracks = items.value.map((t) => ({
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

onMounted(fetchForYou)
</script>
