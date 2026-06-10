export type PerformanceTier = 'low' | 'medium' | 'high'

export function detectPerformanceTier(): PerformanceTier {
  let score = 0

  if (typeof navigator !== 'undefined') {
    if (navigator.hardwareConcurrency) {
      if (navigator.hardwareConcurrency >= 8) score += 2
      else if (navigator.hardwareConcurrency >= 4) score += 1
    }

    if ((navigator as any).deviceMemory) {
      if ((navigator as any).deviceMemory >= 8) score += 2
      else if ((navigator as any).deviceMemory >= 4) score += 1
    }
  }

  if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
    return 'low'
  }

  if (score >= 3) return 'high'
  if (score >= 1) return 'medium'
  return 'low'
}
