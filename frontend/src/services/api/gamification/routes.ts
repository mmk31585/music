import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { GamificationApiRoutes } from './enums'
import type {
  GamificationProfile,
  BadgesResponse,
  ChallengesResponse,
  LeaderboardResponse,
  Badge,
} from './types'

export const useGamificationApi = () => {
  const getProfile = async (config?: UseRequestConfig<GamificationProfile>) => {
    return useRequest<GamificationProfile>(
      GamificationApiRoutes.PROFILE,
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const addXP = async (
    amount: number,
    source: string,
    config?: UseRequestConfig<{ new_balance: number }>,
  ) => {
    return useRequest<{ new_balance: number }>(
      GamificationApiRoutes.ADD_XP,
      { method: 'POST', data: { amount, source } },
      { silent: true, ...config },
    )
  }

  const getBadges = async (config?: UseRequestConfig<BadgesResponse>) => {
    return useRequest<BadgesResponse>(
      GamificationApiRoutes.BADGES,
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const getChallenges = async (config?: UseRequestConfig<ChallengesResponse>) => {
    return useRequest<ChallengesResponse>(
      GamificationApiRoutes.CHALLENGES,
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  const getLeaderboard = async (
    params?: { type?: string; limit?: number },
    config?: UseRequestConfig<LeaderboardResponse>,
  ) => {
    return useRequest<LeaderboardResponse>(
      GamificationApiRoutes.LEADERBOARD,
      { method: 'GET', params },
      { silent: true, ...config },
    )
  }

  const checkBadges = async (config?: UseRequestConfig<Badge[]>) => {
    return useRequest<Badge[]>(
      GamificationApiRoutes.CHECK_BADGES,
      { method: 'POST' },
      { silent: true, ...config },
    )
  }

  return { getProfile, addXP, getBadges, getChallenges, getLeaderboard, checkBadges }
}
