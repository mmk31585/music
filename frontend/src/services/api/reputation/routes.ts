import { useRequest } from '@/composables/useRequest'
import {
  ReputationSummarySchema,
  TrustTierSchema,
  type ReputationSummary,
  type TrustTier,
} from './types'

export const useReputationApi = () => {
  const getUserReputation = async (userId: string) => {
    return useRequest<ReputationSummary>(`/reputation/users/${userId}`, { method: 'GET' }, { schema: ReputationSummarySchema })
  }

  const getTrustTiers = async () => {
    return useRequest<TrustTier>('/reputation/tiers', { method: 'GET' }, { schema: TrustTierSchema, silent: true })
  }

  const getTopContributors = async (limit = 10) => {
    return useRequest<ReputationSummary>(`/reputation/top?limit=${limit}`, { method: 'GET' }, { schema: ReputationSummarySchema, silent: true })
  }

  return {
    getUserReputation,
    getTrustTiers,
    getTopContributors,
  }
}
