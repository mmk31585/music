<template>
  <Teleport to="body">
    <Transition name="ambient-enter">
      <div
        v-if="visible"
        ref="rootEl"
        class="fixed inset-0 z-[200] flex flex-col overflow-hidden select-none"
        :style="bgStyle"
        @keydown="onKeydown"
        tabindex="0"
      >
        <!-- Blurred album art background -->
        <div class="pointer-events-none absolute inset-0 scale-110">
          <img
            v-if="coverUrl"
            :src="coverUrl"
            alt=""
            loading="lazy"
            class="h-full w-full object-cover opacity-70 blur-3xl"
          />
          <div
            v-else
            class="h-full w-full bg-gradient-to-br from-[#0a0a0a] via-[#121212] to-[#1a1a2e]"
          />
        </div>

        <!-- Dark gradient overlay -->
        <div
          class="pointer-events-none absolute inset-0 bg-gradient-to-b from-black/60 via-black/30 to-black/80"
        />

        <!-- Ambient particles layer -->
        <canvas ref="particleCanvas" class="pointer-events-none absolute inset-0 z-0" />

        <!-- Content -->
        <div class="relative z-10 flex flex-1 flex-col">
          <!-- Top bar -->
          <div class="flex items-center justify-between px-5 pt-5 pb-3">
            <button
              type="button"
              class="spring flex h-10 w-10 items-center justify-center rounded-full text-white/60 backdrop-blur-sm transition-all hover:bg-white/10 hover:text-white"
              @click="close"
            >
              <i class="pi pi-chevron-down text-lg" />
            </button>
            <div class="glass flex items-center gap-2 rounded-full px-4 py-2 text-xs text-white/50">
              <span
                class="flex h-2 w-2 rounded-full bg-[#1db954]"
                :class="{ 'glow-spread': isPlaying }"
              />
              <span class="font-semibold tracking-wider uppercase">Ambient Mode</span>
            </div>
            <div class="w-10" />
          </div>

          <!-- Center: Album art + info (minimal) -->
          <div class="flex flex-1 flex-col items-center justify-center gap-5 px-6">
            <div class="relative h-48 w-48 md:h-64 md:w-64 lg:h-80 lg:w-80">
              <div
                class="h-full w-full overflow-hidden rounded-3xl shadow-2xl ring-1 shadow-black/50 ring-white/10"
              >
                <img
                  v-if="coverUrl"
                  :src="coverUrl"
                  :alt="title"
                  class="h-full w-full object-cover transition-transform duration-[3000ms] ease-out"
                  :class="isPlaying ? 'scale-110' : 'scale-100'"
                />
                <div
                  v-else
                  class="flex h-full w-full items-center justify-center bg-gradient-to-br from-[#1db954]/30 to-[#121212]"
                >
                  <i class="pi pi-music text-5xl text-white/20" />
                </div>
              </div>
            </div>

            <div class="text-center">
              <p class="text-2xl font-bold text-white drop-shadow-lg">{{ title }}</p>
              <p class="mt-1 text-sm text-white/50 drop-shadow">{{ artistName }}</p>
            </div>
          </div>

          <!-- Bottom controls (minimal, translucent) -->
          <div class="relative z-10 flex flex-col items-center gap-3 px-6 pb-8">
            <!-- Seekbar -->
            <input
              type="range"
              min="0"
              max="100"
              step="0.1"
              class="fullscreen-range w-full max-w-md"
              :style="progressStyle"
              :value="progressPercent"
              :disabled="!currentTrack"
              @input="onSeek"
            />

            <div class="flex items-center gap-6">
              <button
                type="button"
                class="spring flex h-10 w-10 items-center justify-center rounded-full text-white/40 backdrop-blur-sm transition-all hover:bg-white/15 hover:text-white"
                :class="{ '!text-[#1db954]': shuffleMode }"
                :disabled="!currentTrack"
                @click="toggleShuffle"
              >
                <i class="pi pi-sort-alt text-sm" />
              </button>

              <button
                type="button"
                class="spring flex h-12 w-12 items-center justify-center rounded-full text-white/50 backdrop-blur-sm transition-all hover:bg-white/15 hover:text-white disabled:opacity-20"
                :disabled="!hasPrevious"
                @click="playPrevious"
              >
                <i class="pi pi-step-backward text-xl" />
              </button>

              <button
                type="button"
                class="glow-green spring relative flex h-16 w-16 items-center justify-center rounded-full bg-white/90 text-black shadow-2xl backdrop-blur-sm transition-all hover:scale-105 hover:bg-[#1db954] hover:text-white disabled:opacity-40"
                :class="{ '!bg-[#1db954] !text-white': isPlaying }"
                :disabled="!currentTrack || isLoadingTrack"
                @click="togglePlayPause"
              >
                <i v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-xl" />
                <i
                  v-else
                  :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'"
                  class="text-xl"
                />
                <div
                  v-if="isPlaying"
                  class="absolute -inset-2 animate-ping rounded-full border-2 border-[#1db954]/30"
                />
              </button>

              <button
                type="button"
                class="spring flex h-12 w-12 items-center justify-center rounded-full text-white/50 backdrop-blur-sm transition-all hover:bg-white/15 hover:text-white disabled:opacity-20"
                :disabled="!hasNext"
                @click="playNext"
              >
                <i class="pi pi-step-forward text-xl" />
              </button>

              <button
                type="button"
                class="spring relative flex h-10 w-10 items-center justify-center rounded-full text-white/40 backdrop-blur-sm transition-all hover:bg-white/15 hover:text-white"
                :class="{ '!text-[#1db954]': repeatMode !== 'off' }"
                :disabled="!currentTrack"
                @click="toggleRepeat"
              >
                <i class="pi pi-refresh text-sm" />
                <span
                  v-if="repeatMode === 'one'"
                  class="absolute -top-0.5 -right-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-[#1db954] text-[9px] font-bold text-black"
                  >1</span
                >
              </button>
            </div>
          </div>
        </div>

        <!-- Bottom gradient -->
        <div
          class="pointer-events-none absolute right-0 bottom-0 left-0 z-[1] h-48 bg-gradient-to-t from-black/60 to-transparent"
        />
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { usePlayerControls } from '@/composables/player'
import { useAlbumColors } from '@/composables/useAlbumColors'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{
  'update:visible': [value: boolean]
}>()

