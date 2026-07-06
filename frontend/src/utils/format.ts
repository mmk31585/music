export function formatDuration(seconds: number | null | undefined, locale?: string): string {
  if (seconds == null || !isFinite(seconds) || seconds < 0) return '0:00'
  const loc = locale || (typeof navigator !== 'undefined' ? navigator.language : 'en')
  const nf = new Intl.NumberFormat(loc)
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  const fmtSecs = nf.format(secs)
  // Pad with locale-appropriate zero character
  const zero = nf.format(0)
  const padded = fmtSecs.length >= 2 ? fmtSecs : zero.repeat(2 - fmtSecs.length) + fmtSecs
  return `${nf.format(mins)}:${padded}`
}
