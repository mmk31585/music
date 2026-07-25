import { createApp, readonly, ref, type App } from 'vue'
import { type Pinia } from 'pinia'
import { useDocumentPictureInPicture } from './useDocumentPictureInPicture'
import PiPPlayerContent from '@/components/music/player/PiPPlayerContent.vue'

const initialized = ref(false)
let api: ReturnType<typeof useDocumentPictureInPicture> | null = null
let mountedApp: App | null = null

function getMainPinia(): Pinia | null {
  return (window as any).__PINIA__ ?? null
}

export function usePlayerPiPController() {
  if (!initialized.value) {
    api = useDocumentPictureInPicture({
      width: 380,
      height: 220,
      title: 'Now Playing',
      copyStyles: false, // PiPPlayerContent uses inline styles — no need to copy app CSS
      onClose: () => {
        destroyMountedApp()
      },
    })
    initialized.value = true
  }

  if (!api) {
    throw new Error('PiP controller failed to initialize.')
  }

  function destroyMountedApp() {
    if (mountedApp) {
      mountedApp.unmount()
      mountedApp = null
    }
  }

  async function safeOpen(): Promise<Window> {
    destroyMountedApp()
    const pipWin = await api!.open()

    // Mount PiPPlayerContent directly into the PiP window using createApp.
    // Vue's Teleport doesn't work cross-document, so we mount a mini Vue
    // app that shares the main app's Pinia instance (so stores are shared).
    if (api!.mountEl.value) {
      const pinia = getMainPinia()
      if (!pinia) {
        api!.close()
        throw new Error('Pinia not initialized. The PiP player needs the main app Pinia instance.')
      }
      mountedApp = createApp(PiPPlayerContent, {
        onClose: () => {
          api!.close()
          destroyMountedApp()
        },
      })
      mountedApp.use(pinia)
      mountedApp.mount(api!.mountEl.value)
    }

    return pipWin
  }

  function safeClose() {
    destroyMountedApp()
    api!.close()
  }

  async function toggle() {
    if (api!.isOpen.value) {
      safeClose()
    } else {
      await safeOpen()
    }
  }

  return {
    isSupported: readonly(api.isSupported),
    isOpen: readonly(api.isOpen),
    mountEl: readonly(api.mountEl),
    open: safeOpen,
    close: safeClose,
    toggle,
  }
}
