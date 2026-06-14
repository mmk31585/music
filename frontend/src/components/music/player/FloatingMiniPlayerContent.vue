<template>
  <div
    ref="floatingEl"
    class="popout-player"
    :class="{ 'is-pip': isPiP, dragging: isDragging }"
    :style="{ left: posX + 'px', top: posY + 'px', width: styleWidth, height: styleHeight, willChange: isDragging ? 'transform' : 'auto' }"
    @mousedown="onPointerDown"
    @touchstart.prevent="onPointerDown"
  >
    <div class="popout-backdrop" :style="backdropStyle">
      <div v-if="coverUrl" class="popout-backdrop-image" :style="{ backgroundImage: `url(${coverUrl})` }" />
      <div class="popout-backdrop-overlay" />
    </div>

    <div v-if="mode === 'expanded' || isPiP" class="popout-header">
      <button type="button" class="popout-btn" title="Minimize" @click="mode = 'mini'">
        <i class="pi pi-window-minimize" />
      </button>
      <button type="button" class="popout-btn" title="Full player" @click="openFullscreen">
        <i class="pi pi-expand" />
      </button>
      <button
        type="button"
        class="popout-btn"
        :title="isPiP ? 'Return to app' : 'Pop out'"
        @click="onTogglePiP"
      >
        <i :class="isPiP ? 'pi pi-window-maximize' : 'pi pi-external-link'" />
      </button>
      <button v-if="!isPiP" type="button" class="popout-btn" title="Close" @click="close">
        <i class="pi pi-times" />
      </button>
    </div>

    <div v-if="mode === 'mini' && !isPiP" class="popout-mini" @dblclick="mode = 'expanded'">
      <img v-if="coverUrl" :src="coverUrl" :alt="title" loading="lazy" class="popout-cover" @error="onImgError" />
      <div v-else class="popout-cover-placeholder">
        <i class="pi pi-music" />
      </div>
      <div v-if="isPlaying" class="popout-equalizer">
        <span /><span /><span />
      </div>
      <div class="popout-mini-overlay" @click="togglePlayPause">
        <i :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" class="popout-play-icon" />
      </div>
    </div>

    <div v-else class="popout-expanded">
      <div class="popout-track-info">
        <img v-if="coverUrl" :src="coverUrl" :alt="title" loading="lazy" class="popout-expanded-cover" @error="onImgError" />
        <div v-else class="popout-expanded-cover-placeholder">
          <i class="pi pi-music" />
        </div>
        <div class="popout-meta">
          <p class="popout-title">{{ title || 'No track' }}</p>
          <p class="popout-artist">{{ artistName || '—' }}</p>
        </div>
      </div>

      <div class="popout-progress">
        <div class="popout-progress-bar" ref="progressRef" @click="seekFromEvent">
          <div class="popout-progress-fill" :style="{ transform: `scaleX(${progressPercent / 100})` }" />
        </div>
        <div class="popout-progress-labels">
          <span>{{ currentTimeLabel }}</span>
          <span>{{ durationLabel }}</span>
        </div>
      </div>

      <div class="popout-controls">
        <button type="button" class="popout-ctrl-btn" @click="playPrevious" :disabled="!hasPrevious">
          <i class="pi pi-step-backward" />
        </button>
        <button type="button" class="popout-play-btn" :style="{ background: accentColor }" @click="togglePlayPause">
          <i :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" />
        </button>
        <button type="button" class="popout-ctrl-btn" @click="playNext" :disabled="!hasNext">
          <i class="pi pi-step-forward" />
        </button>
      </div>

      <div class="popout-volume">
        <button type="button" class="popout-ctrl-btn" @click="toggleMute">
          <i :class="volumeIcon" />
        </button>
        <input
          type="range" min="0" max="1" step="0.01"
          class="popout-range"
          :value="muted ? 0 : Number(volume || 0)"
          @input="onVolume"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'
import { usePlayerControls } from '@/composables/player'
import { onImgError } from '@/utils/helpers'
import { useAlbumColors } from '@/composables/useAlbumColors'

const props = withDefaults(
  defineProps<{
    visible?: boolean
    isPiP?: boolean
  }>(),
  { visible: true, isPiP: false },
)

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'open-fullscreen': []
  'toggle-pip': []
}>()

const pc = usePlayerControls()
const {
  isPlaying,
  progressPercent,
  hasNext,
  hasPrevious,
  volume,
  muted,
  volumeIcon,
  playNext,
  playPrevious,
  togglePlayPause,
  toggleMute,
  currentTime,
  duration,
  currentTrack,
} = pc

const mode = ref<'mini' | 'expanded'>('mini')
const floatingEl = ref<HTMLElement | null>(null)
const progressRef = ref<HTMLElement | null>(null)

const title = computed(() => currentTrack.value?.title || '')
const artistName = computed(() => currentTrack.value?.artistName || '')
const coverUrl = computed(() => currentTrack.value?.coverUrl || '')