const rootEl = ref<HTMLElement | null>(null)
const particleCanvas = ref<HTMLCanvasElement | null>(null)

const pc = usePlayerControls()
const currentTrack = pc.currentTrack
const isPlaying = pc.isPlaying
const isBuffering = pc.isBuffering
const isLoadingTrack = pc.isLoadingTrack
const currentTime = pc.currentTime
const duration = pc.duration
const progressPercent = pc.progressPercent
const hasNext = pc.hasNext
const hasPrevious = pc.hasPrevious
const shuffleMode = pc.shuffleMode
const repeatMode = pc.repeatMode

const togglePlayPause = pc.togglePlayPause
const seekPercent = pc.seekPercent
const playNext = pc.playNext
const playPrevious = pc.playPrevious
const toggleShuffle = pc.toggleShuffle
const toggleRepeat = pc.toggleRepeat

const title = computed(() => currentTrack.value?.title || 'No track')
const artistName = computed(() => currentTrack.value?.artistName || '')
const coverUrl = computed(() => currentTrack.value?.coverUrl || '')
const { palette: albumPalette } = useAlbumColors(coverUrl)

const bgStyle = computed(() => {
  if (!coverUrl.value) {
    return {
      background: 'linear-gradient(135deg, #0a0a0a 0%, #121212 100%)',
    }
  }
  return {
    background: `linear-gradient(135deg, ${albumPalette.value.dominant}ee 0%, ${albumPalette.value.dark} 100%)`,
  }
})

function fmtTime(s: number) {
  const total = Math.max(0, Math.floor(Number(s) || 0))
  const m = Math.floor(total / 60)
  const sec = total % 60
  return `${m}:${String(sec).padStart(2, '0')}`
}

