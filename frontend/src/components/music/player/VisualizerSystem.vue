<template>
  <div
    ref="container"
    class="visualizer-container relative h-full w-full overflow-hidden"
  >
    <canvas ref="canvas" class="h-full w-full" />
    <div
      v-if="!hasAudioData && isPlaying"
      class="absolute inset-0 flex items-center justify-center"
    >
      <div class="flex h-12 items-end gap-[3px]">
        <div
          v-for="i in 5"
          :key="i"
          class="equalizer-bar w-[3px] rounded-full bg-white/20"
          :class="`equalizer-bar-${i}`"
          :style="{ height: `${12 + i * 6}px` }"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { audioEngine } from '@/services/player/audio-engine'
import { detectPerformanceTier } from '@/utils/performance'

export type VisualizerMode = 'spectrum' | 'waveform' | 'circular' | 'radialBars' | 'particle' | 'fluid'

const props = withDefaults(
  defineProps<{
    mode?: VisualizerMode
    colorPalette?: {
      vibrant: string
      muted: string
      dark: string
      light: string
      dominant: string
      gradient: string
    }
    isPlaying?: boolean
    barCount?: number
    sensitivity?: number
    glowIntensity?: number
    mirrored?: boolean
    roundedBars?: boolean
    showCenterArt?: boolean
    centerImageUrl?: string
  }>(),
  {
    mode: 'spectrum',
    isPlaying: false,
    barCount: 64,
    sensitivity: 1,
    glowIntensity: 0.5,
    mirrored: true,
    roundedBars: true,
    showCenterArt: false,
  },
)

const container = ref<HTMLElement | null>(null)
const canvas = ref<HTMLCanvasElement | null>(null)

let ctx: CanvasRenderingContext2D | null = null
let analyserNode: AnalyserNode | null = null
let animFrameId: number | null = null
let frequencyData: Uint8Array | null = null
let waveformData: Uint8Array | null = null
let resizeHandler: (() => void) | null = null

const hasAudioData = ref(false)

const perfTier = detectPerformanceTier()
const isLowTier = perfTier === 'low'
const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches

const effectiveMode = computed<VisualizerMode>(() => {
  if (isLowTier && props.mode !== 'spectrum') return 'spectrum'
  return props.mode
})

const colors = computed(() => {
  const p = props.colorPalette
  if (!p) {
    return {
      primary: 'rgba(29, 185, 84, 0.9)',
      secondary: 'rgba(96, 165, 250, 0.7)',
      tertiary: 'rgba(168, 85, 247, 0.5)',
      quaternary: 'rgba(244, 114, 182, 0.4)',
      background: 'rgba(0, 0, 0, 0.3)',
    }
  }
  return {
    primary: p.vibrant,
    secondary: p.light,
    tertiary: p.muted,
    quaternary: p.dominant,
    background: p.dark,
  }
})

function tryConnectAudio() {
  const shared = audioEngine.getAnalyserNode()
  if (!shared) return false

  analyserNode = shared
  analyserNode.fftSize = 256
  analyserNode.smoothingTimeConstant = 0.8

  frequencyData = new Uint8Array(analyserNode.frequencyBinCount) as Uint8Array
  waveformData = new Uint8Array(analyserNode.frequencyBinCount) as Uint8Array
  hasAudioData.value = true
  return true
}

interface SimulatedBar {
  height: number
  target: number
  speed: number
}

let simulatedBars: SimulatedBar[] = []

function initSimulated() {
  simulatedBars = Array.from({ length: props.barCount }, () => ({
    height: 0.08 + Math.random() * 0.15,
    target: 0.08 + Math.random() * 0.15,
    speed: 0.02 + Math.random() * 0.04,
  }))
}

interface Particle {
  x: number
  y: number
  vx: number
  vy: number
  size: number
  alpha: number
  baseAlpha: number
  pulsePhase: number
  hue: number
}

let particles: Particle[] = []

