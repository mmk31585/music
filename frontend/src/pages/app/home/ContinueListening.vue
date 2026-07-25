<template>
  <section v-if="items.length" class="mt-10">
    <div class="mb-4 flex items-end justify-between gap-4">
      <div>
        <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">Jump back in</p>
        <h2 class="mt-1 text-xl font-bold text-white md:text-2xl" style="letter-spacing: -0.02em">Continue Listening</h2>
      </div>
      <RouterLink
        v-if="items.length >= 6"
        to="/recently-played"
        class="shrink-0 whitespace-nowrap text-xs font-medium text-white/30 transition hover:text-white"
      >
        See all
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="ml-0.5 inline h-3 w-3"><polyline points="9 18 15 12 9 6"/></svg>
      </RouterLink>
    </div>
    <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      <button
        v-for="(item, i) in items"
        :key="item.id"
        type="button"
        class="group flex items-center gap-3 rounded-xl border border-white/[3%] bg-white/[2%] p-2 text-left transition hover:bg-white/[6%] hover:border-white/[8%] focus-visible:ring-2 focus-visible:ring-spotify focus-visible:ring-offset-2 focus-visible:outline-hidden"
        :style="{ transitionDelay: `${i * 40}ms` }"
        @click="$emit('play', item)"
        @contextmenu.prevent="openContextMenu($event, item)"
      >
        <div class="relative h-14 w-14 shrink-0 overflow-hidden rounded-lg bg-white/10">
          <img
            v-if="item.cover_url"
            :src="item.cover_url"
            :alt="item.title || ''"
            class="h-full w-full object-cover"
            loading="lazy"
          />
          <div v-else class="flex h-full items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-5 w-5 text-white/30"><path d="M9 18V5l12-2v13"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="16" r="3"/></svg>
          </div>
          <div class="absolute inset-0 flex items-center justify-center bg-black/30 opacity-0 transition group-hover:opacity-100">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="h-6 w-6 text-white"><polygon points="5 3 19 12 5 21 5 3"/></svg>
          </div>
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-semibold text-white">{{ item.title || 'Untitled' }}</p>
          <p class="truncate text-xs text-white/50">{{ item.artist_name || 'Unknown artist' }}</p>
        </div>
        <div class="h-1 w-1 shrink-0 rounded-full bg-spotify" />
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
  items: RecommendationTrack[]
}>()

defineEmits<{
  play: [item: RecommendationTrack]
}>()

// ── Context menu ──────────────────────────────────────────────────
const menuVisible = ref(false)
const menuX = ref(0)
const menuY = ref(0)
const contextTrack = ref<TrackContextItem | null>(null)

function openContextMenu(e: MouseEvent, item: RecommendationTrack) {
  menuX.value = e.clientX
  menuY.value = e.clientY
  contextTrack.value = item as TrackContextItem
  menuVisible.value = true
}

const { sections, header, accentColor } = useTrackContextMenu(
  computed(() => contextTrack.value),
)
</script>
