import { ref } from 'vue'
import { useMediaApi } from '@/services/api'
import type { UploadFieldName } from '@/services/api/media/routes'

type UploadTarget = 'artist-image' | 'album-cover' | 'track-cover'

export function useCatalogImageUpload() {
  const uploadingImage = ref(false)
  const uploadImageError = ref<string | null>(null)

  const { uploadAdminMedia } = useMediaApi()

  async function uploadCatalogImage(file: File, target: UploadTarget): Promise<string> {
    uploadingImage.value = true
    uploadImageError.value = null

    try {
      const fieldName = getFieldName(target)
      const uploaded = await uploadAdminMedia(fieldName, file)

      const url = uploaded.url || uploaded.file_url || uploaded.fileUrl || uploaded.path

      if (!url) {
        throw new Error('Upload response did not include image URL')
      }

      return url
    } catch (error) {
      uploadImageError.value = 'Could not upload image.'
      throw error
    } finally {
      uploadingImage.value = false
    }
  }

  return {
    uploadingImage,
    uploadImageError,
    uploadCatalogImage,
  }
}

function getFieldName(target: UploadTarget): UploadFieldName {
  switch (target) {
    case 'artist-image':
      return 'artistImage'
    case 'album-cover':
      return 'albumCover'
    case 'track-cover':
      return 'trackCover'
  }
}
