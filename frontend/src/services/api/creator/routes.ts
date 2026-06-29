import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { CreatorApiRoutes } from './enums'
import {
  CreatorStatsSchema,
  CreatorDailyStatSchema,
  TrackStatsSchema,
  EarningsBreakdownSchema,
  PayoutSchema,
  PayoutMethodSchema,
  AudienceDataSchema,
  CreatorContentDataSchema,
  type CreatorStats,
  type CreatorDailyStat,
  type TrackStats,
  type EarningsBreakdown,
  type Payout,
  type PayoutMethod,
  type AudienceData,
  type CreatorContentData,
  type TrackUpdateRequest,
  type AlbumUpdateRequest,
} from './types'

export const useCreatorApi = () => {
  const getOverview = async (config?: UseRequestConfig<CreatorStats>) => {
    return useRequest<CreatorStats>(
      CreatorApiRoutes.OVERVIEW,
      { method: 'GET' },
      { schema: CreatorStatsSchema, silent: true, ...config },
    )
  }

  const getDailyStats = async (
    params?: { from?: string; to?: string; limit?: number },
    config?: UseRequestConfig<CreatorDailyStat[]>,
  ) => {
    return useRequest<CreatorDailyStat, true>(
      CreatorApiRoutes.DAILY_STATS,
      { method: 'GET', params },
      { schema: CreatorDailyStatSchema, silent: true, ...config },
    )
  }

  const getTrackStats = async (config?: UseRequestConfig<TrackStats[]>) => {
    return useRequest<TrackStats, true>(
      CreatorApiRoutes.TRACK_STATS,
      { method: 'GET' },
      { schema: TrackStatsSchema, silent: true, ...config },
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
  const getEarnings = async (config?: UseRequestConfig<EarningsBreakdown>) => {
    return useRequest<EarningsBreakdown>(
      CreatorApiRoutes.EARNINGS,
      { method: 'GET' },
      { schema: EarningsBreakdownSchema, silent: true, ...config },
    )
  }

  const getPayouts = async (
    params?: { limit?: number },
    config?: UseRequestConfig<Payout[]>,
  ) => {
    return useRequest<Payout, true>(
      CreatorApiRoutes.PAYOUTS,
      { method: 'GET', params },
      { schema: PayoutSchema, silent: true, ...config },
    )
  }

  const getPayoutMethods = async (config?: UseRequestConfig<PayoutMethod[]>) => {
    return useRequest<PayoutMethod, true>(
      CreatorApiRoutes.PAYOUT_METHODS,
      { method: 'GET' },
      { schema: PayoutMethodSchema, silent: true, ...config },
    )
  }

  // Audience
  const getAudience = async (
    params?: { top_limit?: number },
    config?: UseRequestConfig<AudienceData>,
  ) => {
    return useRequest<AudienceData>(
      CreatorApiRoutes.AUDIENCE,
      { method: 'GET', params },
      { schema: AudienceDataSchema, silent: true, ...config },
    )
  }

  // Content Management
  const getContent = async (config?: UseRequestConfig<CreatorContentData>) => {
    return useRequest<CreatorContentData>(
      CreatorApiRoutes.CONTENT,
      { method: 'GET' },
      { schema: CreatorContentDataSchema, silent: true, ...config },
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
