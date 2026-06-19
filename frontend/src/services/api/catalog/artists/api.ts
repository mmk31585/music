import { useRequest } from '@/composables'
import type { UseRequestConfig } from '@/plugins/client/types'
import { ArtistApiRoutes } from './enums'
import { ArtistSchema, type Artist, type ArtistCreatePayload, type ArtistUpdatePayload } from './types'

export const useArtistsApi = () => {
  // ==================== PUBLIC ====================

  const getArtists = async (config?: UseRequestConfig<Artist[]>) => {
    return useRequest<Artist, true>(
      ArtistApiRoutes.LIST,
      { method: 'GET' },
      {
        schema: ArtistSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getArtist = async (id: string | number, config?: UseRequestConfig<Artist>) => {
    return useRequest<Artist>(
      ArtistApiRoutes.GET.replace(':artistId', String(id)),
      { method: 'GET' },
      {
        schema: ArtistSchema,
        silent: true,
        ...config,
      },
    )
  }

  // ==================== ADMIN ====================

  const adminGetArtists = async (config?: UseRequestConfig<Artist[]>) => {
    return useRequest<Artist, true>(
      ArtistApiRoutes.ADMIN_LIST,
      { method: 'GET' },
      {
        schema: ArtistSchema,
        silent: true,
        ...config,
      },
    )
  }

  const adminCreateArtist = async (
    payload: ArtistCreatePayload,
    config?: UseRequestConfig<Artist>,
  ) => {
    return useRequest<Artist>(
      ArtistApiRoutes.ADMIN_CREATE,
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
    payload: ArtistUpdatePayload,
    config?: UseRequestConfig<Artist>,
  ) => {
    return useRequest<Artist>(
      ArtistApiRoutes.ADMIN_UPDATE.replace(':artistId', String(id)),
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
    return useRequest<void>(
      ArtistApiRoutes.ADMIN_DELETE.replace(':artistId', String(id)),
      { method: 'DELETE' },
      {
        silent: false,
        ...config,
      },
    )
  }

  const adminEnrichArtist = async (
    id: string | number,
    config?: UseRequestConfig<{ data: Artist }>,
  ) => {
    return useRequest<{ data: Artist }>(
      ArtistApiRoutes.ADMIN_ENRICH.replace(':artistId', String(id)),
      { method: 'POST' },
      {
        silent: false,
        ...config,
      },
    )
  }

  return {
    // Public
    getArtists,
    getArtist,
    // Admin
    adminGetArtists,
    adminCreateArtist,
    adminUpdateArtist,
    adminDeleteArtist,
    adminEnrichArtist,
  }
}
