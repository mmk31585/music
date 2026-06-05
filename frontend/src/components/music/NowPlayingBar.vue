<template>
  <div
    class="fixed inset-x-0 bottom-0 z-50 border-t border-white/10 bg-black/80 px-3 py-3 shadow-[0_-20px_60px_rgba(0,0,0,0.65)] backdrop-blur-2xl md:px-4"
  >
    <div
      class="mx-auto grid max-w-7xl grid-cols-1 items-center gap-3 md:grid-cols-[minmax(0,1fr)_minmax(340px,520px)_minmax(0,1fr)] md:gap-5"
    >
      <!-- Track info -->
      <div class="flex min-w-0 items-center gap-3">
        <div
          class="relative h-12 w-12 shrink-0 overflow-hidden rounded-xl bg-white/10 shadow-lg ring-1 ring-white/10 md:h-14 md:w-14"
          :class="currentTrack ? 'animate-art-pop' : ''"
        >
          <img
            v-if="currentTrack?.coverUrl"
            :src="currentTrack.coverUrl"
            :alt="currentTrack.title"
            class="h-full w-full object-cover"
          />

          <div
            v-else
            class="flex h-full w-full items-center justify-center bg-gradient-to-br from-white/10 to-white/5 text-slate-500"
          >
            <i class="pi pi-music" />
          </div>

          <div
            v-if="isPlaying"
            class="absolute inset-0 bg-gradient-to-t from-black/30 via-transparent to-transparent"
          />
        </div>

        <div class="min-w-0">
          <div class="flex min-w-0 items-center gap-2">
            <div class="truncate text-sm font-black text-white">
              {{ currentTrack?.title || 'No track playing' }}
            </div>

            <span
              v-if="isBuffering || isLoadingTrack"
              class="inline-flex shrink-0 items-center gap-1 rounded-full bg-[#1db954]/15 px-2 py-0.5 text-[10px] font-bold text-[#1db954]"
            >
              <i class="pi pi-spin pi-spinner text-[10px]" />
              Loading
            </span>
          </div>

          <div class="mt-0.5 truncate text-xs font-medium text-slate-400">
            {{ currentTrack?.artistName || 'Select a track to start listening' }}
          </div>

          <div v-if="error" class="mt-0.5 max-w-[260px] truncate text-xs text-red-300">
            {{ error }}
          </div>
        </div>

        <button
          type="button"
          class="ml-1 hidden h-9 w-9 shrink-0 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white sm:flex"
          :disabled="!currentTrack"
          title="Like"
        >
          <i class="pi pi-heart" />
        </button>
      </div>

      <!-- Main controls -->
      <div class="flex min-w-0 flex-col items-center gap-2">
        <div class="flex items-center gap-3">
          <button
            type="button"
            class="hidden h-8 w-8 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white disabled:cursor-not-allowed disabled:opacity-30 sm:flex"
            :disabled="!currentTrack"
            title="Shuffle"
          >
            <i class="pi pi-sort-alt text-sm" />
          </button>

          <button
            type="button"
            class="flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white disabled:cursor-not-allowed disabled:opacity-30"
            :disabled="!hasPrevious"
            title="Previous"
            @click="playPrevious"
          >
            <i class="pi pi-step-backward" />
          </button>

          <button
            type="button"
            class="relative flex h-11 w-11 items-center justify-center rounded-full bg-white text-black shadow-lg transition hover:scale-105 hover:bg-[#1db954] disabled:cursor-not-allowed disabled:opacity-40"
            :disabled="!currentTrack || isLoadingTrack"
            title="Play/Pause"
            @click="togglePlayPause"
          >
            <i v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-sm" />
            <i v-else :class="playIcon" class="text-sm" />
          </button>

          <button
            type="button"
            class="flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white disabled:cursor-not-allowed disabled:opacity-30"
            :disabled="!hasNext"
            title="Next"
            @click="playNext"
          >
            <i class="pi pi-step-forward" />
          </button>

          <button
            type="button"
            class="hidden h-8 w-8 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white disabled:cursor-not-allowed disabled:opacity-30 sm:flex"
            :disabled="!currentTrack"
            title="Repeat"
          >
            <i class="pi pi-refresh text-sm" />
          </button>
        </div>

        <div class="flex w-full items-center gap-2">
          <span class="w-10 text-right text-[11px] font-medium text-slate-500 tabular-nums">
            {{ currentTimeLabel }}
          </span>

          <div class="group relative flex flex-1 items-center">
            <input
              type="range"
              min="0"
              max="100"
              step="0.1"
              class="player-range"
              :style="progressStyle"
              :value="progressPercent"
              :disabled="!currentTrack"
              @input="onSeek"
            />
          </div>

          <span class="w-10 text-[11px] font-medium text-slate-500 tabular-nums">
            {{ durationLabel }}
          </span>
        </div>
      </div>

      <!-- Right controls -->
      <div class="hidden min-w-0 justify-end md:flex">
        <div class="flex items-center gap-3">
          <button
            type="button"
            class="flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
            :disabled="!currentTrack"
            title="Queue"
          >
            <i class="pi pi-list" />
          </button>

          <button
            type="button"
            class="flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
            :disabled="!currentTrack"
            title="Lyrics"
          >
            <i class="pi pi-align-left" />
          </button>

          <div class="flex w-40 items-center gap-2">
            <button
              type="button"
              class="flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
              title="Mute"
              @click="toggleMute"
            >
              <i :class="volumeIcon" />
            </button>

            <input
              type="range"
              min="0"
              max="1"
              step="0.01"
              class="player-range"
              :style="volumeStyle"
              :value="volume"
              @input="onVolume"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile volume/actions -->
    <div class="mx-auto mt-2 flex max-w-7xl items-center justify-between md:hidden">
      <button
        type="button"
        class="flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
        :disabled="!currentTrack"
      >
        <i class="pi pi-heart" />
      </button>

      <div class="flex w-40 items-center gap-2">
        <button
          type="button"
          class="flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
          @click="toggleMute"
        >
          <i :class="volumeIcon" />
        </button>

        <input
          type="range"
          min="0"
          max="1"
          step="0.01"
          class="player-range"
          :style="volumeStyle"
          :value="volume"
          @input="onVolume"
        />
      </div>

      <button
        type="button"
        class="flex h-9 w-9 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
        :disabled="!currentTrack"
      >
        <i class="pi pi-list" />
      </button>
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
console.log(currentTrack)
const currentTimeLabel = computed(() => formatTime(currentTime.value))

