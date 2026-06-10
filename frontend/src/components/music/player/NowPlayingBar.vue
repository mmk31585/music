<template>
  <Transition name="bar-slide">
    <div
      v-if="currentTrack"
      ref="barRef"
      role="contentinfo"
      aria-label="Music player"
      class="fixed inset-x-0 bottom-0 z-50 select-none"
    >
      <!-- Glass-darker background base -->
      <div class="absolute inset-0 bg-black/75 backdrop-blur-2xl border-t border-white/5" />
      <div
        class="absolute inset-0 transition-opacity duration-700"
        :style="{ background: bgGradient, opacity: 0.7 }"
      />

      <div class="relative">
        <!-- Thin progress bar -->
        <div class="h-0.5 bg-white/10">
          <div
            class="h-full transition-all duration-150 ease-linear"
            :style="{ width: `${progressPercent}%`, backgroundColor: accentColor, boxShadow: isPlaying ? `0 0 8px ${accentColor}` : 'none' }"
          />
        </div>

        <!-- Main content -->
        <div class="flex h-16 md:h-[72px] items-center px-3 md:px-4">
          <div class="flex w-full items-center gap-2 md:gap-0">
            <!-- LEFT ZONE: Album art + track info (30%) -->
            <div class="flex items-center gap-3 min-w-0 md:w-[30%] md:pr-4">
              <!-- Album art -->
              <div
                class="relative shrink-0 cursor-pointer"
                @click="onTrackInfoClick"
              >
                <div class="h-12 w-12 overflow-hidden rounded-lg shadow-[0_4px_12px_rgba(0,0,0,0.5)] ring-1 ring-white/5">
                  <img
                    v-if="currentTrack?.coverUrl"
                    :src="currentTrack.coverUrl"
                    :alt="currentTrack.title"
                    class="h-full w-full object-cover"
                    loading="lazy"
                    @error="onImgError"
                  />
                  <div v-else class="flex h-full w-full items-center justify-center bg-white/10">
                    <i class="pi pi-headphones text-slate-500" />
                  </div>
                </div>
              </div>

              <!-- Track info (desktop) -->
              <div class="hidden md:block min-w-0 flex-1">
                <p class="truncate text-sm font-semibold text-white">
                  {{ currentTrack?.title || 'No track playing' }}
                </p>
                <div class="flex items-center gap-1.5">
                  <p class="truncate text-xs text-white/50">
                    {{ currentTrack?.artistName || '' }}
                  </p>
                  <button
                    type="button"
                    class="shrink-0 text-white/35 hover:text-white/70 transition-colors"
                    :class="liked ? '!text-[#1db954]' : ''"
                    @click.stop="toggleLike"
                    :aria-label="liked ? 'Unlike' : 'Like'"
                  >
                    <i :class="liked ? 'pi pi-heart-fill' : 'pi pi-heart'" class="text-xs" />
                  </button>
                </div>
              </div>

              <!-- Track info (mobile) -->
              <div
                class="md:hidden min-w-0 flex-1"
                @click="$emit('toggle-mobile-sheet')"
              >
                <p class="truncate text-sm font-semibold text-white">
                  {{ currentTrack?.title || '' }}
                </p>
                <p class="truncate text-xs text-white/50">
                  {{ currentTrack?.artistName || '' }}
                </p>
              </div>

              <MiniEqualizer :is-playing="isPlaying" class="hidden md:flex shrink-0" />
            </div>

            <!-- CENTER ZONE: Controls + seekbar (desktop only, 40%) -->
            <div class="hidden md:flex md:w-[40%] flex-col items-center gap-0.5 px-2">
              <div class="flex items-center gap-5">
                <button
                  type="button"
                  aria-label="Previous track"
                  class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 hover:text-white transition-all disabled:opacity-30 active:scale-90"
                  :disabled="!hasPrevious"
                  @click="playPrevious"
                >
                  <i class="pi pi-step-backward text-sm" />
                </button>

                <button
                  type="button"
                  :aria-label="isPlaying ? 'Pause' : 'Play'"
                  class="flex h-10 w-10 items-center justify-center rounded-full bg-[#1db954] text-white shadow-lg transition-all duration-150 active:scale-95 disabled:opacity-40 hover:brightness-110"
                  :disabled="!currentTrack || isLoadingTrack"
                  @click="togglePlayPause"
                >
                  <i v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-sm" />
                  <i v-else :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" class="text-sm" />
                </button>

                <button
                  type="button"
                  aria-label="Next track"
                  class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 hover:text-white transition-all disabled:opacity-30 active:scale-90"
                  :disabled="!hasNext"
                  @click="playNext"
                >
                  <i class="pi pi-step-forward text-sm" />
                </button>
              </div>

              <div class="flex w-full max-w-[420px] items-center gap-2">
                <span class="w-8 text-right text-[11px] text-white/35 tabular-nums">{{ currentTimeLabel }}</span>
                <div class="relative flex-1 group/seek py-1">
                  <input
                    type="range"
                    min="0"
                    max="100"
                    step="0.1"
                    class="player-range w-full"
                    :style="progressStyle"
                    :value="progressPercent"
                    :disabled="!currentTrack"
                    @input="onSeek"
                  />
                </div>
                <span class="w-8 text-[11px] text-white/35 tabular-nums">{{ durationLabel }}</span>
              </div>
            </div>

            <!-- RIGHT ZONE: Secondary controls (desktop only, 30%) -->
            <div class="hidden md:flex md:w-[30%] items-center justify-end gap-1 pl-4">
              <button
                type="button"
                aria-label="Shuffle"
                class="flex h-9 w-9 items-center justify-center rounded-full hover:bg-white/10 transition-all"
                :class="shuffleMode ? 'text-[#1db954]' : 'text-white/40'"
                :disabled="!currentTrack"
                @click="toggleShuffle"
              >
                <i class="pi pi-sort-alt text-sm" />
              </button>

              <button
                type="button"
                :aria-label="repeatTitle"
                class="relative flex h-9 w-9 items-center justify-center rounded-full hover:bg-white/10 transition-all"
                :class="repeatMode !== 'off' ? 'text-[#1db954]' : 'text-white/40'"
                :disabled="!currentTrack"
                @click="toggleRepeat"
              >
                <i class="pi pi-refresh text-sm" />
                <span
                  v-if="repeatMode === 'one'"
                  class="absolute top-0.5 right-0.5 flex h-3 w-3 items-center justify-center rounded-full bg-[#1db954] text-[7px] font-bold text-black"
                >1</span>
              </button>

              <div class="flex items-center gap-1 ml-1">
                <button
                  type="button"
                  :aria-label="muted ? 'Unmute' : 'Mute'"
                  class="flex h-9 w-9 items-center justify-center rounded-full hover:bg-white/10 transition-all text-white/40 hover:text-white"
                  @click="toggleMute"
                >
                  <i :class="volumeIcon" class="text-sm" />
                </button>
                <input
                  type="range"
                  min="0"
                  max="1"
                  step="0.01"
                  class="player-range volume-range w-20"
                  :style="volumeStyle"
                  :value="volume"
                  @input="onVolume"
                />
              </div>

              <button
                type="button"
                aria-label="Lyrics"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:bg-white/10 hover:text-white transition-all"
                :disabled="!currentTrack"
                @click="$emit('toggle-lyrics')"
              >
                <i class="pi pi-align-left text-sm" />
              </button>

              <button
                type="button"
                aria-label="Expand player"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:bg-white/10 hover:text-white transition-all"
                :disabled="!currentTrack"
                @click="$emit('toggle-fullscreen')"
              >
                <i class="pi pi-window-maximize text-sm" />
              </button>

              <div class="overflow-menu-container relative">
                <button
                  type="button"
                  aria-label="More options"
                  class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:bg-white/10 hover:text-white transition-all"
                  :disabled="!currentTrack"
                  @click="showOverflow = !showOverflow"
                >
                  <i class="pi pi-ellipsis-h text-sm" />
                </button>
                <PlayerOverflowMenu
                  v-if="showOverflow"
                  @close="showOverflow = false"
                />
              </div>
            </div>

            <!-- MOBILE CONTROLS (mobile only) -->
            <div class="flex items-center gap-2 md:hidden ml-auto">
              <button
                type="button"
                aria-label="Previous track"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 disabled:opacity-30"
                :disabled="!hasPrevious"
                @click="playPrevious"
              >
                <i class="pi pi-step-backward text-sm" />
              </button>
              <button
                type="button"
                :aria-label="isPlaying ? 'Pause' : 'Play'"
                class="flex h-9 w-9 items-center justify-center rounded-full bg-[#1db954] text-white shadow-lg disabled:opacity-40 active:scale-90 transition-all"
                :disabled="!currentTrack"
                @click="togglePlayPause"
              >
                <i :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" class="text-sm" />
              </button>
              <button
                type="button"
                aria-label="Next track"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 disabled:opacity-30"
                :disabled="!hasNext"
                @click="playNext"
              >
                <i class="pi pi-step-forward text-sm" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch, nextTick } from 'vue'
