<template>
  <Teleport to="body">
    <Transition name="fullscreen">
      <div v-if="isOpen" class="fixed inset-0 z-[9999] flex flex-col overflow-hidden" :style="dynamicBg">

        <div class="pointer-events-none absolute -inset-5 z-0 scale-110">
          <img
            v-if="currentTrack?.coverUrl"
            :src="currentTrack?.coverUrl"
            class="h-full w-full object-cover opacity-50"
            style="filter: blur(100px) saturate(2)"
          />
        </div>
        <div class="pointer-events-none absolute inset-0 z-[1] bg-gradient-to-b from-black/60 via-black/40 to-black/80" />

        <div class="relative z-10 flex items-center justify-between px-6 pt-4 md:px-10 md:pt-6" style="padding-top: max(1rem, env(safe-area-inset-top))">
          <button class="flex h-11 w-11 items-center justify-center rounded-full text-white/50 transition-all hover:text-white hover:bg-white/10" @click="close">
            <i aria-hidden="true" class="pi pi-chevron-down text-2xl" />
          </button>
          <p class="text-xs font-semibold tracking-[0.2em] text-white/40 uppercase">Now Playing</p>
          <button class="flex h-11 w-11 items-center justify-center rounded-full text-white/50 transition-all hover:text-white hover:bg-white/10" @click="emit('toggle-queue')">
            <i aria-hidden="true" class="pi pi-list text-xl" />
          </button>
        </div>

        <div class="relative z-10 flex flex-1 flex-col items-center justify-center gap-6 md:gap-10 px-6 md:px-10 pb-6 md:pb-10">
          <div class="flex flex-col items-center gap-4">
            <div class="relative">
              <div
                class="w-[min(380px,72vw)] md:w-[min(440px,70vw)] aspect-square overflow-hidden rounded-2xl shadow-2xl transition-all duration-700"
                :class="isPlaying ? 'scale-100 cover-glow' : 'scale-95 opacity-80'"
              >
                <img
                  v-if="currentTrack?.coverUrl"
                  :src="currentTrack?.coverUrl"
                  :alt="currentTrack?.title"
                  class="h-full w-full object-cover"
                />
                <div
                  v-else
                  class="flex h-full w-full items-center justify-center bg-gradient-to-br from-[#1db954]/30 to-[#a855f7]/30"
                >
                  <i aria-hidden="true" class="pi pi-music text-5xl text-white/30" />
                </div>
              </div>
              <div
                v-if="isPlaying"
                class="pointer-events-none absolute -inset-4 animate-pulse rounded-full border-2 border-[#a855f7]/20"
              />

            </div>
          </div>

          <div class="flex w-full max-w-sm md:max-w-md flex-col gap-6 md:gap-8">
            <div class="w-full">
              <div class="flex items-center justify-between gap-3">
                <div class="min-w-0">
                  <h2 class="text-2xl md:text-3xl font-bold text-white truncate">{{ currentTrack?.title }}</h2>
                  <p class="text-base md:text-lg text-white/50 mt-1 truncate">{{ currentTrack?.artistName }}</p>
                </div>
                <button
                  class="shrink-0 flex h-11 w-11 items-center justify-center transition-all hover:scale-110"
                  :class="liked ? 'text-[#f472b6]' : 'text-white/50 hover:text-white'"
                  @click="toggleLike"
                >
                  <i aria-hidden="true" :class="liked ? 'pi pi-heart-fill' : 'pi pi-heart'" class="text-xl" />
                </button>
              </div>
            </div>

            <div class="w-full">
              <div
                class="relative flex h-10 cursor-pointer items-center group"
                ref="progressRef"
                @click="seek"
                @mousedown="startDrag"
              >
                <div class="absolute inset-x-0 h-1.5 rounded-full bg-white/20 group-hover:h-2 transition-all duration-150" />
                <div
                  class="absolute left-0 h-1.5 rounded-full transition-all duration-150 group-hover:h-2"
                  :style="{ width: progressPercent + '%', background: 'linear-gradient(90deg, #1db954, #a855f7, #f472b6)' }"
                />
                <div
                  class="absolute top-1/2 -translate-y-1/2 h-4 w-4 rounded-full bg-white shadow-xl opacity-0 group-hover:opacity-100 transition-all scale-0 group-hover:scale-100"
                  :style="{ left: `calc(${progressPercent}% - 8px)` }"
                />
              </div>
              <div class="flex justify-between mt-1.5">
                <span class="text-xs md:text-sm text-white/40 font-mono tabular-nums tracking-wide">{{ formatTime(currentTime) }}</span>
                <span class="text-xs md:text-sm text-white/40 font-mono tabular-nums tracking-wide">{{ formatTime(duration) }}</span>
              </div>
            </div>

            <div class="flex items-center justify-between px-2">
              <button
                class="flex h-12 w-12 items-center justify-center rounded-full transition-all duration-200 hover:scale-110"
                :class="shuffleMode ? 'text-[#a855f7]' : 'text-white/50 hover:text-white'"
                @click="toggleShuffle"
              >
                <i aria-hidden="true" class="pi pi-sort-alt text-lg" />
              </button>
              <button
                class="flex h-12 w-12 items-center justify-center rounded-full text-white/50 transition-all duration-200 hover:text-white hover:scale-110 disabled:opacity-25"
                :disabled="!hasPrevious"
                @click="playPrevious"
              >
                <i aria-hidden="true" class="pi pi-step-backward text-2xl" />
              </button>
              <button
                class="flex h-20 w-20 items-center justify-center rounded-full bg-white text-black shadow-2xl transition-all duration-200 active:scale-95 hover:scale-105 disabled:opacity-40"
                :disabled="!currentTrack || isLoadingTrack"
                @click="togglePlay"
              >
                <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-3xl" />
                <i aria-hidden="true" v-else :class="isPlaying ? 'pi pi-pause' : 'pi pi-play'" class="text-3xl ml-1" />
              </button>
              <button
                class="flex h-12 w-12 items-center justify-center rounded-full text-white/50 transition-all duration-200 hover:text-white hover:scale-110 disabled:opacity-25"
                :disabled="!hasNext"
                @click="playNext"
              >
                <i aria-hidden="true" class="pi pi-step-forward text-2xl" />
              </button>
              <button
                class="flex h-12 w-12 items-center justify-center rounded-full transition-all duration-200 hover:scale-110"
                :class="repeatMode !== 'off' ? 'text-[#f472b6]' : 'text-white/50 hover:text-white'"
                @click="toggleRepeat"
              >
                <i aria-hidden="true" class="pi pi-refresh text-lg" />
                <span
                  v-if="repeatMode === 'one'"
                  class="absolute -top-0.5 -right-0.5 flex h-5 w-5 items-center justify-center rounded-full bg-[#f472b6] text-[10px] font-bold text-black"
                >1</span>
              </button>
            </div>

            <div class="flex items-center gap-3 md:gap-4 max-w-xs mx-auto w-full">
              <button class="flex h-10 w-10 items-center justify-center text-white/50 transition-colors hover:text-white shrink-0" @click="toggleMute">
                <i aria-hidden="true" :class="muted ? 'pi pi-volume-off' : 'pi pi-volume-up'" class="text-base" />
              </button>
              <div class="relative flex-1 flex items-center group/vol h-5">
                <div class="absolute inset-x-0 h-1.5 rounded-full bg-white/20" />
                <div
                  class="absolute left-0 h-1.5 rounded-full bg-white/60"
                  :style="{ width: `${muted ? 0 : Number(volume) * 100}%` }"
                />
                <input
                  type="range" min="0" max="1" step="0.01"
                  class="absolute inset-0 w-full cursor-pointer opacity-0 z-10"
                  :value="muted ? 0 : volume"
                  @input="onVolume"
                />
              </div>
              <button class="flex h-10 w-10 items-center justify-center text-white/50 transition-colors hover:text-white shrink-0" @click="emit('toggle-lyrics')">
                <i aria-hidden="true" class="pi pi-align-left text-base" />
              </button>
            </div>
          </div>
        </div>

      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { usePlayerControls } from '@/composables/player'
