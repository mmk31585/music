<template>
  <Transition name="slide-right">
    <div
      v-if="visible"
      class="fixed inset-y-0 right-0 z-50 flex w-full max-w-md flex-col border-l border-border-default shadow-floating glass-strong"
    >
      <div class="flex items-center justify-between border-b border-border-default px-5 py-4">
        <h2 class="text-lg font-bold text-primary">Queue</h2>
        <div class="flex items-center gap-1">
          <button
            v-if="queue.length > 0"
            type="button"
            class="flex h-7 items-center gap-1.5 rounded-full px-3 text-[11px] font-medium text-red-400/70 transition bg-surface-active hover:text-red-400 active:scale-90"
            aria-label="Clear queue"
            @click="clearQueue"
          >
            <Trash2 aria-hidden="true" class="text-[10px]"  />
            Clear
          </button>
        <button
            type="button"
            class="flex h-9 w-9 items-center justify-center rounded-full text-secondary transition bg-surface-active hover:text-primary active:scale-90"
            aria-label="Close queue"
            @click="visible = false"
          >
            <X aria-hidden="true" class=""  />
          </button>
        </div>
      </div>

      <div class="flex-1 overflow-y-auto p-4" aria-live="polite">
        <!-- Now Playing -->
        <div v-if="currentTrack" class="mb-6">
          <p class="mb-3 text-[10px] font-semibold tracking-wider text-muted uppercase">
            Now Playing
          </p>
          <div class="flex items-center gap-3 rounded-xl bg-surface-overlay px-4 py-3 ring-1 ring-border-subtle">
            <div class="h-12 w-12 shrink-0 overflow-hidden rounded-lg bg-surface-active ring-1 ring-border-subtle">
              <img
                v-if="currentTrack.coverUrl"
                :src="currentTrack.coverUrl"
                :alt="currentTrack.title"
                loading="lazy"
                class="h-full w-full object-cover"
                @error="onImgError"
              />
              <div v-else class="flex h-full items-center justify-center">
                <Music aria-hidden="true" class="text-muted"  />
              </div>
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-bold text-primary">{{ currentTrack.title }}</p>
              <p class="truncate text-xs text-secondary">{{ currentTrack.artistName }}</p>
            </div>
            <AudioLines aria-hidden="true" class="text-lg text-accent"  />
          </div>
        </div>

        <!-- Next Up -->
        <div v-if="displayQueue.length > 0" class="mb-6">
          <div class="mb-3 flex items-center gap-2 text-[10px] font-semibold tracking-wider text-muted uppercase">
            <ArrowDown aria-hidden="true" class="text-[9px]"  />
            Up Next
            <span class="h-3.5 w-3.5 rounded-full bg-surface-active flex items-center justify-center text-[8px] font-bold text-tertiary">{{ displayQueue.length }}</span>
            <ArrowUpDown v-if="shuffleMode === 'queue'" aria-hidden="true" class="text-[9px] text-aurora-purple ml-auto" title="Shuffle is on — showing playback order"  />
          </div>
          <draggable
            :list="displayQueue"
            :item-key="'id'"
            :handle="canReorder ? '.drag-handle' : undefined"
            :animation="200"
            :disabled="!canReorder"
            ghost-class="opacity-30"
            class="space-y-1"
            @end="onReorder"
          >
            <template #item="{ element: track, index }">
              <button
                :id="'queue-item-' + index"
                type="button"
                class="group flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-left transition-all duration-200"
                :class="canReorder ? 'cursor-pointer' : ''"
                :aria-label="`${track.title} by ${track.artistName}. Press Enter to play, Arrow Up or Down to reorder.`"
                @click="$emit('playFromQueue', index)"
                @contextmenu.prevent="openContextMenu($event, track)"
                @keydown.up.prevent="moveItem(index, -1)"
                @keydown.down.prevent="moveItem(index, 1)"
              >
                <span
                  v-if="canReorder"
                  class="drag-handle flex w-5 items-center justify-center text-secondary transition-colors me-3 cursor-grab active:cursor-grabbing hover:text-secondary"
                  aria-label="Drag to reorder, or use Arrow Up/Down keys"
                >
                  <GripVertical class="text-xs" aria-hidden="true"  />
                </span>
                <span
                  v-else
                  class="flex w-5 items-center justify-center text-secondary me-3 opacity-30"
                >
                  <GripVertical class="text-xs" aria-hidden="true"  />
                </span>
                <div class="h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-surface-active ring-1 ring-border-subtle">
                  <img
                    v-if="track.coverUrl"
                    :src="track.coverUrl"
                    :alt="track.title"
                    loading="lazy"
                    class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-110"
                    @error="onImgError"
                  />
                  <div v-else class="flex h-full items-center justify-center">
                    <Music aria-hidden="true" class="text-xs text-muted"  />
                  </div>
                </div>
                <div class="min-w-0 flex-1">
                  <p class="truncate text-sm font-medium text-secondary group-hover:text-primary transition-colors">{{ track.title }}</p>
                  <p class="truncate text-xs text-tertiary">{{ track.artistName }}</p>
                </div>
                <span class="text-[10px] font-mono tabular-nums text-secondary group-hover:text-secondary transition-colors">{{ formatTime(track.durationSeconds) }}</span>
              </button>
            </template>
          </draggable>
        </div>

        <!-- Empty state -->
        <div
          v-if="!currentTrack && queue.length === 0"
          class="flex flex-col items-center gap-3 pt-16 text-center"
        >
          <List aria-hidden="true" class="text-3xl text-secondary"  />
          <p class="text-sm text-tertiary">Queue is empty — add some tracks to get started</p>
        </div>

        <div
          v-if="currentTrack && queue.length === 0"
          class="flex flex-col items-center gap-3 pt-16 text-center"
        >
          <List aria-hidden="true" class="text-3xl text-secondary"  />
          <p class="text-sm text-tertiary">No upcoming tracks — queue will fill as you play more music</p>
        </div>
      </div>
    </div>
  </Transition>

  <Transition name="fade">
    <div
      v-if="visible"
      class="fixed inset-0 z-40 bg-bg-overlay/60 backdrop-blur-sm"
      @click="visible = false"
    />
  </Transition>
  <ContextMenu
    v-model:visible="menuVisible"
    :sections="sections"
    :header="header"
    :accent-color="accentColor"
    :position="{ x: menuX, y: menuY }"
  />
