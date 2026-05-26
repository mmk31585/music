import CryptoJS from 'crypto-js'
import Cookie from 'js-cookie'
import { v4 as uuidv4 } from 'uuid'
import isObject from 'lodash.isobject'

const COOKIE_NAME = 'ensureSafeDataLocal'

// Get or generate encryption token
const encryptionToken = Cookie.get(COOKIE_NAME) ?? uuidv4()

// Persist token (secure: set true in prod over HTTPS)
Cookie.set(COOKIE_NAME, encryptionToken, { secure: false, expires: 180, sameSite: 'lax' })

export const safeLocalStorage = {
  getItem<T = unknown>(key: string): T | string | null {
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
      return JSON.parse(decrypted) as T
    } catch {
      return decrypted
    }
  },

  setItem<T = unknown>(key: string, value: T): void {
    if (!window) return

    let payload: string

    if (isObject(value) || Array.isArray(value)) {
      payload = JSON.stringify(value)
    } else {
      payload = String(value)
    }

    const encrypted = CryptoJS.AES.encrypt(payload, encryptionToken).toString()
    window.localStorage.setItem(key, encrypted)
  },

  removeItem(key: string): void {
    if (!window) return

    window.localStorage.removeItem(key)
  },
}
