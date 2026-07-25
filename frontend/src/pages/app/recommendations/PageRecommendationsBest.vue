<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <section
      class="rounded-2xl bg-linear-to-br from-yellow-500 via-amber-800 to-surface-base p-8 text-white"
    >
      <p class="text-sm font-bold tracking-[0.35em] text-white/70 uppercase">Curated</p>
      <h1 class="mt-3 text-4xl font-black md:text-6xl">Best Tracks</h1>
      <p class="mt-4 max-w-2xl text-white/80">
        High-quality tracks curated for the best listening experience.
      </p>
    </section>

    <section class="mt-10">
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 10" :key="i" class="h-17 animate-pulse rounded-2xl bg-white/6" />
      </div>

      <div
        v-else-if="items.length === 0"
        class="rounded-3xl border border-white/10 bg-black/20 px-6 py-16 text-center"
      >
        <div
          class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-white/10 text-2xl text-white"
        >
          <Star aria-hidden="true" class=""  />
        </div>
        <h2 class="mt-5 text-xl font-black text-white">No best tracks yet</h2>
        <p class="mt-2 text-sm text-slate-400">Tracks need likes and plays to qualify.</p>
      </div>

      <div
        v-else
        class="overflow-hidden rounded-3xl border border-white/10 bg-black/20 p-2 backdrop-blur-xs"
      >
        <div
          v-for="(track, index) in items"
          :key="track.id"
          class="flex items-center gap-4 rounded-2xl px-4 py-3 transition hover:bg-white/4"
        >
          <span class="w-8 text-center text-sm font-bold text-slate-500">{{ index + 1 }}</span>

          <div
            class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-yellow-400/20 text-lg text-yellow-300"
          >
            <Star aria-hidden="true" class=""  />
          </div>

          <div class="min-w-0 flex-1">
            <p class="truncate font-semibold text-white">{{ track.title }}</p>
            <p class="truncate text-sm text-slate-400">
              {{ track.artist_name || 'Unknown artist' }}
            </p>
          </div>

          <span v-if="track.score" class="shrink-0 text-xs text-yellow-400">
            {{ track.score.toFixed(1) }}
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
import { Star } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { useRecommendationsApi } from '@/services/api/recommendation'
import type { RecommendationTrack } from '@/services/api/recommendation'

const api = useRecommendationsApi()
const items = ref<RecommendationTrack[]>([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const response = await api.getBest({ limit: 50 })
    items.value = response.items
  } catch (err) {
    console.error('Failed to fetch best tracks:', err)
  } finally {
    loading.value = false
  }
})
</script>