const durationLabel = computed(() => {
  return formatTime(duration.value || currentTrack.value?.durationSeconds || 0)
})

const progressStyle = computed(() => {
  return {
    '--range-progress': `${Number(progressPercent.value || 0)}%`,
  }
})

const volumeStyle = computed(() => {
  const value = muted.value ? 0 : Number(volume.value || 0) * 100

  return {
    '--range-progress': `${value}%`,
  }
})

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
  console.log(currentTrack)

  const target = event.target as HTMLInputElement
  setVolume(Number(target.value))
}
</script>

<style scoped>
.player-range {
  --range-progress: 0%;
  width: 100%;
  height: 18px;
  cursor: pointer;
  appearance: none;
  background: transparent;
}

.player-range:disabled {
  cursor: not-allowed;
  opacity: 0.45;
}

.player-range::-webkit-slider-runnable-track {
  height: 4px;
  border-radius: 999px;
  background: linear-gradient(
    to right,
    #1db954 0%,
    #1db954 var(--range-progress),
    rgba(255, 255, 255, 0.18) var(--range-progress),
    rgba(255, 255, 255, 0.18) 100%
  );
}

.player-range::-webkit-slider-thumb {
  width: 12px;
  height: 12px;
  margin-top: -4px;
  border-radius: 999px;
  appearance: none;
  background: #fff;
  opacity: 0;
  transition:
    opacity 160ms ease,
    transform 160ms ease;
}

.player-range:hover::-webkit-slider-thumb {
  opacity: 1;
}

.player-range:active::-webkit-slider-thumb {
  transform: scale(1.15);
}

.player-range::-moz-range-track {
  height: 4px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.18);
}

.player-range::-moz-range-progress {
  height: 4px;
  border-radius: 999px;
  background: #1db954;
}

.player-range::-moz-range-thumb {
  width: 12px;
  height: 12px;
  border: 0;
  border-radius: 999px;
  background: #fff;
  opacity: 0;
  transition:
    opacity 160ms ease,
    transform 160ms ease;
}

.player-range:hover::-moz-range-thumb {
  opacity: 1;
}

@keyframes art-pop {
  from {
    transform: scale(0.96);
    opacity: 0.7;
  }

  to {
    transform: scale(1);
    opacity: 1;
  }
}

.animate-art-pop {
  animation: art-pop 220ms ease-out;
}
</style>
