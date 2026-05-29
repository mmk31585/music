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
  artist_id?: string | number | null
  album_id?: string | number | null
  title?: string
  duration_seconds?: number | null
  audio_url?: string | null
  cover_url?: string | null
  genre_ids?: Array<string | number>
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
    artist_id: payload.artist_id ?? null,
    album_id: payload.album_id ?? null,
    title: payload.title,
    duration_seconds: payload.duration_seconds ?? 0,
    audio_url: payload.audio_url ?? null,
    cover_url: payload.cover_url ?? null,
    genre_ids: payload.genre_id ? [payload.genre_id] : [],
  }
}

export const useCatalogApi = () => {
  // ---------- PUBLIC ----------

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

  const getTrack = async (
    id: string | number,
    config?: UseRequestConfig<Track>,
  ) => {
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

  const searchCatalog = async (
    query: string,
    config?: UseRequestConfig<unknown>,
  ) => {
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

  // ---------- ADMIN TRACKS ----------

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

  const adminCreateTrack = async (
    payload: Partial<Track>,
    config?: UseRequestConfig<Track>,
  ) => {
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
    if (data.artist_id != null) formData.append('artistID', String(data.artist_id))
    if (data.album_id != null) formData.append('albumID', String(data.album_id))
    if (data.duration_seconds != null) {
      formData.append('durationSeconds', String(data.duration_seconds))
    }
    if (data.audio_url) formData.append('audioUrl', data.audio_url)
    if (data.cover_url) formData.append('coverUrl', data.cover_url)

    for (const genreID of data.genre_ids ?? []) {
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

  const adminDeleteTrack = async (
    id: string | number,
    config?: UseRequestConfig<void>,
  ) => {
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

  // ---------- ADMIN ARTISTS ----------

  const adminGetArtists = async (config?: UseRequestConfig<Artist[]>) => {
    return useRequest<Artist, true>(
      CatalogApiRoutes.ADMIN_ARTISTS,
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

  const adminCreateArtist = async (
    payload: Partial<Artist>,
    config?: UseRequestConfig<Artist>,
  ) => {
    return useRequest<Artist>(
      CatalogApiRoutes.ADMIN_ARTISTS,
      {
        method: 'POST',
        data: {
          name: payload.name,
          bio: payload.bio ?? null,
          image_url: payload.image_url ?? null,
        },
      },
      {
        schema: ArtistSchema,
        silent: false,
        ...config,
      },
    )
  }

  const adminUpdateArtist = async (
    id: string | number,
    payload: Partial<Artist>,
    config?: UseRequestConfig<Artist>,
  ) => {
    return useRequest<Artist>(
      `${CatalogApiRoutes.ADMIN_ARTISTS}/${id}`,
      {
        method: 'PATCH',
        data: {
          name: payload.name,
          bio: payload.bio ?? null,
          image_url: payload.image_url ?? null,
        },
      },
      {
        schema: ArtistSchema,
        silent: false,
        ...config,
      },
    )
  }

  const adminDeleteArtist = async (
    id: string | number,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest(
      `${CatalogApiRoutes.ADMIN_ARTISTS}/${id}`,
      {
        method: 'DELETE',
      },
      {
        silent: false,
        ...config,
      },
    )
  }

  // ---------- ADMIN ALBUMS ----------

  const adminGetAlbums = async (config?: UseRequestConfig<Album[]>) => {
    return useRequest<Album, true>(
      CatalogApiRoutes.ADMIN_ALBUMS,
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

  const adminCreateAlbum = async (
    payload: Partial<Album>,
    config?: UseRequestConfig<Album>,
  ) => {
    return useRequest<Album>(
      CatalogApiRoutes.ADMIN_ALBUMS,
      {
        method: 'POST',
        data: {
          title: payload.title,
          cover_url: payload.cover_url ?? null,
          artist_id: payload.artist_id ?? null,
        },
      },
      {
        schema: AlbumSchema,
        silent: false,
        ...config,
      },
    )
  }

  const adminUpdateAlbum = async (
    id: string | number,
    payload: Partial<Album>,
    config?: UseRequestConfig<Album>,
  ) => {
    return useRequest<Album>(
      `${CatalogApiRoutes.ADMIN_ALBUMS}/${id}`,
      {
        method: 'PATCH',
        data: {
          title: payload.title,
          coverUrl: payload.cover_url ?? null,
          artistID: payload.artist_id ?? null,
        },
      },
      {
        schema: AlbumSchema,
        silent: false,
        ...config,
      },
    )
  }

  const adminDeleteAlbum = async (
    id: string | number,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest(
      `${CatalogApiRoutes.ADMIN_ALBUMS}/${id}`,
      {
        method: 'DELETE',
      },
      {
        silent: false,
        ...config,
      },
    )
  }

  // ---------- ADMIN GENRES ----------

  const adminGetGenres = async (config?: UseRequestConfig<Genre[]>) => {
    return useRequest<Genre, true>(
      CatalogApiRoutes.ADMIN_GENRES,
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

  const adminCreateGenre = async (
    payload: Partial<Genre>,
    config?: UseRequestConfig<Genre>,
  ) => {
    return useRequest<Genre>(
      CatalogApiRoutes.ADMIN_GENRES,
      {
        method: 'POST',
        data: {
          name: payload.name,
        },
      },
      {
        schema: GenreSchema,
        silent: false,
        ...config,
      },
    )
  }

  const adminUpdateGenre = async (
    id: string | number,
    payload: Partial<Genre>,
    config?: UseRequestConfig<Genre>,
  ) => {
    return useRequest<Genre>(
      `${CatalogApiRoutes.ADMIN_GENRES}/${id}`,
      {
        method: 'PATCH',
        data: {
          name: payload.name,
        },
      },
      {
        schema: GenreSchema,
        silent: false,
        ...config,
      },
    )
  }

  const adminDeleteGenre = async (
    id: string | number,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest(
      `${CatalogApiRoutes.ADMIN_GENRES}/${id}`,
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

    adminGetArtists,
    adminCreateArtist,
    adminUpdateArtist,
    adminDeleteArtist,

    adminGetAlbums,
    adminCreateAlbum,
    adminUpdateAlbum,
    adminDeleteAlbum,

    adminGetGenres,
    adminCreateGenre,
    adminUpdateGenre,
    adminDeleteGenre,
  }
}
