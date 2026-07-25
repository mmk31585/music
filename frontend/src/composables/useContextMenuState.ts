import { ref } from 'vue'

/**
 * Module-level singleton state for context menu management.
 *
 * Ensures only one context menu (desktop or mobile) is open at any time.
 * When a new menu opens, the previously open one is closed automatically.
 */
let currentMenuId: string | null = null
const isAnyMenuOpen = ref(false)

/**
 * Generates a unique menu ID per component instance.
 */
let menuIdCounter = 0
function nextMenuId(): string {
  return `ctx-menu-${++menuIdCounter}`
}

/**
 * Composable for managing context menu open/close state.
 *
 * Each ContextMenu instance calls `acquire()` on open and `release()` on close.
 * If another menu is already open, it is force-closed first.
 */
export function useContextMenuState() {
  const menuId = nextMenuId()
  const isOpen = ref(false)

  function open() {
    // Close any previously open menu
    if (currentMenuId && currentMenuId !== menuId) {
      window.dispatchEvent(new CustomEvent('ctx-menu:close-all'))
      currentMenuId = null
      isAnyMenuOpen.value = false
    }

    currentMenuId = menuId
    isOpen.value = true
    isAnyMenuOpen.value = true
  }

  function close() {
    if (currentMenuId === menuId) {
      currentMenuId = null
      isOpen.value = false
      isAnyMenuOpen.value = false
    }
  }

  function forceClose() {
    isOpen.value = false
    if (currentMenuId === menuId) {
      currentMenuId = null
      isAnyMenuOpen.value = false
    }
  }

  return {
    /** Unique ID for this menu instance. */
    menuId,
    /** Whether this specific menu is currently open. */
    isOpen,
    /** Whether any context menu in the app is open. */
    isAnyMenuOpen,
    /** Open this menu (closes any other open menu first). */
    open,
    /** Close this menu. */
    close,
    /** Force-close regardless of ownership. */
    forceClose,
  }
}

/**
 * Close all open context menus globally.
 * Useful for global listeners (e.g. route change, scroll).
 */
export function closeAllContextMenus() {
  // Dispatch a custom event that all ContextMenu instances listen to
  window.dispatchEvent(new CustomEvent('ctx-menu:close-all'))
  currentMenuId = null
  isAnyMenuOpen.value = false
}
