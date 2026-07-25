<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <section
      class="rounded-2xl bg-linear-to-br from-sky-500 via-slate-900 to-black p-8 text-white"
    >
      <div class="flex items-center justify-between">
        <div>
          <p class="text-sm font-bold tracking-[0.35em] text-white/70 uppercase">History</p>
          <h1 class="mt-3 text-4xl font-black md:text-6xl">Recently Played</h1>
          <p class="mt-4 max-w-2xl text-white/80">Jump back into tracks you played recently.</p>
        </div>
        <button
          v-if="items.length > 0"
          type="button"
          class="flex h-10 items-center gap-2 rounded-full border border-white/20 px-5 text-sm font-semibold text-white/80 transition hover:border-white/40 hover:text-white"
          @click="clearAll"
        >
          <Trash2 aria-hidden="true" class="text-xs"  />
          Clear All
        </button>
      </div>
    </section>

    <section class="mt-10" aria-live="polite">
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 8" :key="i" class="h-17 animate-pulse rounded-2xl bg-white/6" />
      </div>

      <div
        v-else-if="items.length === 0"
        class="rounded-3xl border border-white/10 bg-black/20 px-6 py-16 text-center"
      >
        <div
          class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-white/10 text-2xl text-white"
        >
          <History aria-hidden="true" class=""  />
        </div>
        <h2 class="mt-5 text-xl font-black text-white">No recent plays yet</h2>
        <p class="mt-2 text-sm text-slate-400">
          Start playing tracks and your history will appear here.
        </p>
        <RouterLink
          to="/discover"
          class="mt-5 inline-flex rounded-full bg-spotify px-6 py-3 text-sm font-bold text-black transition hover:bg-spotify-hover"
        >
          Discover music
        </RouterLink>
      </div>

      <div
        v-else
        class="overflow-hidden rounded-3xl border border-white/10 bg-black/20 p-2 backdrop-blur-xs"
      >
        <div
          v-for="(item, index) in trackRows"
          :key="item.id"
          class="group/track flex items-center gap-2"
        >
          <TrackRow
            :track="item"
            :index="index"
            :queue="trackRows"
            class="min-w-0 flex-1"
          />
          <button
            type="button"
            aria-label="Remove from history"
            class="flex h-8 w-8 items-center justify-center rounded-full text-sm text-slate-500 opacity-0 transition hover:bg-white/10 hover:text-white group-hover/track:opacity-100"
            @click="removeItem(item.historyId)"
          >
            <X aria-hidden="true" class=""  />
          </button>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { History, Trash2, X } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { TrackRow } from '@/components/music'
import { useHistoryApi } from '@/services/api/history'
import { useToast } from 'primevue/usetoast'

const historyApi = useHistoryApi()
const toast = useToast()
const items = ref<Record<string, unknown>[]>([])
const loading = ref(false)

const trackRows = computed(() =>
  items.value.map((item: Record<string, unknown>) => ({
    id: String(item.track_id ?? ''),
    historyId: String(item.id ?? ''),
    title: String(item.track_title ?? ''),
    artist_name: String(item.artist_name ?? ''),
    album_title: String(item.album_title ?? ''),
    cover_url: String(item.track_cover_url ?? ''),
    duration_seconds: Number(item.track_duration ?? 0),
  })),
)

async function removeItem(id: string) {
  try {
    await historyApi.deleteHistoryItem(id)
    items.value = items.value.filter((item) => (item as Record<string, unknown>).id !== id)
    toast.add({ severity: 'success', summary: 'Removed', detail: 'History item removed', life: 3000 })
  } catch {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to remove history item', life: 3000 })
  }
}

async function clearAll() {
  const confirmed = confirm('Are you sure you want to clear all listening history?')
  if (!confirmed) return
  try {
    await historyApi.clearAllHistory()
    items.value = []
    toast.add({ severity: 'success', summary: 'Cleared', detail: 'All history cleared', life: 3000 })
  } catch {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to clear history', life: 3000 })
  }
}

onMounted(async () => {
  loading.value = true
  try {
    const response = await historyApi.getHistory()
    items.value = response?.items ?? []
  } catch (err) {
    console.error('Failed to fetch history:', err)
    items.value = []
  } finally {
    loading.value = false
  }
})
</script>
