import { useRequest } from '@/composables'
import type { UseRequestConfig } from '@/plugins/client/types'

import { AuthApiRoutes } from './enums'
import {
  AuthResponseSchema,
  UserSchema,
  type AuthResponse,
  type RegisterPayload,
  type User,
  type LoginPayload,
} from './types'

export const useAuthApi = () => {
  const login = async (
    payload: LoginPayload,
    config?: UseRequestConfig<AuthResponse>,
  ) => {
    return useRequest<AuthResponse>(
      AuthApiRoutes.LOGIN,
      {
        method: 'POST',
        data: payload
        ,
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

  const logout = async (config?: UseRequestConfig<unknown>) => {
    return useRequest(
      AuthApiRoutes.LOGOUT,
      {
        method: 'POST',
      },
      {
        ...config,
      },
    )
  }

  const me = async (config?: UseRequestConfig<User>) => {
    return useRequest<User>(
      AuthApiRoutes.ME,
      {
        method: 'GET',
      },
      {
        schema: UserSchema,
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
