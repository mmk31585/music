export function formatCount(count?: number | null, locale?: string): string {
  const loc = locale || (typeof navigator !== 'undefined' ? navigator.language : 'en')
  if (!count && count !== 0) return new Intl.NumberFormat(loc).format(0)
  if (count >= 1_000_000) return `${new Intl.NumberFormat(loc).format(count / 1_000_000)}M`
  if (count >= 1_000) return `${new Intl.NumberFormat(loc).format(count / 1_000)}K`
  return new Intl.NumberFormat(loc).format(count)
}

export const numberUtil = {
  clamp: (n: number, min: number, max: number): number => {
    if (!Number.isFinite(n) || !Number.isFinite(min) || !Number.isFinite(max)) {
      throw new TypeError('clampStrict expects finite numbers')
    }
    if (min > max) throw new RangeError('min must be <= max')

    return Math.min(max, Math.max(min, n))
  },

  normDeg: (d: number) => ((d % 360) + 360) % 360,

  formatBytes: (n?: number) => {
    if (typeof n !== 'number' || !Number.isFinite(n) || n < 0) return '—'

    const units = ['B', 'KB', 'MB', 'GB', 'TB']
    let v = n
    let i = 0
    while (v >= 1024 && i < units.length - 1) {
      v /= 1024
      i++
    }

    return `${v.toFixed(i === 0 ? 0 : 1)} ${units[i]}`
  },
}
