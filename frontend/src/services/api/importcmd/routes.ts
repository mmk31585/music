import { useRequest } from '@/composables/useRequest'
import type { SearchResult, ImportResponse, ProgressResponse } from './types'

const BASE = '/admin/import'

export const useImportApi = () => {
  const search = async (q: string) => {
    return useRequest<SearchResult, true>(
      `${BASE}/search?q=${encodeURIComponent(q)}`,
      { method: 'GET' },
      { silent: true },
    )
  }

  const importTrack = async (url: string, source?: string) => {
    return useRequest<ImportResponse>(
      `${BASE}/import`,
      { method: 'POST', data: source ? { url, source } : { url } },
      { silent: false },
    )
  }

  const getProgress = async (jobId: string) => {
    return useRequest<ProgressResponse>(
      `${BASE}/${jobId}/progress`,
      { method: 'GET' },
      { silent: true },
    )
  }

  return { search, importTrack, getProgress }
}
