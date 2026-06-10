import { onMounted, onUnmounted } from 'vue'

interface ShortcutDef {
  key: string
  ctrl?: boolean
  meta?: boolean
  shift?: boolean
  alt?: boolean
  handler: (e: KeyboardEvent) => void
  prevent?: boolean
  description?: string
}

export function useShortcuts(shortcuts: ShortcutDef[]) {
  function onKeydown(e: KeyboardEvent) {
    const tag = (e.target as HTMLElement)?.tagName
    if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return

    for (const s of shortcuts) {
      const ctrlOrMeta = s.ctrl || s.meta
      const matchCtrl = ctrlOrMeta ? (e.ctrlKey || e.metaKey) : true
      const matchShift = s.shift ? e.shiftKey : !e.shiftKey
      const matchAlt = s.alt ? e.altKey : !e.altKey
      const matchKey = e.key.toLowerCase() === s.key.toLowerCase()

      if (matchKey && matchCtrl && matchShift && matchAlt) {
        if (matchCtrl && e.ctrlKey !== e.metaKey && !!s.meta !== !!s.ctrl) continue
        if (s.prevent !== false) e.preventDefault()
        s.handler(e)
        return
      }
    }
  }

  onMounted(() => document.addEventListener('keydown', onKeydown))
  onUnmounted(() => document.removeEventListener('keydown', onKeydown))
}

export function usePlayerShortcuts(handlers: {
  togglePlay?: () => void
  next?: () => void
  previous?: () => void
  seekForward?: () => void
  seekBackward?: () => void
  volumeUp?: () => void
  volumeDown?: () => void
  toggleMute?: () => void
  toggleShuffle?: () => void
  toggleRepeat?: () => void
}) {
  const shortcuts: ShortcutDef[] = [
    { key: ' ', handler: () => handlers.togglePlay?.(), description: 'Play/Pause' },
    { key: 'ArrowRight', handler: () => handlers.seekForward?.(), description: 'Seek forward 5s' },
    { key: 'ArrowLeft', handler: () => handlers.seekBackward?.(), description: 'Seek backward 5s' },
    { key: 'ArrowUp', handler: () => handlers.volumeUp?.(), description: 'Volume up' },
    { key: 'ArrowDown', handler: () => handlers.volumeDown?.(), description: 'Volume down' },
    { key: 'n', handler: () => handlers.next?.(), description: 'Next track' },
    { key: 'p', handler: () => handlers.previous?.(), description: 'Previous track' },
    { key: 'm', handler: () => handlers.toggleMute?.(), description: 'Toggle mute' },
    { key: 's', handler: () => handlers.toggleShuffle?.(), description: 'Toggle shuffle' },
    { key: 'r', handler: () => handlers.toggleRepeat?.(), description: 'Toggle repeat' },
  ].filter((s) => s.handler) as ShortcutDef[]

  useShortcuts(shortcuts)
}
