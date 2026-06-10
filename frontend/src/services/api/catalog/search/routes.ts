import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { SearchApiRoutes } from './enums'
import { SearchResponseSchema, type SearchParams, type SearchResponse } from './types'

export const useSearchApi = () => {
  const search = async (params: SearchParams, config?: UseRequestConfig<SearchResponse>) => {
    return useRequest<SearchResponse>(
      SearchApiRoutes.SEARCH,
      {
        method: 'GET',
        params: {
          q: params.q,
          type: params.type ?? 'all',
          limit: params.limit ?? 20,
        },
      },
      {
        schema: SearchResponseSchema,
        silent: true,
        ...config,
      },
    )
  }

  const suggestions = async (
    q: string,
    limit?: number,
    config?: UseRequestConfig<{ suggestions: string[]; query: string }>,
  ) => {
    return useRequest<{ suggestions: string[]; query: string }>(
      SearchApiRoutes.SUGGESTIONS,
      {
        method: 'GET',
        params: { q, limit: limit ?? 5 },
      },
      { silent: true, ...config },
    )
  }

  return {
    search,
    suggestions,
  }
}
