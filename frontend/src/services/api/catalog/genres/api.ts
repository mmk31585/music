import { useRequest } from '@/composables'
import type { UseRequestConfig } from '@/plugins/client/types'
import { GenreApiRoutes } from './enums'
import { GenreSchema, type Genre, type GenreCreatePayload, type GenreUpdatePayload } from './types'

export const useGenresApi = () => {
  // ==================== PUBLIC ====================

  const getGenres = async (config?: UseRequestConfig<Genre[]>) => {
    return useRequest<Genre, true>(
      GenreApiRoutes.LIST,
      { method: 'GET' },
      {
        schema: GenreSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getGenre = async (id: string | number, config?: UseRequestConfig<Genre>) => {
    return useRequest<Genre>(
      GenreApiRoutes.GET.replace(':genreId', String(id)),
      { method: 'GET' },
      {
        schema: GenreSchema,
        silent: true,
        ...config,
      },
    )
  }

  // ==================== ADMIN ====================

  const adminGetGenres = async (config?: UseRequestConfig<Genre[]>) => {
    return useRequest<Genre, true>(
      GenreApiRoutes.ADMIN_LIST,
      { method: 'GET' },
      {
        schema: GenreSchema,
        silent: true,
        ...config,
      },
    )
  }

  const adminCreateGenre = async (
    payload: GenreCreatePayload,
    config?: UseRequestConfig<Genre>,
  ) => {
    return useRequest<Genre>(
      GenreApiRoutes.ADMIN_CREATE,
      {
        method: 'POST',
        data: { name: payload.name },
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
    payload: GenreUpdatePayload,
    config?: UseRequestConfig<Genre>,
  ) => {
    return useRequest<Genre>(
      GenreApiRoutes.ADMIN_UPDATE.replace(':genreId', String(id)),
      {
        method: 'PATCH',
        data: { name: payload.name },
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
    return useRequest<void>(
      GenreApiRoutes.ADMIN_DELETE.replace(':genreId', String(id)),
      { method: 'DELETE' },
      {
        silent: false,
        ...config,
      },
    )
  }

  return {
    // Public
    getGenres,
    getGenre,
    // Admin
    adminGetGenres,
    adminCreateGenre,
    adminUpdateGenre,
    adminDeleteGenre,
  }
}
