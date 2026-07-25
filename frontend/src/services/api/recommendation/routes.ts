import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { RecommendationApiRoutes } from './enums'
import {
  RecommendationResponseSchema,
  HomeFeedResponseSchema,
  ListeningStatsSchema,
  DiscoverWeeklyResponseSchema,
  type RecommendationResponse,
  type HomeFeedResponse,
  type ListeningStats,
  type DiscoverWeeklyResponse,
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

  const getPersonalized = async (
    params?: { limit?: number },
    config?: UseRequestConfig<RecommendationResponse>,
  ) => {
    return useRequest<RecommendationResponse>(
      RecommendationApiRoutes.PERSONALIZED,
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

  const getHomeFeed = async (
    config?: UseRequestConfig<HomeFeedResponse>,
  ) => {
    return useRequest<HomeFeedResponse>(
      RecommendationApiRoutes.HOME,
      {
        method: 'GET',
      },
      {
        schema: HomeFeedResponseSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getDiscoverWeekly = async (
    config?: UseRequestConfig<DiscoverWeeklyResponse>,
  ) => {
    return useRequest<DiscoverWeeklyResponse>(
      RecommendationApiRoutes.DISCOVER_WEEKLY,
      { method: 'GET' },
      {
        schema: DiscoverWeeklyResponseSchema,
        silent: true,
        ...config,
      },
    )
  }

  const getListeningStats = async (
    params?: { period?: string },
    config?: UseRequestConfig<ListeningStats>,
  ) => {
    return useRequest<ListeningStats>(
      RecommendationApiRoutes.STATS,
      { method: 'GET', params },
      {
        schema: ListeningStatsSchema,
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
    getHomeFeed,
    getPersonalized,
    getDiscoverWeekly,
    getListeningStats,
  }
}
