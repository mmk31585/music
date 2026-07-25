import { useRequest } from '@/composables/useRequest'
import type { AxiosRequestConfig } from 'axios'
import type { UseRequestConfig } from '@/plugins/client/types'
import { MediaApiRoutes } from './enums'
import { UploadResponseSchema, type UploadResponse, type Media } from './types'
import { apiReplaceParams } from '@/utils/api-replace-params'

export type UploadFieldName = 'artistImage' | 'albumCover' | 'trackCover' | 'trackAudio' | 'playlistCover'

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

  const listAdminMedia = async (config?: UseRequestConfig<Media[]>) => {
    return useRequest<Media, true>(
      MediaApiRoutes.ADMIN_LIST,
      { method: 'GET' },
      { silent: false, ...config },
    )
  }

  const deleteAdminMedia = async (id: string, config?: UseRequestConfig<any>) => {
    return useRequest(
      apiReplaceParams(MediaApiRoutes.ADMIN_DELETE, { id }),
      { method: 'DELETE' },
      { silent: false, ...config },
    )
  }

  return {
    uploadAdminMedia,
    listAdminMedia,
    deleteAdminMedia,
  }
}