function initParticles() {
  const w = container.value?.offsetWidth || 800
  const h = container.value?.offsetHeight || 600
  const count = isLowTier ? 30 : 100
  particles = Array.from({ length: count }, () => ({
    x: Math.random() * w,
    y: Math.random() * h,
    vx: (Math.random() - 0.5) * 0.8,
    vy: (Math.random() - 0.5) * 0.4 - 0.3,
    size: 1.5 + Math.random() * 4,
    alpha: 0.1 + Math.random() * 0.4,
    baseAlpha: 0.1 + Math.random() * 0.4,
    pulsePhase: Math.random() * Math.PI * 2,
    hue: Math.random() * 60 + 120,
  }))
}

interface FluidCell {
  density: number
  vx: number
  vy: number
}

let fluidGrid: FluidCell[][] = []
let fluidTime = 0

function initFluid(cols: number, rows: number) {
  fluidGrid = Array.from({ length: cols }, () =>
    Array.from({ length: rows }, () => ({ density: 0, vx: 0, vy: 0 })),
  )
  fluidTime = 0
}

let resizeRaf: number | null = null

function resize() {
  const el = container.value
  const cvs = canvas.value
  if (!el || !cvs) return
  const rect = el.getBoundingClientRect()
  const dpr = Math.min(window.devicePixelRatio, 2)
  cvs.width = rect.width * dpr
  cvs.height = rect.height * dpr
  cvs.style.width = '100%'
  cvs.style.height = '100%'
}

function onResizeThrottled() {
  if (resizeRaf !== null) cancelAnimationFrame(resizeRaf)
  resizeRaf = requestAnimationFrame(resize)
}

let cachedAvgAmplitude = 0

function getFreqData(): Uint8Array | null {
  if (!analyserNode || !frequencyData) return null
  analyserNode.getByteFrequencyData(frequencyData as any as Uint8Array)
  let sum = 0
  for (let i = 0; i < frequencyData.length; i++) {
    sum += frequencyData[i]!
  }
  cachedAvgAmplitude = sum / (frequencyData.length * 255)
  return frequencyData
}

function getWaveData(): Uint8Array | null {
  if (!analyserNode || !waveformData) return null
  analyserNode.getByteTimeDomainData(waveformData as any as Uint8Array)
  return waveformData
}

function lerp(a: number, b: number, t: number) {
  return a + (b - a) * t
}

function clamp(v: number, min = 0, max = 1) {
  return Math.max(min, Math.min(max, v))
}

function hexToRgba(hex: string, alpha = 1): string {
  if (hex.startsWith('rgba')) return hex.replace(/[\d.]+\)$/, `${alpha})`)
  if (hex.startsWith('rgb')) return hex.replace('rgb', 'rgba').replace(')', `, ${alpha})`)
  const h = hex.replace('#', '')
  const r = parseInt(h.substring(0, 2), 16)
  const g = parseInt(h.substring(2, 4), 16)
  const b = parseInt(h.substring(4, 6), 16)
  if (isNaN(r)) return `rgba(255, 255, 255, ${alpha})`
  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}

let smoothFreqs: Float32Array | null = null