import { usePlayerControls } from '@/composables/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import { usePlayerShortcuts } from '@/composables/useShortcuts'
import { useReactionsApi } from '@/services/api/reactions'
import { onImgError } from '@/utils/helpers'
import MiniEqualizer from './MiniEqualizer.vue'
import PlayerOverflowMenu from './PlayerOverflowMenu.vue'

const emit = defineEmits<{
  'toggle-queue': []
  'toggle-fullscreen': []
  'toggle-lyrics': []
  'toggle-mobile-sheet': []
}>()

const liked = ref(false)
const reactionsApi = useReactionsApi()

const showOverflow = ref(false)

async function toggleLike() {
  const track = currentTrack.value
  if (!track?.id) return
  liked.value = !liked.value
  try {
    if (liked.value) {
      await reactionsApi.react({ target_id: track.id, target_type: 'track', type: 'like' })
    } else {
      await reactionsApi.removeReaction('track', track.id)
    }
  } catch {
    liked.value = !liked.value
  }
}

function onTrackInfoClick() {
  emit('toggle-fullscreen')
}

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
  shuffleMode,
  repeatMode,
  playIcon,
  volumeIcon,
  togglePlayPause,
  toggleMute,
  setVolume,
  seekPercent,
  playNext,
  playPrevious,
  toggleShuffle,
  toggleRepeat,
  sleepTimerMinutes,
  setSleepTimer,
  clearSleepTimer,
  crossfadeDuration,
} = usePlayerControls()

