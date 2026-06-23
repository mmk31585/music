import { onMounted, onUnmounted, type Ref } from 'vue'

interface VerticalSwipeOptions {
  element: Ref<HTMLElement | null>
  onSwipeUp?: () => void
  onSwipeDown?: () => void
  threshold?: number
}

/**
 * Detects vertical touch swipes (up/down).
 * Follows the same pattern as useSwipe.ts but for vertical direction.
 * Used by VerticalVideoPlayer for the Reels-style feed.
 */
export function useVerticalSwipe(options: VerticalSwipeOptions) {
  const { element, onSwipeUp, onSwipeDown, threshold = 80 } = options
  let startX = 0
  let startY = 0
  let isTracking = false

  function onTouchStart(e: TouchEvent) {
    if (e.touches.length !== 1) return
    startX = e.touches[0]!.clientX
    startY = e.touches[0]!.clientY
    isTracking = true
  }

  function onTouchMove(e: TouchEvent) {
    if (!isTracking) return
    // Prevent default to stop scroll if we detect vertical swipe
    const dx = e.touches[0]!.clientX - startX
    const dy = e.touches[0]!.clientY - startY
    if (Math.abs(dy) > Math.abs(dx) && Math.abs(dy) > 10) {
      e.preventDefault()
    }
  }

  function onTouchEnd(e: TouchEvent) {
    if (!isTracking) return
    isTracking = false
    const dy = e.changedTouches[0]!.clientY - startY
    if (Math.abs(dy) < threshold) return
    if (dy < 0 && onSwipeUp) onSwipeUp()
    if (dy > 0 && onSwipeDown) onSwipeDown()
  }

  onMounted(() => {
    const el = element.value
    if (!el) return
    el.addEventListener('touchstart', onTouchStart, { passive: true })
    el.addEventListener('touchmove', onTouchMove, { passive: false })
    el.addEventListener('touchend', onTouchEnd, { passive: true })
  })

  onUnmounted(() => {
    const el = element.value
    if (!el) return
    el.removeEventListener('touchstart', onTouchStart)
    el.removeEventListener('touchmove', onTouchMove)
    el.removeEventListener('touchend', onTouchEnd)
  })
}
