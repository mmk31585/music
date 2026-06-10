import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { ModerationApiRoutes } from './enums'
import {
  ContentReportSchema,
  type ContentReport,
  type ContentFlag,
  type ModerationAction,
  type ModerationStats,
} from './types'

export const useModerationApi = () => {
  const reportContent = async (
    payload: {
      target_id: string
      target_type: string
      reason: string
      description?: string | null
    },
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      ModerationApiRoutes.REPORT,
      { method: 'POST', data: payload },
      { silent: true, ...config },
    )
  }

  const getPendingReports = async (
    params?: { limit?: number; offset?: number },
    config?: UseRequestConfig<{ items: ContentReport[]; limit: number; offset: number }>,
  ) => {
    return useRequest<{ items: ContentReport[]; limit: number; offset: number }>(
      ModerationApiRoutes.PENDING,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const getReportsByStatus = async (
    status: string,
    params?: { limit?: number; offset?: number },
    config?: UseRequestConfig<{ items: ContentReport[]; limit: number; offset: number }>,
  ) => {
    return useRequest<{ items: ContentReport[]; limit: number; offset: number }>(
      ModerationApiRoutes.BY_STATUS.replace(':status', status),
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const getReportDetail = async (
    id: string,
    config?: UseRequestConfig<{ report: ContentReport }>,
  ) => {
    return useRequest<{ report: ContentReport }>(
      ModerationApiRoutes.REPORT_DETAIL.replace(':id', id),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const resolveReport = async (
    id: string,
    status?: string,
    note?: string,
    config?: UseRequestConfig<void>,
  ) => {
    let query = ''
    const params = new URLSearchParams()
    if (status) params.set('status', status)
    if (note) params.set('note', note)
    const qs = params.toString()
    if (qs) query = '?' + qs
    return useRequest<void>(
      ModerationApiRoutes.RESOLVE.replace(':id', id) + query,
      { method: 'POST' },
      { silent: false, ...config },
    )
  }

  const bulkAction = async (
    ids: string[],
    action: string,
    note?: string,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      ModerationApiRoutes.BULK,
      { method: 'POST', data: { ids, action, note } },
      { silent: false, ...config },
    )
  }

  const flagContent = async (
    payload: {
      target_id: string
      target_type: string
      flag_type: string
      expires_in_hours?: number
    },
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      ModerationApiRoutes.FLAG,
      { method: 'POST', data: payload },
      { silent: true, ...config },
    )
  }

  const getFlags = async (
    params?: { limit?: number; offset?: number; expired?: boolean },
    config?: UseRequestConfig<{ items: ContentFlag[]; limit: number; offset: number }>,
  ) => {
    return useRequest<{ items: ContentFlag[]; limit: number; offset: number }>(
      ModerationApiRoutes.FLAGS,
      { method: 'GET', params: { ...params, expired: params?.expired ? 'true' : undefined } },
      { silent: true, ...config },
    )
  }

  const getStats = async (config?: UseRequestConfig<{ stats: ModerationStats }>) => {
    return useRequest<{ stats: ModerationStats }>(
      ModerationApiRoutes.STATS,
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const getActions = async (
    params?: { report_id?: string; limit?: number; offset?: number },
    config?: UseRequestConfig<{ items: ModerationAction[]; limit: number; offset: number }>,
  ) => {
    return useRequest<{ items: ModerationAction[]; limit: number; offset: number }>(
      ModerationApiRoutes.ACTIONS,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  return {
    reportContent,
    getPendingReports,
    getReportsByStatus,
    getReportDetail,
    resolveReport,
    bulkAction,
    flagContent,
    getFlags,
    getStats,
    getActions,
  }
}