const { palette } = useAlbumColors(computed(() => currentTrack.value?.coverUrl || null))
const accentColor = computed(() => palette.value.vibrant || '#1db954')

const backdropStyle = computed(() => {
  if (!coverUrl.value) return {}
  return {
    background: `linear-gradient(135deg, ${palette.value.dark} 0%, ${palette.value.dominant} 60%, ${palette.value.muted} 100%)`,
  }
})

let dragRAF: number | null = null

const isDragging = ref(false)
let dragStartX = 0
let dragStartY = 0
let initialX = 0
let initialY = 0
let snapped = false

const posX = ref(typeof window !== 'undefined' ? window.innerWidth - 340 : 1000)
const posY = ref(typeof window !== 'undefined' ? window.innerHeight - 220 : 500)
const offsetX = ref(0)
const offsetY = ref(0)

const styleWidth = computed(() => props.isPiP ? '100vw' : (mode.value === 'mini' ? '72px' : '292px'))
const styleHeight = computed(() => props.isPiP ? '100vh' : (mode.value === 'mini' ? '72px' : 'auto'))

function onPointerDown(e: MouseEvent | TouchEvent) {
  if (props.isPiP) return
  const target = ('touches' in e ? e.target : e.target) as HTMLElement
  if (target.closest('button') || target.closest('input')) return

  isDragging.value = true
  const el = floatingEl.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  const clientX = 'touches' in e ? e.touches[0]!.clientX : e.clientX
  const clientY = 'touches' in e ? e.touches[0]!.clientY : e.clientY

  dragStartX = clientX
  dragStartY = clientY
  initialX = offsetX.value
  initialY = offsetY.value

  const moveEvent = 'touches' in e ? 'touchmove' : 'mousemove'
  const upEvent = 'touches' in e ? 'touchend' : 'mouseup'
  document.addEventListener(moveEvent, onPointerMove, { passive: false })
  document.addEventListener(upEvent, onPointerUp)
}

function onPointerMove(e: MouseEvent | TouchEvent) {
  if (dragRAF !== null) return
  dragRAF = requestAnimationFrame(() => {
    dragRAF = null
    const clientX = 'touches' in e ? (e as TouchEvent).touches[0]!.clientX : (e as MouseEvent).clientX
    const clientY = 'touches' in e ? (e as TouchEvent).touches[0]!.clientY : (e as MouseEvent).clientY
    const dx = clientX - dragStartX
    const dy = clientY - dragStartY
    if (floatingEl.value) {
      floatingEl.value.style.transform = `translate(${initialX + dx}px, ${initialY + dy}px)`
    }
    snapped = false
    if ('touches' in e) e.preventDefault()
  })
}

function onPointerUp() {
  if (dragRAF !== null) {
    cancelAnimationFrame(dragRAF)
    dragRAF = null
  }
  isDragging.value = false
  snapToEdge()
  document.removeEventListener('touchmove', onPointerMove)
  document.removeEventListener('touchend', onPointerUp)
  document.removeEventListener('mousemove', onPointerMove)
  document.removeEventListener('mouseup', onPointerUp)
}

function snapToEdge() {
  const el = floatingEl.value
  if (!el) return
  const w = props.isPiP ? window.innerWidth : el.offsetWidth
  const h = props.isPiP ? window.innerHeight : el.offsetHeight
  const cx = posX.value + offsetX.value
  const cy = posY.value + offsetY.value
  const margin = 12
  const snapDist = 120
  let nx = cx
  let ny = cy

  if (cx < snapDist) nx = margin
  else if (cx + w > window.innerWidth - snapDist) nx = window.innerWidth - w - margin

  ny = Math.max(margin, Math.min(window.innerHeight - h - margin, ny))

  posX.value = nx
  posY.value = ny
  offsetX.value = 0
  offsetY.value = 0
  snapped = true
}

const timeCache = new Map<number, string>()
function fmtTime(s: number) {
  const total = Math.max(0, Math.floor(Number(s) || 0))
  const cached = timeCache.get(total)
  if (cached) return cached
  const str = `${Math.floor(total / 60)}:${String(total % 60).padStart(2, '0')}`
  timeCache.set(total, str)
  if (timeCache.size > 200) timeCache.clear()
  return str
}

const currentTimeLabel = computed(() => fmtTime(currentTime.value))
const durationLabel = computed(() =>
  fmtTime(duration.value || currentTrack.value?.durationSeconds || 0),
)

function seekFromEvent(e: MouseEvent) {
  const el = progressRef.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  const ratio = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width))
  pc.seek(ratio * (duration.value || currentTrack.value?.durationSeconds || 0))
}

function close() { emit('update:visible', false) }
function openFullscreen() { emit('open-fullscreen') }
function onTogglePiP() { emit('toggle-pip') }
function onVolume(e: Event) { pc.setVolume(Number((e.target as HTMLInputElement).value)) }

