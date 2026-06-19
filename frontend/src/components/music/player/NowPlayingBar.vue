<template>
  <!-- DESKTOP BAR -->
  <Transition name="bar-slide">
    <div
      v-if="currentTrack"
      ref="barRef"
      role="contentinfo"
      aria-label="Music player"
      class="fixed inset-x-0 bottom-0 z-50 hidden md:block select-none"
    >
      <div
        class="relative flex flex-col overflow-hidden"
        style="background: #08080A; backdrop-filter: blur(32px); -webkit-backdrop-filter: blur(32px);"
      >
        <div
          v-if="currentTrack.coverUrl"
          class="pointer-events-none absolute inset-0 scale-110"
          aria-hidden="true"
        >
          <img
            :src="currentTrack.coverUrl"
            alt=""
            class="h-full w-full object-cover opacity-[0.35]"
            style="filter: blur(60px) saturate(1.5)"
          />
        </div>
        <div class="pointer-events-none absolute inset-0 bg-gradient-to-t from-black/80 via-black/50 to-transparent" />

        <div class="relative z-10 flex items-center gap-4 px-6 pt-3">
          <div class="flex min-w-0 w-[25%] items-center gap-3">
            <div class="relative shrink-0 cursor-pointer" @click="emit('toggle-fullscreen')">
              <div
                class="h-14 w-14 overflow-hidden rounded-[18px] shadow-[0_16px_32px_rgba(0,0,0,0.5)] ring-1 ring-white/10 transition-all duration-700"
                :class="isPlaying ? 'scale-100' : 'scale-95 opacity-80'"
              >
                <img
                  v-if="currentTrack.coverUrl"
                  :src="currentTrack.coverUrl"
                  :alt="currentTrack.title"
                  class="h-full w-full object-cover"
                  loading="lazy"
                />
                <div v-else class="flex h-full w-full items-center justify-center bg-gradient-to-br from-[#1db954]/30 to-[#a855f7]/30">
                  <i aria-hidden="true" class="pi pi-headphones text-lg text-slate-500" />
                </div>
              </div>
            </div>

            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <p class="truncate text-sm font-bold text-white leading-tight">{{ currentTrack.title }}</p>
                <button
                  type="button"
                  class="shrink-0 flex items-center justify-center transition-all hover:scale-110 active:scale-90"
                  :class="liked ? 'text-[#f472b6]' : 'text-white/30 hover:text-white/60'"
                  @click.stop="liked = !liked"
                >
                  <i aria-hidden="true" :class="liked ? 'pi pi-heart-fill' : 'pi pi-heart'" class="text-xs" />
                </button>
              </div>
              <p class="truncate text-xs text-white/50 mt-0.5 leading-tight">{{ currentTrack.artistName }}</p>
              <div v-if="isPlaying" class="flex items-center gap-1 mt-1">
                <span class="size-1 rounded-full bg-white/60 animate-bounce" style="animation-delay: 0ms" />
                <span class="size-1 rounded-full bg-white/60 animate-bounce" style="animation-delay: 150ms" />
                <span class="size-1 rounded-full bg-white/60 animate-bounce" style="animation-delay: 300ms" />
              </div>
            </div>
          </div>

          <div class="flex flex-1 items-center justify-center gap-1.5">
            <button
              type="button"
              class="flex h-9 w-9 items-center justify-center rounded-full transition-all duration-200"
              :class="shuffleMode ? 'text-[#a855f7]' : 'text-white/40 hover:text-white hover:bg-white/10'"
              :disabled="!currentTrack"
              @click="toggleShuffle"
            >
              <i aria-hidden="true" class="pi pi-sort-alt text-sm" />
            </button>

            <button
              type="button"
              class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 hover:text-white disabled:opacity-25 hover:bg-white/10 transition-all duration-200"
              :disabled="!hasPrevious"
              @click="playPrevious"
            >
              <i aria-hidden="true" class="pi pi-step-backward text-base" />
            </button>

            <GuestPlayGate action="play" @proceed="togglePlayPause" #default="{ proceed }">
              <button
                type="button"
                class="flex h-11 w-11 items-center justify-center rounded-full shadow-xl disabled:opacity-40 hover:scale-110 active:scale-95 transition-all duration-200"
                :style="{ background: progressFill }"
                :disabled="!currentTrack || isLoadingTrack"
                @click="isPlaying ? togglePlayPause() : proceed()"
              >
                <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-base text-white" />
                <i aria-hidden="true" v-else :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" class="ml-0.5 text-base text-white" />
              </button>
            </GuestPlayGate>

            <button
              type="button"
              class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 hover:text-white disabled:opacity-25 hover:bg-white/10 transition-all duration-200"
              :disabled="!hasNext"
              @click="playNext"
            >
              <i aria-hidden="true" class="pi pi-step-forward text-base" />
            </button>

            <button
              type="button"
              class="relative flex h-9 w-9 items-center justify-center rounded-full transition-all duration-200"
              :class="repeatMode !== 'off' ? 'text-[#f472b6]' : 'text-white/40 hover:text-white hover:bg-white/10'"
              :disabled="!currentTrack"
              @click="toggleRepeat"
            >
              <i aria-hidden="true" class="pi pi-refresh text-sm" />
              <span
                v-if="repeatMode === 'one'"
                class="absolute -top-0.5 -right-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-[#f472b6] text-[8px] font-bold text-black"
              >1</span>
            </button>
          </div>

          <div class="flex w-[25%] items-center justify-end gap-1">
            <button
              type="button"
              class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
              :disabled="!currentTrack"
              @click="emit('toggle-queue')"
            >
              <i aria-hidden="true" class="pi pi-list text-sm" />
            </button>

            <button
              type="button"
              class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
              :disabled="!currentTrack"
              @click="emit('toggle-lyrics')"
            >
              <i aria-hidden="true" class="pi pi-align-left text-sm" />
            </button>

            <div class="mx-1 h-6 w-px bg-white/10" />

            <button
              type="button"
              class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
              @click="toggleMute"
            >
              <i aria-hidden="true" :class="volumeIcon" class="text-sm" />
            </button>
            <div class="w-20" dir="ltr">
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

            <div class="mx-1 h-6 w-px bg-white/10" />

            <button
              type="button"
              class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
              :disabled="!currentTrack"
              @click="emit('toggle-fullscreen')"
            >
              <i aria-hidden="true" class="pi pi-arrow-up-right-and-arrow-down-left-from-center text-sm" />
            </button>

            <div class="relative overflow-menu-container">
              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
                :disabled="!currentTrack"
                @click.stop="showOverflow = !showOverflow"
              >
                <i aria-hidden="true" class="pi pi-ellipsis-h text-sm" />
              </button>
              <PlayerOverflowMenu v-if="showOverflow" @close="showOverflow = false" />
            </div>
          </div>
        </div>

        <div class="relative z-10 flex items-center gap-3 px-6 pb-3 pt-1">
          <span class="text-[11px] text-white/40 font-mono tabular-nums w-10 text-right">{{ formatTime(currentTime) }}</span>
          <div
            class="group relative flex-1 h-1.5 cursor-pointer rounded-full bg-white/10 hover:h-2 transition-all duration-150"
            dir="ltr"
            @click="onSeekClick"
          >
            <div
              class="h-full rounded-full transition-[width] duration-100"
              :style="progressStyle"
            />
            <div
              class="absolute top-1/2 -translate-y-1/2 size-3.5 rounded-full bg-white shadow-lg shadow-black/40 opacity-0 group-hover:opacity-100 scale-0 group-hover:scale-100 transition-all duration-200"
              :style="{ left: `calc(${progressPercent}% - 7px)` }"
            />
          </div>
          <span class="text-[11px] text-white/40 font-mono tabular-nums w-10">{{ formatTime(duration) }}</span>
        </div>

        <div
          v-if="playbackError"
          class="relative z-10 flex items-center justify-center gap-2 bg-red-500/10 px-5 py-1.5 text-xs text-red-400"
        >
          <i aria-hidden="true" class="pi pi-exclamation-circle text-xs" />
          <span>{{ playbackError }}</span>
        </div>
      </div>
    </div>
  </Transition>

  <!-- MOBILE BAR - floats above bottom nav -->
  <Transition name="bar-slide">
    <div
      v-if="currentTrack"
      key="mobile-bar"
      role="contentinfo"
      aria-label="Music player"
      class="fixed inset-x-0 z-[45] md:hidden select-none"
      style="bottom: calc(4.5rem + env(safe-area-inset-bottom, 0px))"
    >
      <button
        type="button"
        class="relative mx-3 flex w-[calc(100%-1.5rem)] items-center gap-3 overflow-hidden rounded-2xl text-left shadow-2xl transition-all duration-500 active:scale-[0.98]"
        :style="{ background: coverUrl ? `${palette.dark}dd` : 'rgba(8, 8, 10, 0.92)', backdropFilter: 'blur(24px)', WebkitBackdropFilter: 'blur(24px)' }"
        @click="emit('toggle-fullscreen')"
      >
        <div
          v-if="currentTrack.coverUrl"
          class="absolute inset-0 scale-110 bg-cover bg-center blur-2xl opacity-[0.25] transition-all duration-700"
          :style="{ backgroundImage: `url(${currentTrack.coverUrl})` }"
          aria-hidden="true"
        />
        <div class="absolute inset-0 bg-gradient-to-r from-black/80 via-black/50 to-black/80" />

        <div class="relative shrink-0">
          <div class="h-10 w-10 overflow-hidden rounded-xl shadow-lg ring-1 ring-white/10">
            <img
              v-if="currentTrack.coverUrl"
              :src="currentTrack.coverUrl"
              :alt="currentTrack.title"
              class="h-full w-full object-cover"
              loading="lazy"
            />
            <div v-else class="flex h-full w-full items-center justify-center bg-gradient-to-br from-[#1db954]/30 to-[#a855f7]/30">
              <i aria-hidden="true" class="pi pi-headphones text-sm text-slate-500" />
            </div>
          </div>
        </div>

        <div class="relative min-w-0 flex-1">
          <p class="truncate text-sm font-bold text-white">{{ currentTrack.title }}</p>
          <p class="truncate text-xs text-white/50">{{ currentTrack.artistName }}</p>
        </div>

        <div class="relative flex items-center gap-0.5 pr-1" @click.stop>
          <GuestPlayGate action="play" @proceed="togglePlayPause" #default="{ proceed }">
            <button
              type="button"
              class="flex h-9 w-9 items-center justify-center rounded-full bg-white/80 text-black shadow-lg disabled:opacity-40 active:scale-90 transition-transform"
              :disabled="!currentTrack || isLoadingTrack"
              @click="isPlaying ? togglePlayPause() : proceed()"
            >
              <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-sm" />
              <i aria-hidden="true" v-else :class="isPlaying ? 'pi pi-pause' : 'pi pi-play'" class="ml-0.5 text-sm" />
            </button>
          </GuestPlayGate>
          <button
            type="button"
            class="flex h-9 w-9 items-center justify-center rounded-full text-white/60 hover:text-white active:scale-90 transition-transform"
            :disabled="!hasNext"
            @click.stop="playNext"
          >
            <i aria-hidden="true" class="pi pi-step-forward text-sm" />
          </button>
        </div>
      </button>

      <div
        v-if="playbackError"
        class="mx-3 mt-1 flex items-center justify-center gap-2 rounded-xl bg-red-500/10 px-4 py-1.5 text-xs text-red-400"
      >
        <i aria-hidden="true" class="pi pi-exclamation-circle text-xs" />
        <span>{{ playbackError }}</span>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

