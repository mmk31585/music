<template>
  <div class="fixed right-0 bottom-0 left-0 z-30 border-t border-white/10 bg-[#121212] px-3 py-3">
    <div class="grid items-center gap-3 md:grid-cols-[1fr_1.4fr_1fr]">
      <div class="flex min-w-0 items-center gap-3">
        <div
          class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-white/10 text-slate-300"
        >
          <img
            v-if="track?.cover_url"
            :src="track.cover_url"
            :alt="track.title"
            class="h-full w-full rounded-lg object-cover"
          />
          <i v-else class="pi pi-music" />
        </div>

        <div class="min-w-0">
          <p class="truncate text-sm font-medium text-white">
            {{ track?.title ?? 'No track selected' }}
          </p>
          <p class="truncate text-xs text-slate-400">
            {{ track?.artist_name || (track ? 'Unknown artist' : 'Choose something to play') }}
          </p>
          <p v-if="player.error" class="truncate text-xs text-red-300">{{ player.error }}</p>
        </div>
      </div>

      <div class="min-w-0">
        <div class="flex items-center justify-center gap-2">
          <button
            class="player-icon-button hidden md:inline-flex"
            :class="player.shuffle ? 'text-emerald-300' : 'text-slate-300'"
            aria-label="Shuffle"
            @click="player.toggleShuffle"
          >
            <i class="pi pi-sort-alt" />
          </button>

          <button
            class="player-icon-button"
            :disabled="!player.hasPrevious"
            aria-label="Previous track"
            @click="player.playPrevious"
          >
            <i class="pi pi-step-backward" />
          </button>

          <button
            class="flex h-11 w-11 items-center justify-center rounded-full bg-[#1db954] text-black transition hover:scale-105 disabled:cursor-not-allowed disabled:bg-white/10 disabled:text-slate-500"
            :disabled="!player.canPlayCurrent || player.isBuffering"
            :aria-label="player.isPlaying ? 'Pause' : 'Play'"
            @click="player.togglePlay"
          >
            <i v-if="player.isBuffering" class="pi pi-spin pi-spinner" />
            <i v-else :class="player.isPlaying ? 'pi pi-pause' : 'pi pi-play'" />
          </button>

          <button
            class="player-icon-button"
            :disabled="!player.hasNext"
            aria-label="Next track"
            @click="player.playNext"
          >
            <i class="pi pi-step-forward" />
          </button>

          <button
            class="player-icon-button hidden md:inline-flex"
            :class="player.repeatMode !== 'off' ? 'text-emerald-300' : 'text-slate-300'"
            :aria-label="`Repeat ${player.repeatMode}`"
            @click="player.cycleRepeat"
          >
            <span class="relative">
              <i class="pi pi-refresh" />
              <span
                v-if="player.repeatMode === 'one'"
                class="absolute -right-2 -bottom-2 text-[10px] font-bold"
              >
                1
              </span>
            </span>
          </button>
        </div>

        <div class="mt-2 flex items-center gap-2 text-[11px] text-slate-400">
          <span class="w-10 text-right">{{ formatTime(player.currentTime) }}</span>
          <input
            class="player-range"
            type="range"
            min="0"
            :max="seekMax"
            step="1"
            :value="player.currentTime"
            :disabled="!track"
            aria-label="Track progress"
            @input="onSeek"
          />
          <span class="w-10">{{ formatTime(seekMax) }}</span>
        </div>
      </div>

      <div class="hidden items-center justify-end gap-3 md:flex">
        <span class="text-xs text-slate-400">{{ player.queue.length }} in queue</span>
        <i class="pi pi-volume-up text-slate-400" />
        <input
          class="player-range max-w-28"
          type="range"
          min="0"
          max="1"
          step="0.01"
          :value="player.volume"
          aria-label="Volume"
          @input="onVolume"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { usePlayerStore } from '@/stores'

const player = usePlayerStore()
const track = computed(() => player.currentTrack)
const seekMax = computed(() => player.duration || track.value?.duration_seconds || 0)

onMounted(() => {
  player.restore()
})

function formatTime(value?: number | null) {
  if (!value || !Number.isFinite(value)) return '0:00'
  const total = Math.max(0, Math.floor(value))
  const mins = Math.floor(total / 60)
  const secs = total % 60
  return `${mins}:${String(secs).padStart(2, '0')}`
}

function onSeek(event: Event) {
  const value = Number((event.target as HTMLInputElement).value)
  player.seek(value)
}

function onVolume(event: Event) {
  const value = Number((event.target as HTMLInputElement).value)
  player.setVolume(value)
}
</script>

<style scoped>
.player-icon-button {
  display: inline-flex;
  width: 2.25rem;
  height: 2.25rem;
  align-items: center;
  justify-content: center;
  border-radius: 9999px;
  color: rgb(203 213 225);
  transition:
    color 150ms ease,
    background-color 150ms ease;
}

.player-icon-button:hover:not(:disabled) {
  background: rgb(255 255 255 / 0.1);
  color: #fff;
}

.player-icon-button:disabled {
  cursor: not-allowed;
  color: rgb(71 85 105);
}

.player-range {
  width: 100%;
  height: 0.25rem;
  cursor: pointer;
  appearance: none;
  border-radius: 9999px;
  background: rgb(255 255 255 / 0.15);
  accent-color: #1db954;
}

.player-range:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}
</style>
