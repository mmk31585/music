<template>
  <section v-if="tracks.length" class="mt-12">
    <div class="mb-5 flex items-end justify-between gap-4">
      <div>
        <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">AI powered</p>
        <h2 class="mt-1 flex items-center gap-2 text-xl font-bold text-white md:text-2xl" style="letter-spacing: -0.02em">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 text-amber-400"><path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"/></svg>
          AI Discoveries
        </h2>
      </div>
      <RouterLink
        to="/ai/playlist-generator"
        class="shrink-0 whitespace-nowrap text-xs font-medium text-white/30 transition hover:text-white"
      >
        Generate
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-0.5 inline h-3 w-3"><polyline points="9 18 15 12 9 6"/></svg>
      </RouterLink>
    </div>
    <div class="flex gap-4 overflow-x-auto pb-2 scrollbar-none">
      <button
        v-for="(track, i) in tracks"
        :key="track.id"
        type="button"
        class="group w-44 shrink-0 space-y-2 text-left focus-visible:outline-hidden"
        @click="$emit('play', track)"
        @contextmenu.prevent="openContextMenu($event, track)"
      >
        <div class="relative aspect-square overflow-hidden rounded-xl bg-white/6 ring-1 ring-white/10 transition-all duration-300 hover:-translate-y-0.5 hover:shadow-lg hover:ring-amber-500/30">
          <img
            v-if="track.cover_url"
            :src="track.cover_url"
            :alt="track.title || ''"
            class="h-full w-full object-cover transition duration-500 group-hover:scale-110"
            loading="lazy"
          />
          <div v-else class="flex h-full items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-8 w-8 text-white/30"><path d="M9 18V5l12-2v13"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="16" r="3"/></svg>
          </div>
          <div class="absolute right-2 top-2 rounded-full bg-amber-500/80 px-2 py-0.5 text-[10px] font-bold text-black backdrop-blur-xs">
            AI
          </div>
          <div class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition group-hover:opacity-100">
            <div class="flex h-10 w-10 items-center justify-center rounded-full bg-amber-400 text-black shadow-xl transition-transform group-hover:scale-110">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5"><polygon points="5 3 19 12 5 21 5 3"/></svg>
            </div>
          </div>
        </div>
        <p class="truncate text-sm font-semibold text-white">{{ track.title || 'Untitled' }}</p>
        <p class="truncate text-xs text-white/60">{{ track.artist_name || 'Unknown artist' }}</p>
        <p v-if="track.score" class="text-[11px] text-amber-400/60">{{ Math.round(track.score) }}% match</p>
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
import type { RecommendationTrack } from '@/services/api/recommendation/types'
import type { TrackContextItem } from '@/composables/useTrackContextMenu'
import { useTrackContextMenu } from '@/composables/useTrackContextMenu'
import ContextMenu from '@/components/common/ContextMenu.vue'

defineProps<{
  tracks: RecommendationTrack[]
}>()

defineEmits<{
  play: [track: RecommendationTrack]
}>()

// ── Context menu ──────────────────────────────────────────────────
const menuVisible = ref(false)
const menuX = ref(0)
const menuY = ref(0)
const contextTrack = ref<TrackContextItem | null>(null)

function openContextMenu(e: MouseEvent, track: RecommendationTrack) {
  menuX.value = e.clientX
  menuY.value = e.clientY
  contextTrack.value = track as TrackContextItem
  menuVisible.value = true
}

const { sections, header, accentColor } = useTrackContextMenu(
  computed(() => contextTrack.value),
)
</script>

<style scoped>
.scrollbar-none {
  scrollbar-width: none;
}
.scrollbar-none::-webkit-scrollbar {
  display: none;
}
</style>