import { useAlbumColors } from '@/composables/useAlbumColors'

const props = withDefaults(defineProps<{
  visible: boolean
  initialTab?: string
}>(), {
  initialTab: 'now-playing',
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'toggle-queue': []
  'toggle-lyrics': []
}>()

const isOpen = ref(props.visible)
watch(() => props.visible, (v) => { isOpen.value = v })
watch(isOpen, (v) => { emit('update:visible', v) })

const pc = usePlayerControls()

const {
  currentTrack,
  isPlaying,
  isBuffering,
  isLoadingTrack,
  currentTime,
  duration,
  volume,
  muted,
  hasNext,
  hasPrevious,
  shuffleMode,
  repeatMode,
  togglePlayPause,
  toggleMute,
  setVolume,
  playNext,
  playPrevious,
  toggleShuffle,
  toggleRepeat,
} = pc

const coverUrl = computed(() => currentTrack.value?.coverUrl || null)
const { palette } = useAlbumColors(coverUrl)

const progressRef = ref<HTMLElement>()

const progressPercent = computed(() =>
  duration.value ? (currentTime.value / duration.value) * 100 : 0
)

const dynamicBg = computed(() => {
  const p = palette.value
  return {
    background: coverUrl.value
      ? `radial-gradient(ellipse 80% 60% at 50% 0%, ${p.vibrant}33 0%, transparent 70%), radial-gradient(ellipse 60% 40% at 100% 100%, ${p.muted}44 0%, transparent 60%), ${p.dark}`
      : '#08080A',
  }
})

function seekTo(seconds: number) {
  if (!duration.value) return
  const clamped = Math.max(0, Math.min(seconds, duration.value))
  pc.seek(clamped)

}

function seek(e: MouseEvent) {
  if (!progressRef.value || !duration.value) return
  const rect = progressRef.value.getBoundingClientRect()
  const ratio = (e.clientX - rect.left) / rect.width
  seekTo(ratio * duration.value)
}

function startDrag() {
  const move = (ev: MouseEvent) => seek(ev)
  const up = () => {
    window.removeEventListener('mousemove', move)
    window.removeEventListener('mouseup', up)
  }
  window.addEventListener('mousemove', move)
  window.addEventListener('mouseup', up)
}

function formatTime(s: number) {
  if (!s || !isFinite(s)) return '0:00'
  const m = Math.floor(s / 60)
  const sec = Math.floor(s % 60)
  return `${m}:${sec.toString().padStart(2, '0')}`
}

function togglePlay() {
  togglePlayPause()
}

function close() {
  isOpen.value = false
}

function onVolume(e: Event) {
  setVolume(Number((e.target as HTMLInputElement).value))
}

const rootEl = ref<HTMLElement | null>(null)

onMounted(async () => {
  await nextTick()
  rootEl.value?.focus()
})

onBeforeUnmount(() => {})

const liked = ref(false)
function toggleLike() {
  liked.value = !liked.value
}
</script>

<style scoped>
.fullscreen-enter-active,
.fullscreen-leave-active {
  transition: all 0.5s cubic-bezier(0.19, 1, 0.22, 1);
}
.fullscreen-enter-from,
.fullscreen-leave-to {
  opacity: 0;
  transform: translateY(100%);
}

@keyframes pulse {
  0%, 100% { opacity: 0.4; }
  50% { opacity: 0.15; }
}
.animate-pulse {
  animation: pulse 2s ease-in-out infinite;
}

@keyframes cover-glow {
  0%, 100% { box-shadow: 0 0 60px rgba(168, 85, 247, 0.15), 0 0 120px rgba(29, 185, 84, 0.08); }
  50% { box-shadow: 0 0 80px rgba(168, 85, 247, 0.3), 0 0 160px rgba(29, 185, 84, 0.15); }
}
.cover-glow {
  animation: cover-glow 3s ease-in-out infinite;
}

@keyframes float-anim {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-6px); }
}

@media (prefers-reduced-motion: reduce) {
  .fullscreen-enter-active,
  .fullscreen-leave-active { transition: none; }
  .fullscreen-enter-from,
  .fullscreen-leave-to { transform: none; opacity: 1; }
  .cover-glow { animation: none; }
}
</style>