</template>

<script setup lang="ts">
import { ArrowDown, ArrowUpDown, AudioLines, GripVertical, List, Music, Trash2, X } from 'lucide-vue-next'
import { computed, nextTick, ref, watch } from 'vue'
import draggable from 'vuedraggable'
import { useToast } from 'primevue/usetoast'
import type { TrackContextItem } from '@/composables/useTrackContextMenu'
import { useTrackContextMenu } from '@/composables/useTrackContextMenu'
import ContextMenu from '@/components/common/ContextMenu.vue'
import { onImgError } from '@/utils/helpers'
import { usePlayer } from '@/composables/player'
import { queueManager } from '@/services/player/queue-manager'
import type { PlaybackTrack } from '@/services/api/player'

const visible = defineModel<boolean>('visible', { default: false })

defineEmits<{
  playFromQueue: [index: number]
}>()

const player = usePlayer()

const toast = useToast()

function clearQueue() {
  queueManager.clear()
  player.queue.value = []
  toast.add({ severity: 'info', summary: 'Queue cleared', life: 2000 })
}
const currentTrack = player.currentTrack
const queue = player.queue
const shuffleMode = player.shuffleMode

const localQueue = ref<PlaybackTrack[]>([])

/**
 * B3: When shuffle mode is 'queue', show tracks in the order they'll actually play.
 * For other shuffle modes or when shuffle is off, show the raw queue order.
 */
const displayQueue = computed<PlaybackTrack[]>(() => {
  if (shuffleMode.value === 'queue') {
    return player.remainingShuffledQueue.length > 0
      ? player.remainingShuffledQueue
      : localQueue.value
  }
  return localQueue.value
})

/** Is drag-to-reorder possible? Disabled when shuffle or catalog/similar modes are active. */
const canReorder = computed(() => shuffleMode.value === 'off')

/**
 * B2: Watch the full queue array reference (not just length) so that
 * drag-and-drop reorder and external queue changes both trigger reactively.
 */
watch(() => player.queue.value, (newQueue) => {
  localQueue.value = [...(newQueue || [])]
}, { immediate: true })

/**
 * B3: When shuffle mode changes, rebuild localQueue from the store.
 * For 'queue' mode this uses the shuffled order; for 'off' mode it falls
 * back to the raw queue.
 */
watch(() => player.shuffleMode.value, () => {
  if (player.shuffleMode.value === 'queue') {
    localQueue.value = [...(player.remainingShuffledQueue || [])]
  } else {
    localQueue.value = [...(player.queue.value || [])]
  }
})

function onReorder(event: { oldIndex: number; newIndex: number }) {
  // Apply the drag reorder to the shared QueueManager singleton
  queueManager.reorderQueue(event.oldIndex, event.newIndex)
  // Update store with a new array reference to guarantee reactivity
  const reordered = queueManager.all()
  player.queue.value = [...reordered]
  localQueue.value = reordered
}

/**
 * Keyboard alternative for drag-to-reorder.
 * Arrow Up moves the item up (earlier in queue), Arrow Down moves it down.
 */
function moveItem(currentIndex: number, direction: -1 | 1) {
  const newIndex = currentIndex + direction
  if (newIndex < 0 || newIndex >= localQueue.value.length) return

  queueManager.reorderQueue(currentIndex, newIndex)
  const reordered = queueManager.all()
  player.queue.value = [...reordered]
  localQueue.value = reordered

  // Move focus to the item at its new position
  nextTick(() => {
    const el = document.getElementById('queue-item-' + newIndex)
    el?.focus()
  })
}

function formatTime(seconds?: number | null) {
  if (!seconds) return '0:00'
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

<style scoped>
.slide-right-enter-active,
.slide-right-leave-active {
  transition: all 0.25s ease;
}
.slide-right-enter-from {
  transform: translateX(100%);
}
.slide-right-leave-to {
  transform: translateX(100%);
}
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
