<template>
  <div v-if="items.length > 0">
    <div class="flex items-center justify-between px-1 py-3">
      <h3 class="flex items-center gap-2 text-[11px] font-bold tracking-[0.15em] text-tertiary uppercase">
        <List :size="12" />
        Up Next
        <span class="flex h-4 min-w-4 items-center justify-center rounded-full bg-surface-active px-1.5 text-[9px] font-bold text-secondary">{{ items.length }}</span>
      </h3>
      <button
        type="button"
        class="text-[10px] font-bold text-tertiary transition-all duration-150 hover:text-primary"
        @click="$emit('clear')"
      >
        Clear
      </button>
    </div>

    <div class="space-y-0.5">
      <button
        v-for="(track, idx) in visibleItems"
        :key="track.id"
        type="button"
        class="flex w-full items-center gap-3 rounded-xl px-3 py-2 text-left transition-all duration-150 hover:bg-surface-hover active:scale-[0.99]"
        @click="$emit('play', idx)"
        @contextmenu.prevent="openContextMenu($event, track)"
      >
        <div class="relative h-9 w-9 shrink-0 overflow-hidden rounded-lg bg-surface-active">
          <img
            v-if="track.coverUrl"
            :src="track.coverUrl"
            :alt="track.title"
            loading="lazy"
            class="h-full w-full object-cover"
            @error="onImgError"
          />
          <div v-else class="flex h-full items-center justify-center">
            <Music :size="14" class="text-tertiary" />
          </div>
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-medium text-primary">{{ track.title }}</p>
          <p class="truncate text-xs text-tertiary">{{ track.artistName }}</p>
        </div>
        <span class="shrink-0 text-[10px] font-mono tabular-nums text-tertiary">{{ formatTime(track.durationSeconds) }}</span>
      </button>
    </div>

    <button
      v-if="items.length > maxVisible"
      type="button"
      class="mt-1 flex w-full items-center justify-center gap-1.5 rounded-xl px-3 py-2 text-[11px] font-bold text-tertiary transition-all duration-150 hover:bg-surface-hover hover:text-primary"
      @click="showAll = !showAll"
    >
      <ChevronUp v-if="showAll" :size="12" />
      <ChevronDown v-else :size="12" />
      {{ showAll ? 'Show less' : `Show all (${items.length})` }}
    </button>
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
import { ChevronDown, ChevronUp, List, Music } from 'lucide-vue-next'
import type { PlaybackTrack } from '@/services/api/player/types'
import type { TrackContextItem } from '@/composables/useTrackContextMenu'
import { useTrackContextMenu } from '@/composables/useTrackContextMenu'
import ContextMenu from '@/components/common/ContextMenu.vue'
import { onImgError } from '@/utils/helpers'

const props = defineProps<{
  items: PlaybackTrack[]
}>()

defineEmits<{
  play: [index: number]
  clear: []
}>()

const showAll = ref(false)
const maxVisible = 5

const visibleItems = computed(() => {
  if (showAll.value) return props.items
  return props.items.slice(0, maxVisible)
})

function formatTime(seconds?: number | null) {
  if (typeof seconds !== 'number' || !isFinite(seconds)) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}

// ── Context menu ──────────────────────────────────────────────────
const menuVisible = ref(false)
const menuX = ref(0)
const menuY = ref(0)
const contextTrack = ref<TrackContextItem | null>(null)

function openContextMenu(e: MouseEvent, track: PlaybackTrack) {
  menuX.value = e.clientX
  menuY.value = e.clientY
  contextTrack.value = track as unknown as TrackContextItem
  menuVisible.value = true
}

const { sections, header, accentColor } = useTrackContextMenu(
  computed(() => contextTrack.value),
)
</script>