function drawSpectrum() {
  if (!ctx || !canvas.value) return
  const dpr = window.devicePixelRatio
  const w = canvas.value.width / dpr
  const h = canvas.value.height / dpr
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, w, h)

  const bars = props.barCount
  const freqData = getFreqData()
  const useAudio = freqData && props.isPlaying

  if (!smoothFreqs || smoothFreqs.length !== bars) {
    smoothFreqs = new Float32Array(bars)
  }

  const values = new Float32Array(bars)
  for (let i = 0; i < bars; i++) {
    if (useAudio) {
      const freqIdx = Math.floor((i / bars) * freqData!.length)
      values[i] = freqData![freqIdx]! / 255 * props.sensitivity
    } else {
      const bar = simulatedBars[i]
      if (!bar) continue
      if (props.isPlaying) {
        if (Math.random() < 0.03) {
          bar.target = 0.1 + Math.random() * 0.9
        }
        bar.height += (bar.target - bar.height) * bar.speed
      } else {
        bar.height += (0.05 - bar.height) * 0.03
      }
      values[i] = bar.height
    }
    smoothFreqs[i] = lerp(smoothFreqs[i] ?? 0, values[i] ?? 0, 0.35)
  }

  const gap = Math.max(1, w * 0.003)
  const totalGap = gap * (bars - 1)
  const barW = (w - totalGap) / bars
  const maxH = h * 0.85

  for (let i = 0; i < bars; i++) {
    const val = clamp(smoothFreqs[i] || 0)
    const bh = Math.max(1.5, val * maxH)
    const x = i * (barW + gap)

    let y1: number, h1: number
    if (props.mirrored) {
      y1 = (h - bh) / 2
      h1 = bh
    } else {
      y1 = h - bh
      h1 = bh
    }

    const grad = ctx.createLinearGradient(x, y1, x, y1 + h1)
    const alpha = clamp(val)
    if (val > 0.65) {
      grad.addColorStop(0, hexToRgba(colors.value.primary, alpha))
      grad.addColorStop(0.5, hexToRgba(colors.value.secondary, alpha * 0.8))
      grad.addColorStop(1, hexToRgba(colors.value.tertiary, alpha * 0.2))
    } else if (val > 0.35) {
      grad.addColorStop(0, hexToRgba(colors.value.secondary, alpha))
      grad.addColorStop(1, hexToRgba(colors.value.tertiary, alpha * 0.3))
    } else {
      grad.addColorStop(0, hexToRgba(colors.value.tertiary, alpha * 0.8))
      grad.addColorStop(1, hexToRgba(colors.value.quaternary, alpha * 0.2))
    }

    ctx.beginPath()
    if (props.roundedBars) {
      ctx.roundRect(x, y1, barW, h1, [barW / 2, barW / 2, barW / 2, barW / 2])
    } else {
      ctx.roundRect(x, y1, barW, h1, [1, 1, 1, 1])
    }
    ctx.fillStyle = grad
    ctx.fill()

    if (props.isPlaying && val > 0.5 && props.glowIntensity > 0.1) {
      ctx.shadowColor = hexToRgba(colors.value.primary, 0.5 * props.glowIntensity)
      ctx.shadowBlur = 6 * props.glowIntensity
      ctx.fill()
      ctx.shadowBlur = 0
    }
  }
}

let smoothWaveform: Float32Array | null = null

function drawWaveform() {
  if (!ctx || !canvas.value) return
  const dpr = window.devicePixelRatio
  const w = canvas.value.width / dpr
  const h = canvas.value.height / dpr
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, w, h)

  const waveData = getWaveData()
  const useAudio = waveData && props.isPlaying

  const points = Math.max(30, Math.floor(w / 3))
  if (!smoothWaveform || smoothWaveform.length !== points) {
    smoothWaveform = new Float32Array(points)
  }

  const values = new Float32Array(points)
  for (let i = 0; i < points; i++) {
    if (useAudio) {
      const idx = Math.floor((i / points) * waveData!.length)
      values[i] = (waveData![idx]! - 128) / 128 * props.sensitivity
    } else {
      const bar = simulatedBars[i % simulatedBars.length]
      if (!bar) continue
      values[i] = (bar.height - 0.5) * 2
    }
    smoothWaveform[i] = lerp(smoothWaveform[i] ?? 0, values[i] ?? 0, 0.25)
  }

  const midY = h / 2
  const amp = h * 0.3

  ctx.beginPath()
  ctx.moveTo(0, midY)
  for (let i = 0; i < points; i++) {
    const x = (i / points) * w
    const y = midY + smoothWaveform[i]! * amp
    ctx.lineTo(x, y)
  }
  ctx.lineTo(w, midY)

  const grad = ctx.createLinearGradient(0, 0, 0, h)
  grad.addColorStop(0, hexToRgba(colors.value.primary, 0.6))
  grad.addColorStop(0.3, hexToRgba(colors.value.secondary, 0.8))
  grad.addColorStop(0.7, hexToRgba(colors.value.tertiary, 0.6))
  grad.addColorStop(1, hexToRgba(colors.value.quaternary, 0.3))

  ctx.strokeStyle = grad
  ctx.lineWidth = 2.5
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'
  ctx.stroke()

  if (props.glowIntensity > 0.1) {
    ctx.shadowColor = hexToRgba(colors.value.primary, 0.3 * props.glowIntensity)
    ctx.shadowBlur = 8 * props.glowIntensity
    ctx.stroke()
    ctx.shadowBlur = 0
  }

  ctx.beginPath()
  ctx.moveTo(0, midY)
  for (let i = 0; i < points; i++) {
    const x = (i / points) * w
    const y = midY + smoothWaveform[i]! * amp * 0.5
    ctx.lineTo(x, y)
  }
  ctx.lineTo(w, midY)
  ctx.strokeStyle = hexToRgba(colors.value.primary, 0.2)
  ctx.lineWidth = 1
  ctx.stroke()
}

