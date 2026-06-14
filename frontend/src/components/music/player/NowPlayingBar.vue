<template>
  <Transition name="bar-slide">
    <div
      v-if="currentTrack"
      ref="barRef"
      role="contentinfo"
      aria-label="Music player"
      class="glass-darker fixed inset-x-0 bottom-0 z-50 select-none"
    >
      <!-- Mobile progress line (top edge indicator) -->
      <div class="md:hidden h-0.5 bg-white/10">
        <div
          class="h-full transition-all duration-150 ease-linear"
          :style="{ width: `${progressPercent}%`, backgroundColor: accentColor, boxShadow: isPlaying ? `0 0 8px ${accentColor}` : 'none' }"
        />
      </div>

      <div class="flex h-16 md:h-[72px] items-center px-3 md:px-4">
        <!-- LEFT ZONE: Album art + track info (30% desktop) -->
        <div class="flex items-center gap-3 min-w-0 md:w-[30%] md:pr-4">
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

          <!-- Track info (desktop with heart + equalizer) -->
          <div class="hidden md:flex min-w-0 flex-1 items-center gap-2">
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-semibold text-white">
                {{ currentTrack?.title || 'No track playing' }}
              </p>
              <div class="flex items-center gap-1.5">
                <p class="truncate text-xs text-[rgba(255,255,255,0.6)]">
                  {{ currentTrack?.artistName || '' }}
                </p>
                <button
                  type="button"
                  class="shrink-0 text-[rgba(255,255,255,0.35)] hover:text-white/70 transition-colors"
                  :class="liked ? '!text-[#1db954]' : ''"
                  @click.stop="toggleLike"
                  :aria-label="liked ? 'Unlike' : 'Like'"
                >
                  <i :class="liked ? 'pi pi-heart-fill' : 'pi pi-heart'" class="text-xs" />
                </button>
              </div>
            </div>
            <MiniEqualizer :is-playing="isPlaying" class="shrink-0" />
          </div>

          <!-- Track info (mobile) -->
          <div
            class="md:hidden min-w-0 flex-1"
            @click="$emit('toggle-mobile-sheet')"
          >
            <p class="truncate text-sm font-semibold text-white">
              {{ currentTrack?.title || '' }}
            </p>
            <p class="truncate text-xs text-[rgba(255,255,255,0.6)]">
              {{ currentTrack?.artistName || '' }}
            </p>
          </div>
        </div>

        <!-- CENTER ZONE: Controls + seekbar (40% desktop, hidden mobile) -->
        <div class="hidden md:flex md:w-[40%] flex-col items-center gap-0.5 px-4">
          <div class="flex items-center gap-6">
            <button
              type="button"
              aria-label="Previous track"
              class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 hover:text-white transition-all disabled:opacity-30 active:scale-90"
              :disabled="!hasPrevious"
              @click="playPrevious"
            >
              <i class="pi pi-step-backward text-sm" />
            </button>

            <GuestPlayGate action="play" @proceed="togglePlayPause" #default="{ proceed }">
              <button
                type="button"
                :aria-label="isPlaying ? 'Pause' : 'Play'"
                class="flex h-10 w-10 items-center justify-center rounded-full bg-[#1db954] text-white shadow-lg transition-all duration-150 active:scale-95 disabled:opacity-40 hover:brightness-110"
                :style="{ transitionTimingFunction: 'var(--ease-spring)' }"
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
              class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 hover:text-white transition-all disabled:opacity-30 active:scale-90"
              :disabled="!hasNext"
              @click="playNext"
            >
              <i class="pi pi-step-forward text-sm" />
            </button>
          </div>

          <div class="flex w-full max-w-[420px] items-center gap-2">
            <span class="w-8 text-right text-[11px] text-[rgba(255,255,255,0.35)] tabular-nums">{{ currentTimeLabel }}</span>
            <div class="relative flex-1 group/seekbar">
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
            <span class="w-8 text-[11px] text-[rgba(255,255,255,0.35)] tabular-nums">{{ durationLabel }}</span>
          </div>
        </div>

        <!-- RIGHT ZONE: Secondary controls (30% desktop, hidden mobile) -->
        <div class="hidden md:flex md:w-[30%] items-center justify-end gap-2 pl-4">
          <button
            type="button"
            aria-label="Shuffle"
            class="flex h-9 w-9 items-center justify-center rounded-full hover:bg-white/10 transition-all"
            :class="shuffleMode ? 'text-[#1db954]' : 'text-[rgba(255,255,255,0.35)]'"
            :disabled="!currentTrack"
            @click="toggleShuffle"
          >
            <i class="pi pi-sort-alt text-sm" />
          </button>

          <button
            type="button"
            :aria-label="repeatTitle"
            class="relative flex h-9 w-9 items-center justify-center rounded-full hover:bg-white/10 transition-all"
            :class="repeatMode !== 'off' ? 'text-[#1db954]' : 'text-[rgba(255,255,255,0.35)]'"
            :disabled="!currentTrack"
            @click="toggleRepeat"
          >
            <i class="pi pi-refresh text-sm" />
            <span
              v-if="repeatMode === 'one'"
              class="absolute -top-0.5 -right-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-[#1db954] text-[7px] font-bold text-black"
            >1</span>
          </button>

          <div class="flex items-center gap-1">
            <button
              type="button"
              :aria-label="muted ? 'Unmute' : 'Mute'"
              class="flex h-9 w-9 items-center justify-center rounded-full hover:bg-white/10 transition-all text-[rgba(255,255,255,0.35)] hover:text-white"
              @click="toggleMute"
            >
              <i :class="volumeIcon" class="text-sm" />
            </button>
            <div class="w-20">
              <input
                type="range"
                min="0"
                max="1"
                step="0.01"
                class="player-range volume-range w-full"
                :style="volumeStyle"
                :value="volume"
                @input="onVolume"
              />
            </div>
          </div>

          <div class="relative">
            <button
              type="button"
              v-tooltip="sleepTimerMinutes > 0 ? `تایمر فعال: ${sleepTimerMinutes}دقیقه` : 'تایمر خواب'"
              class="flex h-9 w-9 items-center justify-center rounded-full hover:bg-white/10 transition-all"
              :class="sleepTimerMinutes > 0 ? 'text-[#1db954]' : 'text-[rgba(255,255,255,0.35)]'"
              :disabled="!currentTrack"
              @click="showSleepPopover = !showSleepPopover"
            >
              <i class="pi pi-moon text-sm" />
            </button>
            <Transition name="fade">
              <div
                v-if="showSleepPopover"
                class="absolute bottom-full right-0 mb-2 z-50"
              >
                <div class="backdrop-blur-xl bg-white/5 border border-white/10 rounded-2xl shadow-2xl flex flex-wrap gap-1.5 p-3 min-w-[200px]">
                  <button
                    v-for="opt in sleepOptions"
                    :key="opt.value"
                    type="button"
                    class="rounded-full px-3 py-1.5 text-xs font-medium transition"
                    :class="sleepTimerMinutes === opt.value ? 'bg-[#1db954] text-black' : 'bg-white/10 text-white/50 hover:bg-white/20 hover:text-white/80'"
                    @click="setTimer(opt.value)"
                  >
                    {{ opt.label }}
                  </button>
                </div>
              </div>
            </Transition>
          </div>

          <button
            type="button"
            v-tooltip="'کراس‌فید'"
            class="flex h-9 w-9 items-center justify-center rounded-full hover:bg-white/10 transition-all"
            :class="crossfadeDuration > 0 ? 'text-[#1db954]' : 'text-[rgba(255,255,255,0.35)]'"
            :disabled="!currentTrack"
            @click="toggleCrossfade"
          >
            <i class="pi pi-arrows-h text-sm" />
          </button>

          <div class="overflow-menu-container relative">
            <button
              type="button"
              aria-label="More options"
              :aria-expanded="showOverflow"
              class="flex h-9 w-9 items-center justify-center rounded-full text-[rgba(255,255,255,0.35)] hover:bg-white/10 hover:text-white transition-all"
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
            class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 disabled:opacity-30 active:scale-90"
            :disabled="!hasPrevious"
            @click="playPrevious"
          >
            <i class="pi pi-step-backward text-sm" />
          </button>

          <GuestPlayGate action="play" @proceed="togglePlayPause" #default="{ proceed }">
            <button
              type="button"
              :aria-label="isPlaying ? 'Pause' : 'Play'"
              class="flex h-9 w-9 items-center justify-center rounded-full bg-[#1db954] text-white shadow-lg disabled:opacity-40 transition-all duration-150 active:scale-95"
              :style="{ transitionTimingFunction: 'var(--ease-spring)' }"
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
            class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 disabled:opacity-30 active:scale-90"
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
import { computed, onBeforeUnmount, onMounted, ref, watch, nextTick } from 'vue'
import { usePlayerControls, usePlayer } from '@/composables/player'
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

const showSleepPopover = ref(false)

const sleepOptions = [
  { value: 15, label: '۱۵ دقیقه' },
  { value: 30, label: '۳۰ دقیقه' },
  { value: 45, label: '۴۵ دقیقه' },
  { value: 60, label: '۶۰ دقیقه' },
  { value: 0, label: 'خاموش کردن' },
]

function setTimer(minutes: number) {
  if (minutes === 0) {
    clearSleepTimer()
  } else {
    setSleepTimer(minutes)
  }
  showSleepPopover.value = false
}

function toggleCrossfade() {
  const p = usePlayer()
  p.crossfadeDuration = p.crossfadeDuration > 0 ? 0 : 5
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

const coverUrl = computed(() => currentTrack.value?.coverUrl || null)
const { palette } = useAlbumColors(coverUrl)

const accentColor = computed(() => palette.value.vibrant || '#1db954')

const barRef = ref<HTMLElement | null>(null)

const repeatTitle = computed(() => {
  if (repeatMode === 'off') return 'Repeat: off'
  if (repeatMode === 'all') return 'Repeat: all'
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
