<template>
  <div :style="rootStyle">
    <!-- Background image + gradient overlay (same as FullscreenPlayer) -->
    <div :style="bgImageStyle" />
    <div :style="bgOverlay" />

    <div :style="contentStyle">
      <!-- Top: Cover + Title/Artist -->
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
            <span style="font-size: 22px; opacity: 0.4">♪</span>
          </div>
        </div>
        <div :style="metaStyle">
          <div :style="titleStyle">{{ title || 'No track' }}</div>
          <div :style="artistStyle">{{ artistName || '—' }}</div>
        </div>
      </div>

      <!-- Progress: Matches FullscreenPlayer style (glow + shadow + thumb) -->
      <div :style="progressSectionStyle">
        <div
          ref="progressRef"
          :style="progressTrackStyle"
          role="slider"
          tabindex="0"
          aria-label="Seek"
          :aria-valuemin="0"
          :aria-valuemax="pc.duration.value || 0"
          :aria-valuenow="pc.currentTime.value"
          @click="seekFromEvent"
          @keydown.enter="seekFromEvent"
          @keydown.space.prevent="seekFromEvent"
        >
          <!-- Glow layer behind fill (like FullscreenPlayer's blur layer) -->
          <div v-if="isPlaying" :style="progressGlowStyle" />
          <!-- Main fill with shadow (like FullscreenPlayer) -->
          <div :style="progressFillStyle" />
          <!-- Thumb dot (always visible in PiP, smaller) -->
          <div :style="progressThumbStyle" />
        </div>
        <div :style="progressLabelsStyle">
          <span>{{ currentTimeLabel }}</span>
          <span>{{ durationLabel }}</span>
        </div>
      </div>

      <!-- Controls -->
      <div :style="controlsStyle">
        <button :style="ctrlBtnStyle" aria-label="Previous track" @click="playPrevious" :disabled="hasPrevious!">⏮</button>
        <button
          :style="{ ...playBtnBase, background: accentColor, boxShadow: `0 0 0 3px ${accentColor}22, 0 4px 20px ${accentColor}44` }"
          :aria-label="isPlaying ? 'Pause' : 'Play'"
          @click="togglePlayPause"
        >
          <span :style="playIconStyle">{{ isPlaying ? '⏸' : '▶' }}</span>
        </button>
        <button :style="ctrlBtnStyle" aria-label="Next track" @click="playNext" :disabled="hasNext!">⏭</button>
        <div :style="{ flex: 1 }" />
        <button :style="closeBtnStyle" aria-label="Close PiP player" @click="close">✕</button>
      </div>
    </div>

    <!-- Live indicator -->
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

// ─── Root / Background ──────────────────────────────────────

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
  opacity: 0.15,
  transition: 'opacity 0.5s',
}))

const bgOverlay = {
  position: 'absolute' as const,
  inset: 0,
  background: 'linear-gradient(to top, rgba(0,0,0,0.7) 0%, transparent 55%, rgba(0,0,0,0.3) 100%)',
}

// ─── Content Layout ─────────────────────────────────────────

const contentStyle = {
  position: 'relative' as const,
  zIndex: 1,
  display: 'flex' as const,
  flexDirection: 'column' as const,
  height: '100%',
  padding: '10px 12px 10px',
  boxSizing: 'border-box' as const,
}

// ─── Top: Cover + Meta ──────────────────────────────────────

const topSectionStyle = {
  display: 'flex' as const,
  alignItems: 'center',
  gap: '10px',
  flex: 1,
  minHeight: 0,
}

const coverWrapStyle = {
  width: '60px',
  height: '60px',
  borderRadius: '12px',
  overflow: 'hidden',
  flexShrink: 0,
  boxShadow: '0 3px 16px rgba(0,0,0,0.5)',
}

const coverFallbackStyle = {
  width: '100%',
  height: '100%',
  display: 'flex' as const,
  alignItems: 'center',
  justifyContent: 'center',
  background: 'rgba(255,255,255,0.06)',
  borderRadius: '12px',
}

const metaStyle = { flex: 1, minWidth: 0 }

const titleStyle = {
  fontWeight: 700,
  fontSize: '14px',
  lineHeight: '1.2',
  whiteSpace: 'nowrap' as const,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
}

