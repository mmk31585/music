import { ref, onMounted, onUnmounted } from 'vue'

/**
 * Reactive online/offline detection.
 * Provides a boolean `isOnline` ref that updates when the browser
 * fires `online`/`offline` events.
 *
 * Usage:
 *   const { isOnline } = useOnlineStatus()
 */
export function useOnlineStatus() {
  const isOnline = ref(typeof navigator !== 'undefined' ? navigator.onLine : true)

  function onOnline() { isOnline.value = true }
  function onOffline() { isOnline.value = false }

  onMounted(() => {
    window.addEventListener('online', onOnline)
    window.addEventListener('offline', onOffline)
  })

  onUnmounted(() => {
    window.removeEventListener('online', onOnline)
    window.removeEventListener('offline', onOffline)
  })

  return { isOnline }
}
