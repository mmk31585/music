import { useRequest } from '@/composables/useRequest'
import { apiReplaceParams } from '@/utils/api-replace-params'
import { UserApiRoutes } from './enums'
import type { UserProps } from './types'

export const useUserApi = () => {
  const getPublicUserProfile = async (id: string) => {
    return useRequest<UserProps>(apiReplaceParams(UserApiRoutes.PUBLIC_PROFILE, { id }), {
      method: 'GET',
    })
  }

  return { getPublicUserProfile }
}
