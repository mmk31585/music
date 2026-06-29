import { computed, onBeforeUnmount, ref } from 'vue'

export interface UseDocumentPiPOptions {
  width?: number
  height?: number
  title?: string
  copyStyles?: boolean
  onClose?: () => void
}

function copyDocumentStyles(sourceDoc: Document, targetDoc: Document) {
  const seen = new Set<string>()

  const addElement = (el: Element) => {
    if (el.tagName === 'LINK') {
      const link = el as HTMLLinkElement
      if (!link.href || link.rel !== 'stylesheet') return
      if (seen.has(link.href)) return
      seen.add(link.href)
      const clone = targetDoc.createElement('link')
      clone.rel = 'stylesheet'
      clone.href = link.href
      targetDoc.head.appendChild(clone)
      return
    }

    if (el.tagName === 'STYLE') {
      const style = el as HTMLStyleElement
      const text = style.textContent || ''
      if (!text.trim()) return
      if (seen.has(text)) return
      seen.add(text)
      const clone = targetDoc.createElement('style')
      clone.textContent = text
      targetDoc.head.appendChild(clone)
    }
  }

  for (const el of Array.from(sourceDoc.head.querySelectorAll('link[rel="stylesheet"], style'))) {
    addElement(el)
  }

  const fontLink = sourceDoc.querySelector('link[href*="primeicons"]')
  if (fontLink) addElement(fontLink)

  const meta = targetDoc.createElement('meta')
  meta.name = 'viewport'
  meta.content = 'width=device-width, initial-scale=1.0'
  targetDoc.head.appendChild(meta)
}

export function useDocumentPictureInPicture(options: UseDocumentPiPOptions = {}) {
  const pipWindow = ref<Window | null>(null)
  const mountEl = ref<HTMLElement | null>(null)
  const isOpen = ref(false)

  const isSupported = computed(() => {
    return typeof window !== 'undefined' && !!window.documentPictureInPicture
  })

  let removePageHideListener: (() => void) | null = null
  let removeBeforeUnloadListener: (() => void) | null = null
  let animationFrameId: number | null = null

  function cleanupListeners() {
    removePageHideListener?.()
    removePageHideListener = null
    removeBeforeUnloadListener?.()
    removeBeforeUnloadListener = null
    if (animationFrameId !== null) {
      cancelAnimationFrame(animationFrameId)
      animationFrameId = null
    }
  }

  function cleanupWindow(closeWindow = false) {
    cleanupListeners()

    if (closeWindow && pipWindow.value && !pipWindow.value.closed) {
      pipWindow.value.close()
    }

    mountEl.value = null
    pipWindow.value = null
    isOpen.value = false
  }

  function syncBodyStyles() {
    if (!pipWindow.value || pipWindow.value.closed) return
    const pipBody = pipWindow.value.document.body
    pipBody.style.margin = '0'
    pipBody.style.background = '#121212'
    pipBody.style.overflow = 'hidden'
    pipBody.style.width = '100vw'
    pipBody.style.height = '100vh'
  }

  async function open(): Promise<Window> {
    if (!isSupported.value || !window.documentPictureInPicture) {
      throw new Error('Document Picture-in-Picture is not supported in this browser.')
    }

    if (pipWindow.value && !pipWindow.value.closed) {
      pipWindow.value.focus()
      return pipWindow.value
    }

    const pip = await window.documentPictureInPicture.requestWindow({
      width: options.width ?? 380,
      height: options.height ?? 220,
      preferInitialWindowPlacement: true,
    })

    pipWindow.value = pip
    isOpen.value = true

    pip.document.title = options.title ?? 'Mini Player'
    pip.document.body.innerHTML = ''

    syncBodyStyles()

    // Optionally copy app styles into the PiP window (needed for class-based CSS).
    // PiPPlayerContent uses inline styles so this can be skipped for speed.
    if (options.copyStyles !== false) {
      copyDocumentStyles(document, pip.document)
    }

    const appRoot = pip.document.createElement('div')
    appRoot.id = 'pip-player-root'
    appRoot.style.width = '100vw'
    appRoot.style.height = '100vh'
    pip.document.body.appendChild(appRoot)

    mountEl.value = appRoot

    const handleClose = () => {
      options.onClose?.()
      cleanupWindow(false)
    }

    const handleBeforeUnload = () => {
      options.onClose?.()
      cleanupWindow(false)
    }

    const handleResize = () => {
      if (animationFrameId !== null) cancelAnimationFrame(animationFrameId)
      animationFrameId = requestAnimationFrame(() => {
        syncBodyStyles()
        if (appRoot) {
          appRoot.style.width = `${pip.innerWidth}px`
          appRoot.style.height = `${pip.innerHeight}px`
        }
      })
    }

    pip.addEventListener('pagehide', handleClose)
    pip.addEventListener('beforeunload', handleBeforeUnload)
    pip.addEventListener('resize', handleResize)

    removePageHideListener = () => pip.removeEventListener('pagehide', handleClose)
    removeBeforeUnloadListener = () => pip.removeEventListener('beforeunload', handleBeforeUnload)

    return pip
  }

  function close() {
    cleanupWindow(true)
  }

  onBeforeUnmount(() => {
    cleanupWindow(false)
  })

  return {
    isSupported,
    isOpen,
    pipWindow,
    mountEl,
    open,
    close,
  }
}
