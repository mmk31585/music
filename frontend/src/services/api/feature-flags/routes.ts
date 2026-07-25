import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'
import { FeatureFlagsApiRoutes } from './enums'
import type { FeatureFlags } from './types'

export const useFeatureFlagsApi = () => {
  const list = async (config?: UseRequestConfig<FeatureFlags>) => {
    return useRequest<FeatureFlags>(
      FeatureFlagsApiRoutes.ALL,
      { method: 'GET' },
      { silent: true, ...config },
    )
  }

  return { list }
}
