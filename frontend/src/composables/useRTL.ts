import { computed } from 'vue'
import { useLocaleStore } from '@/stores/locale'

/**
 * Composable for managing RTL (right-to-left) direction state.
 * Derives direction from the locale store (single source of truth).
 *
 * The locale store syncs <html dir>, <html lang>, and PrimeVue locale.
 * This composable provides a convenience interface for components that
 * need RTL-aware logic (e.g., animation direction).
 */
export function useRTL() {
  const localeStore = useLocaleStore()

  const isRTL = computed(() => localeStore.locale === 'fa')
  const dir = computed(() => localeStore.dir())

  function toggleRTL() {
    localeStore.toggleLocale()
  }

  function setRTL(val: boolean) {
    localeStore.setLocale(val ? 'fa' : 'en')
  }

  return {
    isRTL,
    dir,
    toggleRTL,
    setRTL,
  }
}
