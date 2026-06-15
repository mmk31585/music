<template>
  <Transition name="bar-slide">
    <div
      v-if="currentTrack"
      ref="barRef"
      role="contentinfo"
      aria-label="Music player"
      class="glass-darker fixed inset-x-0 bottom-0 z-50 select-none border-t border-white/5"
    >
      <div
        v-if="playbackError"
        class="flex items-center justify-center gap-2 bg-red-500/10 px-4 py-1 text-xs text-red-400 border-b border-red-500/10"
      >
        <i class="pi pi-exclamation-circle text-xs" />
        <span>{{ playbackError }}</span>
      </div>
      <!-- Mobile progress line -->
      <div class="md:hidden h-0.5 bg-white/10">
        <div
          class="h-full transition-all duration-150 ease-linear"
          :style="{ width: `${progressPercent}%`, backgroundColor: accentColor, boxShadow: isPlaying ? `0 0 8px ${accentColor}` : 'none' }"
        />
      </div>

      <div class="flex h-16 md:h-[72px] items-center px-3 md:px-4">
        <!-- LEFT: Cover + Info -->
        <div class="flex items-center gap-3 min-w-0 md:w-[30%] md:pr-4">
          <div class="relative shrink-0 cursor-pointer group/art" @click="onTrackInfoClick">
            <div class="h-12 w-12 overflow-hidden rounded-xl shadow-lg ring-1 ring-white/10 transition-all duration-300 group-hover/art:ring-[#1db954]/50 group-hover/art:shadow-[#1db954]/20">
              <img
                v-if="currentTrack?.coverUrl"
                :src="currentTrack.coverUrl"
                :alt="currentTrack.title"
                class="h-full w-full object-cover"
                loading="lazy"
                @error="onImgError"
              />
              <div v-else class="flex h-full w-full items-center justify-center bg-gradient-to-br from-[#1db954]/20 to-purple-500/20">
                <i class="pi pi-headphones text-slate-500 text-lg" />
              </div>
            </div>
            <div class="absolute inset-0 flex items-center justify-center rounded-xl bg-black/40 opacity-0 transition-opacity duration-200 group-hover/art:opacity-100">
              <i class="pi pi-expand text-white text-xs" />
            </div>
          </div>

          <div class="min-w-0 flex-1">
            <p class="truncate text-sm font-bold text-white leading-tight">
              {{ currentTrack?.title || 'No track playing' }}
            </p>
            <p class="truncate text-xs text-[rgba(255,255,255,0.55)] leading-tight mt-0.5">
              {{ currentTrack?.artistName || '' }}
            </p>
          </div>

          <button
            type="button"
            class="hidden md:flex shrink-0 h-8 w-8 items-center justify-center rounded-full transition-all duration-200"
            :class="liked ? 'text-[#1db954] bg-[#1db954]/10' : 'text-[rgba(255,255,255,0.35)] hover:text-white hover:bg-white/10'"
            @click.stop="toggleLike"
            :aria-label="liked ? 'Unlike' : 'Like'"
          >
            <i :class="liked ? 'pi pi-heart-fill' : 'pi pi-heart'" class="text-xs" />
          </button>
        </div>

        <!-- CENTER: Controls + Seek -->
        <div class="hidden md:flex md:w-[40%] flex-col items-center gap-0.5 px-4">
          <div class="flex items-center gap-5">
            <button
              type="button"
              aria-label="Previous track"
              class="flex h-8 w-8 items-center justify-center rounded-full text-white/40 hover:text-white transition-all disabled:opacity-25 active:scale-90 hover:bg-white/5"
              :disabled="!hasPrevious"
              @click="playPrevious"
            >
              <i class="pi pi-step-backward text-sm" />
            </button>

            <GuestPlayGate action="play" @proceed="togglePlayPause" #default="{ proceed }">
              <button
                type="button"
                :aria-label="isPlaying ? 'Pause' : 'Play'"
                class="flex h-9 w-9 items-center justify-center rounded-full bg-white text-black shadow-lg transition-all duration-150 active:scale-90 disabled:opacity-40 hover:scale-105"
                :disabled="!currentTrack || isLoadingTrack"
                @click="isPlaying ? togglePlayPause() : proceed()"
              >
                <i v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-sm" />
                <i v-else :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" class="text-sm" />
              </button>
            </GuestPlayGate>

            <button
              type="button"
              aria-label="Next track"
              class="flex h-8 w-8 items-center justify-center rounded-full text-white/40 hover:text-white transition-all disabled:opacity-25 active:scale-90 hover:bg-white/5"
              :disabled="!hasNext"
              @click="playNext"
            >
              <i class="pi pi-step-forward text-sm" />
            </button>
          </div>

          <div class="flex w-full max-w-[480px] items-center gap-2">
            <span class="w-8 text-right text-[11px] text-[rgba(255,255,255,0.35)] tabular-nums font-medium">{{ currentTimeLabel }}</span>
            <div class="relative flex-1 group/seekbar">
              <input
                type="range"
                min="0"
                :max="Math.max(1, duration || currentTrack?.durationSeconds || 100)"
                step="0.1"
                class="player-range w-full"
                :style="progressStyle"
                :value="currentTime"
                :disabled="!currentTrack"
                @input="onSeek"
              />
            </div>
            <span class="w-8 text-[11px] text-[rgba(255,255,255,0.35)] tabular-nums font-medium">{{ durationLabel }}</span>
          </div>
        </div>

        <!-- RIGHT: Secondary controls -->
        <div class="hidden md:flex md:w-[30%] items-center justify-end gap-1 pl-4">
          <button
            type="button"
            aria-label="Shuffle"
            class="flex h-8 w-8 items-center justify-center rounded-full transition-all duration-200"
            :class="shuffleMode ? 'text-[#1db954] bg-[#1db954]/10' : 'text-[rgba(255,255,255,0.35)] hover:text-white hover:bg-white/10'"
            :disabled="!currentTrack"
            @click="toggleShuffle"
          >
            <i class="pi pi-sort-alt text-sm" />
          </button>

          <button
            type="button"
            :aria-label="repeatTitle"
            class="relative flex h-8 w-8 items-center justify-center rounded-full transition-all duration-200"
            :class="repeatMode !== 'off' ? 'text-[#1db954] bg-[#1db954]/10' : 'text-[rgba(255,255,255,0.35)] hover:text-white hover:bg-white/10'"
            :disabled="!currentTrack"
            @click="toggleRepeat"
          >
            <i class="pi pi-refresh text-sm" />
            <span
              v-if="repeatMode === 'one'"
              class="absolute -top-0.5 -right-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-[#1db954] text-[7px] font-bold text-black"
            >1</span>
          </button>

          <div class="mx-1 h-6 w-px bg-white/5" />

          <button
            type="button"
            aria-label="Queue"
            class="flex h-8 w-8 items-center justify-center rounded-full text-[rgba(255,255,255,0.35)] hover:text-white hover:bg-white/10 transition-all duration-200"
            :disabled="!currentTrack"
            @click="$emit('toggle-queue')"
          >
            <i class="pi pi-list text-sm" />
          </button>

          <button
            type="button"
            aria-label="Lyrics"
            class="flex h-8 w-8 items-center justify-center rounded-full text-[rgba(255,255,255,0.35)] hover:text-white hover:bg-white/10 transition-all duration-200"
            :disabled="!currentTrack"
            @click="$emit('toggle-lyrics')"
          >
            <i class="pi pi-align-left text-sm" />
          </button>

          <div class="flex items-center gap-1">
            <button
              type="button"
              :aria-label="muted ? 'Unmute' : 'Mute'"
              class="flex h-8 w-8 items-center justify-center rounded-full transition-all duration-200 text-[rgba(255,255,255,0.35)] hover:text-white hover:bg-white/10"
              @click="toggleMute"
            >
              <i :class="volumeIcon" class="text-sm" />
            </button>
            <div class="w-16">
              <input
                type="range"
                min="0"
                max="1"
                step="0.01"
                class="player-range volume-range w-full"
                :style="volumeStyle"
                :value="muted ? 0 : volume"
                @input="onVolume"
              />
            </div>
          </div>

          <button
            type="button"
            aria-label="Fullscreen player"
            class="flex h-8 w-8 items-center justify-center rounded-full text-[rgba(255,255,255,0.35)] hover:text-white hover:bg-white/10 transition-all duration-200"
            :disabled="!currentTrack"
            @click="$emit('toggle-fullscreen')"
          >
            <i class="pi pi-arrow-expand text-sm" />
          </button>

          <div class="overflow-menu-container relative">
            <button
              type="button"
              aria-label="More options"
              :aria-expanded="showOverflow"
              class="flex h-8 w-8 items-center justify-center rounded-full text-[rgba(255,255,255,0.35)] hover:text-white hover:bg-white/10 transition-all duration-200"
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

        <!-- MOBILE CONTROLS -->
        <div class="flex items-center gap-1 md:hidden ml-auto">
          <button
            type="button"
            aria-label="Previous track"
            class="flex h-8 w-8 items-center justify-center rounded-full text-white/40 disabled:opacity-25 active:scale-90"
            :disabled="!hasPrevious"
            @click="playPrevious"
          >
            <i class="pi pi-step-backward text-sm" />
          </button>

          <GuestPlayGate action="play" @proceed="togglePlayPause" #default="{ proceed }">
            <button
              type="button"
              :aria-label="isPlaying ? 'Pause' : 'Play'"
              class="flex h-8 w-8 items-center justify-center rounded-full bg-white text-black shadow-lg disabled:opacity-40 transition-all duration-150 active:scale-90"
              :disabled="!currentTrack"
              @click="isPlaying ? togglePlayPause() : proceed()"
            >
              <i v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-sm" />
              <i v-else :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" class="text-sm" />
            </button>
          </GuestPlayGate>

          <button
            type="button"
            aria-label="Next track"
            class="flex h-8 w-8 items-center justify-center rounded-full text-white/40 disabled:opacity-25 active:scale-90"
            :disabled="!hasNext"
            @click="playNext"
          >
            <i class="pi pi-step-forward text-sm" />
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { usePlayer } from '@/composables/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import { usePlayerShortcuts } from '@/composables/useShortcuts'
import { useReactionsApi } from '@/services/api/reactions'
import { onImgError } from '@/utils/helpers'
import MiniEqualizer from './MiniEqualizer.vue'
import PlayerOverflowMenu from './PlayerOverflowMenu.vue'
import GuestPlayGate from '@/components/common/GuestPlayGate.vue'

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

const player = usePlayer()

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
  error: playbackError,
  hasNext,
  hasPrevious,
  shuffleMode,
  repeatMode,
  volumeIcon,
  togglePlayPause,
  toggleMute,
  setVolume,
  seekPercent,
  playNext,
  playPrevious,
  toggleShuffle,
  toggleRepeat,
} = player

const showSleepPopover = ref(false)

const coverUrl = computed(() => currentTrack.value?.coverUrl || null)
const { palette } = useAlbumColors(coverUrl)

const accentColor = computed(() => palette.value.vibrant || '#1db954')

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
