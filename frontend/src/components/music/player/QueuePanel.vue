<template>
  <Transition name="slide-right">
    <div
      v-if="visible"
      class="fixed inset-y-0 right-0 z-50 flex w-full max-w-md flex-col border-l border-white/8 shadow-[0_0_60px_rgba(0,0,0,0.5)] backdrop-blur-2xl"
      style="backdrop-filter: blur(32px); -webkit-backdrop-filter: blur(32px); background: rgba(8, 8, 10, 0.94);"
    >
      <div class="flex items-center justify-between border-b border-white/8 px-5 py-4">
        <h2 class="text-lg font-bold text-white">Queue</h2>
        <button
          type="button"
          class="flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/8 hover:text-white active:scale-90"
          aria-label="Close queue"
          @click="visible = false"
        >
          <i aria-hidden="true" class="pi pi-times" />
        </button>
      </div>

      <div class="flex-1 overflow-y-auto p-4" aria-live="polite">
        <!-- Now Playing -->
        <div v-if="currentTrack" class="mb-6">
          <p class="mb-3 text-[10px] font-semibold tracking-wider text-white/30 uppercase">
            Now Playing
          </p>
          <div class="flex items-center gap-3 rounded-xl bg-white/6 px-4 py-3 ring-1 ring-white/6">
            <div class="h-12 w-12 shrink-0 overflow-hidden rounded-lg bg-white/10 ring-1 ring-white/6">
              <img
                v-if="currentTrack.coverUrl"
                :src="currentTrack.coverUrl"
                :alt="currentTrack.title"
                loading="lazy"
                class="h-full w-full object-cover"
                @error="onImgError"
              />
              <div v-else class="flex h-full items-center justify-center">
                <i aria-hidden="true" class="pi pi-music text-slate-500" />
              </div>
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-bold text-white">{{ currentTrack.title }}</p>
              <p class="truncate text-xs text-slate-400">{{ currentTrack.artistName }}</p>
            </div>
            <i aria-hidden="true" class="pi pi-waveform text-lg text-spotify" />
          </div>
        </div>

        <!-- Next Up -->
        <div v-if="displayQueue.length > 0" class="mb-6">
          <div class="mb-3 flex items-center gap-2 text-[10px] font-semibold tracking-wider text-white/30 uppercase">
            <i aria-hidden="true" class="pi pi-arrow-down text-[9px]" />
            Up Next
            <span class="h-3.5 w-3.5 rounded-full bg-white/8 flex items-center justify-center text-[8px] font-bold text-white/40">{{ displayQueue.length }}</span>
            <i v-if="shuffleMode === 'queue'" aria-hidden="true" class="pi pi-sort-alt text-[9px] text-aurora-purple ml-auto" title="Shuffle is on — showing playback order" />
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
              <div
                class="group flex items-center gap-3 rounded-xl px-3 py-2.5 transition-all duration-200 cursor-pointer hover:bg-white/6 active:scale-[0.99]"
                @click="$emit('playFromQueue', index)"
              >
                <span
                  class="drag-handle flex w-5 items-center justify-center text-white/20 transition-colors me-3"
                  :class="canReorder ? 'cursor-grab active:cursor-grabbing hover:text-white/60' : 'cursor-default opacity-30'"
                >
                  <i aria-hidden="true" class="pi pi-bars text-xs" />
                </span>
                <div class="h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10 ring-1 ring-white/6">
                  <img
                    v-if="track.coverUrl"
                    :src="track.coverUrl"
                    :alt="track.title"
                    loading="lazy"
                    class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-110"
                    @error="onImgError"
                  />
                  <div v-else class="flex h-full items-center justify-center">
                    <i aria-hidden="true" class="pi pi-music text-xs text-white/30" />
                  </div>
                </div>
                <div class="min-w-0 flex-1">
                  <p class="truncate text-sm font-medium text-white/90 group-hover:text-white transition-colors">{{ track.title }}</p>
                  <p class="truncate text-xs text-white/40">{{ track.artistName }}</p>
                </div>
                <span class="text-[10px] font-mono tabular-nums text-white/25 group-hover:text-white/50 transition-colors">{{ formatTime(track.durationSeconds) }}</span>
              </div>
            </template>
          </draggable>
        </div>

        <!-- Empty state -->
        <div
          v-if="!currentTrack && queue.length === 0"
          class="flex flex-col items-center gap-3 pt-16 text-center"
        >
          <i aria-hidden="true" class="pi pi-list text-3xl text-white/20" />
          <p class="text-sm text-white/40">Queue is empty</p>
          <p class="text-xs text-white/30">Start playing tracks to see them here</p>
        </div>

        <div
          v-if="currentTrack && queue.length === 0"
          class="flex flex-col items-center gap-3 pt-16 text-center"
        >
          <i aria-hidden="true" class="pi pi-list text-3xl text-white/20" />
          <p class="text-sm text-white/40">No upcoming tracks</p>
          <p class="text-xs text-white/30">Queue will fill as you play more music</p>
        </div>
      </div>
    </div>
  </Transition>

  <Transition name="fade">
    <div
      v-if="visible"
      class="fixed inset-0 z-40 bg-black/60 backdrop-blur-sm"
      @click="visible = false"
    />
  </Transition>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import draggable from 'vuedraggable'
import { onImgError } from '@/utils/helpers'
import { usePlayer } from '@/composables/player'
import { queueManager } from '@/services/player/queue-manager'
import type { PlaybackTrack } from '@/services/api/player'

const visible = defineModel<boolean>('visible', { default: false })

defineEmits<{
  playFromQueue: [index: number]
}>()

const player = usePlayer()
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
    return player.remainingShuffledQueue.value.length > 0
      ? player.remainingShuffledQueue.value
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
    localQueue.value = [...(player.remainingShuffledQueue.value || [])]
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

function formatTime(seconds?: number | null) {
  if (!seconds) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}
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