const artistStyle = {
  fontSize: '12px',
  color: 'rgba(255,255,255,0.5)',
  whiteSpace: 'nowrap' as const,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  marginTop: '3px',
}

// ─── Progress Bar (matches FullscreenPlayer style) ──────────

const progressSectionStyle = {
  marginTop: '6px',
  display: 'flex' as const,
  flexDirection: 'column' as const,
  gap: '3px',
}

const progressTrackStyle = {
  position: 'relative' as const,
  height: '4px',
  borderRadius: '3px',
  background: 'rgba(255,255,255,0.1)',
  cursor: 'pointer' as const,
  overflow: 'visible' as const,  // so thumb can overflow
}

/** Glow layer behind the fill (matches FullscreenPlayer's .blur-[3px] layer) */
const progressGlowStyle = computed(() => ({
  position: 'absolute' as const,
  left: 0,
  top: 0,
  height: '100%',
  width: `${progressPercent.value}%`,
  borderRadius: '3px',
  background: accentColor.value,
  opacity: 0.35,
  filter: 'blur(3px)',
  transition: 'width 0.2s linear',
}))

/** Main fill bar with box-shadow glow (matches FullscreenPlayer) */
const progressFillStyle = computed(() => ({
  position: 'absolute' as const,
  left: 0,
  top: 0,
  height: '100%',
  width: `${progressPercent.value}%`,
  borderRadius: '3px',
  background: `linear-gradient(to right, ${accentColor.value}, ${accentColor.value}dd)`,
  boxShadow: `0 0 8px ${accentColor.value}66`,
  transition: 'width 0.2s linear',
}))

/** Thumb dot at progress position (always visible in PiP) */
const progressThumbStyle = computed(() => ({
  position: 'absolute' as const,
  top: '50%',
  left: `calc(${progressPercent.value}% - 5px)`,
  transform: 'translateY(-50%)',
  width: '10px',
  height: '10px',
  borderRadius: '50%',
  background: accentColor.value,
  boxShadow: `0 0 0 2px ${accentColor.value}22, 0 2px 6px rgba(0,0,0,0.5)`,
  transition: 'left 0.2s linear',
}))

const progressLabelsStyle = {
  display: 'flex' as const,
  justifyContent: 'space-between' as const,
  fontSize: '9px',
  color: 'rgba(255,255,255,0.3)',
  fontVariantNumeric: 'tabular-nums' as const,
}

// ─── Controls ───────────────────────────────────────────────

const controlsStyle = {
  display: 'flex' as const,
  alignItems: 'center',
  gap: '14px',
  marginTop: '6px',
}

const ctrlBtnStyle = {
  width: '36px',
  height: '36px',
  borderRadius: '50%',
  border: 'none',
  background: 'rgba(255,255,255,0.08)',
  color: 'rgba(255,255,255,0.7)',
  fontSize: '14px',
  cursor: 'pointer' as const,
  display: 'flex' as const,
  alignItems: 'center' as const,
  justifyContent: 'center' as const,
  transition: 'all 0.15s',
  flexShrink: 0 as const,
  backdropFilter: 'blur(4px)',
}

const playBtnBase = {
  width: '44px',
  height: '44px',
  borderRadius: '50%',
  border: 'none',
  color: '#fff',
  fontSize: '14px',
  cursor: 'pointer' as const,
  display: 'flex' as const,
  alignItems: 'center' as const,
  justifyContent: 'center' as const,
  transition: 'all 0.15s',
  flexShrink: 0 as const,
}

const playIconStyle = {
  position: 'relative' as const,
  left: '1px',       // optical centering for ▶
}

const closeBtnStyle = {
  width: '28px',
  height: '28px',
  borderRadius: '50%',
  border: 'none',
  background: 'rgba(255,255,255,0.08)',
  color: 'rgba(255,255,255,0.5)',
  fontSize: '11px',
  cursor: 'pointer' as const,
  display: 'flex' as const,
  alignItems: 'center' as const,
  justifyContent: 'center' as const,
  transition: 'opacity 0.15s',
  flexShrink: 0 as const,
}

const liveDotStyle = {
  position: 'absolute' as const,
  top: '6px',
  left: '6px',
  width: '6px',
  height: '6px',
  borderRadius: '50%',
  background: '#1db954',
  boxShadow: '0 0 8px rgba(29,185,84,0.7)',
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
