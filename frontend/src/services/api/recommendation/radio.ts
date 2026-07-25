import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import type { PlaybackTrack } from '@/services/api/player'

export interface RadioSession {
  session_id: string
  seed_track_id: string
}

export interface StartRadioResponse {
  session_id: string
  tracks: PlaybackTrack[]
}

export interface NextBatchResponse {
  tracks: PlaybackTrack[]
  has_more: boolean
}

export const useRadioApi = () => {
  const startRadio = async (
    seedTrackId: string,
    seedType?: string,
    config?: UseRequestConfig<StartRadioResponse>,
  ): Promise<StartRadioResponse> => {
    return useRequest<StartRadioResponse>(
      '/radio/start',
      {
        method: 'POST',
        data: { seed_track_id: seedTrackId, seed_type: seedType || 'track' },
      },
      { silent: true, ...config },
    )
  }

  const getNextRadioBatch = async (
    sessionId: string,
    count?: number,
    config?: UseRequestConfig<NextBatchResponse>,
  ): Promise<NextBatchResponse> => {
    const params: Record<string, string | number> = {}
    if (count) params.count = count
    return useRequest<NextBatchResponse>(
      `/radio/${sessionId}/next`,
      {
        method: 'GET',
        params,
      },
      { ...config },
    )
  }

  const endRadio = async (
    sessionId: string,
    config?: UseRequestConfig<any>,
  ): Promise<void> => {
    await useRequest(
      `/radio/${sessionId}/end`,
      { method: 'POST' },
      { silent: true, ...config },
    )
  }

  return {
    startRadio,
    getNextRadioBatch,
    endRadio,
  }
}
