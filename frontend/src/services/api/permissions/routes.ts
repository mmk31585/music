import { useRequest } from '@/composables/useRequest'
import {
  PermissionSchema,
  RoleWithPermissionsSchema,
  UserAccessSchema,
  CheckPermissionResponseSchema,
  UserPermissionOverrideSchema,
  type Permission,
  type RoleWithPermissions,
  type UserAccess,
  type CheckPermissionResponse,
  type UserPermissionOverride,
} from './types'

export const usePermissionsApi = () => {
  const getMyAccess = async () => {
    return useRequest<UserAccess>('/permissions/me', { method: 'GET' }, { schema: UserAccessSchema })
  }

  const listPermissions = async () => {
    return useRequest<Permission>('/permissions', { method: 'GET' }, { schema: PermissionSchema, silent: true })
  }

  const listRoles = async () => {
    return useRequest<RoleWithPermissions>('/roles', { method: 'GET' }, { schema: RoleWithPermissionsSchema, silent: true })
  }

  const getRole = async (slug: string) => {
    return useRequest<RoleWithPermissions>(`/roles/${slug}`, { method: 'GET' }, { schema: RoleWithPermissionsSchema })
  }

  const createRole = async (data: { slug: string; label: string; description?: string; hierarchyLevel: number }) => {
    return useRequest<any>('/roles', { method: 'POST', data }, { silent: false })
  }

  const updateRole = async (id: number, data: Record<string, any>) => {
    return useRequest<any>(`/roles/${id}`, { method: 'PUT', data }, { silent: false })
  }

  const deleteRole = async (id: number) => {
    return useRequest<any>(`/roles/${id}`, { method: 'DELETE' }, { silent: false })
  }

  const setRolePermissions = async (slug: string, permissions: string[]) => {
    return useRequest<any>(`/roles/${slug}/permissions`, { method: 'PUT', data: { permissions } }, { silent: false })
  }

  const getUserAccess = async (userId: string) => {
    return useRequest<UserAccess>(`/users/${userId}/access`, { method: 'GET' }, { schema: UserAccessSchema })
  }

  const checkPermission = async (userId: string, permSlug: string) => {
    return useRequest<CheckPermissionResponse>(`/users/${userId}/access/check`, { method: 'POST', data: { permSlug } }, { schema: CheckPermissionResponseSchema })
  }

  const assignUserRole = async (userId: string, roleSlug: string) => {
    return useRequest<any>(`/users/${userId}/roles`, { method: 'POST', data: { roleSlug } }, { silent: false })
  }

  const removeUserRole = async (userId: string, roleSlug: string) => {
    return useRequest<any>(`/users/${userId}/roles/${roleSlug}`, { method: 'DELETE' }, { silent: false })
  }

  const getUserOverrides = async (userId: string) => {
    return useRequest<UserPermissionOverride>(`/users/${userId}/overrides`, { method: 'GET' }, { schema: UserPermissionOverrideSchema, silent: true })
  }

  const grantPermission = async (userId: string, permSlug: string, reason: string) => {
    return useRequest<any>(`/users/${userId}/overrides/grant`, { method: 'POST', data: { permSlug, reason } }, { silent: false })
  }

  const revokePermission = async (userId: string, permSlug: string) => {
    return useRequest<any>(`/users/${userId}/overrides/revoke`, { method: 'POST', data: { permSlug } }, { silent: false })
  }

  return {
    getMyAccess,
    listPermissions,
    listRoles,
    getRole,
    createRole,
    updateRole,
    deleteRole,
    setRolePermissions,
    getUserAccess,
    checkPermission,
    assignUserRole,
    removeUserRole,
    getUserOverrides,
    grantPermission,
    revokePermission,
  }
}
