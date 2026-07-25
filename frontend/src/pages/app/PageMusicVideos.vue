<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-6 pb-36 md:px-6 lg:px-8">
    <!-- ── Header ── -->
    <div class="mb-8">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-[10px] font-bold tracking-[0.35em] text-white/30 uppercase">Video</p>
          <h1 class="mt-2 text-3xl font-black text-white md:text-4xl">Music Videos</h1>
          <p class="mt-2 text-sm text-slate-400">
            Browse official music videos and fan edits
          </p>
        </div>
      </div>
    </div>

    <!-- ── Type filter chips ── -->
    <div class="mb-8 flex flex-wrap gap-2">
      <button
        v-for="f in filters"
        :key="f.value"
        type="button"
        class="rounded-full px-4 py-2 text-sm font-bold transition-all"
        :class="activeFilter === f.value
          ? 'bg-white text-black'
          : 'bg-white/6 text-white/60 hover:bg-white/10 hover:text-white'"
        @click="activeFilter = f.value; resetVideos()"
      >
        {{ f.label }}
      </button>
    </div>

    <!-- ── Loading state ── -->
    <div v-if="loading && !videos.length && !filteredVideos.length" class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
      <div v-for="i in 10" :key="i" class="animate-pulse">
        <div class="aspect-9/16 w-full rounded-2xl bg-white/4" />
        <div class="mt-2 space-y-1.5 px-1">
          <div class="h-3 w-3/4 rounded bg-white/6" />
          <div class="h-2.5 w-1/2 rounded bg-white/4" />
        </div>
      </div>
    </div>

    <!-- ── Error state ── -->
    <div v-else-if="error && !videos.length && !filteredVideos.length" class="flex flex-col items-center gap-4 py-24 text-center">
      <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/4">
        <AlertCircle aria-hidden="true" class="text-3xl text-slate-500"  />
      </div>
      <h2 class="text-xl font-bold text-white">Failed to load videos</h2>
      <p class="text-sm text-slate-400">Something went wrong. Please try again.</p>
      <button
        type="button"
        class="mt-2 rounded-full bg-white/10 px-6 py-2 text-sm font-bold text-white transition hover:bg-white/20"
        @click="fetchVideos"
      >
        Retry
      </button>
    </div>

    <!-- ── Videos grid ── -->
    <div v-else-if="filteredVideos.length > 0">
      <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
        <VideoCard
          v-for="v in filteredVideos"
          :key="String(v.id)"
          :video="v"
          @open="openVideo(v)"
        />
      </div>

      <!-- Load more -->
      <div v-if="hasMore" class="mt-10 flex justify-center">
        <button
          type="button"
          :disabled="loadingMore"
          class="inline-flex items-center gap-2 rounded-full bg-white/10 px-8 py-3 text-sm font-bold text-white transition hover:bg-white/20 disabled:opacity-50"
          @click="loadMore"
        >
          <Loader2 aria-hidden="true" v-if="loadingMore" class="animate-spin" />
          {{ loadingMore ? 'Loading...' : 'Load More' }}
        </button>
      </div>

      <!-- End of results -->
      <p v-if="!hasMore && filteredVideos.length > 0" class="mt-10 text-center text-sm text-slate-500">
        You've reached the end
      </p>
    </div>

    <!-- ── Empty state (no data OR filter yields nothing) ── -->
    <div v-else class="flex flex-col items-center gap-4 py-24 text-center">
      <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/4">
        <Video aria-hidden="true" class="text-3xl text-slate-500"  />
      </div>
      <h2 class="text-xl font-bold text-white">{{ videos.length ? 'No videos match this filter' : 'No videos yet' }}</h2>
      <p class="text-sm text-slate-400">{{ videos.length ? 'Try selecting a different filter or browse all.' : 'No music videos have been added yet.' }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { AlertCircle, Loader2, Video } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useVideoApi } from '@/services/api/video'
import type { VideoItem } from '@/services/api/video/types'
import VideoCard from '@/components/video/VideoCard.vue'

const router = useRouter()
const videoApi = useVideoApi()

const filters = [
  { label: 'All', value: '' },
  { label: 'Official MV', value: 'official_mv' },
  { label: 'Fan Edits', value: 'user_edit' },
]

const activeFilter = ref('')
const videos = ref<VideoItem[]>([])
const loading = ref(true)
const loadingMore = ref(false)
const error = ref(false)
const offset = ref(0)
const hasMore = ref(true)
const limit = 20

/** Client-side filtered videos based on active type filter */
const filteredVideos = computed(() => {
  if (!activeFilter.value) return videos.value
  return videos.value.filter((v) => v.type === activeFilter.value)
})

function resetVideos() {
  videos.value = []
  offset.value = 0
  hasMore.value = true
  error.value = false
  fetchVideos()
}

async function fetchVideos() {
  loading.value = true
  error.value = false
  try {
    const res = await videoApi.getExploreVideos({ limit, offset: offset.value })
    videos.value = res.items || []
    hasMore.value = (res.items?.length || 0) >= limit
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
}

async function loadMore() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  offset.value += limit
  try {
    const res = await videoApi.getExploreVideos({ limit, offset: offset.value })
    videos.value.push(...(res.items || []))
    hasMore.value = (res.items?.length || 0) >= limit
  } catch {
    offset.value -= limit
  } finally {
    loadingMore.value = false
  }
}

function openVideo(v: VideoItem) {
  router.push(`/music-video/${v.id}`)
}

fetchVideos()
</script>
