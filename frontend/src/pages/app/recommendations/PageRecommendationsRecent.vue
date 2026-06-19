<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <section
      class="rounded-[2rem] bg-gradient-to-br from-sky-400 via-blue-800 to-[#101010] p-8 text-white"
    >
      <p class="text-sm font-bold tracking-[0.35em] text-white/70 uppercase">Fresh</p>
      <h1 class="mt-3 text-4xl font-black md:text-6xl">Recently Added</h1>
      <p class="mt-4 max-w-2xl text-white/80">Fresh tracks recently added to the platform.</p>
    </section>

    <section class="mt-10">
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 10" :key="i" class="h-[68px] animate-pulse rounded-2xl bg-white/[0.06]" />
      </div>

      <div
        v-else-if="items.length === 0"
        class="rounded-3xl border border-white/10 bg-black/20 px-6 py-16 text-center"
      >
        <div
          class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-white/10 text-2xl text-white"
        >
          <i aria-hidden="true" class="pi pi-clock" />
        </div>
        <h2 class="mt-5 text-xl font-black text-white">No recent tracks</h2>
        <p class="mt-2 text-sm text-slate-400">New tracks will appear here as they are added.</p>
      </div>

      <div
        v-else
        class="overflow-hidden rounded-3xl border border-white/10 bg-black/20 p-2 backdrop-blur"
      >
        <div
          v-for="(track, index) in items"
          :key="track.id"
          class="flex items-center gap-4 rounded-2xl px-4 py-3 transition hover:bg-white/[0.04]"
        >
          <span class="w-8 text-center text-sm font-bold text-slate-500">{{ index + 1 }}</span>

          <div
            class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-sky-400/20 text-lg text-sky-300"
          >
            <i aria-hidden="true" class="pi pi-clock" />
          </div>

          <div class="min-w-0 flex-1">
            <p class="truncate font-semibold text-white">{{ track.title }}</p>
            <p class="truncate text-sm text-slate-400">
              {{ track.artist_name || 'Unknown artist' }}
            </p>
          </div>

          <span
            v-if="track.album_title"
            class="hidden max-w-[160px] truncate text-sm text-slate-500 md:block"
          >
            {{ track.album_title }}
          </span>

          <span v-if="track.duration_seconds" class="shrink-0 text-xs text-slate-500">
            {{ Math.floor(track.duration_seconds / 60) }}:{{
              String(track.duration_seconds % 60).padStart(2, '0')
            }}
          </span>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRecommendationsApi } from '@/services/api/recommendation'
import type { RecommendationTrack } from '@/services/api/recommendation'

const api = useRecommendationsApi()
const items = ref<RecommendationTrack[]>([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const response = await api.getRecent({ limit: 50 })
    items.value = response.items
  } catch (err) {
    console.error('Failed to fetch recent tracks:', err)
  } finally {
    loading.value = false
  }
})
</script>
