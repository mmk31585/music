<template>
  <Teleport to="body">
    <Transition name="sheet-slide">
      <div
        v-if="visible"
        ref="sheetEl"
        class="fixed inset-x-0 bottom-0 z-[120] touch-none select-none md:hidden"
        :style="sheetStyle"
      >
        <!-- Background with dynamic album color -->
        <div class="absolute inset-0 backdrop-blur-2xl" :style="{ background: bgGradient }" />
        <div class="absolute inset-0 bg-black/40" />

        <!-- Drag handle (44px min touch target for WCAG 2.5.8) -->
        <div
          class="absolute top-0 left-1/2 z-20 flex w-full -translate-x-1/2 flex-col items-center justify-center min-h-[44px] gap-1"
          @touchstart.prevent="onDragStart"
          @mousedown.prevent="onDragStart"
        >
          <div class="h-1 w-10 rounded-full bg-white/30" />
          <div class="h-0.5 w-6 rounded-full bg-white/15" />
        </div>

        <div ref="sheetRef" class="relative z-10 flex h-full flex-col px-5 pt-12" :style="{ paddingBottom: 'calc(env(safe-area-inset-bottom, 0px) + 1.5rem)' }">
          <!-- Track info + art -->
          <div class="flex items-start gap-4">
            <div
              class="h-20 w-20 shrink-0 overflow-hidden rounded-2xl shadow-xl ring-1 ring-white/10 transition-transform duration-300"
              :class="{ 'scale-105': isPlaying }"
            >
              <img
                v-if="coverUrl"
                :src="coverUrl"
                :alt="title"
                class="h-full w-full object-cover"
                :class="{ 'vinyl-spin': isPlaying }"
                loading="lazy"
                @error="onImgError"
              />
              <div v-else class="flex h-full items-center justify-center bg-white/10">
                <i aria-hidden="true" class="pi pi-music text-white/30" />
              </div>
            </div>

            <div class="min-w-0 flex-1">
              <p class="truncate text-base font-bold text-white">{{ title }}</p>
              <p class="truncate text-sm text-white/50">{{ artistName }}</p>
            </div>

            <button
              type="button"
              class="flex h-10 w-10 items-center justify-center rounded-full text-white/50 transition-all hover:bg-white/10 hover:text-white"
              aria-label="Close"
              @click="close"
            >
              <i aria-hidden="true" class="pi pi-chevron-down text-lg" />
            </button>
          </div>

          <!-- Like button -->
          <div class="mt-3 flex items-center gap-3">
            <button
              type="button"
              class="spring flex items-center gap-1.5 rounded-full px-3 py-2 text-sm font-medium transition-all active:scale-95 min-h-[44px]"
              :class="liked ? '' : 'text-white/40 hover:bg-white/5'"
              :style="liked ? { color: accentColor, backgroundColor: `${accentColor}15` } : {}"
              @click="liked = !liked"
            >
              <i
                :class="liked ? 'pi pi-heart-fill heart-pulse' : 'pi pi-heart'"
                class="text-base"
              />
              <span>{{ liked ? 'Liked' : 'Like' }}</span>
            </button>
          </div>

          <!-- Seekbar -->
          <div class="mt-4">
            <input
              type="range"
              min="0"
              max="100"
              step="0.1"
              class="sheet-range w-full"
              :style="progressStyle"
              :value="progressPercent"
              @input="onSeek"
            />
            <div class="mt-1 flex justify-between text-[10px] text-white/40 tabular-nums">
              <span>{{ currentTimeLabel }}</span>
              <span>{{ durationLabel }}</span>
            </div>
          </div>

          <!-- Main controls -->
          <div
            class="mt-5 flex items-center justify-center gap-7"
            :style="{ '--accent': accentColor }"
          >
            <button
              type="button"
              class="relative flex flex-col items-center gap-1 transition-all hover:text-white active:scale-90"
              :class="shuffleMode !== 'off' ? '' : 'text-white/40'"
              :style="shuffleMode !== 'off' ? { color: accentColor } : {}"
              :title="shuffleMode === 'off' ? 'Shuffle off' : shuffleMode === 'queue' ? 'Shuffle queue' : shuffleMode === 'catalog' ? 'Random catalog tracks' : 'Similar tracks'"
              @click="toggleShuffle"
            >
              <i aria-hidden="true" class="pi pi-sort-alt text-lg" />
              <span class="text-[8px] font-medium">Shuffle</span>
              <span
                v-if="shuffleMode !== 'off'"
                class="absolute -top-0.5 -right-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-[#a855f7] text-[8px] font-bold text-white"
              >{{ shuffleMode === 'queue' ? 'Q' : shuffleMode === 'catalog' ? 'R' : 'S' }}</span>
            </button>

            <button
              type="button"
              class="text-white/50 transition-all hover:text-white active:scale-90 disabled:opacity-30"
              :disabled="!hasPrevious"
              aria-label="Previous track"
              @click="playPrevious"
            >
              <i aria-hidden="true" class="pi pi-step-backward text-2xl" />
            </button>

            <button
              type="button"
              class="spring relative flex h-16 w-16 items-center justify-center rounded-full bg-white text-black shadow-xl transition-all hover:scale-110 hover:text-white active:scale-95"
              :style="
                isPlaying
                  ? {
                      backgroundColor: accentColor,
                      color: 'white',
                      boxShadow: `0 0 30px ${accentColor}66`,
                    }
                  : { '--accent-hover': accentColor }
              "
              :class="!isPlaying ? 'hover:bg-[var(--accent-hover,#1db954)]' : ''"
              :disabled="!currentTrack"
              :aria-label="isPlaying ? 'Pause' : 'Play'"
              @click="togglePlayPause"
            >
              <i aria-hidden="true" v-if="isBuffering" class="pi pi-spin pi-spinner text-xl" />
              <i
                v-else
                :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'"
                class="text-xl"
              />
              <span
                v-if="isPlaying"
                class="absolute -inset-1.5 animate-ping rounded-full"
                :style="{ border: `1px solid ${accentColor}40` }"
              />
            </button>

            <button
              type="button"
              class="text-white/50 transition-all hover:text-white active:scale-90 disabled:opacity-30"
              :disabled="!hasNext"
              aria-label="Next track"
              @click="playNext"
            >
              <i aria-hidden="true" class="pi pi-step-forward text-2xl" />
            </button>

            <button
              type="button"
              class="flex flex-col items-center gap-1 transition-all hover:text-white active:scale-90"
              :class="repeatMode !== 'off' ? '' : 'text-white/40'"
              :style="repeatMode !== 'off' ? { color: accentColor } : {}"
              @click="toggleRepeat"
            >
              <i aria-hidden="true" class="pi pi-refresh text-lg" />
              <span class="text-[8px] font-medium">
                {{ repeatMode === 'one' ? '1' : repeatMode === 'all' ? 'All' : 'Off' }}
              </span>
            </button>
          </div>

          <!-- Bottom row -->
          <div class="mt-5 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <button
                type="button"
                class="flex h-10 w-10 items-center justify-center rounded-full text-white/40 transition-all hover:bg-white/10 hover:text-white"
                :aria-label="muted ? 'Unmute' : 'Mute'"
                @click="toggleMute"
              >
                <i aria-hidden="true" :class="volumeIcon" class="text-base" />
              </button>
              <input
                type="range"
                min="0"
                max="1"
                step="0.01"
                class="sheet-range w-20"
                :style="volumeStyle"
                :value="muted ? 0 : Number(volume || 0)"
                @input="onVolume"
              />
            </div>

            <button
              type="button"
              class="spring flex items-center gap-1.5 rounded-full px-4 py-2 text-xs font-bold transition-all"
              :class="playbackRate !== 1 ? '' : 'text-white/40'"
              :style="
                playbackRate !== 1
                  ? { color: accentColor, backgroundColor: `${accentColor}18` }
                  : {}
              "
              @click="cycleSpeed"
            >
              <i aria-hidden="true" class="pi pi-forward text-[10px]" />
              {{ speedLabel }}
            </button>

            <button
              type="button"
              class="flex h-10 w-10 items-center justify-center rounded-full text-white/40 transition-all hover:bg-white/10 hover:text-white"
              aria-label="Fullscreen"
              @click="$emit('open-fullscreen')"
            >
              <i aria-hidden="true" class="pi pi-expand text-base" />
            </button>
          </div>

        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { usePlayerControls } from '@/composables/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import { onImgError } from '@/utils/helpers'
