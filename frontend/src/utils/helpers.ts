import type { ClassValue } from 'clsx'
import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function isValidInternalRedirectLink(link: string | undefined | null): boolean {
  if (!link) return false
  return link === '/' || /^\/[A-Za-z0-9_-]+(?:\/[A-Za-z0-9_-]+)*$/.test(link)
}

export const getTextColor = (
  backgroundColor: string,
  darkColor: string = 'black',
  lightColor: string = 'white',
) => {
  function hexToRgb(hex: string) {
    // Convert hex color to RGB
    const bigint = parseInt(hex.slice(1), 16)
    return {
      r: (bigint >> 16) & 255,
      g: (bigint >> 8) & 255,
      b: bigint & 255,
    }
  }

  function calculateBrightness(color: string) {
    // Function to calculate brightness from RGB values
    const rgb = hexToRgb(color)
    return (rgb.r * 299 + rgb.g * 587 + rgb.b * 114) / 1000
  }

  // Convert the background color to a brightness value
  const brightness = calculateBrightness(backgroundColor)

  // Use a threshold value to determine whether to use white or black text
  return brightness > 155 ? darkColor : lightColor
}

/**
 * Estimate the reading time of a text.
 * @param content Plain‑text content.
 * @returns Minutes (rounded up) required to read the content.
 */
export const estimateReadTime = (content: string): number => {
  const words = content.trim().split(/\s+/)
  const wordCount = words.length
  const readingSpeed = 200 // words per minute

  let minutes: number

  if (wordCount <= 15_000) {
    const totalWeight = words.reduce((total, w) => {
      const weight = Math.max(1, Math.ceil(w.length / 5))
      return total + weight
    }, 0)
    minutes = totalWeight / readingSpeed
  } else {
    minutes = wordCount / readingSpeed
  }

  return Math.ceil(minutes)
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function onImgError(e: Event) {
  const img = e.currentTarget as HTMLImageElement | null
  if (!img) return

  // If the failed URL is absolute pointing to localhost, retry with relative path
  if (!img.dataset?.retried && img.src.includes('localhost:8080')) {
    img.dataset.retried = 'true'
    try {
      const url = new URL(img.src)
      img.src = url.pathname + url.search
      return
    } catch {
      // fall through to placeholder
    }
  }

  img.src = 'data:image/svg+xml,' + encodeURIComponent(
    '<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200" fill="%231a1a1a"><rect width="200" height="200"/><text x="50%" y="50%" dominant-anchor="central" text-anchor="middle" fill="%23666" font-size="14">No Image</text></svg>'
  )
}

/**
 * Normalize media URLs returned by the API. If the URL is an absolute URL
 * pointing to the local backend (localhost:8080), convert it to a relative
 * path so it works through the Vite proxy or current origin.
 */
export function normalizeMediaUrl(url: string | null | undefined): string | null {
  if (!url) return null
  try {
    const parsed = new URL(url)
    if (parsed.hostname === 'localhost' && parsed.port === '8080') {
      return parsed.pathname + parsed.search
    }
  } catch {
    // Already a relative URL or invalid — return as-is
  }
  return url
}
