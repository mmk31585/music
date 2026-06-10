import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { HistoryApiRoutes } from './enums'
import { HistoryResponseSchema, type HistoryResponse } from './types'

export const useHistoryApi = () => {
  const getHistory = async (
    params?: { limit?: number; offset?: number },
    config?: UseRequestConfig<HistoryResponse>,
  ) => {
    return useRequest<HistoryResponse>(
      HistoryApiRoutes.LIST,
      {
        method: 'GET',
        params,
      },
      {
        schema: HistoryResponseSchema,
        silent: true,
        ...config,
      },
    )
  }

  const addHistory = async (payload: {
    track_id: string
    duration: number
    completed: boolean
  }) => {
    return useRequest<void>(HistoryApiRoutes.LIST, {
      method: 'POST',
      data: payload,
    })
  }

  return {
    getHistory,
    addHistory,
  }
}
