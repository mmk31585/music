<template>
  <div :style="rootStyle">
    <div :style="bgImageStyle" />
    <div :style="bgOverlay" />

    <div :style="contentStyle">
      <div :style="topSectionStyle">
        <div :style="coverWrapStyle">
          <img
            v-if="coverUrl"
            :src="coverUrl"
            :alt="title"
            :style="{ width: '100%', height: '100%', objectFit: 'cover' }"
            @error="onImgError"
          />
          <div v-else :style="coverFallbackStyle">
            <span style="font-size: 18px">♪</span>
          </div>
        </div>
        <div :style="metaStyle">
          <div :style="titleStyle">{{ title || 'No track' }}</div>
          <div :style="artistStyle">{{ artistName || '—' }}</div>
        </div>
      </div>

      <div :style="progressSectionStyle">
        <div :style="progressBarStyle" ref="progressRef" role="button" tabindex="0" @click="seekFromEvent" @keydown.enter="seekFromEvent" @keydown.space.prevent="seekFromEvent">
          <div :style="{ ...progressFillStyle, width: `${progressPercent}%` }" />
        </div>
        <div :style="progressLabelsStyle">
          <span>{{ currentTimeLabel }}</span>
          <span>{{ durationLabel }}</span>
        </div>
      </div>

      <div :style="controlsStyle">
        <button :style="ctrlBtnStyle" aria-label="Previous track" @click="playPrevious" :disabled="!hasPrevious">⏮</button>
        <button
          :style="{ ...playBtnBase, background: accentColor, boxShadow: `0 0 12px ${accentColor}44` }"
          :aria-label="isPlaying ? 'Pause' : 'Play'"
          @click="togglePlayPause"
        >
          {{ isPlaying ? '⏸' : '▶' }}
        </button>
        <button :style="ctrlBtnStyle" aria-label="Next track" @click="playNext" :disabled="!hasNext">⏭</button>
        <div :style="{ flex: 1 }" />
        <button :style="closeBtnStyle" aria-label="Close player" @click="close">✕</button>
      </div>
    </div>

    <div v-if="isPlaying" :style="liveDotStyle" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { usePlayerControls } from '@/composables/player'
import { onImgError } from '@/utils/helpers'
import { useAlbumColors } from '@/composables/useAlbumColors'

const emit = defineEmits<{ close: [] }>()

const pc = usePlayerControls()
const {
  isPlaying,
  progressPercent,
  hasNext,
  hasPrevious,
  playNext,
  playPrevious,
  togglePlayPause,
} = pc
const progressRef = ref<HTMLElement | null>(null)

const title = computed(() => pc.currentTrack.value?.title || '')
const artistName = computed(() => pc.currentTrack.value?.artistName || '')
const coverUrl = computed(() => pc.currentTrack.value?.coverUrl || '')

const { palette } = useAlbumColors(computed(() => pc.currentTrack.value?.coverUrl || null))
const accentColor = computed(() => palette.value.vibrant || '#1db954')
const bgDark = computed(() => palette.value.dark || '#1a1a2e')
const bgDominant = computed(() => palette.value.dominant || '#121212')
const bgMuted = computed(() => palette.value.muted || '#0a0a0a')

const rootStyle = computed(() => ({
  width: '100vw',
  height: '100vh',
  margin: 0,
  padding: 0,
  overflow: 'hidden',
  display: 'flex' as const,
  flexDirection: 'column' as const,
  position: 'relative' as const,
  background: `linear-gradient(135deg, ${bgDark.value} 0%, ${bgDominant.value} 60%, ${bgMuted.value} 100%)`,
  color: '#fff',
  fontFamily: '-apple-system, BlinkMacSystemFont, sans-serif',
  userSelect: 'none' as const,
}))

const bgImageStyle = computed(() => ({
  position: 'absolute' as const,
  inset: 0,
  backgroundImage: coverUrl.value ? `url(${coverUrl.value})` : 'none',
  backgroundSize: 'cover' as const,
  backgroundPosition: 'center',
  opacity: 0.12,
  transition: 'opacity 0.5s',
}))

const bgOverlay = {
  position: 'absolute' as const,
  inset: 0,
  background: 'linear-gradient(to top, rgba(0,0,0,0.6) 0%, transparent 50%, rgba(0,0,0,0.2) 100%)',
}

