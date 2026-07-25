import { useRequest } from '@/composables'
import type { UseRequestConfig } from '@/plugins/client/types'
import { SearchApiRoutes } from './enums'
import { SearchResultSchema, type SearchResult, type SearchParams } from './types'

export const useSearchApi = () => {
  const searchCatalog = async (
    params: SearchParams,
    config?: UseRequestConfig<SearchResult>,
  ) => {
    const queryParams = new URLSearchParams()
    queryParams.set('q', params.query)
    
    if (params.type && params.type !== 'all') {
      queryParams.set('type', params.type)
    }
    if (params.limit) {
      queryParams.set('limit', String(params.limit))
    }

    return useRequest<SearchResult>(
      `${SearchApiRoutes.SEARCH}?${queryParams.toString()}`,
      { method: 'GET' },
      {
        schema: SearchResultSchema,
        silent: true,
        ...config,
      },
    )
  }

  return {
    searchCatalog,
  }
}
