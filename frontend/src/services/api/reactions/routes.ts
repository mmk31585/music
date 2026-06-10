import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { ReactionsApiRoutes } from './enums'
import {
  ReactionSchema,
  CountsResponseSchema,
  type Reaction,
  type CountsResponse,
  type ReactRequest,
} from './types'

export const useReactionsApi = () => {
  const react = async (payload: ReactRequest, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      ReactionsApiRoutes.REACT,
      { method: 'POST', data: payload },
      { silent: false, ...config },
    )
  }

  const removeReaction = async (
    targetType: string,
    targetId: string,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      ReactionsApiRoutes.REMOVE.replace(':targetType', targetType).replace(':targetId', targetId),
      { method: 'DELETE' },
      { silent: false, ...config },
    )
  }

  const getUserReaction = async (
    targetType: string,
    targetId: string,
    config?: UseRequestConfig<{ reaction: Reaction | null }>,
  ) => {
    return useRequest<{ reaction: Reaction | null }>(
      ReactionsApiRoutes.USER_REACTION.replace(':targetType', targetType).replace(
        ':targetId',
        targetId,
      ),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const getCounts = async (
    targetType: string,
    targetId: string,
    config?: UseRequestConfig<CountsResponse>,
  ) => {
    return useRequest<CountsResponse>(
      ReactionsApiRoutes.COUNTS.replace(':targetType', targetType).replace(':targetId', targetId),
      { method: 'GET' },
      { schema: CountsResponseSchema, silent: true, ...config },
    )
  }

  const getLikedTracks = async (
    params?: { user_id?: string; limit?: number; offset?: number },
    config?: UseRequestConfig<{ items: unknown[] }>,
  ) => {
    return useRequest<{ items: unknown[] }>(
      ReactionsApiRoutes.LIKED_TRACKS,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const getLikedAlbums = async (
    params?: { user_id?: string; limit?: number; offset?: number },
    config?: UseRequestConfig<{ items: unknown[] }>,
  ) => {
    return useRequest<{ items: unknown[] }>(
      ReactionsApiRoutes.LIKED_ALBUMS,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  return { react, removeReaction, getUserReaction, getCounts, getLikedTracks, getLikedAlbums }
}
