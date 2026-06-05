import { useRequest } from '@/composables'
import type { UseRequestConfig } from '@/plugins/client/types'

import { AuthApiRoutes } from './enums'
import {
  AuthResponseSchema,
  CurrentUserResponseSchema,
  type AuthResponse,
  type CurrentUserResponse,
  type RegisterPayload,
  type LoginPayload,
  type LogoutPayload,
} from './types'

export const useAuthApi = () => {
  const login = async (payload: LoginPayload, config?: UseRequestConfig<AuthResponse>) => {
    return useRequest<AuthResponse>(
      AuthApiRoutes.LOGIN,
      {
        method: 'POST',
        data: payload,
      },
      {
        schema: AuthResponseSchema,
        ...config,
      },
    )
  }

  const register = async (payload: RegisterPayload, config?: UseRequestConfig<AuthResponse>) => {
    return useRequest<AuthResponse>(
      AuthApiRoutes.REGISTER,
      {
        method: 'POST',
        data: payload,
      },
      {
        schema: AuthResponseSchema,
        ...config,
      },
    )
  }

  const refresh = async (config?: UseRequestConfig<AuthResponse>) => {
    return useRequest<AuthResponse>(
      AuthApiRoutes.REFRESH,
      {
        method: 'POST',
      },
      {
        schema: AuthResponseSchema,
        ...config,
      },
    )
  }

  const logout = async (payload: LogoutPayload, config?: UseRequestConfig<unknown>) => {
    return useRequest(
      AuthApiRoutes.LOGOUT,
      {
        method: 'POST',
        data: payload,
      },
      {
        ...config,
      },
    )
  }

  const me = async (config?: UseRequestConfig<CurrentUserResponse>) => {
    return useRequest<CurrentUserResponse>(
      AuthApiRoutes.ME,
      {
        method: 'GET',
      },
      {
        schema: CurrentUserResponseSchema,
        silent: true,
        ...config,
      },
    )
  }

  return {
    login,
    register,
    refresh,
    logout,
    me,
  }
}
