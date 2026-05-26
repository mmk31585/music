import { useRequest } from '@/composables'
import type { AxiosRequestConfig } from 'axios'
import type { UseRequestConfig } from '@/plugins/client/types'
import { CatalogApiRoutes } from './enums'
import { UploadResponseSchema } from '@/services/api/media'
import {
  ArtistSchema,
  AlbumSchema,
  GenreSchema,
  TrackSchema,
  type Artist,
  type Album,
  type Genre,
  type Track,
} from './types'
import { z } from 'zod'

type TrackMutationPayload = {
  artistId?: string | number | null
  albumId?: string | number | null
  title?: string
  durationSeconds?: number | null
  audioUrl?: string | null
  coverUrl?: string | null
  genreIds?: Array<string | number>
}

export type TrackFormPayload = Partial<Track> & {
  audioFile?: File | null
}

const UploadedTrackResponseSchema = z.object({
  track: TrackSchema,
  media: UploadResponseSchema,
})

function toTrackMutationPayload(payload: Partial<Track>): TrackMutationPayload {
  return {
    artistId: payload.artist_id ?? null,
    albumId: payload.album_id ?? null,
    title: payload.title,
    durationSeconds: payload.duration_seconds ?? 0,
    audioUrl: payload.audio_url ?? null,
    coverUrl: payload.cover_url ?? null,
    genreIds: payload.genre_id ? [payload.genre_id] : [],
  }
}

export const useCatalogApi = () => {
  const getTracks = async (config?: UseRequestConfig<Track[]>) => {
    return useRequest<Track, true>(
      CatalogApiRoutes.TRACKS,
      {
        method: 'GET',
      },
      {
        schema: TrackSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getTrack = async (id: string | number, config?: UseRequestConfig<Track>) => {
    return useRequest<Track>(
      `${CatalogApiRoutes.TRACKS}/${id}`,
      {
        method: 'GET',
      },
      {
        schema: TrackSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getArtists = async (config?: UseRequestConfig<Artist[]>) => {
    return useRequest<Artist, true>(
      CatalogApiRoutes.ARTISTS,
      {
        method: 'GET',
      },
      {
        schema: ArtistSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getAlbums = async (config?: UseRequestConfig<Album[]>) => {
    return useRequest<Album, true>(
      CatalogApiRoutes.ALBUMS,
      {
        method: 'GET',
      },
      {
        schema: AlbumSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getGenres = async (config?: UseRequestConfig<Genre[]>) => {
    return useRequest<Genre, true>(
      CatalogApiRoutes.GENRES,
      {
        method: 'GET',
      },
      {
        schema: GenreSchema,
        silent: true,
        ...config,
      },
    )
  }

  const searchCatalog = async (query: string, config?: UseRequestConfig<unknown>) => {
    return useRequest(
      `${CatalogApiRoutes.SEARCH}?q=${encodeURIComponent(query)}`,
      {
        method: 'GET',
      },
      {
        silent: true,
        ...config,
      },
    )
  }

  // ---------- ADMIN SECTION ----------

  const adminGetTracks = async (config?: UseRequestConfig<Track[]>) => {
    return useRequest<Track, true>(
      CatalogApiRoutes.ADMIN_TRACKS,
      {
        method: 'GET',
      },
      {
        schema: TrackSchema,
        silent: true,
        ...config,
      },
    )
  }

  const adminCreateTrack = async (payload: Partial<Track>, config?: UseRequestConfig<Track>) => {
    return useRequest<Track>(
      CatalogApiRoutes.ADMIN_TRACKS,
      {
        method: 'POST',
        data: toTrackMutationPayload(payload),
      },
      {
        schema: TrackSchema,
        silent: false,
        ...config,
      },
    )
  }

  const adminUploadTrackWithAudio = async (
    payload: TrackFormPayload,
    file: File,
    config?: UseRequestConfig<Track>,
    requestConfig?: Pick<AxiosRequestConfig, 'onUploadProgress' | 'signal'>,
  ) => {
    const data = toTrackMutationPayload(payload)
    const formData = new FormData()

    formData.append('trackAudio', file)
    if (data.title) formData.append('title', data.title)
    if (data.artistId) formData.append('artistId', String(data.artistId))
    if (data.albumId) formData.append('albumId', String(data.albumId))
    if (data.durationSeconds != null) {
      formData.append('durationSeconds', String(data.durationSeconds))
    }
    if (data.audioUrl) formData.append('audioUrl', data.audioUrl)
    if (data.coverUrl) formData.append('coverUrl', data.coverUrl)
    for (const genreID of data.genreIds ?? []) {
      formData.append('genreIds', String(genreID))
    }

    const response = await useRequest<z.infer<typeof UploadedTrackResponseSchema>>(
      CatalogApiRoutes.ADMIN_TRACKS_UPLOAD,
      {
        method: 'POST',
        data: formData,
        ...requestConfig,
      },
      {
        schema: UploadedTrackResponseSchema,
        silent: false,
      },
    )

    config?.success?.(response.track, 1, response)

    return response.track
  }

  const adminUpdateTrack = async (
    id: string | number,
    payload: Partial<Track>,
    config?: UseRequestConfig<Track>,
  ) => {
    return useRequest<Track>(
      `${CatalogApiRoutes.ADMIN_TRACKS}/${id}`,
      {
        method: 'PATCH',
        data: toTrackMutationPayload(payload),
      },
      {
        schema: TrackSchema,
        silent: false,
        ...config,
      },
    )
  }

  const adminDeleteTrack = async (id: string | number, config?: UseRequestConfig<void>) => {
    return useRequest(
      `${CatalogApiRoutes.ADMIN_TRACKS}/${id}`,
      {
        method: 'DELETE',
      },
      {
        silent: false,
        ...config,
      },
    )
  }

  return {
    getTracks,
    getTrack,
    getArtists,
    getAlbums,
    getGenres,
    searchCatalog,

    adminGetTracks,
    adminCreateTrack,
    adminUploadTrackWithAudio,
    adminUpdateTrack,
    adminDeleteTrack,
  }
}