import { usePlayerControls } from '@/composables/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import PlayerOverflowMenu from './PlayerOverflowMenu.vue'
import GuestPlayGate from '@/components/common/GuestPlayGate.vue'

const emit = defineEmits<{
  'toggle-queue': []
  'toggle-fullscreen': []
  'toggle-lyrics': []
  'toggle-mobile-sheet': []
}>()

const showOverflow = ref(false)
const liked = ref(false)

const pc = usePlayerControls()

const {
  currentTrack,
  isPlaying,
  isBuffering,
  isLoadingTrack,
  volume,
  muted,
  currentTime,
  duration,
  error: playbackError,
  hasNext,
  hasPrevious,
  shuffleMode,
  repeatMode,
  volumeIcon,
  togglePlayPause,
  toggleMute,
  setVolume,
  progressPercent,
  seekPercent,
  playNext,
  playPrevious,
  toggleShuffle,
  toggleRepeat,
} = pc

const coverUrl = computed(() => currentTrack.value?.coverUrl || null)
const { palette } = useAlbumColors(coverUrl)

const progressFill = computed(() => {
  return coverUrl.value && palette.value.dominant ? palette.value.dominant : '#1db954'
})

const progressStyle = computed(() => ({
  width: `${progressPercent.value}%`,
  background: progressFill.value,
  boxShadow: `0 0 6px ${progressFill.value}88, 0 0 12px ${progressFill.value}44`,
}))

