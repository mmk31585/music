import { computed, onMounted, onUnmounted } from 'vue'
import { usePlayerStore } from '@/stores/player'
import { useTrackLike } from '@/composables/player'

export interface KeyboardShortcutsOptions {
  /** Opens/closes fullscreen player */
  onToggleFullscreen?: () => void
  /** Opens/closes queue panel */
  onToggleQueue?: () => void
  /** Opens/closes lyrics view */
  onToggleLyrics?: () => void
  /** Toggles the ? shortcuts help panel */
  onToggleShortcuts?: () => void
  /** Opens/closes mobile bottom sheet */
  onToggleMobileSheet?: () => void
}

/**
 * Registers global keyboard shortcuts for player controls.
 * Call once from a root-level component (e.g. PlayerRegion.vue).
 *
 * Shortcuts:
 *   Space     — Play/Pause
 *   N         — Next track
 *   P         — Previous track
 *   ArrowRight — Seek forward 5s
 *   ArrowLeft  — Seek back 5s
 *   ArrowUp    — Volume +0.1
 *   ArrowDown  — Volume -0.1
 *   M         — Toggle mute
 *   S         — Cycle shuffle modes
 *   R         — Toggle repeat mode
 *   L         — Like/unlike current track
 *   F         — Toggle fullscreen player
 *   Q         — Toggle queue panel
 *   ?         — Toggle shortcuts help
 *   Escape    — Close bottom sheet (handled per-component)
 */
export function useKeyboardShortcuts(options: KeyboardShortcutsOptions = {}) {
  const store = usePlayerStore()

  const trackId = computed(() => store.currentTrack?.id)
  const { toggleLike } = useTrackLike(trackId)

  function handleKeydown(e: KeyboardEvent) {
    // Don't intercept when user is typing in an input or textarea
    const tag = (e.target as HTMLElement)?.tagName?.toLowerCase()
    if (tag === 'input' || tag === 'textarea' || tag === 'select') return
    if ((e.target as HTMLElement)?.isContentEditable) return

    // ── Ctrl/Cmd combinations ──
    if (e.ctrlKey || e.metaKey) {
      switch (e.key.toLowerCase()) {
        case 'b':
          // Toggle sidebar — no-op if no sidebar handler
          break
        case 'k':
          // Focus search — handled by global search overlay
          break
        case 'l':
          e.preventDefault()
          options.onToggleLyrics?.()
          break
      }
      return
    }

    // ── Single-key shortcuts ──
    switch (e.key.toLowerCase()) {
      case ' ':
        e.preventDefault()
        store.resume().catch(() => {})
        break

      case 'n':
        store.playNext()
        break

      case 'p':
        store.playPrevious()
        break

      case 'l':
        toggleLike()
        break

      case 'm':
        store.toggleMute()
        break

      case 's':
        store.toggleShuffle()
        break

      case 'r':
        store.toggleRepeat()
        break

      case 'f':
        options.onToggleFullscreen?.()
        break

      case 'q':
        options.onToggleQueue?.()
        break

      case '?':
        e.preventDefault()
        options.onToggleShortcuts?.()
        break

      // ── Arrow keys ──
      case 'arrowright':
        e.preventDefault()
        store.seek(Math.min(store.currentTime + 5, store.duration))
        break

      case 'arrowleft':
        e.preventDefault()
        store.seek(Math.max(store.currentTime - 5, 0))
        break

      case 'arrowup':
        e.preventDefault()
        store.setVolume(Math.min(1, store.volume + 0.1))
        break

      case 'arrowdown':
        e.preventDefault()
        store.setVolume(Math.max(0, store.volume - 0.1))
        break
    }
  }

  onMounted(() => {
    document.addEventListener('keydown', handleKeydown)
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeydown)
  })
}
