<template>
  <div class="flex flex-col w-full items-center gap-1 px-4 pb-2" role="group" aria-label="Playback controls">
    <!-- Playback Controls Row -->
    <div class="flex items-center justify-center gap-1.5" role="group" aria-label="Playback controls">
      <!-- Shuffle -->
      <div class="relative">
        <button
          ref="shuffleBtnRef"
          type="button"
          class="relative flex h-9 w-9 items-center justify-center rounded-full transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-spotify/60 focus-visible:ring-offset-2 focus-visible:ring-offset-surface-base"
          :class="[
            props.shuffleMode !== 'off'
              ? 'text-aurora-purple hover:bg-aurora-purple/15 active:scale-95'
              : 'text-white/40 hover:text-white hover:bg-white/10 active:scale-95'
          ]"
          :disabled="!props.currentTrack"
          :aria-label="shuffleAriaLabel"
          @click="showShuffleMenu = !showShuffleMenu"
          @keydown.escape="showShuffleMenu = false"
        >
          <Shuffle class="h-4 w-4" aria-hidden="true" />
          <span
            v-if="props.shuffleMode !== 'off'"
            class="absolute -top-0.5 -right-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-aurora-purple text-[8px] font-bold text-white"
          >
            {{ shuffleIndicator }}
          </span>
        </button>

        <!-- Shuffle mode selector popup -->
        <Transition name="fade-scale">
          <div
            v-if="showShuffleMenu"
            ref="shuffleMenuRef"
            role="menu"
            aria-label="Select shuffle mode"
            class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 z-60 min-w-[150px] rounded-xl border border-white/10 bg-surface-raised p-1.5 shadow-[0_8px_32px_rgba(0,0,0,0.6)] backdrop-blur-2xl"
            style="backdrop-filter: blur(24px);"
            @keydown="onMenuKeydown($event, 'shuffle')"
          >
            <button
              v-for="(mode, idx) in shuffleModes"
              :key="mode.value"
              :ref="(el) => { if (el) shuffleItemRefs[idx] = el as HTMLElement }"
              type="button"
              role="menuitem"
              :tabindex="idx === 0 ? 0 : -1"
              class="flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-xs font-bold transition hover:bg-white/8 focus-visible:outline-none focus-visible:bg-white/10"
              :class="props.shuffleMode === mode.value ? 'text-aurora-purple bg-white/6' : 'text-slate-400 hover:text-white'"
              @click="setShuffleMode(mode.value)"
            >
              <component :is="mode.icon" class="h-4 w-4" aria-hidden="true" />
              <span class="flex-1 text-left">{{ mode.label }}</span>
              <span v-if="props.shuffleMode === mode.value" class="h-2 w-2 rounded-full bg-aurora-purple" />
            </button>
          </div>
        </Transition>
      </div>

      <!-- Previous -->
      <button
        type="button"
        class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 hover:text-white disabled:opacity-25 hover:bg-white/10 active:scale-95 transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-spotify/60 focus-visible:ring-offset-2 focus-visible:ring-offset-surface-base"
        :disabled="!props.hasPrevious"
        aria-label="Previous track"
        @click="props.playPrevious"
      >
        <SkipBack class="h-5 w-5" aria-hidden="true" />
      </button>

      <!-- Play/Pause - Main Control -->
      <GuestPlayGate action="play" #default="{ proceed }">
        <button
          type="button"
          class="relative flex h-11 w-11 items-center justify-center rounded-full shadow-xl disabled:opacity-40 transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-spotify/60 focus-visible:ring-offset-2 focus-visible:ring-offset-surface-base"
          :style="{ background: props.progressColor }"
          :disabled="!props.currentTrack || props.isLoadingTrack"
          :aria-label="playAriaLabel"
          @click="props.isPlaying ? props.togglePlayPause() : proceed()"
        >
          <span v-if="props.isLoadingTrack || props.isBuffering" class="absolute inset-0 animate-spin-slow">
            <Loader class="h-5 w-5 text-white" aria-hidden="true" />
          </span>
          <Pause v-else-if="props.isPlaying" class="absolute h-5 w-5 text-white" aria-hidden="true" />
          <Play v-else class="absolute h-5 w-5 text-white ml-0.5" aria-hidden="true" />
        </button>
      </GuestPlayGate>

      <!-- Next -->
      <button
        type="button"
        class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 hover:text-white disabled:opacity-25 hover:bg-white/10 active:scale-95 transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-spotify/60 focus-visible:ring-offset-2 focus-visible:ring-offset-surface-base"
        :disabled="!props.hasNext"
        aria-label="Next track"
        @click="props.playNext"
      >
        <SkipForward class="h-5 w-5" aria-hidden="true" />
      </button>

      <!-- Repeat -->
      <div class="relative">
        <button
          ref="repeatBtnRef"
          type="button"
          class="relative flex h-9 w-9 items-center justify-center rounded-full transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-spotify/60 focus-visible:ring-offset-2 focus-visible:ring-offset-surface-base"
          :class="[
            props.repeatMode !== 'off'
              ? 'text-aurora-pink hover:bg-aurora-pink/15 active:scale-95'
              : 'text-white/40 hover:text-white hover:bg-white/10 active:scale-95'
          ]"
          :disabled="!props.currentTrack"
          :aria-label="repeatAriaLabel"
          @click="showRepeatMenu = !showRepeatMenu"
          @keydown.escape="showRepeatMenu = false"
        >
          <Repeat class="h-4 w-4" aria-hidden="true" />
          <span
            v-if="props.repeatMode === 'one'"
            class="absolute -top-0.5 -right-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-aurora-pink text-[8px] font-bold text-black"
          >
            1
          </span>
        </button>

        <!-- Repeat mode selector popup -->
        <Transition name="fade-scale">
          <div
            v-if="showRepeatMenu"
            ref="repeatMenuRef"
            role="menu"
            aria-label="Select repeat mode"
            class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 z-60 min-w-[130px] rounded-xl border border-white/10 bg-surface-raised p-1.5 shadow-[0_8px_32px_rgba(0,0,0,0.6)] backdrop-blur-2xl"
            style="backdrop-filter: blur(24px);"
            @keydown="onMenuKeydown($event, 'repeat')"
          >
            <button
              v-for="(mode, idx) in repeatModes"
              :key="mode.value"
              :ref="(el) => { if (el) repeatItemRefs[idx] = el as HTMLElement }"
              type="button"
              role="menuitem"
              :tabindex="idx === 0 ? 0 : -1"
              class="flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-xs font-bold transition hover:bg-white/8 focus-visible:outline-none focus-visible:bg-white/10"
              :class="props.repeatMode === mode.value ? 'text-aurora-pink bg-white/6' : 'text-slate-400 hover:text-white'"
              @click="setRepeatMode(mode.value)"
            >
            <component :is="mode.icon" class="h-4 w-4" aria-hidden="true" />
            <span class="flex-1 text-left">{{ mode.label }}</span>
            <span v-if="props.repeatMode === mode.value" class="h-2 w-2 rounded-full bg-aurora-pink" />
          </button>
        </div>
      </Transition>
    </div>
  </div>
  
  <!-- Progress Bar / Timeline -->
  <div class="w-full flex justify-center items-center gap-3">
      <span class="text-[11px] text-white/50 font-mono tabular-nums">{{ formatTime(props.currentTime) }}</span>
      <div
        role="slider"
        tabindex="0"
        :aria-valuemin="0"
        :aria-valuemax="props.duration"
        :aria-valuenow="props.currentTime"
        aria-label="Seek"
        class="group relative h-1.5 flex-1 cursor-pointer rounded-full transition-all duration-150"
        :class="props.isPlaying ? 'bg-white/15' : 'bg-white/10'"
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
      v-if="props.isPlaying && props.progressColor"
      class="absolute inset-0 rounded-full opacity-20 blur-xs transition-opacity duration-500"
      :style="{ background: props.progressColor }"
      />
      
      <!-- Progress fill -->
      <div
      class="relative h-full rounded-full"
        :class="props.isPlaying ? 'progress-bar-fill' : ''"
        :style="progressStyle"
        />
        
        <!-- Thumb dot - only visible on hover/focus/drag -->
        <div
        class="absolute top-1/2 -translate-y-1/2 size-4 rounded-full shadow-xl opacity-0 group-hover:opacity-100 group-focus-visible:opacity-100 scale-0 group-hover:scale-100 group-focus-visible:scale-100 transition-all duration-300 ease-spring"
        :style="{
          left: `calc(${props.progressPercent}% - 8px)`,
          background: props.progressColor || 'var(--p-500)',
          boxShadow: props.progressColor ? `0 0 0 3px ${props.progressColor}22, 0 4px 12px rgba(0,0,0,0.5)` : '0 0 0 3px rgba(29,185,84,0.2), 0 4px 12px rgba(0,0,0,0.5)',
        }"
        />
      </div>
      <span class="text-[11px] text-white/50 font-mono tabular-nums">{{ formatTime(props.duration) }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { Ban, Globe, Loader, Pause, Play, Repeat, RotateCcw, Shuffle, SkipBack, SkipForward, Star } from 'lucide-vue-next'
import GuestPlayGate from '@/components/common/GuestPlayGate.vue'

const props = defineProps<{
  currentTrack: any
  isPlaying: boolean
  isBuffering: boolean
  isLoadingTrack: boolean
  hasNext: boolean
  hasPrevious: boolean
  shuffleMode: 'off' | 'queue' | 'catalog' | 'similar'
  repeatMode: 'off' | 'all' | 'one'
  progressColor: string
  currentTime: number
  duration: number
  progressPercent: number
  togglePlayPause: () => void
  playNext: () => void
  playPrevious: () => void
  setShuffleMode: (mode: 'off' | 'queue' | 'catalog' | 'similar') => void
  setRepeatMode: (mode: 'off' | 'all' | 'one') => void
}>()

const emit = defineEmits<{
  'toggle-fullscreen': []
  'toggle-lyrics': []
  seek: [percent: number]
  'seek-relative': [seconds: number]
}>()

const showShuffleMenu = ref(false)
const shuffleBtnRef = ref<HTMLElement | null>(null)
const shuffleMenuRef = ref<HTMLElement | null>(null)
const shuffleItemRefs = ref<HTMLElement[]>([])
const showRepeatMenu = ref(false)
const repeatBtnRef = ref<HTMLElement | null>(null)
const repeatMenuRef = ref<HTMLElement | null>(null)
const repeatItemRefs = ref<HTMLElement[]>([])

watch(showShuffleMenu, async (v) => {
  if (v) { await nextTick(); shuffleItemRefs.value[0]?.focus() }
})
watch(showRepeatMenu, async (v) => {
  if (v) { await nextTick(); repeatItemRefs.value[0]?.focus() }
})

function onMenuKeydown(e: KeyboardEvent, type: 'shuffle' | 'repeat') {
  const items = type === 'shuffle' ? shuffleItemRefs.value : repeatItemRefs.value
  const currentIdx = items.findIndex((el) => el === document.activeElement)

  switch (e.key) {
    case 'ArrowDown': {
      e.preventDefault()
      const next = (currentIdx + 1) % items.length
      items[next]?.focus()
      break
    }
    case 'ArrowUp': {
      e.preventDefault()
      const prev = (currentIdx - 1 + items.length) % items.length
      items[prev]?.focus()
      break
    }
    case 'Home': {
      e.preventDefault()
      items[0]?.focus()
      break
    }
    case 'End': {
      e.preventDefault()
      items[items.length - 1]?.focus()
      break
    }
    case 'Escape':
    case 'Tab': {
      if (type === 'shuffle') showShuffleMenu.value = false
      else showRepeatMenu.value = false
      const trigger = type === 'shuffle' ? shuffleBtnRef.value : repeatBtnRef.value
      trigger?.focus()
      break
    }
  }
}

const shuffleModes = [
  { value: 'off' as const, label: 'Off', icon: Ban },
  { value: 'queue' as const, label: 'Shuffle Queue', icon: Shuffle },
  { value: 'catalog' as const, label: 'Random Catalog', icon: Globe },
  { value: 'similar' as const, label: 'Similar Tracks', icon: Star },
]

const repeatModes = [
  { value: 'off' as const, label: 'No Repeat', icon: RotateCcw },
  { value: 'all' as const, label: 'Repeat All', icon: Repeat },
  { value: 'one' as const, label: 'Repeat One', icon: RotateCcw },
]

function setShuffleMode(mode: 'off' | 'queue' | 'catalog' | 'similar') {
  props.setShuffleMode(mode)
  showShuffleMenu.value = false
}

const playAriaLabel = computed(() => {
  if (props.isLoadingTrack || props.isBuffering) return 'Loading'
  return props.isPlaying ? 'Pause' : 'Play'
})

const shuffleAriaLabel = computed(() => {
  switch (props.shuffleMode) {
    case 'queue': return 'Shuffle queue'
    case 'catalog': return 'Random catalog tracks'
    case 'similar': return 'Similar tracks'
    default: return 'Shuffle off'
  }
})

const repeatAriaLabel = computed(() => {
  switch (props.repeatMode) {
    case 'all': return 'Repeat all'
    case 'one': return 'Repeat one'
    default: return 'Repeat off'
  }
})

const shuffleIndicator = computed(() => {
  switch (props.shuffleMode) {
    case 'queue': return 'Q'
    case 'catalog': return 'R'
    case 'similar': return 'S'
    default: return ''
  }
})

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

// Close popup menus on outside click
onMounted(() => {
  document.addEventListener('click', handleOutsideClick)
})
onUnmounted(() => {
  document.removeEventListener('click', handleOutsideClick)
})
function handleOutsideClick(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (showShuffleMenu.value && shuffleBtnRef.value && !shuffleBtnRef.value.contains(target)) {
    showShuffleMenu.value = false
  }
  if (showRepeatMenu.value && repeatBtnRef.value && !repeatBtnRef.value.contains(target)) {
    showRepeatMenu.value = false
  }
}
</script>

<style scoped>
@keyframes spin-slow {
  to { transform: rotate(360deg); }
}

.animate-spin-slow {
  animation: spin-slow 1.5s linear infinite;
}

.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: opacity 200ms cubic-bezier(0.34, 1.56, 0.64, 1), transform 200ms cubic-bezier(0.34, 1.56, 0.64, 1);
}
.fade-scale-enter-from {
  opacity: 0;
  transform: scale(0.95) translateY(4px);
}
.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

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

.ease-spring {
  transition-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1);
}

@media (prefers-reduced-motion: reduce) {
  .fade-scale-enter-active,
  .fade-scale-leave-active {
    transition: none;
  }
  .fade-scale-enter-from,
  .fade-scale-leave-to {
    opacity: 1;
    transform: none;
  }
  .progress-bar-fill::after {
    animation: none;
    opacity: 0;
  }
  .animate-spin-slow {
    animation: none;
  }
}
</style>