function formatTime(s: number) {
  if (!s || !isFinite(s)) return '0:00'
  const m = Math.floor(s / 60)
  const sec = Math.floor(s % 60)
  return `${m}:${sec.toString().padStart(2, '0')}`
}

const barRef = ref<HTMLElement | null>(null)

const volumeStyle = computed(() => ({
  '--range-progress': `${muted.value ? 0 : Number(volume.value || 0) * 100}%`,
}))

function onVolume(e: Event) {
  setVolume(Number((e.target as HTMLInputElement).value))
}

function onSeekClick(e: MouseEvent) {
  const bar = e.currentTarget as HTMLElement
  const rect = bar.getBoundingClientRect()
  const pct = ((e.clientX - rect.left) / rect.width) * 100
  seekPercent(pct)
}

function onOverflowClickOutside(e: MouseEvent) {
  if (!(e.target as HTMLElement).closest('.overflow-menu-container')) {
    showOverflow.value = false
  }
}

onMounted(() => document.addEventListener('click', onOverflowClickOutside))
onBeforeUnmount(() => document.removeEventListener('click', onOverflowClickOutside))
</script>

<style scoped>
.bar-slide-enter-active {
  transition: transform 400ms cubic-bezier(0.19, 1, 0.22, 1);
}
.bar-slide-leave-active {
  transition: transform 300ms ease-in;
}
.bar-slide-enter-from,
.bar-slide-leave-to {
  transform: translateY(100%);
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}
.animate-bounce {
  animation: bounce 0.6s ease-in-out infinite;
}

@media (prefers-reduced-motion: reduce) {
  .bar-slide-enter-active,
  .bar-slide-leave-active { transition: none; }
  .bar-slide-enter-from,
  .bar-slide-leave-to { transform: none; }
  .animate-bounce { animation: none; }
}
</style>