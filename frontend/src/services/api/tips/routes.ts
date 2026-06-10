import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { TipsApiRoutes } from './enums'
import type { Tip, CreateTipRequest } from './types'

export const useTipsApi = () => {
  const createTip = async (
    payload: CreateTipRequest,
    config?: UseRequestConfig<{ tip: Tip; redirect_url: string }>,
  ) => {
    return useRequest<{ tip: Tip; redirect_url: string }>(
      TipsApiRoutes.CREATE,
      { method: 'POST', data: payload },
      { silent: false, ...config },
    )
  }

  const listSent = async (config?: UseRequestConfig<{ tips: Tip[] }>) => {
    return useRequest<{ tips: Tip[] }>(
      TipsApiRoutes.SENT,
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const listReceived = async (
    artistId: string,
    config?: UseRequestConfig<{ tips: Tip[]; total_cents: number }>,
  ) => {
    return useRequest<{ tips: Tip[]; total_cents: number }>(
      TipsApiRoutes.RECEIVED.replace(':artist_id', artistId),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  return { createTip, listSent, listReceived }
}
