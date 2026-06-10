import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { RecommendationApiRoutes } from './enums'
import {
  RecommendationResponseSchema,
  RecommendationTrackSchema,
  type RecommendationTrack,
  type RecommendationResponse,
} from './types'

export const useRecommendationsApi = () => {
  const getPopular = async (
    params?: { limit?: number },
    config?: UseRequestConfig<RecommendationResponse>,
  ) => {
    return useRequest<RecommendationResponse>(
      RecommendationApiRoutes.POPULAR,
      {
        method: 'GET',
        params,
      },
      {
        schema: RecommendationResponseSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getBest = async (
    params?: { limit?: number },
    config?: UseRequestConfig<RecommendationResponse>,
  ) => {
    return useRequest<RecommendationResponse>(
      RecommendationApiRoutes.BEST,
      {
        method: 'GET',
        params,
      },
      {
        schema: RecommendationResponseSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getRecent = async (
    params?: { limit?: number },
    config?: UseRequestConfig<RecommendationResponse>,
  ) => {
    return useRequest<RecommendationResponse>(
      RecommendationApiRoutes.RECENT,
      {
        method: 'GET',
        params,
      },
      {
        schema: RecommendationResponseSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getForYou = async (
    params?: { limit?: number },
    config?: UseRequestConfig<RecommendationResponse>,
  ) => {
    return useRequest<RecommendationResponse>(
      RecommendationApiRoutes.FOR_YOU,
      {
        method: 'GET',
        params,
      },
      {
        schema: RecommendationResponseSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getSimilar = async (
    trackId: string,
    params?: { limit?: number },
    config?: UseRequestConfig<RecommendationResponse>,
  ) => {
    return useRequest<RecommendationResponse>(
      RecommendationApiRoutes.SIMILAR.replace(':trackId', trackId),
      {
        method: 'GET',
        params,
      },
      {
        schema: RecommendationResponseSchema,
        silent: true,
        ...config,
      },
    )
  }

  return {
    getPopular,
    getBest,
    getRecent,
    getForYou,
    getSimilar,
  }
}