function drawCircular() {
  if (!ctx || !canvas.value) return
  const dpr = window.devicePixelRatio
  const w = canvas.value.width / dpr
  const h = canvas.value.height / dpr
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, w, h)

  const cx = w / 2
  const cy = h / 2
  const radius = Math.min(w, h) * 0.32
  const bars = Math.min(props.barCount, 64)
  const freqData = getFreqData()
  const useAudio = freqData && props.isPlaying

  for (let i = 0; i < bars; i++) {
    let val: number
    if (useAudio) {
      const freqIdx = Math.floor((i / bars) * freqData!.length)
      val = freqData![freqIdx]! / 255 * props.sensitivity
    } else {
      const bar = simulatedBars[i]
      if (!bar) continue
      if (props.isPlaying) {
        if (Math.random() < 0.03) {
          bar.target = 0.1 + Math.random() * 0.9
        }
        bar.height += (bar.target - bar.height) * bar.speed
      } else {
        bar.height += (0.05 - bar.height) * 0.03
      }
      val = bar.height
    }

    val = clamp(val)
    const angle = (i / bars) * Math.PI * 2 - Math.PI / 2
    const barLen = Math.max(2, val * radius * 2)
    const innerR = radius * 0.55 + val * radius * 0.15
    const x1 = cx + Math.cos(angle) * innerR
    const y1 = cy + Math.sin(angle) * innerR
    const x2 = cx + Math.cos(angle) * (innerR + barLen)
    const y2 = cy + Math.sin(angle) * (innerR + barLen)

    ctx.beginPath()
    ctx.moveTo(x1, y1)
    ctx.lineTo(x2, y2)
    ctx.lineWidth = Math.max(2, barLen * 0.06)
    ctx.lineCap = 'round'

    const c = val > 0.6
      ? hexToRgba(colors.value.primary, val)
      : val > 0.35
        ? hexToRgba(colors.value.secondary, val * 0.8)
        : hexToRgba(colors.value.tertiary, val * 0.6)
    ctx.strokeStyle = c
    ctx.stroke()

    if (props.isPlaying && val > 0.5 && props.glowIntensity > 0.1) {
      ctx.shadowColor = hexToRgba(colors.value.primary, 0.4 * props.glowIntensity)
      ctx.shadowBlur = 5 * props.glowIntensity
      ctx.stroke()
      ctx.shadowBlur = 0
    }
  }

  ctx.beginPath()
  ctx.arc(cx, cy, radius * 0.35, 0, Math.PI * 2)
  ctx.fillStyle = 'rgba(255, 255, 255, 0.03)'
  ctx.fill()

  if (props.showCenterArt && props.centerImageUrl) {
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.src = props.centerImageUrl
    ctx.save()
    ctx.beginPath()
    ctx.arc(cx, cy, radius * 0.32, 0, Math.PI * 2)
    ctx.clip()
    ctx.drawImage(img, cx - radius * 0.32, cy - radius * 0.32, radius * 0.64, radius * 0.64)
    ctx.restore()
  }
}