onBeforeUnmount(() => {
  if (dragRAF !== null) cancelAnimationFrame(dragRAF)
  document.removeEventListener('touchmove', onPointerMove)
  document.removeEventListener('touchend', onPointerUp)
  document.removeEventListener('mousemove', onPointerMove)
  document.removeEventListener('mouseup', onPointerUp)
})
</script>

<style scoped>
.popout-player {
  position: fixed;
  z-index: 150;
  overflow: hidden;
  border-radius: 16px;
  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.5);
  transition: border-radius 0.3s ease;
  cursor: grab;
  user-select: none;
  backface-visibility: hidden;
}
.popout-player.is-pip {
  width: 100vw !important;
  height: 100vh !important;
  border-radius: 0;
  cursor: default;
}
.popout-player.dragging {
  cursor: grabbing;
  transition: none;
  box-shadow: 0 16px 64px rgba(0, 0, 0, 0.6);
}

.popout-backdrop {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #1a1a2e, #121212, #0a0a0a);
  transition: background 0.5s ease;
}
.popout-backdrop-image {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
  opacity: 0.15;
}
.popout-backdrop-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.6), transparent, rgba(0, 0, 0, 0.2));
}

.popout-header {
  position: absolute;
  top: 4px;
  right: 4px;
  z-index: 10;
  display: flex;
  gap: 2px;
}
.popout-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: none;
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.5);
  font-size: 10px;
  cursor: pointer;
  transition: all 0.15s;
  backdrop-filter: blur(4px);
}
.popout-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
}

/* MINI */
.popout-mini {
  position: relative;
  width: 72px;
  height: 72px;
  cursor: pointer;
  overflow: hidden;
}
.popout-cover {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}
.popout-mini:hover .popout-cover {
  transform: scale(1.08);
}
.popout-cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.3);
  font-size: 20px;
}
.popout-mini-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
  opacity: 0;
  transition: opacity 0.2s;
}
.popout-mini:hover .popout-mini-overlay {
  opacity: 1;
}
.popout-play-icon {
  font-size: 20px;
  color: #fff;
  filter: drop-shadow(0 2px 8px rgba(0, 0, 0, 0.5));
}

.popout-equalizer {
  position: absolute;
  bottom: 4px;
  left: 4px;
  display: flex;
  gap: 2px;
  align-items: end;
  height: 14px;
}
.popout-equalizer span {
  display: block;
  width: 3px;
  background: #1db954;
  border-radius: 1px;
  animation: eq 0.6s infinite alternate ease-in-out;
}
.popout-equalizer span:nth-child(1) { height: 6px; animation-delay: 0s; }
.popout-equalizer span:nth-child(2) { height: 10px; animation-delay: 0.2s; }
.popout-equalizer span:nth-child(3) { height: 8px; animation-delay: 0.4s; }
@keyframes eq {
  0% { transform: scaleY(0.4); }
  100% { transform: scaleY(1); }
}

/* EXPANDED */
.popout-expanded {
  width: 280px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 14px;
  position: relative;
  z-index: 1;
}
.popout-track-info {
  display: flex;
  align-items: center;
  gap: 10px;
}
.popout-expanded-cover {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  object-fit: cover;
  flex-shrink: 0;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.4);
}
.popout-expanded-cover-placeholder {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.3);
  flex-shrink: 0;
}
.popout-meta {
  min-width: 0;
  flex: 1;
}
.popout-title {
  font-weight: 700;
  font-size: 14px;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.popout-artist {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.popout-progress {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.popout-progress-bar {
  height: 4px;
  border-radius: 2px;
  background: rgba(255, 255, 255, 0.1);
  overflow: hidden;
  cursor: pointer;
}
.popout-progress-fill {
  width: 100%;
  height: 100%;
  border-radius: 2px;
  background: linear-gradient(to right, #1db954, #1ed760);
  transform-origin: left center;
  transition: transform 0.2s linear;
}
.popout-progress-labels {
  display: flex;
  justify-content: space-between;
  font-size: 10px;
  color: rgba(255, 255, 255, 0.3);
  font-variant-numeric: tabular-nums;
}

.popout-controls {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
}
.popout-ctrl-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: none;
  background: transparent;
  color: rgba(255, 255, 255, 0.4);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.15s;
}
.popout-ctrl-btn:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
}
.popout-ctrl-btn:disabled {
  opacity: 0.3;
  cursor: default;
}
.popout-play-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: none;
  color: #fff;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.15s;
  box-shadow: 0 0 16px rgba(0, 0, 0, 0.3);
}
.popout-play-btn:hover {
  transform: scale(1.05);
  filter: brightness(1.1);
}

.popout-volume {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}
.popout-range {
  -webkit-appearance: none;
  appearance: none;
  width: 80px;
  height: 3px;
  border-radius: 2px;
  background: rgba(255, 255, 255, 0.15);
  outline: none;
}
.popout-range::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #fff;
  cursor: pointer;
  border: none;
}
.popout-range::-moz-range-thumb {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #fff;
  cursor: pointer;
  border: none;
}
</style>
