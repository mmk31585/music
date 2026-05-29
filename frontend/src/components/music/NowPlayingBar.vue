<template>
  <div
    class="fixed right-0 bottom-0 left-0 z-50 border-t border-white/10 bg-[#111] px-4 py-3 shadow-2xl"
  >
    <div class="mx-auto grid max-w-7xl grid-cols-[1fr_auto_1fr] items-center gap-4">
      <!-- Track info -->
      <div class="flex min-w-0 items-center gap-3">
        <div class="h-12 w-12 overflow-hidden rounded-lg bg-white/10">
          <img
            v-if="currentTrack?.coverUrl"
            :src="currentTrack.coverUrl"
            :alt="currentTrack.title"
            class="h-full w-full object-cover"
          />

          <div v-else class="flex h-full w-full items-center justify-center text-slate-500">
            <i class="pi pi-music" />
          </div>
        </div>

        <div class="min-w-0">
          <div class="truncate text-sm font-semibold text-white">
            {{ currentTrack?.title || 'No track playing' }}
          </div>

          <div class="truncate text-xs text-slate-400">
            {{ currentTrack?.artistName || 'Select a track' }}
          </div>

          <div v-if="error" class="truncate text-xs text-red-300">
            {{ error }}
          </div>
        </div>
      </div>

      <!-- Main controls -->
      <div class="flex min-w-[360px] flex-col items-center gap-2">
        <div class="flex items-center gap-3">
          <button
            type="button"
            class="text-slate-400 transition hover:text-white disabled:opacity-40"
            :disabled="!hasPrevious"
            @click="playPrevious"
          >
            <i class="pi pi-step-backward" />
          </button>

          <button
            type="button"
            class="flex h-10 w-10 items-center justify-center rounded-full bg-white text-black transition hover:scale-105 disabled:opacity-40"
            :disabled="!currentTrack || isLoadingTrack"
            @click="togglePlayPause"
          >
            <i :class="playIcon" />
          </button>

          <button
            type="button"
            class="text-slate-400 transition hover:text-white disabled:opacity-40"
            :disabled="!hasNext"
            @click="playNext"
          >
            <i class="pi pi-step-forward" />
          </button>
        </div>

        <div class="flex w-full items-center gap-2">
          <span class="w-10 text-right text-[11px] text-slate-500">
            {{ currentTimeLabel }}
          </span>

          <input
            type="range"
            min="0"
            max="100"
            step="0.1"
            class="h-1 w-full cursor-pointer accent-[#1db954]"
            :value="progressPercent"
            :disabled="!currentTrack"
            @input="onSeek"
          />

          <span class="w-10 text-[11px] text-slate-500">
            {{ durationLabel }}
          </span>
        </div>
      </div>

      <!-- Volume -->
      <div class="flex justify-end">
        <div class="flex w-40 items-center gap-2">
          <button
            type="button"
            class="text-slate-400 transition hover:text-white"
            @click="toggleMute"
          >
            <i :class="volumeIcon" />
          </button>

          <input
            type="range"
            min="0"
            max="1"
            step="0.01"
            class="h-1 w-full cursor-pointer accent-[#1db954]"
            :value="volume"
            @input="onVolume"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePlayerControls } from '@/composables/player'

const {
  currentTrack,
  isPlaying,
  isBuffering,
  isLoadingTrack,
  currentTime,
  duration,
  progressPercent,
  volume,
  muted,
  error,
  hasNext,
  hasPrevious,

  playIcon,
  volumeIcon,

  togglePlayPause,
  toggleMute,
  setVolume,
  seekPercent,
  playNext,
  playPrevious,
} = usePlayerControls()

const currentTimeLabel = computed(() => formatTime(currentTime.value))
const durationLabel = computed(() =>
  formatTime(duration.value || currentTrack.value?.durationSeconds || 0),
)

function formatTime(value: number) {
  const total = Math.max(0, Math.floor(Number(value) || 0))
  const minutes = Math.floor(total / 60)
  const seconds = total % 60

  return `${minutes}:${String(seconds).padStart(2, '0')}`
}

function onSeek(event: Event) {
  const target = event.target as HTMLInputElement
  seekPercent(Number(target.value))
}

function onVolume(event: Event) {
  const target = event.target as HTMLInputElement
  setVolume(Number(target.value))
}
</script>
