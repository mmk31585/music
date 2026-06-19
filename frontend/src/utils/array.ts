// utils/arrUtil.ts
export type AnyObject = Record<string, any>

const FORBIDDEN_KEYS = new Set(['__proto__', 'prototype', 'constructor'])

const isSafeKey = (k: string) => k.length > 0 && !FORBIDDEN_KEYS.has(k)

const isPlainObject = (v: any): v is AnyObject =>
  typeof v === 'object' && v !== null && !Array.isArray(v)

const parsePath = (path: string): string[] =>
  path
    .split('.')
    .map((s) => s.trim())
    .filter(Boolean)

export const nestedArray = {
  /**
   * Retrieve a deep value from an object using a dot-separated path.
   * Supports array indexes in the path like: "users.0.name"
   */
  get<T = unknown>(obj: AnyObject, path: string, fallback?: T): T | undefined {
    const segments = parsePath(path)
    let cur: any = obj

    for (const seg of segments) {
      if (!isSafeKey(seg)) return fallback
      if (cur == null) return fallback

      if (Array.isArray(cur)) {
        const idx = Number(seg)
        if (!Number.isInteger(idx)) return fallback
        cur = cur[idx]
        continue
      }

      if (!isPlainObject(cur)) return fallback
      cur = (cur as AnyObject)[seg]
    }

    return cur === undefined ? fallback : (cur as T | undefined)
  },

  /**
   * Set a deep value on an object, creating intermediate objects/arrays as needed.
   * Creates arrays when the next segment looks like an integer index.
   */
  set<T = unknown>(obj: AnyObject, path: string, value: T): AnyObject {
    const segments = parsePath(path)
    if (segments.length === 0) return obj

    let cur: any = obj

    for (let i = 0; i < segments.length; i++) {
      const seg = segments[i] as string
      if (!isSafeKey(seg)) return obj

      const isLast = i === segments.length - 1
      const next = segments[i + 1]
      const nextIsIndex = next !== undefined && Number.isInteger(Number(next))

      // If current container is array
      if (Array.isArray(cur)) {
        const idx = Number(seg)
        if (!Number.isInteger(idx) || idx < 0) return obj

        if (isLast) {
          cur[idx] = value as any
          return obj
        }

        const existing = cur[idx]
        if (isPlainObject(existing) || Array.isArray(existing)) {
          cur = existing
        } else {
          cur[idx] = nextIsIndex ? [] : {}
          cur = cur[idx]
        }
        continue
      }

      // If current container is plain object
      if (!isPlainObject(cur)) return obj

      if (isLast) {
        ;(cur as AnyObject)[seg] = value as any
        return obj
      }

      const existing = (cur as AnyObject)[seg]
      if (isPlainObject(existing) || Array.isArray(existing)) {
        cur = existing
      } else {
        ;(cur as AnyObject)[seg] = nextIsIndex ? [] : {}
        cur = (cur as AnyObject)[seg]
      }
    }

    return obj
  },

  /**
   * Remove a deep property from an object/array.
   * - If parent is array and last is a number: sets that index to undefined (does not reindex).
   *   If you want reindexing, do `arr.splice(index, 1)` yourself.
   */
  remove(obj: AnyObject, path: string): AnyObject {
    const segments = parsePath(path)
    if (segments.length === 0) return obj

    const last = segments.pop()!
    if (!isSafeKey(last)) return obj

    const parentPath = segments.join('.')
    const parent = parentPath ? this.get<any>(obj, parentPath) : obj
    if (parent == null) return obj

    if (Array.isArray(parent)) {
      const idx = Number(last)
      if (!Number.isInteger(idx) || idx < 0) return obj
      parent[idx] = undefined
      return obj
    }

    if (isPlainObject(parent)) {
      if (Object.prototype.hasOwnProperty.call(parent, last)) {
        delete (parent as AnyObject)[last]
      }
    }

    return obj
  },
} as const

export const arrUtil = {
  /**
   * Like Laravel's Arr::get, but works with plain objects/arrays via dot paths.
   */
  get<T = unknown>(obj: AnyObject, path: string, fallback?: T): T | undefined {
    return nestedArray.get(obj, path, fallback)
  },

  /**
   * Like Laravel's Arr::set (dot path).
   */
  set<T = unknown>(obj: AnyObject, path: string, value: T): AnyObject {
    return nestedArray.set(obj, path, value)
  },

  /**
   * Like Laravel's Arr::has (dot path).
   */
  has(obj: AnyObject, path: string): boolean {
    const token = Symbol('missing')
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    return nestedArray.get(obj, path, token as any as any) !== (token as any as any)
  },

  /**
   * Like Laravel's Arr::forget (dot path).
   */
  forget(obj: AnyObject, path: string): AnyObject {
    return nestedArray.remove(obj, path)
  },

  /**
   * Like Laravel's Arr::only (top-level keys).
   */
  only<T extends AnyObject, K extends keyof T>(obj: T, keys: readonly K[]): Pick<T, K> {
    const out = {} as Pick<T, K>
    for (const k of keys) if (k in obj) out[k] = obj[k]
    return out
  },

  /**
   * Like Laravel's Arr::except (top-level keys).
   */
  except<T extends AnyObject, K extends keyof T>(obj: T, keys: readonly K[]): Omit<T, K> {
    const omit = new Set(keys as readonly (keyof T)[])
    const out: AnyObject = {}
    for (const k in obj) if (!omit.has(k as keyof T)) out[k] = obj[k]
    return out as Omit<T, K>
  },

  /**
   * Like Laravel's Arr::pluck for arrays of objects (dot paths supported).
   * - If keyPath is provided => returns a dictionary keyed by that value.
   * - Else => returns a simple array of values.
   */
  pluck<T extends AnyObject, V = unknown>(
    items: T[],
    valuePath: string,
    keyPath?: string,
  ): V[] | Record<string, V> {
    if (!keyPath) return items.map((it) => nestedArray.get<V>(it, valuePath) as V)

    const out: Record<string, V> = {}
    for (const it of items) {
      const k = nestedArray.get<any>(it, keyPath)
      out[String(k)] = nestedArray.get<V>(it, valuePath) as V
    }
    return out
  },

  /**
   * Like Laravel's Arr::first for arrays.
   */
  first<T>(arr: T[], predicate?: (item: T, index: number) => boolean): T | undefined {
    if (!predicate) return arr[0]
    for (let i = 0; i < arr.length; i++) if (predicate(arr[i] as T, i)) return arr[i]
    return undefined
  },

  /**
   * Like Laravel's Arr::last for arrays.
   */
  last<T>(arr: T[]): T | undefined {
    return arr.length ? arr[arr.length - 1] : undefined
  },

  /**
   * Like Laravel's Arr::wrap.
   */
  wrap<T>(value: T | T[] | null | undefined): T[] {
    if (value == null) return []
    return Array.isArray(value) ? value : [value]
  },
} as const