const currentTimeLabel = computed(() => fmtTime(currentTime.value))
const durationLabel = computed(() =>
  fmtTime(duration.value || currentTrack.value?.durationSeconds || 0),
)
const progressStyle = computed(() => ({ '--range-progress': `${progressPercent.value}%` }))

function onSeek(e: Event) {
  seekPercent(Number((e.target as HTMLInputElement).value))
}

function close() {
  emit('update:visible', false)
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') close()
  if (e.key === ' ' && e.target === e.currentTarget) {
    e.preventDefault()
    togglePlayPause()
  }
}

// Particle system
interface AmbParticle {
  x: number
  y: number
  vx: number
  vy: number
  size: number
  alpha: number
  alphaSpeed: number
}

let particles: AmbParticle[] = []
let animFrameId: number | null = null
let particleCtx: CanvasRenderingContext2D | null = null

function initParticles() {
  const canvas = particleCanvas.value
  if (!canvas) return
  particleCtx = canvas.getContext('2d')
  if (!particleCtx) return

  resizeParticles()
  window.addEventListener('resize', resizeParticles)

  const count = 50
  particles = Array.from({ length: count }, () => createAmbParticle())
  startParticleLoop()
}

function resizeParticles() {
  const canvas = particleCanvas.value
  if (!canvas) return
  canvas.width = canvas.offsetWidth * window.devicePixelRatio
  canvas.height = canvas.offsetHeight * window.devicePixelRatio
}

function createAmbParticle(): AmbParticle {
  const canvas = particleCanvas.value
  const w = canvas?.offsetWidth || 800
  const h = canvas?.offsetHeight || 600
  return {
    x: Math.random() * w,
    y: Math.random() * h,
    vx: (Math.random() - 0.5) * 0.15,
    vy: -(0.05 + Math.random() * 0.1),
    size: 1 + Math.random() * 2.5,
    alpha: 0.05 + Math.random() * 0.2,
    alphaSpeed: 0.001 + Math.random() * 0.003,
  }
}

function startParticleLoop() {
  function draw() {
    const canvas = particleCanvas.value
    if (!particleCtx || !canvas) return
    const dpr = window.devicePixelRatio
    const w = canvas.offsetWidth
    const h = canvas.offsetHeight

    particleCtx.clearRect(0, 0, w * dpr, h * dpr)
    particleCtx.setTransform(dpr, 0, 0, dpr, 0, 0)

    particles.forEach((p) => {
      p.x += p.vx
      p.y += p.vy
      p.alpha += p.alphaSpeed
      if (p.alpha > 0.3 || p.alpha < 0.02) p.alphaSpeed *= -1

      if (p.y < -10 || p.x < -10 || p.x > w + 10) {
        Object.assign(p, createAmbParticle())
        p.y = h + 5
      }

      particleCtx!.beginPath()
      particleCtx!.arc(p.x, p.y, p.size, 0, Math.PI * 2)
      particleCtx!.fillStyle = `rgba(255, 255, 255, ${p.alpha})`
      particleCtx!.fill()
    })

    animFrameId = window.requestAnimationFrame(draw)
  }
  animFrameId = window.requestAnimationFrame(draw)
}

function stopParticles() {
  if (animFrameId !== null) {
    cancelAnimationFrame(animFrameId)
    animFrameId = null
  }
  window.removeEventListener('resize', resizeParticles)
  particles = []
  particleCtx = null
}

watch(
  () => props.visible,
  async (v) => {
    if (v) {
      await nextTick()
      initParticles()
    } else {
      stopParticles()
    }
  },
)

onMounted(async () => {
  await nextTick()
  rootEl.value?.focus()
})

onBeforeUnmount(() => {
  stopParticles()
})
</script>

<style scoped>
.ambient-enter-enter-active {
  transition:
    opacity 800ms ease,
    transform 800ms ease;
}
.ambient-enter-leave-active {
  transition:
    opacity 400ms ease,
    transform 400ms ease;
}
.ambient-enter-enter-from {
  opacity: 0;
  transform: scale(1.05);
}
.ambient-enter-leave-to {
  opacity: 0;
  transform: scale(1.05);
}
.spring {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
</style>
