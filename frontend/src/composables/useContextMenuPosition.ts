import { ref, nextTick, type Ref, type CSSProperties } from 'vue'
import { useRTL } from '@/composables/useRTL'

const VIEWPORT_PADDING = 12

export interface PositionOptions {
  /** Cursor position (for right-click context menus). */
  position: Ref<{ x: number; y: number } | undefined>
  /** Element to anchor the menu to (for overflow menus). */
  anchorEl: Ref<HTMLElement | null | undefined>
  /** Ref to the menu DOM element (used to measure dimensions). */
  menuEl: Ref<HTMLElement | null>
}

/**
 * Viewport-aware positioning composable.
 *
 * Measures the menu element, checks all four edges with 12px padding,
 * and flips placement when the menu would overflow the viewport.
 *
 * Returns a reactive CSSProperties style to apply to the menu container.
 */
export function useContextMenuPosition(options: PositionOptions) {
  const { position, anchorEl, menuEl } = options
  const { isRTL } = useRTL()

  const style = ref<CSSProperties>({})
  const placement = ref<'tl' | 'tr' | 'bl' | 'br'>('tl')
  const isReady = ref(false)

  /**
   * Recompute the position based on current state.
   * Must be called after the menu element is mounted and visible.
   */
  async function recompute() {
    await nextTick()

    const el = menuEl.value
    if (!el) {
      isReady.value = false
      return
    }

    const menuRect = el.getBoundingClientRect()
    const menuW = menuRect.width
    const menuH = menuRect.height
    const vpW = window.innerWidth
    const vpH = window.innerHeight

    let x = 0
    let y = 0

    // Determine origin point
    if (anchorEl.value) {
      const anchorRect = anchorEl.value.getBoundingClientRect()
      // Anchor menus: place above-right (or above-left in RTL)
      x = isRTL.value ? anchorRect.left : anchorRect.right - menuW
      y = anchorRect.top - menuH - 4
    } else if (position.value) {
      x = position.value.x
      y = position.value.y
    }

    // Clamp and flip
    // Right edge
    if (x + menuW > vpW - VIEWPORT_PADDING) {
      // Flip left
      if (anchorEl.value) {
        const anchorRect = anchorEl.value.getBoundingClientRect()
        x = isRTL.value ? anchorRect.right : anchorRect.left - menuW
      } else {
        x = position.value ? position.value.x - menuW : vpW - menuW - VIEWPORT_PADDING
      }
    }

    // Left edge
    if (x < VIEWPORT_PADDING) {
      x = VIEWPORT_PADDING
    }

    // Bottom edge
    if (y + menuH > vpH - VIEWPORT_PADDING) {
      // Flip above
      if (anchorEl.value) {
        const anchorRect = anchorEl.value.getBoundingClientRect()
        y = anchorRect.bottom + 4
      } else {
        y = position.value ? position.value.y - menuH : vpH - menuH - VIEWPORT_PADDING
      }
    }

    // Top edge
    if (y < VIEWPORT_PADDING) {
      y = VIEWPORT_PADDING
    }

    // Determine final placement quadrant
    if (anchorEl.value) {
      const anchorRect = anchorEl.value.getBoundingClientRect()
      const isBelow = y > anchorRect.bottom
      const isLeft = x < anchorRect.left
      placement.value = isBelow
        ? (isLeft ? 'bl' : 'br')
        : (isLeft ? 'tl' : 'tr')
    } else {
      placement.value = 'tl'
    }

    style.value = {
      position: 'fixed',
      left: `${x}px`,
      top: `${y}px`,
      zIndex: 300,
    }

    isReady.value = true
  }

  /**
   * Reset state for reuse.
   */
  function reset() {
    style.value = {}
    placement.value = 'tl'
    isReady.value = false
  }

  return {
    /** Reactive CSS style to apply to the menu container. */
    style,
    /** Which corner the menu is anchored to. */
    placement,
    /** Whether positioning has been computed. */
    isReady,
    /** Recompute position (call after mount/visibility change). */
    recompute,
    /** Reset positioning state. */
    reset,
  }
}