import { useSwipe } from '@/composables'

const props = defineProps<{ visible: boolean }>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'open-fullscreen': []
}>()

const liked = ref(false)

const pc = usePlayerControls()

const currentTrack = pc.currentTrack
const isPlaying = pc.isPlaying
const isBuffering = pc.isBuffering
const currentTime = pc.currentTime
const duration = pc.duration
const progressPercent = pc.progressPercent
const volume = pc.volume
const muted = pc.muted
const hasNext = pc.hasNext
const hasPrevious = pc.hasPrevious
const shuffleMode = pc.shuffleMode
const repeatMode = pc.repeatMode
const playbackRate = pc.playbackRate
const volumeIcon = pc.volumeIcon
const speedLabel = pc.speedLabel

const togglePlayPause = pc.togglePlayPause
const toggleMute = pc.toggleMute
const playNext = pc.playNext
const playPrevious = pc.playPrevious
const toggleShuffle = pc.toggleShuffle
const toggleRepeat = pc.toggleRepeat
const seekPercent = pc.seekPercent
const setVolume = pc.setVolume
const setPlaybackRate = pc.setPlaybackRate

const sheetRef = ref<HTMLElement | null>(null)

const speedOptions = [0.5, 0.75, 1, 1.25, 1.5, 2]
function cycleSpeed() {
  const idx = speedOptions.indexOf(playbackRate)
  const nextIdx = (idx + 1) % speedOptions.length
  setPlaybackRate(speedOptions[nextIdx]!)
}

