import { computed, ref, watch } from 'vue'

const STORAGE_KEY = 'muse-rtl'

/**
 * Composable for managing RTL (right-to-left) direction state.
 * Persists to localStorage and updates <html> dir attribute.
 */
function getInitialRTL(): boolean {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored !== null) return stored === 'true'
  // Default based on <html> dir attribute
  return document.documentElement.getAttribute('dir') === 'rtl'
}

const isRTL = ref(getInitialRTL())

// Sync <html> dir attribute
watch(isRTL, (val) => {
  document.documentElement.setAttribute('dir', val ? 'rtl' : 'ltr')
  localStorage.setItem(STORAGE_KEY, String(val))
}, { immediate: true })

export function useRTL() {
  const dir = computed(() => isRTL.value ? 'rtl' : 'ltr')

  function toggleRTL() {
    isRTL.value = !isRTL.value
  }

  function setRTL(val: boolean) {
    isRTL.value = val
  }

  return {
    isRTL,
    dir,
    toggleRTL,
    setRTL,
  }
}