let radialSmooth: Float32Array | null = null

function drawRadialBars() {
  if (!ctx || !canvas.value) return
  const dpr = window.devicePixelRatio
  const w = canvas.value.width / dpr
  const h = canvas.value.height / dpr
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, w, h)

  const cx = w / 2
  const cy = h / 2
  const maxRadius = Math.min(w, h) * 0.42
  const bars = Math.min(props.barCount, 72)
  const freqData = getFreqData()
  const useAudio = freqData && props.isPlaying

  if (!radialSmooth || radialSmooth.length !== bars) {
    radialSmooth = new Float32Array(bars)
  }

  for (let i = 0; i < bars; i++) {
    let val: number
    if (useAudio) {
      const freqIdx = Math.floor((i / bars) * freqData!.length)
      val = freqData![freqIdx]! / 255 * props.sensitivity
    } else {
      const bar = simulatedBars[i]
      if (!bar) continue
      if (props.isPlaying) {
        if (Math.random() < 0.03) {
          bar.target = 0.1 + Math.random() * 0.9
        }
        bar.height += (bar.target - bar.height) * bar.speed
      } else {
        bar.height += (0.05 - bar.height) * 0.03
      }
      val = bar.height
    }
    radialSmooth[i] = lerp(radialSmooth[i]!, clamp(val), 0.3)
  }

  for (let i = 0; i < bars; i++) {
    const val = radialSmooth[i] || 0
    const angle = (i / bars) * Math.PI * 2 - Math.PI / 2
    const barLen = val * maxRadius * 0.9
    const barW = (2 * Math.PI * maxRadius) / bars * 0.8
    const x1 = cx + Math.cos(angle) * (maxRadius - barLen)
    const y1 = cy + Math.sin(angle) * (maxRadius - barLen)
    const x2 = cx + Math.cos(angle) * maxRadius
    const y2 = cy + Math.sin(angle) * maxRadius

    ctx.beginPath()
    ctx.moveTo(x1, y1)
    ctx.lineTo(x2, y2)
    ctx.lineWidth = Math.max(1.5, barW)
    ctx.lineCap = 'round'

    const grad = ctx.createLinearGradient(x1, y1, x2, y2)
    grad.addColorStop(0, hexToRgba(colors.value.tertiary, val * 0.3))
    grad.addColorStop(1, hexToRgba(colors.value.primary, val * 0.9))
    ctx.strokeStyle = grad
    ctx.stroke()

    if (val > 0.4 && props.glowIntensity > 0.1) {
      ctx.shadowColor = hexToRgba(colors.value.primary, 0.3 * props.glowIntensity)
      ctx.shadowBlur = 4 * props.glowIntensity
      ctx.stroke()
      ctx.shadowBlur = 0
    }
  }

  const auraRadius = maxRadius * 0.98
  const auraGrad = ctx.createRadialGradient(cx, cy, 0, cx, cy, auraRadius)
  auraGrad.addColorStop(0, 'rgba(255, 255, 255, 0.02)')
  auraGrad.addColorStop(0.5, 'rgba(255, 255, 255, 0.01)')
  auraGrad.addColorStop(1, 'rgba(255, 255, 255, 0)')
  ctx.beginPath()
  ctx.arc(cx, cy, auraRadius, 0, Math.PI * 2)
  ctx.fillStyle = auraGrad
  ctx.fill()

  if (props.showCenterArt && props.centerImageUrl) {
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.src = props.centerImageUrl
    ctx.save()
    ctx.beginPath()
    ctx.arc(cx, cy, maxRadius * 0.15, 0, Math.PI * 2)
    ctx.clip()
    ctx.drawImage(img, cx - maxRadius * 0.15, cy - maxRadius * 0.15, maxRadius * 0.3, maxRadius * 0.3)
    ctx.restore()
  }
}

