<template>
  <div class="mx-auto w-full max-w-7xl px-4 pb-36 pt-6 md:px-6 lg:px-8">
    <!-- Header -->
    <div class="mb-6">
      <div class="flex items-center gap-3">
        <div class="flex h-10 w-10 items-center justify-center rounded-2xl bg-spotify/10">
          <Compass aria-hidden="true" class="text-spotify"  />
        </div>
        <div>
          <h1 class="text-2xl font-black text-white sm:text-3xl">Explore</h1>
          <p class="mt-0.5 text-sm text-white/40">Official MVs &amp; fan edits from the community</p>
        </div>
      </div>
    </div>

    <!-- Filter chips -->
    <div class="mb-6 flex flex-wrap gap-2">
      <button
        v-for="filter in filters"
        :key="filter.key"
        type="button"
        class="rounded-full px-4 py-2 text-sm font-bold transition-all duration-200"
        :class="activeFilter === filter.key
          ? 'bg-spotify text-black shadow-lg shadow-spotify/20'
          : 'bg-white/5 text-white/60 hover:bg-white/10 hover:text-white'"
        @click="activeFilter = filter.key"
      >
        {{ filter.label }}
      </button>
    </div>

    <!-- Loading skeleton -->
    <div
      v-if="loading && !videos.length"
      class="grid grid-cols-2 gap-3 sm:grid-cols-3 sm:gap-4 lg:grid-cols-4 xl:grid-cols-5"
    >
      <div v-for="i in 10" :key="i">
        <div class="aspect-9/16 w-full animate-pulse rounded-2xl bg-white/5" />
        <div class="mt-2 h-3 w-3/4 animate-pulse rounded bg-white/5" />
        <div class="mt-1.5 h-2.5 w-1/2 animate-pulse rounded bg-white/5" />
      </div>
    </div>

    <!-- Video grid -->
    <div
      v-else-if="filteredVideos.length"
      ref="gridRef"
      class="grid grid-cols-2 gap-3 sm:grid-cols-3 sm:gap-4 lg:grid-cols-4 xl:grid-cols-5"
    >
      <div
        v-for="(video, index) in filteredVideos"
        :key="video.id"
        class="explore-card"
        :style="{ '--i': index }"
      >
        <VideoCard
          :video="video"
          @open="openPlayer(index)"
        />
      </div>
    </div>

    <!-- Empty state -->
    <div
      v-else
      class="flex flex-col items-center justify-center py-24"
    >
      <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/5">
        <Video aria-hidden="true" class="text-2xl text-slate-500"  />
      </div>
      <p class="mt-4 text-sm font-medium text-white/60">
        {{ activeFilter === 'all' ? 'No videos yet' : `No ${activeFilter === 'official_mv' ? 'official MVs' : 'fan edits'} found` }}
      </p>
      <p v-if="activeFilter !== 'all'" class="mt-1 text-xs text-white/30">
        <button class="text-spotify hover:underline" @click="activeFilter = 'all'">Show all videos</button>
      </p>
    </div>

    <!-- Load more sentinel -->
    <div
      v-if="hasMore && filteredVideos.length"
      ref="sentinelRef"
      class="flex items-center justify-center py-8"
    >
      <div v-if="loadingMore" class="flex items-center gap-3 text-sm text-white/40">
        <Loader2 aria-hidden="true" class="animate-spin"  />
        <span>Loading more...</span>
      </div>
      <div v-else class="flex items-center gap-2 text-xs text-white/20">
        <span class="h-px w-8 bg-white/6" />
        <span>Scroll for more</span>
        <span class="h-px w-8 bg-white/6" />
      </div>
    </div>

    <!-- Vertical video player overlay -->
    <VerticalVideoPlayer
      v-if="playerOpen"
      :videos="allLoadedVideos"
      :start-index="playerStartIndex"
      @close="playerOpen = false"
    />
  </div>
</template>

<script setup lang="ts">
import { Compass, Loader2, Video } from 'lucide-vue-next'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import type { VideoItem } from '@/services/api/video/types'
import { useVideoApi } from '@/services/api/video'
import VideoCard from '@/components/video/VideoCard.vue'
import VerticalVideoPlayer from '@/components/video/VerticalVideoPlayer.vue'

const videoApi = useVideoApi()

// ── Filters ───────────────────────────────────────────────────────────

const filters = [
  { key: 'all', label: 'All' },
  { key: 'official_mv', label: 'Official MVs' },
  { key: 'user_edit', label: 'Fan Edits' },
] as const

type FilterKey = (typeof filters)[number]['key']
const activeFilter = ref<FilterKey>('all')

// ── State ─────────────────────────────────────────────────────────────

const videos = ref<VideoItem[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const offset = ref(0)
const hasMore = ref(true)
const LIMIT = 30

const gridRef = ref<HTMLElement | null>(null)
const sentinelRef = ref<HTMLElement | null>(null)

// Player
const playerOpen = ref(false)
const playerStartIndex = ref(0)

// ── Computed ──────────────────────────────────────────────────────────

const filteredVideos = computed(() => {
  if (activeFilter.value === 'all') return videos.value
  return videos.value.filter((v) => v.type === activeFilter.value)
})

const allLoadedVideos = computed(() => videos.value)

// ── Load data ─────────────────────────────────────────────────────────

async function loadVideos() {
  if (loading.value) return
  loading.value = true
  try {
    const response = await videoApi.getExploreVideos({ limit: LIMIT, offset: 0 })
    if (response?.items) {
      videos.value = response.items
      offset.value = response.offset + (response.items.length || 0)
      hasMore.value = (response.items.length || 0) >= LIMIT
    }
  } catch {
    // Silently fail — grid stays empty
  } finally {
    loading.value = false
  }
}

async function loadMore() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  try {
    const response = await videoApi.getExploreVideos({ limit: LIMIT, offset: offset.value })
    if (response?.items?.length) {
      videos.value = [...videos.value, ...response.items]
      offset.value += response.items.length
      hasMore.value = response.items.length >= LIMIT
    } else {
      hasMore.value = false
    }
  } catch {
    hasMore.value = false
  } finally {
    loadingMore.value = false
  }
}

// Reset scroll position when filter changes
watch(activeFilter, () => {
  // scroll to top of grid smoothly
  gridRef.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
})

// ── Infinite scroll ───────────────────────────────────────────────────

let observer: IntersectionObserver | null = null

onMounted(() => {
  loadVideos()

  observer = new IntersectionObserver(
    (entries) => {
      if (entries[0]?.isIntersecting) {
        loadMore()
      }
    },
    { rootMargin: '400px' },
  )

  if (sentinelRef.value) {
    observer.observe(sentinelRef.value)
  }
})

onUnmounted(() => {
  observer?.disconnect()
})

// ── Player ────────────────────────────────────────────────────────────

function openPlayer(index: number) {
  playerStartIndex.value = index
  playerOpen.value = true
}
</script>

<style scoped>
/* Stagger entrance animation for video cards */
.explore-card {
  opacity: 0;
  transform: translateY(12px);
  animation: explore-card-enter 0.35s ease-out forwards;
  animation-delay: calc(var(--i) * 40ms);
}

@keyframes explore-card-enter {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (prefers-reduced-motion: reduce) {
  .explore-card {
    opacity: 1;
    transform: none;
    animation: none;
  }
}
</style>
