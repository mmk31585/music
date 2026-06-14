import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { SearchApiRoutes } from './enums'
import { SearchResultSchema, type SearchParams, type SearchResult } from './types'

export const useSearchApi = () => {
  const search = async (params: SearchParams, config?: UseRequestConfig<SearchResult>) => {
    return useRequest<SearchResult>(
      SearchApiRoutes.SEARCH,
      {
        method: 'GET',
        params: {
          q: params.query,
          type: params.type ?? 'all',
          limit: params.limit ?? 20,
        },
      },
      {
        schema: SearchResultSchema,
        silent: true,
        ...config,
      },
    )
  }

  return {
    search,
  }
}
