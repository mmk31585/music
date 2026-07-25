import { ref, computed, onMounted, onUnmounted, nextTick, type Ref } from 'vue'
import type { ContextMenuAction } from '@/types/context-menu'

interface KeyboardOptions {
  /** All visible, enabled actions (flattened across sections). */
  items: Ref<ContextMenuAction[]>
  /** Ref to the menu root DOM element. */
  menuEl: Ref<HTMLElement | null>
  /** Callback to close the menu. */
  onClose: () => void
  /** Callback when an action is activated (Enter/Space). */
  onActivate: (action: ContextMenuAction) => void
}

/**
 * Shared keyboard navigation for context menus.
 *
 * Supports:
 * - ArrowUp / ArrowDown: cycle through enabled items
 * - Home / End: jump to first / last enabled item
 * - Enter / Space: activate focused item
 * - Escape: close menu
 * - Typeahead: accumulate characters over 500ms, match first-letter
 * - Focus trap: Tab / Shift+Tab cycle within menu
 * - Auto-focus first item on mount
 */
export function useContextMenuKeyboard(options: KeyboardOptions) {
  const { items, menuEl, onClose, onActivate } = options

  const focusedIndex = ref(-1)
  let typeaheadBuffer = ''
  let typeaheadTimer: ReturnType<typeof setTimeout> | null = null

  // ── helpers ────────────────────────────────────────────────────

  /** Get all enabled (non-disabled, non-hidden) item indices. */
  const enabledIndices = computed(() => {
    const indices: number[] = []
    for (let i = 0; i < items.value.length; i++) {
      const item = items.value[i]!
      if (!item.disabled && !item.hidden) {
        indices.push(i)
      }
    }
    return indices
  })

  /** Focus a specific item element by its data attribute. */
  function focusItemByIndex(index: number) {
    nextTick(() => {
      const el = menuEl.value?.querySelector(`[data-menu-item="${index}"]`) as HTMLElement | null
      el?.focus()
    })
  }

  /** Find the nearest enabled index in a direction from current focus. */
  function findNextEnabled(fromIndex: number, direction: 1 | -1): number | null {
    const indices = enabledIndices.value
    if (indices.length === 0) return null

    const currentPos = indices.indexOf(fromIndex)
    if (currentPos === -1) {
      // Current index not in enabled list — jump to first/last
      return direction === 1 ? indices[0]! : indices[indices.length - 1]!
    }

    const nextPos = (currentPos + direction + indices.length) % indices.length
    return indices[nextPos]!
  }

  // ── typeahead ──────────────────────────────────────────────────

  function handleTypeahead(char: string) {
    if (typeaheadTimer) clearTimeout(typeaheadTimer)

    typeaheadBuffer += char.toLowerCase()

    // Find first enabled item whose label starts with the accumulated buffer
    const match = enabledIndices.value.find((idx) => {
      const label = items.value[idx]?.label.toLowerCase() ?? ''
      return label.startsWith(typeaheadBuffer)
    })

    if (match !== undefined) {
      focusedIndex.value = match
      focusItemByIndex(match)
    }

    typeaheadTimer = setTimeout(() => {
      typeaheadBuffer = ''
      typeaheadTimer = null
    }, 500)
  }

  // ── keydown handler ────────────────────────────────────────────

  function onKeydown(e: KeyboardEvent) {
    const currentIdx = focusedIndex.value

    switch (e.key) {
      case 'ArrowDown': {
        e.preventDefault()
        const next = findNextEnabled(currentIdx, 1)
        if (next !== null) {
          focusedIndex.value = next
          focusItemByIndex(next)
        }
        break
      }

      case 'ArrowUp': {
        e.preventDefault()
        const prev = findNextEnabled(currentIdx, -1)
        if (prev !== null) {
          focusedIndex.value = prev
          focusItemByIndex(prev)
        }
        break
      }

      case 'Home': {
        e.preventDefault()
        const first = enabledIndices.value[0]
        if (first !== undefined) {
          focusedIndex.value = first
          focusItemByIndex(first)
        }
        break
      }

      case 'End': {
        e.preventDefault()
        const last = enabledIndices.value[enabledIndices.value.length - 1]
        if (last !== undefined) {
          focusedIndex.value = last
          focusItemByIndex(last)
        }
        break
      }

      case 'Enter':
      case ' ': {
        e.preventDefault()
        if (currentIdx >= 0 && currentIdx < items.value.length) {
          const item = items.value[currentIdx]
          if (item && !item.disabled) {
            onActivate(item)
          }
        }
        break
      }

      case 'Escape': {
        e.preventDefault()
        onClose()
        break
      }

      case 'Tab': {
        // Focus trap — prevent Tab from leaving the menu
        e.preventDefault()
        if (e.shiftKey) {
          const prev = findNextEnabled(currentIdx, -1)
          if (prev !== null) {
            focusedIndex.value = prev
            focusItemByIndex(prev)
          }
        } else {
          const next = findNextEnabled(currentIdx, 1)
          if (next !== null) {
            focusedIndex.value = next
            focusItemByIndex(next)
          }
        }
        break
      }

      default: {
        // Typeahead: single printable character
        if (e.key.length === 1 && !e.ctrlKey && !e.metaKey && !e.altKey) {
          handleTypeahead(e.key)
        }
        break
      }
    }
  }

  // ── lifecycle ──────────────────────────────────────────────────

  function focusFirst() {
    const first = enabledIndices.value[0]
    if (first !== undefined) {
      focusedIndex.value = first
      focusItemByIndex(first)
    }
  }

  function focusLast() {
    const last = enabledIndices.value[enabledIndices.value.length - 1]
    if (last !== undefined) {
      focusedIndex.value = last
      focusItemByIndex(last)
    }
  }

  onMounted(() => {
    nextTick(() => focusFirst())
  })

  onUnmounted(() => {
    if (typeaheadTimer) clearTimeout(typeaheadTimer)
  })

  return {
    /** Index of the currently focused item. */
    focusedIndex,
    /** Keydown event handler to attach to the menu root. */
    onKeydown,
    /** Programmatically focus the first enabled item. */
    focusFirst,
    /** Programmatically focus the last enabled item. */
    focusLast,
  }
}