// Global keyboard shortcuts
usePlayerShortcuts({
  togglePlay: togglePlayPause,
  next: playNext,
  previous: playPrevious,
  toggleMute,
  toggleShuffle,
  toggleRepeat,
  seekBackward: () => seekPercent(Math.max(0, progressPercent.value - 5)),
  seekForward: () => seekPercent(Math.min(100, progressPercent.value + 5)),
  volumeUp: () => setVolume(Math.min(1, volume.value + 0.05)),
  volumeDown: () => setVolume(Math.max(0, volume.value - 0.05)),
})

// Dynamic album colors
const coverUrl = computed(() => currentTrack.value?.coverUrl || null)
const { palette } = useAlbumColors(coverUrl)

const accentColor = computed(() => palette.value.vibrant || '#1db954')

const bgGradient = computed(() => {
  const p = palette.value
  if (!coverUrl.value) return 'rgba(0,0,0,0.85)'
  return `linear-gradient(180deg, ${p.dark}dd 0%, rgba(0,0,0,0.92) 100%)`
})

// Bar ref for cleanup
const barRef = ref<HTMLElement | null>(null)

const repeatTitle = computed(() => {
  if (repeatMode.value === 'off') return 'Repeat: off'
  if (repeatMode.value === 'all') return 'Repeat: all'
  return 'Repeat: one'
})

const currentTimeLabel = computed(() => formatTime(currentTime.value))

const durationLabel = computed(() => formatTime(duration.value || currentTrack.value?.durationSeconds || 0))

const progressStyle = computed(() => ({
  '--range-progress': `${Number(progressPercent.value || 0)}%`,
  '--accent-color': accentColor.value,
}))

const volumeStyle = computed(() => {
  const value = muted.value ? 0 : Number(volume.value || 0) * 100
  return { '--range-progress': `${value}%` }
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
  const target = event.target as HTMLInputElement
  setVolume(Number(target.value))
}

function onOverflowClickOutside(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (!target.closest('.overflow-menu-container')) {
    showOverflow.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', onOverflowClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onOverflowClickOutside)
})
</script>

<style scoped>
.bar-slide-enter-active {
  transition: transform 400ms cubic-bezier(0.19, 1, 0.22, 1);
}
.bar-slide-leave-active {
  transition: transform 300ms ease-in;
}
.bar-slide-enter-from {
  transform: translateY(100%);
}
.bar-slide-leave-to {
  transform: translateY(100%);
}

@media (prefers-reduced-motion: reduce) {
  .bar-slide-enter-active,
  .bar-slide-leave-active {
    transition: none;
  }
  .bar-slide-enter-from,
  .bar-slide-leave-to {
    transform: none;
  }
}
</style>