const contentStyle = {
  position: 'relative' as const,
  zIndex: 1,
  display: 'flex' as const,
  flexDirection: 'column' as const,
  height: '100%',
  padding: '12px',
  boxSizing: 'border-box' as const,
}

const topSectionStyle = {
  display: 'flex' as const,
  alignItems: 'center',
  gap: '10px',
  flex: 1,
  minHeight: 0,
}

const coverWrapStyle = {
  width: '52px',
  height: '52px',
  borderRadius: '10px',
  overflow: 'hidden',
  flexShrink: 0,
  boxShadow: '0 2px 12px rgba(0,0,0,0.4)',
}

const coverFallbackStyle = {
  width: '100%',
  height: '100%',
  display: 'flex' as const,
  alignItems: 'center',
  justifyContent: 'center',
  background: 'rgba(255,255,255,0.05)',
}

const metaStyle = { flex: 1, minWidth: 0 }

const titleStyle = {
  fontWeight: 700,
  fontSize: '13px',
  whiteSpace: 'nowrap' as const,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
}

const artistStyle = {
  fontSize: '11px',
  color: 'rgba(255,255,255,0.5)',
  whiteSpace: 'nowrap' as const,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  marginTop: '2px',
}

const progressSectionStyle = {
  marginTop: '8px',
  display: 'flex' as const,
  flexDirection: 'column' as const,
  gap: '3px',
}

const progressBarStyle = {
  height: '3px',
  borderRadius: '2px',
  background: 'rgba(255,255,255,0.1)',
  overflow: 'hidden',
  cursor: 'pointer' as const,
}

const progressFillStyle = computed(() => ({
  height: '100%',
  borderRadius: '2px',
  background: `linear-gradient(to right, ${accentColor.value}, ${accentColor.value}dd)`,
  transition: 'width 0.2s linear',
}))

const progressLabelsStyle = {
  display: 'flex' as const,
  justifyContent: 'space-between' as const,
  fontSize: '9px',
  color: 'rgba(255,255,255,0.3)',
  fontVariantNumeric: 'tabular-nums' as const,
}

const controlsStyle = {
  display: 'flex' as const,
  alignItems: 'center',
  gap: '12px',
  marginTop: '8px',
}

const ctrlBtnStyle = {
  width: '26px',
  height: '26px',
  borderRadius: '50%',
  border: 'none',
  background: 'rgba(255,255,255,0.08)',
  color: 'rgba(255,255,255,0.6)',
  fontSize: '10px',
  cursor: 'pointer' as const,
  display: 'flex' as const,
  alignItems: 'center' as const,
  justifyContent: 'center' as const,
  transition: 'all 0.15s',
  flexShrink: 0 as const,
}

const playBtnBase = {
  width: '34px',
  height: '34px',
  borderRadius: '50%',
  border: 'none',
  color: '#fff',
  fontSize: '13px',
  cursor: 'pointer' as const,
  display: 'flex' as const,
  alignItems: 'center' as const,
  justifyContent: 'center' as const,
  transition: 'all 0.15s',
  flexShrink: 0 as const,
}

const closeBtnStyle = {
  width: '22px',
  height: '22px',
  borderRadius: '50%',
  border: 'none',
  background: 'rgba(255,255,255,0.08)',
  color: 'rgba(255,255,255,0.4)',
  fontSize: '9px',
  cursor: 'pointer' as const,
  display: 'flex' as const,
  alignItems: 'center' as const,
  justifyContent: 'center' as const,
  transition: 'opacity 0.2s',
  flexShrink: 0 as const,
}

const liveDotStyle = {
  position: 'absolute' as const,
  top: '6px',
  left: '6px',
  width: '5px',
  height: '5px',
  borderRadius: '50%',
  background: '#1db954',
  boxShadow: '0 0 6px rgba(29,185,84,0.6)',
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

const currentTimeLabel = computed(() => fmtTime(pc.currentTime.value))
const durationLabel = computed(() =>
  fmtTime(pc.duration.value || pc.currentTrack.value?.durationSeconds || 0),
)

function seekFromEvent(e: MouseEvent | KeyboardEvent) {
  const el = progressRef.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  const ratio = Math.max(0, Math.min(1, ((e as MouseEvent).clientX - rect.left) / rect.width))
  pc.seek(ratio * (pc.duration.value || pc.currentTrack.value?.durationSeconds || 0))
}

function close() { emit('close') }
</script>
