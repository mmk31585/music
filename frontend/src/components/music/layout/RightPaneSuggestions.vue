<template>
  <div v-if="items.length > 0">
    <h3 class="flex items-center gap-2 px-1 py-3 text-[11px] font-bold tracking-[0.15em] text-slate-500 uppercase">
      <Star :size="12" />
      Suggestions
    </h3>

    <div class="space-y-0.5">
      <div
        v-for="track in items"
        :key="track.id"
        class="group/track flex w-full items-center gap-3 rounded-xl px-3 py-2 transition-all duration-150 hover:bg-white/[0.04]"
        @contextmenu.prevent="openContextMenu($event, track)"
      >
        <div class="h-9 w-9 shrink-0 overflow-hidden rounded-lg bg-white/10">
          <img
            v-if="track.coverUrl"
            :src="track.coverUrl"
            :alt="track.title"
            loading="lazy"
            class="h-full w-full object-cover"
            @error="onImgError"
          />
          <div v-else class="flex h-full items-center justify-center">
            <Music :size="14" class="text-slate-500" />
          </div>
        </div>
        <button
          type="button"
          class="min-w-0 flex-1 text-left"
          @click="$emit('play', track.id)"
        >
          <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
          <p class="truncate text-xs text-slate-500">{{ track.artistName }}</p>
        </button>
        <button
          type="button"
          aria-label="Play suggestion"
          class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-slate-500 opacity-0 transition-all duration-150 hover:bg-white/10 hover:text-primary group-hover/track:opacity-100 active:scale-90"
          @click="$emit('play', track.id)"
        >
          <Play :size="12" />
        </button>
      </div>
    </div>
    <ContextMenu
      v-model:visible="menuVisible"
      :sections="sections"
      :header="header"
      :accent-color="accentColor"
      :position="{ x: menuX, y: menuY }"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Music, Play, Star } from 'lucide-vue-next'
import type { TrackContextItem } from '@/composables/useTrackContextMenu'
import { useTrackContextMenu } from '@/composables/useTrackContextMenu'
import ContextMenu from '@/components/common/ContextMenu.vue'
import { onImgError } from '@/utils/helpers'

export interface SuggestionItem {
  id: string
  title: string
  artistName: string
  coverUrl?: string | null
  durationSeconds?: number | null
}

defineProps<{
  items: SuggestionItem[]
}>()

defineEmits<{
  play: [id: string]
}>()

// ── Context menu ──────────────────────────────────────────────────
const menuVisible = ref(false)
const menuX = ref(0)
const menuY = ref(0)
const contextTrack = ref<TrackContextItem | null>(null)

function openContextMenu(e: MouseEvent, track: SuggestionItem) {
  menuX.value = e.clientX
  menuY.value = e.clientY
  contextTrack.value = {
    id: track.id,
    title: track.title,
    artist_name: track.artistName,
    cover_url: track.coverUrl,
    duration_seconds: track.durationSeconds,
  }
  menuVisible.value = true
}

const { sections, header, accentColor } = useTrackContextMenu(
  computed(() => contextTrack.value),
)
</script>
