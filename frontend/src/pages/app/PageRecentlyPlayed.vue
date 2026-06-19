<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <section
      class="rounded-[2rem] bg-gradient-to-br from-sky-500 via-slate-900 to-black p-8 text-white"
    >
      <p class="text-sm font-bold tracking-[0.35em] text-white/70 uppercase">History</p>
      <h1 class="mt-3 text-4xl font-black md:text-6xl">Recently Played</h1>
      <p class="mt-4 max-w-2xl text-white/80">Jump back into tracks you played recently.</p>
    </section>

    <section class="mt-10" aria-live="polite">
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 8" :key="i" class="h-[68px] animate-pulse rounded-2xl bg-white/[0.06]" />
      </div>

      <div
        v-else-if="items.length === 0"
        class="rounded-3xl border border-white/10 bg-black/20 px-6 py-16 text-center"
      >
        <div
          class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-white/10 text-2xl text-white"
        >
          <i aria-hidden="true" class="pi pi-history" />
        </div>
        <h2 class="mt-5 text-xl font-black text-white">No recent plays yet</h2>
        <p class="mt-2 text-sm text-slate-400">
          Start playing tracks and your history will appear here.
        </p>
        <RouterLink
          to="/discover"
          class="mt-5 inline-flex rounded-full bg-[#1db954] px-6 py-3 text-sm font-bold text-black transition hover:bg-[#1ed760]"
        >
          Discover music
        </RouterLink>
      </div>

      <div
        v-else
        class="overflow-hidden rounded-3xl border border-white/10 bg-black/20 p-2 backdrop-blur"
      >
        <TrackRow
          v-for="(item, index) in trackRows"
          :key="item.id"
          :track="item"
          :index="index"
          :queue="trackRows"
        />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { TrackRow } from '@/components/music'
import { useLibraryApi } from '@/services/api/library'

const libraryApi = useLibraryApi()
const items = ref<Record<string, unknown>[]>([])
const loading = ref(false)

const trackRows = computed(() =>
  items.value.map((item) => ({
    id: item.track_id,
    title: item.title,
    artist_name: item.artist_name,
    album_title: item.album_title,
    cover_url: item.cover_url,
    duration_seconds: item.duration_seconds,
  })),
)

onMounted(async () => {
  loading.value = true
  try {
    const data = await libraryApi.getRecentlyPlayed()
    items.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Failed to fetch recently played:', err)
    items.value = []
  } finally {
    loading.value = false
  }
})
</script>
