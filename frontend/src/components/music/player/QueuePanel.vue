<template>
  <Transition name="slide-right">
    <div
      v-if="visible"
      class="fixed inset-y-0 right-0 z-50 flex w-full max-w-md flex-col border-l border-white/10 bg-[#121212] shadow-2xl"
    >
      <div class="flex items-center justify-between border-b border-white/10 px-5 py-4">
        <h2 class="text-lg font-bold text-white">Queue</h2>
        <button
          type="button"
          class="flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
          aria-label="Close queue"
          @click="visible = false"
        >
          <i aria-hidden="true" class="pi pi-times" />
        </button>
      </div>

      <div class="flex-1 overflow-y-auto p-4" aria-live="polite">
        <div v-if="currentTrack" class="mb-6">
          <p class="mb-3 text-xs font-semibold tracking-wider text-slate-500 uppercase">
            Now Playing
          </p>
          <div class="flex items-center gap-3 rounded-xl bg-[#1db954]/10 px-4 py-3">
            <div class="h-12 w-12 shrink-0 overflow-hidden rounded-lg bg-white/10">
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
            <i aria-hidden="true" class="pi pi-waveform text-lg text-[#1db954]" />
          </div>
        </div>

        <div v-if="queue.length > 0" class="mb-6">
          <p class="mb-3 text-xs font-semibold tracking-wider text-slate-500 uppercase">
            Next Up
          </p>
          <draggable
            :list="localQueue"
            item-key="id"
            handle=".drag-handle"
            animation="200"
            ghost-class="opacity-30"
            class="space-y-1"
            @end="onReorder"
          >
            <template #item="{ element: track, index }">
              <div
                class="flex items-center gap-3 rounded-xl px-3 py-2 transition hover:bg-white/[0.06]"
                @click="$emit('playFromQueue', index)"
              >
                <span class="drag-handle flex w-5 items-center justify-center cursor-grab active:cursor-grabbing text-white/30 hover:text-white/70 transition-colors me-3">
                  <i aria-hidden="true" class="pi pi-bars text-xs" />
                </span>
                <div class="h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10">
                  <img
                    v-if="track.coverUrl"
                    :src="track.coverUrl"
                    :alt="track.title"
                    loading="lazy"
                    class="h-full w-full object-cover"
                    @error="onImgError"
                  />
                  <div v-else class="flex h-full items-center justify-center">
                    <i aria-hidden="true" class="pi pi-music text-xs text-slate-500" />
                  </div>
                </div>
                <div class="min-w-0 flex-1">
                  <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
                  <p class="truncate text-xs text-slate-400">{{ track.artistName }}</p>
                </div>
                <span class="text-xs text-slate-500">{{ formatTime(track.durationSeconds) }}</span>
              </div>
            </template>
          </draggable>
        </div>

        <div
          v-if="!currentTrack && queue.length === 0"
          class="flex flex-col items-center gap-3 pt-16 text-center"
        >
          <i aria-hidden="true" class="pi pi-list text-3xl text-slate-500" />
          <p class="text-sm text-slate-400">Queue is empty</p>
          <p class="text-xs text-slate-500">Start playing tracks to see them here</p>
        </div>
      </div>
    </div>
  </Transition>

  <Transition name="fade">
    <div
      v-if="visible"
      class="fixed inset-0 z-40 bg-black/50 backdrop-blur-sm"
      @click="visible = false"
    />
  </Transition>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
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

const localQueue = ref<PlaybackTrack[]>([])

watch(() => player.queue.value?.length, () => {
  localQueue.value = [...(player.queue.value || [])]
}, { immediate: true })

function onReorder(event: { oldIndex: number; newIndex: number }) {
  queueManager.reorderQueue(event.oldIndex, event.newIndex)
  player.queue.value = queueManager.all()
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
