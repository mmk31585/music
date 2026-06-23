<template>
  <div class="relative h-full w-full overflow-hidden">
    <!-- Particle canvas -->
    <canvas ref="particleCanvas" class="pointer-events-none absolute inset-0 z-0" />

    <div v-if="loading" class="relative z-10 flex h-full items-center justify-center px-6">
      <div class="w-3/4 space-y-5">
        <div
          v-for="i in 6"
          :key="i"
          class="shimmer h-5 rounded bg-white/[0.06]"
          :style="{ width: `${55 + ((i * 7) % 30)}%` }"
        />
      </div>
    </div>

    <div
      v-else-if="!content"
      class="relative z-10 flex h-full flex-col items-center justify-center gap-4 px-6 text-center"
    >
      <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/5">
        <i aria-hidden="true" class="pi pi-align-left text-3xl text-white/15" />
      </div>
      <p class="text-sm text-white/25">No lyrics available</p>
      <p class="text-xs text-white/15">Lyrics will appear here when available</p>
    </div>

    <div v-else ref="containerRef" class="relative z-10 h-full scrollbar-none overflow-y-auto px-6" :dir="isRtl ? 'rtl' : 'ltr'">
      <div class="flex min-h-full flex-col items-center justify-center py-12">
        <div
          v-for="(line, idx) in parsedCache"
          :key="idx"
          ref="lineRefs"
          role="button"
          tabindex="0"
          class="cursor-pointer px-4 py-3 text-center text-lg leading-relaxed transition-all duration-500 ease-out md:text-xl"
          :class="{
            'scale-105 font-bold text-white drop-shadow-[0_0_20px_rgba(255,255,255,0.15)]':
              idx === activeLineIdx,
            'text-white/15 hover:text-white/35': isPast(idx),
            'text-white/25 hover:text-white/50': idx !== activeLineIdx && !isPast(idx),
          }"
          @click="onLineClick(line.timeSeconds)"
          @keydown.enter="onLineClick(line.timeSeconds)"
          @keydown.space.prevent="onLineClick(line.timeSeconds)"
        >
          <template v-if="line.isActive && karaoke">
            <span
              v-for="(word, wIdx) in line.words"
              :key="wIdx"
              class="transition-all duration-[50ms] ease-linear"
              :class="
                isWordActive(line, wIdx)
                  ? 'text-[#1db954] drop-shadow-[0_0_12px_rgba(29,185,84,0.6)]'
                  : 'text-white/40'
              "
            >
              {{ word.text }}<template v-if="wIdx < line.words.length - 1">&nbsp;</template>
            </span>
          </template>
          <template v-else>
            <span
              class="transition-all duration-300"
              :class="
                line.isActive ? 'text-white drop-shadow-[0_0_30px_rgba(255,255,255,0.1)]' : ''
              "
            >
              {{ line.text }}
            </span>
          </template>
        </div>
      </div>
    </div>

    <!-- Karaoke badge -->
    <div v-if="karaoke && content" class="absolute top-4 left-1/2 z-20 -translate-x-1/2">
      <span
        class="inline-flex items-center gap-1.5 rounded-full border border-[#1db954]/20 bg-[#1db954]/10 px-3 py-1 text-[10px] font-bold tracking-wider text-[#1db954] uppercase backdrop-blur-sm"
      >
        <span class="glow-spread flex h-1.5 w-1.5 rounded-full bg-[#1db954]" />
        Karaoke
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { parseLRCLines, parsePlainLines, useLyricsParticles } from '@/composables/lyrics'
import type { ParsedLine } from '@/composables/lyrics'

const props = withDefaults(
  defineProps<{
    content?: string
    type?: string
    language?: string
    currentTime?: number
    loading?: boolean
    karaoke?: boolean
  }>(),
  {
    content: '',
    type: 'plain',
    language: 'en',
    currentTime: 0,
    loading: false,
    karaoke: false,
  },
)

const isRtl = computed(() => {
  const lang = props.language?.toLowerCase()
  return lang === 'fa' || lang === 'far' || lang?.startsWith('fa-') || lang?.startsWith('fa_')
})

const emit = defineEmits<{
  seek: [seconds: number]
}>()

const containerRef = ref<HTMLElement | null>(null)
const lineRefs = ref<HTMLElement[]>([])
const particleCanvas = ref<HTMLCanvasElement | null>(null)
const particleAnim = useLyricsParticles(particleCanvas)

let parsedCache: ParsedLine[] = []

watch(
  () => [props.type, props.content],
  () => {
    if (!props.content) {
      parsedCache = []
      return
    }
    parsedCache =
      props.type === 'lrc' ? parseLRCLines(props.content) : parsePlainLines(props.content)
  },
  { immediate: true },
)

const activeLineIdx = computed(() => {
  const t = props.currentTime
  for (let i = parsedCache.length - 1; i >= 0; i--) {
    if (t >= parsedCache[i]!.timeSeconds) return i
  }
  return -1
})

function isPast(idx: number) {
  return activeLineIdx.value >= 0 && idx < activeLineIdx.value
}

function isWordActive(line: ParsedLine, wordIdx: number) {
  if (!line.words || line.words.length === 0) return true
  const t = props.currentTime
  const word = line.words[wordIdx]
  const nextWord = line.words[wordIdx + 1]
  if (!word) return false
  const start = word.timeSeconds >= 0 ? word.timeSeconds : line.timeSeconds
  const end = nextWord?.timeSeconds ?? line.timeSeconds + 4
  return t >= start && t < end
}

let autoScrollTimer: ReturnType<typeof setTimeout> | null = null

watch(activeLineIdx, (idx) => {
  if (autoScrollTimer) clearTimeout(autoScrollTimer)
  autoScrollTimer = setTimeout(() => {
    if (idx < 0 || !containerRef.value) return
    const target = lineRefs.value[idx]
    target?.scrollIntoView({ block: 'center', behavior: 'smooth' })
  }, 80)
})

function onLineClick(timeSeconds: number) {
  emit('seek', timeSeconds)
}

watch(
  () => props.karaoke,
  async (val) => {
    if (val) {
      await nextTick()
      particleAnim.init()
    } else {
      particleAnim.stop()
    }
  },
)

onBeforeUnmount(() => {
  particleAnim.stop()
  if (autoScrollTimer) clearTimeout(autoScrollTimer)
})
</script>
