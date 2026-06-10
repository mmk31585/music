import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { CreatorApiRoutes } from './enums'
import {
  OverviewResponseSchema,
  DailyStatsResponseSchema,
  TrackStatsResponseSchema,
  type OverviewResponse,
  type DailyStatsResponse,
  type TrackStatsResponse,
  type EarningsBreakdown,
  type Payout,
  type PayoutMethod,
  type AudienceData,
  type CreatorContentData,
  type TrackUpdateRequest,
  type AlbumUpdateRequest,
} from './types'

export const useCreatorApi = () => {
  const getOverview = async (config?: UseRequestConfig<OverviewResponse>) => {
    return useRequest<OverviewResponse>(
      CreatorApiRoutes.OVERVIEW,
      { method: 'GET' },
      { schema: OverviewResponseSchema, silent: true, ...config },
    )
  }

  const getDailyStats = async (
    params?: { from?: string; to?: string; limit?: number },
    config?: UseRequestConfig<DailyStatsResponse>,
  ) => {
    return useRequest<DailyStatsResponse>(
      CreatorApiRoutes.DAILY_STATS,
      { method: 'GET', params },
      { schema: DailyStatsResponseSchema, silent: true, ...config },
    )
  }

  const getTrackStats = async (config?: UseRequestConfig<TrackStatsResponse>) => {
    return useRequest<TrackStatsResponse>(
      CreatorApiRoutes.TRACK_STATS,
      { method: 'GET' },
      { schema: TrackStatsResponseSchema, silent: true, ...config },
    )
  }

  const refreshStats = async (config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      CreatorApiRoutes.REFRESH,
      { method: 'POST' },
      { silent: false, ...config },
    )
  }

  const isCreator = async (config?: UseRequestConfig<{ is_creator: boolean }>) => {
    return useRequest<{ is_creator: boolean }>(
      CreatorApiRoutes.IS_CREATOR,
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  // Earnings
  const getEarnings = async (config?: UseRequestConfig<{ data: EarningsBreakdown }>) => {
    return useRequest<{ data: EarningsBreakdown }>(
      CreatorApiRoutes.EARNINGS,
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const getPayouts = async (
    params?: { limit?: number },
    config?: UseRequestConfig<{ data: Payout[] }>,
  ) => {
    return useRequest<{ data: Payout[] }>(
      CreatorApiRoutes.PAYOUTS,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const getPayoutMethods = async (config?: UseRequestConfig<{ data: PayoutMethod[] }>) => {
    return useRequest<{ data: PayoutMethod[] }>(
      CreatorApiRoutes.PAYOUT_METHODS,
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  // Audience
  const getAudience = async (
    params?: { top_limit?: number },
    config?: UseRequestConfig<{ data: AudienceData }>,
  ) => {
    return useRequest<{ data: AudienceData }>(
      CreatorApiRoutes.AUDIENCE,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  // Content Management
  const getContent = async (config?: UseRequestConfig<{ data: CreatorContentData }>) => {
    return useRequest<{ data: CreatorContentData }>(
      CreatorApiRoutes.CONTENT,
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const updateTrack = async (
    trackId: string,
    data: TrackUpdateRequest,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      CreatorApiRoutes.UPDATE_TRACK.replace(':trackId', trackId),
      { method: 'PUT', data },
      config,
    )
  }

  const updateAlbum = async (
    albumId: string,
    data: AlbumUpdateRequest,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      CreatorApiRoutes.UPDATE_ALBUM.replace(':albumId', albumId),
      { method: 'PUT', data },
      config,
    )
  }

  const deleteTrack = async (trackId: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      CreatorApiRoutes.DELETE_TRACK.replace(':trackId', trackId),
      { method: 'DELETE' },
      config,
    )
  }

  const deleteAlbum = async (albumId: string, config?: UseRequestConfig<void>) => {
    return useRequest<void>(
      CreatorApiRoutes.DELETE_ALBUM.replace(':albumId', albumId),
      { method: 'DELETE' },
      config,
    )
  }

  return {
    getOverview,
    getDailyStats,
    getTrackStats,
    refreshStats,
    isCreator,
    getEarnings,
    getPayouts,
    getPayoutMethods,
    getAudience,
    getContent,
    updateTrack,
    updateAlbum,
    deleteTrack,
    deleteAlbum,
  }
}
