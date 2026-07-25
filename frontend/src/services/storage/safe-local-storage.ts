import CryptoJS from 'crypto-js'
import Cookie from 'js-cookie'
import { v4 as uuidv4 } from 'uuid'

const COOKIE_NAME = 'ensureSafeDataLocal'

// Get or generate encryption token
const encryptionToken = Cookie.get(COOKIE_NAME) ?? uuidv4()

// Persist token (secure: set true in prod over HTTPS)
Cookie.set(COOKIE_NAME, encryptionToken, { secure: false, expires: 180, sameSite: 'lax' })

export const safeLocalStorage = {
  /**
   * Retrieve a typed value from encrypted localStorage.
   * NOTE: The caller is responsible for ensuring the stored data matches type `T`.
   * The `as T` cast is intentionally explicit — JSON.parse returns `unknown`,
   * and we assert the shape through the generic parameter.
   */
  getItem<T = unknown>(key: string): T | null {
    if (!window) return null

    const store = window.localStorage.getItem(key)
    if (!store) return null

    let decrypted: string
    try {
      const bytes = CryptoJS.AES.decrypt(store, encryptionToken)
      decrypted = bytes.toString(CryptoJS.enc.Utf8)
    } catch {
      window.localStorage.removeItem(key)
      return null
    }

    if (!decrypted) return null

    try {
      // Cast through `unknown` first to satisfy strict TypeScript —
      // the caller asserts `T` via the generic parameter.
      const parsed: unknown = JSON.parse(decrypted)
      return parsed as T
    } catch {
      // Stored value is not valid JSON — clean up and return null
      window.localStorage.removeItem(key)
      return null
    }
  },

  setItem<T = unknown>(key: string, value: T): void {
    if (!window) return

    // Always serialize to JSON for consistent deserialization in getItem
    const payload = JSON.stringify(value)
    const encrypted = CryptoJS.AES.encrypt(payload, encryptionToken).toString()
    window.localStorage.setItem(key, encrypted)
  },

  removeItem(key: string): void {
    if (!window) return

    window.localStorage.removeItem(key)
  },
}
