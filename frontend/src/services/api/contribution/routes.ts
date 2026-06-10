import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { ContributionApiRoutes } from './enums'
import type {
  Contribution,
  ContributionHistoryItem,
  ContentVersion,
  ContributorStats,
  CreateContributionPayload,
  ReviewContributionPayload,
  PaginatedResponse,
} from './types'

export const useContributionApi = () => {
  const create = async (
    payload: CreateContributionPayload,
    config?: UseRequestConfig<Contribution>,
  ) => {
    return useRequest<Contribution>(
      ContributionApiRoutes.CREATE,
      { method: 'POST', data: payload },
      { silent: false, ...config },
    )
  }

  const listMy = async (
    params?: { page?: number; page_size?: number },
    config?: UseRequestConfig<PaginatedResponse<Contribution>>,
  ) => {
    return useRequest<PaginatedResponse<Contribution>>(
      ContributionApiRoutes.LIST_MY,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const listPending = async (
    params?: { page?: number; page_size?: number },
    config?: UseRequestConfig<PaginatedResponse<Contribution>>,
  ) => {
    return useRequest<PaginatedResponse<Contribution>>(
      ContributionApiRoutes.LIST_PENDING,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const getById = async (id: string, config?: UseRequestConfig<Contribution>) => {
    return useRequest<Contribution>(
      ContributionApiRoutes.GET_BY_ID.replace(':id', id),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const review = async (
    id: string,
    payload: ReviewContributionPayload,
    config?: UseRequestConfig<Contribution>,
  ) => {
    return useRequest<Contribution>(
      ContributionApiRoutes.REVIEW.replace(':id', id),
      { method: 'POST', data: payload },
      { silent: false, ...config },
    )
  }

  const getHistory = async (id: string, config?: UseRequestConfig<ContributionHistoryItem[]>) => {
    return useRequest<ContributionHistoryItem[]>(
      ContributionApiRoutes.HISTORY.replace(':id', id),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const getLeaderboard = async (
    params?: { limit?: number },
    config?: UseRequestConfig<ContributorStats[]>,
  ) => {
    return useRequest<ContributorStats[]>(
      ContributionApiRoutes.LEADERBOARD,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const getByTarget = async (
    targetType: string,
    targetId: string,
    config?: UseRequestConfig<Contribution[]>,
  ) => {
    return useRequest<Contribution[]>(
      ContributionApiRoutes.BY_TARGET.replace(':targetType', targetType).replace(
        ':targetId',
        targetId,
      ),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const getVersions = async (
    targetType: string,
    targetId: string,
    config?: UseRequestConfig<ContentVersion[]>,
  ) => {
    return useRequest<ContentVersion[]>(
      ContributionApiRoutes.VERSIONS.replace(':targetType', targetType).replace(
        ':targetId',
        targetId,
      ),
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const apply = async (id: string, config?: UseRequestConfig<Contribution>) => {
    return useRequest<Contribution>(
      ContributionApiRoutes.APPLY.replace(':id', id),
      { method: 'POST' },
      { silent: false, ...config },
    )
  }

  const adminListContributions = async (
    params?: { page?: number; page_size?: number; status?: string; type?: string; q?: string },
    config?: UseRequestConfig<PaginatedResponse<Contribution>>,
  ) => {
    return useRequest<PaginatedResponse<Contribution>>(
      ContributionApiRoutes.ADMIN_LIST,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  return {
    create,
    listMy,
    listPending,
    getById,
    review,
    apply,
    adminListContributions,
    getHistory,
    getLeaderboard,
    getByTarget,
    getVersions,
  }
}
