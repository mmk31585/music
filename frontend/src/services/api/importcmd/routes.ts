import { useRequest } from '@/composables/useRequest'
import type { SearchResult, ImportResponse } from './types'

const BASE = '/admin/import'

export const useImportApi = () => {
  const search = async (q: string) => {
    return useRequest<SearchResult[], true>(
      `${BASE}/search?q=${encodeURIComponent(q)}`,
      { method: 'GET' },
      { silent: true },
    )
  }

  const importTrack = async (url: string) => {
    return useRequest<ImportResponse>(
      `${BASE}/import`,
      { method: 'POST', data: { url } },
      { silent: false },
    )
  }

  return { search, importTrack }
}