function drawParticles(_time: number, dt: number) {
  if (!ctx || !canvas.value) return
  const dpr = window.devicePixelRatio
  const w = canvas.value.width / dpr
  const h = canvas.value.height / dpr
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, w, h)

  getFreqData()
  const avgFreq = cachedAvgAmplitude
  const energy = props.isPlaying ? avgFreq || 0.3 : 0.1
  const dtFactor = Math.min(dt, 2)

  particles.forEach((p) => {
    p.x += p.vx * (1 + energy) * dtFactor
    p.y += p.vy * (1 + energy * 0.5) * dtFactor

    if (props.isPlaying) {
      p.pulsePhase += 0.04 * (1 + energy) * dtFactor
      p.alpha = p.baseAlpha + Math.sin(p.pulsePhase) * 0.2 + energy * 0.25
    } else {
      p.alpha += (0.03 - p.alpha) * 0.02 * dtFactor
    }

    if (p.x < -20 || p.x > w + 20 || p.y < -20 || p.y > h + 20) {
      p.x = Math.random() * w
      p.y = Math.random() * h
      p.vx = (Math.random() - 0.5) * 0.8
      p.vy = (Math.random() - 0.5) * 0.4 - 0.2
    }

    const size = p.size * (0.8 + energy * 0.8)
    ctx!.beginPath()
    ctx!.arc(p.x, p.y, size, 0, Math.PI * 2)

    const alpha = clamp(p.alpha)
    ctx!.fillStyle = hexToRgba(colors.value.primary, alpha * 0.6)
    ctx!.fill()

    ctx!.beginPath()
    ctx!.arc(p.x - size * 0.2, p.y - size * 0.2, size * 0.3, 0, Math.PI * 2)
    ctx!.fillStyle = hexToRgba(colors.value.secondary || '#ffffff', alpha * 0.3)
    ctx!.fill()

    if (props.isPlaying && energy > 0.3 && props.glowIntensity > 0.1) {
      ctx!.shadowColor = hexToRgba(colors.value.primary, 0.3 * props.glowIntensity)
      ctx!.shadowBlur = 6 * props.glowIntensity
      ctx!.beginPath()
      ctx!.arc(p.x, p.y, size * 2, 0, Math.PI * 2)
      ctx!.fill()
      ctx!.shadowBlur = 0
    }
  })

  if (particles.length < 100 && Math.random() < 0.15) {
    particles.push({
      x: Math.random() * w,
      y: h + 5,
      vx: (Math.random() - 0.5) * 0.8,
      vy: -(0.2 + Math.random() * 0.4),
      size: 1.5 + Math.random() * 4,
      alpha: 0.1 + Math.random() * 0.4,
      baseAlpha: 0.1 + Math.random() * 0.4,
      pulsePhase: Math.random() * Math.PI * 2,
      hue: Math.random() * 60 + 120,
    })
  }
}

function drawFluid(_time: number, dt: number) {
  if (!ctx || !canvas.value) return
  if (frameCount % 2 !== 0) return
  const dpr = window.devicePixelRatio
  const w = canvas.value.width / dpr
  const h = canvas.value.height / dpr
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  ctx.clearRect(0, 0, w, h)

  const cols = isLowTier ? 20 : 40
  const rows = Math.floor((h / w) * cols)
  const cellW = w / cols
  const cellH = h / rows

  if (fluidGrid.length !== cols || fluidGrid[0]?.length !== rows) {
    initFluid(cols, rows)
  }

  const freqData = getFreqData()
  const energy = props.isPlaying ? cachedAvgAmplitude || 0.3 : 0.05
  const dtFactor = Math.min(dt, 2)
  fluidTime += dt * 0.5 * dtFactor

  for (let x = 0; x < cols; x++) {
    for (let y = 0; y < rows; y++) {
      const freqIdx = Math.floor(((x + y * cols) / (cols * rows)) * (freqData?.length || 1))
      const freqVal = freqData?.[freqIdx] ?? 128
      const normalizedFreq = clamp(freqVal / 255)

      const distX = x / cols - 0.5
      const distY = y / rows - 0.5
      const dist = Math.sqrt(distX * distX + distY * distY)

      const cell = fluidGrid[x]![y]!
      cell.vx += (Math.sin(fluidTime + x * 0.5 + y * 0.3) * 0.02 + distX * 0.001) * (1 + energy)
      cell.vy += (Math.cos(fluidTime + y * 0.4 + x * 0.6) * 0.02 + distY * 0.001) * (1 + energy)
      cell.vx *= 0.97
      cell.vy *= 0.97

      const newDensity = normalizedFreq * 0.8 * (1 - dist * 0.5)
      cell.density += (newDensity - cell.density) * 0.05 * dtFactor

      const alpha = Math.max(0.02, cell.density * 0.5)
      const cellColor = normalizedFreq > 0.5
        ? hexToRgba(colors.value.primary, alpha)
        : hexToRgba(colors.value.secondary, alpha * 0.7)
      ctx.fillStyle = cellColor
      ctx.fillRect(x * cellW, y * cellH, cellW + 1, cellH + 1)
    }
  }
}

