export type UnknownRecord = Record<string, any>

export function isPlainRecord(v: any): v is UnknownRecord {
  return typeof v === 'object' && v !== null && !Array.isArray(v)
}

function toFiniteNumber(v: any): number | null {
  if (typeof v === 'number' && Number.isFinite(v)) return v

  if (typeof v === 'string') {
    const t = v.trim()
    if (!t) return null
    const n = Number(t)
    if (Number.isFinite(n)) return n
  }

  return null
}

export function pickString(d: UnknownRecord, keys: readonly string[]): string | null {
  for (const k of keys) {
    const v = d[k]
    if (typeof v === 'string') {
      const t = v.trim()
      if (t) return t
    }
  }
  return null
}

export function pickNumber(d: UnknownRecord, keys: readonly string[]): number | null {
  for (const k of keys) {
    const n = toFiniteNumber(d[k])
    if (n != null) return n
  }
  return null
}

export function pickFirstKey(d: UnknownRecord, keys: readonly string[], fallback: string): string {
  for (const k of keys) if (k in d) return k
  return fallback
}

export function getAt(obj: any, path: readonly string[]): any {
  let cur: any = obj
  for (const key of path) {
    if (!isPlainRecord(cur)) return undefined
    cur = cur[key]
  }
  return cur
}
