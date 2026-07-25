import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { QueueApiRoutes } from './enums'
import {
  QueueResponseSchema,
  AddTrackResponseSchema,
  QueueItemSchema,
  type QueueResponse,
  type AddTrackResponse,
  type AddTrackPayload,
  type ReorderQueuePayload,
  type QueueItem,
} from './types'

export const useQueueApi = () => {
  const getQueue = async (config?: UseRequestConfig<QueueResponse>) => {
    return useRequest<QueueResponse>(
      QueueApiRoutes.GET,
      { method: 'GET' },
      { schema: QueueResponseSchema, silent: true, ...config },
    )
  }

  const addTrack = async (
    payload: AddTrackPayload,
    config?: UseRequestConfig<AddTrackResponse>,
  ) => {
    return useRequest<AddTrackResponse>(
      QueueApiRoutes.ADD_TRACK,
      { method: 'POST', data: payload },
      { schema: AddTrackResponseSchema, silent: false, ...config },
    )
  }

  const removeTrack = async (id: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      QueueApiRoutes.REMOVE_TRACK.replace(':id', id),
      { method: 'DELETE' },
      { silent: false, ...config },
    )
  }

  const reorder = async (payload: ReorderQueuePayload, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      QueueApiRoutes.REORDER,
      { method: 'PUT', data: payload },
      { silent: false, ...config },
    )
  }

  const clear = async (config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      QueueApiRoutes.CLEAR,
      { method: 'DELETE' },
      { silent: false, ...config },
    )
  }

  return {
    getQueue,
    addTrack,
    removeTrack,
    reorder,
    clear,
  }
}