useSwipe({ element: sheetRef, onSwipeLeft: () => playNext(), onSwipeRight: () => playPrevious(), threshold: 80 })

const title = computed(() => currentTrack.value?.title || '')
const artistName = computed(() => currentTrack.value?.artistName || '')
const coverUrl = computed(() => currentTrack.value?.coverUrl || '')

const { palette } = useAlbumColors(coverUrl)
const accentColor = computed(() => palette.value.vibrant || '#1db954')

const bgGradient = computed(() => {
  const p = palette.value
  if (!coverUrl.value) return 'linear-gradient(135deg, #0a0a0a 0%, #121212 100%)'
  return `linear-gradient(180deg, ${p.dark} 0%, ${p.dominant}88 40%, ${p.muted} 100%)`
})

const sheetEl = ref<HTMLElement | null>(null)
const isDragging = ref(false)

const sheetHeight = ref(0)

function recalcHeight() {
  const vh = window.innerHeight
  const bottomNavH = 64
  const minH = 420
  const maxH = vh - bottomNavH - 16
  sheetHeight.value = Math.max(minH, Math.min(maxH, vh * 0.72))
}

const sheetStyle = computed(() => ({
  height: `${sheetHeight.value}px`,
  bottom: '0',
  transform: `translateY(0px)`,
  transition: 'transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1)',
}))

let cleanupListeners: (() => void) | null = null

function onDragStart(e: TouchEvent | MouseEvent) {
  cleanupListeners?.()
  isDragging.value = true
  const startY = 'touches' in e ? e.touches[0]!.clientY : e.clientY
  const el = sheetEl.value
  if (!el) return

  el.style.transition = 'none'

  const onMove = (ev: TouchEvent | MouseEvent) => {
    const y = 'touches' in ev ? ev.touches[0]!.clientY : ev.clientY
    const delta = Math.max(0, y - startY)
    el.style.transform = `translateY(${delta}px)`
  }

  const onEnd = () => {
    isDragging.value = false
    const match = el.style.transform.match(/translateY\((\d+(?:\.\d+)?)px\)/)
    el.style.transition = ''
    el.style.transform = ''
    if (match && parseFloat(match[1]!) > sheetHeight.value * 0.35) {
      emit('update:visible', false)
    }
    cleanupListeners?.()
    cleanupListeners = null
  }

  document.addEventListener('mousemove', onMove, { passive: true })
  document.addEventListener('mouseup', onEnd)
  document.addEventListener('touchmove', onMove, { passive: true })
  document.addEventListener('touchend', onEnd)

  cleanupListeners = () => {
    document.removeEventListener('mousemove', onMove)
    document.removeEventListener('mouseup', onEnd)
    document.removeEventListener('touchmove', onMove)
    document.removeEventListener('touchend', onEnd)
  }
}

