import { ref } from 'vue'
import { useRegisterSW } from 'virtual:pwa-register/vue'
import { useOnlineStatus } from './useOnlineStatus'

interface BeforeInstallPromptEvent extends Event {
  prompt: () => Promise<void>
  userChoice: Promise<{ outcome: 'accepted' | 'dismissed' }>
}

let deferredPrompt: BeforeInstallPromptEvent | null = null
let isIOS = false
let isInstalled = false

export function usePwa() {
  const { isOnline } = useOnlineStatus()

  const isInstallable = ref(false)
  const showIOSInstructions = ref(false)

  const { needRefresh, updateServiceWorker } = useRegisterSW({
    onRegisteredSW(_swUrl, registration) {
      if (registration) {
        setInterval(() => {
          registration.update().catch(() => {})
        }, 60 * 60 * 1000)
      }
    },
    onRegisterError(error) {
      console.error('SW registration failed:', error)
    },
  })

  function init() {
    if (typeof window === 'undefined') return

    isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent) && !(window as any).MSStream
    isInstalled = window.matchMedia('(display-mode: standalone)').matches || (navigator as any).standalone === true

    if (isIOS && !isInstalled) {
      showIOSInstructions.value = true
    }

    window.addEventListener('beforeinstallprompt', (e: Event) => {
      e.preventDefault()
      deferredPrompt = e as BeforeInstallPromptEvent
      isInstallable.value = true
    })

    window.addEventListener('appinstalled', () => {
      isInstallable.value = false
      isInstalled = true
      deferredPrompt = null
    })
  }

  async function promptInstall() {
    if (!deferredPrompt) return
    deferredPrompt.prompt()
    const { outcome } = await deferredPrompt.userChoice
    if (outcome === 'accepted') {
      isInstallable.value = false
      isInstalled = true
    }
    deferredPrompt = null
  }

  function dismissInstall() {
    isInstallable.value = false
    deferredPrompt = null
    try {
      localStorage.setItem('pwa-install-dismissed', String(Date.now()))
    } catch {}
  }

  function dismissIOSInstructions() {
    showIOSInstructions.value = false
    try {
      localStorage.setItem('pwa-ios-dismissed', String(Date.now()))
    } catch {}
  }

  function updateApp() {
    updateServiceWorker()
  }

  return {
    isInstallable,
    showIOSInstructions,
    isInstalled,
    isOnline,
    needRefresh,
    init,
    promptInstall,
    dismissInstall,
    dismissIOSInstructions,
    updateApp,
  }
}
