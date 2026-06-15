import { useRequest } from '@/composables/useRequest'
import type { UseRequestConfig } from '@/plugins/client/types'

import { AuthApiRoutes } from './enums'
import {
  AuthResponseSchema,
  CurrentUserResponseSchema,
  AdminUserSchema,
  AdminUserListResponseSchema,
  type AuthResponse,
  type CurrentUserResponse,
  type RegisterPayload,
  type LoginPayload,
  type LogoutPayload,
  type AdminUser,
  type AdminUserListResponse,
  type AdminUserUpdatePayload,
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

  const refresh = async (refreshToken?: string, config?: UseRequestConfig<AuthResponse>) => {
    return useRequest<AuthResponse>(
      AuthApiRoutes.REFRESH,
      {
        method: 'POST',
        data: refreshToken ? { refreshToken } : undefined,
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

  // ==================== ADMIN USERS ====================

  const adminListUsers = async (
    params?: Record<string, unknown>,
    config?: UseRequestConfig<AdminUserListResponse>,
  ) => {
    return useRequest<AdminUserListResponse>(
      AuthApiRoutes.ADMIN_USERS,
      {
        method: 'GET',
        params,
      },
      {
        schema: AdminUserListResponseSchema,
        silent: true,
        ...config,
      },
    )
  }

  const adminGetUser = async (
    id: string,
    config?: UseRequestConfig<AdminUser>,
  ) => {
    return useRequest<AdminUser>(
      AuthApiRoutes.ADMIN_USER.replace(':id', String(id)),
      { method: 'GET' },
      {
        schema: AdminUserSchema,
        silent: true,
        ...config,
      },
    )
  }

  const adminUpdateUser = async (
    id: string,
    payload: AdminUserUpdatePayload,
    config?: UseRequestConfig<AdminUser>,
  ) => {
    return useRequest<AdminUser>(
      AuthApiRoutes.ADMIN_USER.replace(':id', String(id)),
      {
        method: 'PATCH',
        data: payload,
      },
      {
        schema: AdminUserSchema,
        ...config,
      },
    )
  }

  const adminDeleteUser = async (
    id: string,
    config?: UseRequestConfig<void>,
  ) => {
    return useRequest<void>(
      AuthApiRoutes.ADMIN_USER.replace(':id', String(id)),
      { method: 'DELETE' },
      {
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
    adminListUsers,
    adminGetUser,
    adminUpdateUser,
    adminDeleteUser,
  }
}
