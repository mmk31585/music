<template>
  <div
    class="relative z-10 flex justify-center items-center gap-3 px-6 pb-3 pt-1"
    :dir="dir"
  >
    <span class="text-[11px] text-white/50 font-mono tabular-nums w-10 text-right">{{ formatTime(currentTime) }}</span>
    
    <div
      role="slider"
      tabindex="0"
      :aria-valuemin="0"
      :aria-valuemax="duration"
      :aria-valuenow="currentTime"
      aria-label="Seek"
      class="group relative w-1/4 h-1.5 cursor-pointer rounded-full transition-all duration-150"
      :class="isPlaying ? 'bg-white/15' : 'bg-white/10'"
      dir="ltr"
      @click="onSeekClick"
      @keydown.enter="onSeekClick"
      @keydown.space.prevent="onSeekClick"
      @keydown.arrow-left.prevent="seekRelative(-5)"
      @keydown.arrow-right.prevent="seekRelative(5)"
      @keydown.home.prevent="seekPercent(0)"
      @keydown.end.prevent="seekPercent(100)"
    >
      <!-- Base track glow when playing -->
      <div
        v-if="isPlaying && progressColor"
        class="absolute inset-0 rounded-full opacity-20 blur-xs transition-opacity duration-500"
        :style="{ background: progressColor }"
      />
      
      <!-- Progress fill -->
      <div
        class="relative h-full rounded-full"
        :class="isPlaying ? 'progress-bar-fill' : ''"
        :style="progressStyle"
      />
      
      <!-- Thumb dot - only visible on hover/focus/drag -->
      <div
        class="absolute top-1/2 -translate-y-1/2 size-4 rounded-full shadow-xl opacity-0 group-hover:opacity-100 group-focus-visible:opacity-100 scale-0 group-hover:scale-100 group-focus-visible:scale-100 transition-all duration-300 ease-spring"
        :style="{
          left: `calc(${progressPercent}% - 8px)`,
          background: progressColor || 'var(--p-500)',
          boxShadow: progressColor ? `0 0 0 3px ${progressColor}22, 0 4px 12px rgba(0,0,0,0.5)` : '0 0 0 3px rgba(29,185,84,0.2), 0 4px 12px rgba(0,0,0,0.5)',
        }"
      />
    </div>
    
    <span class="text-[11px] text-white/50 font-mono tabular-nums w-10">{{ formatTime(duration) }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRTL } from '@/composables/useRTL'

interface Props {
  currentTime: number
  duration: number
  progressPercent: number
  isPlaying: boolean
  progressColor?: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  seek: [percent: number]
  'seek-relative': [seconds: number]
}>()

const { dir } = useRTL()

const progressStyle = computed(() => {
  if (!props.progressColor) return { width: `${props.progressPercent}%` }
  return {
    width: `${props.progressPercent}%`,
    background: props.progressColor,
    boxShadow: `0 0 8px ${props.progressColor}66, 0 0 20px ${props.progressColor}33`,
    transition: 'width 100ms linear, background 0.5s ease, box-shadow 0.3s ease',
  }
})

function formatTime(s: number) {
  if (!isFinite(s)) return '0:00'
  const m = Math.floor(s / 60)
  const sec = Math.floor(s % 60)
  return `${m}:${sec.toString().padStart(2, '0')}`
}

function onSeekClick(e: MouseEvent | KeyboardEvent) {
  const bar = e.currentTarget as HTMLElement
  const rect = bar.getBoundingClientRect()
  const pct = (((e as MouseEvent).clientX - rect.left) / rect.width) * 100
  emit('seek', pct)
}

function seekRelative(seconds: number) {
  emit('seek-relative', seconds)
}

function seekPercent(pct: number) {
  emit('seek', pct)
}
</script>

<style scoped>
/* Progress bar fill glow */
@keyframes progress-glow {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 1; }
}
.progress-bar-fill::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background: inherit;
  filter: blur(6px);
  opacity: 0.35;
  animation: progress-glow 2s ease-in-out infinite;
  z-index: -1;
}

/* Ease-spring utility for the thumb dot */
.ease-spring {
  transition-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1);
}

@media (prefers-reduced-motion: reduce) {
  .progress-bar-fill::after { animation: none; opacity: 0; }
  .group-hover\:opacity-100 { opacity: 1 !important; }
  .group-hover\:scale-100 { transform: translateY(-50%) scale(1) !important; }
  .group-focus-visible\:opacity-100 { opacity: 1 !important; }
  .group-focus-visible\:scale-100 { transform: translateY(-50%) scale(1) !important; }
}
</style>