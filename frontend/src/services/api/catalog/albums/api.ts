import { useRequest } from '@/composables'
import type { UseRequestConfig } from '@/plugins/client/types'
import { AlbumApiRoutes } from './enums'
import { AlbumSchema, type Album, type AlbumCreatePayload, type AlbumUpdatePayload } from './types'

export const useAlbumsApi = () => {
  // ==================== PUBLIC ====================

  const getAlbums = async (params?: Record<string, string | number>, config?: UseRequestConfig<Album[]>) => {
    return useRequest<Album, true>(
      AlbumApiRoutes.LIST,
      { method: 'GET', params },
      {
        schema: AlbumSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getAlbum = async (id: string | number, config?: UseRequestConfig<Album>) => {
    return useRequest<Album>(
      AlbumApiRoutes.GET.replace(':albumId', String(id)),
      { method: 'GET' },
      {
        schema: AlbumSchema,
        silent: true,
        ...config,
      },
    )
  }

  // ==================== ADMIN ====================

  const adminGetAlbums = async (config?: UseRequestConfig<Album[]>) => {
    return useRequest<Album, true>(
      AlbumApiRoutes.ADMIN_LIST,
      { method: 'GET' },
      {
        schema: AlbumSchema,
        silent: true,
        ...config,
      },
    )
  }

  const adminCreateAlbum = async (
    payload: AlbumCreatePayload,
    config?: UseRequestConfig<Album>,
  ) => {
    return useRequest<Album>(
      AlbumApiRoutes.ADMIN_CREATE,
      {
        method: 'POST',
        data: {
          title: payload.title,
          coverUrl: payload.coverUrl ?? null,
          artistId: payload.artistId ?? null,
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
    payload: AlbumUpdatePayload,
    config?: UseRequestConfig<Album>,
  ) => {
    return useRequest<Album>(
      AlbumApiRoutes.ADMIN_UPDATE.replace(':albumId', String(id)),
      {
        method: 'PATCH',
        data: {
          title: payload.title,
          coverUrl: payload.coverUrl ?? null,
          artistId: payload.artistId ?? null,
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
    return useRequest<void>(
      AlbumApiRoutes.ADMIN_DELETE.replace(':albumId', String(id)),
      { method: 'DELETE' },
      {
        silent: false,
        ...config,
      },
    )
  }

  const adminEnrichAlbum = async (
    id: string | number,
    config?: UseRequestConfig<{ data: Album }>,
  ) => {
    return useRequest<{ data: Album }>(
      AlbumApiRoutes.ADMIN_ENRICH.replace(':albumId', String(id)),
      { method: 'POST' },
      {
        silent: false,
        ...config,
      },
    )
  }

  return {
    // Public
    getAlbums,
    getAlbum,
    // Admin
    adminGetAlbums,
    adminCreateAlbum,
    adminUpdateAlbum,
    adminDeleteAlbum,
    adminEnrichAlbum,
  }
}
