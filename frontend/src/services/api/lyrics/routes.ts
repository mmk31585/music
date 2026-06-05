import { useRequest } from '@/composables'
import type { UseRequestConfig } from '@/plugins/client/types'
import { LyricsApiRoutes } from './enums'
import {
  LyricsSchema,
  type Lyrics,
  type CreateLyricsPayload,
  type UpdateLyricsPayload,
} from './types'

export const useLyricsApi = () => {
  const getTrackLyrics = async (
    trackId: string | number,
    lang?: string,
    config?: UseRequestConfig<Lyrics>,
  ) => {
    const query = lang ? `?lang=${encodeURIComponent(lang)}` : ''

    return useRequest<Lyrics>(
      `${LyricsApiRoutes.GET_TRACK_LYRICS.replace(':trackId', String(trackId))}${query}`,
      { method: 'GET' },
      {
        schema: LyricsSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getLyricsByTrackID = async (
    trackId: string | number,
    config?: UseRequestConfig<Lyrics[]>,
  ) => {
    return useRequest<Lyrics, true>(
      LyricsApiRoutes.GET_TRACK_LYRICS_ALL.replace(':trackId', String(trackId)),
      { method: 'GET' },
      {
        schema: LyricsSchema,
        silent: true,
        ...config,
      },
    )
  }

  const adminCreateLyrics = async (
    payload: CreateLyricsPayload,
    config?: UseRequestConfig<Lyrics>,
  ) => {
    return useRequest<Lyrics>(
      LyricsApiRoutes.ADMIN_CREATE,
      {
        method: 'POST',
        data: payload,
      },
      {
        schema: LyricsSchema,
        silent: false,
        ...config,
      },
    )
  }

  const adminUpdateLyrics = async (
    id: string | number,
    payload: UpdateLyricsPayload,
    config?: UseRequestConfig<Lyrics>,
  ) => {
    return useRequest<Lyrics>(
      LyricsApiRoutes.ADMIN_UPDATE.replace(':id', String(id)),
      {
        method: 'PUT',
        data: payload,
      },
      {
        schema: LyricsSchema,
        silent: false,
        ...config,
      },
    )
  }

  const adminDeleteLyrics = async (id: string | number, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      LyricsApiRoutes.ADMIN_DELETE.replace(':id', String(id)),
      { method: 'DELETE' },
      {
        silent: false,
        ...config,
      },
    )
  }

  return {
    getTrackLyrics,
    getLyricsByTrackID,
    adminCreateLyrics,
    adminUpdateLyrics,
    adminDeleteLyrics,
  }
}
