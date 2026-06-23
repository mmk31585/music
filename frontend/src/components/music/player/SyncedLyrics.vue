<template>
  <div
    ref="containerRef"
    class="relative h-full w-full overflow-y-auto scrollbar-none"
  >
    <!-- Empty state -->
    <div
      v-if="lines.length === 0"
      class="flex h-full flex-col items-center justify-center gap-3 px-6 text-center"
    >
      <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-white/[0.03]">
        <i aria-hidden="true" class="pi pi-align-left text-2xl text-white/[0.12]" />
      </div>
      <p class="text-sm font-medium text-white/20">No synced lyrics</p>
      <p class="text-xs text-white/10">Lyrics will appear here when available</p>
    </div>

    <!-- Scrolling lyrics list -->
    <div v-else ref="scrollRef" class="flex min-h-full flex-col items-center justify-center py-16 md:py-20">
      <div
        v-for="(line, idx) in lines"
        :key="idx"
        :ref="(el) => { if (el) lineRefs[idx] = el as HTMLElement }"
        role="button"
        tabindex="0"
        class="w-full max-w-2xl cursor-pointer px-8 py-1.5 select-none"
        :class="isRtl?'text-right':'text-left'"
        :style="lineStyle(idx)"
        @click="emit('seek', line.timeSeconds)"
        @keydown.enter="emit('seek', line.timeSeconds)"
        @keydown.space.prevent="emit('seek', line.timeSeconds)"
      >
        <div
          class="relative inline-block text-center leading-[1.8]"
          :class="{
            'font-bold': idx === activeIdx,
          }"
          :style="{
            fontSize: idx === activeIdx ? '1.45rem' : '1.05rem',
            letterSpacing: idx === activeIdx ? '-0.01em' : '0em',
          }"
        >
          <!-- Active line: glow pill + particle shimmer -->
          <div
            v-if="idx === activeIdx"
            class="absolute -inset-x-6 -inset-y-2 rounded-2xl transition-all duration-700 ease-out"
            :class="isRtl ? 'glow-rtl' : 'glow-ltr'"
            :style="{
              background: `radial-gradient(ellipse at ${isRtl ? '80% 50%' : '20% 50%'}, ${props.activeColor || '#a855f7'}18 0%, transparent 65%)`,
              opacity: 1,
              transform: 'scale(1)',
            }"
          />
          <!-- Label for line number (karaoke-style marker) — flips for RTL -->
          <span
            v-if="idx === activeIdx"
            :class="[
              'absolute top-1/2 -translate-y-1/2 text-[10px] font-bold tracking-wider transition-all duration-500',
              isRtl ? '-right-8' : '-left-8',
            ]"
            :style="{ color: props.activeColor || '#a855f7', opacity: 0.6 }"
          >♪</span>

          <!-- Karaoke word highlight on active line -->
          <template v-if="idx === activeIdx && hasWordTimings(line)">
            <span
              v-for="(word, wIdx) in line.words"
              :key="wIdx"
              class="transition-all duration-[80ms] ease-linear relative"
              :style="wordStyle(line, wIdx)"
            >{{ word.text }}&nbsp;</span>
          </template>
          <template v-else>
            <span class="relative transition-all duration-300">{{ line.text }}</span>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import type { ParsedLine } from '@/composables/lyrics'

const props = withDefaults(defineProps<{
  lines: ParsedLine[]
  currentTime: number
  duration: number
  activeColor?: string
  mutedColor?: string
  language?: string
}>(), {
  language: 'en',
})

const emit = defineEmits<{
  seek: [seconds: number]
}>()

const containerRef = ref<HTMLElement | null>(null)
const lineRefs = ref<HTMLElement[]>([])

// ── Active line index ────────────────────────
const activeIdx = computed(() => {
  const t = props.currentTime
  for (let i = props.lines.length - 1; i >= 0; i--) {
    if (t >= props.lines[i]!.timeSeconds) return i
  }
  // Plain lyrics: estimate position
  if (props.lines.length > 0 && props.duration > 0 && props.lines[0]!.timeSeconds === 0) {
    const progress = t / props.duration
    return Math.min(Math.floor(progress * props.lines.length), props.lines.length - 1)
  }
  return -1
})

// ── Clear refs when lines change ────────────
watch(() => props.lines, () => {
  lineRefs.value = []
}, { flush: 'post' })
const scrollRef = ref<HTMLElement | null>(null)

// ── Auto-scroll to active line (smooth, with debounce) ─────────
let scrollTimer: ReturnType<typeof setTimeout> | null = null

watch(activeIdx, (idx) => {
  if (scrollTimer) clearTimeout(scrollTimer)
  scrollTimer = setTimeout(() => {
    if (idx < 0 || !containerRef.value) return
    const target = lineRefs.value[idx]
    if (target) {
      target.scrollIntoView({ block: 'center', behavior: 'smooth' })
    }
  }, 80)
})

