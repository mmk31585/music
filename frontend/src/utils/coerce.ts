const DIGITS_ONLY = /^\d+$/

export function toText(v: any): string {
  if (v == null) return ''
  if (typeof v === 'string') return v
  if (typeof v === 'number' || typeof v === 'boolean' || typeof v === 'bigint') return String(v)
  try {
    return JSON.stringify(v)
  } catch {
    return String(v)
  }
}

export function parsePositiveInt(raw: any): number | null {
  const s = toText(raw).trim()
  if (!s) return null
  if (!DIGITS_ONLY.test(s)) return null
  const n = Number(s)
  if (!Number.isFinite(n) || !Number.isInteger(n) || n < 1) return null
  return n
}

export function toId(raw: any): number | null {
  const s = toText(raw).trim()
  if (!s) return null
  const n = Number(s)
  return Number.isFinite(n) ? n : null
}
