import { ref, type Ref } from 'vue'
import { watchDebounced } from '@vueuse/core'

export interface AlbumColorPalette {
  vibrant: string
  muted: string
  dark: string
  light: string
  gradient: string
  dominant: string
}

function rgbToHex(r: number, g: number, b: number) {
  const toHex = (n: number) =>
    Math.max(0, Math.min(255, Math.round(n)))
      .toString(16)
      .padStart(2, '0')
  return `#${toHex(r)}${toHex(g)}${toHex(b)}`
}

function extractPalette(imageData: ImageData): AlbumColorPalette {
  const pixels = imageData.data
  const colorBuckets = new Map<string, number>()
  const step = 4

  for (let i = 0; i < pixels.length; i += step * 4) {
    const r = pixels[i]!
    const g = pixels[i + 1]!
    const b = pixels[i + 2]!
    const key = `${Math.round(r / 16) * 16},${Math.round(g / 16) * 16},${Math.round(b / 16) * 16}`
    colorBuckets.set(key, (colorBuckets.get(key) || 0) + 1)
  }

  const sorted = [...colorBuckets.entries()]
    .sort((a, b) => b[1] - a[1])
    .slice(0, 8)
    .map(([key]) => {
      const [r, g, b] = key.split(',').map(Number)
      return { r: r!, g: g!, b: b! }
    })

  if (sorted.length === 0) {
    return {
      vibrant: '#1db954',
      muted: '#1a1a2e',
      dark: '#0a0a0a',
      light: '#404060',
      gradient: 'linear-gradient(135deg, #0a0a0a 0%, #1a1a2e 50%, #121212 100%)',
      dominant: '#1a1a2e',
    }
  }

  const dominant = sorted[0]!

  const vibrant =
    sorted.find((c) => {
      const max = Math.max(c.r, c.g, c.b)
      const min = Math.min(c.r, c.g, c.b)
      return max - min > 60 && max > 100
    }) || dominant

  const dark = sorted.reduce((acc, c) => (c.r + c.g + c.b < acc.r + acc.g + acc.b ? c : acc))

  const light = sorted.reduce((acc, c) => (c.r + c.g + c.b > acc.r + acc.g + acc.b ? c : acc))

  const muted = sorted.length > 2 ? sorted[2]! : dominant

  const gradient = `linear-gradient(135deg, ${rgbToHex(dark.r, dark.g, dark.b)} 0%, ${rgbToHex(dominant.r, dominant.g, dominant.b)} 50%, ${rgbToHex(muted.r, muted.g, muted.b)} 100%)`

  return {
    vibrant: rgbToHex(vibrant.r, vibrant.g, vibrant.b),
    muted: rgbToHex(muted.r, muted.g, muted.b),
    dark: rgbToHex(dark.r, dark.g, dark.b),
    light: rgbToHex(light.r, light.g, light.b),
    gradient,
    dominant: rgbToHex(dominant.r, dominant.g, dominant.b),
  }
}

/**
 * Module-level palette cache shared across all components.
 * Keyed by cover URL so different components with the same URL share cached results.
 */
const paletteCache = new Map<string, AlbumColorPalette>()

const defaultPalette: AlbumColorPalette = {
  vibrant: '#1db954',
  muted: '#1a1a2e',
  dark: '#0a0a0a',
  light: '#404060',
  gradient: 'linear-gradient(135deg, #0a0a0a 0%, #1a1a2e 50%, #121212 100%)',
  dominant: '#1a1a2e',
}

function createExtractFunction(palette: Ref<AlbumColorPalette>, loading: Ref<boolean>) {
  let extractGeneration = 0

  return async function extract(url: string) {
    if (!url) return
    const gen = ++extractGeneration
    const cached = paletteCache.get(url)
    if (cached) {
      palette.value = cached
      return
    }

    loading.value = true

    try {
      const img = new Image()
      img.crossOrigin = 'anonymous'

      await new Promise<void>((resolve, reject) => {
        img.onload = () => resolve()
        img.onerror = () => reject(new Error('Failed to load image'))
        img.src = url
      })

      if (gen !== extractGeneration) return

      const canvas = document.createElement('canvas')
      const size = 64
      canvas.width = size
      canvas.height = size
      const ctx = canvas.getContext('2d')
      if (!ctx) throw new Error('No canvas context')

      ctx.drawImage(img, 0, 0, size, size)
      const imageData = ctx.getImageData(0, 0, size, size)

      const p = extractPalette(imageData)
      paletteCache.set(url, p)
      palette.value = p
    } catch {
      palette.value = { ...defaultPalette }
    } finally {
      loading.value = false
    }
  }
}

export function useAlbumColors(coverUrl: Ref<string | null | undefined>) {
  // Each component gets its own palette + loading refs
  const palette = ref<AlbumColorPalette>({ ...defaultPalette })
  const loading = ref(false)
  const extract = createExtractFunction(palette, loading)

  watchDebounced(coverUrl, (url) => {
    if (url) void extract(url)
  }, { debounce: 300, maxWait: 1000 })

  return {
    palette,
    loading,
    extract,
  }
}