const renderers: Record<VisualizerMode, (time: number, dt: number) => void> = {
  spectrum: drawSpectrum,
  waveform: drawWaveform,
  circular: drawCircular,
  radialBars: drawRadialBars,
  particle: drawParticles,
  fluid: drawFluid,
}

let lastTime = 0
let frameCount = 0

function loop(time: number) {
  if (document.hidden) {
    animFrameId = requestAnimationFrame(loop)
    return
  }

  const dt = Math.min((time - lastTime) / 16, 3)
  lastTime = time
  frameCount++

  const render = renderers[isLowTier || prefersReducedMotion ? 'spectrum' : effectiveMode.value]
  render?.(time, dt)

  if (props.isPlaying) {
    animFrameId = requestAnimationFrame(loop)
  } else {
    animFrameId = null
  }
}

function start() {
  stop()
  lastTime = 0
  animFrameId = requestAnimationFrame(loop)
}

function stop() {
  if (animFrameId !== null) {
    cancelAnimationFrame(animFrameId)
    animFrameId = null
  }
  if (resizeRaf !== null) {
    cancelAnimationFrame(resizeRaf)
    resizeRaf = null
  }
}

function onVisibilityChange() {
  if (document.hidden) {
    if (animFrameId !== null) {
      cancelAnimationFrame(animFrameId)
      animFrameId = null
    }
  } else if (props.isPlaying && animFrameId === null) {
    lastTime = 0
    animFrameId = requestAnimationFrame(loop)
  }
}

function init() {
  const cvs = canvas.value
  if (!cvs) return
  ctx = cvs.getContext('2d')
  if (!ctx) return

  resize()
  window.addEventListener('resize', onResizeThrottled)
  resizeHandler = () => resize()
  document.addEventListener('visibilitychange', onVisibilityChange)

  const connected = tryConnectAudio()
  if (!connected) {
    initSimulated()
  }
  if (!isLowTier) {
    initParticles()
  }
  initFluid(
    40,
    Math.floor(
      ((container.value?.offsetHeight || 600) / (container.value?.offsetWidth || 800)) * 40,
    ),
  )
  start()
}

function cleanup() {
  stop()
  document.removeEventListener('visibilitychange', onVisibilityChange)
  if (resizeHandler) {
    window.removeEventListener('resize', resizeHandler)
    resizeHandler = null
  }
  analyserNode = null
  frequencyData = null
  waveformData = null
  smoothFreqs = null
  smoothWaveform = null
  radialSmooth = null
  ctx = null
  hasAudioData.value = false
}

watch(
  () => props.mode,
  () => {
    resize()
  },
)

watch(
  () => props.barCount,
  () => {
    initSimulated()
  },
)

watch(
  () => props.isPlaying,
  (playing) => {
    if (playing) {
      lastTime = 0
      animFrameId = requestAnimationFrame(loop)
    } else {
      stop()
    }
  },
)

onMounted(() => {
  init()
})

onBeforeUnmount(() => {
  cleanup()
})
</script>

<style scoped>
.visualizer-container canvas {
  display: block;
}
</style>