function fmtTime(s: number) {
  const total = Math.max(0, Math.floor(Number(s) || 0))
  return `${Math.floor(total / 60)}:${String(total % 60).padStart(2, '0')}`
}

const currentTimeLabel = computed(() => fmtTime(currentTime.value))
const durationLabel = computed(() =>
  fmtTime(duration.value || currentTrack.value?.durationSeconds || 0),
)
const progressStyle = computed(() => ({
  '--range-progress': `${progressPercent.value}%`,
  '--accent-color': accentColor.value,
}))
const volumeStyle = computed(() => ({
  '--range-progress': `${muted.value ? 0 : Number(volume.value || 0) * 100}%`,
  '--accent-color': accentColor.value,
}))

function onSeek(e: Event) {
  seekPercent(Number((e.target as HTMLInputElement).value))
}

function onVolume(e: Event) {
  setVolume(Number((e.target as HTMLInputElement).value))
}

function close() {
  emit('update:visible', false)
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') close()
}

watch(
  () => props.visible,
  (val) => {
    if (val) window.history.pushState(null, '')
  },
)

onMounted(() => {
  recalcHeight()
  window.addEventListener('popstate', close)
  window.addEventListener('resize', recalcHeight)
  document.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  window.removeEventListener('popstate', close)
  window.removeEventListener('resize', recalcHeight)
  document.removeEventListener('keydown', onKeydown)
})
</script>

<style scoped>
/* WCAG 2.5.8: minimum target 24×24 CSS px.
   We keep the visual track thin but make the hit area 44px tall for touch. */
.sheet-range {
  --range-progress: 0%;
  --accent-color: #1db954;
  width: 100%;
  height: 44px;
  cursor: pointer;
  appearance: none;
  background: transparent;
  padding: 14px 0;    /* vertically center the 3px track inside 44px */
  box-sizing: border-box;
}
.sheet-range::-webkit-slider-runnable-track {
  height: 3px;
  border-radius: 999px;
  background: linear-gradient(
    to right,
    var(--accent-color) 0%,
    var(--accent-color) var(--range-progress),
    rgba(255, 255, 255, 0.12) var(--range-progress),
    rgba(255, 255, 255, 0.12) 100%
  );
}
.sheet-range::-webkit-slider-thumb {
  width: 24px;
  height: 24px;
  margin-top: -10.5px;   /* (track 3px - thumb 24px) / 2 = -10.5px */
  border-radius: 999px;
  appearance: none;
  background: #fff;
  box-shadow: 0 0 16px color-mix(in srgb, var(--accent-color) 50%, transparent);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}
.sheet-range::-webkit-slider-thumb:hover,
.sheet-range::-webkit-slider-thumb:focus-visible {
  transform: scale(1.2);
  box-shadow: 0 0 24px color-mix(in srgb, var(--accent-color) 70%, transparent);
}
/* Firefox thumb */
.sheet-range::-moz-range-thumb {
  width: 24px;
  height: 24px;
  border: none;
  border-radius: 999px;
  background: #fff;
  box-shadow: 0 0 16px color-mix(in srgb, var(--accent-color) 50%, transparent);
  cursor: pointer;
}
.sheet-slide-enter-active {
  transition: transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.sheet-slide-leave-active {
  transition: transform 0.2s ease-in;
}
.sheet-slide-enter-from {
  transform: translateY(100%);
}
.sheet-slide-leave-to {
  transform: translateY(100%);
}
.spring {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
@keyframes heart-pop {
  0% {
    transform: scale(1);
  }
  30% {
    transform: scale(1.3);
  }
  60% {
    transform: scale(0.9);
  }
  100% {
    transform: scale(1);
  }
}
.heart-pulse {
  animation: heart-pop 400ms cubic-bezier(0.34, 1.56, 0.64, 1);
}
</style>
