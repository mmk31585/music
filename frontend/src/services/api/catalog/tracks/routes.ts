import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import type { AxiosRequestConfig } from 'axios'
import { UploadResponseSchema, type UploadResponse } from '@/services/api/media'
import { TrackApiRoutes } from './enums'
import { TrackSchema, type Track, type TrackCreatePayload, type TrackUpdatePayload } from './types'

export const useTracksApi = () => {
  const getTracks = async (config?: UseRequestConfig<Track[]>) => {
    return useRequest<Track, true>(
      TrackApiRoutes.LIST,
      { method: 'GET' },
      {
        schema: TrackSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getTrack = async (id: string | number, config?: UseRequestConfig<Track>) => {
    return useRequest<Track>(
      TrackApiRoutes.GET.replace(':trackId', String(id)),
      { method: 'GET' },
      {
        schema: TrackSchema,
        silent: true,
        ...config,
      },
    )
  }

  const adminGetTracks = async (config?: UseRequestConfig<Track[]>) => {
    return useRequest<Track, true>(
      TrackApiRoutes.ADMIN_LIST,
      { method: 'GET' },
      {
        schema: TrackSchema,
        silent: true,
        ...config,
      },
    )
  }

  const adminCreateTrack = async (
    payload: TrackCreatePayload,
    config?: UseRequestConfig<Track>,
  ) => {
    return useRequest<Track>(
      TrackApiRoutes.ADMIN_CREATE,
      {
        method: 'POST',
        data: payload,
      },
      {
        schema: TrackSchema,
        silent: false,
        ...config,
      },
    )
  }

  const adminUploadTrackAudio = async (
    file: File,
    config?: UseRequestConfig<UploadResponse>,
    requestConfig?: Pick<AxiosRequestConfig, 'onUploadProgress' | 'signal'>,
  ) => {
    const formData = new FormData()
    formData.append('trackAudio', file)

    return useRequest<UploadResponse>(
      TrackApiRoutes.ADMIN_UPLOAD,
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

  const adminUploadTrackCover = async (
    file: File,
    config?: UseRequestConfig<UploadResponse>,
    requestConfig?: Pick<AxiosRequestConfig, 'onUploadProgress' | 'signal'>,
  ) => {
    const formData = new FormData()
    formData.append('trackCover', file)

    return useRequest<UploadResponse>(
      TrackApiRoutes.ADMIN_UPLOAD,
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

  const adminUpdateTrack = async (
    id: string | number,
    payload: TrackUpdatePayload,
    config?: UseRequestConfig<Track>,
  ) => {
    return useRequest<Track>(
      TrackApiRoutes.ADMIN_UPDATE.replace(':trackId', String(id)),
      {
        method: 'PATCH',
        data: payload,
      },
      {
        schema: TrackSchema,
        silent: false,
        ...config,
      },
    )
  }

  const adminDeleteTrack = async (id: string | number, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      TrackApiRoutes.ADMIN_DELETE.replace(':trackId', String(id)),
      { method: 'DELETE' },
      {
        silent: false,
        ...config,
      },
    )
  }

  return {
    getTracks,
    getTrack,
    adminGetTracks,
    adminCreateTrack,
    adminUploadTrackAudio,
    adminUploadTrackCover,
    adminUpdateTrack,
    adminDeleteTrack,
  }
}
