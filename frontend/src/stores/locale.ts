import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Locale } from '@/locales'
import { DEFAULT_LOCALE, SUPPORTED_LOCALES } from '@/locales'
import { primeLocale as primeLocaleFa } from '@/utils/prime/prime-locale'

/**
 * Locale store — manages the active UI language and direction.
 *
 * Persists the locale preference to localStorage and syncs it with:
 * - vue-i18n (message translations)
 * - PrimeVue locale (component labels)
 * - HTML `lang` and `dir` attributes on <html>
 * - `useRTL()` composable for RTL-aware logic
 */
export const useLocaleStore = defineStore('locale', () => {
  // ── State ─────────────────────────────────────────────────────────────
  const locale = ref<Locale>(loadSavedLocale())

  // ── Computed ──────────────────────────────────────────────────────────
  const isRTL = () => locale.value === 'fa'
  const dir = () => (locale.value === 'fa' ? 'rtl' : 'ltr')

  // ── Actions ───────────────────────────────────────────────────────────
  function setLocale(newLocale: Locale) {
    if (!SUPPORTED_LOCALES.includes(newLocale)) return

    locale.value = newLocale
    saveLocale(newLocale)

    // Sync vue-i18n
    const { locale: i18nLocale } = useI18n()
    i18nLocale.value = newLocale

    // Sync HTML attributes
    document.documentElement.lang = newLocale
    document.documentElement.dir = newLocale === 'fa' ? 'rtl' : 'ltr'

    // Dispatch event for PrimeVue locale update (caught by main.ts listener)
    window.dispatchEvent(
      new CustomEvent('locale-change', { detail: { locale: newLocale } }),
    )
  }

  function toggleLocale() {
    const next: Locale = locale.value === 'fa' ? 'en' : 'fa'
    setLocale(next)
  }

  return {
    locale,
    isRTL,
    dir,
    setLocale,
    toggleLocale,
  }
})

// ── Persistence helpers ──────────────────────────────────────────────────
function loadSavedLocale(): Locale {
  try {
    const saved = localStorage.getItem('muse-locale') as Locale | null
    if (saved && (saved === 'fa' || saved === 'en')) return saved
  } catch {
    // localStorage not available
  }
  return DEFAULT_LOCALE
}

function saveLocale(locale: Locale) {
  try {
    localStorage.setItem('muse-locale', locale)
  } catch {
    // localStorage not available
  }
}

/**
 * Returns the PrimeVue locale object for a given locale.
 * Currently only Persian has a full override; English uses defaults.
 */
export function getPrimeLocale(locale: Locale): Record<string, any> | undefined {
  if (locale === 'fa') return primeLocaleFa
  return undefined
}
