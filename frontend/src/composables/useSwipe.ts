import { onMounted, onUnmounted, type Ref } from 'vue'

interface SwipeOptions {
  element: Ref<HTMLElement | null>
  onSwipeLeft?: () => void
  onSwipeRight?: () => void
  threshold?: number
}

export function useSwipe(options: SwipeOptions) {
  const { element, onSwipeLeft, onSwipeRight, threshold = 50 } = options
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
    // Prevent default to stop scroll if we detect horizontal swipe
    const dx = e.touches[0]!.clientX - startX
    const dy = e.touches[0]!.clientY - startY
    if (Math.abs(dx) > Math.abs(dy) && Math.abs(dx) > 10) {
      e.preventDefault()
    }
  }

  function onTouchEnd(e: TouchEvent) {
    if (!isTracking) return
    isTracking = false
    const dx = e.changedTouches[0]!.clientX - startX
    const dy = e.changedTouches[0]!.clientY - startY
    if (Math.abs(dx) < threshold || Math.abs(dx) < Math.abs(dy) * 1.5) return
    if (dx > 0 && onSwipeRight) onSwipeRight()
    if (dx < 0 && onSwipeLeft) onSwipeLeft()
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
