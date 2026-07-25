<template>
  <section v-if="entries.length" class="mt-12">
    <div class="mb-5 flex items-end justify-between gap-4">
      <div>
        <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">Charts</p>
        <h2 class="mt-1 text-xl font-bold text-white md:text-2xl" style="letter-spacing: -0.02em">Trending Now</h2>
      </div>
      <RouterLink
        to="/recommendations/popular"
        class="shrink-0 whitespace-nowrap text-xs font-medium text-white/30 transition hover:text-white"
      >
        Full chart
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-0.5 inline h-3 w-3"><polyline points="9 18 15 12 9 6"/></svg>
      </RouterLink>
    </div>
    <div class="overflow-hidden rounded-xl border border-white/[4%] bg-white/[1%]">
      <button
        v-for="(entry, i) in entries"
        :key="entry.track.id"
        type="button"
        class="flex w-full cursor-pointer items-center gap-4 border-b border-white/[3%] px-4 py-3 text-left transition hover:bg-white/[3%] last:border-b-0 focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-spotify focus-visible:outline-hidden"
        @click="$emit('play', entry.track)"
        @contextmenu.prevent="openContextMenu($event, entry.track)"
      >
        <span class="w-6 text-center text-sm font-bold" :class="rankClass(i)">{{ entry.rank }}</span>
        <div class="h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10">
          <img
            v-if="entry.track.cover_url"
            :src="entry.track.cover_url"
            :alt="entry.track.title || ''"
            class="h-full w-full object-cover"
            loading="lazy"
          />
          <div v-else class="flex h-full items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4 text-white/30"><path d="M9 18V5l12-2v13"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="16" r="3"/></svg>
          </div>
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-semibold text-white">{{ entry.track.title || 'Untitled' }}</p>
          <p class="truncate text-xs text-white/50">{{ entry.track.artist_name || 'Unknown artist' }}</p>
        </div>
        <div class="flex shrink-0 items-center gap-2">
          <span v-if="entry.isNew" class="rounded-full bg-spotify/20 px-2 py-0.5 text-[10px] font-bold text-spotify">NEW</span>
          <span v-else class="flex items-center gap-1 text-xs" :class="movementClass(entry)">
            <svg v-if="entry.rank < entry.previousRank" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-3 w-3"><polyline points="18 15 12 9 6 15"/></svg>
            <svg v-else-if="entry.rank > entry.previousRank" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-3 w-3"><polyline points="6 9 12 15 18 9"/></svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-3 w-3"><line x1="5" y1="12" x2="19" y2="12"/></svg>
            {{ movementLabel(entry) }}
          </span>
        </div>
      </button>
    </div>
    <ContextMenu
      v-model:visible="menuVisible"
      :sections="sections"
      :header="header"
      :accent-color="accentColor"
      :position="{ x: menuX, y: menuY }"
    />
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { TrendingEntry } from '@/composables/useHomePage'
import type { TrackContextItem } from '@/composables/useTrackContextMenu'
import { useTrackContextMenu } from '@/composables/useTrackContextMenu'
import ContextMenu from '@/components/common/ContextMenu.vue'

defineProps<{
  entries: TrendingEntry[]
}>()

defineEmits<{
  play: [item: TrendingEntry['track']]
}>()

// ── Context menu ──────────────────────────────────────────────────
const menuVisible = ref(false)
const menuX = ref(0)
const menuY = ref(0)
const contextTrack = ref<TrackContextItem | null>(null)

function openContextMenu(e: MouseEvent, track: TrendingEntry['track']) {
  menuX.value = e.clientX
  menuY.value = e.clientY
  contextTrack.value = track as TrackContextItem
  menuVisible.value = true
}

const { sections, header, accentColor } = useTrackContextMenu(
  computed(() => contextTrack.value),
)

function rankClass(index: number) {
  if (index === 0) return 'text-amber-400'
  if (index === 1) return 'text-slate-300'
  if (index === 2) return 'text-amber-600'
  return 'text-white/40'
}

function movementClass(entry: TrendingEntry) {
  if (entry.rank < entry.previousRank) return 'text-spotify'
  if (entry.rank > entry.previousRank) return 'text-red-400'
  return 'text-white/30'
}

function movementLabel(entry: TrendingEntry) {
  const diff = Math.abs(entry.rank - entry.previousRank)
  if (diff === 0) return '—'
  return `${diff}`
}
</script>
