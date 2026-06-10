import { useRequest } from '@/composables/useRequest'
import type { AxiosRequestConfig } from 'axios'
import type { UseRequestConfig } from '@/plugins/client/types'
import { MediaApiRoutes } from './enums'
import { UploadResponseSchema, type UploadResponse } from './types'

export type UploadFieldName = 'artistImage' | 'albumCover' | 'trackCover' | 'trackAudio'

export const useMediaApi = () => {
  const uploadAdminMedia = async (
    fieldName: UploadFieldName,
    file: File,
    config?: UseRequestConfig<UploadResponse>,
    requestConfig?: Pick<AxiosRequestConfig, 'onUploadProgress' | 'signal'>,
  ) => {
    const formData = new FormData()

    formData.append(fieldName, file)

    return useRequest<UploadResponse>(
      MediaApiRoutes.ADMIN_UPLOAD,
      {
        method: 'POST',
        data: formData,
        ...requestConfig,
      },
      {
        schema: UploadResponseSchema,
        silent: false,
        ...config,
      },
    )
  }

  return {
    uploadAdminMedia,
  }
}
