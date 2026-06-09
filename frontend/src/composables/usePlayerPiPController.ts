import { readonly, ref } from 'vue'
import { useDocumentPictureInPicture } from './useDocumentPictureInPicture'

const initialized = ref(false)
let api: ReturnType<typeof useDocumentPictureInPicture> | null = null

export function usePlayerPiPController() {
  if (!initialized.value) {
    api = useDocumentPictureInPicture({
      width: 380,
      height: 220,
      title: 'Now Playing',
    })
    initialized.value = true
  }

  if (!api) {
    throw new Error('PiP controller failed to initialize.')
  }

  return {
    isSupported: readonly(api.isSupported),
    isOpen: readonly(api.isOpen),
    mountEl: readonly(api.mountEl),
    open: api.open,
    close: api.close,
    async toggle() {
      if (api!.isOpen.value) api!.close()
      else await api!.open()
    },
  }
}
