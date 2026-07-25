<template>
  <section class="mt-12">
    <div class="mb-5">
      <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">Curated picks</p>
      <h2 class="mt-1 text-xl font-bold text-white md:text-2xl" style="letter-spacing: -0.02em">Editorial Discoveries</h2>
    </div>
    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div
        v-for="(item, i) in items"
        :key="item.id"
        class="group relative overflow-hidden rounded-2xl border border-white/[4%] bg-white/[2%] transition-all hover:-translate-y-0.5 hover:border-white/[10%]"
        @contextmenu.prevent="openContextMenu($event, item)"
      >
        <div class="flex flex-col sm:flex-row">
          <div class="relative h-32 w-full shrink-0 overflow-hidden sm:h-auto sm:w-32">
            <img
              v-if="item.cover_url"
              :src="item.cover_url"
              :alt="item.title || ''"
              class="h-full w-full object-cover transition duration-500 group-hover:scale-105"
              loading="lazy"
            />
            <div v-else class="flex h-full items-center justify-center bg-white/5">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-8 w-8 text-white/20"><path d="M9 18V5l12-2v13"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="16" r="3"/></svg>
            </div>
          </div>
          <div class="flex flex-1 flex-col justify-center gap-2 p-4">
            <p class="text-[10px] font-bold tracking-[0.2em] text-spotify/70 uppercase">{{ tags[i % tags.length] }}</p>
            <p class="text-sm font-bold text-white leading-snug">{{ item.title || 'Untitled' }}</p>
            <p class="text-xs text-white/40">{{ item.artist_name || 'Unknown artist' }}</p>
            <button
              type="button"
              class="mt-1 inline-flex h-8 w-fit cursor-pointer items-center gap-1.5 rounded-full bg-spotify px-4 text-[11px] font-bold text-black transition hover:bg-spotify-hover"
              @click="$emit('play', item)"
            >
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="h-3 w-3"><polygon points="5 3 19 12 5 21 5 3"/></svg>
              Listen
            </button>
          </div>
        </div>
      </div>
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

const tags = ['Editorial Pick', 'Must Hear', 'New Release', 'Staff Favorite', 'Hidden Gem']

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
