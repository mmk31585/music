import { onMounted, onUnmounted, type Ref } from 'vue'

interface LongPressOptions {
  /** Element to listen on. */
  element: Ref<HTMLElement | null>
  /** Callback fired after the hold duration. */
  onLongPress: (e: TouchEvent) => void
  /** Hold duration in ms before firing (default: 500). */
  delay?: number
  /** Movement threshold in px before cancelling (default: 10). */
  cancelThreshold?: number
}

/**
 * Long press detection composable.
 *
 * Fires `onLongPress` after the user holds a touch for `delay` ms
 * (default 500ms). Cancels if the finger moves more than `cancelThreshold` px.
 * Only listens on touch devices — no-op if no touch support is detected.
 */
export function useLongPress(options: LongPressOptions) {
  const { element, onLongPress, delay = 500, cancelThreshold = 10 } = options

  let timer: ReturnType<typeof setTimeout> | null = null
  let startX = 0
  let startY = 0
  let active = false

  function onTouchStart(e: TouchEvent) {
    if (e.touches.length !== 1) return
    active = true
    startX = e.touches[0]!.clientX
    startY = e.touches[0]!.clientY

    timer = setTimeout(() => {
      if (active) {
        onLongPress(e)
      }
    }, delay)
  }

  function onTouchMove(e: TouchEvent) {
    if (!active || !timer) return
    const dx = e.touches[0]!.clientX - startX
    const dy = e.touches[0]!.clientY - startY
    if (Math.abs(dx) > cancelThreshold || Math.abs(dy) > cancelThreshold) {
      cancel()
    }
  }

  function onTouchEnd() {
    cancel()
  }

  function cancel() {
    active = false
    if (timer) {
      clearTimeout(timer)
      timer = null
    }
  }

  onMounted(() => {
    const el = element.value
    if (!el) return
    // Only attach if touch is supported
    if (!('ontouchstart' in window)) return

    el.addEventListener('touchstart', onTouchStart, { passive: true })
    el.addEventListener('touchmove', onTouchMove, { passive: true })
    el.addEventListener('touchend', onTouchEnd, { passive: true })
    el.addEventListener('touchcancel', cancel, { passive: true })
  })

  onUnmounted(() => {
    const el = element.value
    if (!el) return
    el.removeEventListener('touchstart', onTouchStart)
    el.removeEventListener('touchmove', onTouchMove)
    el.removeEventListener('touchend', onTouchEnd)
    el.removeEventListener('touchcancel', cancel)
    cancel()
  })
}
