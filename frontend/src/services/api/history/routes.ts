import { useRequest } from '@/composables/useRequest'
import { apiReplaceParams } from '@/utils/api-replace-params'
import type { UseRequestConfig } from '@/plugins/client/types'
import { HistoryApiRoutes } from './enums'
import { HistoryResponseSchema, type HistoryResponse } from './types'

export const useHistoryApi = () => {
  /**
   * Get paginated listening history.
   * Backend: GET /history → HistoryResponse { items, pagination }
   */
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

  /**
   * Delete a single history entry.
   * Backend: DELETE /history/:id → 204
   */
  const deleteHistoryItem = async (id: string) => {
    return useRequest<void>(
      apiReplaceParams(HistoryApiRoutes.DELETE_ITEM, { id }),
      { method: 'DELETE' },
    )
  }

  /**
   * Clear all listening history for the current user.
   * Backend: DELETE /history → 204
   */
  const clearAllHistory = async () => {
    return useRequest<void>(
      HistoryApiRoutes.CLEAR_ALL,
      { method: 'DELETE' },
    )
  }

  return {
    getHistory,
    deleteHistoryItem,
    clearAllHistory,
  }
}