onBeforeUnmount(() => {
  if (scrollTimer) clearTimeout(scrollTimer)
})

// ── Line opacity/style based on distance ─────
function lineStyle(idx: number): Record<string, string> {
  const active = activeIdx.value
  if (active < 0) {
    return {
      opacity: '0.12',
      color: '#fff',
      transform: 'scale(0.92) translateY(0)',
      transition: 'all 0.6s cubic-bezier(0.22, 1, 0.36, 1)',
    }
  }

  const dist = Math.abs(idx - active)

  // Active line — bold, bright, dramatic
  if (dist === 0) {
    return {
      opacity: '1',
      color: '#fff',
      textShadow: `0 0 50px ${props.activeColor || 'rgba(168,85,247,0.35)'}, 0 0 100px ${props.activeColor || 'rgba(168,85,247,0.12)'}`,
      transform: 'scale(1) translateY(0)',
      filter: 'brightness(1.35) contrast(1.1)',
      transition: 'all 0.5s cubic-bezier(0.22, 1, 0.36, 1)',
    }
  }

  // Past lines (already sung) — fade out gently
  if (idx < active) {
    const fade = Math.max(0.03, 0.3 - dist * 0.12)
    return {
      opacity: String(fade),
      color: 'rgba(255,255,255,0.7)',
      transform: `scale(${Math.max(0.82, 1 - dist * 0.06)}) translateY(${dist * 2}px)`,
      filter: 'blur(0.5px)',
      transition: 'all 0.7s cubic-bezier(0.22, 1, 0.36, 1)',
    }
  }

  // Upcoming lines — dim, smaller, with slight upward anticipation
  const fade = Math.max(0.04, 0.35 - dist * 0.1)
  return {
    opacity: String(fade),
    color: 'rgba(255,255,255,0.9)',
    transform: `scale(${Math.max(0.84, 1 - dist * 0.05)}) translateY(${-dist * 1.5}px)`,
    transition: 'all 0.6s cubic-bezier(0.22, 1, 0.36, 1)',
  }
}

// ── Word-level karaoke style (smooth gradient reveal) ─────────
function hasWordTimings(line: ParsedLine): boolean {
  return line.words.length > 0 && line.words[0]!.timeSeconds >= 0
}

function wordStyle(line: ParsedLine, wordIdx: number): Record<string, string> {
  const t = props.currentTime
  const word = line.words[wordIdx]
  if (!word) return { color: 'rgba(255,255,255,0.25)' }

  const nextWord = line.words[wordIdx + 1]
  const start = word.timeSeconds >= 0 ? word.timeSeconds : line.timeSeconds
  const end = nextWord?.timeSeconds ?? line.timeSeconds + 4

  const isActive = t >= start && t < end
  const progress = isActive && end > start ? Math.min(1, (t - start) / (end - start)) : (isActive ? 1 : 0)

  // Smooth interpolation from dim to bright with gradient-like feel
  const brightness = isActive ? 1 : 0.3

  return {
    color: isActive ? '#fff' : 'rgba(255,255,255,0.25)',
    textShadow: isActive
      ? `0 0 ${12 + progress * 15}px ${props.activeColor || 'rgba(168,85,247,0.3)'}`
      : 'none',
    transition: 'color 80ms ease-out, text-shadow 80ms ease-out',
  }
}
const isRtl = computed(() => {
  const lang = props.language?.toLowerCase()
  return lang === 'fa' || lang === 'far' || lang?.startsWith('fa-') || lang?.startsWith('fa_')
})
console.log(isRtl)
console.log(props)
</script>

<style scoped>
.scrollbar-none {
  scrollbar-width: none;
  -ms-overflow-style: none;
}
.scrollbar-none::-webkit-scrollbar {
  display: none;
}

/* RTL vs LTR glow pill shifting */
.glow-ltr {
  transform-origin: left center;
}
.glow-rtl {
  transform-origin: right center;
}
.glow-rtl {
  /* Slightly wider glow for Persian text which can be more expansive */
  --glow-scale: 1.02;
}

/* Active line text in RTL gets a different shadow direction */
.glow-rtl + .line-text,
[dir="rtl"] .line-active-text {
  text-shadow: 0 0 50px var(--glow-color, rgba(168,85,247,0.35)), 0 0 100px var(--glow-color, rgba(168,85,247,0.12));
}

/* Word karaoke highlight RTL adjustment */
:deep(.rtl-word-highlight) {
  background: linear-gradient(to left, currentColor 0%, transparent 100%);
  -webkit-background-clip: text;
  background-clip: text;
}

@media (prefers-reduced-motion: reduce) {
  * {
    transition-duration: 0.01ms !important;
    animation-duration: 0.01ms !important;
  }
}
</style>
